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

const TestExtPrivLevel = PrivLevelReadOnly;

type MyData struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Name1 string `json:"name1"`
}

type MyDataResponse struct {
	Response MyData `json:"response"`
}

func testExtHandler(db *sqlx.DB) RegexHandlerFunc {
        return func(w http.ResponseWriter, r *http.Request, p PathParams) {
                handleErr := func(err error, status int) {
                        log.Errorf("%v %v\n", r.RemoteAddr, err)
                        w.WriteHeader(status)
                        fmt.Fprintf(w, http.StatusText(status))
                }

                serverID, err := strconv.ParseInt(p["id"], 10, 64)
		if err != nil {
			handleErr(err, http.StatusBadRequest)
			//w.Header().Set("Content-Type", "application/json")
			//fmt.Fprintf(w, "{\"result\": \"Bad Integer\"}")
			return
		}
	
                resp, err := getTestDataJson(serverID, db)
                if err != nil {
                        handleErr(err, http.StatusInternalServerError)
			//w.Header().Set("Content-Type", "application/json")
			//fmt.Fprintf(w, "{\"result\": \"Error Querying Database\"}")
                        return
                }

                respBts, err := json.Marshal(resp)
                if err != nil {
                        handleErr(err, http.StatusInternalServerError)
			//w.Header().Set("Content-Type", "application/json")
			//fmt.Fprintf(w, "{\"result\": \"Error Marshalling Json\"}")
                        return
                }

                w.Header().Set("Content-Type", "application/json")
                fmt.Fprintf(w, "%s", respBts);
        }
}

func getTestData(db *sqlx.DB, serverID int64) (*MyData, error) {
	query := `SELECT
me.ID as id,
me.NAME as name,
me.NAME1 as name1
FROM my_data me
WHERE me.ID = $1`
	
	rows, err := db.Query(query, serverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	// TODO: what if there are zero rows?
	rows.Next()	
	var id sql.NullString
	var name sql.NullString
	var name1 sql.NullString
	
	if err := rows.Scan(&id, &name, &name1); err != nil {
		return nil, err
	}

	ret := MyData{
		ID: id.String,
		Name: name.String,
		Name1: name1.String,
	}
	return &ret, nil
}

func getTestDataJson(serverID int64, db * sqlx.DB) (*MyDataResponse, error) {
	myData, err := getTestData(db, serverID)
	if err != nil {
		return nil, fmt.Errorf("error getting my data: %v", err)
	}
	
	resp := MyDataResponse{ Response: *myData }
	return &resp, nil
}
