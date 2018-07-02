package user

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
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"

	"github.com/lib/pq"
)

func GetAvailableDSes(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	userTenantID, ok, err := getUserTenantID(inf.Tx.Tx, inf.IntParams["id"])
	if err != nil {
		api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting user tenant: "+err.Error()))
		return
	}
	if !ok {
		api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting user tenant: "+err.Error()))
		return
	}
	if userTenantID == nil {
		userTenantID = util.IntPtr(-1) // set to an invalid tenant, so IsResourceAuthorized will succeed if tenancy is disabled.
	}

	if ok, err := tenant.IsResourceAuthorizedToUserTx(*userTenantID, inf.User, inf.Tx.Tx); err != nil {
		api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("checking if requested user is authorized to current user: "+err.Error()))
		return
	} else if !ok {
		api.HandleErr(w, r, http.StatusForbidden, errors.New("not authorized for that user"), nil)
		return
	}

	userDSes, err := getUserDSes(inf.Tx.Tx, inf.IntParams["id"], *userTenantID)
	if err != nil {
		api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting user delivery services: "+err.Error()))
		return
	}
	*inf.CommitTx = true
	api.WriteResp(w, r, userDSes)
}

// getUserTenantID returns the tenant ID of the given user. May return nil, if the user has no tenant ID in the database.
func getUserTenantID(tx *sql.Tx, userID int) (*int, bool, error) {
	tenantID := (*int)(nil)
	if err := tx.QueryRow(`select tenant_id from tm_user where id = $1`, userID).Scan(&tenantID); err != nil {
		if err == sql.ErrNoRows {
			return util.IntPtr(-1), false, nil
		}
		return util.IntPtr(-1), false, errors.New("querying user tenant: " + err.Error())
	}
	return tenantID, true, nil
}

func getUserDSes(tx *sql.Tx, userID int, userTenantID int) ([]tc.DeliveryServiceAvailableInfo, error) {
	q := `
SELECT ds.xml_id, ds.id, ds.display_name
FROM deliveryservice as ds
WHERE ds.id NOT IN (
  SELECT dsu.deliveryservice FROM deliveryservice_tmuser as dsu WHERE dsu.tm_user_id = $1
)
`
	qParams := []interface{}{userID}
	if tenant.IsTenancyEnabledTx(tx) {
		tenantIDs, err := tenant.GetUserTenantIDListTx(userTenantID, tx)
		if err != nil {
			return nil, errors.New("getting user tenant ID list: " + err.Error())
		}
		q += `
AND ds.tenant_id = ANY($2)
`
		qParams = append(qParams, pq.Array(tenantIDs))
	}
	rows, err := tx.Query(q, qParams...)
	if err != nil {
		return nil, errors.New("querying user available delivery services: " + err.Error())
	}
	defer rows.Close()
	dses := []tc.DeliveryServiceAvailableInfo{}
	for rows.Next() {
		ds := tc.DeliveryServiceAvailableInfo{}
		if err := rows.Scan(&ds.XMLID, &ds.ID, &ds.DisplayName); err != nil {
			return nil, errors.New("scanning user available delivery services: " + err.Error())
		}
		dses = append(dses, ds)
	}
	return dses, nil
}
