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
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// RequestOptions is the set of options commonly available to pass to methods
// of a Session.
type RequestOptions struct {
	// Any and all extra HTTP headers to pass in the request.
	Header http.Header
	// Any and all query parameters to pass in the request.
	QueryParameters url.Values
}

// NewRequestOptions returns a RequestOptions object with initialized, empty Header
// and QueryParameters.
func NewRequestOptions() RequestOptions {
	return RequestOptions{
		Header:          http.Header{},
		QueryParameters: url.Values{},
	}
}

// Login authenticates with Traffic Ops and returns the client object.
//
// Returns the logged in client, the remote address of Traffic Ops which was translated and used to log in, and any error. If the error is not nil, the remote address may or may not be nil, depending whether the error occurred before the login request.
//
// See ClientOpts for details about options, which options are required, and how they behave.
func Login(url, user, pass string, opts Options) (*Session, toclientlib.ReqInf, error) {
	cl, reqInf, err := toclientlib.Login(url, user, pass, opts.ClientOpts, apiVersions())
	if err != nil {
		return nil, reqInf, err
	}
	return &Session{TOClient: *cl}, reqInf, err
}

// Options is the options to configure the creation of the Client.
//
// This exists to allow adding new features without a breaking change to the
// Login function. Users should understand this, and understand that upgrading
// their library may result in new options that their application doesn't know
// to use. New fields should always behave as-before if their value is the
// default.
type Options struct {
	toclientlib.ClientOpts
}

// Session is a Traffic Ops client.
type Session struct {
	toclientlib.TOClient
}

// NewSession constructs a new, unauthenticated Session using the provided information.
func NewSession(user, password, url, userAgent string, client *http.Client, useCache bool) *Session {
	return &Session{
		TOClient: *toclientlib.NewClient(user, password, url, userAgent, client, apiVersions()),
	}
}

// LoginWithAgent creates a new authenticated session with a Traffic Ops
// server. The session cookie should be set automatically in the returned
// Session, so that subsequent calls are properly authenticated without further
// manual intervention.
//
// Returns the logged in client, the remote address of Traffic Ops which was
// translated and used to log in, and any error that occurred along the way. If
// the error is not nil, the remote address may or may not be nil, depending on
// whether the error occurred before the login request.
func LoginWithAgent(toURL string, toUser string, toPasswd string, insecure bool, userAgent string, useCache bool, requestTimeout time.Duration) (*Session, net.Addr, error) {
	cl, ip, err := toclientlib.LoginWithAgent(toURL, toUser, toPasswd, insecure, userAgent, requestTimeout, apiVersions())
	if err != nil {
		return nil, nil, err
	}
	return &Session{TOClient: *cl}, ip, err
}

// LoginWithToken functions identically to LoginWithAgent, but uses token-based
// authentication rather than a username/password pair.
func LoginWithToken(toURL string, token string, insecure bool, userAgent string, useCache bool, requestTimeout time.Duration) (*Session, net.Addr, error) {
	cl, ip, err := toclientlib.LoginWithToken(toURL, token, insecure, userAgent, requestTimeout, apiVersions())
	if err != nil {
		return nil, nil, err
	}
	return &Session{TOClient: *cl}, ip, err
}

// LogoutWithAgent constructs an authenticated Session - exactly like
// LoginWithAgent - and then immediately calls the '/logout' API endpoint to
// end the session.
func LogoutWithAgent(toURL string, toUser string, toPasswd string, insecure bool, userAgent string, useCache bool, requestTimeout time.Duration) (*Session, net.Addr, error) {
	cl, ip, err := toclientlib.LogoutWithAgent(toURL, toUser, toPasswd, insecure, userAgent, requestTimeout, apiVersions())
	if err != nil {
		return nil, nil, err
	}
	return &Session{TOClient: *cl}, ip, err
}

// NewNoAuthSession returns a new Session without logging in.
// this can be used for querying unauthenticated endpoints without requiring a login.
func NewNoAuthSession(toURL string, insecure bool, userAgent string, useCache bool, requestTimeout time.Duration) *Session {
	return &Session{TOClient: *toclientlib.NewNoAuthClient(toURL, insecure, userAgent, requestTimeout, apiVersions())}
}

func (to *Session) get(path string, opts RequestOptions, response interface{}) (toclientlib.ReqInf, error) {
	return to.req(http.MethodGet, path, opts, nil, response)
}

func (to *Session) post(path string, opts RequestOptions, body, response interface{}) (toclientlib.ReqInf, error) {
	return to.req(http.MethodPost, path, opts, body, response)
}

func (to *Session) put(path string, opts RequestOptions, body, response interface{}) (toclientlib.ReqInf, error) {
	return to.req(http.MethodPut, path, opts, body, response)
}

func (to *Session) del(path string, opts RequestOptions, response interface{}) (toclientlib.ReqInf, error) {
	return to.req(http.MethodDelete, path, opts, nil, response)
}

func (to *Session) req(method, path string, opts RequestOptions, body, response interface{}) (toclientlib.ReqInf, error) {
	if len(opts.QueryParameters) > 0 {
		path += "?" + opts.QueryParameters.Encode()
	}
	return to.TOClient.Req(method, path, body, opts.Header, response)
}
