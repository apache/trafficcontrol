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
	"encoding/json"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestMakeMetaConfig(t *testing.T) {
	serverHostName := tc.CacheName("myServer")
	server := &ServerInfo{
		CacheGroupID:                  42,
		CDN:                           "mycdn",
		CDNID:                         43,
		DomainName:                    "myserverdomain.invalid",
		HostName:                      "myserver",
		HTTPSPort:                     443,
		ID:                            44,
		IP:                            "192.168.2.9",
		ParentCacheGroupID:            45,
		ParentCacheGroupType:          "MID_LOC",
		ProfileID:                     46,
		ProfileName:                   "myserverprofile",
		Port:                          80,
		SecondaryParentCacheGroupID:   47,
		SecondaryParentCacheGroupType: "MID_LOC",
		Type:                          "EDGE",
	}

	tmURL := "https://myto.invalid"
	tmReverseProxyURL := "https://myrp.myto.invalid"
	locationParams := map[string]ConfigProfileParams{
		"remap.config": ConfigProfileParams{
			FileNameOnDisk: "remap.config",
			Location:       "/my/location/",
		},
		"regex_revalidate.config": ConfigProfileParams{
			FileNameOnDisk: "regex_revalidate.config",
			Location:       "/my/location/",
			URL:            "http://myurl/remap.config", // cdn-scoped endpoint
			APIURI:         "http://myapi/remap.config",
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
		"unknown.config": ConfigProfileParams{
			FileNameOnDisk: "unknown.config",
			Location:       "/my/location/",
		},
		"custom.config": ConfigProfileParams{
			FileNameOnDisk: "custom.config",
			Location:       "/my/location/",
		},
	}
	uriSignedDSes := []tc.DeliveryServiceName{"mydsname"}

	scopeParams := map[string]string{"custom.config": string(tc.ATSConfigMetaDataConfigFileScopeProfiles)}

	txt := MakeMetaConfig(serverHostName, server, tmURL, tmReverseProxyURL, locationParams, uriSignedDSes, scopeParams)

	cfg := tc.ATSConfigMetaData{}
	if err := json.Unmarshal([]byte(txt), &cfg); err != nil {
		t.Fatalf("MakeMetaConfig returned invalid JSON: " + err.Error())
	}

	if cfg.Info.ProfileID != int(server.ProfileID) {
		t.Errorf("expected Info.ProfileID %v actual %v", server.ProfileID, cfg.Info.ProfileID)
	}

	if cfg.Info.TOReverseProxyURL != tmReverseProxyURL {
		t.Errorf("expected Info.TOReverseProxyURL %v actual %v", tmReverseProxyURL, cfg.Info.TOReverseProxyURL)
	}

	if cfg.Info.TOURL != tmURL {
		t.Errorf("expected Info.TOURL %v actual %v", tmURL, cfg.Info.TOURL)
	}

	if server.IP != cfg.Info.ServerIPv4 {
		t.Errorf("expected Info.ServerIP %v actual %v", server.IP, cfg.Info.ServerIPv4)
	}

	if server.Port != cfg.Info.ServerPort {
		t.Errorf("expected Info.ServerPort %v actual %v", server.Port, cfg.Info.ServerPort)
	}

	if server.HostName != cfg.Info.ServerName {
		t.Errorf("expected Info.ServerName %v actual %v", server.HostName, cfg.Info.ServerName)
	}

	if cfg.Info.CDNID != server.CDNID {
		t.Errorf("expected Info.CDNID %v actual %v", server.CDNID, cfg.Info.CDNID)
	}

	if cfg.Info.CDNName != string(server.CDN) {
		t.Errorf("expected Info.CDNName %v actual %v", server.CDN, cfg.Info.CDNName)
	}
	if cfg.Info.ServerID != server.ID {
		t.Errorf("expected Info.ServerID %v actual %v", server.ID, cfg.Info.ServerID)
	}
	if cfg.Info.ProfileName != server.ProfileName {
		t.Errorf("expected Info.ProfileName %v actual %v", server.ProfileName, cfg.Info.ProfileName)
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
			if expected := "/my/location/"; cf.Location != expected {
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
	txt = MakeMetaConfig(serverHostName, server, tmURL, tmReverseProxyURL, locationParams, uriSignedDSes, scopeParams)
	cfg = tc.ATSConfigMetaData{}
	if err := json.Unmarshal([]byte(txt), &cfg); err != nil {
		t.Fatalf("MakeMetaConfig returned invalid JSON: " + err.Error())
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
}
