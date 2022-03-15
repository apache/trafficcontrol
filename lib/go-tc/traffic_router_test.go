package tc

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

import "fmt"

func ExampleTrafficServer_IPv4() {
	server := TrafficServer{
		Interfaces: []ServerInterfaceInfo{
			{
				IPAddresses: []ServerIPAddress{
					{
						Address:        "192.0.2.1",
						ServiceAddress: true,
					},
					{
						Address:        "192.0.2.2",
						ServiceAddress: false,
					},
				},
				Monitor: true,
			},
			{
				IPAddresses: []ServerIPAddress{
					{
						Address:        "192.0.2.3",
						ServiceAddress: false,
					},
				},
				Monitor: true,
			},
		},
	}

	fmt.Println(server.IPv4())
	// Output: 192.0.2.1
}

func ExampleTrafficServer_IPv6() {
	server := TrafficServer{
		Interfaces: []ServerInterfaceInfo{
			{
				IPAddresses: []ServerIPAddress{
					{
						Address:        "2001:DB8::1",
						ServiceAddress: false,
					},
					{
						Address:        "2001:DB8::2",
						ServiceAddress: false,
					},
				},
				Monitor: true,
			},
			{
				IPAddresses: []ServerIPAddress{
					{
						Address:        "2001:DB8::3",
						ServiceAddress: true,
					},
				},
				Monitor: true,
			},
		},
	}

	fmt.Println(server.IPv6())
	// Output: 2001:DB8::3
}

func ExampleLegacyTrafficServer_Upgrade() {
	lts := LegacyTrafficServer{
		CacheGroup:       "testCG",
		DeliveryServices: []TSDeliveryService{},
		FQDN:             "test.quest",
		HashID:           "test",
		HostName:         "test",
		HTTPSPort:        -1,
		InterfaceName:    "testInterface",
		IP:               "198.0.2.0",
		IP6:              "2001:DB8::1",
		Port:             -1,
		Profile:          "testProfile",
		ServerStatus:     "testStatus",
		Type:             "testType",
	}

	ts := lts.Upgrade()
	fmt.Println("CacheGroup:", ts.CacheGroup)
	fmt.Println("# of DeliveryServices:", len(ts.DeliveryServices))
	fmt.Println("FQDN:", ts.FQDN)
	fmt.Println("HashID:", ts.HashID)
	fmt.Println("HostName:", ts.HostName)
	fmt.Println("HTTPSPort:", ts.HTTPSPort)
	fmt.Println("# of Interfaces:", len(ts.Interfaces))
	fmt.Println("Interface Name:", ts.Interfaces[0].Name)
	fmt.Println("# of Interface IP Addresses:", len(ts.Interfaces[0].IPAddresses))
	fmt.Println("first IP Address:", ts.Interfaces[0].IPAddresses[0].Address)
	fmt.Println("second IP Address:", ts.Interfaces[0].IPAddresses[1].Address)
	fmt.Println("Port:", ts.Port)
	fmt.Println("Profile:", ts.Profile)
	fmt.Println("ServerStatus:", ts.ServerStatus)
	fmt.Println("Type:", ts.Type)

	// Output: CacheGroup: testCG
	// # of DeliveryServices: 0
	// FQDN: test.quest
	// HashID: test
	// HostName: test
	// HTTPSPort: -1
	// # of Interfaces: 1
	// Interface Name: testInterface
	// # of Interface IP Addresses: 2
	// first IP Address: 198.0.2.0
	// second IP Address: 2001:DB8::1
	// Port: -1
	// Profile: testProfile
	// ServerStatus: testStatus
	// Type: testType
}
