package varnishcfg

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
	"net"
	"reflect"
	"testing"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

func TestConfigureAccessControl(t *testing.T) {
	t.Run("edge server", func(t *testing.T) {
		vb := NewVCLBuilder(&t3cutil.ConfigData{
			Server: &atscfg.Server{Type: "EDGE"},
			ServerParams: []tc.ParameterV5{
				{
					ConfigFile: "ip_allow.config",
					Name:       "purge_allow_ip",
					Value:      "1.1.1.1,2.2.2.2,3.3.3.3/16",
				},
				{
					ConfigFile: "other_file",
					Name:       "not relevant",
				},
			},
		})
		vclFile := newVCLFile(defaultVCLVersion)
		vb.configureAccessControl(&vclFile)
		expectedVCLFile := newVCLFile(defaultVCLVersion)
		expectedVCLFile.acls["allow_all"] = []string{
			`"127.0.0.1"`,
			`"::1"`,
			`"1.1.1.1"`,
			`"2.2.2.2"`,
			`"3.3.3.3"/16`,
		}
		expectedVCLFile.subroutines["vcl_recv"] = []string{
			`if ((req.method == "PUSH" || req.method == "PURGE" || req.method == "DELETE") && !client.ip ~ allow_all) {`,
			`	return (synth(405));`,
			`}`,
			`if (req.method == "PURGE") {`,
			`	return (purge);`,
			`}`,
		}
		if !reflect.DeepEqual(vclFile, expectedVCLFile) {
			t.Errorf("got %v want %v", vclFile, expectedVCLFile)
		}
	})
	t.Run("mid server", func(t *testing.T) {
		vb := NewVCLBuilder(&t3cutil.ConfigData{
			Server: &atscfg.Server{
				Type:       "MID",
				HostName:   "server0",
				CacheGroup: "cg0",
			},
			ServerParams: []tc.ParameterV5{
				{
					ConfigFile: "ip_allow.config",
					Name:       atscfg.ParamPurgeAllowIP,
					Value:      "1.1.1.1,2.2.2.2",
				},
			},
			CacheGroups: []tc.CacheGroupNullableV5{
				{Name: util.Ptr("cg0")},
			},
			Servers: []atscfg.Server{
				{
					HostName:   "child0",
					CacheGroup: "childcg",
					Type:       tc.MonitorTypeName,
					Interfaces: []tc.ServerInterfaceInfoV40{
						{
							ServerInterfaceInfo: tc.ServerInterfaceInfo{
								Name:        "eth0",
								IPAddresses: []tc.ServerIPAddress{{Address: "1.1.1.1"}},
							},
						},
					},
				},
				{
					HostName:   "child1",
					CacheGroup: "childcg",
					Type:       tc.MonitorTypeName,
					Interfaces: []tc.ServerInterfaceInfoV40{
						{
							ServerInterfaceInfo: tc.ServerInterfaceInfo{
								Name:        "eth0",
								IPAddresses: []tc.ServerIPAddress{{Address: "2001:DB8:2::2/64"}},
							},
						},
					},
				},
			},
		})
		vclFile := newVCLFile(defaultVCLVersion)
		warnings, err := vb.configureAccessControl(&vclFile)
		if len(warnings) > 0 {
			t.Errorf("got warnings %v", warnings)
		}
		if err != nil {
			t.Errorf("got error while configuring acl: %s", err)
		}
		expectedVCLFile := newVCLFile(defaultVCLVersion)
		expectedVCLFile.acls["allow_all"] = []string{
			`"127.0.0.1"`,
			`"::1"`,
			`"1.1.1.1"`,
			`"2.2.2.2"`,
		}
		expectedVCLFile.acls["allow_all_but_push_purge"] = []string{
			`"1.1.1.1"/32`,
			`"2001:db8:2::"/64`,
			`"10.0.0.0"/8`,
			`"172.16.0.0"/12`,
			`"192.168.0.0"/16`,
		}
		expectedVCLFile.subroutines["vcl_recv"] = []string{
			`if ((req.method == "PUSH" || req.method == "PURGE") && client.ip ~ allow_all_but_push_purge) {`,
			`	return (synth(405));`,
			`}`,
			`if (!client.ip ~ allow_all_but_push_purge && !client.ip ~ allow_all) {`,
			`	return (synth(405));`,
			`}`,
			`if (req.method == "PURGE") {`,
			`	return (purge);`,
			`}`,
		}
		if !reflect.DeepEqual(vclFile, expectedVCLFile) {
			t.Errorf("got %v want %v", vclFile, expectedVCLFile)
		}
	})
}

func TestCIDRsToVarnishCIDRs(t *testing.T) {
	cidrs := make([]*net.IPNet, 0)
	cidrsString := []string{
		"1.0.0.0/10",
		"192.168.1.0/24",
		"2002:0:0:1234::/64",
	}
	for _, cidrString := range cidrsString {
		_, cidr, err := net.ParseCIDR(cidrString)
		if err != nil {
			t.Errorf("failed to parse cidr %s", cidrString)
		}

		cidrs = append(cidrs, cidr)
	}
	gotCIDRs := cidrsToVarnishCIDRs(cidrs)
	expectedCIDRs := []string{
		`"1.0.0.0"/10`,
		`"192.168.1.0"/24`,
		`"2002:0:0:1234::"/64`,
	}
	if !reflect.DeepEqual(gotCIDRs, expectedCIDRs) {
		t.Errorf("got %v want %v", gotCIDRs, expectedCIDRs)
	}
}

func TestConfigureAccessControlForEdge(t *testing.T) {
	acls := make(map[string][]string)
	subroutines := make(map[string][]string)
	allowAllIPs := []string{
		"1.1.1.1",
		"2.2.2.2",
	}
	configureAccessControlForEdge(acls, subroutines, allowAllIPs)

	expectedACLs := map[string][]string{
		"allow_all": {"1.1.1.1", "2.2.2.2"},
	}
	expectedSubroutines := map[string][]string{
		"vcl_recv": {
			`if ((req.method == "PUSH" || req.method == "PURGE" || req.method == "DELETE") && !client.ip ~ allow_all) {`,
			`	return (synth(405));`,
			`}`,
			`if (req.method == "PURGE") {`,
			`	return (purge);`,
			`}`,
		},
	}

	if !reflect.DeepEqual(acls, expectedACLs) {
		t.Errorf("got %v want %v", acls, expectedACLs)
	}
	if !reflect.DeepEqual(subroutines, expectedSubroutines) {
		t.Errorf("got %v want %v", subroutines, expectedSubroutines)
	}

}
func TestConfigureAccessControlForMid(t *testing.T) {
	acls := make(map[string][]string)
	subroutines := make(map[string][]string)
	allowAllIPs := []string{
		"1.1.1.1",
		"2.2.2.2",
	}
	allowAllButPushPurge := []string{
		"3.3.3.3",
		"4.4.4.4",
	}

	configureAccessControlForMid(acls, subroutines, allowAllIPs, allowAllButPushPurge)

	expectedACLs := map[string][]string{
		"allow_all":                {"1.1.1.1", "2.2.2.2"},
		"allow_all_but_push_purge": {"3.3.3.3", "4.4.4.4"},
	}
	expectedSubroutines := map[string][]string{
		"vcl_recv": {
			`if ((req.method == "PUSH" || req.method == "PURGE") && client.ip ~ allow_all_but_push_purge) {`,
			`	return (synth(405));`,
			`}`,
			`if (!client.ip ~ allow_all_but_push_purge && !client.ip ~ allow_all) {`,
			`	return (synth(405));`,
			`}`,
			`if (req.method == "PURGE") {`,
			`	return (purge);`,
			`}`,
		},
	}

	if !reflect.DeepEqual(acls, expectedACLs) {
		t.Errorf("got %v want %v", acls, expectedACLs)
	}
	if !reflect.DeepEqual(subroutines, expectedSubroutines) {
		t.Errorf("got %v want %v", subroutines, expectedSubroutines)
	}

}
