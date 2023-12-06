// Package profile includes logic and handlers for Profile-related API
// endpoints, including /profiles, /profiles/name/{{name}}/parameters,
// /profiles/{{ID}}/parameters,
// /profiles/name/{{New Profile Name}}/copy/{{existing Profile Name}},
// /profiles/import, and /profiles/{{ID}}/export.
package profile

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
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/parameter"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/util/ims"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
)

// Supported (non-pagination) query string parameters for /profiles.
const (
	CDNQueryParam         = "cdn"
	DescriptionQueryParam = "description"
	IDQueryParam          = "id"
	NameQueryParam        = "name"
	ParamQueryParam       = "param"
	TypeQueryParam        = "type"
)

// we need a type alias to define functions on
type TOProfile struct {
	api.APIInfoImpl `json:"-"`
	tc.ProfileNullable
}

func (v *TOProfile) GetLastUpdated() (*time.Time, bool, error) {
	return api.GetLastUpdated(v.APIInfo().Tx, *v.ID, "profile")
}

func (v *TOProfile) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = &t }
func (v *TOProfile) InsertQuery() string           { return insertQuery() }
func (v *TOProfile) UpdateQuery() string           { return updateQuery() }
func (v *TOProfile) DeleteQuery() string           { return deleteQuery() }

func (prof TOProfile) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: IDQueryParam, Func: api.GetIntKey}}
}

// Implementation of the Identifier, Validator interface functions
func (prof TOProfile) GetKeys() (map[string]interface{}, bool) {
	if prof.ID == nil {
		return map[string]interface{}{IDQueryParam: 0}, false
	}
	return map[string]interface{}{IDQueryParam: *prof.ID}, true
}

func (prof *TOProfile) SetKeys(keys map[string]interface{}) {
	i, _ := keys[IDQueryParam].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	prof.ID = &i
}

func (prof *TOProfile) GetAuditName() string {
	if prof.Name != nil {
		return *prof.Name
	}
	if prof.ID != nil {
		return strconv.Itoa(*prof.ID)
	}
	return "unknown"
}

func (prof *TOProfile) GetType() string {
	return "profile"
}

func (prof *TOProfile) Validate() (error, error) {
	errs := validation.Errors{
		NameQueryParam: validation.Validate(prof.Name, validation.By(
			func(value interface{}) error {
				name, ok := value.(*string)
				if !ok {
					return fmt.Errorf("wrong type, need: string, got: %T", value)
				}
				if name == nil || *name == "" {
					return errors.New("required and cannot be blank")
				}
				if strings.Contains(*name, " ") {
					return errors.New("cannot contain spaces")
				}
				return nil
			},
		)),
		DescriptionQueryParam: validation.Validate(prof.Description, validation.Required),
		CDNQueryParam:         validation.Validate(prof.CDNID, validation.Required),
		TypeQueryParam:        validation.Validate(prof.Type, validation.Required),
	}
	if errs != nil {
		return util.JoinErrs(tovalidate.ToErrors(errs)), nil
	}
	return nil, nil
}

func (prof *TOProfile) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	var maxTime time.Time
	var runSecond bool
	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		CDNQueryParam:   dbhelpers.WhereColumnInfo{Column: "c.id", Checker: api.IsInt},
		NameQueryParam:  dbhelpers.WhereColumnInfo{Column: "prof.name"},
		IDQueryParam:    dbhelpers.WhereColumnInfo{Column: "prof.id", Checker: api.IsInt},
		ParamQueryParam: dbhelpers.WhereColumnInfo{Column: "pp.parameter", Checker: api.IsInt},
	}
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(prof.APIInfo().Params, queryParamsToQueryCols)

	query := selectProfilesQuery()
	// Narrow down if the query parameter is 'param'
	// TODO add generic where clause to api.GenericRead
	if paramValue, ok := prof.APIInfo().Params[ParamQueryParam]; ok {
		if len(paramValue) > 0 {
			query += " LEFT JOIN profile_parameter pp ON prof.id = pp.profile"
		}
	}

	if len(errs) > 0 {
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest, nil
	}

	if useIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(prof.APIInfo().Tx, h, queryValues, selectMaxLastUpdatedQuery(where))
		if !runSecond {
			log.Debugln("IMS HIT")
			return []interface{}{}, nil, nil, http.StatusNotModified, &maxTime
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}

	query += where + orderBy + pagination
	log.Debugln("Query is ", query)

	rows, err := prof.ReqInfo.Tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, errors.New("profile read querying: " + err.Error()), http.StatusInternalServerError, nil
	}
	defer rows.Close()

	profiles := []tc.ProfileNullable{}

	for rows.Next() {
		var p tc.ProfileNullable
		if err = rows.StructScan(&p); err != nil {
			return nil, nil, errors.New("profile read scanning: " + err.Error()), http.StatusInternalServerError, nil
		}
		profiles = append(profiles, p)
	}
	rows.Close()
	profileInterfaces := []interface{}{}
	canReadSecureValue := false
	inf := prof.APIInfo()
	if inf != nil && inf.Version != nil {
		if inf.Version.GreaterThanOrEqualTo(&api.Version{Major: 4}) && inf.Config.RoleBasedPermissions {
			canReadSecureValue = inf.User.Can(tc.PermParameterSecureRead)
		} else {
			canReadSecureValue = inf.User.PrivLevel == auth.PrivLevelAdmin
		}
	}
	for _, profile := range profiles {
		// Attach Parameters if the 'id' parameter is sent
		if _, ok := prof.APIInfo().Params[IDQueryParam]; ok {
			profile.Parameters, err = ReadParameters(prof.ReqInfo.Tx, prof.ReqInfo.User, *profile.ID, canReadSecureValue)
			if err != nil {
				return nil, nil, errors.New("profile read reading parameters: " + err.Error()), http.StatusInternalServerError, nil
			}
		}
		profileInterfaces = append(profileInterfaces, profile)
	}

	return profileInterfaces, nil, nil, http.StatusOK, &maxTime

}

func selectMaxLastUpdatedQuery(where string) string {
	return `SELECT max(t) from (
		SELECT max(prof.last_updated) as t FROM profile prof
LEFT JOIN cdn c ON prof.cdn = c.id ` + where +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='profile') as res`
}

func selectProfilesQuery() string {

	query := `SELECT
prof.description,
prof.id,
prof.last_updated,
prof.name,
prof.routing_disabled,
prof.type,
c.id as cdn,
c.name as cdn_name
FROM profile prof
LEFT JOIN cdn c ON prof.cdn = c.id`

	return query
}

func ReadParameters(tx *sqlx.Tx, user *auth.CurrentUser, profileID int, canReadSecureValue bool) ([]tc.ParameterNullable, error) {
	queryValues := make(map[string]interface{})
	queryValues["profile_id"] = profileID

	query := selectParametersQuery()
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	var params []tc.ParameterNullable
	for rows.Next() {
		var param tc.ParameterNullable

		if err = rows.StructScan(&param); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		var isSecure bool
		if param.Secure != nil {
			isSecure = *param.Secure
		}
		if isSecure && !canReadSecureValue {
			param.Value = &parameter.HiddenField
		}
		params = append(params, param)
	}
	return params, nil
}

func selectParametersQuery() string {
	return `SELECT
p.id,
p.name,
p.config_file,
p.value,
p.secure
FROM parameter p
JOIN profile_parameter pp ON pp.parameter = p.id
WHERE pp.profile = :profile_id`
}

func canProfileBeAlteredByCurrentUser(user string, tx *sql.Tx, cName *string, cdnID *int) (error, error, int) {
	var cdnName string
	if cName != nil {
		cdnName = *cName
	} else {
		if cdnID != nil {
			cdn, ok, err := dbhelpers.GetCDNNameFromID(tx, int64(*cdnID))
			if err != nil {
				return nil, err, http.StatusInternalServerError
			} else if !ok {
				return nil, nil, http.StatusNotFound
			}
			cdnName = string(cdn)
		} else {
			return errors.New("no cdn found for this profile"), nil, http.StatusBadRequest
		}
	}
	userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(tx, cdnName, user)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, statusCode
	}
	return nil, nil, http.StatusOK
}

func (pr *TOProfile) Update(h http.Header) (error, error, int) {
	if pr.CDNName != nil || pr.CDNID != nil {
		userErr, sysErr, statusCode := canProfileBeAlteredByCurrentUser(pr.ReqInfo.User.UserName, pr.ReqInfo.Tx.Tx, pr.CDNName, pr.CDNID)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, statusCode
		}
	}
	return api.GenericUpdate(h, pr)
}

func (pr *TOProfile) Create() (error, error, int) {
	if pr.CDNName != nil || pr.CDNID != nil {
		userErr, sysErr, statusCode := canProfileBeAlteredByCurrentUser(pr.ReqInfo.User.UserName, pr.ReqInfo.Tx.Tx, pr.CDNName, pr.CDNID)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, statusCode
		}
	}
	return api.GenericCreate(pr)
}

func (pr *TOProfile) Delete() (error, error, int) {
	if pr.CDNName == nil && pr.CDNID == nil {
		cdnName, err := dbhelpers.GetCDNNameFromProfileID(pr.APIInfo().Tx.Tx, *pr.ID)
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
		pr.CDNName = util.StrPtr(string(cdnName))
	}
	if pr.CDNName != nil || pr.CDNID != nil {
		userErr, sysErr, statusCode := canProfileBeAlteredByCurrentUser(pr.ReqInfo.User.UserName, pr.ReqInfo.Tx.Tx, pr.CDNName, pr.CDNID)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, statusCode
		}
	}
	return api.GenericDelete(pr)
}

func updateQuery() string {
	query := `UPDATE
profile SET
cdn=:cdn,
description=:description,
name=:name,
routing_disabled=:routing_disabled,
type=:type
WHERE id=:id RETURNING last_updated`
	return query
}

func insertQuery() string {
	query := `INSERT INTO profile (
cdn,
description,
name,
routing_disabled,
type) VALUES (
:cdn,
:description,
:name,
:routing_disabled,
:type) RETURNING id,last_updated`
	return query
}

func deleteQuery() string {
	return `DELETE FROM profile WHERE id = :id`
}

// Read gets a list of Profiles for APIv5.
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
		"cdn":   {Column: "c.id", Checker: api.IsInt},
		"name":  {Column: "prof.name", Checker: nil},
		"id":    {Column: "prof.id", Checker: api.IsInt},
		"param": {Column: "pp.parameter", Checker: api.IsInt},
	}

	query := selectProfilesQuery()
	if paramValue, ok := inf.Params["param"]; ok {
		if len(paramValue) > 0 {
			query += " LEFT JOIN profile_parameter pp ON prof.id = pp.profile"
		}
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

	query += where + orderBy + pagination
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("profile read: error getting profile(s): %w", err))
		return
	}
	defer log.Close(rows, "unable to close DB connection")

	profile := tc.ProfileV5{}
	var profileList []tc.ProfileV5
	for rows.Next() {
		if err = rows.StructScan(&profile); err != nil {
			api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("error getting profile(s): %w", err))
			return
		}
		profileList = append(profileList, profile)
	}
	rows.Close()
	profileInterfaces := []interface{}{}

	for _, p := range profileList {
		// Attach Parameters if the 'id' parameter is sent
		if _, ok := inf.Params["id"]; ok {
			p.Parameters, err = ReadParameters(inf.Tx, inf.User, p.ID, inf.User.Can(tc.PermParameterSecureRead))
			if err != nil {
				api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("profile read: error reading parameters for a profile: %w", err))
				return
			}
		}
		profileInterfaces = append(profileInterfaces, p)
	}

	api.WriteResp(w, r, profileInterfaces)
	return
}

// Create a Profile for APIv5.
func Create(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	tx := inf.Tx.Tx

	profile, readValErr := readAndValidateJsonStruct(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	//check if user can modify.
	if len(strings.TrimSpace(profile.CDNName)) != 0 || profile.CDNID != 0 {
		userErr, sysErr, statusCode := canProfileBeAlteredByCurrentUser(inf.User.UserName, inf.Tx.Tx, &profile.CDNName, &profile.CDNID)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, tx, statusCode, userErr, sysErr)
			return
		}
	}

	// check if profile already exists
	var count int
	err := tx.QueryRow("SELECT count(*) from profile where name=$1", profile.Name).Scan(&count)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error: %w, when checking if profile '%s' exists", err, profile.Name))
		return
	}
	if count == 1 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("profile:'%s' already exists", profile.Name), nil)
		return
	}

	// create profile
	query := `INSERT INTO profile (name, cdn, type, routing_disabled, description) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id, last_updated, name, description, (select name FROM cdn where id = $2), cdn, routing_disabled, type`

	err = tx.QueryRow(query, profile.Name, profile.CDNID, profile.Type, profile.RoutingDisabled, profile.Description).
		Scan(&profile.ID, &profile.LastUpdated, &profile.Name, &profile.Description, &profile.CDNName, &profile.CDNID, &profile.RoutingDisabled, &profile.Type)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("error: %w in creating profile:%s", err, profile.Name), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "profile was created.")
	w.Header().Set(rfc.Location, fmt.Sprintf("/api/%s/profiles?id=%d", inf.Version, profile.ID))
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, profile)
	changeLogMsg := fmt.Sprintf("PROFILE: %s, ID:%d, ACTION: Created profile", profile.Name, profile.ID)
	api.CreateChangeLogRawTx(api.Created, changeLogMsg, inf.User, tx)
	return
}

// Update a profile for APIv5.
func Update(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx
	profile, readValErr := readAndValidateJsonStruct(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	//check if user can modify.
	if len(strings.TrimSpace(profile.CDNName)) != 0 || profile.CDNID != 0 {
		userErr, sysErr, statusCode := canProfileBeAlteredByCurrentUser(inf.User.UserName, inf.Tx.Tx, &profile.CDNName, &profile.CDNID)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, tx, statusCode, userErr, sysErr)
			return
		}
	}

	requestedProfileId := inf.IntParams["id"]
	// check if the entity was already updated
	userErr, sysErr, errCode = api.CheckIfUnModified(r.Header, inf.Tx, requestedProfileId, "profile")
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	//update profile
	query := `UPDATE profile SET
		name = $2, 
		cdn = $3,
		type = $4, 
		routing_disabled = $5, 
		description = $6
	WHERE id = $1
	RETURNING id, last_updated, name, description, (select name FROM cdn where id = $3), cdn, routing_disabled, type`

	err := tx.QueryRow(query, requestedProfileId, profile.Name, profile.CDNID, profile.Type, profile.RoutingDisabled, profile.Description).
		Scan(&profile.ID, &profile.LastUpdated, &profile.Name, &profile.Description, &profile.CDNName, &profile.CDNID, &profile.RoutingDisabled, &profile.Type)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("profile: %s not found", profile.Name), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, sysErr)
		return
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "profile was updated")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, profile)
	changeLogMsg := fmt.Sprintf("PROFILE: %s, ID:%d, ACTION: Updated profile", profile.Name, profile.ID)
	api.CreateChangeLogRawTx(api.Updated, changeLogMsg, inf.User, tx)
	return
}

// Delete a profile in APIv5.
func Delete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	id := inf.IntParams["id"]
	cdnName, err := dbhelpers.GetCDNNameFromProfileID(tx, id)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}

	//check if user can modify
	if len(strings.TrimSpace(string(cdnName))) != 0 {
		userErr, sysErr, statusCode := canProfileBeAlteredByCurrentUser(inf.User.UserName, inf.Tx.Tx, util.StrPtr(string(cdnName)), nil)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, tx, statusCode, userErr, sysErr)
			return
		}
	}

	// check if profile exists
	exists, err := dbhelpers.ProfileExists(tx, strconv.Itoa(id))
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	if !exists {
		if id != 0 {
			api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("no profile exists by id: %d", id), nil)
			return
		} else {
			api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("no profile exists for empty id"), nil)
			return
		}
	}

	res, err := tx.Exec("DELETE FROM profile WHERE id=$1", id)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("determining rows affected for delete profile: %w", err))
		return
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("no rows deleted for profile"), nil)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "profile was deleted.")
	api.WriteAlerts(w, r, http.StatusOK, alerts)
	changeLogMsg := fmt.Sprintf("ID:%d, ACTION: Deleted profile", id)
	api.CreateChangeLogRawTx(api.Deleted, changeLogMsg, inf.User, tx)
	return
}

// readAndValidateJsonStruct reads json body and validates json fields
func readAndValidateJsonStruct(r *http.Request) (tc.ProfileV5, error) {
	var profile tc.ProfileV5
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		userErr := fmt.Errorf("error decoding POST request body into ProfilesV5 struct %w", err)
		return profile, userErr
	}

	rule := validation.NewStringRule(tovalidate.IsAlphanumericUnderscoreDash, "must consist of only alphanumeric, dash, or underscore characters")
	// validate JSON body
	errs := tovalidate.ToErrors(validation.Errors{
		"name":            validation.Validate(profile.Name, validation.Required, rule),
		"cdn":             validation.Validate(profile.CDNID, validation.NilOrNotEmpty),
		"type":            validation.Validate(profile.Type, validation.Required, validation.NotNil),
		"routingDisabled": validation.Validate(profile.RoutingDisabled, validation.NotNil),
		"description":     validation.Validate(profile.Description, validation.Required),
	})
	if len(errs) > 0 {
		userErr := util.JoinErrs(errs)
		return profile, userErr
	}
	return profile, nil
}
