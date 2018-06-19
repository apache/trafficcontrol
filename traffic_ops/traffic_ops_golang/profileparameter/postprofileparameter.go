package profileparameter

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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"

	"github.com/lib/pq"
)

func PostProfileParam(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		profileParam := tc.PostProfileParam{}
		if err := json.NewDecoder(r.Body).Decode(&profileParam); err != nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("malformed JSON"), nil)
			return
		}

		if ok, err := profileExists(profileParam.ProfileID, db); err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, fmt.Errorf("checking profile ID %d existence: "+err.Error(), profileParam.ProfileID))
			return
		} else if !ok {
			api.HandleErr(w, r, http.StatusBadRequest, fmt.Errorf("no profile with ID %d exists", profileParam.ProfileID), nil)
			return
		}
		if ok, err := paramsExist(profileParam.ParamIDs, db); err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, fmt.Errorf("checking parameters IDs %v existence: "+err.Error(), profileParam.ParamIDs))
			return
		} else if !ok {
			api.HandleErr(w, r, http.StatusBadRequest, fmt.Errorf("parameters with IDs %v don't all exist", profileParam.ParamIDs), nil)
			return
		}
		if err := insertProfileParameter(profileParam, db); err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("posting profile parameter: "+err.Error()))
			return
		}
		// TODO create helper func
		resp := struct {
			Response tc.PostProfileParam `json:"response"`
			tc.Alerts
		}{profileParam, tc.CreateAlerts(tc.SuccessLevel, fmt.Sprintf("%d parameters were assigned to the %d profile", len(profileParam.ParamIDs), profileParam.ProfileID))}
		respBts, err := json.Marshal(resp)
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("posting profile parameters by name: "+err.Error()))
			return
		}
		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		w.Write(respBts)
	}
}

func profileExists(id int64, db *sql.DB) (bool, error) {
	count := 0
	if err := db.QueryRow(`SELECT count(*) from profile where id = $1`, id).Scan(&count); err != nil {
		return false, errors.New("querying profile existence from id: " + err.Error())
	}
	return count > 0, nil
}

func paramsExist(ids []int64, db *sql.DB) (bool, error) {
	count := 0
	if err := db.QueryRow(`SELECT count(*) from parameter where id = ANY($1)`, pq.Array(ids)).Scan(&count); err != nil {
		return false, errors.New("querying parameters existence from id: " + err.Error())
	}
	return count == len(ids), nil
}

func insertProfileParameter(post tc.PostProfileParam, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.New("beginning transaction: " + err.Error())
	}
	commitTx := false
	defer FinishTx(tx, &commitTx)

	if post.Replace {
		if _, err := tx.Exec(`DELETE FROM profile_parameter WHERE profile = $1`, post.ProfileID); err != nil {
			return errors.New("deleting old profile parameter: " + err.Error())
		}
	}

	q := `
INSERT INTO profile_parameter (profile, parameter)
VALUES ($1, unnest($2::int[]))
ON CONFLICT DO NOTHING;
`
	if _, err := tx.Exec(q, post.ProfileID, pq.Array(post.ParamIDs)); err != nil {
		return errors.New("inserting profile parameter: " + err.Error())
	}
	commitTx = true
	return nil
}
