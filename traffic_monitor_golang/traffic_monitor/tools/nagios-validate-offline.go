package main

import (
	"flag"
	"fmt"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/nagios"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/tmcheck"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

const UserAgent = "tm-offline-validator/0.1"

func main() {
	toURI := flag.String("to", "", "The Traffic Ops URI, whose CRConfig to validate")
	toUser := flag.String("touser", "", "The Traffic Ops user")
	toPass := flag.String("topass", "", "The Traffic Ops password")
	includeOffline := flag.Bool("includeOffline", false, "Whether to include Offline Monitors")
	help := flag.Bool("help", false, "Usage info")
	helpBrief := flag.Bool("h", false, "Usage info")
	flag.Parse()
	if *help || *helpBrief || *toURI == "" {
		fmt.Printf("Usage: ./nagios-validate-offline -to https://traffic-ops.example.net -touser bill -topass thelizard -includeOffline true\n")
		return
	}

	toClient, err := to.LoginWithAgent(*toURI, *toUser, *toPass, true, UserAgent, false, tmcheck.RequestTimeout)
	if err != nil {
		fmt.Printf("Error logging in to Traffic Ops: %v\n", err)
		return
	}

	monitorErrs, err := tmcheck.ValidateAllMonitorsOfflineStates(toClient, *includeOffline)

	if err != nil {
		nagios.Exit(nagios.Critical, fmt.Sprintf("Error validating monitor offline statuses: %v", err))
	}

	errStr := ""
	for monitor, err := range monitorErrs {
		if err != nil {
			errStr += fmt.Sprintf("error validating offline status for monitor %v : %v\n", monitor, err.Error())
		}
	}

	if errStr != "" {
		nagios.Exit(nagios.Critical, errStr)
	}

	nagios.Exit(nagios.Ok, "")
}
