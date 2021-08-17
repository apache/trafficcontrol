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

// AllCacheGroupParametersResponse is a Cache Group Parameter response body.
//
// Deprecated: Cache Group Parameter associations have been removed in APIv4.
type AllCacheGroupParametersResponse struct {
	Response CacheGroupParametersList `json:"response"`
	Alerts
}

// CacheGroupParametersList is the type of the `response` property of Traffic
// Ops's responses to its /cachegroupparameters API endpoint.
//
// Deprecated: Cache Group Parameter associations have been removed in APIv4.
type CacheGroupParametersList struct {
	CacheGroupParameters []CacheGroupParametersResponseNullable `json:"cachegroupParameters"`
}

// CacheGroupParametersNullable is the type of a response from Traffic Ops to a
// POST request made to its /cachegroupparameters API endpoint.
//
// Deprecated: Cache Group Parameter associations have been removed in APIv4.
type CacheGroupParametersNullable struct {
	CacheGroup     *int       `json:"cacheGroupId"  db:"cachegroup"`
	CacheGroupName *string    `json:"cachegroup,omitempty"`
	Parameter      *int       `json:"parameterId"  db:"parameter"`
	LastUpdated    *TimeNoMod `json:"lastUpdated,omitempty"  db:"last_updated"`
}

// CacheGroupParametersResponseNullable is the type of each entry in the
// `cachegroupParameters` property of the response from Traffic Ops to requests
// made to its /cachegroupparameters endpoint.
//
// Deprecated: Cache Group Parameter associations have been removed in APIv4.
type CacheGroupParametersResponseNullable struct {
	CacheGroup  *string    `json:"cachegroup"  db:"cachegroup"`
	Parameter   *int       `json:"parameter"  db:"parameter"`
	LastUpdated *TimeNoMod `json:"last_updated,omitempty"  db:"last_updated"`
}

// FormatForResponse converts a CacheGroupParametersNullable to CacheGroupParametersResponseNullable
// in order to format the output the same as the Perl endpoint.
//
// Deprecated: Cache Group Parameter associations have been removed in APIv4.
func FormatForResponse(param CacheGroupParametersNullable) CacheGroupParametersResponseNullable {
	cgResp := CacheGroupParametersResponseNullable{}
	cgResp.Parameter = param.Parameter
	cgResp.CacheGroup = param.CacheGroupName
	cgResp.LastUpdated = param.LastUpdated
	return cgResp
}
