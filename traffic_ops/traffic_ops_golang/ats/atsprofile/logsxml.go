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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
)

const LogsXMLFileName = "logs_xml.config"

const MaxLogObjects = 10

func GetLogsXML(w http.ResponseWriter, r *http.Request) {
	addHdr := false
	WithProfileDataHdr(w, r, addHdr, tc.ContentTypeTextPlain, makeLogsXML) // TODO change to Content-Type text/xml? Perl uses text/plain.
}

func makeLogsXML(tx *sql.Tx, _ *config.Config, profile ats.ProfileData, _ string) (string, error) {
	profileParamData, err := ats.GetProfileParamData(tx, profile.ID, LogsXMLFileName)
	if err != nil {
		return "", errors.New("getting profile param data: " + err.Error())
	}

	hdrComment, err := ats.HeaderComment(tx, profile.Name)
	if err != nil {
		return "", errors.New("getting header comment: " + err.Error())
	}
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

		logFormatName := profileParamData[logFormatField+".Name"]
		if logFormatName != "" {
			format := profileParamData[logFormatField+".Format"]
			format = strings.Replace(format, `"`, `\"`, -1)

			text += `<LogFormat>
  <Name = "` + logFormatName + `"/>
  <Format = "` + format + `"/>
</LogFormat>
`
		}

		logObjectFileName := profileParamData[logObjectField+".Filename"]
		if logObjectFileName != "" {
			logObjectFormat := profileParamData[logObjectField+".Format"]
			logObjectRollingEnabled := profileParamData[logObjectField+".RollingEnabled"]
			logObjectRollingIntervalSec := profileParamData[logObjectField+".RollingIntervalSec"]
			logObjectRollingOffsetHr := profileParamData[logObjectField+".RollingOffsetHr"]
			logObjectRollingSizeMb := profileParamData[logObjectField+".RollingSizeMb"]
			logObjectHeader := profileParamData[logObjectField+".Header"]

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
	return text, nil
}
