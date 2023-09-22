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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

func TestMakeSNIDotYAMLH2(t *testing.T) {
	opts := &SNIDotYAMLOpts{VerboseComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := "HTTP"
	ds0.Type = &ds0Type
	ds0.Protocol = util.Ptr(int(tc.DSProtocolHTTPAndHTTPS))
	ds0.ProfileName = util.Ptr("ds0profile")
	ds0.QStringIgnore = util.Ptr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.Ptr("http://ds0.example.net")
	ds0.TLSVersions = []string{"1.1", "1.2"}

	ds1 := makeParentDS()
	ds1.ID = util.Ptr(43)
	ds1Type := "DNS"
	ds1.Type = &ds1Type
	ds1.Protocol = util.Ptr(int(tc.DSProtocolHTTPAndHTTPS))
	ds1.RoutingName = "myroutingname"
	ds1.QStringIgnore = util.Ptr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.Ptr("http://ds1.example.net")
	ds1.TLSVersions = []string{"1.1", "1.2"}

	dses := []DeliveryService{*ds0, *ds1}

	parentConfigParams := []tc.ParameterV5{
		tc.ParameterV5{
			Name:       ParentConfigParamQStringHandling,
			ConfigFile: "parent.config",
			Value:      "myQStringHandlingParam",
			Profiles:   []byte(`["serverprofile"]`),
		},
		tc.ParameterV5{
			Name:       ParentConfigRetryKeysDefault.Algorithm,
			ConfigFile: "parent.config",
			Value:      tc.AlgorithmConsistentHash,
			Profiles:   []byte(`["serverprofile"]`),
		},
		tc.ParameterV5{
			Name:       ParentConfigParamQString,
			ConfigFile: "parent.config",
			Value:      "myQstringParam",
			Profiles:   []byte(`["serverprofile"]`),
		},
		tc.ParameterV5{
			Name:       SSLServerNameYAMLParamEnableH2,
			ConfigFile: "parent.config",
			Value:      "true",
			Profiles:   []byte(`["ds0profile"]`),
		},
	}

	server := makeTestParentServer()
	servers := makeTestAnyCastServers()

	mid0 := makeTestParentServer()
	mid0.CacheGroup = "midCG"
	mid0.HostName = "mymid0"
	mid0.ID = 45
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.CacheGroup = "midCG"
	mid1.HostName = "mymid1"
	mid1.ID = 46
	setIP(mid1, "192.168.2.3")

	topologies := []tc.TopologyV5{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	eCG := &tc.CacheGroupNullableV5{}
	eCG.Name = util.Ptr(server.CacheGroup)
	eCG.ID = util.Ptr(server.CacheGroupID)
	eCG.ParentName = util.Ptr(mid0.CacheGroup)
	eCG.ParentCachegroupID = util.Ptr(mid0.CacheGroupID)
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	mCG := &tc.CacheGroupNullableV5{}
	mCG.Name = util.Ptr(mid0.CacheGroup)
	mCG.ID = util.Ptr(mid0.CacheGroupID)
	mCGType := tc.CacheGroupMidTypeName
	mCG.Type = &mCGType

	cgs := []tc.CacheGroupNullableV5{*eCG, *mCG}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          server.ID,
			DeliveryService: *ds0.ID,
		},
		DeliveryServiceServer{
			Server:          server.ID,
			DeliveryService: *ds1.ID,
		},
	}
	cdn := &tc.CDNV5{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	dsr := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: ds0.XMLID,
			Regexes: []tc.DeliveryServiceRegex{
				tc.DeliveryServiceRegex{
					Type:      string(tc.DSMatchTypeHostRegex),
					SetNumber: 0,
					Pattern:   `.*\.ds0\..*`,
				},
			},
		},
	}

	t.Run("sni.yaml http2 param enabled", func(t *testing.T) {
		cfg, err := MakeSNIDotYAML(server, servers, dses, dss, dsr, parentConfigParams, cdn, topologies, cgs, serverCapabilities, dsRequiredCapabilities, opts)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		if !strings.Contains(txt, `fqdn: 'myserver.ds0.cdndomain.example'`) {
			t.Errorf("expected ds0 fqdn, actual ''%+v'' warnings ''%+v''", txt, cfg.Warnings)
		}
		if !strings.Contains(txt, `http2: on`) {
			t.Errorf("expected h2 enabled for ds with parameter, actual ''%+v'' warnings ''%+v''", txt, cfg.Warnings)
		}
		if !strings.Contains(txt, `['TLSv1_1','TLSv1_2']`) {
			t.Errorf("expected TLS 1.1,1.2 for ds with TLSVersions field, actual ''%+v'' warnings ''%+v''", txt, cfg.Warnings)
		}
		if strings.Contains(txt, `TLSv1_3`) {
			t.Errorf("expected no TLS 1.3 for ds with TLSVersions of 1.1,1.2, actual ''%+v'' warnings ''%+v''", txt, cfg.Warnings)
		}
	})

	t.Run("sni.yaml http2 param disabled", func(t *testing.T) {
		parentConfigParams := []tc.ParameterV5{
			tc.ParameterV5{
				Name:       ParentConfigParamQStringHandling,
				ConfigFile: "parent.config",
				Value:      "myQStringHandlingParam",
				Profiles:   []byte(`["serverprofile"]`),
			},
			tc.ParameterV5{
				Name:       ParentConfigRetryKeysDefault.Algorithm,
				ConfigFile: "parent.config",
				Value:      tc.AlgorithmConsistentHash,
				Profiles:   []byte(`["serverprofile"]`),
			},
			tc.ParameterV5{
				Name:       ParentConfigParamQString,
				ConfigFile: "parent.config",
				Value:      "myQstringParam",
				Profiles:   []byte(`["serverprofile"]`),
			},
			tc.ParameterV5{
				Name:       SSLServerNameYAMLParamEnableH2,
				ConfigFile: "parent.config",
				Value:      "false",
				Profiles:   []byte(`["ds0profile"]`),
			},
		}

		cfg, err := MakeSNIDotYAML(server, servers, dses, dss, dsr, parentConfigParams, cdn, topologies, cgs, serverCapabilities, dsRequiredCapabilities, opts)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		if !strings.Contains(txt, `fqdn: 'myserver.ds0.cdndomain.example'`) {
			t.Errorf("expected ds0 fqdn, actual ''%+v'' warnings ''%+v''", txt, cfg.Warnings)
		}
		if !strings.Contains(txt, `http2: off`) {
			t.Errorf("expected h2 enabled for ds with parameter false, actual ''%+v'' warnings ''%+v''", txt, cfg.Warnings)
		}
		if !strings.Contains(txt, `['TLSv1_1','TLSv1_2']`) {
			t.Errorf("expected TLS 1.1,1.2 for ds with TLSVersions field, actual ''%+v'' warnings ''%+v''", txt, cfg.Warnings)
		}
		if strings.Contains(txt, `TLSv1_3`) {
			t.Errorf("expected no TLS 1.3 for ds with TLSVersions of 1.1,1.2, actual ''%+v'' warnings ''%+v''", txt, cfg.Warnings)
		}
	})

	t.Run("sni.yaml http2 param missing default disabled", func(t *testing.T) {
		opts := &SNIDotYAMLOpts{
			VerboseComments: false,
			HdrComment:      "myHeaderComment",
			DefaultEnableH2: false,
		}
		parentConfigParams := []tc.ParameterV5{
			tc.ParameterV5{
				Name:       ParentConfigParamQStringHandling,
				ConfigFile: "parent.config",
				Value:      "myQStringHandlingParam",
				Profiles:   []byte(`["serverprofile"]`),
			},
			tc.ParameterV5{
				Name:       ParentConfigRetryKeysDefault.Algorithm,
				ConfigFile: "parent.config",
				Value:      tc.AlgorithmConsistentHash,
				Profiles:   []byte(`["serverprofile"]`),
			},
			tc.ParameterV5{
				Name:       ParentConfigParamQString,
				ConfigFile: "parent.config",
				Value:      "myQstringParam",
				Profiles:   []byte(`["serverprofile"]`),
			},
		}

		cfg, err := MakeSNIDotYAML(server, servers, dses, dss, dsr, parentConfigParams, cdn, topologies, cgs, serverCapabilities, dsRequiredCapabilities, opts)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		if !strings.Contains(txt, `fqdn: 'myserver.ds0.cdndomain.example'`) {
			t.Errorf("expected ds0 fqdn, actual ''%+v'' warnings ''%+v''", txt, cfg.Warnings)
		}
		if !strings.Contains(txt, `http2: off`) {
			t.Errorf("expected h2 disabled for ds with no parameter and cfg default disabled, actual ''%+v'' warnings ''%+v''", txt, cfg.Warnings)
		}
		if !strings.Contains(txt, `['TLSv1_1','TLSv1_2']`) {
			t.Errorf("expected TLS 1.1,1.2 for ds with TLSVersions field, actual ''%+v'' warnings ''%+v''", txt, cfg.Warnings)
		}
		if strings.Contains(txt, `TLSv1_3`) {
			t.Errorf("expected no TLS 1.3 for ds with TLSVersions of 1.1,1.2, actual ''%+v'' warnings ''%+v''", txt, cfg.Warnings)
		}
	})

	t.Run("sni.yaml http2 param missing default enabled", func(t *testing.T) {
		opts := &SNIDotYAMLOpts{
			VerboseComments: false,
			HdrComment:      "myHeaderComment",
			DefaultEnableH2: true,
		}
		parentConfigParams := []tc.ParameterV5{
			tc.ParameterV5{
				Name:       ParentConfigParamQStringHandling,
				ConfigFile: "parent.config",
				Value:      "myQStringHandlingParam",
				Profiles:   []byte(`["serverprofile"]`),
			},
			tc.ParameterV5{
				Name:       ParentConfigRetryKeysDefault.Algorithm,
				ConfigFile: "parent.config",
				Value:      tc.AlgorithmConsistentHash,
				Profiles:   []byte(`["serverprofile"]`),
			},
			tc.ParameterV5{
				Name:       ParentConfigParamQString,
				ConfigFile: "parent.config",
				Value:      "myQstringParam",
				Profiles:   []byte(`["serverprofile"]`),
			},
		}

		cfg, err := MakeSNIDotYAML(server, servers, dses, dss, dsr, parentConfigParams, cdn, topologies, cgs, serverCapabilities, dsRequiredCapabilities, opts)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		if !strings.Contains(txt, `fqdn: 'myserver.ds0.cdndomain.example'`) {
			t.Errorf("expected ds0 fqdn, actual ''%+v'' warnings ''%+v''", txt, cfg.Warnings)
		}
		if !strings.Contains(txt, `http2: on`) {
			t.Errorf("expected h2 enabled for ds with no parameter and cfg default enabled, actual ''%+v'' warnings ''%+v''", txt, cfg.Warnings)
		}
		if !strings.Contains(txt, `['TLSv1_1','TLSv1_2']`) {
			t.Errorf("expected TLS 1.1,1.2 for ds with TLSVersions field, actual ''%+v'' warnings ''%+v''", txt, cfg.Warnings)
		}
		if strings.Contains(txt, `TLSv1_3`) {
			t.Errorf("expected no TLS 1.3 for ds with TLSVersions of 1.1,1.2, actual ''%+v'' warnings ''%+v''", txt, cfg.Warnings)
		}
	})

}
