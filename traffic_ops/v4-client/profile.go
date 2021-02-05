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
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_PROFILES                 = apiBase + "/profiles"
	API_PROFILES_NAME_PARAMETERS = API_PROFILES + "/name/%s/parameters"
)

// CreateProfile creates a Profile.
func (to *Session) CreateProfile(pl tc.Profile) (tc.Alerts, ReqInf, error) {
	if pl.CDNID == 0 && pl.CDNName != "" {
		cdns, _, err := to.GetCDNByNameWithHdr(pl.CDNName, nil)
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

	var alerts tc.Alerts
	reqInf, err := to.post(API_PROFILES, pl, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) UpdateProfileByIDWithHdr(id int, pl tc.Profile, header http.Header) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_PROFILES, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, pl, header, &alerts)
	return alerts, reqInf, err
}

// UpdateProfileByID updates a Profile by ID.
// Deprecated: UpdateProfileByID will be removed in 6.0. Use UpdateProfileByIDWithHdr.
func (to *Session) UpdateProfileByID(id int, pl tc.Profile) (tc.Alerts, ReqInf, error) {
	return to.UpdateProfileByIDWithHdr(id, pl, nil)
}

func (to *Session) GetParametersByProfileNameWithHdr(profileName string, header http.Header) ([]tc.Parameter, ReqInf, error) {
	route := fmt.Sprintf(API_PROFILES_NAME_PARAMETERS, profileName)
	var data tc.ParametersResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetParametersByProfileName gets all of the Parameters assigned to the Profile named 'profileName'.
// Deprecated: GetParametersByProfileName will be removed in 6.0. Use GetParametersByProfileNameWithHdr.
func (to *Session) GetParametersByProfileName(profileName string) ([]tc.Parameter, ReqInf, error) {
	return to.GetParametersByProfileNameWithHdr(profileName, nil)
}

func (to *Session) GetProfilesWithHdr(header http.Header) ([]tc.Profile, ReqInf, error) {
	var data tc.ProfilesResponse
	reqInf, err := to.get(API_PROFILES, header, &data)
	return data.Response, reqInf, err
}

// GetProfiles returns a list of Profiles.
// Deprecated: GetProfiles will be removed in 6.0. Use GetProfilesWithHdr.
func (to *Session) GetProfiles() ([]tc.Profile, ReqInf, error) {
	return to.GetProfilesWithHdr(nil)
}

func (to *Session) GetProfileByIDWithHdr(id int, header http.Header) ([]tc.Profile, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_PROFILES, id)
	var data tc.ProfilesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetProfileByID GETs a Profile by the Profile ID.
// Deprecated: GetProfileByID will be removed in 6.0. Use GetProfileByIDWithHdr.
func (to *Session) GetProfileByID(id int) ([]tc.Profile, ReqInf, error) {
	return to.GetProfileByIDWithHdr(id, nil)
}

func (to *Session) GetProfileByNameWithHdr(name string, header http.Header) ([]tc.Profile, ReqInf, error) {
	URI := fmt.Sprintf("%s?name=%s", API_PROFILES, url.QueryEscape(name))
	var data tc.ProfilesResponse
	reqInf, err := to.get(URI, header, &data)
	return data.Response, reqInf, err
}

// GetProfileByName GETs a Profile by the Profile name.
// Deprecated: GetProfileByName will be removed in 6.0. Use GetProfileByNameWithHdr.
func (to *Session) GetProfileByName(name string) ([]tc.Profile, ReqInf, error) {
	return to.GetProfileByNameWithHdr(name, nil)
}

func (to *Session) GetProfileByParameterWithHdr(param string, header http.Header) ([]tc.Profile, ReqInf, error) {
	URI := fmt.Sprintf("%s?param=%s", API_PROFILES, url.QueryEscape(param))
	var data tc.ProfilesResponse
	reqInf, err := to.get(URI, header, &data)
	return data.Response, reqInf, err
}

// GetProfileByParameter GETs a Profile by the Profile "param".
// Deprecated: GetProfileByParameter will be removed in 6.0. Use GetProfileByParameterWithHdr.
func (to *Session) GetProfileByParameter(param string) ([]tc.Profile, ReqInf, error) {
	return to.GetProfileByParameterWithHdr(param, nil)
}

func (to *Session) GetProfileByCDNIDWithHdr(cdnID int, header http.Header) ([]tc.Profile, ReqInf, error) {
	URI := fmt.Sprintf("%s?cdn=%d", API_PROFILES, cdnID)
	var data tc.ProfilesResponse
	reqInf, err := to.get(URI, header, &data)
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
	var alerts tc.Alerts
	reqInf, err := to.del(URI, nil, &alerts)
	return alerts, reqInf, err
}

// ExportProfile Returns an exported Profile.
func (to *Session) ExportProfile(id int) (*tc.ProfileExportResponse, ReqInf, error) {
	route := fmt.Sprintf("%s/%d/export", API_PROFILES, id)
	var data tc.ProfileExportResponse
	reqInf, err := to.get(route, nil, &data)
	return &data, reqInf, err
}

// ImportProfile imports an exported Profile.
func (to *Session) ImportProfile(importRequest *tc.ProfileImportRequest) (*tc.ProfileImportResponse, ReqInf, error) {
	route := fmt.Sprintf("%s/import", API_PROFILES)
	var data tc.ProfileImportResponse
	reqInf, err := to.post(route, importRequest, nil, &data)
	return &data, reqInf, err
}

// CopyProfile creates a new profile from an existing profile.
func (to *Session) CopyProfile(p tc.ProfileCopy) (tc.ProfileCopyResponse, ReqInf, error) {
	path := fmt.Sprintf("%s/name/%s/copy/%s", API_PROFILES, p.Name, p.ExistingName)
	resp := tc.ProfileCopyResponse{}
	reqInf, err := to.post(path, p, nil, &resp)
	return resp, reqInf, err
}
