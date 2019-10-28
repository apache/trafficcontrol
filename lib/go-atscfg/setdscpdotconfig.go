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
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func MakeSetDSCPDotConfig(
	cdnName tc.CDNName,
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
	dscpNumStr string,
) string {
	text := GenericHeaderComment(string(cdnName), toToolName, toURL)

	if _, err := strconv.Atoi(dscpNumStr); err != nil {
		// TODO warn? We don't generally warn for client errors, because it can be an attack vector. Provide a more informative error return? Return a 404?
		text = "An error occured generating the DSCP header rewrite file."
		dscpNumStr = "" // emulates Perl
	}

	text += `cond %{REMAP_PSEUDO_HOOK}` + "\n" + `set-conn-dscp ` + dscpNumStr + ` [L]` + "\n"
	return text
}
