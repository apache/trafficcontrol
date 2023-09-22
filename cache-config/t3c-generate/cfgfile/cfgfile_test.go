package cfgfile

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
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/cache-config/t3c-generate/config"
	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/test"
)

func TestWriteConfigs(t *testing.T) {
	buf := &bytes.Buffer{}
	configs := []t3cutil.ATSConfigFile{
		{
			Name:        "config0.txt",
			Path:        "/my/config0/location",
			Text:        "config0",
			Secure:      false,
			ContentType: "text/plain",
		},
		{
			Name:        "config1.txt",
			Path:        "/my/config1/location",
			Text:        "config2,foo",
			Secure:      false,
			ContentType: "text/csv",
		},
	}

	if err := WriteConfigs(configs, buf); err != nil {
		t.Fatalf("WriteConfigs error expected nil, actual: %v", err)
	}

	actual := buf.String()

	expected0 := `[{"name":"config0.txt","path":"/my/config0/location","content_type":"text/plain","line_comment":"","secure":false,"text":"config0","warnings":null},{"name":"config1.txt","path":"/my/config1/location","content_type":"text/csv","line_comment":"","secure":false,"text":"config2,foo","warnings":null}]`

	if !strings.Contains(actual, expected0) {
		t.Errorf("WriteConfigs expected '%v' actual '%v'", expected0, actual)
	}

	expected1 := `[{"name":"config0.txt","path":"/my/config0/location","content_type":"text/plain","line_comment":"","secure":false,"text":"config0","warnings":null},{"name":"config1.txt","path":"/my/config1/location","content_type":"text/csv","line_comment":"","secure":false,"text":"config2,foo","warnings":null}]`
	if !strings.Contains(actual, expected1) {
		t.Errorf("WriteConfigs expected config1 '%v' actual '%v'", expected1, actual)
	}

	expectedPrefix := `[{"name":"config0.txt","path":"/my/config0/location","content_type":"text/plain","line_comment":"","secure":false,"text":"config0","warnings":null},{"name":"config1.txt","path":"/my/config1/location","content_type":"text/csv","line_comment":"","secure":false,"text":"config2,foo","warnings":null}]`
	if !strings.HasPrefix(actual, expectedPrefix) {
		t.Errorf("WriteConfigs expected prefix '%v' actual '%v'", expectedPrefix, actual)
	}
}

// TestGetAllConfigsWriteConfigsDeterministic tests that WriteConfigs(GetAllConfigs()) is Deterministic.
// That is, that for the same input, it always produces the same output.
//
// Because Go map iteration is defined to be random, running it multiple times even on the exact same input could be different, if there's a determinism bug.
// But beyond that, we re-order slices whose order isn't semantically significant (e.g. params) and run it again.
func TestGetAllConfigsWriteConfigsDeterministic(t *testing.T) {
	// TODO expand fake data. Currently, it's only making a remap.config.
	toData := MakeFakeTOData()
	revalOnly := false
	cfgPath := "/etc/trafficserver/"

	cfg := config.Cfg{}
	cfg.Dir = cfgPath
	cfg.RevalOnly = revalOnly

	configs, err := GetAllConfigs(toData, cfg)
	if err != nil {
		t.Fatalf("error getting configs: " + err.Error())
	}
	buf := &bytes.Buffer{}
	if err := WriteConfigs(configs, buf); err != nil {
		t.Fatalf("error writing configs: " + err.Error())
	}
	configStr := buf.String()
	configStr = removeComments(configStr)

	for i := 0; i < 10; i++ {
		configs2, err := GetAllConfigs(toData, cfg)
		if err != nil {
			t.Fatalf("error getting configs2: " + err.Error())
		}
		buf := &bytes.Buffer{}
		if err := WriteConfigs(configs2, buf); err != nil {
			t.Fatalf("error writing configs2: " + err.Error())
		}
		configStr2 := buf.String()

		configStr2 = removeComments(configStr2)

		if configStr != configStr2 {
			// This doesn't actually need to be fatal; but if there are differences, we don't want to spam the error 10 times.
			t.Fatalf("multiple configs with the same data expected to be deterministically the same, actual '''%v''' and '''%v'''", configStr, configStr2)
		}
	}
}

func removeComments(configs string) string {
	lines := strings.Split(configs, "\n")
	newLines := []string{}
	for _, line := range lines {
		if strings.Contains(line, "DO NOT EDIT") {
			continue
		}
		newLines = append(newLines, line)
	}
	return strings.Join(newLines, "\n")
}

func randDSS() tc.DeliveryServiceServerV5 {
	return tc.DeliveryServiceServerV5{
		Server:          util.Ptr(test.RandInt()),
		DeliveryService: util.Ptr(test.RandInt()),
	}
}

func randDS() *atscfg.DeliveryService {
	deepCachingTypeNever := tc.DeepCachingTypeNever
	dsTypeHTTP := util.Ptr("HTTP")
	protocol := tc.DSProtocolHTTP
	ds := atscfg.DeliveryService{}
	ds.EcsEnabled = test.RandBool()
	ds.RangeSliceBlockSize = util.Ptr(test.RandInt())
	ds.ConsistentHashRegex = util.Ptr(test.RandStr())
	ds.ConsistentHashQueryParams = []string{
		test.RandStr(),
		test.RandStr(),
	}
	ds.MaxOriginConnections = util.Ptr(test.RandInt())
	ds.DeepCachingType = deepCachingTypeNever
	ds.FQPacingRate = util.Ptr(test.RandInt())
	ds.SigningAlgorithm = util.Ptr(test.RandStr())
	ds.Tenant = util.Ptr(test.RandStr())
	ds.TRResponseHeaders = util.Ptr(test.RandStr())
	ds.TRRequestHeaders = util.Ptr(test.RandStr())
	ds.Active = randActiveState()
	ds.AnonymousBlockingEnabled = test.RandBool()
	ds.CCRDNSTTL = util.Ptr(test.RandInt())
	ds.CDNID = test.RandInt()
	ds.CDNName = util.Ptr(test.RandStr())
	ds.CheckPath = util.Ptr(test.RandStr())
	ds.DisplayName = test.RandStr()
	ds.DNSBypassCNAME = util.Ptr(test.RandStr())
	ds.DNSBypassIP = util.Ptr(test.RandStr())
	ds.DNSBypassIP6 = util.Ptr(test.RandStr())
	ds.DNSBypassTTL = util.Ptr(test.RandInt())
	ds.DSCP = test.RandInt()
	ds.EdgeHeaderRewrite = util.Ptr(test.RandStr())
	ds.GeoLimit = test.RandInt()
	ds.GeoLimitCountries = nil
	ds.GeoLimitRedirectURL = util.Ptr(test.RandStr())
	ds.GeoProvider = test.RandInt()
	ds.GlobalMaxMBPS = util.Ptr(test.RandInt())
	ds.GlobalMaxTPS = util.Ptr(test.RandInt())
	ds.HTTPBypassFQDN = util.Ptr(test.RandStr())
	ds.ID = util.Ptr(test.RandInt())
	ds.InfoURL = util.Ptr(test.RandStr())
	ds.InitialDispersion = util.Ptr(test.RandInt())
	ds.IPV6RoutingEnabled = util.Ptr(test.RandBool())
	ds.LastUpdated = time.Now()
	ds.LogsEnabled = test.RandBool()
	ds.LongDesc = test.RandStr()
	ds.MatchList = []tc.DeliveryServiceMatch{
		{
			Type:      tc.DSMatchTypeHostRegex,
			SetNumber: 0,
			Pattern:   `\.*foo\.*`,
		},
	}
	ds.MaxDNSAnswers = util.Ptr(test.RandInt())
	ds.MidHeaderRewrite = util.Ptr(test.RandStr())
	ds.MissLat = util.Ptr(test.RandFloat64())
	ds.MissLong = util.Ptr(test.RandFloat64())
	ds.MultiSiteOrigin = test.RandBool()
	ds.OriginShield = util.Ptr(test.RandStr())
	ds.OrgServerFQDN = util.Ptr("http://" + test.RandStr())
	ds.ProfileDesc = util.Ptr(test.RandStr())
	ds.ProfileID = util.Ptr(test.RandInt())
	ds.ProfileName = util.Ptr(test.RandStr())
	ds.Protocol = &protocol
	ds.QStringIgnore = util.Ptr(test.RandInt())
	ds.RangeRequestHandling = util.Ptr(test.RandInt())
	ds.RegexRemap = util.Ptr(test.RandStr())
	ds.RegionalGeoBlocking = test.RandBool()
	ds.RemapText = util.Ptr(test.RandStr())
	ds.RoutingName = test.RandStr()
	ds.Signed = *util.Ptr(test.RandBool())
	ds.SSLKeyVersion = util.Ptr(test.RandInt())
	ds.TenantID = test.RandInt()
	ds.Type = dsTypeHTTP
	ds.TypeID = test.RandInt()
	ds.XMLID = test.RandStr()
	ds.ExampleURLs = []string{
		test.RandStr(),
		test.RandStr(),
	}
	return &ds
}

func randActiveState() tc.DeliveryServiceActiveState {
	states := []tc.DeliveryServiceActiveState{
		tc.DSActiveStateActive,
		tc.DSActiveStateInactive,
		tc.DSActiveStatePrimed,
		tc.DSActiveStateInactive,
		tc.DSActiveStatePrimed,
		tc.DSActiveStateActive,
	}
	n := test.RandIntForActive()
	return states[n]
}
func randServer() *atscfg.Server {
	sv := &atscfg.Server{}
	sv.CacheGroup = test.RandStr()
	sv.CacheGroupID = test.RandInt()
	sv.CDNID = test.RandInt()
	sv.CDN = test.RandStr()
	sv.DomainName = test.RandStr()
	sv.GUID = util.Ptr(test.RandStr())
	sv.HostName = test.RandStr()
	sv.HTTPSPort = util.Ptr(test.RandInt())
	sv.ID = test.RandInt()
	sv.ILOIPAddress = util.Ptr(test.RandStr())
	sv.ILOIPGateway = util.Ptr(test.RandStr())
	sv.ILOIPNetmask = util.Ptr(test.RandStr())
	sv.ILOPassword = util.Ptr(test.RandStr())
	sv.ILOUsername = util.Ptr(test.RandStr())

	sv.Interfaces = []tc.ServerInterfaceInfoV40{}
	{
		ssi := tc.ServerInterfaceInfoV40{}
		ssi.Name = test.RandStr()
		ssi.IPAddresses = []tc.ServerIPAddress{
			tc.ServerIPAddress{
				Address:        test.RandStr(),
				Gateway:        util.Ptr(test.RandStr()),
				ServiceAddress: true,
			},
			tc.ServerIPAddress{
				Address:        test.RandStr(),
				Gateway:        util.Ptr(test.RandStr()),
				ServiceAddress: true,
			},
		}
		sv.Interfaces = append(sv.Interfaces, ssi)
	}

	sv.LastUpdated = time.Now()
	sv.MgmtIPAddress = util.Ptr(test.RandStr())
	sv.MgmtIPGateway = util.Ptr(test.RandStr())
	sv.MgmtIPNetmask = util.Ptr(test.RandStr())
	sv.OfflineReason = util.Ptr(test.RandStr())
	sv.PhysicalLocation = test.RandStr()
	sv.PhysicalLocationID = test.RandInt()
	sv.Profiles = []string{test.RandStr()}
	sv.Rack = util.Ptr(test.RandStr())
	sv.RevalApplyTime = util.Ptr(time.Now())
	sv.RevalUpdateTime = util.Ptr(time.Now())
	sv.RevalUpdateFailed = test.RandBool()
	sv.Status = test.RandStr()
	sv.StatusID = test.RandInt()
	sv.TCPPort = util.Ptr(test.RandInt())
	sv.Type = test.RandStr()
	sv.TypeID = test.RandInt()
	sv.ConfigUpdateTime = util.Ptr(time.Now())
	sv.ConfigApplyTime = util.Ptr(time.Now())
	sv.ConfigUpdateFailed = test.RandBool()
	sv.XMPPID = util.Ptr(test.RandStr())
	sv.XMPPPasswd = util.Ptr(test.RandStr())
	return sv
}

func randCacheGroup() *tc.CacheGroupNullableV5 {
	return &tc.CacheGroupNullableV5{
		ID:        util.Ptr(test.RandInt()),
		Name:      util.Ptr(test.RandStr()),
		ShortName: util.Ptr(test.RandStr()),
		Latitude:  util.Ptr(test.RandFloat64()),
		Longitude: util.Ptr(test.RandFloat64()),
		// ParentName:                  util.StrPtr(test.RandStr()),
		// ParentCachegroupID:          util.IntPtr(test.RandInt()),
		// SecondaryParentName:         util.StrPtr(test.RandStr()),
		// SecondaryParentCachegroupID: util.IntPtr(test.RandInt()),
		FallbackToClosest: util.Ptr(test.RandBool()),
		Type:              util.Ptr(test.RandStr()),
		TypeID:            util.Ptr(test.RandInt()),
		LastUpdated:       util.Ptr(time.Now()),
		Fallbacks: &[]string{
			test.RandStr(),
			test.RandStr(),
		},
	}
}

func randParam() *tc.ParameterV5 {
	return &tc.ParameterV5{
		ConfigFile: test.RandStr(),
		Name:       test.RandStr(),
		Value:      test.RandStr(),
		Profiles:   []byte(`[]`),
	}
}

func randJob() *atscfg.InvalidationJob {
	now := time.Now()
	return &atscfg.InvalidationJob{
		AssetURL:         test.RandStr(),
		CreatedBy:        test.RandStr(),
		StartTime:        now,
		ID:               test.RandUint64(),
		DeliveryService:  test.RandStr(),
		TTLHours:         test.RandUint(),
		InvalidationType: tc.REFRESH,
	}
}

func randCDN() *tc.CDNV5 {
	return &tc.CDNV5{
		DNSSECEnabled: test.RandBool(),
		DomainName:    test.RandStr(),
		Name:          test.RandStr(),
	}
}

func randDSRs() *tc.DeliveryServiceRegexes {
	return &tc.DeliveryServiceRegexes{
		Regexes: []tc.DeliveryServiceRegex{
			*randDSR(),
			*randDSR(),
		},
		DSName: test.RandStr(),
	}
}

func randDSR() *tc.DeliveryServiceRegex {
	return &tc.DeliveryServiceRegex{
		Type:      string(tc.DSMatchTypeHostRegex),
		SetNumber: test.RandInt(),
		Pattern:   `\.*foo\.*`,
	}
}

func randCDNSSLKeys() *tc.CDNSSLKeys {
	return &tc.CDNSSLKeys{
		DeliveryService: test.RandStr(),
		Certificate: tc.CDNSSLKeysCertificate{
			Crt: test.RandStr(),
			Key: test.RandStr(),
		},
		Hostname: test.RandStr(),
	}
}

func MakeFakeTOData() *t3cutil.ConfigData {
	cg0 := *randCacheGroup()
	cg0.ParentName = nil
	cg0.ParentCachegroupID = nil

	cg1 := *randCacheGroup()
	cg1.ParentName = cg0.Name
	cg1.ParentCachegroupID = cg0.ID

	sv0 := randServer()
	sv1 := randServer()
	sv2 := randServer()

	sv0.CacheGroup = *cg0.Name
	sv1.CacheGroup = *cg0.Name
	sv2.CacheGroup = *cg1.Name

	ds0 := *randDS()
	ds1 := *randDS()

	dss := []atscfg.DeliveryServiceServer{
		atscfg.DeliveryServiceServer{
			Server:          sv0.ID,
			DeliveryService: *ds0.ID,
		},
		atscfg.DeliveryServiceServer{
			Server:          sv0.ID,
			DeliveryService: *ds1.ID,
		},
		atscfg.DeliveryServiceServer{
			Server:          sv1.ID,
			DeliveryService: *ds0.ID,
		},
	}

	dsr0 := randDSRs()
	dsr0.DSName = ds0.XMLID
	dsr0.Regexes[0].Pattern = `\.*foo\.*`
	// ds1.Pattern = `\.*bar\.*`

	dsr1 := randDSRs()
	dsr1.DSName = ds1.XMLID

	return &t3cutil.ConfigData{
		CacheGroups: []tc.CacheGroupNullableV5{
			cg0,
			cg1,
		},
		GlobalParams: []tc.ParameterV5{
			*randParam(),
			*randParam(),
			*randParam(),
		},
		ServerParams: []tc.ParameterV5{
			// configLocation := locationParams["remap.config"].Location
			tc.ParameterV5{
				ConfigFile: "remap.config",
				Name:       "location",
				Value:      "/etc/trafficserver",
				Profiles:   []byte(`[]`),
			},
			*randParam(),
			*randParam(),
			*randParam(),
		},
		RemapConfigParams: []tc.ParameterV5{
			*randParam(),
			*randParam(),
			*randParam(),
		},
		ParentConfigParams: []tc.ParameterV5{
			*randParam(),
			*randParam(),
			*randParam(),
		},
		DeliveryServices: []atscfg.DeliveryService{
			ds0,
			ds1,
		},
		Servers: []atscfg.Server{
			*sv1,
			*sv2,
		},
		DeliveryServiceServers: dss,
		Server:                 sv0,
		Jobs: []atscfg.InvalidationJob{
			*randJob(),
			*randJob(),
		},
		CDN: randCDN(),
		DeliveryServiceRegexes: []tc.DeliveryServiceRegexes{
			*dsr0,
			*dsr1,
		},
		URISigningKeys: map[tc.DeliveryServiceName][]byte{
			tc.DeliveryServiceName(test.RandStr()): []byte(test.RandStr()),
			tc.DeliveryServiceName(test.RandStr()): []byte(test.RandStr()),
		},
		URLSigKeys: map[tc.DeliveryServiceName]tc.URLSigKeys{
			tc.DeliveryServiceName(test.RandStr()): map[string]string{
				test.RandStr(): test.RandStr(),
				test.RandStr(): test.RandStr(),
			},
			tc.DeliveryServiceName(test.RandStr()): map[string]string{
				test.RandStr(): test.RandStr(),
				test.RandStr(): test.RandStr(),
			},
		},
		ServerCapabilities: map[int]map[atscfg.ServerCapability]struct{}{
			test.RandInt(): map[atscfg.ServerCapability]struct{}{
				atscfg.ServerCapability(test.RandStr()): struct{}{},
				atscfg.ServerCapability(test.RandStr()): struct{}{},
			},
			test.RandInt(): map[atscfg.ServerCapability]struct{}{
				atscfg.ServerCapability(test.RandStr()): struct{}{},
				atscfg.ServerCapability(test.RandStr()): struct{}{},
			},
		},
		DSRequiredCapabilities: map[int]map[atscfg.ServerCapability]struct{}{
			test.RandInt(): map[atscfg.ServerCapability]struct{}{
				atscfg.ServerCapability(test.RandStr()): struct{}{},
				atscfg.ServerCapability(test.RandStr()): struct{}{},
			},
			test.RandInt(): map[atscfg.ServerCapability]struct{}{
				atscfg.ServerCapability(test.RandStr()): struct{}{},
				atscfg.ServerCapability(test.RandStr()): struct{}{},
			},
		},
		SSLKeys: []tc.CDNSSLKeys{
			*randCDNSSLKeys(),
			*randCDNSSLKeys(),
		},
	}
}
