package tc

import "encoding/json"

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
	ConfigFile string `json:"configFile"`
	Name       string `json:"name"`
	Secure     int    `json:"secure"`
	Value      string `json:"value"`
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
	ProfileID int64   `json:"profileId"`
	ParamIDs  []int64 `json:"paramIds"`
	Replace   bool    `json:"replace"`
}

type PostParamProfile struct {
	ParamID    int64   `json:"paramId"`
	ProfileIDs []int64 `json:"profileIds"`
	Replace    bool    `json:"replace"`
}
