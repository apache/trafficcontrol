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

func ExampleLegacyInterfaceDetails_ToInterfaces() {
	lid := LegacyInterfaceDetails{
		InterfaceMtu:  new(int),
		InterfaceName: new(string),
		IP6Address:    new(string),
		IP6Gateway:    new(string),
		IPAddress:     new(string),
		IPGateway:     new(string),
		IPNetmask:     new(string),
	}
	*lid.InterfaceMtu = 9000
	*lid.InterfaceName = "test"
	*lid.IP6Address = "::14/64"
	*lid.IP6Gateway = "::15"
	*lid.IPAddress = "1.2.3.4"
	*lid.IPGateway = "4.3.2.1"
	*lid.IPNetmask = "255.255.255.252"

	ifaces, err := lid.ToInterfaces(true, false)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	for _, iface := range ifaces {
		fmt.Printf("name=%s, monitor=%t\n", iface.Name, iface.Monitor)
		for _, ip := range iface.IPAddresses {
			fmt.Printf("\taddr=%s, gateway=%s, service address=%t\n", ip.Address, *ip.Gateway, ip.ServiceAddress)
		}
	}
	// Output: name=test, monitor=false
	// 	addr=1.2.3.4/30, gateway=4.3.2.1, service address=true
	// 	addr=::14/64, gateway=::15, service address=false
	//
}
