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
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_SERVERS                         = apiBase + "/servers"
	API_SERVERS_DETAILS                 = apiBase + "/servers/details"
	API_SERVER_ASSIGN_DELIVERY_SERVICES = API_SERVER_DELIVERY_SERVICES + "?replace=%t"
)

func needAndCanFetch(id *int, name *string) bool {
	return (id == nil || *id == 0) && name != nil && *name != ""
}

// CreateServer creates a Server.
func (to *Session) CreateServer(server tc.ServerV40, hdr http.Header) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	var remoteAddr net.Addr
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}

	if needAndCanFetch(server.CachegroupID, server.Cachegroup) {
		cg, _, err := to.GetCacheGroupNullableByNameWithHdr(*server.Cachegroup, nil)
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
		c, _, err := to.GetCDNByNameWithHdr(*server.CDNName, nil)
		if err != nil {
			return alerts, reqInf, fmt.Errorf("no CDN named %s: %v", *server.CDNName, err)
		}
		if len(c) == 0 {
			return alerts, reqInf, fmt.Errorf("no CDN named %s", *server.CDNName)
		}
		server.CDNID = &c[0].ID
	}
	if needAndCanFetch(server.PhysLocationID, server.PhysLocation) {
		ph, _, err := to.GetPhysLocationByNameWithHdr(*server.PhysLocation, nil)
		if err != nil {
			return alerts, reqInf, fmt.Errorf("no physlocation named %s: %v", *server.PhysLocation, err)
		}
		if len(ph) == 0 {
			return alerts, reqInf, fmt.Errorf("no physlocation named %s", *server.PhysLocation)
		}
		server.PhysLocationID = &ph[0].ID
	}
	if needAndCanFetch(server.ProfileID, server.Profile) {
		pr, _, err := to.GetProfileByNameWithHdr(*server.Profile, nil)
		if err != nil {
			return alerts, reqInf, fmt.Errorf("no profile named %s: %v", *server.Profile, err)
		}
		if len(pr) == 0 {
			return alerts, reqInf, fmt.Errorf("no profile named %s", *server.Profile)
		}
		server.ProfileID = &pr[0].ID
	}
	if needAndCanFetch(server.StatusID, server.Status) {
		st, _, err := to.GetStatusByNameWithHdr(*server.Status, nil)
		if err != nil {
			return alerts, reqInf, fmt.Errorf("no status named %s: %v", *server.Status, err)
		}
		if len(st) == 0 {
			return alerts, reqInf, fmt.Errorf("no status named %s", *server.Status)
		}
		server.StatusID = &st[0].ID
	}
	if (server.TypeID == nil || *server.TypeID == 0) && server.Type != "" {
		ty, _, err := to.GetTypeByNameWithHdr(server.Type, nil)
		if err != nil {
			return alerts, reqInf, fmt.Errorf("no type named %s: %v", server.Type, err)
		}
		if len(ty) == 0 {
			return alerts, reqInf, fmt.Errorf("no type named %s", server.Type)
		}
		server.TypeID = &ty[0].ID
	}

	reqInf, err := to.post(API_SERVERS, server, hdr, &alerts)
	return alerts, reqInf, err
}

// UpdateServerByID updates a Server by ID.
func (to *Session) UpdateServerByID(id int, server tc.ServerV40, header http.Header) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	var remoteAddr net.Addr
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}

	reqBody, err := json.Marshal(server)
	if err != nil {
		return alerts, reqInf, err
	}

	route := fmt.Sprintf("%s/%d", API_SERVERS, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody, header)
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
	}
	reqInf.RemoteAddr = remoteAddr
	reqInf.StatusCode = resp.StatusCode
	if err != nil {
		return alerts, reqInf, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, err
}

// GetServersWithHdr retrieves a list of servers using the given optional query
// string parameters and HTTP headers.
func (to *Session) GetServersWithHdr(params *url.Values, header http.Header) (tc.ServersV4Response, ReqInf, error) {
	route := API_SERVERS
	if params != nil {
		route += "?" + params.Encode()
	}

	var data tc.ServersV4Response
	reqInf, err := to.get(route, header, &data)
	return data, reqInf, err
}

// GetServers retrieves a list of servers using the given optional query
// string parameters and HTTP headers.
func (to *Session) GetServers(params *url.Values, header http.Header) ([]tc.ServerV40, ReqInf, error) {
	srvs, inf, err := to.GetServersWithHdr(params, nil)
	if err != nil {
		return []tc.ServerV40{}, inf, err
	}

	servers := make([]tc.ServerV40, 0, len(srvs.Response))
	for _, srv := range srvs.Response {
		servers = append(servers, srv)
	}
	return servers, inf, nil
}

// GetFirstServer returns the first server in a servers GET response.
// If no servers match, an error is returned.
// The 'params' parameter can be used to optionally pass URL "query string
// parameters" in the request.
// It returns, in order, the API response that Traffic Ops returned, a request
// info object, and any error that occurred.
func (to *Session) GetFirstServer(params *url.Values, header http.Header) (tc.ServerV40, ReqInf, error) {
	serversResponse, reqInf, err := to.GetServersWithHdr(params, header)
	var firstServer tc.ServerV40
	if err != nil || reqInf.StatusCode == http.StatusNotModified {
		return firstServer, reqInf, err
	}
	for _, firstServer = range serversResponse.Response {
		return firstServer, reqInf, err
	}

	err = fmt.Errorf("unable to find server matching params %v", *params)
	return firstServer, reqInf, err
}

// GetServerDetailsByHostName GETs Servers by the Server hostname.
func (to *Session) GetServerDetailsByHostName(hostName string, header http.Header) ([]tc.ServerDetailV40, ReqInf, error) {
	v := url.Values{}
	v.Add("hostName", hostName)
	route := API_SERVERS_DETAILS + "?" + v.Encode()
	var data tc.ServersV4DetailResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// DeleteServerByID DELETEs a Server by ID.
func (to *Session) DeleteServerByID(id int, header http.Header) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_SERVERS, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

// GetServerFQDN returns the Fully Qualified Domain Name (FQDN) of the first
// server found to have the Host Name 'n'.
func (to *Session) GetServerFQDN(n string, header http.Header) (string, tc.Alerts, ReqInf, error) {
	// TODO fix to only request one server
	params := url.Values{}
	params.Add("hostName", n)

	resp, reqInf, err := to.GetServersWithHdr(&params, header)
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

// GetServersShortNameSearch returns all of the Host Names of servers that
// contain 'shortname'.
func (to *Session) GetServersShortNameSearch(shortname string, header http.Header) ([]string, tc.Alerts, ReqInf, error) {
	var serverlst []string
	resp, reqInf, err := to.GetServersWithHdr(nil, header)
	if err != nil {
		return serverlst, resp.Alerts, reqInf, err
	}

	for _, server := range resp.Response {
		if server.HostName != nil && strings.Contains(*server.HostName, shortname) {
			serverlst = append(serverlst, *server.HostName)
		}
	}

	if len(serverlst) == 0 {
		err = errors.New("No Servers Found")
	}

	return serverlst, resp.Alerts, reqInf, err
}

// AssignDeliveryServiceIDsToServerID assigns a set of Delivery Services to a
// single server, optionally replacing any and all existing assignments.
// 'server' should be the requested server's ID, 'dsIDs' should be a slice of
// the requested Delivery Services' IDs. If 'replace' is 'true', existing
// assignments to the server will be replaced.
func (to *Session) AssignDeliveryServiceIDsToServerID(server int, dsIDs []int, replace bool, header http.Header) (tc.Alerts, ReqInf, error) {
	// datatypes here match the library tc.Server's and tc.DeliveryService's ID fields
	endpoint := fmt.Sprintf(API_SERVER_ASSIGN_DELIVERY_SERVICES, server, replace)
	var alerts tc.Alerts
	reqInf, err := to.post(endpoint, dsIDs, nil, &alerts)
	return alerts, reqInf, err
}

// GetServerIDDeliveryServices returns all of the Delivery Services assigned to the server identified
// by the integral, unique identifier 'server'.
func (to *Session) GetServerIDDeliveryServices(server int, header http.Header) ([]tc.DeliveryServiceNullable, ReqInf, error) {
	endpoint := fmt.Sprintf(API_SERVER_DELIVERY_SERVICES, server)
	var data tc.DeliveryServicesNullableResponse
	reqInf, err := to.get(endpoint, header, &data)
	return data.Response, reqInf, err
}

// GetServerUpdateStatus GETs the Server Update Status by the Server hostname.
func (to *Session) GetServerUpdateStatus(hostName string, header http.Header) (tc.ServerUpdateStatus, ReqInf, error) {
	path := API_SERVERS + `/` + hostName + `/update_status`
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
