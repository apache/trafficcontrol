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

func TestMakeServerCacheDotConfig(t *testing.T) {
	serverName := "server0"
	hdr := "myHeaderComment"

	server := makeGenericServer()
	server.HostName = serverName
	server.Type = tc.MidTypePrefix + "_CUSTOM"

	makeDS := func(name string, origin string, dsType *string) DeliveryService {
		ds := makeGenericDS()
		ds.XMLID = name
		ds.OrgServerFQDN = util.Ptr(origin)
		ds.Type = dsType
		return *ds
	}

	dses := []DeliveryService{
		makeDS("ds0", "https://ds0.example.test/path", util.Ptr("HTTP")),
		makeDS("ds1", "https://ds1.example.test:4321/path", util.Ptr("DNS")),
		makeDS("ds2", "https://ds2.example.test:4321", util.Ptr("HTTP")),
		makeDS("ds3", "https://ds3.example.test", util.Ptr("HTTP")),
		makeDS("ds4", "https://ds4.example.test/", util.Ptr("HTTP")),
		makeDS("ds5", "http://ds5.example.test:1234/", util.Ptr("HTTP")),
		makeDS("ds6", "ds6.example.test", util.Ptr("HTTP")),
		makeDS("ds7", "ds7.example.test:80", util.Ptr("HTTP")),
		makeDS("ds8", "ds8.example.test:8080/path", util.Ptr("HTTP")),
		makeDS("ds-nocache", "http://ds-nocache.example.test", util.Ptr("HTTP_NO_CACHE")),
	}

	cfg, err := makeCacheDotConfigMid(server, dses, &CacheDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	lines := strings.Split(txt, "\n")

	unexpecteds := map[string]struct{}{
		"http://":  {},
		"https://": {},
		"path":     {},
		"/":        {},
		"ds0":      {},
		"ds1":      {},
		"ds2":      {},
		"ds3":      {},
		"ds4":      {},
		"ds5":      {},
		"ds6":      {},
		"ds7":      {},
		"ds8":      {},
		"4321":     {},
		"1234":     {},
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		for unexpected, _ := range unexpecteds {
			if strings.Contains(line, unexpected) {
				t.Errorf("expected NOT '%v' actual '%v'\n", unexpected, line)
			}
		}

		if !strings.Contains(line, "never-cache") && !strings.HasPrefix(line, "#") {
			t.Errorf("expected '%v' actual '%v'\n", "only nocache DSes", line)
		}
	}
}
