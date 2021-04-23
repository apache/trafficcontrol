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
	// APICoordinates is the API version-relative path for the /coordinates API endpoint.
	APICoordinates = "/coordinates"
)

// CreateCoordinate creates the given Coordinate.
func (to *Session) CreateCoordinate(coordinate tc.Coordinate) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APICoordinates, coordinate, nil, &alerts)
	return alerts, reqInf, err
}

// UpdateCoordinate replaces the Coordinate with the given ID with the one
// provided.
func (to *Session) UpdateCoordinate(id int, coordinate tc.Coordinate, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APICoordinates, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, coordinate, header, &alerts)
	return alerts, reqInf, err
}

// GetCoordinates returns all Coordinates in Traffic Ops.
func (to *Session) GetCoordinates(params url.Values, header http.Header) ([]tc.Coordinate, toclientlib.ReqInf, error) {
	uri := APICoordinates
	if params != nil {
		uri += "?" + params.Encode()
	}
	var data tc.CoordinatesResponse
	reqInf, err := to.get(uri, header, &data)
	return data.Response, reqInf, err
}

// GetCoordinateByID retrieves the Coordinate with the given ID.
func (to *Session) GetCoordinateByID(id int, header http.Header) ([]tc.Coordinate, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APICoordinates, id)
	var data tc.CoordinatesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetCoordinateByName retrieves the Coordinate with the given Name.
func (to *Session) GetCoordinateByName(name string, header http.Header) ([]tc.Coordinate, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", APICoordinates, name)
	var data tc.CoordinatesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// DeleteCoordinateByID deletes the Coordinate with the given ID.
func (to *Session) DeleteCoordinate(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APICoordinates, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
