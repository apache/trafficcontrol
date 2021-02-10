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
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	// APIProfiles is the full path to the /profiles API endpoint.
	APIProfiles = "/profiles"
	// APIProfilesNameParameters is the full path to the
	// /profiles/name/{{name}}/parameters API endpoint.
	APIProfilesNameParameters = APIProfiles + "/name/%s/parameters"
)

// CreateProfile creates the passed Profile.
func (to *Session) CreateProfile(pl tc.Profile) (tc.Alerts, toclientlib.ReqInf, error) {
	if pl.CDNID == 0 && pl.CDNName != "" {
		cdns, _, err := to.GetCDNByName(pl.CDNName, nil)
		if err != nil {
			return tc.Alerts{}, toclientlib.ReqInf{}, err
		}
		if len(cdns) == 0 {
			return tc.Alerts{
					Alerts: []tc.Alert{
						{
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

// UpdateProfile replaces the Profile identified by ID with the one provided.
func (to *Session) UpdateProfile(id int, pl tc.Profile, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APIProfiles, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, pl, header, &alerts)
	return alerts, reqInf, err
}

// GetParametersByProfileName returns all of the Parameters that are assigned
// to the Profile with the given Name.
func (to *Session) GetParametersByProfileName(profileName string, header http.Header) ([]tc.Parameter, toclientlib.ReqInf, error) {
	route := fmt.Sprintf(APIProfilesNameParameters, profileName)
	var data tc.ParametersResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetProfiles returns all Profiles stored in Traffic Ops.
func (to *Session) GetProfiles(header http.Header) ([]tc.Profile, toclientlib.ReqInf, error) {
	var data tc.ProfilesResponse
	reqInf, err := to.get(APIProfiles, header, &data)
	return data.Response, reqInf, err
}

// GetProfileByID retrieves the Profile with the given ID.
func (to *Session) GetProfileByID(id int, header http.Header) ([]tc.Profile, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIProfiles, id)
	var data tc.ProfilesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetProfileByName retrieves the Profile with the given Name.
func (to *Session) GetProfileByName(name string, header http.Header) ([]tc.Profile, toclientlib.ReqInf, error) {
	URI := fmt.Sprintf("%s?name=%s", APIProfiles, url.QueryEscape(name))
	var data tc.ProfilesResponse
	reqInf, err := to.get(URI, header, &data)
	return data.Response, reqInf, err
}

// GetProfilesByParameterID retrieves the Profiles to which the Parameter
// identified by the ID 'param' is assigned.
func (to *Session) GetProfilesByParameterID(param int, header http.Header) ([]tc.Profile, toclientlib.ReqInf, error) {
	URI := fmt.Sprintf("%s?param=%d", APIProfiles, param)
	var data tc.ProfilesResponse
	reqInf, err := to.get(URI, header, &data)
	return data.Response, reqInf, err
}

// GetProfilesByCDNID retrieves all Profiles that are within the scope of the
// CDN identified by 'cdnID'.
func (to *Session) GetProfilesByCDNID(cdnID int, header http.Header) ([]tc.Profile, toclientlib.ReqInf, error) {
	URI := fmt.Sprintf("%s?cdn=%d", APIProfiles, cdnID)
	var data tc.ProfilesResponse
	reqInf, err := to.get(URI, header, &data)
	return data.Response, reqInf, err
}

// DeleteProfile deletes the Profile with the given ID.
func (to *Session) DeleteProfile(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	URI := fmt.Sprintf("%s/%d", APIProfiles, id)
	var alerts tc.Alerts
	reqInf, err := to.del(URI, nil, &alerts)
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
