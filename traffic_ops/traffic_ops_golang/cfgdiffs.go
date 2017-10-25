package main;

import (
        "fmt"
        "net/http"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"
	
	"github.com/apache/incubator-trafficcontrol/lib/go-log"

        "github.com/jmoiron/sqlx"
)

const CfgDiffsPrivLevel = PrivLevelReadOnly;
const CfgDiffsWritePrivLevel = PrivLevelOperations;

type CfgFileDiffs struct {
	FileName string `json:"fileName"`
	DBLinesMissing []string `json:"dbLinesMissing"`
	DiskLinesMissing []string `json:"diskLinesMissing"`
}

func getCfgDiffsHandler(db *sqlx.DB) RegexHandlerFunc {
    return func(w http.ResponseWriter, r *http.Request, p PathParams) {
                handleErr := func(err error, status int) {
                        log.Errorf("%v %v\n", r.RemoteAddr, err)
                        w.WriteHeader(status)
                        fmt.Fprintf(w, http.StatusText(status))
                }

		shortHostName:= p["short-host-name"]

		resp, err := getCfgDiffsJson(shortHostName, db)
		if err != nil {
				handleErr(err, http.StatusInternalServerError)
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

		shortHostName := p["short-host-name"]
		configName := p["cfg"]

		decoder := json.NewDecoder(r.Body)
		var diffs CfgFileDiffs
		err = decoder.Decode(&diffs)
		if err != nil {
			handleErr(err, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
	
				err = putCfgDiffs(db, shortHostName, configName, diffs)
				if err != nil {
						handleErr(err, http.StatusInternalServerError)
						return
		}
	}
}

func cfgDiffsExist(db *sqlx.DB, shortHostName string, configName string) (bool, error) {
	query := `SELECT 
EXISTS(
SELECT 1
FROM config_diffs me 
WHERE me.short_host_name = $1 AND me.config_name = $2)`

	rows, err := db.Query(query, shortHostName, configName)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		var exists sql.NullString
		
		err = rows.Scan(&exists)
		if err != nil {
			return false, err
		}

		if exists.String == "t" {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return false, nil;// this is an issue...
	}
}

func getCfgDiffs(db *sqlx.DB, shortHostName string) ([]CfgFileDiffs, error) {
	query := `SELECT
me.config_name as config_name,
array_to_json(me.db_lines_missing) as db_lines_missing,
array_to_json(me.disk_lines_missing) as disk_lines_missing,
me.last_checked as timestamp
FROM config_diffs me
WHERE me.short_host_name=$1`
	
	rows, err := db.Query(query, shortHostName)
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

		err := json.Unmarshal([]byte(db_lines_missing.String), &db_lines_missing_arr)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(disk_lines_missing.String), &disk_lines_missing_arr)

		configs = append(configs, CfgFileDiffs{
			FileName:    config_name.String,
			DBLinesMissing:     db_lines_missing_arr,
			DiskLinesMissing:  disk_lines_missing_arr,
		})
	}
	return configs, nil
}

func getCfgDiffsJson(shortHostName string, db * sqlx.DB) ([]CfgFileDiffs, error) {
	cfgDiffs, err := getCfgDiffs(db, shortHostName)
	if err != nil {
		return nil, fmt.Errorf("error getting my data: %v", err)
	}
	
	return cfgDiffs, nil
}

func insertCfgDiffs(db *sqlx.DB, shortHostName string, configName string, diffs CfgFileDiffs) ( error) {
	query := `INSERT INTO 
config_diffs(short_host_name, config_name, db_lines_missing, disk_lines_missing, last_checked)
VALUES($1, $2, (SELECT ARRAY(SELECT * FROM json_array_elements_text($3))), (SELECT ARRAY(SELECT * FROM json_array_elements_text($4))), $5)`
		
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
		shortHostName,
		configName, 
		dbLinesMissingJson,
		diskLinesMissingJson,
		time.Now().UTC())

	if err != nil {
		return err
	}
	defer rows.Close()
	
	return nil
}

func updateCfgDiffs(db *sqlx.DB, shortHostName string, configName string, diffs CfgFileDiffs) (bool, error) {
	query := `UPDATE config_diffs SET db_lines_missing=(SELECT ARRAY(SELECT * FROM json_array_elements_text($1))), 
disk_lines_missing=(SELECT ARRAY(SELECT * FROM json_array_elements_text($2))), last_checked=$3 WHERE short_host_name=$4 AND config_name=$5`
		
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
		time.Now().UTC(),
		shortHostName,
		configName)

	if err != nil {
		return false, err
	}
	
	count, err := rows.RowsAffected()
	if err != nil {
		return false, nil
	}

	if count > 0 {
		return true, nil
	}
	
	return false, nil

}

func putCfgDiffs(db *sqlx.DB, shortHostName string, configName string, diffs CfgFileDiffs) (error) {
	
	// Try updating the information first
	updated, err := updateCfgDiffs(db, shortHostName, configName, diffs)
	if err != nil {
		return err
	}
	if updated {
		return nil
	}
	return insertCfgDiffs(db, shortHostName, configName, diffs)
}