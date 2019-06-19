package hwinfo

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
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/jmoiron/sqlx"
)

func Get(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	api.RespWriter(w, r, inf.Tx.Tx)(getHWInfo(inf.Tx, inf.Params))
}

func getHWInfo(tx *sqlx.Tx, params map[string]string) ([]tc.HWInfo, error) {
	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToSQLCols := map[string]dbhelpers.WhereColumnInfo{
		"id":             dbhelpers.WhereColumnInfo{"h.id", api.IsInt},
		"serverHostName": dbhelpers.WhereColumnInfo{"s.host_name", nil},
		"serverId":       dbhelpers.WhereColumnInfo{"s.id", api.IsInt}, // TODO: this can be either s.id or h.serverid not sure what makes the most sense
		"description":    dbhelpers.WhereColumnInfo{"h.description", nil},
		"val":            dbhelpers.WhereColumnInfo{"h.val", nil},
		"lastUpdated":    dbhelpers.WhereColumnInfo{"h.last_updated", nil}, //TODO: this doesn't appear to work needs debugging
	}
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(params, queryParamsToSQLCols)
	if len(errs) > 0 {
		return nil, errors.New("getHWInfo building where clause: " + util.JoinErrsStr(errs))
	}
	query := selectHWInfoQuery() + where + orderBy + pagination
	log.Debugln("Query is ", query)

	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, errors.New("sqlx querying hwInfo: " + err.Error())
	}
	defer rows.Close()

	hwInfo := []tc.HWInfo{}
	for rows.Next() {
		s := tc.HWInfo{}
		if err = rows.StructScan(&s); err != nil {
			return nil, errors.New("sqlx scanning hwInfo: " + err.Error())
		}
		hwInfo = append(hwInfo, s)
	}
	return hwInfo, nil
}

func selectHWInfoQuery() string {

	query := `SELECT
	s.host_name as serverhostname,
    h.id,
    h.serverid,
    h.description,
    h.val,
    h.last_updated

FROM hwInfo h

JOIN server s ON s.id = h.serverid`
	return query
}
