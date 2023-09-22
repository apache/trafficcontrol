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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

const (
	// API_DS_REQUESTS is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
	API_DS_REQUESTS = apiBase + "/deliveryservice_requests"

	APIDSRequests = "/deliveryservice_requests"
)

// CreateDeliveryServiceRequest creates a Delivery Service Request.
func (to *Session) CreateDeliveryServiceRequest(dsr tc.DeliveryServiceRequest) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	if dsr.AssigneeID == 0 && dsr.Assignee != "" {
		res, reqInf, err := to.GetUserByUsernameWithHdr(dsr.Assignee, nil)
		if err != nil {
			return alerts, reqInf, err
		}
		if len(res) == 0 {
			return alerts, reqInf, errors.New("no user with name " + dsr.Assignee)
		}
		dsr.AssigneeID = *res[0].ID
	}

	if dsr.AuthorID == 0 && dsr.Author != "" {
		res, reqInf, err := to.GetUserByUsernameWithHdr(dsr.Author, nil)
		if err != nil {
			return alerts, reqInf, err
		}
		if len(res) == 0 {
			return alerts, reqInf, errors.New("no user with name " + dsr.Author)
		}
		dsr.AuthorID = tc.IDNoMod(*res[0].ID)
	}

	if dsr.DeliveryService.TypeID == 0 && dsr.DeliveryService.Type.String() != "" {
		ty, reqInf, err := to.GetTypeByNameWithHdr(dsr.DeliveryService.Type.String(), nil)
		if err != nil || len(ty) == 0 {
			return alerts, reqInf, errors.New("no type named " + dsr.DeliveryService.Type.String())
		}
		dsr.DeliveryService.TypeID = ty[0].ID
	}

	if dsr.DeliveryService.CDNID == 0 && dsr.DeliveryService.CDNName != "" {
		cdns, reqInf, err := to.GetCDNByNameWithHdr(dsr.DeliveryService.CDNName, nil)
		if err != nil || len(cdns) == 0 {
			return alerts, reqInf, errors.New("no CDN named " + dsr.DeliveryService.CDNName)
		}
		dsr.DeliveryService.CDNID = cdns[0].ID
	}

	if dsr.DeliveryService.ProfileID == 0 && dsr.DeliveryService.ProfileName != "" {
		profiles, reqInf, err := to.GetProfileByNameWithHdr(dsr.DeliveryService.ProfileName, nil)
		if err != nil || len(profiles) == 0 {
			return alerts, reqInf, errors.New("no Profile named " + dsr.DeliveryService.ProfileName)
		}
		dsr.DeliveryService.ProfileID = profiles[0].ID
	}

	if dsr.DeliveryService.TenantID == 0 && dsr.DeliveryService.Tenant != "" {
		ten, reqInf, err := to.TenantByNameWithHdr(dsr.DeliveryService.Tenant, nil)
		if err != nil || ten == nil {
			return alerts, reqInf, errors.New("no Tenant named " + dsr.DeliveryService.Tenant)
		}
		dsr.DeliveryService.TenantID = ten.ID
	}

	reqInf, err := to.post(APIDSRequests, dsr, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) GetDeliveryServiceRequestsWithHdr(header http.Header) ([]tc.DeliveryServiceRequest, toclientlib.ReqInf, error) {
	data := struct {
		Response []tc.DeliveryServiceRequest `json:"response"`
	}{}
	reqInf, err := to.get(APIDSRequests, header, &data)
	return data.Response, reqInf, err
}

// GetDeliveryServiceRequests retrieves all deliveryservices available to session user.
// Deprecated: GetDeliveryServiceRequests will be removed in 6.0. Use GetDeliveryServiceRequestsWithHdr.
func (to *Session) GetDeliveryServiceRequests() ([]tc.DeliveryServiceRequest, toclientlib.ReqInf, error) {
	return to.GetDeliveryServiceRequestsWithHdr(nil)
}

func (to *Session) GetDeliveryServiceRequestByXMLIDWithHdr(XMLID string, header http.Header) ([]tc.DeliveryServiceRequest, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?xmlId=%s", APIDSRequests, url.QueryEscape(XMLID))
	data := struct {
		Response []tc.DeliveryServiceRequest `json:"response"`
	}{}
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GET a DeliveryServiceRequest by the DeliveryServiceRequest XMLID
// Deprecated: GetDeliveryServiceRequestByXMLID will be removed in 6.0. Use GetDeliveryServiceRequestByXMLIDWithHdr.
func (to *Session) GetDeliveryServiceRequestByXMLID(XMLID string) ([]tc.DeliveryServiceRequest, toclientlib.ReqInf, error) {
	return to.GetDeliveryServiceRequestByXMLIDWithHdr(XMLID, nil)
}

func (to *Session) GetDeliveryServiceRequestByIDWithHdr(id int, header http.Header) ([]tc.DeliveryServiceRequest, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIDSRequests, id)
	data := struct {
		Response []tc.DeliveryServiceRequest `json:"response"`
	}{}
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GET a DeliveryServiceRequest by the DeliveryServiceRequest id
// Deprecated: GetDeliveryServiceRequestByID will be removed in 6.0. Use GetDeliveryServiceRequestByIDWithHdr.
func (to *Session) GetDeliveryServiceRequestByID(id int) ([]tc.DeliveryServiceRequest, toclientlib.ReqInf, error) {
	return to.GetDeliveryServiceRequestByIDWithHdr(id, nil)
}

func (to *Session) UpdateDeliveryServiceRequestByIDWithHdr(id int, dsr tc.DeliveryServiceRequest, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIDSRequests, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, dsr, header, &alerts)
	return alerts, reqInf, err
}

// Update a DeliveryServiceRequest by ID
// Deprecated: UpdateDeliveryServiceRequestByID will be removed in 6.0. Use UpdateDeliveryServiceRequestByIDWithHdr.
func (to *Session) UpdateDeliveryServiceRequestByID(id int, dsr tc.DeliveryServiceRequest) (tc.Alerts, toclientlib.ReqInf, error) {
	return to.UpdateDeliveryServiceRequestByIDWithHdr(id, dsr, nil)
}

// DELETE a DeliveryServiceRequest by DeliveryServiceRequest assignee
func (to *Session) DeleteDeliveryServiceRequestByID(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIDSRequests, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
