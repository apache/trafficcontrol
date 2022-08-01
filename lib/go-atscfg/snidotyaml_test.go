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

func TestMakeSNIDotYAMLH2(t *testing.T) {
	opts := &SNIDotYAMLOpts{VerboseComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.Protocol = util.IntPtr(int(tc.DSProtocolHTTPAndHTTPS))
	ds0.ProfileName = util.StrPtr("ds0profile")
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")
	ds0.TLSVersions = []string{"1.1", "1.2"}

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.Protocol = util.IntPtr(int(tc.DSProtocolHTTPAndHTTPS))
	ds1.RoutingName = util.StrPtr("myroutingname")
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.TLSVersions = []string{"1.1", "1.2"}

	dses := []DeliveryService{*ds0, *ds1}

	parentConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       ParentConfigParamQStringHandling,
			ConfigFile: "parent.config",
			Value:      "myQStringHandlingParam",
			Profiles:   []byte(`["serverprofile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigRetryKeysDefault.Algorithm,
			ConfigFile: "parent.config",
			Value:      tc.AlgorithmConsistentHash,
			Profiles:   []byte(`["serverprofile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamQString,
			ConfigFile: "parent.config",
			Value:      "myQstringParam",
			Profiles:   []byte(`["serverprofile"]`),
		},
		tc.Parameter{
			Name:       SSLServerNameYAMLParamEnableH2,
			ConfigFile: "parent.config",
			Value:      "true",
			Profiles:   []byte(`["ds0profile"]`),
		},
	}

	server := makeTestParentServer()

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

	topologies := []tc.Topology{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	eCG := &tc.CacheGroupNullable{}
	eCG.Name = server.Cachegroup
	eCG.ID = server.CachegroupID
	eCG.ParentName = mid0.Cachegroup
	eCG.ParentCachegroupID = mid0.CachegroupID
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	mCG := &tc.CacheGroupNullable{}
	mCG.Name = mid0.Cachegroup
	mCG.ID = mid0.CachegroupID
	mCGType := tc.CacheGroupMidTypeName
	mCG.Type = &mCGType

	cgs := []tc.CacheGroupNullable{*eCG, *mCG}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds0.ID,
		},
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds1.ID,
		},
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	dsr := []tc.DeliveryServiceRegexes{
		tc.DeliveryServiceRegexes{
			DSName: *ds0.XMLID,
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
		cfg, err := MakeSNIDotYAML(server, dses, dss, dsr, parentConfigParams, cdn, topologies, cgs, serverCapabilities, dsRequiredCapabilities, opts)
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
		parentConfigParams := []tc.Parameter{
			tc.Parameter{
				Name:       ParentConfigParamQStringHandling,
				ConfigFile: "parent.config",
				Value:      "myQStringHandlingParam",
				Profiles:   []byte(`["serverprofile"]`),
			},
			tc.Parameter{
				Name:       ParentConfigRetryKeysDefault.Algorithm,
				ConfigFile: "parent.config",
				Value:      tc.AlgorithmConsistentHash,
				Profiles:   []byte(`["serverprofile"]`),
			},
			tc.Parameter{
				Name:       ParentConfigParamQString,
				ConfigFile: "parent.config",
				Value:      "myQstringParam",
				Profiles:   []byte(`["serverprofile"]`),
			},
			tc.Parameter{
				Name:       SSLServerNameYAMLParamEnableH2,
				ConfigFile: "parent.config",
				Value:      "false",
				Profiles:   []byte(`["ds0profile"]`),
			},
		}

		cfg, err := MakeSNIDotYAML(server, dses, dss, dsr, parentConfigParams, cdn, topologies, cgs, serverCapabilities, dsRequiredCapabilities, opts)
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
		parentConfigParams := []tc.Parameter{
			tc.Parameter{
				Name:       ParentConfigParamQStringHandling,
				ConfigFile: "parent.config",
				Value:      "myQStringHandlingParam",
				Profiles:   []byte(`["serverprofile"]`),
			},
			tc.Parameter{
				Name:       ParentConfigRetryKeysDefault.Algorithm,
				ConfigFile: "parent.config",
				Value:      tc.AlgorithmConsistentHash,
				Profiles:   []byte(`["serverprofile"]`),
			},
			tc.Parameter{
				Name:       ParentConfigParamQString,
				ConfigFile: "parent.config",
				Value:      "myQstringParam",
				Profiles:   []byte(`["serverprofile"]`),
			},
		}

		cfg, err := MakeSNIDotYAML(server, dses, dss, dsr, parentConfigParams, cdn, topologies, cgs, serverCapabilities, dsRequiredCapabilities, opts)
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
		parentConfigParams := []tc.Parameter{
			tc.Parameter{
				Name:       ParentConfigParamQStringHandling,
				ConfigFile: "parent.config",
				Value:      "myQStringHandlingParam",
				Profiles:   []byte(`["serverprofile"]`),
			},
			tc.Parameter{
				Name:       ParentConfigRetryKeysDefault.Algorithm,
				ConfigFile: "parent.config",
				Value:      tc.AlgorithmConsistentHash,
				Profiles:   []byte(`["serverprofile"]`),
			},
			tc.Parameter{
				Name:       ParentConfigParamQString,
				ConfigFile: "parent.config",
				Value:      "myQstringParam",
				Profiles:   []byte(`["serverprofile"]`),
			},
		}

		cfg, err := MakeSNIDotYAML(server, dses, dss, dsr, parentConfigParams, cdn, topologies, cgs, serverCapabilities, dsRequiredCapabilities, opts)
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
