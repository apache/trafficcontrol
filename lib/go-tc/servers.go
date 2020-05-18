package tc

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net"
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

// ServersResponse is a list of Servers as a response.
type ServersResponse struct {
	Response []Server `json:"response"`
	Alerts
}

// ServersDetailResponse is the JSON object returned for a single server.
type ServersDetailResponse struct {
	Response Server `json:"response"`
	Alerts
}

// ServersV1Response is a list of Servers for v1 as a response.
type ServersV1Response struct {
	Response []ServerV1 `json:"response"`
	Alerts
}

type ServerDetailV11 struct {
	ServerDetail
	LegacyInterfaceDetails
}

// ServerDetailV30 is the details for a server for API v3
type ServerDetailV30 struct {
	ServerDetail
	ServerInterfaces *[]ServerInterfaceInfo `json:"interfaces"`
}

// ServersV1DetailResponse is the JSON object returned for a single server for v1.
type ServersV1DetailResponse struct {
	Response []ServerDetailV11 `json:"response"`
	Alerts
}

// ServersV3DetailResponse is the JSON object returned for a single server for v3.
type ServersV3DetailResponse struct {
	Response []ServerDetailV30 `json:"response"`
	Alerts
}

// ServerIpAddress is the data associated with a server's interface's IP address.
type ServerIpAddress struct {
	Address        string  `json:"address" db:"address"`
	Gateway        *string `json:"gateway" db:"gateway"`
	ServiceAddress bool    `json:"serviceAddress" db:"service_address"`
}

// ServerInterfaceInfo is the data associated with a server's interface.
type ServerInterfaceInfo struct {
	IPAddresses  []ServerIpAddress `json:"ipAddresses" db:"ip_addresses"`
	MaxBandwidth *uint64           `json:"maxBandwidth" db:"max_bandwidth"`
	Monitor      bool              `json:"monitor" db:"monitor"`
	MTU          *uint64           `json:"mtu" db:"mtu"`
	Name         string            `json:"name" db:"name"`
}

// Value implements the driver.Valuer interface
// marshals struct to json to pass back as a json.RawMessage
func (sii *ServerInterfaceInfo) Value() (driver.Value, error) {
	b, err := json.Marshal(sii)
	return b, err
}

// Scan implements the sql.Scanner interface
// expects json.RawMessage and unmarshals to a deliveryservice struct
func (sii *ServerInterfaceInfo) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("expected deliveryservice in byte array form; got %T", src)
	}

	return json.Unmarshal([]byte(b), sii)
}

// LegacyInterfaceDetails is the details for interfaces on servers for API v1 and v2.
type LegacyInterfaceDetails struct {
	InterfaceMtu  *int    `json:"interfaceMtu" db:"interface_mtu"`
	InterfaceName *string `json:"interfaceName" db:"interface_name"`
	IP6Address    *string `json:"ip6Address" db:"ip6_address"`
	IP6Gateway    *string `json:"ip6Gateway" db:"ip6_gateway"`
	IPAddress     *string `json:"ipAddress" db:"ip_address"`
	IPGateway     *string `json:"ipGateway" db:"ip_gateway"`
	IPNetmask     *string `json:"ipNetmask" db:"ip_netmask"`
}

func (lid *LegacyInterfaceDetails) ToInterfaces(ipv4IsService, ipv6IsService bool) ([]ServerInterfaceInfo, error) {
	var iface ServerInterfaceInfo
	if lid.InterfaceMtu == nil {
		return nil, errors.New("interfaceMtu is null")
	}
	mtu := uint64(*lid.InterfaceMtu)
	iface.MTU = &mtu

	if lid.InterfaceName == nil {
		return nil, errors.New("interfaceName is null")
	}
	iface.Name = *lid.InterfaceName

	var ips []ServerIpAddress
	if lid.IPAddress != nil && *lid.IPAddress != "" {
		if lid.IPGateway != nil && *lid.IPGateway == "" {
			lid.IPGateway = nil
		}

		ipStr := *lid.IPAddress
		if lid.IPNetmask != nil && *lid.IPNetmask != "" {
			mask := net.ParseIP(*lid.IPNetmask).To4()
			if mask == nil {
				return nil, fmt.Errorf("Failed to parse netmask '%s'", *lid.IPNetmask)
			}
			cidr, _ := net.IPv4Mask(mask[0], mask[1], mask[2], mask[3]).Size()
			ipStr = fmt.Sprintf("%s/%d", ipStr, cidr)
		}

		ips = append(ips, ServerIpAddress{
			Address: ipStr,
			Gateway: lid.IPGateway,
			ServiceAddress: ipv4IsService,
		})
	}

	if lid.IP6Address != nil && *lid.IP6Address != "" {
		if lid.IP6Gateway != nil && *lid.IP6Gateway == "" {
			lid.IP6Gateway = nil
		}
		ips = append(ips, ServerIpAddress{
			Address: *lid.IP6Address,
			Gateway: lid.IP6Gateway,
			ServiceAddress: ipv6IsService,
		})
	}

	iface.IPAddresses = ips
	return []ServerInterfaceInfo{iface}, nil
}

// InterfaceInfoToLegacyInterfaces converts a ServerInterfaceInfo to an
// equivalent LegacyInterfaceDetails structure. It does this by creating the
// IP address fields using the "service" interface's IP addresses. All others
// are discarded, as the legacy format is incapable of representing them.
func InterfaceInfoToLegacyInterfaces(serverInterfaces []ServerInterfaceInfo) (LegacyInterfaceDetails, error) {
	var legacyDetails LegacyInterfaceDetails

	for _, intFace := range serverInterfaces {

		for _, addr := range intFace.IPAddresses {
			if !addr.ServiceAddress {
				continue
			}

			address := addr.Address
			gateway := addr.Gateway

			var parsedIp net.IP
			var mask *net.IPNet
			var err error
			parsedIp, mask, err = net.ParseCIDR(address)
			if err != nil {
				parsedIp = net.ParseIP(address)
				if parsedIp == nil {
					return legacyDetails, fmt.Errorf("Failed to parse '%s' as network or CIDR string: %v", address, err)
				}
			}

			if parsedIp.To4() == nil {
				legacyDetails.IP6Address = &address
				legacyDetails.IP6Gateway = gateway
			} else if mask != nil {
				legacyDetails.IPAddress = util.StrPtr(parsedIp.String())
				legacyDetails.IPGateway = gateway
				legacyDetails.IPNetmask = util.StrPtr(fmt.Sprintf("%d.%d.%d.%d", mask.Mask[0], mask.Mask[1], mask.Mask[2], mask.Mask[3]))
			} else {
				legacyDetails.IPAddress = util.StrPtr(parsedIp.String())
				legacyDetails.IPGateway = gateway
			}

			if intFace.MTU != nil {
				legacyDetails.InterfaceMtu = util.IntPtr(int(*intFace.MTU))
			}

			legacyDetails.InterfaceName = &intFace.Name
		}
	}

	return legacyDetails, nil
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
	IP6IsService     bool                `json:"ip6IsService" db:"ip6_address_is_service"`
	IP6Gateway       string              `json:"ip6Gateway" db:"ip6_gateway"`
	IPAddress        string              `json:"ipAddress" db:"ip_address"`
	IPIsService      bool                `json:"ipIsService" db:"ip_address_is_service"`
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

type ServerV1 struct {
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

// CommonServerProperties is just the collection of properties which are
// shared by all servers across API versions.
type CommonServerProperties struct {
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

// ServerNullableV11 is a server as it appeared in API version 1.1.
type ServerNullableV11 struct {
	LegacyInterfaceDetails
	CommonServerProperties
}

// ServerNullableV2 is a server as it appeared in API v2.
type ServerNullableV2 struct {
	ServerNullableV11
	IPIsService  *bool `json:"ipIsService" db:"ip_address_is_service"`
	IP6IsService *bool `json:"ip6IsService" db:"ip6_address_is_service"`
}

// ServerNullable represents an ATC server, as returned by the TO API.
type ServerNullable struct {
	CommonServerProperties
	Interfaces []ServerInterfaceInfo `json:"interfaces" db:"interfaces"`
}

// ToServerV2 converts the server to an equivalent ServerNullableV2 structure,
// if possible. If the conversion could not be performed, an error is returned.
func (s *ServerNullable) ToServerV2() (ServerNullableV2, error) {
	legacyServer := ServerNullableV2{
		ServerNullableV11: ServerNullableV11{
			CommonServerProperties: s.CommonServerProperties,
		},
		IPIsService: new(bool),
		IP6IsService: new(bool),
	}

	var err error
	legacyServer.LegacyInterfaceDetails, err = InterfaceInfoToLegacyInterfaces(s.Interfaces)
	if err != nil {
		return legacyServer, err
	}

	*legacyServer.IPIsService = legacyServer.LegacyInterfaceDetails.IPAddress != nil && *legacyServer.LegacyInterfaceDetails.IPAddress != ""
	*legacyServer.IP6IsService = legacyServer.LegacyInterfaceDetails.IP6Address != nil && *legacyServer.LegacyInterfaceDetails.IP6Address != ""

	return legacyServer, nil
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

// ServerQueueUpdateRequest encodes the request data for the POST
// servers/{{ID}}/queue_update endpoint.
type ServerQueueUpdateRequest struct {
	Action string `json:"action"`
}

// ServerQueueUpdateResponse decodes the full response with alerts from the POST
// servers/{{ID}}/queue_update endpoint.
type ServerQueueUpdateResponse struct {
	Response ServerQueueUpdate `json:"response"`
	Alerts
}

// ServerQueueUpdate decodes the update data from the POST
// servers/{{ID}}/queue_update endpoint.
type ServerQueueUpdate struct {
	ServerID util.JSONIntStr `json:"serverId"`
	Action   string          `json:"action"`
}
