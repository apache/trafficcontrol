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
	"net"
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiOrigins is the full path to the /origins API route.
const apiOrigins = "/origins"

func (to *Session) originIDs(origin *tc.Origin) error {
	if origin == nil {
		return errors.New("invalid call to originIDs; nil origin")
	}

	opts := NewRequestOptions()
	if origin.CachegroupID == nil && origin.Cachegroup != nil {
		opts.QueryParameters.Set("name", *origin.Cachegroup)
		p, _, err := to.GetCacheGroups(opts)
		if err != nil {
			return fmt.Errorf("resolving Cache Group name '%s' to an ID: %w - alerts: %+v", *origin.Cachegroup, err, p.Alerts)
		}
		if len(p.Response) == 0 {
			return fmt.Errorf("no Cache Group named '%s'", *origin.Cachegroup)
		}
		opts.QueryParameters.Del("name")
		origin.CachegroupID = p.Response[0].ID
	}

	if origin.DeliveryServiceID == nil && origin.DeliveryService != nil {
		opts.QueryParameters.Set("xmlId", *origin.DeliveryService)
		dses, _, err := to.GetDeliveryServices(opts)
		if err != nil {
			return fmt.Errorf("resolving Delivery Service XMLID '%s' to an ID: %w - alerts: %+v", *origin.DeliveryService, err, dses.Alerts)
		}
		if len(dses.Response) == 0 {
			return fmt.Errorf("no Delivery Service with XMLID '%s'", *origin.DeliveryService)
		}
		opts.QueryParameters.Del("xmlId")
		origin.DeliveryServiceID = dses.Response[0].ID
	}

	if origin.ProfileID == nil && origin.Profile != nil {
		opts.QueryParameters.Set("name", *origin.Profile)
		profiles, _, err := to.GetProfiles(opts)
		if err != nil {
			return fmt.Errorf("resolving Profile name '%s' to an ID: %w - alerts: %+v", *origin.Profile, err, profiles.Alerts)
		}
		if len(profiles.Response) == 0 {
			return errors.New("no profile with name " + *origin.Profile)
		}
		origin.ProfileID = &profiles.Response[0].ID
	}

	if origin.CoordinateID == nil && origin.Coordinate != nil {
		opts.QueryParameters.Set("name", *origin.Coordinate)
		coordinates, _, err := to.GetCoordinates(opts)
		if err != nil {
			return fmt.Errorf("resolving Coordinates name '%s' to an ID: %w - alerts: %+v", *origin.Coordinate, err, coordinates.Alerts)
		}
		if len(coordinates.Response) == 0 {
			return fmt.Errorf("no coordinate with name '%s'", *origin.Coordinate)
		}
		origin.CoordinateID = &coordinates.Response[0].ID
	}

	if origin.TenantID == nil && origin.Tenant != nil {
		opts.QueryParameters.Set("name", *origin.Tenant)
		tenant, _, err := to.GetTenants(opts)
		if err != nil {
			return fmt.Errorf("resolving Tenant name '%s' to an ID: %w - alerts: %+v", *origin.Tenant, err, tenant.Alerts)
		}
		if len(tenant.Response) == 0 {
			return fmt.Errorf("no Tenant with name '%s'", *origin.Tenant)
		}
		origin.TenantID = &tenant.Response[0].ID
	}

	return nil
}

// CreateOrigin creates the given Origin.
func (to *Session) CreateOrigin(origin tc.Origin, opts RequestOptions) (tc.OriginDetailResponse, toclientlib.ReqInf, error) {
	var originResp tc.OriginDetailResponse
	var remoteAddr net.Addr
	reqInf := toclientlib.ReqInf{CacheHitStatus: toclientlib.CacheHitStatusMiss, RemoteAddr: remoteAddr}

	err := to.originIDs(&origin)
	if err != nil {
		return originResp, reqInf, err
	}
	reqInf, err = to.post(apiOrigins, opts, origin, &originResp)
	return originResp, reqInf, err
}

// UpdateOrigin replaces the Origin identified by 'id' with the passed Origin.
func (to *Session) UpdateOrigin(id int, origin tc.Origin, opts RequestOptions) (tc.OriginDetailResponse, toclientlib.ReqInf, error) {
	var originResp tc.OriginDetailResponse
	var remoteAddr net.Addr
	reqInf := toclientlib.ReqInf{CacheHitStatus: toclientlib.CacheHitStatusMiss, RemoteAddr: remoteAddr}

	err := to.originIDs(&origin)
	if err != nil {
		return originResp, reqInf, err
	}
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("id", strconv.Itoa(id))
	reqInf, err = to.put(apiOrigins, opts, origin, &originResp)
	return originResp, reqInf, err
}

// GetOrigins retrieves Origins from Traffic Ops.
func (to *Session) GetOrigins(opts RequestOptions) (tc.OriginsResponse, toclientlib.ReqInf, error) {
	var data tc.OriginsResponse
	reqInf, err := to.get(apiOrigins, opts, &data)
	return data, reqInf, err
}

// DeleteOrigin deletes the Origin with the given ID.
func (to *Session) DeleteOrigin(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("id", strconv.Itoa(id))
	var alerts tc.Alerts
	reqInf, err := to.del(apiOrigins, opts, &alerts)
	return alerts, reqInf, err
}
