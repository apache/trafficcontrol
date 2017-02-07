package config_files

import (
	"fmt"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"time"
)

func createDropQstringDotConfig(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

	paramMap := createParamsMap(params)
	dropQstring := getParamDefault(paramMap, "drop_qstring.config", "content", "")

	if dropQstring != "" {
		s += fmt.Sprintf("%s\n", dropQstring)
	} else {
		s += "/([^?]+) $s://$t/$1\n"
	}

	return s, nil
}
