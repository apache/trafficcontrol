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

const ATSDotRulesFileName = StorageFileName
const ContentTypeATSDotRules = ContentTypeTextASCII
const LineCommentATSDotRules = LineCommentHash

func MakeATSDotRules(
	server *Server,
	serverParams []tc.Parameter,
	hdrComment string,
) (Cfg, error) {
	warnings := []string{}
	if server.Profile == nil {
		return Cfg{}, makeErr(warnings, "server missing Profile")
	}

	serverParams = filterParams(serverParams, ATSDotRulesFileName, "", "", "location")
	paramData, paramWarns := paramsToMap(serverParams)
	warnings = append(warnings, paramWarns...)
	text := makeHdrComment(hdrComment)

	drivePrefix := strings.TrimPrefix(paramData["Drive_Prefix"], `/dev/`)
	drivePostfix := strings.Split(paramData["Drive_Letters"], ",")

	for _, l := range drivePostfix {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		text += `KERNEL=="` + drivePrefix + l + `", OWNER="ats"` + "\n"
	}
	if ramPrefix, ok := paramData["RAM_Drive_Prefix"]; ok {
		ramPrefix = strings.TrimPrefix(ramPrefix, `/dev/`)
		ramPostfix := strings.Split(paramData["RAM_Drive_Letters"], ",")
		for _, l := range ramPostfix {
			text += `KERNEL=="` + ramPrefix + l + `", OWNER="ats"` + "\n"
		}
	}

	return Cfg{
		Text:        text,
		ContentType: ContentTypeATSDotRules,
		LineComment: LineCommentATSDotRules,
		Warnings:    warnings,
	}, nil
}
