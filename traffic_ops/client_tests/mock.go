package main

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

import "crypto/tls"
import "encoding/json"
import "log"
import "net/http"
import "os"
import "strconv"
import "time"

import "github.com/apache/trafficcontrol/lib/go-tc"
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tocookie"

const API_MIN_MINOR_VERSION = 1
const API_MAX_MINOR_VERSION = 5

// TODO: read this in properly from the VERSION file
const VERSION = "3.0.0"

// This needs to exist for... reasons
const TM_LOGO = `<?xml version="1.0" encoding="UTF-8"?>
<svg width="1391px" height="888px" viewBox="0 0 1391 888" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
    <title>Apache Traffic Control Logo</title>
    <desc>Logo with no text</desc>
    <defs></defs>
    <g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
        <path d="M933.483367,152.57612 L876.912723,209.146764 C820.077911,152.80267 741.853255,118 655.5,118 C481.806446,118 341,258.806446 341,432.5 C341,606.193554 481.806446,747 655.5,747 C742.585973,747 821.404817,711.604215 878.354845,654.414332 L933.923488,709.982975 C862.496219,781.649964 763.677363,826 654.5,826 C436.623666,826 260,649.376334 260,431.5 C260,213.623666 436.623666,37 654.5,37 C763.453296,37 862.09056,81.1681819 933.483367,152.57612 Z" id="Combined-Shape" fill="#F49224"></path>
        <path d="M416.601816,396 L504.202305,396 C501.454569,407.552665 500,419.606466 500,432 C500,517.604136 569.395864,587 655,587 C740.604136,587 810,517.604136 810,432 C810,346.395864 740.604136,277 655,277 C642.963616,277 631.247667,278.371943 620,280.967982 L620,193.456017 C631.266638,191.837544 642.78551,191 654.5,191 C787.324482,191 895,298.675518 895,431.5 C895,564.324482 787.324482,672 654.5,672 C521.675518,672 414,564.324482 414,431.5 C414,419.438696 414.88787,407.584765 416.601816,396 Z" id="Combined-Shape" fill="#F49224"></path>
        <rect id="Rectangle" fill="#F49224" x="620" y="613" width="80" height="161"></rect>
        <rect id="Rectangle" fill="#F49224" x="620" y="77" width="80" height="161"></rect>
        <rect id="Rectangle" fill="#F49224" transform="translate(394.000000, 436.000000) rotate(90.000000) translate(-394.000000, -436.000000) " x="354" y="360" width="80" height="152"></rect>
    </g>
</svg>
`

// Will be used for various `lastUpdated` fields
var CURRENT_TIME *tc.TimeNoMod = tc.NewTimeNoMod()

// Static CDN Fields
var ALL_CDN = "ALL"
var ALL_CDN_ID = 1
var ALL_CDN_DOMAINNAME = "-"
var ALL_CDN_DNSSEC_ENABLED = false
var CDN = "Mock-CDN"
var CDN_ID = 2
var CDN_DOMAIN_NAME = "mock.cdn.test"
var CDN_DNSSEC_ENABLED = false

var CDNS = []tc.CDNNullable{
	tc.CDNNullable{
		DNSSECEnabled: &ALL_CDN_DNSSEC_ENABLED,
		DomainName:    &ALL_CDN_DOMAINNAME,
		ID:            &ALL_CDN_ID,
		LastUpdated:   CURRENT_TIME,
		Name:          &ALL_CDN,
	},
	tc.CDNNullable{
		DNSSECEnabled: &CDN_DNSSEC_ENABLED,
		DomainName:    &CDN_DOMAIN_NAME,
		ID:            &CDN_ID,
		Name:          &CDN,
	},
}

// Static Cache Group fields
var EDGE_CACHEGROUP_ID = 1
var EDGE_CACHEGROUP = "Edge"
var EDGE_CACHEGROUP_SHORT_NAME = "Edge"
var EDGE_CACHEGROUP_LATITUDE = 0.0
var EDGE_CACHEGROUP_LONGITUDE = 0.0
var EDGE_CACHEGROUP_PARENT_NAME = "Mid" // NOTE: This places a hard requirement on the `cachegroups` implementation - must have a `MID_LOC` Cache Group named "Mid"
var EDGE_CACHEGROUP_PARENT_ID = 2       // NOTE: This places a hard requirement on the `cachegroups` implementation - must have a `MID_LOC` Cache Group identified by `2`
var EDGE_CACHEGROUP_FALLBACK_TO_CLOSEST = true
var EDGE_CACHEGROUP_LOCALIZATION_METHODS = []tc.LocalizationMethod{
	tc.LocalizationMethodCZ,
	tc.LocalizationMethodDeepCZ,
	tc.LocalizationMethodGeo,
}
var EDGE_CACHEGROUP_TYPE = "EDGE_LOC"
var EDGE_CACHEGROUP_TYPE_ID = 1 // NOTE: This places a hard requirement on the `types` implementation - must have `EDGE_LOC` == 1

var MID_CACHEGROUP_ID = 2
var MID_CACHEGROUP = "Mid"
var MID_CACHEGROUP_SHORT_NAME = "Mid"
var MID_CACHEGROUP_LATITUDE = 0.0
var MID_CACHEGROUP_LONGITUDE = 0.0
var MID_CACHEGROUP_FALLBACK_TO_CLOSEST = true
var MID_CACHEGROUP_LOCALIZATION_METHODS = []tc.LocalizationMethod{
	tc.LocalizationMethodCZ,
	tc.LocalizationMethodDeepCZ,
	tc.LocalizationMethodGeo,
}
var MID_CACHEGROUP_TYPE = "MID_LOC"
var MID_CACHEGROUP_TYPE_ID = 2 // NOTE: This places a hard requirement on the `types` implementation - must have `MID_LOC` == 2

var CACHEGROUPS = []tc.CacheGroupNullable{
	tc.CacheGroupNullable{
		ID:                          &EDGE_CACHEGROUP_ID,
		Name:                        &EDGE_CACHEGROUP,
		ShortName:                   &EDGE_CACHEGROUP_SHORT_NAME,
		Latitude:                    &EDGE_CACHEGROUP_LATITUDE,
		Longitude:                   &EDGE_CACHEGROUP_LONGITUDE,
		ParentName:                  &EDGE_CACHEGROUP_PARENT_NAME,
		ParentCachegroupID:          &EDGE_CACHEGROUP_PARENT_ID,
		SecondaryParentName:         nil,
		SecondaryParentCachegroupID: nil,
		FallbackToClosest:           &EDGE_CACHEGROUP_FALLBACK_TO_CLOSEST,
		LocalizationMethods:         &EDGE_CACHEGROUP_LOCALIZATION_METHODS,
		Type:                        &EDGE_CACHEGROUP_TYPE,
		TypeID:                      &EDGE_CACHEGROUP_TYPE_ID,
		LastUpdated:                 CURRENT_TIME,
		Fallbacks:                   nil,
	},
	tc.CacheGroupNullable{
		ID:                          &MID_CACHEGROUP_ID,
		Name:                        &MID_CACHEGROUP,
		ShortName:                   &MID_CACHEGROUP_SHORT_NAME,
		Latitude:                    &MID_CACHEGROUP_LATITUDE,
		Longitude:                   &MID_CACHEGROUP_LONGITUDE,
		ParentName:                  nil,
		ParentCachegroupID:          nil,
		SecondaryParentName:         nil,
		SecondaryParentCachegroupID: nil,
		FallbackToClosest:           &MID_CACHEGROUP_FALLBACK_TO_CLOSEST,
		LocalizationMethods:         &MID_CACHEGROUP_LOCALIZATION_METHODS,
		Type:                        &MID_CACHEGROUP_TYPE,
		TypeID:                      &MID_CACHEGROUP_TYPE_ID,
		LastUpdated:                 CURRENT_TIME,
		Fallbacks:                   nil,
	},
}

// Static profile fields
var GLOBAL_PROFILE_ID = 1
var GLOBAL_PROFILE_NAME = "GLOBAL"
var GLOBAL_PROFILE_DESCRIPTION = "Global Traffic Ops profile, DO NOT DELETE"
var GLOBAL_PROFILE_CDN_NAME = ALL_CDN
var GLOBAL_PROFILE_CDN_ID = ALL_CDN_ID
var GLOBAL_PROFILE_ROUTING_DISABLED = false
var GLOBAL_PROFILE_TYPE = "UNK_PROFILE"
var TO_PROFILE_ID = 2
var TO_PROFILE_NAME = "TRAFFIC_OPS"
var TO_PROFILE_DESCRIPTION = "Traffic Ops profile"
var TO_PROFILE_CDN_NAME = ALL_CDN
var TO_PROFILE_CDN_ID = ALL_CDN_ID
var TO_PROFILE_ROUTING_DISABLED = false
var TO_PROFILE_TYPE = "UNK_PROFILE"
var TO_DB_PROFILE_ID = 3
var TO_DB_PROFILE_NAME = "TRAFFIC_OPS_DB"
var TO_DB_PROFILE_DESCRIPTION = "Traffic Ops DB profile"
var TO_DB_PROFILE_CDN_NAME = ALL_CDN
var TO_DB_PROFILE_CDN_ID = ALL_CDN_ID
var TO_DB_PROFILE_ROUTING_DISABLED = false
var TO_DB_PROFILE_TYPE = "UNK_PROFILE"
var EDGE_PROFILE_ID = 4
var EDGE_PROFILE_NAME = "ATS_EDGE_TIER_CACHE"
var EDGE_PROFILE_DESCRIPTION = "Edge Cache - Apache Traffic Server"
var EDGE_PROFILE_CDN_NAME = CDN
var EDGE_PROFILE_CDN_ID = CDN_ID
var EDGE_PROFILE_ROUTING_DISABLED = false
var EDGE_PROFILE_TYPE = "ATS_PROFILE"
var MID_PROFILE_ID = 5
var MID_PROFILE_NAME = "ATS_MID_TIER_CACHE"
var MID_PROFILE_DESCRIPTION = "Mid Cache - Apache Traffic Server"
var MID_PROFILE_CDN_NAME = CDN
var MID_PROFILE_CDN_ID = CDN_ID
var MID_PROFILE_ROUTING_DISABLED = false
var MID_PROFILE_TYPE = "ATS_PROFILE"

var PROFILES = []tc.ProfileNullable{
	tc.ProfileNullable{
		ID:              &GLOBAL_PROFILE_ID,
		LastUpdated:     CURRENT_TIME,
		Name:            &GLOBAL_PROFILE_NAME,
		Description:     &GLOBAL_PROFILE_DESCRIPTION,
		CDNName:         &GLOBAL_PROFILE_CDN_NAME,
		CDNID:           &GLOBAL_PROFILE_CDN_ID,
		RoutingDisabled: &GLOBAL_PROFILE_ROUTING_DISABLED,
		Type:            &GLOBAL_PROFILE_TYPE,
		Parameters:      nil,
	},
	tc.ProfileNullable{
		ID:              &TO_PROFILE_ID,
		LastUpdated:     CURRENT_TIME,
		Name:            &TO_PROFILE_NAME,
		Description:     &TO_PROFILE_DESCRIPTION,
		CDNName:         &TO_PROFILE_CDN_NAME,
		CDNID:           &TO_PROFILE_CDN_ID,
		RoutingDisabled: &TO_PROFILE_ROUTING_DISABLED,
		Type:            &TO_PROFILE_TYPE,
		Parameters:      nil,
	},
	tc.ProfileNullable{
		ID:              &TO_DB_PROFILE_ID,
		LastUpdated:     CURRENT_TIME,
		Name:            &TO_DB_PROFILE_NAME,
		Description:     &TO_DB_PROFILE_DESCRIPTION,
		CDNName:         &TO_DB_PROFILE_CDN_NAME,
		CDNID:           &TO_DB_PROFILE_CDN_ID,
		RoutingDisabled: &TO_DB_PROFILE_ROUTING_DISABLED,
		Type:            &TO_DB_PROFILE_TYPE,
		Parameters:      nil,
	},
	tc.ProfileNullable{
		ID:              &EDGE_PROFILE_ID,
		LastUpdated:     CURRENT_TIME,
		Name:            &EDGE_PROFILE_NAME,
		Description:     &EDGE_PROFILE_DESCRIPTION,
		CDNName:         &EDGE_PROFILE_CDN_NAME,
		CDNID:           &EDGE_PROFILE_CDN_ID,
		RoutingDisabled: &EDGE_PROFILE_ROUTING_DISABLED,
		Type:            &EDGE_PROFILE_TYPE,
		Parameters:      nil,
	},
	tc.ProfileNullable{
		ID:              &MID_PROFILE_ID,
		LastUpdated:     CURRENT_TIME,
		Name:            &MID_PROFILE_NAME,
		Description:     &MID_PROFILE_DESCRIPTION,
		CDNName:         &MID_PROFILE_CDN_NAME,
		CDNID:           &MID_PROFILE_CDN_ID,
		RoutingDisabled: &MID_PROFILE_ROUTING_DISABLED,
		Type:            &MID_PROFILE_TYPE,
		Parameters:      nil,
	},
}

var PARAMETERS = []tc.Parameter{
	tc.Parameter{
		ConfigFile:  "global",
		ID:          1,
		LastUpdated: *CURRENT_TIME,
		Name:        "tm.url",
		Profiles:    json.RawMessage(`["` + GLOBAL_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "https://localhost:443/",
	},
	tc.Parameter{
		ConfigFile:  "global",
		ID:          2,
		LastUpdated: *CURRENT_TIME,
		Name:        "tm.instance_name",
		Profiles:    json.RawMessage(`["` + GLOBAL_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "Mock Traffic Ops CDN",
	},
	tc.Parameter{
		ConfigFile:  "global",
		ID:          3,
		LastUpdated: *CURRENT_TIME,
		Name:        "tm.toolname",
		Profiles:    json.RawMessage(`["` + GLOBAL_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "Mock Traffic Ops",
	},
	tc.Parameter{
		ConfigFile:  "CRConfig.json",
		ID:          4,
		LastUpdated: *CURRENT_TIME,
		Name:        "geolocation.polling.url",
		Profiles:    json.RawMessage(`["` + GLOBAL_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "https://localhost/mock/geo/database.dat",
	},
	tc.Parameter{
		ConfigFile:  "CRConfig.json",
		ID:          5,
		LastUpdated: *CURRENT_TIME,
		Name:        "geolocation6.polling.url",
		Profiles:    json.RawMessage(`["` + GLOBAL_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "https://localhost/mock/geo/database.dat",
	},
	tc.Parameter{
		ConfigFile:  "global",
		ID:          6,
		LastUpdated: *CURRENT_TIME,
		Name:        "tm.logourl",
		Profiles:    json.RawMessage(`["` + GLOBAL_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/logo.svg",
	},
	tc.Parameter{
		ConfigFile:  "global",
		ID:          7,
		LastUpdated: *CURRENT_TIME,
		Name:        "default_geo_miss_latitude",
		Profiles:    json.RawMessage(`["` + GLOBAL_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "0",
	},
	tc.Parameter{
		ConfigFile:  "global",
		ID:          8,
		LastUpdated: *CURRENT_TIME,
		Name:        "default_geo_miss_longitude",
		Profiles:    json.RawMessage(`["` + GLOBAL_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "0",
	},
	tc.Parameter{
		ConfigFile:  "regex_revalidate.config",
		ID:          9,
		LastUpdated: *CURRENT_TIME,
		Name:        "maxRevalDurationDays",
		Profiles:    json.RawMessage(`["` + GLOBAL_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "90",
	},
	tc.Parameter{
		ConfigFile:  "global",
		ID:          10,
		LastUpdated: *CURRENT_TIME,
		Name:        "use_reval_pending",
		Profiles:    json.RawMessage(`["` + GLOBAL_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "1",
	},
	tc.Parameter{
		ConfigFile:  "global",
		ID:          11,
		LastUpdated: *CURRENT_TIME,
		Name:        "use_tenancy",
		Profiles:    json.RawMessage(`["` + GLOBAL_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "1",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          12,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.http.server_ports",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "STRING 80 80:ipv6 443:proto=http:ssl 443:ipv6:proto=http:ssl",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          13,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.reverse_proxy.enabled",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "INT 1",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          14,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.ssl.CA.cert.path",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "STRING /etc/trafficserver/ssl",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          15,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.ssl.client.CA.cert.path",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "STRING /etc/trafficserver/ssl",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          16,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.ssl.client.cert.path",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "STRING /etc/trafficserver/ssl",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          17,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.ssl.client.private_key.path",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "STRING /etc/trafficserver/ssl",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          18,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.ssl.oscp.enabled",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "INT 1",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          19,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.ssl.server.cert.path",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "STRING /etc/trafficserver/ssl",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          20,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.ssl.server.ticket_key.filename",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "STRING NULL",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          21,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.ssl.server.private_key.path",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "STRING /etc/trafficserver/ssl",
	},
	tc.Parameter{
		ConfigFile:  "ssl_multicert.config",
		ID:          22,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          23,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.reverse_proxy.enabled",
		Profiles:    json.RawMessage(`["` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "INT 0",
	},
	tc.Parameter{
		ConfigFile:  "astats.config",
		ID:          24,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver",
	},
	tc.Parameter{
		ConfigFile:  "astats.config",
		ID:          25,
		LastUpdated: *CURRENT_TIME,
		Name:        "allow_ip",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "0.0.0.0/0",
	},
	tc.Parameter{
		ConfigFile:  "astats.config",
		ID:          26,
		LastUpdated: *CURRENT_TIME,
		Name:        "allow_ip6",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "::/0",
	},
	tc.Parameter{
		ConfigFile:  "astats.config",
		ID:          27,
		LastUpdated: *CURRENT_TIME,
		Name:        "path",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "_astats",
	},
	tc.Parameter{
		ConfigFile:  "astats.config",
		ID:          28,
		LastUpdated: *CURRENT_TIME,
		Name:        "record_types",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "122",
	},
	tc.Parameter{
		ConfigFile:  "cache.config",
		ID:          29,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver",
	},
	tc.Parameter{
		ConfigFile:  "chkconfig",
		ID:          30,
		LastUpdated: *CURRENT_TIME,
		Name:        "trafficserver",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "0:off\t1:off\t2:on\t3:on\t4:on\t5:on\t6:off",
	},
	tc.Parameter{
		ConfigFile:  "hosting.config",
		ID:          31,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver",
	},
	tc.Parameter{
		ConfigFile:  "ip_allow.config",
		ID:          32,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver",
	},
	tc.Parameter{
		ConfigFile:  "parent.config",
		ID:          33,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver",
	},
	tc.Parameter{
		ConfigFile:  "plugin.config",
		ID:          34,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver",
	},
	tc.Parameter{
		ConfigFile:  "plugin.config",
		ID:          35,
		LastUpdated: *CURRENT_TIME,
		Name:        "astats_over_http.so",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "",
	},
	tc.Parameter{
		ConfigFile:  "plugin.config",
		ID:          36,
		LastUpdated: *CURRENT_TIME,
		Name:        "regex_revalidate.so",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "--config regex_revalidate.config",
	},
	tc.Parameter{
		ConfigFile:  "rascal.properties",
		ID:          37,
		LastUpdated: *CURRENT_TIME,
		Name:        "health.connection.timeout",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "2000",
	},
	tc.Parameter{
		ConfigFile:  "rascal.properties",
		ID:          38,
		LastUpdated: *CURRENT_TIME,
		Name:        "health.polling.url",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "http://${hostname}/_astats?application=&inf.name=${interface_name}",
	},
	tc.Parameter{
		ConfigFile:  "rascal.properties",
		ID:          39,
		LastUpdated: *CURRENT_TIME,
		Name:        "health.threshold.availableBandwidthInKbps",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       ">1750000",
	},
	tc.Parameter{
		ConfigFile:  "rascal.properties",
		ID:          40,
		LastUpdated: *CURRENT_TIME,
		Name:        "health.threshold.loadavg",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "25.0",
	},
	tc.Parameter{
		ConfigFile:  "rascal.properties",
		ID:          41,
		LastUpdated: *CURRENT_TIME,
		Name:        "health.threshold.queryTime",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "1000",
	},
	tc.Parameter{
		ConfigFile:  "rascal.properties",
		ID:          42,
		LastUpdated: *CURRENT_TIME,
		Name:        "history.count",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "30",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          43,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          44,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.admin.user_id",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "STRING ats",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          45,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.body_factory.template_sets_dir",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "STRING /etc/trafficserver/body_factory",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          46,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.config_dir",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "STRING /etc/trafficserver",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          47,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.diags.debug.enabled",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "INT 1",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          48,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.dns.round_robin_nameservers",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "INT 0",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          49,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.exec_thread.autoconfig",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "INT 0",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          50,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.http.cache.required_headers",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "INT 0",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          51,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.http.enable_http_stats",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "INT 1",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          52,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.http.insert_response_via_str",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "INT 3",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          53,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.http.server_ports",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "STRING 80 80:ipv6",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          54,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.http.slow.log.threshold",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "INT 10000",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          55,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.http.transaction_active_timeout_in",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "INT 0",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          56,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.log.logfile_dir",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "STRING /var/log/trafficserver",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          57,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.http.parent_proxy.retry_time",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "INT 60",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          58,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.http.parent_proxy_routing_enable",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "INT 1",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          59,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.proxy_name",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "STRING __HOSTNAME__",
	},
	tc.Parameter{
		ConfigFile:  "records.config",
		ID:          60,
		LastUpdated: *CURRENT_TIME,
		Name:        "CONFIG proxy.config.url_remap.remap_required",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "INT 0",
	},
	tc.Parameter{
		ConfigFile:  "regex_revalidate.config",
		ID:          61,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver",
	},
	tc.Parameter{
		ConfigFile:  "remap.config",
		ID:          62,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_0.config",
		ID:          63,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_8.config",
		ID:          64,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_10.config",
		ID:          65,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_12.config",
		ID:          66,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_14.config",
		ID:          67,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_16.config",
		ID:          68,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_18.config",
		ID:          69,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_22.config",
		ID:          70,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_24.config",
		ID:          71,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_26.config",
		ID:          72,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_28.config",
		ID:          73,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_30.config",
		ID:          74,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_32.config",
		ID:          75,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_34.config",
		ID:          76,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_36.config",
		ID:          77,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_37.config",
		ID:          78,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_38.config",
		ID:          79,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_40.config",
		ID:          80,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_48.config",
		ID:          81,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "set_dscp_56.config",
		ID:          82,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver/dscp",
	},
	tc.Parameter{
		ConfigFile:  "storage.config",
		ID:          83,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver",
	},
	tc.Parameter{
		ConfigFile:  "storage.config",
		ID:          84,
		LastUpdated: *CURRENT_TIME,
		Name:        "Disk_Volume",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "1",
	},
	tc.Parameter{
		ConfigFile:  "storage.config",
		ID:          85,
		LastUpdated: *CURRENT_TIME,
		Name:        "Drive_Letters",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "cache",
	},
	tc.Parameter{
		ConfigFile:  "storage.config",
		ID:          86,
		LastUpdated: *CURRENT_TIME,
		Name:        "Drive_Prefix",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/var/trafficserver",
	},
	tc.Parameter{
		ConfigFile:  "volume.config",
		ID:          87,
		LastUpdated: *CURRENT_TIME,
		Name:        "location",
		Profiles:    json.RawMessage(`["` + EDGE_PROFILE_NAME + `","` + MID_PROFILE_NAME + `"]`),
		Secure:      false,
		Value:       "/etc/trafficserver",
	},
}

// Static user fields
// (These _should_ be `const`, but you can't take the address of a `const` (for some reason))
var USERNAME = "admin"
var LOCAL_USER = true
var USER_ID = 1
var TENANT = "root"
var TENANT_ID = 1 // NOTE: This places a hard requirement on `tenant` implementation - `root` == `1`
var ROLE = "admin"
var ROLE_ID = 1 // NOTE: This places a hard requirement on `roles` implementation - `admin` == `1`
var NEW_USER = false

var COMMON_USER_FIELDS = tc.CommonUserFields{
	AddressLine1:    nil,
	AddressLine2:    nil,
	City:            nil,
	Company:         nil,
	Country:         nil,
	Email:           nil,
	FullName:        nil,
	GID:             nil,
	ID:              &USER_ID,
	NewUser:         &NEW_USER,
	PhoneNumber:     nil,
	PostalCode:      nil,
	PublicSSHKey:    nil,
	Role:            &ROLE_ID,
	StateOrProvince: nil,
	Tenant:          &TENANT,
	TenantID:        &TENANT_ID,
	UID:             nil,
	LastUpdated:     CURRENT_TIME,
}

var CURRENT_USER = tc.UserCurrent{
	UserName:         &USERNAME,
	LocalUser:        &LOCAL_USER,
	RoleName:         &ROLE,
	CommonUserFields: COMMON_USER_FIELDS,
}

// Just sets some common headers
func common(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", "Traffic Ops/"+VERSION+" (Mock)")
}

func ping(w http.ResponseWriter, r *http.Request) {
	common(w)
	w.Write([]byte("{\"ping\":\"pong\"}\n"))
}

func login(w http.ResponseWriter, r *http.Request) {
	form := auth.PasswordForm{}
	common(w)
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		errBytes, JsonErr := json.Marshal(tc.CreateErrorAlerts(err))
		if JsonErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("Failed to create an alerts structure from '%v': %s\n", err, JsonErr.Error())
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errBytes)
		return
	}

	expiry := time.Now().Add(time.Hour * 6)
	cookie := tocookie.New(form.Username, expiry, "foo")
	httpCookie := http.Cookie{Name: "mojolicious", Value: cookie, Path: "/", Expires: expiry, HttpOnly: true}
	http.SetCookie(w, &httpCookie)
	resp := struct {
		tc.Alerts
	}{tc.CreateAlerts(tc.SuccessLevel, "Successfully logged in.")}
	respBts, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Couldn't marshal /login response: %s", err.Error())
	}
	w.Write(respBts)
}

func cuser(w http.ResponseWriter, r *http.Request) {
	common(w)
	api.WriteResp(w, r, CURRENT_USER)
}

func CDNs(w http.ResponseWriter, r *http.Request) {
	common(w)
	api.WriteResp(w, r, CDNS)
}

func cacheGroups(w http.ResponseWriter, r *http.Request) {
	common(w)
	api.WriteResp(w, r, CACHEGROUPS)
}

func parameters(w http.ResponseWriter, r *http.Request) {
	common(w)
	api.WriteResp(w, r, PARAMETERS)
}

func profiles(w http.ResponseWriter, r *http.Request) {
	common(w)
	api.WriteResp(w, r, PROFILES)
}

func mockDatabase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Traffic Ops/"+VERSION+" (Mock)")
	if r.Method == http.MethodGet {
		w.Write([]byte("Mock Traffic Ops servers don't have databases\n"))
	} else {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func logo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Traffic Ops/"+VERSION+" (Mock)")
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write([]byte(TM_LOGO))
	} else {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	server := &http.Server{
		Addr:      ":443",
		TLSConfig: &tls.Config{InsecureSkipVerify: false},
	}

	for i := API_MIN_MINOR_VERSION; i <= API_MAX_MINOR_VERSION; i += 1 {
		v := "1." + strconv.Itoa(i)
		log.Printf("Loading API v%s\n", v)
		http.HandleFunc("/api/"+v+"/ping", ping)
		http.HandleFunc("/api/"+v+"/user/login", login)
		http.HandleFunc("/api/"+v+"/user/current", cuser)
		http.HandleFunc("/api/"+v+"/cdns", CDNs)
		http.HandleFunc("/api/"+v+"/cachegroups", cacheGroups)
		http.HandleFunc("/api/"+v+"/profiles", profiles)
		http.HandleFunc("/api/"+v+"/parameters", parameters)
	}

	http.HandleFunc("/mock/geo/database.dat", mockDatabase)
	http.HandleFunc("/logo.svg", logo)

	log.Printf("Finished loading API routes at %s, server listening on port 443", CURRENT_TIME.String())

	if err := server.ListenAndServeTLS("./localhost.crt", "./localhost.key"); err != nil {
		log.Fatalf("Server crashed: %v\n", err)
	}
	os.Exit(0)
}
