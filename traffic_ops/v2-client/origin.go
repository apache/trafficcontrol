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

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

const (
	API_ORIGINS = apiBase + "/origins"
)

func originIDs(to *Session, origin *tc.Origin) error {
	if origin.CachegroupID == nil && origin.Cachegroup != nil {
		p, _, err := to.GetCacheGroupNullableByName(*origin.Cachegroup)
		if err != nil {
			return err
		}
		if len(p) == 0 {
			return errors.New("no cachegroup named " + *origin.Cachegroup)
		}
		origin.CachegroupID = p[0].ID
	}

	if origin.DeliveryServiceID == nil && origin.DeliveryService != nil {
		dses, _, err := to.GetDeliveryServiceByXMLIDNullable(*origin.DeliveryService)
		if err != nil {
			return err
		}
		if len(dses) == 0 {
			return errors.New("no deliveryservice with name " + *origin.DeliveryService)
		}
		origin.DeliveryServiceID = dses[0].ID
	}

	if origin.ProfileID == nil && origin.Profile != nil {
		profiles, _, err := to.GetProfileByName(*origin.Profile)
		if err != nil {
			return err
		}
		if len(profiles) == 0 {
			return errors.New("no profile with name " + *origin.Profile)
		}
		origin.ProfileID = &profiles[0].ID
	}

	if origin.CoordinateID == nil && origin.Coordinate != nil {
		coordinates, _, err := to.GetCoordinateByName(*origin.Coordinate)
		if err != nil {
			return err
		}
		if len(coordinates) == 0 {
			return errors.New("no coordinate with name " + *origin.Coordinate)
		}
		origin.CoordinateID = &coordinates[0].ID
	}

	if origin.TenantID == nil && origin.Tenant != nil {
		tenant, _, err := to.TenantByName(*origin.Tenant)
		if err != nil {
			return err
		}
		origin.TenantID = &tenant.ID
	}

	return nil
}

// Create an Origin
func (to *Session) CreateOrigin(origin tc.Origin) (*tc.OriginDetailResponse, ReqInf, error) {
	var remoteAddr net.Addr
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}

	err := originIDs(to, &origin)
	if err != nil {
		return nil, reqInf, err
	}

	reqBody, err := json.Marshal(origin)
	if err != nil {
		return nil, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_ORIGINS, reqBody)
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()
	var originResp tc.OriginDetailResponse
	if err = json.NewDecoder(resp.Body).Decode(&originResp); err != nil {
		return nil, reqInf, err
	}
	return &originResp, reqInf, nil
}

// Update an Origin by ID
func (to *Session) UpdateOriginByID(id int, origin tc.Origin) (*tc.OriginDetailResponse, ReqInf, error) {
	var remoteAddr net.Addr
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}

	err := originIDs(to, &origin)
	if err != nil {
		return nil, reqInf, err
	}

	reqBody, err := json.Marshal(origin)
	if err != nil {
		return nil, reqInf, err
	}
	route := fmt.Sprintf("%s?id=%d", API_ORIGINS, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody)
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()
	var originResp tc.OriginDetailResponse
	if err = json.NewDecoder(resp.Body).Decode(&originResp); err != nil {
		return nil, reqInf, err
	}
	return &originResp, reqInf, nil
}

// GET a list of Origins by a query parameter string
func (to *Session) GetOriginsByQueryParams(queryParams string) ([]tc.Origin, ReqInf, error) {
	URI := API_ORIGINS + queryParams
	resp, remoteAddr, err := to.request(http.MethodGet, URI, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.OriginsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// Returns a list of Origins
func (to *Session) GetOrigins() ([]tc.Origin, ReqInf, error) {
	return to.GetOriginsByQueryParams("")
}

// GET an Origin by the Origin ID
func (to *Session) GetOriginByID(id int) ([]tc.Origin, ReqInf, error) {
	return to.GetOriginsByQueryParams(fmt.Sprintf("?id=%d", id))
}

// GET an Origin by the Origin name
func (to *Session) GetOriginByName(name string) ([]tc.Origin, ReqInf, error) {
	return to.GetOriginsByQueryParams(fmt.Sprintf("?name=%s", url.QueryEscape(name)))
}

// GET a list of Origins by Delivery Service ID
func (to *Session) GetOriginsByDeliveryServiceID(id int) ([]tc.Origin, ReqInf, error) {
	return to.GetOriginsByQueryParams(fmt.Sprintf("?deliveryservice=%d", id))
}

// DELETE an Origin by ID
func (to *Session) DeleteOriginByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_ORIGINS, id)
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
