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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

const (
	// API_COORDINATES is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
	API_COORDINATES = apiBase + "/coordinates"

	APICoordinates = "/coordinates"
)

// Create a Coordinate
func (to *Session) CreateCoordinate(coordinate tc.Coordinate) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APICoordinates, coordinate, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) UpdateCoordinateByIDWithHdr(id int, coordinate tc.Coordinate, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APICoordinates, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, coordinate, header, &alerts)
	return alerts, reqInf, err
}

// Update a Coordinate by ID
// Deprecated: UpdateCoordinateByID will be removed in 6.0. Use UpdateCoordinateByIDWithHdr.
func (to *Session) UpdateCoordinateByID(id int, coordinate tc.Coordinate) (tc.Alerts, toclientlib.ReqInf, error) {
	return to.UpdateCoordinateByIDWithHdr(id, coordinate, nil)
}

func (to *Session) GetCoordinatesWithHdr(header http.Header) ([]tc.Coordinate, toclientlib.ReqInf, error) {
	var data tc.CoordinatesResponse
	reqInf, err := to.get(APICoordinates, header, &data)
	return data.Response, reqInf, err
}

// Returns a list of Coordinates
// Deprecated: GetCoordinates will be removed in 6.0. Use GetCoordinatesWithHdr.
func (to *Session) GetCoordinates() ([]tc.Coordinate, toclientlib.ReqInf, error) {
	return to.GetCoordinatesWithHdr(nil)
}

func (to *Session) GetCoordinateByIDWithHdr(id int, header http.Header) ([]tc.Coordinate, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APICoordinates, id)
	var data tc.CoordinatesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GET a Coordinate by the Coordinate id
// Deprecated: GetCoordinateByID will be removed in 6.0. Use GetCoordinateByIDWithHdr.
func (to *Session) GetCoordinateByID(id int) ([]tc.Coordinate, toclientlib.ReqInf, error) {
	return to.GetCoordinateByIDWithHdr(id, nil)
}

// GET a Coordinate by the Coordinate name
// Deprecated: GetCoordinateByName will be removed in 6.0. Use GetCoordinateByNameWithHdr.
func (to *Session) GetCoordinateByName(name string) ([]tc.Coordinate, toclientlib.ReqInf, error) {
	return to.GetCoordinateByNameWithHdr(name, nil)
}

func (to *Session) GetCoordinateByNameWithHdr(name string, header http.Header) ([]tc.Coordinate, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", APICoordinates, name)
	var data tc.CoordinatesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// DELETE a Coordinate by ID
func (to *Session) DeleteCoordinateByID(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APICoordinates, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
