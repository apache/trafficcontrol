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
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

// These are the API endpoints used by the various Delivery Service-related client methods.
const (
	// API_DELIVERY_SERVICES is the API path on which Traffic Ops serves Delivery Service
	// information. More specific information is typically found on sub-paths of this.
	APIDeliveryServices = "/deliveryservices"

	// APIDeliveryServiceId is the API path on which Traffic Ops serves information about
	// a specific Delivery Service identified by an integral, unique identifier. It is
	// intended to be used with fmt.Sprintf to insert its required path parameter (namely the ID
	// of the Delivery Service of interest).
	APIDeliveryServiceID = APIDeliveryServices + "/%d"

	// APIDeliveryServiceHealth is the API path on which Traffic Ops serves information about
	// the 'health' of a specific Delivery Service identified by an integral, unique identifier. It is
	// intended to be used with fmt.Sprintf to insert its required path parameter (namely the ID
	// of the Delivery Service of interest).
	APIDeliveryServiceHealth = APIDeliveryServiceID + "/health"

	// APIDeliveryServiceCapacity is the API path on which Traffic Ops serves information about
	// the 'capacity' of a specific Delivery Service identified by an integral, unique identifier. It is
	// intended to be used with fmt.Sprintf to insert its required path parameter (namely the ID
	// of the Delivery Service of interest).
	APIDeliveryServiceCapacity = APIDeliveryServiceID + "/capacity"

	// APIDeliveryServiceEligibleServers is the API path on which Traffic Ops serves information about
	// the servers which are eligible to be assigned to a specific Delivery Service identified by an integral,
	// unique identifier. It is intended to be used with fmt.Sprintf to insert its required path parameter
	// (namely the ID of the Delivery Service of interest).
	APIDeliveryServiceEligibleServers = APIDeliveryServiceID + "/servers/eligible"

	// APIDeliveryServicesSafeUpdate is the API path on which Traffic Ops provides the functionality to
	// update the "safe" subset of properties of a Delivery Service identified by an integral, unique
	// identifer. It is intended to be used with fmt.Sprintf to insert its required path parameter
	// (namely the ID of the Delivery Service of interest).
	APIDeliveryServicesSafeUpdate = APIDeliveryServiceID + "/safe"

	// APIDeliveryServiceXMLIDSSLKeys is the API path on which Traffic Ops serves information about
	// and functionality relating to the SSL keys used by a Delivery Service identified by its XMLID. It is
	// intended to be used with fmt.Sprintf to insert its required path parameter (namely the XMLID
	// of the Delivery Service of interest).
	APIDeliveryServiceXMLIDSSLKeys = APIDeliveryServices + "/xmlId/%s/sslkeys"

	// APIDeliveryServiceGenerateSSLKeys is the API path on which Traffic Ops will generate new SSL keys
	APIDeliveryServiceGenerateSSLKeys = APIDeliveryServices + "/sslkeys/generate"

	// APIDeliveryServiceURISigningKeys is the API path on which Traffic Ops serves information
	// about and functionality relating to the URI-signing keys used by a Delivery Service identified
	// by its XMLID. It is intended to be used with fmt.Sprintf to insert its required path parameter
	// (namely the XMLID of the Delivery Service of interest).
	APIDeliveryServicesURISigningKeys = APIDeliveryServices + "/%s/urisignkeys"

	// APIDeliveryServicesURLSigKeys is the API path on which Traffic Ops serves information
	// about and functionality relating to the URL-signing keys used by a Delivery Service identified
	// by its XMLID. It is intended to be used with fmt.Sprintf to insert its required path parameter
	// (namely the XMLID of the Delivery Service of interest).
	APIDeliveryServicesURLSigKeys = APIDeliveryServices + "/xmlId/%s/urlkeys"

	// APIDeliveryServicesRegexes is the API path on which Traffic Ops serves Delivery Service
	// 'regex' (Regular Expression) information.
	APIDeliveryServicesRegexes = "/deliveryservices_regexes"

	// APIServerDeliveryServices is the API path on which Traffic Ops serves functionality
	// related to the associations a specific server and its assigned Delivery Services. It is
	// intended to be used with fmt.Sprintf to insert its required path parameter (namely the ID
	// of the server of interest).
	APIServerDeliveryServices = "/servers/%d/deliveryservices"

	// APIDeliveryServiceServer is the API path on which Traffic Ops serves functionality related
	// to the associations between Delivery Services and their assigned Server(s).
	APIDeliveryServiceServer = "/deliveryserviceserver"

	// APIDeliveryServicesServers is the API path on which Traffic Ops serves functionality related
	// to the associations between a Delivery Service and its assigned Server(s).
	APIDeliveryServicesServers = "/deliveryservices/%s/servers"
)

// GetDeliveryServicesByServer retrieves all Delivery Services assigned to the
// server with the given ID.
func (to *Session) GetDeliveryServicesByServer(id int, header http.Header) ([]tc.DeliveryServiceV4, toclientlib.ReqInf, error) {
	var data tc.DeliveryServicesResponseV4
	reqInf, err := to.get(fmt.Sprintf(APIServerDeliveryServices, id), header, &data)
	return data.Response, reqInf, err
}

// GetDeliveryServices returns all (tenant-visible) Delivery Services that
// satisfy the passed query string parameters. See the API documentation for
// information on the available parameters.
func (to *Session) GetDeliveryServices(header http.Header, params url.Values) ([]tc.DeliveryServiceV4, toclientlib.ReqInf, error) {
	uri := APIDeliveryServices
	if params != nil {
		uri += "?" + params.Encode()
	}
	var data tc.DeliveryServicesResponseV4
	reqInf, err := to.get(uri, header, &data)
	return data.Response, reqInf, err
}

// GetDeliveryServicesByCDNID retrieves all Delivery Services in the CDN with
// the given ID.
func (to *Session) GetDeliveryServicesByCDNID(cdnID int, header http.Header) ([]tc.DeliveryServiceV4, toclientlib.ReqInf, error) {
	var data tc.DeliveryServicesResponseV4
	reqInf, err := to.get(APIDeliveryServices+"?cdn="+strconv.Itoa(cdnID), header, &data)
	return data.Response, reqInf, err
}

// GetDeliveryServiceByID fetches the Delivery Service with the given ID.
func (to *Session) GetDeliveryServiceByID(id string, header http.Header) (*tc.DeliveryServiceV4, toclientlib.ReqInf, error) {
	var data tc.DeliveryServicesResponseV4
	route := fmt.Sprintf("%s?id=%s", APIDeliveryServices, id)
	reqInf, err := to.get(route, header, &data)
	if err != nil {
		return nil, reqInf, err
	}
	if len(data.Response) == 0 {
		return nil, reqInf, nil
	}
	return &data.Response[0], reqInf, nil
}

// GetDeliveryServiceByXMLID fetches all Delivery Services with the given
// XMLID.
func (to *Session) GetDeliveryServiceByXMLID(XMLID string, header http.Header) ([]tc.DeliveryServiceV4, toclientlib.ReqInf, error) {
	var data tc.DeliveryServicesResponseV4
	reqInf, err := to.get(APIDeliveryServices+"?xmlId="+url.QueryEscape(XMLID), header, &data)
	return data.Response, reqInf, err
}

// CreateDeliveryService creates the Delivery Service it's passed.
func (to *Session) CreateDeliveryService(ds tc.DeliveryServiceV4) (tc.DeliveryServiceV4, toclientlib.ReqInf, error) {
	var reqInf toclientlib.ReqInf
	if ds.TypeID == nil && ds.Type != nil {
		ty, _, err := to.GetTypeByName(ds.Type.String(), nil)
		if err != nil {
			return tc.DeliveryServiceV4{}, reqInf, err
		}
		if len(ty) == 0 {
			return tc.DeliveryServiceV4{}, reqInf, fmt.Errorf("no type named %s", ds.Type)
		}
		ds.TypeID = &ty[0].ID
	}

	if ds.CDNID == nil && ds.CDNName != nil {
		cdns, _, err := to.GetCDNByName(*ds.CDNName, nil)
		if err != nil {
			return tc.DeliveryServiceV4{}, reqInf, err
		}
		if len(cdns) == 0 {
			return tc.DeliveryServiceV4{}, reqInf, errors.New("no CDN named " + *ds.CDNName)
		}
		ds.CDNID = &cdns[0].ID
	}

	if ds.ProfileID == nil && ds.ProfileName != nil {
		profiles, _, err := to.GetProfileByName(*ds.ProfileName, nil)
		if err != nil {
			return tc.DeliveryServiceV4{}, reqInf, err
		}
		if len(profiles) == 0 {
			return tc.DeliveryServiceV4{}, reqInf, errors.New("no Profile named " + *ds.ProfileName)
		}
		ds.ProfileID = &profiles[0].ID
	}

	if ds.TenantID == nil && ds.Tenant != nil {
		ten, _, err := to.GetTenantByName(*ds.Tenant, nil)
		if err != nil {
			return tc.DeliveryServiceV4{}, reqInf, err
		}
		ds.TenantID = &ten.ID
	}

	var data tc.DeliveryServicesResponseV4
	reqInf, err := to.post(APIDeliveryServices, ds, nil, &data)
	if err != nil {
		return tc.DeliveryServiceV4{}, reqInf, err
	}
	if len(data.Response) != 1 {
		return tc.DeliveryServiceV4{}, reqInf, fmt.Errorf("failed to create Delivery Service, response indicated %d were created", len(data.Response))
	}

	return data.Response[0], reqInf, nil
}

// UpdateDeliveryService replaces the Delivery Service identified by the
// integral, unique identifier 'id' with the one it's passed.
func (to *Session) UpdateDeliveryService(id int, ds tc.DeliveryServiceV4, header http.Header) (tc.DeliveryServiceV4, toclientlib.ReqInf, error) {
	var data tc.DeliveryServicesResponseV4
	reqInf, err := to.put(fmt.Sprintf(APIDeliveryServiceID, id), ds, header, &data)
	if err != nil {
		return tc.DeliveryServiceV4{}, reqInf, err
	}
	if len(data.Response) != 1 {
		return tc.DeliveryServiceV4{}, reqInf, fmt.Errorf("failed to update Delivery Service #%d; response indicated that %d were updated", id, len(data.Response))
	}
	return data.Response[0], reqInf, nil
}

// DeleteDeliveryService deletes the DeliveryService matching the ID it's passed.
func (to *Session) DeleteDeliveryService(id int) (*tc.DeleteDeliveryServiceResponse, error) {
	var data tc.DeleteDeliveryServiceResponse
	_, err := to.del(fmt.Sprintf(APIDeliveryServiceID, id), nil, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// GetDeliveryServiceHealth gets the 'health' of the Delivery Service identified by the
// integral, unique identifier 'id'.
func (to *Session) GetDeliveryServiceHealth(id int, header http.Header) (*tc.DeliveryServiceHealth, toclientlib.ReqInf, error) {
	var data tc.DeliveryServiceHealthResponse
	reqInf, err := to.get(fmt.Sprintf(APIDeliveryServiceHealth, id), nil, &data)
	if err != nil {
		return nil, reqInf, err
	}

	return &data.Response, reqInf, nil
}

// GetDeliveryServiceCapacity gets the 'capacity' of the Delivery Service identified by the
// integral, unique identifier 'id'.
func (to *Session) GetDeliveryServiceCapacity(id int, header http.Header) (*tc.DeliveryServiceCapacity, toclientlib.ReqInf, error) {
	var data tc.DeliveryServiceCapacityResponse
	reqInf, err := to.get(fmt.Sprintf(APIDeliveryServiceCapacity, id), header, &data)
	if err != nil {
		return nil, reqInf, err
	}
	return &data.Response, reqInf, nil
}

// GenerateSSLKeysForDS generates ssl keys for a given cdn
func (to *Session) GenerateSSLKeysForDS(XMLID string, CDNName string, sslFields tc.SSLKeyRequestFields) (string, toclientlib.ReqInf, error) {
	version := util.JSONIntStr(1)
	request := tc.DeliveryServiceSSLKeysReq{
		BusinessUnit:    sslFields.BusinessUnit,
		CDN:             util.StrPtr(CDNName),
		City:            sslFields.City,
		Country:         sslFields.Country,
		DeliveryService: util.StrPtr(XMLID),
		HostName:        sslFields.HostName,
		Key:             util.StrPtr(XMLID),
		Organization:    sslFields.Organization,
		State:           sslFields.State,
		Version:         &version,
	}
	response := struct {
		Response string `json:"response"`
	}{}
	reqInf, err := to.post(APIDeliveryServiceGenerateSSLKeys, request, nil, &response)
	if err != nil {
		return "", reqInf, err
	}
	return response.Response, reqInf, nil
}

// DeleteDeliveryServiceSSLKeys deletes the SSL Keys used by the Delivery
// Service identified by the passed XMLID.
func (to *Session) DeleteDeliveryServiceSSLKeys(XMLID string) (string, toclientlib.ReqInf, error) {
	resp := struct {
		Response string `json:"response"`
	}{}
	reqInf, err := to.del(fmt.Sprintf(APIDeliveryServiceXMLIDSSLKeys, url.QueryEscape(XMLID)), nil, &resp)
	return resp.Response, reqInf, err
}

// GetDeliveryServiceSSLKeys retrieves the SSL keys of the Delivery Service
// with the given XMLID.
func (to *Session) GetDeliveryServiceSSLKeys(XMLID string, header http.Header) (*tc.DeliveryServiceSSLKeys, toclientlib.ReqInf, error) {
	var data tc.DeliveryServiceSSLKeysResponse
	reqInf, err := to.get(fmt.Sprintf(APIDeliveryServiceXMLIDSSLKeys, url.QueryEscape(XMLID)), header, &data)
	if err != nil {
		return nil, reqInf, err
	}
	return &data.Response, reqInf, nil
}

// GetDeliveryServicesEligible returns the servers eligible for assignment to the Delivery
// Service identified by the integral, unique identifier 'dsID'.
func (to *Session) GetDeliveryServicesEligible(dsID int, header http.Header) ([]tc.DSServer, toclientlib.ReqInf, error) {
	resp := struct {
		Response []tc.DSServer `json:"response"`
	}{Response: []tc.DSServer{}}

	reqInf, err := to.get(fmt.Sprintf(APIDeliveryServiceEligibleServers, dsID), header, &resp)
	return resp.Response, reqInf, err
}

// GetDeliveryServiceURLSigKeys returns the URL-signing keys used by the Delivery Service
// identified by the XMLID 'dsName'.
func (to *Session) GetDeliveryServiceURLSigKeys(dsName string, header http.Header) (tc.URLSigKeys, toclientlib.ReqInf, error) {
	data := struct {
		Response tc.URLSigKeys `json:"response"`
	}{}

	reqInf, err := to.get(fmt.Sprintf(APIDeliveryServicesURLSigKeys, dsName), header, &data)
	if err != nil {
		return tc.URLSigKeys{}, reqInf, err
	}
	return data.Response, reqInf, nil
}

// GetDeliveryServiceURISigningKeys returns the URI-signing keys used by the Delivery Service
// identified by the XMLID 'dsName'. The result is not parsed.
func (to *Session) GetDeliveryServiceURISigningKeys(dsName string, header http.Header) ([]byte, toclientlib.ReqInf, error) {
	data := json.RawMessage{}
	reqInf, err := to.get(fmt.Sprintf(APIDeliveryServicesURISigningKeys, url.QueryEscape(dsName)), header, &data)
	if err != nil {
		return []byte{}, reqInf, err
	}
	return []byte(data), reqInf, nil
}

// SafeDeliveryServiceUpdate updates the "safe" fields of the Delivery
// Service identified by the integral, unique identifier 'id'.
func (to *Session) SafeDeliveryServiceUpdate(id int, r tc.DeliveryServiceSafeUpdateRequest, header http.Header) (tc.DeliveryServiceV4, toclientlib.ReqInf, error) {
	var data tc.DeliveryServiceSafeUpdateResponseV4
	reqInf, err := to.put(fmt.Sprintf(APIDeliveryServicesSafeUpdate, id), r, header, &data)
	if err != nil {
		return tc.DeliveryServiceV4{}, reqInf, err
	}
	if len(data.Response) != 1 {
		return tc.DeliveryServiceV4{}, reqInf, fmt.Errorf("failed to safe update Delivery Service #%d; response indicated that %d were updated", id, len(data.Response))
	}
	return data.Response[0], reqInf, nil
}
