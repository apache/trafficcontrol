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

package tc

// CacheGroupParameterRequest Cache Group Parameter request body
type CacheGroupParameterRequest struct {
	CacheGroupID int `json:"cacheGroupId"`
	ParameterID  int `json:"parameterId"`
}

// CacheGroupParametersPostResponse Response body when Posting to associate a
// Parameter with a Cache Group.
type CacheGroupParametersPostResponse struct {
	Response []CacheGroupParameterRequest `json:"response"`
	Alerts
}

// CacheGroupParametersResponse is a Cache Group Parameter response body.
type CacheGroupParametersResponse struct {
	Response []CacheGroupParameter `json:"response"`
	Alerts
}

// CacheGroupParameter ...
type CacheGroupParameter struct {
	ConfigFile  string    `json:"configFile"`
	ID          int       `json:"id"`
	LastUpdated TimeNoMod `json:"lastUpdated"`
	Name        string    `json:"name"`
	Secure      bool      `json:"secure"`
	Value       string    `json:"value"`
}

// CacheGroupParameterNullable ...
type CacheGroupParameterNullable struct {
	ConfigFile  *string    `json:"configFile" db:"config_file"`
	ID          *int       `json:"id" db:"id"`
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Name        *string    `json:"name" db:"name"`
	Secure      *bool      `json:"secure" db:"secure"`
	Value       *string    `json:"value" db:"value"`
}
