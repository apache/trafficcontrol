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
	"bufio"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func makeTestRemapServer() *Server {
	server := &Server{}
	server.CDNName = util.StrPtr("mycdn")
	server.Cachegroup = util.StrPtr("cg0")
	server.DomainName = util.StrPtr("mydomain")
	server.CDNID = util.IntPtr(43)
	server.HostName = util.StrPtr("server")
	server.HTTPSPort = util.IntPtr(12345)
	server.ID = util.IntPtr(44)
	setIP(server, "192.168.2.4")
	server.ProfileNames = []string{"MyProfile"}
	server.TCPPort = util.IntPtr(1280)
	server.Type = "MID"
	return server
}

func makeTestAnyCastServers() []Server {
	server1 := makeTestRemapServer()
	server1.Type = "EDGE"
	server1.HostName = util.StrPtr("mcastserver1")
	server1.ID = util.IntPtr(45)
	server1.Interfaces = []tc.ServerInterfaceInfoV40{}
	setIPInfo(server1, "lo", "192.168.2.6", "fdf8:f53b:82e4::53")

	server2 := makeTestRemapServer()
	server2.Type = "EDGE"
	server2.HostName = util.StrPtr("mcastserver2")
	server2.ID = util.IntPtr(46)
	server2.Interfaces = []tc.ServerInterfaceInfoV40{}
	setIPInfo(server2, "lo", "192.168.2.6", "fdf8:f53b:82e4::53")

	return []Server{*server1, *server2}
}

// tokenize remap line
func tokenize(txt string) []string {
	tokens := []string{}
	scanner := bufio.NewScanner(strings.NewReader(txt))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		tokens = append(tokens, scanner.Text())
	}
	return tokens
}

// Returns list of plugins from given remap rule.
func pluginsFromTokens(tokens []string, prefix string) []string {
	plugins := []string{}

	for _, token := range tokens {
		if strings.HasPrefix(token, prefix) {
			token = strings.TrimPrefix(token, prefix)
			if 0 < len(token) {
				plugins = append(plugins, token)
			}
		}
	}

	return plugins
}

func TestAnyCastRemapDotConfig(t *testing.T) {
	hdr := "myHeaderComment"
	mappings := map[string]bool{
		"http://dnsroutingname.mypattern1": false,
		"http://myregexpattern1":          false,
		"http://server.mypattern0":        false,
		"https://server.mypattern0":       false,
		"http://mcastserver1.mypattern0":  false,
		"https://mcastserver1.mypattern0": false,
		"http://mcastserver2.mypattern0":  false,
		"https://mcastserver2.mypattern0": false,
		"http://myregexpattern0":          false,
		"https://myregexpattern0":         false,
	}
	server := makeTestRemapServer()
	server.Type = "EDGE"
	server.Interfaces = []tc.ServerInterfaceInfoV40{}
	setIPInfo(server, "lo", "192.168.2.6", "fdf8:f53b:82e4::53")
	servers := makeTestAnyCastServers()
	for _, anyCsstServer := range getAnyCastPartners(server, servers) {
		if len(anyCsstServer) != 2 {
			t.Errorf("expected 2 edges in anycast group, actual '%v'", len(anyCsstServer))
		}
	}

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
	ds.RangeRequestHandling = util.IntPtr(0)
	ds.RemapText = nil
	ds.EdgeHeaderRewrite = util.StrPtr("myedgeheaderrewrite")
	ds.SigningAlgorithm = nil
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(0)
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.FQPacingRate = util.IntPtr(0)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPAndHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	ds1 := DeliveryService{}
	ds1.ID = util.IntPtr(49)
	dsType1 := tc.DSType("DNS")
	ds1.Type = &dsType1
	ds1.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds1.RangeRequestHandling = util.IntPtr(0)
	ds1.RemapText = nil
	ds1.SigningAlgorithm = nil
	ds1.XMLID = util.StrPtr("mydsname1")
	ds1.QStringIgnore = util.IntPtr(0)
	ds1.RegexRemap = util.StrPtr("")
	ds1.FQPacingRate = util.IntPtr(0)
	ds1.DSCP = util.IntPtr(0)
	ds1.RoutingName = util.StrPtr("dnsroutingname")
	ds1.MultiSiteOrigin = util.BoolPtr(false)
	ds1.OriginShield = util.StrPtr("myoriginshield")
	ds1.ProfileID = util.IntPtr(49)
	ds1.ProfileName = util.StrPtr("dsprofile")
	ds1.Protocol = util.IntPtr(int(tc.DSProtocolHTTP))
	ds1.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds1.Active = util.BoolPtr(true)

	dses := []DeliveryService{ds, ds1}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
		},
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds1.ID,
		},
	}
	for _, srv := range servers {
		for _, ds := range dses {
			dssrv := DeliveryServiceServer{
				Server:          *srv.ID,
				DeliveryService: *ds.ID,
			}
			dss = append(dss, dssrv)
		}
	}
	dsRegexes := []tc.DeliveryServiceRegexes{}
	for i, ds := range dses {
		pattern := fmt.Sprintf(`.*\.mypattern%d\..*`, i)
		customRex := fmt.Sprintf("myregexpattern%d", i)
		dsr := tc.DeliveryServiceRegexes{
			DSName: *ds.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   pattern,
				},
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 1,
					Pattern:   customRex,
				},
			},
		}
		dsRegexes = append(dsRegexes, dsr)
	}
	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
	}
	remapConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekey.pparam",
			ConfigFile: "remap.config",
			Value:      "--cachekeykey=cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	//t.Logf("text: %v", txt)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")
	for _, line := range txtLines[2:] {
		switch {
		case strings.Contains(line, "http://dnsroutingname.mypattern1"):
			mappings["http://dnsroutingname.mypattern1"] = true
		case strings.Contains(line, "http://myregexpattern1"):
			mappings["http://myregexpattern1"] = true
		case strings.Contains(line, "http://server.mypattern0"):
			mappings["http://server.mypattern0"] = true
		case strings.Contains(line, "https://server.mypattern0"):
			mappings["https://server.mypattern0"] = true
		case strings.Contains(line, "http://mcastserver1.mypattern0"):
			mappings["http://mcastserver1.mypattern0"] = true
		case strings.Contains(line, "https://mcastserver1.mypattern0"):
			mappings["https://mcastserver1.mypattern0"] = true
		case strings.Contains(line, "http://mcastserver2.mypattern0"):
			mappings["http://mcastserver2.mypattern0"] = true
		case strings.Contains(line, "https://mcastserver2.mypattern0"):
			mappings["https://mcastserver2.mypattern0"] = true
		case strings.Contains(line, "http://myregexpattern0"):
			mappings["http://myregexpattern0"] = true
		case strings.Contains(line, "https://myregexpattern0"):
			mappings["https://myregexpattern0"] = true
		default:
			t.Fatalf("unexpected remap line '%v'", line)
		}
	}
	for key, val := range mappings{
		if !val {
			t.Fatalf("expected to find remap rule for '%v'", key)
		}
	}

	/*if len(txtLines) != 14 {
		t.Log(cfg.Warnings)
		t.Fatalf("expected a total of 12 remap lines for anycast hosts, actual: '%v' count %v", txt, len(txtLines))
	}*/
}
func TestMakeRemapDotConfig0(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"

	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
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
	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
	}

	remapConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekey.pparam",
			ConfigFile: "remap.config",
			Value:      "--cachekeykey=cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	t.Logf("text: %v", txt)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 {
		t.Log(cfg.Warnings)
		t.Fatalf("expected one line for each remap plus a comment and blank, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
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
	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected no remap lines for LIVE local DS, actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigMid(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
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
	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 {
		t.Fatalf("expected one line for each remap plus a comment and blank, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = nil
	ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
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
	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected no remap lines for DS with nil Origin FQDN, actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigEmptyOrigin(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("")
	ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
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
	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected no remap lines for DS with empty Origin FQDN, actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigDuplicateOrigins(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
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

	ds2 := DeliveryService{}
	ds2.ID = util.IntPtr(49)
	dsType2 := tc.DSType("HTTP_LIVE_NATNL")
	ds2.Type = &dsType2
	ds2.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds2.MidHeaderRewrite = util.StrPtr("mymidrewrite2")
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

	dses := []DeliveryService{ds, ds2}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
		},
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds2.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 {
		t.Fatalf("expected a comment header, a blank line, and 1 remap lines for multiple DSes with the same Origin (ATS can't handle multiple remaps with the same origin FQDN), actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigNilMidRewrite(t *testing.T) {
	hdr := "myHeaderComment"
	server := makeTestRemapServer()
	servers := makeTestAnyCastServers()
	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = nil
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected one line, actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigMidHasNoEdgeRewrite(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = nil
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected one line, actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigMidProfileCacheKey(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekey.pparam",
			ConfigFile: "remap.config",
			Value:      "--ckeypp=cvalpp",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "ckeycc",
			ConfigFile: "cachekey.config",
			Value:      "cvalcc",
			Profiles:   []byte(`["dsprofile"]`),
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 {
		t.Errorf("expected one line for each remap plus a comment and blank, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Count(remapLine, "origin.example.test") != 2 {
		t.Errorf("expected to contain origin FQDN twice (Mids remap origins to themselves, as a forward proxy), actual '%v'", txt)
	}

	if 1 != strings.Count(remapLine, "cachekey.so") {
		t.Errorf("expected only single cachekey.so plugin, actual '%v'", txt)
	} else if !strings.Contains(remapLine, "--remove-path=true") {
		t.Errorf("expected cachekey qstring ignore args, actual '%v'", txt)
	} else if !strings.Contains(remapLine, "--ckeypp=cvalpp") {
		t.Errorf("expected to contain cachekey.pparam param, actual '%v'", txt)
	} else if !strings.Contains(remapLine, "--ckeycc=cvalcc") {
		t.Errorf("expected to contain cachekey.config param, actual '%v'", txt)
	}

	if !warningsContains(cfg.Warnings, "Both new cachekey.pparam and old cachekey.config parameters assigned") {
		t.Errorf("expected to contain warning about using both cachekey.config and cachekey.pparam, actual '%v', '%v'", cfg.Warnings, txt)
	}
}

func TestMakeRemapDotConfigMidBgFetchHandling(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingBackgroundFetch))
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 {
		t.Errorf("expected one line for each remap plus a comment and blank, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Count(remapLine, "origin.example.test") != 2 {
		t.Errorf("expected to contain origin FQDN twice (Mids remap origins to themselves, as a forward proxy), actual '%v'", txt)
	}

	if strings.Contains(remapLine, "background_fetch.so") {
		t.Errorf("did not expect to contain background_fetch plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigMidRangeRequestHandling(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 {
		t.Errorf("expected one line for each remap plus a comment and blank, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	ds2 := DeliveryService{}
	ds2.ID = util.IntPtr(48)
	dsType2 := tc.DSType("HTTP_LIVE_NATNL")
	ds2.Type = &dsType2
	ds2.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds2.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds, ds2}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 {
		t.Fatalf("expected a comment header, a blank line, and one remap line for DS with origin, but not DS with empty Origin FQDN, actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigAnyMap(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("ANY_MAP")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	ds2 := DeliveryService{}
	ds2.ID = util.IntPtr(49)
	dsType2 := tc.DSType("ANY_MAP")
	ds2.Type = &dsType2
	ds2.OrgServerFQDN = util.StrPtr("myorigin")
	ds2.MidHeaderRewrite = util.StrPtr("mymidrewrite")
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

	dses := []DeliveryService{ds, ds2}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
		},
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds2.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)
	txt = strings.Replace(txt, "\n\n", "\n", -1)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 {
		t.Fatalf("expected a comment header, a blank line, and one remap line for ANY_MAP DS with remap text, but not ANY_MAP DS with nil remap text, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to any_map to not start with 'map' (should be raw ds.RemapText), actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "myremaptext") {
		t.Errorf("expected to contain ANY_MAP DS remap text, actual '%v'", txt)
	}

}

func TestMakeRemapDotConfigEdgeMissingRemapData(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	dses := []DeliveryService{}
	{ // see regexes - has invalid regex type
		ds := DeliveryService{}
		ds.ID = util.IntPtr(1)
		dsType := tc.DSType("HTTP_LIVE_NATNL")
		ds.Type = &dsType
		ds.OrgServerFQDN = util.StrPtr("myorigin")
		ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
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
		ds := DeliveryService{}
		ds.ID = util.IntPtr(2)
		dsType := tc.DSType("HTTP_LIVE_NATNL")
		ds.Type = &dsType
		ds.OrgServerFQDN = util.StrPtr("myorigin")
		ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
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
		ds := DeliveryService{}
		ds.ID = util.IntPtr(3)
		dsType := tc.DSType("HTTP_LIVE_NATNL")
		ds.Type = &dsType
		ds.OrgServerFQDN = util.StrPtr("myorigin")
		ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
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
		ds := DeliveryService{}
		ds.ID = util.IntPtr(4)
		dsType := tc.DSType("HTTP_LIVE_NATNL")
		ds.Type = &dsType
		ds.OrgServerFQDN = util.StrPtr("myorigin")
		ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
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
		ds := DeliveryService{}
		ds.ID = util.IntPtr(5)
		dsType := tc.DSType("HTTP_LIVE_NATNL")
		ds.Type = &dsType
		ds.OrgServerFQDN = util.StrPtr("myorigin")
		ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
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
		ds := DeliveryService{}
		ds.ID = util.IntPtr(6)
		dsType := tc.DSType("HTTP_LIVE_NATNL")
		ds.Type = &dsType
		ds.OrgServerFQDN = nil // nil origin should not be included
		ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
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
		ds := DeliveryService{}
		ds.ID = util.IntPtr(7)
		dsType := tc.DSType("HTTP_LIVE_NATNL")
		ds.Type = &dsType
		ds.OrgServerFQDN = util.StrPtr("") // empty origin should not be included
		ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
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
		ds := DeliveryService{}
		ds.ID = util.IntPtr(8)
		dsType := tc.DSType("HTTP_LIVE_NATNL")
		ds.Type = &dsType
		ds.OrgServerFQDN = util.StrPtr("") // empty origin should not be included
		ds.MidHeaderRewrite = util.StrPtr("mymidrewrite")
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

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: 1,
		},
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: 2,
		},
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: 3,
		},
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: 4,
		},
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: 5,
		},
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: 6,
		},
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: 7,
		},
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: 8,
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
					Pattern:   "myregexpattern3",
				},
			},
		},
		tc.DeliveryServiceRegexes{
			DSName: "ds4",
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      "",
					SetNumber: 0,
					Pattern:   "myregexpattern4",
				},
			},
		},
		tc.DeliveryServiceRegexes{
			DSName: "ds5",
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      "nonexistenttype",
					SetNumber: 0,
					Pattern:   "myregexpattern5",
				},
			},
		},
		tc.DeliveryServiceRegexes{
			DSName: "ds6",
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern6",
				},
			},
		},
		tc.DeliveryServiceRegexes{
			DSName: "ds7",
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   "myregexpattern7",
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected no remaps from DSes with missing data, actual: '%v' count %v", txt, len(txtLines))
	}

}

func TestMakeRemapDotConfigEdgeHostRegexReplacement(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 4 {
		t.Fatalf("expected a comment header, a blank line, and 2 remaps from HTTP_AND_HTTPS DS, actual: '%v' line count %v", txt, len(txtLines))
	}

	if strings.Count(txt, "mypattern") != 2 {
		t.Errorf("expected 2 pattern occurences from HTTP_AND_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 {
		t.Fatalf("expected a comment header, a blank line, and 1 remap from HTTP DS, actual: '%v' count %v", txt, len(txtLines))
	}

	if strings.Count(txt, "mypattern") != 1 {
		t.Errorf("expected 1 pattern occurences from HTTP DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	if strings.Count(txt, "mypattern") != 1 {
		t.Errorf("expected 1 pattern occurences from HTTP DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	if strings.Count(txt, "mypattern") != 1 {
		t.Errorf("expected 1 pattern occurences from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "url_sig.pparam",
			ConfigFile: "remap.config",
			Value:      "pristine",
			Profiles:   []byte(`["dsprofile"]`),
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if 1 != strings.Count(remapLine, "url_sig.so") {
		t.Errorf("expected remap on edge server with URL Sig to contain url_sig.so, actual '%v'", txt)
	} else if !strings.Contains(remapLine, "url_sig_") {
		t.Errorf("expected remap on edge server with URL Sig to contain url sig file, actual '%v'", txt)
	} else if !strings.Contains(remapLine, "pristine") {
		t.Errorf("expected remap on edge server with URL Sig to contain pristine arg, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "uri_signing") {
		t.Errorf("expected remap on edge server with URL Sig to not contain uri signing file, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeSigningURISigning(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "drop_qstring.config") {
		t.Errorf("expected remap on edge server with qstring drop at edge to contain drop qstring config, actual '%v'", txt)
	}

}

func TestMakeRemapDotConfigEdgeQStringIgnorePassUp(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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

	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekey.pparam",
			ConfigFile: "remap.config",
			Value:      "--ckeypp=",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "ckeycc",
			ConfigFile: "cachekey.config",
			Value:      "",
			Profiles:   []byte(`["dsprofile"]`),
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cachekey.so") {
		t.Errorf("expected remap on edge server with qstring ignore pass up to contain cachekey plugin, actual '%v'", txt)
	}

	if 1 != strings.Count(remapLine, "cachekey.so") {
		t.Errorf("expected only single cachekey.so plugin, actual '%v'", txt)
	} else if !strings.Contains(remapLine, "--remove-all-params=true") {
		t.Errorf("expected remap on edge server with qstring ignore to have cachekey parameters, actual '%v'", txt)
	} else if !strings.Contains(remapLine, "--ckeypp=") {
		t.Errorf("expected cachekey plugin to have '--ckeypp=', actual '%v'", txt)
	} else if !strings.Contains(remapLine, "--ckeycc=") {
		t.Errorf("expected cachekey plugin to have '--ckeycc=', actual '%v'", txt)
	}

	if !warningsContains(cfg.Warnings, "Both new cachekey.pparam and old cachekey.config parameters assigned") {
		t.Errorf("expected to contain warning about using both cachekey.config and cachekey.pparam, actual '%v', '%v'", cfg.Warnings, txt)
	}
}

func TestMakeRemapDotConfigEdgeQStringIgnorePassUpCacheURLParamCacheURL(t *testing.T) {

	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{}
	cgs := []tc.CacheGroupNullable{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "cachekey.so") {
		t.Errorf("expected remap on edge server with ats<5 to not contain cachekey plugin, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "cacheurl.so") {
		t.Errorf("expected remap on edge server with ats<5 to not contain cacheurl plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeCacheKeyParams(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cachekey.pparam",
			ConfigFile: "remap.config",
			Value:      "--ckeypp=cvalpp",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "ckeycc",
			ConfigFile: "cachekey.config",
			Value:      "cvalcc",
			Profiles:   []byte(`["dsprofile"]`),
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cachekey.so") {
		t.Errorf("expected remap on edge server with ds cache key params to contain cachekey plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "--ckeypp=cvalpp") {
		t.Errorf("expected remap on edge server with ds cache key params to contain cachekey.pparam keys, actual '%v'", txt)
	}
	if !strings.Contains(remapLine, "--ckeycc=cvalcc") {
		t.Errorf("expected remap on edge server with ds cache key params to contain cachekey.config vals, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeRegexRemap(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "regex_remap_") {
		t.Errorf("expected remap on edge server with empty ds regex remap to not contain regex remap file, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeRangeRequestNil(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
		tc.Parameter{
			ConfigFile: "cacheurl.config",
			Profiles:   []byte(`["not-dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cachekeykey",
			ConfigFile: "cachekey.config",
			Value:      "cachekeyval",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["dsprofile"]`),
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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

func TestMakeRemapDotConfigEdgeRangeRequestBGFetch(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "background_fetch.pparam",
			ConfigFile: "remap.config",
			Value:      "--log=regex_revalidate.log",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "background_fetch.pparam",
			ConfigFile: "remap.config",
			Value:      "--log=regex_revalidate.log",
			Profiles:   []byte(`["dsprofile"]`),
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}

	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "background_fetch.so") {
		t.Errorf("expected remap on edge server with ds bg-fetch range request handling to contain background fetch plugin, actual '%v'", txt)
	} else if !strings.Contains(remapLine, "@pparam=--log=regex_revalidate.log") {
		t.Errorf("expected remap on edge server to contain background fetch parameter for log, actual '%v'", txt)
	} else if 2 != strings.Count(remapLine, "@pparam=--log=regex_revalidate.log") {
		t.Errorf("expected remap on edge server to contain repeated background fetch parameter for log, actual '%v'", txt)
	}

	if !warningsContains(cfg.Warnings, "Multiple repeated arguments") {
		t.Errorf("expected multiple releated arguments warning, actual '%v'", cfg.Warnings)
	}

	if strings.Contains(remapLine, "cache_range_requests.so") {
		t.Errorf("expected remap on edge server with ds bg-fetch range request handling to not contain cache_range_requests plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeRangeRequestSlice(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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

func TestMakeRemapDotConfigMidRangeRequestSlicePparam(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "MID"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cache_range_requests.pparam",
			ConfigFile: "remap.config",
			Value:      "--consider-ims",
			Profiles:   []byte(`["dsprofile"]`),
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "slice.so") {
		t.Errorf("did not expected remap on mid server with ds slice range request handling to contain slice plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cache_range_requests.so") {
		t.Errorf("expected remap on mid server with ds slice range request handling to contain cache_range_requests plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "--consider-ims") {
		t.Errorf("expected remap on mid server with ds slice range request handling to contain cache_range_requests plugin arg --consider-ims, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "pparam=--blockbytes") {
		t.Errorf("did not expected remap on edge server with ds slice range request handling to contain block size for the slice plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeRangeRequestSlicePparam(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       "cache_range_requests.pparam",
			ConfigFile: "remap.config",
			Value:      "--consider-ims",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "cache_range_requests.pparam",
			ConfigFile: "remap.config",
			Value:      "--no-modify-cachekey",
			Profiles:   []byte(`["dsprofile"]`),
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if 1 != strings.Count(remapLine, "cachekey.so") {
		t.Errorf("expected remap on edge server to contain a single cachekey plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "slice.so") {
		t.Errorf("expected remap on edge server with ds slice range request handling to contain background fetch plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cache_range_requests.so") {
		t.Errorf("expected remap on edge server with ds slice range request handling to contain cache_range_requests plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "--consider-ims") {
		t.Errorf("expected remap on edge server with ds slice range request handling to contain cache_range_requests plugin arg --consider-ims, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "--no-modify-cachekey") {
		t.Errorf("expected remap on edge server with ds slice range request handling to contain cache_range_requests plugin arg --no-modify-cachekey, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "pparam=--blockbytes=262144") {
		t.Errorf("expected remap on edge server with ds slice range request handling to contain block size for the slice plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigRawRemapRangeDirective(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 4 { // 2 remaps plus header comment plus blank
		t.Fatalf("expected a comment header, a blank line, and 2 remaps from HTTP_AND_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if 1 != strings.Count(remapLine, "cachekey.so") {
		t.Errorf("expected remap on edge server to contain a single cachekey plugin, actual '%v'", txt)
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

	if !strings.Contains(remapLine, "@plugin=tslua.so @pparam=my-range-manipulator.lua @plugin=slice.so @pparam=--blockbytes=262144 @plugin=cache_range_requests.so") {
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

func TestMakeRemapDotConfigRawRemapCachekeyDirective(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=uri-manipulator.lua __CACHEKEY_DIRECTIVE__")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 4 { // 2 remaps plus header comment plus blank
		t.Fatalf("expected a comment header, a blank line, and 2 remaps from HTTP_AND_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if 1 != strings.Count(remapLine, "cachekey.so") {
		t.Errorf("expected remap on edge server to contain a single cachekey plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "@plugin=tslua.so @pparam=uri-manipulator.lua @plugin=cachekey.so") {
		t.Errorf("expected cachekey to come after tslua, actual '%v'", txt)
	}
	if strings.Contains(remapLine, "__CACHEKEY_DIRECTIVE__") {
		t.Errorf("expected raw remap cachekey directive to be replaced, actual '%v'", txt)
	}
	if count := strings.Count(remapLine, "cachekey.so"); count != 1 { // Individual line should only have 1 slice.so
		t.Errorf("expected raw remap cachekey directive to be replaced not duplicated, actual count %v '%v'", count, txt)
	}
	if count := strings.Count(txt, "cachekey.so"); count != 2 { // All lines should have 2 slice.so - HTTP and HTTPS lines
		t.Errorf("expected raw remap range directive to have one cachekey.so for HTTP and one for HTTPS remap, actual count %v '%v'", count, txt)
	}
}

func TestMakeRemapDotConfigRawRemapRegexRemapDirective(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=uri-manipulator.lua __REGEX_REMAP_DIRECTIVE__")
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
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTPAndHTTPS))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)
	ds.RangeSliceBlockSize = util.IntPtr(262144)

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 4 { // 2 remaps plus header comment plus blank
		t.Fatalf("expected a comment header, a blank line, and 2 remaps from HTTP_AND_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if 1 != strings.Count(remapLine, "regex_remap.so") {
		t.Errorf("expected remap on edge server to contain a single regex_remap plugin, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "@plugin=tslua.so @pparam=uri-manipulator.lua @plugin=regex_remap.so") {
		t.Errorf("expected regex_remap to come after tslua, actual '%v'", txt)
	}
	if strings.Contains(remapLine, "__REGEX_REMAP_DIRECTIVE__") {
		t.Errorf("expected raw remap regex_remap directive to be replaced, actual '%v'", txt)
	}
	if count := strings.Count(remapLine, "regex_remap.so"); count != 1 { // Individual line should only have 1 regex_remap.so
		t.Errorf("expected raw remap regex_remap directive to be replaced not duplicated, actual count %v '%v'", count, txt)
	}
	if count := strings.Count(txt, "regex_remap.so"); count != 2 { // All lines should have 2 regex_remap.so - HTTP and HTTPS lines
		t.Errorf("expected raw remap regex_remap directive to have one regex_remap.so for HTTP and one for HTTPS remap, actual count %v '%v'", count, txt)
	}
}

func TestMakeRemapDotConfigRawRemapWithoutRangeDirective(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if 1 != strings.Count(remapLine, "cachekey.so") {
		t.Errorf("expected remap on edge server to contain a single cachekey plugin, actual '%v'", txt)
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

	if !strings.HasSuffix(remapLine, "@plugin=tslua.so @pparam=my-range-manipulator.lua # ds 'mydsname' topology ''") {
		t.Errorf("expected raw remap without range directive at end of remap line, actual '%v'", txt)
	}
	if strings.Count(remapLine, "slice.so") != 1 {
		t.Errorf("expected raw remap range directive to not be duplicated, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeRangeRequestCache(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if 1 != strings.Count(remapLine, "cachekey.so") {
		t.Errorf("expected remap on edge server to contain a single cachekey plugin, actual '%v'", txt)
	}

	if strings.Contains(remapLine, "background_fetch.so") {
		t.Errorf("expected remap on edge server with ds cache range request handling to not contain background fetch plugin, actual '%v'", txt)
	}

	if 1 != strings.Count(remapLine, "cache_range_requests.so") {
		t.Errorf("expected remap on edge server with ds cache range request handling to contain cache_range_requests plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeFQPacingNil(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "fq_pacing.so") {
		t.Errorf("expected remap on edge server with ds nil fq pacing to not contain fq_pacing plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeFQPacingNegative(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "fq_pacing.so") {
		t.Errorf("expected remap on edge server with ds negative fq pacing to not contain fq_pacing plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeFQPacingZero(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "fq_pacing.so") {
		t.Errorf("expected remap on edge server with ds zero fq pacing to not contain fq_pacing plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeFQPacingPositive(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

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
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("DNS_LIVE")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "myroutingname") {
		t.Errorf("expected remap on edge server with ds dns to contain routing name, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeDNSNoRoutingName(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("DNS_LIVE")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected no remaps from DNS DS with nil routing name, actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigEdgeRegexTypeNil(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("DNS_LIVE")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected no remaps for DS with nil regex type, actual: '%v' count %v", txt, len(txtLines))
	}

}

func TestMakeRemapDotConfigNoHeaderRewrite(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("DNS_LIVE")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest)
	ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua")
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

	// non-nil default values should not trigger header rewrite plugin directive
	ds.EdgeHeaderRewrite = util.StrPtr("")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.ServiceCategory = util.StrPtr("")
	ds.MaxOriginConnections = util.IntPtr(0)

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if strings.Contains(remapLine, "hdr_rw") {
		t.Errorf("expected remap with default header rewrites to not have header rewrite directive for a file that won't exist, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigMidNoHeaderRewrite(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "MID"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("DNS")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest)
	ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua")
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

	// non-nil default values should not trigger header rewrite plugin directive
	ds.EdgeHeaderRewrite = util.StrPtr("")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.ServiceCategory = util.StrPtr("")
	ds.MaxOriginConnections = util.IntPtr(0)

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if strings.Contains(remapLine, "hdr_rw") {
		t.Errorf("expected remap with default header rewrites to not have header rewrite directive for a file that won't exist, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigMidNoNoCacheRemapLine(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "MID"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_NO_CACHE")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("origin.example.test")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest)
	ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua")
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

	// non-nil default values should not trigger header rewrite plugin directive
	ds.EdgeHeaderRewrite = util.StrPtr("")
	ds.MidHeaderRewrite = util.StrPtr("mid-header-rewrite")
	ds.ServiceCategory = util.StrPtr("")
	ds.MaxOriginConnections = util.IntPtr(0)

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected 0 remaps from HTTP_NO_CACHE DS on Mid, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[0]

	if strings.Contains(remapLine, "hdr_rw_mid_mydsname.config") {
		t.Errorf("expected remap line for HTTP_NO_CACHE to not exist on Mid server, regardless of Mid Header Rewrite, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeHTTPOriginHTTPRemap(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("http://origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "https://origin.example.test") {
		t.Errorf("expected edge->origin http origin to not create https remap target, actual: ''%v''' warn %v", txt, cfg.Warnings)
	}

	if !strings.Contains(remapLine, "http://origin.example.test") {
		t.Errorf("expected edge->origin http origin to create http remap target, actual: ''%v'''", txt)
	}

	if strings.Contains(remapLine, "443") {
		t.Errorf("expected https origin to create http remap target not using 443 (edge->mid communication always uses http), actual: ''%v'''", txt)
	}
}

func TestMakeRemapDotConfigEdgeHTTPSOriginHTTPRemap(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("https://origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "http://origin.example.test") {
		t.Errorf("expected edge->origin https origin to not create http remap target, actual: ''%v''' warn %v", txt, cfg.Warnings)
	}

	if !strings.Contains(remapLine, "https://origin.example.test") {
		t.Errorf("expected edge->origin https origin to create https remap target, actual: ''%v'''", txt)
	}

	if strings.Contains(remapLine, "443") {
		t.Errorf("expected https origin to create http remap target not using 443 (edge->mid communication always uses http), actual: ''%v'''", txt)
	}
}

func TestMakeRemapDotConfigMidHTTPSOriginHTTPRemap(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "MID"
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("DNS")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("https://origin.example.test")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest)
	ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua")
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
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTP))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)

	// non-nil default values should not trigger header rewrite plugin directive
	ds.EdgeHeaderRewrite = util.StrPtr("")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.ServiceCategory = util.StrPtr("")
	ds.MaxOriginConnections = util.IntPtr(0)

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "map http://origin.example.test https://origin.example.test") {
		t.Errorf("expected mid https origin to create remap from http to https (edge->mid communication always uses http, but the origin needs to still be https), actual: ''%v'''", txt)
	}
}

func TestMakeRemapDotConfigEdgeHTTPSOriginHTTPRemapTopology(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "EDGE"
	server.Cachegroup = util.StrPtr("edgeCG")
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("HTTP_LIVE_NATNL")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("https://origin.example.test")
	ds.MidHeaderRewrite = util.StrPtr("")
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
	ds.Topology = util.StrPtr("t0")

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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

	topologies := []tc.Topology{
		{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				{
					Cachegroup: "edgeCG",
					Parents:    []int{1, 2},
				},
				{
					Cachegroup: "midCG",
				},
				{
					Cachegroup: "midCG2",
				},
			},
		},
	}

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.HostName = util.StrPtr("mymid0")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG")
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")

	eCG := &tc.CacheGroupNullable{}
	eCG.Name = server.Cachegroup
	eCG.ID = server.CachegroupID
	eCG.ParentName = mid0.Cachegroup
	eCG.ParentCachegroupID = mid0.CachegroupID
	eCG.SecondaryParentName = mid1.Cachegroup
	eCG.SecondaryParentCachegroupID = mid1.CachegroupID
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	mCG := &tc.CacheGroupNullable{}
	mCG.Name = mid0.Cachegroup
	mCG.ID = mid0.CachegroupID
	mCGType := tc.CacheGroupMidTypeName
	mCG.Type = &mCGType

	mCG2 := &tc.CacheGroupNullable{}
	mCG2.Name = mid1.Cachegroup
	mCG2.ID = mid1.CachegroupID
	mCGType2 := tc.CacheGroupMidTypeName
	mCG2.Type = &mCGType2

	cgs := []tc.CacheGroupNullable{*eCG, *mCG, *mCG2}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if strings.Contains(remapLine, "https://origin.example.test") {
		t.Errorf("expected https origin to create http remap target not https (edge->mid communication always uses http), actual: ''%v'''", txt)
	}

	if !strings.Contains(remapLine, "http://origin.example.test") {
		t.Errorf("expected https origin to create http remap target (edge->mid communication always uses http), actual: ''%v'''", txt)
	}

	if strings.Contains(remapLine, "443") {
		t.Errorf("expected topology edge https origin to create http remap target not using 443 (edge->mid communication always uses http), actual: ''%v'''", txt)
	}
}

func TestMakeRemapDotConfigMidHTTPSOriginHTTPRemapTopology(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "MID"
	server.Cachegroup = util.StrPtr("midCG")
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("DNS")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("https://origin.example.test")
	ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest)
	ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua")
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
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTP))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)
	ds.Topology = util.StrPtr("t0")

	// non-nil default values should not trigger header rewrite plugin directive
	ds.EdgeHeaderRewrite = util.StrPtr("")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.ServiceCategory = util.StrPtr("")
	ds.MaxOriginConnections = util.IntPtr(0)

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
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

	topologies := []tc.Topology{
		{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				{
					Cachegroup: "edgeCG",
					Parents:    []int{1, 2},
				},
				{
					Cachegroup: "midCG",
				},
				{
					Cachegroup: "midCG2",
				},
			},
		},
	}

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.HostName = util.StrPtr("mymid0")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG")
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")

	eCG := &tc.CacheGroupNullable{}
	eCG.Name = server.Cachegroup
	eCG.ID = server.CachegroupID
	eCG.ParentName = mid0.Cachegroup
	eCG.ParentCachegroupID = mid0.CachegroupID
	eCG.SecondaryParentName = mid1.Cachegroup
	eCG.SecondaryParentCachegroupID = mid1.CachegroupID
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	mCG := &tc.CacheGroupNullable{}
	mCG.Name = mid0.Cachegroup
	mCG.ID = mid0.CachegroupID
	mCGType := tc.CacheGroupMidTypeName
	mCG.Type = &mCGType

	mCG2 := &tc.CacheGroupNullable{}
	mCG2.Name = mid1.Cachegroup
	mCG2.ID = mid1.CachegroupID
	mCGType2 := tc.CacheGroupMidTypeName
	mCG2.Type = &mCGType2

	cgs := []tc.CacheGroupNullable{*eCG, *mCG, *mCG2}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}
	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 3 { // comment, blank, and remap
		t.Fatalf("expected a comment header, a blank line, and 1 remaps from HTTP_TO_HTTPS DS, actual: '%v' count %v", txt, len(txtLines))
	}

	remapLine := txtLines[2]

	if !strings.HasPrefix(remapLine, "map") {
		t.Errorf("expected to start with 'map', actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "map http://origin.example.test https://origin.example.test") {
		t.Errorf("expected topology mid https origin to create remap from http to https (edge->mid communication always uses http, but the origin needs to still be https), actual: ''%v'''", txt)
	}
}

func TestMakeRemapDotConfigMidLastRawRemap(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeTestRemapServer()
	server.Type = "MID"
	server.Cachegroup = util.StrPtr("midCG")
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("DNS")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("http://origin.example.test")
	/*
		ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest)
		ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua")
	*/
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
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTP))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)
	ds.Topology = util.StrPtr("t0")

	// non-nil default values should not trigger header rewrite plugin directive
	ds.EdgeHeaderRewrite = util.StrPtr("")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.ServiceCategory = util.StrPtr("")
	ds.MaxOriginConnections = util.IntPtr(0)

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "LastRawRemapPost",
			ConfigFile: "remap.config",
			Value:      "remap http://penraw/ http://penraw0/",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "LastRawRemapPost",
			ConfigFile: "remap.config",
			Value:      "remap http://lastraw/ http://lastraw0/",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "LastRawRemapPre",
			ConfigFile: "remap.config",
			Value:      "map_with_recp_port http://firstraw:8000/ http://firstraw0/",
			Profiles:   []byte(`["dsprofile"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{
		{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				{
					Cachegroup: "edgeCG",
					Parents:    []int{1, 2},
				},
				{
					Cachegroup: "midCG",
				},
				{
					Cachegroup: "midCG2",
				},
			},
		},
	}

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.HostName = util.StrPtr("mymid0")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG")
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")

	eCG := &tc.CacheGroupNullable{}
	eCG.Name = server.Cachegroup
	eCG.ID = server.CachegroupID
	eCG.ParentName = mid0.Cachegroup
	eCG.ParentCachegroupID = mid0.CachegroupID
	eCG.SecondaryParentName = mid1.Cachegroup
	eCG.SecondaryParentCachegroupID = mid1.CachegroupID
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	mCG := &tc.CacheGroupNullable{}
	mCG.Name = mid0.Cachegroup
	mCG.ID = mid0.CachegroupID
	mCGType := tc.CacheGroupMidTypeName
	mCG.Type = &mCGType

	mCG2 := &tc.CacheGroupNullable{}
	mCG2.Name = mid1.Cachegroup
	mCG2.ID = mid1.CachegroupID
	mCGType2 := tc.CacheGroupMidTypeName
	mCG2.Type = &mCGType2

	cgs := []tc.CacheGroupNullable{*eCG, *mCG, *mCG2}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, &RemapDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, hdr)

	txtLines := strings.Split(txt, "\n")

	linesexp := 6

	if len(txtLines) != linesexp {
		t.Fatalf("expected %d lines in comment, actual: '%v' count %v", linesexp, txt, len(txtLines))
	} else {
		var commentstr string

		commentstr = txtLines[0]
		if len(commentstr) == 0 || '#' != commentstr[0] {
			t.Fatalf("expected [0] as comment, actual: \n'%v' got %v", txt, commentstr)
		}

		blankStr := txtLines[1]
		if strings.TrimSpace(blankStr) != "" {
			t.Fatalf("expected [1] as blank line after comment, actual: \n'%v' got %v", txt, commentstr)
		}

		firststr := txtLines[2]
		if !strings.Contains(firststr, "firstraw") {
			t.Fatalf("expected [2] with 'firstraw', actual: '%v' got %v", txt, firststr)
		}
		laststr := txtLines[len(txtLines)-2]
		if !strings.Contains(laststr, "last") {
			t.Fatalf("expected [-2] last with 'lastraw', actual: '%v' got %v", txt, laststr)
		}
		penstr := txtLines[len(txtLines)-1]
		if !strings.Contains(penstr, "penraw") {
			t.Fatalf("expected [-1] with 'penraw', actual: '%v' got %v", txt, penstr)
		}
	}
}

func TestMakeRemapDotConfigStrategies(t *testing.T) {
	opt := &RemapDotConfigOpts{
		HdrComment:    "myHeaderComment",
		UseStrategies: true,
	}

	server := makeTestRemapServer()
	server.Type = "MID"
	server.Cachegroup = util.StrPtr("midCG")
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("DNS")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("http://origin.example.test")
	/*
		ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest)
		ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua")
	*/
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
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTP))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)
	ds.Topology = util.StrPtr("t0")

	// non-nil default values should not trigger header rewrite plugin directive
	ds.EdgeHeaderRewrite = util.StrPtr("")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.ServiceCategory = util.StrPtr("")
	ds.MaxOriginConnections = util.IntPtr(0)

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "LastRawRemapPost",
			ConfigFile: "remap.config",
			Value:      "remap http://penraw/ http://penraw0/",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "LastRawRemapPost",
			ConfigFile: "remap.config",
			Value:      "remap http://lastraw/ http://lastraw0/",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "LastRawRemapPre",
			ConfigFile: "remap.config",
			Value:      "map_with_recp_port http://firstraw:8000/ http://firstraw0/",
			Profiles:   []byte(`["dsprofile"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{
		{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				{
					Cachegroup: "edgeCG",
					Parents:    []int{1, 2},
				},
				{
					Cachegroup: "midCG",
				},
				{
					Cachegroup: "midCG2",
				},
			},
		},
	}

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.HostName = util.StrPtr("mymid0")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG")
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")

	eCG := &tc.CacheGroupNullable{}
	eCG.Name = server.Cachegroup
	eCG.ID = server.CachegroupID
	eCG.ParentName = mid0.Cachegroup
	eCG.ParentCachegroupID = mid0.CachegroupID
	eCG.SecondaryParentName = mid1.Cachegroup
	eCG.SecondaryParentCachegroupID = mid1.CachegroupID
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	mCG := &tc.CacheGroupNullable{}
	mCG.Name = mid0.Cachegroup
	mCG.ID = mid0.CachegroupID
	mCGType := tc.CacheGroupMidTypeName
	mCG.Type = &mCGType

	mCG2 := &tc.CacheGroupNullable{}
	mCG2.Name = mid1.Cachegroup
	mCG2.ID = mid1.CachegroupID
	mCGType2 := tc.CacheGroupMidTypeName
	mCG2.Type = &mCGType2

	cgs := []tc.CacheGroupNullable{*eCG, *mCG, *mCG2}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, opt)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, opt.HdrComment)

	if !strings.Contains(txt, `@plugin=parent_select.so @pparam=/opt/trafficserver/etc/trafficserver/strategies.yaml @pparam=strategy-mydsname`) {
		t.Fatalf("expected parent_select plugin for opt.UseStrategies, actual: %v", txt)
	}
}

func TestMakeRemapDotConfigStrategiesFalseButCoreUnused(t *testing.T) {
	opt := &RemapDotConfigOpts{
		HdrComment:        "myHeaderComment",
		UseStrategies:     false,
		UseStrategiesCore: true,
	}

	server := makeTestRemapServer()
	server.Type = "MID"
	server.Cachegroup = util.StrPtr("midCG")
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("DNS")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("http://origin.example.test")
	/*
		ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest)
		ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua")
	*/
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
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTP))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)
	ds.Topology = util.StrPtr("t0")

	// non-nil default values should not trigger header rewrite plugin directive
	ds.EdgeHeaderRewrite = util.StrPtr("")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.ServiceCategory = util.StrPtr("")
	ds.MaxOriginConnections = util.IntPtr(0)

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "LastRawRemapPost",
			ConfigFile: "remap.config",
			Value:      "remap http://penraw/ http://penraw0/",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "LastRawRemapPost",
			ConfigFile: "remap.config",
			Value:      "remap http://lastraw/ http://lastraw0/",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "LastRawRemapPre",
			ConfigFile: "remap.config",
			Value:      "map_with_recp_port http://firstraw:8000/ http://firstraw0/",
			Profiles:   []byte(`["dsprofile"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{
		{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				{
					Cachegroup: "edgeCG",
					Parents:    []int{1, 2},
				},
				{
					Cachegroup: "midCG",
				},
				{
					Cachegroup: "midCG2",
				},
			},
		},
	}

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.HostName = util.StrPtr("mymid0")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG")
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")

	eCG := &tc.CacheGroupNullable{}
	eCG.Name = server.Cachegroup
	eCG.ID = server.CachegroupID
	eCG.ParentName = mid0.Cachegroup
	eCG.ParentCachegroupID = mid0.CachegroupID
	eCG.SecondaryParentName = mid1.Cachegroup
	eCG.SecondaryParentCachegroupID = mid1.CachegroupID
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	mCG := &tc.CacheGroupNullable{}
	mCG.Name = mid0.Cachegroup
	mCG.ID = mid0.CachegroupID
	mCGType := tc.CacheGroupMidTypeName
	mCG.Type = &mCGType

	mCG2 := &tc.CacheGroupNullable{}
	mCG2.Name = mid1.Cachegroup
	mCG2.ID = mid1.CachegroupID
	mCGType2 := tc.CacheGroupMidTypeName
	mCG2.Type = &mCGType2

	cgs := []tc.CacheGroupNullable{*eCG, *mCG, *mCG2}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, opt)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, opt.HdrComment)

	if strings.Contains(txt, `strategy`) {
		t.Fatalf("expected no strategy core directive for !opt.UseStrategies but useless opt.UseStrategiesCores, actual: %v", txt)
	}

	if strings.Contains(txt, `parent_select`) {
		t.Fatalf("expected no strategy plugin directive for !opt.UseStrategies but useless opt.UseStrategiesCores, actual: %v", txt)
	}

	if !warningsContains(cfg.Warnings, "ot using strategies") {
		t.Errorf("expected warning about not using strategies with useless opt UseStrategiesCore but no UseStrategies, actual '%v'", cfg.Warnings)
	}
}

func TestMakeRemapDotConfigMidCacheParentHTTPSOrigin(t *testing.T) {
	opt := &RemapDotConfigOpts{
		HdrComment:        "myHeaderComment",
		UseStrategies:     false,
		UseStrategiesCore: true,
	}

	server := makeTestRemapServer()
	server.Type = "MID"
	server.Cachegroup = util.StrPtr("midCG")
	servers := makeTestAnyCastServers()

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("DNS")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("https://origin.example.test")
	/*
		ds.RangeRequestHandling = util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest)
		ds.RemapText = util.StrPtr("@plugin=tslua.so @pparam=my-range-manipulator.lua")
	*/
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
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTP))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)
	ds.Topology = util.StrPtr("t0")

	// non-nil default values should not trigger header rewrite plugin directive
	ds.EdgeHeaderRewrite = util.StrPtr("")
	ds.MidHeaderRewrite = util.StrPtr("")
	ds.ServiceCategory = util.StrPtr("")
	ds.MaxOriginConnections = util.IntPtr(0)

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "serverpkgval",
			ConfigFile: "package",
			Value:      "serverpkgval __HOSTNAME__ foo",
			Profiles:   []byte(server.ProfileNames[0]),
		},
		tc.Parameter{
			Name:       "dscp_remap_no",
			ConfigFile: "package",
			Value:      "notused",
			Profiles:   []byte(server.ProfileNames[0]),
		},
	}

	remapConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "LastRawRemapPost",
			ConfigFile: "remap.config",
			Value:      "remap http://penraw/ http://penraw0/",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "LastRawRemapPost",
			ConfigFile: "remap.config",
			Value:      "remap http://lastraw/ http://lastraw0/",
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "LastRawRemapPre",
			ConfigFile: "remap.config",
			Value:      "map_with_recp_port http://firstraw:8000/ http://firstraw0/",
			Profiles:   []byte(`["dsprofile"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{
		{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				{
					Cachegroup: "edgeCG",
					Parents:    []int{1},
				},
				{
					Cachegroup: "midCG",
					Parents:    []int{2},
				},
				{
					Cachegroup: "midCG2",
				},
			},
		},
	}

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.HostName = util.StrPtr("mymid0")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG2")
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")

	eCG := &tc.CacheGroupNullable{}
	eCG.Name = server.Cachegroup
	eCG.ID = server.CachegroupID
	eCG.ParentName = mid0.Cachegroup
	eCG.ParentCachegroupID = mid0.CachegroupID
	eCG.SecondaryParentName = mid1.Cachegroup
	eCG.SecondaryParentCachegroupID = mid1.CachegroupID
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	mCG := &tc.CacheGroupNullable{}
	mCG.Name = mid0.Cachegroup
	mCG.ID = mid0.CachegroupID
	mCGType := tc.CacheGroupMidTypeName
	mCG.Type = &mCGType

	mCG2 := &tc.CacheGroupNullable{}
	mCG2.Name = mid1.Cachegroup
	mCG2.ID = mid1.CachegroupID
	mCGType2 := tc.CacheGroupMidTypeName
	mCG2.Type = &mCGType2

	cgs := []tc.CacheGroupNullable{*eCG, *mCG, *mCG2}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	configDir := `/opt/trafficserver/etc/trafficserver`

	cfg, err := MakeRemapDotConfig(server, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, opt)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	txt = strings.TrimSpace(txt)

	testComment(t, txt, opt.HdrComment)

	if strings.Contains(txt, `https://origin.example.test`) {
		t.Errorf("expected mid with a parent mid cache to set remap target scheme to http and not origin https, actual: %v", txt)
	}
	if !strings.Contains(txt, `http://origin.example.test`) {
		t.Errorf("expected mid with a parent mid cache to set remap target scheme to http, actual: %v", txt)
	}
}

func TestMakeRemapDotConfigRemapTemplate(t *testing.T) {
	opt := &RemapDotConfigOpts{
		HdrComment:        "myHeaderComment",
		UseStrategies:     true,
		UseStrategiesCore: true,
	}

	edge := makeTestRemapServer()
	edge.Type = "EDGE"
	edge.Cachegroup = util.StrPtr("edgeCG")
	servers := makeTestAnyCastServers()

	mid := makeTestParentServer()
	mid.Type = "MID"
	mid.Cachegroup = util.StrPtr("midCG")
	mid.HostName = util.StrPtr("mymid")
	mid.ID = util.IntPtr(45)
	setIP(mid, "192.168.2.2")

	opl := makeTestParentServer()
	opl.Type = "MID"
	opl.Cachegroup = util.StrPtr("oplCG")
	opl.HostName = util.StrPtr("myopl")
	opl.ID = util.IntPtr(46)
	setIP(opl, "192.168.2.3")

	ds := DeliveryService{}
	ds.ID = util.IntPtr(48)
	dsType := tc.DSType("DNS")
	ds.Type = &dsType
	ds.OrgServerFQDN = util.StrPtr("https://origin.example.test")
	ds.SigningAlgorithm = util.StrPtr("foo")
	ds.XMLID = util.StrPtr("mydsname")
	ds.QStringIgnore = util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp))
	ds.RangeRequestHandling = util.IntPtr(int(tc.RangeRequestHandlingBackgroundFetch))
	ds.RegexRemap = util.StrPtr("myregexremap")
	ds.RemapText = util.StrPtr("@plugin=rawtextplugin.so")
	ds.FQPacingRate = util.IntPtr(314159)
	ds.DSCP = util.IntPtr(0)
	ds.RoutingName = util.StrPtr("myroutingname")
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.OriginShield = util.StrPtr("myoriginshield")
	ds.ProfileID = util.IntPtr(49)
	ds.ProfileName = util.StrPtr("dsprofile")
	ds.Protocol = util.IntPtr(int(tc.DSProtocolHTTP))
	ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)
	ds.Topology = util.StrPtr("t0")
	ds.SigningAlgorithm = util.StrPtr("url_sig")

	// non-nil default values should not trigger header rewrite plugin directive
	ds.FirstHeaderRewrite = util.StrPtr("firstfoo")
	ds.InnerHeaderRewrite = util.StrPtr("innerfoo")
	ds.LastHeaderRewrite = util.StrPtr("lastfoo")
	ds.ServiceCategory = util.StrPtr("")
	ds.MaxOriginConnections = util.IntPtr(0)

	dses := []DeliveryService{ds}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *edge.ID,
			DeliveryService: *ds.ID,
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
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
	}

	const FirstTemplateString = `mapfoo {{{Destination}}} {{{Source}}} {{{Signing}}} {{{HeaderRewrite}}} {{{RegexRemap}}} {{{Dscp}}} {{{Signing}}} {{{RawText}}} {{{Cachekey}}} {{{Pacing}}} {{{RangeRequests}}}
map http://foo/ http://bar/`

	const InnerTemplateString = `map {{{Source}}} {{{Destination}}} {{{Cachekey}}} {{{HeaderRewrite}}} {{{RangeRequests}}} @plugin=prefetch.so @pparam=--backend=true`
	const LastTemplateString = `map_with_recv_port {{{Source}}}:3600 {{{Destination}}} {{{HeaderRewrite}}} {{{Cachekey}}}`

	remapConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       "not_location",
			ConfigFile: "cachekey.config",
			Value:      "notinconfig",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       "template.first",
			ConfigFile: "remap.config",
			Value:      FirstTemplateString,
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "template.inner",
			ConfigFile: "remap.config",
			Value:      InnerTemplateString,
			Profiles:   []byte(`["dsprofile"]`),
		},
		tc.Parameter{
			Name:       "template.last",
			ConfigFile: "remap.config",
			Value:      LastTemplateString,
			Profiles:   []byte(`["dsprofile"]`),
		},
	}

	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	topologies := []tc.Topology{
		{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				{
					Cachegroup: "edgeCG",
					Parents:    []int{1},
				},
				{
					Cachegroup: "midCG",
					Parents:    []int{2},
				},
				{
					Cachegroup: "oplCG",
				},
			},
		},
	}

	eCG := &tc.CacheGroupNullable{}
	eCG.Name = edge.Cachegroup
	eCG.ID = edge.CachegroupID
	eCG.ParentName = mid.Cachegroup
	eCG.ParentCachegroupID = mid.CachegroupID
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	mCG := &tc.CacheGroupNullable{}
	mCG.Name = mid.Cachegroup
	mCG.ID = mid.CachegroupID
	mCG.ParentName = opl.Cachegroup
	mCG.ParentCachegroupID = opl.CachegroupID
	mCGType := tc.CacheGroupMidTypeName
	mCG.Type = &mCGType

	oCG := &tc.CacheGroupNullable{}
	oCG.Name = opl.Cachegroup
	oCG.ID = opl.CachegroupID
	oCGType := tc.CacheGroupMidTypeName
	oCG.Type = &oCGType

	cgs := []tc.CacheGroupNullable{*eCG, *mCG, *oCG}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	configDir := `/opt/trafficserver/etc/trafficserver`

	{ // first override
		cfg, err := MakeRemapDotConfig(edge, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text
		txt = strings.TrimSpace(txt)
		testComment(t, txt, opt.HdrComment)

		pluginsExp := []string{"url_sig.so", "header_rewrite.so", "regex_remap.so", "header_rewrite.so", "url_sig.so", "rawtextplugin.so", "cachekey.so", "fq_pacing.so", "background_fetch.so"}
		tokens := tokenize(txt)
		pluginsGot := pluginsFromTokens(tokens, "@plugin=")

		if !reflect.DeepEqual(pluginsGot, pluginsExp) {
			t.Errorf("unexpected plugins order for first override: '%v'\ngot: '%v'\nexp: '%v'", txt, pluginsGot, pluginsExp)
		}

		// check for swapped source and dest and munged name
		if !strings.Contains(txt, `mapfoo http://origin.example.test/ http://myregexpattern/`) {
			t.Errorf("expected munged map and swapped dest/source, got: '%v'", txt)
		}

		// check for tacked on remap rule
		if !strings.Contains(txt, `map http://foo/ http://bar/`) {
			t.Errorf("expected tacked on remap rule, got: '%v'", txt)
		}
	}

	{ // inner override
		cfg, err := MakeRemapDotConfig(mid, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text
		txt = strings.TrimSpace(txt)
		testComment(t, txt, opt.HdrComment)

		pluginsExp := []string{"cachekey.so", "header_rewrite.so", "prefetch.so"}
		tokens := tokenize(txt)
		pluginsGot := pluginsFromTokens(tokens, "@plugin=")

		if !reflect.DeepEqual(pluginsGot, pluginsExp) {
			t.Errorf("unexpected plugins order for inner override: '%v'\ngot: '%v'\nexp: '%v'", txt, pluginsGot, pluginsExp)
		}
	}

	{ // last override
		cfg, err := MakeRemapDotConfig(opl, servers, dses, dss, dsRegexes, serverParams, cdn, remapConfigParams, topologies, cgs, serverCapabilities, dsRequiredCapabilities, configDir, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text
		txt = strings.TrimSpace(txt)
		testComment(t, txt, opt.HdrComment)

		pluginsExp := []string{"header_rewrite.so", "cachekey.so"}
		tokens := tokenize(txt)
		pluginsGot := pluginsFromTokens(tokens, "@plugin=")

		if !reflect.DeepEqual(pluginsGot, pluginsExp) {
			t.Errorf("unexpected plugins order for inner override: '%v'\ngot: '%v'\nexp: '%v'", txt, pluginsGot, pluginsExp)
		}

		if !strings.Contains(txt, `map_with_recv_port http://origin.example.test:3600 https://origin.example.test`) {
			t.Errorf("expected laster to have map_with_recv_port directive, got '%v'", txt)
		}
	}
}
