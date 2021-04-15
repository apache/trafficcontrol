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

// APIDSRequests is the API version-relative path to the
// /deliveryservice_requests API endpoint.
const APIDSRequests = "/deliveryservice_requests"

// CreateDeliveryServiceRequest creates the given Delivery Service Request.
func (to *Session) CreateDeliveryServiceRequest(dsr tc.DeliveryServiceRequestV40, hdr http.Header) (tc.DeliveryServiceRequestV40, tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	var created tc.DeliveryServiceRequestV40
	if dsr.AssigneeID == nil && dsr.Assignee != nil {
		res, reqInf, err := to.GetUserByUsername(*dsr.Assignee, nil)
		if err != nil {
			return created, alerts, reqInf, err
		}
		if len(res) == 0 {
			return created, alerts, reqInf, fmt.Errorf("no user with username '%s'", *dsr.Assignee)
		}
		dsr.AssigneeID = res[0].ID
	}

	if dsr.AuthorID == nil && dsr.Author != "" {
		res, reqInf, err := to.GetUserByUsername(dsr.Author, nil)
		if err != nil {
			return created, alerts, reqInf, err
		}
		if len(res) == 0 {
			return created, alerts, reqInf, fmt.Errorf("no user with name '%s'", dsr.Author)
		}
		dsr.AuthorID = res[0].ID
	}

	var ds *tc.DeliveryServiceV4
	if dsr.ChangeType == tc.DSRChangeTypeDelete {
		ds = dsr.Original
	} else {
		ds = dsr.Requested
	}

	if ds.TypeID == nil && ds.Type.String() != "" {
		ty, reqInf, err := to.GetTypeByName(ds.Type.String(), nil)
		if err != nil || len(ty) == 0 {
			return created, alerts, reqInf, errors.New("no type named " + ds.Type.String())
		}
		ds.TypeID = &ty[0].ID
	}

	if ds.CDNID == nil && ds.CDNName != nil {
		cdns, reqInf, err := to.GetCDNByName(*ds.CDNName, nil)
		if err != nil || len(cdns) == 0 {
			return created, alerts, reqInf, fmt.Errorf("no CDN named '%s'", *ds.CDNName)
		}
		ds.CDNID = &cdns[0].ID
	}

	if ds.ProfileID == nil && ds.ProfileName != nil {
		profiles, reqInf, err := to.GetProfileByName(*ds.ProfileName, nil)
		if err != nil || len(profiles) == 0 {
			return created, alerts, reqInf, fmt.Errorf("no Profile named '%s'", *ds.ProfileName)
		}
		ds.ProfileID = &profiles[0].ID
	}

	if ds.TenantID == nil && ds.Tenant != nil {
		ten, reqInf, err := to.GetTenantByName(*ds.Tenant, nil)
		if err != nil {
			return created, alerts, reqInf, fmt.Errorf("no Tenant named '%s'", *ds.Tenant)
		}
		ds.TenantID = &ten.ID
	}

	var resp struct {
		tc.Alerts
		Response tc.DeliveryServiceRequestV40 `json:"response"`
	}
	reqInf, err := to.post(APIDSRequests, dsr, nil, &resp)
	alerts = resp.Alerts
	created = resp.Response
	return created, alerts, reqInf, err
}

// GetDeliveryServiceRequests retrieves all Delivery Service Requests available to session user.
func (to *Session) GetDeliveryServiceRequests(header http.Header) ([]tc.DeliveryServiceRequestV40, tc.Alerts, toclientlib.ReqInf, error) {
	var data struct {
		tc.Alerts
		Response []tc.DeliveryServiceRequestV40 `json:"response"`
	}
	reqInf, err := to.get(APIDSRequests, header, &data)
	return data.Response, data.Alerts, reqInf, err
}

// GetDeliveryServiceRequestsByXMLID retrives all Delivery Service Requests that
// are requests to create, modify, or delete a Delivery Service with the given
// XMLID.
func (to *Session) GetDeliveryServiceRequestsByXMLID(XMLID string, header http.Header) ([]tc.DeliveryServiceRequestV40, tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?xmlId=%s", APIDSRequests, url.QueryEscape(XMLID))
	var data struct {
		tc.Alerts
		Response []tc.DeliveryServiceRequestV40 `json:"response"`
	}
	reqInf, err := to.get(route, header, &data)
	return data.Response, data.Alerts, reqInf, err
}

// GetDeliveryServiceRequest retrieves the Delivery Service Request with the given ID.
func (to *Session) GetDeliveryServiceRequest(id int, header http.Header) (tc.DeliveryServiceRequestV40, tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIDSRequests, id)

	var data struct {
		tc.Alerts
		Response []tc.DeliveryServiceRequestV40 `json:"response"`
	}
	reqInf, err := to.get(route, header, &data)

	// We presume the cases where an incorrect number of DSRs is returned will
	// be captured in the error returned by to.get
	var ret tc.DeliveryServiceRequestV40
	if len(data.Response) == 1 {
		ret = data.Response[0]
	}

	return ret, data.Alerts, reqInf, err
}

// DeleteDeliveryServiceRequest deletes the Delivery Service Request with the given ID.
func (to *Session) DeleteDeliveryServiceRequest(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIDSRequests, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

// UpdateDeliveryServiceRequest replaces the existing DSR that has the given
// ID with the DSR passed.
func (to *Session) UpdateDeliveryServiceRequest(id int, dsr tc.DeliveryServiceRequestV4, header http.Header) (tc.DeliveryServiceRequestV4, tc.Alerts, toclientlib.ReqInf, error) {

	route := fmt.Sprintf("%s?id=%d", APIDSRequests, id)

	var payload struct {
		tc.Alerts
		Response tc.DeliveryServiceRequestV4 `json:"response"`
	}
	reqInf, err := to.put(route, dsr, header, &payload)

	return payload.Response, payload.Alerts, reqInf, err
}
