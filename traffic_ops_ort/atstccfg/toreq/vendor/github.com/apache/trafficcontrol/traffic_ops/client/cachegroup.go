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
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_v13_CacheGroups = "/api/1.3/cachegroups"
)

// Create a CacheGroup
func (to *Session) CreateCacheGroupNullable(cachegroup tc.CacheGroupNullable) (*tc.CacheGroupDetailResponse, ReqInf, error) {
	if cachegroup.TypeID == nil && cachegroup.Type != nil {
		ty, _, err := to.GetTypeByName(*cachegroup.Type)
		if err != nil {
			return nil, ReqInf{}, err
		}
		if len(ty) == 0 {
			return nil, ReqInf{}, errors.New("no type named " + *cachegroup.Type)
		}
		cachegroup.TypeID = &ty[0].ID
	}

	if cachegroup.ParentCachegroupID == nil && cachegroup.ParentName != nil {
		p, _, err := to.GetCacheGroupByName(*cachegroup.ParentName)
		if err != nil {
			return nil, ReqInf{}, err
		}
		if len(p) == 0 {
			return nil, ReqInf{}, errors.New("no cachegroup named " + *cachegroup.ParentName)
		}
		cachegroup.ParentCachegroupID = &p[0].ID
	}

	if cachegroup.SecondaryParentCachegroupID == nil && cachegroup.SecondaryParentName != nil {
		p, _, err := to.GetCacheGroupByName(*cachegroup.SecondaryParentName)
		if err != nil {
			return nil, ReqInf{}, err
		}
		if len(p) == 0 {
			return nil, ReqInf{}, errors.New("no cachegroup named " + *cachegroup.ParentName)
		}
		cachegroup.SecondaryParentCachegroupID = &p[0].ID
	}

	reqBody, err := json.Marshal(cachegroup)
	if err != nil {
		return nil, ReqInf{}, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_v13_CacheGroups, reqBody)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()
	var cachegroupResp tc.CacheGroupDetailResponse
	if err = json.NewDecoder(resp.Body).Decode(&cachegroupResp); err != nil {
		return nil, reqInf, err
	}
	return &cachegroupResp, reqInf, nil
}

// Create a CacheGroup
// Deprecated: Use CreateCacheGroupNullable
func (to *Session) CreateCacheGroup(cachegroup tc.CacheGroup) (tc.Alerts, ReqInf, error) {
	if cachegroup.TypeID == 0 && cachegroup.Type != "" {
		ty, _, err := to.GetTypeByName(cachegroup.Type)
		if err != nil {
			return tc.Alerts{}, ReqInf{}, err
		}
		if len(ty) == 0 {
			return tc.Alerts{}, ReqInf{}, errors.New("no type named " + cachegroup.Type)
		}
		cachegroup.TypeID = ty[0].ID
	}

	if cachegroup.ParentCachegroupID == 0 && cachegroup.ParentName != "" {
		p, _, err := to.GetCacheGroupByName(cachegroup.ParentName)
		if err != nil {
			return tc.Alerts{}, ReqInf{}, err
		}
		if len(p) == 0 {
			return tc.Alerts{}, ReqInf{}, errors.New("no cachegroup named " + cachegroup.ParentName)
		}
		cachegroup.ParentCachegroupID = p[0].ID
	}

	if cachegroup.SecondaryParentCachegroupID == 0 && cachegroup.SecondaryParentName != "" {
		p, _, err := to.GetCacheGroupByName(cachegroup.SecondaryParentName)
		if err != nil {
			return tc.Alerts{}, ReqInf{}, err
		}
		if len(p) == 0 {
			return tc.Alerts{}, ReqInf{}, errors.New("no cachegroup named " + cachegroup.ParentName)
		}
		cachegroup.SecondaryParentCachegroupID = p[0].ID
	}

	reqBody, err := json.Marshal(cachegroup)
	if err != nil {
		return tc.Alerts{}, ReqInf{}, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_v13_CacheGroups, reqBody)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Update a CacheGroup by ID
func (to *Session) UpdateCacheGroupNullableByID(id int, cachegroup tc.CacheGroupNullable) (*tc.CacheGroupDetailResponse, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(cachegroup)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	route := fmt.Sprintf("%s/%d", API_v13_CacheGroups, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody)
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()
	var cachegroupResp tc.CacheGroupDetailResponse
	if err = json.NewDecoder(resp.Body).Decode(&cachegroupResp); err != nil {
		return nil, reqInf, err
	}
	return &cachegroupResp, reqInf, nil
}

// Update a CacheGroup by ID
// Deprecated: use UpdateCachegroupNullableByID
func (to *Session) UpdateCacheGroupByID(id int, cachegroup tc.CacheGroup) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(cachegroup)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := fmt.Sprintf("%s/%d", API_v13_CacheGroups, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Returns a list of CacheGroups
func (to *Session) GetCacheGroupsNullable() ([]tc.CacheGroupNullable, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_v13_CacheGroups, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CacheGroupsNullableResponse
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}
	return data.Response, reqInf, nil
}

// Returns a list of CacheGroups
// Deprecated: use GetCacheGroupsNullable
func (to *Session) GetCacheGroups() ([]tc.CacheGroup, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_v13_CacheGroups, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CacheGroupsResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// CacheGroups gets the CacheGroups in an array of CacheGroup structs
// (note CacheGroup used to be called location)
// Deprecated: use GetCacheGroups.
func (to *Session) CacheGroups() ([]tc.CacheGroup, error) {
	cgs, _, err := to.GetCacheGroups()
	return cgs, err
}

// GET a CacheGroup by the CacheGroup id
func (to *Session) GetCacheGroupNullableByID(id int) ([]tc.CacheGroupNullable, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_v13_CacheGroups, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CacheGroupsNullableResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GET a CacheGroup by the CacheGroup id
// Deprecated: use GetCacheGroupNullableByID
func (to *Session) GetCacheGroupByID(id int) ([]tc.CacheGroup, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_v13_CacheGroups, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CacheGroupsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GET a CacheGroup by the CacheGroup name
func (to *Session) GetCacheGroupNullableByName(name string) ([]tc.CacheGroupNullable, ReqInf, error) {
	url := fmt.Sprintf("%s?name=%s", API_v13_CacheGroups, url.QueryEscape(name))
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CacheGroupsNullableResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GET a CacheGroup by the CacheGroup name
// Deprecated: use GetCachegroupNullableByName
func (to *Session) GetCacheGroupByName(name string) ([]tc.CacheGroup, ReqInf, error) {
	url := fmt.Sprintf("%s?name=%s", API_v13_CacheGroups, url.QueryEscape(name))
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.CacheGroupsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// DELETE a CacheGroup by ID
func (to *Session) DeleteCacheGroupByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_v13_CacheGroups, id)
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
