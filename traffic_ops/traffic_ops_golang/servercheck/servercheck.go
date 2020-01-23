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

const serverInfoQuery = `
SELECT server.host_name AS hostName,
       server.id AS id,
       profile.name AS profile,
       status.name AS adminState,
       cachegroup.name AS cacheGroup,
       type.name AS type,
       server.upd_pending AS updPending,
       server.reval_pending AS revalPending
FROM server
LEFT JOIN profile ON server.profile = profile.id
LEFT JOIN status ON server.status = status.id
LEFT JOIN cachegroup ON server.cachegroup = cachegroup.id
LEFT JOIN type ON server.type = type.id
WHERE typename LIKE 'MID%' OR type.name LIKE 'EDGE%'
ORDER BY hostName ASCENDING
`

const serverChecksQuery = `
SELECT servercheck.id,
       servercheck.server,
       servercheck.aa,
       servercheck.ab,
       servercheck.ac,
       servercheck.ad,
       servercheck.ae,
       servercheck.af,
       servercheck.ag,
       servercheck.ah,
       servercheck.ai,
       servercheck.aj,
       servercheck.ak,
       servercheck.al,
       servercheck.am,
       servercheck.an,
       servercheck.ao,
       servercheck.ap,
       servercheck.aq,
       servercheck.ar,
       servercheck.bf,
       servercheck.at,
       servercheck.au,
       servercheck.av,
       servercheck.aw,
       servercheck.ax,
       servercheck.ay,
       servercheck.az,
       servercheck.ba,
       servercheck.bb,
       servercheck.bc,
       servercheck.bd,
       servercheck.be
FROM servercheck
`

const extensionsQuery = `
SELECT to_extension.servercheck_column_name,
       to_extension.servercheck_short_name
FROM to_extension
LEFT JOIN type ON type.id = to_extension.type
WHERE (type.name = 'CHECK_EXTENSION_BOOL' OR
      type.name = 'CHECK_EXTENSION_NUM') AND
      to_extension.servercheck_short_name IS NOT NULL AND
      to_extension.servercheck_column_name IS NOT NULL
`

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

// ReadServersChecks is the handler for GET requests for /servers/checks
func ReadServersChecks(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	extensions := make(map[string]string)
	extRows, err := tx.Query(extensionsQuery)
	if err != nil {
		sysErr = fmt.Errorf("querying for extensions: %v", err)
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}
	for extRows.Next() {
		var shortName string
		var columnName string
		if err = extRows.Scan(shortName, columnName); err != nil {
			sysErr = fmt.Errorf("scanning extension: %v", err)
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
			return
		}
		extensions[columnName] = shortName
	}

	colRows, err := inf.Tx.Queryx(serverChecksQuery)
	if err != nil {
		sysErr = fmt.Errorf("Querying server checks columns: %v", err)
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}

	columns := make(map[int]tc.ServerCheckColumns)
	for colRows.Next() {
		var cols tc.ServerCheckColumns
		if err = colRows.StructScan(cols); err != nil {
			sysErr = fmt.Errorf("scanning server checks columns: %v", err)
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
			return
		}

		columns[cols.Server] = cols
	}

	serverRows, err := tx.Query(serverInfoQuery)
	if err != nil {
		sysErr = fmt.Errorf("Querying server info for checks: %v", err)
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}

	data := []tc.GenericServerCheck{}
	for serverRows.Next() {
		var serverInfo tc.GenericServerCheck
		if err = serverRows.Scan(serverInfo.HostName, serverInfo.ID, serverInfo.Profile, serverInfo.AdminState, serverInfo.CacheGroup, serverInfo.Type, serverInfo.UpdPending, serverInfo.RevalPending); err != nil {
			sysErr = fmt.Errorf("scanning server info for checks: %v", err)
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
			return
		}

		serverCheckCols, ok := columns[serverInfo.ID]
		if ok {
			serverInfo.Checks = make(map[string]int)
		} else {
			data = append(data, serverInfo)
			continue
		}

		for colName, checkName := range extensions {
			switch (colName) {
			case "aa":
				if (serverCheckCols.AA != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AA
				}
			case "ab":
				if (serverCheckCols.AB != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AB
				}
			case "ac":
				if (serverCheckCols.AC != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AC
				}
			case "ad":
				if (serverCheckCols.AD != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AD
				}
			case "ae":
				if (serverCheckCols.AE != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AE
				}
			case "af":
				if (serverCheckCols.AF != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AF
				}
			case "ag":
				if (serverCheckCols.AG != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AG
				}
			case "ah":
				if (serverCheckCols.AH != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AH
				}
			case "ai":
				if (serverCheckCols.AI != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AI
				}
			case "aj":
				if (serverCheckCols.AJ != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AJ
				}
			case "ak":
				if (serverCheckCols.AK != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AK
				}
			case "al":
				if (serverCheckCols.AL != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AL
				}
			case "am":
				if (serverCheckCols.AM != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AM
				}
			case "an":
				if (serverCheckCols.AN != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AN
				}
			case "ao":
				if (serverCheckCols.AO != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AO
				}
			case "ap":
				if (serverCheckCols.AP != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AP
				}
			case "aq":
				if (serverCheckCols.AQ != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AQ
				}
			case "ar":
				if (serverCheckCols.AR != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AR
				}
			case "at":
				if (serverCheckCols.AT != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AT
				}
			case "au":
				if (serverCheckCols.AU != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AU
				}
			case "av":
				if (serverCheckCols.AV != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AV
				}
			case "aw":
				if (serverCheckCols.AW != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AW
				}
			case "ax":
				if (serverCheckCols.AX != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AX
				}
			case "ay":
				if (serverCheckCols.AY != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AY
				}
			case "az":
				if (serverCheckCols.AZ != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.AZ
				}
			case "ba":
				if (serverCheckCols.BA != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.BA
				}
			case "bb":
				if (serverCheckCols.BB != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.BB
				}
			case "bc":
				if (serverCheckCols.BC != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.BC
				}
			case "bd":
				if (serverCheckCols.BD != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.BD
				}
			case "be":
				if (serverCheckCols.BE != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.BE
				}
			case "bf":
				if (serverCheckCols.BF != nil) {
					serverInfo.Checks[checkName] = *serverCheckCols.BF
				}
			}
		}

		data = append(data, serverInfo)
	}

	api.WriteResp(w, r, data)
}
