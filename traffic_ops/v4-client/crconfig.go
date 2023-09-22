package client

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

import (
	"encoding/json"
	"errors"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiSnapshot is the API version-relative path for the /snapshot API endpoint.
const apiSnapshot = "/snapshot"

// OuterResponse is the most basic type of a Traffic Ops API response,
// containing some kind of JSON-encoded 'response' object.
type OuterResponse struct {
	Response json.RawMessage `json:"response"`
}

// GetCRConfig returns the Snapshot for the given CDN from Traffic Ops.
func (to *Session) GetCRConfig(cdn string, opts RequestOptions) (tc.SnapshotResponse, toclientlib.ReqInf, error) {
	uri := `/cdns/` + cdn + `/snapshot`
	var resp tc.SnapshotResponse
	reqInf, err := to.get(uri, opts, &resp)
	return resp, reqInf, err
}

// SnapshotCRConfig creates a new Snapshot for the CDN with the given Name -
// NOT just a new CRConfig!
func (to *Session) SnapshotCRConfig(opts RequestOptions) (tc.PutSnapshotResponse, toclientlib.ReqInf, error) {
	var resp tc.PutSnapshotResponse
	if opts.QueryParameters == nil || (opts.QueryParameters.Get("cdn") == "" && opts.QueryParameters.Get("cdnID") == "") {
		return resp, toclientlib.ReqInf{}, errors.New("cannot take Snapshot of unidentified CDN - set 'cdn' or 'cdnID' query parameter")
	}
	reqInf, err := to.put(apiSnapshot, opts, nil, &resp)
	return resp, reqInf, err
}

// GetCRConfigNew returns the *new* Snapshot for the given CDN from Traffic
// Ops.
func (to *Session) GetCRConfigNew(cdn string, opts RequestOptions) (tc.SnapshotResponse, toclientlib.ReqInf, error) {
	uri := `/cdns/` + cdn + `/snapshot/new`
	var resp tc.SnapshotResponse
	reqInf, err := to.get(uri, opts, &resp)
	return resp, reqInf, err
}
