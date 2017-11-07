package main;

import (
	"fmt"
    "net/http"
	"database/sql"
	"encoding/json"
	"errors"
	
	"github.com/apache/incubator-trafficcontrol/lib/go-log"

        "github.com/jmoiron/sqlx"
)

const CfgDiffsPrivLevel = PrivLevelReadOnly;
const CfgDiffsWritePrivLevel = PrivLevelOperations;

type CfgFileDiffs struct {
	FileName string `json:"fileName"`
	DBLinesMissing []string `json:"dbLinesMissing"`
	DiskLinesMissing []string `json:"diskLinesMissing"`
	ReportTimestamp string `json:"timestamp"`
}

type CfgFileDiffsResponse struct {
	Response []CfgFileDiffs `json:"response"`
}

func getCfgDiffsHandler(db *sqlx.DB) RegexHandlerFunc {
    return func(w http.ResponseWriter, r *http.Request, p PathParams) {
		handleErr := func(err error, status int) {
				log.Errorf("%v %v\n", r.RemoteAddr, err)
				w.WriteHeader(status)
				fmt.Fprintf(w, http.StatusText(status))
		}

		hostName:= p["host-name"]

		resp, err := getCfgDiffsJson(hostName, db)
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
		fmt.Fprintf(w, "%s", respBts);
	}
}

func putCfgDiffsHandler(db *sqlx.DB) AuthRegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p PathParams, username string, privLevel int) {
		handleErr := func(err error, status int) {
				log.Errorf("%v %v\n", r.RemoteAddr, err)
				w.WriteHeader(status)
				fmt.Fprintf(w, http.StatusText(status))
		}

		hostName := p["host-name"]
		configName := p["cfg-file-name"]

		decoder := json.NewDecoder(r.Body)
		var diffs CfgFileDiffs
		err := decoder.Decode(&diffs)
		if err != nil {
			handleErr(err, http.StatusBadRequest)
			return
		}

		defer r.Body.Close()
	
		result, err := putCfgDiffs(db, hostName, configName, diffs)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}
		
		// Not found (invalid hostname)
		if result == 0 { // This keeps happening
			w.WriteHeader(404)
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

func serverExists(db *sqlx.DB, hostName string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM server me WHERE me.host_name=$1)`
	rows, err := db.Query(query, hostName)
	if err != nil {
		return false, err
	}

	defer rows.Close()
	
	for rows.Next() {
		var exists sql.NullString
		
		err = rows.Scan(&exists)
		if err != nil {
			return false, err
		}

		log.Infof(exists.String)

		if exists.String == "true" {
			return true, nil
		} else {
			return false, nil
		}
		break
	} //else {
		return false, errors.New("Failed to load row!") // What does this mean?
	//}
}

func getCfgDiffs(db *sqlx.DB, hostName string) ([]CfgFileDiffs, error) {
	query := `SELECT
me.config_name as config_name,
array_to_json(me.db_lines_missing) as db_lines_missing,
array_to_json(me.disk_lines_missing) as disk_lines_missing,
me.last_checked as timestamp
FROM config_diffs me
WHERE me.server_id=(SELECT server.id FROM server WHERE host_name=$1)`
	
	rows, err := db.Query(query, hostName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	configs := []CfgFileDiffs{}
	
	// TODO: what if there are zero rows?
	for rows.Next()	{
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
			FileName:    config_name.String,
			DBLinesMissing:     db_lines_missing_arr,
			DiskLinesMissing:  disk_lines_missing_arr,
			ReportTimestamp: timestamp.String,
		})
	}
	return configs, nil
}

func getCfgDiffsJson(hostName string, db * sqlx.DB) (*CfgFileDiffsResponse, error) {
	cfgDiffs, err := getCfgDiffs(db, hostName)
	if err != nil {
		return nil, fmt.Errorf("error getting my data: %v", err)
	}

	response := CfgFileDiffsResponse{
		Response: cfgDiffs,
	}
	
	return &response, nil
}

func insertCfgDiffs(db *sqlx.DB, hostName string, configName string, diffs CfgFileDiffs) ( error) {
	query := `INSERT INTO 
config_diffs(server_id, config_name, db_lines_missing, disk_lines_missing, last_checked)
VALUES((SELECT server.id FROM server WHERE host_name=$1), $2, (SELECT ARRAY(SELECT * FROM json_array_elements_text($3))), (SELECT ARRAY(SELECT * FROM json_array_elements_text($4))), $5)`
		
	dbLinesMissingJson, err := json.Marshal(diffs.DBLinesMissing)
	if err != nil {
		return err
	}
	diskLinesMissingJson, err := json.Marshal(diffs.DiskLinesMissing)
	if err != nil {
		return err
	}

	//NOTE: if the serverID doesn't match a server, this error will appear like a 500-type error
	rows, err := db.Query(query,
		hostName,
		configName, 
		dbLinesMissingJson,
		diskLinesMissingJson,
		diffs.ReportTimestamp)

	if err != nil {
		return err
	}
	defer rows.Close()
	
	return nil
}

func updateCfgDiffs(db *sqlx.DB, hostName string, configName string, diffs CfgFileDiffs) (bool, error) {
	query := `UPDATE config_diffs SET db_lines_missing=(SELECT ARRAY(SELECT * FROM json_array_elements_text($1))), 
disk_lines_missing=(SELECT ARRAY(SELECT * FROM json_array_elements_text($2))), last_checked=$3 WHERE server_id=(SELECT server.id FROM server WHERE host_name=$4) AND config_name=$5`
		
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
		hostName,
		configName)

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

func putCfgDiffs(db *sqlx.DB, hostName string, configName string, diffs CfgFileDiffs) (int, error) {
	
	sExists, err := serverExists(db, hostName)
	if err != nil {
		return -1, err
	}
	if sExists == false {
		return 0, nil
	}

	// Try updating the information first
	updated, err := updateCfgDiffs(db, hostName, configName, diffs)
	if err != nil {
		return -1, err
	}
	if updated {
		return 2, nil
	}
	return 1, insertCfgDiffs(db, hostName, configName, diffs)
}