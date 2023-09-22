// toreq implements a Traffic Ops client for features in the latest version.
//
// This should only be used if an endpoint or field needed for config gen is in the latest.
//
// Users should always check the returned bool, and if it's false, call the vendored toreq client and set the proper defaults for the new feature(s).
//
// All TOClient functions should check for 404 or 503 and return a bool false if so.
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

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil/toreq/toreqold"
	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil/toreq/torequtil"
	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

type TOClient struct {
	// c is the "new" Traffic Ops client, for the latest major API.
	//
	// This MUST NOT be accessed without checking for nil.
	// If the Traffic Ops server doesn't support the latest API, it will fall back,
	// and c will be nil and old will not, and old must be used.
	//
	// **WARNING** This MUST NOT be accessed without checking for nil. See above.
	c *toclient.Session

	// old is the "old" Traffic Ops client, for the previous major API.
	// This will be nil unless the Traffic Ops Server doesn't support the latest API,
	// in which case this will exist and c will be nil.
	old *toreqold.TOClient

	// NumRetries is the number of times to retry Traffic Ops server failures
	// before giving up and returning an error to the caller.
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

func (cl *TOClient) HTTPClient() *http.Client {
	if cl.c == nil {
		return cl.old.HTTPClient()
	}
	return cl.c.Client
}

// New logs into Traffic Ops, returning the TOClient which contains the logged-in client.
func New(url *url.URL, user string, pass string, insecure bool, timeout time.Duration, userAgent string) (*TOClient, error) {
	log.Infoln("URL: '" + url.String() + "' User: '" + user + "' Pass len: '" + strconv.Itoa(len(pass)) + "'")

	cookiePath := torequtil.CookieCachePath(user)

	fsCookie, err := torequtil.GetFsCookie(cookiePath)
	if err != nil {
		log.Infof("Failed to retrieve cached cookie for user '%v' at '%v', using password login: %v", user, cookiePath, err)
		return newWithPassword(url, user, pass, insecure, timeout, userAgent)
	}

	if fsCookie.Cookies == nil {
		log.Infof("Cached cookie for user '%v' at '%v' not found, using password login", user, cookiePath)
		return newWithPassword(url, user, pass, insecure, timeout, userAgent)
	}

	log.Infof("Cached cookie for user '%v' at '%v' found, attempting to reuse cookie to avoid login", user, cookiePath)
	return newWithCookie(url, user, pass, insecure, timeout, userAgent, fsCookie)
}

func newWithPassword(url *url.URL, user string, pass string, insecure bool, timeout time.Duration, userAgent string) (*TOClient, error) {
	opts := toclient.Options{}
	opts.Insecure = insecure
	opts.UserAgent = userAgent
	opts.RequestTimeout = timeout

	toURLStr := makeTOURLStr(url)
	log.Infoln("Traffic Ops URL string: '" + toURLStr + "'")

	toClient, inf, err := toclient.Login(toURLStr, user, pass, opts)
	if err != nil {
		if errIsUnsupportedVersion := inf.StatusCode == 404 || inf.StatusCode == 501; errIsUnsupportedVersion {
			log.Infof("toreqnew.New logging into Traffic Ops '%v': got %v, falling back to older client\n", torequtil.MaybeIPStr(inf.RemoteAddr), inf.StatusCode)
			return checkLatestAndFallBack(nil, url, user, pass, insecure, timeout, userAgent)
		}
		return nil, fmt.Errorf("Logging in to Traffic Ops '%v' code %v: %v", torequtil.MaybeIPStr(inf.RemoteAddr), inf.StatusCode, err)
	}

	// we successfully logged in, but the login may not have used the latest API,
	// double-check the client's API is supported.
	return checkLatestAndFallBack(toClient, url, user, pass, insecure, timeout, userAgent)
}

func newWithCookie(url *url.URL, user string, pass string, insecure bool, timeout time.Duration, userAgent string, fsCookie torequtil.FsCookie) (*TOClient, error) {
	toURLStr := makeTOURLStr(url)
	log.Infoln("Traffic Ops URL string: '" + toURLStr + "'")

	toClient := toclient.NewSession(user, pass, toURLStr, userAgent, &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
	}, false)
	err := error(nil)
	toClient.Client.Jar, err = cookiejar.New(nil)
	if err != nil {
		return nil, errors.New("error creating cookie jar: " + err.Error())
	}
	toClient.Client.Jar.SetCookies(url, fsCookie.GetHTTPCookies())
	return checkLatestAndFallBack(toClient, url, user, pass, insecure, timeout, userAgent)
}

// checkLatestAndFallBack takes a client and checks if it supports the latest Traffic ops API.
// If not, it attempts to fallback.
//
// The passed client.c may be nil if there was an error logging in, in which case it will be assumed that
// the latest API isn't supported and fallback will be tried.
//
// Returns a TOClient which is the latest if supported or has fallen back to the previous API if not, and any error.
func checkLatestAndFallBack(client *toclient.Session, url *url.URL, user string, pass string, insecure bool, timeout time.Duration, userAgent string) (*TOClient, error) {
	latestSupported, toAddr, err := IsLatestSupported(client)
	if err != nil {
		return nil, errors.New("checking Traffic Ops '" + torequtil.MaybeIPStr(toAddr) + "' support: " + err.Error())
	}

	if latestSupported {
		log.Infof("Traffic Ops '%v' supports this client's latest API version %v, using latest client\n", torequtil.MaybeIPStr(toAddr), client.APIVersion())
		return &TOClient{c: client}, nil
	}

	log.Warnf("Traffic Ops '%v' does not support the latest client API version %v, falling back to the previous\n", torequtil.MaybeIPStr(toAddr), LatestKnownAPIVersion())

	oldClient, err := toreqold.New(url, user, pass, insecure, timeout, userAgent)
	if err != nil {
		return nil, errors.New("logging into old client: " + err.Error())
	}

	log.Warnf("Latest Traffic Ops client version %v not supported, falling back to %v\n", LatestKnownAPIVersion(), oldClient.APIVersion())

	return &TOClient{old: oldClient}, nil
}

func LatestKnownAPIVersion() string {
	newClient := toclient.NewNoAuthSession("", false, "", false, 0) // created for the version, because a real toClient could be nil if it had an error
	return newClient.APIVersion()
}

// makeTOURLStr creates the Traffic Ops client URL string from uri.
// It specifically returns the scheme and host, but drops any path in uri.
// The uri must not be nil, and if the scheme or host is malformed the returned value will be malformed.
func makeTOURLStr(uri *url.URL) string {
	return uri.Scheme + "://" + uri.Host
}

// FellBack() returns whether the client fell back to the previous major version, because Traffic Ops didn't support the latest.
func (cl *TOClient) FellBack() bool {
	return cl.c == nil
}

// IsLatestSupported returns whether toClient supports the latest API, the address of the Traffic Ops connected to (which may be nil), and any error.
//
// A nil toClient may be passed, in which case it will be assumed that there was an error creating it and the latest isn't supported,
// and false will be returned with no address and no error.
func IsLatestSupported(toClient *toclient.Session) (bool, net.Addr, error) {
	if toClient == nil {
		return false, nil, nil
	}
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
func RequestInfoStr(inf toclientlib.ReqInf, reqPath string) string {
	return fmt.Sprintf(`requestinfo path=%v ip=%v code=%v, date="%v" age=%v`,
		reqPath,
		torequtil.MaybeIPStr(inf.RemoteAddr),
		inf.StatusCode,
		torequtil.MaybeHdrStr(inf.RespHeaders, rfc.Date),
		torequtil.MaybeHdrStr(inf.RespHeaders, rfc.Age),
	)
}
