package dsuser

/*
 * LICENSED to the Apache Software Foundation (ASF) under one
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
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"

	"github.com/lib/pq"
)

func Post(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	dsu := tc.DeliveryServiceUserPost{}
	if err := api.Parse(r.Body, inf.Tx.Tx, &dsu); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("decoding: "+err.Error()), nil)
		return
	}

	if userErr, sysErr, errCode := tenant.CheckDSIDs(inf.Tx.Tx, inf.User, *dsu.DeliveryServices); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	if userErr, sysErr, errCode := tenant.CheckUser(inf.Tx.Tx, inf.User, *dsu.UserID); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	if *dsu.Replace {
		if err := deleteUserDSes(inf.Tx.Tx, *dsu.UserID); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("deleting user delivery services (to replace): "+err.Error()))
			return
		}
	}

	numAssigned, err := createUserDSes(inf.Tx.Tx, *dsu.UserID, *dsu.DeliveryServices)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("creating user delivery services: "+err.Error()))
		return
	}

	api.CreateChangeLogRawTx(api.ApiChange, strconv.Itoa(numAssigned)+" delivery services were assigned to user "+strconv.Itoa(*dsu.UserID), inf.User, inf.Tx.Tx)
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Delivery service assignments complete.", dsu)
}

func deleteUserDSes(tx *sql.Tx, userID int) error {
	if _, err := tx.Exec(`DELETE FROM deliveryservice_tmuser WHERE tm_user_id = $1`, userID); err != nil {
		return errors.New("deleting user delivery services: " + err.Error())
	}
	return nil
}

// createUserDSes creates the given user delivery service assignments.
// Returns the number of delivery services assigned, and any error.
func createUserDSes(tx *sql.Tx, userID int, dsIDs []int) (int, error) {
	result, err := tx.Exec(`INSERT INTO deliveryservice_tmuser (tm_user_id, deliveryservice) values ($1, UNNEST($2::bigint[]))`, userID, pq.Array(dsIDs))
	if err != nil {
		return 0, errors.New("inserting user delivery services assignments: " + err.Error())
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.New("getting rows affected: " + err.Error())
	}
	return int(rowsAffected), nil
}
