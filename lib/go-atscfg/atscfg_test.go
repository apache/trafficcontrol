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
	server.ProfileNames = &[]string{"serverprofile"}
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
