package toreq

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptrace"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	toclient "github.com/apache/trafficcontrol/traffic_ops/client"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/config"
)

// TrafficOpsRequest makes a request to Traffic Ops for the given method, url, and body.
// If it gets an Unauthorized or Forbidden, it tries to log in again and makes the request again.
func TrafficOpsRequest(cfg config.TCCfg, method string, url string, body []byte) (string, int, error) {
	return trafficOpsRequestWithRetry(cfg, method, url, body, cfg.NumRetries)
}

func trafficOpsRequestWithRetry(
	cfg config.TCCfg,
	method string,
	url string,
	body []byte,
	retryNum int,
) (string, int, error) {
	currentRetry := 0
	for {
		body, code, err := trafficOpsRequestWithLogin(cfg, method, url, body)
		if err == nil || currentRetry == retryNum {
			return body, code, err
		}

		sleepSeconds := config.RetryBackoffSeconds(currentRetry)
		log.Errorf("getting '%v' '%v', sleeping for %v seconds: %v\n", method, url, sleepSeconds, err)
		currentRetry++
		time.Sleep(time.Second * time.Duration(sleepSeconds))
	}
}

func trafficOpsRequestWithLogin(
	cfg config.TCCfg,
	method string,
	url string,
	body []byte,
) (string, int, error) {
	resp, toIP, err := rawTrafficOpsRequest(*cfg.TOClient, method, url, body)
	if err != nil {
		toIPStr := ""
		if toIP != nil {
			toIPStr = toIP.String()
		}
		return "", 1, errors.New("requesting from Traffic Ops '" + toIPStr + "': " + err.Error())
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		resp.Body.Close()

		log.Infoln("TrafficOpsRequest got unauthorized/forbidden, logging in again")
		log.Infof("TrafficOpsRequest url '%v' user '%v' pass '%v'\n", (*cfg.TOClient).URL, (*cfg.TOClient).UserName, (*cfg.TOClient).Password)

		useCache := false
		newTOClient, toIP, err := toclient.LoginWithAgent((*cfg.TOClient).URL, (*cfg.TOClient).UserName, (*cfg.TOClient).Password, cfg.TOInsecure, config.UserAgent, useCache, cfg.TOTimeout)
		if err != nil {
			toIPStr := ""
			if toIP != nil {
				toIPStr = toIP.String()
			}
			return "", 1, errors.New("logging in to Traffic Ops IP '" + toIPStr + "': " + err.Error())
		}
		*cfg.TOClient = newTOClient

		resp, toIP, err = rawTrafficOpsRequest(*cfg.TOClient, method, url, body)
		if err != nil {
			toIPStr := ""
			if toIP != nil {
				toIPStr = toIP.String()
			}
			return "", 1, errors.New("requesting from Traffic Ops '" + toIPStr + "': " + err.Error())
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			bts = []byte("(read failure)") // if it's a non-200 and the body read fails, don't error, just note the read fail in the error
		}
		return "", resp.StatusCode, errors.New("Traffic Ops returned non-200 code '" + strconv.Itoa(resp.StatusCode) + "' body '" + string(bts) + "'")
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		toIPStr := ""
		if toIP != nil {
			toIPStr = toIP.String()
		}
		return "", resp.StatusCode, errors.New("reading body from Traffic Ops '" + toIPStr + "': " + err.Error())
	}

	if err := IntegrityCheck(resp.Header, bts, url); err != nil {
		return "", resp.StatusCode, errors.New("integrity check failed for url '" + url + "': " + err.Error())
	}

	return string(bts), http.StatusOK, nil
}

func IntegrityCheck(hdr http.Header, body []byte, url string) error {
	if hdrSHA := hdr.Get("Whole-Content-SHA512"); hdrSHA != "" {
		realSHA := sha512.Sum512(body)
		realSHAStr := base64.StdEncoding.EncodeToString(realSHA[:])
		if realSHAStr != hdrSHA {
			return errors.New("Body does not match header Whole-Content-SHA512")
		}
		log.Infoln("Integrity check for url '" + url + "' passed (sha match)")
		return nil
	}
	if hdrLenStr := hdr.Get("Content-Length"); hdrLenStr != "" {
		hdrLen, err := strconv.Atoi(hdrLenStr)
		if err != nil {
			return errors.New("No Whole-Content-SHA512, and Content-Length '" + hdrLenStr + "' is not an integer")
		}
		if hdrLen != len(body) {
			return errors.New("No Whole-Content-SHA512, and Content-Length '" + hdrLenStr + "' does not match body length")
		}
		log.Infoln("Integrity check for url '" + url + "' passed (length match)\n")
		return nil
	}
	return errors.New("No Whole-Content-SHA512, and no Content-Length, cannot verify content")
}

// rawTrafficOpsRequest makes a request to Traffic Ops for the given method, url, and body.
// If it gets an Unauthorized or Forbidden, it tries to log in again and makes the request again.
func rawTrafficOpsRequest(toClient *toclient.Session, method string, url string, body []byte) (*http.Response, net.Addr, error) {
	bodyReader := io.Reader(nil)
	if len(body) > 0 {
		bodyReader = bytes.NewBuffer(body)
	}

	remoteAddr := net.Addr(nil)
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, remoteAddr, err
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			remoteAddr = connInfo.Conn.RemoteAddr()
		},
	}))

	req.Header.Set("User-Agent", toClient.UserAgentStr)

	resp, err := toClient.Client.Do(req)
	if err != nil {
		return nil, remoteAddr, err
	}

	return resp, remoteAddr, nil
}

// MaybeIPStr returns the Traffic Ops IP string if it isn't nil, or the empty string if it is.
func MaybeIPStr(addr net.Addr) string {
	if addr != nil {
		return addr.String()
	}
	return ""
}
