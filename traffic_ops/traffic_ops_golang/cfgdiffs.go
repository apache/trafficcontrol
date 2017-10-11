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
	ServerID int64 `json:"serverId"`
	FileName string `json:"FileName"`
	DBLines []string `json:"dbLines"`
	LocalLines []string `json:"localLines"`
}

func cfgDiffsHandler(db *sqlx.DB) RegexHandlerFunc {
        return func(w http.ResponseWriter, r *http.Request, p PathParams) {
                handleErr := func(err error, status int) {
                        log.Errorf("%v %v\n", r.RemoteAddr, err)
                        w.WriteHeader(status)
                        fmt.Fprintf(w, http.StatusText(status))
                }

                serverID, err := strconv.ParseInt(p["id"], 10, 64)
		if err != nil {
			handleErr(err, http.StatusBadRequest)
			return
		}
	
                resp, err := getCfgDiffsJson(serverID, db)
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

func postCfgDiffsHandler(db *sqlx.DB) AuthRegexHandlerFunc {
        return func(w http.ResponseWriter, r *http.Request, p PathParams, username string, privLevel int) {
                handleErr := func(err error, status int) {
                        log.Errorf("%v %v\n", r.RemoteAddr, err)
                        w.WriteHeader(status)
                        fmt.Fprintf(w, http.StatusText(status))
                }

                serverID, err := strconv.ParseInt(p["id"], 10, 64)
		if err != nil {
			handleErr(err, http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var diffs CfgFileDiffs
		err = decoder.Decode(&diffs)
		if err != nil {
			handleErr(err, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
	
                err = postCfgDiffs(db, serverID, diffs)
                if err != nil {
                        handleErr(err, http.StatusInternalServerError)
                        return
                }
        }
}

func cfgDiffsExist(db *sqlx.DB, serverID int64, configName string) (bool, error) {
	query := `SELECT 
EXISTS(
SELECT 1
FROM config_diffs me 
WHERE me.server_id = $1 AND me.config_name = $2)`

	rows, err := db.Query(query, serverID, configName)
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

func getCfgDiffs(db *sqlx.DB, serverID int64) ([]CfgFileDiffs, error) {
	query := `SELECT
me.config_name as config_name,
array_to_json(me.db_lines) as db_lines,
array_to_json(me.local_lines) as local_lines,
me.last_checked as timestamp
FROM config_diffs me
WHERE me.server_id = $1`
	
	rows, err := db.Query(query, serverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	configs := []CfgFileDiffs{}
	
	// TODO: what if there are zero rows?
	for rows.Next()	{
		var config_name sql.NullString
		var db_lines sql.NullString
		var local_lines sql.NullString
		var timestamp sql.NullString
		
		var db_lines_arr []string
		var local_lines_arr []string

		if err := rows.Scan(&config_name, &db_lines, &local_lines, &timestamp); err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(db_lines.String), &db_lines_arr)
		json.Unmarshal([]byte(local_lines.String), &local_lines_arr)

		configs = append(configs, CfgFileDiffs{
			ServerID:    serverID,
			FileName:    config_name.String,
			DBLines:     db_lines_arr,
			LocalLines:  local_lines_arr,
		})
	}
	return configs, nil
}

func getCfgDiffsJson(serverID int64, db * sqlx.DB) ([]CfgFileDiffs, error) {
	cfgDiffs, err := getCfgDiffs(db, serverID)
	if err != nil {
		return nil, fmt.Errorf("error getting my data: %v", err)
	}
	
	return cfgDiffs, nil
}

func postCfgDiffs(db *sqlx.DB, serverID int64, diffs CfgFileDiffs) (error) {
	query := `INSERT INTO 
config_diffs(server_id, config_name, db_lines, local_lines, last_checked)
VALUES($1, $2, json_array_elements_text($3), json_array_elements_text($4), $5)`
		
	dbLinesJson, err := json.Marshal(diffs.DBLines)
	if err != nil {
		return err
	}
	localLinesJson, err := json.Marshal(diffs.LocalLines)
	if err != nil {
		return err
	}

	rows, err := db.Query(query, 
		serverID, 
		diffs.FileName, 
		dbLinesJson,
		localLinesJson,
		time.Now().UTC())

	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}