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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// ATSDotRulesFileName is the name of a "rules" configuration file.
//
// TODO: This isn't actually correct. This is just the ConfigFile value that a
// cache server's Profile's Drive_Prefix, Drive_Letters, RAM_Drive_Prefix, and
// RAM_Drive_Letters Parameters must have for generation to be successful (in
// fact it seems to panic if those aren't found), but the actual name of the
// file for which MakeATSDotRules outputs content is '50-ats.rules'. Is this
// misleading? Maybe it just doesn't even need to be exported? Or exist at all?
const ATSDotRulesFileName = StorageFileName

// ContentTypeATSDotRules is the MIME type of the contents of a 50-ats.rules
// file.
const ContentTypeATSDotRules = ContentTypeTextASCII

// LineCommentATSDotRules is the string used by parsers of 50-ats.rules to
// determine that the rest of the current line's content is a comment.
const LineCommentATSDotRules = LineCommentHash

// ATSDotRulesOpts contains settings to configure generation options.
type ATSDotRulesOpts struct {
	// HdrComment is the header comment to include at the beginning of the file.
	// This should be the text desired, without comment syntax (like # or //). The file's comment syntax will be added.
	// To omit the header comment, pass the empty string.
	HdrComment string
}

// MakeATSDotRules constructs a '50-ats.rules' file for the given server with
// the given parameters and header comment content.
func MakeATSDotRules(
	server *Server,
	serverParams []tc.ParameterV5,
	opt *ATSDotRulesOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &ATSDotRulesOpts{}
	}
	warnings := []string{}
	if len(server.Profiles) == 0 {
		return Cfg{}, makeErr(warnings, "server missing Profile")
	}

	serverParams = filterParams(serverParams, ATSDotRulesFileName, "", "", "location")
	paramData, paramWarns := paramsToMap(serverParams)
	warnings = append(warnings, paramWarns...)
	text := makeHdrComment(opt.HdrComment)

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
