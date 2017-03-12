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
	"encoding/json"
	"fmt"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"strings"
	"time"
)

func createSslMulticertDotConfig(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

	deliveryServices, err := toClient.DeliveryServices()
	if err != nil {
		return "", fmt.Errorf("error getting delivery services: %v", err)
	}

	server, err := getServer(toClient, trafficServerHost)
	if err != nil {
		return "", fmt.Errorf("error getting server: %v", err)
	}

	serverDSes, err := getServerDeliveryServices(toClient, server)
	if err != nil {
		return "", fmt.Errorf("error getting server delivery services: %v", err)
	}

	crc, err := getCRConfig(toClient, server)
	if err != nil {
		return "", fmt.Errorf("error getting CRConfig: %v", err)
	}

	crcServer, crcServerExists := crc.ContentServers[trafficServerHost]
	if !crcServerExists {
		return "", fmt.Errorf("error CRConfig server self '%v' doesn't exist", trafficServerHost)
	}

	for _, ds := range deliveryServices {
		if _, ok := serverDSes[ds.XMLID]; !ok {
			continue
		}
		if ds.Protocol <= 0 {
			continue
		}

		crcDsFqdns, crcDsExists := crcServer.DeliveryServices[ds.XMLID]
		if !crcDsExists {
			return "", fmt.Errorf("delivery service %s not found in CRConfig for server self %s", ds.XMLID, trafficServerHost)
		}
		if len(crcDsFqdns) == 0 {
			return "", fmt.Errorf("delivery service %s has no FQDNs in CRConfig for server self %s", ds.XMLID, trafficServerHost)
		}
		crcDsFqdn := crcDsFqdns[0]

		crcDsFqdn = strings.Replace(crcDsFqdn, trafficServerHost, "ccr", -1) // TODO fix this terrible hack, and get the real FQDN for DNS DSes

		keyName := fmt.Sprintf("%s.key", crcDsFqdn)

		uscoreFqdn := strings.Replace(crcDsFqdn, ".", "_", -1)
		certName := fmt.Sprintf("%s_cert.cer", uscoreFqdn)

		s += fmt.Sprintf("ssl_cert_name=%s\t ssl_key_name=%s\n", certName, keyName)
	}

	return s, nil
}

func getCRConfig(toClient towrap.ITrafficOpsSession, server to.Server) (CRConfig, error) {
	crConfigBytes, err := toClient.CRConfigRaw(server.CDNName)
	if err != nil {
		return CRConfig{}, fmt.Errorf("fetching: %v", err)
	}

	crConfig := CRConfig{}
	if err := json.Unmarshal(crConfigBytes, &crConfig); err != nil {
		return CRConfig{}, fmt.Errorf("unmarshalling: %v", err)
	}
	return crConfig, nil
}
