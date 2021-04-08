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
	// APICachegroupParameters is the API version-relative path for the /cachegroupparameters API endpoint.
	APICachegroupParameters = "/cachegroupparameters"
)

// GetCacheGroupParameters gets all Parameters for the identified Cache Group.
func (to *Session) GetCacheGroupParameters(cacheGroupID int, queryParameters url.Values, header http.Header) ([]tc.CacheGroupParameter, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d/parameters", APICachegroups, cacheGroupID)
	if len(queryParameters) > 0 {
		route = fmt.Sprintf("%s?%s", route, queryParameters.Encode())
	}
	var data tc.CacheGroupParametersResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetAllCacheGroupParameters Gets all Cachegroup-to-Parameter associations.
func (to *Session) GetAllCacheGroupParameters(header http.Header) ([]tc.CacheGroupParametersResponseNullable, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/", APICachegroupParameters)
	var data tc.AllCacheGroupParametersResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response.CacheGroupParameters, reqInf, err
}

// DeleteCacheGroupParameter de-associates a Parameter with a Cache Group.
func (to *Session) DeleteCacheGroupParameter(cacheGroupID, parameterID int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d/%d", APICachegroupParameters, cacheGroupID, parameterID)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

// CreateCacheGroupParameter associates a Parameter with a Cache Group.
func (to *Session) CreateCacheGroupParameter(cacheGroupID, parameterID int) (*tc.CacheGroupParametersPostResponse, toclientlib.ReqInf, error) {
	cacheGroupParameterReq := tc.CacheGroupParameterRequest{
		CacheGroupID: cacheGroupID,
		ParameterID:  parameterID,
	}
	var data tc.CacheGroupParametersPostResponse
	reqInf, err := to.post(APICachegroupParameters, cacheGroupParameterReq, nil, &data)
	return &data, reqInf, err
}
