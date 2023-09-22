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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiPhysLocations is the full path to the /phys_locations API route.
const apiPhysLocations = "/phys_locations"

// CreatePhysLocation creates the passed Physical Location.
func (to *Session) CreatePhysLocation(pl tc.PhysLocation, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if pl.RegionID == 0 && pl.RegionName != "" {
		regionOpts := NewRequestOptions()
		regionOpts.QueryParameters.Set("name", pl.RegionName)
		regions, reqInf, err := to.GetRegions(regionOpts)
		if err != nil {
			err = fmt.Errorf("resolving Region name '%s' to an ID", pl.RegionName)
			return regions.Alerts, reqInf, err
		}
		if len(regions.Response) == 0 {
			return regions.Alerts, reqInf, fmt.Errorf("no region with name '%s'", pl.RegionName)
		}
		pl.RegionID = regions.Response[0].ID
	}
	var alerts tc.Alerts
	reqInf, err := to.post(apiPhysLocations, opts, pl, &alerts)
	return alerts, reqInf, err
}

// UpdatePhysLocation replaces the Physical Location identified by 'id' with
// the given Physical Location structure.
func (to *Session) UpdatePhysLocation(id int, pl tc.PhysLocation, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", apiPhysLocations, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, opts, pl, &alerts)
	return alerts, reqInf, err
}

// GetPhysLocations retrieves Physical Locations from Traffic Ops.
func (to *Session) GetPhysLocations(opts RequestOptions) (tc.PhysLocationsResponse, toclientlib.ReqInf, error) {
	var data tc.PhysLocationsResponse
	reqInf, err := to.get(apiPhysLocations, opts, &data)
	return data, reqInf, err
}

// DeletePhysLocation deletes the Physical Location with the given ID.
func (to *Session) DeletePhysLocation(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", apiPhysLocations, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, opts, &alerts)
	return alerts, reqInf, err
}
