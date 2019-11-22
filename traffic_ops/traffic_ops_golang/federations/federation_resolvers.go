package federations

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

// GetFederationFederationResolvers returns a subset of federation_resolvers belonging to the federation ID supplied.
func GetFederationFederationResolvers(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	fedID := inf.IntParams["id"]
	frs, _, err := dbhelpers.GetFederationResolversByFederationID(inf.Tx.Tx, fedID)
	te := tc.APIError{
		Err:         err,
		Action:      "federations.federation_resolvers.GetFederationFederationResolvers",
		Description: "getting federation_resolvers from federation ID",
	}
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, te)
		return
	}

	frsr := tc.FederationFederationResolversResponse{Response: frs}
	api.WriteResp(w, r, frsr)
}

// AssignFederationResolversToFederation associates one or more federation_resolver to the federation ID supplied.
func AssignFederationResolversToFederation(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id", "fedResolverIds"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	ffrr := tc.AssignFederationResolversRequest{}
	if err := json.NewDecoder(r.Body).Decode(&ffrr); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}

	userErr, sysErr, errCode = addFederationResolversToFederation(inf.Tx.Tx, ffrr.ID, ffrr.FedResolverIDs)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	name, _, err := dbhelpers.GetFederationNameFromID(ffrr.ID, inf.Tx.Tx)
	te := tc.APIError{
		Err:         err,
		Action:      "federations.federation_resolvers.AssignFederationResolversToFederation",
		Description: "getting federation name from ID",
	}
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, te)
		return
	}

	api.WriteRespAlertObj(
		w,
		r,
		tc.SuccessLevel,
		fmt.Sprintf("%d resolver(s) were assigned to the %s federation", len(ffrr.FedResolverIDs), name),
		tc.AssignFederationFederationResolversResponse{
			Response: ffrr,
		},
	)

	return
}

func addFederationResolversToFederation(tx *sql.Tx, fedID int, frIDs []int) (error, error, int) {
	return nil, nil, 0
}
