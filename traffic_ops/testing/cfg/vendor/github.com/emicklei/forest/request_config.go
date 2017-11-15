package forest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RequestConfig holds additional information to construct a Http request.
type RequestConfig struct {
	URI            string
	BodyReader     io.Reader
	HeaderMap      http.Header
	Values         url.Values
	User, Password string
}

// Path is an alias for NewConfig
var Path = NewConfig

// NewConfig returns a new RequestConfig with initialized empty headers and query parameters.
// See Path for an explanation of the function parameters.
func NewConfig(pathTemplate string, pathParams ...interface{}) *RequestConfig {
	cfg := &RequestConfig{
		HeaderMap: http.Header{},
		Values:    url.Values{},
	}
	cfg.Path(pathTemplate, pathParams...)
	return cfg
}

// Do calls the one-argument function parameter with the receiver.
// This allows for custom convenience functions without breaking the fluent programming style.
func (r *RequestConfig) Do(block func(config *RequestConfig)) *RequestConfig {
	block(r)
	return r
}

// Path sets the URL path with optional path parameters.
// format example: /v1/persons/{param}/ + 42 => /v1/persons/42
// format example: /v1/persons/:param/ + 42 => /v1/persons/42
// format example: /v1/assets/*rest +  js/some/file.js => /v1/assets/js/some/file.js
func (r *RequestConfig) Path(pathTemplate string, pathParams ...interface{}) *RequestConfig {
	var uri bytes.Buffer
	p := 0
	tokens := strings.Split(pathTemplate, "/")
	for i, each := range tokens {
		if len(each) == 0 && i == 0 { // skip leading space
			continue
		}
		uri.WriteString("/")

		if strings.HasPrefix(each, "*") {
			// treat remainder as is
			uri.WriteString(fmt.Sprintf("%v", pathParams[p]))
			break
		}

		if strings.HasPrefix(each, ":") ||
			(strings.HasPrefix(each, "{") && strings.HasSuffix(each, "}")) {
			if p == len(pathParams) {
				// abort
				r.URI = pathTemplate
				return r
			}
			uri.WriteString(fmt.Sprintf("%v", pathParams[p]))
			p++
		} else {
			uri.WriteString(each)
		}
	}
	// need to do path encoding
	r.URI = URLPathEncode(uri.String())
	return r
}

// BasicAuth sets the credentials for Basic Authentication (if username is not empty)
func (r *RequestConfig) BasicAuth(username, password string) *RequestConfig {
	r.User = username
	r.Password = password
	return r
}

// Query adds a name=value pair to the list of query parameters.
func (r *RequestConfig) Query(name string, value interface{}) *RequestConfig {
	r.Values.Add(name, fmt.Sprintf("%v", value))
	return r
}

// Header adds a name=value pair to the list of header parameters.
func (r *RequestConfig) Header(name, value string) *RequestConfig {
	r.HeaderMap.Add(name, value)
	return r
}

// Body sets the playload as is. No content type is set.
// It sets the BodyReader field of the RequestConfig.
func (r *RequestConfig) Body(body string) *RequestConfig {
	r.BodyReader = strings.NewReader(body)
	return r
}

func (r *RequestConfig) pathAndQuery() string {
	if len(r.Values) == 0 {
		return r.URI
	}
	return fmt.Sprintf("%s?%s", r.URI, r.Values.Encode())
}

// Content encodes (marshals) the payload conform the content type given.
// If the payload is already a string (JSON,XML,plain) then it is used as is.
// Supported Content-Type values for marshalling: application/json, application/xml, text/plain
// Payload can also be a slice of bytes; use application/octet-stream in that case.
// It sets the BodyReader field of the RequestConfig.
func (r *RequestConfig) Content(payload interface{}, contentType string) *RequestConfig {
	r.Header("Content-Type", contentType)
	if payloadAsIs, ok := payload.(string); ok {
		r.BodyReader = strings.NewReader(payloadAsIs)
		return r
	}
	if strings.Index(contentType, "application/json") != -1 {
		data, err := json.Marshal(payload)
		if err != nil {
			r.Body(fmt.Sprintf("json marshal failed:%v", err))
			return r
		}
		r.BodyReader = bytes.NewReader(data)
		return r
	}
	if strings.Index(contentType, "application/xml") != -1 {
		data, err := xml.Marshal(payload)
		if err != nil {
			r.Body(fmt.Sprintf("xml marshal failed:%v", err))
			return r
		}
		r.BodyReader = bytes.NewReader(data)
		return r
	}
	if strings.Index(contentType, "text/plain") != -1 {
		content, ok := payload.(string)
		if !ok {
			r.Body(fmt.Sprintf("content is not a string:%v", payload))
			return r
		}
		r.BodyReader = strings.NewReader(content)
		return r
	}
	bits, ok := payload.([]byte)
	if ok {
		r.BodyReader = bytes.NewReader(bits)
		return r
	}
	r.Body(fmt.Sprintf("cannot encode payload, unknown content type:%s", contentType))
	return r
}

// Read sets the BodyReader for content to send with the request.
func (r *RequestConfig) Read(bodyReader io.Reader) *RequestConfig {
	r.BodyReader = bodyReader
	return r
}
