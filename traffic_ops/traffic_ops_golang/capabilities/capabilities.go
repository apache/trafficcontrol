package capabilities

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

import "database/sql"
import "encoding/json"
import "errors"
import "fmt"
import "net/http"

import "github.com/apache/trafficcontrol/lib/go-tc"
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"

const readQuery = `
SELECT description,
       last_updated,
       name
FROM capability
`

const createQuery = `
INSERT INTO capability (name, description)
VALUES ($1, $2)
RETURNING description, last_updated, name
`

const replaceQuery = `
UPDATE capability
SET name=$1, description=$2
WHERE name=$3
RETURNING description, last_updated, name
`

const deleteQuery = `
DELETE FROM capability
WHERE name=$1
RETURNING description, last_updated, name
`

func Read(w http.ResponseWriter, r *http.Request) {
	inf, sysErr, userErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()


	var rows *sql.Rows
	var err error
	if name, ok := inf.Params["name"]; ok {
		rows, err = tx.Query(readQuery + "WHERE name=$1", name)
	} else {
		rows, err = tx.Query(readQuery)
	}
	if err != nil && err != sql.ErrNoRows {
		errCode = http.StatusInternalServerError
		sysErr = fmt.Errorf("querying capabilities: %v", err)
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}
	defer rows.Close()

	caps := []tc.Capability{}
	for rows.Next() {
		cap := tc.Capability{}
		if err := rows.Scan(&cap.Description, &cap.LastUpdated, &cap.Name); err != nil {
			errCode = http.StatusInternalServerError
			sysErr = fmt.Errorf("Parsing database response: %v", err)
			api.HandleErr(w, r, tx, errCode, nil, sysErr)
			return
		}

		caps = append(caps, cap)
	}

	api.WriteResp(w, r, caps)
}

func Create(w http.ResponseWriter, r *http.Request) {
	inf, sysErr, userErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	decoder := json.NewDecoder(r.Body)
	var cap tc.Capability
	if err := decoder.Decode(&cap); err != nil {
		sysErr = fmt.Errorf("Decoding request body: %v")
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}

	if cap.Name == "" {
		userErr = errors.New("'name' must be defined! (and not empty)")
		errCode = http.StatusBadRequest
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	if cap.Description == "" {
		userErr = errors.New("'description' must be defined! (and not empty)")
		errCode = http.StatusBadRequest
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	if ok, err := capabilityNameExists(cap.Name, tx); err != nil {
		sysErr = fmt.Errorf("Checking for capability %s's existence: %v", cap.Name, err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	} else if ok {
		userErr = fmt.Errorf("Capability '%s' already exists!", cap.Name)
		errCode = http.StatusConflict
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	row := tx.QueryRow(createQuery, cap.Name, cap.Description)
	if err := row.Scan(&cap.Description, &cap.LastUpdated, &cap.Name); err != nil {
		sysErr = fmt.Errorf("Inserting capability: %v", err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}

	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Capability created.", cap)
	api.CreateChangeLogRawTx(api.ApiChange, fmt.Sprintf("CAPABILITY: %s, ACTION: Created", cap.Name), inf.User, tx)
}

func Replace(w http.ResponseWriter, r *http.Request) {
	inf, sysErr, userErr, errCode := api.NewInfo(r, []string{"name"}, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	decoder := json.NewDecoder(r.Body)
	var cap tc.Capability
	if err := decoder.Decode(&cap); err != nil {
		userErr = fmt.Errorf("Couldn't parse capacity: %v", err)
		errCode = http.StatusBadRequest
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	if cap.Name != inf.Params["name"] {
		if ok, err := capabilityNameExists(cap.Name, tx); err != nil {
			sysErr = fmt.Errorf("Checking for capability %s's existence: %v", cap.Name, err)
			errCode = http.StatusInternalServerError
			api.HandleErr(w, r, tx, errCode, nil, sysErr)
			return
		} else if ok {
			errCode = http.StatusConflict
			userErr = fmt.Errorf("A capability named '%s' already exists!", cap.Name)
			api.HandleErr(w, r, tx, errCode, userErr, nil)
			return
		}
	}

	row := tx.QueryRow(replaceQuery, cap.Name, cap.Description, inf.Params["name"])
	if err := row.Scan(&cap.Description, &cap.LastUpdated, &cap.Name); err != nil {
		if err == sql.ErrNoRows {
			errCode = http.StatusNotFound
			userErr = fmt.Errorf("No capability '%s' found!", inf.Params["name"])
		} else {
			errCode = http.StatusInternalServerError
			sysErr = fmt.Errorf("Replacing capability '%s' with '%s': %v", inf.Params["name"], cap.Name, err)
		}
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	msg := "CAPABILITY: %s, ACTION: Replaced with capability '%s' (%s)'"
	msg = fmt.Sprintf(msg, inf.Params["name"], cap.Name, cap.Description)
	api.CreateChangeLogRawTx(api.ApiChange, msg, inf.User, tx)
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Capability was updated.", cap)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	inf, sysErr, userErr, errCode := api.NewInfo(r, []string{"name"}, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	capName := inf.Params["name"]

	row := tx.QueryRow(`SELECT COUNT(*) FROM api_capability WHERE capability=$1`, capName)
	var num uint
	if err := row.Scan(&num); err != nil && err != sql.ErrNoRows {
		sysErr = fmt.Errorf("Checking for API routes linked to capability %s: %v", capName, err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}

	if num > 0 {
		userErr = fmt.Errorf("Capability '%s' is used by %d API routes; cannot be deleted!", capName, num)
		errCode = http.StatusConflict
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	row = tx.QueryRow(deleteQuery, capName)
	var cap tc.Capability
	if err := row.Scan(&cap.Description, &cap.LastUpdated, &cap.Name); err != nil {
		if err == sql.ErrNoRows {
			errCode = http.StatusNotFound
			userErr = fmt.Errorf("No capability '%s' found!", capName)
		} else {
			sysErr = fmt.Errorf("Deleting capability %s: %v", capName, err)
			errCode = http.StatusInternalServerError
		}
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Capability deleted.", cap)
	api.CreateChangeLogRawTx(api.ApiChange, fmt.Sprintf("CAPABILITY: %s, ACTION: Deleted", capName), inf.User, tx)
}

func capabilityNameExists(c string, tx *sql.Tx) (bool, error) {
	row := tx.QueryRow(`SELECT name FROM capability WHERE name=$1`, c)
	var n string
	if err := row.Scan(&n); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
