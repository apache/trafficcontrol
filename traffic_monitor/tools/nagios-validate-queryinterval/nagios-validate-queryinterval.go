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

package main

import (
	"flag"
	"fmt"
	"github.com/apache/trafficcontrol/v8/lib/go-nagios"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/tmcheck"
	to "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

const UserAgent = "tm-queryinterval-validator/0.1"

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

	toClient, _, err := to.LoginWithAgent(*toURI, *toUser, *toPass, true, UserAgent, false, tmcheck.RequestTimeout)
	if err != nil {
		fmt.Printf("Error logging in to Traffic Ops: %v\n", err)
		return
	}

	monitorErrs, err := tmcheck.ValidateAllMonitorsQueryInterval(toClient, *includeOffline)

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
