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

	"github.com/apache/trafficcontrol/cache-config/t3c-generate/config"
	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/test"
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

func randDSS() tc.DeliveryServiceServer {
	return tc.DeliveryServiceServer{
		Server:          util.IntPtr(test.RandInt()),
		DeliveryService: util.IntPtr(test.RandInt()),
	}
}

func randDS() *atscfg.DeliveryService {
	deepCachingTypeNever := tc.DeepCachingTypeNever
	dsTypeHTTP := tc.DSTypeHTTP
	protocol := tc.DSProtocolHTTP
	ds := atscfg.DeliveryService{}
	ds.EcsEnabled = test.RandBool()
	ds.RangeSliceBlockSize = util.IntPtr(test.RandInt())
	ds.ConsistentHashRegex = util.StrPtr(test.RandStr())
	ds.ConsistentHashQueryParams = []string{
		test.RandStr(),
		test.RandStr(),
	}
	ds.MaxOriginConnections = util.IntPtr(test.RandInt())
	ds.DeepCachingType = &deepCachingTypeNever
	ds.FQPacingRate = util.IntPtr(test.RandInt())
	ds.SigningAlgorithm = util.StrPtr(test.RandStr())
	ds.Tenant = util.StrPtr(test.RandStr())
	ds.TRResponseHeaders = util.StrPtr(test.RandStr())
	ds.TRRequestHeaders = util.StrPtr(test.RandStr())
	ds.Active = util.BoolPtr(test.RandBool())
	ds.AnonymousBlockingEnabled = util.BoolPtr(test.RandBool())
	ds.TRDNSTTL = util.IntPtr(test.RandInt())
	ds.CDNID = util.IntPtr(test.RandInt())
	ds.CDNName = util.StrPtr(test.RandStr())
	ds.CheckPath = util.StrPtr(test.RandStr())
	ds.DisplayName = util.StrPtr(test.RandStr())
	ds.DNSBypassCNAME = util.StrPtr(test.RandStr())
	ds.DNSBypassIP = util.StrPtr(test.RandStr())
	ds.DNSBypassIP6 = util.StrPtr(test.RandStr())
	ds.DNSBypassTTL = util.IntPtr(test.RandInt())
	ds.DSCP = util.IntPtr(test.RandInt())
	ds.EdgeHeaderRewrite = util.StrPtr(test.RandStr())
	ds.GeoLimit = util.IntPtr(test.RandInt())
	ds.GeoLimitCountries = nil
	ds.GeoLimitRedirectURL = util.StrPtr(test.RandStr())
	ds.GeoProvider = util.IntPtr(test.RandInt())
	ds.GlobalMaxMBPS = util.IntPtr(test.RandInt())
	ds.GlobalMaxTPS = util.IntPtr(test.RandInt())
	ds.HTTPBypassFQDN = util.StrPtr(test.RandStr())
	ds.ID = util.IntPtr(test.RandInt())
	ds.InfoURL = util.StrPtr(test.RandStr())
	ds.InitialDispersion = util.IntPtr(test.RandInt())
	ds.IPV6RoutingEnabled = util.BoolPtr(test.RandBool())
	ds.LastUpdated = &tc.TimeNoMod{Time: time.Now()}
	ds.LogsEnabled = util.BoolPtr(test.RandBool())
	ds.LongDesc = util.StrPtr(test.RandStr())
	ds.LongDesc1 = util.StrPtr(test.RandStr())
	ds.LongDesc2 = util.StrPtr(test.RandStr())
	ds.MatchList = &[]tc.DeliveryServiceMatch{
		{
			Type:      tc.DSMatchTypeHostRegex,
			SetNumber: 0,
			Pattern:   `\.*foo\.*`,
		},
	}
	ds.MaxDNSAnswers = util.IntPtr(test.RandInt())
	ds.MidHeaderRewrite = util.StrPtr(test.RandStr())
	ds.MissLat = util.FloatPtr(test.RandFloat64())
	ds.MissLong = util.FloatPtr(test.RandFloat64())
	ds.MultiSiteOrigin = util.BoolPtr(test.RandBool())
	ds.OriginShield = util.StrPtr(test.RandStr())
	ds.OrgServerFQDN = util.StrPtr("http://" + test.RandStr())
	ds.ProfileDesc = util.StrPtr(test.RandStr())
	ds.ProfileID = util.IntPtr(test.RandInt())
	ds.ProfileName = util.StrPtr(test.RandStr())
	ds.Protocol = &protocol
	ds.QStringIgnore = util.IntPtr(test.RandInt())
	ds.RangeRequestHandling = util.IntPtr(test.RandInt())
	ds.RegexRemap = util.StrPtr(test.RandStr())
	ds.RegionalGeoBlocking = util.BoolPtr(test.RandBool())
	ds.RemapText = util.StrPtr(test.RandStr())
	ds.RoutingName = util.StrPtr(test.RandStr())
	ds.Signed = *util.BoolPtr(test.RandBool())
	ds.SSLKeyVersion = util.IntPtr(test.RandInt())
	ds.TenantID = util.IntPtr(test.RandInt())
	ds.Type = &dsTypeHTTP
	ds.TypeID = util.IntPtr(test.RandInt())
	ds.XMLID = util.StrPtr(test.RandStr())
	ds.ExampleURLs = []string{
		test.RandStr(),
		test.RandStr(),
	}
	return &ds
}

func randServer() *atscfg.Server {
	sv := &atscfg.Server{}
	sv.Cachegroup = util.StrPtr(test.RandStr())
	sv.CachegroupID = util.IntPtr(test.RandInt())
	sv.CDNID = util.IntPtr(test.RandInt())
	sv.CDNName = util.StrPtr(test.RandStr())
	sv.DomainName = util.StrPtr(test.RandStr())
	sv.FQDN = util.StrPtr(test.RandStr())
	sv.FqdnTime = time.Now()
	sv.GUID = util.StrPtr(test.RandStr())
	sv.HostName = util.StrPtr(test.RandStr())
	sv.HTTPSPort = util.IntPtr(test.RandInt())
	sv.ID = util.IntPtr(test.RandInt())
	sv.ILOIPAddress = util.StrPtr(test.RandStr())
	sv.ILOIPGateway = util.StrPtr(test.RandStr())
	sv.ILOIPNetmask = util.StrPtr(test.RandStr())
	sv.ILOPassword = util.StrPtr(test.RandStr())
	sv.ILOUsername = util.StrPtr(test.RandStr())

	sv.Interfaces = []tc.ServerInterfaceInfoV40{}
	{
		ssi := tc.ServerInterfaceInfoV40{}
		ssi.Name = test.RandStr()
		ssi.IPAddresses = []tc.ServerIPAddress{
			tc.ServerIPAddress{
				Address:        test.RandStr(),
				Gateway:        util.StrPtr(test.RandStr()),
				ServiceAddress: true,
			},
			tc.ServerIPAddress{
				Address:        test.RandStr(),
				Gateway:        util.StrPtr(test.RandStr()),
				ServiceAddress: true,
			},
		}
		sv.Interfaces = append(sv.Interfaces, ssi)
	}

	sv.LastUpdated = &tc.TimeNoMod{Time: time.Now()}
	sv.MgmtIPAddress = util.StrPtr(test.RandStr())
	sv.MgmtIPGateway = util.StrPtr(test.RandStr())
	sv.MgmtIPNetmask = util.StrPtr(test.RandStr())
	sv.OfflineReason = util.StrPtr(test.RandStr())
	sv.PhysLocation = util.StrPtr(test.RandStr())
	sv.PhysLocationID = util.IntPtr(test.RandInt())
	sv.ProfileNames = []string{test.RandStr()}
	sv.Rack = util.StrPtr(test.RandStr())
	sv.RevalPending = util.BoolPtr(test.RandBool())
	sv.Status = util.StrPtr(test.RandStr())
	sv.StatusID = util.IntPtr(test.RandInt())
	sv.TCPPort = util.IntPtr(test.RandInt())
	sv.Type = test.RandStr()
	sv.TypeID = util.IntPtr(test.RandInt())
	sv.UpdPending = util.BoolPtr(test.RandBool())
	sv.XMPPID = util.StrPtr(test.RandStr())
	sv.XMPPPasswd = util.StrPtr(test.RandStr())
	return sv
}

func randCacheGroup() *tc.CacheGroupNullable {
	return &tc.CacheGroupNullable{
		ID:        util.IntPtr(test.RandInt()),
		Name:      util.StrPtr(test.RandStr()),
		ShortName: util.StrPtr(test.RandStr()),
		Latitude:  util.FloatPtr(test.RandFloat64()),
		Longitude: util.FloatPtr(test.RandFloat64()),
		// ParentName:                  util.StrPtr(test.RandStr()),
		// ParentCachegroupID:          util.IntPtr(test.RandInt()),
		// SecondaryParentName:         util.StrPtr(test.RandStr()),
		// SecondaryParentCachegroupID: util.IntPtr(test.RandInt()),
		FallbackToClosest: util.BoolPtr(test.RandBool()),
		Type:              util.StrPtr(test.RandStr()),
		TypeID:            util.IntPtr(test.RandInt()),
		LastUpdated:       &tc.TimeNoMod{Time: time.Now()},
		Fallbacks: &[]string{
			test.RandStr(),
			test.RandStr(),
		},
	}
}

func randParam() *tc.Parameter {
	return &tc.Parameter{
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

func randCDN() *tc.CDN {
	return &tc.CDN{
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

	sv0.Cachegroup = cg0.Name
	sv1.Cachegroup = cg0.Name
	sv2.Cachegroup = cg1.Name

	ds0 := *randDS()
	ds1 := *randDS()

	dss := []atscfg.DeliveryServiceServer{
		atscfg.DeliveryServiceServer{
			Server:          *sv0.ID,
			DeliveryService: *ds0.ID,
		},
		atscfg.DeliveryServiceServer{
			Server:          *sv0.ID,
			DeliveryService: *ds1.ID,
		},
		atscfg.DeliveryServiceServer{
			Server:          *sv1.ID,
			DeliveryService: *ds0.ID,
		},
	}

	dsr0 := randDSRs()
	dsr0.DSName = *ds0.XMLID
	dsr0.Regexes[0].Pattern = `\.*foo\.*`
	// ds1.Pattern = `\.*bar\.*`

	dsr1 := randDSRs()
	dsr1.DSName = *ds1.XMLID

	return &t3cutil.ConfigData{
		CacheGroups: []tc.CacheGroupNullable{
			cg0,
			cg1,
		},
		GlobalParams: []tc.Parameter{
			*randParam(),
			*randParam(),
			*randParam(),
		},
		ServerParams: []tc.Parameter{
			// configLocation := locationParams["remap.config"].Location
			tc.Parameter{
				ConfigFile: "remap.config",
				Name:       "location",
				Value:      "/etc/trafficserver",
				Profiles:   []byte(`[]`),
			},
			*randParam(),
			*randParam(),
			*randParam(),
		},
		RemapConfigParams: []tc.Parameter{
			*randParam(),
			*randParam(),
			*randParam(),
		},
		ParentConfigParams: []tc.Parameter{
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
