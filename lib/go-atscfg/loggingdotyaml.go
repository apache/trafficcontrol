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
	"fmt"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// LoggingYAMLFileName is the name of a logging configuration file used by ATS
// version 8 and later.
const LoggingYAMLFileName = "logging.yaml"

// ContentTypeLoggingDotYAML is the content type of the contents of a logging
// configuration file used by ATS version 8 and later.
//
// Note YAML has no IANA standard mime type. This is one of several common
// usages, and is likely to be the standardized value. If you're reading this,
// please check IANA to see if YAML has been added, and change this to the IANA
// definition if so. Also note we include 'charset=us-ascii' because YAML is
// commonly UTF-8, but ATS is likely to be unable to handle UTF.
const ContentTypeLoggingDotYAML = ContentTypeYAML

// LineCommentLoggingDotYAML is the string used to indicate the beginning of a
// line comment in the grammar of a logging configuration file used by ATS
// version 8 and later.
const LineCommentLoggingDotYAML = LineCommentHash

// LoggingDotYAMLOpts contains settings to configure generation options.
type LoggingDotYAMLOpts struct {
	// HdrComment is the header comment to include at the beginning of the file.
	// This should be the text desired, without comment syntax (like # or //). The file's comment syntax will be added.
	// To omit the header comment, pass the empty string.
	HdrComment string

	// ATSMajorVersion is the integral major version of Apache Traffic server,
	// used to generate the proper config for the proper version.
	//
	// If omitted or 0, the major version will be read from the Server's Profile Parameter config file 'package' name 'trafficserver'. If no such Parameter exists, the ATS version will default to 5.
	// This was the old Traffic Control behavior, before the version was specifiable externally.
	//
	ATSMajorVersion uint
}

// MakeLoggingDotYAML creates a logging.yaml for a given ATS Profile.
//
// serverParams is expected to be the map of Parameter Names to Values of all
// Parameters assigned to the given Profile, that have the ConfigFile
// "logging.config". That is, they must already be filtered BEFORE being passed
// in here.
func MakeLoggingDotYAML(
	server *Server,
	serverParams []tc.ParameterV5,
	opt *LoggingDotYAMLOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &LoggingDotYAMLOpts{}
	}
	warnings := []string{}
	requiredIndent := 0

	if server.HostName == "" {
		return Cfg{}, makeErr(warnings, "this server missing HostName")
	}

	if len(server.Profiles) == 0 {
		return Cfg{}, makeErr(warnings, "this server missing Profile")
	}

	paramData, paramWarns := paramsToMap(filterParams(serverParams, LoggingYAMLFileName, "", "", "location"))
	warnings = append(warnings, paramWarns...)

	hdr := makeHdrComment(opt.HdrComment)

	atsMajorVersion := getATSMajorVersion(opt.ATSMajorVersion, serverParams, &warnings)

	// note we use the same const as logs.xml - this isn't necessarily a requirement, and we may want to make separate variables in the future.
	maxLogObjects := MaxLogObjects

	text := hdr
	if atsMajorVersion >= 9 {
		text += "\nlogging:"
		requiredIndent += 2
	}

	indentSpaces := strings.Repeat(" ", requiredIndent)
	text += "\n" + indentSpaces + "formats: \n"
	for i := 0; i < maxLogObjects; i++ {
		logFormatField := "LogFormat"
		if i > 0 {
			logFormatField += strconv.Itoa(i)
		}
		logFormatName := paramData[logFormatField+".Name"]
		if logFormatName != "" {
			format := paramData[logFormatField+".Format"]
			if format == "" {
				// TODO determine if the line should be excluded. Perl includes it anyway, without checking.
				warnings = append(warnings, fmt.Sprintf("server '%v' profile has logging.yaml format '%v' Name Parameter but no Format Parameter. Setting blank Format!\n", server.HostName, logFormatField))
			}
			text += indentSpaces + " - name: " + logFormatName + " \n"
			text += indentSpaces + "   format: '" + format + "'\n"
		}
	}

	text += indentSpaces + "filters:\n"
	for i := 0; i < maxLogObjects; i++ {
		logFilterField := "LogFilter"
		if i > 0 {
			logFilterField += strconv.Itoa(i)
		}
		if logFilterName := paramData[logFilterField+".Name"]; logFilterName != "" {
			filter := paramData[logFilterField+".Filter"]
			if filter == "" {
				// TODO determine if the line should be excluded. Perl includes it anyway, without checking.
				warnings = append(warnings, fmt.Sprintf("server '%v' profile has logging.yaml filter '%v' Name Parameter but no Filter Parameter. Setting blank Filter!\n", server.HostName, logFilterField))
			}
			logFilterType := paramData[logFilterField+".Type"]
			if logFilterType == "" {
				logFilterType = "accept"
			}
			text += indentSpaces + " - name: " + logFilterName + "\n"
			text += indentSpaces + "   action: " + logFilterType + "\n"
			text += indentSpaces + "   condition: " + filter + "\n"
		}
	}

	var firstObject = true
	for i := 0; i < maxLogObjects; i++ {
		logObjectField := "LogObject"
		if i > 0 {
			logObjectField += strconv.Itoa(i)
		}

		if logObjectFilename := paramData[logObjectField+".Filename"]; logObjectFilename != "" {
			logObjectType := paramData[logObjectField+".Type"]
			if logObjectType == "" {
				logObjectType = "ascii"
			}
			logObjectFormat := paramData[logObjectField+".Format"]
			logObjectRollingEnabled := paramData[logObjectField+".RollingEnabled"]
			logObjectRollingIntervalSec := paramData[logObjectField+".RollingIntervalSec"]
			logObjectRollingOffsetHr := paramData[logObjectField+".RollingOffsetHr"]
			logObjectRollingSizeMb := paramData[logObjectField+".RollingSizeMb"]
			logObjectFilters := paramData[logObjectField+".Filters"]

			if firstObject {
				text += "\n" + indentSpaces + "logs:\n"
				firstObject = false
			}
			text += indentSpaces + " - mode: " + logObjectType + "\n"
			text += indentSpaces + "   filename: " + logObjectFilename + "\n"
			text += indentSpaces + "   format: " + logObjectFormat + "\n"

			if logObjectType != "pipe" {
				if logObjectRollingEnabled != "" {
					text += indentSpaces + "   rolling_enabled: " + logObjectRollingEnabled + "\n"
				}
				if logObjectRollingIntervalSec != "" {
					text += indentSpaces + "   rolling_interval_sec: " + logObjectRollingIntervalSec + "\n"
				}
				if logObjectRollingOffsetHr != "" {
					text += indentSpaces + "   rolling_offset_hr: " + logObjectRollingOffsetHr + "\n"
				}
				if logObjectRollingSizeMb != "" {
					text += indentSpaces + "   rolling_size_mb: " + logObjectRollingSizeMb + "\n"
				}
			}
			if logObjectFilters != "" {
				logObjectFilters = strings.Replace(logObjectFilters, "\v", "", -1)
				text += indentSpaces + "   filters: [" + logObjectFilters + "]\n"
			}
		}
	}

	return Cfg{
		Text:        text,
		ContentType: ContentTypeLoggingDotYAML,
		LineComment: LineCommentLoggingDotYAML,
		Warnings:    warnings,
	}, nil
}
