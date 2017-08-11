package fixtures

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

import "github.com/apache/incubator-trafficcontrol/traffic_ops/client"

// TrafficMonitorConfig returns a default TMConfigResponse to be used for testing.
func TrafficMonitorConfig() *client.TMConfigResponse {
	return &client.TMConfigResponse{
		Response: client.TrafficMonitorConfig{
			TrafficServers: []client.TrafficServer{
				client.TrafficServer{
					Profile:       "TR_CDN",
					IP:            "10.10.10.10",
					Status:        "OFFLINE",
					CacheGroup:    "tr-chicago",
					IP6:           "2001:123:abc12:64:22:::17/64",
					Port:          80,
					HostName:      "tr-chi-05",
					FQDN:          "tr-chi-05.kabletown.com",
					InterfaceName: "eth0",
					Type:          "TR",
					HashID:        "chi-tr-05",
				},
				client.TrafficServer{
					Profile:       "EDGE1_CDN",
					IP:            "3.3.3.3",
					Status:        "OFFLINE",
					CacheGroup:    "philadelphia",
					IP6:           "2009:123:456::2/64",
					Port:          80,
					HostName:      "edge-test-01",
					FQDN:          "edge-test-01.kabletown.com",
					InterfaceName: "bond0",
					Type:          "EDGE",
					HashID:        "edge-test-01",
				},
			},
			CacheGroups: []client.TMCacheGroup{
				client.TMCacheGroup{
					Name: "philadelphia",
					Coordinates: client.Coordinates{
						Longitude: 5,
						Latitude:  55,
					},
				},
				client.TMCacheGroup{
					Name: "tr-chicago",
					Coordinates: client.Coordinates{
						Longitude: 9,
						Latitude:  99,
					},
				},
			},
			Config: map[string]interface{}{
				"health.event-count":     200,
				"hack.ttl":               30,
				"health.threadPool":      4,
				"peers.polling.interval": 1000,
			},
			TrafficMonitors: []client.TrafficMonitor{
				client.TrafficMonitor{
					Port:     80,
					IP6:      "",
					IP:       "1.1.1.1",
					HostName: "traffic-monitor-01",
					FQDN:     "traffic-monitor-01@example.com",
					Profile:  "tm-123",
					Location: "philadelphia",
					Status:   "ONLINE",
				},
			},
			DeliveryServices: []client.TMDeliveryService{
				client.TMDeliveryService{
					XMLID:              "ds-05",
					TotalTPSThreshold:  0,
					Status:             "REPORTED",
					TotalKbpsThreshold: 0,
				},
			},
			Profiles: []client.TMProfile{
				client.TMProfile{
					Name: "tm-123",
					Type: "TR",
					Parameters: client.TMParameters{
						HealthConnectionTimeout: 2000,
						HealthPollingURL:        "http://${hostname}/_astats?application=&inf.name=${interface_name}",
						HistoryCount:            30,
						MinFreeKbps:             11500000,
						Thresholds: map[string]client.HealthThreshold{
							"queryTime": client.HealthThreshold{
								Val:        1000,
								Comparator: "<",
							},
							"availableBandwidthInKbps": client.HealthThreshold{
								Val:        1750000,
								Comparator: "<",
							},
							"loadavg": client.HealthThreshold{
								Val:        25.0,
								Comparator: "<",
							},
						},
					},
				},
			},
		},
	}
}
