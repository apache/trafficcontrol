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
	"sort"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const ContentTypeURLSig = ContentTypeTextASCII
const LineCommentURLSig = LineCommentHash

func MakeURLSigConfig(
	profileName string,
	urlSigKeys tc.URLSigKeys,
	paramData map[string]string, // GetProfileParamData(tx, profile.ID, StorageFileName)
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
) string {
	hdr := GenericHeaderComment(profileName, toToolName, toURL)

	sep := " = "

	text := hdr

	paramLines := []string{}
	for paramName, paramVal := range paramData {
		if len(urlSigKeys) == 0 || !strings.HasPrefix(paramName, "key") {
			paramLines = append(paramLines, paramName+sep+paramVal+"\n")
		}
	}
	sort.Strings(paramLines)
	text += strings.Join(paramLines, "")

	keyLines := []string{}
	for key, val := range urlSigKeys {
		keyLines = append(keyLines, key+sep+val+"\n")
	}
	sort.Strings(keyLines)
	text += strings.Join(keyLines, "")

	return text
}
