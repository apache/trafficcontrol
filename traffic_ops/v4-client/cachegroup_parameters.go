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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

// apiCachegroupParameters is the API version-relative path for the /cachegroupparameters API endpoint.
const apiCachegroupParameters = "/cachegroupparameters"

// GetCacheGroupParameters gets all Parameters for the identified Cache Group.
func (to *Session) GetCacheGroupParameters(cacheGroupID int, opts RequestOptions) (tc.CacheGroupParametersResponse, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d/parameters", apiCachegroups, cacheGroupID)
	var data tc.CacheGroupParametersResponse
	reqInf, err := to.get(route, opts, &data)
	return data, reqInf, err
}

// GetAllCacheGroupParameters Gets all Cachegroup-to-Parameter associations.
func (to *Session) GetAllCacheGroupParameters(opts RequestOptions) (tc.AllCacheGroupParametersResponse, toclientlib.ReqInf, error) {
	var data tc.AllCacheGroupParametersResponse
	reqInf, err := to.get(apiCachegroupParameters, opts, &data)
	return data, reqInf, err
}

// DeleteCacheGroupParameter de-associates a Parameter with a Cache Group.
func (to *Session) DeleteCacheGroupParameter(cacheGroupID, parameterID int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d/%d", apiCachegroupParameters, cacheGroupID, parameterID)
	var alerts tc.Alerts
	reqInf, err := to.del(route, opts, &alerts)
	return alerts, reqInf, err
}

// CreateCacheGroupParameter associates a Parameter with a Cache Group.
func (to *Session) CreateCacheGroupParameter(cacheGroupID, parameterID int, opts RequestOptions) (tc.CacheGroupParametersPostResponse, toclientlib.ReqInf, error) {
	cacheGroupParameterReq := tc.CacheGroupParameterRequest{
		CacheGroupID: cacheGroupID,
		ParameterID:  parameterID,
	}
	var data tc.CacheGroupParametersPostResponse
	reqInf, err := to.post(apiCachegroupParameters, opts, cacheGroupParameterReq, &data)
	return data, reqInf, err
}

// CreateMultipleCacheGroupParameter associates multiple parameter with multiple Cache Group.
func (to *Session) CreateMultipleCacheGroupParameter(pps []tc.CacheGroupParameterRequest, opts RequestOptions) (tc.CacheGroupParametersPostResponse, toclientlib.ReqInf, error) {
	var data tc.CacheGroupParametersPostResponse
	reqInf, err := to.post(apiCachegroupParameters, opts, pps, &data)
	return data, reqInf, err
}
