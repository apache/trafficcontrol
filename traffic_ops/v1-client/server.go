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
	API_v13_Servers                        = "/api/1.3/servers"
	API_v14_Server_Assign_DeliveryServices = "/api/1.4/servers/%d/deliveryservices?replace=%t"
	API_v14_Server_DeliveryServices        = "/api/1.4/servers/%d/deliveryservices"
)

// Create a Server
func (to *Session) CreateServer(server tc.Server) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}

	if server.CachegroupID == 0 && server.Cachegroup != "" {
		cg, _, err := to.GetCacheGroupNullableByName(server.Cachegroup)
		if err != nil {
			return tc.Alerts{}, ReqInf{}, errors.New("no cachegroup named " + server.Cachegroup + ":" + err.Error())
		}
		if len(cg) == 0 {
			return tc.Alerts{}, ReqInf{}, errors.New("no cachegroup named " + server.Cachegroup)
		}
		if cg[0].ID == nil {
			return tc.Alerts{}, ReqInf{}, errors.New("Cachegroup named " + server.Cachegroup + " has a nil ID")
		}
		server.CachegroupID = *cg[0].ID
	}
	if server.CDNID == 0 && server.CDNName != "" {
		c, _, err := to.GetCDNByName(server.CDNName)
		if err != nil {
			return tc.Alerts{}, ReqInf{}, errors.New("no CDN named " + server.CDNName + ":" + err.Error())
		}
		if len(c) == 0 {
			return tc.Alerts{}, ReqInf{}, errors.New("no CDN named " + server.CDNName)
		}
		server.CDNID = c[0].ID
	}
	if server.PhysLocationID == 0 && server.PhysLocation != "" {
		ph, _, err := to.GetPhysLocationByName(server.PhysLocation)
		if err != nil {
			return tc.Alerts{}, ReqInf{}, errors.New("no physlocation named " + server.PhysLocation + ":" + err.Error())
		}
		if len(ph) == 0 {
			return tc.Alerts{}, ReqInf{}, errors.New("no physlocation named " + server.PhysLocation)
		}
		server.PhysLocationID = ph[0].ID
	}
	if server.ProfileID == 0 && server.Profile != "" {
		pr, _, err := to.GetProfileByName(server.Profile)
		if err != nil {
			return tc.Alerts{}, ReqInf{}, errors.New("no profile named " + server.Profile + ":" + err.Error())
		}
		if len(pr) == 0 {
			return tc.Alerts{}, ReqInf{}, errors.New("no profile named " + server.Profile)
		}
		server.ProfileID = pr[0].ID
	}
	if server.StatusID == 0 && server.Status != "" {
		st, _, err := to.GetStatusByName(server.Status)
		if err != nil {
			return tc.Alerts{}, ReqInf{}, errors.New("no status named " + server.Status + ":" + err.Error())
		}
		if len(st) == 0 {
			return tc.Alerts{}, ReqInf{}, errors.New("no status named " + server.Status)
		}
		server.StatusID = st[0].ID
	}
	if server.TypeID == 0 && server.Type != "" {
		ty, _, err := to.GetTypeByName(server.Type)
		if err != nil {
			return tc.Alerts{}, ReqInf{}, errors.New("no type named " + server.Type + ":" + err.Error())
		}
		if len(ty) == 0 {
			return tc.Alerts{}, ReqInf{}, errors.New("no type named " + server.Type)
		}
		server.TypeID = ty[0].ID
	}
	reqBody, err := json.Marshal(server)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}

	resp, remoteAddr, err := to.request(http.MethodPost, API_v13_Servers, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Update a Server by ID
func (to *Session) UpdateServerByID(id int, server tc.Server) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(server)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := fmt.Sprintf("%s/%d", API_v13_Servers, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Returns a list of Servers
func (to *Session) GetServers() ([]tc.Server, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_v13_Servers, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ServersResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// GET a Server by the Server ID
func (to *Session) GetServerByID(id int) ([]tc.Server, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_v13_Servers, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ServersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GET a Server by the Server hostname
func (to *Session) GetServerByHostName(hostName string) ([]tc.Server, ReqInf, error) {
	url := fmt.Sprintf("%s?hostName=%s", API_v13_Servers, hostName)
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ServersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// DELETE a Server by ID
func (to *Session) DeleteServerByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_v13_Servers, id)
	resp, remoteAddr, err := to.request(http.MethodDelete, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

func (to *Session) GetServersByType(qparams url.Values) ([]tc.Server, ReqInf, error) {
	url := fmt.Sprintf("%s.json?%s", API_v13_Servers, qparams.Encode())
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ServersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

func (to *Session) GetServerFQDN(n string) (string, ReqInf, error) {
	// TODO fix to only request one server
	fdn := ""
	servers, reqInf, err := to.GetServers()
	if err != nil {
		return "Error", reqInf, err
	}

	for _, server := range servers {
		if server.HostName == n {
			fdn = fmt.Sprintf("%s.%s", server.HostName, server.DomainName)
		}
	}
	if fdn == "" {
		return "Error", reqInf, fmt.Errorf("No Server %s found", n)
	}
	return fdn, reqInf, nil
}

func (to *Session) GetServersShortNameSearch(shortname string) ([]string, ReqInf, error) {
	var serverlst []string
	servers, reqInf, err := to.GetServers()
	if err != nil {
		serverlst = append(serverlst, "N/A")
		return serverlst, reqInf, err
	}
	for _, server := range servers {
		if strings.Contains(server.HostName, shortname) {
			serverlst = append(serverlst, server.HostName)
		}
	}
	if len(serverlst) == 0 {
		serverlst = append(serverlst, "N/A")
		return serverlst, reqInf, fmt.Errorf("No Servers Found")
	}
	return serverlst, reqInf, nil
}

// AssignDeliveryServiceIDsToServerID assigns a set of Delivery Services to a single server, optionally
// replacing any and all existing assignments. 'server' should be the requested server's ID, 'dsIDs'
// should be a slice of the requested Delivery Services' IDs. If 'replace' is 'true', existing
// assignments to the server will be replaced.
func (to *Session) AssignDeliveryServiceIDsToServerID(server int, dsIDs []int, replace bool) (tc.Alerts, ReqInf, error) {
	// datatypes here match the library tc.Server's and tc.DeliveryService's ID fields

	var remoteAddr net.Addr
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}

	endpoint := fmt.Sprintf(API_v14_Server_Assign_DeliveryServices, server, replace)

	reqBody, err := json.Marshal(dsIDs)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}

	resp, remoteAddr, err := to.request(http.MethodPost, endpoint, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	reqInf.RemoteAddr = remoteAddr
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, err
}

func (to *Session) GetServerIDDeliveryServices(server int) ([]tc.DeliveryServiceNullable, ReqInf, error) {
	endpoint := fmt.Sprintf(API_v14_Server_DeliveryServices, server)

	resp, remoteAddr, err := to.request(http.MethodGet, endpoint, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.DeliveryServicesNullableResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, err
}
