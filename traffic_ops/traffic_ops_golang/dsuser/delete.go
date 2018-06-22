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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
)

func Delete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"dsid", "userid"}, []string{"dsid", "userid"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	dsID := inf.IntParams["dsid"]
	userID := inf.IntParams["userid"]

	userErr, sysErr, errCode = tenant.CheckID(inf.Tx.Tx, inf.User, dsID)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	ok, err := deleteDSUser(inf.Tx.Tx, dsID, userID)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("deleting delivery service user: "+err.Error()))
		return
	}
	if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, nil)
		return
	}
	api.WriteRespAlert(w, r, tc.SuccessLevel, "User unlinked from delivery service.")
}

// deleteDSUser deletes the given deliveryservice_user. Returns whether the association existed, and any error.
func deleteDSUser(tx *sql.Tx, dsID int, userID int) (bool, error) {
	deletedUserID := 0
	if err := tx.QueryRow(`DELETE FROM deliveryservice_tmuser WHERE deliveryservice = $1 AND tm_user_id = $2 RETURNING tm_user_id`, dsID, userID).Scan(&deletedUserID); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.New("deleting delivery service user: " + err.Error())
	}
	return true, nil
}
