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
)

const ContentTypeUnknownConfig = ContentTypeTextASCII

func MakeUnknownConfig(
	profileName string,
	paramData map[string]string, // GetProfileParamData(tx, profile.ID, AstatsFileName)
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
) string {
	hdr := GenericHeaderComment(profileName, toToolName, toURL)

	lines := []string{}
	for paramName, paramVal := range paramData {
		if paramName == "header" {
			if paramVal == "none" {
				hdr = ""
			} else {
				hdr = paramVal + "\n"
			}
		} else {
			lines = append(lines, paramVal+"\n")
		}
	}
	sort.Strings(lines)
	text := strings.Join(lines, "")
	text = strings.Replace(text, "__RETURN__", "\n", -1)
	return hdr + text
}

// GetUnknownConfigCommentType takes the same data as MakeUnknownConfig and returns the comment type for that config.
// In particular, it returns # unless there is a 'header' parameter, in which case it returns an empty string.
// Wwe don't actually know that the first characters of a custom header are a comment, or how many characters it might be.
func GetUnknownConfigCommentType(
	profileName string,
	paramData map[string]string,
	toToolName string,
	toURL string,
) string {
	for paramName, _ := range paramData {
		if paramName == "header" {
			return ""
		}
	}
	return LineCommentHash
}
