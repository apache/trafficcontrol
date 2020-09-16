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

	"github.com/apache/trafficcontrol/lib/go-log"
)

const LoggingYAMLFileName = "logging.yaml"
const ContentTypeLoggingDotYAML = "application/yaml; charset=us-ascii" // Note YAML has no IANA standard mime type. This is one of several common usages, and is likely to be the standardized value. If you're reading this, please check IANA to see if YAML has been added, and change this to the IANA definition if so. Also note we include 'charset=us-ascii' because YAML is commonly UTF-8, but ATS is likely to be unable to handle UTF.
const LineCommentLoggingDotYAML = LineCommentHash

func MakeLoggingDotYAML(
	profileName string,
	paramData map[string]string, // GetProfileParamData(tx, profile.ID, LoggingYAMLFileName)
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
) string {
	hdr := GenericHeaderComment(profileName, toToolName, toURL)

	// note we use the same const as logs.xml - this isn't necessarily a requirement, and we may want to make separate variables in the future.
	maxLogObjects := MaxLogObjects

	text := hdr
	text += "\nformats: \n"
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
				log.Errorf("Profile '%v' has logging.yaml format '%v' Name Parameter but no Format Parameter. Setting blank Format!\n", profileName, logFormatField)
			}
			text += " - name: " + logFormatName + " \n"
			text += "   format: '" + format + "'\n"
		}
	}

	text += "filters:\n"
	for i := 0; i < maxLogObjects; i++ {
		logFilterField := "LogFilter"
		if i > 0 {
			logFilterField += strconv.Itoa(i)
		}
		if logFilterName := paramData[logFilterField+".Name"]; logFilterName != "" {
			filter := paramData[logFilterField+".Filter"]
			if filter == "" {
				// TODO determine if the line should be excluded. Perl includes it anyway, without checking.
				log.Errorf("Profile '%v' has logging.yaml filter '%v' Name Parameter but no Filter Parameter. Setting blank Filter!\n", profileName, logFilterField)
			}
			logFilterType := paramData[logFilterField+".Type"]
			if logFilterType == "" {
				logFilterType = "accept"
			}
			text += "- name: " + logFilterName + "\n"
			text += "  action: " + logFilterType + "\n"
			text += "  condition: " + filter + "\n"
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
				text += "\nlogs:\n"
				firstObject = false
			}
			text += "- mode: " + logObjectType + "\n"
			text += "  filename: " + logObjectFilename + "\n"
			text += "  format: " + logObjectFormat + "\n"

			if logObjectType != "pipe" {
				if logObjectRollingEnabled != "" {
					text += "  rolling_enabled: " + logObjectRollingEnabled + "\n"
				}
				if logObjectRollingIntervalSec != "" {
					text += "  rolling_interval_sec: " + logObjectRollingIntervalSec + "\n"
				}
				if logObjectRollingOffsetHr != "" {
					text += "  rolling_offset_hr: " + logObjectRollingOffsetHr + "\n"
				}
				if logObjectRollingSizeMb != "" {
					text += "  rolling_size_mb: " + logObjectRollingSizeMb + "\n"
				}
			}
			if logObjectFilters != "" {
				logObjectFilters = strings.Replace(logObjectFilters, "\v", "", -1)
				text += "  filters: [" + logObjectFilters + "]\n"
			}
		}
	}

	return text
}
