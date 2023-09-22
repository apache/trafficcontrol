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
	"fmt"
	"strconv"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiFederations is the API version-relative path to the /federations API route.
const apiFederations = "/federations"

// apiFederationsAll is the API version-relative path to the /federations/all
// API route.
const apiFederationsAll = "/federations/all"

// Federations gets all Delivery Service-to-Federation mappings in Traffic Ops
// that are assigned to the current user.
func (to *Session) Federations(opts RequestOptions) (tc.FederationsResponse, toclientlib.ReqInf, error) {
	var data tc.FederationsResponse
	inf, err := to.get(apiFederations, opts, &data)
	return data, inf, err
}

// AllFederations gets all Delivery Service-to-Federation mappings in Traffic
// Ops.
func (to *Session) AllFederations(opts RequestOptions) (tc.FederationsResponse, toclientlib.ReqInf, error) {
	var data tc.FederationsResponse
	inf, err := to.get(apiFederationsAll, opts, &data)
	return data, inf, err
}

// CreateFederationDeliveryServices assigns the Delivery Services identified in
// 'deliveryServiceIDs' with the Federation identified by 'federationID'. If
// 'replace' is true, existing assignments for the Federation are overwritten.
func (to *Session) CreateFederationDeliveryServices(
	federationID int,
	deliveryServiceIDs []int,
	replace bool,
	opts RequestOptions,
) (tc.Alerts, toclientlib.ReqInf, error) {
	req := tc.FederationDSPost{DSIDs: deliveryServiceIDs, Replace: &replace}
	var alerts tc.Alerts
	inf, err := to.post(`federations/`+strconv.Itoa(federationID)+`/deliveryservices`, opts, req, &alerts)
	return alerts, inf, err
}

// GetFederationDeliveryServices returns the Delivery Services assigned to the
// Federation identified by 'federationID'.
func (to *Session) GetFederationDeliveryServices(federationID int, opts RequestOptions) (tc.FederationDeliveryServicesResponse, toclientlib.ReqInf, error) {
	var data tc.FederationDeliveryServicesResponse
	inf, err := to.get(fmt.Sprintf("federations/%d/deliveryservices", federationID), opts, &data)
	return data, inf, err
}

// DeleteFederationDeliveryService unassigns the Delivery Service identified by
// 'deliveryServiceID' from the Federation identified by 'federationID'.
func (to *Session) DeleteFederationDeliveryService(federationID, deliveryServiceID int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("federations/%d/deliveryservices/%d", federationID, deliveryServiceID)
	var alerts tc.Alerts
	reqInf, err := to.del(route, opts, &alerts)
	return alerts, reqInf, err
}

// CreateFederationUsers assigns the Federation identified by 'federationID' to
// the user(s) identified in 'userIDs'. If 'replace' is true, all existing user
// assignments for the Federation are overwritten.
func (to *Session) CreateFederationUsers(federationID int, userIDs []int, replace bool, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	req := tc.FederationUserPost{IDs: userIDs, Replace: &replace}
	var alerts tc.Alerts
	inf, err := to.post(fmt.Sprintf("federations/%d/users", federationID), opts, req, &alerts)
	return alerts, inf, err
}

// GetFederationUsers retrieves all users to whom the Federation identified by
// 'federationID' is assigned.
func (to *Session) GetFederationUsers(federationID int, opts RequestOptions) (tc.FederationUsersResponse, toclientlib.ReqInf, error) {
	var data tc.FederationUsersResponse
	inf, err := to.get(fmt.Sprintf("federations/%d/users", federationID), opts, &data)
	return data, inf, err
}

// DeleteFederationUser unassigns the Federation identified by 'federationID'
// from the user identified by 'userID'.
func (to *Session) DeleteFederationUser(federationID, userID int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("federations/%d/users/%d", federationID, userID)
	var alerts tc.Alerts
	reqInf, err := to.del(route, opts, &alerts)
	return alerts, reqInf, err
}

// AddFederationResolverMappingsForCurrentUser adds Federation Resolver mappings to one or more
// Delivery Services for the current user.
func (to *Session) AddFederationResolverMappingsForCurrentUser(
	mappings tc.DeliveryServiceFederationResolverMappingRequest,
	opts RequestOptions,
) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiFederations, opts, mappings, &alerts)
	return alerts, reqInf, err
}

// DeleteFederationResolverMappingsForCurrentUser removes ALL Federation Resolver mappings for ALL
// Federations assigned to the currently authenticated user, as well as deleting ALL of the
// Federation Resolvers themselves.
func (to *Session) DeleteFederationResolverMappingsForCurrentUser(opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.del(apiFederations, opts, &alerts)
	return alerts, reqInf, err
}

// ReplaceFederationResolverMappingsForCurrentUser replaces any and all Federation Resolver mappings
// on all Federations assigned to the currently authenticated user. This will first remove ALL
// Federation Resolver mappings for ALL Federations assigned to the currently authenticated user, as
// well as deleting ALL of the Federation Resolvers themselves. In other words, calling this is
// equivalent to a call to DeleteFederationResolverMappingsForCurrentUser followed by a call to
// AddFederationResolverMappingsForCurrentUser .
func (to *Session) ReplaceFederationResolverMappingsForCurrentUser(
	mappings tc.DeliveryServiceFederationResolverMappingRequest,
	opts RequestOptions,
) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.put(apiFederations, opts, mappings, &alerts)
	return alerts, reqInf, err
}
