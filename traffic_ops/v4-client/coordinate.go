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
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiCoordinates is the API version-relative path for the /coordinates API endpoint.
const apiCoordinates = "/coordinates"

// CreateCoordinate creates the given Coordinate.
func (to *Session) CreateCoordinate(coordinate tc.Coordinate, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiCoordinates, opts, coordinate, &alerts)
	return alerts, reqInf, err
}

// UpdateCoordinate replaces the Coordinate with the given ID with the one
// provided.
func (to *Session) UpdateCoordinate(id int, coordinate tc.Coordinate, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("id", strconv.Itoa(id))
	var alerts tc.Alerts
	reqInf, err := to.put(apiCoordinates, opts, coordinate, &alerts)
	return alerts, reqInf, err
}

// GetCoordinates returns all Coordinates in Traffic Ops.
func (to *Session) GetCoordinates(opts RequestOptions) (tc.CoordinatesResponse, toclientlib.ReqInf, error) {
	var data tc.CoordinatesResponse
	reqInf, err := to.get(apiCoordinates, opts, &data)
	return data, reqInf, err
}

// DeleteCoordinate deletes the Coordinate with the given ID.
func (to *Session) DeleteCoordinate(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("id", strconv.Itoa(id))
	var alerts tc.Alerts
	reqInf, err := to.del(apiCoordinates, opts, &alerts)
	return alerts, reqInf, err
}
