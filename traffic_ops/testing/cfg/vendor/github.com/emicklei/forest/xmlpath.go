package forest

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"gopkg.in/xmlpath.v2"
)

// XMLPath returns the value found by following the xpath expression in a XML document (payload of response).
func XMLPath(t T, r *http.Response, xpath string) interface{} {
	if r == nil {
		logfatal(t, sfatalf("XMLPath: no response to read body from"))
		return nil
	}
	if r.Body == nil {
		logfatal(t, sfatalf("XMLPath: no response body to read"))
		return nil
	}
	path, err := xmlpath.Compile(xpath)
	if err != nil {
		logerror(t, serrorf("XMLPath: invalid xpath expression:%v", err))
		return nil
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logerror(t, serrorf("XMLPath: unable to read response body"))
		return nil
	}
	root, err := xmlpath.Parse(bytes.NewReader(data))
	// put the body back for re-reads
	r.Body = ioutil.NopCloser(bytes.NewReader(data))

	if err != nil {
		logerror(t, serrorf("XMLPath: unable to parse xml:%v", err))
		return nil
	}
	if value, ok := path.String(root); ok {
		return value
	}
	logerror(t, serrorf("XMLPath: no value for path: %s", xpath))
	return nil
}
