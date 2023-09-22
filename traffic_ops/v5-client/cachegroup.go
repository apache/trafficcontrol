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
	"errors"
	"fmt"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiCachegroups is the API version-relative path for the /cachegroups API endpoint.
const apiCachegroups = "/cachegroups"

// CreateCacheGroup creates the given Cache Group.
func (to *Session) CreateCacheGroup(cachegroup tc.CacheGroupNullableV5, opts RequestOptions) (tc.CacheGroupDetailResponseV5, toclientlib.ReqInf, error) {
	var resp tc.CacheGroupDetailResponseV5
	if cachegroup.TypeID == nil && cachegroup.Type != nil {
		opts := NewRequestOptions()
		opts.QueryParameters.Set("name", *cachegroup.Type)
		ty, _, err := to.GetTypes(opts)
		if err != nil {
			return resp, toclientlib.ReqInf{}, fmt.Errorf("resolving Type name '%s' to an ID: %w - alerts: %+v", *cachegroup.Name, err, ty.Alerts)
		}
		if len(ty.Response) == 0 {
			return resp, toclientlib.ReqInf{}, errors.New("no type named " + *cachegroup.Type)
		}
		cachegroup.TypeID = &ty.Response[0].ID
	}

	if cachegroup.ParentCachegroupID == nil && cachegroup.ParentName != nil {
		opts := NewRequestOptions()
		opts.QueryParameters.Set("name", *cachegroup.ParentName)
		p, _, err := to.GetCacheGroups(opts)
		if err != nil {
			resp.Alerts = p.Alerts
			return resp, toclientlib.ReqInf{}, err
		}
		if len(p.Response) == 0 {
			resp.Alerts = p.Alerts
			return resp, toclientlib.ReqInf{}, errors.New("no cachegroup named " + *cachegroup.ParentName)
		}
		cachegroup.ParentCachegroupID = p.Response[0].ID
	}

	if cachegroup.SecondaryParentCachegroupID == nil && cachegroup.SecondaryParentName != nil {
		opts := NewRequestOptions()
		opts.QueryParameters.Set("name", *cachegroup.SecondaryParentName)
		p, _, err := to.GetCacheGroups(opts)
		if err != nil {
			resp.Alerts = p.Alerts
			return resp, toclientlib.ReqInf{}, err
		}
		if len(p.Response) == 0 {
			resp.Alerts = p.Alerts
			return resp, toclientlib.ReqInf{}, errors.New("no cachegroup named " + *cachegroup.ParentName)
		}
		cachegroup.SecondaryParentCachegroupID = p.Response[0].ID
	}

	reqInf, err := to.post(apiCachegroups, opts, cachegroup, &resp)
	return resp, reqInf, err
}

// UpdateCacheGroup replaces the Cache Group identified by the given ID with
// the given Cache Group.
func (to *Session) UpdateCacheGroup(id int, cachegroup tc.CacheGroupNullableV5, opts RequestOptions) (tc.CacheGroupDetailResponseV5, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", apiCachegroups, id)
	var cachegroupResp tc.CacheGroupDetailResponseV5
	reqInf, err := to.put(route, opts, cachegroup, &cachegroupResp)
	return cachegroupResp, reqInf, err
}

// GetCacheGroups retrieves Cache Groups configured in Traffic Ops.
func (to *Session) GetCacheGroups(opts RequestOptions) (tc.CacheGroupsNullableResponseV5, toclientlib.ReqInf, error) {
	var data tc.CacheGroupsNullableResponseV5
	reqInf, err := to.get(apiCachegroups, opts, &data)
	return data, reqInf, err
}

// DeleteCacheGroup deletes the Cache Group with the given ID.
func (to *Session) DeleteCacheGroup(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", apiCachegroups, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, opts, &alerts)
	return alerts, reqInf, err
}

// SetCacheGroupDeliveryServices assigns all of the assignable Cache Servers in
// the identified Cache Group to all of the identified the Delivery Services.
func (to *Session) SetCacheGroupDeliveryServices(cgID int, dsIDs []int, opts RequestOptions) (tc.CacheGroupPostDSRespResponse, toclientlib.ReqInf, error) {
	uri := fmt.Sprintf("%s/%d/deliveryservices", apiCachegroups, cgID)
	req := tc.CachegroupPostDSReq{DeliveryServices: dsIDs}
	var resp tc.CacheGroupPostDSRespResponse
	reqInf, err := to.post(uri, opts, req, &resp)
	return resp, reqInf, err
}
