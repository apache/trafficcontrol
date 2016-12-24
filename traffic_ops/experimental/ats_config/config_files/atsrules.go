package config_files

import (
	"fmt"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"strings"
	"time"
)

func createAtsDotRules(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

	paramMap := createParamsMap(params)

	if _, ok := paramMap["storage.config"]; !ok {
		return "", fmt.Errorf("No storage config parameters")
	}

	storageConfigParams := paramMap["storage.config"]

	drivePrefix := storageConfigParams["Drive_Prefix"] // TODO handle nonexistent param
	drivePostfix := strings.Split(storageConfigParams["Drive_Letters"], ",")
	for _, letter := range drivePostfix {
		drivePrefix := strings.Replace(drivePrefix, "/dev/", "", -1) // TODO verify; put outside loop?
		s += fmt.Sprintf(`KERNEL=="%s%s" OWNER="ats"
`, drivePrefix, letter)
	}

	if drivePrefix, ok := storageConfigParams["RAM_Drive_Prefix"]; ok {
		drivePostfix := strings.Split(storageConfigParams["RAM_Drive_Letters"], ",")
		for _, letter := range drivePostfix {
			drivePrefix := strings.Replace(drivePrefix, "/dev/", "", -1) // TODO verify; put outside loop?
			s += fmt.Sprintf(`KERNEL=="%s%s" OWNER="ats"
`, drivePrefix, letter)
		}
	}

	return s, nil
}
