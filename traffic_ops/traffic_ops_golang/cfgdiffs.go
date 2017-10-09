package main;

import (
        "fmt"
        "net/http"
	"database/sql"
	"encoding/json"
	"strconv"
	
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
