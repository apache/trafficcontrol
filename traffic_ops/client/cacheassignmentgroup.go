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
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"net"
	"net/http"
)

const (
	API_V14_CacheAssignmentGroups = "/api/1.4/cacheassignmentgroups/"
)

func (to *Session) CreateCacheAssignmentGroup(cag tc.CacheAssignmentGroup) (tc.Alerts, ReqInf, error) {
	var remoteAddr net.Addr
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}

	reqBody, err := json.Marshal(cag)

	resp, remoteAddr, err := to.request(http.MethodPost, API_V14_CacheAssignmentGroups, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}

	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

func (to *Session) UpdateCacheAssignmentGroupByID(id int, cag tc.CacheAssignmentGroup) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(cag)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := fmt.Sprintf("%s?id=%d", API_V14_CacheAssignmentGroups, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}


func (to *Session) GetCacheAssignmentGroups() ([]tc.CacheAssignmentGroup, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_V14_CacheAssignmentGroups, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CacheAssignmentGroupResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}


func (to *Session) GetCacheAssignmentGroupByID(id int) ([]tc.CacheAssignmentGroup, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_V14_CacheAssignmentGroups, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CacheAssignmentGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}


func (to *Session) GetCacheAssignmentGroupByName(name string) ([]tc.CacheAssignmentGroup, ReqInf, error) {
	url := fmt.Sprintf("%s?name=%s", API_V14_CacheAssignmentGroups, name)
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CacheAssignmentGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}


func (to *Session) DeleteCacheAssignmentGroupByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_V14_CacheAssignmentGroups, id)
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
