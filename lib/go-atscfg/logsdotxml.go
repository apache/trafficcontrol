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
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// LogsXMLFileName is the name of a logging configuration file used by ATS
// version 6 and earlier.
//
// Deprecated: The versions of ATS that use this file are not supported, either
// by ATC or the ATS project itself.
const LogsXMLFileName = "logs_xml.config"

// ContentTypeLogsDotXML is the MIME content type of the logging configuration
// file used by ATS version 6 and earlier.
//
// Deprecated: The versions of ATS that use this file are not supported, either
// by ATC or the ATS project itself.
const ContentTypeLogsDotXML = `text/xml`

// LineCommentLogsDotXML is the string used in the grammar of a logging
// configuration file used by ATS version 6 and earlier to indicate the
// beginning of a comment node.
//
// Note that XML does not have "line" comments, only comment nodes - which MUST
// have the associated closing tag (-->).
//
// Deprecated: The versions of ATS that use this file are not supported, either
// by ATC or the ATS project itself.
const LineCommentLogsDotXML = `<!--`

// LogsXMLDotConfigOpts contains settings to configure generation options.
type LogsXMLDotConfigOpts struct {
	// HdrComment is the header comment to include at the beginning of the file.
	// This should be the text desired, without comment syntax (like # or //). The file's comment syntax will be added.
	// To omit the header comment, pass the empty string.
	HdrComment string
}

// MakeLogsXMLDotConfig constructs a logging configuration file used by ATS
// version 6 and earlier for the given server.
//
// serverParams is expected to be the map of Parameter Names to Values of all
// Parameters assigned to the given Profile, that have the ConfigFile
// "logging.config". That is, they must already be filtered BEFORE being passed
// in here.
//
// Deprecated: The versions of ATS that use this file are not supported, either
// by ATC or the ATS project itself.
func MakeLogsXMLDotConfig(
	server *Server,
	serverParams []tc.ParameterV5,
	opt *LogsXMLDotConfigOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &LogsXMLDotConfigOpts{}
	}
	warnings := []string{}

	if len(server.Profiles) == 0 {
		return Cfg{}, makeErr(warnings, "this server missing Profile")
	}

	paramData, paramWarns := paramsToMap(filterParams(serverParams, LogsXMLFileName, "", "", "location"))
	warnings = append(warnings, paramWarns...)

	// Note LineCommentLogsDotXML must be a single-line comment!
	// But this file doesn't have a single-line format, so we use <!-- for the header and promise it's on a single line
	// Note! if this file is ever changed to have multi-line comments, LineCommentLogsDotXML will have to be changed to the empty string.
	hdrComment := makeHdrComment(opt.HdrComment)
	hdrComment = strings.Replace(hdrComment, `# `, ``, -1)
	hdrComment = strings.Replace(hdrComment, "\n", ``, -1)
	text := "<!-- " + hdrComment + " -->\n"

	for i := 0; i < MaxLogObjects; i++ {
		logFormatField := "LogFormat"
		logObjectField := "LogObject"
		if i > 0 {
			iStr := strconv.Itoa(i)
			logFormatField += iStr
			logObjectField += iStr
		}

		logFormatName := paramData[logFormatField+".Name"]
		if logFormatName != "" {
			format := paramData[logFormatField+".Format"]
			format = strings.Replace(format, `"`, `\"`, -1)

			text += `<LogFormat>
  <Name = "` + logFormatName + `"/>
  <Format = "` + format + `"/>
</LogFormat>
`
		}

		if logObjectFileName := paramData[logObjectField+".Filename"]; logObjectFileName != "" {
			logObjectFormat := paramData[logObjectField+".Format"]
			logObjectRollingEnabled := paramData[logObjectField+".RollingEnabled"]
			logObjectRollingIntervalSec := paramData[logObjectField+".RollingIntervalSec"]
			logObjectRollingOffsetHr := paramData[logObjectField+".RollingOffsetHr"]
			logObjectRollingSizeMb := paramData[logObjectField+".RollingSizeMb"]
			logObjectHeader := paramData[logObjectField+".Header"]

			text += `<LogObject>
  <Format = "` + logObjectFormat + `"/>
  <Filename = "` + logObjectFileName + `"/>
`
			if logObjectRollingEnabled != "" {
				text += `  <RollingEnabled = ` + logObjectRollingEnabled + `/>
`
			}
			text += `  <RollingIntervalSec = ` + logObjectRollingIntervalSec + `/>
  <RollingOffsetHr = ` + logObjectRollingOffsetHr + `/>
  <RollingSizeMb = ` + logObjectRollingSizeMb + `/>
`
			if logObjectHeader != "" {
				text += `  <Header = "` + logObjectHeader + `"/>
`
			}
			text += `</LogObject>
`
		}
	}

	return Cfg{
		Text:        text,
		ContentType: ContentTypeLogsDotXML,
		LineComment: LineCommentLogsDotXML,
		Warnings:    warnings,
	}, nil
}
