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

func ExampleFederationResolver_Validate() {
	var typeID uint = 1
	var IPAddress string = "0.0.0.0"
	fr := FederationResolver{
		TypeID:    &typeID,
		IPAddress: &IPAddress,
	}

	fmt.Printf("%v\n", fr.Validate(nil))

	IPAddress = "0.0.0.0/24"
	fmt.Printf("%v\n", fr.Validate(nil))

	IPAddress = "::1"
	fmt.Printf("%v\n", fr.Validate(nil))

	IPAddress = "dead::babe/63"
	fmt.Printf("%v\n", fr.Validate(nil))

	IPAddress = "1.2.3.4/33"
	fmt.Printf("%v\n", fr.Validate(nil))

	IPAddress = "w.x.y.z"
	fmt.Printf("%v\n", fr.Validate(nil))

	IPAddress = "::f1d0:f00d/129"
	fmt.Printf("%v\n", fr.Validate(nil))

	IPAddress = "test::quest"
	fmt.Printf("%v\n", fr.Validate(nil))

	IPAddress = ""
	fmt.Printf("%v\n", fr.Validate(nil))

	fr.IPAddress = nil
	fmt.Printf("%v\n", fr.Validate(nil))

	fr.TypeID = nil
	fmt.Printf("%v\n", fr.Validate(nil))

	// Output:
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// ipAddress: invalid network IP or CIDR-notation subnet.
	// ipAddress: invalid network IP or CIDR-notation subnet.
	// ipAddress: invalid network IP or CIDR-notation subnet.
	// ipAddress: invalid network IP or CIDR-notation subnet.
	// ipAddress: cannot be blank.
	// ipAddress: cannot be blank.
	// ipAddress: cannot be blank; typeId: cannot be blank.
}
