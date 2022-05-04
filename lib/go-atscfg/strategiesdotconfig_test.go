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

func TestMakeStrategiesDotConfig(t *testing.T) {
	opt := &StrategiesYAMLOpts{VerboseComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0.XMLID = util.StrPtr("ds0")
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")

	ds1 := makeParentDS()
	ds1.XMLID = util.StrPtr("ds1")
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")

	dses := []DeliveryService{*ds0, *ds1}

	parentConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       ParentConfigParamQStringHandling,
			ConfigFile: "parent.config",
			Value:      "myQStringHandlingParam",
			Profiles:   []byte(`["serverprofile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamAlgorithm,
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

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
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

	servers := []Server{*server, *mid0, *mid1}

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

	cfg, err := MakeStrategiesDotYAML(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, opt.HdrComment)

	expecteds := []string{
		"strategy:'strategy-ds0'",
		"strategy:'strategy-ds1'",
		"host__ds0__parent__mymid0-dot-mydomain-dot-example-dot-net__80\nhost:mymid0.mydomain.example.net",
		"host__ds0__parent__mymid1-dot-mydomain-dot-example-dot-net__80\nhost:mymid1.mydomain.example.net",
		"host__ds1__parent__mymid0-dot-mydomain-dot-example-dot-net__80\nhost:mymid0.mydomain.example.net",
		"host__ds1__parent__mymid1-dot-mydomain-dot-example-dot-net__80\nhost:mymid1.mydomain.example.net",
	}
	txt = strings.Replace(txt, " ", "", -1)
	for _, expected := range expecteds {
		if !strings.Contains(txt, expected) {
			t.Errorf("expected parent '''%v''', actual: '''%v'''", expected, txt)
		}
	}
	if !warningsContains(cfg.Warnings, "myQStringHandlingParam") {
		t.Errorf("expected malformed qstring 'myQstringParam' in warnings, actual: '%v' val '%v'", cfg.Warnings, txt)
	}
}

func TestMakeStrategiesTopologiesParams(t *testing.T) {
	opt := &StrategiesYAMLOpts{VerboseComments: false, HdrComment: "myHeaderComment"}

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.Topology = util.StrPtr("t0")
	ds1.ProfileName = util.StrPtr("ds1Profile")
	ds1.ProfileID = util.IntPtr(994)
	ds1.MultiSiteOrigin = util.BoolPtr(true)

	dses := []DeliveryService{*ds1}

	parentConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       ParentConfigParamQStringHandling,
			ConfigFile: "parent.config",
			Value:      "myQStringHandlingParam",
			Profiles:   []byte(`["serverprofile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamAlgorithm,
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
			Name:       ParentConfigParamAlgorithm,
			ConfigFile: "parent.config",
			Value:      "consistent_hash",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamParentRetry,
			ConfigFile: "parent.config",
			Value:      "both",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamUnavailableServerRetryResponses,
			ConfigFile: "parent.config",
			Value:      `"400,503"`,
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamMaxSimpleRetries,
			ConfigFile: "parent.config",
			Value:      "14",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamMaxUnavailableServerRetries,
			ConfigFile: "parent.config",
			Value:      "9",
			Profiles:   []byte(`["ds1Profile"]`),
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "8",
			Profiles:   []byte(`["global"]`),
		},
	}

	server := makeTestParentServer()
	server.Cachegroup = util.StrPtr("edgeCG")
	server.CachegroupID = util.IntPtr(400)

	origin0 := makeTestParentServer()
	origin0.Cachegroup = util.StrPtr("originCG")
	origin0.CachegroupID = util.IntPtr(500)
	origin0.HostName = util.StrPtr("myorigin0")
	origin0.ID = util.IntPtr(45)
	setIP(origin0, "192.168.2.2")
	origin0.Type = tc.OriginTypeName
	origin0.TypeID = util.IntPtr(991)

	origin1 := makeTestParentServer()
	origin1.Cachegroup = util.StrPtr("originCG")
	origin1.CachegroupID = util.IntPtr(500)
	origin1.HostName = util.StrPtr("myorigin1")
	origin1.ID = util.IntPtr(46)
	setIP(origin1, "192.168.2.3")
	origin1.Type = tc.OriginTypeName
	origin1.TypeID = util.IntPtr(991)

	servers := []Server{*server, *origin0, *origin1}

	topologies := []tc.Topology{
		tc.Topology{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				tc.TopologyNode{
					Cachegroup: "edgeCG",
					Parents:    []int{1},
				},
				tc.TopologyNode{
					Cachegroup: "originCG",
				},
			},
		},
	}

	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	eCG := &tc.CacheGroupNullable{}
	eCG.Name = server.Cachegroup
	eCG.ID = server.CachegroupID
	eCG.ParentName = origin0.Cachegroup
	eCG.ParentCachegroupID = origin0.CachegroupID
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	oCG := &tc.CacheGroupNullable{}
	oCG.Name = origin0.Cachegroup
	oCG.ID = origin0.CachegroupID
	oCGType := tc.CacheGroupOriginTypeName
	oCG.Type = &oCGType

	cgs := []tc.CacheGroupNullable{*eCG, *oCG}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *origin0.ID,
			DeliveryService: *ds1.ID,
		},
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	cfg, err := MakeStrategiesDotYAML(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, opt.HdrComment)

	txt = strings.Replace(txt, " ", "", -1)

	expecteds := []string{
		"strategy:'strategy-ds1'",
		"max_simple_retries:14",
		"max_unavailable_retries:9",
		"response_codes:\n-404",
		"markdown_codes:\n-400\n-503",
		"host__ds1__parent__myorigin0-dot-mydomain-dot-example-dot-net__80\nhost:myorigin0.mydomain.example.net",
	}

	for _, expected := range expecteds {
		if !strings.Contains(txt, expected) {
			t.Errorf("expected parent '''%v''', actual: '''%v'''", expected, txt)
		}
	}
}

func TestMakeStrategiesHTTPSOrigin(t *testing.T) {
	opt := &StrategiesYAMLOpts{VerboseComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0.XMLID = util.StrPtr("ds0")
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("https://ds0.example.net")

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")

	dses := []DeliveryService{*ds0, *ds1}

	parentConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       ParentConfigParamQStringHandling,
			ConfigFile: "parent.config",
			Value:      "myQStringHandlingParam",
			Profiles:   []byte(`["serverprofile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamAlgorithm,
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

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
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

	servers := []Server{*server, *mid0, *mid1}

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

	cfg, err := MakeStrategiesDotYAML(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, opt.HdrComment)

	txt = strings.Replace(txt, " ", "", -1)

	// this is an Edge, and all traffic is internal to a Mid
	// So even though the Origin is HTTPS, all traffic should be 80=HTTP

	if !strings.Contains(txt, "protocol:\n-port:80") {
		t.Errorf("expected edge parent.config of https origin to use internal http port 80 (not https/443), actual: '%v'", txt)
	}
	if strings.Contains(txt, "port: 443") {
		t.Errorf("expected edge parent.config of https origin to use internal http port 80 and not https/443, actual: '%v'", txt)
	}
}

func TestMakeStrategiesPeeringRing(t *testing.T) {
	opt := &StrategiesYAMLOpts{VerboseComments: false, HdrComment: "myHeaderComment"}

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.Topology = util.StrPtr("t0")
	ds1.ProfileName = util.StrPtr("ds1Profile")
	ds1.ProfileID = util.IntPtr(994)
	ds1.MultiSiteOrigin = util.BoolPtr(false)

	dses := []DeliveryService{*ds1}

	parentConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       ParentConfigParamQStringHandling,
			ConfigFile: "parent.config",
			Value:      "myQStringHandlingParam",
			Profiles:   []byte(`["serverprofile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamAlgorithm,
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
			Name:       ParentConfigParamAlgorithm,
			ConfigFile: "parent.config",
			Value:      "consistent_hash",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamParentRetry,
			ConfigFile: "parent.config",
			Value:      "both",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamUnavailableServerRetryResponses,
			ConfigFile: "parent.config",
			Value:      `"400,503"`,
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamMaxSimpleRetries,
			ConfigFile: "parent.config",
			Value:      "14",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamMaxUnavailableServerRetries,
			ConfigFile: "parent.config",
			Value:      "9",
			Profiles:   []byte(`["ds1Profile"]`),
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "8",
			Profiles:   []byte(`["global"]`),
		},
	}

	edge0 := makeTestParentServer()
	edge0.ID = util.IntPtr(12)
	edge0.HostName = util.StrPtr("edge0")
	edge0.Cachegroup = util.StrPtr("edgeCG")
	edge0.CachegroupID = util.IntPtr(400)

	edge1 := makeTestParentServer()
	edge1.ID = util.IntPtr(13)
	edge1.HostName = util.StrPtr("edge1")
	edge1.Cachegroup = util.StrPtr("edgeCG")
	edge1.CachegroupID = util.IntPtr(400)

	origin0 := makeTestParentServer()
	origin0.Cachegroup = util.StrPtr("originCG")
	origin0.CachegroupID = util.IntPtr(500)
	origin0.HostName = util.StrPtr("myorigin0")
	origin0.ID = util.IntPtr(45)
	setIP(origin0, "192.168.2.2")
	origin0.Type = tc.OriginTypeName
	origin0.TypeID = util.IntPtr(991)

	origin1 := makeTestParentServer()
	origin1.Cachegroup = util.StrPtr("originCG")
	origin1.CachegroupID = util.IntPtr(500)
	origin1.HostName = util.StrPtr("myorigin1")
	origin1.ID = util.IntPtr(46)
	setIP(origin1, "192.168.2.3")
	origin1.Type = tc.OriginTypeName
	origin1.TypeID = util.IntPtr(991)

	servers := []Server{*edge0, *edge1, *origin0, *origin1}

	topologies := []tc.Topology{
		tc.Topology{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				tc.TopologyNode{
					Cachegroup: "edgeCG",
					Parents:    []int{1},
				},
				tc.TopologyNode{
					Cachegroup: "originCG",
				},
			},
		},
	}

	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	eCG := &tc.CacheGroupNullable{}
	eCG.Name = edge0.Cachegroup
	eCG.ID = edge0.CachegroupID
	eCG.ParentName = origin0.Cachegroup
	eCG.ParentCachegroupID = origin0.CachegroupID
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	oCG := &tc.CacheGroupNullable{}
	oCG.Name = origin0.Cachegroup
	oCG.ID = origin0.CachegroupID
	oCGType := tc.CacheGroupOriginTypeName
	oCG.Type = &oCGType

	cgs := []tc.CacheGroupNullable{*eCG, *oCG}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *origin0.ID,
			DeliveryService: *ds1.ID,
		},
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	t.Run("peering ring true", func(t *testing.T) {
		parentConfigParamsPR := make([]tc.Parameter, len(parentConfigParams), len(parentConfigParams))
		copy(parentConfigParamsPR, parentConfigParams)
		parentConfigParamsPR = append(parentConfigParamsPR, tc.Parameter{
			Name:       StrategyConfigUsePeering,
			ConfigFile: "parent.config",
			Value:      "true",
			Profiles:   []byte(`["ds1Profile"]`),
		})

		cfg, err := MakeStrategiesDotYAML(dses, edge0, servers, topologies, serverParams, parentConfigParamsPR, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, opt.HdrComment)

		txt = strings.Replace(txt, " ", "", -1)

		expecteds := []string{
			"strategy:'strategy-ds1'",
			"max_simple_retries:14",
			"max_unavailable_retries:9",
			"response_codes:\n-404",
			"markdown_codes:\n-400\n-503",
			"host__ds1__parent__ds1-dot-example-dot-net__80\nhost:ds1.example.net",
			"groups:\n-*peers_group\n-*group_parents_ds1", // peer ring group before parent group, param 'true'
		}

		for _, expected := range expecteds {
			if !strings.Contains(txt, expected) {
				t.Errorf("expected parent '''%v''', actual: '''%v'''", expected, txt)
			}
		}
	})

	t.Run("peering ring false", func(t *testing.T) {
		parentConfigParamsPR := make([]tc.Parameter, len(parentConfigParams), len(parentConfigParams))
		copy(parentConfigParamsPR, parentConfigParams)
		parentConfigParamsPR = append(parentConfigParamsPR, tc.Parameter{
			Name:       StrategyConfigUsePeering,
			ConfigFile: "parent.config",
			Value:      "false",
			Profiles:   []byte(`["ds1Profile"]`),
		})

		cfg, err := MakeStrategiesDotYAML(dses, edge0, servers, topologies, serverParams, parentConfigParamsPR, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, opt.HdrComment)

		txt = strings.Replace(txt, " ", "", -1)

		expecteds := []string{
			"strategy:'strategy-ds1'",
			"max_simple_retries:14",
			"max_unavailable_retries:9",
			"response_codes:\n-404",
			"markdown_codes:\n-400\n-503",
			"host__ds1__parent__ds1-dot-example-dot-net__80\nhost:ds1.example.net",
			"groups:\n-*group_parents_ds1\nfailover:", // no peer ring group, param is not 'true'
		}

		for _, expected := range expecteds {
			if !strings.Contains(txt, expected) {
				t.Errorf("expected parent '''%v''', actual: '''%v'''", expected, txt)
			}
		}
	})

	t.Run("peering ring nonexistent", func(t *testing.T) {
		cfg, err := MakeStrategiesDotYAML(dses, edge0, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, opt.HdrComment)

		txt = strings.Replace(txt, " ", "", -1)

		expecteds := []string{
			"strategy:'strategy-ds1'",
			"max_simple_retries:14",
			"max_unavailable_retries:9",
			"response_codes:\n-404",
			"markdown_codes:\n-400\n-503",
			"host__ds1__parent__ds1-dot-example-dot-net__80\nhost:ds1.example.net",
			"groups:\n-*group_parents_ds1\nfailover:", // no peer ring group, no parameter
		}

		for _, expected := range expecteds {
			if !strings.Contains(txt, expected) {
				t.Errorf("expected parent '''%v''', actual: '''%v'''", expected, txt)
			}
		}
	})
}

func TestMakeStrategiesPeeringRingMSO(t *testing.T) {
	opt := &StrategiesYAMLOpts{VerboseComments: false, HdrComment: "myHeaderComment"}

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.Topology = util.StrPtr("t0")
	ds1.ProfileName = util.StrPtr("ds1Profile")
	ds1.ProfileID = util.IntPtr(994)
	ds1.MultiSiteOrigin = util.BoolPtr(true)

	dses := []DeliveryService{*ds1}

	parentConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       ParentConfigParamQStringHandling,
			ConfigFile: "parent.config",
			Value:      "myQStringHandlingParam",
			Profiles:   []byte(`["serverprofile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamAlgorithm,
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
			Name:       ParentConfigParamAlgorithm,
			ConfigFile: "parent.config",
			Value:      "consistent_hash",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamParentRetry,
			ConfigFile: "parent.config",
			Value:      "both",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamUnavailableServerRetryResponses,
			ConfigFile: "parent.config",
			Value:      `"400,503"`,
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamMaxSimpleRetries,
			ConfigFile: "parent.config",
			Value:      "14",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamMaxUnavailableServerRetries,
			ConfigFile: "parent.config",
			Value:      "9",
			Profiles:   []byte(`["ds1Profile"]`),
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "8",
			Profiles:   []byte(`["global"]`),
		},
	}

	edge0 := makeTestParentServer()
	edge0.ID = util.IntPtr(12)
	edge0.HostName = util.StrPtr("edge0")
	edge0.Cachegroup = util.StrPtr("edgeCG")
	edge0.CachegroupID = util.IntPtr(400)

	edge1 := makeTestParentServer()
	edge1.ID = util.IntPtr(13)
	edge1.HostName = util.StrPtr("edge1")
	edge1.Cachegroup = util.StrPtr("edgeCG")
	edge1.CachegroupID = util.IntPtr(400)

	origin0 := makeTestParentServer()
	origin0.Cachegroup = util.StrPtr("originCG")
	origin0.CachegroupID = util.IntPtr(500)
	origin0.HostName = util.StrPtr("myorigin0")
	origin0.ID = util.IntPtr(45)
	setIP(origin0, "192.168.2.2")
	origin0.Type = tc.OriginTypeName
	origin0.TypeID = util.IntPtr(991)

	origin1 := makeTestParentServer()
	origin1.Cachegroup = util.StrPtr("originCG")
	origin1.CachegroupID = util.IntPtr(500)
	origin1.HostName = util.StrPtr("myorigin1")
	origin1.ID = util.IntPtr(46)
	setIP(origin1, "192.168.2.3")
	origin1.Type = tc.OriginTypeName
	origin1.TypeID = util.IntPtr(991)

	servers := []Server{*edge0, *edge1, *origin0, *origin1}

	topologies := []tc.Topology{
		tc.Topology{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				tc.TopologyNode{
					Cachegroup: "edgeCG",
					Parents:    []int{1},
				},
				tc.TopologyNode{
					Cachegroup: "originCG",
				},
			},
		},
	}

	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	eCG := &tc.CacheGroupNullable{}
	eCG.Name = edge0.Cachegroup
	eCG.ID = edge0.CachegroupID
	eCG.ParentName = origin0.Cachegroup
	eCG.ParentCachegroupID = origin0.CachegroupID
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	oCG := &tc.CacheGroupNullable{}
	oCG.Name = origin0.Cachegroup
	oCG.ID = origin0.CachegroupID
	oCGType := tc.CacheGroupOriginTypeName
	oCG.Type = &oCGType

	cgs := []tc.CacheGroupNullable{*eCG, *oCG}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *origin0.ID,
			DeliveryService: *ds1.ID,
		},
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	t.Run("peering ring true", func(t *testing.T) {
		parentConfigParamsPR := make([]tc.Parameter, len(parentConfigParams), len(parentConfigParams))
		copy(parentConfigParamsPR, parentConfigParams)
		parentConfigParamsPR = append(parentConfigParamsPR, tc.Parameter{
			Name:       StrategyConfigUsePeering,
			ConfigFile: "parent.config",
			Value:      "true",
			Profiles:   []byte(`["ds1Profile"]`),
		})

		cfg, err := MakeStrategiesDotYAML(dses, edge0, servers, topologies, serverParams, parentConfigParamsPR, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, opt.HdrComment)

		txt = strings.Replace(txt, " ", "", -1)

		expecteds := []string{
			"strategy:'strategy-ds1'",
			"max_simple_retries:14",
			"max_unavailable_retries:9",
			"response_codes:\n-404",
			"markdown_codes:\n-400\n-503",
			"host__ds1__parent__myorigin0-dot-mydomain-dot-example-dot-net__80\nhost:myorigin0.mydomain.example.net",
			"groups:\n-*peers_group\n-*group_parents_ds1", // peer ring group before parent group, param 'true'
		}

		for _, expected := range expecteds {
			if !strings.Contains(txt, expected) {
				t.Errorf("expected parent '''%v''', actual: '''%v'''", expected, txt)
			}
		}
	})

	t.Run("peering ring false", func(t *testing.T) {
		parentConfigParamsPR := make([]tc.Parameter, len(parentConfigParams), len(parentConfigParams))
		copy(parentConfigParamsPR, parentConfigParams)
		parentConfigParamsPR = append(parentConfigParamsPR, tc.Parameter{
			Name:       StrategyConfigUsePeering,
			ConfigFile: "parent.config",
			Value:      "false",
			Profiles:   []byte(`["ds1Profile"]`),
		})

		cfg, err := MakeStrategiesDotYAML(dses, edge0, servers, topologies, serverParams, parentConfigParamsPR, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, opt.HdrComment)

		txt = strings.Replace(txt, " ", "", -1)

		expecteds := []string{
			"strategy:'strategy-ds1'",
			"max_simple_retries:14",
			"max_unavailable_retries:9",
			"response_codes:\n-404",
			"markdown_codes:\n-400\n-503",
			"host__ds1__parent__myorigin0-dot-mydomain-dot-example-dot-net__80\nhost:myorigin0.mydomain.example.net",
			"groups:\n-*group_parents_ds1\nfailover:", // no peer ring group, param is not 'true'
		}

		for _, expected := range expecteds {
			if !strings.Contains(txt, expected) {
				t.Errorf("expected parent '''%v''', actual: '''%v'''", expected, txt)
			}
		}
	})

	t.Run("peering ring nonexistent", func(t *testing.T) {
		cfg, err := MakeStrategiesDotYAML(dses, edge0, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, opt.HdrComment)

		txt = strings.Replace(txt, " ", "", -1)

		expecteds := []string{
			"strategy:'strategy-ds1'",
			"max_simple_retries:14",
			"max_unavailable_retries:9",
			"response_codes:\n-404",
			"markdown_codes:\n-400\n-503",
			"host__ds1__parent__myorigin0-dot-mydomain-dot-example-dot-net__80\nhost:myorigin0.mydomain.example.net",
			"groups:\n-*group_parents_ds1\nfailover:", // no peer ring group, no parameter
		}

		for _, expected := range expecteds {
			if !strings.Contains(txt, expected) {
				t.Errorf("expected parent '''%v''', actual: '''%v'''", expected, txt)
			}
		}
	})
}

func TestMakeStrategiesPeeringRingNonTopology(t *testing.T) {
	opt := &StrategiesYAMLOpts{VerboseComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0.XMLID = util.StrPtr("ds0")
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")

	ds1 := makeParentDS()
	ds1.XMLID = util.StrPtr("ds1")
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.ProfileName = util.StrPtr("ds1Profile")

	dses := []DeliveryService{*ds0, *ds1}

	parentConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       ParentConfigParamQStringHandling,
			ConfigFile: "parent.config",
			Value:      "myQStringHandlingParam",
			Profiles:   []byte(`["serverprofile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigParamAlgorithm,
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

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
	}

	edge0 := makeTestParentServer()
	edge0.ID = util.IntPtr(12)
	edge0.HostName = util.StrPtr("edge0")
	edge0.Cachegroup = util.StrPtr("edgeCG")
	edge0.CachegroupID = util.IntPtr(400)

	edge1 := makeTestParentServer()
	edge1.ID = util.IntPtr(13)
	edge1.HostName = util.StrPtr("edge1")
	edge1.Cachegroup = util.StrPtr("edgeCG")
	edge1.CachegroupID = util.IntPtr(400)

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

	servers := []Server{*edge0, *edge1, *mid0, *mid1}

	topologies := []tc.Topology{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	eCG := &tc.CacheGroupNullable{}
	eCG.Name = edge0.Cachegroup
	eCG.ID = edge0.CachegroupID
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
			Server:          *edge0.ID,
			DeliveryService: *ds0.ID,
		},
		DeliveryServiceServer{
			Server:          *edge0.ID,
			DeliveryService: *ds1.ID,
		},
		DeliveryServiceServer{
			Server:          *edge1.ID,
			DeliveryService: *ds0.ID,
		},
		DeliveryServiceServer{
			Server:          *edge1.ID,
			DeliveryService: *ds1.ID,
		},
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	t.Run("peering ring true", func(t *testing.T) {
		parentConfigParamsPR := make([]tc.Parameter, len(parentConfigParams), len(parentConfigParams))
		copy(parentConfigParamsPR, parentConfigParams)
		parentConfigParamsPR = append(parentConfigParamsPR, tc.Parameter{
			Name:       StrategyConfigUsePeering,
			ConfigFile: "parent.config",
			Value:      "true",
			Profiles:   []byte(`["ds1Profile"]`),
		})

		cfg, err := MakeStrategiesDotYAML(dses, edge0, servers, topologies, serverParams, parentConfigParamsPR, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, opt.HdrComment)

		expecteds := []string{
			"strategy:'strategy-ds0'",
			"strategy:'strategy-ds1'",
			"host__ds0__parent__mymid0-dot-mydomain-dot-example-dot-net__80\nhost:mymid0.mydomain.example.net",
			"host__ds0__parent__mymid1-dot-mydomain-dot-example-dot-net__80\nhost:mymid1.mydomain.example.net",
			"host__ds1__parent__mymid0-dot-mydomain-dot-example-dot-net__80\nhost:mymid0.mydomain.example.net",
			"host__ds1__parent__mymid1-dot-mydomain-dot-example-dot-net__80\nhost:mymid1.mydomain.example.net",
			"groups:\n-*peers_group\n-*group_parents_ds1", // peer ring group before parent group, param 'true'
		}
		txt = strings.Replace(txt, " ", "", -1)
		for _, expected := range expecteds {
			if !strings.Contains(txt, expected) {
				t.Errorf("expected parent '''%v''', actual: '''%v'''", expected, txt)
			}
		}
		if !warningsContains(cfg.Warnings, "myQStringHandlingParam") {
			t.Errorf("expected malformed qstring 'myQstringParam' in warnings, actual: '%v' val '%v'", cfg.Warnings, txt)
		}
	})
	t.Run("peering ring false", func(t *testing.T) {
		parentConfigParamsPR := make([]tc.Parameter, len(parentConfigParams), len(parentConfigParams))
		copy(parentConfigParamsPR, parentConfigParams)
		parentConfigParamsPR = append(parentConfigParamsPR, tc.Parameter{
			Name:       StrategyConfigUsePeering,
			ConfigFile: "parent.config",
			Value:      "false",
			Profiles:   []byte(`["ds1Profile"]`),
		})

		cfg, err := MakeStrategiesDotYAML(dses, edge0, servers, topologies, serverParams, parentConfigParamsPR, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, opt.HdrComment)

		expecteds := []string{
			"strategy:'strategy-ds0'",
			"strategy:'strategy-ds1'",
			"host__ds0__parent__mymid0-dot-mydomain-dot-example-dot-net__80\nhost:mymid0.mydomain.example.net",
			"host__ds0__parent__mymid1-dot-mydomain-dot-example-dot-net__80\nhost:mymid1.mydomain.example.net",
			"host__ds1__parent__mymid0-dot-mydomain-dot-example-dot-net__80\nhost:mymid0.mydomain.example.net",
			"host__ds1__parent__mymid1-dot-mydomain-dot-example-dot-net__80\nhost:mymid1.mydomain.example.net",
			"groups:\n-*group_parents_ds1", // peer ring group before parent group, param 'true'
		}
		txt = strings.Replace(txt, " ", "", -1)
		for _, expected := range expecteds {
			if !strings.Contains(txt, expected) {
				t.Errorf("expected parent '''%v''', actual: '''%v'''", expected, txt)
			}
		}
		if !warningsContains(cfg.Warnings, "myQStringHandlingParam") {
			t.Errorf("expected malformed qstring 'myQstringParam' in warnings, actual: '%v' val '%v'", cfg.Warnings, txt)
		}
	})
	t.Run("peering ring nonexistent", func(t *testing.T) {
		cfg, err := MakeStrategiesDotYAML(dses, edge0, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, opt.HdrComment)

		expecteds := []string{
			"strategy:'strategy-ds0'",
			"strategy:'strategy-ds1'",
			"host__ds0__parent__mymid0-dot-mydomain-dot-example-dot-net__80\nhost:mymid0.mydomain.example.net",
			"host__ds0__parent__mymid1-dot-mydomain-dot-example-dot-net__80\nhost:mymid1.mydomain.example.net",
			"host__ds1__parent__mymid0-dot-mydomain-dot-example-dot-net__80\nhost:mymid0.mydomain.example.net",
			"host__ds1__parent__mymid1-dot-mydomain-dot-example-dot-net__80\nhost:mymid1.mydomain.example.net",
			"groups:\n-*group_parents_ds1", // peer ring group before parent group, param 'true'
		}
		txt = strings.Replace(txt, " ", "", -1)
		for _, expected := range expecteds {
			if !strings.Contains(txt, expected) {
				t.Errorf("expected parent '''%v''', actual: '''%v'''", expected, txt)
			}
		}
		if !warningsContains(cfg.Warnings, "myQStringHandlingParam") {
			t.Errorf("expected malformed qstring 'myQstringParam' in warnings, actual: '%v' val '%v'", cfg.Warnings, txt)
		}
	})
}
