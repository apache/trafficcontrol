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
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/util/ims"

	"github.com/asaskevich/govalidator"
	validation "github.com/go-ozzo/ozzo-validation"
)

// TOCDN is the struct needed for the CRUDer
type TOCDN struct {
	api.APIInfoImpl `json:"-"`
	tc.CDNNullable
}

func Read(w http.ResponseWriter, r *http.Request) {
	var runSecond bool
	var maxTime time.Time
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	// Query Parameters to Database Query column mappings
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"domainName":    dbhelpers.WhereColumnInfo{Column: "domain_name"},
		"dnssecEnabled": dbhelpers.WhereColumnInfo{Column: "dnssec_enabled"},
		"id":            dbhelpers.WhereColumnInfo{Column: "id", Checker: api.IsInt},
		"name":          dbhelpers.WhereColumnInfo{Column: "name"},
		"ttlOverride":   dbhelpers.WhereColumnInfo{Column: "ttl_override", Checker: api.IsInt},
	}

	if _, ok := inf.Params["orderby"]; !ok {
		inf.Params["orderby"] = "name"
	}
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, queryParamsToQueryCols)
	if len(errs) > 0 {
		api.HandleErr(w, r, tx.Tx, http.StatusBadRequest, util.JoinErrs(errs), nil)
		return
	}

	if inf.Config.UseIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(tx, r.Header, queryValues, SelectMaxLastUpdatedQuery(where))
		if !runSecond {
			log.Debugln("IMS HIT")
			api.AddLastModifiedHdr(w, maxTime)
			w.WriteHeader(http.StatusNotModified)
			return
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}

	query := selectQuery(inf.Version) + where + orderBy + pagination
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx.Tx, http.StatusNotFound, nil, fmt.Errorf("cdn get: error getting cdn(s): %w", err))
		return
	}
	defer log.Close(rows, "unable to close DB connection")

	cdn := tc.CDNV5{}
	cdns := []tc.CDNV5{}
	for rows.Next() {
		if err = rows.Scan(&cdn.DNSSECEnabled, &cdn.DomainName, &cdn.ID, &cdn.LastUpdated, &cdn.TTLOverride, &cdn.Name); err != nil {
			api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("error getting cdn(s): %w", err))
			return
		}
		cdns = append(cdns, cdn)
	}

	api.WriteResp(w, r, cdns)
	return
}

func Create(w http.ResponseWriter, r *http.Request) {
	cdn := tc.CDNV5{}

	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	tx := inf.Tx.Tx

	cdn, err := validateRequest(r, inf.Version)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	var exists bool
	if err = tx.QueryRow("SELECT EXISTS(SELECT id from cdn where name = $1)", cdn.Name).Scan(&exists); err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error: %w, when checking if cdn with name %s exists", err, cdn.Name))
		return
	}
	if exists {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("cdn name '%s' already exists.", cdn.Name), nil)
		return
	}

	cdn.DomainName = strings.ToLower(cdn.DomainName)

	rows, err := inf.Tx.NamedQuery(insertQuery(inf.Version), cdn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("error: %w in creating cdn with name: %s", err, cdn.Name), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&cdn.ID, &cdn.LastUpdated); err != nil {
			usrErr, sysErr, code := api.ParseDBError(err)
			api.HandleErr(w, r, tx, code, usrErr, sysErr)
			return
		}
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "cdn was created.")
	w.Header().Set(rfc.Location, fmt.Sprintf("/api/%s/cdns?name=%s", inf.Version, cdn.Name))
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, cdn)
	changeLogMsg := fmt.Sprintf("CDN: %s, ID:%d, ACTION: Created cdn", cdn.Name, cdn.ID)
	api.CreateChangeLogRawTx(api.Created, changeLogMsg, inf.User, tx)
	return
}

func Update(w http.ResponseWriter, r *http.Request) {
	cdn := tc.CDNV5{}

	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	tx := inf.Tx.Tx

	cdn, err := validateRequest(r, inf.Version)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	id, err := strconv.Atoi(inf.Params["id"])
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	userErr, sysErr, errCode = dbhelpers.CheckIfCurrentUserCanModifyCDNWithID(tx, int64(id), inf.User.UserName)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	userErr, sysErr, errCode = api.CheckIfUnModified(r.Header, inf.Tx, id, "cdn")
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	cdn.DomainName = strings.ToLower(cdn.DomainName)

	query := `UPDATE
cdn SET
dnssec_enabled=$1,
domain_name=$2,
name=$3,
ttl_override=$4
WHERE id=$5 RETURNING last_updated, id`
	err = tx.QueryRow(query, cdn.DNSSECEnabled, cdn.DomainName, cdn.Name, cdn.TTLOverride, inf.Params["id"]).Scan(&cdn.LastUpdated, &cdn.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("cdn with id: %s not found", inf.Params["id"]), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "cdn was updated.")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, cdn)
	changeLogMsg := fmt.Sprintf("CDN: %s, ID:%d, ACTION: Updated cdn", cdn.Name, cdn.ID)
	api.CreateChangeLogRawTx(api.Updated, changeLogMsg, inf.User, tx)
	return
}

func Delete(w http.ResponseWriter, r *http.Request) {
	cdn := tc.CDNV5{}

	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	tx := inf.Tx.Tx

	var exists bool
	if err := tx.QueryRow("SELECT EXISTS(SELECT id from cdn where id = $1)", inf.Params["id"]).Scan(&exists); err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error: %w, when checking if cdn with name %s exists", err, cdn.Name))
		return
	}
	if !exists {
		api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("cdn id '%d' does not exist", cdn.ID), nil)
		return
	}

	id, err := strconv.Atoi(inf.Params["id"])
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	userErr, sysErr, errCode = dbhelpers.CheckIfCurrentUserCanModifyCDNWithID(tx, int64(id), inf.User.UserName)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	res, err := tx.Exec(`DELETE FROM cdn WHERE id=$1`, id)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	if rows, err := res.RowsAffected(); err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("unable to determine rows affected for deletion of cdn: %w", err))
		return
	} else if rows == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("no rows deleted for cdn"))
		return
	}

	api.WriteAlerts(w, r, http.StatusOK, tc.CreateAlerts(tc.SuccessLevel, "cdn was deleted."))
	changeLogMsg := fmt.Sprintf("ID:%d, ACTION: Deleted cdn", id)
	api.CreateChangeLogRawTx(api.Deleted, changeLogMsg, inf.User, tx)
	return
}
func validateRequest(r *http.Request, v *api.Version) (tc.CDNV5, error) {
	var cdn tc.CDNV5
	if err := json.NewDecoder(r.Body).Decode(&cdn); err != nil {
		return cdn, fmt.Errorf("error decoding POST request body into CDN struct %w", err)
	}

	validName := validation.NewStringRule(IsValidCDNName, "invalid characters found - Use alphanumeric . or - .")
	validDomainName := validation.NewStringRule(govalidator.IsDNSName, "not a valid domain name")
	errs := validation.Errors{
		"name":        validation.Validate(cdn.Name, validation.Required, validName),
		"domainName":  validation.Validate(cdn.DomainName, validation.Required, validDomainName),
		"ttlOverride": validation.Validate(cdn.TTLOverride, validation.By(tovalidate.IsGreaterThanZero)),
	}
	return cdn, util.JoinErrs(tovalidate.ToErrors(errs))
}

func SelectMaxLastUpdatedQuery(where string) string {
	return `SELECT max(t) from (
		SELECT max(last_updated) as t from cdn c  ` + where +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='cdn') as res`
}

func (cdn *TOCDN) GetLastUpdated() (*time.Time, bool, error) {
	return api.GetLastUpdated(cdn.APIInfo().Tx, *cdn.ID, "cdn")
}

func (cdn *TOCDN) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(last_updated) as t from ` + tableName + ` c ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='` + tableName + `') as res`
}

func (cdn *TOCDN) SetLastUpdated(t tc.TimeNoMod) { cdn.LastUpdated = &t }
func (cdn *TOCDN) InsertQuery() string           { return insertQuery(cdn.APIInfo().Version) }
func (cdn *TOCDN) NewReadObj() interface{}       { return &tc.CDNNullable{} }
func (cdn *TOCDN) SelectQuery() string           { return selectQuery(cdn.APIInfo().Version) }
func (cdn *TOCDN) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	columnInfo := map[string]dbhelpers.WhereColumnInfo{
		"domainName":    dbhelpers.WhereColumnInfo{Column: "domain_name"},
		"dnssecEnabled": dbhelpers.WhereColumnInfo{Column: "dnssec_enabled"},
		"id":            dbhelpers.WhereColumnInfo{Column: "id", Checker: api.IsInt},
		"name":          dbhelpers.WhereColumnInfo{Column: "name"},
	}
	if cdn.APIInfo().Version.GreaterThanOrEqualTo(&api.Version{Major: 4, Minor: 1}) {
		columnInfo["ttlOverride"] = dbhelpers.WhereColumnInfo{Column: "ttl_override", Checker: api.IsInt}
	}
	return columnInfo
}
func (cdn *TOCDN) UpdateQuery() string { return updateQuery(cdn.APIInfo().Version) }
func (cdn *TOCDN) DeleteQuery() string { return deleteQuery() }

func (cdn *TOCDN) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

// Implementation of the Identifier, Validator interface functions
func (cdn *TOCDN) GetKeys() (map[string]interface{}, bool) {
	if cdn.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *cdn.ID}, true
}

func (cdn *TOCDN) GetAuditName() string {
	if cdn.Name != nil {
		return *cdn.Name
	}
	if cdn.ID != nil {
		return strconv.Itoa(*cdn.ID)
	}
	return "0"
}

func (cdn *TOCDN) GetType() string {
	return "cdn"
}

func (cdn *TOCDN) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	cdn.ID = &i
}

// Validate fulfills the api.Validator interface.
func (cdn *TOCDN) Validate() (error, error) {
	validName := validation.NewStringRule(IsValidCDNName, "invalid characters found - Use alphanumeric . or - .")
	validDomainName := validation.NewStringRule(govalidator.IsDNSName, "not a valid domain name")
	errs := validation.Errors{
		"name":       validation.Validate(cdn.Name, validation.Required, validName),
		"domainName": validation.Validate(cdn.DomainName, validation.Required, validDomainName),
	}
	if cdn.APIInfo().Version.GreaterThanOrEqualTo(&api.Version{Major: 4, Minor: 1}) {
		errs["ttlOverride"] = validation.Validate(cdn.TTLOverride, validation.By(tovalidate.IsGreaterThanZero))
	}
	return util.JoinErrs(tovalidate.ToErrors(errs)), nil
}

func (cdn *TOCDN) Create() (error, error, int) {
	*cdn.DomainName = strings.ToLower(*cdn.DomainName)
	if cdn.APIInfo().Version.LessThan(&api.Version{Major: 4, Minor: 1}) {
		cdn.TTLOverride = nil
	}
	return api.GenericCreate(cdn)
}

func (cdn *TOCDN) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	api.DefaultSort(cdn.APIInfo(), "name")
	return api.GenericRead(h, cdn, useIMS)
}

func (cdn *TOCDN) Update(h http.Header) (error, error, int) {
	if cdn.ID != nil {
		userErr, sysErr, errCode := dbhelpers.CheckIfCurrentUserCanModifyCDNWithID(cdn.APIInfo().Tx.Tx, int64(*cdn.ID), cdn.APIInfo().User.UserName)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
	}
	*cdn.DomainName = strings.ToLower(*cdn.DomainName)
	if cdn.APIInfo().Version.LessThan(&api.Version{Major: 4, Minor: 1}) {
		cdn.TTLOverride = nil
	}
	return api.GenericUpdate(h, cdn)
}

func (cdn *TOCDN) Delete() (error, error, int) {
	if cdn.ID != nil {
		userErr, sysErr, errCode := dbhelpers.CheckIfCurrentUserCanModifyCDNWithID(cdn.APIInfo().Tx.Tx, int64(*cdn.ID), cdn.APIInfo().User.UserName)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
	}
	return api.GenericDelete(cdn)
}

func isValidCDNchar(r rune) bool {
	if r >= 'a' && r <= 'z' {
		return true
	}
	if r >= 'A' && r <= 'Z' {
		return true
	}
	if r >= '0' && r <= '9' {
		return true
	}
	if r == '.' || r == '-' {
		return true
	}
	return false
}

// IsValidCDNName returns true if the name contains only characters valid for a CDN name
func IsValidCDNName(str string) bool {
	i := strings.IndexFunc(str, func(r rune) bool { return !isValidCDNchar(r) })
	return i == -1
}

func formatQueryByAPIVersion(apiVersion *api.Version, minimumAPIVersion *api.Version, queryFormatString string, columnStrs []string, lowAPIVersionColumnStrs []string) string {
	if apiVersion.LessThan(&api.Version{Major: 4, Minor: 1}) {
		for index, _ := range columnStrs {
			columnStrs[index] = lowAPIVersionColumnStrs[index]
		}
	}
	columnStrArgs := make([]interface{}, len(columnStrs))
	for index, _ := range columnStrs {
		columnStrArgs[index] = columnStrs[index]
	}
	query := fmt.Sprintf(queryFormatString, columnStrArgs...)
	return query
}

func selectQuery(apiVersion *api.Version) string {
	query := `SELECT
dnssec_enabled,
domain_name,
id,
last_updated,
%s
name

FROM cdn c`
	return formatQueryByAPIVersion(apiVersion, &api.Version{Major: 4, Minor: 1}, query, []string{`
			ttl_override,
`}, []string{``})
}

func updateQuery(apiVersion *api.Version) string {
	query := `UPDATE
cdn SET
dnssec_enabled=:dnssec_enabled,
domain_name=:domain_name,
name=:name
%s
WHERE id=:id RETURNING last_updated`
	return formatQueryByAPIVersion(apiVersion, &api.Version{Major: 4, Minor: 1}, query, []string{`,
ttl_override=:ttl_override
`}, []string{``})
}

func insertQuery(apiVersion *api.Version) string {
	query := `INSERT INTO cdn (
dnssec_enabled,
domain_name,
name
%s
) VALUES (
:dnssec_enabled,
:domain_name,
:name
%s
) RETURNING id,last_updated`
	return formatQueryByAPIVersion(apiVersion, &api.Version{Major: 4, Minor: 1}, query, []string{`,
ttl_override
`, `,
:ttl_override
`}, []string{``, ``})
}

func deleteQuery() string {
	query := `DELETE FROM cdn
WHERE id=:id`
	return query
}
