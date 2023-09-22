package main

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/v8/traffic_monitor/tools/testcaches/fakesrvr"
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
	_, err := fakesrvr.News(*portStart, *numPorts, remaps)
	if err != nil {
		fmt.Println("Error making FakeServers: " + err.Error())
		return
	}
	for {
		// TODO handle sighup to die
		time.Sleep(time.Hour)
	}
}
