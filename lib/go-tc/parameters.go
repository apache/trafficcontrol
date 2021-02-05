package tc

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-util"

	"github.com/lib/pq"
)

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

// ParametersResponse ...
type ParametersResponse struct {
	Response []Parameter `json:"response"`
	Alerts
}

// Parameter ...
type Parameter struct {
	ConfigFile  string          `json:"configFile" db:"config_file"`
	ID          int             `json:"id" db:"id"`
	LastUpdated TimeNoMod       `json:"lastUpdated" db:"last_updated"`
	Name        string          `json:"name" db:"name"`
	Profiles    json.RawMessage `json:"profiles" db:"profiles"`
	Secure      bool            `json:"secure" db:"secure"`
	Value       string          `json:"value" db:"value"`
}

// ParameterNullable - a struct version that allows for all fields to be null, mostly used by the API side
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

type ProfileParameterByName struct {
	ConfigFile  string    `json:"configFile"`
	ID          int       `json:"id"`
	LastUpdated TimeNoMod `json:"lastUpdated"`
	Name        string    `json:"name"`
	Secure      bool      `json:"secure"`
	Value       string    `json:"value"`
}

type ProfileParameterByNamePost struct {
	ConfigFile *string `json:"configFile"`
	Name       *string `json:"name"`
	Secure     *int    `json:"secure"`
	Value      *string `json:"value"`
}

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

type ProfileParameterPostRespObj struct {
	ProfileParameterByNamePost
	ID int64 `json:"id"`
}

type ProfileParameterPostResp struct {
	Parameters  []ProfileParameterPostRespObj `json:"parameters"`
	ProfileID   int                           `json:"profileId"`
	ProfileName string                        `json:"profileName"`
}

type PostProfileParam struct {
	ProfileID *int64   `json:"profileId"`
	ParamIDs  *[]int64 `json:"paramIds"`
	Replace   *bool    `json:"replace"`
}

func (pp *PostProfileParam) Sanitize(tx *sql.Tx) {
	if pp.Replace == nil {
		pp.Replace = util.BoolPtr(false)
	}
}

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
	} else if ok, err := ParamsExist(*pp.ParamIDs, tx); err != nil {
		errs = append(errs, errors.New(fmt.Sprintf("checking parameter IDs %v existence: "+err.Error(), *pp.ParamIDs)))
	} else if !ok {
		errs = append(errs, errors.New(fmt.Sprintf("parameters with IDs %v don't all exist", *pp.ParamIDs)))
	}
	if len(errs) > 0 {
		return util.JoinErrs(errs)
	}
	return nil
}

type PostParamProfile struct {
	ParamID    *int64   `json:"paramId"`
	ProfileIDs *[]int64 `json:"profileIds"`
	Replace    *bool    `json:"replace"`
}

func (pp *PostParamProfile) Sanitize(tx *sql.Tx) {
	if pp.Replace == nil {
		pp.Replace = util.BoolPtr(false)
	}
}

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
func ParamsExist(ids []int64, tx *sql.Tx) (bool, error) {
	count := 0
	if err := tx.QueryRow(`SELECT count(*) from parameter where id = ANY($1)`, pq.Array(ids)).Scan(&count); err != nil {
		return false, errors.New("querying parameters existence from id: " + err.Error())
	}
	return count == len(ids), nil
}

// ProfileParametersNullable is an object of the form returned by the Traffic Ops /profileparameters endpoint.
type ProfileParametersNullable struct {
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Profile     *string    `json:"profile" db:"profile"`
	Parameter   *int       `json:"parameter" db:"parameter_id"`
}

type ProfileParametersNullableResponse struct {
	Response []ProfileParametersNullable `json:"response"`
}

// ProfileExportImportParameterNullable is an object of the form used by Traffic Ops
// to represent parameters for exported and imported profiles.
type ProfileExportImportParameterNullable struct {
	ConfigFile *string `json:"config_file"`
	Name       *string `json:"name"`
	Value      *string `json:"value"`
}
