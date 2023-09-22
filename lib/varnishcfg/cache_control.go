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
	"fmt"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

func (v VCLBuilder) configureUncacheableDSes(vclFile *vclFile, profileDSes []atscfg.ProfileDS) []string {
	warnings := make([]string, 0)

	uncacheableHostsConditions := make([]string, 0)
	for _, ds := range profileDSes {
		if ds.Type != tc.DSTypeHTTPNoCache {
			continue
		}

		if ds.OriginFQDN == "" {
			warnings = append(warnings, fmt.Sprintf("ds %s has no origin fqdn, skipping!", ds.Name))
			continue
		}

		host, _ := atscfg.GetHostPortFromURI(ds.OriginFQDN)
		uncacheableHostsConditions = append(uncacheableHostsConditions, fmt.Sprintf(`bereq.http.host == "%s"`, host))
	}
	if len(uncacheableHostsConditions) == 0 {
		return warnings
	}
	berespLines := make([]string, 0)
	berespLines = append(berespLines, fmt.Sprintf(`if (%s) {`, strings.Join(uncacheableHostsConditions, " || ")))
	berespLines = append(berespLines, fmt.Sprint(`	set beresp.uncacheable = true;`))
	berespLines = append(berespLines, fmt.Sprint(`}`))
	vclFile.subroutines["vcl_backend_response"] = append(vclFile.subroutines["vcl_backend_response"], berespLines...)

	return warnings
}
