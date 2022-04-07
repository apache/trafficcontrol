// toreq implements a Traffic Ops client for features in the latest version.
//
// This should only be used if an endpoint or field needed for config gen is in the latest.
//
// Users should always check the returned bool, and if it's false, call the vendored toreq client and set the proper defaults for the new feature(s).
//
// All TOClient functions should check for 404 or 503 and return a bool false if so.
//
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
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/cache-config/t3cutil/toreq/toreqold"
	"github.com/apache/trafficcontrol/cache-config/t3cutil/toreq/torequtil"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

type TOClient struct {
	c          *toclient.Session
	old        *toreqold.TOClient
	NumRetries int
}

func (cl *TOClient) URL() string {
	if cl.c == nil {
		return cl.old.URL()
	}
	return cl.c.URL
}

func (cl *TOClient) SetURL(newURL string) {
	if cl.c == nil {
		cl.old.SetURL(newURL)
	} else {
		cl.c.URL = newURL
	}
}

const FsCookiePath = `/var/lib/trafficcontrol-cache-config/`

// New logs into Traffic Ops, returning the TOClient which contains the logged-in client.
func New(url *url.URL, user string, pass string, insecure bool, timeout time.Duration, userAgent string) (*TOClient, error) {
	log.Infoln("URL: '" + url.String() + "' User: '" + user + "' Pass len: '" + strconv.Itoa(len(pass)) + "'")

	client := &TOClient{}
	fsCookie, err := torequtil.GetFsCookie(FsCookiePath + user + ".cookie")
	if err != nil {
		log.Infoln("Error retrieving cookie: ", err)
	}
	toURLStr := url.Scheme + "://" + url.Host
	log.Infoln("TO URL string: '" + toURLStr + "'")
	log.Infoln("TO URL: '" + url.String() + "'")

	if fsCookie.Cookies != nil {
		toIP, err := net.LookupIP(url.Hostname())
		if err != nil {
			log.Warnln("error getting traffic ops IP: ", err)
		}
		log.Infof("Logging in to Traffic Ops '%s' with Cookie", toIP)
		cookies := []*http.Cookie{}
		for _, cookie := range fsCookie.Cookies {
			cookies = append(cookies, cookie.Cookie)
		}
		toClient := toclient.NewSession(user, pass, toURLStr, userAgent, &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
			},
		}, false)
		toClient.Client.Jar, err = cookiejar.New(nil)
		if err != nil {
			log.Warnln("error creating cookie jar ", err)
		}
		toClient.Client.Jar.SetCookies(url, cookies)
		client = &TOClient{c: toClient}
	} else {
		opts := toclient.Options{}
		opts.Insecure = insecure
		opts.UserAgent = userAgent
		opts.RequestTimeout = timeout
		toClient, inf, err := toclient.Login(toURLStr, user, pass, opts)
		latestSupported := inf.StatusCode != 404 && inf.StatusCode != 501

		if err != nil && latestSupported {
			return nil, fmt.Errorf("Logging in to Traffic Ops '%v' code %v: %v", torequtil.MaybeIPStr(inf.RemoteAddr), inf.StatusCode, err)
		} else if !latestSupported {
			log.Infof("toreqnew.New Logged into in to Traffic Ops: got %v, falling back to older client\n", inf.StatusCode)
		} else {
			log.Infoln("toreqnew.New Logged into in to Traffic Ops '" + torequtil.MaybeIPStr(inf.RemoteAddr) + "'")
		}
		if latestSupported {
			toAddr := net.Addr(nil)
			latestSupported, toAddr, err = IsLatestSupported(toClient)
			if err != nil {
				return nil, errors.New("checking Traffic Ops '" + torequtil.MaybeIPStr(toAddr) + "' support: " + err.Error())
			}
		}

		client = &TOClient{c: toClient}
		if !latestSupported {
			log.Warnln("toreqnew.New Traffic Ops '" + torequtil.MaybeIPStr(inf.RemoteAddr) + "' does not support the latest client, falling back ot the previous")

			oldClient, err := toreqold.New(url, user, pass, insecure, timeout, userAgent)
			if err != nil {
				return nil, errors.New("logging into old client: " + err.Error())
			}
			client.c = nil
			client.old = oldClient

			{
				newClient := toclient.NewNoAuthSession("", false, "", false, 0) // only used for the version, because toClient could be nil if it had an error
				log.Infof("cache-config Traffic Ops client: %v not supported, falling back to %v\n", newClient.APIVersion(), oldClient.APIVersion())
			}
		}
	}
	return client, nil
}

// FellBack() returns whether the client fell back to the previous major version, because Traffic Ops didn't support the latest.
func (cl *TOClient) FellBack() bool {
	return cl.c == nil
}

func IsLatestSupported(toClient *toclient.Session) (bool, net.Addr, error) {
	_, inf, err := toClient.Ping(toclient.RequestOptions{})
	if err != nil {
		if errS := strings.ToLower(err.Error()); strings.Contains(errS, "not found") || strings.Contains(errS, "not implemented") {
			return false, inf.RemoteAddr, nil
		}
		return false, inf.RemoteAddr, err
	}
	return true, inf.RemoteAddr, nil
}

// RequestInfoStr returns a loggable string with info about the Traffic Ops request.
//
// The returned string does not have a trailing newline, nor anything in the standard
// logger prefix (time, level, etc).
// If the string isn't going to be logged via lib/go-log, it's advisable to add a timestamp.
//
// This is safe to call even if the function returning a ReqInf returned an error;
// it checks for nil values in all cases, and the TO Client guarantees even if a non-nil
// error is returned, all ReqInf values are either nil or valid.
//
func RequestInfoStr(inf toclientlib.ReqInf, reqPath string) string {
	return fmt.Sprintf(`requestinfo path=%v ip=%v code=%v, date="%v" age=%v`,
		reqPath,
		torequtil.MaybeIPStr(inf.RemoteAddr),
		inf.StatusCode,
		torequtil.MaybeHdrStr(inf.RespHeaders, rfc.Date),
		torequtil.MaybeHdrStr(inf.RespHeaders, rfc.Age),
	)
}
