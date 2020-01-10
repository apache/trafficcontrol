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
	"github.com/apache/trafficcontrol/lib/go-tc/enum"
	"strings"
	"testing"
)

func TestMakeParentDotConfig(t *testing.T) {
	atsMajorVer := 7
	serverName := "myserver"
	toolName := "myToolName"
	toURL := "https://myto.example.net"

	parentConfigDSes := []ParentConfigDSTopLevel{
		ParentConfigDSTopLevel{
			ParentConfigDS: ParentConfigDS{
				Name:            "ds0",
				QStringIgnore:   enum.QStringIgnoreUseInCacheKeyAndPassUp,
				OriginFQDN:      "http://ds0.example.net",
				MultiSiteOrigin: false,
				Type:            enum.DSTypeHTTP,
				QStringHandling: "ds0qstringhandling",
			},
		},
		ParentConfigDSTopLevel{
			ParentConfigDS: ParentConfigDS{
				Name:            "ds1",
				QStringIgnore:   enum.QStringIgnoreDrop,
				OriginFQDN:      "http://ds1.example.net",
				MultiSiteOrigin: false,
				Type:            enum.DSTypeDNS,
				QStringHandling: "ds1qstringhandling",
			},
		},
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "myCDN",
		CDNID:                         43,
		DomainName:                    "serverdomain.example.net",
		HostName:                      "myserver",
		ID:                            44,
		IP:                            "192.168.2.1",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "myParentCGType",
		ProfileID:                     46,
		ProfileName:                   "MyProfileName",
		Port:                          80,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCGType",
		Type:                          "EDGE",
	}

	serverParams := map[string]string{
		ParentConfigParamQStringHandling: "myQStringHandlingParam",
		ParentConfigParamAlgorithm:       enum.AlgorithmConsistentHash,
		ParentConfigParamQString:         "myQstringParam",
	}

	parentInfos := map[OriginHost][]ParentInfo{
		"ds1.example.net": []ParentInfo{
			ParentInfo{
				Host:            "my-parent-0",
				Port:            80,
				Domain:          "my-parent-0-domain",
				Weight:          "1",
				UseIP:           false,
				Rank:            1,
				IP:              "192.168.2.2",
				PrimaryParent:   true,
				SecondaryParent: true,
			},
		},
	}

	txt := MakeParentDotConfig(serverInfo, atsMajorVer, toolName, toURL, parentConfigDSes, serverParams, parentInfos)

	testComment(t, txt, serverName, toolName, toURL)

	if !strings.Contains(txt, "dest_domain=ds0.example.net") {
		t.Errorf("expected parent 'dest_domain=ds0.example.net', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "dest_domain=ds1.example.net") {
		t.Errorf("expected parent 'dest_domain=ds0.example.net', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "qstring=myQStringHandlingParam") {
		t.Errorf("expected qstring from param 'qstring=myQStringHandlingParam', actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigCapabilities(t *testing.T) {
	atsMajorVer := 7
	serverName := "myserver"
	toolName := "myToolName"
	toURL := "https://myto.example.net"

	parentConfigDSes := []ParentConfigDSTopLevel{
		ParentConfigDSTopLevel{
			ParentConfigDS: ParentConfigDS{
				Name:            "ds0",
				QStringIgnore:   enum.QStringIgnoreUseInCacheKeyAndPassUp,
				OriginFQDN:      "http://ds0.example.net",
				MultiSiteOrigin: false,
				Type:            enum.DSTypeHTTP,
				QStringHandling: "ds0qstringhandling",
				RequiredCapabilities: map[ServerCapability]struct{}{
					"FOO": {},
				},
			},
		},
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "myCDN",
		CDNID:                         43,
		DomainName:                    "serverdomain.example.net",
		HostName:                      "myserver",
		ID:                            44,
		IP:                            "192.168.2.1",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "myParentCGType",
		ProfileID:                     46,
		ProfileName:                   "MyProfileName",
		Port:                          80,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MySecondaryParentCGType",
		Type:                          "EDGE",
	}

	serverParams := map[string]string{
		ParentConfigParamQStringHandling: "myQStringHandlingParam",
		ParentConfigParamAlgorithm:       enum.AlgorithmConsistentHash,
		ParentConfigParamQString:         "myQstringParam",
	}

	parentInfos := map[OriginHost][]ParentInfo{
		DeliveryServicesAllParentsKey: []ParentInfo{
			ParentInfo{
				Host:            "my-parent-nocaps",
				Port:            80,
				Domain:          "my-parent-nocaps-domain",
				Weight:          "1",
				UseIP:           false,
				Rank:            1,
				IP:              "192.168.2.1",
				PrimaryParent:   true,
				SecondaryParent: true,
				Capabilities:    map[ServerCapability]struct{}{},
			},
			ParentInfo{
				Host:            "my-parent-fooonly",
				Port:            80,
				Domain:          "my-parent-fooonly-domain",
				Weight:          "1",
				UseIP:           false,
				Rank:            1,
				IP:              "192.168.2.2",
				PrimaryParent:   true,
				SecondaryParent: true,
				Capabilities: map[ServerCapability]struct{}{
					"FOO": {},
				},
			},
			ParentInfo{
				Host:            "my-parent-foobar",
				Port:            80,
				Domain:          "my-parent-foobar-domain",
				Weight:          "1",
				UseIP:           false,
				Rank:            1,
				IP:              "192.168.2.2",
				PrimaryParent:   true,
				SecondaryParent: true,
				Capabilities: map[ServerCapability]struct{}{
					"FOO": {},
					"BAR": {},
				},
			},
		},
	}

	txt := MakeParentDotConfig(serverInfo, atsMajorVer, toolName, toURL, parentConfigDSes, serverParams, parentInfos)

	testComment(t, txt, serverName, toolName, toURL)

	lines := strings.Split(txt, "\n")

	if len(lines) != 4 {
		t.Fatalf("expected 4 lines (comment, ds, dot remap, and empty newline), actual: '%+v'", len(lines))
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue // skip empty newline
		}
		if strings.HasPrefix(line, `dest_domain=.`) {
			continue // skip dot remap, which has all parents irrespective of capability
		}
		if strings.HasPrefix(line, `#`) {
			continue // skip comment
		}

		if !strings.Contains(line, "dest_domain=ds0.example.net") {
			t.Errorf("expected parent 'dest_domain=ds0.example.net', actual: '%v'", line)
		}
		if !strings.Contains(line, "foobar") {
			t.Errorf("expected parent with all capabilities, actual: '%v'", line)
		}
		if !strings.Contains(line, "fooonly") {
			t.Errorf("expected parent with required capabilities, actual: '%v'", line)
		}
		if strings.Contains(line, "nocaps") {
			t.Errorf("expected not to contain parent with no capabilities, actual line: '%v'", line)
		}
	}
}

func TestMakeParentDotConfigMSOSecondaryParent(t *testing.T) {
	atsMajorVer := 7
	serverName := "myserver"
	toolName := "myToolName"
	toURL := "https://myto.example.net"

	parentConfigDSes := []ParentConfigDSTopLevel{
		ParentConfigDSTopLevel{
			ParentConfigDS: ParentConfigDS{
				Name:            "ds0",
				QStringIgnore:   enum.QStringIgnoreUseInCacheKeyAndPassUp,
				OriginFQDN:      "http://ds0.example.net",
				MultiSiteOrigin: true,
				Type:            enum.DSTypeHTTP,
				QStringHandling: "ds0qstringhandling",
			},
			MSOAlgorithm: ParentConfigDSParamDefaultMSOAlgorithm,
		},
	}

	serverInfo := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "myCDN",
		CDNID:                         43,
		DomainName:                    "serverdomain.example.net",
		HostName:                      "myserver",
		ID:                            44,
		IP:                            "192.168.2.1",
		ParentCacheGroupID:            InvalidID,
		ParentCacheGroupType:          "myParentCGType",
		ProfileID:                     46,
		ProfileName:                   "MyProfileName",
		Port:                          80,
		SecondaryParentCacheGroupID:   InvalidID,
		SecondaryParentCacheGroupType: "MySecondaryParentCGType",
		Type:                          "EDGE",
	}

	serverParams := map[string]string{
		ParentConfigParamQStringHandling: "myQStringHandlingParam",
		ParentConfigParamAlgorithm:       enum.AlgorithmConsistentHash,
		ParentConfigParamQString:         "myQstringParam",
	}

	parentInfos := map[OriginHost][]ParentInfo{
		"ds0.example.net": []ParentInfo{
			ParentInfo{
				Host:            "my-parent-0",
				Port:            80,
				Domain:          "my-parent-0-domain",
				Weight:          "1",
				UseIP:           false,
				Rank:            1,
				IP:              "192.168.2.2",
				PrimaryParent:   true,
				SecondaryParent: false,
			},
			ParentInfo{
				Host:            "my-parent-1",
				Port:            80,
				Domain:          "my-parent-1-domain",
				Weight:          "1",
				UseIP:           false,
				Rank:            1,
				IP:              "192.168.2.3",
				PrimaryParent:   false,
				SecondaryParent: true,
			},
		},
	}

	if !serverInfo.IsTopLevelCache() {
		t.Fatal("server should have been top level, was not; cannot test MSO Secondary Parent")
	}

	txt := MakeParentDotConfig(serverInfo, atsMajorVer, toolName, toURL, parentConfigDSes, serverParams, parentInfos)

	testComment(t, txt, serverName, toolName, toURL)

	txt = strings.Replace(txt, " ", "", -1)

	if !strings.Contains(txt, `secondary_parent="my-parent-1.my-parent-1-domain`) {
		t.Errorf("expected secondary parent 'my-parent-1.my-parent-1-domain', actual: '%v'", txt)
	}
}
