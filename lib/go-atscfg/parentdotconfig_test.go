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

func TestMakeParentDotConfig(t *testing.T) {
	atsMajorVer := 7
	serverName := "myserver"
	toolName := "myToolName"
	toURL := "https://myto.example.net"

	parentConfigDSes := []ParentConfigDSTopLevel{
		ParentConfigDSTopLevel{
			ParentConfigDS: ParentConfigDS{
				Name:            "ds0",
				QStringIgnore:   tc.QStringIgnoreUseInCacheKeyAndPassUp,
				OriginFQDN:      "http://ds0.example.net",
				MultiSiteOrigin: false,
				Type:            tc.DSTypeHTTP,
				QStringHandling: "ds0qstringhandling",
			},
		},
		ParentConfigDSTopLevel{
			ParentConfigDS: ParentConfigDS{
				Name:            "ds1",
				QStringIgnore:   tc.QStringIgnoreDrop,
				OriginFQDN:      "http://ds1.example.net",
				MultiSiteOrigin: false,
				Type:            tc.DSTypeDNS,
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
		ParentConfigParamAlgorithm:       tc.AlgorithmConsistentHash,
		ParentConfigParamQString:         "myQstringParam",
	}

	parentInfos := map[string][]ParentInfo{
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
