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
			Name:       ParentConfigRetryKeysDefault.Algorithm,
			ConfigFile: "parent.config",
			Value:      tc.AlgorithmConsistentHash,
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigRetryKeysDefault.ParentRetry,
			ConfigFile: "parent.config",
			Value:      "both",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigRetryKeysDefault.MaxSimpleRetries,
			ConfigFile: "parent.config",
			Value:      "14",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigRetryKeysDefault.MaxUnavailableRetries,
			ConfigFile: "parent.config",
			Value:      "9",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigRetryKeysDefault.UnavailableRetryResponses,
			ConfigFile: "parent.config",
			Value:      `"400,503"`,
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
		"response_codes:[404]",
		"markdown_codes:[400,503]",
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

	// checks for properly-formed merge keys
	if strings.Contains(cfg.Text, "<< ") {
		t.Errorf("expected yaml merge keys to be '<<: ', actual malformed '<< ': %v", cfg.Text)
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
			Name:       ParentConfigRetryKeysDefault.Algorithm,
			ConfigFile: "parent.config",
			Value:      tc.AlgorithmConsistentHash,
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigRetryKeysDefault.ParentRetry,
			ConfigFile: "parent.config",
			Value:      "both",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigRetryKeysDefault.MaxSimpleRetries,
			ConfigFile: "parent.config",
			Value:      "14",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigRetryKeysDefault.MaxUnavailableRetries,
			ConfigFile: "parent.config",
			Value:      "9",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigRetryKeysDefault.UnavailableRetryResponses,
			ConfigFile: "parent.config",
			Value:      `"400,503"`,
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
			"response_codes:[404]",
			"markdown_codes:[400,503]",
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
			"response_codes:[404]",
			"markdown_codes:[400,503]",
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
			"response_codes:[404]",
			"markdown_codes:[400,503]",
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
			Name:       ParentConfigRetryKeysDefault.Algorithm,
			ConfigFile: "parent.config",
			Value:      tc.AlgorithmConsistentHash,
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigRetryKeysDefault.ParentRetry,
			ConfigFile: "parent.config",
			Value:      "both",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigRetryKeysDefault.MaxSimpleRetries,
			ConfigFile: "parent.config",
			Value:      "14",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigRetryKeysDefault.MaxUnavailableRetries,
			ConfigFile: "parent.config",
			Value:      "9",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigRetryKeysDefault.UnavailableRetryResponses,
			ConfigFile: "parent.config",
			Value:      `"400,503"`,
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
			"response_codes:[404]",
			"markdown_codes:[400,503]",
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
			"response_codes:[404]",
			"markdown_codes:[400,503]",
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
			"response_codes:[404]",
			"markdown_codes:[400,503]",
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

func TestMakeStrategiesDotYAMLFirstLastNoTopoParams(t *testing.T) {
	opt := &StrategiesYAMLOpts{VerboseComments: false, HdrComment: "myHeaderComment"}

	// Non Toplogy
	ds0 := makeParentDS()
	ds0.ID = util.IntPtr(42)
	ds0Type := tc.DSTypeDNS
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")
	ds0.ProfileID = util.IntPtr(310)
	ds0.ProfileName = util.StrPtr("ds0Profile")

	// Non Toplogy, MSO
	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.ProfileID = util.IntPtr(310)
	ds1.ProfileName = util.StrPtr("ds0Profile")
	ds1.MultiSiteOrigin = util.BoolPtr(true)

	dsesall := []DeliveryService{*ds0, *ds1}

	parentConfigParams := []tc.Parameter{
		{
			Name:       ParentConfigParamQStringHandling,
			ConfigFile: "parent.config",
			Value:      "myQStringHandlingParam",
			Profiles:   []byte(`["serverprofile"]`),
		},
		{
			Name:       ParentConfigRetryKeysDefault.Algorithm,
			ConfigFile: "parent.config",
			Value:      tc.AlgorithmConsistentHash,
			Profiles:   []byte(`["serverprofile"]`),
		},
		{
			Name:       ParentConfigParamQString,
			ConfigFile: "parent.config",
			Value:      "myQstringParam",
			Profiles:   []byte(`["serverprofile"]`),
		},
	}

	// Create set of DS params
	params := map[string]string{
		ParentConfigRetryKeysDefault.Algorithm: "strict",
		ParentConfigRetryKeysMSO.Algorithm:     "strict",
		ParentConfigRetryKeysFirst.Algorithm:   "true",
		ParentConfigRetryKeysInner.Algorithm:   "latched",
		ParentConfigRetryKeysLast.Algorithm:    "true",

		ParentConfigRetryKeysDefault.SecondaryMode: "exhaust",
		ParentConfigRetryKeysMSO.SecondaryMode:     "exhaust",
		ParentConfigRetryKeysFirst.SecondaryMode:   "alternate",
		ParentConfigRetryKeysInner.SecondaryMode:   "peering",
		ParentConfigRetryKeysLast.SecondaryMode:    "alternate",

		ParentConfigRetryKeysMSO.ParentRetry:   "unavailable_server_retry",
		ParentConfigRetryKeysFirst.ParentRetry: "both",
		ParentConfigRetryKeysInner.ParentRetry: "both",
		ParentConfigRetryKeysLast.ParentRetry:  "both",

		ParentConfigRetryKeysDefault.MaxSimpleRetries: "11",
		ParentConfigRetryKeysMSO.MaxSimpleRetries:     "11",
		ParentConfigRetryKeysFirst.MaxSimpleRetries:   "12",
		ParentConfigRetryKeysInner.MaxSimpleRetries:   "13",
		ParentConfigRetryKeysLast.MaxSimpleRetries:    "14",

		ParentConfigRetryKeysDefault.SimpleRetryResponses: `"401"`,
		ParentConfigRetryKeysMSO.SimpleRetryResponses:     `"401"`,
		ParentConfigRetryKeysFirst.SimpleRetryResponses:   `"401,402"`,
		ParentConfigRetryKeysInner.SimpleRetryResponses:   `"401,403"`,
		ParentConfigRetryKeysLast.SimpleRetryResponses:    `"401,404"`,

		ParentConfigRetryKeysDefault.MaxUnavailableRetries: "21",
		ParentConfigRetryKeysMSO.MaxUnavailableRetries:     "21",
		ParentConfigRetryKeysFirst.MaxUnavailableRetries:   "22",
		ParentConfigRetryKeysInner.MaxUnavailableRetries:   "23",
		ParentConfigRetryKeysLast.MaxUnavailableRetries:    "24",

		ParentConfigRetryKeysDefault.UnavailableRetryResponses: `"501"`,
		ParentConfigRetryKeysMSO.UnavailableRetryResponses:     `"501"`,
		ParentConfigRetryKeysFirst.UnavailableRetryResponses:   `"501,502"`,
		ParentConfigRetryKeysInner.UnavailableRetryResponses:   `"501,503"`,
		ParentConfigRetryKeysLast.UnavailableRetryResponses:    `"501,504"`,
	}

	// Assign them to the profile
	for key, val := range params {
		tcparam := tc.Parameter{
			Name:       key,
			ConfigFile: "parent.config",
			Value:      val,
			Profiles:   []byte(`["ds0Profile"]`),
		}
		parentConfigParams = append(parentConfigParams, tcparam)
	}

	serverParams := []tc.Parameter{
		{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
	}

	edge := makeTestParentServer()
	edge.Cachegroup = util.StrPtr("edgeCG")
	edge.CachegroupID = util.IntPtr(400)

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG0")
	mid0.CachegroupID = util.IntPtr(500)
	mid0.HostName = util.StrPtr("mymid0")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")
	mid0.Type = tc.CacheGroupMidTypeName
	mid0.TypeID = util.IntPtr(990)

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG1")
	mid1.CachegroupID = util.IntPtr(501)
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")
	mid1.Type = tc.CacheGroupMidTypeName
	mid1.TypeID = util.IntPtr(990)

	org0 := makeTestParentServer()
	org0.Cachegroup = util.StrPtr("orgCG0")
	org0.CachegroupID = util.IntPtr(502)
	org0.HostName = util.StrPtr("myorg0")
	org0.ID = util.IntPtr(48)
	setIP(org0, "192.168.2.4")
	org0.Type = tc.OriginTypeName
	org0.TypeID = util.IntPtr(991)

	org1 := makeTestParentServer()
	org1.Cachegroup = util.StrPtr("orgCG1")
	org1.CachegroupID = util.IntPtr(503)
	org1.HostName = util.StrPtr("myorg1")
	org1.ID = util.IntPtr(49)
	setIP(org1, "192.168.2.5")
	org1.Type = tc.OriginTypeName
	org1.TypeID = util.IntPtr(991)

	servers := []Server{*edge, *mid0, *mid1, *org0, *org1}

	topologies := []tc.Topology{
		{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				{
					Cachegroup: "edgeCG",
					Parents:    []int{1, 2},
				},
				{
					Cachegroup: "midCG0",
					Parents:    []int{3, 4},
				},
				{
					Cachegroup: "midCG1",
					Parents:    []int{3, 4},
				},
				{
					Cachegroup: "orgCG0",
				},
				{
					Cachegroup: "orgCG1",
				},
			},
		},
	}

	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	eCG := &tc.CacheGroupNullable{}
	eCG.Name = edge.Cachegroup
	eCG.ID = edge.CachegroupID
	eCG.ParentName = mid0.Cachegroup
	eCG.ParentCachegroupID = mid0.CachegroupID
	eCG.SecondaryParentName = mid1.Cachegroup
	eCG.SecondaryParentCachegroupID = mid1.CachegroupID
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	mCG0 := &tc.CacheGroupNullable{}
	mCG0.Name = mid0.Cachegroup
	mCG0.ID = mid0.CachegroupID
	mCG0.ParentName = org0.Cachegroup
	mCG0.ParentCachegroupID = org0.CachegroupID
	mCG0.SecondaryParentName = org1.Cachegroup
	mCG0.SecondaryParentCachegroupID = org1.CachegroupID
	mCGType0 := tc.CacheGroupMidTypeName
	mCG0.Type = &mCGType0

	mCG1 := &tc.CacheGroupNullable{}
	mCG1.Name = mid1.Cachegroup
	mCG1.ID = mid1.CachegroupID
	mCG1.ParentName = org1.Cachegroup
	mCG1.ParentCachegroupID = org1.CachegroupID
	mCG1.SecondaryParentName = org0.Cachegroup
	mCG1.SecondaryParentCachegroupID = org0.CachegroupID
	mCGType1 := tc.CacheGroupMidTypeName
	mCG1.Type = &mCGType1

	oCG0 := &tc.CacheGroupNullable{}
	oCG0.Name = org0.Cachegroup
	oCG0.ID = org0.CachegroupID
	oCGType0 := tc.CacheGroupOriginTypeName
	oCG0.Type = &oCGType0

	oCG1 := &tc.CacheGroupNullable{}
	oCG1.Name = org1.Cachegroup
	oCG1.ID = org1.CachegroupID
	oCGType1 := tc.CacheGroupOriginTypeName
	oCG1.Type = &oCGType1

	cgs := []tc.CacheGroupNullable{*eCG, *mCG0, *mCG1, *oCG0, *oCG1}

	dss := []DeliveryServiceServer{
		{Server: *edge.ID, DeliveryService: *ds0.ID},
		{Server: *mid0.ID, DeliveryService: *ds0.ID},
		{Server: *mid1.ID, DeliveryService: *ds0.ID},
		{Server: *org0.ID, DeliveryService: *ds0.ID},
		{Server: *org1.ID, DeliveryService: *ds0.ID},

		{Server: *edge.ID, DeliveryService: *ds1.ID},
		{Server: *mid0.ID, DeliveryService: *ds1.ID},
		{Server: *mid1.ID, DeliveryService: *ds1.ID},
		{Server: *org0.ID, DeliveryService: *ds1.ID},
		{Server: *org1.ID, DeliveryService: *ds1.ID},
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	// edge config
	/*
		for _, ds := range dsesall {
			dses := []DeliveryService{ds}
			cfg, err := MakeStrategiesDotYAML(dses, edge, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
			if err != nil {
				t.Fatal(err)
			}
			txt := cfg.Text

			testComment(t, txt, opt.HdrComment)

			needs := []string{
				` policy: consistent_hash`,
				` go_direct: false`,
				` max_simple_retries: 12`,
				` max_unavailable_retries: 22`,
				` response_codes: [ 401, 402 ]`,
				` markdown_codes: [ 501, 502 ]`,
			}

			missing := missingFrom(txt, needs)
			if 0 < len(missing) {
				t.Errorf("Missing required string(s) from line: %v\n%v", missing, txt)
			}
		}
	*/

	// test mid config, MS only
	{
		ds := dsesall[1]
		dses := []DeliveryService{ds}
		cfg, err := MakeStrategiesDotYAML(dses, mid0, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, opt.HdrComment)

		needs := []string{
			` policy: rr_ip`,
			` go_direct: true`,
			` max_simple_retries: 14`,
			` max_unavailable_retries: 24`,
			` response_codes: [ 401, 404 ]`,
			` markdown_codes: [ 501, 504 ]`,
		}

		missing := missingFrom(txt, needs)
		if 0 < len(missing) {
			t.Errorf("Missing required string(s) from ds/line: %s/%v\n%v", *ds.XMLID, missing, txt)
		}

		excludes := []string{
			`hash_key`,
		}

		excluding := missingFrom(txt, excludes)
		if 1 != len(excludes) {
			t.Errorf("Excluded required string(s) from ds/line: %s/%v\n%v", *ds.XMLID, excluding, txt)
		}
	}
}

func TestMakeStrategiesDotYamlMSONoTopologyNoMid(t *testing.T) {
	opt := &StrategiesYAMLOpts{VerboseComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")
	ds0.MultiSiteOrigin = util.BoolPtr(true)
	ds0.ProfileName = util.StrPtr("dsprofile")
	dses := []DeliveryService{*ds0}

	parentConfigParams := []tc.Parameter{}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
	}

	edge := makeTestParentServer()

	origin0 := makeTestParentServer()
	origin0.Cachegroup = util.StrPtr("originCG")
	origin0.CachegroupID = util.IntPtr(500)
	origin0.HostName = util.StrPtr("myorigin0")
	origin0.ID = util.IntPtr(45)
	setIP(origin0, "192.168.2.2")
	origin0.Type = tc.OriginTypeName
	origin0.TypeID = util.IntPtr(991)

	servers := []Server{*edge, *origin0}

	topologies := []tc.Topology{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	eCG := &tc.CacheGroupNullable{}
	eCG.Name = edge.Cachegroup
	eCG.ID = edge.CachegroupID
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
		{
			Server:          *edge.ID,
			DeliveryService: *ds0.ID,
		},
		{
			Server:          *origin0.ID,
			DeliveryService: *ds0.ID,
		},
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	cfg, err := MakeStrategiesDotYAML(dses, edge, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, opt.HdrComment)

	txtx := strings.Replace(txt, " ", "", -1)

	needs := []string{
		`host:myorigin0.mydomain.example.net`,
	}

	missing := missingFrom(txtx, needs)
	if 0 < len(missing) {
		t.Errorf("Missing required string(s) from ds/line: %s/%v\n%v", *ds0.XMLID, missing, txt)
	}
}

// Test for mso non topology where mid cache group has no primary/secondary
// parents assigned, just any arbitrary servers.
func TestMakeStrategiesDotYamlMSONoTopoMultiCG(t *testing.T) {
	opt := &StrategiesYAMLOpts{VerboseComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")
	ds0.MultiSiteOrigin = util.BoolPtr(true)

	dses := []DeliveryService{*ds0}

	parentConfigParams := []tc.Parameter{}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
	}

	edge := makeTestParentServer()
	edge.Cachegroup = util.StrPtr("edgeCG")
	edge.CachegroupID = util.IntPtr(400)

	mid := makeTestParentServer()
	mid.Cachegroup = util.StrPtr("midCG")
	mid.CachegroupID = util.IntPtr(500)
	mid.HostName = util.StrPtr("mid0")
	mid.ID = util.IntPtr(45)
	setIP(mid, "192.168.2.2")

	org0 := makeTestParentServer()
	org0.Cachegroup = util.StrPtr("orgCG0")
	org0.CachegroupID = util.IntPtr(501)
	org0.HostName = util.StrPtr("org0")
	org0.ID = util.IntPtr(46)
	setIP(org0, "192.168.2.3")
	org0.Type = tc.OriginTypeName
	org0.TypeID = util.IntPtr(991)

	org1 := makeTestParentServer()
	org1.Cachegroup = util.StrPtr("orgCG1")
	org1.CachegroupID = util.IntPtr(502)
	org1.HostName = util.StrPtr("org1")
	org1.ID = util.IntPtr(47)
	setIP(org1, "192.168.2.4")
	org1.Type = tc.OriginTypeName
	org1.TypeID = util.IntPtr(991)

	servers := []Server{*edge, *mid, *org0, *org1}

	topologies := []tc.Topology{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	eCG := &tc.CacheGroupNullable{}
	eCG.Name = edge.Cachegroup
	eCG.ID = edge.CachegroupID
	eCG.ParentName = mid.Cachegroup
	eCG.ParentCachegroupID = mid.CachegroupID
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	// NOTE: no parent cache groups specified
	mCG := &tc.CacheGroupNullable{}
	mCG.Name = mid.Cachegroup
	mCG.ID = mid.CachegroupID
	mCGType := tc.CacheGroupMidTypeName
	mCG.Type = &mCGType

	oCG0 := &tc.CacheGroupNullable{}
	oCG0.Name = org0.Cachegroup
	oCG0.ID = org0.CachegroupID
	oCG0Type := tc.CacheGroupOriginTypeName
	oCG0.Type = &oCG0Type

	oCG1 := &tc.CacheGroupNullable{}
	oCG1.Name = org1.Cachegroup
	oCG1.ID = org1.CachegroupID
	oCG1Type := tc.CacheGroupOriginTypeName
	oCG1.Type = &oCG1Type

	cgs := []tc.CacheGroupNullable{*eCG, *mCG, *oCG0, *oCG1}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *edge.ID,
			DeliveryService: *ds0.ID,
		},
		DeliveryServiceServer{
			Server:          *org0.ID,
			DeliveryService: *ds0.ID,
		},
		DeliveryServiceServer{
			Server:          *org1.ID,
			DeliveryService: *ds0.ID,
		},
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}
	cfg, err := MakeStrategiesDotYAML(dses, mid, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, opt.HdrComment)

	txtx := strings.Replace(txt, " ", "", -1)

	needs := []string{
		`host:org0.mydomain.example.net`,
		`host:org1.mydomain.example.net`,
		`policy:consistent_hash`,
	}

	missing := missingFrom(txtx, needs)
	if 0 < len(missing) {
		t.Errorf("Missing required string(s) from line: %v\n%v", missing, txtx)
	}
}

func TestMakeStrategiesDotYAMLFirstInnerLastParams(t *testing.T) {
	opt := &StrategiesYAMLOpts{VerboseComments: false, HdrComment: "myHeaderComment"}

	// Toplogy ds, MSO
	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")
	ds0.ProfileID = util.IntPtr(311)
	ds0.ProfileName = util.StrPtr("ds0Profile")
	ds0.MultiSiteOrigin = util.BoolPtr(true)
	ds0.Topology = util.StrPtr("t0")

	// Toplogy ds, non MSO
	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.ProfileID = util.IntPtr(311)
	ds1.ProfileName = util.StrPtr("ds0Profile")
	ds1.Topology = util.StrPtr("t0")

	dsesall := []DeliveryService{*ds0, *ds1}

	parentConfigParams := []tc.Parameter{
		{
			Name:       ParentConfigParamQStringHandling,
			ConfigFile: "parent.config",
			Value:      "myQStringHandlingParam",
			Profiles:   []byte(`["serverprofile"]`),
		},
		{
			Name:       ParentConfigRetryKeysDefault.Algorithm,
			ConfigFile: "parent.config",
			Value:      tc.AlgorithmConsistentHash,
			Profiles:   []byte(`["serverprofile"]`),
		},
		{
			Name:       ParentConfigParamQString,
			ConfigFile: "parent.config",
			Value:      "myQstringParam",
			Profiles:   []byte(`["serverprofile"]`),
		},
	}

	// Create set of DS params
	params := map[string]string{
		ParentConfigRetryKeysDefault.Algorithm: "strict",
		ParentConfigRetryKeysMSO.Algorithm:     "strict",
		ParentConfigRetryKeysFirst.Algorithm:   "true",
		ParentConfigRetryKeysInner.Algorithm:   "latched",
		ParentConfigRetryKeysLast.Algorithm:    "true",

		ParentConfigRetryKeysDefault.SecondaryMode: "exhaust",
		ParentConfigRetryKeysMSO.SecondaryMode:     "exhaust",
		ParentConfigRetryKeysFirst.SecondaryMode:   "alternate",
		ParentConfigRetryKeysInner.SecondaryMode:   "peering",
		ParentConfigRetryKeysLast.SecondaryMode:    "alternate",

		ParentConfigRetryKeysMSO.ParentRetry:   "unavailable_server_retry",
		ParentConfigRetryKeysFirst.ParentRetry: "both",
		ParentConfigRetryKeysInner.ParentRetry: "both",
		ParentConfigRetryKeysLast.ParentRetry:  "both",

		ParentConfigRetryKeysDefault.MaxSimpleRetries: "11",
		ParentConfigRetryKeysMSO.MaxSimpleRetries:     "11",
		ParentConfigRetryKeysFirst.MaxSimpleRetries:   "12",
		ParentConfigRetryKeysInner.MaxSimpleRetries:   "13",
		ParentConfigRetryKeysLast.MaxSimpleRetries:    "14",

		ParentConfigRetryKeysDefault.SimpleRetryResponses: `"401"`,
		ParentConfigRetryKeysMSO.SimpleRetryResponses:     `"401"`,
		ParentConfigRetryKeysFirst.SimpleRetryResponses:   `"401,402"`,
		ParentConfigRetryKeysInner.SimpleRetryResponses:   `"401,403"`,
		ParentConfigRetryKeysLast.SimpleRetryResponses:    `"401,404"`,

		ParentConfigRetryKeysDefault.MaxUnavailableRetries: "21",
		ParentConfigRetryKeysMSO.MaxUnavailableRetries:     "21",
		ParentConfigRetryKeysFirst.MaxUnavailableRetries:   "22",
		ParentConfigRetryKeysInner.MaxUnavailableRetries:   "23",
		ParentConfigRetryKeysLast.MaxUnavailableRetries:    "24",

		ParentConfigRetryKeysDefault.UnavailableRetryResponses: `"501"`,
		ParentConfigRetryKeysMSO.UnavailableRetryResponses:     `"501"`,
		ParentConfigRetryKeysFirst.UnavailableRetryResponses:   `"501,502"`,
		ParentConfigRetryKeysInner.UnavailableRetryResponses:   `"501,503"`,
		ParentConfigRetryKeysLast.UnavailableRetryResponses:    `"501,504"`,
	}

	// Assign them to the profile
	for key, val := range params {
		tcparam := tc.Parameter{
			Name:       key,
			ConfigFile: "parent.config",
			Value:      val,
			Profiles:   []byte(`["ds0Profile"]`),
		}
		parentConfigParams = append(parentConfigParams, tcparam)
	}

	serverParams := []tc.Parameter{
		{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
	}

	edge := makeTestParentServer()
	edge.Cachegroup = util.StrPtr("edgeCG")
	edge.CachegroupID = util.IntPtr(400)

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG0")
	mid0.CachegroupID = util.IntPtr(500)
	mid0.HostName = util.StrPtr("mymid0")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG1")
	mid1.CachegroupID = util.IntPtr(501)
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")

	opl0 := makeTestParentServer()
	opl0.Cachegroup = util.StrPtr("oplCG0")
	opl0.CachegroupID = util.IntPtr(502)
	opl0.HostName = util.StrPtr("myopl0")
	opl0.ID = util.IntPtr(46)
	setIP(opl0, "192.168.2.4")

	opl1 := makeTestParentServer()
	opl1.Cachegroup = util.StrPtr("oplCG1")
	opl1.CachegroupID = util.IntPtr(503)
	opl1.HostName = util.StrPtr("myopl1")
	opl1.ID = util.IntPtr(47)
	setIP(opl1, "192.168.2.5")

	org0 := makeTestParentServer()
	org0.Cachegroup = util.StrPtr("orgCG0")
	org0.CachegroupID = util.IntPtr(504)
	org0.HostName = util.StrPtr("myorg0")
	org0.ID = util.IntPtr(48)
	setIP(org0, "192.168.2.6")
	org0.Type = tc.OriginTypeName
	org0.TypeID = util.IntPtr(991)

	org1 := makeTestParentServer()
	org1.Cachegroup = util.StrPtr("orgCG1")
	org1.CachegroupID = util.IntPtr(505)
	org1.HostName = util.StrPtr("myorg1")
	org1.ID = util.IntPtr(49)
	setIP(org1, "192.168.2.7")
	org1.Type = tc.OriginTypeName
	org1.TypeID = util.IntPtr(991)

	servers := []Server{*edge, *mid0, *mid1, *opl0, *opl1, *org0, *org1}

	topologies := []tc.Topology{
		{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				{
					Cachegroup: "edgeCG",
					Parents:    []int{1, 2},
				},
				{
					Cachegroup: "midCG0",
					Parents:    []int{3, 4},
				},
				{
					Cachegroup: "midCG1",
					Parents:    []int{3, 4},
				},
				{
					Cachegroup: "oplCG0",
					Parents:    []int{5, 6},
				},
				{
					Cachegroup: "oplCG1",
					Parents:    []int{5, 6},
				},
				{
					Cachegroup: "orgCG0",
				},
				{
					Cachegroup: "orgCG1",
				},
			},
		},
	}

	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	eCG := &tc.CacheGroupNullable{}
	eCG.Name = edge.Cachegroup
	eCG.ID = edge.CachegroupID
	eCG.ParentName = mid0.Cachegroup
	eCG.ParentCachegroupID = mid0.CachegroupID
	eCG.SecondaryParentName = mid1.Cachegroup
	eCG.SecondaryParentCachegroupID = mid1.CachegroupID
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	mCG0 := &tc.CacheGroupNullable{}
	mCG0.Name = mid0.Cachegroup
	mCG0.ID = mid0.CachegroupID
	mCG0.ParentName = opl0.Cachegroup
	mCG0.ParentCachegroupID = opl0.CachegroupID
	mCG0.SecondaryParentName = opl1.Cachegroup
	mCG0.SecondaryParentCachegroupID = opl1.CachegroupID
	mCGType0 := tc.CacheGroupMidTypeName
	mCG0.Type = &mCGType0

	mCG1 := &tc.CacheGroupNullable{}
	mCG1.Name = mid1.Cachegroup
	mCG1.ID = mid1.CachegroupID
	mCG1.ParentName = opl1.Cachegroup
	mCG1.ParentCachegroupID = opl1.CachegroupID
	mCG1.SecondaryParentName = opl0.Cachegroup
	mCG1.SecondaryParentCachegroupID = opl0.CachegroupID
	mCGType1 := tc.CacheGroupMidTypeName
	mCG1.Type = &mCGType1

	oplCG0 := &tc.CacheGroupNullable{}
	oplCG0.Name = opl0.Cachegroup
	oplCG0.ID = opl0.CachegroupID
	oplCG0.ParentName = org0.Cachegroup
	oplCG0.ParentCachegroupID = org0.CachegroupID
	oplCG0.SecondaryParentName = org1.Cachegroup
	oplCG0.SecondaryParentCachegroupID = org1.CachegroupID
	oplCGType0 := tc.CacheGroupMidTypeName
	oplCG0.Type = &oplCGType0

	oplCG1 := &tc.CacheGroupNullable{}
	oplCG1.Name = opl1.Cachegroup
	oplCG1.ID = opl1.CachegroupID
	oplCG1.ParentName = org1.Cachegroup
	oplCG1.ParentCachegroupID = org1.CachegroupID
	oplCG1.SecondaryParentName = org0.Cachegroup
	oplCG1.SecondaryParentCachegroupID = org0.CachegroupID
	oplCGType1 := tc.CacheGroupMidTypeName
	oplCG1.Type = &oplCGType1

	oCG0 := &tc.CacheGroupNullable{}
	oCG0.Name = org0.Cachegroup
	oCG0.ID = org0.CachegroupID
	oCGType0 := tc.CacheGroupOriginTypeName
	oCG0.Type = &oCGType0

	oCG1 := &tc.CacheGroupNullable{}
	oCG1.Name = org1.Cachegroup
	oCG1.ID = org1.CachegroupID
	oCGType1 := tc.CacheGroupOriginTypeName
	oCG1.Type = &oCGType1

	cgs := []tc.CacheGroupNullable{*eCG, *mCG0, *mCG1, *oplCG0, *oplCG1, *oCG0, *oCG1}

	dss := []DeliveryServiceServer{
		{Server: *org0.ID, DeliveryService: *ds0.ID},
		{Server: *org1.ID, DeliveryService: *ds0.ID},

		{Server: *edge.ID, DeliveryService: *ds1.ID},
		{Server: *mid0.ID, DeliveryService: *ds1.ID},
		{Server: *mid1.ID, DeliveryService: *ds1.ID},
		{Server: *opl0.ID, DeliveryService: *ds1.ID},
		{Server: *opl1.ID, DeliveryService: *ds1.ID},
		{Server: *org0.ID, DeliveryService: *ds1.ID},
		{Server: *org1.ID, DeliveryService: *ds1.ID},
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	// edge config
	for _, ds := range dsesall {
		dses := []DeliveryService{ds}
		cfg, err := MakeStrategiesDotYAML(dses, edge, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, opt.HdrComment)

		needs := []string{
			` policy: consistent_hash`,
			` hash_key: path`,
			` go_direct: false`,
			` max_simple_retries: 12`,
			` max_unavailable_retries: 22`,
			` response_codes: [ 401, 402 ]`,
			` markdown_codes: [ 501, 502 ]`,
		}

		missing := missingFrom(txt, needs)
		if 0 < len(missing) {
			t.Errorf("Missing required string(s) from line: %v\n%v", missing, txt)
		}
	}

	// test mid config
	for _, ds := range dsesall {
		dses := []DeliveryService{ds}
		cfg, err := MakeStrategiesDotYAML(dses, mid0, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, opt.HdrComment)

		needs := []string{
			` policy: consistent_hash`,
			` hash_key: path`,
			` go_direct: false`,
			` max_simple_retries: 13`,
			` max_unavailable_retries: 23`,
			` response_codes: [ 401, 403 ]`,
			` markdown_codes: [ 501, 503 ]`,
		}

		missing := missingFrom(txt, needs)
		if 0 < len(missing) {
			t.Errorf("Missing required string(s) from ds/line: %s/%v\n%v", *ds.XMLID, missing, txt)
		}
	}

	// test opl config
	for _, ds := range dsesall {
		dses := []DeliveryService{ds}
		cfg, err := MakeStrategiesDotYAML(dses, opl0, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, opt.HdrComment)

		needs := []string{
			` policy: rr_ip`,
			` go_direct: true`,
			` max_simple_retries: 14`,
			` max_unavailable_retries: 24`,
			` response_codes: [ 401, 404 ]`,
			` markdown_codes: [ 501, 504 ]`,
		}

		missing := missingFrom(txt, needs)
		if 0 < len(missing) {
			t.Errorf("Missing required string(s) from line: %v\n%v", missing, txt)
		}

		excludes := []string{
			`hash_key`,
		}

		excluding := missingFrom(txt, excludes)
		if 1 != len(excludes) {
			t.Errorf("Excluded required string(s) from ds/line: %s/%v\n%v", *ds.XMLID, excluding, txt)
		}
	}
}
