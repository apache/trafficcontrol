package servers

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
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
)

// Delete handler for deleting the association between a Delivery Service and a Server
func Delete(w http.ResponseWriter, r *http.Request) {
	delete(w, r)
}

const lastServerQuery = `
SELECT
(SELECT (CASE WHEN t.name LIKE '` + string(tc.EdgeTypePrefix) + `%' THEN TRUE ELSE FALSE END) AS available
FROM type t
JOIN server s ON s.type = t.id
WHERE s.id = $2)
AND
(SELECT COUNT(*) = 0 AS available
FROM deliveryservice_server
JOIN server s ON deliveryservice_server.server = s.id
JOIN type t ON t.id = s.type
JOIN status st ON st.id = s.status
WHERE (st.name = '` + string(tc.CacheStatusOnline) + `' OR st.name = '` + string(tc.CacheStatusReported) + `')
AND t.name LIKE '` + string(tc.EdgeTypePrefix) + `%'
AND deliveryservice = $1
AND server <> $2)
`

// checkLastServer checks if the given Server ID identifies the last server
// assigned to the Delivery Service identified by its passed ID. It returns -
// in order - an HTTP status code (useful only if an error occurs), an error
// suitable for reporting back to the user, and an error that must not be shown
// to the user. If the server is, in fact, the last server assigned to the
// Delivery Service, the "user error" will be set to an appropriate, non-nil
// value.
func checkLastServer(dsID, serverID int, tx *sql.Tx) (int, error, error) {
	var isLast bool
	if tx == nil {
		return http.StatusInternalServerError, nil, errors.New("nil transaction")
	}
	err := tx.QueryRow(lastServerQuery, dsID, serverID).Scan(&isLast)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("checking if server #%d is the last one assigned to DS #%d: %v", serverID, dsID, err)
	}
	if isLast {
		return http.StatusConflict, fmt.Errorf("removing server #%d from active Delivery Service #%d would leave it with no REPORTED/ONLINE assigned servers", serverID, dsID), nil
	}
	return http.StatusOK, nil, nil
}

func delete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"serverid", "dsid"}, []string{"serverid", "dsid"})
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	if tx == nil {
		errCode = http.StatusInternalServerError
		sysErr = errors.New("nil transaction")
		api.HandleErr(w, r, inf.Tx.Tx, errCode, nil, sysErr)
		return
	}

	dsID := inf.IntParams["dsid"]
	serverID := inf.IntParams["serverid"]

	userErr, sysErr, errCode = tenant.CheckID(tx, inf.User, dsID)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	query := deliveryservice.SelectDeliveryServicesQuery + " WHERE ds.id=:dsid"
	vals := map[string]interface{}{"dsid": dsID}
	dses, userErr, sysErr, errCode := deliveryservice.GetDeliveryServices(query, vals, inf.Tx)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	if len(dses) < 1 {
		errCode = http.StatusNotFound
		userErr = fmt.Errorf("no such Delivery Service: #%d", dsID)
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}
	if len(dses) > 1 {
		errCode = http.StatusInternalServerError
		sysErr = fmt.Errorf("too many Delivery Services with ID %d: %d", dsID, len(dses))
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}

	ds := dses[0]
	if ds.Active == nil {
		errCode = http.StatusInternalServerError
		sysErr = fmt.Errorf("Delivery Service #%d had nil Active", dsID)
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}

	if *ds.Active {
		errCode, userErr, sysErr = checkLastServer(dsID, serverID, tx)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
			return
		}
	}
	if ds.XMLID == nil {
		errCode = http.StatusInternalServerError
		sysErr = fmt.Errorf("Delivery Service #%d had nil XMLID", dsID)
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}
	dsName := *ds.XMLID

	if ds.CDNName != nil {
		userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(inf.Tx.Tx, *ds.CDNName, inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
			return
		}
	}
	serverName, exists, err := dbhelpers.GetServerNameFromID(tx, serverID)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting server name from id: "+err.Error()))
		return
	} else if !exists {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("server not found"), nil)
		return
	}

	ok, err := deleteDSServer(tx, dsID, serverID)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("deleting delivery service server: "+err.Error()))
		return
	}
	if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, nil)
		return
	}
	api.CreateChangeLogRawTx(api.ApiChange, "DS: "+string(dsName)+", ID: "+strconv.Itoa(dsID)+", ACTION: Remove server "+string(serverName)+" from delivery service", inf.User, inf.Tx.Tx)
	api.WriteRespAlert(w, r, tc.SuccessLevel, "Server unlinked from delivery service.")
}

// deleteDSServer deletes the given deliveryservice_server. Returns whether the server existed, and any error.
func deleteDSServer(tx *sql.Tx, dsID int, serverID int) (bool, error) {
	deletedServerID := 0
	if err := tx.QueryRow(`DELETE FROM deliveryservice_server WHERE deliveryservice = $1 AND server = $2 RETURNING server`, dsID, serverID).Scan(&deletedServerID); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.New("deleting delivery service server: " + err.Error())
	}
	return true, nil
}
