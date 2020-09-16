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

package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-log"
)

func get(to *Session, endpoint string, respStruct interface{}, header http.Header) (ReqInf, error) {
	return makeReq(to, "GET", endpoint, nil, respStruct, header)
}

func post(to *Session, endpoint string, body []byte, respStruct interface{}) (ReqInf, error) {
	return makeReq(to, "POST", endpoint, body, respStruct, nil)
}

func put(to *Session, endpoint string, body []byte, respStruct interface{}) (ReqInf, error) {
	return makeReq(to, "PUT", endpoint, body, respStruct, nil)
}

func del(to *Session, endpoint string, respStruct interface{}) (ReqInf, error) {
	return makeReq(to, "DELETE", endpoint, nil, respStruct, nil)
}

func makeReq(to *Session, method, endpoint string, body []byte, respStruct interface{}, header http.Header) (ReqInf, error) {
	resp, remoteAddr, err := to.request(method, endpoint, body, header) // TODO change to getBytesWithTTL
	reqInf := ReqInf{RemoteAddr: remoteAddr, CacheHitStatus: CacheHitStatusMiss}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return reqInf, nil
		}
	}
	if err != nil {
		return reqInf, err
	}
	defer log.Close(resp.Body, "unable to close response body")

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return reqInf, errors.New("reading body: " + err.Error())
	}

	if err := json.Unmarshal(bts, respStruct); err != nil {
		return reqInf, errors.New("unmarshalling body '" + string(body) + "': " + err.Error())
	}

	return reqInf, nil
}

func mapToQueryParameters(params map[string]string) string {
	path := ""
	for key, value := range params {
		if path == "" {
			path += "?"
		} else {
			path += "&"
		}
		path += key + "=" + value
	}
	return path
}
