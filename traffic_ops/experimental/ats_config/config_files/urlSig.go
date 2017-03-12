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
	"time"
)

// TODO pass riak creds to all funcs? Change all funcs to use lambdas?
func createUrlSigDotConfigFunc(riakUser string, riakPass string) ConfigFileCreatorFunc {
	return func(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {

		separator := " = "

		s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

		insecure := true // TODO fix riak certs, and change to false
		keys, err := RiakGetURLSigKeys(toClient, riakUser, riakPass, filename, insecure)
		if err != nil {
			return "", fmt.Errorf("error getting keys from Riak: %v", err)
		}

		for keyName, keyVal := range keys {
			s += fmt.Sprintf("%s%s%s\n", keyName, separator, keyVal)
		}

		return s, nil
	}
}
