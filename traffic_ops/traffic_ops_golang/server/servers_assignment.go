package server

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
	"strconv"
	"strings"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/ats"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func AssignDeliveryServicesToServerHandler(db *sqlx.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		params, err := api.GetCombinedParams(r)
		if err != nil {
			log.Errorf("unable to get parameters from request: %s", err)
			handleErrs(http.StatusInternalServerError, err)
		}

		var dsList []int

		err = json.NewDecoder(r.Body).Decode(&dsList)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		replaceQueryParameter := params["replace"]
		replace, err := strconv.ParseBool(replaceQueryParameter) //accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False. for replace url parameter documentation
		if err != nil {
			handleErrs(http.StatusBadRequest, err)
			return
		}

		serverPathParameter := params["id"]
		server, err := strconv.Atoi(serverPathParameter)
		if err != nil {
			handleErrs(http.StatusBadRequest, err)
			return
		}

		assignedDSes, err := assignDeliveryServicesToServer(server, dsList, replace, db)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		type AssignedDsResponse struct {
			ServerID int   `json:"serverId"`
			DSIds    []int `json:"dsIds"`
			Replace  bool  `json:"replace"`
		}

		assignResp := AssignedDsResponse{server, assignedDSes, replace}

		resp := struct {
			Response AssignedDsResponse `json:"response"`
			tc.Alerts
		}{assignResp, tc.CreateAlerts(tc.SuccessLevel, "successfully assigned dses to server")}
		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		fmt.Fprintf(w, "%s", respBts)
	}
}

func assignDeliveryServicesToServer(server int, dses []int, replace bool, db *sqlx.DB) ([]int, error) {
	//transaction rollback in this functions requires err to be set to the proper error or nil before returning
	tx, err := db.Beginx()
	defer func() {
		if tx == nil {
			return
		}
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	if err != nil {
		log.Error.Printf("could not begin transaction: %v\n", err)
		return nil, tc.DBError
	}

	if replace {
		//delete currently assigned dses from server
		deleteCurrent, err := tx.Prepare("DELETE FROM deliveryservice_server WHERE server = $1")
		if err != nil {
			log.Error.Printf("could not prepare deliveryservice_server delete statement: %s\n", err)
			return nil, tc.DBError
		}
		_, err = deleteCurrent.Exec(server)
		if err != nil {
			log.Error.Printf("could not delete old deliveryservice_server associations for server: %s\n", err)
			return nil, tc.DBError
		}
	}

	//assign new dses
	dsPqArray := pq.Array(dses)
	// The common table expressions (CTEs) used below allow for inserting every deliveryService in dses with the server without a loop
	// the result of the select (for server = 100, dses = [1,2,3]) used by the insert essentially looks like:
	//          | server
	//    ---------------
	//      1   |   100
	//      2   |   100
	//      3   |   100
	// UNNEST is used to turn the array of ds values into, essentially, * FROM ( VALUES (1),(2),(3) )
	// this allows for a single insert query instead of a loop over the dses.
	// This pattern is used for the other bulk inserts as well.
	bulkInsert, err := tx.Prepare(`INSERT INTO deliveryservice_server (deliveryservice, server)
	WITH
	q1 AS (SELECT UNNEST($1::bigint[])),
	q2 AS ( SELECT * FROM (VALUES ($2::bigint)) AS server )
	SELECT * FROM q1,q2 ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Error.Printf("could not prepare deliveryservice_server bulk insert: %s\n", err)
		return nil, tc.DBError
	}
	_, err = bulkInsert.Exec(dsPqArray, server)
	if err != nil {
		log.Error.Printf("could not execute deliveryservice_server bulk insert: %s\n", err)
		return nil, tc.DBError
	}

	//need remap config location
	row := tx.QueryRow("SELECT value FROM parameter WHERE name = 'location' AND config_file = '" + ats.RemapFile + "'")
	var atsConfigLocation string
	row.Scan(&atsConfigLocation)
	if strings.HasSuffix(atsConfigLocation, "/") {
		atsConfigLocation = atsConfigLocation[:len(atsConfigLocation)-1]
	}

	//we need dses: xmlids and edge_header_rewrite, regex_remap, and cache_url
	selectDsFieldsQuery, err := tx.Prepare("SELECT xml_id, edge_header_rewrite, regex_remap, cacheurl FROM deliveryservice WHERE id = ANY($1::bigint[])")
	if err != nil {
		log.Error.Printf("could not prepare ds fields query: %s\n", err)
		return nil, tc.DBError
	}
	rows, err := selectDsFieldsQuery.Query(dsPqArray)
	if err != nil {
		log.Error.Printf("could not execute ds fields select query: %s\n", err)
		return nil, tc.DBError
	}
	defer rows.Close()

	//create new parameters here as necessary:
	//loop over ds results and build file parameters we need to insert / select
	//for all of: header rewrite, regex_remap, cache_url
	// if ds has it add parameter to insert list
	// other wise add to delete list
	//TODO: DylanVolz this may need to be extended refactored, there are potentially other parameters that need this like urlSigKeys...
	insert := []string{}
	delete := []string{}
	for rows.Next() {
		var xmlID sql.NullString
		var edgeHeaderRewrite sql.NullString
		var regexRemap sql.NullString
		var cacheURL sql.NullString

		if err := rows.Scan(&xmlID, &edgeHeaderRewrite, &regexRemap, &cacheURL); err != nil {
			log.Error.Printf("could not scan ds fields row: %s\n", err)
			return nil, tc.DBError
		}
		if xmlID.Valid && len(xmlID.String) > 0 {
			//param := "hdr_rw_" + xmlID.String + ".config"
			param := ats.GetConfigFile(ats.HeaderRewritePrefix, xmlID.String)
			if edgeHeaderRewrite.Valid && len(edgeHeaderRewrite.String) > 0 {
				insert = append(insert, param)
			} else {
				delete = append(delete, param)
			}
			param = ats.GetConfigFile(ats.RegexRemapPrefix, xmlID.String)
			if regexRemap.Valid && len(regexRemap.String) > 0 {
				insert = append(insert, param)
			} else {
				delete = append(delete, param)
			}
			param = ats.GetConfigFile(ats.CacheUrlPrefix, xmlID.String)
			if cacheURL.Valid && len(cacheURL.String) > 0 {
				insert = append(insert, param)
			} else {
				delete = append(delete, param)
			}
		}

	}

	//insert the parameters we selected above:
	insertParams, err := tx.Prepare(`INSERT INTO parameter (config_file, name, value)
	WITH
	q1 AS (SELECT UNNEST($1::text[]) AS config_file),
	q2 AS (SELECT * FROM (VALUES ($2) ) AS name),
	q3 AS (SELECT * FROM (VALUES ($3) ) AS value)
	 SELECT * FROM q1,q2,q3 ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Error.Printf("could not prepare parameter bulk insert query: %s\n", err)
		return nil, tc.DBError
	}
	fileNamePqArray := pq.Array(insert)
	if _, err = insertParams.Exec(fileNamePqArray, "location", atsConfigLocation); err != nil {
		log.Error.Printf("could not execute parameter bulk insert: %s\n", err)
		return nil, tc.DBError
	}

	//select the ids associated with the parameters we created above (may be able to get them from insert above to optimize)
	selectParameterIds, err := tx.Prepare("SELECT id FROM parameter WHERE name = 'location' AND config_file IN ($1)")
	if err != nil {
		log.Error.Printf("could not prepare parameter id select query: %s\n", err)
		return nil, tc.DBError
	}
	rows, err = selectParameterIds.Query(fileNamePqArray)
	if err != nil {
		log.Error.Printf("could not execute parameter id select query: %s\n", err)
		return nil, tc.DBError
	}
	defer rows.Close()

	parameterIds := []int64{}
	for rows.Next() {
		var ID int64
		if err := rows.Scan(&ID); err != nil {
			log.Error.Printf("could not scan parameter ID: %s\n", err)
			return nil, tc.DBError
		}
		parameterIds = append(parameterIds, ID)
	}

	//associate all parameter ids with the profiles associated with all servers associated with assigned dses.
	insertProfileParams, err := tx.Prepare(`INSERT INTO profile_parameter (profile, parameter)
	WITH
	q1 AS ( SELECT DISTINCT profile FROM server LEFT JOIN deliveryservice_server ON server.id = deliveryservice_server.server WHERE deliveryservice_server.deliveryservice = ANY($1::bigint[]) ),
	q2 AS (SELECT UNNEST($2::bigint[]) AS parameter)
	SELECT * FROM q1,q2
	ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Error.Printf("could not prepare profile_parameter bulk insert: %s\n", err)
		return nil, tc.DBError
	}
	if _, err = insertProfileParams.Exec(dsPqArray, pq.Array(parameterIds)); err != nil {
		log.Error.Printf("could not execute profile_parameter bulk insert: %s\n", err)
		return nil, tc.DBError
	}

	//process delete list
	deleteTx, err := tx.Prepare(`DELETE FROM parameter WHERE name = 'location' AND config_file = ANY($1)`)
	if err != nil {
		log.Error.Printf("could not prepare parameter delete query: %s\n", err)
		return nil, tc.DBError
	}
	if _, err = deleteTx.Exec(pq.Array(delete)); err != nil {
		log.Error.Printf("could not execute parameter delete query: %s\n", err)
		return nil, tc.DBError
	}

	return dses, nil
}
