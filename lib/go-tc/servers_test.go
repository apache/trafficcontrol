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
import "testing"
import "time"

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

	if nullable.Cachegroup == nil {
		t.Error("nullable conversion gave nil Cachegroup")
	} else if *nullable.Cachegroup != srv.Cachegroup {
		t.Errorf("Incorrect Cachegroup after nullable conversion; want: '%s', got: '%s'", srv.Cachegroup, *nullable.Cachegroup)
	}

	if nullable.CachegroupID == nil {
		t.Error("nullable conversion gave nil CachegroupID")
	} else if *nullable.CachegroupID != srv.CachegroupID {
		t.Errorf("Incorrect CachegroupID after nullable conversion; want: %d, got: %d", srv.CachegroupID, *nullable.CachegroupID)
	}

	if nullable.CDNID == nil {
		t.Error("nullable conversion gave nil CDNID")
	} else if *nullable.CDNID != srv.CDNID {
		t.Errorf("Incorrect CDNID after nullable conversion; want: %d, got: %d", srv.CDNID, *nullable.CDNID)
	}

	if nullable.CDNName == nil {
		t.Error("nullable conversion gave nil CDNName")
	} else if *nullable.CDNName != srv.CDNName {
		t.Errorf("Incorrect CDNName after nullable conversion; want: '%s', got: '%s'", srv.CDNName, *nullable.CDNName)
	}

	if nullable.DeliveryServices == nil {
		t.Error("nullable conversion gave nil DeliveryServices")
	} else if len(*nullable.DeliveryServices) != len(srv.DeliveryServices) {
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

	if nullable.DomainName == nil {
		t.Error("nullable conversion gave nil DomainName")
	} else if *nullable.DomainName != srv.DomainName {
		t.Errorf("Incorrect DomainName after nullable conversion; want: '%s', got: '%s'", srv.DomainName, *nullable.DomainName)
	}

	if nullable.FQDN == nil {
		t.Error("nullable conversion gave nil FQDN")
	} else if *nullable.FQDN != fqdn {
		t.Errorf("Incorrect FQDN after nullable conversion; want: '%s', got: '%s'", fqdn, *nullable.FQDN)
	}

	if nullable.FqdnTime != srv.FqdnTime {
		t.Errorf("Incorrect FqdnTime after nullable conversion; want: '%s', got: '%s'", srv.FqdnTime, nullable.FqdnTime)
	}

	if nullable.GUID == nil {
		t.Error("nullable conversion gave nil GUID")
	} else if *nullable.GUID != srv.GUID {
		t.Errorf("Incorrect GUID after nullable conversion; want: '%s', got: '%s'", srv.GUID, *nullable.GUID)
	}

	if nullable.HostName == nil {
		t.Error("nullable conversion gave nil HostName")
	} else if *nullable.HostName != srv.HostName {
		t.Errorf("Incorrect HostName after nullable conversion; want: '%s', got: '%s'", srv.HostName, *nullable.HostName)
	}

	if nullable.HTTPSPort == nil {
		t.Error("nullable conversion gave nil HTTPSPort")
	} else if *nullable.HTTPSPort != srv.HTTPSPort {
		t.Errorf("Incorrect HTTPSPort after nullable conversion; want: %d, got: %d", srv.HTTPSPort, *nullable.HTTPSPort)
	}

	if nullable.ID == nil {
		t.Error("nullable conversion gave nil ID")
	} else if *nullable.ID != srv.ID {
		t.Errorf("Incorrect ID after nullable conversion; want: %d, got: %d", srv.ID, *nullable.ID)
	}

	if nullable.ILOIPAddress == nil {
		t.Error("nullable conversion gave nil ILOIPAddress")
	} else if *nullable.ILOIPAddress != srv.ILOIPAddress {
		t.Errorf("Incorrect ILOIPAddress after nullable conversion; want: '%s', got: '%s'", srv.ILOIPAddress, *nullable.ILOIPAddress)
	}

	if nullable.ILOIPGateway == nil {
		t.Error("nullable conversion gave nil ILOIPGateway")
	} else if *nullable.ILOIPGateway != srv.ILOIPGateway {
		t.Errorf("Incorrect ILOIPGateway after nullable conversion; want: '%s', got: '%s'", srv.ILOIPGateway, *nullable.ILOIPGateway)
	}

	if nullable.ILOIPNetmask == nil {
		t.Error("nullable conversion gave nil ILOIPNetmask")
	} else if *nullable.ILOIPNetmask != srv.ILOIPNetmask {
		t.Errorf("Incorrect ILOIPNetmask after nullable conversion; want: '%s', got: '%s'", srv.ILOIPNetmask, *nullable.ILOIPNetmask)
	}

	if nullable.ILOPassword == nil {
		t.Error("nullable conversion gave nil ILOPassword")
	} else if *nullable.ILOPassword != srv.ILOPassword {
		t.Errorf("Incorrect ILOPassword after nullable conversion; want: '%s', got: '%s'", srv.ILOPassword, *nullable.ILOPassword)
	}

	if nullable.ILOUsername == nil {
		t.Error("nullable conversion gave nil ILOUsername")
	} else if *nullable.ILOUsername != srv.ILOUsername {
		t.Errorf("Incorrect ILOUsername after nullable conversion; want: '%s', got: '%s'", srv.ILOUsername, *nullable.ILOUsername)
	}

	if nullable.InterfaceMtu == nil {
		t.Error("nullable conversion gave nil InterfaceMtu")
	} else if *nullable.InterfaceMtu != srv.InterfaceMtu {
		t.Errorf("Incorrect InterfaceMtu after nullable conversion; want: %d, got: %d", srv.InterfaceMtu, *nullable.InterfaceMtu)
	}

	if nullable.InterfaceName == nil {
		t.Error("nullable conversion gave nil InterfaceName")
	} else if *nullable.InterfaceName != srv.InterfaceName {
		t.Errorf("Incorrect InterfaceName after nullable conversion; want: '%s', got: '%s'", srv.InterfaceName, *nullable.InterfaceName)
	}

	if nullable.IP6Address == nil {
		t.Error("nullable conversion gave nil IP6Address")
	} else if *nullable.IP6Address != srv.IP6Address {
		t.Errorf("Incorrect IP6Address after nullable conversion; want: '%s', got: '%s'", srv.IP6Address, *nullable.IP6Address)
	}

	if nullable.IP6IsService == nil {
		t.Error("nullable conversion gave nil IP6IsService")
	} else if *nullable.IP6IsService != srv.IP6IsService {
		t.Errorf("Incorrect IP6IsService after nullable conversion; want: %t, got: %t", srv.IP6IsService, *nullable.IP6IsService)
	}

	if nullable.IP6Gateway == nil {
		t.Error("nullable conversion gave nil IP6Gateway")
	} else if *nullable.IP6Gateway != srv.IP6Gateway {
		t.Errorf("Incorrect IP6Gateway after nullable conversion; want: '%s', got: '%s'", srv.IP6Gateway, *nullable.IP6Gateway)
	}

	if nullable.IPAddress == nil {
		t.Error("nullable conversion gave nil IPAddress")
	} else if *nullable.IPAddress != srv.IPAddress {
		t.Errorf("Incorrect IPAddress after nullable conversion; want: '%s', got: '%s'", srv.IPAddress, *nullable.IPAddress)
	}

	if nullable.IPIsService == nil {
		t.Error("nullable conversion gave nil IPIsService")
	} else if *nullable.IPIsService != srv.IPIsService {
		t.Errorf("Incorrect IPIsService after nullable conversion; want: %t, got: %t", srv.IPIsService, *nullable.IPIsService)
	}

	if nullable.IPGateway == nil {
		t.Error("nullable conversion gave nil IPGateway")
	} else if *nullable.IPGateway != srv.IPGateway {
		t.Errorf("Incorrect IPGateway after nullable conversion; want: '%s', got: '%s'", srv.IPGateway, *nullable.IPGateway)
	}

	if nullable.IPNetmask == nil {
		t.Error("nullable conversion gave nil IPNetmask")
	} else if *nullable.IPNetmask != srv.IPNetmask {
		t.Errorf("Incorrect IPNetmask after nullable conversion; want: '%s', got: '%s'", srv.IPNetmask, *nullable.IPNetmask)
	}

	if nullable.LastUpdated == nil {
		t.Error("nullable conversion gave nil LastUpdated")
	} else if *nullable.LastUpdated != srv.LastUpdated {
		t.Errorf("Incorrect LastUpdated after nullable conversion; want: '%s', got: '%s'", srv.LastUpdated, *nullable.LastUpdated)
	}

	if nullable.MgmtIPAddress == nil {
		t.Error("nullable conversion gave nil MgmtIPAddress")
	} else if *nullable.MgmtIPAddress != srv.MgmtIPAddress {
		t.Errorf("Incorrect MgmtIPAddress after nullable conversion; want: '%s', got: '%s'", srv.MgmtIPAddress, *nullable.MgmtIPAddress)
	}

	if nullable.MgmtIPGateway == nil {
		t.Error("nullable conversion gave nil MgmtIPGateway")
	} else if *nullable.MgmtIPGateway != srv.MgmtIPGateway {
		t.Errorf("Incorrect MgmtIPGateway after nullable conversion; want: '%s', got: '%s'", srv.MgmtIPGateway, *nullable.MgmtIPGateway)
	}

	if nullable.MgmtIPNetmask == nil {
		t.Error("nullable conversion gave nil MgmtIPNetmask")
	} else if *nullable.MgmtIPNetmask != srv.MgmtIPNetmask {
		t.Errorf("Incorrect MgmtIPNetmask after nullable conversion; want: '%s', got: '%s'", srv.MgmtIPNetmask, *nullable.MgmtIPNetmask)
	}

	if nullable.OfflineReason == nil {
		t.Error("nullable conversion gave nil OfflineReason")
	} else if *nullable.OfflineReason != srv.OfflineReason {
		t.Errorf("Incorrect OfflineReason after nullable conversion; want: '%s', got: '%s'", srv.OfflineReason, *nullable.OfflineReason)
	}

	if nullable.PhysLocation == nil {
		t.Error("nullable conversion gave nil PhysLocation")
	} else if *nullable.PhysLocation != srv.PhysLocation {
		t.Errorf("Incorrect PhysLocation after nullable conversion; want: '%s', got: '%s'", srv.PhysLocation, *nullable.PhysLocation)
	}

	if nullable.PhysLocationID == nil {
		t.Error("nullable conversion gave nil PhysLocationID")
	} else if *nullable.PhysLocationID != srv.PhysLocationID {
		t.Errorf("Incorrect PhysLocationID after nullable conversion; want: %d, got: %d", srv.PhysLocationID, *nullable.PhysLocationID)
	}

	if nullable.Profile == nil {
		t.Error("nullable conversion gave nil Profile")
	} else if *nullable.Profile != srv.Profile {
		t.Errorf("Incorrect Profile after nullable conversion; want: '%s', got: '%s'", srv.Profile, *nullable.Profile)
	}

	if nullable.ProfileDesc == nil {
		t.Error("nullable conversion gave nil ProfileDesc")
	} else if *nullable.ProfileDesc != srv.ProfileDesc {
		t.Errorf("Incorrect ProfileDesc after nullable conversion; want: '%s', got: '%s'", srv.ProfileDesc, *nullable.ProfileDesc)
	}

	if nullable.ProfileID == nil {
		t.Error("nullable conversion gave nil ProfileID")
	} else if *nullable.ProfileID != srv.ProfileID {
		t.Errorf("Incorrect ProfileID after nullable conversion; want: %d, got: %d", srv.ProfileID, *nullable.ProfileID)
	}

	if nullable.Rack == nil {
		t.Error("nullable conversion gave nil Rack")
	} else if *nullable.Rack != srv.Rack {
		t.Errorf("Incorrect Rack after nullable conversion; want: '%s', got: '%s'", srv.Rack, *nullable.Rack)
	}

	if nullable.RevalPending == nil {
		t.Error("nullable conversion gave nil RevalPending")
	} else if *nullable.RevalPending != srv.RevalPending {
		t.Errorf("Incorrect RevalPending after nullable conversion; want: %t, got: %t", srv.RevalPending, *nullable.RevalPending)
	}

	if nullable.RouterHostName == nil {
		t.Error("nullable conversion gave nil RouterHostName")
	} else if *nullable.RouterHostName != srv.RouterHostName {
		t.Errorf("Incorrect RouterHostName after nullable conversion; want: '%s', got: '%s'", srv.RouterHostName, *nullable.RouterHostName)
	}

	if nullable.RouterPortName == nil {
		t.Error("nullable conversion gave nil RouterPortName")
	} else if *nullable.RouterPortName != srv.RouterPortName {
		t.Errorf("Incorrect RouterPortName after nullable conversion; want: '%s', got: '%s'", srv.RouterPortName, *nullable.RouterPortName)
	}

	if nullable.Status == nil {
		t.Error("nullable conversion gave nil Status")
	} else if *nullable.Status != srv.Status {
		t.Errorf("Incorrect Status after nullable conversion; want: '%s', got: '%s'", srv.Status, *nullable.Status)
	}

	if nullable.StatusID == nil {
		t.Error("nullable conversion gave nil StatusID")
	} else if *nullable.StatusID != srv.StatusID {
		t.Errorf("Incorrect StatusID after nullable conversion; want: %d, got: %d", srv.StatusID, *nullable.StatusID)
	}

	if nullable.TCPPort == nil {
		t.Error("nullable conversion gave nil TCPPort")
	} else if *nullable.TCPPort != srv.TCPPort {
		t.Errorf("Incorrect TCPPort after nullable conversion; want: %d, got: %d", srv.TCPPort, *nullable.TCPPort)
	}

	if nullable.Type != srv.Type {
		t.Errorf("Incorrect Type after nullable conversion; want: '%s', got: '%s'", srv.Type, nullable.Type)
	}

	if nullable.TypeID == nil {
		t.Error("nullable conversion gave nil TypeID")
	} else if *nullable.TypeID != srv.TypeID {
		t.Errorf("Incorrect TypeID after nullable conversion; want: %d, got: %d", srv.TypeID, *nullable.TypeID)
	}

	if nullable.UpdPending == nil {
		t.Error("nullable conversion gave nil UpdPending")
	} else if *nullable.UpdPending != srv.UpdPending {
		t.Errorf("Incorrect UpdPending after nullable conversion; want: %t, got: %t", srv.UpdPending, *nullable.UpdPending)
	}

	if nullable.XMPPID == nil {
		t.Error("nullable conversion gave nil XMPPID")
	} else if *nullable.XMPPID != srv.XMPPID {
		t.Errorf("Incorrect XMPPID after nullable conversion; want: '%s', got: '%s'", srv.XMPPID, *nullable.XMPPID)
	}

	if nullable.XMPPPasswd == nil {
		t.Error("nullable conversion gave nil XMPPPasswd")
	} else if *nullable.XMPPPasswd != srv.XMPPPasswd {
		t.Errorf("Incorrect XMPPPasswd after nullable conversion; want: '%s', got: '%s'", srv.XMPPPasswd, *nullable.XMPPPasswd)
	}
}
