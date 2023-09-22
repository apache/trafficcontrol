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

func TestMakeStrategiesDotConfig(t *testing.T) {
	opt := &StrategiesYAMLOpts{VerboseComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0.XMLID = "ds0"
	ds0Type := "HTTP"
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.Ptr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.Ptr("http://ds0.example.net")

	ds1 := makeParentDS()
	ds1.XMLID = "ds1"
	ds1.ID = util.Ptr(43)
	ds1Type := "DNS"
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.Ptr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.Ptr("http://ds1.example.net")

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
	}

	serverParams := []tc.ParameterV5{
		tc.ParameterV5{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
	}

	server := makeTestParentServer()

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

	servers := []Server{*server, *mid0, *mid1}

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
	ds1.ID = util.Ptr(43)
	ds1Type := "DNS"
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.Ptr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.Ptr("http://ds1.example.net")
	ds1.Topology = util.Ptr("t0")
	ds1.ProfileName = util.Ptr("ds1Profile")
	ds1.ProfileID = util.Ptr(994)
	ds1.MultiSiteOrigin = true

	dses := []DeliveryService{*ds1}

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
			Name:       ParentConfigRetryKeysDefault.Algorithm,
			ConfigFile: "parent.config",
			Value:      tc.AlgorithmConsistentHash,
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.ParameterV5{
			Name:       ParentConfigRetryKeysDefault.ParentRetry,
			ConfigFile: "parent.config",
			Value:      "both",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.ParameterV5{
			Name:       ParentConfigRetryKeysDefault.MaxSimpleRetries,
			ConfigFile: "parent.config",
			Value:      "14",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.ParameterV5{
			Name:       ParentConfigRetryKeysDefault.MaxUnavailableRetries,
			ConfigFile: "parent.config",
			Value:      "9",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.ParameterV5{
			Name:       ParentConfigRetryKeysDefault.UnavailableRetryResponses,
			ConfigFile: "parent.config",
			Value:      `"400,503"`,
			Profiles:   []byte(`["ds1Profile"]`),
		},
	}

	serverParams := []tc.ParameterV5{
		tc.ParameterV5{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "8",
			Profiles:   []byte(`["global"]`),
		},
	}

	server := makeTestParentServer()
	server.CacheGroup = "edgeCG"
	server.CacheGroupID = 400

	origin0 := makeTestParentServer()
	origin0.CacheGroup = "originCG"
	origin0.CacheGroupID = 500
	origin0.HostName = "myorigin0"
	origin0.ID = 45
	setIP(origin0, "192.168.2.2")
	origin0.Type = tc.OriginTypeName
	origin0.TypeID = 991

	origin1 := makeTestParentServer()
	origin1.CacheGroup = "originCG"
	origin1.CacheGroupID = 500
	origin1.HostName = "myorigin1"
	origin1.ID = 46
	setIP(origin1, "192.168.2.3")
	origin1.Type = tc.OriginTypeName
	origin1.TypeID = 991

	servers := []Server{*server, *origin0, *origin1}

	topologies := []tc.TopologyV5{
		tc.TopologyV5{
			Name: "t0",
			Nodes: []tc.TopologyNodeV5{
				tc.TopologyNodeV5{
					Cachegroup: "edgeCG",
					Parents:    []int{1},
				},
				tc.TopologyNodeV5{
					Cachegroup: "originCG",
				},
			},
		},
	}

	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	eCG := &tc.CacheGroupNullableV5{}
	eCG.Name = util.Ptr(server.CacheGroup)
	eCG.ID = util.Ptr(server.CacheGroupID)
	eCG.ParentName = util.Ptr(origin0.CacheGroup)
	eCG.ParentCachegroupID = util.Ptr(origin0.CacheGroupID)
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	oCG := &tc.CacheGroupNullableV5{}
	oCG.Name = util.Ptr(origin0.CacheGroup)
	oCG.ID = util.Ptr(origin0.CacheGroupID)
	oCGType := tc.CacheGroupOriginTypeName
	oCG.Type = &oCGType

	cgs := []tc.CacheGroupNullableV5{*eCG, *oCG}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          origin0.ID,
			DeliveryService: *ds1.ID,
		},
	}
	cdn := &tc.CDNV5{
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
	ds0.XMLID = "ds0"
	ds0Type := "HTTP"
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.Ptr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.Ptr("https://ds0.example.net")

	ds1 := makeParentDS()
	ds1.ID = util.Ptr(43)
	ds1Type := "DNS"
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.Ptr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.Ptr("http://ds1.example.net")

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
	}

	serverParams := []tc.ParameterV5{
		tc.ParameterV5{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
	}

	server := makeTestParentServer()

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

	servers := []Server{*server, *mid0, *mid1}

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
	ds1.ID = util.Ptr(43)
	ds1Type := "DNS"
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.Ptr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.Ptr("http://ds1.example.net")
	ds1.Topology = util.Ptr("t0")
	ds1.ProfileName = util.Ptr("ds1Profile")
	ds1.ProfileID = util.Ptr(994)
	ds1.MultiSiteOrigin = false

	dses := []DeliveryService{*ds1}

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
			Name:       ParentConfigRetryKeysDefault.Algorithm,
			ConfigFile: "parent.config",
			Value:      tc.AlgorithmConsistentHash,
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.ParameterV5{
			Name:       ParentConfigRetryKeysDefault.ParentRetry,
			ConfigFile: "parent.config",
			Value:      "both",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.ParameterV5{
			Name:       ParentConfigRetryKeysDefault.MaxSimpleRetries,
			ConfigFile: "parent.config",
			Value:      "14",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.ParameterV5{
			Name:       ParentConfigRetryKeysDefault.MaxUnavailableRetries,
			ConfigFile: "parent.config",
			Value:      "9",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.ParameterV5{
			Name:       ParentConfigRetryKeysDefault.UnavailableRetryResponses,
			ConfigFile: "parent.config",
			Value:      `"400,503"`,
			Profiles:   []byte(`["ds1Profile"]`),
		},
	}

	serverParams := []tc.ParameterV5{
		tc.ParameterV5{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "8",
			Profiles:   []byte(`["global"]`),
		},
	}

	edge0 := makeTestParentServer()
	edge0.ID = 12
	edge0.HostName = "edge0"
	edge0.CacheGroup = "edgeCG"
	edge0.CacheGroupID = 400

	edge1 := makeTestParentServer()
	edge1.ID = 13
	edge1.HostName = "edge1"
	edge1.CacheGroup = "edgeCG"
	edge1.CacheGroupID = 400

	origin0 := makeTestParentServer()
	origin0.CacheGroup = "originCG"
	origin0.CacheGroupID = 500
	origin0.HostName = "myorigin0"
	origin0.ID = 45
	setIP(origin0, "192.168.2.2")
	origin0.Type = tc.OriginTypeName
	origin0.TypeID = 991

	origin1 := makeTestParentServer()
	origin1.CacheGroup = "originCG"
	origin1.CacheGroupID = 500
	origin1.HostName = "myorigin1"
	origin1.ID = 46
	setIP(origin1, "192.168.2.3")
	origin1.Type = tc.OriginTypeName
	origin1.TypeID = 991

	servers := []Server{*edge0, *edge1, *origin0, *origin1}

	topologies := []tc.TopologyV5{
		tc.TopologyV5{
			Name: "t0",
			Nodes: []tc.TopologyNodeV5{
				tc.TopologyNodeV5{
					Cachegroup: "edgeCG",
					Parents:    []int{1},
				},
				tc.TopologyNodeV5{
					Cachegroup: "originCG",
				},
			},
		},
	}

	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	eCG := &tc.CacheGroupNullableV5{}
	eCG.Name = util.Ptr(edge0.CacheGroup)
	eCG.ID = util.Ptr(edge0.CacheGroupID)
	eCG.ParentName = util.Ptr(origin0.CacheGroup)
	eCG.ParentCachegroupID = util.Ptr(origin0.CacheGroupID)
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	oCG := &tc.CacheGroupNullableV5{}
	oCG.Name = util.Ptr(origin0.CacheGroup)
	oCG.ID = util.Ptr(origin0.CacheGroupID)
	oCGType := tc.CacheGroupOriginTypeName
	oCG.Type = &oCGType

	cgs := []tc.CacheGroupNullableV5{*eCG, *oCG}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          origin0.ID,
			DeliveryService: *ds1.ID,
		},
	}
	cdn := &tc.CDNV5{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	t.Run("peering ring true", func(t *testing.T) {
		parentConfigParamsPR := make([]tc.ParameterV5, len(parentConfigParams), len(parentConfigParams))
		copy(parentConfigParamsPR, parentConfigParams)
		parentConfigParamsPR = append(parentConfigParamsPR, tc.ParameterV5{
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
		parentConfigParamsPR := make([]tc.ParameterV5, len(parentConfigParams), len(parentConfigParams))
		copy(parentConfigParamsPR, parentConfigParams)
		parentConfigParamsPR = append(parentConfigParamsPR, tc.ParameterV5{
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
	ds1.ID = util.Ptr(43)
	ds1Type := "DNS"
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.Ptr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.Ptr("http://ds1.example.net")
	ds1.Topology = util.Ptr("t0")
	ds1.ProfileName = util.Ptr("ds1Profile")
	ds1.ProfileID = util.Ptr(994)
	ds1.MultiSiteOrigin = true

	dses := []DeliveryService{*ds1}

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
			Name:       ParentConfigRetryKeysDefault.Algorithm,
			ConfigFile: "parent.config",
			Value:      tc.AlgorithmConsistentHash,
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.ParameterV5{
			Name:       ParentConfigRetryKeysDefault.ParentRetry,
			ConfigFile: "parent.config",
			Value:      "both",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.ParameterV5{
			Name:       ParentConfigRetryKeysDefault.MaxSimpleRetries,
			ConfigFile: "parent.config",
			Value:      "14",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.ParameterV5{
			Name:       ParentConfigRetryKeysDefault.MaxUnavailableRetries,
			ConfigFile: "parent.config",
			Value:      "9",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.ParameterV5{
			Name:       ParentConfigRetryKeysDefault.UnavailableRetryResponses,
			ConfigFile: "parent.config",
			Value:      `"400,503"`,
			Profiles:   []byte(`["ds1Profile"]`),
		},
	}

	serverParams := []tc.ParameterV5{
		tc.ParameterV5{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "8",
			Profiles:   []byte(`["global"]`),
		},
	}

	edge0 := makeTestParentServer()
	edge0.ID = 12
	edge0.HostName = "edge0"
	edge0.CacheGroup = "edgeCG"
	edge0.CacheGroupID = 400

	edge1 := makeTestParentServer()
	edge1.ID = 13
	edge1.HostName = "edge1"
	edge1.CacheGroup = "edgeCG"
	edge1.CacheGroupID = 400

	origin0 := makeTestParentServer()
	origin0.CacheGroup = "originCG"
	origin0.CacheGroupID = 500
	origin0.HostName = "myorigin0"
	origin0.ID = 45
	setIP(origin0, "192.168.2.2")
	origin0.Type = tc.OriginTypeName
	origin0.TypeID = 991

	origin1 := makeTestParentServer()
	origin1.CacheGroup = "originCG"
	origin1.CacheGroupID = 500
	origin1.HostName = "myorigin1"
	origin1.ID = 46
	setIP(origin1, "192.168.2.3")
	origin1.Type = tc.OriginTypeName
	origin1.TypeID = 991

	servers := []Server{*edge0, *edge1, *origin0, *origin1}

	topologies := []tc.TopologyV5{
		tc.TopologyV5{
			Name: "t0",
			Nodes: []tc.TopologyNodeV5{
				tc.TopologyNodeV5{
					Cachegroup: "edgeCG",
					Parents:    []int{1},
				},
				tc.TopologyNodeV5{
					Cachegroup: "originCG",
				},
			},
		},
	}

	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	eCG := &tc.CacheGroupNullableV5{}
	eCG.Name = util.Ptr(edge0.CacheGroup)
	eCG.ID = util.Ptr(edge0.CacheGroupID)
	eCG.ParentName = util.Ptr(origin0.CacheGroup)
	eCG.ParentCachegroupID = util.Ptr(origin0.CacheGroupID)
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	oCG := &tc.CacheGroupNullableV5{}
	oCG.Name = util.Ptr(origin0.CacheGroup)
	oCG.ID = util.Ptr(origin0.CacheGroupID)
	oCGType := tc.CacheGroupOriginTypeName
	oCG.Type = &oCGType

	cgs := []tc.CacheGroupNullableV5{*eCG, *oCG}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          origin0.ID,
			DeliveryService: *ds1.ID,
		},
	}
	cdn := &tc.CDNV5{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	t.Run("peering ring true", func(t *testing.T) {
		parentConfigParamsPR := make([]tc.ParameterV5, len(parentConfigParams), len(parentConfigParams))
		copy(parentConfigParamsPR, parentConfigParams)
		parentConfigParamsPR = append(parentConfigParamsPR, tc.ParameterV5{
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
		parentConfigParamsPR := make([]tc.ParameterV5, len(parentConfigParams), len(parentConfigParams))
		copy(parentConfigParamsPR, parentConfigParams)
		parentConfigParamsPR = append(parentConfigParamsPR, tc.ParameterV5{
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
	ds0.XMLID = "ds0"
	ds0Type := "HTTP"
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.Ptr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.Ptr("http://ds0.example.net")

	ds1 := makeParentDS()
	ds1.XMLID = "ds1"
	ds1.ID = util.Ptr(43)
	ds1Type := "DNS"
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.Ptr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.Ptr("http://ds1.example.net")
	ds1.ProfileName = util.Ptr("ds1Profile")

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
	}

	serverParams := []tc.ParameterV5{
		tc.ParameterV5{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
	}

	edge0 := makeTestParentServer()
	edge0.ID = 12
	edge0.HostName = "edge0"
	edge0.CacheGroup = "edgeCG"
	edge0.CacheGroupID = 400

	edge1 := makeTestParentServer()
	edge1.ID = 13
	edge1.HostName = "edge1"
	edge1.CacheGroup = "edgeCG"
	edge1.CacheGroupID = 400

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

	servers := []Server{*edge0, *edge1, *mid0, *mid1}

	topologies := []tc.TopologyV5{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	eCG := &tc.CacheGroupNullableV5{}
	eCG.Name = util.Ptr(edge0.CacheGroup)
	eCG.ID = util.Ptr(edge0.CacheGroupID)
	eCG.ParentName = util.Ptr(mid0.CacheGroup)
	eCG.ParentCachegroupID = util.Ptr(mid0.CacheGroupID)
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	mCG := &tc.CacheGroupNullableV5{}
	mCG.Name = util.Ptr(mid0.CacheGroup)
	mCG.ID = util.Ptr(mid0.CacheGroupID)
	mCGType := util.Ptr(tc.CacheGroupMidTypeName)
	mCG.Type = mCGType

	cgs := []tc.CacheGroupNullableV5{*eCG, *mCG}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          edge0.ID,
			DeliveryService: *ds0.ID,
		},
		DeliveryServiceServer{
			Server:          edge0.ID,
			DeliveryService: *ds1.ID,
		},
		DeliveryServiceServer{
			Server:          edge1.ID,
			DeliveryService: *ds0.ID,
		},
		DeliveryServiceServer{
			Server:          edge1.ID,
			DeliveryService: *ds1.ID,
		},
	}
	cdn := &tc.CDNV5{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	t.Run("peering ring true", func(t *testing.T) {
		parentConfigParamsPR := make([]tc.ParameterV5, len(parentConfigParams), len(parentConfigParams))
		copy(parentConfigParamsPR, parentConfigParams)
		parentConfigParamsPR = append(parentConfigParamsPR, tc.ParameterV5{
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
		parentConfigParamsPR := make([]tc.ParameterV5, len(parentConfigParams), len(parentConfigParams))
		copy(parentConfigParamsPR, parentConfigParams)
		parentConfigParamsPR = append(parentConfigParamsPR, tc.ParameterV5{
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
	opt := &StrategiesYAMLOpts{VerboseComments: false, HdrComment: "myHeaderComment", GoDirect: "true"}

	// Non Toplogy
	ds0 := makeParentDS()
	ds0.ID = util.Ptr(42)
	ds0Type := "DNS"
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.Ptr(int(tc.QStringIgnoreDrop))
	ds0.OrgServerFQDN = util.Ptr("http://ds0.example.net")
	ds0.ProfileID = util.Ptr(310)
	ds0.ProfileName = util.Ptr("ds0Profile")

	// Non Toplogy, MSO
	ds1 := makeParentDS()
	ds1.ID = util.Ptr(43)
	ds1Type := "DNS"
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.Ptr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.Ptr("http://ds1.example.net")
	ds1.ProfileID = util.Ptr(310)
	ds1.ProfileName = util.Ptr("ds0Profile")
	ds1.MultiSiteOrigin = true

	dsesall := []DeliveryService{*ds0, *ds1}

	parentConfigParams := []tc.ParameterV5{
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
		tcparam := tc.ParameterV5{
			Name:       key,
			ConfigFile: "parent.config",
			Value:      val,
			Profiles:   []byte(`["ds0Profile"]`),
		}
		parentConfigParams = append(parentConfigParams, tcparam)
	}

	serverParams := []tc.ParameterV5{
		{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
	}

	edge := makeTestParentServer()
	edge.CacheGroup = "edgeCG"
	edge.CacheGroupID = 400

	mid0 := makeTestParentServer()
	mid0.CacheGroup = "midCG0"
	mid0.CacheGroupID = 500
	mid0.HostName = "mymid0"
	mid0.ID = 45
	setIP(mid0, "192.168.2.2")
	mid0.Type = tc.CacheGroupMidTypeName
	mid0.TypeID = 990

	mid1 := makeTestParentServer()
	mid1.CacheGroup = "midCG1"
	mid1.CacheGroupID = 501
	mid1.HostName = "mymid1"
	mid1.ID = 46
	setIP(mid1, "192.168.2.3")
	mid1.Type = tc.CacheGroupMidTypeName
	mid1.TypeID = 990

	org0 := makeTestParentServer()
	org0.CacheGroup = "orgCG0"
	org0.CacheGroupID = 502
	org0.HostName = "myorg0"
	org0.ID = 48
	setIP(org0, "192.168.2.4")
	org0.Type = tc.OriginTypeName
	org0.TypeID = 991

	org1 := makeTestParentServer()
	org1.CacheGroup = "orgCG1"
	org1.CacheGroupID = 503
	org1.HostName = "myorg1"
	org1.ID = 49
	setIP(org1, "192.168.2.5")
	org1.Type = tc.OriginTypeName
	org1.TypeID = 991

	servers := []Server{*edge, *mid0, *mid1, *org0, *org1}

	topologies := []tc.TopologyV5{
		{
			Name: "t0",
			Nodes: []tc.TopologyNodeV5{
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

	eCG := &tc.CacheGroupNullableV5{}
	eCG.Name = util.Ptr(edge.CacheGroup)
	eCG.ID = util.Ptr(edge.CacheGroupID)
	eCG.ParentName = util.Ptr(mid0.CacheGroup)
	eCG.ParentCachegroupID = util.Ptr(mid0.CacheGroupID)
	eCG.SecondaryParentName = util.Ptr(mid1.CacheGroup)
	eCG.SecondaryParentCachegroupID = util.Ptr(mid1.CacheGroupID)
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	mCG0 := &tc.CacheGroupNullableV5{}
	mCG0.Name = util.Ptr(mid0.CacheGroup)
	mCG0.ID = util.Ptr(mid0.CacheGroupID)
	mCG0.ParentName = util.Ptr(org0.CacheGroup)
	mCG0.ParentCachegroupID = util.Ptr(org0.CacheGroupID)
	mCG0.SecondaryParentName = util.Ptr(org1.CacheGroup)
	mCG0.SecondaryParentCachegroupID = util.Ptr(org1.CacheGroupID)
	mCGType0 := tc.CacheGroupMidTypeName
	mCG0.Type = &mCGType0

	mCG1 := &tc.CacheGroupNullableV5{}
	mCG1.Name = util.Ptr(mid1.CacheGroup)
	mCG1.ID = util.Ptr(mid1.CacheGroupID)
	mCG1.ParentName = util.Ptr(org1.CacheGroup)
	mCG1.ParentCachegroupID = util.Ptr(org1.CacheGroupID)
	mCG1.SecondaryParentName = util.Ptr(org0.CacheGroup)
	mCG1.SecondaryParentCachegroupID = util.Ptr(org0.CacheGroupID)
	mCGType1 := tc.CacheGroupMidTypeName
	mCG1.Type = &mCGType1

	oCG0 := &tc.CacheGroupNullableV5{}
	oCG0.Name = util.Ptr(org0.CacheGroup)
	oCG0.ID = util.Ptr(org0.CacheGroupID)
	oCGType0 := tc.CacheGroupOriginTypeName
	oCG0.Type = &oCGType0

	oCG1 := &tc.CacheGroupNullableV5{}
	oCG1.Name = util.Ptr(org1.CacheGroup)
	oCG1.ID = util.Ptr(org1.CacheGroupID)
	oCGType1 := tc.CacheGroupOriginTypeName
	oCG1.Type = &oCGType1

	cgs := []tc.CacheGroupNullableV5{*eCG, *mCG0, *mCG1, *oCG0, *oCG1}

	dss := []DeliveryServiceServer{
		{Server: edge.ID, DeliveryService: *ds0.ID},
		{Server: mid0.ID, DeliveryService: *ds0.ID},
		{Server: mid1.ID, DeliveryService: *ds0.ID},
		{Server: org0.ID, DeliveryService: *ds0.ID},
		{Server: org1.ID, DeliveryService: *ds0.ID},

		{Server: edge.ID, DeliveryService: *ds1.ID},
		{Server: mid0.ID, DeliveryService: *ds1.ID},
		{Server: mid1.ID, DeliveryService: *ds1.ID},
		{Server: org0.ID, DeliveryService: *ds1.ID},
		{Server: org1.ID, DeliveryService: *ds1.ID},
	}
	cdn := &tc.CDNV5{
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
			t.Errorf("Missing required string(s) from ds/line: %s/%v\n%v", ds.XMLID, missing, txt)
		}

		excludes := []string{
			`hash_key`,
		}

		excluding := missingFrom(txt, excludes)
		if 1 != len(excludes) {
			t.Errorf("Excluded required string(s) from ds/line: %s/%v\n%v", ds.XMLID, excluding, txt)
		}
	}
}

func TestMakeStrategiesDotYamlMSONoTopologyNoMid(t *testing.T) {
	opt := &StrategiesYAMLOpts{VerboseComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := "HTTP"
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.Ptr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.Ptr("http://ds0.example.net")
	ds0.MultiSiteOrigin = true
	ds0.ProfileName = util.Ptr("dsprofile")
	dses := []DeliveryService{*ds0}

	parentConfigParams := []tc.ParameterV5{}

	serverParams := []tc.ParameterV5{
		tc.ParameterV5{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
	}

	edge := makeTestParentServer()

	origin0 := makeTestParentServer()
	origin0.CacheGroup = "originCG"
	origin0.CacheGroupID = 500
	origin0.HostName = "myorigin0"
	origin0.ID = 45
	setIP(origin0, "192.168.2.2")
	origin0.Type = tc.OriginTypeName
	origin0.TypeID = 991

	servers := []Server{*edge, *origin0}

	topologies := []tc.TopologyV5{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	eCG := &tc.CacheGroupNullableV5{}
	eCG.Name = util.Ptr(edge.CacheGroup)
	eCG.ID = util.Ptr(edge.CacheGroupID)
	eCG.ParentName = util.Ptr(origin0.CacheGroup)
	eCG.ParentCachegroupID = util.Ptr(origin0.CacheGroupID)
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	oCG := &tc.CacheGroupNullableV5{}
	oCG.Name = util.Ptr(origin0.CacheGroup)
	oCG.ID = util.Ptr(origin0.CacheGroupID)
	oCGType := tc.CacheGroupOriginTypeName
	oCG.Type = &oCGType

	cgs := []tc.CacheGroupNullableV5{*eCG, *oCG}

	dss := []DeliveryServiceServer{
		{
			Server:          edge.ID,
			DeliveryService: *ds0.ID,
		},
		{
			Server:          origin0.ID,
			DeliveryService: *ds0.ID,
		},
	}
	cdn := &tc.CDNV5{
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
		t.Errorf("Missing required string(s) from ds/line: %s/%v\n%v", ds0.XMLID, missing, txt)
	}
}

// Test for mso non topology where mid cache group has no primary/secondary
// parents assigned, just any arbitrary servers.
func TestMakeStrategiesDotYamlMSONoTopoMultiCG(t *testing.T) {
	opt := &StrategiesYAMLOpts{VerboseComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := "HTTP"
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.Ptr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.Ptr("http://ds0.example.net")
	ds0.MultiSiteOrigin = true

	dses := []DeliveryService{*ds0}

	parentConfigParams := []tc.ParameterV5{}

	serverParams := []tc.ParameterV5{
		tc.ParameterV5{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
	}

	edge := makeTestParentServer()
	edge.CacheGroup = "edgeCG"
	edge.CacheGroupID = 400

	mid := makeTestParentServer()
	mid.CacheGroup = "midCG"
	mid.CacheGroupID = 500
	mid.HostName = "mid0"
	mid.ID = 45
	setIP(mid, "192.168.2.2")

	org0 := makeTestParentServer()
	org0.CacheGroup = "orgCG0"
	org0.CacheGroupID = 501
	org0.HostName = "org0"
	org0.ID = 46
	setIP(org0, "192.168.2.3")
	org0.Type = tc.OriginTypeName
	org0.TypeID = 991

	org1 := makeTestParentServer()
	org1.CacheGroup = "orgCG1"
	org1.CacheGroupID = 502
	org1.HostName = "org1"
	org1.ID = 47
	setIP(org1, "192.168.2.4")
	org1.Type = tc.OriginTypeName
	org1.TypeID = 991

	servers := []Server{*edge, *mid, *org0, *org1}

	topologies := []tc.TopologyV5{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

	eCG := &tc.CacheGroupNullableV5{}
	eCG.Name = util.Ptr(edge.CacheGroup)
	eCG.ID = util.Ptr(edge.CacheGroupID)
	eCG.ParentName = util.Ptr(mid.CacheGroup)
	eCG.ParentCachegroupID = util.Ptr(mid.CacheGroupID)
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	// NOTE: no parent cache groups specified
	mCG := &tc.CacheGroupNullableV5{}
	mCG.Name = util.Ptr(mid.CacheGroup)
	mCG.ID = util.Ptr(mid.CacheGroupID)
	mCGType := tc.CacheGroupMidTypeName
	mCG.Type = &mCGType

	oCG0 := &tc.CacheGroupNullableV5{}
	oCG0.Name = util.Ptr(org0.CacheGroup)
	oCG0.ID = util.Ptr(org0.CacheGroupID)
	oCG0Type := tc.CacheGroupOriginTypeName
	oCG0.Type = &oCG0Type

	oCG1 := &tc.CacheGroupNullableV5{}
	oCG1.Name = util.Ptr(org1.CacheGroup)
	oCG1.ID = util.Ptr(org1.CacheGroupID)
	oCG1Type := tc.CacheGroupOriginTypeName
	oCG1.Type = &oCG1Type

	cgs := []tc.CacheGroupNullableV5{*eCG, *mCG, *oCG0, *oCG1}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          edge.ID,
			DeliveryService: *ds0.ID,
		},
		DeliveryServiceServer{
			Server:          org0.ID,
			DeliveryService: *ds0.ID,
		},
		DeliveryServiceServer{
			Server:          org1.ID,
			DeliveryService: *ds0.ID,
		},
	}
	cdn := &tc.CDNV5{
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
	opt := &StrategiesYAMLOpts{VerboseComments: false, HdrComment: "myHeaderComment", GoDirect: "true"}

	// Toplogy ds, MSO
	ds0 := makeParentDS()
	ds0Type := "HTTP"
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.Ptr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.Ptr("http://ds0.example.net")
	ds0.ProfileID = util.Ptr(311)
	ds0.ProfileName = util.Ptr("ds0Profile")
	ds0.MultiSiteOrigin = true
	ds0.Topology = util.Ptr("t0")

	// Toplogy ds, non MSO
	ds1 := makeParentDS()
	ds1.ID = util.Ptr(43)
	ds1Type := "DNS"
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.Ptr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.Ptr("http://ds1.example.net")
	ds1.ProfileID = util.Ptr(311)
	ds1.ProfileName = util.Ptr("ds0Profile")
	ds1.Topology = util.Ptr("t0")

	dsesall := []DeliveryService{*ds0, *ds1}

	parentConfigParams := []tc.ParameterV5{
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
		tcparam := tc.ParameterV5{
			Name:       key,
			ConfigFile: "parent.config",
			Value:      val,
			Profiles:   []byte(`["ds0Profile"]`),
		}
		parentConfigParams = append(parentConfigParams, tcparam)
	}

	serverParams := []tc.ParameterV5{
		{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "9",
			Profiles:   []byte(`["global"]`),
		},
	}

	edge := makeTestParentServer()
	edge.CacheGroup = "edgeCG"
	edge.CacheGroupID = 400

	mid0 := makeTestParentServer()
	mid0.CacheGroup = "midCG0"
	mid0.CacheGroupID = 500
	mid0.HostName = "mymid0"
	mid0.ID = 45
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.CacheGroup = "midCG1"
	mid1.CacheGroupID = 501
	mid1.HostName = "mymid1"
	mid1.ID = 46
	setIP(mid1, "192.168.2.3")

	opl0 := makeTestParentServer()
	opl0.CacheGroup = "oplCG0"
	opl0.CacheGroupID = 502
	opl0.HostName = "myopl0"
	opl0.ID = 46
	setIP(opl0, "192.168.2.4")

	opl1 := makeTestParentServer()
	opl1.CacheGroup = "oplCG1"
	opl1.CacheGroupID = 503
	opl1.HostName = "myopl1"
	opl1.ID = 47
	setIP(opl1, "192.168.2.5")

	org0 := makeTestParentServer()
	org0.CacheGroup = "orgCG0"
	org0.CacheGroupID = 504
	org0.HostName = "myorg0"
	org0.ID = 48
	setIP(org0, "192.168.2.6")
	org0.Type = tc.OriginTypeName
	org0.TypeID = 991

	org1 := makeTestParentServer()
	org1.CacheGroup = "orgCG1"
	org1.CacheGroupID = 505
	org1.HostName = "myorg1"
	org1.ID = 49
	setIP(org1, "192.168.2.7")
	org1.Type = tc.OriginTypeName
	org1.TypeID = 991

	servers := []Server{*edge, *mid0, *mid1, *opl0, *opl1, *org0, *org1}

	topologies := []tc.TopologyV5{
		{
			Name: "t0",
			Nodes: []tc.TopologyNodeV5{
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

	eCG := &tc.CacheGroupNullableV5{}
	eCG.Name = util.Ptr(edge.CacheGroup)
	eCG.ID = util.Ptr(edge.CacheGroupID)
	eCG.ParentName = util.Ptr(mid0.CacheGroup)
	eCG.ParentCachegroupID = util.Ptr(mid0.CacheGroupID)
	eCG.SecondaryParentName = util.Ptr(mid1.CacheGroup)
	eCG.SecondaryParentCachegroupID = util.Ptr(mid1.CacheGroupID)
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	mCG0 := &tc.CacheGroupNullableV5{}
	mCG0.Name = util.Ptr(mid0.CacheGroup)
	mCG0.ID = util.Ptr(mid0.CacheGroupID)
	mCG0.ParentName = util.Ptr(opl0.CacheGroup)
	mCG0.ParentCachegroupID = util.Ptr(opl0.CacheGroupID)
	mCG0.SecondaryParentName = util.Ptr(opl1.CacheGroup)
	mCG0.SecondaryParentCachegroupID = util.Ptr(opl1.CacheGroupID)
	mCGType0 := tc.CacheGroupMidTypeName
	mCG0.Type = &mCGType0

	mCG1 := &tc.CacheGroupNullableV5{}
	mCG1.Name = util.Ptr(mid1.CacheGroup)
	mCG1.ID = util.Ptr(mid1.CacheGroupID)
	mCG1.ParentName = util.Ptr(opl1.CacheGroup)
	mCG1.ParentCachegroupID = util.Ptr(opl1.CacheGroupID)
	mCG1.SecondaryParentName = util.Ptr(opl0.CacheGroup)
	mCG1.SecondaryParentCachegroupID = util.Ptr(opl0.CacheGroupID)
	mCGType1 := tc.CacheGroupMidTypeName
	mCG1.Type = &mCGType1

	oplCG0 := &tc.CacheGroupNullableV5{}
	oplCG0.Name = util.Ptr(opl0.CacheGroup)
	oplCG0.ID = util.Ptr(opl0.CacheGroupID)
	oplCG0.ParentName = util.Ptr(org0.CacheGroup)
	oplCG0.ParentCachegroupID = util.Ptr(org0.CacheGroupID)
	oplCG0.SecondaryParentName = util.Ptr(org1.CacheGroup)
	oplCG0.SecondaryParentCachegroupID = util.Ptr(org1.CacheGroupID)
	oplCGType0 := tc.CacheGroupMidTypeName
	oplCG0.Type = &oplCGType0

	oplCG1 := &tc.CacheGroupNullableV5{}
	oplCG1.Name = util.Ptr(opl1.CacheGroup)
	oplCG1.ID = util.Ptr(opl1.CacheGroupID)
	oplCG1.ParentName = util.Ptr(org1.CacheGroup)
	oplCG1.ParentCachegroupID = util.Ptr(org1.CacheGroupID)
	oplCG1.SecondaryParentName = util.Ptr(org0.CacheGroup)
	oplCG1.SecondaryParentCachegroupID = util.Ptr(org0.CacheGroupID)
	oplCGType1 := tc.CacheGroupMidTypeName
	oplCG1.Type = &oplCGType1

	oCG0 := &tc.CacheGroupNullableV5{}
	oCG0.Name = util.Ptr(org0.CacheGroup)
	oCG0.ID = util.Ptr(org0.CacheGroupID)
	oCGType0 := tc.CacheGroupOriginTypeName
	oCG0.Type = &oCGType0

	oCG1 := &tc.CacheGroupNullableV5{}
	oCG1.Name = util.Ptr(org1.CacheGroup)
	oCG1.ID = util.Ptr(org1.CacheGroupID)
	oCGType1 := tc.CacheGroupOriginTypeName
	oCG1.Type = &oCGType1

	cgs := []tc.CacheGroupNullableV5{*eCG, *mCG0, *mCG1, *oplCG0, *oplCG1, *oCG0, *oCG1}

	dss := []DeliveryServiceServer{
		{Server: org0.ID, DeliveryService: *ds0.ID},
		{Server: org1.ID, DeliveryService: *ds0.ID},

		{Server: edge.ID, DeliveryService: *ds1.ID},
		{Server: mid0.ID, DeliveryService: *ds1.ID},
		{Server: mid1.ID, DeliveryService: *ds1.ID},
		{Server: opl0.ID, DeliveryService: *ds1.ID},
		{Server: opl1.ID, DeliveryService: *ds1.ID},
		{Server: org0.ID, DeliveryService: *ds1.ID},
		{Server: org1.ID, DeliveryService: *ds1.ID},
	}
	cdn := &tc.CDNV5{
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
			t.Errorf("Missing required string(s) from ds/line: %s/%v\n%v", ds.XMLID, missing, txt)
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
			t.Errorf("Excluded required string(s) from ds/line: %s/%v\n%v", ds.XMLID, excluding, txt)
		}
	}
}
