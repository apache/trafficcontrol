package tcstructs

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
	Response []Server `json:"response"`
}

type Server struct {
	Cachegroup     string `json:"cachegroup"`
	CachegroupID   int    `json:"cachegroupId"`
	CdnID          int    `json:"cdnId"`
	CdnName        string `json:"cdnName"`
	DomainName     string `json:"domainName"`
	GUID           string `json:"guid"`
	HostName       string `json:"hostName"`
	HTTPSPort      int    `json:"httpsPort"`
	ID             int    `json:"id"`
	ILOIPAddress   string `json:"iloIpAddress"`
	ILOIPGateway   string `json:"iloIpGateway"`
	ILOIPNetmask   string `json:"iloIpNetmask"`
	ILOPassword    string `json:"iloPassword"`
	ILOUsername    string `json:"iloUsername"`
	InterfaceMTU   int    `json:"interfaceMtu"`
	InterfaceName  string `json:"interfaceName"`
	IP6Address     string `json:"ip6Address"`
	IP6Gateway     string `json:"ip6Gateway"`
	IPAddress      string `json:"ipAddress"`
	IPGateway      string `json:"ipGateway"`
	IPNetmask      string `json:"ipNetmask"`
	LastUpdated    string `json:"lastUpdated"`
	MgmtIPAddress  string `json:"mgmtIpAddress"`
	MgmtIPGateway  string `json:"mgmtIpGateway"`
	MgmtIPNetmask  string `json:"mgmtIpNetmask"`
	OfflineReason  string `json:"offlineReason"`
	PhysLocation   string `json:"physLocation"`
	PhysLocationID int    `json:"physLocationId"`
	Profile        string `json:"profile"`
	ProfileDesc    string `json:"profileDesc"`
	ProfileID      int    `json:"profileId"`
	Rack           string `json:"rack"`
	RevalPending   bool   `json:"revalPending"`
	RouterHostName string `json:"routerHostName"`
	RouterPortName string `json:"routerPortName"`
	Status         string `json:"status"`
	StatusID       int    `json:"statusId"`
	TCPPort        int    `json:"tcpPort"`
	ServerType     string `json:"type"`
	ServerTypeID   int    `json:"typeId"`
	UpdPending     bool   `json:"updPending"`
	XMPPID         string `json:"xmppId"`
	XMPPPasswd     string `json:"xmppPasswd"`
}
