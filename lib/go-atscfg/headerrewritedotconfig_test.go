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
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestMakeHeaderRewriteDotConfig(t *testing.T) {
	xmlID := "xml-id"
	fileName := "hdr_rw_" + xmlID + ".config"
	cdnName := "mycdn"
	hdr := "myHeaderComment"

	server := makeGenericServer()
	server.CDNName = &cdnName

	server.HostName = util.StrPtr("my-edge")
	server.ID = util.IntPtr(990)
	serverStatus := string(tc.CacheStatusReported)
	server.Status = &serverStatus
	server.CDNName = &cdnName

	ds := makeGenericDS()
	ds.EdgeHeaderRewrite = util.StrPtr("edgerewrite")
	ds.ID = util.IntPtr(240)
	ds.XMLID = &xmlID
	ds.MaxOriginConnections = util.IntPtr(42)
	ds.MidHeaderRewrite = util.StrPtr("midrewrite")
	ds.CDNName = &cdnName
	dsType := tc.DSTypeHTTPLive
	ds.Type = &dsType
	ds.ServiceCategory = util.StrPtr("servicecategory")

	sv1 := makeGenericServer()
	sv1.HostName = util.StrPtr("my-edge-1")
	sv1.CDNName = &cdnName
	sv1.ID = util.IntPtr(991)
	sv1Status := string(tc.CacheStatusOnline)
	sv1.Status = &sv1Status

	sv2 := makeGenericServer()
	sv2.HostName = util.StrPtr("my-edge-2")
	sv2.CDNName = &cdnName
	sv2.ID = util.IntPtr(992)
	sv2Status := string(tc.CacheStatusOffline)
	sv2.Status = &sv2Status

	servers := []Server{*server, *sv1, *sv2}
	dses := []tc.DeliveryServiceNullableV30{*ds}

	dss := makeDSS(servers, dses)

	cfg, err := MakeHeaderRewriteDotConfig(fileName, dses, dss, server, servers, hdr)

	if err != nil {
		t.Errorf("error expected nil, actual '%v'\n", err)
	}

	txt := cfg.Text

	if !strings.Contains(txt, "edgerewrite") {
		t.Errorf("expected 'edgerewrite' actual '%v'\n", txt)
	}

	if strings.Contains(txt, "midrewrite") {
		t.Errorf("expected no 'midrewrite' actual '%v'\n", txt)
	}

	if !strings.Contains(txt, "origin_max_connections") {
		t.Errorf("expected origin_max_connections on edge header rewrite that doesn't use the mids, actual '%v'\n", txt)
	}

	if !strings.Contains(txt, "21") { // 21, because max is 42, and there are 2 not-offline mids, so 42/2=21
		t.Errorf("expected origin_max_connections of 21, actual '%v'\n", txt)
	}

	if !strings.Contains(txt, "xml-id|servicecategory") {
		t.Errorf("expected 'xml-id|servicecategory' actual '%v'\n", txt)
	}
}

func TestMakeHeaderRewriteDotConfigNoMaxOriginConnections(t *testing.T) {
	xmlID := "xml-id"
	fileName := "hdr_rw_" + xmlID + ".config"
	cdnName := "mycdn"
	hdr := "myHeaderComment"

	server := makeGenericServer()
	server.CDNName = &cdnName

	server.HostName = util.StrPtr("my-edge")
	server.ID = util.IntPtr(990)
	serverStatus := string(tc.CacheStatusReported)
	server.Status = &serverStatus
	server.CDNName = &cdnName

	ds := makeGenericDS()
	ds.EdgeHeaderRewrite = util.StrPtr("edgerewrite")
	ds.ID = util.IntPtr(240)
	ds.XMLID = &xmlID
	ds.MaxOriginConnections = util.IntPtr(42)
	ds.MidHeaderRewrite = util.StrPtr("midrewrite")
	ds.CDNName = &cdnName
	dsType := tc.DSTypeHTTP
	ds.Type = &dsType
	ds.ServiceCategory = util.StrPtr("servicecategory")

	sv1 := makeGenericServer()
	sv1.HostName = util.StrPtr("my-edge-1")
	sv1.CDNName = &cdnName
	sv1.ID = util.IntPtr(991)
	sv1Status := string(tc.CacheStatusOnline)
	sv1.Status = &sv1Status

	sv2 := makeGenericServer()
	sv2.HostName = util.StrPtr("my-edge-2")
	sv2.CDNName = &cdnName
	sv2.ID = util.IntPtr(992)
	sv2Status := string(tc.CacheStatusOffline)
	sv2.Status = &sv2Status

	servers := []Server{*server, *sv1, *sv2}
	dses := []tc.DeliveryServiceNullableV30{*ds}

	dss := makeDSS(servers, dses)

	cfg, err := MakeHeaderRewriteDotConfig(fileName, dses, dss, server, servers, hdr)

	if err != nil {
		t.Errorf("error expected nil, actual '%v'\n", err)
	}

	txt := cfg.Text

	if strings.Contains(txt, "origin_max_connections") {
		t.Errorf("expected no origin_max_connections on DS that uses the mid, actual '%v'\n", txt)
	}
}
