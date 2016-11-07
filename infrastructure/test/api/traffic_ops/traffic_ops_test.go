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


import (
	"github.com/Comcast/traffic_control/infrastructure/test/apitest"
	"strconv"
	"testing"
)

// TestExistingRegions tests the `regions.json` endpoint for values
// which should exist in a newly installed Traffic Ops.
func TestExistingRegions(t *testing.T) {
	at := NewApiTester(t)
	err := apitest.TestJSONContains(at, "regions.json", map[string]interface{}{
		"response": []interface{}{
			map[string]interface{}{
				"name": "Eastish",
				"id":   "19", // TODO fix TestJSONContains to not require ID
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
}

// ApiTestExistingCdns tests the `cdns.json` endpoint for values
// which should exist in a newly installed Traffic Ops.
func TestExistingCdns(t *testing.T) {
	at := NewApiTester(t)
	err := apitest.TestJSONContains(at, "cdns.json", map[string]interface{}{
		"response": []interface{}{
			map[string]interface{}{
				"name": "cdn",
				"id":   "1",
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
}

// ApiTestNewCdn tests creating a new CDN with the `/cdn/create` path
// and tests the new CDN exists at the `cdns.json` API endpoint.
func TestNewCdn(t *testing.T) {
	at := NewApiTester(t)

	testCdnName := "testcdn"
	createEndpoint := "/cdn/create"
	apiGetEndpoint := "cdns.json"
	err := apitest.DoPOST(at, createEndpoint, map[string]string{
		"cdn_data.name": testCdnName,
	})
	if err != nil {
		t.Error(err)
	}

	testCdnId, err := apitest.GetJSONID(at, apiGetEndpoint, testCdnName)
	if err != nil {
		t.Errorf("POST %s to %s succeeded, but GET %s didn't contain posted CDN", testCdnName, createEndpoint, apiGetEndpoint)
	}

	err = apitest.TestJSONContains(at, apiGetEndpoint, map[string]interface{}{
		"response": []interface{}{
			map[string]interface{}{
				"id":   strconv.Itoa(testCdnId),
				"name": testCdnName,
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
}
