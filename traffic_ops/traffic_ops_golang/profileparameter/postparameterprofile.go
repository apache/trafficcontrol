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
	"errors"
	"fmt"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	"github.com/lib/pq"
)

func PostParamProfile(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	paramProfile := tc.PostParamProfile{}
	if err := api.Parse(r.Body, inf.Tx.Tx, &paramProfile); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("parse error: "+err.Error()), nil)
		return
	}
	if err := insertParameterProfile(paramProfile, inf.Tx.Tx); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("posting parameter profile: "+err.Error()))
		return
	}
	paramName, _, _ := dbhelpers.GetParamNameFromID(int(*paramProfile.ParamID), inf.Tx.Tx)
	api.CreateChangeLogRawTx(api.ApiChange, fmt.Sprintf("PARAM: %v, ID: %v, ACTION: Assigned %v profiles to parameter", paramName, *paramProfile.ParamID, len(*paramProfile.ProfileIDs)), inf.User, inf.Tx.Tx)
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, fmt.Sprintf("%d profiles were assigned to the %d parameter", len(*paramProfile.ProfileIDs), *paramProfile.ParamID), paramProfile)
}

func insertParameterProfile(post tc.PostParamProfile, tx *sql.Tx) error {
	if *post.Replace {
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
	return nil
}
