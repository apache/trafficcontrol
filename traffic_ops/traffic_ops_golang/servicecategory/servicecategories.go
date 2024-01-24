package servicecategory

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
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/util/ims"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
)

type TOServiceCategory struct {
	api.APIInfoImpl `json:"-"`
	tc.ServiceCategory
}

func (v *TOServiceCategory) GetLastUpdated() (*time.Time, bool, error) {
	return api.GetLastUpdatedByName(v.APIInfo().Tx, v.Name, "service_category")
}

func (v *TOServiceCategory) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = t }
func (v *TOServiceCategory) InsertQuery() string           { return insertQuery() }
func (v *TOServiceCategory) NewReadObj() interface{}       { return &tc.ServiceCategory{} }
func (v *TOServiceCategory) SelectQuery() string           { return selectQuery() }
func (v *TOServiceCategory) UpdateQuery() string           { return updateQuery() }
func (v *TOServiceCategory) DeleteQuery() string           { return deleteQuery() }

func (serviceCategory TOServiceCategory) GetAuditName() string {
	if serviceCategory.Name != "" {
		return serviceCategory.Name
	}
	return "unknown"
}

func (serviceCategory TOServiceCategory) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "name", Func: api.GetStringKey}}
}

// Implementation of the Identifier, Validator interface functions
func (serviceCategory TOServiceCategory) GetKeys() (map[string]interface{}, bool) {
	if serviceCategory.Name == "" {
		return map[string]interface{}{"name": ""}, false
	}
	return map[string]interface{}{"name": serviceCategory.Name}, true
}

func (serviceCategory *TOServiceCategory) SetKeys(keys map[string]interface{}) {
	n, _ := keys["name"].(string)
	serviceCategory.Name = n
}

func (serviceCategory TOServiceCategory) GetType() string {
	return "serviceCategory"
}

func (serviceCategory *TOServiceCategory) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"name": dbhelpers.WhereColumnInfo{Column: "sc.name"},
	}
}

func (serviceCategory *TOServiceCategory) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(last_updated) as t from service_category sc ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='service_category') as res`
}

func (serviceCategory TOServiceCategory) Validate() (error, error) {
	nameRule := validation.NewStringRule(tovalidate.IsAlphanumericDash, "must consist of only alphanumeric or dash characters.")
	errs := validation.Errors{
		"name": validation.Validate(serviceCategory.Name, validation.Required, nameRule),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs)), nil
}

func (serviceCategory *TOServiceCategory) Create() (error, error, int) {
	return api.GenericCreateNameBasedID(serviceCategory)
}

func (serviceCategory *TOServiceCategory) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	api.DefaultSort(serviceCategory.APIInfo(), "name")
	serviceCategories, userErr, sysErr, errCode, maxTime := api.GenericRead(h, serviceCategory, useIMS)
	if userErr != nil || sysErr != nil {
		return nil, userErr, sysErr, errCode, nil
	}

	return serviceCategories, nil, nil, errCode, maxTime
}

func Update(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"name"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	name := inf.Params["name"]

	var newSC TOServiceCategory
	if err := json.NewDecoder(r.Body).Decode(&newSC); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
		return
	}

	if userErr, sysErr := newSC.Validate(); userErr != nil || sysErr != nil {
		code := http.StatusBadRequest
		if sysErr != nil {
			code = http.StatusInternalServerError
		}
		api.HandleErr(w, r, inf.Tx.Tx, code, userErr, sysErr)
		return
	}

	var origSC TOServiceCategory
	if err := inf.Tx.QueryRow(`SELECT name, last_updated FROM service_category WHERE name = $1`, name).Scan(&origSC.Name, &origSC.LastUpdated); err != nil {
		if err == sql.ErrNoRows {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("no service category found with name "+name), nil)
			return
		}
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	if !api.IsUnmodified(r.Header, origSC.LastUpdated.Time) {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusPreconditionFailed, errors.New("service category could not be modified because the precondition failed"), nil)
		return
	}

	resp, err := inf.Tx.Tx.Exec(updateQuery(), newSC.Name, name)
	if err != nil {
		userErr, sysErr, errCode = api.ParseDBError(err)
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	api.CreateChangeLogRawTx(api.ApiChange, api.Updated+" Service Category from "+name+" to "+newSC.Name, inf.User, inf.Tx.Tx)
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Service Category update from "+name+" to "+newSC.Name+" was successful.", resp)
}

func (serviceCategory *TOServiceCategory) Delete() (error, error, int) {
	return api.GenericDelete(serviceCategory)
}

func insertQuery() string {
	return `INSERT INTO service_category (name) VALUES (:name) RETURNING name, last_updated`
}

func selectQuery() string {
	return `SELECT
sc.last_updated,
sc.name
FROM service_category as sc`
}

func updateQuery() string {
	return `UPDATE
service_category SET
name=$1
WHERE name=$2 RETURNING last_updated`
}

func deleteQuery() string {
	return `DELETE FROM service_category WHERE name=:name`
}

// Get [Version : V5] function Process the *http.Request and writes the response. It uses GetServiceCategory function.
func Get(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	code := http.StatusOK
	useIMS := false
	config, e := api.GetConfig(r.Context())
	if e == nil && config != nil {
		useIMS = config.UseIMS
	} else {
		log.Warnf("Couldn't get config %v", e)
	}

	var maxTime time.Time
	var usrErr error
	var syErr error

	var scList []tc.ServiceCategoryV5

	tx := inf.Tx

	scList, maxTime, code, usrErr, syErr = GetServiceCategory(tx, inf.Params, useIMS, r.Header)
	if code == http.StatusNotModified {
		w.WriteHeader(code)
		api.WriteResp(w, r, []tc.ServiceCategoryV5{})
		return
	}

	if code == http.StatusBadRequest {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, usrErr, nil)
		return
	}

	if sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, syErr)
		return
	}

	if maxTime != (time.Time{}) && api.SetLastModifiedHeader(r, useIMS) {
		api.AddLastModifiedHdr(w, maxTime)
	}

	api.WriteResp(w, r, scList)
}

// GetServiceCategory [Version : V5] receives transactions from Get function and returns service_categories list.
func GetServiceCategory(tx *sqlx.Tx, params map[string]string, useIMS bool, header http.Header) ([]tc.ServiceCategoryV5, time.Time, int, error, error) {
	var runSecond bool
	var maxTime time.Time
	scList := []tc.ServiceCategoryV5{}

	selectQuery := `SELECT name, last_updated FROM service_category as sc`

	// Query Parameters to Database Query column mappings
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"name": {Column: "sc.name", Checker: nil},
	}
	if _, ok := params["orderby"]; !ok {
		params["orderby"] = "name"
	}
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(params, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, time.Time{}, http.StatusBadRequest, util.JoinErrs(errs), nil
	}

	if useIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(tx, header, queryValues, SelectMaxLastUpdatedQuery(where))
		if !runSecond {
			log.Debugln("IMS HIT")
			return scList, maxTime, http.StatusNotModified, nil, nil
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}
	query := selectQuery + where + orderBy + pagination
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, time.Time{}, http.StatusInternalServerError, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		sc := tc.ServiceCategoryV5{}
		if err = rows.Scan(&sc.Name, &sc.LastUpdated); err != nil {
			return nil, time.Time{}, http.StatusInternalServerError, nil, err
		}
		scList = append(scList, sc)
	}

	return scList, maxTime, http.StatusOK, nil, nil
}

// CreateServiceCategory [Version : V5] function creates the service category with the passed name.
func CreateServiceCategory(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	tx := inf.Tx.Tx

	sc, readValErr := readAndValidateJsonStruct(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	// check if service category already exists
	var exists bool
	err := tx.QueryRow(`SELECT EXISTS(SELECT * from service_category where name = $1)`, sc.Name).Scan(&exists)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error: %w, when checking if service category with name %s exists", err, sc.Name))
		return
	}
	if exists {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("service category name '%s' already exists.", sc.Name), nil)
		return
	}

	// create service category
	query := `INSERT INTO service_category (name) VALUES ($1) RETURNING name, last_updated`
	err = tx.QueryRow(query, sc.Name).Scan(&sc.Name, &sc.LastUpdated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("error: %w in creating service category with name: %s", err, sc.Name), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "service category was created.")
	w.Header().Set(rfc.Location, fmt.Sprintf("/api/%s/service_category?name=%s", inf.Version, sc.Name))
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, sc)
	changeLogMsg := fmt.Sprintf("SERVICECATEGORY: %s ACTION: Created serviceCategory", sc.Name)
	api.CreateChangeLogRawTx(api.Created, changeLogMsg, inf.User, tx)
	return
}

// UpdateServiceCategory [Version : V5] function updates the name of the service category passed.
func UpdateServiceCategory(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx
	sc, readValErr := readAndValidateJsonStruct(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	requestedName := inf.Params["name"]
	// check if the entity was already updated
	userErr, sysErr, errCode = api.CheckIfUnModifiedByName(r.Header, inf.Tx, requestedName, "service_category")
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	//update name of a service category
	query := `UPDATE service_category sc SET
		name = $1
	WHERE sc.name = $2
	RETURNING sc.name, sc.last_updated`

	err := tx.QueryRow(query, sc.Name, requestedName).Scan(&sc.Name, &sc.LastUpdated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("service category with name: %s not found", requestedName), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "service category was updated")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, sc)
	changeLogMsg := fmt.Sprintf("SERVICECATEGORY: %s, ACTION: Updated serviceCategory", sc.Name)
	api.CreateChangeLogRawTx(api.Updated, changeLogMsg, inf.User, tx)
	return
}

// DeleteServiceCategory [Version : V5] function deletes the service category passed.
func DeleteServiceCategory(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	name := inf.Params["name"]
	exists, err := dbhelpers.ServiceCategoryExists(tx, name)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	if !exists {
		api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("no service category exists for name: %s", name), nil)
		return

	}

	assignedDeliveryService := 0
	if err := inf.Tx.Get(&assignedDeliveryService, "SELECT count(service_category) FROM deliveryservice d WHERE d.service_category=$1", name); err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("service category delete, counting assigned Delivery Service(s): %w", err))
		return
	} else if assignedDeliveryService != 0 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("can not delete a service category with %d assigned Delivery Service(s)", assignedDeliveryService), nil)
		return
	}

	res, err := tx.Exec("DELETE FROM service_category AS sc WHERE sc.name=$1", name)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("determining rows affected for delete service_category: %w", err))
		return
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("no rows deleted for service_category"))
		return
	}

	alertMessage := fmt.Sprintf("%s was deleted.", name)
	alerts := tc.CreateAlerts(tc.SuccessLevel, alertMessage)
	api.WriteAlerts(w, r, http.StatusOK, alerts)
	changeLogMsg := fmt.Sprintf("SERVICECATEGORY: %s, ACTION: Deleted serviceCategory", name)
	api.CreateChangeLogRawTx(api.Deleted, changeLogMsg, inf.User, tx)
	return
}

func readAndValidateJsonStruct(r *http.Request) (tc.ServiceCategoryV5, error) {
	var sc tc.ServiceCategoryV5
	if err := json.NewDecoder(r.Body).Decode(&sc); err != nil {
		userErr := fmt.Errorf("error decoding POST request body into ServiceCategoryV5 struct %w", err)
		return sc, userErr
	}

	// validate JSON body
	rule := validation.NewStringRule(tovalidate.IsAlphanumericUnderscoreDash, "must consist of only alphanumeric, dash, or underscore characters")
	errs := tovalidate.ToErrors(validation.Errors{
		"name": validation.Validate(sc.Name, validation.Required, rule),
	})
	if len(errs) > 0 {
		userErr := util.JoinErrs(errs)
		return sc, userErr
	}
	return sc, nil
}

// SelectMaxLastUpdatedQuery used for TryIfModifiedSinceQuery()
func SelectMaxLastUpdatedQuery(where string) string {
	return `SELECT max(t) from (
		SELECT max(last_updated) as t from service_category sc ` + where +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='service_category') as res`
}
