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
	"strconv"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

// These are the API endpoints used by the various Delivery Service-related client methods.
const (
	// API_DELIVERY_SERVICES is the API path on which Traffic Ops serves Delivery Service
	// information. More specific information is typically found on sub-paths of this.
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v2/deliveryservices.html
	API_DELIVERY_SERVICES = apiBase + "/deliveryservices"

	// API_DELIVERY_SERVICE_ID is the API path on which Traffic Ops serves information about
	// a specific Delivery Service identified by an integral, unique identifier. It is
	// intended to be used with fmt.Sprintf to insert its required path parameter (namely the ID
	// of the Delivery Service of interest).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v2/deliveryservices_id.html
	API_DELIVERY_SERVICE_ID = API_DELIVERY_SERVICES + "/%v"

	// API_DELIVERY_SERVICE_HEALTH is the API path on which Traffic Ops serves information about
	// the 'health' of a specific Delivery Service identified by an integral, unique identifier. It is
	// intended to be used with fmt.Sprintf to insert its required path parameter (namely the ID
	// of the Delivery Service of interest).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v2/deliveryservices_id_health.html
	API_DELIVERY_SERVICE_HEALTH = API_DELIVERY_SERVICE_ID + "/health"

	// API_DELIVERY_SERVICE_CAPACITY is the API path on which Traffic Ops serves information about
	// the 'capacity' of a specific Delivery Service identified by an integral, unique identifier. It is
	// intended to be used with fmt.Sprintf to insert its required path parameter (namely the ID
	// of the Delivery Service of interest).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v2/deliveryservices_id_capacity.html
	API_DELIVERY_SERVICE_CAPACITY = API_DELIVERY_SERVICE_ID + "/capacity"

	// API_DELIVERY_SERVICE_ELIGIBLE_SERVERS is the API path on which Traffic Ops serves information about
	// the servers which are eligible to be assigned to a specific Delivery Service identified by an integral,
	// unique identifier. It is intended to be used with fmt.Sprintf to insert its required path parameter
	// (namely the ID of the Delivery Service of interest).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v2/deliveryservices_id_servers_eligible.html
	API_DELIVERY_SERVICE_ELIGIBLE_SERVERS = API_DELIVERY_SERVICE_ID + "/servers/eligible"

	// API_DELIVERY_SERVICES_SAFE_UPDATE is the API path on which Traffic Ops provides the functionality to
	// update the "safe" subset of properties of a Delivery Service identified by an integral, unique
	// identifer. It is intended to be used with fmt.Sprintf to insert its required path parameter
	// (namely the ID of the Delivery Service of interest).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v2/deliveryservices_id_safe.html
	API_DELIVERY_SERVICES_SAFE_UPDATE = API_DELIVERY_SERVICE_ID + "/safe"

	// API_DELIVERY_SERVICE_XMLID_SSL_KEYS is the API path on which Traffic Ops serves information about
	// and functionality relating to the SSL keys used by a Delivery Service identified by its XMLID. It is
	// intended to be used with fmt.Sprintf to insert its required path parameter (namely the XMLID
	// of the Delivery Service of interest).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v2/deliveryservices_xmlid_xmlid_sslkeys.html
	API_DELIVERY_SERVICE_XMLID_SSL_KEYS = API_DELIVERY_SERVICES + "/xmlid/%s/sslkeys"

	// API_DELIVERY_SERVICE_URI_SIGNING_KEYS is the API path on which Traffic Ops serves information
	// about and functionality relating to the URI-signing keys used by a Delivery Service identified
	// by its XMLID. It is intended to be used with fmt.Sprintf to insert its required path parameter
	// (namely the XMLID of the Delivery Service of interest).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v2/deliveryservices_xmlid_urisignkeys.html
	API_DELIVERY_SERVICES_URI_SIGNING_KEYS = API_DELIVERY_SERVICES + "/%s/urisignkeys"

	// API_DELIVERY_SERVICES_URL_SIGNING_KEYS is the API path on which Traffic Ops serves information
	// about and functionality relating to the URL-signing keys used by a Delivery Service identified
	// by its XMLID. It is intended to be used with fmt.Sprintf to insert its required path parameter
	// (namely the XMLID of the Delivery Service of interest).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v2/deliveryservices_xmlid_xmlid_urlkeys.html
	API_DELIVERY_SERVICES_URL_SIGNING_KEYS = API_DELIVERY_SERVICES + "/xmlId/%s/urlkeys"

	// API_DELIVERY_SERVICES_REGEXES is the API path on which Traffic Ops serves Delivery Service
	// 'regex' (Regular Expression) information.
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v2/deliveryservices_regexes.html
	API_DELIVERY_SERVICES_REGEXES = apiBase + "/deliveryservices_regexes"

	// API_SERVER_DELIVERY_SERVICES is the API path on which Traffic Ops serves functionality
	// related to the associations a specific server and its assigned Delivery Services. It is
	// intended to be used with fmt.Sprintf to insert its required path parameter (namely the ID
	// of the server of interest).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v2/servers_id_deliveryservices.html
	API_SERVER_DELIVERY_SERVICES = apiBase + "/servers/%d/deliveryservices"

	// API_DELIVERY_SERVICE_SERVER is the API path on which Traffic Ops serves functionality related
	// to the associations between Delivery Services and their assigned Server(s).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v2/deliveryserviceserver.html
	API_DELIVERY_SERVICE_SERVER = apiBase + "/deliveryserviceserver"
)

// GetDeliveryServicesByServer returns all of the (tenant-visible) Delivery Services assigned to
// the server identified by the integral, unique identifier 'id'.
func (to *Session) GetDeliveryServicesByServer(id int) ([]tc.DeliveryServiceNullable, ReqInf, error) {
	var data tc.DeliveryServicesNullableResponse

	reqInf, err := get(to, fmt.Sprintf(API_SERVER_DELIVERY_SERVICES, id), &data)
	if err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetDeliveryServicesNullable returns a slice of Delivery Services.
func (to *Session) GetDeliveryServicesNullable() ([]tc.DeliveryServiceNullable, ReqInf, error) {
	data := struct {
		Response []tc.DeliveryServiceNullable `json:"response"`
	}{}
	reqInf, err := get(to, API_DELIVERY_SERVICES, &data)
	if err != nil {
		return nil, reqInf, err
	}
	return data.Response, reqInf, nil
}

// GetDeliveryServicesByCDNID returns the (tenant-visible) Delivery Services within the CDN identified
// by the integral, unique identifier 'cdnID'.
func (to *Session) GetDeliveryServicesByCDNID(cdnID int) ([]tc.DeliveryServiceNullable, ReqInf, error) {
	data := struct {
		Response []tc.DeliveryServiceNullable `json:"response"`
	}{}
	reqInf, err := get(to, API_DELIVERY_SERVICES+"?cdn="+strconv.Itoa(cdnID), &data)
	if err != nil {
		return nil, reqInf, err
	}
	return data.Response, reqInf, nil
}

// GetDeliveryServiceNullable returns the Delivery Service identified by the integral, unique identifier
// 'id' (which must be passed as a string).
func (to *Session) GetDeliveryServiceNullable(id string) (*tc.DeliveryServiceNullable, ReqInf, error) {
	data := struct {
		Response []tc.DeliveryServiceNullable `json:"response"`
	}{}
	reqInf, err := get(to, API_DELIVERY_SERVICES+"?id="+id, &data)
	if err != nil {
		return nil, reqInf, err
	}
	if len(data.Response) == 0 {
		return nil, reqInf, nil
	}
	return &data.Response[0], reqInf, nil
}

// GetDeliveryServiceByXMLIDNullable returns the Delivery Service identified by the passed XMLID.
// The length of the returned slice should always be 1 when the request is succesful - if it isn't
// something very wicked has happened to Traffic Ops.
func (to *Session) GetDeliveryServiceByXMLIDNullable(XMLID string) ([]tc.DeliveryServiceNullable, ReqInf, error) {
	var data tc.DeliveryServicesNullableResponse
	reqInf, err := get(to, API_DELIVERY_SERVICES+"?xmlId="+XMLID, &data)
	if err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// CreateDeliveryServiceNullable creates the DeliveryService it's passed.
func (to *Session) CreateDeliveryServiceNullable(ds *tc.DeliveryServiceNullable) (*tc.CreateDeliveryServiceNullableResponse, error) {
	if ds.TypeID == nil && ds.Type != nil {
		ty, _, err := to.GetTypeByName(ds.Type.String())
		if err != nil {
			return nil, err
		}
		if len(ty) == 0 {
			return nil, errors.New("no type named " + ds.Type.String())
		}
		ds.TypeID = &ty[0].ID
	}

	if ds.CDNID == nil && ds.CDNName != nil {
		cdns, _, err := to.GetCDNByName(*ds.CDNName)
		if err != nil {
			return nil, err
		}
		if len(cdns) == 0 {
			return nil, errors.New("no CDN named " + *ds.CDNName)
		}
		ds.CDNID = &cdns[0].ID
	}

	if ds.ProfileID == nil && ds.ProfileName != nil {
		profiles, _, err := to.GetProfileByName(*ds.ProfileName)
		if err != nil {
			return nil, err
		}
		if len(profiles) == 0 {
			return nil, errors.New("no Profile named " + *ds.ProfileName)
		}
		ds.ProfileID = &profiles[0].ID
	}

	if ds.TenantID == nil && ds.Tenant != nil {
		ten, _, err := to.TenantByName(*ds.Tenant)
		if err != nil {
			return nil, err
		}
		ds.TenantID = &ten.ID
	}

	var data tc.CreateDeliveryServiceNullableResponse
	jsonReq, err := json.Marshal(ds)
	if err != nil {
		return nil, err
	}
	_, err = post(to, API_DELIVERY_SERVICES, jsonReq, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// UpdateDeliveryServiceNullable updates the DeliveryService matching the ID it's
// passed with the DeliveryService it is passed.
func (to *Session) UpdateDeliveryServiceNullable(id string, ds *tc.DeliveryServiceNullable) (*tc.UpdateDeliveryServiceNullableResponse, error) {
	var data tc.UpdateDeliveryServiceNullableResponse
	jsonReq, err := json.Marshal(ds)
	if err != nil {
		return nil, err
	}
	_, err = put(to, fmt.Sprintf(API_DELIVERY_SERVICE_ID, id), jsonReq, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// DeleteDeliveryService deletes the DeliveryService matching the ID it's passed.
func (to *Session) DeleteDeliveryService(id string) (*tc.DeleteDeliveryServiceResponse, error) {
	var data tc.DeleteDeliveryServiceResponse
	_, err := del(to, fmt.Sprintf(API_DELIVERY_SERVICE_ID, id), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GetDeliveryServiceHealth gets the 'health' of the Delivery Service identified by the
// integral, unique identifier 'id' (which must be passed as a string).
func (to *Session) GetDeliveryServiceHealth(id string) (*tc.DeliveryServiceHealth, ReqInf, error) {
	var data tc.DeliveryServiceHealthResponse
	reqInf, err := get(to, fmt.Sprintf(API_DELIVERY_SERVICE_HEALTH, id), &data)
	if err != nil {
		return nil, reqInf, err
	}

	return &data.Response, reqInf, nil
}

// GetDeliveryServiceCapacity gets the 'capacity' of the Delivery Service identified by the
// integral, unique identifier 'id' (which must be passed as a string).
func (to *Session) GetDeliveryServiceCapacity(id string) (*tc.DeliveryServiceCapacity, ReqInf, error) {
	var data tc.DeliveryServiceCapacityResponse
	reqInf, err := get(to, fmt.Sprintf(API_DELIVERY_SERVICE_CAPACITY, id), &data)
	if err != nil {
		return nil, reqInf, err
	}

	return &data.Response, reqInf, nil
}

// GetDeliveryServiceServer returns associations between Delivery Services and servers using the
// provided pagination controls.
func (to *Session) GetDeliveryServiceServer(page, limit string) ([]tc.DeliveryServiceServer, ReqInf, error) {
	var data tc.DeliveryServiceServerResponse
	reqInf, err := get(to, API_DELIVERY_SERVICE_SERVER+"?page="+page+"&limit="+limit, &data)
	if err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetDeliveryServiceRegexes returns the "Regexes" (Regular Expressions) used by all (tenant-visible)
// Delivery Services.
func (to *Session) GetDeliveryServiceRegexes() ([]tc.DeliveryServiceRegexes, ReqInf, error) {
	var data tc.DeliveryServiceRegexResponse
	reqInf, err := get(to, API_DELIVERY_SERVICES_REGEXES, &data)
	if err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetDeliveryServiceSSLKeysByID returns information about the SSL Keys used by the Delivery
// Service identified by the passed XMLID.
func (to *Session) GetDeliveryServiceSSLKeysByID(XMLID string) (*tc.DeliveryServiceSSLKeys, ReqInf, error) {
	var data tc.DeliveryServiceSSLKeysResponse
	reqInf, err := get(to, fmt.Sprintf(API_DELIVERY_SERVICE_XMLID_SSL_KEYS, XMLID), &data)
	if err != nil {
		return nil, reqInf, err
	}

	return &data.Response, reqInf, nil
}

// GetDeliveryServicesEligible returns the servers eligible for assignment to the Delivery
// Service identified by the integral, unique identifier 'dsID'.
func (to *Session) GetDeliveryServicesEligible(dsID int) ([]tc.DSServerV11, ReqInf, error) {
	resp := struct {
		Response []tc.DSServerV11 `json:"response"`
	}{Response: []tc.DSServerV11{}}

	reqInf, err := get(to, fmt.Sprintf(API_DELIVERY_SERVICE_ELIGIBLE_SERVERS, dsID), &resp)
	if err != nil {
		return nil, reqInf, err
	}
	return resp.Response, reqInf, nil
}

// GetDeliveryServiceURLSigKeys returns the URL-signing keys used by the Delivery Service
// identified by the XMLID 'dsName'.
func (to *Session) GetDeliveryServiceURLSigKeys(dsName string) (tc.URLSigKeys, ReqInf, error) {
	data := struct {
		Response tc.URLSigKeys `json:"response"`
	}{}

	reqInf, err := get(to, fmt.Sprintf(API_DELIVERY_SERVICES_URL_SIGNING_KEYS, dsName), &data)
	if err != nil {
		return tc.URLSigKeys{}, reqInf, err
	}
	return data.Response, reqInf, nil
}

// GetDeliveryServiceURISigningKeys returns the URI-signing keys used by the Delivery Service
// identified by the XMLID 'dsName'. The result is not parsed.
func (to *Session) GetDeliveryServiceURISigningKeys(dsName string) ([]byte, ReqInf, error) {
	data := json.RawMessage{}
	reqInf, err := get(to, fmt.Sprintf(API_DELIVERY_SERVICES_URI_SIGNING_KEYS, dsName), &data)
	if err != nil {
		return []byte{}, reqInf, err
	}
	return []byte(data), reqInf, nil
}

// UpdateDeliveryServiceSafe updates the given Delivery Service identified by 'id' with the given "safe" fields.
func (to *Session) UpdateDeliveryServiceSafe(id int, ds tc.DeliveryServiceSafeUpdateRequest) ([]tc.DeliveryServiceNullable, ReqInf, error) {
	var reqInf ReqInf
	var resp tc.DeliveryServiceSafeUpdateResponse

	req, err := json.Marshal(ds)
	if err != nil {
		return resp.Response, reqInf, err
	}

	if reqInf, err = put(to, fmt.Sprintf(API_DELIVERY_SERVICES_SAFE_UPDATE, id), req, &resp); err != nil {
		return resp.Response, reqInf, err
	}

	if len(resp.Response) < 1 {
		err = errors.New("Traffic Ops returned success, but response was missing the Delivery Service")
	}
	return resp.Response, reqInf, err
}

// GetAccessibleDeliveryServicesByTenant gets all delivery services associated with the given tenant, and all of
// it's children.
func (to *Session) GetAccessibleDeliveryServicesByTenant(tenantId int) ([]tc.DeliveryServiceNullable, ReqInf, error) {
	data := tc.DeliveryServicesNullableResponse{}
	reqInf, err := get(to, fmt.Sprintf("%v?accessibleTo=%v", API_DELIVERY_SERVICES, tenantId), &data)
	return data.Response, reqInf, err
}
