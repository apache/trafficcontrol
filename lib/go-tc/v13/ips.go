package v13

import tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"

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

type IPsResponse struct {
	Response []IP `json:"response"`
}

// IPResponse is the JSON object returned for a single secondary IP
type IPResponse struct {
	Response IP `json:"response"`
}

type IP struct {
	ID               int          `json:"id" db:"id"`
	Server           string       `json:"server" db:"server"`
	ServerID         int          `json:"serverId" db:"server_id"`
	Type             string       `json:"type" db:"type"`
	TypeID           int          `json:"typeId" db:"type_id"`
	Interface        string       `json:"interface" db:"interface"`
	InterfaceID      int          `json:"interfaceId" db:"interface_id"`
	IP6Address       string       `json:"ip6Address" db:"ipv6"`
	IP6Gateway       string       `json:"ip6Gateway" db:"ipv6_gateway"`
	IPAddress        string       `json:"ipAddress" db:"ip_address"`
	IPNetmask        string       `json:"ipNetmask" db:"ip_netmask"`
	IPAddressNetmask string       `json:"-" db:"ipv4"`
	IPGateway        string       `json:"ipGateway" db:"ipv4_gateway"`
	LastUpdated      tc.TimeNoMod `json:"lastUpdated" db:"last_updated"`
}

type IPNullable struct {
	ID               *int          `json:"id" db:"id"`
	Server           *string       `json:"server" db:"server"`
	ServerID         *int          `json:"serverId" db:"server_id"`
	Type             *string       `json:"type" db:"type"`
	TypeID           *int          `json:"typeId" db:"type_id"`
	Interface        *string       `json:"interface" db:"interface"`
	InterfaceID      *int          `json:"interfaceId" db:"interface_id"`
	IP6Address       *string       `json:"ip6Address" db:"ipv6"`
	IP6Gateway       *string       `json:"ip6Gateway" db:"ipv6_gateway"`
	IPAddress        *string       `json:"ipAddress" db:"ip_address"`
	IPNetmask        *string       `json:"ipNetmask" db:"ip_netmask"`
	IPAddressNetmask *string       `json:"-" db:"ipv4"`
	IPGateway        *string       `json:"ipGateway" db:"ipv4_gateway"`
	LastUpdated      *tc.TimeNoMod `json:"lastUpdated" db:"last_updated"`
}
