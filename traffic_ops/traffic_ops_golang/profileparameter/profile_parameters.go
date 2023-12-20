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
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/util/ims"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	ProfileIDQueryParam   = "profileId"
	ParameterIDQueryParam = "parameterId"
)

// we need a type alias to define functions on
type TOProfileParameter struct {
	api.APIInfoImpl `json:"-"`
	tc.ProfileParameterNullable
}

// AllowMultipleCreates indicates whether an array can be POSTed using the shared Create handler
func (v *TOProfileParameter) AllowMultipleCreates() bool { return true }
func (v *TOProfileParameter) NewReadObj() interface{}    { return &tc.ProfileParametersNullable{} }
func (v *TOProfileParameter) SelectQuery() string        { return selectQuery() }
func (v *TOProfileParameter) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"profileId":   dbhelpers.WhereColumnInfo{Column: "pp.profile"},
		"parameterId": dbhelpers.WhereColumnInfo{Column: "pp.parameter"},
		"lastUpdated": dbhelpers.WhereColumnInfo{Column: "pp.last_updated"},
	}
}
func (v *TOProfileParameter) DeleteQuery() string { return deleteQuery() }

func (pp TOProfileParameter) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: ProfileIDQueryParam, Func: api.GetIntKey}, {Field: ParameterIDQueryParam, Func: api.GetIntKey}}
}

// Implementation of the Identifier, Validator interface functions
func (pp TOProfileParameter) GetKeys() (map[string]interface{}, bool) {
	if pp.ProfileID == nil {
		return map[string]interface{}{ProfileIDQueryParam: 0}, false
	}
	if pp.ParameterID == nil {
		return map[string]interface{}{ParameterIDQueryParam: 0}, false
	}
	keys := make(map[string]interface{})
	profileID := *pp.ProfileID
	parameterID := *pp.ParameterID

	keys[ProfileIDQueryParam] = profileID
	keys[ParameterIDQueryParam] = parameterID
	return keys, true
}

func (pp *TOProfileParameter) GetAuditName() string {
	if pp.ProfileID != nil {
		return strconv.Itoa(*pp.ProfileID) + "-" + strconv.Itoa(*pp.ParameterID)
	}
	return "unknown"
}

func (pp *TOProfileParameter) GetType() string {
	return "profileParameter"
}

func (pp *TOProfileParameter) SetKeys(keys map[string]interface{}) {
	profId, _ := keys[ProfileIDQueryParam].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	pp.ProfileID = &profId

	paramId, _ := keys[ParameterIDQueryParam].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	pp.ParameterID = &paramId
}

// Validate fulfills the api.Validator interface.
func (pp *TOProfileParameter) Validate() (error, error) {

	errs := validation.Errors{
		"profileId":   validation.Validate(pp.ProfileID, validation.Required),
		"parameterId": validation.Validate(pp.ParameterID, validation.Required),
	}

	return util.JoinErrs(tovalidate.ToErrors(errs)), nil
}

// The TOProfileParameter implementation of the Creator interface
// all implementations of Creator should use transactions and return the proper errorType
// ParsePQUniqueConstraintError is used to determine if a profileparameter with conflicting values exists
// if so, it will return an errorType of DataConflict and the type should be appended to the
// generic error message returned
// The insert sql returns the profile and lastUpdated values of the newly inserted profileparameter and have
// to be added to the struct
func (pp *TOProfileParameter) Create() (error, error, int) {
	if pp.ProfileID != nil {
		cdnName, err := dbhelpers.GetCDNNameFromProfileID(pp.ReqInfo.Tx.Tx, *pp.ProfileID)
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
		userErr, sysErr, errCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(pp.ReqInfo.Tx.Tx, string(cdnName), pp.ReqInfo.User.UserName)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
	} else {
		return errors.New("no profile ID in request"), nil, http.StatusBadRequest
	}
	resultRows, err := pp.APIInfo().Tx.NamedQuery(insertQuery(), pp)
	if err != nil {
		return api.ParseDBError(err)
	}
	defer resultRows.Close()

	var profile int
	var parameter int
	var lastUpdated tc.TimeNoMod
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&profile, &parameter, &lastUpdated); err != nil {
			return nil, errors.New("profileparameter create scanning: " + err.Error()), http.StatusInternalServerError
		}
	}
	if rowsAffected == 0 {
		return nil, errors.New("profileparameter create returned no rows"), http.StatusInternalServerError
	}
	if rowsAffected > 1 {
		return nil, errors.New("profileparameter create returned multiple rows"), http.StatusInternalServerError
	}

	pp.SetKeys(map[string]interface{}{ProfileIDQueryParam: profile, ParameterIDQueryParam: parameter})
	return nil, nil, http.StatusOK
}

func insertQuery() string {
	return `INSERT INTO profile_parameter (
profile,
parameter) VALUES (
:profile_id,
:parameter_id) RETURNING profile, parameter, last_updated`
}

func (pp *TOProfileParameter) Update(h http.Header) (error, error, int) {
	return nil, nil, http.StatusNotImplemented
}
func (pp *TOProfileParameter) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	api.DefaultSort(pp.APIInfo(), "parameter")
	return api.GenericRead(h, pp, useIMS)
}
func (pp *TOProfileParameter) Delete() (error, error, int) {
	if pp.ProfileID != nil {
		cdnName, err := dbhelpers.GetCDNNameFromProfileID(pp.ReqInfo.Tx.Tx, *pp.ProfileID)
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
		userErr, sysErr, errCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(pp.ReqInfo.Tx.Tx, string(cdnName), pp.ReqInfo.User.UserName)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
	} else {
		return errors.New("no profile ID in request"), nil, http.StatusBadRequest
	}
	return api.GenericDelete(pp)
}
func (v *TOProfileParameter) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(pp.last_updated) as t FROM profile_parameter pp
JOIN profile prof ON prof.id = pp.profile
JOIN parameter param ON param.id = pp.parameter ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='profile_parameter') as res`
}

func selectQuery() string {

	query := `SELECT
pp.last_updated,
pp.parameter parameter_id,
prof.name profile
FROM profile_parameter pp
JOIN profile prof ON prof.id = pp.profile
JOIN parameter param ON param.id = pp.parameter`
	return query
}

func deleteQuery() string {
	query := `DELETE FROM profile_parameter
	WHERE profile=:profile_id and parameter=:parameter_id`
	return query
}

func GetProfileParameter(w http.ResponseWriter, r *http.Request) {
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
		"profileId":   {Column: "pp.profile"},
		"parameterId": {Column: "pp.parameter"},
		"lastUpdated": {Column: "pp.last_updated"},
	}
	if _, ok := inf.Params["orderby"]; !ok {
		inf.Params["orderby"] = "parameter"
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

	query := selectQuery() + where + orderBy + pagination
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("Profile Parameter read: error getting Profile Parameter(s): %w", err))
		return
	}
	defer log.Close(rows, "unable to close DB connection")

	profileParams := tc.ProfileParametersNullableV5{}
	profileParamsList := []tc.ProfileParametersNullableV5{}
	for rows.Next() {
		if err = rows.Scan(&profileParams.LastUpdated, &profileParams.Parameter, &profileParams.Profile); err != nil {
			api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("error getting profile parameter(s): %w", err))
			return
		}

		profileParamsList = append(profileParamsList, profileParams)
	}

	api.WriteResp(w, r, profileParamsList)
	return
}

func CreateProfileParameter(w http.ResponseWriter, r *http.Request) {
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
	var profileParams []tc.ProfileParameterCreationRequest
	if err := json.Unmarshal(body, &profileParams); err != nil {
		var profileParam tc.ProfileParameterCreationRequest
		if err := json.Unmarshal(body, &profileParam); err != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("error unmarshalling single object"), nil)
			return
		}
		profileParams = append(profileParams, profileParam)
	}

	// Validate all objects of the every profile parameter from the request slice
	for _, profileParameter := range profileParams {
		readValErr := validateRequestProfileParameter(profileParameter)
		if readValErr != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
			return
		}
	}

	// Check user Permissions on all Profiles requested
	for _, profileParameter := range profileParams {
		cdnName, err := dbhelpers.GetCDNNameFromProfileID(tx, profileParameter.ProfileID)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, err, nil)
			return
		}
		userErr, sysErr, errCode = dbhelpers.CheckIfCurrentUserCanModifyCDN(tx, string(cdnName), inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, tx, errCode, userErr, sysErr)
			return
		}
	}

	// Check if any of the profile parameter from the request slice already exists
	for _, profileParameter := range profileParams {
		var count int
		err = tx.QueryRow("SELECT count(*) from profile_parameter where profile = $1 and parameter = $2", profileParameter.ProfileID, profileParameter.ParameterID).Scan(&count)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error: %w, when checking if profile parameter with profile_id %d and parameter_id %d exists", err, profileParameter.ProfileID, profileParameter.ParameterID))
			return
		}
		if count == 1 {
			api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("profile parameter with profile_id %d and parameter_id %d already exists", profileParameter.ProfileID, profileParameter.ParameterID), nil)
			return
		}
	}

	// Create all profile parameters from the request slice
	var objProfileParams []tc.ProfileParameterV5
	for _, profileParameter := range profileParams {
		query := `
				INSERT INTO profile_parameter (
				profile, 
				parameter
				) VALUES (
				$1, $2
				) RETURNING profile, parameter, last_updated
`
		var objProfileParam tc.ProfileParameterV5
		err = tx.QueryRow(
			query,
			profileParameter.ProfileID,
			profileParameter.ParameterID,
		).Scan(
			&objProfileParam.ProfileID,
			&objProfileParam.ParameterID,
			&objProfileParam.LastUpdated,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("error: %w in profile_parameter with with profile_id %d and parameter_id %d", err, profileParameter.ProfileID, profileParameter.ParameterID), nil)
				return
			}
			usrErr, sysErr, code := api.ParseDBError(err)
			api.HandleErr(w, r, tx, code, usrErr, sysErr)
			return
		}

		// Fetch the Profile Name from ID to insert in type ProfileParameterV5
		profileName, ok, err := dbhelpers.GetProfileNameFromID(profileParameter.ProfileID, tx)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting profile name from id: "+err.Error()))
			return
		} else if !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("profile not found"), nil)
			return
		}

		// Fetch the Parameter Name from ID to insert in type ProfileParameterV5
		parameterName, ok, err := dbhelpers.GetParamNameByID(tx, profileParameter.ParameterID)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting parameter name from id: "+err.Error()))
			return
		} else if !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("parameter not found"), nil)
			return
		}

		objProfileParam.Profile = profileName
		objProfileParam.Parameter = parameterName
		objProfileParams = append(objProfileParams, objProfileParam)
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "All Requested ProfileParameters were created.")
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, objProfileParams)
	for _, profileParam := range profileParams {
		changeLogMsg := fmt.Sprintf("PROFILEPARAMETER Profile ID: %d, ParameterID:%d, ACTION: Created profileParameter", profileParam.ProfileID, profileParam.ParameterID)
		api.CreateChangeLogRawTx(api.Created, changeLogMsg, inf.User, tx)
	}
	return

}

func DeleteProfileParameter(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"profileId"})
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	profileID := inf.Params["profileId"]
	parameterID := inf.Params["parameterId"]
	intProfileID := inf.IntParams["profileId"]

	if profileID == "" || parameterID == "" {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("couldn't delete Profile_Parameter. profileID  & parameterID Cannot be empty for Delete Operation"), nil)
		return
	}

	cdnName, err := dbhelpers.GetCDNNameFromProfileID(tx, intProfileID)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, err, nil)
		return
	}
	userErr, sysErr, errCode = dbhelpers.CheckIfCurrentUserCanModifyCDN(tx, string(cdnName), inf.User.UserName)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	exists, err := dbhelpers.ProfileParameterExists(tx, profileID, parameterID)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	if !exists {
		api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("no profile_parameter exists by profile_id: %s & parameter_id: %s", profileID, parameterID), nil)
		return
	}

	res, err := tx.Exec("DELETE FROM profile_parameter AS pp WHERE pp.profile=$1 and parameter=$2", profileID, parameterID)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error determining rows affected for delete profile_parameter: %w", err))
		return
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("no rows deleted for profile_parameter"), nil)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "profile_parameter"+
		" was deleted.")
	api.WriteAlerts(w, r, http.StatusOK, alerts)
	changeLogMsg := fmt.Sprintf("PROFILEPARAMETER Profile ID: %s, ParameterID:%s, ACTION: Deleted profileParameter", profileID, parameterID)
	api.CreateChangeLogRawTx(api.Deleted, changeLogMsg, inf.User, tx)
	return
}

func selectMaxLastUpdatedQuery(where string) string {
	return `
        SELECT max(t) from (
            SELECT max(pp.last_updated) as t FROM profile_parameter pp
            JOIN profile prof ON prof.id = pp.profile
            JOIN parameter param ON param.id = pp.parameter ` + where +
		` UNION ALL
            SELECT max(last_updated) as t FROM last_deleted l WHERE l.table_name = 'profile_parameter'
        ) as res
    `
}

// validateRequestProfileParameter validate the JSON objects
func validateRequestProfileParameter(profileParameter tc.ProfileParameterCreationRequest) error {
	errs := make(map[string]error)

	errs[ProfileIDQueryParam] = validation.Validate(profileParameter.ProfileID, validation.Required)
	errs[ParameterIDQueryParam] = validation.Validate(profileParameter.ParameterID, validation.Required)

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
