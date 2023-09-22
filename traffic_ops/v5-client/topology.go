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
	"net/url"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiTopologies is the API version-relative path to the /topologies API endpoint.
const apiTopologies = "/topologies"

// CreateTopology creates the passed Topology.
func (to *Session) CreateTopology(top tc.TopologyV5, opts RequestOptions) (tc.TopologyResponseV5, toclientlib.ReqInf, error) {
	var resp tc.TopologyResponseV5
	reqInf, err := to.post(apiTopologies, opts, top, &resp)
	return resp, reqInf, err
}

// GetTopologies returns all Topologies stored in Traffic Ops.
func (to *Session) GetTopologies(opts RequestOptions) (tc.TopologiesResponseV5, toclientlib.ReqInf, error) {
	var data tc.TopologiesResponseV5
	reqInf, err := to.get(apiTopologies, opts, &data)
	return data, reqInf, err
}

// UpdateTopology updates a Topology by name.
func (to *Session) UpdateTopology(name string, t tc.TopologyV5, opts RequestOptions) (tc.TopologyResponseV5, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("name", name)
	var response tc.TopologyResponseV5
	reqInf, err := to.put(apiTopologies, opts, t, &response)
	return response, reqInf, err
}

// DeleteTopology deletes the Topology with the given name.
func (to *Session) DeleteTopology(name string, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("name", name)
	var alerts tc.Alerts
	reqInf, err := to.del(apiTopologies, opts, &alerts)
	return alerts, reqInf, err
}
