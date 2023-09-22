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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

const (
	// API_ORIGINS is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
	API_ORIGINS = apiBase + "/origins"

	APIOrigins = "/origins"
)

func originIDs(to *Session, origin *tc.Origin) error {
	if origin.CachegroupID == nil && origin.Cachegroup != nil {
		p, _, err := to.GetCacheGroupNullableByNameWithHdr(*origin.Cachegroup, nil)
		if err != nil {
			return err
		}
		if len(p) == 0 {
			return errors.New("no cachegroup named " + *origin.Cachegroup)
		}
		origin.CachegroupID = p[0].ID
	}

	if origin.DeliveryServiceID == nil && origin.DeliveryService != nil {
		dses, _, err := to.GetDeliveryServiceByXMLIDNullableWithHdr(*origin.DeliveryService, nil)
		if err != nil {
			return err
		}
		if len(dses) == 0 {
			return errors.New("no deliveryservice with name " + *origin.DeliveryService)
		}
		origin.DeliveryServiceID = dses[0].ID
	}

	if origin.ProfileID == nil && origin.Profile != nil {
		profiles, _, err := to.GetProfileByNameWithHdr(*origin.Profile, nil)
		if err != nil {
			return err
		}
		if len(profiles) == 0 {
			return errors.New("no profile with name " + *origin.Profile)
		}
		origin.ProfileID = &profiles[0].ID
	}

	if origin.CoordinateID == nil && origin.Coordinate != nil {
		coordinates, _, err := to.GetCoordinateByNameWithHdr(*origin.Coordinate, nil)
		if err != nil {
			return err
		}
		if len(coordinates) == 0 {
			return errors.New("no coordinate with name " + *origin.Coordinate)
		}
		origin.CoordinateID = &coordinates[0].ID
	}

	if origin.TenantID == nil && origin.Tenant != nil {
		tenant, _, err := to.TenantByNameWithHdr(*origin.Tenant, nil)
		if err != nil {
			return err
		}
		origin.TenantID = &tenant.ID
	}

	return nil
}

// Create an Origin
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

func (to *Session) UpdateOriginByIDWithHdr(id int, origin tc.Origin, header http.Header) (*tc.OriginDetailResponse, toclientlib.ReqInf, error) {
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

// Update an Origin by ID
// Deprecated: UpdateOriginByID will be removed in 6.0. Use UpdateOriginByIDWithHdr.
func (to *Session) UpdateOriginByID(id int, origin tc.Origin) (*tc.OriginDetailResponse, toclientlib.ReqInf, error) {
	return to.UpdateOriginByIDWithHdr(id, origin, nil)
}

// GET a list of Origins by a query parameter string
func (to *Session) GetOriginsByQueryParams(queryParams string) ([]tc.Origin, toclientlib.ReqInf, error) {
	uri := APIOrigins + queryParams
	var data tc.OriginsResponse
	reqInf, err := to.get(uri, nil, &data)
	return data.Response, reqInf, err
}

// Returns a list of Origins
func (to *Session) GetOrigins() ([]tc.Origin, toclientlib.ReqInf, error) {
	return to.GetOriginsByQueryParams("")
}

// GET an Origin by the Origin ID
func (to *Session) GetOriginByID(id int) ([]tc.Origin, toclientlib.ReqInf, error) {
	return to.GetOriginsByQueryParams(fmt.Sprintf("?id=%d", id))
}

// GET an Origin by the Origin name
func (to *Session) GetOriginByName(name string) ([]tc.Origin, toclientlib.ReqInf, error) {
	return to.GetOriginsByQueryParams(fmt.Sprintf("?name=%s", url.QueryEscape(name)))
}

// GET a list of Origins by Delivery Service ID
func (to *Session) GetOriginsByDeliveryServiceID(id int) ([]tc.Origin, toclientlib.ReqInf, error) {
	return to.GetOriginsByQueryParams(fmt.Sprintf("?deliveryservice=%d", id))
}

// DELETE an Origin by ID
func (to *Session) DeleteOriginByID(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIOrigins, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
