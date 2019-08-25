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

func TestMakeServerCacheDotConfig(t *testing.T) {
	serverName := tc.CacheName("server0")
	toToolName := "to0"
	toURL := "trafficops.example.net"

	dses := map[tc.DeliveryServiceName]ServerCacheConfigDS{
		"ds0": ServerCacheConfigDS{
			OrgServerFQDN: "https://ds0.example.test/path",
			Type:          tc.DSTypeHTTP,
		},
		"ds1": ServerCacheConfigDS{
			OrgServerFQDN: "https://ds1.example.test:42/path",
			Type:          tc.DSTypeDNS,
		},
		"ds2": ServerCacheConfigDS{
			OrgServerFQDN: "https://ds2.example.test:42",
			Type:          tc.DSTypeHTTP,
		},
		"ds3": ServerCacheConfigDS{
			OrgServerFQDN: "https://ds3.example.test",
			Type:          tc.DSTypeHTTP,
		},
		"ds4": ServerCacheConfigDS{
			OrgServerFQDN: "https://ds4.example.test/",
			Type:          tc.DSTypeHTTP,
		},
		"ds5": ServerCacheConfigDS{
			OrgServerFQDN: "http://ds5.example.test:1234/",
			Type:          tc.DSTypeHTTP,
		},
		"ds6": ServerCacheConfigDS{
			OrgServerFQDN: "ds6.example.test",
			Type:          tc.DSTypeHTTP,
		},
		"ds7": ServerCacheConfigDS{
			OrgServerFQDN: "ds7.example.test:80",
			Type:          tc.DSTypeHTTP,
		},
		"ds8": ServerCacheConfigDS{
			OrgServerFQDN: "ds8.example.test:8080/path",
			Type:          tc.DSTypeHTTP,
		},
		"ds-nocache": ServerCacheConfigDS{
			OrgServerFQDN: "http://ds-nocache.example.test",
			Type:          tc.DSTypeHTTPNoCache,
		},
	}

	txt := MakeServerCacheDotConfig(serverName, toToolName, toURL, dses)

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
		"42":       {},
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
