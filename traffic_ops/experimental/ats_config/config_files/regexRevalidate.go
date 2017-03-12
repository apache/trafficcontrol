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
)

func createRegexRevalidateDotConfig(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	// TODO add Jobs endpoint to TO, implement
	return "", fmt.Errorf("cannot create regex_revalidate.config - Traffic Ops has not Jobs API")
	// s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

	// paramMap := createParamsMap(params)

	// server, err := getServer(toClient, trafficServerHost)
	// if err != nil {
	// 	return "", fmt.Errorf("getting server %s: %v", trafficServerHost, err)
	// }

	// maxDays := getParamDefault("regex_revalidate.config", "maxRevalDurationDays", "")
	// interval := fmt.Sprintf(`> now() - interval '%s day'`, maxDays) // postgres

	// return s, nil
}
