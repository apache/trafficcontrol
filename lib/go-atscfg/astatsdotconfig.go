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
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// AstatsSeparator is the string used by astats.config to separate
// configuration option names from values.
const AstatsSeparator = "="

// AstatsFileName is the name of the file used to contain configuration for the
// astats_over_http Traffic Server plugin.
const AstatsFileName = "astats.config"

// ContentTypeAstatsDotConfig is the MIME type of the content of astats.config
// file.
const ContentTypeAstatsDotConfig = ContentTypeTextASCII

// LineCommentAstatsDotConfig is the string used to mark the beginning of a
// line comment as understood by parsers of astats.config files.
const LineCommentAstatsDotConfig = LineCommentHash

// AStatsDotConfigOpts contains settings to configure generation options.
type AStatsDotConfigOpts struct {
	// HdrComment is the header comment to include at the beginning of the file.
	// This should be the text desired, without comment syntax (like # or //). The file's comment syntax will be added.
	// To omit the header comment, pass the empty string.
	HdrComment string
}

// MakeAStatsDotConfig constructs an astats.config file for the given server
// with the given Parameters and header comment content.
//
// TODO: Rename to 'MakeAstatsDotConfig' for consistency with other exported
// symbols?
func MakeAStatsDotConfig(
	server *Server,
	serverParams []tc.ParameterV5,
	opt *AStatsDotConfigOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &AStatsDotConfigOpts{}
	}
	warnings := []string{}

	if len(server.Profiles) == 0 {
		return Cfg{}, makeErr(warnings, "server missing Profile")
	}

	serverParams = filterParams(serverParams, AstatsFileName, "", "", "location")
	paramData, paramWarns := paramsToMap(serverParams)
	warnings = append(warnings, paramWarns...)
	hdr := makeHdrComment(opt.HdrComment)
	txt := genericProfileConfig(paramData, AstatsSeparator)
	if txt == "" {
		txt = "\n" // If no params exist, don't send "not found," but an empty file. We know the profile exists.
	}
	txt = hdr + txt

	return Cfg{
		Text:        txt,
		ContentType: ContentTypeAstatsDotConfig,
		LineComment: LineCommentAstatsDotConfig,
		Warnings:    warnings,
	}, nil
}
