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
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

const (
	API_COORDINATES = apiBase + "/coordinates"
)

// Create a Coordinate
func (to *Session) CreateCoordinate(coordinate tc.Coordinate) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(coordinate)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_COORDINATES, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Update a Coordinate by ID
func (to *Session) UpdateCoordinateByID(id int, coordinate tc.Coordinate) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(coordinate)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := fmt.Sprintf("%s?id=%d", API_COORDINATES, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Returns a list of Coordinates
func (to *Session) GetCoordinates() ([]tc.Coordinate, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_COORDINATES, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CoordinatesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// GET a Coordinate by the Coordinate id
func (to *Session) GetCoordinateByID(id int) ([]tc.Coordinate, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_COORDINATES, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CoordinatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GET a Coordinate by the Coordinate name
func (to *Session) GetCoordinateByName(name string) ([]tc.Coordinate, ReqInf, error) {
	url := fmt.Sprintf("%s?name=%s", API_COORDINATES, name)
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CoordinatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// DELETE a Coordinate by ID
func (to *Session) DeleteCoordinateByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_COORDINATES, id)
	resp, remoteAddr, err := to.request(http.MethodDelete, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

func (to *Session) SetCachegroupDeliveryServices(cgID int, dsIDs []int) (tc.CacheGroupPostDSRespResponse, ReqInf, error) {
	uri := apiBase + `/cachegroups/` + strconv.Itoa(cgID) + `/deliveryservices`
	req := tc.CachegroupPostDSReq{DeliveryServices: dsIDs}
	reqBody, err := json.Marshal(req)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss}
	if err != nil {
		return tc.CacheGroupPostDSRespResponse{}, reqInf, err
	}
	reqResp, _, err := to.request(http.MethodPost, uri, reqBody)
	if err != nil {
		return tc.CacheGroupPostDSRespResponse{}, reqInf, errors.New("requesting from Traffic Ops: " + err.Error())
	}
	defer reqResp.Body.Close()

	resp := tc.CacheGroupPostDSRespResponse{}
	if err := json.NewDecoder(reqResp.Body).Decode(&resp); err != nil {
		return tc.CacheGroupPostDSRespResponse{}, reqInf, errors.New("decoding response: " + err.Error())
	}
	return resp, reqInf, nil
}
