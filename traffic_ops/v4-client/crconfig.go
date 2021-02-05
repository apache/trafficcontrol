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
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_SNAPSHOT = apiBase + "/snapshot"
)

type OuterResponse struct {
	Response json.RawMessage `json:"response"`
}

// GetCRConfig returns the raw JSON bytes of the CRConfig from Traffic Ops, and whether the bytes were from the client's internal cache.
func (to *Session) GetCRConfig(cdn string) ([]byte, ReqInf, error) {
	uri := apiBase + `/cdns/` + cdn + `/snapshot`
	bts, reqInf, err := to.getBytesWithTTL(uri, tmPollingInterval)
	if err != nil {
		return nil, reqInf, err
	}

	resp := OuterResponse{}
	if err := json.Unmarshal(bts, &resp); err != nil {
		return nil, reqInf, err
	}
	return resp.Response, reqInf, nil
}

func (to *Session) SnapshotCRConfigWithHdr(cdn string, header http.Header) (ReqInf, error) {
	uri := fmt.Sprintf("%s?cdn=%s", API_SNAPSHOT, url.QueryEscape(cdn))
	resp := OuterResponse{}
	reqInf, err := to.put(uri, nil, header, &resp)
	return reqInf, err
}

// GetCRConfigNew returns the raw JSON bytes of the latest CRConfig from Traffic Ops, and whether the bytes were from the client's internal cache.
func (to *Session) GetCRConfigNew(cdn string) ([]byte, ReqInf, error) {
	uri := apiBase + `/cdns/` + cdn + `/snapshot/new`
	bts, reqInf, err := to.getBytesWithTTL(uri, tmPollingInterval)
	if err != nil {
		return nil, reqInf, err
	}

	resp := OuterResponse{}
	if err := json.Unmarshal(bts, &resp); err != nil {
		return nil, reqInf, err
	}
	return resp.Response, reqInf, nil
}

// SnapshotCRConfig snapshots a CDN by name.
// Deprecated: SnapshotCRConfig will be removed in 6.0. Use SnapshotCRConfigWithHdr.
func (to *Session) SnapshotCRConfig(cdn string) (ReqInf, error) {
	return to.SnapshotCRConfigWithHdr(cdn, nil)
}

// SnapshotCDNByID snapshots a CDN by ID.
func (to *Session) SnapshotCRConfigByID(id int) (tc.Alerts, ReqInf, error) {
	url := fmt.Sprintf("%s?cdnID=%d", API_SNAPSHOT, id)
	var alerts tc.Alerts
	reqInf, err := to.put(url, nil, nil, &alerts)
	return alerts, reqInf, err
}
