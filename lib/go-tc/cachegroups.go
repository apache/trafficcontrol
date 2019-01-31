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
	"github.com/apache/trafficcontrol/lib/go-util"
)

// CacheGroupResponse ...
type CacheGroupsResponse struct {
	Response []CacheGroup `json:"response"`
}

type CacheGroupsNullableResponse struct {
	Response []CacheGroupNullable `json:"response"`
}

// CacheGroupDetailResponse is the JSON object returned for a single CacheGroup
type CacheGroupDetailResponse struct {
	Response CacheGroupNullable `json:"response"`
	Alerts
}

// CacheGroup contains information about a given Cachegroup in Traffic Ops.
type CacheGroup struct {
	ID                          int                  `json:"id" db:"id"`
	Name                        string               `json:"name" db:"name"`
	ShortName                   string               `json:"shortName" db:"short_name"`
	Latitude                    float64              `json:"latitude" db:"latitude"`
	Longitude                   float64              `json:"longitude" db:"longitude"`
	ParentName                  string               `json:"parentCachegroupName" db:"parent_cachegroup_name"`
	ParentCachegroupID          int                  `json:"parentCachegroupId" db:"parent_cachegroup_id"`
	SecondaryParentName         string               `json:"secondaryParentCachegroupName" db:"secondary_parent_cachegroup_name"`
	SecondaryParentCachegroupID int                  `json:"secondaryParentCachegroupId" db:"secondary_parent_cachegroup_id"`
	FallbackToClosest           bool                 `json:"fallbackToClosest" db:"fallback_to_closest"`
	LocalizationMethods         []LocalizationMethod `json:"localizationMethods" db:"localization_methods"`
	Type                        string               `json:"typeName" db:"type_name"` // aliased to type_name to disambiguate struct scans due to join on 'type' table
	TypeID                      int                  `json:"typeId" db:"type_id"`     // aliased to type_id to disambiguate struct scans due join on 'type' table
	LastUpdated                 TimeNoMod            `json:"lastUpdated" db:"last_updated"`
	Fallbacks                   []string             `json:"fallbacks" db:"fallbacks"`
}

type CacheGroupNullable struct {
	ID                          *int                  `json:"id" db:"id"`
	Name                        *string               `json:"name" db:"name"`
	ShortName                   *string               `json:"shortName" db:"short_name"`
	Latitude                    *float64              `json:"latitude" db:"latitude"`
	Longitude                   *float64              `json:"longitude" db:"longitude"`
	ParentName                  *string               `json:"parentCachegroupName" db:"parent_cachegroup_name"`
	ParentCachegroupID          *int                  `json:"parentCachegroupId" db:"parent_cachegroup_id"`
	SecondaryParentName         *string               `json:"secondaryParentCachegroupName" db:"secondary_parent_cachegroup_name"`
	SecondaryParentCachegroupID *int                  `json:"secondaryParentCachegroupId" db:"secondary_parent_cachegroup_id"`
	FallbackToClosest           *bool                 `json:"fallbackToClosest" db:"fallback_to_closest"`
	LocalizationMethods         *[]LocalizationMethod `json:"localizationMethods" db:"localization_methods"`
	Type                        *string               `json:"typeName" db:"type_name"` // aliased to type_name to disambiguate struct scans due to join on 'type' table
	TypeID                      *int                  `json:"typeId" db:"type_id"`     // aliased to type_id to disambiguate struct scans due join on 'type' table
	LastUpdated                 *TimeNoMod            `json:"lastUpdated" db:"last_updated"`
	Fallbacks                   *[]string             `json:"fallbacks" db:"fallbacks"`
}

type CachegroupTrimmedName struct {
	Name string `json:"name"`
}

type CachegroupQueueUpdatesRequest struct {
	Action string           `json:"action"`
	CDN    *CDNName         `json:"cdn"`
	CDNID  *util.JSONIntStr `json:"cdnId"`
}
