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
	"time"

	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

// Login authenticates with Traffic Ops and returns the client object.
//
// Returns the logged in client, the remote address of Traffic Ops which was translated and used to log in, and any error. If the error is not nil, the remote address may or may not be nil, depending whether the error occurred before the login request.
//
// See ClientOpts for details about options, which options are required, and how they behave.
//
func Login(url, user, pass string, opts ClientOpts) (*Session, toclientlib.ReqInf, error) {
	cl, ip, err := toclientlib.Login(url, user, pass, opts.ClientOpts, apiVersions())
	if err != nil {
		return nil, toclientlib.ReqInf{}, err
	}
	return &Session{TOClient: *cl}, ip, err
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
//     to := traffic_ops.Login("user", "passwd", true)
// subsequent calls like to.GetData("datadeliveryservice") will be authenticated.
// Returns the logged in client, the remote address of Traffic Ops which was translated and used to log in, and any error. If the error is not nil, the remote address may or may not be nil, depending whether the error occurred before the login request.
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

// Logout of traffic_ops
func LogoutWithAgent(toURL string, toUser string, toPasswd string, insecure bool, userAgent string, useCache bool, requestTimeout time.Duration) (*Session, net.Addr, error) {
	cl, ip, err := toclientlib.LogoutWithAgent(toURL, toUser, toPasswd, insecure, userAgent, requestTimeout, apiVersions())
	if err != nil {
		return nil, nil, err
	}
	return &Session{TOClient: *cl}, ip, err
}

// NewNoAuthSession returns a new Session without logging in
// this can be used for querying unauthenticated endpoints without requiring a login
func NewNoAuthSession(toURL string, insecure bool, userAgent string, useCache bool, requestTimeout time.Duration) *Session {
	return &Session{TOClient: *toclientlib.NewNoAuthClient(toURL, insecure, userAgent, requestTimeout, apiVersions())}
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
