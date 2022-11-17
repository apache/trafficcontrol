package deliveryservice

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

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
)

const safeUpdateQuery = `
UPDATE deliveryservice
SET display_name=$1,
    info_url=$2,
    long_desc=$3,
    long_desc_1=$4
WHERE id = $5
RETURNING id
`

const safeUpdateQueryWithoutLD1 = `
UPDATE deliveryservice
SET display_name=$1,
    info_url=$2,
    long_desc=$3
WHERE id = $4
RETURNING id
`

// UpdateSafe is the handler for PUT requests to /deliveryservices/{{ID}}/safe.
//
// The only fields which are "safe" to modify are the displayName, infoURL, longDesc, and longDesc1.
func UpdateSafe(w http.ResponseWriter, r *http.Request) {
	var ok bool
	var err error
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	version := inf.Version
	if version == nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("TODeliveryService.UpdateSafe called with nil API version"))
		return
	}
	if version.Major == 1 && version.Minor < 1 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("TODeliveryService.UpdateSafe called with invalid API version: %d.%d", version.Major, version.Minor))
		return
	}
	dsID := inf.IntParams["id"]

	userErr, sysErr, errCode = tenant.CheckID(tx, inf.User, dsID)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	var dsr tc.DeliveryServiceSafeUpdateRequest
	if err := api.Parse(r.Body, tx, &dsr); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("decoding: %s", err), nil)
		return
	}

	_, cdn, _, err := dbhelpers.GetDSNameAndCDNFromID(inf.Tx.Tx, dsID)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("deliveryservice update safe: getting CDN from DS ID "+err.Error()))
		return
	}
	userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(inf.Tx.Tx, string(cdn), inf.User.UserName)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
		return
	}
	if version.Major > 3 && version.Minor >= 0 {
		if dsr.LongDesc1 != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("the longDesc1 field is no longer supported in API 4.0 onwards"), nil)
			return
		}
		ok, err = updateDSSafe(tx, dsID, dsr, true)
	} else {
		ok, err = updateDSSafe(tx, dsID, dsr, false)
	}
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("Updating Delivery Service (safe): %s", err))
		return
	} else if !ok {
		userErr = fmt.Errorf("no Delivery Service exists by ID '%d'", dsID)
		api.HandleErr(w, r, tx, http.StatusNotFound, userErr, nil)
		return
	}
	useIMS := false
	config, e := api.GetConfig(r.Context())
	if e == nil && config != nil {
		useIMS = config.UseIMS
	} else {
		log.Warnf("Couldn't get config %v", e)
	}
	dses, userErr, sysErr, errCode, _ := readGetDeliveryServices(r.Header, inf.Params, inf.Tx, inf.User, useIMS, *version)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	if len(dses) != 1 {
		sysErr = fmt.Errorf("Updating Delivery Service (safe): expected one Delivery Service returned from read, got %v", len(dses))
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}

	ds := dses[0].DS
	alertMsg := "Delivery Service safe update successful."
	if inf.Version == nil {
		log.Warnln("API version found to be null in DS safe update")
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, alertMsg, dses)
	} else {
		switch inf.Version.Major {
		default:
			fallthrough
		case 5:
			api.WriteRespAlertObj(w, r, tc.SuccessLevel, alertMsg, ds)
		case 4:
			if inf.Version.Minor >= 1 {
				api.WriteRespAlertObj(w, r, tc.SuccessLevel, alertMsg, []tc.DeliveryServiceV41{ds.Downgrade()})
			} else {
				api.WriteRespAlertObj(w, r, tc.SuccessLevel, alertMsg, []tc.DeliveryServiceV40{ds.Downgrade().DeliveryServiceV40})
			}
		case 3:
			legacyDS := ds.Downgrade()
			legacyDS.LongDesc1 = dses[0].LongDesc1
			legacyDS.LongDesc2 = dses[0].LongDesc2
			ret := legacyDS.DowngradeToV31()
			if inf.Version.Minor >= 1 {
				api.WriteRespAlertObj(w, r, tc.SuccessLevel, alertMsg, []tc.DeliveryServiceV31{tc.DeliveryServiceV31(ret)})
			} else {
				api.WriteRespAlertObj(w, r, tc.SuccessLevel, alertMsg, []tc.DeliveryServiceV30{ret.DeliveryServiceV30})
			}
		}
	}

	api.CreateChangeLogRawTx(api.ApiChange, fmt.Sprintf("DS: %s, ID: %d, ACTION: Updated safe fields", ds.XMLID, *ds.ID), inf.User, tx)
}

// updateDSSafe updates the given delivery service in the database. Returns whether the DS existed, and any error.
func updateDSSafe(tx *sql.Tx, dsID int, ds tc.DeliveryServiceSafeUpdateRequest, omitLD1Field bool) (bool, error) {
	var err error
	var res sql.Result
	if omitLD1Field {
		res, err = tx.Exec(safeUpdateQueryWithoutLD1, ds.DisplayName, ds.InfoURL, ds.LongDesc, dsID)
	} else {
		res, err = tx.Exec(safeUpdateQuery, ds.DisplayName, ds.InfoURL, ds.LongDesc, ds.LongDesc1, dsID)
	}
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("Checking rows affected: %s", err)
	}
	if rowsAffected < 1 {
		return false, nil
	}
	if rowsAffected > 1 {
		return false, fmt.Errorf("Too many rows affected: %v", rowsAffected)
	}
	return true, nil
}
