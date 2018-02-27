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
	"net"
	"net/http"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

const (
	API_v13_PHYS_LOCATIONS = "/api/1.3/phys_locations"
)

// Create a PhysLocation
func (to *Session) CreatePhysLocation(physLocation tc.PhysLocation) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(physLocation)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_v13_PHYS_LOCATIONS, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Update a PhysLocation by ID
func (to *Session) UpdatePhysLocationByID(id int, physLocation tc.PhysLocation) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(physLocation)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := fmt.Sprintf("%s/%d", API_v13_PHYS_LOCATIONS, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Returns a list of physLocations
func (to *Session) GetPhysLocations() ([]tc.PhysLocation, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_v13_PHYS_LOCATIONS, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.PhysLocationsResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// GET a PhysLocation by the PhysLocation id
func (to *Session) GetPhysLocationByID(id int) ([]tc.PhysLocation, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_v13_PHYS_LOCATIONS, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.PhysLocationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GET a PhysLocation by the PhysLocation name
func (to *Session) GetPhysLocationByName(name string) ([]tc.PhysLocation, ReqInf, error) {
	url := fmt.Sprintf("%s?name=%s", API_v13_PHYS_LOCATIONS, name)
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.PhysLocationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// DELETE a PhysLocation by ID
func (to *Session) DeletePhysLocationByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_v13_PHYS_LOCATIONS, id)
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
