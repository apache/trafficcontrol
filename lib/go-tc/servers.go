package tc

import (
	"time"

	"github.com/apache/trafficcontrol/lib/go-util"
)

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

// ServerDetailResponse is the JSON object returned for a single server
type ServersDetailResponse struct {
	Response Server `json:"response"`
}

type Server struct {
	Cachegroup       string              `json:"cachegroup" db:"cachegroup"`
	CachegroupID     int                 `json:"cachegroupId" db:"cachegroup_id"`
	CDNID            int                 `json:"cdnId" db:"cdn_id"`
	CDNName          string              `json:"cdnName" db:"cdn_name"`
	DeliveryServices map[string][]string `json:"deliveryServices,omitempty"`
	DomainName       string              `json:"domainName" db:"domain_name"`
	FQDN             *string             `json:"fqdn,omitempty"`
	FqdnTime         time.Time           `json:"-"`
	GUID             string              `json:"guid" db:"guid"`
	HostName         string              `json:"hostName" db:"host_name"`
	HTTPSPort        int                 `json:"httpsPort" db:"https_port"`
	ID               int                 `json:"id" db:"id"`
	ILOIPAddress     string              `json:"iloIpAddress" db:"ilo_ip_address"`
	ILOIPGateway     string              `json:"iloIpGateway" db:"ilo_ip_gateway"`
	ILOIPNetmask     string              `json:"iloIpNetmask" db:"ilo_ip_netmask"`
	ILOPassword      string              `json:"iloPassword" db:"ilo_password"`
	ILOUsername      string              `json:"iloUsername" db:"ilo_username"`
	InterfaceMtu     int                 `json:"interfaceMtu" db:"interface_mtu"`
	InterfaceName    string              `json:"interfaceName" db:"interface_name"`
	IP6Address       string              `json:"ip6Address" db:"ip6_address"`
	IP6Gateway       string              `json:"ip6Gateway" db:"ip6_gateway"`
	IPAddress        string              `json:"ipAddress" db:"ip_address"`
	IPGateway        string              `json:"ipGateway" db:"ip_gateway"`
	IPNetmask        string              `json:"ipNetmask" db:"ip_netmask"`
	LastUpdated      TimeNoMod           `json:"lastUpdated" db:"last_updated"`
	MgmtIPAddress    string              `json:"mgmtIpAddress" db:"mgmt_ip_address"`
	MgmtIPGateway    string              `json:"mgmtIpGateway" db:"mgmt_ip_gateway"`
	MgmtIPNetmask    string              `json:"mgmtIpNetmask" db:"mgmt_ip_netmask"`
	OfflineReason    string              `json:"offlineReason" db:"offline_reason"`
	PhysLocation     string              `json:"physLocation" db:"phys_location"`
	PhysLocationID   int                 `json:"physLocationId" db:"phys_location_id"`
	Profile          string              `json:"profile" db:"profile"`
	ProfileDesc      string              `json:"profileDesc" db:"profile_desc"`
	ProfileID        int                 `json:"profileId" db:"profile_id"`
	Rack             string              `json:"rack" db:"rack"`
	RevalPending     bool                `json:"revalPending" db:"reval_pending"`
	RouterHostName   string              `json:"routerHostName" db:"router_host_name"`
	RouterPortName   string              `json:"routerPortName" db:"router_port_name"`
	Status           string              `json:"status" db:"status"`
	StatusID         int                 `json:"statusId" db:"status_id"`
	TCPPort          int                 `json:"tcpPort" db:"tcp_port"`
	Type             string              `json:"type" db:"server_type"`
	TypeID           int                 `json:"typeId" db:"server_type_id"`
	UpdPending       bool                `json:"updPending" db:"upd_pending"`
	XMPPID           string              `json:"xmppId" db:"xmpp_id"`
	XMPPPasswd       string              `json:"xmppPasswd" db:"xmpp_passwd"`
}

type ServerNullable struct {
	Cachegroup       *string              `json:"cachegroup" db:"cachegroup"`
	CachegroupID     *int                 `json:"cachegroupId" db:"cachegroup_id"`
	CDNID            *int                 `json:"cdnId" db:"cdn_id"`
	CDNName          *string              `json:"cdnName" db:"cdn_name"`
	DeliveryServices *map[string][]string `json:"deliveryServices,omitempty"`
	DomainName       *string              `json:"domainName" db:"domain_name"`
	FQDN             *string              `json:"fqdn,omitempty"`
	FqdnTime         time.Time            `json:"-"`
	GUID             *string              `json:"guid" db:"guid"`
	HostName         *string              `json:"hostName" db:"host_name"`
	HTTPSPort        *int                 `json:"httpsPort" db:"https_port"`
	ID               *int                 `json:"id" db:"id"`
	ILOIPAddress     *string              `json:"iloIpAddress" db:"ilo_ip_address"`
	ILOIPGateway     *string              `json:"iloIpGateway" db:"ilo_ip_gateway"`
	ILOIPNetmask     *string              `json:"iloIpNetmask" db:"ilo_ip_netmask"`
	ILOPassword      *string              `json:"iloPassword" db:"ilo_password"`
	ILOUsername      *string              `json:"iloUsername" db:"ilo_username"`
	InterfaceMtu     *int                 `json:"interfaceMtu" db:"interface_mtu"`
	InterfaceName    *string              `json:"interfaceName" db:"interface_name"`
	IP6Address       *string              `json:"ip6Address" db:"ip6_address"`
	IP6Gateway       *string              `json:"ip6Gateway" db:"ip6_gateway"`
	IPAddress        *string              `json:"ipAddress" db:"ip_address"`
	IPGateway        *string              `json:"ipGateway" db:"ip_gateway"`
	IPNetmask        *string              `json:"ipNetmask" db:"ip_netmask"`
	LastUpdated      *TimeNoMod           `json:"lastUpdated" db:"last_updated"`
	MgmtIPAddress    *string              `json:"mgmtIpAddress" db:"mgmt_ip_address"`
	MgmtIPGateway    *string              `json:"mgmtIpGateway" db:"mgmt_ip_gateway"`
	MgmtIPNetmask    *string              `json:"mgmtIpNetmask" db:"mgmt_ip_netmask"`
	OfflineReason    *string              `json:"offlineReason" db:"offline_reason"`
	PhysLocation     *string              `json:"physLocation" db:"phys_location"`
	PhysLocationID   *int                 `json:"physLocationId" db:"phys_location_id"`
	Profile          *string              `json:"profile" db:"profile"`
	ProfileDesc      *string              `json:"profileDesc" db:"profile_desc"`
	ProfileID        *int                 `json:"profileId" db:"profile_id"`
	Rack             *string              `json:"rack" db:"rack"`
	RevalPending     *bool                `json:"revalPending" db:"reval_pending"`
	RouterHostName   *string              `json:"routerHostName" db:"router_host_name"`
	RouterPortName   *string              `json:"routerPortName" db:"router_port_name"`
	Status           *string              `json:"status" db:"status"`
	StatusID         *int                 `json:"statusId" db:"status_id"`
	TCPPort          *int                 `json:"tcpPort" db:"tcp_port"`
	Type             string               `json:"type" db:"server_type"`
	TypeID           *int                 `json:"typeId" db:"server_type_id"`
	UpdPending       *bool                `json:"updPending" db:"upd_pending"`
	XMPPID           *string              `json:"xmppId" db:"xmpp_id"`
	XMPPPasswd       *string              `json:"xmppPasswd" db:"xmpp_passwd"`
}

type ServerUpdateStatus struct {
	HostName           string `json:"host_name"`
	UpdatePending      bool   `json:"upd_pending"`
	RevalPending       bool   `json:"reval_pending"`
	UseRevalPending    bool   `json:"use_reval_pending"`
	HostId             int    `json:"host_id"`
	Status             string `json:"status"`
	ParentPending      bool   `json:"parent_pending"`
	ParentRevalPending bool   `json:"parent_reval_pending"`
}

type ServerPutStatus struct {
	Status        util.JSONNameOrIDStr `json:"status"`
	OfflineReason *string              `json:"offlineReason"`
}

type ServerInfo struct {
	CachegroupID int    `json:"cachegroupId" db:"cachegroup_id"`
	CDNID        int    `json:"cdnId" db:"cdn_id"`
	DomainName   string `json:"domainName" db:"domain_name"`
	HostName     string `json:"hostName" db:"host_name"`
	Type         string `json:"type" db:"server_type"`
}

type ServerDetail struct {
	CacheGroup         *string           `json:"cachegroup" db:"cachegroup"`
	CDNName            *string           `json:"cdnName" db:"cdn_name"`
	DeliveryServiceIDs []int64           `json:"deliveryservices,omitempty"`
	DomainName         *string           `json:"domainName" db:"domain_name"`
	GUID               *string           `json:"guid" db:"guid"`
	HardwareInfo       map[string]string `json:"hardwareInfo"`
	HostName           *string           `json:"hostName" db:"host_name"`
	HTTPSPort          *int              `json:"httpsPort" db:"https_port"`
	ID                 *int              `json:"id" db:"id"`
	ILOIPAddress       *string           `json:"iloIpAddress" db:"ilo_ip_address"`
	ILOIPGateway       *string           `json:"iloIpGateway" db:"ilo_ip_gateway"`
	ILOIPNetmask       *string           `json:"iloIpNetmask" db:"ilo_ip_netmask"`
	ILOPassword        *string           `json:"iloPassword" db:"ilo_password"`
	ILOUsername        *string           `json:"iloUsername" db:"ilo_username"`
	InterfaceMTU       *int              `json:"interfaceMtu" db:"interface_mtu"`
	InterfaceName      *string           `json:"interfaceName" db:"interface_name"`
	IP6Address         *string           `json:"ip6Address" db:"ip6_address"`
	IP6Gateway         *string           `json:"ip6Gateway" db:"ip6_gateway"`
	IPAddress          *string           `json:"ipAddress" db:"ip_address"`
	IPGateway          *string           `json:"ipGateway" db:"ip_gateway"`
	IPNetmask          *string           `json:"ipNetmask" db:"ip_netmask"`
	MgmtIPAddress      *string           `json:"mgmtIpAddress" db:"mgmt_ip_address"`
	MgmtIPGateway      *string           `json:"mgmtIpGateway" db:"mgmt_ip_gateway"`
	MgmtIPNetmask      *string           `json:"mgmtIpNetmask" db:"mgmt_ip_netmask"`
	OfflineReason      *string           `json:"offlineReason" db:"offline_reason"`
	PhysLocation       *string           `json:"physLocation" db:"phys_location"`
	Profile            *string           `json:"profile" db:"profile"`
	ProfileDesc        *string           `json:"profileDesc" db:"profile_desc"`
	Rack               *string           `json:"rack" db:"rack"`
	RouterHostName     *string           `json:"routerHostName" db:"router_host_name"`
	RouterPortName     *string           `json:"routerPortName" db:"router_port_name"`
	Status             *string           `json:"status" db:"status"`
	TCPPort            *int              `json:"tcpPort" db:"tcp_port"`
	Type               string            `json:"type" db:"server_type"`
	XMPPID             *string           `json:"xmppId" db:"xmpp_id"`
	XMPPPasswd         *string           `json:"xmppPasswd" db:"xmpp_passwd"`
}
