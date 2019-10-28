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
	"math"
	"regexp"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const HeaderRewriteMidPrefix = "hdr_rw_mid_"

func MakeHeaderRewriteMidDotConfig(
	cdnName tc.CDNName,
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
	ds HeaderRewriteDS,
	assignedMids []HeaderRewriteServer, // the mids assigned to ds (mids whose cachegroup is the parent of the cachegroup of any edge assigned to this ds)
) string {
	text := GenericHeaderComment(string(cdnName), toToolName, toURL)

	// write a header rewrite rule if maxOriginConnections > 0 and the ds DOES use mids
	if ds.MaxOriginConnections > 0 && ds.Type.UsesMidCache() {
		dsOnlineMidCount := 0
		for _, sv := range assignedMids {
			if sv.Status == tc.CacheStatusReported || sv.Status == tc.CacheStatusOnline {
				dsOnlineMidCount++
			}
		}
		if dsOnlineMidCount > 0 {
			maxOriginConnectionsPerMid := int(math.Round(float64(ds.MaxOriginConnections) / float64(dsOnlineMidCount)))
			text += "cond %{REMAP_PSEUDO_HOOK}\nset-config proxy.config.http.origin_max_connections " + strconv.Itoa(maxOriginConnectionsPerMid)
			if ds.MidHeaderRewrite == "" {
				text += " [L]"
			} else {
				text += "\n"
			}
		}
	}

	// write the contents of ds.MidHeaderRewrite to hdr_rw_mid_xml-id.config replacing any instances of __RETURN__ (surrounded by spaces or not) with \n
	if ds.MidHeaderRewrite != "" {
		re := regexp.MustCompile(`\s*__RETURN__\s*`)
		text += re.ReplaceAllString(ds.MidHeaderRewrite, "\n")
	}

	text += "\n"
	return text
}
