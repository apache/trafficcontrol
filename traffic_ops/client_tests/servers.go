package main

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
import "net/http"
import "time"

import "github.com/apache/trafficcontrol/lib/go-tc"
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"

// For "ILO*", "Mgmt*", IPv6 stuff and other things that need to be empty strings because the database doesn't model them as NULL
var EMPTY_STRING = ""

// These can be constant for all servers
var GATEWAY = "192.0.2.0"
var NETMASK = "255.255.0.0"

var TO_SERVER_DOMAIN_NAME = "cdn.test"
var TO_SERVER_HOSTNAME = "trafficops"
var TO_SERVER_HTTPS_PORT = 443
var TO_SERVER_ID = 1
var TO_SERVER_INTERFACE_NAME = "eth0"
var TO_SERVER_INTERFACE_MTU = 1500
var TO_SERVER_IP = "192.0.2.1"
var TO_SERVER_REVAL_PENDING = false
var TO_SERVER_UPD_PENDING = false
var TO_SERVER_TCP_PORT = 80
var TO_SERVER_XMPPID = "trafficops"

var TODB_SERVER_DOMAIN_NAME = "cdn.test"
var TODB_SERVER_HOSTNAME = "todb"
var TODB_SERVER_HTTPS_PORT = 5432
var TODB_SERVER_ID = 2
var TODB_SERVER_INTERFACE_NAME = "eth0"
var TODB_SERVER_INTERFACE_MTU = 1500
var TODB_SERVER_IP = "192.0.2.2"
var TODB_SERVER_REVAL_PENDING = false
var TODB_SERVER_UPD_PENDING = false
var TODB_SERVER_TCP_PORT = 5432
var TODB_SERVER_XMPPID = "todb"

var EDGE_SERVER_DOMAIN_NAME = "cdn.test"
var EDGE_SERVER_HOSTNAME = "edge"
var EDGE_SERVER_HTTPS_PORT = 443
var EDGE_SERVER_ID = 3
var EDGE_SERVER_INTERFACE_NAME = "eth0"
var EDGE_SERVER_INTERFACE_MTU = 1500
var EDGE_SERVER_IP = "192.0.2.3"
var EDGE_SERVER_REVAL_PENDING = true
var EDGE_SERVER_UPD_PENDING = true
var EDGE_SERVER_TCP_PORT = 80
var EDGE_SERVER_XMPPID = "edge"

var MID_SERVER_DOMAIN_NAME = "cdn.test"
var MID_SERVER_HOSTNAME = "mid"
var MID_SERVER_HTTPS_PORT = 443
var MID_SERVER_ID = 4
var MID_SERVER_INTERFACE_NAME = "eth0"
var MID_SERVER_INTERFACE_MTU = 1500
var MID_SERVER_IP = "192.0.2.4"
var MID_SERVER_REVAL_PENDING = true
var MID_SERVER_UPD_PENDING = true
var MID_SERVER_TCP_PORT = 80
var MID_SERVER_XMPPID = "mid"

var SERVERS = []tc.ServerNullable {
	tc.ServerNullable {
		Cachegroup:       &EDGE_CACHEGROUP,
		CachegroupID:     &EDGE_CACHEGROUP_ID,
		CDNID:            &ALL_CDN_ID,
		CDNName:          &ALL_CDN,
		DeliveryServices: nil,
		DomainName:       &TO_SERVER_DOMAIN_NAME,
		FQDN:             nil,
		FqdnTime:         time.Time{},
		GUID:             nil, // it's because of this little bastard that I have to use "Nullable" - the real API output `null`, so it's gotta be a pointer
		HostName:         &TO_SERVER_HOSTNAME,
		HTTPSPort:        &TO_SERVER_HTTPS_PORT,
		ID:               &TO_SERVER_ID,
		ILOIPAddress:     &EMPTY_STRING,
		ILOIPGateway:     &EMPTY_STRING,
		ILOIPNetmask:     &EMPTY_STRING,
		ILOPassword:      &EMPTY_STRING,
		ILOUsername:      &EMPTY_STRING,
		InterfaceMtu:     &TO_SERVER_INTERFACE_MTU,
		InterfaceName:    &TO_SERVER_INTERFACE_NAME,
		IP6Address:       nil,
		IP6Gateway:       &EMPTY_STRING,
		IPAddress:        &TO_SERVER_IP,
		IPGateway:        &GATEWAY,
		IPNetmask:        &NETMASK,
		LastUpdated:      CURRENT_TIME,
		MgmtIPAddress:    &EMPTY_STRING,
		MgmtIPGateway:    &EMPTY_STRING,
		MgmtIPNetmask:    &EMPTY_STRING,
		OfflineReason:    &EMPTY_STRING,
		PhysLocation:     &(LOCATION.Name),
		PhysLocationID:   &(LOCATION.ID),
		Profile:          &TO_PROFILE_NAME,
		ProfileDesc:      &TO_PROFILE_DESCRIPTION,
		ProfileID:        &TO_PROFILE_ID,
		Rack:             &EMPTY_STRING,
		RevalPending:     &TO_SERVER_REVAL_PENDING,
		RouterHostName:   &EMPTY_STRING,
		RouterPortName:   &EMPTY_STRING,
		Status:           &(STATUS_ONLINE.Name),
		StatusID:         &(STATUS_ONLINE.ID),
		TCPPort:          &TO_SERVER_TCP_PORT,
		Type:             TYPE_TRAFFIC_OPS.Name,
		TypeID:           &(TYPE_TRAFFIC_OPS.ID),
		UpdPending:       &TO_SERVER_UPD_PENDING,
		XMPPID:           &TO_SERVER_XMPPID,
		XMPPPasswd:       &EMPTY_STRING,
	},
	tc.ServerNullable {
		Cachegroup:       &EDGE_CACHEGROUP,
		CachegroupID:     &EDGE_CACHEGROUP_ID,
		CDNID:            &ALL_CDN_ID,
		CDNName:          &ALL_CDN,
		DeliveryServices: nil,
		DomainName:       &TODB_SERVER_DOMAIN_NAME,
		FQDN:             nil,
		FqdnTime:         time.Time{},
		GUID:             nil, // it's because of this little bastard that I have to use "Nullable" - the real API output `null`, so it's gotta be a pointer
		HostName:         &TODB_SERVER_HOSTNAME,
		HTTPSPort:        &TODB_SERVER_HTTPS_PORT,
		ID:               &TODB_SERVER_ID,
		ILOIPAddress:     &EMPTY_STRING,
		ILOIPGateway:     &EMPTY_STRING,
		ILOIPNetmask:     &EMPTY_STRING,
		ILOPassword:      &EMPTY_STRING,
		ILOUsername:      &EMPTY_STRING,
		InterfaceMtu:     &TODB_SERVER_INTERFACE_MTU,
		InterfaceName:    &TODB_SERVER_INTERFACE_NAME,
		IP6Address:       nil,
		IP6Gateway:       &EMPTY_STRING,
		IPAddress:        &TODB_SERVER_IP,
		IPGateway:        &GATEWAY,
		IPNetmask:        &NETMASK,
		LastUpdated:      CURRENT_TIME,
		MgmtIPAddress:    &EMPTY_STRING,
		MgmtIPGateway:    &EMPTY_STRING,
		MgmtIPNetmask:    &EMPTY_STRING,
		OfflineReason:    &EMPTY_STRING,
		PhysLocation:     &(LOCATION.Name),
		PhysLocationID:   &(LOCATION.ID),
		Profile:          &TO_DB_PROFILE_NAME,
		ProfileDesc:      &TO_DB_PROFILE_DESCRIPTION,
		ProfileID:        &TO_DB_PROFILE_ID,
		Rack:             &EMPTY_STRING,
		RevalPending:     &TODB_SERVER_REVAL_PENDING,
		RouterHostName:   &EMPTY_STRING,
		RouterPortName:   &EMPTY_STRING,
		Status:           &(STATUS_ONLINE.Name),
		StatusID:         &(STATUS_ONLINE.ID),
		TCPPort:          &TODB_SERVER_TCP_PORT,
		Type:             TYPE_TRAFFIC_OPS_DB.Name,
		TypeID:           &(TYPE_TRAFFIC_OPS_DB.ID),
		UpdPending:       &TODB_SERVER_UPD_PENDING,
		XMPPID:           &TODB_SERVER_XMPPID,
		XMPPPasswd:       &EMPTY_STRING,
	},
	tc.ServerNullable {
		Cachegroup:       &EDGE_CACHEGROUP,
		CachegroupID:     &EDGE_CACHEGROUP_ID,
		CDNID:            &CDN_ID,
		CDNName:          &CDN,
		DeliveryServices: nil,
		DomainName:       &EDGE_SERVER_DOMAIN_NAME,
		FQDN:             nil,
		FqdnTime:         time.Time{},
		GUID:             nil, // it's because of this little bastard that I have to use "Nullable" - the real API output `null`, so it's gotta be a pointer
		HostName:         &EDGE_SERVER_HOSTNAME,
		HTTPSPort:        &EDGE_SERVER_HTTPS_PORT,
		ID:               &EDGE_SERVER_ID,
		ILOIPAddress:     &EMPTY_STRING,
		ILOIPGateway:     &EMPTY_STRING,
		ILOIPNetmask:     &EMPTY_STRING,
		ILOPassword:      &EMPTY_STRING,
		ILOUsername:      &EMPTY_STRING,
		InterfaceMtu:     &EDGE_SERVER_INTERFACE_MTU,
		InterfaceName:    &EDGE_SERVER_INTERFACE_NAME,
		IP6Address:       nil,
		IP6Gateway:       &EMPTY_STRING,
		IPAddress:        &EDGE_SERVER_IP,
		IPGateway:        &GATEWAY,
		IPNetmask:        &NETMASK,
		LastUpdated:      CURRENT_TIME,
		MgmtIPAddress:    &EMPTY_STRING,
		MgmtIPGateway:    &EMPTY_STRING,
		MgmtIPNetmask:    &EMPTY_STRING,
		OfflineReason:    &EMPTY_STRING,
		PhysLocation:     &(LOCATION.Name),
		PhysLocationID:   &(LOCATION.ID),
		Profile:          &EDGE_PROFILE_NAME,
		ProfileDesc:      &EDGE_PROFILE_DESCRIPTION,
		ProfileID:        &EDGE_PROFILE_ID,
		Rack:             &EMPTY_STRING,
		RevalPending:     &EDGE_SERVER_REVAL_PENDING,
		RouterHostName:   &EMPTY_STRING,
		RouterPortName:   &EMPTY_STRING,
		Status:           &(STATUS_REPORTED.Name),
		StatusID:         &(STATUS_REPORTED.ID),
		TCPPort:          &EDGE_SERVER_TCP_PORT,
		Type:             TYPE_EDGE.Name,
		TypeID:           &(TYPE_EDGE.ID),
		UpdPending:       &EDGE_SERVER_UPD_PENDING,
		XMPPID:           &EDGE_SERVER_XMPPID,
		XMPPPasswd:       &EMPTY_STRING,
	},
	tc.ServerNullable {
		Cachegroup:       &MID_CACHEGROUP,
		CachegroupID:     &MID_CACHEGROUP_ID,
		CDNID:            &CDN_ID,
		CDNName:          &CDN,
		DeliveryServices: nil,
		DomainName:       &MID_SERVER_DOMAIN_NAME,
		FQDN:             nil,
		FqdnTime:         time.Time{},
		GUID:             nil, // it's because of this little bastard that I have to use "Nullable" - the real API output `null`, so it's gotta be a pointer
		HostName:         &MID_SERVER_HOSTNAME,
		HTTPSPort:        &MID_SERVER_HTTPS_PORT,
		ID:               &MID_SERVER_ID,
		ILOIPAddress:     &EMPTY_STRING,
		ILOIPGateway:     &EMPTY_STRING,
		ILOIPNetmask:     &EMPTY_STRING,
		ILOPassword:      &EMPTY_STRING,
		ILOUsername:      &EMPTY_STRING,
		InterfaceMtu:     &MID_SERVER_INTERFACE_MTU,
		InterfaceName:    &MID_SERVER_INTERFACE_NAME,
		IP6Address:       nil,
		IP6Gateway:       &EMPTY_STRING,
		IPAddress:        &MID_SERVER_IP,
		IPGateway:        &GATEWAY,
		IPNetmask:        &NETMASK,
		LastUpdated:      CURRENT_TIME,
		MgmtIPAddress:    &EMPTY_STRING,
		MgmtIPGateway:    &EMPTY_STRING,
		MgmtIPNetmask:    &EMPTY_STRING,
		OfflineReason:    &EMPTY_STRING,
		PhysLocation:     &(LOCATION.Name),
		PhysLocationID:   &(LOCATION.ID),
		Profile:          &MID_PROFILE_NAME,
		ProfileDesc:      &MID_PROFILE_DESCRIPTION,
		ProfileID:        &MID_PROFILE_ID,
		Rack:             &EMPTY_STRING,
		RevalPending:     &MID_SERVER_REVAL_PENDING,
		RouterHostName:   &EMPTY_STRING,
		RouterPortName:   &EMPTY_STRING,
		Status:           &(STATUS_REPORTED.Name),
		StatusID:         &(STATUS_REPORTED.ID),
		TCPPort:          &MID_SERVER_TCP_PORT,
		Type:             TYPE_MID.Name,
		TypeID:           &(TYPE_MID.ID),
		UpdPending:       &MID_SERVER_UPD_PENDING,
		XMPPID:           &MID_SERVER_XMPPID,
		XMPPPasswd:       &EMPTY_STRING,
	},
}

func servers(w http.ResponseWriter, r *http.Request) {
	common(w)
	if (r.Method == http.MethodGet) {
		api.WriteResp(w, r, SERVERS)
	} else {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(`{"alerts":[{"level":"error","text":"This method hasn't yet been implemented."}]}`))
	}
}
