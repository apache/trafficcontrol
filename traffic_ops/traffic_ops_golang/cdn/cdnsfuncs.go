package cdn

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
	"errors"
	"net/http"
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/lib/go-util"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc/v13"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tovalidate"

	"github.com/asaskevich/govalidator"
	validation "github.com/go-ozzo/ozzo-validation"
)

func Get(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, _, userErr, sysErr, errCode := api.AllParams(r, nil)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, errCode, userErr, sysErr)
			return
		}
		orderBy := DefaultCDNOrderBy
		if paramOrderBy, ok := params["orderby"]; ok {
			validCols := orderByCols()
			if sqlCol, ok := validCols[paramOrderBy]; ok {
				orderBy = sqlCol
			} else {
				api.HandleErr(w, r, http.StatusBadRequest, errors.New("invalid orderby parameter"), nil)
				return
			}
		}
		api.RespWriter(w, r)(getCDNs(db, orderBy))
	}
}

func Post(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		user, err := auth.GetCurrentUser(r.Context())
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting user: "+err.Error()))
			return
		}
		cdn := v13.CDNNullable{}
		if err := json.NewDecoder(r.Body).Decode(&cdn); err != nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("malformed JSON: "+err.Error()), nil)
			return
		}

		if errs := validatePost(&cdn); len(errs) > 0 {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("invalid request: "+util.JoinErrs(errs).Error()), nil)
			return
		}
		newCDN, err := insert(db, &cdn)
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("inserting CDN: "+err.Error()))
			return
		}
		api.CreateChangeLogRaw(api.ApiChange, "Created CDN with id:" + strconv.Itoa(*cdn.ID) + " and name: " + *cdn.Name, *user, db)
		api.WriteResp(w, r, newCDN)
	}
}

func Put(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		_, intParams, userErr, sysErr, errCode := api.AllParams(r, []string{"id"})
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, errCode, userErr, sysErr)
			return
		}
		user, err := auth.GetCurrentUser(r.Context())
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting user: "+err.Error()))
			return
		}
		cdn := v13.CDNNullable{}
		if err := json.NewDecoder(r.Body).Decode(&cdn); err != nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("malformed JSON: "+err.Error()), nil)
			return
		}
		id := intParams["id"]
		cdn.ID = &id

		if errs := validatePost(&cdn); len(errs) > 0 {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("invalid request: "+util.JoinErrs(errs).Error()), nil)
			return
		}
		newCDN, err := update(db, &cdn)
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("updating CDN: "+err.Error()))
			return
		}
		api.CreateChangeLogRaw(api.ApiChange, "Updated CDN name " + *cdn.Name + " for id: " + strconv.Itoa(*cdn.ID), *user, db)
		api.WriteResp(w, r, newCDN)
	}
}

func Delete(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, intParams, userErr, sysErr, errCode := api.AllParams(r, []string{"id"})
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, errCode, userErr, sysErr)
			return
		}
		user, err := auth.GetCurrentUser(r.Context())
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting user: "+err.Error()))
			return
		}
		cdnName, ok, err := cdnNameByID(db, intParams["id"])
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("checking CDN existence: "+err.Error()))
			return
		} else if !ok {
			api.HandleErr(w, r, http.StatusNotFound, nil, nil)
			return
		}
		if ok, err := cdnUnused(db, cdnName); err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("checking CDN usage: "+err.Error()))
			return
		} else if !ok {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("Failed to delete cdn name = "+string(cdnName)+" has delivery services or servers"), nil)
			return
		}
		if err := deleteCDNByName(db, tc.CDNName(cdnName)); err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("deleting CDN: "+err.Error()))
			return
		}
		api.CreateChangeLogRaw(api.ApiChange, "Delete cdn " + string(cdnName), *user, db)
		api.WriteRespAlert(w, r, tc.SuccessLevel, "cdn was deleted.")
	}
}

// func isValidCDNchar(r rune) bool {
// 	if r >= 'a' && r <= 'z' {
// 		return true
// 	}
// 	if r >= 'A' && r <= 'Z' {
// 		return true
// 	}
// 	if r >= '0' && r <= '9' {
// 		return true
// 	}
// 	if r == '.' || r == '-' {
// 		return true
// 	}
// 	return false
// }

// // IsValidCDNName returns true if the name contains only characters valid for a CDN name
// func IsValidCDNName(str string) bool {
// 	i := strings.IndexFunc(str, func(r rune) bool { return !isValidCDNchar(r) })
// 	return i == -1
// }

func orderByCols() map[string]string {
	return map[string]string {
		"id": "id",
		"name": "name",
		"domainName": "domain_name",
		"dnssecEnabled":"dnssec_enabled",
		"lastUpdated":"lastUpdated",
	}
}

const DefaultCDNOrderBy = "name"

func validatePost(cdn *v13.CDNNullable) []error {
	validName := validation.NewStringRule(IsValidCDNName, "invalid characters found - Use alphanumeric . or - .")
	validDomainName := validation.NewStringRule(govalidator.IsDNSName, "not a valid domain name")
	errs := validation.Errors{
		"name":       validation.Validate(cdn.Name, validation.Required, validName),
		"domainName": validation.Validate(cdn.DomainName, validation.Required, validDomainName),
	}
	return tovalidate.ToErrors(errs)
}

func getCDNs(db *sql.DB, orderBy string) ([]v13.CDNNullable, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.New("beginning transaction: " + err.Error())
	}
	commitTx := false
	defer dbhelpers.FinishTx(tx, &commitTx)

	rows, err := tx.Query(`SELECT id, name, domain_name, dnssec_enabled, last_updated FROM cdn ORDER BY $1`, orderBy)
	if err != nil {
		return nil, errors.New("querying cdns: " + err.Error())
	}
	cdns := []v13.CDNNullable{}
	defer rows.Close()
	for rows.Next() {
		c := v13.CDNNullable{}
		if err := rows.Scan(&c.ID, &c.Name, &c.DomainName, &c.DNSSECEnabled, &c.LastUpdated); err != nil {
			return nil, errors.New("scanning cdns: " + err.Error())
		}
		cdns = append(cdns, c)
	}
	commitTx = true
	return cdns, nil
}

func cdnNameByID(db *sql.DB, id int) (tc.CDNName, bool, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", false, errors.New("beginning transaction: " + err.Error())
	}
	commitTx := false
	defer dbhelpers.FinishTx(tx, &commitTx)

	name := ""
	if err := tx.QueryRow(`SELECT name FROM cdn WHERE id = $1`, name).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying cdn existence: " + err.Error())
	}
	commitTx = true
	return tc.CDNName(name), true, nil
}

func insert(db *sql.DB, cdn *v13.CDNNullable) (*v13.CDNNullable, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.New("beginning transaction: " + err.Error())
	}
	commitTx := false
	defer dbhelpers.FinishTx(tx, &commitTx)

	err = tx.QueryRow(`
INSERT INTO cdn (name, domain_name, dnssec_enabled)
VALUES ($1, $2, $3)
RETURNING id, last_updated
`, cdn.Name, cdn.DomainName, cdn.DNSSECEnabled).Scan(&cdn.ID, &cdn.LastUpdated)
	if err != nil {
		return nil, errors.New("inserting cdn: " + err.Error())
	}
	commitTx = true
	return cdn, nil
}

func update(db *sql.DB, cdn *v13.CDNNullable) (*v13.CDNNullable, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.New("beginning transaction: " + err.Error())
	}
	commitTx := false
	defer dbhelpers.FinishTx(tx, &commitTx)

	err = tx.QueryRow(`
UPDATE cdn set name=$1, domain_name=$2, dnssec_enabled=$3 WHERE id=$4
RETURNING last_updated
`, cdn.Name, cdn.DomainName, cdn.DNSSECEnabled, cdn.ID).Scan(&cdn.LastUpdated)
	if err != nil {
		return nil, errors.New("updating cdn: " + err.Error())
	}
	commitTx = true
	return cdn, nil
}
