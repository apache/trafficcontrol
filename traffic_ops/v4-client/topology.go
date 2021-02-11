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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

// APITopologies is the API version-relative path to the /topologies API endpoint.
const APITopologies = "/topologies"

// CreateTopology creates the passed Topology.
func (to *Session) CreateTopology(top tc.Topology) (tc.TopologyResponse, toclientlib.ReqInf, error) {
	var resp tc.TopologyResponse
	reqInf, err := to.post(APITopologies, top, nil, &resp)
	return resp, reqInf, err
}

// GetTopologies returns all Topologies stored in Traffic Ops.
func (to *Session) GetTopologies(header http.Header) ([]tc.Topology, toclientlib.ReqInf, error) {
	var data tc.TopologiesResponse
	reqInf, err := to.get(APITopologies, header, &data)
	return data.Response, reqInf, err
}

// GetTopology returns the Topology with the given Name.
func (to *Session) GetTopology(name string, header http.Header) (tc.Topology, toclientlib.ReqInf, error) {
	reqURL := fmt.Sprintf("%s?name=%s", APITopologies, url.QueryEscape(name))
	var data tc.TopologiesResponse
	reqInf, err := to.get(reqURL, header, &data)
	if err != nil {
		return tc.Topology{}, reqInf, err
	}
	if len(data.Response) == 1 {
		return data.Response[0], reqInf, nil
	}
	return tc.Topology{}, reqInf, fmt.Errorf("expected one topology in response, instead got %d", len(data.Response))
}

// UpdateTopology updates a Topology by name.
func (to *Session) UpdateTopology(name string, t tc.Topology, header http.Header) (tc.TopologyResponse, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", APITopologies, name)
	var response tc.TopologyResponse
	reqInf, err := to.put(route, t, header, &response)
	return response, reqInf, err
}

// DeleteTopology deletes the Topology with the given name.
func (to *Session) DeleteTopology(name string) (tc.Alerts, toclientlib.ReqInf, error) {
	reqURL := fmt.Sprintf("%s?name=%s", APITopologies, url.QueryEscape(name))
	var alerts tc.Alerts
	reqInf, err := to.del(reqURL, nil, &alerts)
	return alerts, reqInf, err
}
