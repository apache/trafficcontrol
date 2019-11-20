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
	"errors"
	"fmt"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
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
	frs, err := dbhelpers.GetFederationResolversByFederationID(inf.Tx.Tx, fedID)

	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("federations.federation_resolvers.Get getting federation resolvers from federation: "+err.Error()))
		return
	}
	frsr := tc.FederationFederationResolversResponse{Response: frs}
	api.WriteResp(w, r, frsr)
}

// AssignFederationResolversToFederation associates one or more federation_resolver to the federation ID supplied.
func AssignFederationResolversToFederation(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id", "fedResolverIds"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Txt.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	fedID := inf.IntParams["id"]
	replace := inf.Params["replace"]
	frs := inf.Params["fedResolverIds"]

	userErr, sysErr, errCode = addFederationResolverMappings(inf.User, tx, mappings)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	msg := fmt.Sprintf("%s successfully created federation resolvers.", inf.User.UserName)
	if inf.Version.Major <= 1 && inf.Version.Minor <= 3 {
		api.WriteResp(w, r, msg)
	} else {
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, msg, msg)
	}
	return
}
