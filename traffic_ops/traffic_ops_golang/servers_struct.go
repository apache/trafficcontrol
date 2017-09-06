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

type ServersResponse struct {
	Version  string   `json:"version"`
	Response []Server `json:"response"`
}

type ServerShort struct {
	Cachegroup string `json:"cachegroup" db:"cachegroup"`
	HostName   string `json:"hostName" db:"host_name"`
}
type Server struct {
	Cachegroup     string `json:"cachegroup" db:"cachegroup"`
	CachegroupId   int    `json:"cachegroupId" db:"cachegroup_id"`
	CdnId          int    `json:"cdnId" db:"cdn_id"`
	CdnName        string `json:"cdnName" db:"cdn_name"`
	DomainName     string `json:"domainName" db:"domain_name"`
	Guid           string `json:"guid" db:"guid"`
	HostName       string `json:"hostName" db:"host_name"`
	HttpsPort      int    `json:"httpsPort" db:"https_port"`
	Id             int    `json:"id" db:"id"`
	IloIpAddress   string `json:"iloIpAddress" db:"ilo_ip_address"`
	IloIpGateway   string `json:"iloIpGateway" db:"ilo_ip_gateway"`
	IloIpNetmask   string `json:"iloIpNetmask" db:"ilo_ip_netmask"`
	IloPassword    string `json:"iloPassword" db:"ilo_password"`
	IloUsername    string `json:"iloUsername" db:"ilo_username"`
	InterfaceMtu   string `json:"interfaceMtu" db:"interface_mtu"`
	InterfaceName  string `json:"interfaceName" db:"interface_name"`
	Ip6Address     string `json:"ip6Address" db:"ip6_address"`
	Ip6Gateway     string `json:"ip6Gateway" db:"ip6_gateway"`
	IpAddress      string `json:"ipAddress" db:"ip_address"`
	IpGateway      string `json:"ipGateway" db:"ip_gateway"`
	IpNetmask      string `json:"ipNetmask" db:"ip_netmask"`
	LastUpdated    string `json:"lastUpdated" db:"last_updated"`
	MgmtIpAddress  string `json:"mgmtIpAddress" db:"mgmt_ip_address"`
	MgmtIpGateway  string `json:"mgmtIpGateway" db:"mgmt_ip_gateway"`
	MgmtIpNetmask  string `json:"mgmtIpNetmask" db:"mgmt_ip_netmask"`
	OfflineReason  string `json:"offlineReason" db:"offline_reason"`
	PhysLocation   string `json:"physLocation" db:"phys_location"`
	PhysLocationId int    `json:"physLocationId" db:"phys_location_id"`
	Profile        string `json:"profile" db:"profile"`
	ProfileDesc    string `json:"profileDesc" db:"profile_desc"`
	ProfileId      int    `json:"profileId" db:"profile_id"`
	Rack           string `json:"rack" db:"rack"`
	RouterHostName string `json:"routerHostName" db:"router_host_name"`
	RouterPortName string `json:"routerPortName" db:"router_port_name"`
	Status         string `json:"status" db:"status"`
	StatusId       int    `json:"statusId" db:"status_id"`
	TcpPort        int    `json:"tcpPort" db:"tcp_port"`
	ServerType     string `json:"type" db:"server_type"`
	ServerTypeId   int    `json:"typeId" db:"server_type_id"`
	UpdPending     bool   `json:"updPending" db:"upd_pending"`
	XmppId         string `json:"xmppId" db:"xmpp_id"`
	XmppPasswd     string `json:"xmppPasswd" db:"xmpp_passwd"`
}
