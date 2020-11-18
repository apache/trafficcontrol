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

	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestMakeCacheURLDotConfigWithDS(t *testing.T) {
	server := makeGenericServer()
	cdnName := "mycdn"
	server.CDNName = &cdnName

	hdr := "myHeaderComment"

	fileName := "cacheurl_myds.config"

	ds0 := makeGenericDS()
	ds0.ID = util.IntPtr(420)
	ds0.XMLID = util.StrPtr("myds")
	ds0.OrgServerFQDN = util.StrPtr("http://myorigin.example.net")
	ds0.QStringIgnore = util.IntPtr(0)
	ds0.CacheURL = util.StrPtr("http://mycacheurl.net")

	servers := []Server{*server}
	dses := []DeliveryService{*ds0}
	dss := makeDSS(servers, dses)

	cfg, err := MakeCacheURLDotConfig(fileName, server, dses, dss, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	if !strings.Contains(txt, hdr) {
		t.Errorf("expected: header comment text '"+hdr+"', actual: %v", txt)
	}

	if !strings.HasPrefix(strings.TrimSpace(txt), "#") {
		t.Errorf("expected: header comment, actual: %v", txt)
	}

	if !strings.Contains(txt, "mycacheurl") {
		t.Errorf("expected: contains cacheurl, actual: %v", txt)
	}
}

func TestMakeCacheURLDotConfigGlobalFile(t *testing.T) {
	server := makeGenericServer()
	cdnName := "mycdn"
	server.CDNName = &cdnName

	hdr := "myHeaderComment"

	fileName := "cacheurl.config"

	ds0 := makeGenericDS()
	ds0.ID = util.IntPtr(420)
	ds0.XMLID = util.StrPtr("ds0")
	ds0.OrgServerFQDN = util.StrPtr("http://myorigin.example.net")
	ds0.QStringIgnore = util.IntPtr(1)
	ds0.CacheURL = util.StrPtr("http://mycacheurl.net")

	servers := []Server{*server}
	dses := []DeliveryService{*ds0}
	dss := makeDSS(servers, dses)

	cfg, err := MakeCacheURLDotConfig(fileName, server, dses, dss, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	if !strings.Contains(txt, hdr) {
		t.Errorf("expected: header comment text '"+hdr+"', actual: %v", txt)
	}

	if !strings.HasPrefix(strings.TrimSpace(txt), "#") {
		t.Errorf("expected: header comment, actual: %v", txt)
	}

	if !strings.Contains(txt, "myorigin") {
		t.Errorf("expected: contains origin, actual: %v", txt)
	}

	if strings.Contains(txt, "mycacheurl") {
		t.Errorf("expected: global file to NOT contain cacheurl, actual: contains cacheurl")
	}
}

func TestMakeCacheURLDotConfigGlobalFileNoQStringIgnore(t *testing.T) {
	server := makeGenericServer()
	cdnName := "mycdn"
	server.CDNName = &cdnName

	hdr := "myHeaderComment"

	fileName := "cacheurl.config"

	ds0 := makeGenericDS()
	ds0.ID = util.IntPtr(420)
	ds0.XMLID = util.StrPtr("ds0")
	ds0.OrgServerFQDN = util.StrPtr("http://myorigin.example.net")
	ds0.QStringIgnore = util.IntPtr(0)
	ds0.CacheURL = util.StrPtr("http://mycacheurl.net")

	servers := []Server{*server}
	dses := []DeliveryService{*ds0}
	dss := makeDSS(servers, dses)

	cfg, err := MakeCacheURLDotConfig(fileName, server, dses, dss, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	if !strings.Contains(txt, hdr) {
		t.Errorf("expected: header comment text '" + hdr + "', actual: missing")
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
	server := makeGenericServer()
	cdnName := "mycdn"
	server.CDNName = &cdnName

	hdr := "myHeaderComment"

	fileName := "cacheurl_qstring.config"

	ds0 := makeGenericDS()
	ds0.ID = util.IntPtr(420)
	ds0.XMLID = util.StrPtr("ds0")
	ds0.OrgServerFQDN = util.StrPtr("http://myorigin.example.net")
	ds0.QStringIgnore = util.IntPtr(0)
	ds0.CacheURL = util.StrPtr("http://mycacheurl.net")

	servers := []Server{*server}
	dses := []DeliveryService{*ds0}
	dss := makeDSS(servers, dses)

	cfg, err := MakeCacheURLDotConfig(fileName, server, dses, dss, hdr)
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	if !strings.Contains(txt, hdr) {
		t.Errorf("expected: header comment '" + hdr + "', actual: missing")
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
