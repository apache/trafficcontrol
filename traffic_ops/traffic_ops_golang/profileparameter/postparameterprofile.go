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

func PostParamProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		paramProfile := tc.PostParamProfile{}
		if err := json.NewDecoder(r.Body).Decode(&paramProfile); err != nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("malformed JSON"), nil)
			return
		}

		if ok, err := paramExists(paramProfile.ParamID, db); err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, fmt.Errorf("checking param ID %d existence: "+err.Error(), paramProfile.ParamID))
			return
		} else if !ok {
			api.HandleErr(w, r, http.StatusBadRequest, fmt.Errorf("no parameter with ID %d exists", paramProfile.ParamID), nil)
			return
		}
		if ok, err := profilesExist(paramProfile.ProfileIDs, db); err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, fmt.Errorf("checking profiles IDs %v existence: "+err.Error(), paramProfile.ProfileIDs))
			return
		} else if !ok {
			api.HandleErr(w, r, http.StatusBadRequest, fmt.Errorf("profiles with IDs %v don't all exist", paramProfile.ProfileIDs), nil)
			return
		}
		if err := insertParameterProfile(paramProfile, db); err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("posting parameter profile: "+err.Error()))
			return
		}
		// TODO create helper func
		resp := struct {
			Response tc.PostParamProfile `json:"response"`
			tc.Alerts
		}{paramProfile, tc.CreateAlerts(tc.SuccessLevel, fmt.Sprintf("%d profiles were assigned to the %d parameter", len(paramProfile.ProfileIDs), paramProfile.ParamID))}
		respBts, err := json.Marshal(resp)
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("posting parameter profiles: "+err.Error()))
			return
		}
		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		w.Write(respBts)
	}
}

func paramExists(id int64, db *sql.DB) (bool, error) {
	count := 0
	if err := db.QueryRow(`SELECT count(*) from parameter where id = $1`, id).Scan(&count); err != nil {
		return false, errors.New("querying param existence from id: " + err.Error())
	}
	return count > 0, nil
}

func profilesExist(ids []int64, db *sql.DB) (bool, error) {
	count := 0
	if err := db.QueryRow(`SELECT count(*) from profile where id = ANY($1)`, pq.Array(ids)).Scan(&count); err != nil {
		return false, errors.New("querying profiles existence from id: " + err.Error())
	}
	return count == len(ids), nil
}

func insertParameterProfile(post tc.PostParamProfile, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.New("beginning transaction: " + err.Error())
	}
	commitTx := false
	defer FinishTx(tx, &commitTx)

	if post.Replace {
		if _, err := tx.Exec(`DELETE FROM profile_parameter WHERE parameter = $1`, post.ParamID); err != nil {
			return errors.New("deleting old parameter profile: " + err.Error())
		}
	}

	q := `
INSERT INTO profile_parameter (profile, parameter)
VALUES (unnest($1::int[]), $2)
ON CONFLICT DO NOTHING;
`
	if _, err := tx.Exec(q, pq.Array(post.ProfileIDs), post.ParamID); err != nil {
		return errors.New("inserting parameter profile: " + err.Error())
	}
	commitTx = true
	return nil
}
