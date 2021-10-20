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
	"strconv"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

func (to *Session) Federations() ([]tc.AllDeliveryServiceFederationsMapping, ReqInf, error) {
	type FederationResponse struct {
		Response []tc.AllDeliveryServiceFederationsMapping `json:"response"`
	}
	data := FederationResponse{}
	inf, err := get(to, apiBase+"/federations", &data)
	return data.Response, inf, err
}

func (to *Session) AllFederations() ([]tc.AllDeliveryServiceFederationsMapping, ReqInf, error) {
	type FederationResponse struct {
		Response []tc.AllDeliveryServiceFederationsMapping `json:"response"`
	}
	data := FederationResponse{}
	inf, err := get(to, apiBase+"/federations/all", &data)
	return data.Response, inf, err
}

func (to *Session) AllFederationsForCDN(cdnName string) ([]tc.AllDeliveryServiceFederationsMapping, ReqInf, error) {
	// because the Federations JSON array is heterogeneous (array members may be a AllFederation or AllFederationCDN), we have to try decoding each separately.
	type FederationResponse struct {
		Response []json.RawMessage `json:"response"`
	}
	data := FederationResponse{}
	inf, err := get(to, apiBase+"/federations/all?cdnName="+cdnName, &data)
	if err != nil {
		return nil, inf, err
	}

	feds := []tc.AllDeliveryServiceFederationsMapping{}
	for _, raw := range data.Response {
		fed := tc.AllDeliveryServiceFederationsMapping{}
		if err := json.Unmarshal([]byte(raw), &fed); err != nil {
			// we don't actually need the CDN, but we want to return an error if we got something unexpected
			cdnFed := tc.AllFederationCDN{}
			if err := json.Unmarshal([]byte(raw), &cdnFed); err != nil {
				return nil, inf, errors.New("Traffic Ops returned an unexpected object: '" + string(raw) + "'")
			}
		}
		feds = append(feds, fed)
	}
	return feds, inf, nil
}

func (to *Session) CreateFederationDeliveryServices(federationID int, deliveryServiceIDs []int, replace bool) (ReqInf, error) {
	req := tc.FederationDSPost{DSIDs: deliveryServiceIDs, Replace: &replace}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return ReqInf{CacheHitStatus: CacheHitStatusMiss}, err
	}
	resp := map[string]interface{}{}
	inf, err := makeReq(to, http.MethodPost, apiBase+`/federations/`+strconv.Itoa(federationID)+`/deliveryservices`, jsonReq, &resp)
	return inf, err
}

// GetFederationDeliveryServices Returns a given Federation's Delivery Services
func (to *Session) GetFederationDeliveryServices(federationID int) ([]tc.FederationDeliveryServiceNullable, ReqInf, error) {
	type FederationDSesResponse struct {
		Response []tc.FederationDeliveryServiceNullable `json:"response"`
	}
	data := FederationDSesResponse{}
	inf, err := get(to, fmt.Sprintf("%s/federations/%v/deliveryservices", apiBase, federationID), &data)
	return data.Response, inf, err
}

// DeleteFederationDeliveryService Deletes a given Delivery Service from a Federation
func (to *Session) DeleteFederationDeliveryService(federationID, deliveryServiceID int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/federations/%v/deliveryservices/%v", apiBase, federationID, deliveryServiceID)
	resp, remoteAddr, err := to.request(http.MethodDelete, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	if err = json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return tc.Alerts{}, reqInf, err
	}
	return alerts, reqInf, nil
}

// GetFederationUsers Associates the given Users' IDs to a Federation
func (to *Session) CreateFederationUsers(federationID int, userIDs []int, replace bool) (tc.Alerts, ReqInf, error) {
	req := tc.FederationUserPost{IDs: userIDs, Replace: &replace}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return tc.Alerts{}, ReqInf{CacheHitStatus: CacheHitStatusMiss}, err
	}
	var alerts tc.Alerts
	inf, err := makeReq(to, http.MethodPost, fmt.Sprintf("%s/federations/%v/users", apiBase, federationID), jsonReq, &alerts)
	return alerts, inf, err
}

// GetFederationUsers Returns a given Federation's Users
func (to *Session) GetFederationUsers(federationID int) ([]tc.FederationUser, ReqInf, error) {
	type FederationUsersResponse struct {
		Response []tc.FederationUser `json:"response"`
	}
	data := FederationUsersResponse{}
	inf, err := get(to, fmt.Sprintf("%s/federations/%v/users", apiBase, federationID), &data)
	return data.Response, inf, err
}

// DeleteFederationUser Deletes a given User from a Federation
func (to *Session) DeleteFederationUser(federationID, userID int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/federations/%v/users/%v", apiBase, federationID, userID)
	resp, remoteAddr, err := to.request(http.MethodDelete, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	if err = json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return tc.Alerts{}, reqInf, err
	}
	return alerts, reqInf, nil
}

// AddFederationResolverMappingsForCurrentUser adds Federation Resolver mappings to one or more
// Delivery Services for the current user.
func (to *Session) AddFederationResolverMappingsForCurrentUser(mappings tc.DeliveryServiceFederationResolverMappingRequest) (tc.Alerts, ReqInf, error) {
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss}
	var alerts tc.Alerts

	bts, err := json.Marshal(mappings)
	if err != nil {
		return alerts, reqInf, err
	}

	resp, remoteAddr, err := to.request(http.MethodPost, apiBase+"/federations", bts)
	reqInf.RemoteAddr = remoteAddr
	if err != nil {
		return alerts, reqInf, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, err
}

// DeleteFederationResolverMappingsForCurrentUser removes ALL Federation Resolver mappings for ALL
// Federations assigned to the currently authenticated user, as well as deleting ALL of the
// Federation Resolvers themselves.
func (to *Session) DeleteFederationResolverMappingsForCurrentUser() (tc.Alerts, ReqInf, error) {
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss}
	var alerts tc.Alerts

	resp, remoteAddr, err := to.request(http.MethodDelete, apiBase+"/federations", nil)
	reqInf.RemoteAddr = remoteAddr
	if err != nil {
		return alerts, reqInf, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, err
}

// ReplaceFederationResolverMappingsForCurrentUser replaces any and all Federation Resolver mappings
// on all Federations assigned to the currently authenticated user. This will first remove ALL
// Federation Resolver mappings for ALL Federations assigned to the currently authenticated user, as
// well as deleting ALL of the Federation Resolvers themselves. In other words, calling this is
// equivalent to a call to DeleteFederationResolverMappingsForCurrentUser followed by a call to
// AddFederationResolverMappingsForCurrentUser .
func (to *Session) ReplaceFederationResolverMappingsForCurrentUser(mappings tc.DeliveryServiceFederationResolverMappingRequest) (tc.Alerts, ReqInf, error) {
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss}
	var alerts tc.Alerts

	bts, err := json.Marshal(mappings)
	if err != nil {
		return alerts, reqInf, err
	}

	resp, remoteAddr, err := to.request(http.MethodPut, apiBase+"/federations", bts)
	reqInf.RemoteAddr = remoteAddr
	if err != nil {
		return alerts, reqInf, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, err
}
