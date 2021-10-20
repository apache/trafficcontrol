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
	"github.com/apache/trafficcontrol/v6/lib/go-tc"
	"net/http"
	"net/url"
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
	return []byte(resp.Response), reqInf, nil
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
	return []byte(resp.Response), reqInf, nil
}

// SnapshotCRConfig snapshots a CDN by name.
func (to *Session) SnapshotCRConfig(cdn string) (ReqInf, error) {
	uri := fmt.Sprintf("%s?cdn=%s", API_SNAPSHOT, url.QueryEscape(cdn))
	_, remoteAddr, err := to.request(http.MethodPut, uri, nil)
	reqInf := ReqInf{RemoteAddr: remoteAddr, CacheHitStatus: CacheHitStatusMiss}
	return reqInf, err
}

// SnapshotCDNByID snapshots a CDN by ID.
func (to *Session) SnapshotCRConfigByID(id int) (tc.Alerts, ReqInf, error) {
	url := fmt.Sprintf("%s?cdnID=%d", API_SNAPSHOT, id)
	resp, remoteAddr, err := to.request(http.MethodPut, url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}
