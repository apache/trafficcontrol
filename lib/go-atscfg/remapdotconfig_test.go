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
