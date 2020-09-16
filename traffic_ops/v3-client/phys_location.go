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

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_PHYS_LOCATIONS = apiBase + "/phys_locations"
)

// CreatePhysLocation creates a PhysLocation.
func (to *Session) CreatePhysLocation(pl tc.PhysLocation) (tc.Alerts, ReqInf, error) {
	if pl.RegionID == 0 && pl.RegionName != "" {
		regions, _, err := to.GetRegionByName(pl.RegionName)
		if err != nil {
			return tc.Alerts{}, ReqInf{}, err
		}
		if len(regions) == 0 {
			return tc.Alerts{}, ReqInf{}, errors.New("no region with name " + pl.RegionName)
		}
		pl.RegionID = regions[0].ID
	}
	var remoteAddr net.Addr
	reqBody, err := json.Marshal(pl)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_PHYS_LOCATIONS, reqBody, nil)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Update a PhysLocation by ID
func (to *Session) UpdatePhysLocationByID(id int, pl tc.PhysLocation) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(pl)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := fmt.Sprintf("%s/%d", API_PHYS_LOCATIONS, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody, nil)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

func (to *Session) GetPhysLocationsWithHdr(params map[string]string, header http.Header) ([]tc.PhysLocation, ReqInf, error) {
	path := API_PHYS_LOCATIONS + mapToQueryParameters(params)
	resp, remoteAddr, err := to.request(http.MethodGet, path, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.PhysLocation{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.PhysLocationsResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// Returns a list of PhysLocations with optional query parameters applied
// Deprecated: GetPhysLocations will be removed in 6.0. Use GetPhysLocationsWithHdr.
func (to *Session) GetPhysLocations(params map[string]string) ([]tc.PhysLocation, ReqInf, error) {
	return to.GetPhysLocationsWithHdr(params, nil)
}

func (to *Session) GetPhysLocationByIDWithHdr(id int, header http.Header) ([]tc.PhysLocation, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_PHYS_LOCATIONS, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.PhysLocation{}, reqInf, nil
		}
	}
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

// GET a PhysLocation by the PhysLocation ID
// Deprecated: GetPhysLocationByID will be removed in 6.0. Use GetPhysLocationByIDWithHdr.
func (to *Session) GetPhysLocationByID(id int) ([]tc.PhysLocation, ReqInf, error) {
	return to.GetPhysLocationByIDWithHdr(id, nil)
}

func (to *Session) GetPhysLocationByNameWithHdr(name string, header http.Header) ([]tc.PhysLocation, ReqInf, error) {
	url := fmt.Sprintf("%s?name=%s", API_PHYS_LOCATIONS, url.QueryEscape(name))
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.PhysLocation{}, reqInf, nil
		}
	}
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
// Deprecated: GetPhysLocationByName will be removed in 6.0. Use GetPhysLocationByNameWithHdr.
func (to *Session) GetPhysLocationByName(name string) ([]tc.PhysLocation, ReqInf, error) {
	return to.GetPhysLocationByNameWithHdr(name, nil)
}

// DELETE a PhysLocation by ID
func (to *Session) DeletePhysLocationByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_PHYS_LOCATIONS, id)
	resp, remoteAddr, err := to.request(http.MethodDelete, route, nil, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}
