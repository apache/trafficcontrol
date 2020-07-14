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
