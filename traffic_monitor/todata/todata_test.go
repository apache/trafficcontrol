package todata

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
	"github.com/apache/trafficcontrol/v8/lib/go-tc"

	"reflect"
	"testing"
)

func TestGetDeliveryServiceServersWithTopologyBasedDeliveryService(t *testing.T) {
	topologyCrConfig := CRConfig{
		ContentServers: map[tc.CacheName]struct {
			DeliveryServices map[tc.DeliveryServiceName][]string `json:"deliveryServices"`
			CacheGroup       string                              `json:"cacheGroup"`
			Type             string                              `json:"type"`
		}{
			tc.CacheName("edge"): {
				CacheGroup: "CDN_in_a_Box_Edge",
				Type:       "EDGE",
			},
			tc.CacheName("mid"): {
				CacheGroup: "CDN_in_a_Box_Mid",
				Type:       "MID",
			},
		},
		DeliveryServices: map[tc.DeliveryServiceName]struct {
			Topology  tc.TopologyName `json:"topology"`
			Matchsets []struct {
				Protocol  string `json:"protocol"`
				MatchList []struct {
					Regex string `json:"regex"`
					Type  string `json:"match-type"`
				} `json:"matchlist"`
			} `json:"matchsets"`
		}{
			"demo1": {
				Topology: "demo1-top",
				Matchsets: []struct {
					Protocol  string `json:"protocol"`
					MatchList []struct {
						Regex string `json:"regex"`
						Type  string `json:"match-type"`
					} `json:"matchlist"`
				}{{
					Protocol: "HTTP",
					MatchList: []struct {
						Regex string `json:"regex"`
						Type  string `json:"match-type"`
					}{{Regex: `.*\.demo1\..*`}},
				}},
			}},
		Topologies: map[tc.TopologyName]struct {
			Nodes []string `json:"nodes"`
		}{
			"demo1-top": {Nodes: []string{"CDN_in_a_Box_Edge"}},
		}}

	expectedTopologiesTOData := TOData{
		DeliveryServiceServers: map[tc.DeliveryServiceName][]tc.CacheName{"demo1": {"edge"}},
		ServerDeliveryServices: map[tc.CacheName][]tc.DeliveryServiceName{"edge": {"demo1"}},
	}

	topologiesTOData := TOData{}
	topologiesTOData.DeliveryServiceServers, topologiesTOData.ServerDeliveryServices = getDeliveryServiceServers(topologyCrConfig, tc.TrafficMonitorConfigMap{})
	if !reflect.DeepEqual(expectedTopologiesTOData, topologiesTOData) {
		t.Fatalf("getDeliveryServiceServers with topology-based delivery service expected: %+v actual: %+v", expectedTopologiesTOData, topologiesTOData)
	}
}

func TestGetDeliveryServiceServersWithNonTopologyBasedDeliveryService(t *testing.T) {
	nonTopologyCrConfig := CRConfig{
		ContentServers: map[tc.CacheName]struct {
			DeliveryServices map[tc.DeliveryServiceName][]string `json:"deliveryServices"`
			CacheGroup       string                              `json:"cacheGroup"`
			Type             string                              `json:"type"`
		}{
			"edge": {
				DeliveryServices: map[tc.DeliveryServiceName][]string{
					"demo2": {"edge.demo2.mycdn.ciab.test"},
				},
				CacheGroup: "CDN_in_a_Box_Edge",
				Type:       "EDGE",
			},
			"mid": {
				CacheGroup: "CDN_in_a_Box_Mid",
				Type:       "MID",
			}},
		DeliveryServices: map[tc.DeliveryServiceName]struct {
			Topology  tc.TopologyName `json:"topology"`
			Matchsets []struct {
				Protocol  string `json:"protocol"`
				MatchList []struct {
					Regex string `json:"regex"`
					Type  string `json:"match-type"`
				} `json:"matchlist"`
			} `json:"matchsets"`
		}{
			"demo2": {
				Matchsets: []struct {
					Protocol  string `json:"protocol"`
					MatchList []struct {
						Regex string `json:"regex"`
						Type  string `json:"match-type"`
					} `json:"matchlist"`
				}{{
					Protocol: "HTTP",
					MatchList: []struct {
						Regex string `json:"regex"`
						Type  string `json:"match-type"`
					}{{Regex: `.*\.demo2\..*`}},
				}},
			},
		}}

	expectedNonTopologiesTOData := TOData{
		DeliveryServiceServers: map[tc.DeliveryServiceName][]tc.CacheName{"demo2": {"edge"}},
		ServerDeliveryServices: map[tc.CacheName][]tc.DeliveryServiceName{"edge": {"demo2"}},
	}

	nonTopologiesTOData := TOData{}
	nonTopologiesTOData.DeliveryServiceServers, nonTopologiesTOData.ServerDeliveryServices = getDeliveryServiceServers(nonTopologyCrConfig, tc.TrafficMonitorConfigMap{})
	if !reflect.DeepEqual(expectedNonTopologiesTOData, nonTopologiesTOData) {
		t.Fatalf("getDeliveryServiceServers with non-topology-based delivery service expected: %+v actual: %+v", expectedNonTopologiesTOData, nonTopologiesTOData)
	}
}
