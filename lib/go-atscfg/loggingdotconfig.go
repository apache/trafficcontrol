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

const MaxLogObjects = 10

const LoggingFileName = "logging.config"

// MakeStorageDotConfig creates storage.config for a given ATS Profile.
// The paramData is the map of parameter names to values, for all parameters assigned to the given profile, with the config_file "storage.config".
func MakeLoggingDotConfig(
	profileName string,
	paramData map[string]string, // GetProfileParamData(tx, profile.ID, StorageFileName)
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
) string {

	hdrComment := GenericHeaderComment(profileName, toToolName, toURL)
	// This is an LUA file, so we need to massage the header a bit for LUA commenting.
	hdrComment = strings.Replace(hdrComment, `# `, ``, -1)
	hdrComment = strings.Replace(hdrComment, "\n", ``, -1)
	text := "-- " + hdrComment + " --\n"

	for i := 0; i < MaxLogObjects; i++ {
		logFormatField := "LogFormat"
		if i > 0 {
			logFormatField += strconv.Itoa(i)
		}
		if logFormatName := paramData[logFormatField+".Name"]; logFormatName != "" {
			format := paramData[logFormatField+".Format"]
			if format == "" {
				// TODO determine if the line should be excluded. Perl includes it anyway, without checking.
				log.Errorf("Profile '%v' has logging.config format '%v' Name Parameter but no Format Parameter. Setting blank Format!\n", profileName, logFormatField)
			}
			format = strings.Replace(format, `"`, `\"`, -1)
			text += logFormatName + " = format {\n"
			text += "	Format = '" + format + " '\n"
			text += "}\n"
		}
	}

	for i := 0; i < MaxLogObjects; i++ {
		logFilterField := "LogFilter"
		if i > 0 {
			logFilterField += strconv.Itoa(i)
		}

		if logFilterName := paramData[logFilterField+".Name"]; logFilterName != "" {
			filter := paramData[logFilterField+".Filter"]
			if filter == "" {
				// TODO determine if the line should be excluded. Perl includes it anyway, without checking.
				log.Errorf("Profile '%v' has logging.config format '%v' Name Parameter but no Filter Parameter. Setting blank Filter!\n", profileName, logFilterField)
			}

			filter = strings.Replace(filter, `\`, `\\`, -1)
			filter = strings.Replace(filter, `'`, `\'`, -1)

			logFilterType := paramData[logFilterField+".Type"]
			if logFilterType == "" {
				logFilterType = "accept"
			}
			text += logFilterName + " = filter." + logFilterType + "('" + filter + "')\n"
		}
	}

	for i := 0; i < MaxLogObjects; i++ {
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

			text += "\nlog." + logObjectType + " {\n"
			text += "  Format = " + logObjectFormat + ",\n"
			text += "  Filename = '" + logObjectFilename + "'"
			if logObjectType != "pipe" {
				text += ",\n"
				text += "  RollingEnabled = " + logObjectRollingEnabled + ",\n"
				text += "  RollingIntervalSec = " + logObjectRollingIntervalSec + ",\n"
				text += "  RollingOffsetHr = " + logObjectRollingOffsetHr + ",\n"
				text += "  RollingSizeMb = " + logObjectRollingSizeMb
			}
			if logObjectFilters != "" {
				text += ",\n  Filters = { " + logObjectFilters + " }"
			}
			text += "\n}\n"
		}
	}

	return text
}
