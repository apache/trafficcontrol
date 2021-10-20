/*

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

// Package client provides Go bindings to the Traffic Ops RPC API.
package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptrace"
	"strconv"
	"strings"
	"sync"
	"time"

	tc "github.com/apache/trafficcontrol/v6/lib/go-tc"

	"golang.org/x/net/publicsuffix"
)

// Session ...
type Session struct {
	UserName     string
	Password     string
	URL          string
	Client       *http.Client
	cache        map[string]CacheEntry
	cacheMutex   *sync.RWMutex
	useCache     bool
	UserAgentStr string
}

func NewSession(user, password, url, userAgent string, client *http.Client, useCache bool) *Session {
	return &Session{
		UserName:     user,
		Password:     password,
		URL:          url,
		Client:       client,
		cache:        map[string]CacheEntry{},
		cacheMutex:   &sync.RWMutex{},
		useCache:     useCache,
		UserAgentStr: userAgent,
	}
}

const DefaultTimeout = time.Second * time.Duration(30)

// HTTPError is returned on Update Session failure.
type HTTPError struct {
	HTTPStatusCode int
	HTTPStatus     string
	URL            string
	Body           string
}

// Error implements the error interface for our customer error type.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("%s[%d] - Error requesting Traffic Ops %s %s", e.HTTPStatus, e.HTTPStatusCode, e.URL, e.Body)
}

// CacheEntry ...
type CacheEntry struct {
	Entered    int64
	Bytes      []byte
	RemoteAddr net.Addr
}

// TODO JvD
const tmPollingInterval = 60

// loginCreds gathers login credentials for Traffic Ops.
func loginCreds(toUser string, toPasswd string) ([]byte, error) {
	credentials := tc.UserCredentials{
		Username: toUser,
		Password: toPasswd,
	}

	js, err := json.Marshal(credentials)
	if err != nil {
		err := fmt.Errorf("Error creating login json: %v", err)
		return nil, err
	}
	return js, nil
}

// loginToken gathers token login credentials for Traffic Ops.
func loginToken(token string) ([]byte, error) {
	form := tc.UserToken{
		Token: token,
	}

	j, e := json.Marshal(form)
	if e != nil {
		e := fmt.Errorf("Error creating token login json: %v", e)
		return nil, e
	}
	return j, nil
}

// login tries to log in to Traffic Ops, and set the auth cookie in the Session. Returns the IP address of the remote Traffic Ops.
func (to *Session) login() (net.Addr, error) {
	credentials, err := loginCreds(to.UserName, to.Password)
	if err != nil {
		return nil, errors.New("creating login credentials: " + err.Error())
	}

	path := apiBase + "/user/login"
	resp, remoteAddr, err := to.RawRequest("POST", path, credentials)
	resp, remoteAddr, err = to.ErrUnlessOK(resp, remoteAddr, err, path)
	if err != nil {
		return remoteAddr, errors.New("requesting: " + err.Error())
	}
	defer resp.Body.Close()

	var alerts tc.Alerts
	if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return remoteAddr, errors.New("decoding response JSON: " + err.Error())
	}

	success := false
	for _, alert := range alerts.Alerts {
		if alert.Level == "success" && alert.Text == "Successfully logged in." {
			success = true
			break
		}
	}

	if !success {
		return remoteAddr, fmt.Errorf("Login failed, alerts string: %+v", alerts)
	}

	return remoteAddr, nil
}

func (to *Session) loginWithToken(token []byte) (net.Addr, error) {
	path := apiBase + "/user/login/token"
	resp, remoteAddr, err := to.RawRequest(http.MethodPost, path, token)
	resp, remoteAddr, err = to.ErrUnlessOK(resp, remoteAddr, err, path)
	if err != nil {
		return remoteAddr, fmt.Errorf("requesting: %v", err)
	}
	defer resp.Body.Close()

	var alerts tc.Alerts
	if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return remoteAddr, fmt.Errorf("decoding response JSON: %v", err)
	}

	for _, alert := range alerts.Alerts {
		if alert.Level == tc.SuccessLevel.String() && alert.Text == "Successfully logged in." {
			return remoteAddr, nil
		}
	}

	return remoteAddr, fmt.Errorf("Login failed, alerts string: %+v", alerts)
}

// logout of Traffic Ops
func (to *Session) logout() (net.Addr, error) {
	credentials, err := loginCreds(to.UserName, to.Password)
	if err != nil {
		return nil, errors.New("creating login credentials: " + err.Error())
	}

	path := apiBase + "/user/logout"
	resp, remoteAddr, err := to.RawRequest("POST", path, credentials)
	resp, remoteAddr, err = to.ErrUnlessOK(resp, remoteAddr, err, path)
	if err != nil {
		return remoteAddr, errors.New("requesting: " + err.Error())
	}
	defer resp.Body.Close()

	var alerts tc.Alerts
	if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return remoteAddr, errors.New("decoding response JSON: " + err.Error())
	}

	success := false
	for _, alert := range alerts.Alerts {
		if alert.Level == "success" && alert.Text == "Successfully logged in." {
			success = true
			break
		}
	}

	if !success {
		return remoteAddr, fmt.Errorf("Logout failed, alerts string: %+v", alerts)
	}

	return remoteAddr, nil
}

// Login to traffic_ops, the response should set the cookie for this session
// automatically. Start with
//     to := traffic_ops.Login("user", "passwd", true)
// subsequent calls like to.GetData("datadeliveryservice") will be authenticated.
// Returns the logged in client, the remote address of Traffic Ops which was translated and used to log in, and any error. If the error is not nil, the remote address may or may not be nil, depending whether the error occurred before the login request.
func LoginWithAgent(toURL string, toUser string, toPasswd string, insecure bool, userAgent string, useCache bool, requestTimeout time.Duration) (*Session, net.Addr, error) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	jar, err := cookiejar.New(&options)
	if err != nil {
		return nil, nil, err
	}

	to := NewSession(toUser, toPasswd, toURL, userAgent, &http.Client{
		Timeout: requestTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
		Jar: jar,
	}, useCache)

	remoteAddr, err := to.login()
	if err != nil {
		return nil, remoteAddr, errors.New("logging in: " + err.Error())
	}
	return to, remoteAddr, nil
}

func LoginWithToken(toURL string, token string, insecure bool, userAgent string, useCache bool, requestTimeout time.Duration) (*Session, net.Addr, error) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	jar, err := cookiejar.New(&options)
	if err != nil {
		return nil, nil, err
	}

	client := http.Client{
		Timeout: requestTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
		Jar: jar,
	}

	to := NewSession("", "", toURL, userAgent, &client, useCache)
	tBts, err := loginToken(token)
	if err != nil {
		return nil, nil, fmt.Errorf("logging in: %v", err)
	}

	remoteAddr, err := to.loginWithToken(tBts)
	if err != nil {
		return nil, remoteAddr, fmt.Errorf("logging in: %v", err)
	}
	return to, remoteAddr, nil
}

// Logout of traffic_ops
func LogoutWithAgent(toURL string, toUser string, toPasswd string, insecure bool, userAgent string, useCache bool, requestTimeout time.Duration) (*Session, net.Addr, error) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	jar, err := cookiejar.New(&options)
	if err != nil {
		return nil, nil, err
	}

	to := NewSession(toUser, toPasswd, toURL, userAgent, &http.Client{
		Timeout: requestTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
		Jar: jar,
	}, useCache)

	remoteAddr, err := to.logout()
	if err != nil {
		return nil, remoteAddr, errors.New("logging out: " + err.Error())
	}
	return to, remoteAddr, nil
}

// NewNoAuthSession returns a new Session without logging in
// this can be used for querying unauthenticated endpoints without requiring a login
func NewNoAuthSession(toURL string, insecure bool, userAgent string, useCache bool, requestTimeout time.Duration) *Session {
	return NewSession("", "", toURL, userAgent, &http.Client{
		Timeout: requestTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
	}, useCache)
}

// ErrUnlessOk returns nil and an error if the given Response's status code is anything but 200 OK. This includes reading the Response.Body and Closing it. Otherwise, the given response and error are returned unchanged.
func (to *Session) ErrUnlessOK(resp *http.Response, remoteAddr net.Addr, err error, path string) (*http.Response, net.Addr, error) {
	if err != nil {
		return resp, remoteAddr, err
	}
	if resp.StatusCode < 300 {
		return resp, remoteAddr, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotImplemented {
		return nil, remoteAddr, errors.New("Traffic Ops Server returned 'Not Implemented', this client is probably newer than Traffic Ops, and you probably need to either upgrade Traffic Ops, or use a client whose version matches your Traffic Ops version.")
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, remoteAddr, readErr
	}
	return nil, remoteAddr, errors.New(resp.Status + "[" + strconv.Itoa(resp.StatusCode) + "] - Error requesting Traffic Ops " + to.getURL(path) + " " + string(body))
}

func (to *Session) getURL(path string) string { return to.URL + path }

// request performs the HTTP request to Traffic Ops, trying to refresh the cookie if an Unauthorized or Forbidden code is received. It only tries once. If the login fails, the original Unauthorized/Forbidden response is returned. If the login succeeds and the subsequent re-request fails, the re-request's response is returned even if it's another Unauthorized/Forbidden.
// Returns the response, the remote address of the Traffic Ops instance used, and any error.
// The returned net.Addr is guaranteed to be either nil or valid, even if the returned error is not nil. Callers are encouraged to check and use the net.Addr if an error is returned, and use the remote address in their own error messages. This violates the Go idiom that a non-nil error implies all other values are undefined, but it's more straightforward than alternatives like typecasting.
func (to *Session) request(method, path string, body []byte) (*http.Response, net.Addr, error) {
	r, remoteAddr, err := to.RawRequest(method, path, body)
	if err != nil {
		return r, remoteAddr, err
	}
	if r.StatusCode != http.StatusUnauthorized && r.StatusCode != http.StatusForbidden {
		return to.ErrUnlessOK(r, remoteAddr, err, path)
	}
	if _, lerr := to.login(); lerr != nil {
		return to.ErrUnlessOK(r, remoteAddr, err, path) // if re-logging-in fails, return the original request's response
	}

	// return second request, even if it's another Unauthorized or Forbidden.
	r, remoteAddr, err = to.RawRequest(method, path, body)
	return to.ErrUnlessOK(r, remoteAddr, err, path)
}

// RawRequest performs the actual HTTP request to Traffic Ops, simply, without trying to refresh the cookie if an Unauthorized code is returned.
// Returns the response, the remote address of the Traffic Ops instance used, and any error.
// The returned net.Addr is guaranteed to be either nil or valid, even if the returned error is not nil. Callers are encouraged to check and use the net.Addr if an error is returned, and use the remote address in their own error messages. This violates the Go idiom that a non-nil error implies all other values are undefined, but it's more straightforward than alternatives like typecasting.
func (to *Session) RawRequest(method, path string, body []byte) (*http.Response, net.Addr, error) {
	url := to.getURL(path)

	var req *http.Request
	var err error
	remoteAddr := net.Addr(nil)

	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
		if err != nil {
			return nil, remoteAddr, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, remoteAddr, err
		}
	}

	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			remoteAddr = connInfo.Conn.RemoteAddr()
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	req.Header.Set("User-Agent", to.UserAgentStr)

	resp, err := to.Client.Do(req)
	if err != nil {
		return nil, remoteAddr, err
	}

	return resp, remoteAddr, nil
}

type ReqInf struct {
	CacheHitStatus CacheHitStatus
	RemoteAddr     net.Addr
}

type CacheHitStatus string

const CacheHitStatusHit = CacheHitStatus("hit")
const CacheHitStatusExpired = CacheHitStatus("expired")
const CacheHitStatusMiss = CacheHitStatus("miss")
const CacheHitStatusInvalid = CacheHitStatus("")

func (s CacheHitStatus) String() string {
	return string(s)
}

func StringToCacheHitStatus(s string) CacheHitStatus {
	s = strings.ToLower(s)
	switch s {
	case "hit":
		return CacheHitStatusHit
	case "expired":
		return CacheHitStatusExpired
	case "miss":
		return CacheHitStatusMiss
	default:
		return CacheHitStatusInvalid
	}
}

// setCache Sets the given cache key and value. This is threadsafe for multiple goroutines.
func (to *Session) setCache(path string, entry CacheEntry) {
	if !to.useCache {
		return
	}
	to.cacheMutex.Lock()
	defer to.cacheMutex.Unlock()
	to.cache[path] = entry
}

// getCache gets the cache value at the given key, or false if it doesn't exist. This is threadsafe for multiple goroutines.
func (to *Session) getCache(path string) (CacheEntry, bool) {
	to.cacheMutex.RLock()
	defer to.cacheMutex.RUnlock()
	cacheEntry, ok := to.cache[path]
	return cacheEntry, ok
}

//if cacheEntry, ok := to.Cache[path]; ok {

// getBytesWithTTL gets the path, and caches in the session. Returns bytes from the cache, if found and the TTL isn't expired. Otherwise, gets it and store it in cache
func (to *Session) getBytesWithTTL(path string, ttl int64) ([]byte, ReqInf, error) {
	var body []byte
	var err error
	var cacheHitStatus CacheHitStatus
	var remoteAddr net.Addr

	getFresh := false
	if cacheEntry, ok := to.getCache(path); ok {
		if cacheEntry.Entered > time.Now().Unix()-ttl {
			cacheHitStatus = CacheHitStatusHit
			body = cacheEntry.Bytes
			remoteAddr = cacheEntry.RemoteAddr
		} else {
			cacheHitStatus = CacheHitStatusExpired
			getFresh = true
		}
	} else {
		cacheHitStatus = CacheHitStatusMiss
		getFresh = true
	}

	if getFresh {
		body, remoteAddr, err = to.getBytes(path)
		if err != nil {
			return nil, ReqInf{CacheHitStatus: CacheHitStatusInvalid, RemoteAddr: remoteAddr}, err
		}

		newEntry := CacheEntry{
			Entered:    time.Now().Unix(),
			Bytes:      body,
			RemoteAddr: remoteAddr,
		}
		to.setCache(path, newEntry)
	}

	return body, ReqInf{CacheHitStatus: cacheHitStatus, RemoteAddr: remoteAddr}, nil
}

// GetBytes - get []bytes array for a certain path on the to session.
// returns the raw body, the remote address the Traffic Ops URL resolved to, or any error. If the error is not nil, the RemoteAddr may or may not be nil, depending whether the error occurred before the request was executed.
func (to *Session) getBytes(path string) ([]byte, net.Addr, error) {
	resp, remoteAddr, err := to.request("GET", path, nil)
	if err != nil {
		return nil, remoteAddr, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, remoteAddr, err
	}

	return body, remoteAddr, nil
}
