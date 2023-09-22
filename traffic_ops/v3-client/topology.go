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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// ApiTopologies is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
const ApiTopologies = apiBase + "/topologies"

const APITopologies = "/topologies"

// CreateTopology creates a topology and returns the response.
func (to *Session) CreateTopology(top tc.Topology) (*tc.TopologyResponse, toclientlib.ReqInf, error) {
	resp := new(tc.TopologyResponse)
	reqInf, err := to.post(APITopologies, top, nil, resp)
	return resp, reqInf, err
}

func (to *Session) GetTopologiesWithHdr(header http.Header) ([]tc.Topology, toclientlib.ReqInf, error) {
	var data tc.TopologiesResponse
	reqInf, err := to.get(APITopologies, header, &data)
	return data.Response, reqInf, err
}

// GetTopologies returns all topologies.
// Deprecated: GetTopologies will be removed in 6.0. Use GetTopologiesWithHdr.
func (to *Session) GetTopologies() ([]tc.Topology, toclientlib.ReqInf, error) {
	return to.GetTopologiesWithHdr(nil)
}

func (to *Session) GetTopologyWithHdr(name string, header http.Header) (*tc.Topology, toclientlib.ReqInf, error) {
	reqUrl := fmt.Sprintf("%s?name=%s", APITopologies, url.QueryEscape(name))
	var data tc.TopologiesResponse
	reqInf, err := to.get(reqUrl, header, &data)
	if err != nil {
		return nil, reqInf, err
	}
	if len(data.Response) == 1 {
		return &data.Response[0], reqInf, nil
	}
	return nil, reqInf, fmt.Errorf("expected one topology in response, instead got %d", len(data.Response))
}

// GetTopology returns the given topology by name.
// Deprecated: GetTopology will be removed in 6.0. Use GetTopologyWithHdr.
func (to *Session) GetTopology(name string) (*tc.Topology, toclientlib.ReqInf, error) {
	return to.GetTopologyWithHdr(name, nil)
}

// UpdateTopology updates a Topology by name.
func (to *Session) UpdateTopology(name string, t tc.Topology) (*tc.TopologyResponse, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", APITopologies, name)
	var response = new(tc.TopologyResponse)
	reqInf, err := to.put(route, t, nil, &response)
	return response, reqInf, err
}

// DeleteTopology deletes the given topology by name.
func (to *Session) DeleteTopology(name string) (tc.Alerts, toclientlib.ReqInf, error) {
	reqUrl := fmt.Sprintf("%s?name=%s", APITopologies, url.QueryEscape(name))
	var alerts tc.Alerts
	reqInf, err := to.del(reqUrl, nil, &alerts)
	return alerts, reqInf, err
}
