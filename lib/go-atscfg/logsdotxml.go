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
)

const LogsXMLFileName = "logs_xml.config"

func MakeLogsXMLDotConfig(
	profileName string,
	paramData map[string]string, // GetProfileParamData(tx, profile.ID, LoggingYAMLFileName)
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
) string {
	hdrComment := GenericHeaderComment(profileName, toToolName, toURL)
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
	return text
}
