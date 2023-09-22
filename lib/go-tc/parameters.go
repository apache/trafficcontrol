package tc

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
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-util"

	"github.com/lib/pq"
)

// ParametersResponse is the type of the response from Traffic Ops to GET
// requests made to the /parameters and /profiles/name/{{Name}}/parameters
// endpoints of its API.
type ParametersResponse struct {
	Response []Parameter `json:"response"`
	Alerts
}

// A Parameter defines some configuration setting (which is usually but
// definitely not always a line in a configuration file) used by some Profile
// or Cache Group.
type Parameter struct {
	ConfigFile  string          `json:"configFile" db:"config_file"`
	ID          int             `json:"id" db:"id"`
	LastUpdated TimeNoMod       `json:"lastUpdated" db:"last_updated"`
	Name        string          `json:"name" db:"name"`
	Profiles    json.RawMessage `json:"profiles" db:"profiles"`
	Secure      bool            `json:"secure" db:"secure"`
	Value       string          `json:"value" db:"value"`
}

// ParameterNullable is exactly like Parameter except that its properties are
// reference values, so they can be nil.
type ParameterNullable struct {
	//
	// NOTE: the db: struct tags are used for testing to map to their equivalent database column (if there is one)
	//
	ConfigFile  *string         `json:"configFile" db:"config_file"`
	ID          *int            `json:"id" db:"id"`
	LastUpdated *TimeNoMod      `json:"lastUpdated" db:"last_updated"`
	Name        *string         `json:"name" db:"name"`
	Profiles    json.RawMessage `json:"profiles" db:"profiles"`
	Secure      *bool           `json:"secure" db:"secure"`
	Value       *string         `json:"value" db:"value"`
}

// ParametersResponseV5 is an alias for the latest minor version for the major version 5.
type ParametersResponseV5 = ParametersResponseV50

// ParametersResponseV50 is the type of the response from Traffic Ops to GET
// requests made to the /parameters and /profiles/name/{{Name}}/parameters
// endpoints of its API.
type ParametersResponseV50 struct {
	Response []ParameterV5 `json:"response"`
	Alerts
}

// ParameterV5 is an alias for the latest minor version for the major version 5.
type ParameterV5 = ParameterV50

// A ParameterV50 defines some configuration setting (which is usually but
// definitely not always a line in a configuration file) used by some Profile
// or Cache Group.
type ParameterV50 struct {
	ConfigFile  string          `json:"configFile" db:"config_file"`
	ID          int             `json:"id" db:"id"`
	LastUpdated time.Time       `json:"lastUpdated" db:"last_updated"`
	Name        string          `json:"name" db:"name"`
	Profiles    json.RawMessage `json:"profiles" db:"profiles"`
	Secure      bool            `json:"secure" db:"secure"`
	Value       string          `json:"value" db:"value"`
}

// ParameterNullableV5 is an alias for the latest minor version for the major version 5.
type ParameterNullableV5 = ParameterNullableV50

// ParameterNullableV50 is exactly like Parameter except that its properties are
// reference values, so they can be nil.
type ParameterNullableV50 struct {
	//
	// NOTE: the db: struct tags are used for testing to map to their equivalent database column (if there is one)
	//
	ConfigFile  *string         `json:"configFile" db:"config_file"`
	ID          *int            `json:"id" db:"id"`
	LastUpdated *time.Time      `json:"lastUpdated" db:"last_updated"`
	Name        *string         `json:"name" db:"name"`
	Profiles    json.RawMessage `json:"profiles" db:"profiles"`
	Secure      *bool           `json:"secure" db:"secure"`
	Value       *string         `json:"value" db:"value"`
}

// ProfileParameterByName is a structure that's used to represent a Parameter
// in a context where they are associated with some Profile specified by a
// client of the Traffic Ops API by its Name.
type ProfileParameterByName struct {
	ConfigFile  string    `json:"configFile"`
	ID          int       `json:"id"`
	LastUpdated TimeNoMod `json:"lastUpdated"`
	Name        string    `json:"name"`
	Secure      bool      `json:"secure"`
	Value       string    `json:"value"`
}

// ProfileParameterByNameV5 is the alias to the latest minor version of major version 5
type ProfileParameterByNameV5 ProfileParameterByNameV50

// ProfileParameterByNameV50 is a structure that's used to represent a Parameter
// in a context where they are associated with some Profile specified by a
// client of the Traffic Ops API by its Name.
type ProfileParameterByNameV50 struct {
	ConfigFile  string    `json:"configFile"`
	ID          int       `json:"id"`
	LastUpdated time.Time `json:"lastUpdated"`
	Name        string    `json:"name"`
	Secure      bool      `json:"secure"`
	Value       string    `json:"value"`
}

// ProfileParameterByNamePost is a structure that's only used internally to
// represent a Parameter that has been requested by a client of the Traffic Ops
// API to be associated with some Profile which was specified by Name.
//
// TODO: This probably shouldn't exist, or at least not be in the lib.
type ProfileParameterByNamePost struct {
	ConfigFile *string `json:"configFile"`
	Name       *string `json:"name"`
	Secure     *int    `json:"secure"`
	Value      *string `json:"value"`
}

// Validate implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
func (p *ProfileParameterByNamePost) Validate(tx *sql.Tx) []error {
	return validateProfileParamPostFields(p.ConfigFile, p.Name, p.Value, p.Secure)
}

func validateProfileParamPostFields(configFile, name, value *string, secure *int) []error {
	errs := []string{}
	if configFile == nil || *configFile == "" {
		errs = append(errs, "configFile")
	}
	if name == nil || *name == "" {
		errs = append(errs, "name")
	}
	if secure == nil {
		errs = append(errs, "secure")
	}
	if value == nil {
		errs = append(errs, "value")
	}
	if len(errs) > 0 {
		return []error{errors.New("required fields missing: " + strings.Join(errs, ", "))}
	}
	return nil
}

// ProfileParametersByNamePost is the object posted to profile/name/parameter endpoints. This object may be posted as either a single JSON object, or an array of objects. Either will unmarshal into this object; a single object will unmarshal into an array of 1 element.
type ProfileParametersByNamePost []ProfileParameterByNamePost

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (pp *ProfileParametersByNamePost) UnmarshalJSON(bts []byte) error {
	bts = bytes.TrimLeft(bts, " \n\t\r")
	if len(bts) == 0 {
		return errors.New("no body")
	}
	if bts[0] == '[' {
		ppArr := []ProfileParameterByNamePost{}
		err := json.Unmarshal(bts, &ppArr)
		*pp = ppArr
		return err
	}
	p := ProfileParameterByNamePost{}
	if err := json.Unmarshal(bts, &p); err != nil {
		return err
	}
	*pp = append(*pp, p)
	return nil
}

// Validate implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
func (pp *ProfileParametersByNamePost) Validate(tx *sql.Tx) error {
	errs := []error{}
	ppArr := ([]ProfileParameterByNamePost)(*pp)
	for i, profileParam := range ppArr {
		if ppErrs := profileParam.Validate(tx); len(ppErrs) > 0 {
			for _, err := range ppErrs {
				errs = append(errs, errors.New("object "+strconv.Itoa(i)+": "+err.Error()))
			}
		}
	}
	if len(errs) > 0 {
		return util.JoinErrs(errs)
	}
	return nil
}

// ProfileParameterPostRespObj is a single Parameter in the Parameters slice of
// a ProfileParameterPostResp.
type ProfileParameterPostRespObj struct {
	ProfileParameterByNamePost
	ID int64 `json:"id"`
}

// ProfileParameterPostResp is the type of the `response` property of responses
// from Traffic Ops to POST requests made to its
// /profiles/name/{{Name}}/parameters API endpoint.
type ProfileParameterPostResp struct {
	Parameters  []ProfileParameterPostRespObj `json:"parameters"`
	ProfileID   int                           `json:"profileId"`
	ProfileName string                        `json:"profileName"`
}

// A PostProfileParam is a request to associate zero or more Parameters with a
// particular Profile.
type PostProfileParam struct {
	ProfileID *int64   `json:"profileId"`
	ParamIDs  *[]int64 `json:"paramIds"`
	Replace   *bool    `json:"replace"`
}

// Sanitize ensures that Replace is not nil, setting it to false if it is.
//
// TODO: Figure out why this is taking a db transaction - should this be moved?
func (pp *PostProfileParam) Sanitize(tx *sql.Tx) {
	if pp.Replace == nil {
		pp.Replace = util.BoolPtr(false)
	}
}

// Validate implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
func (pp *PostProfileParam) Validate(tx *sql.Tx) error {
	pp.Sanitize(tx)
	errs := []error{}
	if pp.ProfileID == nil {
		errs = append(errs, errors.New("profileId missing"))
	} else if ok, err := ProfileExistsByID(*pp.ProfileID, tx); err != nil {
		errs = append(errs, errors.New("checking profile ID "+strconv.Itoa(int(*pp.ProfileID))+" existence: "+err.Error()))
	} else if !ok {
		errs = append(errs, errors.New("no profile with ID "+strconv.Itoa(int(*pp.ProfileID))+" exists"))
	}
	if pp.ParamIDs == nil {
		errs = append(errs, errors.New("paramIds missing"))
	} else if nonExistingIDs, err := ParamsExist(*pp.ParamIDs, tx); err != nil {
		errs = append(errs, fmt.Errorf("checking parameter IDs %v existence: %w", *pp.ParamIDs, err))
	} else if len(nonExistingIDs) >= 1 {
		errs = append(errs, fmt.Errorf("parameters with IDs %v don't exist", nonExistingIDs))
	}
	if len(errs) > 0 {
		return util.JoinErrs(errs)
	}
	return nil
}

// A PostParamProfile is a request to associate a particular Parameter with
// zero or more Profiles.
type PostParamProfile struct {
	ParamID    *int64   `json:"paramId"`
	ProfileIDs *[]int64 `json:"profileIds"`
	Replace    *bool    `json:"replace"`
}

// Sanitize ensures that Replace is not nil, setting it to false if it is.
//
// TODO: Figure out why this is taking a db transaction - should this be moved?
func (pp *PostParamProfile) Sanitize(tx *sql.Tx) {
	if pp.Replace == nil {
		pp.Replace = util.BoolPtr(false)
	}
}

// Validate implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
func (pp *PostParamProfile) Validate(tx *sql.Tx) error {
	pp.Sanitize(tx)

	errs := []error{}
	if pp.ParamID == nil {
		errs = append(errs, errors.New("paramId missing"))
	} else if ok, err := ParamExists(*pp.ParamID, tx); err != nil {
		errs = append(errs, errors.New("checking param ID "+strconv.Itoa(int(*pp.ParamID))+" existence: "+err.Error()))
	} else if !ok {
		errs = append(errs, errors.New("no parameter with ID "+strconv.Itoa(int(*pp.ParamID))+" exists"))
	}
	if pp.ProfileIDs == nil {
		errs = append(errs, errors.New("profileIds missing"))
	} else if ok, err := ProfilesExistByIDs(*pp.ProfileIDs, tx); err != nil {
		errs = append(errs, errors.New(fmt.Sprintf("checking profiles IDs %v existence: "+err.Error(), *pp.ProfileIDs)))
	} else if !ok {
		errs = append(errs, errors.New(fmt.Sprintf("profiles with IDs %v don't all exist", *pp.ProfileIDs)))
	}
	if len(errs) > 0 {
		return util.JoinErrs(errs)
	}
	return nil
}

// ParamExists returns whether a parameter with the given id exists, and any error.
// TODO move to helper package.
func ParamExists(id int64, tx *sql.Tx) (bool, error) {
	count := 0
	if err := tx.QueryRow(`SELECT count(*) from parameter where id = $1`, id).Scan(&count); err != nil {
		return false, errors.New("querying param existence from id: " + err.Error())
	}
	return count > 0, nil
}

// ParamsExist returns whether parameters exist for all the given ids, and any error.
// TODO move to helper package.
func ParamsExist(ids []int64, tx *sql.Tx) ([]int64, error) {
	var nonExistingIDs []int64
	if err := tx.QueryRow(`SELECT ARRAY_AGG(id) FROM UNNEST($1::INT[]) AS id WHERE id NOT IN (SELECT id FROM parameter)`, pq.Array(ids)).Scan(pq.Array(&nonExistingIDs)); err != nil {
		return nil, fmt.Errorf("querying parameters existence from id: %w", err)
	}
	if len(nonExistingIDs) >= 1 {
		return nonExistingIDs, nil
	}
	return nil, nil
}

// ProfileParametersNullable is an object of the form returned by the Traffic Ops /profileparameters endpoint.
type ProfileParametersNullable struct {
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Profile     *string    `json:"profile" db:"profile"`
	Parameter   *int       `json:"parameter" db:"parameter_id"`
}

// ProfileParametersNullableV5 is the latest minor version of the major version 5
type ProfileParametersNullableV5 ProfileParametersNullableV50

// ProfileParametersNullableV50 is an object of the form returned by the Traffic Ops /profileparameters endpoint.
type ProfileParametersNullableV50 struct {
	LastUpdated *time.Time `json:"lastUpdated" db:"last_updated"`
	Profile     *string    `json:"profile" db:"profile"`
	Parameter   *int       `json:"parameter" db:"parameter_id"`
}

// ProfileParametersNullableResponse is the structure of a response from
// Traffic Ops to GET requests made to its /profileparameters API endpoint.
//
// TODO: This is only used internally in a v3 client method (not its call
// signature) - deprecate? Remove?
type ProfileParametersNullableResponse struct {
	Response []ProfileParametersNullable `json:"response"`
}

// ProfileParam is a relationship between a Profile and some Parameter
// assigned to it as it appears in the Traffic Ops API's responses to the
// /profileparameters endpoint.
type ProfileParam struct {
	// Parameter is the ID of the Parameter.
	Parameter int `json:"parameter"`
	// Profile is the name of the Profile to which the Parameter is assigned.
	Profile     string     `json:"profile"`
	LastUpdated *TimeNoMod `json:"lastUpdated"`
}

// ProfileParamV5 is the latest minor version of the major version 5
type ProfileParamV5 ProfileParamV50

// ProfileParamV50 is a relationship between a Profile and some Parameter
// assigned to it as it appears in the Traffic Ops API's responses to the
// /profileparameters endpoint.
type ProfileParamV50 struct {
	Parameter   int        `json:"parameter"`
	Profile     string     `json:"profile"`
	LastUpdated *time.Time `json:"lastUpdated"`
}

// ProfileParameterCreationRequest is the type of data accepted by Traffic
// Ops as payloads in POST requests to its /profileparameters endpoint.
type ProfileParameterCreationRequest struct {
	ParameterID int `json:"parameterId"`
	ProfileID   int `json:"profileId"`
}

// ProfileParametersAPIResponse is the type of a response from Traffic Ops to
// requests made to its /profileparameters endpoint.
type ProfileParametersAPIResponse struct {
	Response []ProfileParam `json:"response"`
	Alerts
}

// ProfileParametersAPIResponseV5 is the latest minor version of the major version 5
type ProfileParametersAPIResponseV5 ProfileParametersAPIResponseV50

// ProfileParametersAPIResponseV50 is the type of a response from Traffic Ops to
// requests made to its /profileparameters endpoint.
type ProfileParametersAPIResponseV50 struct {
	Response []ProfileParamV5 `json:"response"`
	Alerts
}

// ProfileExportImportParameterNullable is an object of the form used by Traffic Ops
// to represent parameters for exported and imported profiles.
type ProfileExportImportParameterNullable struct {
	ConfigFile *string `json:"config_file"`
	Name       *string `json:"name"`
	Value      *string `json:"value"`
}
