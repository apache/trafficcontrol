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
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	client "github.com/apache/trafficcontrol/traffic_ops/v3-client"
	"github.com/apache/trafficcontrol/traffic_ops_ort/atstccfg/config"
	"github.com/apache/trafficcontrol/traffic_ops_ort/atstccfg/toreq"
)

func TestWriteConfigs(t *testing.T) {
	buf := &bytes.Buffer{}
	configs := []config.ATSConfigFile{
		{
			Name:        "config0.txt",
			Path:        "/my/config0/location",
			Text:        "config0",
			ContentType: "text/plain",
		},
		{
			Name:        "config1.txt",
			Path:        "/my/config1/location",
			Text:        "config2,foo",
			ContentType: "text/csv",
		},
	}

	if err := WriteConfigs(configs, buf); err != nil {
		t.Fatalf("WriteConfigs error expected nil, actual: %v", err)
	}

	actual := buf.String()

	expected0 := "Content-Type: text/plain\r\nLine-Comment: \r\nPath: /my/config0/location/config0.txt\r\n\r\nconfig0\r\n"

	if !strings.Contains(actual, expected0) {
		t.Errorf("WriteConfigs expected '%v' actual '%v'", expected0, actual)
	}

	expected1 := "Content-Type: text/csv\r\nLine-Comment: \r\nPath: /my/config1/location/config1.txt\r\n\r\nconfig2,foo\r\n"
	if !strings.Contains(actual, expected1) {
		t.Errorf("WriteConfigs expected config1 '%v' actual '%v'", expected1, actual)
	}

	expectedPrefix := "MIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary="
	if !strings.HasPrefix(actual, expectedPrefix) {
		t.Errorf("WriteConfigs expected prefix '%v' actual '%v'", expectedPrefix, actual)
	}
}

func TestPreprocessConfigFile(t *testing.T) {
	// the TCP port replacement is fundamentally different for 80 vs non-80, so test both
	{
		server := &atscfg.Server{}
		server.TCPPort = util.IntPtr(8080)
		server.Interfaces = []tc.ServerInterfaceInfo{
			tc.ServerInterfaceInfo{
				IPAddresses: []tc.ServerIPAddress{
					tc.ServerIPAddress{
						Address:        "127.0.2.1",
						ServiceAddress: true,
					},
				},
			},
		}
		server.HostName = util.StrPtr("my-edge")
		server.DomainName = util.StrPtr("example.net")
		cfgFile := "abc__SERVER_TCP_PORT__def__CACHE_IPV4__ghi __RETURN__  \t __HOSTNAME__ jkl __FULL_HOSTNAME__ \n__SOMETHING__ __ELSE__\nmno\r\n"

		actual := PreprocessConfigFile(server, cfgFile)

		expected := "abc8080def127.0.2.1ghi\nmy-edge jkl my-edge.example.net \n__SOMETHING__ __ELSE__\nmno\r\n"

		if expected != actual {
			t.Errorf("PreprocessConfigFile expected '%v' actual '%v'", expected, actual)
		}
	}

	{
		server := &atscfg.Server{}
		server.TCPPort = util.IntPtr(80)
		server.Interfaces = []tc.ServerInterfaceInfo{
			tc.ServerInterfaceInfo{
				IPAddresses: []tc.ServerIPAddress{
					tc.ServerIPAddress{
						Address:        "127.0.2.1",
						ServiceAddress: true,
					},
				},
			},
		}
		server.HostName = util.StrPtr("my-edge")
		server.DomainName = util.StrPtr("example.net")

		cfgFile := "abc__SERVER_TCP_PORT__def__CACHE_IPV4__ghi __RETURN__  \t __HOSTNAME__ jkl __FULL_HOSTNAME__ \n__SOMETHING__ __ELSE__\nmno:__SERVER_TCP_PORT__\r\n"

		actual := PreprocessConfigFile(server, cfgFile)

		expected := "abc__SERVER_TCP_PORT__def127.0.2.1ghi\nmy-edge jkl my-edge.example.net \n__SOMETHING__ __ELSE__\nmno\r\n"

		if expected != actual {
			t.Errorf("PreprocessConfigFile expected '%v' actual '%v'", expected, actual)
		}
	}
}

// TestGetAllConfigsWriteConfigsDeterministic tests that WriteConfigs(GetAllConfigs()) is Deterministic.
// That is, that for the same input, it always produces the same output.
//
// Because Go map iteration is defined to be random, running it multiple times even on the exact same input could be different, if there's a determinism bug.
// But beyond that, we re-order slices whose order isn't semantically significant (e.g. params) and run it again.
//
func TestGetAllConfigsWriteConfigsDeterministic(t *testing.T) {
	// TODO expand fake data. Currently, it's only making a remap.config.
	toData := MakeFakeTOData()
	revalOnly := false
	cfgPath := "/etc/trafficserver/"

	cfg := config.TCCfg{}
	cfg.Dir = cfgPath
	cfg.RevalOnly = revalOnly
	cfg.TOClient = &toreq.TOClient{}
	cfg.TOClient.C = &client.Session{}
	cfg.TOClient.C.URL = ""

	configs, err := GetAllConfigs(toData, "", nil, cfg)
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
		configs2, err := GetAllConfigs(toData, "", nil, cfg)
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

func randBool() *bool {
	b := rand.Int()%2 == 0
	return &b
}

func randStr() *string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_"
	num := 100
	s := ""
	for i := 0; i < num; i++ {
		s += string(chars[rand.Intn(len(chars))])
	}
	return &s
}

func randInt() *int {
	i := rand.Int()
	return &i
}

func randInt64() *int64 {
	i := int64(rand.Int63())
	return &i
}

func randFloat64() *float64 {
	f := rand.Float64()
	return &f
}

func randDSS() tc.DeliveryServiceServer {
	return tc.DeliveryServiceServer{
		Server:          randInt(),
		DeliveryService: randInt(),
	}
}

func randDS() *atscfg.DeliveryService {
	deepCachingTypeNever := tc.DeepCachingTypeNever
	dsTypeHTTP := tc.DSTypeHTTP
	protocol := tc.DSProtocolHTTP
	ds := atscfg.DeliveryService{}
	ds.EcsEnabled = *randBool()
	ds.RangeSliceBlockSize = randInt()
	ds.ConsistentHashRegex = randStr()
	ds.ConsistentHashQueryParams = []string{
		*randStr(),
		*randStr(),
	}
	ds.MaxOriginConnections = randInt()
	ds.DeepCachingType = &deepCachingTypeNever
	ds.FQPacingRate = randInt()
	ds.SigningAlgorithm = randStr()
	ds.Tenant = randStr()
	ds.TRResponseHeaders = randStr()
	ds.TRRequestHeaders = randStr()
	ds.Active = randBool()
	ds.AnonymousBlockingEnabled = randBool()
	ds.CCRDNSTTL = randInt()
	ds.CDNID = randInt()
	ds.CDNName = randStr()
	ds.CheckPath = randStr()
	ds.DisplayName = randStr()
	ds.DNSBypassCNAME = randStr()
	ds.DNSBypassIP = randStr()
	ds.DNSBypassIP6 = randStr()
	ds.DNSBypassTTL = randInt()
	ds.DSCP = randInt()
	ds.EdgeHeaderRewrite = randStr()
	ds.GeoLimit = randInt()
	ds.GeoLimitCountries = randStr()
	ds.GeoLimitRedirectURL = randStr()
	ds.GeoProvider = randInt()
	ds.GlobalMaxMBPS = randInt()
	ds.GlobalMaxTPS = randInt()
	ds.HTTPBypassFQDN = randStr()
	ds.ID = randInt()
	ds.InfoURL = randStr()
	ds.InitialDispersion = randInt()
	ds.IPV6RoutingEnabled = randBool()
	ds.LastUpdated = &tc.TimeNoMod{Time: time.Now()}
	ds.LogsEnabled = randBool()
	ds.LongDesc = randStr()
	ds.LongDesc1 = randStr()
	ds.LongDesc2 = randStr()
	ds.MatchList = &[]tc.DeliveryServiceMatch{
		{
			Type:      tc.DSMatchTypeHostRegex,
			SetNumber: 0,
			Pattern:   `\.*foo\.*`,
		},
	}
	ds.MaxDNSAnswers = randInt()
	ds.MidHeaderRewrite = randStr()
	ds.MissLat = randFloat64()
	ds.MissLong = randFloat64()
	ds.MultiSiteOrigin = randBool()
	ds.OriginShield = randStr()
	ds.OrgServerFQDN = util.StrPtr("http://" + *(randStr()))
	ds.ProfileDesc = randStr()
	ds.ProfileID = randInt()
	ds.ProfileName = randStr()
	ds.Protocol = &protocol
	ds.QStringIgnore = randInt()
	ds.RangeRequestHandling = randInt()
	ds.RegexRemap = randStr()
	ds.RegionalGeoBlocking = randBool()
	ds.RemapText = randStr()
	ds.RoutingName = randStr()
	ds.Signed = *randBool()
	ds.SSLKeyVersion = randInt()
	ds.TenantID = randInt()
	ds.Type = &dsTypeHTTP
	ds.TypeID = randInt()
	ds.XMLID = randStr()
	ds.ExampleURLs = []string{
		*randStr(),
		*randStr(),
	}
	return &ds
}

func randServer() *atscfg.Server {
	sv := &atscfg.Server{}
	sv.Cachegroup = randStr()
	sv.CachegroupID = randInt()
	sv.CDNID = randInt()
	sv.CDNName = randStr()
	sv.DomainName = randStr()
	sv.FQDN = randStr()
	sv.FqdnTime = time.Now()
	sv.GUID = randStr()
	sv.HostName = randStr()
	sv.HTTPSPort = randInt()
	sv.ID = randInt()
	sv.ILOIPAddress = randStr()
	sv.ILOIPGateway = randStr()
	sv.ILOIPNetmask = randStr()
	sv.ILOPassword = randStr()
	sv.ILOUsername = randStr()

	sv.Interfaces = []tc.ServerInterfaceInfo{
		tc.ServerInterfaceInfo{
			Name: *randStr(),
			IPAddresses: []tc.ServerIPAddress{
				tc.ServerIPAddress{
					Address:        *randStr(),
					Gateway:        randStr(),
					ServiceAddress: true,
				},
				tc.ServerIPAddress{
					Address:        *randStr(),
					Gateway:        randStr(),
					ServiceAddress: true,
				},
			},
		},
	}

	sv.LastUpdated = &tc.TimeNoMod{Time: time.Now()}
	sv.MgmtIPAddress = randStr()
	sv.MgmtIPGateway = randStr()
	sv.MgmtIPNetmask = randStr()
	sv.OfflineReason = randStr()
	sv.PhysLocation = randStr()
	sv.PhysLocationID = randInt()
	sv.Profile = randStr()
	sv.ProfileDesc = randStr()
	sv.ProfileID = randInt()
	sv.Rack = randStr()
	sv.RevalPending = randBool()
	sv.Status = randStr()
	sv.StatusID = randInt()
	sv.TCPPort = randInt()
	sv.Type = *randStr()
	sv.TypeID = randInt()
	sv.UpdPending = randBool()
	sv.XMPPID = randStr()
	sv.XMPPPasswd = randStr()
	return sv
}

func randCacheGroup() *tc.CacheGroupNullable {
	return &tc.CacheGroupNullable{
		ID:        randInt(),
		Name:      randStr(),
		ShortName: randStr(),
		Latitude:  randFloat64(),
		Longitude: randFloat64(),
		// ParentName:                  randStr(),
		// ParentCachegroupID:          randInt(),
		// SecondaryParentName:         randStr(),
		// SecondaryParentCachegroupID: randInt(),
		FallbackToClosest: randBool(),
		Type:              randStr(),
		TypeID:            randInt(),
		LastUpdated:       &tc.TimeNoMod{Time: time.Now()},
		Fallbacks: &[]string{
			*randStr(),
			*randStr(),
		},
	}
}

func randParam() *tc.Parameter {
	return &tc.Parameter{
		ConfigFile: *randStr(),
		Name:       *randStr(),
		Value:      *randStr(),
		Profiles:   []byte(`[]`),
	}
}

func randJob() *tc.Job {
	return &tc.Job{
		Parameters:      *randStr(),
		Keyword:         *randStr(),
		AssetURL:        *randStr(),
		CreatedBy:       *randStr(),
		StartTime:       *randStr(),
		ID:              *randInt64(),
		DeliveryService: *randStr(),
	}
}

func randCDN() *tc.CDN {
	return &tc.CDN{
		DNSSECEnabled: *randBool(),
		DomainName:    *randStr(),
		Name:          *randStr(),
	}
}

func randDSRs() *tc.DeliveryServiceRegexes {
	return &tc.DeliveryServiceRegexes{
		Regexes: []tc.DeliveryServiceRegex{
			*randDSR(),
			*randDSR(),
		},
		DSName: *randStr(),
	}
}

func randDSR() *tc.DeliveryServiceRegex {
	return &tc.DeliveryServiceRegex{
		Type:      string(tc.DSMatchTypeHostRegex),
		SetNumber: *randInt(),
		Pattern:   `\.*foo\.*`,
	}
}

func randCDNSSLKeys() *tc.CDNSSLKeys {
	return &tc.CDNSSLKeys{
		DeliveryService: *randStr(),
		Certificate: tc.CDNSSLKeysCertificate{
			Crt: *randStr(),
			Key: *randStr(),
		},
		Hostname: *randStr(),
	}
}

func MakeFakeTOData() *config.TOData {
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

	dss := []tc.DeliveryServiceServer{
		tc.DeliveryServiceServer{
			Server:          sv0.ID,
			DeliveryService: ds0.ID,
		},
		tc.DeliveryServiceServer{
			Server:          sv0.ID,
			DeliveryService: ds1.ID,
		},
		tc.DeliveryServiceServer{
			Server:          sv1.ID,
			DeliveryService: ds0.ID,
		},
	}

	dsr0 := randDSRs()
	dsr0.DSName = *ds0.XMLID
	dsr0.Regexes[0].Pattern = `\.*foo\.*`
	// ds1.Pattern = `\.*bar\.*`

	dsr1 := randDSRs()
	dsr1.DSName = *ds1.XMLID

	return &config.TOData{
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
		CacheKeyParams: []tc.Parameter{
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
		Jobs: []tc.Job{
			*randJob(),
			*randJob(),
		},
		CDN: randCDN(),
		DeliveryServiceRegexes: []tc.DeliveryServiceRegexes{
			*dsr0,
			*dsr1,
		},
		URISigningKeys: map[tc.DeliveryServiceName][]byte{
			tc.DeliveryServiceName(*randStr()): []byte(*randStr()),
			tc.DeliveryServiceName(*randStr()): []byte(*randStr()),
		},
		URLSigKeys: map[tc.DeliveryServiceName]tc.URLSigKeys{
			tc.DeliveryServiceName(*randStr()): map[string]string{
				*randStr(): *randStr(),
				*randStr(): *randStr(),
			},
			tc.DeliveryServiceName(*randStr()): map[string]string{
				*randStr(): *randStr(),
				*randStr(): *randStr(),
			},
		},
		ServerCapabilities: map[int]map[atscfg.ServerCapability]struct{}{
			*randInt(): map[atscfg.ServerCapability]struct{}{
				atscfg.ServerCapability(*randStr()): struct{}{},
				atscfg.ServerCapability(*randStr()): struct{}{},
			},
			*randInt(): map[atscfg.ServerCapability]struct{}{
				atscfg.ServerCapability(*randStr()): struct{}{},
				atscfg.ServerCapability(*randStr()): struct{}{},
			},
		},
		DSRequiredCapabilities: map[int]map[atscfg.ServerCapability]struct{}{
			*randInt(): map[atscfg.ServerCapability]struct{}{
				atscfg.ServerCapability(*randStr()): struct{}{},
				atscfg.ServerCapability(*randStr()): struct{}{},
			},
			*randInt(): map[atscfg.ServerCapability]struct{}{
				atscfg.ServerCapability(*randStr()): struct{}{},
				atscfg.ServerCapability(*randStr()): struct{}{},
			},
		},
		SSLKeys: []tc.CDNSSLKeys{
			*randCDNSSLKeys(),
			*randCDNSSLKeys(),
		},
	}
}
