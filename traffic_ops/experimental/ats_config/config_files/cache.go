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

func createCacheDotConfig(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

	// TODO get only DSes for this server?
	deliveryServices, err := toClient.DeliveryServices()
	if err != nil {
		return "", fmt.Errorf("error getting delivery services: %v", err)
	}
	for _, deliveryService := range deliveryServices {
		if deliveryService.Type != "HTTP_NO_CACHE" {
			continue
		}
		s += fmt.Sprintf("dest_domain=%s scheme=http action=never-cache\n", stripProtocol(deliveryService.OrgServerFQDN))
	}
	return s, nil
}
