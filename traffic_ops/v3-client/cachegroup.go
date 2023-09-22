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
	"strconv"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

const (
	// DEPRECATED: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
	API_CACHEGROUPS = apiBase + "/cachegroups"

	APICachegroups = "/cachegroups"
)

// Create a CacheGroup.
func (to *Session) CreateCacheGroupNullable(cachegroup tc.CacheGroupNullable) (*tc.CacheGroupDetailResponse, toclientlib.ReqInf, error) {
	if cachegroup.TypeID == nil && cachegroup.Type != nil {
		ty, _, err := to.GetTypeByNameWithHdr(*cachegroup.Type, nil)
		if err != nil {
			return nil, toclientlib.ReqInf{}, err
		}
		if len(ty) == 0 {
			return nil, toclientlib.ReqInf{}, errors.New("no type named " + *cachegroup.Type)
		}
		cachegroup.TypeID = &ty[0].ID
	}

	if cachegroup.ParentCachegroupID == nil && cachegroup.ParentName != nil {
		p, _, err := to.GetCacheGroupNullableByNameWithHdr(*cachegroup.ParentName, nil)
		if err != nil {
			return nil, toclientlib.ReqInf{}, err
		}
		if len(p) == 0 {
			return nil, toclientlib.ReqInf{}, errors.New("no cachegroup named " + *cachegroup.ParentName)
		}
		cachegroup.ParentCachegroupID = p[0].ID
	}

	if cachegroup.SecondaryParentCachegroupID == nil && cachegroup.SecondaryParentName != nil {
		p, _, err := to.GetCacheGroupNullableByNameWithHdr(*cachegroup.SecondaryParentName, nil)
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

func (to *Session) UpdateCacheGroupNullableByIDWithHdr(id int, cachegroup tc.CacheGroupNullable, h http.Header) (*tc.CacheGroupDetailResponse, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APICachegroups, id)
	var cachegroupResp tc.CacheGroupDetailResponse
	reqInf, err := to.put(route, cachegroup, h, &cachegroupResp)
	return &cachegroupResp, reqInf, err
}

// Update a CacheGroup by ID.
// Deprecated: UpdateCacheGroupNullableByID will be removed in 6.0. Use UpdateCacheGroupNullableByIDWithHdr.
func (to *Session) UpdateCacheGroupNullableByID(id int, cachegroup tc.CacheGroupNullable) (*tc.CacheGroupDetailResponse, toclientlib.ReqInf, error) {
	return to.UpdateCacheGroupNullableByIDWithHdr(id, cachegroup, nil)
}

func (to *Session) GetCacheGroupsNullableWithHdr(header http.Header) ([]tc.CacheGroupNullable, toclientlib.ReqInf, error) {
	var data tc.CacheGroupsNullableResponse
	reqInf, err := to.get(APICachegroups, header, &data)
	return data.Response, reqInf, err
}

// Returns a list of CacheGroups.
// Deprecated: GetCacheGroupsNullable will be removed in 6.0. Use GetCacheGroupsNullableWithHdr.
func (to *Session) GetCacheGroupsNullable() ([]tc.CacheGroupNullable, toclientlib.ReqInf, error) {
	return to.GetCacheGroupsNullableWithHdr(nil)
}

func (to *Session) GetCacheGroupNullableByIDWithHdr(id int, header http.Header) ([]tc.CacheGroupNullable, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%v", APICachegroups, id)
	var data tc.CacheGroupsNullableResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GET a CacheGroup by the CacheGroup ID.
// Deprecated: GetCacheGroupNullableByID will be removed in 6.0. Use GetCacheGroupNullableByIDWithHdr.
func (to *Session) GetCacheGroupNullableByID(id int) ([]tc.CacheGroupNullable, toclientlib.ReqInf, error) {
	return to.GetCacheGroupNullableByIDWithHdr(id, nil)
}

func (to *Session) GetCacheGroupNullableByNameWithHdr(name string, header http.Header) ([]tc.CacheGroupNullable, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", APICachegroups, url.QueryEscape(name))
	var data tc.CacheGroupsNullableResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GET a CacheGroup by the CacheGroup name.
// Deprecated: GetCacheGroupNullableByName will be removed in 6.0. Use GetCacheGroupNullableByNameWithHdr.
func (to *Session) GetCacheGroupNullableByName(name string) ([]tc.CacheGroupNullable, toclientlib.ReqInf, error) {
	return to.GetCacheGroupNullableByNameWithHdr(name, nil)
}

func (to *Session) GetCacheGroupNullableByShortNameWithHdr(shortName string, header http.Header) ([]tc.CacheGroupNullable, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?shortName=%s", APICachegroups, url.QueryEscape(shortName))
	var data tc.CacheGroupsNullableResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GET a CacheGroup by the CacheGroup short name.
// Deprecated: GetCacheGroupNullableByShortName will be removed in 6.0. Use GetCacheGroupNullableByShortNameWithHdr.
func (to *Session) GetCacheGroupNullableByShortName(shortName string) ([]tc.CacheGroupNullable, toclientlib.ReqInf, error) {
	return to.GetCacheGroupNullableByShortNameWithHdr(shortName, nil)
}

// DELETE a CacheGroup by ID.
func (to *Session) DeleteCacheGroupByID(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APICachegroups, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

// GetCacheGroupsByQueryParams gets cache groups by the given query parameters.
// Deprecated: GetCacheGroupsByQueryParams will be removed in 6.0. Use GetCacheGroupsByQueryParamsWithHdr.
func (to *Session) GetCacheGroupsByQueryParams(qparams url.Values) ([]tc.CacheGroupNullable, toclientlib.ReqInf, error) {
	return to.GetCacheGroupsByQueryParamsWithHdr(qparams, nil)
}

func (to *Session) GetCacheGroupsByQueryParamsWithHdr(qparams url.Values, header http.Header) ([]tc.CacheGroupNullable, toclientlib.ReqInf, error) {
	route := APICachegroups
	if len(qparams) > 0 {
		route += "?" + qparams.Encode()
	}
	var data tc.CacheGroupsNullableResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

func (to *Session) SetCachegroupDeliveryServices(cgID int, dsIDs []int) (tc.CacheGroupPostDSRespResponse, toclientlib.ReqInf, error) {
	uri := `/cachegroups/` + strconv.Itoa(cgID) + `/deliveryservices`
	req := tc.CachegroupPostDSReq{DeliveryServices: dsIDs}
	resp := tc.CacheGroupPostDSRespResponse{}
	reqInf, err := to.post(uri, req, nil, &resp)
	return resp, reqInf, err
}
