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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	// APIPhysLocations is the full path to the /phys_locations API route.
	APIPhysLocations = "/phys_locations"
)

// CreatePhysLocation creates the passed Physical Location.
func (to *Session) CreatePhysLocation(pl tc.PhysLocation) (tc.Alerts, toclientlib.ReqInf, error) {
	if pl.RegionID == 0 && pl.RegionName != "" {
		regions, _, err := to.GetRegionByName(pl.RegionName, nil)
		if err != nil {
			return tc.Alerts{}, toclientlib.ReqInf{}, err
		}
		if len(regions) == 0 {
			return tc.Alerts{}, toclientlib.ReqInf{}, errors.New("no region with name " + pl.RegionName)
		}
		pl.RegionID = regions[0].ID
	}
	var alerts tc.Alerts
	reqInf, err := to.post(APIPhysLocations, pl, nil, &alerts)
	return alerts, reqInf, err
}

// UpdatePhysLocation replaces the Physical Location identified by 'id' with
// the given Physical Location structure.
func (to *Session) UpdatePhysLocation(id int, pl tc.PhysLocation, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APIPhysLocations, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, pl, header, &alerts)
	return alerts, reqInf, err
}

// GetPhysLocations retrieves Physical Locations from Traffic Ops.
func (to *Session) GetPhysLocations(params url.Values, header http.Header) ([]tc.PhysLocation, toclientlib.ReqInf, error) {
	path := APIPhysLocations
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	var data tc.PhysLocationsResponse
	reqInf, err := to.get(path, header, &data)
	return data.Response, reqInf, err
}

// GetPhysLocationByID returns the Physical Location with the given ID.
func (to *Session) GetPhysLocationByID(id int, header http.Header) ([]tc.PhysLocation, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIPhysLocations, id)
	var data tc.PhysLocationsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetPhysLocationByName returns the Physical Location with the given Name.
func (to *Session) GetPhysLocationByName(name string, header http.Header) ([]tc.PhysLocation, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", APIPhysLocations, url.QueryEscape(name))
	var data tc.PhysLocationsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// DeletePhysLocation deletes the Physical Location with the given ID.
func (to *Session) DeletePhysLocation(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APIPhysLocations, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
