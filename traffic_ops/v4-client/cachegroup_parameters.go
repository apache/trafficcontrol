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

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_CACHEGROUPPARAMETERS = apiBase + "/cachegroupparameters"
)

func (to *Session) GetCacheGroupParametersWithHdr(cacheGroupID int, header http.Header) ([]tc.CacheGroupParameter, ReqInf, error) {
	route := fmt.Sprintf("%s/%d/parameters", API_CACHEGROUPS, cacheGroupID)
	return to.getCacheGroupParameters(route, "", header)
}

// GetCacheGroupParameters Gets a Cache Group's Parameters
// Deprecated: GetCacheGroupParameters will be removed in 6.0. Use GetCacheGroupParametersWithHdr.
func (to *Session) GetCacheGroupParameters(cacheGroupID int) ([]tc.CacheGroupParameter, ReqInf, error) {
	return to.GetCacheGroupParametersWithHdr(cacheGroupID, nil)
}

func (to *Session) GetCacheGroupParametersByQueryParamsWithHdr(cacheGroupID int, queryParams string, header http.Header) ([]tc.CacheGroupParameter, ReqInf, error) {
	route := fmt.Sprintf("%s/%d/parameters", API_CACHEGROUPS, cacheGroupID)
	return to.getCacheGroupParameters(route, queryParams, header)
}

// GetCacheGroupParametersByQueryParams Gets a Cache Group's Parameters with query parameters
// Deprecated: GetCacheGroupParametersByQueryParams will be removed in 6.0. Use GetCacheGroupParametersByQueryParamsWithHdr.
func (to *Session) GetCacheGroupParametersByQueryParams(cacheGroupID int, queryParams string) ([]tc.CacheGroupParameter, ReqInf, error) {
	return to.GetCacheGroupParametersByQueryParamsWithHdr(cacheGroupID, queryParams, nil)
}

func (to *Session) getCacheGroupParameters(route, queryParams string, header http.Header) ([]tc.CacheGroupParameter, ReqInf, error) {
	r := fmt.Sprintf("%s%s", route, queryParams)
	var data tc.CacheGroupParametersResponse
	reqInf, err := to.get(r, header, &data)
	return data.Response, reqInf, err
}

func (to *Session) GetAllCacheGroupParametersWithHdr(header http.Header) ([]tc.CacheGroupParametersResponseNullable, ReqInf, error) {
	route := fmt.Sprintf("%s/", API_CACHEGROUPPARAMETERS)
	var data tc.AllCacheGroupParametersResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response.CacheGroupParameters, reqInf, err
}

// GetAllCacheGroupParameters Gets all Cachegroup Parameter associations
// Deprecated: GetAllCacheGroupParameters will be removed in 6.0. Use GetAllCacheGroupParametersWithHdr.
func (to *Session) GetAllCacheGroupParameters() ([]tc.CacheGroupParametersResponseNullable, ReqInf, error) {
	return to.GetAllCacheGroupParametersWithHdr(nil)
}

// DeleteCacheGroupParameter Deassociates a Parameter with a Cache Group
func (to *Session) DeleteCacheGroupParameter(cacheGroupID, parameterID int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d/%d", API_CACHEGROUPPARAMETERS, cacheGroupID, parameterID)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

// CreateCacheGroupParameter Associates a Parameter with a Cache Group
func (to *Session) CreateCacheGroupParameter(cacheGroupID, parameterID int) (*tc.CacheGroupParametersPostResponse, ReqInf, error) {
	cacheGroupParameterReq := tc.CacheGroupParameterRequest{
		CacheGroupID: cacheGroupID,
		ParameterID:  parameterID,
	}
	var data tc.CacheGroupParametersPostResponse
	reqInf, err := to.post(API_CACHEGROUPPARAMETERS, cacheGroupParameterReq, nil, &data)
	return &data, reqInf, err
}
