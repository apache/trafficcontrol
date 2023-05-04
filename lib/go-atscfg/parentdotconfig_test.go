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

func TestMakeParentDotConfig(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")

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

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_domain=ds0.example.net") {
		t.Errorf("expected parent 'dest_domain=ds0.example.net', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "dest_domain=ds1.example.net") {
		t.Errorf("expected parent 'dest_domain=ds0.example.net', actual: '%v'", txt)
	}
	if !warningsContains(cfg.Warnings, "myQStringHandlingParam") {
		t.Errorf("expected malformed qstring 'myQstringParam' in warnings, actual: '%v' val '%v'", cfg.Warnings, txt)
	}
}

func TestMakeParentDotConfigCapabilities(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")

	dses := []DeliveryService{*ds0}

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
	mid0.HostName = util.StrPtr("my-parent-nocaps")
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.HostName = util.StrPtr("mymid0")
	mid0.ID = util.IntPtr(45)
	mid0.CachegroupID = util.IntPtr(423)
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.HostName = util.StrPtr("my-parent-fooonly")
	mid1.Cachegroup = util.StrPtr("midCG")
	mid1.ID = util.IntPtr(46)
	mid1.CachegroupID = util.IntPtr(423)
	setIP(mid1, "192.168.2.3")

	mid2 := makeTestParentServer()
	mid2.HostName = util.StrPtr("my-parent-foobar")
	mid2.Cachegroup = util.StrPtr("midCG")
	mid2.ID = util.IntPtr(47)
	mid2.CachegroupID = util.IntPtr(423)
	setIP(mid1, "192.168.2.4")

	servers := []Server{*server, *mid0, *mid1, *mid2}

	topologies := []tc.Topology{}
	serverCapabilities := map[int]map[ServerCapability]struct{}{
		*server.ID: map[ServerCapability]struct{}{"FOO": {}},
		*mid1.ID:   map[ServerCapability]struct{}{"FOO": {}},
		*mid2.ID:   map[ServerCapability]struct{}{"FOO": {}, "BAR": {}},
	}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{
		*ds0.ID: map[ServerCapability]struct{}{"FOO": {}},
	}

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
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	lines := strings.Split(txt, "\n")

	if expectedLines := 5; len(lines) != expectedLines {
		t.Fatalf("expected %v lines (comment, blank, ds, dot remap, and empty newline), actual: '%+v' text %v", expectedLines, len(lines), txt)
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
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")
	ds0.MultiSiteOrigin = util.BoolPtr(true)
	dses := []DeliveryService{*ds0}

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
	mid0.Cachegroup = util.StrPtr("midCG0")
	mid0.CachegroupID = util.IntPtr(500)
	mid0.HostName = util.StrPtr("my-parent-0")
	mid0.DomainName = util.StrPtr("my-parent-0-domain")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG1")
	mid1.CachegroupID = util.IntPtr(501)
	mid1.HostName = util.StrPtr("my-parent-1")
	mid1.DomainName = util.StrPtr("my-parent-1-domain")
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
	eCG.SecondaryParentName = mid1.Cachegroup
	eCG.SecondaryParentCachegroupID = mid1.CachegroupID
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	mCG := &tc.CacheGroupNullable{}
	mCG.Name = mid0.Cachegroup
	mCG.ID = mid0.CachegroupID
	mCGType := tc.CacheGroupMidTypeName
	mCG.Type = &mCGType

	mCG1 := &tc.CacheGroupNullable{}
	mCG1.Name = mid1.Cachegroup
	mCG1.ID = mid1.CachegroupID
	mCGType1 := tc.CacheGroupMidTypeName
	mCG1.Type = &mCGType1

	cgs := []tc.CacheGroupNullable{*eCG, *mCG, *mCG1}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds0.ID,
		},
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	txtx := strings.Replace(txt, " ", "", -1)

	if !strings.Contains(txtx, `secondary_parent="my-parent-1.my-parent-1-domain`) {
		t.Errorf("expected secondary parent 'my-parent-1.my-parent-1-domain', actual: '%v'", txt)
	}

	if strings.Contains(txtx, "parent_retry") {
		t.Errorf("Did not expect parent_retry parameter at edge/inner: '%v'", txt)
	}
}

func TestMakeParentDotConfigMSONoPrimaryParent(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")
	ds0.MultiSiteOrigin = util.BoolPtr(true)
	dses := []DeliveryService{*ds0}

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
	mid0.Cachegroup = util.StrPtr("midCG0")
	mid0.CachegroupID = util.IntPtr(500)
	mid0.HostName = util.StrPtr("my-parent-0")
	mid0.DomainName = util.StrPtr("my-parent-0-domain")
	mid0.Status = util.StrPtr(string(tc.CacheStatusAdminDown))
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG1")
	mid1.CachegroupID = util.IntPtr(501)
	mid1.HostName = util.StrPtr("my-parent-1")
	mid1.DomainName = util.StrPtr("my-parent-1-domain")
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
	eCG.SecondaryParentName = mid1.Cachegroup
	eCG.SecondaryParentCachegroupID = mid1.CachegroupID
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	mCG := &tc.CacheGroupNullable{}
	mCG.Name = mid0.Cachegroup
	mCG.ID = mid0.CachegroupID
	mCGType := tc.CacheGroupMidTypeName
	mCG.Type = &mCGType

	mCG1 := &tc.CacheGroupNullable{}
	mCG1.Name = mid1.Cachegroup
	mCG1.ID = mid1.CachegroupID
	mCGType1 := tc.CacheGroupMidTypeName
	mCG1.Type = &mCGType1

	cgs := []tc.CacheGroupNullable{*eCG, *mCG, *mCG1}

	dss := []DeliveryServiceServer{
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds0.ID,
		},
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	txtx := strings.Replace(txt, " ", "", -1)

	if !strings.Contains(txtx, `parent="my-parent-1.my-parent-1-domain:80|0.999`) {
		t.Errorf("expected primary parent 'my-parent-1.my-parent-1-domain', actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigMSONoTopologyNoMid(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")
	ds0.MultiSiteOrigin = util.BoolPtr(true)
	ds0.ProfileName = util.StrPtr("dsprofile")
	dses := []DeliveryService{*ds0}

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
			Name:       "mso.parent_retry",
			ConfigFile: "parent.config",
			Value:      "both",
			Profiles:   []byte(`["` + *ds0.ProfileName + `"]`),
		},
		tc.Parameter{
			Name:       "mso.algorithm",
			ConfigFile: "parent.config",
			Value:      "consistent_hash",
			Profiles:   []byte(`["` + *ds0.ProfileName + `"]`),
		},
		tc.Parameter{
			Name:       "mso.unavailable_server_retry_responses",
			ConfigFile: "parent.config",
			Value:      `"500,502,503,542"`,
			Profiles:   []byte(`["` + *ds0.ProfileName + `"]`),
		},
		tc.Parameter{
			Name:       "mso.max_simple_retries",
			ConfigFile: "parent.config",
			Value:      "2",
			Profiles:   []byte(`["` + *ds0.ProfileName + `"]`),
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

	cfg, err := MakeParentDotConfig(dses, edge, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, `parent="myorigin0.mydomain.example.net:80`) {
		t.Errorf("expected parent myorigin0, actual: '%v'", txt)
	}
	if !strings.Contains(txt, `unavailable_server_retry_responses="500,502,503,542`) {
		t.Errorf(`expected unavailable_server_retry_repsonse 500,502,503,542 from DS params, actual: '%v'`, txt)
	}
	if !strings.Contains(txt, `max_simple_retries=2`) {
		t.Errorf(`expected max_simple_retries=2 from DS params, actual: '%v'`, txt)
	}
}

// Test for mso non topology where mid cache group has no primary/secondary
// parents assigned, just any arbitrary servers.
func TestMakeParentDotConfigMSONoTopoMultiCG(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

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

	cfg, err := MakeParentDotConfig(dses, mid, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if strings.Contains(txt, `secondary_parent="`) {
		t.Errorf("Did not expect secondary parent', actual: '%v'", txt)
	}

	if !strings.Contains(txt, ` parent="org0.mydomain.example.net:80|0.999;org1.mydomain.example.net:80|0.999"`) {
		t.Errorf("Expected parent with org0 and org1 both listed, actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigTopologies(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.Topology = util.StrPtr("t0")

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
	server.Cachegroup = util.StrPtr("edgeCG")
	server.CachegroupID = util.IntPtr(400)

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.CachegroupID = util.IntPtr(500)
	mid0.HostName = util.StrPtr("mymid")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG")
	mid1.CachegroupID = util.IntPtr(500)
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")

	servers := []Server{*server, *mid0, *mid1}

	topologies := []tc.Topology{
		tc.Topology{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				tc.TopologyNode{
					Cachegroup: "edgeCG",
					Parents:    []int{1},
				},
				tc.TopologyNode{
					Cachegroup: "midCG",
				},
			},
		},
	}

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

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_domain=ds0.example.net") {
		t.Errorf("expected parent 'dest_domain=ds0.example.net', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "dest_domain=ds1.example.net") {
		t.Errorf("expected parent 'dest_domain=ds1.example.net', actual: '%v'", txt)
	}
	if !warningsContains(cfg.Warnings, "myQStringHandlingParam") {
		t.Errorf("expected warning for malformed myQStringHandlingParam', actual: '%+v'", cfg.Warnings)
	}
	if strings.Contains(txt, "# topology") {
		// ATS doesn't support inline comments in parent.config
		t.Errorf("expected: no inline '# topology' comment, actual: '%v'", txt)
	}
}

// TestMakeParentDotConfigNotInTopologies tests when a given edge is NOT in a Topology, that it doesn't add a remap line.
func TestMakeParentDotConfigNotInTopologies(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")
	ds0.Topology = util.StrPtr("t0")

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

	topologies := []tc.Topology{
		tc.Topology{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				tc.TopologyNode{
					Cachegroup: "otherEdgeCG",
					Parents:    []int{1},
				},
				tc.TopologyNode{
					Cachegroup: "midCG",
				},
			},
		},
	}

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

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if strings.Contains(txt, "dest_domain=ds0.example.net") {
		t.Errorf("expected parent 'dest_domain=ds0.example.net' to NOT contain Topology DS without this edge: '%v'", txt)
	}
	if !strings.Contains(txt, "dest_domain=ds1.example.net") {
		t.Errorf("expected parent 'dest_domain=ds0.example.net', actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigTopologiesCapabilities(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0.ID = util.IntPtr(42)
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")
	ds0.Topology = util.StrPtr("t0")

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.Topology = util.StrPtr("t0")

	ds2 := makeParentDS()
	ds2.ID = util.IntPtr(44)
	ds2Type := tc.DSTypeDNS
	ds2.Type = &ds2Type
	ds2.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds2.OrgServerFQDN = util.StrPtr("http://ds2.example.net")
	ds2.Topology = util.StrPtr("t0")

	dses := []DeliveryService{*ds0, *ds1, *ds2}

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
	server.Cachegroup = util.StrPtr("edgeCG")
	server.CachegroupID = util.IntPtr(400)

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.CachegroupID = util.IntPtr(500)
	mid0.HostName = util.StrPtr("mymid0")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG")
	mid1.CachegroupID = util.IntPtr(500)
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")

	servers := []Server{*server, *mid0, *mid1}

	topologies := []tc.Topology{
		tc.Topology{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				tc.TopologyNode{
					Cachegroup: "edgeCG",
					Parents:    []int{1},
				},
				tc.TopologyNode{
					Cachegroup: "midCG",
				},
			},
		},
	}

	serverCapabilities := map[int]map[ServerCapability]struct{}{
		44: map[ServerCapability]struct{}{"FOO": {}},
		45: map[ServerCapability]struct{}{"FOO": {}},
		46: map[ServerCapability]struct{}{"FOO": {}},
	}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{
		*ds1.ID: map[ServerCapability]struct{}{"FOO": {}},
		*ds2.ID: map[ServerCapability]struct{}{"BAR": {}},
	}

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
		DeliveryServiceServer{
			Server:          *server.ID,
			DeliveryService: *ds2.ID,
		},
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_domain=ds0.example.net") {
		t.Errorf("expected parent 'dest_domain=ds0.example.net' without required capabilities: '%v'", txt)
	}
	if !strings.Contains(txt, "dest_domain=ds1.example.net") {
		t.Errorf("expected parent 'dest_domain=ds1.example.net' with necessary required capabilities, actual: '%v'", txt)
	}
	if strings.Contains(txt, "dest_domain=ds2.example.net") {
		t.Errorf("expected no parent 'dest_domain=ds2.example.net' without necessary required capabilities, actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigTopologiesOmitOfflineParents(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.Topology = util.StrPtr("t0")

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
	server.Cachegroup = util.StrPtr("edgeCG")
	server.CachegroupID = util.IntPtr(400)

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.CachegroupID = util.IntPtr(500)
	mid0.HostName = util.StrPtr("mymid-should-omit")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")
	statusOffline := string(tc.CacheStatusOffline)
	mid0.Status = &statusOffline

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG")
	mid1.CachegroupID = util.IntPtr(500)
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")

	servers := []Server{*server, *mid0, *mid1}

	topologies := []tc.Topology{
		tc.Topology{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				tc.TopologyNode{
					Cachegroup: "edgeCG",
					Parents:    []int{1},
				},
				tc.TopologyNode{
					Cachegroup: "midCG",
				},
			},
		},
	}

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

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_domain=ds0.example.net") {
		t.Errorf("expected parent 'dest_domain=ds0.example.net', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "dest_domain=ds1.example.net") {
		t.Errorf("expected parent 'dest_domain=ds1.example.net', actual: '%v'", txt)
	}
	if !warningsContains(cfg.Warnings, "myQStringHandlingParam") {
		t.Errorf("expected warning for malformed myQStringHandlingParam', actual: '%+v'", cfg.Warnings)
	}

	if strings.Contains(txt, "should-omit") {
		t.Errorf("Topology expected to omit OFFLINE mid, actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigTopologiesOmitDifferentCDNParents(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.Topology = util.StrPtr("t0")

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
	server.Cachegroup = util.StrPtr("edgeCG")
	server.CachegroupID = util.IntPtr(400)

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.CachegroupID = util.IntPtr(500)
	mid0.HostName = util.StrPtr("mymid-should-omit")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")
	mid0.CDNName = util.StrPtr("myCDN-different-than-edge")
	mid0CDNID := *server.CDNID + 1
	mid0.CDNID = &mid0CDNID

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG")
	mid1.CachegroupID = util.IntPtr(500)
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")

	servers := []Server{*server, *mid0, *mid1}

	topologies := []tc.Topology{
		tc.Topology{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				tc.TopologyNode{
					Cachegroup: "edgeCG",
					Parents:    []int{1},
				},
				tc.TopologyNode{
					Cachegroup: "midCG",
				},
			},
		},
	}

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

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_domain=ds0.example.net") {
		t.Errorf("expected parent 'dest_domain=ds0.example.net', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "dest_domain=ds1.example.net") {
		t.Errorf("expected parent 'dest_domain=ds1.example.net', actual: '%v'", txt)
	}
	if !warningsContains(cfg.Warnings, "myQStringHandlingParam") {
		t.Errorf("expected warning for malformed myQStringHandlingParam', actual: '%+v'", cfg.Warnings)
	}

	if strings.Contains(txt, "should-omit") {
		t.Errorf("Topology expected to omit parent with a different CDN, actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigTopologiesMSO(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment", GoDirect: true}

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.Topology = util.StrPtr("t0")
	ds1.MultiSiteOrigin = util.BoolPtr(true)
	ds1.ProfileName = util.StrPtr("dsprofile")

	dses := []DeliveryService{*ds1}

	parentConfigParams := []tc.Parameter{
		tc.Parameter{
			Name:       ParentConfigRetryKeysDefault.Algorithm,
			ConfigFile: "parent.config",
			Value:      tc.AlgorithmConsistentHash,
			Profiles:   []byte(`["serverprofile"]`),
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

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_domain=ds1.example.net") {
		t.Errorf("expected parent 'dest_domain=ds1.example.net', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "myorigin0") {
		t.Errorf("expected origin0 with DeliveryServiceServer assigned to this DS, actual: '%v'", txt)
	}
	if strings.Contains(txt, "myorigin1") {
		t.Errorf("expected no origin1 without DeliveryServiceServer assigned to this DS, actual: '%v'", txt)
	}

	if !strings.Contains(txt, "go_direct=true") {
		t.Errorf("expected MSO Topologies to Origin to go_direct=true, actual: '%v'", txt)
	}

	if !strings.Contains(txt, "parent_is_proxy=false") {
		t.Errorf("expected MSO Topologies to Origin to parent_is_proxy=false, actual: '%v'", txt)
	}

	t.Run("MSO topologoies default qstring=ignore", func(t *testing.T) {
		cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(cfg.Text, "qstring=ignore") {
			t.Errorf("expected MSO Topologies to Origin to default to qstring=ignore, actual: '%v'", cfg.Text)
		}
	})

	t.Run("MSO topologoies param qstring=ignore", func(t *testing.T) {
		parentConfigParamsWithQstr := append(parentConfigParams, tc.Parameter{
			Name:       ParentConfigParamQString,
			ConfigFile: "parent.config",
			Value:      "ignore",
			Profiles:   []byte(`["serverprofile"]`),
		})

		cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParamsWithQstr, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(cfg.Text, "qstring=ignore") {
			t.Errorf("expected MSO Topologies to Origin to default to qstring=ignore, actual: '%v'", cfg.Text)
		}
	})

	t.Run("MSO topologoies param qstring=consider", func(t *testing.T) {
		parentConfigParamsWithQstr := append(parentConfigParams, tc.Parameter{
			Name:       ParentConfigParamQStringHandling,
			ConfigFile: "parent.config",
			Value:      "consider",
			Profiles:   []byte(`["` + *ds1.ProfileName + `"]`),
		})

		cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParamsWithQstr, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(cfg.Text, "qstring=consider") {
			t.Errorf("expected MSO Topologies to Origin with param to qstring=consider, actual: '''%v''' warnings '''%+v'''", cfg.Text, cfg.Warnings)
		}
	})

	t.Run("MSO topologoies param ds qstring consider", func(t *testing.T) {
		ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
		dses := []DeliveryService{*ds1}

		cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(cfg.Text, "qstring=consider") {
			t.Errorf("expected MSO Topologies to Origin with param to qstring=consider, actual: '''%v''' warnings '''%+v'''", cfg.Text, cfg.Warnings)
		}
	})
}

func TestMakeParentDotConfigTopologiesMSOWithCapabilities(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.Topology = util.StrPtr("t0")
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
	server.Cachegroup = util.StrPtr("edgeCG")
	server.CachegroupID = util.IntPtr(400)
	server.ID = util.IntPtr(44)

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

	serverCapabilities := map[int]map[ServerCapability]struct{}{
		*server.ID: {
			ServerCapability("FOO"): struct{}{},
		},
	}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{
		*ds1.ID: {
			ServerCapability("FOO"): struct{}{},
		},
	}

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

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_domain=ds1.example.net") {
		t.Errorf("expected parent 'dest_domain=ds1.example.net', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "myorigin0") {
		t.Errorf("expected origin0 with DeliveryServiceServer assigned to this DS, actual: '%v'", txt)
	}
	if strings.Contains(txt, "myorigin1") {
		t.Errorf("expected no origin1 without DeliveryServiceServer assigned to this DS, actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigMSOWithCapabilities(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
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
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
	}

	mid := makeTestParentServer()
	mid.Cachegroup = util.StrPtr("midCG")
	mid.Type = "MID"
	mid.CachegroupID = util.IntPtr(400)
	mid.ID = util.IntPtr(44)

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

	servers := []Server{*mid, *origin0, *origin1}

	serverCapabilities := map[int]map[ServerCapability]struct{}{
		*mid.ID: {
			ServerCapability("FOO"): struct{}{},
		},
	}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{
		*ds1.ID: {
			ServerCapability("FOO"): struct{}{},
		},
	}

	midCG := &tc.CacheGroupNullable{}
	midCG.Name = mid.Cachegroup
	midCG.ID = mid.CachegroupID
	midCG.ParentName = origin0.Cachegroup
	midCG.ParentCachegroupID = origin0.CachegroupID
	midCGType := tc.CacheGroupMidTypeName
	midCG.Type = &midCGType

	oCG := &tc.CacheGroupNullable{}
	oCG.Name = origin0.Cachegroup
	oCG.ID = origin0.CachegroupID
	oCGType := tc.CacheGroupOriginTypeName
	oCG.Type = &oCGType

	cgs := []tc.CacheGroupNullable{*midCG, *oCG}

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
	topologies := []tc.Topology{}

	cfg, err := MakeParentDotConfig(dses, mid, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_domain=ds1.example.net") {
		t.Errorf("expected parent 'dest_domain=ds1.example.net', actual: '%v' warnings %+v", txt, cfg.Warnings)
	}
	if !strings.Contains(txt, "myorigin0") {
		t.Errorf("expected origin0 with DeliveryServiceServer assigned to this DS, actual: '%v'", txt)
	}
	if strings.Contains(txt, "myorigin1") {
		t.Errorf("expected no origin1 without DeliveryServiceServer assigned to this DS, actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigTopologiesMSOParams(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

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
			Name:       ParentConfigRetryKeysMSO.Algorithm,
			ConfigFile: "parent.config",
			Value:      "consistent_hash",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigRetryKeysMSO.ParentRetry,
			ConfigFile: "parent.config",
			Value:      "both",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigRetryKeysMSO.MaxSimpleRetries,
			ConfigFile: "parent.config",
			Value:      "14",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigRetryKeysMSO.MaxUnavailableRetries,
			ConfigFile: "parent.config",
			Value:      "9",
			Profiles:   []byte(`["ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigRetryKeysMSO.UnavailableRetryResponses,
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

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_domain=ds1.example.net") {
		t.Errorf("expected parent 'dest_domain=ds1.example.net', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "myorigin0") {
		t.Errorf("expected origin0 with DeliveryServiceServer assigned to this DS, actual: '%v'", txt)
	}
	if strings.Contains(txt, "myorigin1") {
		t.Errorf("expected no origin1 without DeliveryServiceServer assigned to this DS, actual: '%v'", txt)
	}
	if !strings.Contains(txt, "parent_retry=both") {
		t.Errorf("expected DS MSO parent_retry param 'both', actual: '%v'", txt)
	}
	if !strings.Contains(txt, `unavailable_server_retry_responses="400,503"`) {
		t.Errorf(`expected DS MSO unavailable_server_retry_responses param '"400,503"'', actual: '%v'`, txt)
	}
	if !strings.Contains(txt, "max_simple_retries=14") {
		t.Errorf("expected DS MSO max_simple_retries param '14', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "max_unavailable_server_retries=9") {
		t.Errorf("expected DS MSO max_unavailable_server_retries param '9', actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigTopologiesParams(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

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
			Value:      "consistent_hash",
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

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_domain=ds1.example.net") {
		t.Errorf("expected parent 'dest_domain=ds1.example.net', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "myorigin0") {
		t.Errorf("expected origin0 with DeliveryServiceServer assigned to this DS, actual: '%v'", txt)
	}
	if strings.Contains(txt, "myorigin1") {
		t.Errorf("expected no origin1 without DeliveryServiceServer assigned to this DS, actual: '%v'", txt)
	}
	if !strings.Contains(txt, "parent_retry=both") {
		t.Errorf("expected DS MSO parent_retry param 'both', actual: '%v'", txt)
	}
	if !strings.Contains(txt, `unavailable_server_retry_responses="400,503"`) {
		t.Errorf(`expected DS MSO unavailable_server_retry_responses param '"400,503"'', actual: '%v'`, txt)
	}
	if !strings.Contains(txt, "max_simple_retries=14") {
		t.Errorf("expected DS MSO max_simple_retries param '14', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "max_unavailable_server_retries=9") {
		t.Errorf("expected DS MSO max_unavailable_server_retries param '9', actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigTopologiesNonStandardServerTypes(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.Topology = util.StrPtr("t0")

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
	server.Cachegroup = util.StrPtr("edgeCG")
	server.CachegroupID = util.IntPtr(400)

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.CachegroupID = util.IntPtr(500)
	mid0.HostName = util.StrPtr("mymid")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")
	mid0.Type = "MIDSOMETHING-RANDOM"

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG")
	mid1.CachegroupID = util.IntPtr(500)
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	mid1.Type = "MID_SOMETHING_ELSE"
	setIP(mid1, "192.168.2.3")

	servers := []Server{*server, *mid0, *mid1}

	topologies := []tc.Topology{
		tc.Topology{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				tc.TopologyNode{
					Cachegroup: "edgeCG",
					Parents:    []int{1},
				},
				tc.TopologyNode{
					Cachegroup: "midCG",
				},
			},
		},
	}

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

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_domain=ds0.example.net") {
		t.Errorf("expected parent 'dest_domain=ds0.example.net', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "dest_domain=ds1.example.net") {
		t.Errorf("expected parent 'dest_domain=ds1.example.net', actual: '%v'", txt)
	}
	if !warningsContains(cfg.Warnings, "myQStringHandlingParam") {
		t.Errorf("expected warning for malformed myQStringHandlingParam', actual: '%+v'", cfg.Warnings)
	}
	if strings.Contains(txt, "# topology") {
		// ATS doesn't support inline comments in parent.config
		t.Errorf("expected: no inline '# topology' comment, actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigSecondaryMode(t *testing.T) {

	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")
	ds0.ProfileID = util.IntPtr(311)
	ds0.ProfileName = util.StrPtr("ds0Profile")

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.Topology = util.StrPtr("t0")
	ds1.ProfileID = util.IntPtr(312)
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
		tc.Parameter{
			Name:       ParentConfigRetryKeysDefault.SecondaryMode,
			ConfigFile: "parent.config",
			Value:      "",
			Profiles:   []byte(`["ds0Profile","ds1Profile"]`),
		},
		tc.Parameter{
			Name:       ParentConfigRetryKeysFirst.SecondaryMode,
			ConfigFile: "parent.config",
			Value:      "",
			Profiles:   []byte(`["ds0Profile","ds1Profile"]`),
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

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.CachegroupID = util.IntPtr(500)
	mid0.HostName = util.StrPtr("mymid")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG2")
	mid1.CachegroupID = util.IntPtr(501)
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")

	servers := []Server{*server, *mid0, *mid1}

	topologies := []tc.Topology{
		tc.Topology{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				tc.TopologyNode{
					Cachegroup: "edgeCG",
					Parents:    []int{1, 2},
				},
				tc.TopologyNode{
					Cachegroup: "midCG",
				},
				tc.TopologyNode{
					Cachegroup: "midCG2",
				},
			},
		},
	}

	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

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

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_domain=ds0.example.net") {
		t.Errorf("expected parent 'dest_domain=ds0.example.net', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "dest_domain=ds1.example.net") {
		t.Errorf("expected parent 'dest_domain=ds1.example.net', actual: '%v'", txt)
	}
	if !warningsContains(cfg.Warnings, "myQStringHandlingParam") {
		t.Errorf("expected warning for malformed myQStringHandlingParam', actual: '%+v'", cfg.Warnings)
	}
	if strings.Count(txt, "secondary_mode=2") != 2 {
		t.Errorf("expected secondary_mode=2 for both Topology and DSS DSes with ParentConfigParamSecondaryMode parameter and secondary parents, actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigNoSecondaryMode(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")
	ds0.ProfileID = util.IntPtr(311)
	ds0.ProfileName = util.StrPtr("ds0Profile")

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.Topology = util.StrPtr("t0")
	ds1.ProfileID = util.IntPtr(312)
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
			Value:      "8",
			Profiles:   []byte(`["global"]`),
		},
	}

	server := makeTestParentServer()
	server.Cachegroup = util.StrPtr("edgeCG")
	server.CachegroupID = util.IntPtr(400)

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.CachegroupID = util.IntPtr(500)
	mid0.HostName = util.StrPtr("mymid")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG2")
	mid1.CachegroupID = util.IntPtr(501)
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")

	servers := []Server{*server, *mid0, *mid1}

	topologies := []tc.Topology{
		tc.Topology{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				tc.TopologyNode{
					Cachegroup: "edgeCG",
					Parents:    []int{1, 2},
				},
				tc.TopologyNode{
					Cachegroup: "midCG",
				},
				tc.TopologyNode{
					Cachegroup: "midCG2",
				},
			},
		},
	}

	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

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

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_domain=ds0.example.net") {
		t.Errorf("expected parent 'dest_domain=ds0.example.net', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "dest_domain=ds1.example.net") {
		t.Errorf("expected parent 'dest_domain=ds1.example.net', actual: '%v'", txt)
	}
	if !warningsContains(cfg.Warnings, "myQStringHandlingParam") {
		t.Errorf("expected warning for malformed myQStringHandlingParam', actual: '%+v'", cfg.Warnings)
	}
	if !strings.Contains(txt, "secondary_mode=1") {
		t.Errorf("expected default secondary_mode=1 for DSes without ParentConfigParamSecondaryMode parameter, actual: '%v'", txt)
	}

	if strings.Contains(txt, `topology 't0'`) {
		t.Errorf("expected no comment with topology name, actual: '%v'", txt)
	}
	if strings.Contains(txt, `ds 'ds1'`) {
		t.Errorf("expected no comment with delivery service name, actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigComments(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: true, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")

	dses := []DeliveryService{*ds0, *ds1}

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

	serverParams := []tc.Parameter{
		{
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
		{
			Server:          *server.ID,
			DeliveryService: *ds0.ID,
		},
		{
			Server:          *server.ID,
			DeliveryService: *ds1.ID,
		},
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_domain=ds0.example.net") {
		t.Errorf("expected parent 'dest_domain=ds0.example.net', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "dest_domain=ds1.example.net") {
		t.Errorf("expected parent 'dest_domain=ds0.example.net', actual: '%v'", txt)
	}
	if !warningsContains(cfg.Warnings, "myQstringParam") {
		t.Errorf("expected warning for malformed myQstringParam', actual: '%+v'", cfg.Warnings)
	}
	if !strings.Contains(txt, "# ds 'ds1'") {
		t.Errorf("expected comment with delivery service name, actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigCommentTopology(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: true, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")
	ds0.ProfileID = util.IntPtr(311)
	ds0.ProfileName = util.StrPtr("ds0Profile")

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.Topology = util.StrPtr("t0")
	ds1.ProfileID = util.IntPtr(312)
	ds1.ProfileName = util.StrPtr("ds1Profile")

	dses := []DeliveryService{*ds0, *ds1}

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

	serverParams := []tc.Parameter{
		{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "8",
			Profiles:   []byte(`["global"]`),
		},
	}

	server := makeTestParentServer()
	server.Cachegroup = util.StrPtr("edgeCG")
	server.CachegroupID = util.IntPtr(400)

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.CachegroupID = util.IntPtr(500)
	mid0.HostName = util.StrPtr("mymid")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG2")
	mid1.CachegroupID = util.IntPtr(501)
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")

	servers := []Server{*server, *mid0, *mid1}

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

	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

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

	dss := []DeliveryServiceServer{
		{
			Server:          *server.ID,
			DeliveryService: *ds0.ID,
		},
		{
			Server:          *server.ID,
			DeliveryService: *ds1.ID,
		},
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_domain=ds0.example.net") {
		t.Errorf("expected parent 'dest_domain=ds0.example.net', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "dest_domain=ds1.example.net") {
		t.Errorf("expected parent 'dest_domain=ds1.example.net', actual: '%v'", txt)
	}
	if !warningsContains(cfg.Warnings, "myQStringHandlingParam") {
		t.Errorf("expected warning for malformed myQStringHandlingParam', actual: '%+v'", cfg.Warnings)
	}

	if !strings.Contains(txt, "secondary_mode=1") {
		t.Errorf("expected default secondary_mode=1 for DSes without ParentConfigParamSecondaryMode parameter, actual: '%v'", txt)
	}
	if !strings.Contains(txt, `# ds 'ds1' topology 't0'`) {
		t.Errorf("expected comment with delivery service and topology, actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigHTTPSOrigin(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
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

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_domain=ds0.example.net port=80") {
		t.Errorf("expected edge parent.config of https origin to use internal http port 80 (not https/443), actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigHTTPSOriginTopology(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: true, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("https://ds0.example.net")
	ds0.ProfileID = util.IntPtr(311)
	ds0.ProfileName = util.StrPtr("ds0Profile")

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.Topology = util.StrPtr("t0")
	ds1.ProfileID = util.IntPtr(312)
	ds1.ProfileName = util.StrPtr("ds1Profile")

	dses := []DeliveryService{*ds0, *ds1}

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

	serverParams := []tc.Parameter{
		{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "8",
			Profiles:   []byte(`["global"]`),
		},
	}

	server := makeTestParentServer()
	server.Cachegroup = util.StrPtr("edgeCG")
	server.CachegroupID = util.IntPtr(400)

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.CachegroupID = util.IntPtr(500)
	mid0.HostName = util.StrPtr("mymid")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG2")
	mid1.CachegroupID = util.IntPtr(501)
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")

	servers := []Server{*server, *mid0, *mid1}

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

	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

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

	dss := []DeliveryServiceServer{
		{
			Server:          *server.ID,
			DeliveryService: *ds0.ID,
		},
		{
			Server:          *server.ID,
			DeliveryService: *ds1.ID,
		},
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_domain=ds0.example.net port=80") {
		t.Errorf("expected topology parent.config of https origin to be http/80 not https/443, actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigNoParentNoTopology(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment", GoDirect: true}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTPLive
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0-origin.example.net")

	dses := []DeliveryService{*ds0}

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
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_domain=ds0-origin.example.net") {
		t.Errorf("expected parent 'dest_domain=ds0-origin.example.net', actual: '%v'", txt)
	}

	lines := strings.Split(txt, "\n")
	for _, line := range lines {
		if !strings.Contains(line, "dest_domain=ds0-origin.example.net") {
			continue
		}
		if !strings.Contains(line, `parent="ds0-origin.example.net:80`) {
			t.Errorf("expected non-topology DS of type not using parents to have parent=origin directive, actual: '%v'", txt)
		}
		if !strings.Contains(line, `go_direct=true`) {
			t.Errorf("expected non-topology DS of type not using parents to have go_direct=true directive, actual: '%v'", txt)
		}
		if !strings.Contains(line, `parent_is_proxy=false`) {
			t.Errorf("expected non-topology DS of type not using parents to have parent_is_proxy=false directive, actual: '%v'", txt)
		}
	}
}

func TestMakeParentDotConfigHTTPSOriginTopologyNoPrimaryParent(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: true, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("https://ds0.example.net")
	ds0.ProfileID = util.IntPtr(311)
	ds0.ProfileName = util.StrPtr("ds0Profile")

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.Topology = util.StrPtr("t0")
	ds1.ProfileID = util.IntPtr(312)
	ds1.ProfileName = util.StrPtr("ds1Profile")

	dses := []DeliveryService{*ds0, *ds1}

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

	serverParams := []tc.Parameter{
		{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "8",
			Profiles:   []byte(`["global"]`),
		},
	}

	server := makeTestParentServer()
	server.Cachegroup = util.StrPtr("edgeCG")
	server.CachegroupID = util.IntPtr(400)

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.CachegroupID = util.IntPtr(500)
	mid0.HostName = util.StrPtr("mymid")
	mid0.ID = util.IntPtr(45)
	mid0.Status = util.StrPtr(string(tc.CacheStatusAdminDown))
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG2")
	mid1.CachegroupID = util.IntPtr(501)
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")

	servers := []Server{*server, *mid0, *mid1}

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

	serverCapabilities := map[int]map[ServerCapability]struct{}{}
	dsRequiredCapabilities := map[int]map[ServerCapability]struct{}{}

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

	dss := []DeliveryServiceServer{
		{
			Server:          *server.ID,
			DeliveryService: *ds0.ID,
		},
		{
			Server:          *server.ID,
			DeliveryService: *ds1.ID,
		},
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, `parent="mymid1.mydomain.example.net:80|0.999"`) {
		t.Errorf("expected topology parent.config withparent=\"mymid1.mydomain.example.net:80|0.999\", actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigMergeParentGroupTopology(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: true, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://ds0.example.net")
	ds0.ProfileID = util.IntPtr(311)
	ds0.ProfileName = util.StrPtr("ds0Profile")
	ds0.Topology = util.StrPtr("t0")

	dses := []DeliveryService{*ds0}

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
		{
			Name:       ParentConfigParamMergeGroups,
			ConfigFile: "parent.config",
			Value:      "oplCG0 oplCG1",
			Profiles:   []byte(`["ds0Profile"]`),
		},
	}

	serverParams := []tc.Parameter{
		{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "8",
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
	setIP(mid0, "192.168.2.3")

	opl0 := makeTestParentServer()
	opl0.Cachegroup = util.StrPtr("oplCG0")
	opl0.CachegroupID = util.IntPtr(600)
	opl0.HostName = util.StrPtr("myopl0")
	opl0.ID = util.IntPtr(47)
	setIP(opl0, "192.168.2.4")

	opl1 := makeTestParentServer()
	opl1.Cachegroup = util.StrPtr("oplCG1")
	opl1.CachegroupID = util.IntPtr(601)
	opl1.HostName = util.StrPtr("myopl1")
	opl1.ID = util.IntPtr(48)
	setIP(opl0, "192.168.2.5")

	servers := []Server{*edge, *mid0, *mid1, *opl0, *opl1}

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
				},
				{
					Cachegroup: "oplCG1",
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

	oCG0 := &tc.CacheGroupNullable{}
	oCG0.Name = opl0.Cachegroup
	oCG0.ID = opl0.CachegroupID
	oCGType0 := tc.CacheGroupMidTypeName
	oCG0.Type = &oCGType0

	oCG1 := &tc.CacheGroupNullable{}
	oCG1.Name = opl1.Cachegroup
	oCG1.ID = opl1.CachegroupID
	oCGType1 := tc.CacheGroupMidTypeName
	oCG1.Type = &oCGType1

	cgs := []tc.CacheGroupNullable{*eCG, *mCG0, *mCG1, *oCG0, *oCG1}

	dss := []DeliveryServiceServer{
		{
			Server:          *edge.ID,
			DeliveryService: *ds0.ID,
		},
		{
			Server:          *mid0.ID,
			DeliveryService: *ds0.ID,
		},
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	{ // test edge config
		cfg, err := MakeParentDotConfig(dses, edge, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, hdr.HdrComment)

		if !strings.Contains(txt, `dest_domain=ds0.example.net port=80 parent="mymid0.mydomain.example.net:80|0.999" secondary_parent="mymid1.mydomain.example.net:80|0.999"`) {
			t.Errorf("expected topology parent.config of ds0 edge to have parent only: '%v'", txt)
		}

		if strings.Count(txt, "secondary_parent") != 2 {
			t.Errorf("expected 2 secondary parents for edge (dest_domain=.): '%v'", txt)
		}
	}

	{ // test mid config
		cfg, err := MakeParentDotConfig(dses, mid0, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, hdr.HdrComment)

		if !strings.Contains(txt, `dest_domain=ds0.example.net port=80 parent="myopl0.mydomain.example.net:80|0.999;myopl1.mydomain.example.net:80|0.999"`) {
			t.Errorf("expected topology parent.config of ds0 mid1 to have parent only: '%v'", txt)
		} else if strings.Count(txt, "secondary_parent") != 1 {
			t.Errorf("expected one secondary parent for mid1 (dest_domain=.): '%v'", txt)
		}
	}
}

func TestMakeParentDotConfigTopologiesServerMultipleProfileParams(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

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
			Value:      "consistent_hash",
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
		tc.Parameter{
			Name:       ParentConfigCacheParamWeight,
			ConfigFile: "parent.config",
			Value:      "100",
			Profiles:   []byte(`["serverprofile1"]`),
		},
		tc.Parameter{
			Name:       ParentConfigCacheParamWeight,
			ConfigFile: "parent.config",
			Value:      "200",
			Profiles:   []byte(`["serverprofile0"]`),
		},
	}

	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "8",
			Profiles:   []byte(`["global"]`),
		},
		tc.Parameter{
			Name:       ParentConfigCacheParamWeight,
			ConfigFile: "parent.config",
			Value:      "100",
			Profiles:   []byte(`["serverprofile0"]`),
		},
		tc.Parameter{
			Name:       ParentConfigCacheParamWeight,
			ConfigFile: "parent.config",
			Value:      "200",
			Profiles:   []byte(`["serverprofile1"]`),
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
	origin0.ProfileNames = []string{"serverprofile0", "serverprofile1"}

	origin1 := makeTestParentServer()
	origin1.Cachegroup = util.StrPtr("originCG")
	origin1.CachegroupID = util.IntPtr(500)
	origin1.HostName = util.StrPtr("myorigin1")
	origin1.ID = util.IntPtr(46)
	setIP(origin1, "192.168.2.3")
	origin1.Type = tc.OriginTypeName
	origin1.TypeID = util.IntPtr(991)
	origin1.ProfileNames = []string{"serverprofile1", "serverprofile0"}

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
		DeliveryServiceServer{
			Server:          *origin1.ID,
			DeliveryService: *ds1.ID,
		},
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "myorigin0.mydomain.example.net:80|200") {
		t.Errorf("expected origin 0 with profiles [0,1] to have weight 200 from profile 1, actual: '%v'", txt)
	}
	if !strings.Contains(txt, "myorigin1.mydomain.example.net:80|100") {
		t.Errorf("expected origin 0 with profiles [1,0] to have weight 100 from profile 0, actual: '%v'", txt)
	}
}

func TestMakeParentDotConfigFirstLastNoTopo(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: true, HdrComment: "myHeaderComment", GoDirect: true}

	// Non Toplogy ds
	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("https://ds0.example.net")
	ds0.ProfileID = util.IntPtr(311)
	ds0.ProfileName = util.StrPtr("ds0Profile")

	// Non Toplogy ds, MSO
	ds1 := makeParentDS()
	ds1Type := tc.DSTypeHTTP
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds1.OrgServerFQDN = util.StrPtr("https://ds1.example.net")
	ds1.ProfileID = util.IntPtr(312)
	ds1.ProfileName = util.StrPtr("ds0Profile")
	ds1.MultiSiteOrigin = util.BoolPtr(true)

	dses := []DeliveryService{*ds0, *ds1}

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
		ParentConfigGoDirectEdge:                               "false",
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
	setIP(mid0, "192.168.2.3")

	org0 := makeTestParentServer()
	org0.Cachegroup = util.StrPtr("orgCG0")
	org0.CachegroupID = util.IntPtr(502)
	org0.HostName = util.StrPtr("myorg0")
	org0.ID = util.IntPtr(47)
	setIP(org0, "192.168.2.4")
	org0.Type = tc.OriginTypeName
	org0.TypeID = util.IntPtr(991)

	org1 := makeTestParentServer()
	org1.Cachegroup = util.StrPtr("orgCG1")
	org1.CachegroupID = util.IntPtr(503)
	org1.HostName = util.StrPtr("myorg1")
	org1.ID = util.IntPtr(48)
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

	{ // test edge config
		cfg, err := MakeParentDotConfig(dses, edge, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, hdr.HdrComment)

		needs := []string{
			` secondary_mode=2`,
			` round_robin=consistent_hash`,
			` go_direct=false`,
			` parent_is_proxy=true`,
			` parent_retry=both`,
			` max_simple_retries=12`,
			` max_unavailable_server_retries=22`,
			` simple_server_retry_responses="401,402"`,
			` unavailable_server_retry_responses="501,502"`,
		}

		dsstrs := []string{
			`dest_domain=ds0.example.net `,
			`dest_domain=ds1.example.net `,
		}

		for _, dsstr := range dsstrs {
			cnt := strings.Count(txt, dsstr)
			if 1 != cnt {
				t.Errorf("Expected one entry for %s got %d\n%v", dsstr, cnt, txt)
			} else {
				lines := strings.Split(txt, "\n")
				dsline := lineWhichContains(lines, dsstr)
				missing := missingFrom(dsline, needs)
				if 0 < len(missing) {
					t.Errorf("Missing required string(s) from line: %v\n%v (warnings: %v)", missing, dsline, cfg.Warnings)
				}
			}
		}
	}

	{ // test mid config
		cfg, err := MakeParentDotConfig(dses, mid0, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, hdr.HdrComment)

		needs := []string{
			` round_robin=true`,
			` go_direct=true`,
			` parent_is_proxy=false`,
			` parent_retry=both`,
			` max_simple_retries=14`,
			` max_unavailable_server_retries=24`,
			` simple_server_retry_responses="401,404"`,
			` unavailable_server_retry_responses="501,504"`,
		}

		{
			dsstr := "dest_domain=ds1.example.net"
			cnt := strings.Count(txt, dsstr)
			if 1 != cnt {
				t.Errorf("Expected one entry for %s got %d\n%v", dsstr, cnt, txt)
			} else {
				lines := strings.Split(txt, "\n")
				dsline := lineWhichContains(lines, dsstr)
				missing := missingFrom(dsline, needs)
				if 0 < len(missing) {
					t.Errorf("Missing required string(s) from line: %v\n%v (warnings: %v)", missing, dsline, cfg.Warnings)
				}
			}
		}

		// Check parent ordering (based on cache group prim/sec parents)
		if !strings.Contains(txt, `parent="myorg0`) {
			t.Errorf("Incorrect parent ordering, got %v", txt)
		}
	}

	{
		cfg, err := MakeParentDotConfig(dses, mid1, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		// Check parent ordering (based on cache group prim/sec parents)
		if !strings.Contains(txt, `parent="myorg1`) {
			t.Errorf("Incorrect parent ordering, got %v", txt)
		}
	}
}

func TestMakeParentDotConfigFirstInnerLastTopology(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: true, HdrComment: "myHeaderComment", GoDirect: true}

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
	ds1.ID = util.IntPtr(44)
	ds1Type := tc.DSTypeHTTP
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.ProfileID = util.IntPtr(311)
	ds1.ProfileName = util.StrPtr("ds0Profile")
	ds1.Topology = util.StrPtr("t0")

	dses := []DeliveryService{*ds0, *ds1}

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
		ParentConfigGoDirectFirst:                              "false",
		ParentConfigGoDirectInner:                              "false",
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
	edge.HostName = util.StrPtr("edge")

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG0")
	mid0.CachegroupID = util.IntPtr(500)
	mid0.HostName = util.StrPtr("mid0")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG1")
	mid1.CachegroupID = util.IntPtr(501)
	mid1.HostName = util.StrPtr("mid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")

	opl0 := makeTestParentServer()
	opl0.Cachegroup = util.StrPtr("oplCG0")
	opl0.CachegroupID = util.IntPtr(502)
	opl0.HostName = util.StrPtr("opl0")
	opl0.ID = util.IntPtr(47)
	setIP(opl0, "192.168.2.4")

	opl1 := makeTestParentServer()
	opl1.Cachegroup = util.StrPtr("oplCG1")
	opl1.CachegroupID = util.IntPtr(503)
	opl1.HostName = util.StrPtr("opl1")
	opl1.ID = util.IntPtr(48)
	setIP(opl1, "192.168.2.5")

	org0 := makeTestParentServer()
	org0.Cachegroup = util.StrPtr("orgCG0")
	org0.CachegroupID = util.IntPtr(504)
	org0.HostName = util.StrPtr("org0")
	org0.ID = util.IntPtr(49)
	setIP(org0, "192.168.2.6")
	org0.Type = tc.OriginTypeName
	org0.TypeID = util.IntPtr(991)

	org1 := makeTestParentServer()
	org1.Cachegroup = util.StrPtr("orgCG1")
	org1.CachegroupID = util.IntPtr(505)
	org1.HostName = util.StrPtr("org1")
	org1.ID = util.IntPtr(50)
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
		/*
			{Server: *edge.ID, DeliveryService: *ds0.ID},
			{Server: *mid0.ID, DeliveryService: *ds0.ID},
			{Server: *mid1.ID, DeliveryService: *ds0.ID},
			{Server: *opl0.ID, DeliveryService: *ds0.ID},
			{Server: *opl1.ID, DeliveryService: *ds0.ID},
			{Server: *org0.ID, DeliveryService: *ds0.ID},
			{Server: *org1.ID, DeliveryService: *ds0.ID},

			{Server: *edge.ID, DeliveryService: *ds1.ID},
			{Server: *mid0.ID, DeliveryService: *ds1.ID},
			{Server: *mid1.ID, DeliveryService: *ds1.ID},
			{Server: *opl0.ID, DeliveryService: *ds1.ID},
			{Server: *opl1.ID, DeliveryService: *ds1.ID},
			{Server: *org0.ID, DeliveryService: *ds1.ID},
			{Server: *org1.ID, DeliveryService: *ds1.ID},
		*/
	}
	cdn := &tc.CDN{
		DomainName: "cdndomain.example",
		Name:       "my-cdn-name",
	}

	dsstrs := []string{
		`dest_domain=ds0.example.net `,
		`dest_domain=ds1.example.net `,
	}

	{ // test edge config
		cfg, err := MakeParentDotConfig(dses, edge, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, hdr.HdrComment)

		needs := []string{
			` secondary_mode=2`,
			` round_robin=consistent_hash`,
			` go_direct=false`,
			` parent_is_proxy=true`,
			` parent_retry=both`,
			` max_simple_retries=12`,
			` max_unavailable_server_retries=22`,
			` simple_server_retry_responses="401,402"`,
			` unavailable_server_retry_responses="501,502"`,
		}

		for _, dsstr := range dsstrs {
			cnt := strings.Count(txt, dsstr)
			if 1 != cnt {
				t.Errorf("Expected one entry for %s got %d\n%v", dsstr, cnt, txt)
			} else {
				lines := strings.Split(txt, "\n")
				dsline := lineWhichContains(lines, dsstr)
				missing := missingFrom(dsline, needs)
				if 0 < len(missing) {
					t.Errorf("Missing required string(s) from line: %v\n%v", missing, dsline)
				}
			}
		}
	}

	{ // test mid config
		cfg, err := MakeParentDotConfig(dses, mid0, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, hdr.HdrComment)

		needs := []string{
			` round_robin=consistent_hash`,
			` go_direct=false`,
			` parent_is_proxy=true`,
			` parent_retry=both`,
			` max_simple_retries=13`,
			` max_unavailable_server_retries=23`,
			` simple_server_retry_responses="401,403"`,
			` unavailable_server_retry_responses="501,503"`,
		}

		for _, dsstr := range dsstrs {
			cnt := strings.Count(txt, dsstr)
			if 1 != cnt {
				t.Errorf("Expected one entry for %s got %d\n%v", dsstr, cnt, txt)
			} else {
				lines := strings.Split(txt, "\n")
				dsline := lineWhichContains(lines, dsstr)
				missing := missingFrom(dsline, needs)
				if 0 < len(missing) {
					t.Errorf("Missing required string(s) from line: %v\n%v", missing, dsline)
				}
			}
		}
	}

	{ // test opl config
		cfg, err := MakeParentDotConfig(dses, opl0, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, hdr.HdrComment)

		needs := []string{
			` round_robin=true`,
			` go_direct=true`,
			` parent_is_proxy=false`,
			` parent_retry=both`,
			` max_simple_retries=14`,
			` max_unavailable_server_retries=24`,
			` simple_server_retry_responses="401,404"`,
			` unavailable_server_retry_responses="501,504"`,
		}

		dsstr := `dest_domain=ds1.example.net `
		cnt := strings.Count(txt, dsstr)
		if 1 != cnt {
			t.Errorf("Expected one entry for %s got %d\n%v", dsstr, cnt, txt)
		} else {
			lines := strings.Split(txt, "\n")
			dsline := lineWhichContains(lines, dsstr)
			missing := missingFrom(dsline, needs)
			if 0 < len(missing) {
				t.Errorf("Missing required string(s) from line: %v\n%v", missing, dsline)
			}
		}
	}
}

func TestMakeParentDotConfigOptVersion(t *testing.T) {
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
			Value:      "consistent_hash",
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

	// unavailable_server_retry_responses is not available as a feature in ATS 5, but is in ATS 9.

	t.Run("Package Parameter 9 with no Opt ATSVersion has ATS 9 feature", func(t *testing.T) {
		serverParams := []tc.Parameter{
			tc.Parameter{
				Name:       "trafficserver",
				ConfigFile: "package",
				Value:      "9",
				Profiles:   []byte(`["global"]`),
			},
		}

		opt := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

		cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, opt.HdrComment)

		if !strings.Contains(txt, `unavailable_server_retry_responses="400,503"`) {
			t.Errorf(`expected Package Parameter ATS 9 with no Opt ATS Version to have unavailable_server_retry_responses feature', actual: '%v'`, txt)
		}
	})

	t.Run("Package Parameter 5 with no Opt ATSVersion does not have the feature it shouldn't", func(t *testing.T) {
		serverParams := []tc.Parameter{
			tc.Parameter{
				Name:       "trafficserver",
				ConfigFile: "package",
				Value:      "5",
				Profiles:   []byte(`["global"]`),
			},
		}

		opt := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

		cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, opt.HdrComment)

		if strings.Contains(txt, `unavailable_server_retry_responses`) {
			t.Errorf(`expected Package Parameter ATS 5 with no Opt ATS Version to not have unavailable_server_retry_responses feature', actual: '%v'`, txt)
		}
	})

	t.Run("Package Parameter 5 with Opt ATSVersion 9 uses Opt not Param.", func(t *testing.T) {
		serverParams := []tc.Parameter{
			tc.Parameter{
				Name:       "trafficserver",
				ConfigFile: "package",
				Value:      "5",
				Profiles:   []byte(`["global"]`),
			},
		}

		opt := &ParentConfigOpts{
			AddComments:     false,
			HdrComment:      "myHeaderComment",
			ATSMajorVersion: 9,
		}

		cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, opt.HdrComment)

		if !strings.Contains(txt, `unavailable_server_retry_responses`) {
			t.Errorf(`expected Package Parameter ATS 5 with Opt ATS Version 9 to use Opt not Parameter with ATS 9 unavailable_server_retry_responses feature', actual: '%v'`, txt)
		}
	})

	t.Run("Package Parameter 9 with Opt ATSVersion 5 uses Opt not Param.", func(t *testing.T) {
		serverParams := []tc.Parameter{
			tc.Parameter{
				Name:       "trafficserver",
				ConfigFile: "package",
				Value:      "9",
				Profiles:   []byte(`["global"]`),
			},
		}

		opt := &ParentConfigOpts{
			AddComments:     false,
			HdrComment:      "myHeaderComment",
			ATSMajorVersion: 5,
		}

		cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, opt)
		if err != nil {
			t.Fatal(err)
		}
		txt := cfg.Text

		testComment(t, txt, opt.HdrComment)

		if strings.Contains(txt, `unavailable_server_retry_responses`) {
			t.Errorf(`expected Package Parameter ATS 9 with Opt ATS Version 5 to use Opt not Parameter with no ATS 9 unavailable_server_retry_responses feature', actual: '%v'`, txt)
		}
	})
}

func TestMakeParentDotConfigOriginIP(t *testing.T) {
	hdr := &ParentConfigOpts{AddComments: false, HdrComment: "myHeaderComment"}

	ds0 := makeParentDS()
	ds0Type := tc.DSTypeHTTP
	ds0.Type = &ds0Type
	ds0.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreUseInCacheKeyAndPassUp))
	ds0.OrgServerFQDN = util.StrPtr("http://192.0.2.42")

	ds1 := makeParentDS()
	ds1.ID = util.IntPtr(43)
	ds1Type := tc.DSTypeDNS
	ds1.Type = &ds1Type
	ds1.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds1.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	ds1.Topology = util.StrPtr("t0")

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
	server.Cachegroup = util.StrPtr("edgeCG")
	server.CachegroupID = util.IntPtr(400)

	mid0 := makeTestParentServer()
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.CachegroupID = util.IntPtr(500)
	mid0.HostName = util.StrPtr("mymid")
	mid0.ID = util.IntPtr(45)
	setIP(mid0, "192.168.2.2")

	mid1 := makeTestParentServer()
	mid1.Cachegroup = util.StrPtr("midCG")
	mid1.CachegroupID = util.IntPtr(500)
	mid1.HostName = util.StrPtr("mymid1")
	mid1.ID = util.IntPtr(46)
	setIP(mid1, "192.168.2.3")

	servers := []Server{*server, *mid0, *mid1}

	topologies := []tc.Topology{
		tc.Topology{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				tc.TopologyNode{
					Cachegroup: "edgeCG",
					Parents:    []int{1},
				},
				tc.TopologyNode{
					Cachegroup: "midCG",
				},
			},
		},
	}

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

	cfg, err := MakeParentDotConfig(dses, server, servers, topologies, serverParams, parentConfigParams, serverCapabilities, dsRequiredCapabilities, cgs, dss, cdn, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr.HdrComment)

	if !strings.Contains(txt, "dest_ip=192.0.2.42") {
		t.Errorf("expected parent 'dest_ip=192.0.2.42', actual: '%v'", txt)
	}
	if strings.Contains(txt, "dest_domain=192.0.2.42") {
		t.Errorf("expected parent to not contain dest_domain for an IP, actual: '%v'", txt)
	}
	if !strings.Contains(txt, "dest_domain=ds1.example.net") {
		t.Errorf("expected parent 'dest_domain=ds1.example.net', actual: '%v'", txt)
	}
	if !warningsContains(cfg.Warnings, "myQStringHandlingParam") {
		t.Errorf("expected warning for malformed myQStringHandlingParam', actual: '%+v'", cfg.Warnings)
	}
	if strings.Contains(txt, "# topology") {
		// ATS doesn't support inline comments in parent.config
		t.Errorf("expected: no inline '# topology' comment, actual: '%v'", txt)
	}
}

// returns which elements in "needs" is missing from "line"
func missingFrom(line string, needs []string) []string {
	misses := []string{}
	for _, need := range needs {
		if !strings.Contains(line, need) {
			misses = append(misses, need)
		}
	}
	return misses
}

// returns first line which contains
func lineWhichContains(lines []string, str string) (res string) {
	for _, line := range lines {
		if strings.Contains(line, str) {
			res = line
			break
		}
	}
	return res
}

// warningsContains returns whether the given warnings has str as a substring of any warning.
// Note this is different than lib/go-util.ContainsStr, which only returns if the array has the exact value as one of its values.
func warningsContains(warnings []string, str string) bool {
	for _, warn := range warnings {
		if strings.Contains(warn, str) {
			return true
		}
	}
	return false
}

func makeTestParentServer() *Server {
	server := &Server{}
	server.CDNName = util.StrPtr("myCDN")
	server.Cachegroup = util.StrPtr("cg0")
	server.CachegroupID = util.IntPtr(422)
	server.DomainName = util.StrPtr("mydomain.example.net")
	server.CDNID = util.IntPtr(43)
	server.HostName = util.StrPtr("myserver")
	server.HTTPSPort = util.IntPtr(12443)
	server.ID = util.IntPtr(44)
	setIP(server, "192.168.2.1")
	server.ProfileNames = []string{"serverprofile"}
	server.TCPPort = util.IntPtr(80)
	server.Type = "EDGE"
	server.TypeID = util.IntPtr(91)
	status := string(tc.CacheStatusReported)
	server.Status = &status
	server.StatusID = util.IntPtr(99)
	return server
}

func makeParentDS() *DeliveryService {
	ds := &DeliveryService{}
	ds.ID = util.IntPtr(42)
	ds.XMLID = util.StrPtr("ds1")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	dsType := tc.DSTypeDNS
	ds.Type = &dsType
	ds.MultiSiteOrigin = util.BoolPtr(false)
	return ds
}
