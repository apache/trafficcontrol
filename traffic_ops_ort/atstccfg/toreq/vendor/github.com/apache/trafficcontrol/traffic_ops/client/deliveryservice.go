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
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

// DeliveryServices gets an array of DeliveryServices
// Deprecated: use GetDeliveryServices
func (to *Session) DeliveryServices() ([]tc.DeliveryService, error) {
	dses, _, err := to.GetDeliveryServices()
	return dses, err
}

func (to *Session) GetDeliveryServices() ([]tc.DeliveryService, ReqInf, error) {
	var data tc.DeliveryServicesResponse
	reqInf, err := get(to, deliveryServicesEp(), &data)
	if err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// DeliveryServicesByServer gets an array of all DeliveryServices with the given server ID assigend.
// Deprecated: use GetDeliveryServicesByServer
func (to *Session) DeliveryServicesByServer(id int) ([]tc.DeliveryService, error) {
	dses, _, err := to.GetDeliveryServicesByServer(id)
	return dses, err
}

func (to *Session) GetDeliveryServicesByServer(id int) ([]tc.DeliveryService, ReqInf, error) {
	var data tc.DeliveryServicesResponse

	reqInf, err := get(to, deliveryServicesByServerEp(strconv.Itoa(id)), &data)
	if err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

func (to *Session) GetDeliveryServiceByXMLID(XMLID string) ([]tc.DeliveryService, ReqInf, error) {
	var data tc.GetDeliveryServiceResponse
	reqInf, err := get(to, deliveryServicesByXMLID(XMLID), &data)
	if err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// DeliveryService gets the DeliveryService for the ID it's passed
// Deprecated: use GetDeliveryService
func (to *Session) DeliveryService(id string) (*tc.DeliveryService, error) {
	ds, _, err := to.GetDeliveryService(id)
	return ds, err
}

func (to *Session) GetDeliveryService(id string) (*tc.DeliveryService, ReqInf, error) {
	var data tc.DeliveryServicesResponse
	reqInf, err := get(to, deliveryServiceEp(id), &data)
	if err != nil {
		return nil, reqInf, err
	}
	if len(data.Response) == 0 {
		return nil, reqInf, nil
	}
	return &data.Response[0], reqInf, nil
}

func (to *Session) GetDeliveryServicesNullable() ([]tc.DeliveryServiceNullable, ReqInf, error) {
	data := struct {
		Response []tc.DeliveryServiceNullable `json:"response"`
	}{}
	reqInf, err := get(to, deliveryServicesEp(), &data)
	if err != nil {
		return nil, reqInf, err
	}
	return data.Response, reqInf, nil
}

func (to *Session) GetDeliveryServicesByCDNID(cdnID int) ([]tc.DeliveryServiceNullable, ReqInf, error) {
	data := struct {
		Response []tc.DeliveryServiceNullable `json:"response"`
	}{}
	reqInf, err := get(to, apiBase+dsPath+"?cdn="+strconv.Itoa(cdnID), &data)
	if err != nil {
		return nil, reqInf, err
	}
	return data.Response, reqInf, nil
}

func (to *Session) GetDeliveryServiceNullable(id string) (*tc.DeliveryServiceNullable, ReqInf, error) {
	data := struct {
		Response []tc.DeliveryServiceNullable `json:"response"`
	}{}
	reqInf, err := get(to, deliveryServiceEp(id), &data)
	if err != nil {
		return nil, reqInf, err
	}
	if len(data.Response) == 0 {
		return nil, reqInf, nil
	}
	return &data.Response[0], reqInf, nil
}

// CreateDeliveryService creates the DeliveryService it's passed
func (to *Session) CreateDeliveryService(ds *tc.DeliveryService) (*tc.CreateDeliveryServiceResponse, error) {
	if ds.TypeID == 0 && ds.Type.String() != "" {
		ty, _, err := to.GetTypeByName(ds.Type.String())
		if err != nil {
			return nil, err
		}
		if len(ty) == 0 {
			return nil, errors.New("no type named " + ds.Type.String())
		}
		ds.TypeID = ty[0].ID
	}

	if ds.CDNID == 0 && ds.CDNName != "" {
		cdns, _, err := to.GetCDNByName(ds.CDNName)
		if err != nil {
			return nil, err
		}
		if len(cdns) == 0 {
			return nil, errors.New("no CDN named " + ds.CDNName)
		}
		ds.CDNID = cdns[0].ID
	}

	if ds.ProfileID == 0 && ds.ProfileName != "" {
		profiles, _, err := to.GetProfileByName(ds.ProfileName)
		if err != nil {
			return nil, err
		}
		if len(profiles) == 0 {
			return nil, errors.New("no Profile named " + ds.ProfileName)
		}
		ds.ProfileID = profiles[0].ID
	}

	if ds.TenantID == 0 && ds.Tenant != "" {
		ten, _, err := to.TenantByName(ds.Tenant)
		if err != nil {
			return nil, err
		}
		ds.TenantID = ten.ID
	}

	var data tc.CreateDeliveryServiceResponse
	jsonReq, err := json.Marshal(ds)
	if err != nil {
		return nil, err
	}
	_, err = post(to, deliveryServicesEp(), jsonReq, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// CreateDeliveryServiceNullable creates the DeliveryService it's passed
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
	_, err = post(to, deliveryServicesEp(), jsonReq, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// UpdateDeliveryService updates the DeliveryService matching the ID it's passed with
// the DeliveryService it is passed
func (to *Session) UpdateDeliveryService(id string, ds *tc.DeliveryService) (*tc.UpdateDeliveryServiceResponse, error) {
	var data tc.UpdateDeliveryServiceResponse
	jsonReq, err := json.Marshal(ds)
	if err != nil {
		return nil, err
	}
	_, err = put(to, deliveryServiceEp(id), jsonReq, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (to *Session) UpdateDeliveryServiceNullable(id string, ds *tc.DeliveryServiceNullable) (*tc.UpdateDeliveryServiceResponse, error) {
	var data tc.UpdateDeliveryServiceResponse
	jsonReq, err := json.Marshal(ds)
	if err != nil {
		return nil, err
	}
	_, err = put(to, deliveryServiceEp(id), jsonReq, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// DeleteDeliveryService deletes the DeliveryService matching the ID it's passed
func (to *Session) DeleteDeliveryService(id string) (*tc.DeleteDeliveryServiceResponse, error) {
	var data tc.DeleteDeliveryServiceResponse
	_, err := del(to, deliveryServiceEp(id), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// DeliveryServiceState gets the DeliveryServiceState for the ID it's passed
// Deprecated: use GetDeliveryServiceState
func (to *Session) DeliveryServiceState(id string) (*tc.DeliveryServiceState, error) {
	dss, _, err := to.GetDeliveryServiceState(id)
	return dss, err
}

func (to *Session) GetDeliveryServiceState(id string) (*tc.DeliveryServiceState, ReqInf, error) {
	var data tc.DeliveryServiceStateResponse
	reqInf, err := get(to, deliveryServiceStateEp(id), &data)
	if err != nil {
		return nil, reqInf, err
	}

	return &data.Response, reqInf, nil
}

// DeliveryServiceHealth gets the DeliveryServiceHealth for the ID it's passed
// Deprecated: use GetDeliveryServiceHealth
func (to *Session) DeliveryServiceHealth(id string) (*tc.DeliveryServiceHealth, error) {
	dsh, _, err := to.GetDeliveryServiceHealth(id)
	return dsh, err
}

func (to *Session) GetDeliveryServiceHealth(id string) (*tc.DeliveryServiceHealth, ReqInf, error) {
	var data tc.DeliveryServiceHealthResponse
	reqInf, err := get(to, deliveryServiceHealthEp(id), &data)
	if err != nil {
		return nil, reqInf, err
	}

	return &data.Response, reqInf, nil
}

// DeliveryServiceCapacity gets the DeliveryServiceCapacity for the ID it's passed
// Deprecated: use GetDeliveryServiceCapacity
func (to *Session) DeliveryServiceCapacity(id string) (*tc.DeliveryServiceCapacity, error) {
	dsc, _, err := to.GetDeliveryServiceCapacity(id)
	return dsc, err
}

func (to *Session) GetDeliveryServiceCapacity(id string) (*tc.DeliveryServiceCapacity, ReqInf, error) {
	var data tc.DeliveryServiceCapacityResponse
	reqInf, err := get(to, deliveryServiceCapacityEp(id), &data)
	if err != nil {
		return nil, reqInf, err
	}

	return &data.Response, reqInf, nil
}

// DeliveryServiceRouting gets the DeliveryServiceRouting for the ID it's passed
// Deprecated: use GetDeliveryServiceRouting
func (to *Session) DeliveryServiceRouting(id string) (*tc.DeliveryServiceRouting, error) {
	dsr, _, err := to.GetDeliveryServiceRouting(id)
	return dsr, err
}

func (to *Session) GetDeliveryServiceRouting(id string) (*tc.DeliveryServiceRouting, ReqInf, error) {
	var data tc.DeliveryServiceRoutingResponse
	reqInf, err := get(to, deliveryServiceRoutingEp(id), &data)
	if err != nil {
		return nil, reqInf, err
	}

	return &data.Response, reqInf, nil
}

// DeliveryServiceServer gets the DeliveryServiceServer
// Deprecated: use GetDeliveryServiceServer
func (to *Session) DeliveryServiceServer(page, limit string) ([]tc.DeliveryServiceServer, error) {
	dss, _, err := to.GetDeliveryServiceServer(page, limit)
	return dss, err
}

func (to *Session) GetDeliveryServiceServer(page, limit string) ([]tc.DeliveryServiceServer, ReqInf, error) {
	var data tc.DeliveryServiceServerResponse
	reqInf, err := get(to, deliveryServiceServerEp(page, limit), &data)
	if err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// DeliveryServiceRegexes gets the DeliveryService regexes
// Deprecated: use GetDeliveryServiceRegexes
func (to *Session) DeliveryServiceRegexes() ([]tc.DeliveryServiceRegexes, error) {
	dsrs, _, err := to.GetDeliveryServiceRegexes()
	return dsrs, err
}
func (to *Session) GetDeliveryServiceRegexes() ([]tc.DeliveryServiceRegexes, ReqInf, error) {
	var data tc.DeliveryServiceRegexResponse
	reqInf, err := get(to, deliveryServiceRegexesEp(), &data)
	if err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// DeliveryServiceSSLKeysByID gets the DeliveryServiceSSLKeys by ID
// Deprecated: use GetDeliveryServiceSSLKeysByID
func (to *Session) DeliveryServiceSSLKeysByID(id string) (*tc.DeliveryServiceSSLKeys, error) {
	dsks, _, err := to.GetDeliveryServiceSSLKeysByID(id)
	return dsks, err
}

func (to *Session) GetDeliveryServiceSSLKeysByID(id string) (*tc.DeliveryServiceSSLKeys, ReqInf, error) {
	var data tc.DeliveryServiceSSLKeysResponse
	reqInf, err := get(to, deliveryServiceSSLKeysByIDEp(id), &data)
	if err != nil {
		return nil, reqInf, err
	}

	return &data.Response, reqInf, nil
}

// DeliveryServiceSSLKeysByHostname gets the DeliveryServiceSSLKeys by Hostname
// Deprecated: use GetDeliveryServiceSSLKeysByHostname
func (to *Session) DeliveryServiceSSLKeysByHostname(hostname string) (*tc.DeliveryServiceSSLKeys, error) {
	dsks, _, err := to.GetDeliveryServiceSSLKeysByHostname(hostname)
	return dsks, err
}

func (to *Session) GetDeliveryServiceSSLKeysByHostname(hostname string) (*tc.DeliveryServiceSSLKeys, ReqInf, error) {
	var data tc.DeliveryServiceSSLKeysResponse
	reqInf, err := get(to, deliveryServiceSSLKeysByHostnameEp(hostname), &data)
	if err != nil {
		return nil, reqInf, err
	}

	return &data.Response, reqInf, nil
}

func (to *Session) GetDeliveryServiceMatches() ([]tc.DeliveryServicePatterns, ReqInf, error) {
	uri := apiBase + `/deliveryservice_matches`
	resp := tc.DeliveryServiceMatchesResponse{}
	reqInf, err := get(to, uri, &resp)
	if err != nil {
		return nil, reqInf, err
	}
	return resp.Response, reqInf, nil
}

func (to *Session) GetDeliveryServicesEligible(dsID int) ([]tc.DSServer, ReqInf, error) {
	resp := struct {
		Response []tc.DSServer `json:"response"`
	}{Response: []tc.DSServer{}}
	uri := apiBase + `/deliveryservices/` + strconv.Itoa(dsID) + `/servers/eligible`
	reqInf, err := get(to, uri, &resp)
	if err != nil {
		return nil, reqInf, err
	}
	return resp.Response, reqInf, nil
}

func (to *Session) GetDeliveryServiceURLSigKeys(dsName string) (tc.URLSigKeys, ReqInf, error) {
	data := struct {
		Response tc.URLSigKeys `json:"response"`
	}{}
	path := apiBase + `/deliveryservices/xmlId/` + dsName + `/urlkeys.json`
	reqInf, err := get(to, path, &data)
	if err != nil {
		return tc.URLSigKeys{}, reqInf, err
	}
	return data.Response, reqInf, nil
}

func (to *Session) GetDeliveryServiceURISigningKeys(dsName string) ([]byte, ReqInf, error) {
	path := apiBase + `/deliveryservices/` + dsName + `/urisignkeys`
	data := json.RawMessage{}
	reqInf, err := get(to, path, &data)
	if err != nil {
		return []byte{}, reqInf, err
	}
	return []byte(data), reqInf, nil
}
