package config_files

import (
	"fmt"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"time"
)

func createVolumeDotConfig(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	// # DO NOT EDIT - Generated for my-edge-0 by Twelve Monkeys (https://tm.example.net/) on Thu Dec 22 23:33:13 UTC 2016
	// 	# 12M NOTE: This is running with forced volumes - the size is irrelevant
	// 	volume=1 scheme=http size=50%
	// 	volume=2 scheme=http size=50%

	s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"
	s += "# 12M NOTE: This is running with forced volumes - the size is irrelevant\n"

	paramMap := createParamsMap(params)

	if _, ok := paramMap["storage.config"]; !ok {
		return "", fmt.Errorf("No storage config parameters")
	}

	storageConfigParams := paramMap["storage.config"]

	volumePrefixes := []string{"", "RAM_", "SSD_"}

	numVolumes := 0
	for _, prefix := range volumePrefixes {
		if _, ok := storageConfigParams[prefix+"Drive_Prefix"]; ok {
			numVolumes++
		}
	}

	volumeText := func(volumeNum int, numVolumes int) string {
		return fmt.Sprintf("volume=%d scheme=http size=%d%%\n", volumeNum, 100/numVolumes)
	}

	nextVolumeNum := 1
	for _, prefix := range volumePrefixes {
		if _, hasDrivePrefix := storageConfigParams[prefix+"Drive_Prefix"]; hasDrivePrefix {
			s += volumeText(nextVolumeNum, numVolumes)
			nextVolumeNum++
		}
	}
	return s, nil
}
