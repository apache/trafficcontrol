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
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

// CacheGroupsResponse is a list of CacheGroups as a response.
type CacheGroupsResponse struct {
	Response []CacheGroup `json:"response"`
	Alerts
}

// CacheGroupsNullableResponse is a response with a list of CacheGroupNullables.
// Traffic Ops API responses instead uses an interface hold a list of
// TOCacheGroups.
type CacheGroupsNullableResponse struct {
	Response []CacheGroupNullable `json:"response"`
	Alerts
}

// CacheGroupDetailResponse is the JSON object returned for a single Cache Group.
type CacheGroupDetailResponse struct {
	Response CacheGroupNullable `json:"response"`
	Alerts
}

// CacheGroup contains information about a given cache group in Traffic Ops.
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

// CacheGroupNullable contains information about a given cache group in Traffic Ops.
// Unlike CacheGroup, CacheGroupNullable's fields are nullable.
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

// CachegroupTrimmedName is useful when the only info about a cache group you
// want to return is its name.
type CachegroupTrimmedName struct {
	Name string `json:"name"`
}

// CachegroupQueueUpdatesRequest holds info relating to the
// cachegroups/{{ID}}/queue_update TO route.
type CachegroupQueueUpdatesRequest struct {
	Action string           `json:"action"`
	CDN    *CDNName         `json:"cdn"`
	CDNID  *util.JSONIntStr `json:"cdnId"`
}

// CacheGroupsResponseV50 is a list of CacheGroups as a response.
type CacheGroupsResponseV50 struct {
	Response []CacheGroupV50 `json:"response"`
	Alerts
}

// CacheGroupsNullableResponseV50 is a response with a list of CacheGroupNullables.
// Traffic Ops API responses instead uses an interface hold a list of
// TOCacheGroups.
type CacheGroupsNullableResponseV50 struct {
	Response []CacheGroupNullableV50 `json:"response"`
	Alerts
}

// CacheGroupDetailResponseV50 is the JSON object returned for a single Cache Group.
type CacheGroupDetailResponseV50 struct {
	Response CacheGroupNullableV50 `json:"response"`
	Alerts
}

// CacheGroupV50 contains information about a given cache group in Traffic Ops.
type CacheGroupV50 struct {
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
	LastUpdated                 time.Time            `json:"lastUpdated" db:"last_updated"`
	Fallbacks                   []string             `json:"fallbacks" db:"fallbacks"`
}

// CacheGroupNullableV50 contains information about a given cache group in Traffic Ops.
// Unlike CacheGroup, CacheGroupNullable's fields are nullable.
type CacheGroupNullableV50 struct {
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
	LastUpdated                 *time.Time            `json:"lastUpdated" db:"last_updated"`
	Fallbacks                   *[]string             `json:"fallbacks" db:"fallbacks"`
}

// CachegroupTrimmedNameV50 is useful when the only info about a cache group you
// want to return is its name.
type CachegroupTrimmedNameV50 struct {
	Name string `json:"name"`
}

// CachegroupQueueUpdatesRequestV50 holds info relating to the
// cachegroups/{{ID}}/queue_update TO route.
type CachegroupQueueUpdatesRequestV50 struct {
	Action string           `json:"action"`
	CDN    *CDNName         `json:"cdn"`
	CDNID  *util.JSONIntStr `json:"cdnId"`
}

// CacheGroupsResponseV5 is the type of response from the cachegroups
// Traffic Ops endpoint.
// It always points to the type for the latest minor version of CacheGroupsResponseV5x APIv5.
type CacheGroupsResponseV5 = CacheGroupsResponseV50

// CacheGroupsNullableResponseV5 is the type of response from the cachegroups
// Traffic Ops endpoint.
// It always points to the type for the latest minor version of CacheGroupsNullableResponseV5x APIv5.
type CacheGroupsNullableResponseV5 = CacheGroupsNullableResponseV50

// CacheGroupDetailResponseV5 is the type of response from the cachegroups
// Traffic Ops endpoint.
// It always points to the type for the latest minor version of CacheGroupDetailResponseV5x APIv5.
type CacheGroupDetailResponseV5 = CacheGroupDetailResponseV50

// CacheGroupV5 always points to the type for the latest minor version of CacheGroupV5x APIv5.
type CacheGroupV5 = CacheGroupV50

// CacheGroupNullableV5 always points to the type for the latest minor version of CacheGroupNullableV5x APIv5.
type CacheGroupNullableV5 = CacheGroupNullableV50

// CachegroupTrimmedNameV5 always points to the type for the latest minor version of CachegroupTrimmedNameV5x APIv5.
type CachegroupTrimmedNameV5 = CachegroupTrimmedNameV50

// CachegroupQueueUpdatesRequestV5 is the type of request to the cachegroups
// Traffic Ops endpoint.
// It always points to the type for the latest minor version of CachegroupQueueUpdatesRequestV5x APIv5.
type CachegroupQueueUpdatesRequestV5 = CachegroupQueueUpdatesRequestV50
