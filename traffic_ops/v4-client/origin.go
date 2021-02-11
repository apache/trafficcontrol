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
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	// APIOrigins is the full path to the /origins API route.
	APIOrigins = "/origins"
)

func originIDs(to *Session, origin *tc.Origin) error {
	if origin.CachegroupID == nil && origin.Cachegroup != nil {
		p, _, err := to.GetCacheGroupByName(*origin.Cachegroup, nil)
		if err != nil {
			return err
		}
		if len(p) == 0 {
			return errors.New("no cachegroup named " + *origin.Cachegroup)
		}
		origin.CachegroupID = p[0].ID
	}

	if origin.DeliveryServiceID == nil && origin.DeliveryService != nil {
		dses, _, err := to.GetDeliveryServiceByXMLID(*origin.DeliveryService, nil)
		if err != nil {
			return err
		}
		if len(dses) == 0 {
			return errors.New("no deliveryservice with name " + *origin.DeliveryService)
		}
		origin.DeliveryServiceID = dses[0].ID
	}

	if origin.ProfileID == nil && origin.Profile != nil {
		profiles, _, err := to.GetProfileByName(*origin.Profile, nil)
		if err != nil {
			return err
		}
		if len(profiles) == 0 {
			return errors.New("no profile with name " + *origin.Profile)
		}
		origin.ProfileID = &profiles[0].ID
	}

	if origin.CoordinateID == nil && origin.Coordinate != nil {
		coordinates, _, err := to.GetCoordinateByName(*origin.Coordinate, nil)
		if err != nil {
			return err
		}
		if len(coordinates) == 0 {
			return errors.New("no coordinate with name " + *origin.Coordinate)
		}
		origin.CoordinateID = &coordinates[0].ID
	}

	if origin.TenantID == nil && origin.Tenant != nil {
		tenant, _, err := to.GetTenantByName(*origin.Tenant, nil)
		if err != nil {
			return err
		}
		origin.TenantID = &tenant.ID
	}

	return nil
}

// CreateOrigin creates the given Origin.
func (to *Session) CreateOrigin(origin tc.Origin) (*tc.OriginDetailResponse, toclientlib.ReqInf, error) {
	var remoteAddr net.Addr
	reqInf := toclientlib.ReqInf{CacheHitStatus: toclientlib.CacheHitStatusMiss, RemoteAddr: remoteAddr}

	err := originIDs(to, &origin)
	if err != nil {
		return nil, reqInf, err
	}
	var originResp tc.OriginDetailResponse
	reqInf, err = to.post(APIOrigins, origin, nil, &originResp)
	return &originResp, reqInf, err
}

// UpdateOrigin replaces the Origin identified by 'id' with the passed Origin.
func (to *Session) UpdateOrigin(id int, origin tc.Origin, header http.Header) (*tc.OriginDetailResponse, toclientlib.ReqInf, error) {
	var remoteAddr net.Addr
	reqInf := toclientlib.ReqInf{CacheHitStatus: toclientlib.CacheHitStatusMiss, RemoteAddr: remoteAddr}

	err := originIDs(to, &origin)
	if err != nil {
		return nil, reqInf, err
	}
	route := fmt.Sprintf("%s?id=%d", APIOrigins, id)
	var originResp tc.OriginDetailResponse
	reqInf, err = to.put(route, origin, header, &originResp)
	return &originResp, reqInf, err
}

// GetOrigins retrieves Origins from Traffic Ops.
func (to *Session) GetOrigins(queryParams url.Values) ([]tc.Origin, toclientlib.ReqInf, error) {
	URI := APIOrigins
	if len(queryParams) > 0 {
		URI += "?" + queryParams.Encode()
	}
	var data tc.OriginsResponse
	reqInf, err := to.get(URI, nil, &data)
	return data.Response, reqInf, err
}

// GetOriginByID retrieves the Origin with the given ID.
func (to *Session) GetOriginByID(id int) ([]tc.Origin, toclientlib.ReqInf, error) {
	params := url.Values{}
	params.Set("id", strconv.Itoa(id))
	return to.GetOrigins(params)
}

// GetOriginByName retrieves the Origin with the given Name.
func (to *Session) GetOriginByName(name string) ([]tc.Origin, toclientlib.ReqInf, error) {
	params := url.Values{}
	params.Set("name", name)
	return to.GetOrigins(params)
}

// GetOriginsByDeliveryServiceID retrieves all Origins assigned to the Delivery
// Service with the given ID.
func (to *Session) GetOriginsByDeliveryServiceID(id int) ([]tc.Origin, toclientlib.ReqInf, error) {
	params := url.Values{}
	params.Set("deliveryservice", strconv.Itoa(id))
	return to.GetOrigins(params)
}

// DeleteOrigin deletes the Origin with the given ID.
func (to *Session) DeleteOrigin(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIOrigins, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
