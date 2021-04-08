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
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	// APICachegroups is the API version-relative path for the /cachegroups API endpoint.
	APICachegroups = "/cachegroups"
)

// CreateCacheGroup creates the given Cache Group.
func (to *Session) CreateCacheGroup(cachegroup tc.CacheGroupNullable) (*tc.CacheGroupDetailResponse, toclientlib.ReqInf, error) {
	if cachegroup.TypeID == nil && cachegroup.Type != nil {
		ty, _, err := to.GetTypeByName(*cachegroup.Type, nil)
		if err != nil {
			return nil, toclientlib.ReqInf{}, err
		}
		if len(ty) == 0 {
			return nil, toclientlib.ReqInf{}, errors.New("no type named " + *cachegroup.Type)
		}
		cachegroup.TypeID = &ty[0].ID
	}

	if cachegroup.ParentCachegroupID == nil && cachegroup.ParentName != nil {
		p, _, err := to.GetCacheGroupByName(*cachegroup.ParentName, nil)
		if err != nil {
			return nil, toclientlib.ReqInf{}, err
		}
		if len(p) == 0 {
			return nil, toclientlib.ReqInf{}, errors.New("no cachegroup named " + *cachegroup.ParentName)
		}
		cachegroup.ParentCachegroupID = p[0].ID
	}

	if cachegroup.SecondaryParentCachegroupID == nil && cachegroup.SecondaryParentName != nil {
		p, _, err := to.GetCacheGroupByName(*cachegroup.SecondaryParentName, nil)
		if err != nil {
			return nil, toclientlib.ReqInf{}, err
		}
		if len(p) == 0 {
			return nil, toclientlib.ReqInf{}, errors.New("no cachegroup named " + *cachegroup.ParentName)
		}
		cachegroup.SecondaryParentCachegroupID = p[0].ID
	}

	var cachegroupResp tc.CacheGroupDetailResponse
	reqInf, err := to.post(APICachegroups, cachegroup, nil, &cachegroupResp)
	return &cachegroupResp, reqInf, err
}

// UpdateCacheGroup replaces the Cache Group identified by the given ID with
// the given Cache Group.
func (to *Session) UpdateCacheGroup(id int, cachegroup tc.CacheGroupNullable, h http.Header) (*tc.CacheGroupDetailResponse, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APICachegroups, id)
	var cachegroupResp tc.CacheGroupDetailResponse
	reqInf, err := to.put(route, cachegroup, h, &cachegroupResp)
	return &cachegroupResp, reqInf, err
}

// GetCacheGroups retrieves all of the Cache Groups configured in Traffic Ops.
func (to *Session) GetCacheGroups(params url.Values, header http.Header) ([]tc.CacheGroupNullable, toclientlib.ReqInf, error) {
	route := APICachegroups
	if len(params) > 0 {
		route += "?" + params.Encode()
	}

	var data tc.CacheGroupsNullableResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetCacheGroupByID retrieves the Cache Group with the given ID.
func (to *Session) GetCacheGroupByID(id int, header http.Header) ([]tc.CacheGroupNullable, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%v", APICachegroups, id)
	var data tc.CacheGroupsNullableResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetCacheGroupByName retrieves the Cache Group with the given Name.
func (to *Session) GetCacheGroupByName(name string, header http.Header) ([]tc.CacheGroupNullable, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", APICachegroups, url.QueryEscape(name))
	var data tc.CacheGroupsNullableResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetCacheGroupByShortName retrieves the Cache Group with the given Short Name.
func (to *Session) GetCacheGroupByShortName(shortName string, header http.Header) ([]tc.CacheGroupNullable, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?shortName=%s", APICachegroups, url.QueryEscape(shortName))
	var data tc.CacheGroupsNullableResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// DeleteCacheGroupByID deletes the Cache Group with the given ID.
func (to *Session) DeleteCacheGroupByID(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APICachegroups, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

// SetCacheGroupDeliveryServices assigns all of the assignable Cache Servers in
// the identified Cache Group to all of the identified the Delivery Services.
func (to *Session) SetCacheGroupDeliveryServices(cgID int, dsIDs []int) (tc.CacheGroupPostDSRespResponse, toclientlib.ReqInf, error) {
	uri := fmt.Sprintf("%s/%d/deliveryservices", APICachegroups, cgID)
	req := tc.CachegroupPostDSReq{DeliveryServices: dsIDs}
	resp := tc.CacheGroupPostDSRespResponse{}
	reqInf, err := to.post(uri, req, nil, &resp)
	return resp, reqInf, err
}
