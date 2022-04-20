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

import "strings"

const ContentType12MFacts = ContentTypeTextASCII
const LineComment12MFacts = LineCommentHash

// Config12MFactsOpts contains settings to configure generation options.
type Config12MFactsOpts struct {
	// HdrComment is the header comment to include at the beginning of the file.
	// This should be the text desired, without comment syntax (like # or //). The file's comment syntax will be added.
	// To omit the header comment, pass the empty string.
	HdrComment string
}

func Make12MFacts(
	server *Server,
	opt *Config12MFactsOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &Config12MFactsOpts{}
	}
	warnings := []string{}

	if len(server.ProfileNames) == 0 {
		return Cfg{}, makeErr(warnings, "this server missing Profile")
	}

	hdr := makeHdrComment(opt.HdrComment)
	txt := hdr
	txt += "profiles:" + strings.Join(server.ProfileNames, ", ") + "\n"

	return Cfg{
		Text:        txt,
		ContentType: ContentType12MFacts,
		LineComment: LineComment12MFacts,
		Warnings:    warnings,
	}, nil
}
