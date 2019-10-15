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
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE",
			OriginFQDN:               util.StrPtr("origin.example.test"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

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
		t.Errorf("expected to contain routing name, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "http://myregexpattern") {
		t.Errorf("expected to contain routing name, actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "origin.example.test") {
		t.Errorf("expected to contain origin FQDN, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigMidLiveLocalExcluded(t *testing.T) {
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "MID",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE",
			OriginFQDN:               util.StrPtr("origin.example.test"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected no remap lines for LIVE local DS, actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigMid(t *testing.T) {
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "MID",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("origin.example.test"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

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

	if strings.Count(remapLine, "origin.example.test") != 2 {
		t.Errorf("expected to contain origin FQDN twice (Mids remap origins to themselves, as a forward proxy), actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "hdr_rw_mid_"+"mydsname"+".config") {
		t.Errorf("expected to contain header rewrite for DS, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigNilOrigin(t *testing.T) {
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "MID",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr(""),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected no remap lines for DS with nil Origin FQDN, actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigEmptyOrigin(t *testing.T) {
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "MID",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               nil,
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected no remap lines for DS with empty Origin FQDN, actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigDuplicateOrigins(t *testing.T) {
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "MID",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("origin.example.test"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
		RemapConfigDSData{
			ID:                       49,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("origin.example.test"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite2"),
			CacheURL:                 util.StrPtr("mycacheurl2"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname2": "cachekeyparamval2"},
			RemapText:                util.StrPtr("myremaptext2"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite2"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname2",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap2"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname2"),
			MultiSiteOrigin:          util.StrPtr("mymso2"),
			Pattern:                  util.StrPtr("myregexpattern2"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain2"),
			RegexSetNumber:           util.StrPtr("myregexsetnum2"),
			OriginShield:             util.StrPtr("myoriginshield2"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected 1 remap lines for multiple DSes with the same Origin (ATS can't handle multiple remaps with the same origin FQDN), actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigNilMidRewrite(t *testing.T) {
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "MID",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("origin.example.test"),
			MidHeaderRewrite:         nil,
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "MID",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("origin.example.test"),
			MidHeaderRewrite:         util.StrPtr(""),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

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

	if strings.Count(remapLine, "origin.example.test") != 2 {
		t.Errorf("expected to contain origin FQDN twice (Mids remap origins to themselves, as a forward proxy), actual '%v'", txt)
	}

	if strings.Contains(remapLine, "hdr_rw_mid_") {
		t.Errorf("expected no 'hdr_rw_mid_' for nil mid header rewrite on DS, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigMidQStringPassUpATS7CacheKey(t *testing.T) {
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 6

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "MID",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("origin.example.test"),
			MidHeaderRewrite:         util.StrPtr(""),
			CacheURL:                 util.StrPtr(""), // no cacheurl, so we can see if the qstring puts one in
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 5

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "MID",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("origin.example.test"),
			MidHeaderRewrite:         util.StrPtr(""),
			CacheURL:                 util.StrPtr(""), // no cacheurl, so we can see if the qstring puts one in
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		49: map[string]string{
			"cachekeykey": "cachekeyval",
		},
		42: map[string]string{
			"shouldnotexist": "shouldnotexisteither",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "MID",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("origin.example.test"),
			MidHeaderRewrite:         util.StrPtr(""),
			CacheURL:                 util.StrPtr("mycacheurl"), // cacheurl, so it gets added twice, also with qstring
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		49: map[string]string{
			"cachekeykey": "cachekeyval",
		},
		42: map[string]string{
			"shouldnotexist": "shouldnotexisteither",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "MID",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("origin.example.test"),
			MidHeaderRewrite:         util.StrPtr(""),
			CacheURL:                 util.StrPtr("mycacheurl"), // cacheurl, so it gets added twice, also with qstring
			RangeRequestHandling:     util.IntPtr(int(tc.RangeRequestHandlingCacheRangeRequest)),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

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

	if strings.Count(remapLine, "origin.example.test") != 2 {
		t.Errorf("expected to contain origin FQDN twice (Mids remap origins to themselves, as a forward proxy), actual '%v'", txt)
	}

	if !strings.Contains(remapLine, "cache_range_requests.so") {
		t.Errorf("expected to contain range request handling plugin, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigFirstExcludedSecondIncluded(t *testing.T) {
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "MID",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               nil, // this DS should not be included
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"), // this DS should be included
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 2 {
		t.Fatalf("expected one remap line for DS with origin, but not DS with empty Origin FQDN, actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigAnyMap(t *testing.T) {
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "ANY_MAP",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                nil, // this DS shouldn't be included - anymap with nil remap text
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
		RemapConfigDSData{
			ID:                       48,
			Type:                     "ANY_MAP",
			OriginFQDN:               util.StrPtr("myorigin"), // this DS should be included
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                nil, // nil regex should not be included
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypePathRegex)), // non-host regex should not be included
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeSteeringRegex)), // non-host regex should not be included
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHeaderRegex)), // non-host regex should not be included
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(""), // non-host regex should not be included
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr("nonexistenttype"), // non-host regex should not be included
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               nil, // nil origin should not be included
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr(""),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  nil, // nil pattern shouldn't be included
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 nil, // nil protocol shouldn't be included
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr("myregexpattern"),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   nil, // nil domain shouldn't be included
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(0),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected no remaps from DSes with missing data, actual: '%v' count %v", txt, len(txtLines))
	}

}

func TestMakeRemapDotConfigEdgeHostRegexReplacement(t *testing.T) {
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`.*\.mypattern\..*`), // common host regex syntax, should be replaced
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPAndHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`.*\.mypattern\..*`), // common host regex syntax, should be replaced
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTP)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`.*\.mypattern\..*`), // common host regex syntax, should be replaced
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`.*\.mypattern\..*`), // common host regex syntax, should be replaced
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"serverpkgval": "serverpkgval __HOSTNAME__ foo",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`myliteralpattern__http__foo`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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

	if !strings.Contains(remapLine, "myliteralpattern"+serverInfo.HostName+"foo") {
		t.Errorf("expected literal pattern to replace __http__ with server name, actual '%v'", txt)
	}
}

func TestMakeRemapDotConfigEdgeDSCPRemap(t *testing.T) {
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr("myedgeheaderrewrite"),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        util.StrPtr(""),
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("url_sig"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("uri_signing"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         nil,
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr(""),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(0),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreDropAtEdge)),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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

	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"not_location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		49: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 5 // ats <6 uses cacheurl, not cache key

	cacheURLConfigParams := map[string]string{
		"notlocation": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 nil,
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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

	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 5 // ats <6 uses cacheurl, not cache key

	cacheURLConfigParams := map[string]string{
		"notlocation": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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

	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 5 // ats <6 uses cacheurl, not cache key

	cacheURLConfigParams := map[string]string{
		"notlocation": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "MID",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		46: map[string]string{
			"cachekeykey": "cachekeyval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr("mycacheurl"),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		49: map[string]string{
			"cachekeykey": "cachekeyval",
		},
		44: map[string]string{
			"shouldnotincludeotherprofile": "shouldnotincludeotherprofileval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr(""),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		49: map[string]string{
			"cachekeykey": "cachekeyval",
		},
		44: map[string]string{
			"shouldnotincludeotherprofile": "shouldnotincludeotherprofileval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr(""),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr("myregexremap"),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		49: map[string]string{
			"cachekeykey": "cachekeyval",
		},
		44: map[string]string{
			"shouldnotincludeotherprofile": "shouldnotincludeotherprofileval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr(""),
			RangeRequestHandling:     util.IntPtr(0),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr(""),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		49: map[string]string{
			"cachekeykey": "cachekeyval",
		},
		44: map[string]string{
			"shouldnotincludeotherprofile": "shouldnotincludeotherprofileval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr(""),
			RangeRequestHandling:     nil,
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr(""),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		49: map[string]string{
			"cachekeykey": "cachekeyval",
		},
		44: map[string]string{
			"shouldnotincludeotherprofile": "shouldnotincludeotherprofileval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr(""),
			RangeRequestHandling:     util.IntPtr(tc.RangeRequestHandlingDontCache),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr(""),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		49: map[string]string{
			"cachekeykey": "cachekeyval",
		},
		44: map[string]string{
			"shouldnotincludeotherprofile": "shouldnotincludeotherprofileval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr(""),
			RangeRequestHandling:     util.IntPtr(tc.RangeRequestHandlingBackgroundFetch),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr(""),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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

func TestMakeRemapDotConfigEdgeRangeRequestCache(t *testing.T) {
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		49: map[string]string{
			"cachekeykey": "cachekeyval",
		},
		44: map[string]string{
			"shouldnotincludeotherprofile": "shouldnotincludeotherprofileval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr(""),
			RangeRequestHandling:     util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr(""),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		49: map[string]string{
			"cachekeykey": "cachekeyval",
		},
		44: map[string]string{
			"shouldnotincludeotherprofile": "shouldnotincludeotherprofileval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr(""),
			RangeRequestHandling:     util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr(""),
			FQPacingRate:             nil,
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		49: map[string]string{
			"cachekeykey": "cachekeyval",
		},
		44: map[string]string{
			"shouldnotincludeotherprofile": "shouldnotincludeotherprofileval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr(""),
			RangeRequestHandling:     util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr(""),
			FQPacingRate:             util.IntPtr(-42),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		49: map[string]string{
			"cachekeykey": "cachekeyval",
		},
		44: map[string]string{
			"shouldnotincludeotherprofile": "shouldnotincludeotherprofileval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr(""),
			RangeRequestHandling:     util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr(""),
			FQPacingRate:             util.IntPtr(0),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		49: map[string]string{
			"cachekeykey": "cachekeyval",
		},
		44: map[string]string{
			"shouldnotincludeotherprofile": "shouldnotincludeotherprofileval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "HTTP_LIVE_NATNL",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr(""),
			RangeRequestHandling:     util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr(""),
			FQPacingRate:             util.IntPtr(314159),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`mypattern`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		49: map[string]string{
			"cachekeykey": "cachekeyval",
		},
		44: map[string]string{
			"shouldnotincludeotherprofile": "shouldnotincludeotherprofileval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "DNS_LIVE",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr(""),
			RangeRequestHandling:     util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr(""),
			FQPacingRate:             util.IntPtr(314159),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`.*\.mypattern\..*`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

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
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		49: map[string]string{
			"cachekeykey": "cachekeyval",
		},
		44: map[string]string{
			"shouldnotincludeotherprofile": "shouldnotincludeotherprofileval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "DNS_LIVE",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr(""),
			RangeRequestHandling:     util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr(""),
			FQPacingRate:             util.IntPtr(314159),
			DSCP:                     0,
			RoutingName:              nil,
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`.*\.mypattern\..*`),
			RegexType:                util.StrPtr(string(tc.DSMatchTypeHostRegex)),
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected no remaps from DNS DS with nil routing name, actual: '%v' count %v", txt, len(txtLines))
	}
}

func TestMakeRemapDotConfigEdgeRegexTypeNil(t *testing.T) {
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"
	atsMajorVersion := 7

	cacheURLConfigParams := map[string]string{
		"location": "notinconfig",
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{
		49: map[string]string{
			"cachekeykey": "cachekeyval",
		},
		44: map[string]string{
			"shouldnotincludeotherprofile": "shouldnotincludeotherprofileval",
		},
	}

	serverPackageParamData := map[string]string{
		"dscp_remap_no": "notused",
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "mydomain",
		HostName:                      "myhost",
		HTTPSPort:                     12443,
		ID:                            44,
		IP:                            "192.168.2.4",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "CGType4",
		ProfileID:                     46,
		ProfileName:                   "MyProfile",
		Port:                          12080,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCG",
		Type:                          "EDGE",
	}

	remapDSData := []RemapConfigDSData{
		RemapConfigDSData{
			ID:                       48,
			Type:                     "DNS_LIVE",
			OriginFQDN:               util.StrPtr("myorigin"),
			MidHeaderRewrite:         util.StrPtr("mymidrewrite"),
			CacheURL:                 util.StrPtr(""),
			RangeRequestHandling:     util.IntPtr(tc.RangeRequestHandlingCacheRangeRequest),
			CacheKeyConfigParams:     map[string]string{"cachekeyparamname": "cachekeyparamval"},
			RemapText:                util.StrPtr("myremaptext"),
			EdgeHeaderRewrite:        nil,
			SigningAlgorithm:         util.StrPtr("foo"),
			Name:                     "mydsname",
			QStringIgnore:            util.IntPtr(int(tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp)),
			RegexRemap:               util.StrPtr(""),
			FQPacingRate:             util.IntPtr(314159),
			DSCP:                     0,
			RoutingName:              util.StrPtr("myroutingname"),
			MultiSiteOrigin:          util.StrPtr("mymso"),
			Pattern:                  util.StrPtr(`.*\.mypattern\..*`),
			RegexType:                nil,
			Domain:                   util.StrPtr("mydomain"),
			RegexSetNumber:           util.StrPtr("myregexsetnum"),
			OriginShield:             util.StrPtr("myoriginshield"),
			ProfileID:                util.IntPtr(49),
			Protocol:                 util.IntPtr(int(tc.DSProtocolHTTPToHTTPS)),
			AnonymousBlockingEnabled: util.BoolPtr(false),
			Active:                   true,
		},
	}

	txt := MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	txt = strings.TrimSpace(txt)

	testComment(t, txt, string(serverName), toToolName, toURL)

	txtLines := strings.Split(txt, "\n")

	if len(txtLines) != 1 {
		t.Fatalf("expected no remaps for DS with nil regex type, actual: '%v' count %v", txt, len(txtLines))
	}

}
