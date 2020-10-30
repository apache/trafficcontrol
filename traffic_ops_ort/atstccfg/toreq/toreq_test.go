package toreq

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
	"encoding/json"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"testing"
)

func TestServersToNullable(t *testing.T) {
	expectedNullableServer := tc.ServerNullable{
		CommonServerProperties: tc.CommonServerProperties{
			Cachegroup:     util.StrPtr("cachegroup1"),
			CachegroupID:   util.IntPtr(0),
			CDNID:          util.IntPtr(0),
			CDNName:        util.StrPtr("cdn1"),
			DomainName:     util.StrPtr("ga.atlanta.kabletown.net"),
			GUID:           util.StrPtr(""),
			HostName:       util.StrPtr("atlanta-edge-01"),
			HTTPSPort:      util.IntPtr(443),
			ID:             util.IntPtr(0),
			ILOIPAddress:   util.StrPtr("2.2.2.2"),
			ILOIPGateway:   util.StrPtr("3.3.3.3"),
			ILOIPNetmask:   util.StrPtr("255.255.0.0"),
			ILOPassword:    util.StrPtr("noonewillguessthis"),
			ILOUsername:    util.StrPtr("ilo"),
			LastUpdated:    &tc.TimeNoMod{},
			MgmtIPAddress:  util.StrPtr("0.0.0.0"),
			MgmtIPGateway:  util.StrPtr("1.1.1.1"),
			MgmtIPNetmask:  util.StrPtr("255.255.255.255"),
			OfflineReason:  util.StrPtr(""),
			PhysLocation:   util.StrPtr("Denver"),
			PhysLocationID: util.IntPtr(0),
			Profile:        util.StrPtr("EDGE1"),
			ProfileDesc:    util.StrPtr(""),
			ProfileID:      util.IntPtr(0),
			Rack:           util.StrPtr("RR 119.02"),
			RevalPending:   util.BoolPtr(false),
			RouterHostName: util.StrPtr(""),
			RouterPortName: util.StrPtr(""),
			Status:         util.StrPtr("REPORTED"),
			StatusID:       util.IntPtr(0),
			TCPPort:        util.IntPtr(80),
			Type:           "EDGE",
			TypeID:         util.IntPtr(0),
			UpdPending:     util.BoolPtr(false),
			XMPPID:         util.StrPtr("atlanta-edge-01\\@ocdn.kabletown.net"),
			XMPPPasswd:     util.StrPtr("X"),
		},
		Interfaces: []tc.ServerInterfaceInfo{{
			IPAddresses: []tc.ServerIPAddress{{
				Address:        "127.0.0.21",
				Gateway:        util.StrPtr("127.0.0.21"),
				ServiceAddress: true,
			}, {
				Address:        "2345:1234:12:8::1/64",
				Gateway:        util.StrPtr("2345:1234:12:8::1"),
				ServiceAddress: true,
			}},
			MaxBandwidth: nil,
			Monitor:      true,
			MTU:          util.UInt64Ptr(9000),
			Name:         "bond0",
		}},
	}

	serverV1 := tc.ServerV1{
		Cachegroup:     "cachegroup1",
		CachegroupID:   0,
		CDNID:          0,
		CDNName:        "cdn1",
		DomainName:     "ga.atlanta.kabletown.net",
		GUID:           "",
		HostName:       "atlanta-edge-01",
		HTTPSPort:      443,
		ID:             0,
		ILOIPAddress:   "2.2.2.2",
		ILOIPGateway:   "3.3.3.3",
		ILOIPNetmask:   "255.255.0.0",
		ILOPassword:    "noonewillguessthis",
		ILOUsername:    "ilo",
		InterfaceMtu:   9000,
		InterfaceName:  "bond0",
		LastUpdated:    tc.TimeNoMod{},
		MgmtIPAddress:  "0.0.0.0",
		MgmtIPGateway:  "1.1.1.1",
		MgmtIPNetmask:  "255.255.255.255",
		OfflineReason:  "",
		PhysLocation:   "Denver",
		PhysLocationID: 0,
		Profile:        "EDGE1",
		ProfileDesc:    "",
		ProfileID:      0,
		Rack:           "RR 119.02",
		RevalPending:   false,
		RouterHostName: "",
		RouterPortName: "",
		Status:         "REPORTED",
		StatusID:       0,
		TCPPort:        80,
		Type:           "EDGE",
		TypeID:         0,
		UpdPending:     false,
		XMPPID:         "atlanta-edge-01\\@ocdn.kabletown.net",
		XMPPPasswd:     "X",
		IPAddress:      "127.0.0.21",
		IPGateway:      "127.0.0.21",
		IP6Address:     "2345:1234:12:8::1/64",
		IP6Gateway:     "2345:1234:12:8::1",
	}
	bytes, err := json.Marshal(expectedNullableServer)
	if err != nil {
		t.Fatalf("marshalling expectedNullableServer: %s", err.Error())
	}
	expectedJSON := string(bytes)

	actualNullableServer := ServersToNullable([]tc.ServerV1{serverV1})[0]
	bytes, err = json.Marshal(actualNullableServer)
	if err != nil {
		t.Fatalf("marshalling nullable server: %s", err.Error())
	}
	actualJSON := string(bytes)

	if expectedJSON != actualJSON {
		t.Fatalf("servers did not match, expected: %s actual: %s", expectedJSON, actualJSON)
	}
}
