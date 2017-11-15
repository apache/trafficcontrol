package forest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Dump is a convenient method to log the full contents of a request and its response.
func Dump(t T, resp *http.Response) {
	// dump request
	var buffer bytes.Buffer
	buffer.WriteString("\n")
	buffer.WriteString(fmt.Sprintf("%v %v\n", resp.Request.Method, resp.Request.URL))
	for k, v := range resp.Request.Header {
		if len(k) > 0 {
			buffer.WriteString(fmt.Sprintf("%s : %v\n", k, strings.Join(v, ",")))
		}
	}
	if resp == nil {
		buffer.WriteString("-- no response --")
		Logf(t, buffer.String())
		return
	}
	// dump response
	buffer.WriteString(fmt.Sprintf("\n%s\n", resp.Status))
	for k, v := range resp.Header {
		if len(k) > 0 {
			buffer.WriteString(fmt.Sprintf("%s : %v\n", k, strings.Join(v, ",")))
		}
	}
	if resp.Body != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			if resp.StatusCode/100 == 3 {
				// redirect closes body ; nothing to read
				buffer.WriteString("\n")
			} else {
				buffer.WriteString(fmt.Sprintf("unable to read body:%v", err))
			}
		} else {
			if len(body) > 0 {
				buffer.WriteString("\n")
			}
			buffer.WriteString(string(body))
		}
		resp.Body.Close()
		// put the body back for re-reads
		resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	}
	buffer.WriteString("\n")
	Logf(t, buffer.String())
}

type skippeable interface {
	Skipf(string, ...interface{})
}

// SkipUnless will Skip the test unless the LABELS environment variable includes any of the provided labels.
//
//	LABELS=integration,nightly go test -v
//
func SkipUnless(t skippeable, labels ...string) {
	env := strings.Split(os.Getenv("LABELS"), ",")
	for _, each := range labels {
		for _, other := range env {
			if each == other {
				return
			}
		}
	}
	t.Skipf("skipped because provided LABELS=%v does not include any of %v", env, labels)
}

func headersString(h http.Header) string {
	if len(h) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", h)
}

func copyHeaders(from, to http.Header) {
	for k, list := range from {
		for _, v := range list {
			to.Set(k, v)
		}
	}
}

func setBasicAuth(config *RequestConfig, req *http.Request) {
	if len(config.User) > 0 {
		req.SetBasicAuth(config.User, config.Password)
	}
}

func ensureResponse(req *http.Request, resp *http.Response) *http.Response {
	if resp != nil {
		return resp
	}
	// wrap the request into an empty response ; Body must be nil
	wrap := new(http.Response)
	wrap.Request = req
	return wrap
}

func URLPathEncode(path string) string {
	buf := new(bytes.Buffer)
	for _, each := range path {
		switch each {
		case '+':
			buf.WriteString("%20")
		case ' ':
			buf.WriteString("%20")
		default:
			buf.WriteRune(each)
		}
	}
	return buf.String()
}
