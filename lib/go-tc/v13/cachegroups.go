package v13

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

import tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"

// CacheGroupResponse ...
type CacheGroupsResponse struct {
	Response []CacheGroup `json:"response"`
}

// CacheGroup contains information about a given Cachegroup in Traffic Ops.
type CacheGroup struct {
	ID                          int          `json:"id" db:"id"`
	Name                        string       `json:"name" db:"name"`
	ShortName                   string       `json:"shortName" db:"short_name"`
	Latitude                    float64      `json:"latitude" db:"latitude"`
	Longitude                   float64      `json:"longitude" db:"longitude"`
	ParentName                  string       `json:"parentCachegroupName"`
	ParentCachegroupID          int          `json:"parentCachegroupId" db:"parent_cachegroup_id"`
	SecondaryParentName         string       `json:"secondaryParentCachegroupName"`
	SecondaryParentCachegroupID int          `json:"secondaryParentCachegroupId" db:"secondary_parent_cachegroup_id"`
	Type                        string       `json:"typeName" db:"type_name"` // aliased to type_name to disambiguate struct scans due to join on 'type' table
	TypeID                      int          `json:"typeId" db:"type_id"`     // aliased to type_id to disambiguate struct scans due join on 'type' table
	LastUpdated                 tc.TimeNoMod `json:"lastUpdated" db:"last_updated"`
}

type CacheGroupNullable struct {
	ID                          *int          `json:"id" db:"id"`
	Name                        *string       `json:"name" db:"name"`
	ShortName                   *string       `json:"shortName" db:"short_name"`
	Latitude                    *float64      `json:"latitude" db:"latitude"`
	Longitude                   *float64      `json:"longitude"db:"longitude"`
	ParentName                  *string       `json:"parentCachegroupName"`
	ParentCachegroupID          *int          `json:"parentCachegroupId" db:"parent_cachegroup_id"`
	SecondaryParentName         *string       `json:"secondaryParentCachegroupName"`
	SecondaryParentCachegroupID *int          `json:"secondaryParentCachegroupId" db:"secondary_parent_cachegroup_id"`
	Type                        *string       `json:"typeName" db:"type_name"` // aliased to type_name to disambiguate struct scans due to join on 'type' table
	TypeID                      *int          `json:"typeId" db:"type_id"`     // aliased to type_id to disambiguate struct scans due join on 'type' table
	LastUpdated                 *tc.TimeNoMod `json:"lastUpdated" db:"last_updated"`
}
