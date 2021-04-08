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
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

// APIFederations is the API version-relative path to the /federations API route.
const APIFederations = "/federations"

// Federations gets all Delivery Service-to-Federation mappings in Traffic Ops
// that are assigned to the current user.
func (to *Session) Federations(header http.Header) ([]tc.AllDeliveryServiceFederationsMapping, toclientlib.ReqInf, error) {
	type FederationResponse struct {
		Response []tc.AllDeliveryServiceFederationsMapping `json:"response"`
	}
	data := FederationResponse{}
	inf, err := to.get(APIFederations, header, &data)
	return data.Response, inf, err
}

// AllFederations gets all Delivery Service-to-Federation mappings in Traffic
// Ops.
func (to *Session) AllFederations(header http.Header) ([]tc.AllDeliveryServiceFederationsMapping, toclientlib.ReqInf, error) {
	type FederationResponse struct {
		Response []tc.AllDeliveryServiceFederationsMapping `json:"response"`
	}
	data := FederationResponse{}
	inf, err := to.get("/federations/all", header, &data)
	return data.Response, inf, err
}

// AllFederationsByCDN gets all Delivery Service-to-Federation mappings for
// Federations in the CDN with the given name.
func (to *Session) AllFederationsByCDN(cdnName string, header http.Header) ([]tc.AllDeliveryServiceFederationsMapping, toclientlib.ReqInf, error) {
	// because the Federations JSON array is heterogeneous (array members may be a AllFederation or AllFederationCDN), we have to try decoding each separately.
	type FederationResponse struct {
		Response []json.RawMessage `json:"response"`
	}
	data := FederationResponse{}
	inf, err := to.get("/federations/all?cdnName="+url.QueryEscape(cdnName), header, &data)
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

// CreateFederationDeliveryServices assigns the Delivery Services identified in
// 'deliveryServiceIDs' with the Federation identified by 'federationID'. If
// 'replace' is true, existing assignments for the Federation are overwritten.
func (to *Session) CreateFederationDeliveryServices(federationID int, deliveryServiceIDs []int, replace bool) (toclientlib.ReqInf, error) {
	req := tc.FederationDSPost{DSIDs: deliveryServiceIDs, Replace: &replace}
	resp := map[string]interface{}{}
	inf, err := to.post(`federations/`+strconv.Itoa(federationID)+`/deliveryservices`, req, nil, &resp)
	return inf, err
}

// GetFederationDeliveryServices returns the Delivery Services assigned to the
// Federation identified by 'federationID'.
func (to *Session) GetFederationDeliveryServices(federationID int, header http.Header) ([]tc.FederationDeliveryServiceNullable, toclientlib.ReqInf, error) {
	type FederationDSesResponse struct {
		Response []tc.FederationDeliveryServiceNullable `json:"response"`
	}
	data := FederationDSesResponse{}
	inf, err := to.get(fmt.Sprintf("federations/%d/deliveryservices", federationID), header, &data)
	return data.Response, inf, err
}

// DeleteFederationDeliveryService unassigns the Delivery Service identified by
// 'deliveryServiceID' from the Federation identified by 'federationID'.
func (to *Session) DeleteFederationDeliveryService(federationID, deliveryServiceID int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("federations/%d/deliveryservices/%d", federationID, deliveryServiceID)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

// CreateFederationUsers assigns the Federation identified by 'federationID' to
// the user(s) identified in 'userIDs'. If 'replace' is true, all existing user
// assignments for the Federation are overwritten.
func (to *Session) CreateFederationUsers(federationID int, userIDs []int, replace bool) (tc.Alerts, toclientlib.ReqInf, error) {
	req := tc.FederationUserPost{IDs: userIDs, Replace: &replace}
	var alerts tc.Alerts
	inf, err := to.post(fmt.Sprintf("federations/%d/users", federationID), req, nil, &alerts)
	return alerts, inf, err
}

// GetFederationUsers retrieves all users to whom the Federation identified by
// 'federationID' is assigned.
func (to *Session) GetFederationUsers(federationID int, header http.Header) ([]tc.FederationUser, toclientlib.ReqInf, error) {
	type FederationUsersResponse struct {
		Response []tc.FederationUser `json:"response"`
	}
	data := FederationUsersResponse{}
	inf, err := to.get(fmt.Sprintf("federations/%d/users", federationID), header, &data)
	return data.Response, inf, err
}

// DeleteFederationUser unassigns the Federation identified by 'federationID'
// from the user identified by 'userID'.
func (to *Session) DeleteFederationUser(federationID, userID int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("federations/%d/users/%d", federationID, userID)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

// AddFederationResolverMappingsForCurrentUser adds Federation Resolver mappings to one or more
// Delivery Services for the current user.
func (to *Session) AddFederationResolverMappingsForCurrentUser(mappings tc.DeliveryServiceFederationResolverMappingRequest) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APIFederations, mappings, nil, &alerts)
	return alerts, reqInf, err
}

// DeleteFederationResolverMappingsForCurrentUser removes ALL Federation Resolver mappings for ALL
// Federations assigned to the currently authenticated user, as well as deleting ALL of the
// Federation Resolvers themselves.
func (to *Session) DeleteFederationResolverMappingsForCurrentUser() (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.del(APIFederations, nil, &alerts)
	return alerts, reqInf, err
}

// ReplaceFederationResolverMappingsForCurrentUser replaces any and all Federation Resolver mappings
// on all Federations assigned to the currently authenticated user. This will first remove ALL
// Federation Resolver mappings for ALL Federations assigned to the currently authenticated user, as
// well as deleting ALL of the Federation Resolvers themselves. In other words, calling this is
// equivalent to a call to DeleteFederationResolverMappingsForCurrentUser followed by a call to
// AddFederationResolverMappingsForCurrentUser .
func (to *Session) ReplaceFederationResolverMappingsForCurrentUser(mappings tc.DeliveryServiceFederationResolverMappingRequest) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.put(APIFederations, mappings, nil, &alerts)
	return alerts, reqInf, err
}
