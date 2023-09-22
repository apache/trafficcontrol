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
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// Login authenticates with Traffic Ops and returns the client object.
//
// Returns the logged in client, the remote address of Traffic Ops which was translated and used to log in, and any error. If the error is not nil, the remote address may or may not be nil, depending whether the error occurred before the login request.
//
// See ClientOpts for details about options, which options are required, and how they behave.
func Login(url, user, pass string, opts ClientOpts) (*Session, toclientlib.ReqInf, error) {
	cl, reqInf, err := toclientlib.Login(url, user, pass, opts.ClientOpts, apiVersions())
	if err != nil {
		return nil, reqInf, err
	}
	return &Session{TOClient: *cl}, reqInf, err
}

type ClientOpts struct {
	toclientlib.ClientOpts
}

// Session is a Traffic Ops client.
type Session struct {
	toclientlib.TOClient
}

func NewSession(user, password, url, userAgent string, client *http.Client, useCache bool) *Session {
	return &Session{
		TOClient: *toclientlib.NewClient(user, password, url, userAgent, client, apiVersions()),
	}
}

// Login to traffic_ops, the response should set the cookie for this session
// automatically. Start with
//
//	to := traffic_ops.Login("user", "passwd", true)
//
// subsequent calls like to.GetData("datadeliveryservice") will be authenticated.
// Returns the logged in client, the remote address of Traffic Ops which was translated and used to log in, and any error. If the error is not nil, the remote address may or may not be nil, depending whether the error occurred before the login request.
// The useCache argument is ignored. It exists to avoid breaking compatibility, and does not exist in newer functions.
func LoginWithAgent(toURL string, toUser string, toPasswd string, insecure bool, userAgent string, useCache bool, requestTimeout time.Duration) (*Session, net.Addr, error) {
	cl, ip, err := toclientlib.LoginWithAgent(toURL, toUser, toPasswd, insecure, userAgent, requestTimeout, apiVersions())
	if err != nil {
		return nil, nil, err
	}
	return &Session{TOClient: *cl}, ip, err
}

func LoginWithToken(toURL string, token string, insecure bool, userAgent string, useCache bool, requestTimeout time.Duration) (*Session, net.Addr, error) {
	cl, ip, err := toclientlib.LoginWithToken(toURL, token, insecure, userAgent, requestTimeout, apiVersions())
	if err != nil {
		return nil, nil, err
	}
	return &Session{TOClient: *cl}, ip, err
}

// Logout of Traffic Ops.
// The useCache argument is ignored. It exists to avoid breaking compatibility, and does not exist in newer functions.
func LogoutWithAgent(toURL string, toUser string, toPasswd string, insecure bool, userAgent string, useCache bool, requestTimeout time.Duration) (*Session, net.Addr, error) {
	cl, ip, err := toclientlib.LogoutWithAgent(toURL, toUser, toPasswd, insecure, userAgent, requestTimeout, apiVersions())
	if err != nil {
		return nil, nil, err
	}
	return &Session{TOClient: *cl}, ip, err
}

// NewNoAuthSession returns a new Session without logging in
// this can be used for querying unauthenticated endpoints without requiring a login
// The useCache argument is ignored. It exists to avoid breaking compatibility, and does not exist in newer functions.
func NewNoAuthSession(toURL string, insecure bool, userAgent string, useCache bool, requestTimeout time.Duration) *Session {
	return NewSession("", "", toURL, userAgent, &http.Client{
		Timeout: requestTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
	}, useCache)
}

// ErrIsNotImplemented checks if err ultimately arose at least in part because
// the Traffic Ops server did not support the requested API version.
func ErrIsNotImplemented(err error) bool {
	return err != nil && strings.Contains(err.Error(), ErrNotImplemented.Error()) // use string.Contains in case context was added to the error
}

// ErrNotImplemented is returned when Traffic Ops returns a 501 Not Implemented
// Users should check ErrIsNotImplemented rather than comparing directly, in case context was added.
var ErrNotImplemented = errors.New("Traffic Ops Server returned 'Not Implemented', this client is probably newer than Traffic Ops, and you probably need to either upgrade Traffic Ops, or use a client whose version matches your Traffic Ops version")

// errUnlessOKOrNotModified returns the response, the remote address, and an error if the given Response's status code is anything
// but 200 OK/ 304 Not Modified. This includes reading the Response.Body and Closing it. Otherwise, the given response, the remote
// address, and a nil error are returned.
func (to *Session) errUnlessOKOrNotModified(resp *http.Response, remoteAddr net.Addr, err error, path string) (*http.Response, net.Addr, error) {
	if err != nil {
		return resp, remoteAddr, err
	}
	if resp.StatusCode < 300 || resp.StatusCode == 304 {
		return resp, remoteAddr, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotImplemented {
		return resp, remoteAddr, ErrNotImplemented
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return resp, remoteAddr, readErr
	}
	return resp, remoteAddr, errors.New(resp.Status + "[" + strconv.Itoa(resp.StatusCode) + "] - Error requesting Traffic Ops " + to.getURL(path) + " " + string(body))
}

func (to *Session) getURL(path string) string { return to.URL + path }

type ReqF func(to *Session, method string, path string, body interface{}, header http.Header, response interface{}) (toclientlib.ReqInf, error)

type MidReqF func(ReqF) ReqF

// composeReqFuncs takes an initial request func and middleware, and
// returns a single ReqFunc to be called,
func composeReqFuncs(reqF ReqF, middleware []MidReqF) ReqF {
	// compose in reverse-order, which causes them to be applied in forward-order.
	for i := len(middleware) - 1; i >= 0; i-- {
		reqF = middleware[i](reqF)
	}
	return reqF
}

// makeRequestWithHeader marshals the response body (if non-nil), performs the HTTP request,
// and decodes the response into the given response pointer.
//
// Note processing on the following codes:
// 304 http.StatusNotModified  - Will return the 304 in ReqInf, a nil error, and a nil response.
//
// 401 http.StatusUnauthorized - Via to.request(), Same as 403 Forbidden.
// 403 http.StatusForbidden    - Via to.request()
//
//	Will try to log in again, and try the request again.
//	The second request is returned, even if it fails.
//
// To request the bytes without deserializing, pass a *[]byte response.
func makeRequestWithHeader(to *Session, method, path string, body interface{}, header http.Header, response interface{}) (toclientlib.ReqInf, error) {
	var remoteAddr net.Addr
	reqInf := toclientlib.ReqInf{CacheHitStatus: toclientlib.CacheHitStatusMiss, RemoteAddr: remoteAddr}
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return reqInf, errors.New("marshalling request body: " + err.Error())
		}
	}
	reqInf, err = to.TOClient.Req(method, path, reqBody, header, &response)
	reqInf.RemoteAddr = remoteAddr
	if err != nil {
		return reqInf, errors.New("requesting from Traffic Ops: " + err.Error())
	}
	return reqInf, nil
}

func (to *Session) get(path string, header http.Header, response interface{}) (toclientlib.ReqInf, error) {
	return to.TOClient.Req(http.MethodGet, path, nil, header, response)
}

func (to *Session) post(path string, body interface{}, header http.Header, response interface{}) (toclientlib.ReqInf, error) {
	return to.TOClient.Req(http.MethodPost, path, body, header, response)
}

func (to *Session) put(path string, body interface{}, header http.Header, response interface{}) (toclientlib.ReqInf, error) {
	return to.TOClient.Req(http.MethodPut, path, body, header, response)
}

func (to *Session) del(path string, header http.Header, response interface{}) (toclientlib.ReqInf, error) {
	return to.TOClient.Req(http.MethodDelete, path, nil, header, response)
}
