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

type Server struct {
	Cachegroup     string `json:"cachegroup"`
	CachegroupId   int    `json:"cachegroupId"`
	CdnId          int    `json:"cdnId"`
	CdnName        string `json:"cdnName"`
	DomainName     string `json:"domainName"`
	Guid           string `json:"guid"`
	HostName       string `json:"hostName"`
	HttpsPort      int    `json:"httpsPort"`
	Id             int    `json:"id"`
	IloIpAddress   string `json:"iloIpAddress"`
	IloIpGateway   string `json:"iloIpGateway"`
	IloIpNetmask   string `json:"iloIpNetmask"`
	IloPassword    string `json:"iloPassword"`
	IloUsername    string `json:"iloUsername"`
	InterfaceMtu   string `json:"interfaceMtu"`
	InterfaceName  string `json:"interfaceName"`
	Ip6Address     string `json:"ip6Address"`
	Ip6Gateway     string `json:"ip6Gateway"`
	IpAddress      string `json:"ipAddress"`
	IpGateway      string `json:"ipGateway"`
	IpNetmask      string `json:"ipNetmask"`
	LastUpdated    string `json:"lastUpdated"`
	MgmtIpAddress  string `json:"mgmtIpAddress"`
	MgmtIpGateway  string `json:"mgmtIpGateway"`
	MgmtIpNetmask  string `json:"mgmtIpNetmask"`
	OfflineReason  string `json:"offlineReason"`
	PhysLocation   string `json:"physLocation"`
	PhysLocationId int    `json:"physLocationId"`
	Profile        string `json:"profile"`
	ProfileDesc    string `json:"profileDesc"`
	ProfileId      int    `json:"profileId"`
	Rack           string `json:"rack"`
	RouterHostName string `json:"routerHostName"`
	RouterPortName string `json:"routerPortName"`
	Status         string `json:"status"`
	StatusId       int    `json:"statusId"`
	TcpPort        int    `json:"tcpPort"`
	ServerType     string `json:"type"`
	ServerTypeId   int    `json:"typeId"`
	UpdPending     bool   `json:"updPending"`
	XmppId         string `json:"xmppId"`
	XmppPasswd     string `json:"xmppPasswd"`
}
