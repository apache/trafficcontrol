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
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_REGIONS = apiBase + "/regions"
)

// CreateRegion creates a Region.
func (to *Session) CreateRegion(region tc.Region) (tc.Alerts, ReqInf, error) {
	if region.Division == 0 && region.DivisionName != "" {
		divisions, _, err := to.GetDivisionByNameWithHdr(region.DivisionName, nil)
		if err != nil {
			return tc.Alerts{}, ReqInf{}, err
		}
		if len(divisions) == 0 {
			return tc.Alerts{}, ReqInf{}, errors.New("no division with name " + region.DivisionName)
		}
		region.Division = divisions[0].ID
	}
	var alerts tc.Alerts
	reqInf, err := to.post(API_REGIONS, region, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) UpdateRegionByIDWithHdr(id int, region tc.Region, header http.Header) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_REGIONS, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, region, header, &alerts)
	return alerts, reqInf, err
}

// UpdateRegionByID updates a Region by ID.
// Deprecated: UpdateRegionByID will be removed in 6.0. Use UpdateRegionByIDWithHdr.
func (to *Session) UpdateRegionByID(id int, region tc.Region) (tc.Alerts, ReqInf, error) {
	return to.UpdateRegionByIDWithHdr(id, region, nil)
}

func (to *Session) GetRegionsWithHdr(header http.Header) ([]tc.Region, ReqInf, error) {
	var data tc.RegionsResponse
	reqInf, err := to.get(API_REGIONS, header, &data)
	return data.Response, reqInf, err
}

// GetRegions returns a list of regions.
// Deprecated: GetRegions will be removed in 6.0. Use GetRegionsWithHdr.
func (to *Session) GetRegions() ([]tc.Region, ReqInf, error) {
	return to.GetRegionsWithHdr(nil)
}

func (to *Session) GetRegionByIDWithHdr(id int, header http.Header) ([]tc.Region, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_REGIONS, id)
	var data tc.RegionsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetRegionByID GETs a Region by the Region ID.
// Deprecated: GetRegionByID will be removed in 6.0. Use GetRegionByIDWithHdr.
func (to *Session) GetRegionByID(id int) ([]tc.Region, ReqInf, error) {
	return to.GetRegionByIDWithHdr(id, nil)
}

func (to *Session) GetRegionByNameWithHdr(name string, header http.Header) ([]tc.Region, ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", API_REGIONS, url.QueryEscape(name))
	var data tc.RegionsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetRegionByName GETs a Region by the Region name.
// Deprecated: GetRegionByName will be removed in 6.0. Use GetRegionByNameHdr.
func (to *Session) GetRegionByName(name string) ([]tc.Region, ReqInf, error) {
	return to.GetRegionByNameWithHdr(name, nil)
}

// DeleteRegionByID DELETEs a Region by ID.
func (to *Session) DeleteRegionByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_REGIONS, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

// DeleteRegion lets you DELETE a Region. Only 1 parameter is required, not both.
func (to *Session) DeleteRegion(id *int, name *string) (tc.Alerts, ReqInf, error) {
	v := url.Values{}
	if id != nil {
		v.Add("id", strconv.Itoa(*id))
	}
	if name != nil {
		v.Add("name", *name)
	}
	URI := apiBase + "/regions"
	if qStr := v.Encode(); len(qStr) > 0 {
		URI = fmt.Sprintf("%s?%s", URI, qStr)
	}

	var alerts tc.Alerts
	reqInf, err := to.del(URI, nil, &alerts)
	return alerts, reqInf, err
}
