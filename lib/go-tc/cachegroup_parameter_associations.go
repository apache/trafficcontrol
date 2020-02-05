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

// CacheGroupParametersResponse is a Cache Group Parameter response body.
type AllCacheGroupParametersResponse struct {
	Response CacheGroupParametersList `json:"response"`
	Alerts
}

// CacheGroupParametersList ...
type CacheGroupParametersList struct {
	CacheGroupParameters []CacheGroupParametersResponseNullable `json:"cachegroupParameters"`
}

// CacheGroupParameter ...
type CacheGroupParameters struct {
	CacheGroup     int       `json:"cacheGroupId"  db:"cachegroup"`
	CacheGroupName string    `json:"cachegroup,omitempty"`
	Parameter      int       `json:"parameterId"  db:"parameter"`
	LastUpdated    TimeNoMod `json:"lastUpdated,omitempty"  db:"last_updated"`
}

// CacheGroupParameterNullable ...
type CacheGroupParametersNullable struct {
	CacheGroup     *int       `json:"cacheGroupId"  db:"cachegroup"`
	CacheGroupName *string    `json:"cachegroup,omitempty"`
	Parameter      *int       `json:"parameterId"  db:"parameter"`
	LastUpdated    *TimeNoMod `json:"lastUpdated,omitempty"  db:"last_updated"`
}

// CacheGroupParameterResponseNullable ...
type CacheGroupParametersResponseNullable struct {
	CacheGroup  *string    `json:"cachegroup"  db:"cachegroup"`
	Parameter   *int       `json:"parameter"  db:"parameter"`
	LastUpdated *TimeNoMod `json:"last_updated,omitempty"  db:"last_updated"`
}

// FormatForResponse converts a CacheGroupParametersNullable to CacheGroupParametersResponseNullable
// in order to format the output the same as the Perl endpoint.
func FormatForResponse(param CacheGroupParametersNullable) CacheGroupParametersResponseNullable {
	cgResp := CacheGroupParametersResponseNullable{}
	cgResp.Parameter = param.Parameter
	cgResp.CacheGroup = param.CacheGroupName
	cgResp.LastUpdated = param.LastUpdated
	return cgResp
}
