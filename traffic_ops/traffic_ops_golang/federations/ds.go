package federations

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
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	"github.com/lib/pq"
)

func PostDSes(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	fedID := inf.IntParams["id"]
	fedName, ok, err := dbhelpers.GetFedNameByID(inf.Tx.Tx, fedID)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting federation cname from ID '"+string(fedID)+"': "+err.Error()))
	} else if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("federation not found: "+err.Error()), nil)
	}

	post := tc.FederationDSPost{}
	if err := api.Parse(r.Body, inf.Tx.Tx, &post); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("parse error: "+err.Error()), nil)
		return
	}

	if post.Replace != nil && *post.Replace {
		if len(post.DSIDs) < 1 {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("A federation must have at least one delivery service assigned"), nil)
			return
		}
		if err := deleteDSFeds(inf.Tx.Tx, fedID); err != nil {
			userErr, sysErr, errCode := api.ParseDBError(err)
			api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
			return
		}
	}

	if len(post.DSIDs) > 0 {
		// there might be no DSes, if the user is trying to clear the assignments
		if err := insertDSFeds(inf.Tx.Tx, fedID, post.DSIDs); err != nil {
			userErr, sysErr, errCode := api.ParseDBError(err)
			api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
			return
		}
	}
	api.CreateChangeLogRawTx(api.ApiChange, fmt.Sprintf("FEDERATION: %v, ID: %v, ACTION: Assign DSes to federation", fedName, fedID), inf.User, inf.Tx.Tx)
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, strconv.Itoa(len(post.DSIDs))+" delivery service(s) were assigned to the federation "+strconv.Itoa(fedID), post)
}

func deleteDSFeds(tx *sql.Tx, fedID int) error {
	qry := `DELETE FROM federation_deliveryservice WHERE federation = $1`
	_, err := tx.Exec(qry, fedID)
	return err
}

func insertDSFeds(tx *sql.Tx, fedID int, dsIDs []int) error {
	qry := `
INSERT INTO federation_deliveryservice (federation, deliveryservice)
VALUES ($1, unnest($2::integer[]))
`
	_, err := tx.Exec(qry, fedID, pq.Array(dsIDs))
	return err
}
