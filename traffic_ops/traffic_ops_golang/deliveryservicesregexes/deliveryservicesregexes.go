package deliveryservicesregexes

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
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"

	"github.com/jmoiron/sqlx"
)

func Get(dbx *sqlx.DB) http.HandlerFunc {
	db := dbx.DB
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		user, err := auth.GetCurrentUser(r.Context())
		if err != nil {
			log.Errorf("unable to retrieve current user from context: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		q := `
SELECT ds.xml_id, ds.tenant_id, dsr.set_number, r.pattern, rt.name as type
FROM deliveryservice_regex as dsr
JOIN deliveryservice as ds ON dsr.deliveryservice = ds.id
JOIN regex as r ON dsr.regex = r.id
JOIN type as rt ON r.type = rt.id
`
		rows, err := db.Query(q)
		if err != nil {
			handleErrs(http.StatusInternalServerError, errors.New("querying: "+err.Error()))
			return
		}
		dsRegexes := map[tc.DeliveryServiceName][]tc.DeliveryServiceRegex{}
		for rows.Next() {
			// if (!$tenant_utils->is_ds_resource_accessible($tenants_data, $ds->tenant_id)) {
			// 	next;
			// }
			dsName := ""
			dsTenantID := 0
			setNumber := 0
			pattern := ""
			typeName := ""
			if err = rows.Scan(&dsName, &dsTenantID, &setNumber, &pattern, &typeName); err != nil {
				handleErrs(http.StatusInternalServerError, errors.New("querying: "+err.Error()))
				return
			}
			if ok, err := tenant.IsResourceAuthorizedToUser(dsTenantID, user, dbx); !ok {
				continue
			} else if err != nil {
				handleErrs(http.StatusInternalServerError, errors.New("checking tenancy: "+err.Error()))
				return
			}
			dsRegexes[tc.DeliveryServiceName(dsName)] = append(dsRegexes[tc.DeliveryServiceName(dsName)], tc.DeliveryServiceRegex{
				Type:      typeName,
				SetNumber: setNumber,
				Pattern:   pattern,
			})
		}
		respRegexes := []tc.DeliveryServiceRegexes{}
		for dsName, regexes := range dsRegexes {
			respRegexes = append(respRegexes, tc.DeliveryServiceRegexes{DSName: string(dsName), Regexes: regexes})
		}
		respBts, err := json.Marshal(&tc.DeliveryServiceRegexResponse{Response: respRegexes})
		if err != nil {
			handleErrs(http.StatusInternalServerError, errors.New("marshalling JSON: "+err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(respBts)
	}
}

func DSGet(dbx *sqlx.DB) http.HandlerFunc {
	db := dbx.DB
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		user, err := auth.GetCurrentUser(r.Context())
		if err != nil {
			log.Errorf("unable to retrieve current user from context: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		params, err := api.GetCombinedParams(r)
		if err != nil {
			log.Errorf("unable to get parameters from request: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		dsIDStr, ok := params["dsid"]
		if !ok {
			log.Errorln("no delivery service ID")
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		dsID, err := strconv.Atoi(dsIDStr)
		if err != nil {
			log.Errorln("Delivery service ID '" + dsIDStr + "' not an integer")
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		q := `
SELECT ds.tenant_id, dsr.set_number, r.id, r.pattern, rt.id as type, rt.name as type_name
FROM deliveryservice_regex as dsr
JOIN deliveryservice as ds ON dsr.deliveryservice = ds.id
JOIN regex as r ON dsr.regex = r.id
JOIN type as rt ON r.type = rt.id
WHERE ds.ID = $1
ORDER BY dsr.set_number ASC
`
		rows, err := db.Query(q, dsID)
		if err != nil {
			handleErrs(http.StatusInternalServerError, errors.New("querying: "+err.Error()))
			return
		}
		regexes := []tc.DeliveryServiceIDRegex{}
		for rows.Next() {
			dsTenantID := 0
			setNumber := 0
			id := 0
			pattern := ""
			typeID := 0
			typeName := ""
			if err = rows.Scan(&dsTenantID, &setNumber, &id, &pattern, &typeID, &typeName); err != nil {
				handleErrs(http.StatusInternalServerError, errors.New("querying: "+err.Error()))
				return
			}
			if ok, err := tenant.IsResourceAuthorizedToUser(dsTenantID, user, dbx); !ok {
				continue
			} else if err != nil {
				handleErrs(http.StatusInternalServerError, errors.New("checking tenancy: "+err.Error()))
				return
			}
			regexes = append(regexes, tc.DeliveryServiceIDRegex{
				ID:        id,
				Type:      typeID,
				TypeName:  typeName,
				SetNumber: setNumber,
				Pattern:   pattern,
			})
		}
		respBts, err := json.Marshal(&tc.DeliveryServiceIDRegexResponse{Response: regexes})
		if err != nil {
			handleErrs(http.StatusInternalServerError, errors.New("marshalling JSON: "+err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(respBts)
	}
}

func DSGetID(dbx *sqlx.DB) http.HandlerFunc {
	db := dbx.DB
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		user, err := auth.GetCurrentUser(r.Context())
		if err != nil {
			log.Errorf("unable to retrieve current user from context: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		params, err := api.GetCombinedParams(r)
		if err != nil {
			log.Errorf("unable to get parameters from request: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		dsIDStr, ok := params["dsid"]
		if !ok {
			log.Errorln("no delivery service ID")
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		dsID, err := strconv.Atoi(dsIDStr)
		if err != nil {
			log.Errorln("Delivery service ID '" + dsIDStr + "' not an integer")
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		regexIDStr, ok := params["regexid"]
		if !ok {
			log.Errorln("no regex ID")
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		regexID, err := strconv.Atoi(regexIDStr)
		if err != nil {
			log.Errorln("Regex ID '" + regexIDStr + "' not an integer")
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		q := `
SELECT ds.tenant_id, dsr.set_number, r.id, r.pattern, rt.id as type, rt.name as type_name
FROM deliveryservice_regex as dsr
JOIN deliveryservice as ds ON dsr.deliveryservice = ds.id
JOIN regex as r ON dsr.regex = r.id
JOIN type as rt ON r.type = rt.id
WHERE ds.ID = $1
AND r.ID = $2
ORDER BY dsr.set_number ASC
`
		rows, err := db.Query(q, dsID, regexID)
		if err != nil {
			handleErrs(http.StatusInternalServerError, errors.New("querying: "+err.Error()))
			return
		}
		regexes := []tc.DeliveryServiceIDRegex{}
		for rows.Next() {
			dsTenantID := 0
			setNumber := 0
			id := 0
			pattern := ""
			typeID := 0
			typeName := ""
			if err = rows.Scan(&dsTenantID, &setNumber, &id, &pattern, &typeID, &typeName); err != nil {
				handleErrs(http.StatusInternalServerError, errors.New("querying: "+err.Error()))
				return
			}
			if ok, err := tenant.IsResourceAuthorizedToUser(dsTenantID, user, dbx); !ok {
				continue
			} else if err != nil {
				handleErrs(http.StatusInternalServerError, errors.New("checking tenancy: "+err.Error()))
				return
			}
			regexes = append(regexes, tc.DeliveryServiceIDRegex{
				ID:        id,
				Type:      typeID,
				TypeName:  typeName,
				SetNumber: setNumber,
				Pattern:   pattern,
			})
		}
		respBts, err := json.Marshal(&tc.DeliveryServiceIDRegexResponse{Response: regexes})
		if err != nil {
			handleErrs(http.StatusInternalServerError, errors.New("marshalling JSON: "+err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(respBts)
	}
}

func Post(dbx *sqlx.DB) http.HandlerFunc {
	db := dbx.DB
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		user, err := auth.GetCurrentUser(r.Context())
		if err != nil {
			handleErrs(http.StatusInternalServerError, errors.New("unable to retrieve current user from context: "+err.Error()))
			return
		}
		params, err := api.GetCombinedParams(r)
		if err != nil {
			handleErrs(http.StatusInternalServerError, errors.New("unable to get parameters from request: "+err.Error()))
			return
		}
		dsIDStr, ok := params["dsid"]
		if !ok {
			handleErrs(http.StatusInternalServerError, errors.New("no deliveryservice ID"))
			return
		}
		dsID, err := strconv.Atoi(dsIDStr)
		if err != nil {
			handleErrs(http.StatusInternalServerError, errors.New("deliveryservice ID not an integer"))
			return
		}
		dsTenantID := 0
		if err := db.QueryRow(`SELECT tenant_id from deliveryservice where id = $1`, dsID).Scan(&dsTenantID); err != nil {
			if err != sql.ErrNoRows {
				log.Errorln("getting deliveryservice name: " + err.Error())
			}
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		if ok, err := tenant.IsResourceAuthorizedToUser(dsTenantID, user, dbx); !ok {
			handleErrs(http.StatusInternalServerError, errors.New("unauthorized"))
			return
		} else if err != nil {
			handleErrs(http.StatusInternalServerError, errors.New("checking tenancy: "+err.Error()))
			return
		}
		dsr := tc.DeliveryServiceRegexPost{}
		if err := json.NewDecoder(r.Body).Decode(&dsr); err != nil {
			log.Errorln("failed to parse body: " + err.Error())
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		regexID := 0
		if err := db.QueryRow(`INSERT INTO regex (pattern, type) VALUES ($1, $2) RETURNING id`, dsr.Pattern, dsr.Type).Scan(&regexID); err != nil {
			log.Errorln("inserting regex: " + err.Error())
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		if _, err := db.Exec(`INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) values ($1, $2, $3)`, dsID, regexID, dsr.SetNumber); err != nil {
			log.Errorln("inserting deliveryservice_regex: " + err.Error())
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		typeName := ""
		if err := db.QueryRow(`SELECT name from type where id = $1`, dsr.Type).Scan(&typeName); err != nil {
			if err != sql.ErrNoRows {
				log.Errorln("getting regex type: " + err.Error())
			}
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		respObj := tc.DeliveryServiceIDRegex{
			ID:        regexID,
			Pattern:   dsr.Pattern,
			Type:      dsr.Type,
			TypeName:  typeName,
			SetNumber: dsr.SetNumber,
		}
		resp := struct {
			Response tc.DeliveryServiceIDRegex `json:"response"`
			tc.Alerts
		}{respObj, tc.CreateAlerts(tc.SuccessLevel, "Delivery service regex creation was successful.")}

		respBts, err := json.Marshal(&resp)
		if err != nil {
			handleErrs(http.StatusInternalServerError, errors.New("marshalling JSON: "+err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(respBts)
	}
}

func Put(dbx *sqlx.DB) http.HandlerFunc {
	db := dbx.DB
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		user, err := auth.GetCurrentUser(r.Context())
		if err != nil {
			log.Errorf("unable to retrieve current user from context: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		params, err := api.GetCombinedParams(r)
		if err != nil {
			log.Errorf("unable to get parameters from request: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		dsIDStr, ok := params["dsid"]
		if !ok {
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		dsID, err := strconv.Atoi(dsIDStr)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		regexIDStr, ok := params["regexid"]
		if !ok {
			handleErrs(http.StatusInternalServerError, errors.New("no regex ID"))
			return
		}
		regexID, err := strconv.Atoi(regexIDStr)
		if err != nil {
			handleErrs(http.StatusInternalServerError, errors.New("Regex ID '"+regexIDStr+"' not an integer"))
			return
		}
		dsTenantID := 0
		if err := db.QueryRow(`SELECT tenant_id from deliveryservice where id = $1`, dsID).Scan(&dsTenantID); err != nil {
			if err != sql.ErrNoRows {
				log.Errorln("getting deliveryservice name: " + err.Error())
			}
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		if ok, err := tenant.IsResourceAuthorizedToUser(dsTenantID, user, dbx); !ok {
			handleErrs(http.StatusInternalServerError, errors.New("unauthorized"))
			return
		} else if err != nil {
			handleErrs(http.StatusInternalServerError, errors.New("checking tenancy: "+err.Error()))
			return
		}
		dsr := tc.DeliveryServiceRegexPost{} // PUT uses same format as POST
		if err := json.NewDecoder(r.Body).Decode(&dsr); err != nil {
			log.Errorln("failed to parse body: " + err.Error())
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		if _, err := db.Exec(`UPDATE regex SET pattern=$1, type=$2 WHERE id=$3`, dsr.Pattern, dsr.Type, regexID); err != nil {
			log.Errorln("deliveryservicesregexes.Put: updating regex: " + err.Error())
			handleErrs(http.StatusInternalServerError, errors.New("server error"))
			return
		}
		if _, err := db.Exec(`UPDATE deliveryservice_regex SET set_number=$1 WHERE deliveryservice=$2 AND regex=$3`, dsr.SetNumber, dsID, regexID); err != nil {
			log.Errorln("deliveryservicesregexes.Put: updating deliveryservice_regex: " + err.Error())
			handleErrs(http.StatusInternalServerError, errors.New("server error"))
			return
		}
		typeName := ""
		if err := db.QueryRow(`SELECT name from type where id = $1`, dsr.Type).Scan(&typeName); err != nil {
			if err != sql.ErrNoRows {
				log.Errorln("getting regex type: " + err.Error())
			}
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		respObj := tc.DeliveryServiceIDRegex{
			ID:        regexID,
			Pattern:   dsr.Pattern,
			Type:      dsr.Type,
			TypeName:  typeName,
			SetNumber: dsr.SetNumber,
		}
		resp := struct {
			Response tc.DeliveryServiceIDRegex `json:"response"`
			tc.Alerts
		}{respObj, tc.CreateAlerts(tc.SuccessLevel, "Delivery service regex creation was successful.")}
		respBts, err := json.Marshal(&resp)
		if err != nil {
			handleErrs(http.StatusInternalServerError, errors.New("marshalling JSON: "+err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(respBts)
	}
}

func Delete(dbx *sqlx.DB) http.HandlerFunc {
	db := dbx.DB
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		user, err := auth.GetCurrentUser(r.Context())
		if err != nil {
			log.Errorf("unable to retrieve current user from context: %s", err)
			handleErrs(http.StatusInternalServerError, errors.New("server error"))
			return
		}
		params, err := api.GetCombinedParams(r)
		if err != nil {
			log.Errorf("unable to get parameters from request: %s", err)
			handleErrs(http.StatusInternalServerError, errors.New("server error"))
			return
		}
		dsIDStr, ok := params["dsid"]
		if !ok {
			log.Errorf("deliveryservicesregexes.Delete: no dsid parameter")
			handleErrs(http.StatusInternalServerError, errors.New("internal server error"))
			return
		}
		dsID, err := strconv.Atoi(dsIDStr)
		if err != nil {
			handleErrs(http.StatusBadRequest, errors.New("delivery service ID is not an integer"))
			return
		}
		regexIDStr, ok := params["regexid"]
		if !ok {
			handleErrs(http.StatusBadRequest, errors.New("no regex ID"))
			return
		}
		regexID, err := strconv.Atoi(regexIDStr)
		if err != nil {
			handleErrs(http.StatusBadRequest, errors.New("Regex ID '"+regexIDStr+"' not an integer"))
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Errorln("could not begin transaction: " + err.Error())
			handleErrs(http.StatusInternalServerError, errors.New("could not begin transaction: "+err.Error()))
			return
		}
		commitTx := false
		defer FinishTx(tx, &commitTx)

		count := 0
		if err := db.QueryRow(`SELECT count(*) from deliveryservice_regex where deliveryservice = $1`, dsID).Scan(&count); err != nil {
			if err != sql.ErrNoRows {
				handleErrs(http.StatusNotFound, errors.New("not found"))
				return
			}
			log.Errorln(errors.New("getting deliveryservice regex count: " + err.Error()))
			handleErrs(http.StatusInternalServerError, errors.New("internal server error"))
			return
		}
		if count < 2 {
			handleErrs(http.StatusBadRequest, errors.New("a delivery service must have at least one regex"))
			return
		}

		dsTenantID := 0
		if err := db.QueryRow(`SELECT tenant_id from deliveryservice where id = $1`, dsID).Scan(&dsTenantID); err != nil {
			if err != sql.ErrNoRows {
				log.Errorln("getting deliveryservice name: " + err.Error())
			}
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		if ok, err := tenant.IsResourceAuthorizedToUser(dsTenantID, user, dbx); !ok {
			handleErrs(http.StatusUnauthorized, errors.New("unauthorized"))
			return
		} else if err != nil {
			log.Errorln("checking tenancy: " + err.Error())
			handleErrs(http.StatusInternalServerError, errors.New("server error"))
			return
		}

		result, err := tx.Exec(`DELETE FROM deliveryservice_regex WHERE deliveryservice = $1 and regex = $2`, dsID, regexID)
		if err != nil {
			log.Errorln("deliveryservicesregexes.Delete deleting delivery service regexes: " + err.Error())
			handleErrs(http.StatusInternalServerError, errors.New("server error"))
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Errorln("deliveryservicesregexes.Delete delete error: " + err.Error())
			handleErrs(http.StatusInternalServerError, errors.New("server error"))
			return
		}
		if rowsAffected != 1 {
			if rowsAffected < 1 {
				handleErrs(http.StatusNotFound, errors.New("not found"))
				return
			}
			log.Errorf("this create affected too many rows: %d", rowsAffected)
			handleErrs(http.StatusInternalServerError, errors.New("server error"))
			return
		}

		log.Debugf("changelog for delete on object")
		api.CreateChangeLogRaw(api.ApiChange, fmt.Sprintf(`deleted deliveryservice_regex {"ds": %d, "regex": %d}`, dsID, regexID), user, dbx.DB)
		resp := struct {
			tc.Alerts
		}{tc.CreateAlerts(tc.SuccessLevel, "deliveryservice_regex was deleted.")}
		respBts, err := json.Marshal(resp)
		if err != nil {
			log.Errorln("deliveryservicesregexes.Delete creating JSON: " + err.Error())
			handleErrs(http.StatusInternalServerError, errors.New("server error"))
			return
		}
		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		w.Write(respBts)
		commitTx = true
	}
}

// FinishTx commits the transaction if commit is true when it's called, otherwise it rolls back the transaction. This is designed to be called in a defer.
func FinishTx(tx *sql.Tx, commit *bool) {
	if tx == nil {
		return
	}
	if !*commit {
		tx.Rollback()
		return
	}
	tx.Commit()
}
