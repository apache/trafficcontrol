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
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}

	frsr := tc.FederationFederationResolversResponse{Response: frs}
	api.WriteResp(w, r, frsr)
}

// AssignFederationResolversToFederation associates one or more federation_resolver to the federation ID supplied.
func AssignFederationResolversToFederation(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	ffrr := tc.AssignFederationResolversRequest{}
	fedID := inf.IntParams["id"]
	if err := json.NewDecoder(r.Body).Decode(&ffrr); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	// TODO: "Fed Resolver IDs must be an array"

	name, _, err := dbhelpers.GetFederationNameFromID(fedID, inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
		return
	}

	if ffrr.Replace {
		if _, err := inf.Tx.Tx.Exec(deleteFederationFederationResolversQuery, fedID); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
			return
		}
	}
	for _, id := range ffrr.FedResolverIDs {
		if _, err := inf.Tx.Tx.Exec(associateFederationWithResolverQuery, fedID, id); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
			return
		}
	}

	api.WriteRespAlertObj(
		w, r, tc.SuccessLevel,
		fmt.Sprintf("%d resolver(s) were assigned to the %s federation", len(ffrr.FedResolverIDs), name),
		tc.AssignFederationFederationResolversResponse{Response: ffrr},
	)
}
