package forest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

var verboseOnFailure = false

// VerboseOnFailure (default is false) will produce more information about the request and response when a failure is detected on an expectation.
// This setting is not the same but related to the value of testing.Verbose().
func VerboseOnFailure(verbose bool) {
	verboseOnFailure = verbose
}

// ExpectStatus inspects the response status code.
// If the value is not expected, the complete request, response is logged (iff verboseOnFailure) and the test is aborted.
// Return true if the status is as expected.
func ExpectStatus(t T, r *http.Response, status int) bool {
	if r == nil {
		logerror(t, serrorf("ExpectStatus: got nil but want Http response"))
		return false
	}
	if r.StatusCode != status {
		if verboseOnFailure {
			Dump(t, r)
		}
		logfatal(t, serrorf("ExpectStatus: got status %d but want %d, %s %v", r.StatusCode, status, r.Request.Method, r.Request.URL))
		return false
	}
	return true
}

// CheckError simply tests the error and fail is not undefined.
// This is implicity called after sending a Http request.
// Return true if there was an error.
func CheckError(t T, err error) bool {
	if err != nil {
		logerror(t, serrorf("CheckError: did not expect to receive err: %v", err))
	}
	return err != nil
}

// ExpectHeader inspects the header of the response.
// Return true if the header matches.
func ExpectHeader(t T, r *http.Response, name, value string) bool {
	if r == nil {
		logerror(t, serrorf("ExpectHeader: got nil but want a Http response"))
		return false
	}
	rname := r.Header.Get(name)
	if rname != value {
		logerror(t, serrorf("ExpectHeader: got header %s=%s but want %s", name, rname, value))
	}
	return rname == value
}

// ExpectJSONHash tries to unmarshal the response body into a Go map callback parameter.
// Fail if the body could not be read or if unmarshalling was not possible.
// Returns true if the callback was executed with a map.
func ExpectJSONHash(t T, r *http.Response, callback func(hash map[string]interface{})) bool {
	if r == nil {
		logerror(t, serrorf("ExpectJSONHash: no response available"))
		return false
	}
	if r.Body == nil {
		logerror(t, serrorf("ExpectJSONHash: no body to read"))
		return false
	}
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		if verboseOnFailure {
			Dump(t, r)
		}
		logerror(t, serrorf("ExpectJSONHash: unable to read response body:%v", err))
		return false
	}
	// put the body back for re-reads
	r.Body = ioutil.NopCloser(bytes.NewReader(data))

	dict := map[string]interface{}{}
	err = json.Unmarshal(data, &dict)
	if err != nil {
		if verboseOnFailure {
			Dump(t, r)
		}
		logerror(t, serrorf("ExpectJSONHash: unable to unmarshal Json:%v", err))
		return false
	}
	callback(dict)
	return true
}

// ExpectJSONArray tries to unmarshal the response body into a Go slice callback parameter.
// Fail if the body could not be read or if unmarshalling was not possible.
// Returns true if the callback was executed with an array.
func ExpectJSONArray(t T, r *http.Response, callback func(array []interface{})) bool {
	if r == nil {
		logerror(t, serrorf("ExpectJSONArray: no response available"))
		return false
	}
	if r.Body == nil {
		logerror(t, serrorf("ExpectJSONArray: no body to read"))
		return false
	}
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		if verboseOnFailure {
			Dump(t, r)
		}
		logerror(t, serrorf("ExpectJSONArray: unable to read response body:%v", err))
		return false
	}
	// put the body back for re-reads
	r.Body = ioutil.NopCloser(bytes.NewReader(data))

	slice := []interface{}{}
	err = json.Unmarshal(data, &slice)
	if err != nil {
		if verboseOnFailure {
			Dump(t, r)
		}
		logerror(t, serrorf("ExpectJSONArray: unable to unmarshal Json:%v", err))
		return false
	}
	callback(slice)
	return true
}

// ExpectString reads the response body into a Go string callback parameter.
// Fail if the body could not be read or unmarshalled.
// Returns true if a response body was read.
func ExpectString(t T, r *http.Response, callback func(content string)) bool {
	if r == nil {
		logerror(t, serrorf("ExpectString: no response available"))
		return false
	}
	if r.Body == nil {
		logerror(t, serrorf("ExpectString: no body to read"))
		return false
	}
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		if verboseOnFailure {
			Dump(t, r)
		}
		logerror(t, serrorf("ExpectString: unable to read response body:%v", err))
		return false
	}
	// put the body back for re-reads
	r.Body = ioutil.NopCloser(bytes.NewReader(data))

	callback(string(data))
	return true
}

// ExpectXMLDocument tries to unmarshal the response body into fields of the provided document (struct).
// Fail if the body could not be read or unmarshalled.
// Returns true if a document could be unmarshalled.
func ExpectXMLDocument(t T, r *http.Response, doc interface{}) bool {
	if r == nil {
		logerror(t, serrorf("ExpectXMLDocument: no response available"))
		return false
	}
	if r.Body == nil {
		logerror(t, serrorf("ExpectXMLDocument: no body to read"))
		return false
	}
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		if verboseOnFailure {
			Dump(t, r)
		}
		logerror(t, serrorf("ExpectXMLDocument: unable to read response body:%v", err))
		return false
	}
	// put the body back for re-reads
	r.Body = ioutil.NopCloser(bytes.NewReader(data))

	err = xml.Unmarshal(data, doc)
	if err != nil {
		if verboseOnFailure {
			Dump(t, r)
		}
		logerror(t, serrorf("ExpectXMLDocument: unable to unmarshal Xml:%v", err))
	}
	return err == nil
}

// ExpectJSONDocument tries to unmarshal the response body into fields of the provided document (struct).
// Fail if the body could not be read or unmarshalled.
// Returns true if a document could be unmarshalled.
func ExpectJSONDocument(t T, r *http.Response, doc interface{}) bool {
	if r == nil {
		logerror(t, serrorf("ExpectJSONDocument: no response available"))
		return false
	}
	if r.Body == nil {
		logerror(t, serrorf("ExpectJSONDocument: no body to read"))
		return false
	}
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logerror(t, serrorf("ExpectJSONDocument: unable to read response body :%v", err))
		return false
	}
	// put the body back for re-reads
	r.Body = ioutil.NopCloser(bytes.NewReader(data))

	err = json.Unmarshal(data, doc)
	if err != nil {
		if verboseOnFailure {
			Dump(t, r)
		}
		logerror(t, serrorf("ExpectJSONDocument: unable to unmarshal Json:%v", err))
	}
	return err == nil
}
