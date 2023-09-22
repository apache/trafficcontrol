package varnishcfg

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
	"reflect"
	"testing"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

func TestConfigureUncacheableDSes(t *testing.T) {
	testCases := []struct {
		name                string
		dses                []atscfg.ProfileDS
		expectedSubroutines map[string][]string
	}{
		{
			name: "no uncacheable DSes",
			dses: []atscfg.ProfileDS{
				{Name: "ds1", Type: tc.DSTypeHTTP, OriginFQDN: "ds1.example.com"},
				{Name: "ds2", Type: tc.DSTypeHTTP, OriginFQDN: "ds2.example.com"},
			},
			expectedSubroutines: make(map[string][]string),
		},
		{
			name: "single uncacheable DS",
			dses: []atscfg.ProfileDS{
				{Name: "ds1", Type: tc.DSTypeHTTP, OriginFQDN: "ds1.example.com"},
				{Name: "ds2", Type: tc.DSTypeHTTPNoCache, OriginFQDN: "ds2.example.com"},
			},
			expectedSubroutines: map[string][]string{
				"vcl_backend_response": {
					`if (bereq.http.host == "ds2.example.com") {`,
					`	set beresp.uncacheable = true;`,
					`}`,
				},
			},
		},
		{
			name: "multiple uncacheable DS",
			dses: []atscfg.ProfileDS{
				{Name: "ds1", Type: tc.DSTypeHTTP, OriginFQDN: "ds1.example.com"},
				{Name: "ds2", Type: tc.DSTypeHTTPNoCache, OriginFQDN: "ds2.example.com"},
				{Name: "ds3", Type: tc.DSTypeHTTPNoCache, OriginFQDN: "ds3.example.com"},
				{Name: "ds4", Type: tc.DSTypeHTTPNoCache, OriginFQDN: "ds4.example.com"},
				{Name: "ds5", Type: tc.DSTypeDNS, OriginFQDN: "ds5.example.com"},
			},
			expectedSubroutines: map[string][]string{
				"vcl_backend_response": {
					`if (bereq.http.host == "ds2.example.com" || bereq.http.host == "ds3.example.com" || bereq.http.host == "ds4.example.com") {`,
					`	set beresp.uncacheable = true;`,
					`}`,
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			vb := NewVCLBuilder(&t3cutil.ConfigData{})
			vclFile := newVCLFile(defaultVCLVersion)
			warnings := vb.configureUncacheableDSes(&vclFile, tC.dses)
			if len(warnings) != 0 {
				t.Errorf("got unexpected warnings %v", warnings)
			}
			if !reflect.DeepEqual(vclFile.subroutines, tC.expectedSubroutines) {
				t.Errorf("got %v want %v", vclFile.subroutines, tC.expectedSubroutines)
			}
		})
	}
}
