package client

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

import (
	"fmt"
	"net/url"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

const (
	// apiProfiles is the full path to the /profiles API endpoint.
	apiProfiles = "/profiles"
	// apiProfilesNameParameters is the full path to the
	// /profiles/name/{{name}}/parameters API endpoint.
	apiProfilesNameParameters = apiProfiles + "/name/%s/parameters"
)

// CreateProfile creates the passed Profile.
func (to *Session) CreateProfile(pl tc.ProfileV5, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if pl.CDNID == 0 && pl.CDNName != "" {
		cdnOpts := NewRequestOptions()
		cdnOpts.QueryParameters.Set("name", pl.CDNName)
		cdns, _, err := to.GetCDNs(cdnOpts)
		if err != nil {
			return tc.Alerts{}, toclientlib.ReqInf{}, fmt.Errorf("resolving Profile's CDN Name '%s' to an ID: %w", pl.CDNName, err)
		}
		if len(cdns.Response) == 0 {
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
		pl.CDNID = cdns.Response[0].ID
	}

	var alerts tc.Alerts
	reqInf, err := to.post(apiProfiles, opts, pl, &alerts)
	return alerts, reqInf, err
}

// UpdateProfile replaces the Profile identified by ID with the one provided.
func (to *Session) UpdateProfile(id int, pl tc.ProfileV5, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", apiProfiles, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, opts, pl, &alerts)
	return alerts, reqInf, err
}

// GetParametersByProfileName returns all of the Parameters that are assigned
// to the Profile with the given Name.
func (to *Session) GetParametersByProfileName(profileName string, opts RequestOptions) (tc.ParametersResponseV5, toclientlib.ReqInf, error) {
	route := fmt.Sprintf(apiProfilesNameParameters, profileName)
	var data tc.ParametersResponseV5
	reqInf, err := to.get(route, opts, &data)
	return data, reqInf, err
}

// GetProfiles returns all Profiles stored in Traffic Ops.
func (to *Session) GetProfiles(opts RequestOptions) (tc.ProfilesResponseV5, toclientlib.ReqInf, error) {
	var data tc.ProfilesResponseV5
	reqInf, err := to.get(apiProfiles, opts, &data)
	return data, reqInf, err
}

// DeleteProfile deletes the Profile with the given ID.
func (to *Session) DeleteProfile(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	URI := fmt.Sprintf("%s/%d", apiProfiles, id)
	var alerts tc.Alerts
	reqInf, err := to.del(URI, opts, &alerts)
	return alerts, reqInf, err
}

// ExportProfile Returns an exported Profile.
func (to *Session) ExportProfile(id int, opts RequestOptions) (tc.ProfileExportResponse, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d/export", apiProfiles, id)
	var data tc.ProfileExportResponse
	reqInf, err := to.get(route, opts, &data)
	return data, reqInf, err
}

// ImportProfile imports an exported Profile.
func (to *Session) ImportProfile(importRequest tc.ProfileImportRequest, opts RequestOptions) (tc.ProfileImportResponse, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/import", apiProfiles)
	var data tc.ProfileImportResponse
	reqInf, err := to.post(route, opts, importRequest, &data)
	return data, reqInf, err
}

// CopyProfile creates a new profile from an existing profile.
func (to *Session) CopyProfile(p tc.ProfileCopy, opts RequestOptions) (tc.ProfileCopyResponse, toclientlib.ReqInf, error) {
	path := fmt.Sprintf("%s/name/%s/copy/%s", apiProfiles, url.PathEscape(p.Name), url.PathEscape(p.ExistingName))
	var resp tc.ProfileCopyResponse
	reqInf, err := to.post(path, opts, p, &resp)
	return resp, reqInf, err
}
