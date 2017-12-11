package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/tools/testcaches/fakesrvr"
)

func makeFakeRemaps(n int) []string {
	remaps := []string{}
	for i := 0; i < n; i++ {
		remaps = append(remaps, "num"+strconv.Itoa(i)+".example.net")
	}
	return remaps
}

func main() {
	portStart := flag.Int("portStart", 40000, "Starting port in range")
	numPorts := flag.Int("numPorts", 1000, "Number of ports to serve")
	numRemaps := flag.Int("numRemaps", 1000, "Number of remaps to serve")
	flag.Parse()
	if *portStart < 0 || *portStart > 65535 {
		fmt.Println("portStart must be 0-65535")
		return
	} else if *numPorts < 0 || *portStart+*numPorts > 65535 {
		fmt.Println("numPorts must be > 0 and portStart+numPorts < 65535")
		return
	} else if *numRemaps < 0 {
		fmt.Println("numRemaps must be > 0")
		return
	}

	remaps := makeFakeRemaps(*numRemaps)
	servers, err := fakesrvr.News(*portStart, *numPorts, remaps)
	if err != nil {
		fmt.Println("Error making FakeServers: " + err.Error())
		return
	}
	servers = servers // debug
	for {
		// TODO handle sighup to die
		time.Sleep(time.Hour)
	}
}
