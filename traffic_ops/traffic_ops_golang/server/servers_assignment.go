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
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/lib/pq"
)

func AssignDeliveryServicesToServerHandler(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	dsList := []int{}
	if err := json.NewDecoder(r.Body).Decode(&dsList); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}

	replaceQueryParameter := inf.Params["replace"]
	replace, err := strconv.ParseBool(replaceQueryParameter) //accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False. for replace url parameter documentation
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
		return
	}

	serverPathParameter := inf.Params["id"]
	server, err := strconv.Atoi(serverPathParameter)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
		return
	}
	serverName, ok, err := dbhelpers.GetServerNameFromID(inf.Tx.Tx, server)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting server name from ID: "+err.Error()))
		return
	} else if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("no server with that ID found"), nil)
		return
	}

	usrErr, sysErr, status := ValidateDSCapabilities(dsList, serverName, inf.Tx.Tx)
	if usrErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, status, usrErr, sysErr)
		return
	}

	assignedDSes, err := assignDeliveryServicesToServer(server, dsList, replace, inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}

	api.CreateChangeLogRawTx(api.ApiChange, "SERVER: "+serverName+", ID: "+strconv.Itoa(server)+", ACTION: Assigned "+strconv.Itoa(len(assignedDSes))+" DSes to server", inf.User, inf.Tx.Tx)
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "successfully assigned dses to server", tc.AssignedDsResponse{server, assignedDSes, replace})
}

// ValidateDSCapabilities checks that the server meets the requirements of each delivery service to be assigned.
func ValidateDSCapabilities(dsIDs []int, serverName string, tx *sql.Tx) (error, error, int) {
	var dsCaps []string
	sCaps, err := dbhelpers.GetServerCapabilitiesFromName(serverName, tx)

	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	for _, id := range dsIDs {
		dsCaps, err = dbhelpers.GetDSRequiredCapabilitiesFromID(id, tx)
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
		for _, dsc := range dsCaps {
			if !util.ContainsStr(sCaps, dsc) {
				return errors.New(fmt.Sprintf("Caching server cannot assign this delivery service without having the required delivery service capabilities: [%v] for server %s", dsCaps, serverName)), nil, http.StatusBadRequest
			}
		}
	}

	return nil, nil, 0
}

func assignDeliveryServicesToServer(server int, dses []int, replace bool, tx *sql.Tx) ([]int, error) {
	if replace {
		//delete currently assigned dses from server
		if _, err := tx.Exec(`DELETE FROM deliveryservice_server WHERE server = $1`, server); err != nil {
			return nil, errors.New("could not delete old deliveryservice_server associations for server: " + err.Error())
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
	q := `
INSERT INTO deliveryservice_server (deliveryservice, server)
	WITH
	q1 AS (SELECT UNNEST($1::bigint[])),
	q2 AS ( SELECT * FROM (VALUES ($2::bigint)) AS server )
	SELECT * FROM q1,q2 ON CONFLICT DO NOTHING
`
	if _, err := tx.Exec(q, dsPqArray, server); err != nil {
		return nil, errors.New("inserting deliveryservice_server: " + err.Error())
	}

	//need remap config location
	var atsConfigLocation string
	if err := tx.QueryRow("SELECT value FROM parameter WHERE name = 'location' AND config_file = '" + atscfg.RemapFile + "'").Scan(&atsConfigLocation); err != nil {
		return nil, errors.New("scanning location parameter: " + err.Error())
	}
	if strings.HasSuffix(atsConfigLocation, "/") {
		atsConfigLocation = atsConfigLocation[:len(atsConfigLocation)-1]
	}

	//we need dses: xmlids and edge_header_rewrite, regex_remap, and cache_url
	rows, err := tx.Query(`SELECT xml_id, edge_header_rewrite, regex_remap, cacheurl FROM deliveryservice WHERE id = ANY($1::bigint[])`, dsPqArray)
	if err != nil {
		return nil, errors.New("querying deliveryservice: " + err.Error())
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
			return nil, errors.New("scanning deliveryservice: " + err.Error())
		}
		if xmlID.Valid && len(xmlID.String) > 0 {
			//param := "hdr_rw_" + xmlID.String + ".config"
			param := atscfg.GetConfigFile(atscfg.HeaderRewritePrefix, xmlID.String)
			if edgeHeaderRewrite.Valid && len(edgeHeaderRewrite.String) > 0 {
				insert = append(insert, param)
			} else {
				delete = append(delete, param)
			}
			param = atscfg.GetConfigFile(atscfg.RegexRemapPrefix, xmlID.String)
			if regexRemap.Valid && len(regexRemap.String) > 0 {
				insert = append(insert, param)
			} else {
				delete = append(delete, param)
			}
			param = atscfg.GetConfigFile(atscfg.CacheUrlPrefix, xmlID.String)
			if cacheURL.Valid && len(cacheURL.String) > 0 {
				insert = append(insert, param)
			} else {
				delete = append(delete, param)
			}
		}

	}

	//insert the parameters we selected above:
	q = `
INSERT INTO parameter (config_file, name, value)
	WITH
	q1 AS (SELECT UNNEST($1::text[]) AS config_file),
	q2 AS (SELECT * FROM (VALUES ($2) ) AS name),
	q3 AS (SELECT * FROM (VALUES ($3) ) AS value)
	 SELECT * FROM q1,q2,q3 ON CONFLICT DO NOTHING
`
	fileNamePqArray := pq.Array(insert)
	if _, err = tx.Exec(q, fileNamePqArray, "location", atsConfigLocation); err != nil {
		return nil, errors.New("inserting parameters: " + err.Error())
	}

	//select the ids associated with the parameters we created above (may be able to get them from insert above to optimize)
	rows, err = tx.Query(`SELECT id FROM parameter WHERE name = 'location' AND config_file IN ($1)`, fileNamePqArray)
	if err != nil {
		return nil, errors.New("selecting location parameter after insert: " + err.Error())
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
	q = `
INSERT INTO profile_parameter (profile, parameter)
	WITH
	q1 AS ( SELECT DISTINCT profile FROM server LEFT JOIN deliveryservice_server ON server.id = deliveryservice_server.server WHERE deliveryservice_server.deliveryservice = ANY($1::bigint[]) ),
	q2 AS (SELECT UNNEST($2::bigint[]) AS parameter)
	SELECT * FROM q1,q2
	ON CONFLICT DO NOTHING
`
	if _, err = tx.Exec(q, dsPqArray, pq.Array(parameterIds)); err != nil {
		return nil, errors.New("inserting profile_parameter: " + err.Error())
	}

	//process delete list
	if _, err = tx.Exec(`DELETE FROM parameter WHERE name = 'location' AND config_file = ANY($1)`, pq.Array(delete)); err != nil {
		return nil, errors.New("deleting parameters: " + err.Error())
	}

	return dses, nil
}
