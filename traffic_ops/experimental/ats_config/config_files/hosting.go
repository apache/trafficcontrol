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

// getServerTypeStr returns the type string from Traffic Ops for the given Traffic Server (e.g. EDGE, MID)
func getServerTypeStr(toClient towrap.ITrafficOpsSession, serverToFind string) (string, error) {
	// TODO add TO endpoint to get a single server's data, for efficiency.
	servers, err := toClient.Servers()
	if err != nil {
		return "", fmt.Errorf("error getting servers from Traffic Ops: %v", err)
	}
	for _, server := range servers {
		if server.HostName == serverToFind {
			return server.Type, nil
		}
	}
	return "", fmt.Errorf("not found")
}

func createHostingDotConfig(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {

	// TODO put in common func, remove duplicates
	s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

	paramMap := createParamsMap(params)
	if _, ok := paramMap["storage.config"]; !ok {
		return "", fmt.Errorf("No storage config parameters")
	}
	storageConfigParams := paramMap["storage.config"]

	if _, hasRamDrivePrefix := storageConfigParams["RAM_Drive_Prefix"]; !hasRamDrivePrefix {
		diskVolumeNum := 1 // TODO verify correct (mirrors old TO func)
		s += fmt.Sprintf("hostname=*   volume=%d\n", diskVolumeNum)
		return s, nil
	}

	nextVolumeNum := 1
	if _, hasDrivePrefix := storageConfigParams["Drive_Prefix"]; hasDrivePrefix {
		diskVolumeNum := nextVolumeNum
		nextVolumeNum++
		s += fmt.Sprintf("# 12M NOTE: volume %v is the Disk volume\n", diskVolumeNum)
	}

	ramVolumeNum := nextVolumeNum
	nextVolumeNum++
	s += fmt.Sprintf("# 12M NOTE: volume %v is the RAM volume\n", ramVolumeNum)

	deliveryServices, err := toClient.DeliveryServices()
	if err != nil {
		return "", fmt.Errorf("error getting delivery services: %v", err)
	}

	serverTypeStr, err := getServerTypeStr(toClient, trafficServerHost)
	if err != nil {
		return "", fmt.Errorf("error getting server '%v' type: %v", trafficServerHost, err)
	}

	listed := map[string]struct{}{}
	for _, deliveryService := range deliveryServices {
		if _, alreadyListed := listed[deliveryService.XMLID]; alreadyListed {
			continue
		}
		isEdge := strings.HasPrefix(serverTypeStr, "EDGE")
		isMid := strings.HasPrefix(serverTypeStr, "MID")
		isLive := strings.HasSuffix(deliveryService.Type, "_LIVE")
		isLiveNatl := strings.HasSuffix(deliveryService.Type, "_LIVE_NATNL")
		if (isEdge && (isLive || isLiveNatl)) || (isMid && isLiveNatl) {
			orgFqdn := stripProtocol(deliveryService.OrgServerFQDN)
			s += fmt.Sprintf("hostname=%s volume=%d\n", orgFqdn, ramVolumeNum)
			listed[deliveryService.XMLID] = struct{}{}
		}
	}
	diskVolume := 1 // TODO verify correct (mirrors old TO func)
	s += fmt.Sprintf("hostname=*   volume=%d\n", diskVolume)
	return s, nil
}
