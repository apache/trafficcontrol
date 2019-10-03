package main

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
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptrace"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/net/publicsuffix"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	toclient "github.com/apache/trafficcontrol/traffic_ops/client"
)

// GetClient returns a TO Client, using a cached cookie if it exists, or logging in otherwise
func GetClient(toURL string, toUser string, toPass string, tempDir string, cacheFileMaxAge time.Duration, toTimeout time.Duration, toInsecure bool) (*toclient.Session, error) {
	cookies, err := GetCookiesFromFile(tempDir, cacheFileMaxAge)
	if err != nil {
		log.Infoln("failed to get cookies from cache file (trying real TO): " + err.Error())
		cookies = ""
	}

	if cookies == "" {
		err := error(nil)
		cookies, err = GetCookiesFromTO(toURL, toUser, toPass, tempDir, toTimeout, toInsecure)
		if err != nil {
			return nil, errors.New("getting cookies from Traffic Ops: " + err.Error())
		}
		log.Infoln("using cookies from TO")
	} else {
		log.Infoln("using cookies from cache file")
	}

	useCache := false
	toClient := toclient.NewNoAuthSession(toURL, toInsecure, UserAgent, useCache, toTimeout)
	toClient.UserName = toUser
	toClient.Password = toPass

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, errors.New("making cookie jar: " + err.Error())
	}
	toClient.Client.Jar = jar

	toURLParsed, err := url.Parse(toURL)
	if err != nil {
		return nil, errors.New("parsing Traffic Ops URL '" + toURL + "': " + err.Error())
	}

	toClient.Client.Jar.SetCookies(toURLParsed, StringToCookies(cookies))
	return toClient, nil
}

// GetCookies gets the cookies from logging in to Traffic Ops.
// If this succeeds, it also writes the cookies to TempSubdir/TempCookieFileName.
func GetCookiesFromTO(toURL string, toUser string, toPass string, tempDir string, toTimeout time.Duration, toInsecure bool) (string, error) {
	toURLParsed, err := url.Parse(toURL)
	if err != nil {
		return "", errors.New("parsing Traffic Ops URL '" + toURL + "': " + err.Error())
	}

	toUseCache := false
	toClient, toIP, err := toclient.LoginWithAgent(toURL, toUser, toPass, toInsecure, UserAgent, toUseCache, toTimeout)
	if err != nil {
		toIPStr := ""
		if toIP != nil {
			toIPStr = toIP.String()
		}
		return "", errors.New("logging in to Traffic Ops IP '" + toIPStr + "': " + err.Error())
	}

	cookiesStr := CookiesToString(toClient.Client.Jar.Cookies(toURLParsed))
	WriteCookiesToFile(cookiesStr, tempDir)

	return cookiesStr, nil
}

// TrafficOpsRequest makes a request to Traffic Ops for the given method, url, and body.
// If it gets an Unauthorized or Forbidden, it tries to log in again and makes the request again.
func TrafficOpsRequest(cfg TCCfg, method string, url string, body []byte) (string, int, error) {
	return trafficOpsRequestWithRetry(cfg, method, url, body, cfg.NumRetries)
}

func trafficOpsRequestWithRetry(
	cfg TCCfg,
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

		sleepSeconds := RetryBackoffSeconds(currentRetry)
		log.Errorf("getting '%v' '%v', sleeping for %v seconds: %v\n", method, url, sleepSeconds, err)
		currentRetry++
		time.Sleep(time.Second * time.Duration(sleepSeconds))
	}
}

func trafficOpsRequestWithLogin(
	cfg TCCfg,
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
		newTOClient, toIP, err := toclient.LoginWithAgent((*cfg.TOClient).URL, (*cfg.TOClient).UserName, (*cfg.TOClient).Password, cfg.TOInsecure, UserAgent, useCache, cfg.TOTimeout)
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
func MaybeIPStr(reqInf toclient.ReqInf) string {
	if reqInf.RemoteAddr != nil {
		return reqInf.RemoteAddr.String()
	}
	return ""
}

// TCParamsToParamsWithProfiles unmarshals the Profiles that the tc struct doesn't.
func TCParamsToParamsWithProfiles(tcParams []tc.Parameter) ([]ParameterWithProfiles, error) {
	params := make([]ParameterWithProfiles, 0, len(tcParams))
	for _, tcParam := range tcParams {
		param := ParameterWithProfiles{Parameter: tcParam}

		profiles := []string{}
		if err := json.Unmarshal(tcParam.Profiles, &profiles); err != nil {
			return nil, errors.New("unmarshalling JSON from parameter '" + strconv.Itoa(param.ID) + "': " + err.Error())
		}
		param.ProfileNames = profiles
		param.Profiles = nil
		params = append(params, param)
	}
	return params, nil
}

type ParameterWithProfiles struct {
	tc.Parameter
	ProfileNames []string
}

type ParameterWithProfilesMap struct {
	tc.Parameter
	ProfileNames map[string]struct{}
}

func ParameterWithProfilesToMap(tcParams []ParameterWithProfiles) []ParameterWithProfilesMap {
	params := []ParameterWithProfilesMap{}
	for _, tcParam := range tcParams {
		param := ParameterWithProfilesMap{Parameter: tcParam.Parameter, ProfileNames: map[string]struct{}{}}
		for _, profile := range tcParam.ProfileNames {
			param.ProfileNames[profile] = struct{}{}
		}
		params = append(params, param)
	}
	return params
}
