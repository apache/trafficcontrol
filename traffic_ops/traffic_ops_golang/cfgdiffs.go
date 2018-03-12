package main

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

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"

	"github.com/jmoiron/sqlx"
)

type CfgFileDiffs struct {
	FileName         string   `json:"fileName"`
	DBLinesMissing   []string `json:"dbLinesMissing"`
	DiskLinesMissing []string `json:"diskLinesMissing"`
	ReportTimestamp  string   `json:"timestamp"`
}

type cfgFileDiffsResponse struct {
	Response []CfgFileDiffs `json:"response"`
}

type updateCfgDiffsMethod func(db *sqlx.DB, serverId int64, diffs CfgFileDiffs) (bool, error)
type insertCfgDiffsMethod func(db *sqlx.DB, serverId int64, diffs CfgFileDiffs) error
type getCfgDiffsMethod func(db *sqlx.DB, serverId int64) ([]CfgFileDiffs, error)

func getCfgDiffsHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErr := func(err error, status int) {
			log.Errorf("%v %v\n", r.RemoteAddr, err)
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
		}
		
		pathParams, err := api.GetCombinedParams(r)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		hostName := pathParams["host-name"]
		domainName := pathParams["domain-name"]

		serverId, err := getServerId(db, hostName, domainName)

		// server not found
		if err == sql.ErrNoRows {
			handleErr(err, http.StatusNotFound)
			return
		}

		// some other error (querying/scanning)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}


		resp, err := getCfgDiffsJson(db, serverId, getCfgDiffs)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		// if the response has a length of zero, no results were found for that server
		if len(resp.Response) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", respBts)
	}
}

func putCfgDiffsHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErr := func(err error, status int) {
			log.Errorf("%v %v\n", r.RemoteAddr, err)
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
		}
		
		pathParams, err := api.GetCombinedParams(r)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		hostName := pathParams["host-name"]
		domainName := pathParams["domain-name"]
		configName := pathParams["cfg-file-name"]

		decoder := json.NewDecoder(r.Body)
		var diffs CfgFileDiffs
		err = decoder.Decode(&diffs)
		if err != nil {
			handleErr(err, http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		diffs.FileName = configName

		serverId, err := getServerId(db, hostName, domainName)

		// server not found
		if err == sql.ErrNoRows {
			handleErr(err, http.StatusNotFound)
			return
		}

		// some other error (querying/scanning)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		result, err := putCfgDiffs(db, serverId, diffs, updateCfgDiffs, insertCfgDiffs)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		// Created (newly added)
		if result == 1 {
			w.WriteHeader(201)
			return
		}
		// Updated (already existed)
		if result == 2 {
			w.WriteHeader(202)
			return
		}
	}
}

func getServerId(db *sqlx.DB, hostName string, domainName string) (int64, error) {
	query := `SELECT id FROM server me WHERE me.host_name=$1 AND me.domain_name=$2`

	var id sql.NullInt64
	err := db.QueryRow(query, hostName, domainName).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id.Int64, nil
}

func getCfgDiffs(db *sqlx.DB, serverId int64) ([]CfgFileDiffs, error) {
	query := `SELECT
me.config_name as config_name,
array_to_json(me.db_lines_missing) as db_lines_missing,
array_to_json(me.disk_lines_missing) as disk_lines_missing,
me.last_checked as timestamp
FROM config_diffs me
WHERE me.server=$1`

	rows, err := db.Query(query, serverId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	configs := []CfgFileDiffs{}

	// TODO: what if there are zero rows?
	for rows.Next() {
		var config_name sql.NullString
		var db_lines_missing sql.NullString
		var disk_lines_missing sql.NullString
		var timestamp sql.NullString

		var db_lines_missing_arr []string
		var disk_lines_missing_arr []string

		if err := rows.Scan(&config_name, &db_lines_missing, &disk_lines_missing, &timestamp); err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(db_lines_missing.String), &db_lines_missing_arr)
		if err != nil {
			return nil, err
		}

		err := json.Unmarshal([]byte(disk_lines_missing.String), &disk_lines_missing_arr)
		if err != nil {
			return nil, err
		}

		configs = append(configs, CfgFileDiffs{
			FileName:         config_name.String,
			DBLinesMissing:   db_lines_missing_arr,
			DiskLinesMissing: disk_lines_missing_arr,
			ReportTimestamp:  timestamp.String,
		})
	}
	return configs, nil
}

func getCfgDiffsJson(db *sqlx.DB, serverId int64, getMethod getCfgDiffsMethod) (*cfgFileDiffsResponse, error) {
	cfgDiffs, err := getMethod(db, serverId)
	if err != nil {
		return nil, fmt.Errorf("error getting my data: %v", err)
	}

	response := cfgFileDiffsResponse{
		Response: cfgDiffs,
	}

	return &response, nil
}

func insertCfgDiffs(db *sqlx.DB, serverId int64, diffs CfgFileDiffs) error {
	query := `INSERT INTO 
config_diffs(server, config_name, db_lines_missing, disk_lines_missing, last_checked)
VALUES($1, $2, (SELECT ARRAY(SELECT * FROM json_array_elements_text($3))), (SELECT ARRAY(SELECT * FROM json_array_elements_text($4))), $5)`

	dbLinesMissingJson, err := json.Marshal(diffs.DBLinesMissing)
	if err != nil {
		return err
	}
	diskLinesMissingJson, err := json.Marshal(diffs.DiskLinesMissing)
	if err != nil {
		return err
	}

	_, err = db.Exec(query,
		serverId,
		diffs.FileName,
		dbLinesMissingJson,
		diskLinesMissingJson,
		diffs.ReportTimestamp)

	if err != nil {
		return err
	}

	return nil
}

func updateCfgDiffs(db *sqlx.DB, serverId int64, diffs CfgFileDiffs) (bool, error) {
	query := `UPDATE config_diffs 
SET db_lines_missing=(SELECT ARRAY(SELECT * FROM json_array_elements_text($1))), 
disk_lines_missing=(SELECT ARRAY(SELECT * FROM json_array_elements_text($2))), 
last_checked=$3 
WHERE server=$4 AND config_name=$5`

	dbLinesMissingJson, err := json.Marshal(diffs.DBLinesMissing)
	if err != nil {
		return false, err
	}
	diskLinesMissingJson, err := json.Marshal(diffs.DiskLinesMissing)
	if err != nil {
		return false, err
	}

	rows, err := db.Exec(query,
		dbLinesMissingJson,
		diskLinesMissingJson,
		diffs.ReportTimestamp,
		serverId,
		diffs.FileName)

	if err != nil {
		return false, err
	}

	count, err := rows.RowsAffected()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil

}

func putCfgDiffs(db *sqlx.DB, serverId int64, diffs CfgFileDiffs, updateMethod updateCfgDiffsMethod, insertMethod insertCfgDiffsMethod) (int, error) {

	// Try updating the information first
	updated, err := updateMethod(db, serverId, diffs)
	if err != nil {
		return 2, err
	}
	if updated {
		return 2, nil
	}
	return 1, insertMethod(db, serverId, diffs)
}
