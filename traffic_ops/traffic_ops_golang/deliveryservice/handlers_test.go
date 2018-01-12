package deliveryservice

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
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/utils"
)

// TestValidateErrors ...
func TestValidateErrors(t *testing.T) {

	ds := GetRefType()
	if err := json.Unmarshal([]byte(errorTestCase()), &ds); err != nil {
		fmt.Printf("err ---> %v\n", err)
		return
	}

	errors := ds.Validate(nil)
	errorStrs := utils.ErrorsToStrings(errors)
	sort.Strings(errorStrs)
	errorsFmt, _ := json.MarshalIndent(errorStrs, "", "  ")
	fmt.Printf("returned errors ---> %v\n", string(errorsFmt))

	expected := []string{
		"'active' is required",
		"'cdnId' is required",
		"'initialDispersion' must be greater than zero",
		"'routingName' the length must be between 1 and 48",
		"'xmlId' the length must be between 1 and 48",
	}
	sort.Strings(expected)
	expectedFmt, _ := json.MarshalIndent(expected, "", "  ")

	same := reflect.DeepEqual(expected, errorStrs)
	if !same {
		t.Errorf("\nExpected %s \n Actual %v", string(expectedFmt), string(errorsFmt))
	}

}

func errorTestCase() string {

	routingName := strings.Repeat("X", 49)

	// Test the xmlId length
	xmlId := strings.Repeat("X", 49)

	displayName := strings.Repeat("X", 49)

	errorTestCase := `
{
   "ccrDnsTtl": 1,
   "checkPath": "/crossdomain.xml",
   "displayName": "` + displayName + `",
   "dnsBypassCname": "cname",
   "dnsBypassIp": "127.0.0.1",
   "dnsBypassIp6": "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
   "dnsBypassTTL": 10,
   "dscp": 0,
   "edgeHeaderRewrite": "cond %{REMAP_PSEUDO_HOOK} __RETURN__ set-config proxy.config.http.transaction_active_timeout_in 10800 [L]",
   "geoLimit": 0,
   "geoLimitCountries": "Can,Mex",
   "geoRedirectURL": "http://localhost/redirect",
   "geoProvider": 0,
   "globalMaxMBPS": 0,
   "globalMaxTPS": 0,
   "httpBypassFqdn": "http://bypass",
   "id": 1,
   "initialDispersion": 0,
   "infoUrl": "http://info.url",
   "ipv6RoutingEnabled": false,
   "lastUpdated": "2017-01-05 15:04:05+00",
   "logsEnabled": true,
   "longDesc": "longdesc",
   "longDesc1": "longdesc1",
   "longDesc2": "longdesc2",
   "maxDnsAnswers": 5,
   "midHeaderRewrite": "cond %{REMAP_PSEUDO_HOOK} __RETURN__ set-config proxy.config.http.cache.ignore_authentication 1 __RETURN__ set-config proxy.config.http.auth_server_session_private 0 __RETURN__ set-config proxy.config.http.transaction_no_activity_timeout_out 10 __RETURN__ set-config proxy.config.http.transaction_active_timeout_out 10  [L] __RETURN__",
   "missLat": -2.0,
   "missLong": -1.0,
   "multiSiteOrigin": false,
   "multiSiteOriginAlgorithm": 1,
   "orgServerFqdn": "http://localhost",
   "profile": 1,
   "protocol": 2,
   "qstringIgnore": 1,
   "rangeRequestHandling": 1,
   "regexRemap": "^/([^\/]+)/(.*) http://$1.foo.com/$2",
   "regionalGeoBlocking": false,
   "remapText": "@action=allow @src_ip=127.0.0.1-127.0.0.1",
   "routingName": "` + routingName + `",
   "signingAlgorithm": "url_sig",
   "sslKeyVersion": 1,
   "tenantId": 1,
   "trRequestHeaders": "xyz",
   "trResponseHeaders": "Access-Control-Allow-Origin: *",
   "typeId": 1,
   "xmlId": "` + xmlId + `"
 }
`
	return errorTestCase
}

func goodTestCase() string {

	goodTestCase := `
{
   "active": true,
   "ccrDnsTtl": 1,
   "cdnId": 1,
   "checkPath": "disp1",
   "displayName": "/crossdomain.xml",
   "dnsBypassCname": "cname",
   "dnsBypassIp": "127.0.0.1",
   "dnsBypassIp6": "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
   "dnsBypassTTL": 10,
   "dscp": 0,
   "edgeHeaderRewrite": "cond %{REMAP_PSEUDO_HOOK} __RETURN__ set-config proxy.config.http.transaction_active_timeout_in 10800 [L]",
   "geoLimit": 0,
   "geoLimitCountries": "Can,Mex",
   "geoRedirectURL": "http://localhost/redirect",
   "geoProvider": 0,
   "globalMaxMBPS": 0,
   "globalMaxTPS": 0,
   "httpBypassFqdn": "http://bypass",
   "id": 1,
   "infoUrl": "http://info.url",
   "ipv6RoutingEnabled": false,
   "lastUpdated": "2017-01-05 15:04:05+00",
   "logsEnabled": true,
   "longDesc": "longdesc",
   "longDesc1": "longdesc1",
   "longDesc2": "longdesc2",
   "maxDnsAnswers": 5,
   "midHeaderRewrite": "cond %{REMAP_PSEUDO_HOOK} __RETURN__ set-config proxy.config.http.cache.ignore_authentication 1 __RETURN__ set-config proxy.config.http.auth_server_session_private 0 __RETURN__ set-config proxy.config.http.transaction_no_activity_timeout_out 10 __RETURN__ set-config proxy.config.http.transaction_active_timeout_out 10  [L] __RETURN__",
   "missLat": -2.0,
   "missLong": -1.0,
   "multiSiteOrigin": false,
   "multiSiteOriginAlgorithm": 1,
   "orgServerFqdn": "http://localhost",
   "profile": 1,
   "protocol": 2,
   "qstringIgnore": 1,
   "rangeRequestHandling": 1,
   "regexRemap": "^/([^\/]+)/(.*) http://$1.foo.com/$2",
   "regionalGeoBlocking": false,
   "remapText": "@action=allow @src_ip=127.0.0.1-127.0.0.1",
   "routingName": "ccr",
   "signingAlgorithm": "url_sig",
   "sslKeyVersion": 1,
   "tenantId": 1,
   "trRequestHeaders": "xyz",
   "trResponseHeaders": "Access-Control-Allow-Origin: *",
   "typeId": 1,
   "xmlId": "ds1"
}
`
	return goodTestCase
}
