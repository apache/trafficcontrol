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
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"regexp"
	"time"
)

// my $separator ||= {
// 	"records.config"  => " ",
// 	"plugin.config"   => " ",
// 	"sysctl.conf"     => " = ",
// 	"url_sig_.config" => " = ",
// 	"astats.config"   => "=",
// };

func createGenericDotConfigFunc(separator string) ConfigFileCreatorFunc {
	underscoreDigitSuffixRegex := regexp.MustCompile("__[0-9]+$")

	return func(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
		s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

		paramMap := createParamsMap(params)
		fileParams := paramMap[filename]
		for name, val := range fileParams {
			name := underscoreDigitSuffixRegex.ReplaceAllString(name, "")
			s += name + separator + val + "\n"
		}
		return s, nil
	}
}
