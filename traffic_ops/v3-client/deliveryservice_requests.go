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
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_DS_REQUESTS = apiBase + "/deliveryservice_requests"
)

// CreateDeliveryServiceRequest creates a Delivery Service Request.
//
// Deprecated: Please use versioned client methods from now on - in this case, CreateDeliveryServiceRequestV30
func (to *Session) CreateDeliveryServiceRequest(dsr tc.DeliveryServiceRequest) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	var remoteAddr net.Addr
	if dsr.AssigneeID == 0 && dsr.Assignee != "" {
		res, reqInf, err := to.GetUserByUsername(dsr.Assignee)
		if err != nil {
			return alerts, reqInf, err
		}
		if len(res) == 0 {
			return alerts, reqInf, errors.New("no user with name " + dsr.Assignee)
		}
		dsr.AssigneeID = *res[0].ID
	}

	if dsr.AuthorID == 0 && dsr.Author != "" {
		res, reqInf, err := to.GetUserByUsername(dsr.Author)
		if err != nil {
			return alerts, reqInf, err
		}
		if len(res) == 0 {
			return alerts, reqInf, errors.New("no user with name " + dsr.Author)
		}
		dsr.AuthorID = tc.IDNoMod(*res[0].ID)
	}

	if dsr.DeliveryService.TypeID == 0 && dsr.DeliveryService.Type.String() != "" {
		ty, reqInf, err := to.GetTypeByName(dsr.DeliveryService.Type.String())
		if err != nil || len(ty) == 0 {
			return alerts, reqInf, errors.New("no type named " + dsr.DeliveryService.Type.String())
		}
		dsr.DeliveryService.TypeID = ty[0].ID
	}

	if dsr.DeliveryService.CDNID == 0 && dsr.DeliveryService.CDNName != "" {
		cdns, reqInf, err := to.GetCDNByName(dsr.DeliveryService.CDNName)
		if err != nil || len(cdns) == 0 {
			return alerts, reqInf, errors.New("no CDN named " + dsr.DeliveryService.CDNName)
		}
		dsr.DeliveryService.CDNID = cdns[0].ID
	}

	if dsr.DeliveryService.ProfileID == 0 && dsr.DeliveryService.ProfileName != "" {
		profiles, reqInf, err := to.GetProfileByName(dsr.DeliveryService.ProfileName)
		if err != nil || len(profiles) == 0 {
			return alerts, reqInf, errors.New("no Profile named " + dsr.DeliveryService.ProfileName)
		}
		dsr.DeliveryService.ProfileID = profiles[0].ID
	}

	if dsr.DeliveryService.TenantID == 0 && dsr.DeliveryService.Tenant != "" {
		ten, reqInf, err := to.TenantByName(dsr.DeliveryService.Tenant)
		if err != nil || ten == nil {
			return alerts, reqInf, errors.New("no Tenant named " + dsr.DeliveryService.Tenant)
		}
		dsr.DeliveryService.TenantID = ten.ID
	}

	reqBody, err := json.Marshal(dsr)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return alerts, reqInf, err
	}
	resp, remoteAddr, err := to.RawRequest(http.MethodPost, API_DS_REQUESTS, reqBody)
	defer resp.Body.Close()

	if err == nil {
		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			return alerts, reqInf, readErr
		}
		if err = json.Unmarshal(body, &alerts); err != nil {
			return alerts, reqInf, err
		}
	}

	return alerts, reqInf, err
}

// CreateDeliveryServiceRequestV30 creates a Delivery Service Request.
func (to *Session) CreateDeliveryServiceRequestV30(dsr tc.DeliveryServiceRequestV30, header http.Header) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	var remoteAddr net.Addr
	if dsr.AssigneeID == nil && dsr.Assignee != nil {
		res, reqInf, err := to.GetUserByUsername(*dsr.Assignee)
		if err != nil {
			return alerts, reqInf, err
		}
		if len(res) == 0 {
			return alerts, reqInf, fmt.Errorf("no user with name '%s'", *dsr.Assignee)
		}
		dsr.AssigneeID = res[0].ID
	}

	if dsr.AuthorID == nil && dsr.Author != "" {
		res, reqInf, err := to.GetUserByUsername(dsr.Author)
		if err != nil {
			return alerts, reqInf, err
		}
		if len(res) == 0 {
			return alerts, reqInf, fmt.Errorf("no user with name '%s'", dsr.Author)
		}
		dsr.AuthorID = res[0].ID
	}

	if dsr.ChangeType == tc.DSRChangeTypeDelete && dsr.Original != nil {
		if dsr.Original.TypeID == nil && dsr.Original.Type != nil {
			ty, reqInf, err := to.GetTypeByName(dsr.Original.Type.String())
			if err != nil || len(ty) == 0 {
				return alerts, reqInf, fmt.Errorf("no type named '%s'", dsr.Original.Type)
			}
			dsr.Original.TypeID = &ty[0].ID
		}

		if dsr.Original.CDNID == nil && dsr.Original.CDNName != nil {
			cdns, reqInf, err := to.GetCDNByName(*dsr.Original.CDNName)
			if err != nil || len(cdns) == 0 {
				return alerts, reqInf, fmt.Errorf("no CDN named '%s'", *dsr.Original.CDNName)
			}
			dsr.Original.CDNID = &cdns[0].ID
		}

		if dsr.Original.ProfileID == nil && dsr.Original.ProfileName != nil {
			profiles, reqInf, err := to.GetProfileByName(*dsr.Original.ProfileName)
			if err != nil || len(profiles) == 0 {
				return alerts, reqInf, fmt.Errorf("no Profile named '%s'", *dsr.Original.ProfileName)
			}
			dsr.Original.ProfileID = &profiles[0].ID
		}

		if dsr.Original.TenantID == nil && dsr.Original.Tenant != nil {
			ten, reqInf, err := to.TenantByName(*dsr.Original.Tenant)
			if err != nil || ten == nil {
				return alerts, reqInf, fmt.Errorf("no Tenant named '%s'", *dsr.Original.Tenant)
			}
			dsr.Original.TenantID = &ten.ID
		}
	}

	reqBody, err := json.Marshal(dsr)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return alerts, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_DS_REQUESTS, reqBody, header)
	defer resp.Body.Close()

	if err == nil {
		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			return alerts, reqInf, readErr
		}
		if err = json.Unmarshal(body, &alerts); err != nil {
			return alerts, reqInf, err
		}
	}

	return alerts, reqInf, err
}

// GetDeliveryServiceRequestsV30 retrieves DSRs based on the given HTTP header
// and query string parameters.
func (to *Session) GetDeliveryServiceRequestsV30(header http.Header, params url.Values) ([]tc.DeliveryServiceRequestV30, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_DS_REQUESTS, nil, header)

	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.DeliveryServiceRequestV30{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	data := struct {
		Response []tc.DeliveryServiceRequestV30 `json:"response"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetDeliveryServiceRequests retrieves all deliveryservices available to session user.
//
// Deprecated: GetDeliveryServiceRequests will be removed in 6.0. Use GetDeliveryServiceRequestsV30.
func (to *Session) GetDeliveryServiceRequests() ([]tc.DeliveryServiceRequest, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_DS_REQUESTS, nil, nil)

	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.DeliveryServiceRequest{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	data := struct {
		Response []tc.DeliveryServiceRequest `json:"response"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GET a DeliveryServiceRequest by the DeliveryServiceRequest XMLID
//
// Deprecated: GetDeliveryServiceRequestByXMLID will be removed in 6.0. Use GetDeliveryServiceRequestsV30.
func (to *Session) GetDeliveryServiceRequestByXMLID(XMLID string) ([]tc.DeliveryServiceRequest, ReqInf, error) {
	route := fmt.Sprintf("%s?xmlId=%s", API_DS_REQUESTS, url.QueryEscape(XMLID))
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil, nil)

	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.DeliveryServiceRequest{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	data := struct {
		Response []tc.DeliveryServiceRequest `json:"response"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GET a DeliveryServiceRequest by the DeliveryServiceRequest id
// Deprecated: GetDeliveryServiceRequestByID will be removed in 6.0. Use GetDeliveryServiceRequestByIDWithHdr.
func (to *Session) GetDeliveryServiceRequestByID(id int) ([]tc.DeliveryServiceRequest, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_DS_REQUESTS, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.DeliveryServiceRequest{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	data := struct {
		Response []tc.DeliveryServiceRequest `json:"response"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// Update a DeliveryServiceRequest by ID
func (to *Session) UpdateDeliveryServiceRequestByID(id int, dsr tc.DeliveryServiceRequest) (tc.Alerts, ReqInf, error) {

	route := fmt.Sprintf("%s?id=%d", API_DS_REQUESTS, id)

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(dsr)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody, nil)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// DELETE a DeliveryServiceRequest by DeliveryServiceRequest assignee
func (to *Session) DeleteDeliveryServiceRequestByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_DS_REQUESTS, id)
	resp, remoteAddr, err := to.RawRequest(http.MethodDelete, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}
