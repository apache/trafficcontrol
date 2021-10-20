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
	"net/http"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

const (
	API_CACHEGROUPPARAMETERS = apiBase + "/cachegroupparameters"
)

// GetCacheGroupParameters Gets a Cache Group's Parameters
func (to *Session) GetCacheGroupParameters(cacheGroupID int) ([]tc.CacheGroupParameter, ReqInf, error) {
	route := fmt.Sprintf("%s/%d/parameters", API_CACHEGROUPS, cacheGroupID)
	return to.getCacheGroupParameters(route, "")
}

// GetCacheGroupParametersByQueryParams Gets a Cache Group's Parameters with query parameters
func (to *Session) GetCacheGroupParametersByQueryParams(cacheGroupID int, queryParams string) ([]tc.CacheGroupParameter, ReqInf, error) {
	route := fmt.Sprintf("%s/%d/parameters", API_CACHEGROUPS, cacheGroupID)
	return to.getCacheGroupParameters(route, queryParams)
}

func (to *Session) getCacheGroupParameters(route, queryParams string) ([]tc.CacheGroupParameter, ReqInf, error) {
	r := fmt.Sprintf("%s%s", route, queryParams)
	resp, remoteAddr, err := to.request(http.MethodGet, r, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CacheGroupParametersResponse
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}
	return data.Response, reqInf, nil
}

// GetAllCacheGroupParameters Gets all Cachegroup Parameter associations
func (to *Session) GetAllCacheGroupParameters() ([]tc.CacheGroupParametersResponseNullable, ReqInf, error) {
	route := fmt.Sprintf("%s/", API_CACHEGROUPPARAMETERS)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.AllCacheGroupParametersResponse
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}
	return data.Response.CacheGroupParameters, reqInf, nil
}

// DeleteCacheGroupParameter Deassociates a Parameter with a Cache Group
func (to *Session) DeleteCacheGroupParameter(cacheGroupID, parameterID int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d/%d", API_CACHEGROUPPARAMETERS, cacheGroupID, parameterID)
	resp, remoteAddr, err := to.request(http.MethodDelete, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()

	var alerts tc.Alerts
	if err = json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return tc.Alerts{}, reqInf, err
	}
	return alerts, reqInf, nil
}

// CreateCacheGroupParameter Associates a Parameter with a Cache Group
func (to *Session) CreateCacheGroupParameter(cacheGroupID, parameterID int) (*tc.CacheGroupParametersPostResponse, ReqInf, error) {
	cacheGroupParameterReq := tc.CacheGroupParameterRequest{
		CacheGroupID: cacheGroupID,
		ParameterID:  parameterID,
	}
	reqBody, err := json.Marshal(cacheGroupParameterReq)
	if err != nil {
		return nil, ReqInf{}, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_CACHEGROUPPARAMETERS, reqBody)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CacheGroupParametersPostResponse
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}
	return &data, reqInf, nil
}
