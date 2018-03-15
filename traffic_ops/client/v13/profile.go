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

package v13

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc/common"
)

const (
	API_v13_Profiles = "/api/1.3/profiles"
)

// Create a Profile
func (to *Session) CreateProfile(pl tc.Profile) (common.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(pl)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return common.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_v13_Profiles, reqBody)
	if err != nil {
		return common.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts common.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Update a Profile by ID
func (to *Session) UpdateProfileByID(id int, pl tc.Profile) (common.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(pl)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return common.Alerts{}, reqInf, err
	}
	route := fmt.Sprintf("%s/%d", API_v13_Profiles, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody)
	if err != nil {
		return common.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts common.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Returns a list of Profiles
func (to *Session) GetProfiles() ([]tc.Profile, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_v13_Profiles, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfilesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// GET a Profile by the Profile ID
func (to *Session) GetProfileByID(id int) ([]tc.Profile, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_v13_Profiles, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfilesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GET a Profile by the Profile name
func (to *Session) GetProfileByName(name string) ([]tc.Profile, ReqInf, error) {
	URI := API_v13_Profiles + "?name=" + url.QueryEscape(name)
	resp, remoteAddr, err := to.request(http.MethodGet, URI, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfilesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// DELETE a Profile by ID
func (to *Session) DeleteProfileByID(id int) (common.Alerts, ReqInf, error) {
	URI := fmt.Sprintf("%s/%d", API_v13_Profiles, id)
	resp, remoteAddr, err := to.request(http.MethodDelete, URI, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return common.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts common.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}
