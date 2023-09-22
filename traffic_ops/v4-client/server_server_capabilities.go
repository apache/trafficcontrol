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
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiServerServerCapabilities is the API version-relative path to the
// /server_server_capabilities API endpoint.
const apiServerServerCapabilities = "/server_server_capabilities"

// apiMultipleServersCapabilities is the API version-relative path to the /multiple_servers_capabilities API endpoint.
const apiMultipleServersCapabilities = "/multiple_servers_capabilities"

// CreateServerServerCapability assigns a Server Capability to a Server.
func (to *Session) CreateServerServerCapability(ssc tc.ServerServerCapability, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiServerServerCapabilities, opts, ssc, &alerts)
	return alerts, reqInf, err
}

// DeleteServerServerCapability unassigns a Server Capability from a Server.
func (to *Session) DeleteServerServerCapability(serverID int, serverCapability string, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("serverId", strconv.Itoa(serverID))
	opts.QueryParameters.Set("serverCapability", serverCapability)
	var alerts tc.Alerts
	reqInf, err := to.del(apiServerServerCapabilities, opts, &alerts)
	return alerts, reqInf, err
}

// GetServerServerCapabilities retrieves a list of Server Capabilities that are
// assigned to Servers.
func (to *Session) GetServerServerCapabilities(opts RequestOptions) (tc.ServerServerCapabilitiesResponse, toclientlib.ReqInf, error) {
	var resp tc.ServerServerCapabilitiesResponse
	reqInf, err := to.get(apiServerServerCapabilities, opts, &resp)
	return resp, reqInf, err
}

// AssignMultipleServersCapabilities assigns multiple server capabilities to a server or multiple servers to a capability.
func (to *Session) AssignMultipleServersCapabilities(mssc tc.MultipleServersCapabilities, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiMultipleServersCapabilities, opts, mssc, &alerts)
	return alerts, reqInf, err
}

// DeleteMultipleServersCapabilities deletes multiple assigned server capabilities to a server or multiple servers to a capability.
func (to *Session) DeleteMultipleServersCapabilities(mssc tc.MultipleServersCapabilities, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.req(http.MethodDelete, apiMultipleServersCapabilities, opts, mssc, &alerts)
	return alerts, reqInf, err
}
