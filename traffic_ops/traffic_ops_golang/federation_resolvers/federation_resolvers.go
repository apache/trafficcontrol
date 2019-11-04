package federation_resolvers

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

import "net/http"

import "github.com/apache/trafficcontrol/lib/go-tc"

import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"

const insertFederationResolverQuery = `
INSERT INTO federation_resolver (ip_address, type)
VALUES ($1, $2)
RETURNING (
	federation_resolver.id,
	federation_resolver.ip_address,
	(
		SELECT type.name
		FROM type
		WHERE type.id = federation_resolver.type
	) AS type
)
`

func Create(w http.ResponseWriter, r *http.Request) {
	inf, sysErr, userErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if sysErr != nil || userErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	var fr tc.FederationResolver
	if userErr = api.Parse(r.Body, tx, &fr); userErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, userErr, nil)
		return
	}

	if err := tx.QueryRow(insertFederationResolverQuery, fr.IPAddress, fr.TypeID).Scan(&fr); err != nil {
		userErr, sysErr errCode = api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	fr.TypeID = nil
	if inf.Version.Major < 2 && inf.Version.Minor < 4 {
		fr.LastUpdated = nil
	}

	changeLogMsg := fmt.Sprintf("FEDERATION_RESOLVER: %s, ID: %d, ACTION: Created", fr.IPAddress, fr.ID)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)

	alertMsg := fmt.Sprintf("Federation Resolver created [ IP = %s ] with id: %d", fr.IPAddress, fr.ID)
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, alertMsg, fr)
}
