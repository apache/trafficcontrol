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
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

const (
	// apiServers is the API version-relative path to the /servers API
	// endpoint.
	apiServers = "/servers"
)

// CreateServer creates the given Server.
func (to *Session) CreateServer(server tc.ServerV5, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	var remoteAddr net.Addr
	reqInf := toclientlib.ReqInf{CacheHitStatus: toclientlib.CacheHitStatusMiss, RemoteAddr: remoteAddr}

	if server.CacheGroupID <= 0 && server.CacheGroup != "" {
		innerOpts := NewRequestOptions()
		innerOpts.QueryParameters.Set("name", server.CacheGroup)
		cg, reqInf, err := to.GetCacheGroups(innerOpts)
		if err != nil {
			return cg.Alerts, reqInf, fmt.Errorf("no Cache Group named %s: %w", server.CacheGroup, err)
		}
		if len(cg.Response) == 0 {
			return cg.Alerts, reqInf, fmt.Errorf("no Cache Group named %s", server.CacheGroup)
		}
		if cg.Response[0].ID == nil {
			return cg.Alerts, reqInf, fmt.Errorf("Cache Group named %s has a nil ID", server.CacheGroup)
		}
		server.CacheGroupID = *cg.Response[0].ID
	}
	if server.CDNID <= 0 && server.CDN != "" {
		innerOpts := NewRequestOptions()
		innerOpts.QueryParameters.Set("name", server.CDN)
		c, reqInf, err := to.GetCDNs(innerOpts)
		if err != nil {
			return c.Alerts, reqInf, fmt.Errorf("no CDN named %s: %w", server.CDN, err)
		}
		if len(c.Response) == 0 {
			return c.Alerts, reqInf, fmt.Errorf("no CDN named %s", server.CDN)
		}
		server.CDNID = c.Response[0].ID
	}
	if server.PhysicalLocationID <= 0 && server.PhysicalLocation != "" {
		innerOpts := NewRequestOptions()
		innerOpts.QueryParameters.Set("name", server.PhysicalLocation)
		ph, reqInf, err := to.GetPhysLocations(innerOpts)
		if err != nil {
			return ph.Alerts, reqInf, fmt.Errorf("no Physical Location named %s: %w", server.PhysicalLocation, err)
		}
		if len(ph.Response) == 0 {
			return ph.Alerts, reqInf, fmt.Errorf("no Physical Location named %s", server.PhysicalLocation)
		}
		server.PhysicalLocationID = ph.Response[0].ID
	}
	if server.StatusID <= 0 && server.Status != "" {
		innerOpts := NewRequestOptions()
		innerOpts.QueryParameters.Set("name", server.Status)
		st, reqInf, err := to.GetStatuses(innerOpts)
		if err != nil {
			return st.Alerts, reqInf, fmt.Errorf("no Status named %s: %w", server.Status, err)
		}
		if len(st.Response) == 0 {
			return alerts, reqInf, fmt.Errorf("no Status named %s", server.Status)
		}
		server.StatusID = st.Response[0].ID
	}
	if server.TypeID <= 0 && server.Type != "" {
		innerOpts := NewRequestOptions()
		innerOpts.QueryParameters.Set("name", server.Type)
		ty, _, err := to.GetTypes(innerOpts)
		if err != nil {
			return ty.Alerts, reqInf, fmt.Errorf("no Type named '%s': %w", server.Type, err)
		}
		if len(ty.Response) == 0 {
			return ty.Alerts, reqInf, fmt.Errorf("no Type named %s", server.Type)
		}
		server.TypeID = ty.Response[0].ID
	}

	reqInf, err := to.post(apiServers, opts, server, &alerts)
	return alerts, reqInf, err
}

// UpdateServer replaces the Server identified by ID with the provided one.
func (to *Session) UpdateServer(id int, server tc.ServerV5, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	route := fmt.Sprintf("%s/%d", apiServers, id)
	reqInf, err := to.put(route, opts, server, &alerts)
	return alerts, reqInf, err
}

// GetServers retrieves Servers from Traffic Ops.
func (to *Session) GetServers(opts RequestOptions) (tc.ServersV5Response, toclientlib.ReqInf, error) {
	var data tc.ServersV5Response
	reqInf, err := to.get(apiServers, opts, &data)
	return data, reqInf, err
}

// DeleteServer deletes the Server with the given ID.
func (to *Session) DeleteServer(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", apiServers, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, opts, &alerts)
	return alerts, reqInf, err
}

// AssignDeliveryServiceIDsToServerID assigns a set of Delivery Services to a
// single server, optionally replacing any and all existing assignments.
// 'server' should be the requested server's ID, 'dsIDs' should be a slice of
// the requested Delivery Services' IDs. If 'replace' is 'true', existing
// assignments to the server will be replaced.
func (to *Session) AssignDeliveryServiceIDsToServerID(server int, dsIDs []int, replace bool, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	// datatypes here match the library tc.Server's and tc.DeliveryService's ID fields
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("replace", strconv.FormatBool(replace))

	endpoint := fmt.Sprintf(apiServerDeliveryServices, server)

	var alerts tc.Alerts
	reqInf, err := to.post(endpoint, opts, dsIDs, &alerts)
	return alerts, reqInf, err
}

// GetServerIDDeliveryServices returns all of the Delivery Services assigned to the server identified
// by the integral, unique identifier 'server'.
func (to *Session) GetServerIDDeliveryServices(server int, opts RequestOptions) (tc.DeliveryServicesResponseV5, toclientlib.ReqInf, error) {
	endpoint := fmt.Sprintf(apiServerDeliveryServices, server)
	var data tc.DeliveryServicesResponseV5
	reqInf, err := to.get(endpoint, opts, &data)
	return data, reqInf, err
}

// GetServerUpdateStatus retrieves the Server Update Status of the Server with
// the given (short) hostname.
func (to *Session) GetServerUpdateStatus(hostName string, opts RequestOptions) (tc.ServerUpdateStatusResponseV5, toclientlib.ReqInf, error) {
	path := apiServers + `/` + url.PathEscape(hostName) + `/update_status`
	var data tc.ServerUpdateStatusResponseV5
	reqInf, err := to.get(path, opts, &data)
	return data, reqInf, err
}
