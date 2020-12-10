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
)

// These are the API endpoints used by the various Delivery Service-related client methods.
const (
	// API_DELIVERY_SERVICES is the API path on which Traffic Ops serves Delivery Service
	// information. More specific information is typically found on sub-paths of this.
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v3/deliveryservices.html
	API_DELIVERY_SERVICES = apiBase + "/deliveryservices"

	// API_DELIVERY_SERVICE_ID is the API path on which Traffic Ops serves information about
	// a specific Delivery Service identified by an integral, unique identifier. It is
	// intended to be used with fmt.Sprintf to insert its required path parameter (namely the ID
	// of the Delivery Service of interest).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v3/deliveryservices_id.html
	API_DELIVERY_SERVICE_ID = API_DELIVERY_SERVICES + "/%v"

	// API_DELIVERY_SERVICE_HEALTH is the API path on which Traffic Ops serves information about
	// the 'health' of a specific Delivery Service identified by an integral, unique identifier. It is
	// intended to be used with fmt.Sprintf to insert its required path parameter (namely the ID
	// of the Delivery Service of interest).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v3/deliveryservices_id_health.html
	API_DELIVERY_SERVICE_HEALTH = API_DELIVERY_SERVICE_ID + "/health"

	// API_DELIVERY_SERVICE_CAPACITY is the API path on which Traffic Ops serves information about
	// the 'capacity' of a specific Delivery Service identified by an integral, unique identifier. It is
	// intended to be used with fmt.Sprintf to insert its required path parameter (namely the ID
	// of the Delivery Service of interest).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v3/deliveryservices_id_capacity.html
	API_DELIVERY_SERVICE_CAPACITY = API_DELIVERY_SERVICE_ID + "/capacity"

	// API_DELIVERY_SERVICE_ELIGIBLE_SERVERS is the API path on which Traffic Ops serves information about
	// the servers which are eligible to be assigned to a specific Delivery Service identified by an integral,
	// unique identifier. It is intended to be used with fmt.Sprintf to insert its required path parameter
	// (namely the ID of the Delivery Service of interest).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v3/deliveryservices_id_servers_eligible.html
	API_DELIVERY_SERVICE_ELIGIBLE_SERVERS = API_DELIVERY_SERVICE_ID + "/servers/eligible"

	// API_DELIVERY_SERVICES_SAFE_UPDATE is the API path on which Traffic Ops provides the functionality to
	// update the "safe" subset of properties of a Delivery Service identified by an integral, unique
	// identifer. It is intended to be used with fmt.Sprintf to insert its required path parameter
	// (namely the ID of the Delivery Service of interest).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v3/deliveryservices_id_safe.html
	API_DELIVERY_SERVICES_SAFE_UPDATE = API_DELIVERY_SERVICE_ID + "/safe"

	// API_DELIVERY_SERVICE_XMLID_SSL_KEYS is the API path on which Traffic Ops serves information about
	// and functionality relating to the SSL keys used by a Delivery Service identified by its XMLID. It is
	// intended to be used with fmt.Sprintf to insert its required path parameter (namely the XMLID
	// of the Delivery Service of interest).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v3/deliveryservices_xmlid_xmlid_sslkeys.html
	API_DELIVERY_SERVICE_XMLID_SSL_KEYS = API_DELIVERY_SERVICES + "/xmlId/%s/sslkeys"

	// API_DELIVERY_SERVICE_GENERATE_SSL_KEYS is the API path on which Traffic Ops will generate new SSL keys
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v3/deliveryservices_sslkeys_generate.html
	API_DELIVERY_SERVICE_GENERATE_SSL_KEYS = API_DELIVERY_SERVICES + "/sslkeys/generate"

	// API_DELIVERY_SERVICE_URI_SIGNING_KEYS is the API path on which Traffic Ops serves information
	// about and functionality relating to the URI-signing keys used by a Delivery Service identified
	// by its XMLID. It is intended to be used with fmt.Sprintf to insert its required path parameter
	// (namely the XMLID of the Delivery Service of interest).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v3/deliveryservices_xmlid_urisignkeys.html
	API_DELIVERY_SERVICES_URI_SIGNING_KEYS = API_DELIVERY_SERVICES + "/%s/urisignkeys"

	// API_DELIVERY_SERVICES_URL_SIGNING_KEYS is the API path on which Traffic Ops serves information
	// about and functionality relating to the URL-signing keys used by a Delivery Service identified
	// by its XMLID. It is intended to be used with fmt.Sprintf to insert its required path parameter
	// (namely the XMLID of the Delivery Service of interest).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v3/deliveryservices_xmlid_xmlid_urlkeys.html
	API_DELIVERY_SERVICES_URL_SIGNING_KEYS = API_DELIVERY_SERVICES + "/xmlid/%s/urlkeys"

	// API_DELIVERY_SERVICES_REGEXES is the API path on which Traffic Ops serves Delivery Service
	// 'regex' (Regular Expression) information.
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v3/deliveryservices_regexes.html
	API_DELIVERY_SERVICES_REGEXES = apiBase + "/deliveryservices_regexes"

	// API_SERVER_DELIVERY_SERVICES is the API path on which Traffic Ops serves functionality
	// related to the associations a specific server and its assigned Delivery Services. It is
	// intended to be used with fmt.Sprintf to insert its required path parameter (namely the ID
	// of the server of interest).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v3/servers_id_deliveryservices.html
	API_SERVER_DELIVERY_SERVICES = apiBase + "/servers/%d/deliveryservices"

	// API_DELIVERY_SERVICE_SERVER is the API path on which Traffic Ops serves functionality related
	// to the associations between Delivery Services and their assigned Server(s).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v3/deliveryserviceserver.html
	API_DELIVERY_SERVICE_SERVER = apiBase + "/deliveryserviceserver"

	// API_DELIVERY_SERVICES_SERVERS is the API path on which Traffic Ops serves functionality related
	// to the associations between a Delivery Service and its assigned Server(s).
	// See Also: https://traffic-control-cdn.readthedocs.io/en/latest/api/v3/deliveryservices_xmlid_servers.html
	API_DELIVERY_SERVICES_SERVERS = apiBase + "/deliveryservices/%s/servers"
)

func (to *Session) GetDeliveryServicesByServerV30WithHdr(id int, header http.Header) ([]tc.DeliveryServiceNullableV30, ReqInf, error) {
	var data tc.DeliveryServicesResponseV30
	reqInf, err := get(to, fmt.Sprintf(API_SERVER_DELIVERY_SERVICES, id), &data, header)
	return data.Response, reqInf, err
}

// GetDeliveryServicesByServer returns all of the (tenant-visible) Delivery Services assigned to
// the server identified by the integral, unique identifier 'id'.
//
// Warning: This method coerces its returned data into an APIv1.5 format.
//
// Deprecated: Please used versioned library imports in the future, and
// versioned methods, specifically, for API v3.0 - in this case,
// GetDeliveryServicesByServerV30WithHdr.
func (to *Session) GetDeliveryServicesByServer(id int) ([]tc.DeliveryServiceNullable, ReqInf, error) {
	return to.GetDeliveryServicesByServerWithHdr(id, nil)
}

func (to *Session) GetDeliveryServicesByServerWithHdr(id int, header http.Header) ([]tc.DeliveryServiceNullable, ReqInf, error) {
	var data tc.DeliveryServicesNullableResponse

	reqInf, err := get(to, fmt.Sprintf(API_SERVER_DELIVERY_SERVICES, id), &data, header)
	if err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetDeliveryServicesV31 returns all (tenant-visible) Delivery Services that
// satisfy the passed query string parameters. See the API documentation for
// information on the available parameters.
func (to *Session) GetDeliveryServicesV31(header http.Header, params url.Values) ([]tc.DeliveryServiceV31, ReqInf, error) {
	uri := API_DELIVERY_SERVICES
	if params != nil {
		uri += "?" + params.Encode()
	}

	var data tc.DeliveryServicesResponseV31

	reqInf, err := get(to, uri, &data, header)
	return data.Response, reqInf, err
}

// GetDeliveryServicesV30WithHdr returns all (tenant-visible) Delivery Services that
// satisfy the passed query string parameters. See the API documentation for
// information on the available parameters.
func (to *Session) GetDeliveryServicesV30WithHdr(header http.Header, params url.Values) ([]tc.DeliveryServiceNullableV30, ReqInf, error) {
	uri := API_DELIVERY_SERVICES
	if params != nil {
		uri += "?" + params.Encode()
	}

	var data tc.DeliveryServicesResponseV30

	reqInf, err := get(to, uri, &data, header)
	return data.Response, reqInf, err
}

func (to *Session) GetDeliveryServicesNullableWithHdr(header http.Header) ([]tc.DeliveryServiceNullable, ReqInf, error) {
	data := struct {
		Response []tc.DeliveryServiceNullable `json:"response"`
	}{}
	reqInf, err := get(to, API_DELIVERY_SERVICES, &data, header)
	if err != nil {
		return nil, reqInf, err
	}
	return data.Response, reqInf, nil
}

// GetDeliveryServicesNullable returns a slice of Delivery Services.
//
// Warning: This method coerces its returned data into an APIv1.5 format.
//
// Deprecated: Please used versioned library imports in the future, and
// versioned methods, specifically, for API v3.0 - in this case,
// GetDeliveryServicesV30WithHdr.
func (to *Session) GetDeliveryServicesNullable() ([]tc.DeliveryServiceNullable, ReqInf, error) {
	return to.GetDeliveryServicesNullableWithHdr(nil)
}

func (to *Session) GetDeliveryServicesByCDNIDWithHdr(cdnID int, header http.Header) ([]tc.DeliveryServiceNullable, ReqInf, error) {
	data := struct {
		Response []tc.DeliveryServiceNullable `json:"response"`
	}{}
	reqInf, err := get(to, API_DELIVERY_SERVICES+"?cdn="+strconv.Itoa(cdnID), &data, header)
	if err != nil {
		return nil, reqInf, err
	}
	return data.Response, reqInf, nil
}

// GetDeliveryServicesByCDNID returns the (tenant-visible) Delivery Services within the CDN identified
// by the integral, unique identifier 'cdnID'.
//
// Warning: This method coerces its returned data into an APIv1.5 format.
//
// Deprecated: Please used versioned library imports in the future, and
// versioned methods, specifically, for API v3.0 - in this case,
// GetDeliveryServicesV30WithHdr.
func (to *Session) GetDeliveryServicesByCDNID(cdnID int) ([]tc.DeliveryServiceNullable, ReqInf, error) {
	return to.GetDeliveryServicesByCDNIDWithHdr(cdnID, nil)
}

func (to *Session) GetDeliveryServiceNullableWithHdr(id string, header http.Header) (*tc.DeliveryServiceNullable, ReqInf, error) {
	data := struct {
		Response []tc.DeliveryServiceNullable `json:"response"`
	}{}
	reqInf, err := get(to, API_DELIVERY_SERVICES+"?id="+id, &data, header)
	if err != nil {
		return nil, reqInf, err
	}
	if len(data.Response) == 0 {
		return nil, reqInf, nil
	}
	return &data.Response[0], reqInf, nil
}

// GetDeliveryServiceNullable returns the Delivery Service identified by the integral, unique identifier
// 'id' (which must be passed as a string).
//
// Warning: This method coerces its returned data into an APIv1.5 format.
//
// Deprecated: Please used versioned library imports in the future, and
// versioned methods, specifically, for API v3.0 - in this case,
// GetDeliveryServicesV30WithHdr.
func (to *Session) GetDeliveryServiceNullable(id string) (*tc.DeliveryServiceNullable, ReqInf, error) {
	return to.GetDeliveryServiceNullableWithHdr(id, nil)
}

func (to *Session) GetDeliveryServiceByXMLIDNullableWithHdr(XMLID string, header http.Header) ([]tc.DeliveryServiceNullable, ReqInf, error) {
	var data tc.DeliveryServicesNullableResponse
	reqInf, err := get(to, API_DELIVERY_SERVICES+"?xmlId="+XMLID, &data, header)
	if err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetDeliveryServiceByXMLIDNullable returns the Delivery Service identified by the passed XMLID.
// The length of the returned slice should always be 1 when the request is succesful - if it isn't
// something very wicked has happened to Traffic Ops.
//
// Warning: This method coerces its returned data into an APIv1.5 format.
//
// Deprecated: Please used versioned library imports in the future, and
// versioned methods, specifically, for API v3.0 - in this case,
// GetDeliveryServicesV30WithHdr.
func (to *Session) GetDeliveryServiceByXMLIDNullable(XMLID string) ([]tc.DeliveryServiceNullable, ReqInf, error) {
	return to.GetDeliveryServiceByXMLIDNullableWithHdr(XMLID, nil)
}

// CreateDeliveryServiceV30 creates the Delivery Service it's passed.
func (to *Session) CreateDeliveryServiceV30(ds tc.DeliveryServiceNullableV30) (tc.DeliveryServiceNullableV30, ReqInf, error) {
	var reqInf ReqInf
	if ds.TypeID == nil && ds.Type != nil {
		ty, _, err := to.GetTypeByName(ds.Type.String())
		if err != nil {
			return tc.DeliveryServiceNullableV30{}, reqInf, err
		}
		if len(ty) == 0 {
			return tc.DeliveryServiceNullableV30{}, reqInf, fmt.Errorf("no type named %s", ds.Type)
		}
		ds.TypeID = &ty[0].ID
	}

	if ds.CDNID == nil && ds.CDNName != nil {
		cdns, _, err := to.GetCDNByName(*ds.CDNName)
		if err != nil {
			return tc.DeliveryServiceNullableV30{}, reqInf, err
		}
		if len(cdns) == 0 {
			return tc.DeliveryServiceNullableV30{}, reqInf, errors.New("no CDN named " + *ds.CDNName)
		}
		ds.CDNID = &cdns[0].ID
	}

	if ds.ProfileID == nil && ds.ProfileName != nil {
		profiles, _, err := to.GetProfileByName(*ds.ProfileName)
		if err != nil {
			return tc.DeliveryServiceNullableV30{}, reqInf, err
		}
		if len(profiles) == 0 {
			return tc.DeliveryServiceNullableV30{}, reqInf, errors.New("no Profile named " + *ds.ProfileName)
		}
		ds.ProfileID = &profiles[0].ID
	}

	if ds.TenantID == nil && ds.Tenant != nil {
		ten, _, err := to.TenantByName(*ds.Tenant)
		if err != nil {
			return tc.DeliveryServiceNullableV30{}, reqInf, err
		}
		ds.TenantID = &ten.ID
	}

	bts, err := json.Marshal(ds)
	if err != nil {
		return tc.DeliveryServiceNullableV30{}, reqInf, nil
	}

	var data tc.DeliveryServicesResponseV30
	reqInf, err = post(to, API_DELIVERY_SERVICES, bts, &data)
	if err != nil {
		return tc.DeliveryServiceNullableV30{}, reqInf, err
	}
	if len(data.Response) != 1 {
		return tc.DeliveryServiceNullableV30{}, reqInf, fmt.Errorf("failed to create Delivery Service, response indicated %d were created", len(data.Response))
	}

	return data.Response[0], reqInf, nil
}

// CreateDeliveryServiceNullable creates the DeliveryService it's passed.
//
// Warning: This method coerces its returned data into an APIv1.5 format, and
// only accepts input in an APIv1.5 format.
//
// Deprecated: Please used versioned library imports in the future, and
// versioned methods, specifically, for API v3.0 - in this case,
// CreateDeliveryServiceV30.
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

// UpdateDeliveryServiceV30WithHdr replaces the Delivery Service identified by the
// integral, unique identifier 'id' with the one it's passed.
func (to *Session) UpdateDeliveryServiceV30WithHdr(id int, ds tc.DeliveryServiceNullableV30, header http.Header) (tc.DeliveryServiceNullableV30, ReqInf, error) {
	var reqInf ReqInf
	bts, err := json.Marshal(ds)
	if err != nil {
		return tc.DeliveryServiceNullableV30{}, reqInf, err
	}

	var data tc.DeliveryServicesResponseV30
	reqInf, err = put(to, fmt.Sprintf(API_DELIVERY_SERVICE_ID, id), bts, &data, header)
	if err != nil {
		return tc.DeliveryServiceNullableV30{}, reqInf, err
	}
	if len(data.Response) != 1 {
		return tc.DeliveryServiceNullableV30{}, reqInf, fmt.Errorf("failed to update Delivery Service #%d; response indicated that %d were updated", id, len(data.Response))
	}
	return data.Response[0], reqInf, nil

}

// UpdateDeliveryServiceV31 replaces the Delivery Service identified by the
// integral, unique identifier 'id' with the one it's passed.
func (to *Session) UpdateDeliveryServiceV31(id int, ds tc.DeliveryServiceV31, header http.Header) (tc.DeliveryServiceV31, ReqInf, error) {
	var reqInf ReqInf
	bts, err := json.Marshal(ds)
	if err != nil {
		return tc.DeliveryServiceV31{}, reqInf, err
	}

	var data tc.DeliveryServicesResponseV31
	reqInf, err = put(to, fmt.Sprintf(API_DELIVERY_SERVICE_ID, id), bts, &data, header)
	if err != nil {
		return tc.DeliveryServiceV31{}, reqInf, err
	}
	if len(data.Response) != 1 {
		return tc.DeliveryServiceV31{}, reqInf, fmt.Errorf("failed to update Delivery Service #%d; response indicated that %d were updated", id, len(data.Response))
	}
	return data.Response[0], reqInf, nil

}

// UpdateDeliveryServiceNullable updates the DeliveryService matching the ID it's
// passed with the DeliveryService it is passed.
//
// Warning: This method coerces its returned data into an APIv1.5 format, and
// only accepts input in an APIv1.5 format.
//
// Deprecated: Please used versioned library imports in the future, and
// versioned methods, specifically, for API v3.0 - in this case,
// UpdateDeliveryServiceV30WithHdr.
func (to *Session) UpdateDeliveryServiceNullable(id string, ds *tc.DeliveryServiceNullable) (*tc.UpdateDeliveryServiceNullableResponse, error) {
	return to.UpdateDeliveryServiceNullableWithHdr(id, ds, nil)
}

func (to *Session) UpdateDeliveryServiceNullableWithHdr(id string, ds *tc.DeliveryServiceNullable, header http.Header) (*tc.UpdateDeliveryServiceNullableResponse, error) {
	var data tc.UpdateDeliveryServiceNullableResponse
	jsonReq, err := json.Marshal(ds)
	if err != nil {
		return nil, err
	}
	_, err = put(to, fmt.Sprintf(API_DELIVERY_SERVICE_ID, id), jsonReq, &data, header)
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

func (to *Session) GetDeliveryServiceHealthWithHdr(id string, header http.Header) (*tc.DeliveryServiceHealth, ReqInf, error) {
	var data tc.DeliveryServiceHealthResponse
	reqInf, err := get(to, fmt.Sprintf(API_DELIVERY_SERVICE_HEALTH, id), &data, nil)
	if err != nil {
		return nil, reqInf, err
	}

	return &data.Response, reqInf, nil
}

// GetDeliveryServiceHealth gets the 'health' of the Delivery Service identified by the
// integral, unique identifier 'id' (which must be passed as a string).
// Deprecated: GetDeliveryServiceHealth will be removed in 6.0. Use GetDeliveryServiceHealthWithHdr.
func (to *Session) GetDeliveryServiceHealth(id string) (*tc.DeliveryServiceHealth, ReqInf, error) {
	return to.GetDeliveryServiceHealthWithHdr(id, nil)
}

func (to *Session) GetDeliveryServiceCapacityWithHdr(id string, header http.Header) (*tc.DeliveryServiceCapacity, ReqInf, error) {
	var data tc.DeliveryServiceCapacityResponse
	reqInf, err := get(to, fmt.Sprintf(API_DELIVERY_SERVICE_CAPACITY, id), &data, header)
	if err != nil {
		return nil, reqInf, err
	}

	return &data.Response, reqInf, nil
}

// GetDeliveryServiceCapacity gets the 'capacity' of the Delivery Service identified by the
// integral, unique identifier 'id' (which must be passed as a string).
// Deprecated: GetDeliveryServiceCapacity will be removed in 6.0. Use GetDeliveryServiceCapacityWithHdr.
func (to *Session) GetDeliveryServiceCapacity(id string) (*tc.DeliveryServiceCapacity, ReqInf, error) {
	return to.GetDeliveryServiceCapacityWithHdr(id, nil)
}

// GetDeliveryServiceServer returns associations between Delivery Services and servers using the
// provided pagination controls.
// Deprecated: GetDeliveryServiceServer will be removed in 6.0. Use GetDeliveryServiceServerWithHdr.
func (to *Session) GetDeliveryServiceServer(page, limit string) ([]tc.DeliveryServiceServer, ReqInf, error) {
	return to.GetDeliveryServiceServerWithHdr(page, limit, nil)
}

func (to *Session) GetDeliveryServiceServerWithHdr(page, limit string, header http.Header) ([]tc.DeliveryServiceServer, ReqInf, error) {
	var data tc.DeliveryServiceServerResponse
	reqInf, err := get(to, API_DELIVERY_SERVICE_SERVER+"?page="+page+"&limit="+limit, &data, header)
	if err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetDeliveryServiceRegexes returns the "Regexes" (Regular Expressions) used by all (tenant-visible)
// Delivery Services.
// Deprecated: GetDeliveryServiceRegexes will be removed in 6.0. Use GetDeliveryServiceRegexesWithHdr.
func (to *Session) GetDeliveryServiceRegexes() ([]tc.DeliveryServiceRegexes, ReqInf, error) {
	return to.GetDeliveryServiceRegexesWithHdr(nil)
}

func (to *Session) GetDeliveryServiceRegexesWithHdr(header http.Header) ([]tc.DeliveryServiceRegexes, ReqInf, error) {
	var data tc.DeliveryServiceRegexResponse
	reqInf, err := get(to, API_DELIVERY_SERVICES_REGEXES, &data, header)
	if err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GenerateSSLKeysForDS generates ssl keys for a given cdn
func (to *Session) GenerateSSLKeysForDS(XMLID string, CDNName string, sslFields tc.SSLKeyRequestFields) (string, ReqInf, error) {
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
	jsonReq, err := json.Marshal(request)
	if err != nil {
		return "", ReqInf{}, err
	}
	response := struct {
		Response string `json:"response"`
	}{}
	reqInf, err := post(to, API_DELIVERY_SERVICE_GENERATE_SSL_KEYS, jsonReq, &response)
	if err != nil {
		return "", reqInf, err
	}
	return response.Response, reqInf, nil
}

func (to *Session) DeleteDeliveryServiceSSLKeysByID(XMLID string) (string, ReqInf, error) {
	resp := struct {
		Response string `json:"resposne"`
	}{}
	reqInf, err := del(to, fmt.Sprintf(API_DELIVERY_SERVICE_XMLID_SSL_KEYS, XMLID), &resp)
	if err != nil {
		return "", reqInf, err
	}
	return resp.Response, reqInf, nil
}

// GetDeliveryServiceSSLKeysByID returns information about the SSL Keys used by the Delivery
// Service identified by the passed XMLID.
// Deprecated: GetDeliveryServiceSSLKeysByID will be removed in 6.0. Use GetDeliveryServiceSSLKeysByIDWithHdr.
func (to *Session) GetDeliveryServiceSSLKeysByID(XMLID string) (*tc.DeliveryServiceSSLKeys, ReqInf, error) {
	return to.GetDeliveryServiceSSLKeysByIDWithHdr(XMLID, nil)
}

func (to *Session) GetDeliveryServiceSSLKeysByIDWithHdr(XMLID string, header http.Header) (*tc.DeliveryServiceSSLKeys, ReqInf, error) {
	var data tc.DeliveryServiceSSLKeysResponse
	reqInf, err := get(to, fmt.Sprintf(API_DELIVERY_SERVICE_XMLID_SSL_KEYS, XMLID), &data, header)
	if err != nil {
		return nil, reqInf, err
	}

	return &data.Response, reqInf, nil
}

func (to *Session) GetDeliveryServicesEligibleWithHdr(dsID int, header http.Header) ([]tc.DSServer, ReqInf, error) {
	resp := struct {
		Response []tc.DSServer `json:"response"`
	}{Response: []tc.DSServer{}}

	reqInf, err := get(to, fmt.Sprintf(API_DELIVERY_SERVICE_ELIGIBLE_SERVERS, dsID), &resp, header)
	if err != nil {
		return nil, reqInf, err
	}
	return resp.Response, reqInf, nil
}

// GetDeliveryServicesEligible returns the servers eligible for assignment to the Delivery
// Service identified by the integral, unique identifier 'dsID'.
// Deprecated: GetDeliveryServicesEligible will be removed in 6.0. Use GetDeliveryServicesEligibleWithHdr.
func (to *Session) GetDeliveryServicesEligible(dsID int) ([]tc.DSServer, ReqInf, error) {
	return to.GetDeliveryServicesEligibleWithHdr(dsID, nil)
}

// GetDeliveryServiceURLSigKeys returns the URL-signing keys used by the Delivery Service
// identified by the XMLID 'dsName'.
// Deprecated: GetDeliveryServiceURLSigKeys will be removed in 6.0. Use GetDeliveryServiceURLSigKeysWithHdr.
func (to *Session) GetDeliveryServiceURLSigKeys(dsName string) (tc.URLSigKeys, ReqInf, error) {
	return to.GetDeliveryServiceURLSigKeysWithHdr(dsName, nil)
}

func (to *Session) GetDeliveryServiceURLSigKeysWithHdr(dsName string, header http.Header) (tc.URLSigKeys, ReqInf, error) {
	data := struct {
		Response tc.URLSigKeys `json:"response"`
	}{}

	reqInf, err := get(to, fmt.Sprintf(API_DELIVERY_SERVICES_URL_SIGNING_KEYS, dsName), &data, header)
	if err != nil {
		return tc.URLSigKeys{}, reqInf, err
	}
	return data.Response, reqInf, nil
}

// Deprecated: GetDeliveryServiceURISigningKeys will be removed in 6.0. Use GetDeliveryServiceURISigningKeysWithHdr.
func (to *Session) GetDeliveryServiceURISigningKeys(dsName string) ([]byte, ReqInf, error) {
	return to.GetDeliveryServiceURISigningKeysWithHdr(dsName, nil)
}

// GetDeliveryServiceURISigningKeys returns the URI-signing keys used by the Delivery Service
// identified by the XMLID 'dsName'. The result is not parsed.
func (to *Session) GetDeliveryServiceURISigningKeysWithHdr(dsName string, header http.Header) ([]byte, ReqInf, error) {
	data := json.RawMessage{}
	reqInf, err := get(to, fmt.Sprintf(API_DELIVERY_SERVICES_URI_SIGNING_KEYS, dsName), &data, header)
	if err != nil {
		return []byte{}, reqInf, err
	}
	return []byte(data), reqInf, nil
}

// SafeDeliveryServiceUpdateV31 updates the "safe" fields of the Delivery
// Service identified by the integral, unique identifier 'id'.
func (to *Session) SafeDeliveryServiceUpdateV31(id int, r tc.DeliveryServiceSafeUpdateRequest, header http.Header) (tc.DeliveryServiceV31, ReqInf, error) {
	var reqInf ReqInf
	req, err := json.Marshal(r)
	if err != nil {
		return tc.DeliveryServiceV31{}, reqInf, err
	}

	var data tc.DeliveryServiceSafeUpdateResponseV31
	reqInf, err = put(to, fmt.Sprintf(API_DELIVERY_SERVICES_SAFE_UPDATE, id), req, &data, header)
	if err != nil {
		return tc.DeliveryServiceV31{}, reqInf, err
	}
	if len(data.Response) != 1 {
		return tc.DeliveryServiceV31{}, reqInf, fmt.Errorf("failed to safe update Delivery Service #%d; response indicated that %d were updated", id, len(data.Response))
	}
	return data.Response[0], reqInf, nil
}

// SafeDeliveryServiceUpdateV30WithHdr updates the "safe" fields of the Delivery
// Service identified by the integral, unique identifier 'id'.
func (to *Session) SafeDeliveryServiceUpdateV30WithHdr(id int, r tc.DeliveryServiceSafeUpdateRequest, header http.Header) (tc.DeliveryServiceNullableV30, ReqInf, error) {
	var reqInf ReqInf
	req, err := json.Marshal(r)
	if err != nil {
		return tc.DeliveryServiceNullableV30{}, reqInf, err
	}

	var data tc.DeliveryServiceSafeUpdateResponseV30
	reqInf, err = put(to, fmt.Sprintf(API_DELIVERY_SERVICES_SAFE_UPDATE, id), req, &data, header)
	if err != nil {
		return tc.DeliveryServiceNullableV30{}, reqInf, err
	}
	if len(data.Response) != 1 {
		return tc.DeliveryServiceNullableV30{}, reqInf, fmt.Errorf("failed to safe update Delivery Service #%d; response indicated that %d were updated", id, len(data.Response))
	}
	return data.Response[0], reqInf, nil
}

// UpdateDeliveryServiceSafe updates the given Delivery Service identified by 'id' with the given "safe" fields.
//
// Warning: This method coerces its returned data into an APIv1.5 format.
//
// Deprecated: Please used versioned library imports in the future, and
// versioned methods, specifically, for API v3.0 - in this case,
// SafeDeliveryServiceUpdateV30WithHdr.
func (to *Session) UpdateDeliveryServiceSafe(id int, ds tc.DeliveryServiceSafeUpdateRequest) ([]tc.DeliveryServiceNullable, ReqInf, error) {
	var reqInf ReqInf
	var resp tc.DeliveryServiceSafeUpdateResponse

	req, err := json.Marshal(ds)
	if err != nil {
		return resp.Response, reqInf, err
	}

	if reqInf, err = put(to, fmt.Sprintf(API_DELIVERY_SERVICES_SAFE_UPDATE, id), req, &resp, nil); err != nil {
		return resp.Response, reqInf, err
	}

	if len(resp.Response) < 1 {
		err = errors.New("Traffic Ops returned success, but response was missing the Delivery Service")
	}
	return resp.Response, reqInf, err
}

// GetAccessibleDeliveryServicesByTenant gets all delivery services associated with the given tenant and all of
// its children.
//
// Warning: This method coerces its returned data into an APIv1.5 format.
//
// Deprecated: Please used versioned library imports in the future, and
// versioned methods, specifically, for API v3.0 - in this case,
// GetDeliveryServicesV30WithHdr.
func (to *Session) GetAccessibleDeliveryServicesByTenant(tenantId int) ([]tc.DeliveryServiceNullable, ReqInf, error) {
	data := tc.DeliveryServicesNullableResponse{}
	reqInf, err := get(to, fmt.Sprintf("%v?accessibleTo=%v", API_DELIVERY_SERVICES, tenantId), &data, nil)
	return data.Response, reqInf, err
}
