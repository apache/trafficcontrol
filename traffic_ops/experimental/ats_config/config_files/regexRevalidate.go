package config_files

import (
	"fmt"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/trafficopswrapper"
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
