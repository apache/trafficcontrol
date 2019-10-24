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

func TestMakeCacheURLDotConfigWithDS(t *testing.T) {
	cdnName := tc.CDNName("mycdn")
	toToolName := "my-to"
	toURL := "my-to.example.net"

	fileName := "cacheurl_myds.config"

	dses := map[tc.DeliveryServiceName]CacheURLDS{
		"myds": CacheURLDS{
			OrgServerFQDN: "https://myorigin.example.net", // DS "origin_server_fqdn" is actually a URL including the scheme, the name is wrong.
			QStringIgnore: 0,
			CacheURL:      "https://mycacheurl.net",
		},
	}

	txt := MakeCacheURLDotConfig(cdnName, toToolName, toURL, fileName, dses)

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

	if !strings.Contains(txt, "mycacheurl") {
		t.Errorf("expected: contains cacheurl, actual: missing")
	}
}

func TestMakeCacheURLDotConfigGlobalFile(t *testing.T) {
	cdnName := tc.CDNName("mycdn")
	toToolName := "my-to"
	toURL := "my-to.example.net"

	fileName := "cacheurl.config"

	dses := map[tc.DeliveryServiceName]CacheURLDS{
		"myds": CacheURLDS{
			OrgServerFQDN: "https://myorigin.example.net", // DS "origin_server_fqdn" is actually a URL including the scheme, the name is wrong.
			QStringIgnore: 1,
			CacheURL:      "https://mycacheurl.net",
		},
	}

	txt := MakeCacheURLDotConfig(cdnName, toToolName, toURL, fileName, dses)

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

	if !strings.Contains(txt, "myorigin") {
		t.Errorf("expected: contains origin, actual: missing")
	}

	if strings.Contains(txt, "mycacheurl") {
		t.Errorf("expected: global file to NOT contain cacheurl, actual: contains cacheurl")
	}
}

func TestMakeCacheURLDotConfigGlobalFileNoQStringIgnore(t *testing.T) {
	cdnName := tc.CDNName("mycdn")
	toToolName := "my-to"
	toURL := "my-to.example.net"

	fileName := "cacheurl.config"

	dses := map[tc.DeliveryServiceName]CacheURLDS{
		"myds": CacheURLDS{
			OrgServerFQDN: "https://myorigin.example.net", // DS "origin_server_fqdn" is actually a URL including the scheme, the name is wrong.
			QStringIgnore: 0,
			CacheURL:      "https://mycacheurl.net",
		},
	}

	txt := MakeCacheURLDotConfig(cdnName, toToolName, toURL, fileName, dses)

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

	if strings.Contains(txt, "myorigin") {
		t.Errorf("expected: qstring ignore 0 to omit DS, actual: '%v'", txt)
	}

	if strings.Contains(txt, "mycacheurl") {
		t.Errorf("expected: global file to NOT contain cacheurl, actual: contains cacheurl")
	}
}

func TestMakeCacheURLDotConfigQStringFile(t *testing.T) {
	cdnName := tc.CDNName("mycdn")
	toToolName := "my-to"
	toURL := "my-to.example.net"

	fileName := "cacheurl_qstring.config"

	dses := map[tc.DeliveryServiceName]CacheURLDS{
		"myds": CacheURLDS{
			OrgServerFQDN: "https://myorigin.example.net", // DS "origin_server_fqdn" is actually a URL including the scheme, the name is wrong.
			QStringIgnore: 0,
			CacheURL:      "https://mycacheurl.net",
		},
	}

	txt := MakeCacheURLDotConfig(cdnName, toToolName, toURL, fileName, dses)

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

	if strings.Contains(txt, "myorigin") {
		t.Errorf("expected: qstring file to NOT contain origin, actual: '%v'", txt)
	}

	if strings.Contains(txt, "mycacheurl") {
		t.Errorf("expected: qstring file to NOT contain cacheurl, actual: '%v'", txt)
	}
}
