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
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestGenericHeaderComment(t *testing.T) {
	commentTxt := "foo"
	txt := makeHdrComment(commentTxt)
	testComment(t, txt, commentTxt)
}

func testComment(t *testing.T, txt string, commentTxt string) {
	commentLine := strings.SplitN(txt, "\n", 2)[0] // SplitN always returns at least 1 element, no need to check len before indexing

	if !strings.HasPrefix(strings.TrimSpace(commentLine), "#") {
		t.Errorf("expected comment on first line, actual: '" + commentLine + "'")
	}
	if !strings.Contains(commentLine, commentTxt) {
		t.Errorf("expected comment text '" + commentTxt + "' in comment, actual: '" + commentLine + "'")
	}
}

func TestTrimParamUnderscoreNumSuffix(t *testing.T) {
	inputExpected := map[string]string{
		``:                         ``,
		`a`:                        `a`,
		`_`:                        `_`,
		`foo__`:                    `foo__`,
		`foo__1`:                   `foo`,
		`foo__1234567890`:          `foo`,
		`foo_1234567890`:           `foo_1234567890`,
		`foo__1234__1234567890`:    `foo__1234`,
		`foo__1234__1234567890_`:   `foo__1234__1234567890_`,
		`foo__1234__1234567890a`:   `foo__1234__1234567890a`,
		`foo__1234__1234567890__`:  `foo__1234__1234567890__`,
		`foo__1234__1234567890__a`: `foo__1234__1234567890__a`,
		`__`:                       `__`,
		`__9`:                      ``,
		`_9`:                       `_9`,
		`__35971234789124`:         ``,
		`a__35971234789124`:        `a`,
		`1234`:                     `1234`,
		`foo__asdf_1234`:           `foo__asdf_1234`,
	}

	for input, expected := range inputExpected {
		if actual := trimParamUnderscoreNumSuffix(input); expected != actual {
			t.Errorf("Expected '%v' Actual '%v'", expected, actual)
		}
	}
}

func TestGetATSMajorVersionFromATSVersion(t *testing.T) {
	inputExpected := map[string]int{
		`7.1.2-34.56abcde.el7.centos.x86_64`:    7,
		`8`:                                     8,
		`8.1`:                                   8,
		`10.1`:                                  10,
		`1234.1.2-34.56abcde.el7.centos.x86_64`: 1234,
	}
	errExpected := []string{
		"a7.1.2-34.56abcde.el7.centos.x86_64",
		`-7.1.2-34.56abcde.el7.centos.x86_64`,
		".7.1.2-34.56abcde.el7.centos.x86_64",
		"7a.1.2-34.56abcde.el7.centos.x86_64",
		"7-a.1.2-34.56abcde.el7.centos.x86_64",
		"7-2.1.2-34.56abcde.el7.centos.x86_64",
		"100-2.1.2-34.56abcde.el7.centos.x86_64",
		"7a",
		"",
		"-",
		".",
	}

	for input, expected := range inputExpected {
		if actual, err := getATSMajorVersionFromATSVersion(input); err != nil {
			t.Errorf("expected %v actual: error '%v'", expected, err)
		} else if actual != expected {
			t.Errorf("expected %v actual: %v", expected, actual)
		}
	}
	for _, input := range errExpected {
		if actual, err := getATSMajorVersionFromATSVersion(input); err == nil {
			t.Errorf("input %v expected: error, actual: nil error '%v'", input, actual)
		}
	}
}

func TestLayerProfiles(t *testing.T) {
	profileNames := []string{
		"FOO",
		"BAR",
		"BAZ",
	}

	allParams := []tc.Parameter{
		{
			ConfigFile: "cfg_a",
			ID:         1000,
			Name:       "param_a",
			Profiles:   json.RawMessage(`["FOO"]`),
			Value:      "alpha",
		},
		{
			ConfigFile: "cfg_a",
			ID:         1000,
			Name:       "param_a",
			Profiles:   json.RawMessage(`["BAR"]`),
			Value:      "beta",
		},
		{
			ConfigFile: "cfg_b",
			ID:         1000,
			Name:       "param_a",
			Profiles:   json.RawMessage(`["BAR"]`),
			Value:      "gamma",
		},
		{
			ConfigFile: "cfg_c",
			ID:         1000,
			Name:       "param_c",
			Profiles:   json.RawMessage(`["BAZ"]`),
			Value:      "epsilon",
		},
		{
			ConfigFile: "cfg_c",
			ID:         1000,
			Name:       "param_c",
			Profiles:   json.RawMessage(`["BAR"]`),
			Value:      "delta",
		},
		{
			ConfigFile: "cfg_a",
			ID:         1000,
			Name:       "param_b",
			Profiles:   json.RawMessage(`["BAR"]`),
			Value:      "zeta",
		},
		{
			ConfigFile: "cfg_d",
			ID:         1000,
			Name:       "param_d",
			Profiles:   json.RawMessage(`["BAR"]`),
			Value:      "eta",
		},
		{
			ConfigFile: "cfg_d",
			ID:         1000,
			Name:       "param_d",
			Profiles:   json.RawMessage(`["BAZ"]`),
			Value:      "theta",
		},
		{
			ConfigFile: "cfg_d",
			ID:         1000,
			Name:       "param_d",
			Profiles:   json.RawMessage(`["FOO"]`),
			Value:      "iota",
		},
	}

	// per the params, on the layered profiles "FOO,BAR,BAZ" (in that order):
	// (1) beta in BAR should override alpha FOO,
	//    because they share the ConfigFile+Name key and BAR is later in the layering
	// (2) gamma should be added, but not override
	//    because the key is ConfigFile+Name, which isn't on any previous profile
	// (3) epsilon in BAZ should override delta in BAR
	//    because they share the ConfigFile+Name key and BAZ is later in the layering
	//    - this tests the parameters being in a different order than (1)
	// (4) zeta should be added, but not override
	//    because the key is ConfigFile+Name, which isn't on any previous profile
	//    - this tests the name matching but not config file, reverse of (2).
	// (5) theta in BAZ should override eta and iota in FOO and BAR
	//    because they share the ConfigFile+Name key and BAZ is last in the layering
	//    - this tests multiple overrides

	layeredParams, err := LayerProfiles(profileNames, allParams)
	if err != nil {
		t.Fatalf("expected LayerProfiles nil error, actual: %+v", err)
	}

	vals := map[string]struct{}{}
	for _, param := range layeredParams {
		vals[param.Value] = struct{}{}
	}

	if _, ok := vals["alpha"]; ok {
		t.Errorf("expected: param 'beta' to override 'alpha', actual: alpha in layered parameters")
	}
	if _, ok := vals["beta"]; !ok {
		t.Errorf("expected: param 'beta' to override 'alpha', actual: beta not in layered parameters")
	}
	if _, ok := vals["gamma"]; !ok {
		t.Errorf("expected: param 'gamma' with no ConfigFile+Name in another profile to be in layered parameters, actual: gamma not in layered parameters")
	}
	if _, ok := vals["delta"]; ok {
		t.Errorf("expected: param 'epsilon' to override 'delta' in prior profile, actual: delta in layered parameters")
	}
	if _, ok := vals["epsilon"]; !ok {
		t.Errorf("expected: param 'epsilon' to override 'delta' in prior profile, actual: epsilon not in layered parameters")
	}
	if _, ok := vals["zeta"]; !ok {
		t.Errorf("expected: param 'zeta' with no ConfigFile+Name in another profile to be in layered parameters, actual: zeta not in layered parameters")
	}
	if _, ok := vals["theta"]; !ok {
		t.Errorf("expected: param 'theta' to override 'eta' and 'iota' in prior profile, actual: theta not in layered parameters")
	}
	if _, ok := vals["eta"]; ok {
		t.Errorf("expected: param 'theta' to override 'eta' and 'iota' in prior profile, actual: eta in layered parameters")
	}
	if _, ok := vals["iota"]; ok {
		t.Errorf("expected: param 'theta' to override 'eta' and 'iota' in prior profile, actual: iota in layered parameters")
	}
}

func setIP(sv *Server, ipAddress string) {
	setIPInfo(sv, "", ipAddress, "")
}

func setIP6(sv *Server, ip6Address string) {
	setIPInfo(sv, "", "", ip6Address)
}

func setIPInfo(sv *Server, interfaceName string, ipAddress string, ip6Address string) {
	if len(sv.Interfaces) == 0 {
		sv.Interfaces = []tc.ServerInterfaceInfoV40{}
		{
			si := tc.ServerInterfaceInfoV40{}
			si.Name = interfaceName
			sv.Interfaces = append(sv.Interfaces, si)
		}
	}
	if ipAddress != "" {
		sv.Interfaces[0].IPAddresses = append(sv.Interfaces[0].IPAddresses, tc.ServerIPAddress{
			Address:        ipAddress,
			Gateway:        nil,
			ServiceAddress: true,
		})
	}
	if ip6Address != "" {
		sv.Interfaces[0].IPAddresses = append(sv.Interfaces[0].IPAddresses, tc.ServerIPAddress{
			Address:        ip6Address,
			Gateway:        nil,
			ServiceAddress: true,
		})
	}
}

func makeGenericServer() *Server {
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

func makeGenericDS() *DeliveryService {
	ds := &DeliveryService{}
	ds.ID = util.IntPtr(42)
	ds.XMLID = util.StrPtr("ds1")
	ds.QStringIgnore = util.IntPtr(int(tc.QStringIgnoreDrop))
	ds.OrgServerFQDN = util.StrPtr("http://ds1.example.net")
	dsType := tc.DSTypeDNS
	ds.Type = &dsType
	ds.MultiSiteOrigin = util.BoolPtr(false)
	ds.Active = util.BoolPtr(true)
	return ds
}

// makeDSS creates DSS as an outer product of every server and ds given.
// The given servers and dses must all have non-nil, unique IDs.
func makeDSS(servers []Server, dses []DeliveryService) []DeliveryServiceServer {
	dss := []DeliveryServiceServer{}
	for _, sv := range servers {
		for _, ds := range dses {
			dss = append(dss, DeliveryServiceServer{
				Server:          *sv.ID,
				DeliveryService: *ds.ID,
			})
		}
	}
	return dss
}

func makeParamsFromMapArr(profile string, configFile string, paramM map[string][]string) []tc.Parameter {
	params := []tc.Parameter{}
	for name, vals := range paramM {
		for _, val := range vals {
			params = append(params, tc.Parameter{
				Name:       name,
				ConfigFile: configFile,
				Value:      val,
				Profiles:   []byte(`["` + profile + `"]`),
			})
		}
	}
	return params
}

func makeParamsFromMap(profile string, configFile string, paramM map[string]string) []tc.Parameter {
	params := []tc.Parameter{}
	for name, val := range paramM {
		params = append(params, tc.Parameter{
			Name:       name,
			ConfigFile: configFile,
			Value:      val,
			Profiles:   []byte(`["` + profile + `"]`),
		})
	}
	return params
}
