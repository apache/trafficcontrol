package parameter

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
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/util/ims"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	NameQueryParam       = "name"
	SecureQueryParam     = "secure"
	ConfigFileQueryParam = "configFile"
	IDQueryParam         = "id"
	ValueQueryParam      = "value"
)

var (
	HiddenField = "********"
)

// we need a type alias to define functions on
type TOParameter struct {
	api.APIInfoImpl `json:"-"`
	tc.ParameterNullable
}

func (v *TOParameter) GetLastUpdated() (*time.Time, bool, error) {
	return api.GetLastUpdated(v.APIInfo().Tx, *v.ID, "parameter")
}

// AllowMultipleCreates indicates whether an array can be POSTed using the shared Create handler
func (v *TOParameter) AllowMultipleCreates() bool    { return true }
func (v *TOParameter) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = &t }
func (v *TOParameter) InsertQuery() string           { return insertQuery() }
func (v *TOParameter) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(p.last_updated) as t FROM parameter p
LEFT JOIN profile_parameter pp ON p.id = pp.parameter
LEFT JOIN profile pr ON pp.profile = pr.id ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='` + tableName + `') as res`
}

func (v *TOParameter) NewReadObj() interface{} { return &tc.ParameterNullable{} }
func (v *TOParameter) SelectQuery() string     { return selectQuery() }
func (v *TOParameter) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		ConfigFileQueryParam: {Column: "p.config_file"},
		IDQueryParam:         {Column: "p.id", Checker: api.IsInt},
		NameQueryParam:       {Column: "p.name"},
		SecureQueryParam:     {Column: "p.secure", Checker: api.IsBool},
		ValueQueryParam:      {Column: "p.value"},
	}
}
func (v *TOParameter) UpdateQuery() string { return updateQuery() }
func (v *TOParameter) DeleteQuery() string { return deleteQuery() }

func (param TOParameter) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: IDQueryParam, Func: api.GetIntKey}}
}

// Implementation of the Identifier, Validator interface functions
func (param TOParameter) GetKeys() (map[string]interface{}, bool) {
	if param.ID == nil {
		return map[string]interface{}{IDQueryParam: 0}, false
	}
	return map[string]interface{}{IDQueryParam: *param.ID}, true
}

func (param *TOParameter) SetKeys(keys map[string]interface{}) {
	i, _ := keys[IDQueryParam].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	param.ID = &i
}

func (param *TOParameter) GetAuditName() string {
	if param.Name != nil {
		return *param.Name
	}
	if param.ID != nil {
		return strconv.Itoa(*param.ID)
	}
	return "unknown"
}

func (param *TOParameter) GetType() string {
	return "param"
}

// Validate fulfills the api.Validator interface
func (param TOParameter) Validate() (error, error) {
	// Test
	// - Secure Flag is always set to either 1/0
	// - Admin rights only
	// - Do not allow duplicate parameters by name+config_file+value
	// - Client can send NOT NULL constraint on 'value' so removed it's validation as .Required
	errs := validation.Errors{
		NameQueryParam:       validation.Validate(param.Name, validation.Required),
		ConfigFileQueryParam: validation.Validate(param.ConfigFile, validation.Required),
	}
	if *param.ConfigFile == atscfg.ParentConfigFileName && *param.Name == atscfg.ParentConfigCacheParamWeight {
		errs[atscfg.ParentConfigFileName+" "+atscfg.ParentConfigCacheParamWeight] = validation.Validate(*param.Value, tovalidate.StringIsValidFloat())
	}

	return util.JoinErrs(tovalidate.ToErrors(errs)), nil
}

func (pa *TOParameter) Create() (error, error, int) {
	if pa.Value == nil {
		pa.Value = util.StrPtr("")
	}
	return api.GenericCreate(pa)
}

func (param *TOParameter) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	var maxTime time.Time
	var runSecond bool
	code := http.StatusOK
	queryParamsToQueryCols := param.ParamColumns()
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(param.APIInfo().Params, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest, nil
	}
	if useIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(param.APIInfo().Tx, h, queryValues, param.SelectMaxLastUpdatedQuery(where, orderBy, pagination, "parameter"))
		if !runSecond {
			log.Debugln("IMS HIT")
			return []interface{}{}, nil, nil, http.StatusNotModified, &maxTime
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}
	query := selectQuery() + where + ParametersGroupBy() + orderBy + pagination
	rows, err := param.ReqInfo.Tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, errors.New("querying " + param.GetType() + ": " + err.Error()), http.StatusInternalServerError, nil
	}
	defer rows.Close()

	params := []interface{}{}
	for rows.Next() {
		var p tc.ParameterNullable
		if err = rows.StructScan(&p); err != nil {
			return nil, nil, errors.New("scanning " + param.GetType() + ": " + err.Error()), http.StatusInternalServerError, nil
		}
		if p.Secure != nil && *p.Secure {
			if param.ReqInfo.Version.Major >= 5 {
				if !param.ReqInfo.User.Can("PARAMETER-SECURE:READ") {
					p.Value = &HiddenField
				}
			} else if param.ReqInfo.Version.Major == 4 {
				if param.ReqInfo.Config.RoleBasedPermissions {
					if !param.ReqInfo.User.Can("PARAMETER-SECURE:READ") {
						p.Value = &HiddenField
					}
				} else if param.ReqInfo.User.PrivLevel < auth.PrivLevelAdmin {
					p.Value = &HiddenField
				}
			} else if param.ReqInfo.User.PrivLevel < auth.PrivLevelAdmin {
				p.Value = &HiddenField
			}
		}
		params = append(params, p)
	}

	return params, nil, nil, code, &maxTime
}

func (pa *TOParameter) Update(h http.Header) (error, error, int) {
	if pa.Value == nil {
		pa.Value = util.StrPtr("")
	}
	return api.GenericUpdate(h, pa)
}

func (pa *TOParameter) Delete() (error, error, int) { return api.GenericDelete(pa) }

func insertQuery() string {
	query := `INSERT INTO parameter (
name,
config_file,
value,
secure) VALUES (
:name,
:config_file,
:value,
:secure) RETURNING id,last_updated`
	return query
}

func selectQuery() string {

	query := `SELECT
p.config_file,
p.id,
p.last_updated,
p.name,
p.value,
p.secure,
COALESCE(array_to_json(array_agg(pr.name) FILTER (WHERE pr.name IS NOT NULL)), '[]') AS profiles
FROM parameter p
LEFT JOIN profile_parameter pp ON p.id = pp.parameter
LEFT JOIN profile pr ON pp.profile = pr.id`
	return query
}

func updateQuery() string {
	query := `UPDATE
parameter SET
config_file=:config_file,
id=:id,
name=:name,
value=:value,
secure=:secure
WHERE id=:id RETURNING last_updated`
	return query
}

// ParametersGroupBy ...
func ParametersGroupBy() string {
	groupBy := ` GROUP BY p.config_file, p.id, p.last_updated, p.name, p.value, p.secure`
	return groupBy
}

func deleteQuery() string {
	query := `DELETE FROM parameter
WHERE id=:id`
	return query
}

func GetParameters(w http.ResponseWriter, r *http.Request) {
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
		ConfigFileQueryParam: {Column: "p.config_file"},
		IDQueryParam:         {Column: "p.id", Checker: api.IsInt},
		NameQueryParam:       {Column: "p.name"},
		SecureQueryParam:     {Column: "p.secure", Checker: api.IsBool},
		ValueQueryParam:      {Column: "p.value"},
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
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(tx, r.Header, queryValues, selectMaxLastUpdatedQuery(where))
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

	query := selectQuery() + where + ParametersGroupBy() + orderBy + pagination
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("Parameter read: error getting Parameter(s): %w", err))
		return
	}
	defer log.Close(rows, "unable to close DB connection")

	params := tc.ParameterNullableV5{}
	paramsList := []tc.ParameterNullableV5{}
	for rows.Next() {
		if err = rows.Scan(&params.ConfigFile, &params.ID, &params.LastUpdated, &params.Name, &params.Value, &params.Secure, &params.Profiles); err != nil {
			api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("error getting parameter(s): %w", err))
			return
		}
		if params.Secure != nil && *params.Secure {
			if inf.Version.Major >= 4 &&
				inf.Config.RoleBasedPermissions &&
				!inf.User.Can("PARAMETER-SECURE:READ") {
				params.Value = &HiddenField
			} else if inf.User.PrivLevel < auth.PrivLevelAdmin {
				params.Value = &HiddenField
			}
		}

		paramsList = append(paramsList, params)
	}

	api.WriteResp(w, r, paramsList)
	return
}

func CreateParameter(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	tx := inf.Tx.Tx

	body, err := io.ReadAll(r.Body)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("error reading request body"), nil)
		return
	}
	defer r.Body.Close()

	// Initial Unmarshal to validate request body
	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("invalid request format"), nil)
		return
	}

	// This code block decides if the request body is a slice of parameters or a single object and unmarshal it.
	var params []tc.ParameterV5
	if err := json.Unmarshal(body, &params); err != nil {
		var param tc.ParameterV5
		if err := json.Unmarshal(body, &param); err != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("error unmarshalling single object"), nil)
			return
		}
		params = append(params, param)
	}
	// Validate all objects of the every parameter from the request slice
	for _, parameter := range params {
		readValErr := validateRequestParameter(parameter)
		if readValErr != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
			return
		}
	}

	// Create all parameters from the request slice
	var objParams []tc.ParameterV5
	var objParam tc.ParameterV5
	for _, parameter := range params {
		query := `
			INSERT INTO parameter (
			    name,
			    config_file,
			    value,
			    secure
			    ) VALUES (
			        $1, $2, $3, $4
			    ) RETURNING id, name, config_file, value, last_updated, secure 
		`
		err = tx.QueryRow(
			query,
			parameter.Name,
			parameter.ConfigFile,
			parameter.Value,
			parameter.Secure,
		).Scan(

			&objParam.ID,
			&objParam.Name,
			&objParam.ConfigFile,
			&objParam.Value,
			&objParam.LastUpdated,
			&objParam.Secure,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("error: %w in parameter  with name: %s", err, parameter.Name), nil)
				return
			}
			usrErr, sysErr, code := api.ParseDBError(err)
			api.HandleErr(w, r, tx, code, usrErr, sysErr)
			return
		}

		objParams = append(objParams, objParam)
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "All Requested Parameters were created.")
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, objParam)
	for _, param := range params {
		changeLogMsg := fmt.Sprintf("PARAMETER: %s, ID:%d, ACTION: Created parameter", param.Name, param.ID)
		api.CreateChangeLogRawTx(api.Created, changeLogMsg, inf.User, tx)
	}
	return
}

func UpdateParameter(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx
	parameter, readValErr := readAndValidateJsonStruct(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	requestedID := inf.Params["id"]
	intRequestId := inf.IntParams["id"]

	// check if the entity was already updated
	userErr, sysErr, errCode = api.CheckIfUnModified(r.Header, inf.Tx, intRequestId, "parameter")
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	//update name and description of a phys_location
	query := `
	UPDATE parameter p
	SET
		config_file = $1,
		name = $2,
		value = $3,
		secure = $4
	WHERE
		p.id = $5
	RETURNING
		p.id,
		p.last_updated
`

	err := tx.QueryRow(
		query,
		parameter.ConfigFile,
		parameter.Name,
		parameter.Value,
		parameter.Secure,
		requestedID,
	).Scan(
		&parameter.ID,
		&parameter.LastUpdated,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("parameter with ID: %v not found", parameter.ID), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "parameter was updated")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, parameter)
	changeLogMsg := fmt.Sprintf("PARAMETER: %s, ID:%d, ACTION: Updated parameter", parameter.Name, parameter.ID)
	api.CreateChangeLogRawTx(api.Updated, changeLogMsg, inf.User, tx)
	return
}

func DeleteParameter(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	id := inf.Params["id"]
	if id == "" {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("couldn't delete Parameter. Invalid ID. Id Cannot be empty for Delete Operation"), nil)
		return
	}
	exists, err := dbhelpers.ParameterExists(tx, id)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	if !exists {
		api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("no Parameter exists by id: %s", id), nil)
		return
	}

	res, err := tx.Exec("DELETE FROM parameter AS p WHERE p.id=$1", id)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error determining rows affected for delete parameter: %w", err))
		return
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("no rows deleted for parameter"), nil)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "parameter"+
		" was deleted.")
	api.WriteAlerts(w, r, http.StatusOK, alerts)
	changeLogMsg := fmt.Sprintf("ID:%s, ACTION: Deleted parameter", id)
	api.CreateChangeLogRawTx(api.Deleted, changeLogMsg, inf.User, tx)
	return
}

// selectMaxLastUpdatedQuery used for TryIfModifiedSinceQuery()
func selectMaxLastUpdatedQuery(where string) string {
	return `
SELECT max(t) from (
	SELECT max(p.last_updated) as t FROM parameter p
	LEFT JOIN profile_parameter pp ON p.id = pp.parameter
	LEFT JOIN profile pr ON pp.profile = pr.id ` + where +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='parameter'
) as res`
}

func readAndValidateJsonStruct(r *http.Request) (tc.ParameterV5, error) {
	var parameter tc.ParameterV5
	if err := json.NewDecoder(r.Body).Decode(&parameter); err != nil {
		userErr := fmt.Errorf("error decoding POST request body into ParameterV5 struct %w", err)
		return parameter, userErr
	}

	// validate JSON body
	errs := make(map[string]error)

	errs[NameQueryParam] = validation.Validate(parameter.Name, validation.Required)
	errs[ConfigFileQueryParam] = validation.Validate(parameter.ConfigFile, validation.Required)

	if parameter.ConfigFile == atscfg.ParentConfigFileName && parameter.Name == atscfg.ParentConfigCacheParamWeight {
		key := atscfg.ParentConfigFileName + " " + atscfg.ParentConfigCacheParamWeight
		errs[key] = validation.Validate(parameter.Value, tovalidate.StringIsValidFloat())
	}

	if len(errs) > 0 {
		var errorSlice []error
		for _, err := range errs {
			errorSlice = append(errorSlice, err)
		}
		userErr := util.JoinErrs(errorSlice)
		return parameter, userErr
	}
	return parameter, nil
}

// Unlike usual readAndValidateJsonStruct function, this function does not decode the JSON data
// but only validate the JSON objects
func validateRequestParameter(parameter tc.ParameterV5) error {
	errs := make(map[string]error)

	errs[NameQueryParam] = validation.Validate(parameter.Name, validation.Required)
	errs[ConfigFileQueryParam] = validation.Validate(parameter.ConfigFile, validation.Required)

	if parameter.ConfigFile == atscfg.ParentConfigFileName && parameter.Name == atscfg.ParentConfigCacheParamWeight {
		key := atscfg.ParentConfigFileName + " " + atscfg.ParentConfigCacheParamWeight
		errs[key] = validation.Validate(parameter.Value, tovalidate.StringIsValidFloat())
	}

	if len(errs) > 0 {
		var errorSlice []error
		for _, err := range errs {
			errorSlice = append(errorSlice, err)
		}
		userErr := util.JoinErrs(errorSlice)
		return userErr
	}
	return nil
}
