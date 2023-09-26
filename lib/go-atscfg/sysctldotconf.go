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

// SysctlSeparator is the string used to separate Parameter Names from their
// Values on lines of a sysctl.conf ATS configuration file.
const SysctlSeparator = " = "

// SysctlFileName is the ConfigFile of Parameters which, if found on a server's
// Profile, specify lines in the sysctl.conf ATS configuration file.
const SysctlFileName = "sysctl.conf"

// ContentTypeSysctlDotConf is the MIME type of the contents of a sysctl.conf
// ATS configuration file.
const ContentTypeSysctlDotConf = ContentTypeTextASCII

// LineCommentSysctlDotConf is the string understood by parses of sysctl.conf to
// be the beginning of a line comment.
const LineCommentSysctlDotConf = LineCommentHash

// SysCtlDotConfOpts contains settings to configure generation options.
type SysCtlDotConfOpts struct {
	// HdrComment is the header comment to include at the beginning of the file.
	// This should be the text desired, without comment syntax (like # or //). The file's comment syntax will be added.
	// To omit the header comment, pass the empty string.
	HdrComment string
}

// MakeSysCtlDotConf generates a sysctl.conf ATS configuration file for the
// given server with the given Parameters.
func MakeSysCtlDotConf(
	server *Server,
	serverParams []tc.ParameterV5,
	opt *SysCtlDotConfOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &SysCtlDotConfOpts{}
	}
	warnings := []string{}
	if len(server.Profiles) == 0 {
		return Cfg{}, makeErr(warnings, "server missing Profiles")
	}

	paramData, paramWarns := paramsToMap(filterParams(serverParams, SysctlFileName, "", "", "location"))
	warnings = append(warnings, paramWarns...)

	hdr := makeHdrComment(opt.HdrComment)
	txt := genericProfileConfig(paramData, SysctlSeparator)
	if txt == "" {
		txt = "\n" // If no params exist, don't send "not found," but an empty file. We know the profile exists.
	}
	txt = hdr + txt

	return Cfg{
		Text:        txt,
		ContentType: ContentTypeSysctlDotConf,
		LineComment: LineCommentSysctlDotConf,
		Warnings:    warnings,
	}, nil
}
