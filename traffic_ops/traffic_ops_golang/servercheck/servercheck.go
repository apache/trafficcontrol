package servercheck

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

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

// CreateUpdateServercheck handles creating or updating an existing servercheck
func CreateUpdateServercheck(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	if inf.User.UserName != "extension" {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusForbidden, errors.New("invalid user for this API. Only the \"extension\" user can use this"), nil)
		return
	}

	serverCheckReq := tc.ServercheckRequestNullable{}

	if err := api.Parse(r.Body, inf.Tx.Tx, &serverCheckReq); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
		return
	}

	id, exists, err := getServerID(serverCheckReq.ID, serverCheckReq.HostName, inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting server id: "+err.Error()))
		return
	}
	if !exists {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("Server not found"), nil)
		return
	}

	col, exists, err := getColName(serverCheckReq.Name, inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting servercheck column name: "+err.Error()))
		return
	}
	if !exists {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("Server Check Extension %v not found - Do you need to install it?", *serverCheckReq.Name), nil)
		return
	}

	err = createUpdateServerCheck(id, col, *serverCheckReq.Value, inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("updating servercheck: "+err.Error()))
		return
	}

	successMsg := "Server Check was successfully updated"
	api.CreateChangeLogRawTx(api.ApiChange, successMsg, inf.User, inf.Tx.Tx)
	api.WriteRespAlert(w, r, tc.SuccessLevel, successMsg)
}

func getServerID(id *int, hostname *string, tx *sql.Tx) (int, bool, error) {
	if id != nil {
		_, exists, err := dbhelpers.GetServerNameFromID(tx, *id)
		return *id, exists, err
	}
	sID, exists, err := dbhelpers.GetServerIDFromName(*hostname, tx)
	return sID, exists, err
}

func getColName(shortName *string, tx *sql.Tx) (string, bool, error) {
	col := ""
	log.Infoln(*shortName)
	if err := tx.QueryRow(`SELECT servercheck_column_name FROM to_extension WHERE servercheck_short_name = $1`, *shortName).Scan(&col); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying servercheck column name: " + err.Error())
	}
	return col, true, nil
}

func createUpdateServerCheck(sid int, colName string, value int, tx *sql.Tx) error {

	insertUpdateQuery := fmt.Sprintf(`
INSERT INTO servercheck (server, %[1]s)
VALUES ($1, $2) 
ON CONFLICT (server) 
DO UPDATE SET %[1]s = EXCLUDED.%[1]s`, colName)

	result, err := tx.Exec(insertUpdateQuery, sid, value)
	if err != nil {
		return errors.New("insert server check: " + err.Error())
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("reading rows affected after server check insert : " + err.Error())
	}
	if rowsAffected == 0 {
		return errors.New("server check was inserted, no rows were affected : " + err.Error())
	}

	return nil
}
