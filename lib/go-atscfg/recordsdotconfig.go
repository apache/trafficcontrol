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

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const RecordsSeparator = " "
const RecordsFileName = "records.config"
const ContentTypeRecordsDotConfig = ContentTypeTextASCII
const LineCommentRecordsDotConfig = LineCommentHash

func MakeRecordsDotConfig(
	server *tc.ServerNullable,
	serverParams []tc.Parameter,
	hdrComment string,
) (Cfg, error) {
	warnings := []string{}
	if server.Profile == nil {
		return Cfg{}, makeErr(warnings, "server profile missing")
	}

	params, paramWarns := ParamsToMap(FilterParams(serverParams, RecordsFileName, "", "", "location"))
	warnings = append(warnings, paramWarns...)

	hdr := makeHdrComment(hdrComment)
	txt := GenericProfileConfig(params, RecordsSeparator)
	if txt == "" {
		txt = "\n" // If no params exist, don't send "not found," but an empty file. We know the profile exists.
	}
	txt = replaceLineSuffixes(txt, "STRING __HOSTNAME__", "STRING __FULL_HOSTNAME__")
	txt = hdr + txt

	txt, overrideWarns := addRecordsDotConfigOverrides(txt, server)
	warnings = append(warnings, overrideWarns...)

	return Cfg{
		Text:        txt,
		ContentType: ContentTypeRecordsDotConfig,
		LineComment: LineCommentRecordsDotConfig,
		Warnings:    warnings,
	}, nil
}

// addRecordsDotConfigOverrides modifies the records.config text and adds any overrides.
// Returns the modified text and any warnings.
func addRecordsDotConfigOverrides(txt string, server *tc.ServerNullable) (string, []string) {
	warnings := []string{}
	txt, ipWarns := addRecordsDotConfigOutgoingIP(txt, server)
	warnings = append(warnings, ipWarns...)
	return txt, warnings
}

// addRecordsDotConfigOutgoingIP returns the outgoing IP added to the config text, and any warnings.
func addRecordsDotConfigOutgoingIP(txt string, server *tc.ServerNullable) (string, []string) {
	warnings := []string{}

	outgoingIPConfig := `proxy.local.outgoing_ip_to_bind`
	if strings.Contains(txt, outgoingIPConfig) {
		warnings = append(warnings, "records.config had a proxy.local.outgoing_ip_to_bind Parameter! Using Parameter, not setting Outgoing IP from Server")
		return txt, warnings
	}

	v4, v6 := getServiceAddresses(server)
	if v4 == nil {
		warnings = append(warnings, "server had no IPv4 service address, cannot set "+outgoingIPConfig+"!")
		return txt, warnings
	}

	txt = txt + `LOCAL ` + outgoingIPConfig + ` STRING ` + v4.String()
	if v6 != nil {
		txt += ` [` + v6.String() + `]`
	}
	txt += "\n"
	return txt, warnings
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
