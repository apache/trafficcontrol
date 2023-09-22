package crstatespoller

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
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/availableservers"
	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/crconfig"
	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/crstates"
	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/fetch"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

func updateAvailableServers(crcThs crconfig.Ths, crsThs crstates.Ths, as availableservers.AvailableServers) {
	newAS := map[tc.DeliveryServiceName]map[tc.CacheGroupName][]tc.CacheName{}
	crc := crcThs.Get()
	crs := crsThs.Get()
	for serverNameStr, server := range crc.ContentServers {
		serverName := tc.CacheName(serverNameStr)
		if !crs.Caches[serverName].IsAvailable {
			continue
		}
		if server.CacheGroup == nil {
			fmt.Println("ERROR updateAvailableServers CRConfig server " + serverNameStr + " cachegroup is nil")
			continue
		}
		cgName := tc.CacheGroupName(*server.CacheGroup)
		for dsNameStr, _ := range server.DeliveryServices {
			dsName := tc.DeliveryServiceName(dsNameStr)
			if newAS[dsName] == nil {
				newAS[dsName] = map[tc.CacheGroupName][]tc.CacheName{}
			}
			newAS[dsName][cgName] = append(newAS[dsName][cgName], serverName)
		}
	}
	fmt.Println("updateAvailableServers setting new ", newAS)
	as.Set(newAS)
}

// TODO implement HTTP poller
func Start(fetcher fetch.Fetcher, interval time.Duration, crc crconfig.Ths) (crstates.Ths, availableservers.AvailableServers, error) {
	thsCrs := crstates.NewThs()
	availableServers := availableservers.New()
	prevBts := []byte{}

	get := func() {
		newBts, err := fetcher.Fetch()
		if err != nil {
			fmt.Println("ERROR CRStates read error: " + err.Error())
			return
		}

		if bytes.Equal(newBts, prevBts) {
			fmt.Println("INFO CRStates unchanged.")
			return
		}

		fmt.Println("INFO CRStates changed.")
		crs := &tc.CRStates{}
		if err := json.Unmarshal(newBts, crs); err != nil {
			fmt.Println("ERROR CRStates unmarshalling: " + err.Error())
			return
		}

		thsCrs.Set(crs)
		prevBts = newBts

		updateAvailableServers(crc, thsCrs, availableServers) // TODO update AvailableServers when CRStates OR CRConfig is update, via channel and manager goroutine?

		fmt.Println("INFO CRStates set new")
		// TODO update AvailableServers
	}

	get()

	go func() {
		for {
			time.Sleep(interval)
			get()
		}
	}()
	return thsCrs, availableServers, nil
}
