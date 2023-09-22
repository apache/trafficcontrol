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
	// API_CACHEGROUPPARAMETERS is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
	API_CACHEGROUPPARAMETERS = apiBase + "/cachegroupparameters"

	APICachegroupParameters = "/cachegroupparameters"
)

func (to *Session) GetCacheGroupParametersWithHdr(cacheGroupID int, header http.Header) ([]tc.CacheGroupParameter, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d/parameters", APICachegroups, cacheGroupID)
	return to.getCacheGroupParameters(route, "", header)
}

// GetCacheGroupParameters Gets a Cache Group's Parameters
// Deprecated: GetCacheGroupParameters will be removed in 6.0. Use GetCacheGroupParametersWithHdr.
func (to *Session) GetCacheGroupParameters(cacheGroupID int) ([]tc.CacheGroupParameter, toclientlib.ReqInf, error) {
	return to.GetCacheGroupParametersWithHdr(cacheGroupID, nil)
}

func (to *Session) GetCacheGroupParametersByQueryParamsWithHdr(cacheGroupID int, queryParams string, header http.Header) ([]tc.CacheGroupParameter, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d/parameters", APICachegroups, cacheGroupID)
	return to.getCacheGroupParameters(route, queryParams, header)
}

// GetCacheGroupParametersByQueryParams Gets a Cache Group's Parameters with query parameters
// Deprecated: GetCacheGroupParametersByQueryParams will be removed in 6.0. Use GetCacheGroupParametersByQueryParamsWithHdr.
func (to *Session) GetCacheGroupParametersByQueryParams(cacheGroupID int, queryParams string) ([]tc.CacheGroupParameter, toclientlib.ReqInf, error) {
	return to.GetCacheGroupParametersByQueryParamsWithHdr(cacheGroupID, queryParams, nil)
}

func (to *Session) getCacheGroupParameters(route, queryParams string, header http.Header) ([]tc.CacheGroupParameter, toclientlib.ReqInf, error) {
	r := fmt.Sprintf("%s%s", route, queryParams)
	var data tc.CacheGroupParametersResponse
	reqInf, err := to.get(r, header, &data)
	return data.Response, reqInf, err
}

func (to *Session) GetAllCacheGroupParametersWithHdr(header http.Header) ([]tc.CacheGroupParametersResponseNullable, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/", APICachegroupParameters)
	var data tc.AllCacheGroupParametersResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response.CacheGroupParameters, reqInf, err
}

// GetAllCacheGroupParameters Gets all Cachegroup Parameter associations
// Deprecated: GetAllCacheGroupParameters will be removed in 6.0. Use GetAllCacheGroupParametersWithHdr.
func (to *Session) GetAllCacheGroupParameters() ([]tc.CacheGroupParametersResponseNullable, toclientlib.ReqInf, error) {
	return to.GetAllCacheGroupParametersWithHdr(nil)
}

// DeleteCacheGroupParameter Deassociates a Parameter with a Cache Group
func (to *Session) DeleteCacheGroupParameter(cacheGroupID, parameterID int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d/%d", APICachegroupParameters, cacheGroupID, parameterID)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

// CreateCacheGroupParameter Associates a Parameter with a Cache Group
func (to *Session) CreateCacheGroupParameter(cacheGroupID, parameterID int) (*tc.CacheGroupParametersPostResponse, toclientlib.ReqInf, error) {
	cacheGroupParameterReq := tc.CacheGroupParameterRequest{
		CacheGroupID: cacheGroupID,
		ParameterID:  parameterID,
	}
	var data tc.CacheGroupParametersPostResponse
	reqInf, err := to.post(APICachegroupParameters, cacheGroupParameterReq, nil, &data)
	return &data, reqInf, err
}
