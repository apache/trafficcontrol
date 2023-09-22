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
 *
 */

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/availableservers"
	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/config"
	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/coveragezone"
	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/crconfigpoller"
	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/crstatespoller"
	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/fetch"
	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/httpsrvr"
	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/toutil"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	to "github.com/apache/trafficcontrol/v8/traffic_ops/v3-client"
)

// Version is the utility's version.
const Version = "0.1"

// UserAgent is the string passed by the utility in the HTTP User-Agent header.
const UserAgent = "traffic_router_golang/" + Version

// DefaultConfigFile is the relative file path used as the default
// configuration file location
const DefaultConfigFile = "./cfg.json"

// CRConfigPath is the Traffic Monitor API path used to retrieve the CDN
// Snapshot.
const CRConfigPath = "/publish/CrConfig"

// CRStatesPath is the Traffic Monitor API path used to retrieve the states of
// the monitored cache servers.
const CRStatesPath = "/publish/CrStates"

func main() {
	cfgFile := flag.String("cfg", DefaultConfigFile, "The config file path")
	flag.Parse()

	cfg, err := config.Load(*cfgFile)
	if err != nil {
		fmt.Println("Error loading config file '" + *cfgFile + "': " + err.Error())
		os.Exit(1)
	}
	if err := log.InitCfg(cfg); err != nil {
		fmt.Println("Error creating log writers: " + err.Error())
		os.Exit(1)
	}

	log.Infof("Starting with config %+v\n", cfg)

	cz, err := coveragezone.Load(cfg.CoverageZoneFile)
	if err != nil {
		fmt.Println("Error loading coverage zone file '" + cfg.CoverageZoneFile + "': " + err.Error())
		os.Exit(1)
	}

	toURIStr := (*url.URL)(cfg.TrafficOpsURI).String()
	log.Infof("TO URI Str: " + toURIStr + "\n")

	toClient, toAddr, err := to.LoginWithAgent(toURIStr, cfg.TrafficOpsUser, cfg.TrafficOpsPass, cfg.TrafficOpsInsecure, UserAgent, cfg.TrafficOpsClientCache, time.Duration(cfg.TrafficOpsTimeout))
	if err != nil {
		log.Errorf("logging in to Traffic Ops (%v): %v\n", toAddr, err)
		return
	}
	log.Infoln("Connected to Traffic Ops " + toURIStr + " (" + toAddr.String() + ")")

	monitors, err := toutil.GetMonitorURIs(toClient, cfg.CDN) // TODO re-fetch monitors on interval
	if err != nil {
		log.Errorf("getting monitors from Traffic Ops (%v): %v\n", toAddr, err)
		return
	}
	log.Infof("Got Traffic Ops Monitors: %+v\n", monitors)

	availableservers.Test() // debug

	crconfigFetcher := fetch.NewHTTPRoundRobin(monitors, CRConfigPath, time.Duration(cfg.ReqTimeout), UserAgent)
	crstatesFetcher := fetch.NewHTTPRoundRobin(monitors, CRStatesPath, time.Duration(cfg.ReqTimeout), UserAgent)

	// debug
	// crconfigFetcher := fetch.NewFile("./crconfig.json")
	// crstatesFetcher := fetch.NewFile("./crstates.json")

	thsCRConfig, thsCRConfigRegexes, thsCGSearcher, thsNextCacher, err := crconfigpoller.Start(crconfigFetcher, time.Duration(cfg.CRConfigInterval))
	if err != nil {
		fmt.Println("Could not get initial CRConfig: ", err)
	}

	thsCRStates, availableServers, err := crstatespoller.Start(crstatesFetcher, time.Duration(cfg.CRStatesInterval), thsCRConfig)
	if err != nil {
		fmt.Println("Could not get initial CRStates from: ", err)
	}

	httpsrvr.Start(thsCRConfigRegexes, availableServers, thsCGSearcher, thsNextCacher, cz, cfg.Port)

	// debug
	for {
		time.Sleep(time.Second * 10)
		crc := thsCRConfig.Get()
		if crc == nil {
			fmt.Println("CRConfig nil")
		} else if crc.Stats.CDNName == nil {
			fmt.Println("CRConfig no CDN Name")
		} else if crc.Stats.DateUnixSeconds == nil {
			fmt.Println("CRConfig no Date")
		} else {
			fmt.Println("CDN: "+*crc.Stats.CDNName+" Date: ", *crc.Stats.DateUnixSeconds)
		}

		crs := thsCRStates.Get()
		if crs == nil {
			fmt.Println("CRStates nil")
		} else {
			fmt.Println("CRStates: deliveryservices: ", len(crs.DeliveryService), " caches: ", len(crs.Caches))
		}

		srvs, err := availableServers.Get("my-delivery-service-xmlid", "my-cachegroup-name")
		if err != nil {
			fmt.Println("availableServers err", err)
		} else {
			fmt.Println("availableServers", srvs)
		}
	}
}
