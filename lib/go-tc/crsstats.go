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

// CRSStats is the returned data from TRs stats endpoint.
type CRSStats struct {
	App   CRSStatsApp   `json:"app"`
	Stats CRSStatsStats `json:"stats"`
}

// CRSStatsApp represents metadata about a given TR.
type CRSStatsApp struct {
	BuildTimestamp string `json:"buildTimestamp"`
	Name           string `json:"name"`
	DeployDir      string `json:"deploy-dir"`
	GitRevision    string `json:"git-revision"`
	Version        string `json:"version"`
}

// CRSStatsStats represents stats about a given TR.
type CRSStatsStats struct {
	DNSMap           map[string]CRSStatsStat
	HTTPMap          map[string]CRSStatsStat
	TotalDNSCount    uint64                `json:"totalDnsCount"`
	TotalHTTPCount   uint64                `json:"totalHttpCount"`
	TotalDSMissCount uint64                `json:"totalDsMissCount"`
	AppStartTime     uint64                `json:"appStartTime"`
	AverageDnsTime   uint64                `json:"averageDnsTime"`
	AverageHttpTime  uint64                `json:"averageHttpTime"`
	UpdateTracker    CRSStatsUpdateTracker `json:"updateTracker"`
}

// CRSStatsStat represents an individual stat.
type CRSStatsStat struct {
	CZCount                uint64 `json:"czCount"`
	GeoCount               uint64 `json:"geoCount"`
	DeepCZCount            uint64 `json:"deepCzCount"`
	MissCount              uint64 `json:"missCount"`
	DSRCount               uint64 `json:"dsrCount"`
	ErrCount               uint64 `json:"errCount"`
	StaticRouteCount       uint64 `json:"staticRouteCount"`
	FedCount               uint64 `json:"fedCount"`
	RegionalDeniedCount    uint64 `json:"regionalDeniedCount"`
	RegionalAlternateCount uint64 `json:"regionalAlternateCount"`
}

// CRSStatsUpdateTracker is a collection of miscellaneous statistics about a
// Traffic Router instance.
//
// This is named with "CRS" because of legacy naming conventions for Traffic
// Router.
type CRSStatsUpdateTracker struct {
	LastHttpsCertificatesCheck           uint64 `json:"lastHttpsCertificatesCheck"`
	LastGeolocationDatabaseUpdaterUpdate uint64 `json:"lastGeolocationDatabaseUpdaterUpdate"`
	LastCacheStateCheck                  uint64 `json:"lastCacheStateCheck"`
	LastCacheStateChange                 uint64 `json:"lastCacheStateChange"`
	LastNetworkUpdaterUpdate             uint64 `json:"lastNetworkUpdaterUpdate"`
	LastHTTPSCertificatesUpdate          uint64 `json:"lastHttpsCertificatesUpdate"`
	LastConfigCheck                      uint64 `json:"lastConfigCheck"`
	LastConfigChange                     uint64 `json:"lastConfigChange"`
	LastHTTPSCertificatesFetchFail       uint64 `json:"lastHttpsCertificatesFetchFail"`
	LastNetworkUpdaterCheck              uint64 `json:"lastNetworkUpdaterCheck"`
	NewDNSSECKeysFound                   uint64 `json:"newDnsSecKeysFound"`
	LastGeolocationDatabaseUpdaterCheck  uint64 `json:"lastGeolocationDatabaseUpdaterCheck"`
	LastHTTPSCertificatesFetchSuccess    uint64 `json:"lastHttpsCertificatesFetchSuccess"`
	LastSteeringWatcherCheck             uint64 `json:"lastSteeringWatcherCheck"`
	LastDNSSECKeysCheck                  uint64 `json:"lastDnsSecKeysCheck"`
	LastFederationsWatcherCheck          uint64 `json:"lastFederationsWatcherCheck"`
	LastHTTPSCertificatesFetchAttempt    uint64 `json:"lastHttpsCertificatesFetchAttempt"`
}

// Routing represents the aggregated routing percentages across CDNs or for a DS.
type Routing struct {
	StaticRoute       float64 `json:"staticRoute"`
	Geo               float64 `json:"geo"`
	Err               float64 `json:"err"`
	Fed               float64 `json:"fed"`
	CZ                float64 `json:"cz"`
	DeepCZ            float64 `json:"deepCz"`
	RegionalAlternate float64 `json:"regionalAlternate"`
	DSR               float64 `json:"dsr"`
	Miss              float64 `json:"miss"`
	RegionalDenied    float64 `json:"regionalDenied"`
}
