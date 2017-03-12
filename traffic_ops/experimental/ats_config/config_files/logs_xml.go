// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config_files

import (
	"fmt"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"strings"
	"time"
)

func createLogsXmlDotConfig(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	// <!-- Generated for my-edge-0 by Twelve Monkeys (https://to.example.net/) - Do not edit!! -->
	// 	<LogFormat>
	// 		<Name = "custom_ats_2"/>
	// 		<Format = "%<cqtq> chi=%<chi> phn=%<phn> php=%<php> shn=%<shn> url=%<cquuc> cqhm=%<cqhm> cqhv=%<cqhv> pssc=%<pssc> ttms=%<ttms> b=%<pscl> sssc=%<sssc>"/>
	// 		</LogFormat>
	// 		<LogObject>
	// 		<Format = "custom_ats_2"/>
	// 		<Filename = "custom_ats_2"/>
	// 		<RollingEnabled = 2/>
	// 		<RollingIntervalSec = 3600/>
	// 		<RollingOffsetHr = 4/>
	// 		<RollingSizeMb = 128/>
	// 		</LogObject>

	paramMap := createParamsMap(params)

	if _, ok := paramMap["logs_xml.config"]; !ok {
		return "", fmt.Errorf("No logs_xml config parameters")
	}

	configParams := paramMap["logs_xml.config"]

	return fmt.Sprintf(`<!-- Generated for %s by Traffic Ops (%s) on %s - Do not edit!! -->
<LogFormat>
  <Name = "%s"/>
  <Format = "%s"/>
</LogFormat>
<LogObject>
  <Format = "%s"/>
  <Filename = "%s"/>
  <RollingEnabled = %s/>
  <RollingIntervalSec = %s/>
  <RollingOffsetHr = %s/>
  <RollingSizeMb = %s/>
</LogObject>
`,
		trafficServerHost,
		trafficOpsHost,
		time.Now().String(),
		configParams["LogFormat.Name"],
		strings.Replace(configParams["LogFormat.Format"], `"`, `\"`, -1),
		configParams["LogObject.Format"],
		configParams["LogObject.Filename"],
		configParams["LogObject.RollingEnabled"],
		configParams["LogObject.RollingIntervalSec"],
		configParams["LogObject.RollingOffsetHr"],
		configParams["LogObject.RollingSizeMb"]), nil
}
