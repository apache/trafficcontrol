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
	"errors"
	"fmt"
	"net/url"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiRegions is the API version-relative path to the /regions API endpoint.
const apiRegions = "/regions"

// CreateRegion creates the given Region.
func (to *Session) CreateRegion(region tc.Region, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if region.Division == 0 && region.DivisionName != "" {
		divisionOpts := NewRequestOptions()
		divisionOpts.QueryParameters.Set("name", region.DivisionName)
		divisions, reqInf, err := to.GetDivisions(divisionOpts)
		if err != nil {
			return divisions.Alerts, reqInf, err
		}
		if len(divisions.Response) == 0 {
			return divisions.Alerts, reqInf, errors.New("no division with name " + region.DivisionName)
		}
		region.Division = divisions.Response[0].ID
	}
	var alerts tc.Alerts
	reqInf, err := to.post(apiRegions, opts, region, &alerts)
	return alerts, reqInf, err
}

// UpdateRegion replaces the Region identified by ID with the one provided.
func (to *Session) UpdateRegion(id int, region tc.Region, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", apiRegions, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, opts, region, &alerts)
	return alerts, reqInf, err
}

// GetRegions returns all Regions in Traffic Ops.
func (to *Session) GetRegions(opts RequestOptions) (tc.RegionsResponse, toclientlib.ReqInf, error) {
	var data tc.RegionsResponse
	reqInf, err := to.get(apiRegions, opts, &data)
	return data, reqInf, err
}

// DeleteRegion lets you delete a Region. Regions can be deleted by ID instead
// of by name if the ID is provided in the request options and the name is an
// empty string.
func (to *Session) DeleteRegion(name string, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}

	if name != "" || opts.QueryParameters.Get("id") == "" {
		opts.QueryParameters.Set("name", name)
	}

	var alerts tc.Alerts
	reqInf, err := to.del(apiRegions, opts, &alerts)
	return alerts, reqInf, err
}
