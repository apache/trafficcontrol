package config_files

import (
	"fmt"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"time"
)

func createFacts(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

	server, err := getServer(toClient, trafficServerHost)
	if err != nil {
		return "", fmt.Errorf("getting server %s: %v", trafficServerHost, err)
	}

	s += fmt.Sprintf("profile:%s\n", server.Profile)
	return s, nil
}
