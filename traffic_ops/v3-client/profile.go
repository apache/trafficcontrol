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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

const (
	// API_PROFILES is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
	API_PROFILES = apiBase + "/profiles"

	// API_PROFILES_NAME_PARAMETERS is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
	API_PROFILES_NAME_PARAMETERS = API_PROFILES + "/name/%s/parameters"

	APIProfiles               = "/profiles"
	APIProfilesNameParameters = APIProfiles + "/name/%s/parameters"
)

// CreateProfile creates a Profile.
func (to *Session) CreateProfile(pl tc.Profile) (tc.Alerts, toclientlib.ReqInf, error) {
	if pl.CDNID == 0 && pl.CDNName != "" {
		cdns, _, err := to.GetCDNByNameWithHdr(pl.CDNName, nil)
		if err != nil {
			return tc.Alerts{}, toclientlib.ReqInf{}, err
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
				toclientlib.ReqInf{},
				fmt.Errorf("no CDN with name %s", pl.CDNName)
		}
		pl.CDNID = cdns[0].ID
	}

	var alerts tc.Alerts
	reqInf, err := to.post(APIProfiles, pl, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) UpdateProfileByIDWithHdr(id int, pl tc.Profile, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APIProfiles, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, pl, header, &alerts)
	return alerts, reqInf, err
}

// UpdateProfileByID updates a Profile by ID.
// Deprecated: UpdateProfileByID will be removed in 6.0. Use UpdateProfileByIDWithHdr.
func (to *Session) UpdateProfileByID(id int, pl tc.Profile) (tc.Alerts, toclientlib.ReqInf, error) {
	return to.UpdateProfileByIDWithHdr(id, pl, nil)
}

func (to *Session) GetParametersByProfileNameWithHdr(profileName string, header http.Header) ([]tc.Parameter, toclientlib.ReqInf, error) {
	route := fmt.Sprintf(APIProfilesNameParameters, profileName)
	var data tc.ParametersResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetParametersByProfileName gets all of the Parameters assigned to the Profile named 'profileName'.
// Deprecated: GetParametersByProfileName will be removed in 6.0. Use GetParametersByProfileNameWithHdr.
func (to *Session) GetParametersByProfileName(profileName string) ([]tc.Parameter, toclientlib.ReqInf, error) {
	return to.GetParametersByProfileNameWithHdr(profileName, nil)
}

func (to *Session) GetProfilesWithHdr(header http.Header) ([]tc.Profile, toclientlib.ReqInf, error) {
	var data tc.ProfilesResponse
	reqInf, err := to.get(APIProfiles, header, &data)
	return data.Response, reqInf, err
}

// GetProfiles returns a list of Profiles.
// Deprecated: GetProfiles will be removed in 6.0. Use GetProfilesWithHdr.
func (to *Session) GetProfiles() ([]tc.Profile, toclientlib.ReqInf, error) {
	return to.GetProfilesWithHdr(nil)
}

func (to *Session) GetProfileByIDWithHdr(id int, header http.Header) ([]tc.Profile, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIProfiles, id)
	var data tc.ProfilesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetProfileByID GETs a Profile by the Profile ID.
// Deprecated: GetProfileByID will be removed in 6.0. Use GetProfileByIDWithHdr.
func (to *Session) GetProfileByID(id int) ([]tc.Profile, toclientlib.ReqInf, error) {
	return to.GetProfileByIDWithHdr(id, nil)
}

func (to *Session) GetProfileByNameWithHdr(name string, header http.Header) ([]tc.Profile, toclientlib.ReqInf, error) {
	URI := fmt.Sprintf("%s?name=%s", APIProfiles, url.QueryEscape(name))
	var data tc.ProfilesResponse
	reqInf, err := to.get(URI, header, &data)
	return data.Response, reqInf, err
}

// GetProfileByName GETs a Profile by the Profile name.
// Deprecated: GetProfileByName will be removed in 6.0. Use GetProfileByNameWithHdr.
func (to *Session) GetProfileByName(name string) ([]tc.Profile, toclientlib.ReqInf, error) {
	return to.GetProfileByNameWithHdr(name, nil)
}

func (to *Session) GetProfileByParameterWithHdr(param string, header http.Header) ([]tc.Profile, toclientlib.ReqInf, error) {
	URI := fmt.Sprintf("%s?param=%s", APIProfiles, url.QueryEscape(param))
	var data tc.ProfilesResponse
	reqInf, err := to.get(URI, header, &data)
	return data.Response, reqInf, err
}

// GetProfileByParameter GETs a Profile by the Profile "param".
// Deprecated: GetProfileByParameter will be removed in 6.0. Use GetProfileByParameterWithHdr.
func (to *Session) GetProfileByParameter(param string) ([]tc.Profile, toclientlib.ReqInf, error) {
	return to.GetProfileByParameterWithHdr(param, nil)
}

// GetProfileByParameterIdWithHdr GETs a Profile by the ParameterID and Header.
func (to *Session) GetProfileByParameterIdWithHdr(param int, header http.Header) ([]tc.Profile, toclientlib.ReqInf, error) {
	URI := fmt.Sprintf("%s?param=%d", APIProfiles, param)
	var data tc.ProfilesResponse
	reqInf, err := to.get(URI, header, &data)
	return data.Response, reqInf, err
}

// GetProfileByParameterId GETs a Profile by the Profile "param".
// Deprecated: GetProfileByParameterId will be removed in 6.0. Use GetProfileByParameterIdWithHdr.
func (to *Session) GetProfileByParameterId(param int) ([]tc.Profile, toclientlib.ReqInf, error) {
	return to.GetProfileByParameterIdWithHdr(param, nil)
}

func (to *Session) GetProfileByCDNIDWithHdr(cdnID int, header http.Header) ([]tc.Profile, toclientlib.ReqInf, error) {
	uri := fmt.Sprintf("%s?cdn=%d", APIProfiles, cdnID)
	var data tc.ProfilesResponse
	reqInf, err := to.get(uri, header, &data)
	return data.Response, reqInf, err
}

// GetProfileByCDNID GETs a Profile by the Profile CDN ID.
// Deprecated: GetProfileByCDNID will be removed in 6.0. Use GetProfileByCDNIDWithHdr.
func (to *Session) GetProfileByCDNID(cdnID int) ([]tc.Profile, toclientlib.ReqInf, error) {
	return to.GetProfileByCDNIDWithHdr(cdnID, nil)
}

// DeleteProfileByID DELETEs a Profile by ID.
func (to *Session) DeleteProfileByID(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	uri := fmt.Sprintf("%s/%d", APIProfiles, id)
	var alerts tc.Alerts
	reqInf, err := to.del(uri, nil, &alerts)
	return alerts, reqInf, err
}

// ExportProfile Returns an exported Profile.
func (to *Session) ExportProfile(id int) (*tc.ProfileExportResponse, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d/export", APIProfiles, id)
	var data tc.ProfileExportResponse
	reqInf, err := to.get(route, nil, &data)
	return &data, reqInf, err
}

// ImportProfile imports an exported Profile.
func (to *Session) ImportProfile(importRequest *tc.ProfileImportRequest) (*tc.ProfileImportResponse, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/import", APIProfiles)
	var data tc.ProfileImportResponse
	reqInf, err := to.post(route, importRequest, nil, &data)
	return &data, reqInf, err
}

// CopyProfile creates a new profile from an existing profile.
func (to *Session) CopyProfile(p tc.ProfileCopy) (tc.ProfileCopyResponse, toclientlib.ReqInf, error) {
	path := fmt.Sprintf("%s/name/%s/copy/%s", APIProfiles, p.Name, p.ExistingName)
	resp := tc.ProfileCopyResponse{}
	reqInf, err := to.post(path, p, nil, &resp)
	return resp, reqInf, err
}
