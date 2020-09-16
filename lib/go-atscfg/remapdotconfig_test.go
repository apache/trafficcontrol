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

func TestMakeRemapDotConfig(t *testing.T) {
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(0)
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = util.StrPtr("myedgeheaderrewrite")
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(0)
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.Protocol = util.IntPtr(0)
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)
	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern",
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeyparamname",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyparamval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Errorf("expected one line for each remap plus a comment, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "http://myregexpattern") {
		t.Errorf("expected to contain regex pattern, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "origin.example.test") {
		t.Errorf("expected to contain origin FQDN, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigMidLiveLocalExcluded(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(0)
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = util.StrPtr("myedgeheaderrewrite")
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(0)
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.Protocol = util.IntPtr(0)
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)
	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern",
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeyparamname",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyparamval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected no remap lines for LIVE local DS, actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigMid(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(0)
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = util.StrPtr("myedgeheaderrewrite")
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(0)
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.Protocol = util.IntPtr(0)
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)
	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern",
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeyparamname",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyparamval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected one line for each remap plus a comment, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Count(remapLine, "origin.example.test") != 2 {
		t.Errorf("expected to contain origin FQDN twice (Mids remap origins to themselves, as a forward proxy), actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "hdr_rw_mid_"+"mydsname"+".config") {
		t.Errorf("expected to contain header rewrite for DS, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigNilOrigin(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = nil
	ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(0)
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = util.StrPtr("myedgeheaderrewrite")
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(0)
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.Protocol = util.IntPtr(0)
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)
	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern",
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeyparamname",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyparamval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected no remap lines for DS with nil Origin FQDN, actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigEmptyOrigin(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("")
	ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(0)
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = util.StrPtr("myedgeheaderrewrite")
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(0)
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.Protocol = util.IntPtr(0)
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)
	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern",
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeyparamname",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyparamval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected no remap lines for DS with empty Origin FQDN, actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigDuplicateOrigins(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(0)
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = util.StrPtr("myedgeheaderrewrite")
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(0)
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.Protocol = util.IntPtr(0)
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	ds2 := tc.DeliveryServiceNullableV30{}
	ds2.ID = util.IntPtr(49)
	dsType2 := tc.DSType("HTTP_LIVE_NATNL")
	ds2.Type = &dsType2
	ds2.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds2.MidHeaderRewrite = util.StrPtr("mymidrewrite2")
	ds2.CacheURL = util.StrPtr("mycacheurl2")
	ds2.RangeRequestHandling = util.IntPtr(0)
	ds2.RemapText = util.StrPtr("myremaptext")
	ds2.EdgeHeaderRewrite = util.StrPtr("myedgeheaderrewrite")
	ds2.SigningAlgorithm = util.StrPtr("url_sig")
	ds2.XMLID = util.StrPtr("mydsname")
	ds2.QStringIgnore = util.IntPtr(0)
	ds2.RegexRemap = util.StrPtr("myregexremap")
	ds2.FQPacingRate = util.IntPtr(0)
	ds2.DSCP = util.IntPtr(0)
	ds2.RoutingName = util.StrPtr("myroutingname")
	ds2.MultiSiteOrigin = util.BoolPtr(false)
	ds2.OriginShield = util.StrPtr("myoriginshield")
	ds2.ProfileID = util.IntPtr(49)
	ds2.Protocol = util.IntPtr(0)
	ds2.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds2.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds, ds2}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds2.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern",
				},
			},
		},
		tc.DeliveryServiceRegexes{
			DSName: *ds2.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern2",
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeyparamname",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyparamval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remap lines for multiple DSes with the same Origin (ATS can't handle multiple remaps with the same origin FQDN), actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigNilMidRewrite(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = nil
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(0)
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = util.StrPtr("myedgeheaderrewrite")
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(0)
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.Protocol = util.IntPtr(0)
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern",
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeyparamname",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyparamval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Errorf("expected one line for each remap plus a comment, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Count(remapLine, "origin.example.test") != 2 {
		t.Errorf("expected to contain origin FQDN twice (Mids remap origins to themselves, as a forward proxy), actual '%v'", txt)
	}

	if strings.Contains(remapLine, "hdr_rw_mid_") {
		t.Errorf("expected no 'hdr_rw_mid_' for nil mid header rewrite on DS, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "myedgeheaderrewrite") {
		t.Errorf("expected no edge header rewrite text for mid server, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "hdr_rw_") {
		t.Errorf("expected no edge header rewrite for mid server, actual '%v'", txt)
	}

}

func TestMakeRemapDotConfigMidHasNoEdgeRewrite(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = nil
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(0)
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(0)
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.Protocol = util.IntPtr(0)
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern",
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeyparamname",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyparamval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Errorf("expected one line for each remap plus a comment, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Count(remapLine, "origin.example.test") != 2 {
		t.Errorf("expected to contain origin FQDN twice (Mids remap origins to themselves, as a forward proxy), actual '%v'", txt)
	}

	if strings.Contains(remapLine, "hdr_rw_mid_") {
		t.Errorf("expected no 'hdr_rw_mid_' for nil mid header rewrite on DS, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigMidQStringPassUpATS7CacheKey(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(0)
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.Protocol = util.IntPtr(0)
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern",
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "6",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeyparamname",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyparamval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Errorf("expected one line for each remap plus a comment, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Count(remapLine, "origin.example.test") != 2 {
		t.Errorf("expected to contain origin FQDN twice (Mids remap origins to themselves, as a forward proxy), actual '%v'", txt)
	}

	if strings.Contains(remapLine, "hdr_rw_mid_") {
		t.Errorf("expected no 'hdr_rw_mid_' for nil mid header rewrite on DS, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cachekey") {
		t.Errorf("expected 'cachekey' for qstring pass up and ATS 6+, actual '%v'", txt)
	}
	if strings.Contains(remapLine, "cacheurl") {
		t.Errorf("expected no 'cacheurl' for ATS 6+, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigMidQStringPassUpATS5CacheURL(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(0)
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.Protocol = util.IntPtr(0)
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern",
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "5",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeyparamname",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyparamval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Errorf("expected one line for each remap plus a comment, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Count(remapLine, "origin.example.test") != 2 {
		t.Errorf("expected to contain origin FQDN twice (Mids remap origins to themselves, as a forward proxy), actual '%v'", txt)
	}

	if strings.Contains(remapLine, "hdr_rw_mid_") {
		t.Errorf("expected no 'hdr_rw_mid_' for nil mid header rewrite on DS, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cacheurl") {
		t.Errorf("expected 'cacheurl' for qstring pass up and ATS <=6, actual '%v'", txt)
	}
	if strings.Contains(remapLine, "cachekey") {
		t.Errorf("expected no 'cachekey' for ATS <=6, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigMidProfileCacheKey(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(0)
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(0)
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern",
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Errorf("expected one line for each remap plus a comment, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Count(remapLine, "origin.example.test") != 2 {
		t.Errorf("expected to contain origin FQDN twice (Mids remap origins to themselves, as a forward proxy), actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cachekeykey") {
		t.Errorf("expected to contain cachekey parameter, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cachekeyval") {
		t.Errorf("expected to contain cachekey parameter value, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "shouldnotexist") {
		t.Errorf("expected to not contain cachekey parameter for different DS profile, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigMidRangeRequestHandling(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(0)
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern",
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Errorf("expected one line for each remap plus a comment, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Count(remapLine, "origin.example.test") != 2 {
		t.Errorf("expected to contain origin FQDN twice (Mids remap origins to themselves, as a forward proxy), actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cache_range_requests.so") {
		t.Errorf("expected to contain range request handling plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigMidSlicePluginRangeRequestHandling(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingSlice)
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(0)
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern",
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected one line for each remap plus a comment, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Count(remapLine, "origin.example.test") != 2 {
		t.Errorf("expected to contain origin FQDN twice (Mids remap origins to themselves, as a forward proxy), actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cache_range_requests.so") {
		t.Errorf("expected to contain range request handling plugin, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "slice.so") {
		t.Errorf("expected to not contain range request handling slice plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigFirstExcludedSecondIncluded(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("") // should be excluded
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingSlice)
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(0)
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	ds2 := tc.DeliveryServiceNullableV30{}
	ds2.ID = util.IntPtr(48)
	dsType2 := tc.DSType("HTTP_LIVE_NATNL")
	ds2.Type = &dsType2
	ds2.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds2.MidHeaderRewrite = util.StrPtr("")
	ds2.CacheURL = util.StrPtr("")
	ds2.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingSlice)
	ds2.RemapText = util.StrPtr("myremaptext")
	ds2.EdgeHeaderRewrite = nil
	ds2.SigningAlgorithm = util.StrPtr("url_sig")
	ds2.XMLID = util.StrPtr("mydsname")
	ds2.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds2.RegexRemap = util.StrPtr("myregexremap")
	ds2.FQPacingRate = util.IntPtr(0)
	ds2.DSCP = util.IntPtr(0)
	ds2.RoutingName = util.StrPtr("myroutingname")
	ds2.MultiSiteOrigin = util.BoolPtr(false)
	ds2.OriginShield = util.StrPtr("myoriginshield")
	ds2.ProfileID = util.IntPtr(49)
	ds2.ProfileName = util.StrPtr("dsprofile")
	ds2.Protocol = util.IntPtr(0)
	ds2.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds2.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds, ds2}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern",
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected one remap line for DS with origin, but not DS with empty Origin FQDN, actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigAnyMap(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("ANY_MAP")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingSlice)
	ds.RemapText = util.StrPtr("") // should not be included, any map requires remap text
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(0)
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	ds2 := tc.DeliveryServiceNullableV30{}
	ds2.ID = util.IntPtr(49)
	dsType2 := tc.DSType("ANY_MAP")
	ds2.Type = &dsType2
	ds2.OrgServerFQDN = util.StrPtr("myorigin")
	ds2.MidHeaderRewrite = util.StrPtr("mymidrewrite")
	ds2.CacheURL = util.StrPtr("mycacheurl")
	ds2.RangeRequestHandling = util.IntPtr(0)
	ds2.RemapText = util.StrPtr("myremaptext")
	ds2.EdgeHeaderRewrite = util.StrPtr("myedgerewrite")
	ds2.SigningAlgorithm = util.StrPtr("url_sig")
	ds2.XMLID = util.StrPtr("mydsname2")
	ds2.QStringIgnore = util.IntPtr(0)
	ds2.RegexRemap = util.StrPtr("myregexremap")
	ds2.FQPacingRate = util.IntPtr(0)
	ds2.DSCP = util.IntPtr(0)
	ds2.RoutingName = util.StrPtr("myroutingname")
	ds2.MultiSiteOrigin = util.BoolPtr(false)
	ds2.OriginShield = util.StrPtr("myoriginshield")
	ds2.ProfileID = util.IntPtr(49)
	ds2.ProfileName = util.StrPtr("dsprofile")
	ds2.Protocol = util.IntPtr(0)
	ds2.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds2.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds, ds2}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds2.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern",
				},
			},
		},
		tc.DeliveryServiceRegexes{
			DSName: *ds2.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern2",
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)
	txt = strings.Replace(txt, "\n\n", "\n", -1)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected one remap line for ANY_MAP DS with remap text, but not ANY_MAP DS with nil remap text, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to any_map to not start with 'map' (should be raw ds.RemapText), actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "myremaptext") {
		t.Errorf("expected to contain ANY_MAP DS remap text, actual '%v'", txt)
	}

}

func TestMakeRemapDotConfigEdgeMissingRemapData(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	dses := []tc.DeliveryServiceNullableV30{}
	{ // see regexes - has invalid regex type
		ds := tc.DeliveryServiceNullableV30{}
		ds.ID = util.IntPtr(1)
		dsType := tc.DSType("HTTP_LIVE_NATNL")
		ds.Type = &dsType
		ds.OrgServerFQDN = util.StrPtr("myorigin")
		ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
		ds.CacheURL = util.StrPtr("mycacheurl")
		ds.RangeRequestHandling = util.IntPtr(0)
		ds.RemapText = util.StrPtr("myreamptext")
		ds.EdgeHeaderRewrite = util.StrPtr("myedgeheaderrewrite")
		ds.SigningAlgorithm = util.StrPtr("url_sig")
		ds.XMLID = util.StrPtr("ds")
		ds.QStringIgnore = util.IntPtr(0)
		ds.RegexRemap = util.StrPtr("myregexremap")
		ds.FQPacingRate = util.IntPtr(0)
		ds.DSCP = util.IntPtr(0)
		ds.RoutingName = util.StrPtr("myroutingname")
		ds.MultiSiteOrigin = util.BoolPtr(false)
		ds.OriginShield = util.StrPtr("myoriginshield")
		ds.ProfileID = util.IntPtr(49)
		ds.ProfileName = util.StrPtr("dsprofile")
		ds.Protocol = util.IntPtr(0)
		ds.AnonymousBlockingEnabled = util.BoolPtr(false)
		ds.Active = util.BoolPtr(true)
		dses = append(dses, ds)
	}
	{ // see regexes - has invalid regex type
		ds := tc.DeliveryServiceNullableV30{}
		ds.ID = util.IntPtr(2)
		dsType := tc.DSType("HTTP_LIVE_NATNL")
		ds.Type = &dsType
		ds.OrgServerFQDN = util.StrPtr("myorigin")
		ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
		ds.CacheURL = util.StrPtr("mycacheurl")
		ds.RangeRequestHandling = util.IntPtr(0)
		ds.RemapText = util.StrPtr("myremaptext")
		ds.EdgeHeaderRewrite = util.StrPtr("myedgeheaderrewrite")
		ds.SigningAlgorithm = util.StrPtr("url_sig")
		ds.XMLID = util.StrPtr("ds2")
		ds.QStringIgnore = util.IntPtr(0)
		ds.RegexRemap = util.StrPtr("myregexremap")
		ds.FQPacingRate = util.IntPtr(0)
		ds.DSCP = util.IntPtr(0)
		ds.RoutingName = util.StrPtr("myroutingname")
		ds.MultiSiteOrigin = util.BoolPtr(false)
		ds.OriginShield = util.StrPtr("myoriginshield")
		ds.ProfileID = util.IntPtr(49)
		ds.ProfileName = util.StrPtr("dsprofile")
		ds.Protocol = util.IntPtr(0)
		ds.AnonymousBlockingEnabled = util.BoolPtr(false)
		ds.Active = util.BoolPtr(true)
		dses = append(dses, ds)
	}
	{ // see regexes - has invalid regex type
		ds := tc.DeliveryServiceNullableV30{}
		ds.ID = util.IntPtr(3)
		dsType := tc.DSType("HTTP_LIVE_NATNL")
		ds.Type = &dsType
		ds.OrgServerFQDN = util.StrPtr("myorigin")
		ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
		ds.CacheURL = util.StrPtr("mycacheurl")
		ds.RangeRequestHandling = util.IntPtr(0)
		ds.RemapText = util.StrPtr("myremaptext")
		ds.EdgeHeaderRewrite = util.StrPtr("myedgeheaderrewrite")
		ds.SigningAlgorithm = util.StrPtr("url_sig")
		ds.XMLID = util.StrPtr("ds3")
		ds.QStringIgnore = util.IntPtr(0)
		ds.RegexRemap = util.StrPtr("myregexremap")
		ds.FQPacingRate = util.IntPtr(0)
		ds.DSCP = util.IntPtr(0)
		ds.RoutingName = util.StrPtr("myroutingname")
		ds.MultiSiteOrigin = util.BoolPtr(false)
		ds.OriginShield = util.StrPtr("myoriginshield")
		ds.ProfileID = util.IntPtr(49)
		ds.ProfileName = util.StrPtr("dsprofile")
		ds.Protocol = util.IntPtr(0)
		ds.AnonymousBlockingEnabled = util.BoolPtr(false)
		ds.Active = util.BoolPtr(true)
		dses = append(dses, ds)
	}
	{ // see regexes - has invalid regex type
		ds := tc.DeliveryServiceNullableV30{}
		ds.ID = util.IntPtr(4)
		dsType := tc.DSType("HTTP_LIVE_NATNL")
		ds.Type = &dsType
		ds.OrgServerFQDN = util.StrPtr("myorigin")
		ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
		ds.CacheURL = util.StrPtr("mycacheurl")
		ds.RangeRequestHandling = util.IntPtr(0)
		ds.RemapText = util.StrPtr("myremaptext")
		ds.EdgeHeaderRewrite = util.StrPtr("myedgeheaderrewrite")
		ds.SigningAlgorithm = util.StrPtr("url_sig")
		ds.XMLID = util.StrPtr("ds4")
		ds.QStringIgnore = util.IntPtr(0)
		ds.RegexRemap = util.StrPtr("myregexremap")
		ds.FQPacingRate = util.IntPtr(0)
		ds.DSCP = util.IntPtr(0)
		ds.RoutingName = util.StrPtr("myroutingname")
		ds.MultiSiteOrigin = util.BoolPtr(false)
		ds.OriginShield = util.StrPtr("myoriginshield")
		ds.ProfileID = util.IntPtr(49)
		ds.ProfileName = util.StrPtr("dsprofile")
		ds.Protocol = util.IntPtr(0)
		ds.AnonymousBlockingEnabled = util.BoolPtr(false)
		ds.Active = util.BoolPtr(true)
		dses = append(dses, ds)
	}
	{ // see regexes - has invalid regex type
		ds := tc.DeliveryServiceNullableV30{}
		ds.ID = util.IntPtr(5)
		dsType := tc.DSType("HTTP_LIVE_NATNL")
		ds.Type = &dsType
		ds.OrgServerFQDN = util.StrPtr("myorigin")
		ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
		ds.CacheURL = util.StrPtr("mycacheurl")
		ds.RangeRequestHandling = util.IntPtr(0)
		ds.RemapText = util.StrPtr("myremaptext")
		ds.EdgeHeaderRewrite = util.StrPtr("myedgeheaderrewrite")
		ds.SigningAlgorithm = util.StrPtr("url_sig")
		ds.XMLID = util.StrPtr("ds5")
		ds.QStringIgnore = util.IntPtr(0)
		ds.RegexRemap = util.StrPtr("myregexremap")
		ds.FQPacingRate = util.IntPtr(0)
		ds.DSCP = util.IntPtr(0)
		ds.RoutingName = util.StrPtr("myroutingname")
		ds.MultiSiteOrigin = util.BoolPtr(false)
		ds.OriginShield = util.StrPtr("myoriginshield")
		ds.ProfileID = util.IntPtr(49)
		ds.ProfileName = util.StrPtr("dsprofile")
		ds.Protocol = util.IntPtr(0)
		ds.AnonymousBlockingEnabled = util.BoolPtr(false)
		ds.Active = util.BoolPtr(true)
		dses = append(dses, ds)
	}
	{
		ds := tc.DeliveryServiceNullableV30{}
		ds.ID = util.IntPtr(6)
		dsType := tc.DSType("HTTP_LIVE_NATNL")
		ds.Type = &dsType
		ds.OrgServerFQDN = nil // nil origin should not be included
		ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
		ds.CacheURL = util.StrPtr("mycacheurl")
		ds.RangeRequestHandling = util.IntPtr(0)
		ds.RemapText = util.StrPtr("myremaptext")
		ds.EdgeHeaderRewrite = util.StrPtr("myedgeheaderrewrite")
		ds.SigningAlgorithm = util.StrPtr("url_sig")
		ds.XMLID = util.StrPtr("ds6")
		ds.QStringIgnore = util.IntPtr(0)
		ds.RegexRemap = util.StrPtr("myregexremap")
		ds.FQPacingRate = util.IntPtr(0)
		ds.DSCP = util.IntPtr(0)
		ds.RoutingName = util.StrPtr("myroutingname")
		ds.MultiSiteOrigin = util.BoolPtr(false)
		ds.OriginShield = util.StrPtr("myoriginshield")
		ds.ProfileID = util.IntPtr(49)
		ds.ProfileName = util.StrPtr("dsprofile")
		ds.Protocol = util.IntPtr(0)
		ds.AnonymousBlockingEnabled = util.BoolPtr(false)
		ds.Active = util.BoolPtr(true)
		dses = append(dses, ds)
	}
	{
		ds := tc.DeliveryServiceNullableV30{}
		ds.ID = util.IntPtr(7)
		dsType := tc.DSType("HTTP_LIVE_NATNL")
		ds.Type = &dsType
		ds.OrgServerFQDN = util.StrPtr("") // empty origin should not be included
		ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
		ds.CacheURL = util.StrPtr("mycacheurl")
		ds.RangeRequestHandling = util.IntPtr(0)
		ds.RemapText = util.StrPtr("myremaptext")
		ds.EdgeHeaderRewrite = util.StrPtr("myedgeheaderrewrite")
		ds.SigningAlgorithm = util.StrPtr("url_sig")
		ds.XMLID = util.StrPtr("ds7")
		ds.QStringIgnore = util.IntPtr(0)
		ds.RegexRemap = util.StrPtr("myregexremap")
		ds.FQPacingRate = util.IntPtr(0)
		ds.DSCP = util.IntPtr(0)
		ds.RoutingName = util.StrPtr("myroutingname")
		ds.MultiSiteOrigin = util.BoolPtr(false)
		ds.OriginShield = util.StrPtr("myoriginshield")
		ds.ProfileID = util.IntPtr(49)
		ds.ProfileName = util.StrPtr("dsprofile")
		ds.Protocol = util.IntPtr(0)
		ds.AnonymousBlockingEnabled = util.BoolPtr(false)
		ds.Active = util.BoolPtr(true)
		dses = append(dses, ds)
	}
	{ // see regexes - nil pattern
		ds := tc.DeliveryServiceNullableV30{}
		ds.ID = util.IntPtr(8)
		dsType := tc.DSType("HTTP_LIVE_NATNL")
		ds.Type = &dsType
		ds.OrgServerFQDN = util.StrPtr("") // empty origin should not be included
		ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
		ds.CacheURL = util.StrPtr("mycacheurl")
		ds.RangeRequestHandling = util.IntPtr(0)
		ds.RemapText = util.StrPtr("myremaptext")
		ds.EdgeHeaderRewrite = util.StrPtr("myedgeheaderrewrite")
		ds.SigningAlgorithm = util.StrPtr("url_sig")
		ds.XMLID = util.StrPtr("ds8")
		ds.QStringIgnore = util.IntPtr(0)
		ds.RegexRemap = util.StrPtr("myregexremap")
		ds.FQPacingRate = util.IntPtr(0)
		ds.DSCP = util.IntPtr(0)
		ds.RoutingName = util.StrPtr("myroutingname")
		ds.MultiSiteOrigin = util.BoolPtr(false)
		ds.OriginShield = util.StrPtr("myoriginshield")
		ds.ProfileID = util.IntPtr(49)
		ds.ProfileName = util.StrPtr("dsprofile")
		ds.Protocol = nil // nil protocol shouldn't be included
		ds.AnonymousBlockingEnabled = util.BoolPtr(false)
		ds.Active = util.BoolPtr(true)
		dses = append(dses, ds)
	}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(1),
		},
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(2),
		},
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(3),
		},
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(4),
		},
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(5),
		},
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(6),
		},
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(7),
		},
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(8),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: "ds",
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypePathRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern",
				},
			},
		},
		tc.DeliveryServiceRegexes{
			DSName: "ds2",
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeSteeringRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern2",
				},
			},
		},
		tc.DeliveryServiceRegexes{
			DSName: "ds3",
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHeaderRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern2",
				},
			},
		},
		tc.DeliveryServiceRegexes{
			DSName: "ds4",
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      "",
					SetNumber: 0,
					Pattern:   "myregexpattern2",
				},
			},
		},
		tc.DeliveryServiceRegexes{
			DSName: "ds5",
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      "nonexistenttype",
					SetNumber: 0,
					Pattern:   "myregexpattern2",
				},
			},
		},
		tc.DeliveryServiceRegexes{
			DSName: "ds6",
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern2",
				},
			},
		},
		tc.DeliveryServiceRegexes{
			DSName: "ds7",
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern2",
				},
			},
		},
		tc.DeliveryServiceRegexes{
			DSName: "ds8",
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "",
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected no remaps from DSes with missing data, actual: '%v' count %v", txt, len(txtLines))
	}

}

func TestMakeRemapDotConfigEdgeHostRegexReplacement(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPAndHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `.*\.mypattern\..*`, // common host regex syntax, should be replaced
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 {
		t.Fatalf("expected 3 remaps from HTTP_AND_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	if strings.Count(txt, "mypattern") != 2 {
		t.Errorf("expected 2 pattern occurences from HTTP_AND_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(txt, `\`) {
		t.Errorf(`expected regex pattern '\' to be replaced and not exist in remap, actual: '%v'`, txtLines)
	}

	if strings.Contains(txt, `.*`) {
		t.Errorf(`expected regex pattern '.*' to be replaced and not exist in remap, actual: '%v'`, txtLines)
	}
}

func TestMakeRemapDotConfigEdgeHostRegexReplacementHTTP(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTP))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `.*\.mypattern\..*`, // common host regex syntax, should be replaced
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remap from HTTP DS, actual: '%v' count %v", txt, len(txtLines))
	}

	if strings.Count(txt, "mypattern") != 1 {
		t.Errorf("expected 1 pattern occurences from HTTP DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(txt, `\`) {
		t.Errorf(`expected regex pattern '\' to be replaced and not exist in remap, actual: '%v'`, txtLines)
	}

	if strings.Contains(txt, `.*`) {
		t.Errorf(`expected regex pattern '.*' to be replaced and not exist in remap, actual: '%v'`, txtLines)
	}
}

func TestMakeRemapDotConfigEdgeHostRegexReplacementHTTPS(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `.*\.mypattern\..*`, // common host regex syntax, should be replaced
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	if strings.Count(txt, "mypattern") != 1 {
		t.Errorf("expected 1 pattern occurences from HTTP DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(txt, `\`) {
		t.Errorf(`expected regex pattern '\' to be replaced and not exist in remap, actual: '%v'`, txtLines)
	}

	if strings.Contains(txt, `.*`) {
		t.Errorf(`expected regex pattern '.*' to be replaced and not exist in remap, actual: '%v'`, txtLines)
	}
}

func TestMakeRemapDotConfigEdgeHostRegexReplacementHTTPToHTTPS(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `.*\.mypattern\..*`, // common host regex syntax, should be replaced
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	if strings.Count(txt, "mypattern") != 1 {
		t.Errorf("expected 1 pattern occurences from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(txt, `\`) {
		t.Errorf(`expected regex pattern '\' to be replaced and not exist in remap, actual: '%v'`, txtLines)
	}

	if strings.Contains(txt, `.*`) {
		t.Errorf(`expected regex pattern '.*' to be replaced and not exist in remap, actual: '%v'`, txtLines)
	}
}

func TestMakeRemapDotConfigEdgeRemapUnderscoreHTTPReplace(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "__http__") {
		t.Errorf("expected literal pattern to replace '__http__', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "myliteralpattern"+*server.HostName+"foo") {
		t.Errorf("expected literal pattern to replace __http__ with server name, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeDSCPRemap(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "set_dscp_") {
		t.Errorf("expected remap with dscp_remap parameter to not have set_dscp header rewrite, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "dscp_remap") {
		t.Errorf("expected remap with dscp_remap parameter to have dscp_remap text, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeNoDSCPRemap(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "set_dscp_") {
		t.Errorf("expected remap with no dscp_remap parameter to have set_dscp header rewrite, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "dscp_remap") {
		t.Errorf("expected remap with no dscp_remap parameter to not have dscp_remap text, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeHeaderRewrite(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = util.StrPtr("myedgeheaderrewrite")
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "hdr_rw_") {
		t.Errorf("expected remap on edge server with edge header rewrite to contain rewrite file, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "mymidrewrite") {
		t.Errorf("expected remap on edge server to not contain mid rewrite, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "hdr_rw_mid") {
		t.Errorf("expected remap on edge server to not contain mid rewrite file, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeHeaderRewriteEmpty(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = util.StrPtr("")
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "hdr_rw_") {
		t.Errorf("expected remap on edge server with empty edge header rewrite to not contain rewrite file, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "mymidrewrite") {
		t.Errorf("expected remap on edge server to not contain mid rewrite, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "hdr_rw_mid") {
		t.Errorf("expected remap on edge server to not contain mid rewrite file, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeHeaderRewriteNil(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "hdr_rw_") {
		t.Errorf("expected remap on edge server with nil edge header rewrite to not contain rewrite file, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "mymidrewrite") {
		t.Errorf("expected remap on edge server to not contain mid rewrite, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "hdr_rw_mid") {
		t.Errorf("expected remap on edge server to not contain mid rewrite file, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeSigningURLSig(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("url_sig")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "url_sig_") {
		t.Errorf("expected remap on edge server with URL Sig to contain url sig file, actual '%v'", txt)
	}
	if strings.Contains(remapLine, "uri_signing") {
		t.Errorf("expected remap on edge server with URL Sig to not contain uri signing file, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeSigningURISigning(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("uri_signing")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "url_sig_") {
		t.Errorf("expected remap on edge server with URL Sig to not contain url sig file, actual '%v'", txt)
	}
	if !strings.Contains(remapLine, "uri_signing") {
		t.Errorf("expected remap on edge server with URL Sig to contain uri signing file, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeSigningNone(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = nil
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "url_sig") {
		t.Errorf("expected remap on edge server with nil signing to not contain url sig file, actual '%v'", txt)
	}
	if strings.Contains(remapLine, "uri_signing") {
		t.Errorf("expected remap on edge server with nil signing to not contain uri signing file, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeSigningEmpty(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "url_sig") {
		t.Errorf("expected remap on edge server with empty signing to not contain url sig file, actual '%v'", txt)
	}
	if strings.Contains(remapLine, "uri_signing") {
		t.Errorf("expected remap on edge server with empty signing to not contain uri signing file, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeSigningWrong(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "url_sig") {
		t.Errorf("expected remap on edge server with wrong signing to not contain url sig file, actual '%v'", txt)
	}
	if strings.Contains(remapLine, "uri_signing") {
		t.Errorf("expected remap on edge server with wrong signing to not contain uri signing file, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeQStringDropAtEdge(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreDropAtEdge))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "drop_qstring.config") {
		t.Errorf("expected remap on edge server with qstring drop at edge to contain drop qstring config, actual '%v'", txt)
	}

}

func TestMakeRemapDotConfigEdgeQStringIgnorePassUp(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cachekey.so") {
		t.Errorf("expected remap on edge server with qstring ignore pass up to contain cachekey plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeQStringIgnorePassUpWithCacheKeyParameter(t *testing.T) {
	// ATS doesn't allow multiple inclusions of the same plugin.
	// Currently, if there's both a QString cachekey inclusion, and a cache key parameter,
	// the make func adds both, and logs a warning.

	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cachekey.so") {
		t.Errorf("expected remap on edge server with qstring ignore pass up to contain cachekey plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cachekeykey") {
		t.Errorf("expected remap on edge server with qstring ignore pass up and cachekey param to include both, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeQStringIgnorePassUpCacheURLParam(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cacheKeyParams := []tc.Parameter{}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "cachekey.so") {
		t.Errorf("expected remap on edge server with qstring ignore pass up but also cacheurl parameter to not contain cachekey plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeQStringIgnorePassUpCacheURLParamCacheURL(t *testing.T) {

	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "5",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "cachekey.so") {
		t.Errorf("expected remap on edge server with ats<5 to not contain cachekey plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cacheurl.so") {
		t.Errorf("expected remap on edge server with ats<5 to contain cacheurl  plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeQStringIgnorePassUpCacheURLParamCacheURLAndDSCacheURL(t *testing.T) {
	// Currently, the make func should log an error if the QString results in a cacheurl plugin, and there's also a cacheurl, but it should generate it anyway.

	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "5",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "cachekey.so") {
		t.Errorf("expected remap on edge server with ats<5 to not contain cachekey plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cacheurl.so") {
		t.Errorf("expected remap on edge server with ats<5 to contain cacheurl  plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cacheurl_") {
		t.Errorf("expected remap on edge server with ds qstring cacheurl and ds cacheurl to generate both, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigMidQStringIgnorePassUpCacheURLParamCacheURLAndDSCacheURL(t *testing.T) {

	// Currently, the make func should log an error if the QString results in a cacheurl plugin, and there's also a cacheurl, but it should generate it anyway.

	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "5",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "cachekey.so") {
		t.Errorf("expected remap on edge server with ats<5 to not contain cachekey plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cacheurl.so") {
		t.Errorf("expected remap on edge server with ats<5 to contain cacheurl  plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cacheurl_") {
		t.Errorf("expected remap on edge server with ds qstring cacheurl and ds cacheurl to generate both, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeCacheURL(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cacheurl_") {
		t.Errorf("expected remap on edge server with ds cacheurl to contain cacheurl plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeCacheKeyParams(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cachekey.so") {
		t.Errorf("expected remap on edge server with ds cache key params to contain cachekey plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cachekeykey") {
		t.Errorf("expected remap on edge server with ds cache key params to contain param keys, actual '%v'", txt)
	}
	if !strings.Contains(remapLine, "cachekeyval") {
		t.Errorf("expected remap on edge server with ds cache key params to contain param vals, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "shouldnotinclude") {
		t.Errorf("expected remap on edge server to not include different ds cache key params, actual '%v'", txt)
	}

}

func TestMakeRemapDotConfigEdgeRegexRemap(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "regex_remap_") {
		t.Errorf("expected remap on edge server with ds regex remap to contain regex remap file, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "myregexremap") {
		t.Errorf("expected remap on edge server with ds regex remap to contain regex remap file, but not actual text, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeRegexRemapEmpty(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest))
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "regex_remap_") {
		t.Errorf("expected remap on edge server with empty ds regex remap to not contain regex remap file, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeRangeRequestNil(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = nil
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "background_fetch.so") {
		t.Errorf("expected remap on edge server with ds nil range request handling to not contain background fetch plugin, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "cache_range_requests.so") {
		t.Errorf("expected remap on edge server with ds nil range request handling to not contain cache_range_requests plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeRangeRequestDontCache(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingDontCache)
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "background_fetch.so") {
		t.Errorf("expected remap on edge server with ds dont-cache range request handling to not contain background fetch plugin, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "cache_range_requests.so") {
		t.Errorf("expected remap on edge server with ds dont-cache range request handling to not contain cache_range_requests plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeRangeRequestBGFetch(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingBackgroundFetch)
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "background_fetch.so") {
		t.Errorf("expected remap on edge server with ds bg-fetch range request handling to contain background fetch plugin, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "cache_range_requests.so") {
		t.Errorf("expected remap on edge server with ds bg-fetch range request handling to not contain cache_range_requests plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeRangeRequestSlice(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingSlice)
	ds.RemapText = util.StrPtr("myremaptext")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)
	ds.RangeSliceBlockSize = util.IntPtr(262144)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "slice.so") {
		t.Errorf("expected remap on edge server with ds slice range request handling to contain background fetch plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cache_range_requests.so") {
		t.Errorf("expected remap on edge server with ds slice range request handling to contain cache_range_requests plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "pparam=--blockbytes=262144") {
		t.Errorf("expected remap on edge server with ds slice range request handling to contain block size for the slice plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigRawRemapRangeDirective(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingSlice)
	ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua __RANGE_DIRECTIVE__")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPAndHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)
	ds.RangeSliceBlockSize = util.IntPtr(262144)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // 2 remaps plus header comment
		t.Fatalf("expected 2 remaps from HTTP_AND_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "slice.so") {
		t.Errorf("expected remap on edge server with ds slice range request handling to contain background fetch plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cache_range_requests.so") {
		t.Errorf("expected remap on edge server with ds slice range request handling to contain cache_range_requests plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "pparam=--blockbytes=262144") {
		t.Errorf("expected remap on edge server with ds slice range request handling to contain block size for the slice plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "@plugin=tslua.so @pparam=my-range-manipulator.lua  @plugin=slice.so @pparam=--blockbytes=262144 @plugin=cache_range_requests.so") {
		t.Errorf("expected raw remap to come after range directive, actual '%v'", txt)
	}
	if strings.Contains(remapLine, "__RANGE_DIRECTIVE__") {
		t.Errorf("expected raw remap range directive to be replaced, actual '%v'", txt)
	}
	if count := strings.Count(remapLine, "slice.so"); count != 1 { // Individual line should only have 1 slice.so
		t.Errorf("expected raw remap range directive to be replaced not duplicated, actual count %v '%v'", count, txt)
	}
	if count := strings.Count(txt, "slice.so"); count != 2 { // All lines should have 2 slice.so - HTTP and HTTPS lines
		t.Errorf("expected raw remap range directive to have one slice.so for HTTP and one for HTTPS remap, actual count %v '%v'", count, txt)
	}
}

func TestMakeRemapDotConfigRawRemapWithoutRangeDirective(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingSlice)
	ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)
	ds.RangeSliceBlockSize = util.IntPtr(262144)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "slice.so") {
		t.Errorf("expected remap on edge server with ds slice range request handling to contain background fetch plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cache_range_requests.so") {
		t.Errorf("expected remap on edge server with ds slice range request handling to contain cache_range_requests plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "pparam=--blockbytes=262144") {
		t.Errorf("expected remap on edge server with ds slice range request handling to contain block size for the slice plugin, actual '%v'", txt)
	}

	if !strings.HasSuffix(remapLine, "@plugin=tslua.so @pparam=my-range-manipulator.lua") {
		t.Errorf("expected raw remap without range directive at end of remap line, actual '%v'", txt)
	}
	if strings.Count(remapLine, "slice.so") != 1 {
		t.Errorf("expected raw remap range directive to not be duplicated, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeRangeRequestCache(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest)
	ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "background_fetch.so") {
		t.Errorf("expected remap on edge server with ds cache range request handling to not contain background fetch plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cache_range_requests.so") {
		t.Errorf("expected remap on edge server with ds cache range request handling to contain cache_range_requests plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeFQPacingNil(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest)
	ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("")
	ds.FQPacingRate = nil
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "fq_pacing.so") {
		t.Errorf("expected remap on edge server with ds nil fq pacing to not contain fq_pacing plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeFQPacingNegative(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest)
	ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("")
	ds.FQPacingRate = util.IntPtr(-42)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "fq_pacing.so") {
		t.Errorf("expected remap on edge server with ds negative fq pacing to not contain fq_pacing plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeFQPacingZero(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest)
	ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "fq_pacing.so") {
		t.Errorf("expected remap on edge server with ds zero fq pacing to not contain fq_pacing plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeFQPacingPositive(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest)
	ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("")
	ds.FQPacingRate = util.IntPtr(314159)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `myliteralpattern__http__foo`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "fq_pacing.so") {
		t.Errorf("expected remap on edge server with ds positive fq pacing to contain fq_pacing plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "314159") {
		t.Errorf("expected remap on edge server with ds positive fq pacing to contain fq_pacing number, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeDNS(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("DNS_LIVE")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest)
	ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("")
	ds.FQPacingRate = util.IntPtr(314159)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `.*\.mypattern\..*`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[1]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "myroutingname") {
		t.Errorf("expected remap on edge server with ds dns to contain routing name, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeDNSNoRoutingName(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("DNS_LIVE")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest)
	ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("")
	ds.FQPacingRate = util.IntPtr(314159)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = nil
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `.*\.mypattern\..*`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected no remaps from DNS DS with nil routing name, actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigEdgeRegexTypeNil(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	ds := tc.DeliveryServiceNullableV30{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("DNS_LIVE")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.CacheURL = util.StrPtr("mycacheurl")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest)
	ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua")
	ds.EdgeHeaderRewrite = nil
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RegexRemap = util.StrPtr("")
	ds.FQPacingRate = util.IntPtr(314159)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = nil
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPToHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{ds}

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          util.IntPtr(*server.ID),
			DeliveryService: util.IntPtr(*ds.ID),
		},
	}

	dsRegexes := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      "",
					SetNumber: 0,
					Pattern:   `.*\.mypattern\..*`,
				},
			},
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(*server.Profile),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(*server.Profile),
		},
	}

	cacheKeyParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "shouldnotexist",
			ConfigFile: "cacheurl.config",
			Value:      "shouldnotexisteither",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cacheurl.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cacheurl.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	txt := MakeRemapDotConfig(server, dses, dss, dsRegexes, serverParams, cdn, toToolName, toURL, cacheKeyParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, *server.HostName, toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected no remaps for DS with nil regex type, actual: '%v' count %v", txt, len(txtLines))
	}

}

func makeTestRemapServer() *tc.ServerNullable {
	server := &tc.ServerNullable{}
	server.ProfileID = util.IntPtr(42)
	server.CDNName = util.StrPtr("mycdn")
	server.Cachegroup = util.StrPtr("cg0")
	server.DomainName = util.StrPtr("mydomain")
	server.CDNID = util.IntPtr(43)
	server.HostName = util.StrPtr("server0")
	server.HTTPSPort = util.IntPtr(12443)
	server.ID = util.IntPtr(44)
	setIP(server, "192.168.2.4")
	server.ProfileID = util.IntPtr(46)
	server.Profile = util.StrPtr("MyProfile")
	server.TCPPort = util.IntPtr(12080)
	server.Type = "MID"
	return server
}
