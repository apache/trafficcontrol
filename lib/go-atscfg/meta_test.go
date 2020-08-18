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

func TestMakeMetaConfig(t *testing.T) {
	server := &tc.ServerNullable{}
	server.CachegroupID = util.IntPtr(42)
	server.Cachegroup = util.StrPtr("cg0")
	server.CDNName = util.StrPtr("mycdn")
	server.CDNID = util.IntPtr(43)
	server.DomainName = util.StrPtr("myserverdomain.invalid")
	server.HostName = util.StrPtr("myserver")
	server.HTTPSPort = util.IntPtr(443)
	server.ID = util.IntPtr(44)
	ip := "192.168.2.9"
	setIP(server, ip)
	// server.ParentCacheGroupID=            45
	// server.ParentCacheGroupType=          "MID_LOC"
	server.ProfileID = util.IntPtr(46)
	server.Profile = util.StrPtr("myserverprofile")
	server.TCPPort = util.IntPtr(80)
	// server.SecondaryParentCacheGroupID=   47
	// server.SecondaryParentCacheGroupType= "MID_LOC"
	server.Type = "EDGE"

	tmURL := "https://myto.invalid"
	tmReverseProxyURL := "https://myrp.myto.invalid"
	locationParams := map[string]ConfigProfileParams{
		"regex_revalidate.config": ConfigProfileParams{
			FileNameOnDisk: "regex_revalidate.config",
			Location:       "/my/location/",
			URL:            "http://myurl/remap.config", // cdn-scoped endpoint
		},
		"cache.config": ConfigProfileParams{
			FileNameOnDisk: "cache.config", // cache.config on mids is server-scoped
			Location:       "/my/location/",
		},
		"ip_allow.config": ConfigProfileParams{
			FileNameOnDisk: "ip_allow.config",
			Location:       "/my/location/",
		},
		"volume.config": ConfigProfileParams{
			FileNameOnDisk: "volume.config",
			Location:       "/my/location/",
		},
		"ssl_multicert.config": ConfigProfileParams{
			FileNameOnDisk: "ssl_multicert.config",
			Location:       "/my/location/",
		},
		"uri_signing_mydsname.config": ConfigProfileParams{
			FileNameOnDisk: "uri_signing_mydsname.config",
			Location:       "/my/location/",
		},
		"uri_signing_nonexistentds.config": ConfigProfileParams{
			FileNameOnDisk: "uri_signing_nonexistentds.config",
			Location:       "/my/location/",
		},
		"regex_remap_nonexistentds.config": ConfigProfileParams{
			FileNameOnDisk: "regex_remap_nonexistentds.config",
			Location:       "/my/location/",
		},
		"url_sig_nonexistentds.config": ConfigProfileParams{
			FileNameOnDisk: "url_sig_nonexistentds.config",
			Location:       "/my/location/",
		},
		"hdr_rw_nonexistentds.config": ConfigProfileParams{
			FileNameOnDisk: "hdr_rw_nonexistentds.config",
			Location:       "/my/location/",
		},
		"hdr_rw_mid_nonexistentds.config": ConfigProfileParams{
			FileNameOnDisk: "hdr_rw_mid_nonexistentds.config",
			Location:       "/my/location/",
		},
		"unknown.config": ConfigProfileParams{
			FileNameOnDisk: "unknown.config",
			Location:       "/my/location/",
		},
		"custom.config": ConfigProfileParams{
			FileNameOnDisk: "custom.config",
			Location:       "/my/location/",
		},
		"external.config": ConfigProfileParams{
			FileNameOnDisk: "external.config",
			Location:       "/my/location/",
			URL:            "http://myurl/remap.config",
		},
	}
	uriSignedDSes := []tc.DeliveryServiceName{"mydsname"}
	dses := map[tc.DeliveryServiceName]tc.DeliveryServiceNullableV30{"mydsname": {}}

	scopeParams := map[string]string{"custom.config": string(tc.ATSConfigMetaDataConfigFileScopeProfiles)}

	cgs := []tc.CacheGroupNullable{}
	topologies := []tc.Topology{}

	cfgPath := "/etc/foo/trafficserver"

	cfg, err := MakeMetaObj(server, tmURL, tmReverseProxyURL, locationParams, uriSignedDSes, scopeParams, dses, cgs, topologies, cfgPath)
	if err != nil {
		t.Fatalf("MakeMetaObj: " + err.Error())
	}

	if cfg.Info.ProfileID != int(*server.ProfileID) {
		t.Errorf("expected Info.ProfileID %v actual %v", server.ProfileID, cfg.Info.ProfileID)
	}

	if cfg.Info.TOReverseProxyURL != tmReverseProxyURL {
		t.Errorf("expected Info.TOReverseProxyURL %v actual %v", tmReverseProxyURL, cfg.Info.TOReverseProxyURL)
	}

	if cfg.Info.TOURL != tmURL {
		t.Errorf("expected Info.TOURL %v actual %v", tmURL, cfg.Info.TOURL)
	}

	if *server.TCPPort != cfg.Info.ServerPort {
		t.Errorf("expected Info.ServerPort %v actual %v", server.TCPPort, cfg.Info.ServerPort)
	}

	if *server.HostName != cfg.Info.ServerName {
		t.Errorf("expected Info.ServerName %v actual %v", *server.HostName, cfg.Info.ServerName)
	}

	if cfg.Info.CDNID != *server.CDNID {
		t.Errorf("expected Info.CDNID %v actual %v", *server.CDNID, cfg.Info.CDNID)
	}

	if cfg.Info.CDNName != *server.CDNName {
		t.Errorf("expected Info.CDNName %v actual %v", server.CDNName, cfg.Info.CDNName)
	}
	if cfg.Info.ServerID != *server.ID {
		t.Errorf("expected Info.ServerID %v actual %v", *server.ID, cfg.Info.ServerID)
	}
	if cfg.Info.ProfileName != *server.Profile {
		t.Errorf("expected Info.ProfileName %v actual %v", *server.Profile, cfg.Info.ProfileName)
	}

	expectedConfigs := map[string]func(cf tc.ATSConfigMetaDataConfigFile){
		"cache.config": func(cf tc.ATSConfigMetaDataConfigFile) {
			if expected := "/my/location/"; cf.Location != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Location)
			}
			if expected := string(tc.ATSConfigMetaDataConfigFileScopeProfiles); cf.Scope != expected {
				t.Errorf("expected scope '%v', actual '%v'", expected, cf.Scope)
			}
		},
		"ip_allow.config": func(cf tc.ATSConfigMetaDataConfigFile) {
			if expected := "/my/location/"; cf.Location != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Location)
			}
			if expected := string(tc.ATSConfigMetaDataConfigFileScopeServers); cf.Scope != expected {
				t.Errorf("expected scope '%v', actual '%v'", expected, cf.Scope)
			}
		},
		"volume.config": func(cf tc.ATSConfigMetaDataConfigFile) {
			if expected := "/my/location/"; cf.Location != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Location)
			}
			if expected := string(tc.ATSConfigMetaDataConfigFileScopeProfiles); cf.Scope != expected {
				t.Errorf("expected scope '%v', actual '%v'", expected, cf.Scope)
			}
		},
		"ssl_multicert.config": func(cf tc.ATSConfigMetaDataConfigFile) {
			if expected := "/my/location/"; cf.Location != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Location)
			}
			if expected := string(tc.ATSConfigMetaDataConfigFileScopeCDNs); cf.Scope != expected {
				t.Errorf("expected scope '%v', actual '%v'", expected, cf.Scope)
			}
		},
		"uri_signing_mydsname.config": func(cf tc.ATSConfigMetaDataConfigFile) {
			if expected := "/my/location/"; cf.Location != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Location)
			}
			if expected := string(tc.ATSConfigMetaDataConfigFileScopeProfiles); cf.Scope != expected {
				t.Errorf("expected scope '%v', actual '%v'", expected, cf.Scope)
			}
		},
		"unknown.config": func(cf tc.ATSConfigMetaDataConfigFile) {
			if expected := "/my/location/"; cf.Location != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Location)
			}
			if expected := string(tc.ATSConfigMetaDataConfigFileScopeServers); cf.Scope != expected {
				t.Errorf("expected scope '%v', actual '%v'", expected, cf.Scope)
			}
		},
		"custom.config": func(cf tc.ATSConfigMetaDataConfigFile) {
			if expected := "/my/location/"; cf.Location != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Location)
			}
			if expected := string(tc.ATSConfigMetaDataConfigFileScopeProfiles); cf.Scope != expected {
				t.Errorf("expected scope '%v', actual '%v'", expected, cf.Scope)
			}
		},
		"remap.config": func(cf tc.ATSConfigMetaDataConfigFile) {
			if expected := cfgPath; cf.Location != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Location)
			}
			if expected := string(tc.ATSConfigMetaDataConfigFileScopeServers); cf.Scope != expected {
				t.Errorf("expected scope '%v', actual '%v'", expected, cf.Scope)
			}
		},
		"regex_revalidate.config": func(cf tc.ATSConfigMetaDataConfigFile) {
			if expected := "/my/location/"; cf.Location != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Location)
			}
			if expected := string(tc.ATSConfigMetaDataConfigFileScopeCDNs); cf.Scope != expected {
				t.Errorf("expected scope '%v', actual '%v'", expected, cf.Scope)
			}
		},
		"external.config": func(cf tc.ATSConfigMetaDataConfigFile) {
			if expected := "/my/location/"; cf.Location != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Location)
			}
			if expected := string(tc.ATSConfigMetaDataConfigFileScopeCDNs); cf.Scope != expected {
				t.Errorf("expected scope '%v', actual '%v'", expected, cf.Scope)
			}
		},
		"hosting.config": func(cf tc.ATSConfigMetaDataConfigFile) {
			if expected := cfgPath; cf.Location != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Location)
			}
			if expected := string(tc.ATSConfigMetaDataConfigFileScopeServers); cf.Scope != expected {
				t.Errorf("expected scope for %v is '%v', actual '%v'", cf.FileNameOnDisk, expected, cf.Scope)
			}
		},
		"parent.config": func(cf tc.ATSConfigMetaDataConfigFile) {
			if expected := cfgPath; cf.Location != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Location)
			}
			if expected := string(tc.ATSConfigMetaDataConfigFileScopeServers); cf.Scope != expected {
				t.Errorf("expected scope for %v is '%v', actual '%v'", cf.FileNameOnDisk, expected, cf.Scope)
			}
		},
		"plugin.config": func(cf tc.ATSConfigMetaDataConfigFile) {
			if expected := cfgPath; cf.Location != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Location)
			}
			if expected := string(tc.ATSConfigMetaDataConfigFileScopeProfiles); cf.Scope != expected {
				t.Errorf("expected scope for %v is '%v', actual '%v'", cf.FileNameOnDisk, expected, cf.Scope)
			}
		},
		"records.config": func(cf tc.ATSConfigMetaDataConfigFile) {
			if expected := cfgPath; cf.Location != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Location)
			}
			if expected := string(tc.ATSConfigMetaDataConfigFileScopeProfiles); cf.Scope != expected {
				t.Errorf("expected scope for %v is '%v', actual '%v'", cf.FileNameOnDisk, expected, cf.Scope)
			}
		},
		"storage.config": func(cf tc.ATSConfigMetaDataConfigFile) {
			if expected := cfgPath; cf.Location != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Location)
			}
			if expected := string(tc.ATSConfigMetaDataConfigFileScopeProfiles); cf.Scope != expected {
				t.Errorf("expected scope for %v is '%v', actual '%v'", cf.FileNameOnDisk, expected, cf.Scope)
			}
		},
	}

	for _, cfgFile := range cfg.ConfigFiles {
		if testF, ok := expectedConfigs[cfgFile.FileNameOnDisk]; !ok {
			t.Errorf("unexpected config '" + cfgFile.FileNameOnDisk + "'")
		} else {
			testF(cfgFile)
			delete(expectedConfigs, cfgFile.FileNameOnDisk)
		}
	}

	server.Type = "MID"
	cfg, err = MakeMetaObj(server, tmURL, tmReverseProxyURL, locationParams, uriSignedDSes, scopeParams, dses, cgs, topologies, cfgPath)
	if err != nil {
		t.Fatalf("MakeMetaObj: " + err.Error())
	}
	for _, cfgFile := range cfg.ConfigFiles {
		if cfgFile.FileNameOnDisk != "cache.config" {
			continue
		}
		if expected := string(tc.ATSConfigMetaDataConfigFileScopeServers); cfgFile.Scope != expected {
			t.Errorf("expected cache.config on a Mid to be scope '%v', actual '%v'", expected, cfgFile.Scope)
		}
		break
	}
	for _, fi := range cfg.ConfigFiles {
		if strings.Contains(fi.FileNameOnDisk, "nonexistentds") {
			t.Errorf("expected location parameters for nonexistent delivery services to not be added to config, actual '%v'", fi.FileNameOnDisk)
		}
		if fi.FileNameOnDisk == "external.config" {
			if fi.APIURI != "" {
				t.Errorf("expected: apiUri field to be omitted for external.config, actual: present")
			}
			if fi.URL != "" {
				t.Errorf("expected: url field to be present for external.config, actual: omitted")
			}
		}
	}
}
