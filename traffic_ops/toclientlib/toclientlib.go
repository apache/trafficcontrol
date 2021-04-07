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

// Package toclientlib provides shared symbols for Traffic Ops Go clients.
package toclientlib

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptrace"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"golang.org/x/net/publicsuffix"
)

// Login authenticates with Traffic Ops and returns the client object.
//
// Returns the logged in client, the remote address of Traffic Ops which was translated and used to log in, and any error. If the error is not nil, the remote address may or may not be nil, depending whether the error occurred before the login request.
//
// apiVersions is the list of API versions to be supported. This should generally be provided by the specific client version wrapping this library.
//
// See ClientOpts for details about options, which options are required, and how they behave.
//
func Login(url, user, pass string, opts ClientOpts, apiVersions []string) (*TOClient, ReqInf, error) {
	if strings.TrimSpace(opts.UserAgent) == "" {
		return nil, ReqInf{}, errors.New("opts.UserAgent is required")
	}
	if opts.RequestTimeout == 0 {
		opts.RequestTimeout = DefaultTimeout
	}
	if opts.APIVersionCheckInterval == 0 {
		opts.APIVersionCheckInterval = DefaultAPIVersionCheckInterval
	}

	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, ReqInf{}, errors.New("creating cookie jar: " + err.Error())
	}

	to := NewClient(user, pass, url, opts.UserAgent, &http.Client{
		Timeout: opts.RequestTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: opts.Insecure},
		},
		Jar: jar,
	}, apiVersions)

	if !opts.ForceLatestAPI {
		to.latestSupportedAPI = apiVersions[0]
	}

	to.forceLatestAPI = opts.ForceLatestAPI
	to.apiVerCheckInterval = opts.APIVersionCheckInterval

	remoteAddr, err := to.login()
	if err != nil {
		return nil, ReqInf{RemoteAddr: remoteAddr}, errors.New("logging in: " + err.Error())
	}
	return to, ReqInf{RemoteAddr: remoteAddr}, nil
}

// ClientOpts is the options to configure the creation of the Client.
//
// This exists to allow adding new features without a breaking change to the Login function.
// Users should understand this, and understand that upgrading their library may result in new options that their application doesn't know to use.
// New fields should always behave as-before if their value is the default.
type ClientOpts struct {
	// ForceLatestAPI will cause Login to return an error if the latest minor API in the client
	// is not supported by the Traffic Ops Server.
	//
	// Note this was the behavior of all Traffic Ops client functions prior to the Login function.
	//
	// If this is false or unset, login will determine the latest minor version supported, and use that for all requests.
	//
	// Be aware, this means client fields unknown to the server will always be default-initialized.
	// For example, suppose the field Foo was added in API 3.1, the client code is 3.1, and the server is 3.0.
	// Then, the field Foo will always be nil or the default value.
	// Client applications must understand this, and code processing the new feature Foo must be able to
	// process default or nil values, understanding they may indicate a server version before the feature was added.
	//
	ForceLatestAPI bool

	// Insecure is whether to ignore HTTPS certificate errors with Traffic Ops.
	// Setting this on production systems is strongly discouraged.
	Insecure bool

	// RequestTimeout is the HTTP timeout for Traffic Ops requests.
	// If 0 or not explicitly set, DefaultTimeout will be used.
	RequestTimeout time.Duration

	// UserAgent is the HTTP User Agent to use set when communicating with Traffic Ops.
	// This field is required, and Login will fail if it is unset or the empty string.
	UserAgent string

	// APIVersionCheckInterval is how often to try a newer Traffic Ops API Version.
	// This allows clients to get new Traffic Ops features after Traffic Ops is upgraded
	// without requiring a restart or new client.
	//
	// If 0 or not explicitly set, DefaultAPIVersionCheckInterval will be used.
	// To disable, set to a very high value (like 100 years).
	//
	// This has no effect if ForceLatestAPI is true.
	APIVersionCheckInterval time.Duration
}

// Client is a Traffic Ops client, with generic functions to be used by any specific client.
type TOClient struct {
	UserName     string
	Password     string
	URL          string
	Client       *http.Client
	UserAgentStr string

	latestSupportedAPI string
	// forceLatestAPI is whether to forcibly always use the latest API version known to this client.
	// This should only ever be set by ClientOpts.ForceLatestAPI.
	forceLatestAPI bool
	// lastAPIVerCheck is the last time the Client tried to get a newer API version from TO.
	// Used internally to decide whether to try again.
	lastAPIVerCheck time.Time
	// apiVerCheckInterval is how often to try a newer Traffic Ops API, in case Traffic Ops was upgraded.
	// This should only ever be set by ClientOpts.APIVersionCheckInterval.
	apiVerCheckInterval time.Duration

	// apiVersions is the list of support Traffic Ops versions.
	// This must be provided on construction, typically by the client wrapping this lib.
	apiVersions []string
}

func NewClient(user, password, url, userAgent string, client *http.Client, apiVersions []string) *TOClient {
	return &TOClient{
		UserName:     user,
		Password:     password,
		URL:          url,
		Client:       client,
		UserAgentStr: userAgent,
		apiVersions:  apiVersions,
	}
}

const DefaultTimeout = time.Second * 30
const DefaultAPIVersionCheckInterval = time.Second * 60

// HTTPError is returned on Update Client failure.
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

// login tries to log in to Traffic Ops, and set the auth cookie in the Client. Returns the IP address of the remote Traffic Ops.
func (to *TOClient) login() (net.Addr, error) {
	path := "/user/login"
	body := tc.UserCredentials{Username: to.UserName, Password: to.Password}
	alerts := tc.Alerts{}

	// Can't use req() because it retries login failures, which would be an infinite loop.
	reqF := composeReqFuncs(makeRequestWithHeader, []MidReqF{reqTryLatest, reqFallback, reqAPI})

	reqInf, err := reqF(to, http.MethodPost, path, body, nil, &alerts)
	if err != nil {
		return reqInf.RemoteAddr, fmt.Errorf("Login error %v, alerts string: %+v", err, alerts)
	}

	success := false
	for _, alert := range alerts.Alerts {
		if alert.Level == "success" && alert.Text == "Successfully logged in." {
			success = true
			break
		}
	}

	if !success {
		return reqInf.RemoteAddr, fmt.Errorf("Login failed, alerts string: %+v", alerts)
	}

	return reqInf.RemoteAddr, nil
}

func (to *TOClient) loginWithToken(token []byte) (net.Addr, error) {
	path := to.APIBase() + "/user/login/token"
	resp, remoteAddr, err := to.RawRequestWithHdr(http.MethodPost, path, token, nil)
	resp, remoteAddr, err = to.errUnlessOKOrNotModified(resp, remoteAddr, err, path)
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
func (to *TOClient) logout() (net.Addr, error) {
	credentials, err := loginCreds(to.UserName, to.Password)
	if err != nil {
		return nil, errors.New("creating login credentials: " + err.Error())
	}

	path := to.APIBase() + "/user/logout"
	resp, remoteAddr, err := to.RawRequestWithHdr("POST", path, credentials, nil)
	resp, remoteAddr, err = to.errUnlessOKOrNotModified(resp, remoteAddr, err, path)
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

// Login to traffic_ops, the response should set the cookie for this client
// automatically. Start with
//     to := traffic_ops.Login("user", "passwd", true)
// subsequent calls like to.GetData("datadeliveryservice") will be authenticated.
// Returns the logged in client, the remote address of Traffic Ops which was translated and used to log in, and any error. If the error is not nil, the remote address may or may not be nil, depending whether the error occurred before the login request.
//
// The apiVersions is the list of API versions supported in this client. This should generally be provided by the client package wrapping this package.
func LoginWithAgent(
	toURL string,
	toUser string,
	toPasswd string,
	insecure bool,
	userAgent string,
	requestTimeout time.Duration,
	apiVersions []string,
) (*TOClient, net.Addr, error) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	jar, err := cookiejar.New(&options)
	if err != nil {
		return nil, nil, err
	}

	to := NewClient(toUser, toPasswd, toURL, userAgent, &http.Client{
		Timeout: requestTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
		Jar: jar,
	}, apiVersions)

	remoteAddr, err := to.login()
	if err != nil {
		return nil, remoteAddr, errors.New("logging in: " + err.Error())
	}
	return to, remoteAddr, nil
}

// The apiVersions is the list of API versions supported in this client. This should generally be provided by the client package wrapping this package.
func LoginWithToken(
	toURL string,
	token string,
	insecure bool,
	userAgent string,
	requestTimeout time.Duration,
	apiVersions []string,
) (*TOClient, net.Addr, error) {
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

	to := NewClient("", "", toURL, userAgent, &client, apiVersions)
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

// Logout of Traffic Ops.
//
// The apiVersions is the list of API versions supported in this client. This should generally be provided by the client package wrapping this package.
//
func LogoutWithAgent(
	toURL string,
	toUser string,
	toPasswd string,
	insecure bool,
	userAgent string,
	requestTimeout time.Duration,
	apiVersions []string,
) (*TOClient, net.Addr, error) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	jar, err := cookiejar.New(&options)
	if err != nil {
		return nil, nil, err
	}

	to := NewClient(toUser, toPasswd, toURL, userAgent, &http.Client{
		Timeout: requestTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
		Jar: jar,
	}, apiVersions)

	remoteAddr, err := to.logout()
	if err != nil {
		return nil, remoteAddr, errors.New("logging out: " + err.Error())
	}
	return to, remoteAddr, nil
}

// NewNoAuthClient returns a new Client without logging in
// this can be used for querying unauthenticated endpoints without requiring a login
// The apiVersions is the list of API versions supported in this client. This should generally be provided by the client package wrapping this package.
func NewNoAuthClient(
	toURL string,
	insecure bool,
	userAgent string,
	requestTimeout time.Duration,
	apiVersions []string,
) *TOClient {
	return NewClient("", "", toURL, userAgent, &http.Client{
		Timeout: requestTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
	}, apiVersions)
}

func ErrIsNotImplemented(err error) bool {
	return err != nil && strings.Contains(err.Error(), ErrNotImplemented.Error()) // use string.Contains in case context was added to the error
}

// ErrNotImplemented is returned when Traffic Ops returns a 501 Not Implemented
// Users should check ErrIsNotImplemented rather than comparing directly, in case context was added.
var ErrNotImplemented = errors.New("Traffic Ops Server returned 'Not Implemented', this client is probably newer than Traffic Ops, and you probably need to either upgrade Traffic Ops, or use a client whose version matches your Traffic Ops version.")

// errUnlessOKOrNotModified returns the response, the remote address, and an error if the given Response's status code is anything
// but 200 OK/ 304 Not Modified. This includes reading the Response.Body and Closing it. Otherwise, the given response, the remote
// address, and a nil error are returned.
func (to *TOClient) errUnlessOKOrNotModified(resp *http.Response, remoteAddr net.Addr, err error, path string) (*http.Response, net.Addr, error) {
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

func (to *TOClient) getURL(path string) string { return to.URL + path }

type ReqF func(to *TOClient, method string, path string, body interface{}, header http.Header, response interface{}) (ReqInf, error)

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

// reqTryLatest will re-set to.latestSupportedAPI to the latest, if it's less than the latest and to.apiVerCheckInterval has passed.
// This does not fallback, so it should generally be composed with reqFallback.
func reqTryLatest(reqF ReqF) ReqF {
	return func(to *TOClient, method string, path string, body interface{}, header http.Header, response interface{}) (ReqInf, error) {
		if to.apiVerCheckInterval == 0 {
			// Client could have been default-initialized rather than created with a func, so we need to check here, not just in login funcs.
			to.apiVerCheckInterval = DefaultAPIVersionCheckInterval
		}

		if !to.forceLatestAPI && time.Since(to.lastAPIVerCheck) >= to.apiVerCheckInterval {
			// if it's been apiVerCheckInterval since the last version check,
			// set the version to the latest (and fall back again, if necessary)
			to.latestSupportedAPI = to.apiVersions[0]

			// Set the last version check to far in the future, and then
			// defer setting the last check until this function returns,
			// so that if fallback takes longer than the interval,
			// the recursive calls to this function don't end up forever retrying the latest.
			to.lastAPIVerCheck = time.Now().Add(time.Hour * 24 * 365)
			defer func() { to.lastAPIVerCheck = time.Now() }()
		}
		return reqF(to, method, path, body, header, response)
	}
}

// reqLogin makes the request, and if it fails, tries to log in again then makes the request again.
// If the login fails, the original response is returned.
// If the second request fails, it's returned. Login is only tried once.
// This is designed to handle expired sessions, when the time between requests is longer than the session expiration;
// it does not do perpetual retry.
func reqLogin(reqF ReqF) ReqF {
	return func(to *TOClient, method string, path string, body interface{}, header http.Header, response interface{}) (ReqInf, error) {
		inf, err := reqF(to, method, path, body, header, response)
		if inf.StatusCode != http.StatusUnauthorized && inf.StatusCode != http.StatusForbidden {
			return inf, err
		}
		if _, lerr := to.login(); lerr != nil {
			return inf, err
		}
		return reqF(to, method, path, body, header, response)
	}
}

// reqFallback calls reqF, and if Traffic Ops doesn't support the latest version,
// falls back to the previous and retries, recursively.
// If all supported versions fail, the last response error is returned.
func reqFallback(reqF ReqF) ReqF {
	var fallbackFunc func(to *TOClient, method string, path string, body interface{}, header http.Header, response interface{}) (ReqInf, error)
	fallbackFunc = func(to *TOClient, method string, path string, body interface{}, header http.Header, response interface{}) (ReqInf, error) {
		inf, err := reqF(to, method, path, body, header, response)
		if err == nil {
			return inf, err
		}
		if !ErrIsNotImplemented(err) ||
			to.forceLatestAPI {
			return inf, err
		}

		apiVersions := to.apiVersions

		nextAPIVerI := int(math.MaxInt32) - 1
		for verI, ver := range apiVersions {
			if to.latestSupportedAPI == ver {
				nextAPIVerI = verI
				break
			}
		}
		nextAPIVerI = nextAPIVerI + 1
		if nextAPIVerI >= len(apiVersions) {
			return inf, err // we're already on the oldest minor supported, and the server doesn't support it.
		}
		to.latestSupportedAPI = apiVersions[nextAPIVerI]

		return fallbackFunc(to, method, path, body, header, response)
	}
	return fallbackFunc
}

// reqAPI calls reqF with a path not including the /api/x prefix,
// and adds the API version from the Client.
//
// For example, path should be like '/deliveryservices'
// and this will request '/api/3.1/deliveryservices'.
//
func reqAPI(reqF ReqF) ReqF {
	return func(to *TOClient, method string, path string, body interface{}, header http.Header, response interface{}) (ReqInf, error) {
		path = to.APIBase() + path
		return reqF(to, method, path, body, header, response)
	}
}

// makeRequestWithHeader marshals the response body (if non-nil), performs the HTTP request,
// and decodes the response into the given response pointer.
//
// Note processing on the following codes:
// 304 http.StatusNotModified  - Will return the 304 in ReqInf, a nil error, and a nil response.
//
// 401 http.StatusUnauthorized - Via to.request(), Same as 403 Forbidden.
// 403 http.StatusForbidden    - Via to.request()
//                               Will try to log in again, and try the request again.
//                               The second request is returned, even if it fails.
//
// To request the bytes without deserializing, pass a *[]byte response.
//
func makeRequestWithHeader(to *TOClient, method, path string, body interface{}, header http.Header, response interface{}) (ReqInf, error) {
	var remoteAddr net.Addr
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return reqInf, errors.New("marshalling request body: " + err.Error())
		}
	}
	resp, remoteAddr, err := to.request(method, path, reqBody, header)
	reqInf.RemoteAddr = remoteAddr
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return reqInf, nil
		}
		defer log.Close(resp.Body, "unable to close response body")
	}
	if err != nil {
		return reqInf, errors.New("requesting from Traffic Ops: " + err.Error())
	}

	if btsPtr, isBytes := response.(*[]byte); isBytes {
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return reqInf, errors.New("reading response body: " + err.Error())
		}
		*btsPtr = bts
		return reqInf, nil
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return reqInf, errors.New("decoding response body: " + err.Error())
	}
	return reqInf, nil
}

func (to *TOClient) Req(method string, path string, body interface{}, header http.Header, response interface{}) (ReqInf, error) {
	reqF := composeReqFuncs(makeRequestWithHeader, []MidReqF{reqTryLatest, reqFallback, reqAPI, reqLogin})
	return reqF(to, method, path, body, header, response)
}

// request performs the HTTP request to Traffic Ops, trying to refresh the cookie if an Unauthorized or Forbidden code is received. It only tries once. If the login fails, the original Unauthorized/Forbidden response is returned. If the login succeeds and the subsequent re-request fails, the re-request's response is returned even if it's another Unauthorized/Forbidden.
// Returns the response, the remote address of the Traffic Ops instance used, and any error.
// The returned net.Addr is guaranteed to be either nil or valid, even if the returned error is not nil. Callers are encouraged to check and use the net.Addr if an error is returned, and use the remote address in their own error messages. This violates the Go idiom that a non-nil error implies all other values are undefined, but it's more straightforward than alternatives like typecasting.
func (to *TOClient) request(method, path string, body []byte, header http.Header) (*http.Response, net.Addr, error) {
	r, remoteAddr, err := to.RawRequestWithHdr(method, path, body, header)
	if err != nil {
		return r, remoteAddr, err
	}
	if r.StatusCode != http.StatusUnauthorized && r.StatusCode != http.StatusForbidden {
		return to.errUnlessOKOrNotModified(r, remoteAddr, err, path)
	}
	if _, lerr := to.login(); lerr != nil {
		return to.errUnlessOKOrNotModified(r, remoteAddr, err, path) // if re-logging-in fails, return the original request's response
	}

	// return second request, even if it's another Unauthorized or Forbidden.
	r, remoteAddr, err = to.RawRequestWithHdr(method, path, body, header)
	return to.errUnlessOKOrNotModified(r, remoteAddr, err, path)
}

func (to *TOClient) RawRequestWithHdr(method, path string, body []byte, header http.Header) (*http.Response, net.Addr, error) {
	url := to.getURL(path)

	var req *http.Request
	var err error
	remoteAddr := net.Addr(nil)

	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
		if err != nil {
			return nil, remoteAddr, err
		}
		if header != nil {
			req.Header = header.Clone()
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, remoteAddr, err
		}
		if header != nil {
			req.Header = header.Clone()
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
	return resp, remoteAddr, err
}

// RawRequest performs the actual HTTP request to Traffic Ops, simply, without trying to refresh the cookie if an Unauthorized code is returned.
// Returns the response, the remote address of the Traffic Ops instance used, and any error.
// The returned net.Addr is guaranteed to be either nil or valid, even if the returned error is not nil. Callers are encouraged to check and use the net.Addr if an error is returned, and use the remote address in their own error messages. This violates the Go idiom that a non-nil error implies all other values are undefined, but it's more straightforward than alternatives like typecasting.
// Deprecated: RawRequest will be removed in 6.0. Use RawRequestWithHdr.
func (to *TOClient) RawRequest(method, path string, body []byte) (*http.Response, net.Addr, error) {
	return to.RawRequestWithHdr(method, path, body, nil)
}

type ReqInf struct {
	// CacheHitStatus is deprecated and will be removed in the next major version.
	CacheHitStatus CacheHitStatus
	RemoteAddr     net.Addr
	StatusCode     int
}

// CacheHitStatus is deprecated and will be removed in the next major version.
type CacheHitStatus string

// CacheHitStatusHit is deprecated and will be removed in the next major version.
const CacheHitStatusHit = CacheHitStatus("hit")

// CacheHitStatusExpired is deprecated and will be removed in the next major version.
const CacheHitStatusExpired = CacheHitStatus("expired")

// CacheHitStatusMiss is deprecated and will be removed in the next major version.
const CacheHitStatusMiss = CacheHitStatus("miss")

// CacheHitStatusInvalid is deprecated and will be removed in the next major version.
const CacheHitStatusInvalid = CacheHitStatus("")

// String is deprecated and will be removed in the next major version.
func (s CacheHitStatus) String() string {
	return string(s)
}

// StringToCacheHitStatus is deprecated and will be removed in the next major version.
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
