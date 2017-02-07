package config_files

import (
	"fmt"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"strings"
	"time"
)

func createRegexRemapDotConfig(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

	filenamePrefix := "regex_remap_"
	filenameSuffix := ".config"
	dsXmlId := strings.TrimSuffix(strings.TrimPrefix(filename, filenamePrefix), filenameSuffix)

	deliveryServices, err := toClient.DeliveryServices()
	if err != nil {
		return "", fmt.Errorf("error getting delivery services: %v", err)
	}

	ds, err := getDeliveryService(deliveryServices, dsXmlId)
	if err != nil {
		return "", fmt.Errorf("error getting delivery service '%v': %v", dsXmlId, err)
	}

	s += fmt.Sprintf("%s\n", ds.RegexRemap)

	s = strings.Replace(s, "__RETURN__", "\n", -1)
	return s, nil
}
