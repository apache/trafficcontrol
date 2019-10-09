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
)

const RecordsSeparator = " "
const RecordsFileName = "records.config"

func MakeRecordsDotConfig(
	profileName string,
	paramData map[string]string, // GetProfileParamData(tx, profile.ID, StorageFileName)
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
) string {
	hdr := GenericHeaderComment(profileName, toToolName, toURL)
	txt := GenericProfileConfig(paramData, RecordsSeparator)
	if txt == "" {
		txt = "\n" // If no params exist, don't send "not found," but an empty file. We know the profile exists.
	}
	txt = replaceLineSuffixes(txt, "STRING __HOSTNAME__", "STRING __FULL_HOSTNAME__")
	txt = hdr + txt
	return txt
}

func replaceLineSuffixes(txt string, suffix string, newSuffix string) string {
	lines := strings.Split(txt, "\n")
	newLines := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.HasSuffix(line, suffix) {
			line = line[:len(line)-len(suffix)]
			line += newSuffix
		}
		newLines = append(newLines, line)
	}
	return strings.Join(newLines, "\n")
}
