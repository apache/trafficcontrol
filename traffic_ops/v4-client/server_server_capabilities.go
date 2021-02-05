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
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_SERVER_SERVER_CAPABILITIES = apiBase + "/server_server_capabilities"
)

// CreateServerServerCapability assigns a Server Capability to a Server
func (to *Session) CreateServerServerCapability(ssc tc.ServerServerCapability) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(API_SERVER_SERVER_CAPABILITIES, ssc, nil, &alerts)
	return alerts, reqInf, err
}

// DeleteServerServerCapability unassigns a Server Capability from a Server
func (to *Session) DeleteServerServerCapability(serverID int, serverCapability string) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	v := url.Values{}
	v.Add("serverId", strconv.Itoa(serverID))
	v.Add("serverCapability", serverCapability)
	qStr := v.Encode()
	queryURL := fmt.Sprintf("%s?%s", API_SERVER_SERVER_CAPABILITIES, qStr)
	reqInf, err := to.del(queryURL, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) GetServerServerCapabilitiesWithHdr(serverID *int, serverHostName, serverCapability *string, header http.Header) ([]tc.ServerServerCapability, ReqInf, error) {
	v := url.Values{}
	if serverID != nil {
		v.Add("serverId", strconv.Itoa(*serverID))
	}
	if serverHostName != nil {
		v.Add("serverHostName", *serverHostName)
	}
	if serverCapability != nil {
		v.Add("serverCapability", *serverCapability)
	}
	queryURL := API_SERVER_SERVER_CAPABILITIES
	if qStr := v.Encode(); len(qStr) > 0 {
		queryURL = fmt.Sprintf("%s?%s", queryURL, qStr)
	}

	resp := struct {
		Response []tc.ServerServerCapability `json:"response"`
	}{}
	reqInf, err := to.get(queryURL, header, &resp)
	return resp.Response, reqInf, err
}

// GetServerServerCapabilities retrieves a list of Server Capabilities that are assigned to a Server
// Callers can filter the results by server id, server host name and/or server capability via the optional parameters
// Deprecated: GetServerServerCapabilities will be removed in 6.0. Use GetServerServerCapabilitiesWithHdr.
func (to *Session) GetServerServerCapabilities(serverID *int, serverHostName, serverCapability *string) ([]tc.ServerServerCapability, ReqInf, error) {
	return to.GetServerServerCapabilitiesWithHdr(serverID, serverHostName, serverCapability, nil)
}
