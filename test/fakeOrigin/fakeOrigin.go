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
	"os"

	"github.com/apache/trafficcontrol/v8/test/fakeOrigin/endpoint"
	"github.com/apache/trafficcontrol/v8/test/fakeOrigin/httpService"
	"github.com/apache/trafficcontrol/v8/test/fakeOrigin/transcode"
)

func printUsage() {
	fmt.Print(`Usage:
	fakeOrigin -cfg config.json
`)
}

func main() {
	cfg := flag.String("cfg", endpoint.DefaultConfigFile, "config file location")
	printVersion := flag.Bool("version", false, "print fakeOrigin version")
	flag.Parse()

	if *printVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	config := endpoint.Config{}
	err := error(nil)

	fmt.Println("using config: " + *cfg)
	if config, err = endpoint.LoadAndGenerateDefaultConfig(*cfg); err != nil {
		fmt.Printf("An error occurred while loading configuration '%v': %s\n", *cfg, err)
		os.Exit(1)
	}

	for i, ep := range config.Endpoints {
		if ep.EndpointType == endpoint.Static || ep.EndpointType == endpoint.Dir {
			continue
		}
		var cmd string
		var args []string
		cmd, args, err = endpoint.GetTranscoderCommand(config.Endpoints[i])
		if err != nil {
			fmt.Printf("An error occurred while fetching transcoder commands: %s\n", err)
			os.Exit(1)
		}
		if cmd == "" {
			fmt.Println("Skipping Transcode for endpoint: " + config.Endpoints[i].ID)
		} else if err = transcode.Do(&config.Endpoints[i], cmd, args); err != nil {
			fmt.Printf("An error occurred while performing transcoder commands: %s\n", err)
			os.Exit(1)
		}
	}

	routes, err := httpService.GetRoutes(config)
	if err != nil {
		fmt.Println("Error getting routes: " + err.Error())
		os.Exit(1)
	}
	httpService.PrintRoutes(os.Stdout, routes, "", "Serving ", false)

	go func() {
		if err := httpService.StartHTTPSListener(config, routes); err != nil {
			fmt.Printf("Error serving HTTPS: %s\n", err)
			os.Exit(1)
		}
	}()

	if err := httpService.StartHTTPListener(config, routes); err != nil {
		fmt.Printf("Error serving HTTP: %s\n", err)
		os.Exit(1)
	}
}
