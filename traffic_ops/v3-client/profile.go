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
	API_PROFILES                 = apiBase + "/profiles"
	API_PROFILES_NAME_PARAMETERS = API_PROFILES + "/name/%v/parameters"
)

// CreateProfile creates a Profile.
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

	resp, remoteAddr, err := to.request(http.MethodPost, API_PROFILES, reqBody, nil)
	reqInf.RemoteAddr = remoteAddr
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()

	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)

	return alerts, reqInf, err
}

func (to *Session) UpdateProfileByIDWithHdr(id int, pl tc.Profile, header http.Header) (tc.Alerts, ReqInf, error) {
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss}

	reqBody, err := json.Marshal(pl)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}

	route := fmt.Sprintf("%s/%d", API_PROFILES, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody, header)
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
	}
	reqInf.RemoteAddr = remoteAddr
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()

	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)

	return alerts, reqInf, err
}

// UpdateProfileByID updates a Profile by ID.
// Deprecated: UpdateProfileByID will be removed in 6.0. Use UpdateProfileByIDWithHdr.
func (to *Session) UpdateProfileByID(id int, pl tc.Profile) (tc.Alerts, ReqInf, error) {
	return to.UpdateProfileByIDWithHdr(id, pl, nil)
}

func (to *Session) GetParametersByProfileNameWithHdr(profileName string, header http.Header) ([]tc.Parameter, ReqInf, error) {
	url := fmt.Sprintf(API_PROFILES_NAME_PARAMETERS, profileName)
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Parameter{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ParametersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetParametersByProfileName gets all of the Parameters assigned to the Profile named 'profileName'.
// Deprecated: GetParametersByProfileName will be removed in 6.0. Use GetParametersByProfileNameWithHdr.
func (to *Session) GetParametersByProfileName(profileName string) ([]tc.Parameter, ReqInf, error) {
	return to.GetParametersByProfileNameWithHdr(profileName, nil)
}

func (to *Session) GetProfilesWithHdr(header http.Header) ([]tc.Profile, ReqInf, error) {
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss}

	resp, remoteAddr, err := to.request(http.MethodGet, API_PROFILES, nil, header)
	reqInf.RemoteAddr = remoteAddr
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Profile{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfilesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)

	return data.Response, reqInf, err
}

// GetProfiles returns a list of Profiles.
// Deprecated: GetProfiles will be removed in 6.0. Use GetProfilesWithHdr.
func (to *Session) GetProfiles() ([]tc.Profile, ReqInf, error) {
	return to.GetProfilesWithHdr(nil)
}

func (to *Session) GetProfileByIDWithHdr(id int, header http.Header) ([]tc.Profile, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_PROFILES, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Profile{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfilesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)

	return data.Response, reqInf, err
}

// GetProfileByID GETs a Profile by the Profile ID.
// Deprecated: GetProfileByID will be removed in 6.0. Use GetProfileByIDWithHdr.
func (to *Session) GetProfileByID(id int) ([]tc.Profile, ReqInf, error) {
	return to.GetProfileByIDWithHdr(id, nil)
}

func (to *Session) GetProfileByNameWithHdr(name string, header http.Header) ([]tc.Profile, ReqInf, error) {
	URI := fmt.Sprintf("%s?name=%s", API_PROFILES, url.QueryEscape(name))
	resp, remoteAddr, err := to.request(http.MethodGet, URI, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Profile{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfilesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)

	return data.Response, reqInf, err
}

// GetProfileByName GETs a Profile by the Profile name.
// Deprecated: GetProfileByName will be removed in 6.0. Use GetProfileByNameWithHdr.
func (to *Session) GetProfileByName(name string) ([]tc.Profile, ReqInf, error) {
	return to.GetProfileByNameWithHdr(name, nil)
}

func (to *Session) GetProfileByParameterWithHdr(param string, header http.Header) ([]tc.Profile, ReqInf, error) {
	URI := fmt.Sprintf("%s?param=%s", API_PROFILES, url.QueryEscape(param))
	resp, remoteAddr, err := to.request(http.MethodGet, URI, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Profile{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfilesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)

	return data.Response, reqInf, err
}

// GetProfileByParameter GETs a Profile by the Profile "param".
// Deprecated: GetProfileByParameter will be removed in 6.0. Use GetProfileByParameterWithHdr.
func (to *Session) GetProfileByParameter(param string) ([]tc.Profile, ReqInf, error) {
	return to.GetProfileByParameterWithHdr(param, nil)
}

func (to *Session) GetProfileByCDNIDWithHdr(cdnID int, header http.Header) ([]tc.Profile, ReqInf, error) {
	URI := fmt.Sprintf("%s?cdn=%s", API_PROFILES, strconv.Itoa(cdnID))
	resp, remoteAddr, err := to.request(http.MethodGet, URI, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Profile{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfilesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)

	return data.Response, reqInf, err
}

// GetProfileByCDNID GETs a Profile by the Profile CDN ID.
// Deprecated: GetProfileByCDNID will be removed in 6.0. Use GetProfileByCDNIDWithHdr.
func (to *Session) GetProfileByCDNID(cdnID int) ([]tc.Profile, ReqInf, error) {
	return to.GetProfileByCDNIDWithHdr(cdnID, nil)
}

// DeleteProfileByID DELETEs a Profile by ID.
func (to *Session) DeleteProfileByID(id int) (tc.Alerts, ReqInf, error) {
	URI := fmt.Sprintf("%s/%d", API_PROFILES, id)
	resp, remoteAddr, err := to.request(http.MethodDelete, URI, nil, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()

	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)

	return alerts, reqInf, err
}

// ExportProfile Returns an exported Profile.
func (to *Session) ExportProfile(id int) (*tc.ProfileExportResponse, ReqInf, error) {
	route := fmt.Sprintf("%s/%d/export", API_PROFILES, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfileExportResponse
	err = json.NewDecoder(resp.Body).Decode(&data)

	return &data, reqInf, err
}

// ImportProfile imports an exported Profile.
func (to *Session) ImportProfile(importRequest *tc.ProfileImportRequest) (*tc.ProfileImportResponse, ReqInf, error) {
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss}

	route := fmt.Sprintf("%s/import", API_PROFILES)
	reqBody, err := json.Marshal(importRequest)
	if err != nil {
		return nil, reqInf, err
	}

	resp, remoteAddr, err := to.request(http.MethodPost, route, reqBody, nil)
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
		path   = fmt.Sprintf("%s/name/%s/copy/%s", API_PROFILES, p.Name, p.ExistingName)
		reqInf = ReqInf{CacheHitStatus: CacheHitStatusMiss}
		resp   tc.ProfileCopyResponse
	)

	reqBody, err := json.Marshal(p)
	if err != nil {
		return tc.ProfileCopyResponse{}, ReqInf{}, err
	}

	httpResp, remoteAddr, err := to.request(http.MethodPost, path, reqBody, nil)
	reqInf.RemoteAddr = remoteAddr
	if err != nil {
		return tc.ProfileCopyResponse{}, reqInf, err
	}
	defer httpResp.Body.Close()

	err = json.NewDecoder(httpResp.Body).Decode(&resp)

	return resp, reqInf, err
}
