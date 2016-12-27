package config_files

import (
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"time"
)

func createBgFetchDotConfig(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	// This is a temporary workaround until we have real partial object caching support in ATS, so hardcoding for now
	s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"
	s += "include User-Agent *\n"
	return s, nil
}
