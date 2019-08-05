package atsprofile

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
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
)

const LoggingFileName = "logging.config"

func GetLogging(w http.ResponseWriter, r *http.Request) {
	addHdr := false
	WithProfileDataHdr(w, r, addHdr, tc.ContentTypeTextPlain, makeLogging) // TODO change to Content-Type text/x-lua? Perl uses text/plain.
}

func makeLogging(tx *sql.Tx, _ *config.Config, profile ats.ProfileData, _ string) (string, error) {
	profileParamData, err := ats.GetProfileParamData(tx, profile.ID, LoggingFileName)

	if err != nil {
		return "", errors.New("getting profile param data: " + err.Error())
	}

	hdrComment, err := ats.HeaderComment(tx, profile.Name)
	if err != nil {
		return "", errors.New("getting header comment: " + err.Error())
	}
	// This is an LUA file, so we need to massage the header a bit for LUA commenting.
	hdrComment = strings.Replace(hdrComment, `# `, ``, -1)
	hdrComment = strings.Replace(hdrComment, "\n", ``, -1)
	text := "-- " + hdrComment + " --\n"

	for i := 0; i < MaxLogObjects; i++ {
		logFormatField := "LogFormat"
		if i > 0 {
			logFormatField += strconv.Itoa(i)
		}
		if logFormatName := profileParamData[logFormatField+".Name"]; logFormatName != "" {
			format := profileParamData[logFormatField+".Format"]
			if format == "" {
				// TODO determine if the line should be excluded. Perl includes it anyway, without checking.
				log.Errorf("Profile '%v' has logging.config format '%v' Name Parameter but no Format Parameter. Setting blank Format!\n", profile.Name, logFormatField)
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

		if logFilterName := profileParamData[logFilterField+".Name"]; logFilterName != "" {
			filter := profileParamData[logFilterField+".Filter"]
			if filter == "" {
				// TODO determine if the line should be excluded. Perl includes it anyway, without checking.
				log.Errorf("Profile '%v' has logging.config format '%v' Name Parameter but no Filter Parameter. Setting blank Filter!\n", profile.Name, logFilterField)
			}

			filter = strings.Replace(filter, `\`, `\\`, -1)
			filter = strings.Replace(filter, `'`, `\'`, -1)

			logFilterType := profileParamData[logFilterField+".Type"]
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

		if logObjectFilename := profileParamData[logObjectField+".Filename"]; logObjectFilename != "" {
			logObjectType := profileParamData[logObjectField+".Type"]
			if logObjectType == "" {
				logObjectType = "ascii"
			}
			logObjectFormat := profileParamData[logObjectField+".Format"]
			logObjectRollingEnabled := profileParamData[logObjectField+".RollingEnabled"]
			logObjectRollingIntervalSec := profileParamData[logObjectField+".RollingIntervalSec"]
			logObjectRollingOffsetHr := profileParamData[logObjectField+".RollingOffsetHr"]
			logObjectRollingSizeMb := profileParamData[logObjectField+".RollingSizeMb"]
			logObjectFilters := profileParamData[logObjectField+".Filters"]

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

	return text, nil
}
