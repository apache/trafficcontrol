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
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	// APIServers is the API version-relative path to the /servers API
	// endpoint.
	APIServers = "/servers"
	// APIServersDetails is the API version-relative path to the
	// /servers/details API endpoint.
	APIServersDetails = "/servers/details"
	// APIServerAssignDeliveryServices is the API version-relative path to the
	// /servers/{{ID}}/deliveryservices API endpoint with the 'replace' query
	// string parameter.
	APIServerAssignDeliveryServices = APIServerDeliveryServices + "?replace=%t"
)

func needAndCanFetch(id *int, name *string) bool {
	return (id == nil || *id == 0) && name != nil && *name != ""
}

// CreateServer creates the given Server.
func (to *Session) CreateServer(server tc.ServerV40, hdr http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	var remoteAddr net.Addr
	reqInf := toclientlib.ReqInf{CacheHitStatus: toclientlib.CacheHitStatusMiss, RemoteAddr: remoteAddr}

	if needAndCanFetch(server.CachegroupID, server.Cachegroup) {
		cg, _, err := to.GetCacheGroupByName(*server.Cachegroup, nil)
		if err != nil {
			return alerts, reqInf, fmt.Errorf("no cachegroup named %s: %v", *server.Cachegroup, err)
		}
		if len(cg) == 0 {
			return alerts, reqInf, fmt.Errorf("no cachegroup named %s", *server.Cachegroup)
		}
		if cg[0].ID == nil {
			return alerts, reqInf, fmt.Errorf("Cachegroup named %s has a nil ID", *server.Cachegroup)
		}
		server.CachegroupID = cg[0].ID
	}
	if needAndCanFetch(server.CDNID, server.CDNName) {
		c, _, err := to.GetCDNByName(*server.CDNName, nil)
		if err != nil {
			return alerts, reqInf, fmt.Errorf("no CDN named %s: %v", *server.CDNName, err)
		}
		if len(c) == 0 {
			return alerts, reqInf, fmt.Errorf("no CDN named %s", *server.CDNName)
		}
		server.CDNID = &c[0].ID
	}
	if needAndCanFetch(server.PhysLocationID, server.PhysLocation) {
		ph, _, err := to.GetPhysLocationByName(*server.PhysLocation, nil)
		if err != nil {
			return alerts, reqInf, fmt.Errorf("no physlocation named %s: %v", *server.PhysLocation, err)
		}
		if len(ph) == 0 {
			return alerts, reqInf, fmt.Errorf("no physlocation named %s", *server.PhysLocation)
		}
		server.PhysLocationID = &ph[0].ID
	}
	if needAndCanFetch(server.ProfileID, server.Profile) {
		pr, _, err := to.GetProfileByName(*server.Profile, nil)
		if err != nil {
			return alerts, reqInf, fmt.Errorf("no profile named %s: %v", *server.Profile, err)
		}
		if len(pr) == 0 {
			return alerts, reqInf, fmt.Errorf("no profile named %s", *server.Profile)
		}
		server.ProfileID = &pr[0].ID
	}
	if needAndCanFetch(server.StatusID, server.Status) {
		st, _, err := to.GetStatusByName(*server.Status, nil)
		if err != nil {
			return alerts, reqInf, fmt.Errorf("no status named %s: %v", *server.Status, err)
		}
		if len(st) == 0 {
			return alerts, reqInf, fmt.Errorf("no status named %s", *server.Status)
		}
		server.StatusID = &st[0].ID
	}
	if (server.TypeID == nil || *server.TypeID == 0) && server.Type != "" {
		ty, _, err := to.GetTypeByName(server.Type, nil)
		if err != nil {
			return alerts, reqInf, fmt.Errorf("no type named %s: %v", server.Type, err)
		}
		if len(ty) == 0 {
			return alerts, reqInf, fmt.Errorf("no type named %s", server.Type)
		}
		server.TypeID = &ty[0].ID
	}

	reqInf, err := to.post(APIServers, server, hdr, &alerts)
	return alerts, reqInf, err
}

// UpdateServer replaces the Server identified by ID with the provided one.
func (to *Session) UpdateServer(id int, server tc.ServerV40, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	route := fmt.Sprintf("%s/%d", APIServers, id)
	reqInf, err := to.put(route, server, header, &alerts)
	return alerts, reqInf, err
}

// GetServers retrieves Servers from Traffic Ops.
func (to *Session) GetServers(params url.Values, header http.Header) (tc.ServersV4Response, toclientlib.ReqInf, error) {
	route := APIServers
	if params != nil {
		route += "?" + params.Encode()
	}

	var data tc.ServersV4Response
	reqInf, err := to.get(route, header, &data)
	return data, reqInf, err
}

// GetServerDetailsByHostName retrieves the Server Details of the Server with
// the given (short) Hostname.
func (to *Session) GetServerDetailsByHostName(hostName string, header http.Header) ([]tc.ServerDetailV40, toclientlib.ReqInf, error) {
	v := url.Values{}
	v.Add("hostName", hostName)
	route := APIServersDetails + "?" + v.Encode()
	var data tc.ServersV4DetailResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// DeleteServer deletes the Server with the given ID.
func (to *Session) DeleteServer(id int, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APIServers, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

// GetServerFQDN returns the Fully Qualified Domain Name (FQDN) of the first
// server found to have the Host Name 'n'.
func (to *Session) GetServerFQDN(n string, header http.Header) (string, tc.Alerts, toclientlib.ReqInf, error) {
	// TODO fix to only request one server
	params := url.Values{}
	params.Add("hostName", n)

	resp, reqInf, err := to.GetServers(params, header)
	if err != nil {
		return "", resp.Alerts, reqInf, err
	}

	var fdn string
	for _, server := range resp.Response {
		if server.HostName != nil && server.DomainName != nil {
			fdn = fmt.Sprintf("%s.%s", *server.HostName, *server.DomainName)
		}
	}

	if fdn == "" {
		err = fmt.Errorf("No Server %s found", n)
	}

	return fdn, resp.Alerts, reqInf, err
}

// AssignDeliveryServiceIDsToServerID assigns a set of Delivery Services to a
// single server, optionally replacing any and all existing assignments.
// 'server' should be the requested server's ID, 'dsIDs' should be a slice of
// the requested Delivery Services' IDs. If 'replace' is 'true', existing
// assignments to the server will be replaced.
func (to *Session) AssignDeliveryServiceIDsToServerID(server int, dsIDs []int, replace bool, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	// datatypes here match the library tc.Server's and tc.DeliveryService's ID fields
	endpoint := fmt.Sprintf(APIServerAssignDeliveryServices, server, replace)
	var alerts tc.Alerts
	reqInf, err := to.post(endpoint, dsIDs, nil, &alerts)
	return alerts, reqInf, err
}

// GetServerIDDeliveryServices returns all of the Delivery Services assigned to the server identified
// by the integral, unique identifier 'server'.
func (to *Session) GetServerIDDeliveryServices(server int, header http.Header) ([]tc.DeliveryServiceNullable, toclientlib.ReqInf, error) {
	endpoint := fmt.Sprintf(APIServerDeliveryServices, server)
	var data tc.DeliveryServicesNullableResponse
	reqInf, err := to.get(endpoint, header, &data)
	return data.Response, reqInf, err
}

// GetServerUpdateStatus retrieves the Server Update Status of the Server with
// the given (short) hostname.
func (to *Session) GetServerUpdateStatus(hostName string, header http.Header) (tc.ServerUpdateStatus, toclientlib.ReqInf, error) {
	path := APIServers + `/` + hostName + `/update_status`
	data := []tc.ServerUpdateStatus{}
	reqInf, err := to.get(path, header, &data)
	if err != nil {
		return tc.ServerUpdateStatus{}, reqInf, err
	}
	if len(data) == 0 {
		return tc.ServerUpdateStatus{}, reqInf, errors.New("Traffic Ops returned no update statuses for that server")
	}
	return data[0], reqInf, nil
}
