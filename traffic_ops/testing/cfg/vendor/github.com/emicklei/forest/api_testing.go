package forest

import "net/http"

// APITesting provides functions to call a REST api and validate its responses.
type APITesting struct {
	BaseURL string
	client  *http.Client
}

// NewClient returns a new ApiTesting that can be used to send Http requests.
func NewClient(baseURL string, httpClient *http.Client) *APITesting {
	return &APITesting{
		BaseURL: baseURL,
		client:  httpClient,
	}
}

// GET sends a Http request using a config (headers,...)
// The request is logged and any sending error will fail the test.
func (a *APITesting) GET(t T, config *RequestConfig) *http.Response {
	httpReq, err := http.NewRequest("GET", a.BaseURL+config.pathAndQuery(), nil)
	if err != nil {
		logfatal(t, sfatalf("GET: invalid Url:%s", a.BaseURL+config.pathAndQuery()))
	}
	setBasicAuth(config, httpReq)
	copyHeaders(config.HeaderMap, httpReq.Header)
	Logf(t, "\n%v %v %v", httpReq.Method, httpReq.URL, headersString(httpReq.Header))
	resp, err := a.client.Do(httpReq)
	CheckError(t, err)
	return ensureResponse(httpReq, resp)
}

// POST sends a Http request using a config (headers,body,...)
// The request is logged and any sending error will fail the test.
func (a *APITesting) POST(t T, config *RequestConfig) *http.Response {
	httpReq, err := http.NewRequest("POST", a.BaseURL+config.pathAndQuery(), config.BodyReader)
	if err != nil {
		logfatal(t, sfatalf("POST: invalid Url:%s", a.BaseURL+config.pathAndQuery()))
	}
	setBasicAuth(config, httpReq)
	copyHeaders(config.HeaderMap, httpReq.Header)
	Logf(t, "\n%v %v %v", httpReq.Method, httpReq.URL, headersString(httpReq.Header))
	resp, err := a.client.Do(httpReq)
	CheckError(t, err)
	return ensureResponse(httpReq, resp)
}

// PUT sends a Http request using a config (headers,body,...)
// The request is logged and any sending error will fail the test.
func (a *APITesting) PUT(t T, config *RequestConfig) *http.Response {
	httpReq, err := http.NewRequest("PUT", a.BaseURL+config.pathAndQuery(), config.BodyReader)
	if err != nil {
		logfatal(t, sfatalf("PUT: invalid Url:%s", a.BaseURL+config.pathAndQuery()))
	}
	setBasicAuth(config, httpReq)
	copyHeaders(config.HeaderMap, httpReq.Header)
	Logf(t, "\n%v %v %v", httpReq.Method, httpReq.URL, headersString(httpReq.Header))
	resp, err := a.client.Do(httpReq)
	CheckError(t, err)
	return ensureResponse(httpReq, resp)
}

// DELETE sends a Http request using a config (headers,...)
// The request is logged and any sending error will fail the test.
func (a *APITesting) DELETE(t T, config *RequestConfig) *http.Response {
	httpReq, err := http.NewRequest("DELETE", a.BaseURL+config.pathAndQuery(), nil)
	if err != nil {
		logfatal(t, sfatalf("DELETE: invalid Url:%s", a.BaseURL+config.pathAndQuery()))
	}
	setBasicAuth(config, httpReq)
	copyHeaders(config.HeaderMap, httpReq.Header)
	Logf(t, "\n%v %v %v", httpReq.Method, httpReq.URL, headersString(httpReq.Header))
	resp, err := a.client.Do(httpReq)
	CheckError(t, err)
	return ensureResponse(httpReq, resp)
}

// PATCH sends a Http request using a config (headers,...)
// The request is logged and any sending error will fail the test.
func (a *APITesting) PATCH(t T, config *RequestConfig) *http.Response {
	httpReq, err := http.NewRequest("PATCH", a.BaseURL+config.pathAndQuery(), config.BodyReader)
	if err != nil {
		logfatal(t, sfatalf("PATCH: invalid Url:%s", a.BaseURL+config.pathAndQuery()))
	}
	setBasicAuth(config, httpReq)
	copyHeaders(config.HeaderMap, httpReq.Header)
	Logf(t, "\n%v %v %v", httpReq.Method, httpReq.URL, headersString(httpReq.Header))
	resp, err := a.client.Do(httpReq)
	CheckError(t, err)
	return ensureResponse(httpReq, resp)
}

// Do sends a Http request using a Http method (GET,PUT,POST,....) and config (headers,...)
// The request is not logged and any URL build error or send error will be returned.
func (a *APITesting) Do(method string, config *RequestConfig) (*http.Response, error) {
	httpReq, err := http.NewRequest(method, a.BaseURL+config.pathAndQuery(), config.BodyReader)
	if err != nil {
		return nil, err
	}
	setBasicAuth(config, httpReq)
	copyHeaders(config.HeaderMap, httpReq.Header)
	resp, err := a.client.Do(httpReq)
	return ensureResponse(httpReq, resp), err
}
