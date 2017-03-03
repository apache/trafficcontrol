package main

import (
	"flag"
	"fmt"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/nagios"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/tmcheck"
)

const UserAgent = "tm-peerpoller-validator/0.1"

func main() {
	tmURI := flag.String("tm", "", "The Traffic Monitor URI, whose Peer Poller to validate")
	// toUser := flag.String("touser", "", "The Traffic Ops user")
	// toPass := flag.String("topass", "", "The Traffic Ops password")
	// includeOffline := flag.Bool("includeOffline", false, "Whether to include Offline Monitors")
	help := flag.Bool("help", false, "Usage info")
	helpBrief := flag.Bool("h", false, "Usage info")
	flag.Parse()
	if *help || *helpBrief {
		fmt.Printf("Usage: ./nagios-validate-peerpoller -to https://traffic-ops.example.net -touser bill -topass thelizard -includeOffline true\n")
		return
	}

	// toClient, err := to.LoginWithAgent(*toURI, *toUser, *toPass, true, UserAgent, false, tmcheck.RequestTimeout)
	// if err != nil {
	// 	fmt.Printf("Error logging in to Traffic Ops: %v\n", err)
	// 	return
	// }

	err := tmcheck.ValidatePeerPoller(*tmURI)
	if err != nil {
		nagios.Exit(nagios.Critical, fmt.Sprintf("Error validating monitor peer poller: %v", err))
	}
	nagios.Exit(nagios.Ok, "")
}
