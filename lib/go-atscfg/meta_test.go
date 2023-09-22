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

func TestMakeMetaConfig(t *testing.T) {
	server := &Server{}
	server.CacheGroupID = 42
	server.CacheGroup = "cg0"
	server.CDN = "mycdn"
	server.CDNID = 43
	server.DomainName = "myserverdomain.invalid"
	server.HostName = "myserver"
	server.HTTPSPort = util.Ptr(443)
	server.ID = 44
	ip := "192.168.2.9"
	setIP(server, ip)
	// server.ParentCacheGroupID=            45
	// server.ParentCacheGroupType=          "MID_LOC"
	//server.ProfileID = util.Ptr(46)
	server.Profiles = []string{"myserverprofile"}
	server.TCPPort = util.Ptr(80)
	// server.SecondaryParentCacheGroupID=   47
	// server.SecondaryParentCacheGroupType= "MID_LOC"
	server.Type = "EDGE"

	// uriSignedDSes := []tc.DeliveryServiceName{"mydsname"}
	// dses := map[tc.DeliveryServiceName]DeliveryService{"mydsname": {}}

	cgs := []tc.CacheGroupNullableV5{}
	topologies := []tc.TopologyV5{}

	cfgPath := "/etc/foo/trafficserver"

	deliveryServices := []DeliveryService{}
	dss := []DeliveryServiceServer{}
	globalParams := []tc.ParameterV5{}

	makeLocationParam := func(name string) tc.ParameterV5 {
		return tc.ParameterV5{
			Name:       "location",
			ConfigFile: name,
			Value:      "/my/location/",
			Profiles:   []byte(`["` + server.Profiles[0] + `"]`),
		}
	}

	serverParams := []tc.ParameterV5{
		makeLocationParam("ssl_multicert.config"),
		makeLocationParam("volume.config"),
		makeLocationParam("ip_allow.config"),
		makeLocationParam("cache.config"),
		makeLocationParam("regex_revalidate.config"),
		makeLocationParam("uri_signing_mydsname.config"),
		makeLocationParam("uri_signing_nonexistentds.config"),
		makeLocationParam("regex_remap_nonexistentds.config"),
		makeLocationParam("url_sig_nonexistentds.config"),
		makeLocationParam("hdr_rw_nonexistentds.config"),
		makeLocationParam("hdr_rw_mid_nonexistentds.config"),
		makeLocationParam("unknown.config"),
		makeLocationParam("custom.config"),
		makeLocationParam("external.config"),
	}

	cfg, _, err := MakeConfigFilesList(cfgPath, server, serverParams, deliveryServices, dss, globalParams, cgs, topologies, &ConfigFilesListOpts{})
	if err != nil {
		t.Fatalf("MakeConfigFilesList: " + err.Error())
	}

	expectedConfigs := map[string]func(cf CfgMeta){
		"cache.config": func(cf CfgMeta) {
			if expected := "/my/location/"; cf.Path != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Path)
			}
		},
		"ip_allow.config": func(cf CfgMeta) {
			if expected := "/my/location/"; cf.Path != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Path)
			}
		},
		"volume.config": func(cf CfgMeta) {
			if expected := "/my/location/"; cf.Path != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Path)
			}
		},
		"ssl_multicert.config": func(cf CfgMeta) {
			if expected := "/my/location/"; cf.Path != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Path)
			}
		},
		"uri_signing_mydsname.config": func(cf CfgMeta) {
			if expected := "/my/location/"; cf.Path != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Path)
			}
		},
		"unknown.config": func(cf CfgMeta) {
			if expected := "/my/location/"; cf.Path != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Path)
			}
		},
		"custom.config": func(cf CfgMeta) {
			if expected := "/my/location/"; cf.Path != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Path)
			}
		},
		"remap.config": func(cf CfgMeta) {
			if expected := cfgPath; cf.Path != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Path)
			}
		},
		"regex_revalidate.config": func(cf CfgMeta) {
			if expected := "/my/location/"; cf.Path != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Path)
			}
		},
		"external.config": func(cf CfgMeta) {
			if expected := "/my/location/"; cf.Path != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Path)
			}
		},
		"hosting.config": func(cf CfgMeta) {
			if expected := cfgPath; cf.Path != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Path)
			}
		},
		"parent.config": func(cf CfgMeta) {
			if expected := cfgPath; cf.Path != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Path)
			}
		},
		"plugin.config": func(cf CfgMeta) {
			if expected := cfgPath; cf.Path != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Path)
			}
		},
		"records.config": func(cf CfgMeta) {
			if expected := cfgPath; cf.Path != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Path)
			}
		},
		"storage.config": func(cf CfgMeta) {
			if expected := cfgPath; cf.Path != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Path)
			}
		},
		"ssl_server_name.yaml": func(cf CfgMeta) {
			if expected := cfgPath; cf.Path != expected {
				t.Errorf("expected location '%v', actual '%v'", expected, cf.Path)
			}
		},
	}

	for _, cfgFile := range cfg {
		if testF, ok := expectedConfigs[cfgFile.Name]; !ok {
			t.Errorf("unexpected config '" + cfgFile.Name + "'")
		} else {
			testF(cfgFile)
			delete(expectedConfigs, cfgFile.Name)
		}
	}

	server.Type = "MID"
	cfg, _, err = MakeConfigFilesList(cfgPath, server, serverParams, deliveryServices, dss, globalParams, cgs, topologies, &ConfigFilesListOpts{})
	if err != nil {
		t.Fatalf("MakeConfigFilesList: " + err.Error())
	}
	for _, cfgFile := range cfg {
		if cfgFile.Name != "cache.config" {
			continue
		}
		break
	}
	for _, fi := range cfg {
		if strings.Contains(fi.Name, "nonexistentds") {
			t.Errorf("expected location parameters for nonexistent delivery services to not be added to config, actual '%v'", fi.Name)
		}
	}
}
