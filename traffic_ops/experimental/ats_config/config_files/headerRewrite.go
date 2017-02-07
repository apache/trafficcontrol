package config_files

import (
	"fmt"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"strings"
	"time"
)

// TODO combine duplicate logic in createMidHeaderRewriteDotConfig
func createHeaderRewriteDotConfig(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	filenamePrefix := "hdr_rw_"
	filenameSuffix := ".config"
	dsXmlId := strings.TrimSuffix(strings.TrimPrefix(filename, filenamePrefix), filenameSuffix)

	// TODO change config_files dispatch to be ordered, and directly dispatch to createHeaderRewriteMidDotConfig
	if strings.HasPrefix(dsXmlId, "mid_") {
		return createHeaderRewriteMidDotConfig(toClient, filename, trafficOpsHost, trafficServerHost, params)
	}

	s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

	server, err := getServer(toClient, trafficServerHost)
	if err != nil {
		return "", fmt.Errorf("getting server %s: %v", trafficServerHost, err)
	}

	deliveryServices, err := toClient.DeliveryServices()
	if err != nil {
		return "", fmt.Errorf("error getting delivery services: %v", err)
	}

	ds, err := getDeliveryService(deliveryServices, dsXmlId)
	if err != nil {
		return "", fmt.Errorf("error getting delivery service '%v': %v", dsXmlId, err)
	}

	actions := ds.EdgeHeaderRewrite
	s += fmt.Sprintf("%s\n", actions)

	s = strings.Replace(s, "__RETURN__", "\n", -1)
	s = strings.Replace(s, "__CACHE_IPV4__", server.IPAddress, -1)
	return s, nil
}

func createHeaderRewriteMidDotConfig(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	filenamePrefix := "hdr_rw_mid_"
	filenameSuffix := ".config"
	dsXmlId := strings.TrimSuffix(strings.TrimPrefix(filename, filenamePrefix), filenameSuffix)

	s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

	server, err := getServer(toClient, trafficServerHost)
	if err != nil {
		return "", fmt.Errorf("getting server %s: %v", trafficServerHost, err)
	}

	deliveryServices, err := toClient.DeliveryServices()
	if err != nil {
		return "", fmt.Errorf("error getting delivery services: %v", err)
	}

	ds, err := getDeliveryService(deliveryServices, dsXmlId)
	if err != nil {
		return "", fmt.Errorf("error getting delivery service '%v': %v", dsXmlId, err)
	}
	actions := ds.MidHeaderRewrite
	s += fmt.Sprintf("%s\n", actions)

	s = strings.Replace(s, "__RETURN__", "\n", -1)
	s = strings.Replace(s, "__CACHE_IPV4__", server.IPAddress, -1)
	return s, nil
}

func getDeliveryService(dses []to.DeliveryService, name string) (to.DeliveryService, error) {
	for _, ds := range dses {
		if ds.XMLID == name {
			return ds, nil
		}
	}
	return to.DeliveryService{}, fmt.Errorf("not found")
}
