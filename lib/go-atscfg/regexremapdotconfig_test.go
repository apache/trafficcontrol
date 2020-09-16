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
)

func TestMakeRegexRemapDotConfig(t *testing.T) {
	cdnName := tc.CDNName("mycdn")
	toToolName := "my-to"
	toURL := "my-to.example.net"

	fileName := "regex_remap_myds.config"

	dses := map[tc.DeliveryServiceName]CDNDS{
		"myds": CDNDS{
			OrgServerFQDN: "https://myorigin.example.net", // DS "origin_server_fqdn" is actually a URL including the scheme, the name is wrong.
			QStringIgnore: 0,
			CacheURL:      "https://mycacheurl.net",
			RegexRemap:    "myregexremap",
		},
	}

	txt := MakeRegexRemapDotConfig(cdnName, toToolName, toURL, fileName, dses)

	if !strings.Contains(txt, string(cdnName)) {
		t.Errorf("expected: cdnName '" + string(cdnName) + "', actual: missing")
	}
	if !strings.Contains(txt, toToolName) {
		t.Errorf("expected: toToolName '" + toToolName + "', actual: missing")
	}
	if !strings.Contains(txt, toURL) {
		t.Errorf("expected: toURL '" + toURL + "', actual: missing")
	}
	if !strings.HasPrefix(strings.TrimSpace(txt), "#") {
		t.Errorf("expected: header comment, actual: missing")
	}

	if strings.Contains(txt, "mycacheurl") {
		t.Errorf("expected: regex remap to not contain cacheurl, actual: '%v'", txt)
	}
	if strings.Contains(txt, "myorigin") {
		t.Errorf("expected: regex remap to not contain org server fqdn, actual: '%v'", txt)
	}
	if !strings.Contains(txt, "myregexremap") {
		t.Errorf("expected: regex remap to contain regex remap, actual: '%v'", txt)
	}
}

func TestMakeRegexRemapDotConfigUnusedDS(t *testing.T) {
	cdnName := tc.CDNName("mycdn")
	toToolName := "my-to"
	toURL := "my-to.example.net"

	fileName := "regex_remap_myds.config"

	dses := map[tc.DeliveryServiceName]CDNDS{
		"myds": CDNDS{
			OrgServerFQDN: "https://myorigin.example.net", // DS "origin_server_fqdn" is actually a URL including the scheme, the name is wrong.
			QStringIgnore: 0,
			CacheURL:      "https://mycacheurl.net",
			RegexRemap:    "myregexremap",
		},
		"otherds": CDNDS{
			OrgServerFQDN: "https://otherorigin.example.net", // DS "origin_server_fqdn" is actually a URL including the scheme, the name is wrong.
			QStringIgnore: 0,
			CacheURL:      "https://othercacheurl.net",
			RegexRemap:    "otherregexremap",
		},
	}

	txt := MakeRegexRemapDotConfig(cdnName, toToolName, toURL, fileName, dses)

	if !strings.Contains(txt, string(cdnName)) {
		t.Errorf("expected: cdnName '" + string(cdnName) + "', actual: missing")
	}
	if !strings.Contains(txt, toToolName) {
		t.Errorf("expected: toToolName '" + toToolName + "', actual: missing")
	}
	if !strings.Contains(txt, toURL) {
		t.Errorf("expected: toURL '" + toURL + "', actual: missing")
	}
	if !strings.HasPrefix(strings.TrimSpace(txt), "#") {
		t.Errorf("expected: header comment, actual: missing")
	}

	if strings.Contains(txt, "mycacheurl") {
		t.Errorf("expected: regex remap to not contain cacheurl, actual: '%v'", txt)
	}
	if strings.Contains(txt, "myorigin") {
		t.Errorf("expected: regex remap to not contain org server fqdn, actual: '%v'", txt)
	}
	if !strings.Contains(txt, "myregexremap") {
		t.Errorf("expected: regex remap to contain regex remap, actual: '%v'", txt)
	}

	if strings.Contains(txt, "mycacheurl") {
		t.Errorf("expected: regex remap to not contain other cacheurl, actual: '%v'", txt)
	}
	if strings.Contains(txt, "myorigin") {
		t.Errorf("expected: regex remap to not contain other org server fqdn, actual: '%v'", txt)
	}
	if strings.Contains(txt, "otherregexremap") {
		t.Errorf("expected: regex remap to contain other regex remap, actual: '%v'", txt)
	}
}

func TestMakeRegexRemapDotConfigReplaceReturns(t *testing.T) {
	cdnName := tc.CDNName("mycdn")
	toToolName := "my-to"
	toURL := "my-to.example.net"

	fileName := "regex_remap_myds.config"

	dses := map[tc.DeliveryServiceName]CDNDS{
		"myds": CDNDS{
			OrgServerFQDN: "https://myorigin.example.net", // DS "origin_server_fqdn" is actually a URL including the scheme, the name is wrong.
			QStringIgnore: 0,
			CacheURL:      "https://mycacheurl.net",
			RegexRemap:    "myregexremap__RETURN__mypostnewline",
		},
	}

	txt := MakeRegexRemapDotConfig(cdnName, toToolName, toURL, fileName, dses)

	if !strings.Contains(txt, string(cdnName)) {
		t.Errorf("expected: cdnName '" + string(cdnName) + "', actual: missing")
	}
	if !strings.Contains(txt, toToolName) {
		t.Errorf("expected: toToolName '" + toToolName + "', actual: missing")
	}
	if !strings.Contains(txt, toURL) {
		t.Errorf("expected: toURL '" + toURL + "', actual: missing")
	}
	if !strings.HasPrefix(strings.TrimSpace(txt), "#") {
		t.Errorf("expected: header comment, actual: missing")
	}

	if strings.Contains(txt, "mycacheurl") {
		t.Errorf("expected: regex remap to not contain cacheurl, actual: '%v'", txt)
	}
	if strings.Contains(txt, "myorigin") {
		t.Errorf("expected: regex remap to not contain org server fqdn, actual: '%v'", txt)
	}
	if !strings.Contains(txt, "myregexremap") {
		t.Errorf("expected: regex remap to contain regex remap, actual: '%v'", txt)
	}

	if strings.Contains(txt, "__RETURN__") {
		t.Errorf("expected: regex remap to replace __RETURN__, actual: '%v'", txt)
	}
	if !strings.Contains(txt, "myregexremap\nmypostnewline") {
		t.Errorf("expected: regex remap to replace __RETURN__ with newline, actual: '%v'", txt)
	}
}
