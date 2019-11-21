package atscfg

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
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestMakeHeaderRewriteDotConfig(t *testing.T) {
	cdnName := tc.CDNName("mycdn")
	toToolName := "my-to"
	toURL := "my-to.example.net"

	ds := HeaderRewriteDS{
		EdgeHeaderRewrite:    "edgerewrite",
		ID:                   24,
		MaxOriginConnections: 42,
		MidHeaderRewrite:     "midrewrite",
		Type:                 tc.DSTypeHTTPLive,
	}
	assignedEdges := []HeaderRewriteServer{
		HeaderRewriteServer{
			Status: tc.CacheStatusReported,
		},
		HeaderRewriteServer{
			Status: tc.CacheStatusOnline,
		},
		HeaderRewriteServer{
			Status: tc.CacheStatusOffline,
		},
	}

	txt := MakeHeaderRewriteDotConfig(cdnName, toToolName, toURL, ds, assignedEdges)

	if !strings.Contains(txt, "edgerewrite") {
		t.Errorf("expected 'edgerewrite' actual '%v'\n", txt)
	}

	if strings.Contains(txt, "midrewrite") {
		t.Errorf("expected no 'midrewrite' actual '%v'\n", txt)
	}

	if !strings.Contains(txt, "origin_max_connections") {
		t.Errorf("expected origin_max_connections on edge header rewrite that uses the mids, actual '%v'\n", txt)
	}

	if !strings.Contains(txt, "21") { // 21, because max is 42, and there are 2 not-offline mids, so 42/2=21
		t.Errorf("expected origin_max_connections of 21, actual '%v'\n", txt)
	}
}

func TestMakeHeaderRewriteDotConfigNoMaxOriginConnections(t *testing.T) {
	cdnName := tc.CDNName("mycdn")
	toToolName := "my-to"
	toURL := "my-to.example.net"

	ds := HeaderRewriteDS{
		EdgeHeaderRewrite:    "edgerewrite",
		ID:                   24,
		MaxOriginConnections: 42,
		MidHeaderRewrite:     "midrewrite",
		Type:                 tc.DSTypeHTTP,
	}
	assignedEdges := []HeaderRewriteServer{
		HeaderRewriteServer{
			Status: tc.CacheStatusReported,
		},
		HeaderRewriteServer{
			Status: tc.CacheStatusOnline,
		},
		HeaderRewriteServer{
			Status: tc.CacheStatusOffline,
		},
	}

	txt := MakeHeaderRewriteDotConfig(cdnName, toToolName, toURL, ds, assignedEdges)

	if strings.Contains(txt, "origin_max_connections") {
		t.Errorf("expected no origin_max_connections on DS that uses the mid, actual '%v'\n", txt)
	}
}
