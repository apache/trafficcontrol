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
	API_v13_REGIONS = "/api/1.3/regions"
)

// Create a Region
func (to *Session) CreateRegion(region tc.Region) (tc.Alerts, ReqInf, error) {
	if region.Division == 0 && region.DivisionName != "" {
		divisions, _, err := to.GetDivisionByName(region.DivisionName)
		if err != nil {
			return tc.Alerts{}, ReqInf{}, err
		}
		if len(divisions) == 0 {
			return tc.Alerts{}, ReqInf{}, errors.New("no division with name " + region.DivisionName)
		}
		region.Division = divisions[0].ID
	}

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(region)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_v13_REGIONS, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Update a Region by ID
func (to *Session) UpdateRegionByID(id int, region tc.Region) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(region)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := fmt.Sprintf("%s/%d", API_v13_REGIONS, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Returns a list of regions
func (to *Session) GetRegions() ([]tc.Region, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_v13_REGIONS, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.RegionsResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// GET a Region by the Region id
func (to *Session) GetRegionByID(id int) ([]tc.Region, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_v13_REGIONS, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.RegionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GET a Region by the Region name
func (to *Session) GetRegionByName(name string) ([]tc.Region, ReqInf, error) {
	url := fmt.Sprintf("%s?name=%s", API_v13_REGIONS, url.QueryEscape(name))
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.RegionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// DELETE a Region by ID
func (to *Session) DeleteRegionByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_v13_REGIONS, id)
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

// GetRegionByNamePath gets a region by name, using the /api/version/region/name path. This gets the same data as GetRegionByName, but uses a different API path to get the same data, and returns a slightly different format.
func (to *Session) GetRegionByNamePath(name string) ([]tc.RegionName, ReqInf, error) {
	url := apiBase + `/regions/name/` + url.QueryEscape(name)
	reqResp, remoteAddr, err := to.request(http.MethodGet, url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer reqResp.Body.Close()

	resp := tc.RegionNameResponse{}
	if err := json.NewDecoder(reqResp.Body).Decode(&resp); err != nil {
		return nil, reqInf, err
	}
	return resp.Response, reqInf, nil
}
