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
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	v13Profiles = "/api/1.3/profiles"
	v14Profiles = "/api/1.4/profiles"
)

// CreateProfile creates a Profile
func (to *Session) CreateProfile(pl tc.Profile) (tc.Alerts, ReqInf, error) {
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss}

	if pl.CDNID == 0 && pl.CDNName != "" {
		cdns, _, err := to.GetCDNByName(pl.CDNName)
		if err != nil {
			return tc.Alerts{}, ReqInf{}, err
		}
		if len(cdns) == 0 {
			return tc.Alerts{
					Alerts: []tc.Alert{
						tc.Alert{
							Text:  fmt.Sprintf("no CDN with name %s", pl.CDNName),
							Level: "error",
						},
					},
				},
				ReqInf{},
				fmt.Errorf("no CDN with name %s", pl.CDNName)
		}
		pl.CDNID = cdns[0].ID
	}

	reqBody, err := json.Marshal(pl)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}

	resp, remoteAddr, err := to.request(http.MethodPost, v13Profiles, reqBody)
	reqInf.RemoteAddr = remoteAddr
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()

	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)

	return alerts, reqInf, err
}

// UpdateProfileByID updates a Profile by ID
func (to *Session) UpdateProfileByID(id int, pl tc.Profile) (tc.Alerts, ReqInf, error) {
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss}

	reqBody, err := json.Marshal(pl)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}

	route := fmt.Sprintf("%s/%d", v13Profiles, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody)
	reqInf.RemoteAddr = remoteAddr
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()

	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)

	return alerts, reqInf, err
}

// Profiles gets an array of Profiles
// Deprecated: use GetProfiles
func (to *Session) Profiles() ([]tc.Profile, error) {
	ps, _, err := to.GetProfiles()
	return ps, err
}

// GetProfiles returns a list of Profiles
func (to *Session) GetProfiles() ([]tc.Profile, ReqInf, error) {
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss}

	resp, remoteAddr, err := to.request(http.MethodGet, v13Profiles, nil)
	reqInf.RemoteAddr = remoteAddr
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfilesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)

	return data.Response, reqInf, err
}

// GetProfileByID gET a Profile by the Profile ID
func (to *Session) GetProfileByID(id int) ([]tc.Profile, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", v13Profiles, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfilesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)

	return data.Response, reqInf, err
}

// GetProfileByName GET a Profile by the Profile name
func (to *Session) GetProfileByName(name string) ([]tc.Profile, ReqInf, error) {
	URI := fmt.Sprintf("%s?name=%s", v13Profiles, url.QueryEscape(name))
	resp, remoteAddr, err := to.request(http.MethodGet, URI, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfilesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)

	return data.Response, reqInf, err
}

// GetProfileByParameter GETs a Profile by the Profile "param"
func (to *Session) GetProfileByParameter(param string) ([]tc.Profile, ReqInf, error) {
	URI := fmt.Sprintf("%s?param=%s", v13Profiles, url.QueryEscape(param))
	resp, remoteAddr, err := to.request(http.MethodGet, URI, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfilesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)

	return data.Response, reqInf, err
}

// GetProfileByCDNID GETs a Profile by the Profile cdn id
func (to *Session) GetProfileByCDNID(cdnID int) ([]tc.Profile, ReqInf, error) {
	URI := fmt.Sprintf("%s?cdn=%s", v13Profiles, strconv.Itoa(cdnID))
	resp, remoteAddr, err := to.request(http.MethodGet, URI, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfilesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)

	return data.Response, reqInf, err
}

// DeleteProfileByID DELETEs a Profile by ID
func (to *Session) DeleteProfileByID(id int) (tc.Alerts, ReqInf, error) {
	URI := fmt.Sprintf("%s/%d", v13Profiles, id)
	resp, remoteAddr, err := to.request(http.MethodDelete, URI, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()

	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)

	return alerts, reqInf, err
}

// ExportProfile Returns an exported Profile
func (to *Session) ExportProfile(id int) (*tc.ProfileExportResponse, ReqInf, error) {
	route := fmt.Sprintf("%s/%d/export", v13Profiles, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfileExportResponse
	err = json.NewDecoder(resp.Body).Decode(&data)

	return &data, reqInf, err
}

// ImportProfile imports an exported Profile
func (to *Session) ImportProfile(importRequest *tc.ProfileImportRequest) (*tc.ProfileImportResponse, ReqInf, error) {
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss}

	route := fmt.Sprintf("%s/import", v13Profiles)
	reqBody, err := json.Marshal(importRequest)
	if err != nil {
		return nil, reqInf, err
	}

	resp, remoteAddr, err := to.request(http.MethodPost, route, reqBody)
	reqInf.RemoteAddr = remoteAddr
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfileImportResponse
	err = json.NewDecoder(resp.Body).Decode(&data)

	return &data, reqInf, err
}

// CopyProfile creates a new profile from an existing profile.
func (to *Session) CopyProfile(p tc.ProfileCopy) (tc.ProfileCopyResponse, ReqInf, error) {
	var (
		path   = fmt.Sprintf("%s/name/%s/copy/%s", v14Profiles, p.Name, p.ExistingName)
		reqInf = ReqInf{CacheHitStatus: CacheHitStatusMiss}
		resp   tc.ProfileCopyResponse
	)

	reqBody, err := json.Marshal(p)
	if err != nil {
		return tc.ProfileCopyResponse{}, ReqInf{}, err
	}

	httpResp, remoteAddr, err := to.request(http.MethodPost, path, reqBody)
	reqInf.RemoteAddr = remoteAddr
	if err != nil {
		return tc.ProfileCopyResponse{}, reqInf, err
	}
	defer httpResp.Body.Close()

	err = json.NewDecoder(httpResp.Body).Decode(&resp)

	return resp, reqInf, err
}
