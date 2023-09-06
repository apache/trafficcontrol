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

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

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
	// Output: name=test, monitor=true
	// 	addr=1.2.3.4/30, gateway=4.3.2.1, service address=true
	// 	addr=::14/64, gateway=::15, service address=false
	//
}

func ExampleLegacyInterfaceDetails_String() {
	ipv4 := "192.0.2.0"
	ipv6 := "2001:DB8::/64"
	name := "test"
	mtu := 9000

	lid := LegacyInterfaceDetails{
		InterfaceMtu:  &mtu,
		InterfaceName: &name,
		IP6Address:    &ipv6,
		IP6Gateway:    nil,
		IPAddress:     &ipv4,
		IPGateway:     nil,
		IPNetmask:     nil,
	}

	fmt.Println(lid.String())

	// Output: LegacyInterfaceDetails(InterfaceMtu=9000, InterfaceName='test', IP6Address='2001:DB8::/64', IP6Gateway=nil, IPAddress='192.0.2.0', IPGateway=nil, IPNetmask=nil)
}

func ExampleServerIPAddress_Copy() {
	ip := ServerIPAddress{
		Address:        "test",
		Gateway:        new(string),
		ServiceAddress: false,
	}

	*ip.Gateway = "not a gateway, but who cares?"
	ip2 := ip.Copy()
	fmt.Println(*ip.Gateway == *ip2.Gateway)

	*ip.Gateway = "something different"
	fmt.Println(*ip.Gateway == *ip2.Gateway)

	// Output: true
	// false
}

func ExampleServerInterfaceInfoV40_Copy() {
	inf := ServerInterfaceInfoV40{
		ServerInterfaceInfo: ServerInterfaceInfo{
			IPAddresses: []ServerIPAddress{
				{
					Address:        "test",
					Gateway:        new(string),
					ServiceAddress: false,
				},
			},
			MaxBandwidth: new(uint64),
			Monitor:      false,
			MTU:          new(uint64),
			Name:         "eth0",
		},
		RouterHostName: "router host",
		RouterPortName: "router port",
	}

	*inf.IPAddresses[0].Gateway = "not a gateway, but who cares?"
	inf2 := inf.Copy()

	fmt.Println(*inf.IPAddresses[0].Gateway == *inf2.IPAddresses[0].Gateway)
	*inf.IPAddresses[0].Gateway = "something different"
	fmt.Println(*inf.IPAddresses[0].Gateway == *inf2.IPAddresses[0].Gateway)

	// Output: true
	// false
}

func TestServerV5DowngradeUpgrade(t *testing.T) {
	serverV5 := ServerV50{
		CacheGroup:         "Cache Group",
		CacheGroupID:       1,
		CDNID:              2,
		CDN:                "CDN",
		DomainName:         "domain",
		GUID:               nil,
		HostName:           "host",
		HTTPSPort:          nil,
		ID:                 3,
		ILOIPAddress:       nil,
		ILOIPGateway:       nil,
		ILOIPNetmask:       nil,
		ILOPassword:        nil,
		ILOUsername:        nil,
		LastUpdated:        time.Time{}.Add(time.Hour),
		MgmtIPAddress:      nil,
		MgmtIPGateway:      nil,
		MgmtIPNetmask:      nil,
		OfflineReason:      nil,
		PhysicalLocation:   "physical location",
		PhysicalLocationID: 4,
		Profiles:           []string{"test", "quest"},
		Rack:               nil,
		Status:             "Status",
		StatusID:           5,
		TCPPort:            nil,
		Type:               "type",
		TypeID:             6,
		XMPPID:             nil,
		XMPPPasswd:         nil,
		Interfaces: []ServerInterfaceInfoV40{
			{
				ServerInterfaceInfo: ServerInterfaceInfo{
					IPAddresses: []ServerIPAddress{
						{
							Address:        "192.0.0.1/12",
							Gateway:        nil,
							ServiceAddress: true,
						},
					},
					MaxBandwidth: nil,
					Monitor:      false,
					MTU:          nil,
					Name:         "eth0",
				},
				RouterHostName: "router host",
				RouterPortName: "router port",
			},
		},
		StatusLastUpdated:  nil,
		ConfigUpdateTime:   nil,
		ConfigApplyTime:    nil,
		ConfigUpdateFailed: false,
		RevalUpdateTime:    nil,
		RevalApplyTime:     nil,
		RevalUpdateFailed:  false,
	}

	serverV4 := serverV5.Downgrade()
	if fqdn := serverV5.HostName + "." + serverV5.DomainName; serverV4.FQDN == nil || *serverV4.FQDN != fqdn {
		t.Errorf("incorrectly calculated FQDN; want: %s, got: %v", fqdn, serverV4.FQDN)
	}

	if !reflect.DeepEqual(serverV4.Upgrade(), serverV5) {
		t.Error("server not equal after downgrading then upgrading")
	}
}

type interfaceTest struct {
	ExpectedIPv4        string
	ExpectedIPv4Gateway string
	ExpectedIPv6        string
	ExpectedIPv6Gateway string
	ExpectedMTU         *uint64
	ExpectedName        string
	ExpectedNetmask     string
	Interfaces          []ServerInterfaceInfo
}

// tests a set of interfaces' conversion to legacy format against expected
// values.
// Note: This doesn't distinguish between nil and pointer-to-empty-string values
// when a value is not expected. That's because all ATC components treat null
// and empty-string values the same, so it's not important which is returned by
// the conversion process (and in fact expecting one or the other could
// potentially break some applications).
func testInfs(expected interfaceTest, t *testing.T) {
	lid, err := InterfaceInfoToLegacyInterfaces(expected.Interfaces)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if lid.InterfaceName == nil {
		t.Error("Unexpectedly nil Interface Name")
	} else if *lid.InterfaceName != expected.ExpectedName {
		t.Errorf("Incorrect Interface Name; want: '%s', got: '%s'", expected.ExpectedName, *lid.InterfaceName)
	}

	if expected.ExpectedMTU != nil {
		if lid.InterfaceMtu == nil {
			t.Error("Unexpectedly nil Interface MTU")
		} else if uint64(*lid.InterfaceMtu) != *expected.ExpectedMTU {
			t.Errorf("Incorrect Interface MTU; want: %d, got: %d", *expected.ExpectedMTU, *lid.InterfaceMtu)
		}
	} else if lid.InterfaceMtu != nil {
		t.Error("Unexpectedly non-nil Interface MTU")
	}

	if expected.ExpectedIPv4 != "" {
		if lid.IPAddress == nil {
			t.Error("Unexpectedly nil IPv4 Address")
		} else if *lid.IPAddress != expected.ExpectedIPv4 {
			t.Errorf("Incorrect IPv4 Address; want: '%s', got: '%s'", expected.ExpectedIPv4, *lid.IPAddress)
		}
	} else if lid.IPAddress != nil && *lid.IPAddress != "" {
		t.Error("Unexpectedly non-empty IPv4 Address")
	}

	if expected.ExpectedIPv4Gateway != "" {
		if lid.IPGateway == nil {
			t.Error("Unexpectedly nil IPv4 Gateway")
		} else if *lid.IPGateway != expected.ExpectedIPv4Gateway {
			t.Errorf("Incorrect IPv4 Gateway; want: '%s', got: '%s'", expected.ExpectedIPv4Gateway, *lid.IPGateway)
		}
	} else if lid.IPGateway != nil && *lid.IPGateway != "" {
		t.Error("Unexpectedly non-empty IPv4 Gateway")
	}

	if expected.ExpectedNetmask != "" {
		if lid.IPNetmask == nil {
			t.Error("Unexpectedly nil IPv4 Netmask")
		} else if *lid.IPNetmask != expected.ExpectedNetmask {
			t.Errorf("Incorrect IPv4 Netmask; want: '%s', got: '%s'", expected.ExpectedNetmask, *lid.IPNetmask)
		}
	} else if lid.IPNetmask != nil && *lid.IPNetmask != "" {
		t.Error("Unexpectedly non-empty IPv4 Netmask")
	}

	if expected.ExpectedIPv6 != "" {
		if lid.IP6Address == nil {
			t.Error("Unexpectedly nil IPv6 Address")
		} else if *lid.IP6Address != expected.ExpectedIPv6 {
			t.Errorf("Incorrect IPv6 Address; want: '%s', got: '%s'", expected.ExpectedIPv6, *lid.IP6Address)
		}
	} else if lid.IP6Address != nil && *lid.IP6Address != "" {
		t.Error("Unexpectedly non-empty IPv6 Address")
	}

	if expected.ExpectedIPv6Gateway != "" {
		if lid.IP6Gateway == nil {
			t.Error("Unexpectedly nil IPv6 Gateway")
		} else if *lid.IP6Gateway != expected.ExpectedIPv6Gateway {
			t.Errorf("Incorrect IPv6 Gateway; want: '%s', got: '%s'", expected.ExpectedIPv6Gateway, *lid.IP6Gateway)
		}
	} else if lid.IP6Gateway != nil && *lid.IP6Gateway != "" {
		t.Error("Unexpectedly non-empty IPv6 Gateway")
	}
}

func TestInterfaceInfoToLegacyInterfaces(t *testing.T) {
	var mtu uint64 = 9000
	ipv4Gateway := "192.0.2.2"
	ipv6Gateway := "2001:DB8::2"

	cases := map[string]interfaceTest{
		"single interface, IPv4 only, no gateway, MTU, or netmask": interfaceTest{
			ExpectedIPv4:        "192.0.2.0",
			ExpectedIPv4Gateway: "",
			ExpectedIPv6:        "",
			ExpectedIPv6Gateway: "",
			ExpectedMTU:         nil,
			ExpectedName:        "test",
			ExpectedNetmask:     "",
			Interfaces: []ServerInterfaceInfo{
				ServerInterfaceInfo{
					MTU:  nil,
					Name: "test",
					IPAddresses: []ServerIPAddress{
						ServerIPAddress{
							Address:        "192.0.2.0",
							Gateway:        nil,
							ServiceAddress: true,
						},
					},
				},
			},
		},
		"single interface, IPv4 only, no gateway or netmask": interfaceTest{
			ExpectedIPv4:        "192.0.2.0",
			ExpectedIPv4Gateway: "",
			ExpectedIPv6:        "",
			ExpectedIPv6Gateway: "",
			ExpectedMTU:         &mtu,
			ExpectedName:        "test",
			ExpectedNetmask:     "",
			Interfaces: []ServerInterfaceInfo{
				ServerInterfaceInfo{
					MTU:  &mtu,
					Name: "test",
					IPAddresses: []ServerIPAddress{
						ServerIPAddress{
							Address:        "192.0.2.0",
							Gateway:        nil,
							ServiceAddress: true,
						},
					},
				},
			},
		},
		"single interface, IPv4 only, no netmask": interfaceTest{ // Final Destination
			ExpectedIPv4:        "192.0.2.0",
			ExpectedIPv4Gateway: ipv4Gateway,
			ExpectedIPv6:        "",
			ExpectedIPv6Gateway: "",
			ExpectedMTU:         &mtu,
			ExpectedName:        "test",
			ExpectedNetmask:     "",
			Interfaces: []ServerInterfaceInfo{
				ServerInterfaceInfo{
					MTU:  &mtu,
					Name: "test",
					IPAddresses: []ServerIPAddress{
						ServerIPAddress{
							Address:        "192.0.2.0",
							Gateway:        &ipv4Gateway,
							ServiceAddress: true,
						},
					},
				},
			},
		},
		"single interface, IPv4 only": interfaceTest{
			ExpectedIPv4:        "192.0.2.0",
			ExpectedIPv4Gateway: ipv4Gateway,
			ExpectedIPv6:        "",
			ExpectedIPv6Gateway: "",
			ExpectedMTU:         &mtu,
			ExpectedName:        "test",
			ExpectedNetmask:     "255.255.255.0",
			Interfaces: []ServerInterfaceInfo{
				ServerInterfaceInfo{
					MTU:  &mtu,
					Name: "test",
					IPAddresses: []ServerIPAddress{
						ServerIPAddress{
							Address:        "192.0.2.0/24",
							Gateway:        &ipv4Gateway,
							ServiceAddress: true,
						},
					},
				},
			},
		},
		"single interface, no gateway, MTU, or netmask": interfaceTest{
			ExpectedIPv4:        "192.0.2.0",
			ExpectedIPv4Gateway: "",
			ExpectedIPv6:        "2001:DB8::1",
			ExpectedIPv6Gateway: "",
			ExpectedMTU:         nil,
			ExpectedName:        "test",
			ExpectedNetmask:     "",
			Interfaces: []ServerInterfaceInfo{
				ServerInterfaceInfo{
					MTU:  nil,
					Name: "test",
					IPAddresses: []ServerIPAddress{
						ServerIPAddress{
							Address:        "192.0.2.0",
							Gateway:        nil,
							ServiceAddress: true,
						},
						ServerIPAddress{
							Address:        "2001:DB8::1",
							Gateway:        nil,
							ServiceAddress: true,
						},
					},
				},
			},
		},
		"single interface": interfaceTest{
			ExpectedIPv4:        "192.0.2.0",
			ExpectedIPv4Gateway: ipv4Gateway,
			ExpectedIPv6:        "2001:DB8::1",
			ExpectedIPv6Gateway: ipv6Gateway,
			ExpectedMTU:         &mtu,
			ExpectedName:        "test",
			ExpectedNetmask:     "255.255.255.0",
			Interfaces: []ServerInterfaceInfo{
				ServerInterfaceInfo{
					MTU:  &mtu,
					Name: "test",
					IPAddresses: []ServerIPAddress{
						ServerIPAddress{
							Address:        "192.0.2.0/24",
							Gateway:        &ipv4Gateway,
							ServiceAddress: true,
						},
						ServerIPAddress{
							Address:        "2001:DB8::1",
							Gateway:        &ipv6Gateway,
							ServiceAddress: true,
						},
					},
				},
			},
		},
		"single interface, extra IP addresses": interfaceTest{
			ExpectedIPv4:        "192.0.2.0",
			ExpectedIPv4Gateway: ipv4Gateway,
			ExpectedIPv6:        "2001:DB8::1",
			ExpectedIPv6Gateway: ipv6Gateway,
			ExpectedMTU:         &mtu,
			ExpectedName:        "test",
			ExpectedNetmask:     "255.255.255.0",
			Interfaces: []ServerInterfaceInfo{
				ServerInterfaceInfo{
					MTU:  &mtu,
					Name: "test",
					IPAddresses: []ServerIPAddress{
						ServerIPAddress{
							Address:        "192.0.2.1/5",
							Gateway:        nil,
							ServiceAddress: false,
						},
						ServerIPAddress{
							Address:        "192.0.2.0/24",
							Gateway:        &ipv4Gateway,
							ServiceAddress: true,
						},
						ServerIPAddress{
							Address:        "2001:DB8::2",
							Gateway:        nil,
							ServiceAddress: false,
						},
						ServerIPAddress{
							Address:        "2001:DB8::1",
							Gateway:        &ipv6Gateway,
							ServiceAddress: true,
						},
						ServerIPAddress{
							Address:        "192.0.2.2/20",
							Gateway:        nil,
							ServiceAddress: false,
						},
					},
				},
			},
		},
		"multiple interfaces, IPv4 only, no netmask": interfaceTest{
			ExpectedIPv4:        "192.0.2.1",
			ExpectedIPv4Gateway: ipv4Gateway,
			ExpectedMTU:         &mtu,
			ExpectedName:        "test",
			ExpectedNetmask:     "",
			ExpectedIPv6:        "",
			ExpectedIPv6Gateway: "",
			Interfaces: []ServerInterfaceInfo{
				{
					IPAddresses: []ServerIPAddress{
						{
							Address:        "192.0.2.1",
							Gateway:        &ipv4Gateway,
							ServiceAddress: true,
						},
					},
					MaxBandwidth: nil,
					Monitor:      true,
					MTU:          &mtu,
					Name:         "test",
				},
				{
					IPAddresses: []ServerIPAddress{
						{
							Address:        "192.0.2.2",
							Gateway:        nil,
							ServiceAddress: false,
						},
					},
					MaxBandwidth: nil,
					Monitor:      false,
					MTU:          &mtu,
					Name:         "invalid",
				},
			},
		},
		"multiple interfaces": interfaceTest{
			ExpectedIPv4:        "192.0.2.0",
			ExpectedIPv4Gateway: ipv4Gateway,
			ExpectedIPv6:        "2001:DB8::1",
			ExpectedIPv6Gateway: ipv6Gateway,
			ExpectedMTU:         &mtu,
			ExpectedName:        "test",
			ExpectedNetmask:     "255.255.255.0",
			Interfaces: []ServerInterfaceInfo{
				ServerInterfaceInfo{
					MTU:  nil,
					Name: "invalid1",
					IPAddresses: []ServerIPAddress{
						ServerIPAddress{
							Address:        "192.0.2.1/5",
							Gateway:        nil,
							ServiceAddress: false,
						},
						ServerIPAddress{
							Address:        "2001:DB8::2",
							Gateway:        nil,
							ServiceAddress: false,
						},
					},
				},
				ServerInterfaceInfo{
					MTU:  &mtu,
					Name: "test",
					IPAddresses: []ServerIPAddress{
						ServerIPAddress{
							Address:        "192.0.2.0/24",
							Gateway:        &ipv4Gateway,
							ServiceAddress: true,
						},
						ServerIPAddress{
							Address:        "2001:DB8::1",
							Gateway:        &ipv6Gateway,
							ServiceAddress: true,
						},
					},
				},
				ServerInterfaceInfo{
					MTU:  nil,
					Name: "invalid2",
					IPAddresses: []ServerIPAddress{
						ServerIPAddress{
							Address:        "192.0.2.2/7",
							Gateway:        nil,
							ServiceAddress: false,
						},
						ServerIPAddress{
							Address:        "2001:DB8::3/12",
							Gateway:        nil,
							ServiceAddress: false,
						},
					},
				},
			},
		},
	}

	for description, test := range cases {
		t.Run(description, func(t *testing.T) { testInfs(test, t) })
	}
}

func ensureNoNulls(s ServerNullableV2, t *testing.T) {
	if s.Cachegroup == nil {
		t.Error("nullable conversion gave nil Cachegroup")
	}

	if s.CachegroupID == nil {
		t.Error("nullable conversion gave nil CachegroupID")
	}

	if s.CDNID == nil {
		t.Error("nullable conversion gave nil CDNID")
	}

	if s.CDNName == nil {
		t.Error("nullable conversion gave nil CDNName")
	}

	if s.DeliveryServices == nil {
		t.Error("nullable conversion gave nil DeliveryServices")
	}

	if s.DomainName == nil {
		t.Error("nullable conversion gave nil DomainName")
	}

	if s.FQDN == nil {
		t.Error("nullable conversion gave nil FQDN")
	}

	if s.GUID == nil {
		t.Error("nullable conversion gave nil GUID")
	}

	if s.HostName == nil {
		t.Error("nullable conversion gave nil HostName")
	}

	if s.HTTPSPort == nil {
		t.Error("nullable conversion gave nil HTTPSPort")
	}

	if s.ID == nil {
		t.Error("nullable conversion gave nil ID")
	}

	if s.ILOIPAddress == nil {
		t.Error("nullable conversion gave nil ILOIPAddress")
	}

	if s.ILOIPGateway == nil {
		t.Error("nullable conversion gave nil ILOIPGateway")
	}

	if s.ILOIPNetmask == nil {
		t.Error("nullable conversion gave nil ILOIPNetmask")
	}

	if s.ILOPassword == nil {
		t.Error("nullable conversion gave nil ILOPassword")
	}

	if s.ILOUsername == nil {
		t.Error("nullable conversion gave nil ILOUsername")
	}

	if s.InterfaceMtu == nil {
		t.Error("nullable conversion gave nil InterfaceMtu")
	}

	if s.InterfaceName == nil {
		t.Error("nullable conversion gave nil InterfaceName")
	}

	if s.IP6Address == nil {
		t.Error("nullable conversion gave nil IP6Address")
	}

	if s.IP6IsService == nil {
		t.Error("nullable conversion gave nil IP6IsService")
	}

	if s.IP6Gateway == nil {
		t.Error("nullable conversion gave nil IP6Gateway")
	}

	if s.IPAddress == nil {
		t.Error("nullable conversion gave nil IPAddress")
	}

	if s.IPIsService == nil {
		t.Error("nullable conversion gave nil IPIsService")
	}

	if s.IPGateway == nil {
		t.Error("nullable conversion gave nil IPGateway")
	}

	if s.IPNetmask == nil {
		t.Error("nullable conversion gave nil IPNetmask")
	}

	if s.LastUpdated == nil {
		t.Error("nullable conversion gave nil LastUpdated")
	}

	if s.MgmtIPAddress == nil {
		t.Error("nullable conversion gave nil MgmtIPAddress")
	}

	if s.MgmtIPGateway == nil {
		t.Error("nullable conversion gave nil MgmtIPGateway")
	}

	if s.MgmtIPNetmask == nil {
		t.Error("nullable conversion gave nil MgmtIPNetmask")
	}

	if s.OfflineReason == nil {
		t.Error("nullable conversion gave nil OfflineReason")
	}

	if s.PhysLocation == nil {
		t.Error("nullable conversion gave nil PhysLocation")
	}

	if s.PhysLocationID == nil {
		t.Error("nullable conversion gave nil PhysLocationID")
	}

	if s.Profile == nil {
		t.Error("nullable conversion gave nil Profile")
	}

	if s.ProfileDesc == nil {
		t.Error("nullable conversion gave nil ProfileDesc")
	}

	if s.ProfileID == nil {
		t.Error("nullable conversion gave nil ProfileID")
	}

	if s.Rack == nil {
		t.Error("nullable conversion gave nil Rack")
	}

	if s.RevalPending == nil {
		t.Error("nullable conversion gave nil RevalPending")
	}

	if s.RouterHostName == nil {
		t.Error("nullable conversion gave nil RouterHostName")
	}

	if s.RouterPortName == nil {
		t.Error("nullable conversion gave nil RouterPortName")
	}

	if s.Status == nil {
		t.Error("nullable conversion gave nil Status")
	}

	if s.StatusID == nil {
		t.Error("nullable conversion gave nil StatusID")
	}

	if s.TCPPort == nil {
		t.Error("nullable conversion gave nil TCPPort")
	}

	if s.TypeID == nil {
		t.Error("nullable conversion gave nil TypeID")
	}

	if s.UpdPending == nil {
		t.Error("nullable conversion gave nil UpdPending")
	}

	if s.XMPPID == nil {
		t.Error("nullable conversion gave nil XMPPID")
	}

	if s.XMPPPasswd == nil {
		t.Error("nullable conversion gave nil XMPPPasswd")
	}
}

func TestServer_ToNullable(t *testing.T) {
	fqdn := "testFQDN"
	srv := Server{
		Cachegroup:       "testCachegroup",
		CachegroupID:     42,
		CDNID:            43,
		CDNName:          "testCDNName",
		DeliveryServices: map[string][]string{"test": []string{"quest"}},
		DomainName:       "testDomainName",
		FQDN:             &fqdn,
		FqdnTime:         time.Now(),
		GUID:             "testGUID",
		HostName:         "testHostName",
		HTTPSPort:        -1,
		ID:               44,
		ILOIPAddress:     "testILOIPAddress",
		ILOIPGateway:     "testILOIPGateway",
		ILOIPNetmask:     "testILOIPNetmask",
		ILOPassword:      "testILOPassword",
		ILOUsername:      "testILOUsername",
		InterfaceMtu:     -2,
		InterfaceName:    "testInterfaceName",
		IP6Address:       "testIP6Address",
		IP6IsService:     true,
		IP6Gateway:       "testIP6Gateway",
		IPAddress:        "testIPAddress",
		IPIsService:      false,
		IPGateway:        "testIPGateway",
		IPNetmask:        "testIPNetmask",
		LastUpdated:      TimeNoMod(Time{Time: time.Now().Add(time.Minute), Valid: true}),
		MgmtIPAddress:    "testMgmtIPAddress",
		MgmtIPGateway:    "testMgmtIPGateway",
		MgmtIPNetmask:    "testMgmtIPNetmask",
		OfflineReason:    "testOfflineReason",
		PhysLocation:     "testPhysLocation",
		PhysLocationID:   45,
		Profile:          "testProfile",
		ProfileDesc:      "testProfileDesc",
		ProfileID:        46,
		Rack:             "testRack",
		RevalPending:     true,
		RouterHostName:   "testRouterHostName",
		RouterPortName:   "testRouterPortName",
		Status:           "testStatus",
		StatusID:         47,
		TCPPort:          -3,
		Type:             "testType",
		TypeID:           48,
		UpdPending:       false,
		XMPPID:           "testXMPPID",
		XMPPPasswd:       "testXMPPasswd",
	}

	nullable := srv.ToNullable()

	if nullable.FqdnTime != srv.FqdnTime {
		t.Errorf("Incorrect FqdnTime after nullable conversion; want: '%s', got: '%s'", srv.FqdnTime, nullable.FqdnTime)
	}

	if nullable.Type != srv.Type {
		t.Errorf("Incorrect Type after nullable conversion; want: '%s', got: '%s'", srv.Type, nullable.Type)
	}

	noNulls := t.Run("check nullable-converted server for null values", func(t *testing.T) { ensureNoNulls(nullable, t) })
	if !noNulls {
		t.Fatal("Cannot check nullable server field referred-to values - null values found")
	}

	if *nullable.Cachegroup != srv.Cachegroup {
		t.Errorf("Incorrect Cachegroup after nullable conversion; want: '%s', got: '%s'", srv.Cachegroup, *nullable.Cachegroup)
	}

	if *nullable.CachegroupID != srv.CachegroupID {
		t.Errorf("Incorrect CachegroupID after nullable conversion; want: %d, got: %d", srv.CachegroupID, *nullable.CachegroupID)
	}

	if *nullable.CDNID != srv.CDNID {
		t.Errorf("Incorrect CDNID after nullable conversion; want: %d, got: %d", srv.CDNID, *nullable.CDNID)
	}

	if *nullable.CDNName != srv.CDNName {
		t.Errorf("Incorrect CDNName after nullable conversion; want: '%s', got: '%s'", srv.CDNName, *nullable.CDNName)
	}

	if len(*nullable.DeliveryServices) != len(srv.DeliveryServices) {
		t.Errorf("Incorrect number of DeliveryServices after nullable conversion; want: %d, got: %d", len(srv.DeliveryServices), len(*nullable.DeliveryServices))
	} else {
		for k, v := range srv.DeliveryServices {
			nullableV, ok := (*nullable.DeliveryServices)[k]
			if !ok {
				t.Errorf("Missing Delivery Service '%s' after nullable conversion", k)
				continue
			}
			if len(nullableV) != len(v) {
				t.Errorf("Delivery Service '%s' has incorrect length after nullable conversion; want: %d, got: %d", k, len(v), len(nullableV))
			}
			for i, ds := range v {
				nullableDS := nullableV[i]
				if nullableDS != ds {
					t.Errorf("Incorrect value at position %d in Delivery Service '%s' after nullable conversion; want: '%s', got: '%s'", i, k, ds, nullableDS)
				}
			}
		}
	}

	if *nullable.DomainName != srv.DomainName {
		t.Errorf("Incorrect DomainName after nullable conversion; want: '%s', got: '%s'", srv.DomainName, *nullable.DomainName)
	}

	if *nullable.FQDN != fqdn {
		t.Errorf("Incorrect FQDN after nullable conversion; want: '%s', got: '%s'", fqdn, *nullable.FQDN)
	}

	if *nullable.GUID != srv.GUID {
		t.Errorf("Incorrect GUID after nullable conversion; want: '%s', got: '%s'", srv.GUID, *nullable.GUID)
	}

	if *nullable.HostName != srv.HostName {
		t.Errorf("Incorrect HostName after nullable conversion; want: '%s', got: '%s'", srv.HostName, *nullable.HostName)
	}

	if *nullable.HTTPSPort != srv.HTTPSPort {
		t.Errorf("Incorrect HTTPSPort after nullable conversion; want: %d, got: %d", srv.HTTPSPort, *nullable.HTTPSPort)
	}

	if *nullable.ID != srv.ID {
		t.Errorf("Incorrect ID after nullable conversion; want: %d, got: %d", srv.ID, *nullable.ID)
	}

	if *nullable.ILOIPAddress != srv.ILOIPAddress {
		t.Errorf("Incorrect ILOIPAddress after nullable conversion; want: '%s', got: '%s'", srv.ILOIPAddress, *nullable.ILOIPAddress)
	}

	if *nullable.ILOIPGateway != srv.ILOIPGateway {
		t.Errorf("Incorrect ILOIPGateway after nullable conversion; want: '%s', got: '%s'", srv.ILOIPGateway, *nullable.ILOIPGateway)
	}

	if *nullable.ILOIPNetmask != srv.ILOIPNetmask {
		t.Errorf("Incorrect ILOIPNetmask after nullable conversion; want: '%s', got: '%s'", srv.ILOIPNetmask, *nullable.ILOIPNetmask)
	}

	if *nullable.ILOPassword != srv.ILOPassword {
		t.Errorf("Incorrect ILOPassword after nullable conversion; want: '%s', got: '%s'", srv.ILOPassword, *nullable.ILOPassword)
	}

	if *nullable.ILOUsername != srv.ILOUsername {
		t.Errorf("Incorrect ILOUsername after nullable conversion; want: '%s', got: '%s'", srv.ILOUsername, *nullable.ILOUsername)
	}

	if *nullable.InterfaceMtu != srv.InterfaceMtu {
		t.Errorf("Incorrect InterfaceMtu after nullable conversion; want: %d, got: %d", srv.InterfaceMtu, *nullable.InterfaceMtu)
	}

	if *nullable.InterfaceName != srv.InterfaceName {
		t.Errorf("Incorrect InterfaceName after nullable conversion; want: '%s', got: '%s'", srv.InterfaceName, *nullable.InterfaceName)
	}

	if *nullable.IP6Address != srv.IP6Address {
		t.Errorf("Incorrect IP6Address after nullable conversion; want: '%s', got: '%s'", srv.IP6Address, *nullable.IP6Address)
	}

	if *nullable.IP6IsService != srv.IP6IsService {
		t.Errorf("Incorrect IP6IsService after nullable conversion; want: %t, got: %t", srv.IP6IsService, *nullable.IP6IsService)
	}

	if *nullable.IP6Gateway != srv.IP6Gateway {
		t.Errorf("Incorrect IP6Gateway after nullable conversion; want: '%s', got: '%s'", srv.IP6Gateway, *nullable.IP6Gateway)
	}

	if *nullable.IPAddress != srv.IPAddress {
		t.Errorf("Incorrect IPAddress after nullable conversion; want: '%s', got: '%s'", srv.IPAddress, *nullable.IPAddress)
	}

	if *nullable.IPIsService != srv.IPIsService {
		t.Errorf("Incorrect IPIsService after nullable conversion; want: %t, got: %t", srv.IPIsService, *nullable.IPIsService)
	}

	if *nullable.IPGateway != srv.IPGateway {
		t.Errorf("Incorrect IPGateway after nullable conversion; want: '%s', got: '%s'", srv.IPGateway, *nullable.IPGateway)
	}

	if *nullable.IPNetmask != srv.IPNetmask {
		t.Errorf("Incorrect IPNetmask after nullable conversion; want: '%s', got: '%s'", srv.IPNetmask, *nullable.IPNetmask)
	}

	if *nullable.LastUpdated != srv.LastUpdated {
		t.Errorf("Incorrect LastUpdated after nullable conversion; want: '%s', got: '%s'", srv.LastUpdated, *nullable.LastUpdated)
	}

	if *nullable.MgmtIPAddress != srv.MgmtIPAddress {
		t.Errorf("Incorrect MgmtIPAddress after nullable conversion; want: '%s', got: '%s'", srv.MgmtIPAddress, *nullable.MgmtIPAddress)
	}

	if *nullable.MgmtIPGateway != srv.MgmtIPGateway {
		t.Errorf("Incorrect MgmtIPGateway after nullable conversion; want: '%s', got: '%s'", srv.MgmtIPGateway, *nullable.MgmtIPGateway)
	}

	if *nullable.MgmtIPNetmask != srv.MgmtIPNetmask {
		t.Errorf("Incorrect MgmtIPNetmask after nullable conversion; want: '%s', got: '%s'", srv.MgmtIPNetmask, *nullable.MgmtIPNetmask)
	}

	if *nullable.OfflineReason != srv.OfflineReason {
		t.Errorf("Incorrect OfflineReason after nullable conversion; want: '%s', got: '%s'", srv.OfflineReason, *nullable.OfflineReason)
	}

	if *nullable.PhysLocation != srv.PhysLocation {
		t.Errorf("Incorrect PhysLocation after nullable conversion; want: '%s', got: '%s'", srv.PhysLocation, *nullable.PhysLocation)
	}

	if *nullable.PhysLocationID != srv.PhysLocationID {
		t.Errorf("Incorrect PhysLocationID after nullable conversion; want: %d, got: %d", srv.PhysLocationID, *nullable.PhysLocationID)
	}

	if *nullable.Profile != srv.Profile {
		t.Errorf("Incorrect Profile after nullable conversion; want: '%s', got: '%s'", srv.Profile, *nullable.Profile)
	}

	if *nullable.ProfileDesc != srv.ProfileDesc {
		t.Errorf("Incorrect ProfileDesc after nullable conversion; want: '%s', got: '%s'", srv.ProfileDesc, *nullable.ProfileDesc)
	}

	if *nullable.ProfileID != srv.ProfileID {
		t.Errorf("Incorrect ProfileID after nullable conversion; want: %d, got: %d", srv.ProfileID, *nullable.ProfileID)
	}

	if *nullable.Rack != srv.Rack {
		t.Errorf("Incorrect Rack after nullable conversion; want: '%s', got: '%s'", srv.Rack, *nullable.Rack)
	}

	if *nullable.RevalPending != srv.RevalPending {
		t.Errorf("Incorrect RevalPending after nullable conversion; want: %t, got: %t", srv.RevalPending, *nullable.RevalPending)
	}

	if *nullable.RouterHostName != srv.RouterHostName {
		t.Errorf("Incorrect RouterHostName after nullable conversion; want: '%s', got: '%s'", srv.RouterHostName, *nullable.RouterHostName)
	}

	if *nullable.RouterPortName != srv.RouterPortName {
		t.Errorf("Incorrect RouterPortName after nullable conversion; want: '%s', got: '%s'", srv.RouterPortName, *nullable.RouterPortName)
	}

	if *nullable.Status != srv.Status {
		t.Errorf("Incorrect Status after nullable conversion; want: '%s', got: '%s'", srv.Status, *nullable.Status)
	}

	if *nullable.StatusID != srv.StatusID {
		t.Errorf("Incorrect StatusID after nullable conversion; want: %d, got: %d", srv.StatusID, *nullable.StatusID)
	}

	if *nullable.TCPPort != srv.TCPPort {
		t.Errorf("Incorrect TCPPort after nullable conversion; want: %d, got: %d", srv.TCPPort, *nullable.TCPPort)
	}

	if *nullable.TypeID != srv.TypeID {
		t.Errorf("Incorrect TypeID after nullable conversion; want: %d, got: %d", srv.TypeID, *nullable.TypeID)
	}

	if *nullable.UpdPending != srv.UpdPending {
		t.Errorf("Incorrect UpdPending after nullable conversion; want: %t, got: %t", srv.UpdPending, *nullable.UpdPending)
	}

	if *nullable.XMPPID != srv.XMPPID {
		t.Errorf("Incorrect XMPPID after nullable conversion; want: '%s', got: '%s'", srv.XMPPID, *nullable.XMPPID)
	}

	if *nullable.XMPPPasswd != srv.XMPPPasswd {
		t.Errorf("Incorrect XMPPPasswd after nullable conversion; want: '%s', got: '%s'", srv.XMPPPasswd, *nullable.XMPPPasswd)
	}
}

func TestServerNullableV2_Upgrade(t *testing.T) {
	fqdn := "testFQDN"
	srv := Server{
		Cachegroup:       "testCachegroup",
		CachegroupID:     42,
		CDNID:            43,
		CDNName:          "testCDNName",
		DeliveryServices: map[string][]string{"test": []string{"quest"}},
		DomainName:       "testDomainName",
		FQDN:             &fqdn,
		FqdnTime:         time.Now(),
		GUID:             "testGUID",
		HostName:         "testHostName",
		HTTPSPort:        -1,
		ID:               44,
		ILOIPAddress:     "testILOIPAddress",
		ILOIPGateway:     "testILOIPGateway",
		ILOIPNetmask:     "testILOIPNetmask",
		ILOPassword:      "testILOPassword",
		ILOUsername:      "testILOUsername",
		InterfaceMtu:     2,
		InterfaceName:    "testInterfaceName",
		IP6Address:       "::1/64",
		IP6IsService:     true,
		IP6Gateway:       "::2",
		IPAddress:        "0.0.0.1",
		IPIsService:      false,
		IPGateway:        "0.0.0.2",
		IPNetmask:        "255.255.255.0",
		LastUpdated:      TimeNoMod(Time{Time: time.Now().Add(time.Minute), Valid: true}),
		MgmtIPAddress:    "testMgmtIPAddress",
		MgmtIPGateway:    "testMgmtIPGateway",
		MgmtIPNetmask:    "testMgmtIPNetmask",
		OfflineReason:    "testOfflineReason",
		PhysLocation:     "testPhysLocation",
		PhysLocationID:   45,
		Profile:          "testProfile",
		ProfileDesc:      "testProfileDesc",
		ProfileID:        46,
		Rack:             "testRack",
		RevalPending:     true,
		RouterHostName:   "testRouterHostName",
		RouterPortName:   "testRouterPortName",
		Status:           "testStatus",
		StatusID:         47,
		TCPPort:          -3,
		Type:             "testType",
		TypeID:           48,
		UpdPending:       false,
		XMPPID:           "testXMPPID",
		XMPPPasswd:       "testXMPPasswd",
	}

	// this is so much easier than double the lines to manually construct a
	// nullable v2 server
	nullable := srv.ToNullable()

	noNulls := t.Run("check nullable-converted server for null values", func(t *testing.T) { ensureNoNulls(nullable, t) })
	if !noNulls {
		t.Fatal("Cannot check nullable server field referred-to values - null values found")
	}

	upgraded, err := nullable.Upgrade()
	if err != nil {
		t.Fatalf("Unexpected error upgrading server: %v", err)
	}

	if upgraded.Cachegroup == nil {
		t.Error("upgraded conversion gave nil Cachegroup")
	} else if *upgraded.Cachegroup != *nullable.Cachegroup {
		t.Errorf("Incorrect Cachegroup after upgraded conversion; want: '%s', got: '%s'", *nullable.Cachegroup, *upgraded.Cachegroup)
	}

	if upgraded.CachegroupID == nil {
		t.Error("upgraded conversion gave nil CachegroupID")
	} else if *upgraded.CachegroupID != *nullable.CachegroupID {
		t.Errorf("Incorrect CachegroupID after upgraded conversion; want: %d, got: %d", *nullable.CachegroupID, *upgraded.CachegroupID)
	}

	if upgraded.CDNID == nil {
		t.Error("upgraded conversion gave nil CDNID")
	} else if *upgraded.CDNID != *nullable.CDNID {
		t.Errorf("Incorrect CDNID after upgraded conversion; want: %d, got: %d", *nullable.CDNID, *upgraded.CDNID)
	}

	if upgraded.CDNName == nil {
		t.Error("upgraded conversion gave nil CDNName")
	} else if *upgraded.CDNName != *nullable.CDNName {
		t.Errorf("Incorrect CDNName after upgraded conversion; want: '%s', got: '%s'", *nullable.CDNName, *upgraded.CDNName)
	}

	if upgraded.DeliveryServices == nil {
		t.Error("upgraded conversion gave nil DeliveryServices")
	} else if len(*upgraded.DeliveryServices) != len(*nullable.DeliveryServices) {
		t.Errorf("Incorrect number of DeliveryServices after upgraded conversion; want: %d, got: %d", len(*nullable.DeliveryServices), len(*upgraded.DeliveryServices))
	} else {
		for k, v := range *nullable.DeliveryServices {
			upgradedV, ok := (*upgraded.DeliveryServices)[k]
			if !ok {
				t.Errorf("Missing Delivery Service '%s' after upgraded conversion", k)
				continue
			}
			if len(upgradedV) != len(v) {
				t.Errorf("Delivery Service '%s' has incorrect length after upgraded conversion; want: %d, got: %d", k, len(v), len(upgradedV))
			}
			for i, ds := range v {
				upgradedDS := upgradedV[i]
				if upgradedDS != ds {
					t.Errorf("Incorrect value at position %d in Delivery Service '%s' after upgraded conversion; want: '%s', got: '%s'", i, k, ds, upgradedDS)
				}
			}
		}
	}

	if upgraded.DomainName == nil {
		t.Error("upgraded conversion gave nil DomainName")
	} else if *upgraded.DomainName != *nullable.DomainName {
		t.Errorf("Incorrect DomainName after upgraded conversion; want: '%s', got: '%s'", *nullable.DomainName, *upgraded.DomainName)
	}

	if upgraded.FQDN == nil {
		t.Error("upgraded conversion gave nil FQDN")
	} else if *upgraded.FQDN != fqdn {
		t.Errorf("Incorrect FQDN after upgraded conversion; want: '%s', got: '%s'", fqdn, *upgraded.FQDN)
	}

	if upgraded.FqdnTime != nullable.FqdnTime {
		t.Errorf("Incorrect FqdnTime after upgraded conversion; want: '%s', got: '%s'", nullable.FqdnTime, upgraded.FqdnTime)
	}

	if upgraded.GUID == nil {
		t.Error("upgraded conversion gave nil GUID")
	} else if *upgraded.GUID != *nullable.GUID {
		t.Errorf("Incorrect GUID after upgraded conversion; want: '%s', got: '%s'", *nullable.GUID, *upgraded.GUID)
	}

	if upgraded.HostName == nil {
		t.Error("upgraded conversion gave nil HostName")
	} else if *upgraded.HostName != *nullable.HostName {
		t.Errorf("Incorrect HostName after upgraded conversion; want: '%s', got: '%s'", *nullable.HostName, *upgraded.HostName)
	}

	if upgraded.HTTPSPort == nil {
		t.Error("upgraded conversion gave nil HTTPSPort")
	} else if *upgraded.HTTPSPort != *nullable.HTTPSPort {
		t.Errorf("Incorrect HTTPSPort after upgraded conversion; want: %d, got: %d", *nullable.HTTPSPort, *upgraded.HTTPSPort)
	}

	if upgraded.ID == nil {
		t.Error("upgraded conversion gave nil ID")
	} else if *upgraded.ID != *nullable.ID {
		t.Errorf("Incorrect ID after upgraded conversion; want: %d, got: %d", *nullable.ID, *upgraded.ID)
	}

	if upgraded.ILOIPAddress == nil {
		t.Error("upgraded conversion gave nil ILOIPAddress")
	} else if *upgraded.ILOIPAddress != *nullable.ILOIPAddress {
		t.Errorf("Incorrect ILOIPAddress after upgraded conversion; want: '%s', got: '%s'", *nullable.ILOIPAddress, *upgraded.ILOIPAddress)
	}

	if upgraded.ILOIPGateway == nil {
		t.Error("upgraded conversion gave nil ILOIPGateway")
	} else if *upgraded.ILOIPGateway != *nullable.ILOIPGateway {
		t.Errorf("Incorrect ILOIPGateway after upgraded conversion; want: '%s', got: '%s'", *nullable.ILOIPGateway, *upgraded.ILOIPGateway)
	}

	if upgraded.ILOIPNetmask == nil {
		t.Error("upgraded conversion gave nil ILOIPNetmask")
	} else if *upgraded.ILOIPNetmask != *nullable.ILOIPNetmask {
		t.Errorf("Incorrect ILOIPNetmask after upgraded conversion; want: '%s', got: '%s'", *nullable.ILOIPNetmask, *upgraded.ILOIPNetmask)
	}

	if upgraded.ILOPassword == nil {
		t.Error("upgraded conversion gave nil ILOPassword")
	} else if *upgraded.ILOPassword != *nullable.ILOPassword {
		t.Errorf("Incorrect ILOPassword after upgraded conversion; want: '%s', got: '%s'", *nullable.ILOPassword, *upgraded.ILOPassword)
	}

	if upgraded.ILOUsername == nil {
		t.Error("upgraded conversion gave nil ILOUsername")
	} else if *upgraded.ILOUsername != *nullable.ILOUsername {
		t.Errorf("Incorrect ILOUsername after upgraded conversion; want: '%s', got: '%s'", *nullable.ILOUsername, *upgraded.ILOUsername)
	}

	infLen := len(upgraded.Interfaces)
	if infLen < 1 {
		t.Error("Expected exactly one interface after upgrade, got: 0")
	} else {
		if infLen > 1 {
			t.Errorf("Expected exactly one interface after upgrade, got: %d", infLen)
		}

		inf := upgraded.Interfaces[0]
		if inf.Name != *nullable.InterfaceName {
			t.Errorf("Incorrect interface name after upgrade; want: '%s', got: '%s'", *nullable.InterfaceName, inf.Name)
		}

		if inf.MTU == nil {
			t.Error("Unexpectedly nil Interface MTU after upgrade")
		} else if *inf.MTU != uint64(*nullable.InterfaceMtu) {
			t.Errorf("Incorrect Interface MTU after upgrade; want: %d, got: %d", *nullable.InterfaceMtu, *inf.MTU)
		}

		if !inf.Monitor {
			t.Error("Incorrect Interface Monitor after upgrade; want: true, got: false")
		}

		if inf.MaxBandwidth != nil {
			t.Error("Unexpectedly non-nil Interface MaxBandwidth after upgrade")
		}

		if len(inf.IPAddresses) != 2 {
			t.Errorf("Incorrect number of IP Addresses after upgrade; want: 2, got: %d", len(inf.IPAddresses))
		} else {
			ip := inf.IPAddresses[0]
			cidrIndex := strings.Index(ip.Address, "/")
			addr := ip.Address
			if cidrIndex >= 0 {
				addr = addr[:cidrIndex]
			}

			// TODO: calculate and verify netmask
			if addr == *nullable.IPAddress {
				if ip.Gateway == nil {
					t.Error("Unexpectedly nil IPv4 Gateway after upgrade")
				} else if *ip.Gateway != *nullable.IPGateway {
					t.Errorf("Incorrect IPv4 Gateway after upgrade; want: '%s', got: '%s'", *nullable.IPGateway, *ip.Gateway)
				}

				if ip.ServiceAddress != *nullable.IPIsService {
					t.Errorf("Incorrect IPv4 ServiceAddress value after upgrade; want: %t, got: %t", *nullable.IPIsService, ip.ServiceAddress)
				}

				secondIP := inf.IPAddresses[1]
				if secondIP.Address != *nullable.IP6Address {
					t.Errorf("Incorrect IPv6 Address after upgrade; want: '%s', got: '%s'", *nullable.IP6Address, secondIP.Address)
				} else {
					if secondIP.Gateway == nil {
						t.Error("Unexpectedly nil IPv6 Gateway after upgrade")
					} else if *secondIP.Gateway != *nullable.IP6Gateway {
						t.Errorf("Incorrect IPv6 Gateway after upgrade; want: '%s', got: '%s'", *nullable.IP6Gateway, *secondIP.Gateway)
					}

					if secondIP.ServiceAddress != *nullable.IP6IsService {
						t.Errorf("Incorrect IPv6 ServiceAddress value after upgrade; want: %t, got: %t", *nullable.IP6IsService, secondIP.ServiceAddress)
					}
				}
			} else if ip.Address == *nullable.IP6Address {
				if ip.Gateway == nil {
					t.Error("Unexpectedly nil IPv6 Gateway after upgrade")
				} else if *ip.Gateway != *nullable.IP6Gateway {
					t.Errorf("Incorrect IPv6 Gateway after upgrade; want: '%s', got: '%s'", *nullable.IP6Gateway, *ip.Gateway)
				}

				if ip.ServiceAddress != *nullable.IP6IsService {
					t.Errorf("Incorrect IPv6 ServiceAddress value after upgrade; want: %t, got: %t", *nullable.IP6IsService, ip.ServiceAddress)
				}

				secondIP := inf.IPAddresses[1]
				cidrIndex = strings.Index(secondIP.Address, "/")
				addr = secondIP.Address
				if cidrIndex >= 0 {
					addr = addr[:cidrIndex]
				}
				// TODO: calculate and verify netmask
				if addr != *nullable.IPAddress {
					t.Errorf("Incorrect IPv4 Address after upgrade; want: '%s', got: '%s'", *nullable.IPAddress, secondIP.Address)
				} else {
					if secondIP.Gateway == nil {
						t.Error("Unexpectedly nil IPv4 Gateway after upgrade")
					} else if *secondIP.Gateway != *nullable.IPGateway {
						t.Errorf("Incorrect IPv4 Gateway after upgrade; want: '%s', got: '%s'", *nullable.IPGateway, *secondIP.Gateway)
					}

					if secondIP.ServiceAddress != *nullable.IPIsService {
						t.Errorf("Incorrect IPv4 ServiceAddress value after upgrade; want: %t, got: %t", *nullable.IPIsService, secondIP.ServiceAddress)
					}
				}

			} else {
				t.Errorf("Unknown IP address '%s' found in interface after upgrade", ip.Address)
				ip = inf.IPAddresses[1]
				cidrIndex = strings.Index(ip.Address, "/")
				addr = ip.Address
				if cidrIndex >= 0 {
					addr = addr[:cidrIndex]
				}

				if addr == *nullable.IPAddress {
					t.Error("Missing IPv6 address after upgrade")
					if ip.Gateway == nil {
						t.Error("Unexpectedly nil IPv4 Gateway after upgrade")
					} else if *ip.Gateway != *nullable.IPGateway {
						t.Errorf("Incorrect IPv4 Gateway after upgrade; want: '%s', got: '%s'", *nullable.IPGateway, *ip.Gateway)
					}

					if ip.ServiceAddress != *nullable.IPIsService {
						t.Errorf("Incorrect IPv4 ServiceAddress value after upgrade; want: %t, got: %t", *nullable.IPIsService, ip.ServiceAddress)
					}
				} else if ip.Address == *nullable.IP6Address {
					t.Error("Missing IPv4 address after upgrade")
					if ip.Gateway == nil {
						t.Error("Unexpectedly nil IPv6 Gateway after upgrade")
					} else if *ip.Gateway != *nullable.IP6Gateway {
						t.Errorf("Incorrect IPv6 Gateway after upgrade; want: '%s', got: '%s'", *nullable.IP6Gateway, *ip.Gateway)
					}

					if ip.ServiceAddress != *nullable.IP6IsService {
						t.Errorf("Incorrect IPv6 ServiceAddress value after upgrade; want: %t, got: %t", *nullable.IP6IsService, ip.ServiceAddress)
					}
				} else {
					t.Errorf("Unknown IP address '%s' found in interface after upgrade", ip.Address)
					t.Error("Missing both IPv4 and IPv6 address after upgrade")
				}
			}
		}
	}

	if upgraded.LastUpdated == nil {
		t.Error("upgraded conversion gave nil LastUpdated")
	} else if *upgraded.LastUpdated != *nullable.LastUpdated {
		t.Errorf("Incorrect LastUpdated after upgraded conversion; want: '%s', got: '%s'", *nullable.LastUpdated, *upgraded.LastUpdated)
	}

	if upgraded.MgmtIPAddress == nil {
		t.Error("upgraded conversion gave nil MgmtIPAddress")
	} else if *upgraded.MgmtIPAddress != *nullable.MgmtIPAddress {
		t.Errorf("Incorrect MgmtIPAddress after upgraded conversion; want: '%s', got: '%s'", *nullable.MgmtIPAddress, *upgraded.MgmtIPAddress)
	}

	if upgraded.MgmtIPGateway == nil {
		t.Error("upgraded conversion gave nil MgmtIPGateway")
	} else if *upgraded.MgmtIPGateway != *nullable.MgmtIPGateway {
		t.Errorf("Incorrect MgmtIPGateway after upgraded conversion; want: '%s', got: '%s'", *nullable.MgmtIPGateway, *upgraded.MgmtIPGateway)
	}

	if upgraded.MgmtIPNetmask == nil {
		t.Error("upgraded conversion gave nil MgmtIPNetmask")
	} else if *upgraded.MgmtIPNetmask != *nullable.MgmtIPNetmask {
		t.Errorf("Incorrect MgmtIPNetmask after upgraded conversion; want: '%s', got: '%s'", *nullable.MgmtIPNetmask, *upgraded.MgmtIPNetmask)
	}

	if upgraded.OfflineReason == nil {
		t.Error("upgraded conversion gave nil OfflineReason")
	} else if *upgraded.OfflineReason != *nullable.OfflineReason {
		t.Errorf("Incorrect OfflineReason after upgraded conversion; want: '%s', got: '%s'", *nullable.OfflineReason, *upgraded.OfflineReason)
	}

	if upgraded.PhysLocation == nil {
		t.Error("upgraded conversion gave nil PhysLocation")
	} else if *upgraded.PhysLocation != *nullable.PhysLocation {
		t.Errorf("Incorrect PhysLocation after upgraded conversion; want: '%s', got: '%s'", *nullable.PhysLocation, *upgraded.PhysLocation)
	}

	if upgraded.PhysLocationID == nil {
		t.Error("upgraded conversion gave nil PhysLocationID")
	} else if *upgraded.PhysLocationID != *nullable.PhysLocationID {
		t.Errorf("Incorrect PhysLocationID after upgraded conversion; want: %d, got: %d", *nullable.PhysLocationID, *upgraded.PhysLocationID)
	}

	if upgraded.Profile == nil {
		t.Error("upgraded conversion gave nil Profile")
	} else if *upgraded.Profile != *nullable.Profile {
		t.Errorf("Incorrect Profile after upgraded conversion; want: '%s', got: '%s'", *nullable.Profile, *upgraded.Profile)
	}

	if upgraded.ProfileDesc == nil {
		t.Error("upgraded conversion gave nil ProfileDesc")
	} else if *upgraded.ProfileDesc != *nullable.ProfileDesc {
		t.Errorf("Incorrect ProfileDesc after upgraded conversion; want: '%s', got: '%s'", *nullable.ProfileDesc, *upgraded.ProfileDesc)
	}

	if upgraded.ProfileID == nil {
		t.Error("upgraded conversion gave nil ProfileID")
	} else if *upgraded.ProfileID != *nullable.ProfileID {
		t.Errorf("Incorrect ProfileID after upgraded conversion; want: %d, got: %d", *nullable.ProfileID, *upgraded.ProfileID)
	}

	if upgraded.Rack == nil {
		t.Error("upgraded conversion gave nil Rack")
	} else if *upgraded.Rack != *nullable.Rack {
		t.Errorf("Incorrect Rack after upgraded conversion; want: '%s', got: '%s'", *nullable.Rack, *upgraded.Rack)
	}

	if upgraded.RevalPending == nil {
		t.Error("upgraded conversion gave nil RevalPending")
	} else if *upgraded.RevalPending != *nullable.RevalPending {
		t.Errorf("Incorrect RevalPending after upgraded conversion; want: %t, got: %t", *nullable.RevalPending, *upgraded.RevalPending)
	}

	if upgraded.RouterHostName == nil {
		t.Error("upgraded conversion gave nil RouterHostName")
	} else if *upgraded.RouterHostName != *nullable.RouterHostName {
		t.Errorf("Incorrect RouterHostName after upgraded conversion; want: '%s', got: '%s'", *nullable.RouterHostName, *upgraded.RouterHostName)
	}

	if upgraded.RouterPortName == nil {
		t.Error("upgraded conversion gave nil RouterPortName")
	} else if *upgraded.RouterPortName != *nullable.RouterPortName {
		t.Errorf("Incorrect RouterPortName after upgraded conversion; want: '%s', got: '%s'", *nullable.RouterPortName, *upgraded.RouterPortName)
	}

	if upgraded.Status == nil {
		t.Error("upgraded conversion gave nil Status")
	} else if *upgraded.Status != *nullable.Status {
		t.Errorf("Incorrect Status after upgraded conversion; want: '%s', got: '%s'", *nullable.Status, *upgraded.Status)
	}

	if upgraded.StatusID == nil {
		t.Error("upgraded conversion gave nil StatusID")
	} else if *upgraded.StatusID != *nullable.StatusID {
		t.Errorf("Incorrect StatusID after upgraded conversion; want: %d, got: %d", *nullable.StatusID, *upgraded.StatusID)
	}

	if upgraded.TCPPort == nil {
		t.Error("upgraded conversion gave nil TCPPort")
	} else if *upgraded.TCPPort != *nullable.TCPPort {
		t.Errorf("Incorrect TCPPort after upgraded conversion; want: %d, got: %d", *nullable.TCPPort, *upgraded.TCPPort)
	}

	if upgraded.Type != nullable.Type {
		t.Errorf("Incorrect Type after upgraded conversion; want: '%s', got: '%s'", nullable.Type, upgraded.Type)
	}

	if upgraded.TypeID == nil {
		t.Error("upgraded conversion gave nil TypeID")
	} else if *upgraded.TypeID != *nullable.TypeID {
		t.Errorf("Incorrect TypeID after upgraded conversion; want: %d, got: %d", *nullable.TypeID, *upgraded.TypeID)
	}

	if upgraded.UpdPending == nil {
		t.Error("upgraded conversion gave nil UpdPending")
	} else if *upgraded.UpdPending != *nullable.UpdPending {
		t.Errorf("Incorrect UpdPending after upgraded conversion; want: %t, got: %t", *nullable.UpdPending, *upgraded.UpdPending)
	}

	if upgraded.XMPPID == nil {
		t.Error("upgraded conversion gave nil XMPPID")
	} else if *upgraded.XMPPID != *nullable.XMPPID {
		t.Errorf("Incorrect XMPPID after upgraded conversion; want: '%s', got: '%s'", *nullable.XMPPID, *upgraded.XMPPID)
	}

	if upgraded.XMPPPasswd == nil {
		t.Error("upgraded conversion gave nil XMPPPasswd")
	} else if *upgraded.XMPPPasswd != *nullable.XMPPPasswd {
		t.Errorf("Incorrect XMPPPasswd after upgraded conversion; want: '%s', got: '%s'", *nullable.XMPPPasswd, *upgraded.XMPPPasswd)
	}
}

func ExampleServerV50_UpdatePending() {
	s := ServerV50{
		ConfigApplyTime:  new(time.Time),
		ConfigUpdateTime: new(time.Time),
	}

	*s.ConfigApplyTime = time.Now()
	*s.ConfigUpdateTime = s.ConfigApplyTime.Add(-time.Hour)
	fmt.Println(s.UpdatePending())

	*s.ConfigUpdateTime = s.ConfigUpdateTime.Add(2 * time.Hour)
	fmt.Println(s.UpdatePending())

	// Output: false
	// true
}

func ExampleServerV50_RevalidationPending() {
	s := ServerV50{
		RevalApplyTime:  new(time.Time),
		RevalUpdateTime: new(time.Time),
	}

	*s.RevalApplyTime = time.Now()
	*s.RevalUpdateTime = s.RevalApplyTime.Add(-time.Hour)
	fmt.Println(s.RevalidationPending())

	*s.RevalUpdateTime = s.RevalUpdateTime.Add(2 * time.Hour)
	fmt.Println(s.RevalidationPending())

	// Output: false
	// true
}
