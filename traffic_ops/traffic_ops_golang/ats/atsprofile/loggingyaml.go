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
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
)

const LoggingYAMLFileName = "logging.yaml"

func GetLoggingYAML(w http.ResponseWriter, r *http.Request) {
	WithProfileData(w, r, makeLoggingYAML) // TODO change to Content-Type text/yaml? Perl uses text/plain.
}

func makeLoggingYAML(tx *sql.Tx, _ *config.Config, profile ats.ProfileData, _ string) (string, error) {
	profileParamData, err := ats.GetProfileParamData(tx, profile.ID, LoggingYAMLFileName)
	if err != nil {
		return "", errors.New("getting profile param data: " + err.Error())
	}

	// note we use the same const as logs.xml - this isn't necessarily a requirement, and we may want to make separate variables in the future.
	maxLogObjects := MaxLogObjects

	text := "\nformats: \n"
	for i := 0; i < maxLogObjects; i++ {
		logFormatField := "LogFormat"
		if i > 0 {
			logFormatField += strconv.Itoa(i)
		}
		logFormatName := profileParamData[logFormatField+".Name"]
		if logFormatName != "" {
			format := profileParamData[logFormatField+".Format"]
			if format == "" {
				// TODO determine if the line should be excluded. Perl includes it anyway, without checking.
				log.Errorf("Profile '%v' has logging.yaml format '%v' Name Parameter but no Format Parameter. Setting blank Format!\n", profile.Name, logFormatField)
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
		if logFilterName := profileParamData[logFilterField+".Name"]; logFilterName != "" {
			filter := profileParamData[logFilterField+".Filter"]
			if filter == "" {
				// TODO determine if the line should be excluded. Perl includes it anyway, without checking.
				log.Errorf("Profile '%v' has logging.yaml filter '%v' Name Parameter but no Filter Parameter. Setting blank Filter!\n", profile.Name, logFilterField)
			}
			logFilterType := profileParamData[logFilterField+".Type"]
			if logFilterType == "" {
				logFilterType = "accept"
			}
			text += "- name: " + logFilterName + "\n"
			text += "  action: " + logFilterType + "\n"
			text += "  condition: " + filter + "\n"
		}
	}

	for i := 0; i < maxLogObjects; i++ {
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

			text += "\nlogs:\n"
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
				text += "  filters: [" + logObjectFilters + "]"
			}
		}
	}
	return text, nil
}
