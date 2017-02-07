package config_files

import (
	"fmt"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"strings"
	"time"
)

func createSetDscpDotConfig(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

	filenamePrefix := "set_dscp_"
	filenameSuffix := ".config"
	dscpDecimal := strings.TrimSuffix(strings.TrimPrefix(filename, filenamePrefix), filenameSuffix)

	s += fmt.Sprintf(`cond %%{REMAP_PSEUDO_HOOK}
set-conn-dscp %s [L]
`, dscpDecimal)
	return s, nil
}
