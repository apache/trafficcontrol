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

package datareq

import (
	"github.com/apache/trafficcontrol/v7/lib/go-tc"
	"github.com/apache/trafficcontrol/v7/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/v7/traffic_monitor/todata"
	"github.com/apache/trafficcontrol/v7/traffic_ops/traffic_ops_golang/test"
	"testing"
)

func setMockTOData(tod *todata.TOData) {
	numCaches := 100
	numDSes := 100
	numCacheDSes := numDSes / 3
	numCGs := 20

	types := []tc.CacheType{tc.CacheTypeEdge, tc.CacheTypeEdge, tc.CacheTypeEdge, tc.CacheTypeEdge, tc.CacheTypeEdge, tc.CacheTypeMid}

	caches := []tc.CacheName{}
	for i := 0; i < numCaches; i++ {
		caches = append(caches, tc.CacheName(test.RandStr()))
	}

	dses := []tc.DeliveryServiceName{}
	for i := 0; i < numDSes; i++ {
		dses = append(dses, tc.DeliveryServiceName(test.RandStr()))
	}

	cgs := []tc.CacheGroupName{}
	for i := 0; i < numCGs; i++ {
		cgs = append(cgs, tc.CacheGroupName(test.RandStr()))
	}

	serverDSes := map[tc.CacheName][]tc.DeliveryServiceName{}
	for _, ca := range caches {
		for i := 0; i < numCacheDSes; i++ {
			serverDSes[ca] = append(serverDSes[ca], dses[test.RandIntn(len(dses))])
		}
	}

	dsServers := map[tc.DeliveryServiceName][]tc.CacheName{}
	for server, dses := range serverDSes {
		for _, ds := range dses {
			dsServers[ds] = append(dsServers[ds], server)
		}
	}

	serverCGs := map[tc.CacheName]tc.CacheGroupName{}
	for _, cache := range caches {
		serverCGs[cache] = cgs[test.RandIntn(len(cgs))]
	}

	serverTypes := map[tc.CacheName]tc.CacheType{}
	for _, cache := range caches {
		serverTypes[cache] = types[test.RandIntn(len(types))]
	}

	tod.DeliveryServiceServers = dsServers
	tod.ServerDeliveryServices = serverDSes
	tod.ServerTypes = serverTypes
	tod.ServerCachegroups = serverCGs
}

func TestUpdateStatusSameIpServers(t *testing.T) {
	toDataTS := todata.NewThreadsafe()
	toData := todata.New()
	setMockTOData(toData)

	toData.SameIpServers = map[tc.CacheName]map[tc.CacheName]bool{}
	toData.SameIpServers["server1_ip1_up"] = map[tc.CacheName]bool{}
	toData.SameIpServers["server1_ip1_up"]["server2_ip1_down"] = true
	toData.SameIpServers["server2_ip1_down"] = map[tc.CacheName]bool{}
	toData.SameIpServers["server2_ip1_down"]["server1_ip1_up"] = true

	toData.SameIpServers["server3_ip3_up"] = map[tc.CacheName]bool{}
	toData.SameIpServers["server3_ip3_up"]["server4_ip3_up"] = true
	toData.SameIpServers["server4_ip3_up"] = map[tc.CacheName]bool{}
	toData.SameIpServers["server4_ip3_up"]["server3_ip3_up"] = true

	localStates := peer.NewCRStatesThreadsafe()
	localStates.AddCache("server1_ip1_up", tc.IsAvailable{IsAvailable: true, Ipv4Available: true, Ipv6Available: true, Status: string(tc.CacheStatusReported)})
	localStates.AddCache("server2_ip1_down", tc.IsAvailable{IsAvailable: false, Ipv4Available: false, Ipv6Available: false, Status: string(tc.CacheStatusReported) + "too high"})
	localStates.AddCache("server3_ip3_up", tc.IsAvailable{IsAvailable: true, Ipv4Available: true, Ipv6Available: true, Status: string(tc.CacheStatusReported)})
	localStates.AddCache("server4_ip3_up", tc.IsAvailable{IsAvailable: true, Ipv4Available: true, Ipv6Available: true, Status: string(tc.CacheStatusReported)})
	localStates.AddCache("server5_ip5_up", tc.IsAvailable{IsAvailable: true, Ipv4Available: true, Ipv6Available: true, Status: string(tc.CacheStatusReported)})

	toDataTS.SetForTest(*toData)

	localStatesC := updateStatusSameIpServers(localStates, toDataTS)

	if localStatesC.Caches["server1_ip1_up"].IsAvailable == true ||
		localStatesC.Caches["server1_ip1_up"].Ipv4Available == true ||
		localStatesC.Caches["server1_ip1_up"].Ipv6Available == true {
		t.Error("expected server1_ip1_up to be false for IsAvailable Ipv4Available Ipv6Available")
	}
	if localStatesC.Caches["server2_ip1_down"].IsAvailable != false ||
		localStatesC.Caches["server2_ip1_down"].Ipv4Available != false ||
		localStatesC.Caches["server2_ip1_down"].Ipv6Available != false {
		t.Error("expected server2_ip1_up to be false for IsAvailable Ipv4Available Ipv6Available")
	}
	if localStatesC.Caches["server3_ip3_up"].IsAvailable != true ||
		localStatesC.Caches["server3_ip3_up"].Ipv4Available != true ||
		localStatesC.Caches["server3_ip3_up"].Ipv6Available != true {
		t.Error("expected server3_ip3_up to be true for IsAvailable Ipv4Available Ipv6Available")
	}
	if localStatesC.Caches["server4_ip3_up"].IsAvailable != true ||
		localStatesC.Caches["server4_ip3_up"].Ipv4Available != true ||
		localStatesC.Caches["server4_ip3_up"].Ipv6Available != true {
		t.Error("expected server4_ip3_up to be true for IsAvailable Ipv4Available Ipv6Available")
	}
	if localStatesC.Caches["server5_ip5_up"].IsAvailable != true ||
		localStatesC.Caches["server5_ip5_up"].Ipv4Available != true ||
		localStatesC.Caches["server5_ip5_up"].Ipv6Available != true {
		t.Error("expected server5_ip5_up to be true for IsAvailable Ipv4Available Ipv6Available")
	}
}
