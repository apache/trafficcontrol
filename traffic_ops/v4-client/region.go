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
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	// APIRegions is the API version-relative path to the /regions API endpoint.
	APIRegions = "/regions"
)

// CreateRegion creates the given Region.
func (to *Session) CreateRegion(region tc.Region) (tc.Alerts, toclientlib.ReqInf, error) {
	if region.Division == 0 && region.DivisionName != "" {
		divisions, _, err := to.GetDivisionByName(region.DivisionName, nil)
		if err != nil {
			return tc.Alerts{}, toclientlib.ReqInf{}, err
		}
		if len(divisions) == 0 {
			return tc.Alerts{}, toclientlib.ReqInf{}, errors.New("no division with name " + region.DivisionName)
		}
		region.Division = divisions[0].ID
	}
	var alerts tc.Alerts
	reqInf, err := to.post(APIRegions, region, nil, &alerts)
	return alerts, reqInf, err
}

// UpdateRegion replaces the Region identified by ID with the one provided.
func (to *Session) UpdateRegion(id int, region tc.Region, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APIRegions, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, region, header, &alerts)
	return alerts, reqInf, err
}

// GetRegions returns all Regions in Traffic Ops.
func (to *Session) GetRegions(params url.Values, header http.Header) ([]tc.Region, toclientlib.ReqInf, error) {
	uri := APIRegions
	if params != nil {
		uri += "?" + params.Encode()
	}
	var data tc.RegionsResponse
	reqInf, err := to.get(uri, header, &data)
	return data.Response, reqInf, err
}

// GetRegionByID returns the Region with the given ID.
func (to *Session) GetRegionByID(id int, header http.Header) ([]tc.Region, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIRegions, id)
	var data tc.RegionsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetRegionByName retrieves the Region with the given Name.
func (to *Session) GetRegionByName(name string, header http.Header) ([]tc.Region, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", APIRegions, url.QueryEscape(name))
	var data tc.RegionsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetRegionByDivision retrieves the Region with the given Name.
func (to *Session) GetRegionByDivision(id int, header http.Header) ([]tc.Region, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?division=%d", APIRegions, id)
	var data tc.RegionsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// DeleteRegion lets you DELETE a Region. Only 1 parameter is required, not both.
func (to *Session) DeleteRegion(id *int, name *string) (tc.Alerts, toclientlib.ReqInf, error) {
	v := url.Values{}
	if id != nil {
		v.Add("id", strconv.Itoa(*id))
	}
	if name != nil {
		v.Add("name", *name)
	}
	URI := "/regions"
	if qStr := v.Encode(); len(qStr) > 0 {
		URI = fmt.Sprintf("%s?%s", URI, qStr)
	}

	var alerts tc.Alerts
	reqInf, err := to.del(URI, nil, &alerts)
	return alerts, reqInf, err
}
