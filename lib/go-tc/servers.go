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
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

// ServersV4Response is the format of a response to a GET request for API v4.x /servers.
type ServersV4Response struct {
	Response []ServerV40 `json:"response"`
	Summary  struct {
		Count uint64 `json:"count"`
	} `json:"summary"`
	Alerts
}

// ServersV3Response is the format of a response to a GET request for /servers.
type ServersV3Response struct {
	Response []ServerV30 `json:"response"`
	Summary  struct {
		Count uint64 `json:"count"`
	} `json:"summary"`
	Alerts
}

// ServersResponse is a list of Servers as a response to an API v2 request.
// This can't change because it will break ORT. Unfortunately.
type ServersResponse struct {
	Response []Server `json:"response"`
	Alerts
}

// ServersDetailResponse is the JSON object returned for a single server.
type ServersDetailResponse struct {
	Response Server `json:"response"`
	Alerts
}

// ServerDetailV11 is the type of each entry in the `response` array property
// of responses from Traffic Ops to GET requests made to its /servers/details
// API endpoint in API version 2.0.
//
// The reason it's named with "V11" is because it was originally used to model
// the response in API version 1.1, and the structure simply hasn't changed
// between then and 2.0.
type ServerDetailV11 struct {
	ServerDetail
	LegacyInterfaceDetails
	RouterHostName *string `json:"routerHostName" db:"router_host_name"`
	RouterPortName *string `json:"routerPortName" db:"router_port_name"`
}

// ServerDetailV30 is the details for a server for API v3.
type ServerDetailV30 struct {
	ServerDetail
	ServerInterfaces *[]ServerInterfaceInfo `json:"interfaces"`
	RouterHostName   *string                `json:"routerHostName" db:"router_host_name"`
	RouterPortName   *string                `json:"routerPortName" db:"router_port_name"`
}

// ServerDetailV40 is the details for a server for API v4.
type ServerDetailV40 struct {
	CacheGroup         *string                  `json:"cachegroup" db:"cachegroup"`
	CDNName            *string                  `json:"cdnName" db:"cdn_name"`
	DeliveryServiceIDs []int64                  `json:"deliveryservices,omitempty"`
	DomainName         *string                  `json:"domainName" db:"domain_name"`
	GUID               *string                  `json:"guid" db:"guid"`
	HardwareInfo       map[string]string        `json:"hardwareInfo"`
	HostName           *string                  `json:"hostName" db:"host_name"`
	HTTPSPort          *int                     `json:"httpsPort" db:"https_port"`
	ID                 *int                     `json:"id" db:"id"`
	ILOIPAddress       *string                  `json:"iloIpAddress" db:"ilo_ip_address"`
	ILOIPGateway       *string                  `json:"iloIpGateway" db:"ilo_ip_gateway"`
	ILOIPNetmask       *string                  `json:"iloIpNetmask" db:"ilo_ip_netmask"`
	ILOPassword        *string                  `json:"iloPassword" db:"ilo_password"`
	ILOUsername        *string                  `json:"iloUsername" db:"ilo_username"`
	MgmtIPAddress      *string                  `json:"mgmtIpAddress" db:"mgmt_ip_address"`
	MgmtIPGateway      *string                  `json:"mgmtIpGateway" db:"mgmt_ip_gateway"`
	MgmtIPNetmask      *string                  `json:"mgmtIpNetmask" db:"mgmt_ip_netmask"`
	OfflineReason      *string                  `json:"offlineReason" db:"offline_reason"`
	PhysLocation       *string                  `json:"physLocation" db:"phys_location"`
	ProfileNames       []string                 `json:"profileNames" db:"profile_name"`
	Rack               *string                  `json:"rack" db:"rack"`
	Status             *string                  `json:"status" db:"status"`
	TCPPort            *int                     `json:"tcpPort" db:"tcp_port"`
	Type               string                   `json:"type" db:"server_type"`
	XMPPID             *string                  `json:"xmppId" db:"xmpp_id"`
	XMPPPasswd         *string                  `json:"xmppPasswd" db:"xmpp_passwd"`
	ServerInterfaces   []ServerInterfaceInfoV40 `json:"interfaces"`
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

// ServersV4DetailResponse is the JSON object returned for a single server for v4.
type ServersV4DetailResponse struct {
	Response []ServerDetailV40 `json:"response"`
	Alerts
}

// ServerIPAddress is the data associated with a server's interface's IP address.
type ServerIPAddress struct {
	Address        string  `json:"address" db:"address"`
	Gateway        *string `json:"gateway" db:"gateway"`
	ServiceAddress bool    `json:"serviceAddress" db:"service_address"`
}

// Copy creates a deep copy of the IP address.
func (ip ServerIPAddress) Copy() ServerIPAddress {
	return ServerIPAddress{
		Address:        ip.Address,
		Gateway:        util.CopyIfNotNil(ip.Gateway),
		ServiceAddress: ip.ServiceAddress,
	}
}

// ServerInterfaceInfo is the data associated with a server's interface.
type ServerInterfaceInfo struct {
	IPAddresses  []ServerIPAddress `json:"ipAddresses" db:"ip_addresses"`
	MaxBandwidth *uint64           `json:"maxBandwidth" db:"max_bandwidth"`
	Monitor      bool              `json:"monitor" db:"monitor"`
	MTU          *uint64           `json:"mtu" db:"mtu"`
	Name         string            `json:"name" db:"name"`
}

// Copy creates a deep copy of the Server Interface.
func (inf ServerInterfaceInfo) Copy() ServerInterfaceInfo {
	newInf := ServerInterfaceInfo{
		IPAddresses:  make([]ServerIPAddress, len(inf.IPAddresses)),
		MaxBandwidth: util.CopyIfNotNil(inf.MaxBandwidth),
		Monitor:      inf.Monitor,
		MTU:          util.CopyIfNotNil(inf.MTU),
		Name:         inf.Name,
	}
	for i, ip := range inf.IPAddresses {
		newInf.IPAddresses[i] = ip.Copy()
	}

	return newInf
}

// ServerInterfaceInfoV40 is the data associated with a V40 server's interface.
type ServerInterfaceInfoV40 struct {
	ServerInterfaceInfo
	RouterHostName string `json:"routerHostName" db:"router_host_name"`
	RouterPortName string `json:"routerPortName" db:"router_port_name"`
}

// Copy creates a deep copy of the Server Interface.
func (inf ServerInterfaceInfoV40) Copy() ServerInterfaceInfoV40 {
	return ServerInterfaceInfoV40{
		ServerInterfaceInfo: inf.ServerInterfaceInfo.Copy(),
		RouterHostName:      inf.RouterHostName,
		RouterPortName:      inf.RouterPortName,
	}
}

// GetDefaultAddress returns the IPv4 and IPv6 service addresses of the
// interface - without any subnet/masking that may or may not be present on said
// address(es).
func (i *ServerInterfaceInfo) GetDefaultAddress() (string, string) {
	ipv4, ipv6 := i.GetDefaultAddressOrCIDR()
	address, _, err := net.ParseCIDR(ipv4)
	if address != nil && err == nil {
		ipv4 = address.String()
	}
	address, _, err = net.ParseCIDR(ipv6)
	if address != nil && err == nil {
		ipv6 = address.String()
	}
	return ipv4, ipv6
}

// GetDefaultAddressOrCIDR returns the IPv4 and IPv6 service addresses of the interface,
// including a subnet, if one exists.
func (i *ServerInterfaceInfo) GetDefaultAddressOrCIDR() (string, string) {
	var ipv4, ipv6 string
	var err error
	for _, ip := range i.IPAddresses {
		if ip.ServiceAddress {
			address := net.ParseIP(ip.Address)
			if address == nil {
				address, _, err = net.ParseCIDR(ip.Address)
				if err != nil || address == nil {
					continue
				}
			}
			if address.To4() != nil {
				ipv4 = ip.Address
			} else if address.To16() != nil {
				ipv6 = ip.Address
			}

			if ipv4 != "" && ipv6 != "" {
				break
			}
		}
	}
	return ipv4, ipv6
}

// Value implements the driver.Valuer interface
// marshals struct to json to pass back as a json.RawMessage.
func (i *ServerInterfaceInfo) Value() (driver.Value, error) {
	b, err := json.Marshal(i)
	return b, err
}

// Scan implements the sql.Scanner interface.
//
// This expects src to be a json.RawMessage and unmarshals it into the
// ServerInterfaceInfo.
func (i *ServerInterfaceInfo) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("expected deliveryservice in byte array form; got %T", src)
	}

	return json.Unmarshal(b, i)
}

// LegacyInterfaceDetails is the details for interfaces on servers for API v2.
//
// Deprecated: Traffic Ops API version 2 is deprecated, upgrade to Server
// representations that support interfaces (i.e. ServerInterfaceInfoV40
// slices).
type LegacyInterfaceDetails struct {
	InterfaceMtu  *int    `json:"interfaceMtu" db:"interface_mtu"`
	InterfaceName *string `json:"interfaceName" db:"interface_name"`
	IP6Address    *string `json:"ip6Address" db:"ip6_address"`
	IP6Gateway    *string `json:"ip6Gateway" db:"ip6_gateway"`
	IPAddress     *string `json:"ipAddress" db:"ip_address"`
	IPGateway     *string `json:"ipGateway" db:"ip_gateway"`
	IPNetmask     *string `json:"ipNetmask" db:"ip_netmask"`
}

// ToInterfaces converts a LegacyInterfaceDetails to a slice of
// ServerInterfaceInfo structures. Only one interface is expected and will be marked for monitoring.
// It will generate service addresses according to the passed indicators
// for each address family.
//
// Deprecated: LegacyInterfaceDetails is deprecated, and this will be removed
// with it.
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

	// default to true since there should only be one interface from legacy API versions
	// if Monitor is false on all interfaces, then TM will see the server as unhealthy
	iface.Monitor = true

	var ips []ServerIPAddress
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

		ips = append(ips, ServerIPAddress{
			Address:        ipStr,
			Gateway:        lid.IPGateway,
			ServiceAddress: ipv4IsService,
		})
	}

	if lid.IP6Address != nil && *lid.IP6Address != "" {
		if lid.IP6Gateway != nil && *lid.IP6Gateway == "" {
			lid.IP6Gateway = nil
		}
		ips = append(ips, ServerIPAddress{
			Address:        *lid.IP6Address,
			Gateway:        lid.IP6Gateway,
			ServiceAddress: ipv6IsService,
		})
	}

	iface.IPAddresses = ips
	return []ServerInterfaceInfo{iface}, nil
}

// ToInterfacesV4 upgrades server interfaces from their APIv3 representation to
// an APIv4 representation.
//
// The passed routerName and routerPort can be nil, which will leave the
// upgradeded interfaces' RouterHostName and RouterPortName unset, or they can
// be pointers to string values to which ALL of the interfaces will have their
// RouterHostName and RouterPortName set.
func ToInterfacesV4(oldInterfaces []ServerInterfaceInfo, routerName, routerPort *string) ([]ServerInterfaceInfoV40, error) {
	v4Interfaces := make([]ServerInterfaceInfoV40, 0)
	var v4Int ServerInterfaceInfoV40
	for _, i := range oldInterfaces {
		v4Int.ServerInterfaceInfo = i
		if routerName != nil {
			v4Int.RouterHostName = *routerName
		}
		if routerPort != nil {
			v4Int.RouterPortName = *routerPort
		}
		v4Interfaces = append(v4Interfaces, v4Int)
	}
	return v4Interfaces, nil
}

// ToInterfacesV4 converts a LegacyInterfaceDetails to a slice of
// ServerInterfaceInfoV40 structures.
//
// Only one interface is expected and will be marked for monitoring. This will
// generate service addresses according to the passed indicators for each
// address family.
//
// The passed routerName and routerPort can be nil, which will leave the
// upgradeded interfaces' RouterHostName and RouterPortName unset, or they can
// be pointers to string values to which ALL of the interfaces will have their
// RouterHostName and RouterPortName set.
//
// Deprecated: LegacyInterfaceDetails is deprecated, and this will be removed
// with it.
func (lid *LegacyInterfaceDetails) ToInterfacesV4(ipv4IsService, ipv6IsService bool, routerName, routerPort *string) ([]ServerInterfaceInfoV40, error) {
	var iface ServerInterfaceInfoV40
	if lid.InterfaceMtu == nil {
		return nil, errors.New("interfaceMtu is null")
	}
	mtu := uint64(*lid.InterfaceMtu)
	iface.MTU = &mtu

	if lid.InterfaceName == nil {
		return nil, errors.New("interfaceName is null")
	}
	iface.Name = *lid.InterfaceName

	// default to true since there should only be one interface from legacy API versions
	// if Monitor is false on all interfaces, then TM will see the server as unhealthy
	iface.Monitor = true

	var ips []ServerIPAddress
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

		ips = append(ips, ServerIPAddress{
			Address:        ipStr,
			Gateway:        lid.IPGateway,
			ServiceAddress: ipv4IsService,
		})
	}

	if lid.IP6Address != nil && *lid.IP6Address != "" {
		if lid.IP6Gateway != nil && *lid.IP6Gateway == "" {
			lid.IP6Gateway = nil
		}
		ips = append(ips, ServerIPAddress{
			Address:        *lid.IP6Address,
			Gateway:        lid.IP6Gateway,
			ServiceAddress: ipv6IsService,
		})
	}

	iface.IPAddresses = ips
	if routerName != nil {
		iface.RouterHostName = *routerName
	}
	if routerPort != nil {
		iface.RouterPortName = *routerPort
	}
	return []ServerInterfaceInfoV40{iface}, nil
}

// String implements the fmt.Stringer interface.
func (lid LegacyInterfaceDetails) String() string {
	var b strings.Builder
	b.Write([]byte("LegacyInterfaceDetails(InterfaceMtu="))

	if lid.InterfaceMtu == nil {
		b.Write([]byte("nil"))
	} else {
		b.WriteString(strconv.FormatInt(int64(*lid.InterfaceMtu), 10))
	}

	b.Write([]byte(", InterfaceName="))
	if lid.InterfaceName != nil {
		b.WriteRune('\'')
		b.WriteString(*lid.InterfaceName)
		b.WriteRune('\'')
	} else {
		b.Write([]byte("nil"))
	}

	b.Write([]byte(", IP6Address="))
	if lid.IP6Address != nil {
		b.WriteRune('\'')
		b.WriteString(*lid.IP6Address)
		b.WriteRune('\'')
	} else {
		b.Write([]byte("nil"))
	}

	b.Write([]byte(", IP6Gateway="))
	if lid.IP6Gateway != nil {
		b.WriteRune('\'')
		b.WriteString(*lid.IP6Gateway)
		b.WriteRune('\'')
	} else {
		b.Write([]byte("nil"))
	}

	b.Write([]byte(", IPAddress="))
	if lid.IPAddress != nil {
		b.WriteRune('\'')
		b.WriteString(*lid.IPAddress)
		b.WriteRune('\'')
	} else {
		b.Write([]byte("nil"))
	}

	b.Write([]byte(", IPGateway="))
	if lid.IPGateway != nil {
		b.WriteRune('\'')
		b.WriteString(*lid.IPGateway)
		b.WriteRune('\'')
	} else {
		b.Write([]byte("nil"))
	}

	b.Write([]byte(", IPNetmask="))
	if lid.IPNetmask != nil {
		b.WriteRune('\'')
		b.WriteString(*lid.IPNetmask)
		b.WriteRune('\'')
	} else {
		b.Write([]byte("nil"))
	}

	b.WriteRune(')')

	return b.String()
}

// V4InterfaceInfoToV3Interfaces downgrades a set of ServerInterfaceInfoV40s to
// ServerInterfaceInfos.
func V4InterfaceInfoToV3Interfaces(serverInterfaces []ServerInterfaceInfoV40) ([]ServerInterfaceInfo, error) {
	var interfaces []ServerInterfaceInfo

	for _, intFace := range serverInterfaces {
		var interfaceV3 ServerInterfaceInfo
		interfaceV3.IPAddresses = intFace.IPAddresses
		interfaceV3.Monitor = intFace.Monitor
		interfaceV3.MaxBandwidth = intFace.MaxBandwidth
		interfaceV3.MTU = intFace.MTU
		interfaceV3.Name = intFace.Name
		interfaces = append(interfaces, interfaceV3)
	}

	return interfaces, nil
}

// V4InterfaceInfoToLegacyInterfaces downgrades a set of
// ServerInterfaceInfoV40s to a LegacyInterfaceDetails.
//
// Deprecated: LegacyInterfaceDetails is deprecated, and this will be removed
// with it.
func V4InterfaceInfoToLegacyInterfaces(serverInterfaces []ServerInterfaceInfoV40) (LegacyInterfaceDetails, error) {
	var legacyDetails LegacyInterfaceDetails

	for _, intFace := range serverInterfaces {

		foundServiceInterface := false

		for _, addr := range intFace.IPAddresses {
			if !addr.ServiceAddress {
				continue
			}

			foundServiceInterface = true

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
				legacyDetails.IPNetmask = new(string)
			}

			if intFace.MTU != nil {
				legacyDetails.InterfaceMtu = util.IntPtr(int(*intFace.MTU))
			}

			// This should no longer matter now that short-circuiting is better,
			// but this temporary variable is necessary because the 'intFace'
			// variable is referential, so taking '&intFace.Name' would cause
			// problems when intFace is reassigned.
			name := intFace.Name
			legacyDetails.InterfaceName = &name

			// we can jump out here since servers can only legally have one
			// IPv4 and one IPv6 service address
			if legacyDetails.IPAddress != nil && *legacyDetails.IPAddress != "" && legacyDetails.IP6Address != nil && *legacyDetails.IP6Address != "" {
				break
			}
		}

		if foundServiceInterface {
			return legacyDetails, nil
		}
	}

	return legacyDetails, errors.New("no service addresses found")
}

// InterfaceInfoToLegacyInterfaces converts a ServerInterfaceInfo to an
// equivalent LegacyInterfaceDetails structure. It does this by creating the
// IP address fields using the "service" interface's IP addresses. All others
// are discarded, as the legacy format is incapable of representing them.
//
// Deprecated: LegacyInterfaceDetails is deprecated, and this will be removed
// with it.
func InterfaceInfoToLegacyInterfaces(serverInterfaces []ServerInterfaceInfo) (LegacyInterfaceDetails, error) {
	var legacyDetails LegacyInterfaceDetails

	for _, intFace := range serverInterfaces {

		foundServiceInterface := false

		for _, addr := range intFace.IPAddresses {
			if !addr.ServiceAddress {
				continue
			}

			foundServiceInterface = true

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
				legacyDetails.IPNetmask = new(string)
			}

			if intFace.MTU != nil {
				legacyDetails.InterfaceMtu = util.IntPtr(int(*intFace.MTU))
			}

			// This should no longer matter now that short-circuiting is better,
			// but this temporary variable is necessary because the 'intFace'
			// variable is referential, so taking '&intFace.Name' would cause
			// problems when intFace is reassigned.
			name := intFace.Name
			legacyDetails.InterfaceName = &name

			// we can jump out here since servers can only legally have one
			// IPv4 and one IPv6 service address
			if legacyDetails.IPAddress != nil && *legacyDetails.IPAddress != "" && legacyDetails.IP6Address != nil && *legacyDetails.IP6Address != "" {
				break
			}
		}

		if foundServiceInterface {
			return legacyDetails, nil
		}
	}

	return legacyDetails, errors.New("no service addresses found")
}

// Server is a non-"nullable" representation of a Server as it appeared in API
// version 2.0
//
// Deprecated: Please use versioned and nullable structures from now on.
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
	// Deprecated: In the future, management interfaces must be configured as
	// interfaces within the Interfaces of the server, not separately on these
	// properties.
	MgmtIPAddress *string `json:"mgmtIpAddress" db:"mgmt_ip_address"`
	// Deprecated: In the future, management interfaces must be configured as
	// interfaces within the Interfaces of the server, not separately on these
	// properties.
	MgmtIPGateway *string `json:"mgmtIpGateway" db:"mgmt_ip_gateway"`
	// Deprecated: In the future, management interfaces must be configured as
	// interfaces within the Interfaces of the server, not separately on these
	// properties.
	MgmtIPNetmask  *string `json:"mgmtIpNetmask" db:"mgmt_ip_netmask"`
	OfflineReason  *string `json:"offlineReason" db:"offline_reason"`
	PhysLocation   *string `json:"physLocation" db:"phys_location"`
	PhysLocationID *int    `json:"physLocationId" db:"phys_location_id"`
	Profile        *string `json:"profile" db:"profile"`
	// Deprecated: In API versions 4 and later, Profile descriptions must be
	// taken from the Profiles themselves, and Servers only contain identifying
	// information for their Profiles.
	ProfileDesc *string `json:"profileDesc" db:"profile_desc"`
	// Deprecated: In API versions 4 and later, Servers identify their Profiles
	// by Name, not ID.
	ProfileID *int    `json:"profileId" db:"profile_id"`
	Rack      *string `json:"rack" db:"rack"`
	// Deprecated: In APIv5 and later, this extraneous field is not calculated
	// by Traffic Ops; the information is available by comparing RevalUpdateTime
	// to RevalApplyTime.
	RevalPending *bool   `json:"revalPending" db:"reval_pending"`
	Status       *string `json:"status" db:"status"`
	StatusID     *int    `json:"statusId" db:"status_id"`
	TCPPort      *int    `json:"tcpPort" db:"tcp_port"`
	Type         string  `json:"type" db:"server_type"`
	TypeID       *int    `json:"typeId" db:"server_type_id"`
	// Deprecated: In APIv5 and later, this extraneous field is not calculated
	// by Traffic Ops; the information is available by comparing
	// ConfigUpdateTime to ConfigApplyTime.
	UpdPending *bool   `json:"updPending" db:"upd_pending"`
	XMPPID     *string `json:"xmppId" db:"xmpp_id"`
	XMPPPasswd *string `json:"xmppPasswd" db:"xmpp_passwd"`
}

// ServerNullableV11 is a server as it appeared in API version 1.1.
type ServerNullableV11 struct {
	LegacyInterfaceDetails
	CommonServerProperties
	RouterHostName *string `json:"routerHostName" db:"router_host_name"`
	RouterPortName *string `json:"routerPortName" db:"router_port_name"`
}

// ServerNullableV2 is a server as it appeared in API v2.
//
// Deprecated: Traffic Ops API version 2 is deprecated, new code should use
// ServerV40 or newer structures.
type ServerNullableV2 struct {
	ServerNullableV11
	IPIsService  *bool `json:"ipIsService" db:"ip_address_is_service"`
	IP6IsService *bool `json:"ip6IsService" db:"ip6_address_is_service"`
}

// ToNullable converts the Server to an equivalent, "nullable" structure.
//
// Note that "zero" values (e.g. the empty string "") are NOT coerced to actual
// null values. In particular, the only fields that will possibly be nil are
// FQDN - if the original server had a nil FQDN - and DeliveryServices - which
// will actually be a pointer to a nil map if the original server had a nil
// DeliveryServices map.
// Further note that this makes "shallow" copies of member properties; if
// reference types (map, slice, pointer etc.) are altered on the original after
// conversion, the changes WILL affect the nullable copy.
//
// Deprecated: Traffic Ops API version 2 is deprecated, new code should use
// ServerV40 or newer structures.
func (s Server) ToNullable() ServerNullableV2 {
	return ServerNullableV2{
		ServerNullableV11: ServerNullableV11{
			CommonServerProperties: CommonServerProperties{
				Cachegroup:       &s.Cachegroup,
				CachegroupID:     &s.CachegroupID,
				CDNID:            &s.CDNID,
				CDNName:          &s.CDNName,
				DeliveryServices: &s.DeliveryServices,
				DomainName:       &s.DomainName,
				FQDN:             s.FQDN,
				FqdnTime:         s.FqdnTime,
				GUID:             &s.GUID,
				HostName:         &s.HostName,
				HTTPSPort:        &s.HTTPSPort,
				ID:               &s.ID,
				ILOIPAddress:     &s.ILOIPAddress,
				ILOIPGateway:     &s.ILOIPGateway,
				ILOIPNetmask:     &s.ILOIPNetmask,
				ILOPassword:      &s.ILOPassword,
				ILOUsername:      &s.ILOUsername,
				LastUpdated:      &s.LastUpdated,
				MgmtIPAddress:    &s.MgmtIPAddress,
				MgmtIPGateway:    &s.MgmtIPGateway,
				MgmtIPNetmask:    &s.MgmtIPNetmask,
				OfflineReason:    &s.OfflineReason,
				PhysLocation:     &s.PhysLocation,
				PhysLocationID:   &s.PhysLocationID,
				Profile:          &s.Profile,
				ProfileDesc:      &s.ProfileDesc,
				ProfileID:        &s.ProfileID,
				Rack:             &s.Rack,
				RevalPending:     &s.RevalPending,
				Status:           &s.Status,
				StatusID:         &s.StatusID,
				TCPPort:          &s.TCPPort,
				Type:             s.Type,
				TypeID:           &s.TypeID,
				UpdPending:       &s.UpdPending,
				XMPPID:           &s.XMPPID,
				XMPPPasswd:       &s.XMPPPasswd,
			},
			LegacyInterfaceDetails: LegacyInterfaceDetails{
				InterfaceMtu:  &s.InterfaceMtu,
				InterfaceName: &s.InterfaceName,
				IPAddress:     &s.IPAddress,
				IPGateway:     &s.IPGateway,
				IPNetmask:     &s.IPNetmask,
				IP6Address:    &s.IP6Address,
				IP6Gateway:    &s.IP6Gateway,
			},
			RouterHostName: &s.RouterHostName,
			RouterPortName: &s.RouterPortName,
		},
		IPIsService:  &s.IPIsService,
		IP6IsService: &s.IP6IsService,
	}
}

// ToNonNullable converts the ServerNullableV2 safely to a Server structure.
//
// Deprecated: Traffic Ops API version 2 is deprecated, new code should use
// ServerV40 or newer structures.
func (s ServerNullableV2) ToNonNullable() Server {
	ret := Server{
		Cachegroup:     util.CoalesceToDefault(s.Cachegroup),
		CachegroupID:   util.CoalesceToDefault(s.CachegroupID),
		CDNID:          util.CoalesceToDefault((s.CDNID)),
		CDNName:        util.CoalesceToDefault(s.CDNName),
		DomainName:     util.CoalesceToDefault(s.DomainName),
		FQDN:           s.FQDN,
		FqdnTime:       s.FqdnTime,
		GUID:           util.CoalesceToDefault(s.GUID),
		HostName:       util.CoalesceToDefault(s.HostName),
		HTTPSPort:      util.CoalesceToDefault(s.HTTPSPort),
		ID:             util.CoalesceToDefault(s.ID),
		ILOIPAddress:   util.CoalesceToDefault(s.ILOIPAddress),
		ILOIPGateway:   util.CoalesceToDefault(s.ILOIPGateway),
		ILOIPNetmask:   util.CoalesceToDefault(s.ILOIPNetmask),
		ILOPassword:    util.CoalesceToDefault(s.ILOPassword),
		ILOUsername:    util.CoalesceToDefault(s.ILOUsername),
		InterfaceMtu:   util.CoalesceToDefault(s.InterfaceMtu),
		InterfaceName:  util.CoalesceToDefault(s.InterfaceName),
		IP6Address:     util.CoalesceToDefault(s.IP6Address),
		IP6IsService:   util.CoalesceToDefault(s.IP6IsService),
		IP6Gateway:     util.CoalesceToDefault(s.IP6Gateway),
		IPAddress:      util.CoalesceToDefault(s.IPAddress),
		IPIsService:    util.CoalesceToDefault(s.IPIsService),
		IPGateway:      util.CoalesceToDefault(s.IPGateway),
		IPNetmask:      util.CoalesceToDefault(s.IPNetmask),
		MgmtIPAddress:  util.CoalesceToDefault(s.MgmtIPAddress),
		MgmtIPGateway:  util.CoalesceToDefault(s.MgmtIPGateway),
		MgmtIPNetmask:  util.CoalesceToDefault(s.MgmtIPNetmask),
		OfflineReason:  util.CoalesceToDefault(s.OfflineReason),
		PhysLocation:   util.CoalesceToDefault(s.PhysLocation),
		PhysLocationID: util.CoalesceToDefault(s.PhysLocationID),
		Profile:        util.CoalesceToDefault(s.Profile),
		ProfileDesc:    util.CoalesceToDefault(s.ProfileDesc),
		ProfileID:      util.CoalesceToDefault(s.ProfileID),
		Rack:           util.CoalesceToDefault(s.Rack),
		RevalPending:   util.CoalesceToDefault(s.RevalPending),
		RouterHostName: util.CoalesceToDefault(s.RouterHostName),
		RouterPortName: util.CoalesceToDefault(s.RouterPortName),
		Status:         util.CoalesceToDefault(s.Status),
		StatusID:       util.CoalesceToDefault(s.StatusID),
		TCPPort:        util.CoalesceToDefault(s.TCPPort),
		Type:           s.Type,
		TypeID:         util.CoalesceToDefault(s.TypeID),
		UpdPending:     util.CoalesceToDefault(s.UpdPending),
		XMPPID:         util.CoalesceToDefault(s.XMPPID),
		XMPPPasswd:     util.CoalesceToDefault(s.XMPPPasswd),
	}

	if s.DeliveryServices == nil {
		ret.DeliveryServices = nil
	} else {
		ret.DeliveryServices = *s.DeliveryServices
	}

	if s.LastUpdated == nil {
		ret.LastUpdated = TimeNoMod{}
	} else {
		ret.LastUpdated = *s.LastUpdated
	}

	return ret
}

// Upgrade upgrades the ServerNullableV2 to the new ServerNullable structure.
//
// Note that this makes "shallow" copies of all underlying data, so changes to
// the original will affect the upgraded copy.
//
// Deprecated: Traffic Ops API versions 2 and 3 are both deprecated, new code
// should use ServerV40 or newer structures.
func (s ServerNullableV2) Upgrade() (ServerV30, error) {
	ipv4IsService := false
	if s.IPIsService != nil {
		ipv4IsService = *s.IPIsService
	}
	ipv6IsService := false
	if s.IP6IsService != nil {
		ipv6IsService = *s.IP6IsService
	}

	upgraded := ServerV30{
		CommonServerProperties: s.CommonServerProperties,
		RouterHostName:         s.RouterHostName,
		RouterPortName:         s.RouterPortName,
	}

	infs, err := s.LegacyInterfaceDetails.ToInterfaces(ipv4IsService, ipv6IsService)
	if err != nil {
		return upgraded, err
	}
	upgraded.Interfaces = infs
	return upgraded, nil
}

// UpgradeToV40 upgrades the ServerV30 to a ServerV40.
//
// This makes a "shallow" copy of the structure's properties.
//
// Deprecated: Traffic Ops API version 3 is deprecated, new code should use
// ServerV40 or newer structures.
func (s ServerV30) UpgradeToV40(profileNames []string) (ServerV40, error) {
	upgraded := UpdateServerPropertiesV40(profileNames, s.CommonServerProperties)
	upgraded.StatusLastUpdated = s.StatusLastUpdated
	infs, err := ToInterfacesV4(s.Interfaces, s.RouterHostName, s.RouterPortName)
	if err != nil {
		return upgraded, err
	}
	upgraded.Interfaces = infs
	return upgraded, nil
}

// UpgradeToV40 upgrades the ServerNullableV2 to a ServerV40.
//
// This makes a "shallow" copy of the structure's properties.
//
// Deprecated: Traffic Ops API version 2 is gone, new code should use
// ServerV40 or newer structures.
func (s ServerNullableV2) UpgradeToV40(profileNames []string) (ServerV40, error) {
	ipv4IsService := false
	if s.IPIsService != nil {
		ipv4IsService = *s.IPIsService
	}
	ipv6IsService := false
	if s.IP6IsService != nil {
		ipv6IsService = *s.IP6IsService
	}
	upgraded := UpdateServerPropertiesV40(profileNames, s.CommonServerProperties)

	infs, err := s.LegacyInterfaceDetails.ToInterfacesV4(ipv4IsService, ipv6IsService, s.RouterHostName, s.RouterPortName)
	if err != nil {
		return upgraded, err
	}
	upgraded.Interfaces = infs
	return upgraded, nil
}

// UpdateServerPropertiesV40 updates CommonServerProperties of V2 and V3 to
// ServerV40.
//
// Deprecated: Traffic Ops API version 3 is deprecated, new code should use
// ServerV40 or newer structures.
func UpdateServerPropertiesV40(profileNames []string, properties CommonServerProperties) ServerV40 {
	return ServerV40{
		Cachegroup:       properties.Cachegroup,
		CachegroupID:     properties.CachegroupID,
		CDNID:            properties.CDNID,
		CDNName:          properties.CDNName,
		DeliveryServices: properties.DeliveryServices,
		DomainName:       properties.DomainName,
		FQDN:             properties.FQDN,
		FqdnTime:         properties.FqdnTime,
		GUID:             properties.GUID,
		HostName:         properties.HostName,
		HTTPSPort:        properties.HTTPSPort,
		ID:               properties.ID,
		ILOIPAddress:     properties.ILOIPAddress,
		ILOIPGateway:     properties.ILOIPGateway,
		ILOIPNetmask:     properties.ILOIPNetmask,
		ILOPassword:      properties.ILOPassword,
		ILOUsername:      properties.ILOUsername,
		LastUpdated:      properties.LastUpdated,
		MgmtIPAddress:    properties.MgmtIPAddress,
		MgmtIPGateway:    properties.MgmtIPGateway,
		MgmtIPNetmask:    properties.MgmtIPNetmask,
		OfflineReason:    properties.OfflineReason,
		ProfileNames:     profileNames,
		PhysLocation:     properties.PhysLocation,
		PhysLocationID:   properties.PhysLocationID,
		Rack:             properties.Rack,
		RevalPending:     properties.RevalPending,
		Status:           properties.Status,
		StatusID:         properties.StatusID,
		TCPPort:          properties.TCPPort,
		Type:             properties.Type,
		TypeID:           properties.TypeID,
		UpdPending:       properties.UpdPending,
		XMPPID:           properties.XMPPID,
		XMPPPasswd:       properties.XMPPPasswd,
	}
}

// ServerV40 is the representation of a Server in version 4.0 of the Traffic Ops API.
type ServerV40 struct {
	Cachegroup   *string `json:"cachegroup" db:"cachegroup"`
	CachegroupID *int    `json:"cachegroupId" db:"cachegroup_id"`
	CDNID        *int    `json:"cdnId" db:"cdn_id"`
	CDNName      *string `json:"cdnName" db:"cdn_name"`
	// Deprecated: this has no known purpose, doesn't appear in any known API
	// responses, and it doesn't exist in the V5 version of this structure.
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
	// Deprecated: In the future, management interfaces must be configured as
	// interfaces within the Interfaces of the server, not separately on these
	// properties.
	MgmtIPAddress *string `json:"mgmtIpAddress" db:"mgmt_ip_address"`
	// Deprecated: In the future, management interfaces must be configured as
	// interfaces within the Interfaces of the server, not separately on these
	// properties.
	MgmtIPGateway *string `json:"mgmtIpGateway" db:"mgmt_ip_gateway"`
	// Deprecated: In the future, management interfaces must be configured as
	// interfaces within the Interfaces of the server, not separately on these
	// properties.
	MgmtIPNetmask  *string  `json:"mgmtIpNetmask" db:"mgmt_ip_netmask"`
	OfflineReason  *string  `json:"offlineReason" db:"offline_reason"`
	PhysLocation   *string  `json:"physLocation" db:"phys_location"`
	PhysLocationID *int     `json:"physLocationId" db:"phys_location_id"`
	ProfileNames   []string `json:"profileNames" db:"profile_name"`
	Rack           *string  `json:"rack" db:"rack"`
	// Deprecated: In APIv5 and later, this extraneous field is not calculated
	// by Traffic Ops; the information is available by comparing RevalUpdateTime
	// to RevalApplyTime.
	RevalPending *bool   `json:"revalPending" db:"reval_pending"`
	Status       *string `json:"status" db:"status"`
	StatusID     *int    `json:"statusId" db:"status_id"`
	TCPPort      *int    `json:"tcpPort" db:"tcp_port"`
	Type         string  `json:"type" db:"server_type"`
	TypeID       *int    `json:"typeId" db:"server_type_id"`
	// Deprecated: In APIv5 and later, this extraneous field is not calculated
	// by Traffic Ops; the information is available by comparing
	// ConfigUpdateTime to ConfigApplyTime.
	UpdPending        *bool                    `json:"updPending" db:"upd_pending"`
	XMPPID            *string                  `json:"xmppId" db:"xmpp_id"`
	XMPPPasswd        *string                  `json:"xmppPasswd" db:"xmpp_passwd"`
	Interfaces        []ServerInterfaceInfoV40 `json:"interfaces" db:"interfaces"`
	StatusLastUpdated *time.Time               `json:"statusLastUpdated" db:"status_last_updated"`
	ConfigUpdateTime  *time.Time               `json:"configUpdateTime" db:"config_update_time"`
	ConfigApplyTime   *time.Time               `json:"configApplyTime" db:"config_apply_time"`
	RevalUpdateTime   *time.Time               `json:"revalUpdateTime" db:"revalidate_update_time"`
	RevalApplyTime    *time.Time               `json:"revalApplyTime" db:"revalidate_apply_time"`
}

// ServerV4 is the representation of a Server in the latest minor version of
// version 4 of the Traffic Ops API.
type ServerV4 = ServerV40

// Upgrade upgrades to an APIv5 representation of a Server.
func (s ServerV4) Upgrade() ServerV50 {
	upgraded := ServerV50{
		CacheGroup:         util.CoalesceToDefault(s.Cachegroup),
		CacheGroupID:       util.CoalesceToDefault(s.CachegroupID),
		CDNID:              util.CoalesceToDefault(s.CDNID),
		CDN:                util.CoalesceToDefault(s.CDNName),
		DomainName:         util.CoalesceToDefault(s.DomainName),
		GUID:               util.CopyIfNotNil(s.GUID),
		HostName:           util.CoalesceToDefault(s.HostName),
		HTTPSPort:          util.CopyIfNotNil(s.HTTPSPort),
		ID:                 util.CoalesceToDefault(s.ID),
		ILOIPAddress:       util.CopyIfNotNil(s.ILOIPAddress),
		ILOIPGateway:       util.CopyIfNotNil(s.ILOIPGateway),
		ILOIPNetmask:       util.CopyIfNotNil(s.ILOIPNetmask),
		ILOPassword:        util.CopyIfNotNil(s.ILOPassword),
		ILOUsername:        util.CopyIfNotNil(s.ILOUsername),
		LastUpdated:        util.CoalesceToDefault(s.LastUpdated).Time,
		MgmtIPAddress:      util.CopyIfNotNil(s.MgmtIPAddress),
		MgmtIPGateway:      util.CopyIfNotNil(s.MgmtIPGateway),
		MgmtIPNetmask:      util.CopyIfNotNil(s.MgmtIPNetmask),
		OfflineReason:      util.CopyIfNotNil(s.OfflineReason),
		PhysicalLocation:   util.CoalesceToDefault(s.PhysLocation),
		PhysicalLocationID: util.CoalesceToDefault(s.PhysLocationID),
		Profiles:           make([]string, len(s.ProfileNames)),
		Rack:               util.CopyIfNotNil(s.Rack),
		Status:             util.CoalesceToDefault(s.Status),
		StatusID:           util.CoalesceToDefault(s.StatusID),
		TCPPort:            util.CopyIfNotNil(s.TCPPort),
		Type:               s.Type,
		TypeID:             util.CoalesceToDefault(s.TypeID),
		XMPPID:             util.CopyIfNotNil(s.XMPPID),
		XMPPPasswd:         util.CopyIfNotNil(s.XMPPPasswd),
		Interfaces:         make([]ServerInterfaceInfoV40, len(s.Interfaces)),
		StatusLastUpdated:  util.CopyIfNotNil(s.StatusLastUpdated),
		ConfigUpdateTime:   util.CopyIfNotNil(s.ConfigUpdateTime),
		ConfigApplyTime:    util.CopyIfNotNil(s.ConfigApplyTime),
		RevalUpdateTime:    util.CopyIfNotNil(s.RevalUpdateTime),
		RevalApplyTime:     util.CopyIfNotNil(s.RevalApplyTime),
	}

	copy(upgraded.Profiles, s.ProfileNames)

	for i, inf := range s.Interfaces {
		upgraded.Interfaces[i] = inf.Copy()
	}

	return upgraded
}

// ServerV30 is the representation of a Server in version 3 of the Traffic Ops API.
//
// Deprecated: Traffic Ops API version 3 is deprecated, new code should use
// ServerV40 or newer structures.
type ServerV30 struct {
	CommonServerProperties
	RouterHostName    *string               `json:"routerHostName" db:"router_host_name"`
	RouterPortName    *string               `json:"routerPortName" db:"router_port_name"`
	Interfaces        []ServerInterfaceInfo `json:"interfaces" db:"interfaces"`
	StatusLastUpdated *time.Time            `json:"statusLastUpdated" db:"status_last_updated"`
}

// ServerNullable represents an ATC server, as returned by the TO API.
//
// Deprecated: Traffic Ops API version 3 is deprecated, new code should use
// ServerV40 or newer structures.
type ServerNullable ServerV30

// ToServerV2 converts the server to an equivalent ServerNullableV2 structure,
// if possible. If the conversion could not be performed, an error is returned.
func (s *ServerNullable) ToServerV2() (ServerNullableV2, error) {
	nullable := ServerV30(*s)
	return nullable.ToServerV2()
}

// ToServerV2 converts the server to an equivalent ServerNullableV2 structure,
// if possible. If the conversion could not be performed, an error is returned.
//
// Deprecated: Traffic Ops API version 2 is deprecated, new code should use
// ServerV40 or newer structures.
func (s *ServerV30) ToServerV2() (ServerNullableV2, error) {
	legacyServer := ServerNullableV2{
		ServerNullableV11: ServerNullableV11{
			CommonServerProperties: s.CommonServerProperties,
		},
		IPIsService:  new(bool),
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

// ToServerV3FromV4 downgrades the ServerV40 to a ServerV30.
//
// This makes a "shallow" copy of most of the structure's properties.
//
// Deprecated: Traffic Ops API version 3 is deprecated, new code should use
// ServerV40 or newer structures.
func (s *ServerV40) ToServerV3FromV4(csp CommonServerProperties) (ServerV30, error) {
	routerHostName := ""
	routerPortName := ""
	interfaces := make([]ServerInterfaceInfo, 0)
	i := ServerInterfaceInfo{}
	for _, in := range s.Interfaces {
		i.Name = in.Name
		i.MTU = in.MTU
		i.MaxBandwidth = in.MaxBandwidth
		i.Monitor = in.Monitor
		i.IPAddresses = in.IPAddresses
		for _, ip := range i.IPAddresses {
			if ip.ServiceAddress {
				routerHostName = in.RouterHostName
				routerPortName = in.RouterPortName
			}
		}
		interfaces = append(interfaces, i)
	}
	serverV30 := ServerV30{
		CommonServerProperties: csp,
		Interfaces:             interfaces,
		StatusLastUpdated:      s.StatusLastUpdated,
	}
	if len(s.Interfaces) != 0 {
		serverV30.RouterHostName = &routerHostName
		serverV30.RouterPortName = &routerPortName
	}
	return serverV30, nil
}

// ToServerV2FromV4 downgrades the ServerV40 to a ServerNullableV2.
//
// This makes a "shallow" copy of most of the structure's properties.
//
// Deprecated: Traffic Ops API version 2 is deprecated, new code should use
// ServerV40 or newer structures.
func (s *ServerV40) ToServerV2FromV4(csp CommonServerProperties) (ServerNullableV2, error) {
	routerHostName := ""
	routerPortName := ""
	legacyServer := ServerNullableV2{
		ServerNullableV11: ServerNullableV11{
			CommonServerProperties: csp,
		},
		IPIsService:  new(bool),
		IP6IsService: new(bool),
	}

	interfaces := make([]ServerInterfaceInfo, 0)
	i := ServerInterfaceInfo{}
	for _, in := range s.Interfaces {
		i.Name = in.Name
		i.MTU = in.MTU
		i.MaxBandwidth = in.MaxBandwidth
		i.Monitor = in.Monitor
		i.IPAddresses = in.IPAddresses
		for _, ip := range i.IPAddresses {
			if ip.ServiceAddress {
				routerHostName = in.RouterHostName
				routerPortName = in.RouterPortName
			}
		}
		interfaces = append(interfaces, i)
	}
	var err error
	legacyServer.LegacyInterfaceDetails, err = InterfaceInfoToLegacyInterfaces(interfaces)
	if err != nil {
		return legacyServer, err
	}

	*legacyServer.IPIsService = legacyServer.LegacyInterfaceDetails.IPAddress != nil && *legacyServer.LegacyInterfaceDetails.IPAddress != ""
	*legacyServer.IP6IsService = legacyServer.LegacyInterfaceDetails.IP6Address != nil && *legacyServer.LegacyInterfaceDetails.IP6Address != ""
	if len(s.Interfaces) != 0 {
		legacyServer.RouterHostName = &routerHostName
		legacyServer.RouterPortName = &routerPortName
	}
	return legacyServer, nil
}

// ServerV50 is the representation of a Server in version 5.0 of the Traffic Ops
// API.
type ServerV50 struct {
	CacheGroup   string `json:"cacheGroup" db:"cachegroup"`
	CacheGroupID int    `json:"cacheGroupID" db:"cachegroup_id"`
	CDNID        int    `json:"cdnID" db:"cdn_id"`
	CDN          string `json:"cdn" db:"cdn_name"`
	// The time at which configuration updates were last applied for this server
	// by t3c.
	ConfigApplyTime *time.Time `json:"configApplyTime,omitempty" db:"config_apply_time"`
	// The time at which configuration updates were last queued for this server.
	ConfigUpdateTime *time.Time `json:"configUpdateTime,omitempty" db:"config_update_time"`
	// If the last config apply failed for this server
	ConfigUpdateFailed bool   `json:"configUpdateFailed" db:"config_update_failed"`
	DomainName         string `json:"domainName" db:"domain_name"`
	// Deprecated: This property has unknown purpose and should not be used so
	// that we can get rid of it.
	GUID         *string                  `json:"guid" db:"guid"`
	HostName     string                   `json:"hostName" db:"host_name"`
	HTTPSPort    *int                     `json:"httpsPort,omitempty" db:"https_port"`
	ID           int                      `json:"id" db:"id"`
	ILOIPAddress *string                  `json:"iloIpAddress,omitempty" db:"ilo_ip_address"`
	ILOIPGateway *string                  `json:"iloIpGateway,omitempty" db:"ilo_ip_gateway"`
	ILOIPNetmask *string                  `json:"iloIpNetmask,omitempty" db:"ilo_ip_netmask"`
	ILOPassword  *string                  `json:"iloPassword,omitempty" db:"ilo_password"`
	ILOUsername  *string                  `json:"iloUsername,omitempty" db:"ilo_username"`
	Interfaces   []ServerInterfaceInfoV40 `json:"interfaces" db:"interfaces"`
	LastUpdated  time.Time                `json:"lastUpdated" db:"last_updated"`
	// Deprecated: In the future, management interfaces must be configured as
	// interfaces within the Interfaces of the server, not separately on these
	// properties.
	MgmtIPAddress *string `json:"mgmtIpAddress,omitempty" db:"mgmt_ip_address"`
	// Deprecated: In the future, management interfaces must be configured as
	// interfaces within the Interfaces of the server, not separately on these
	// properties.
	MgmtIPGateway *string `json:"mgmtIpGateway,omitempty" db:"mgmt_ip_gateway"`
	// Deprecated: In the future, management interfaces must be configured as
	// interfaces within the Interfaces of the server, not separately on these
	// properties.
	MgmtIPNetmask      *string  `json:"mgmtIpNetmask,omitempty" db:"mgmt_ip_netmask"`
	OfflineReason      *string  `json:"offlineReason" db:"offline_reason"`
	PhysicalLocation   string   `json:"physicalLocation" db:"phys_location"`
	PhysicalLocationID int      `json:"physicalLocationID" db:"phys_location_id"`
	Profiles           []string `json:"profiles,omitempty" db:"profile_name"`
	// Deprecated: This property has unknown purpose and should not be used so
	// that we can get rid of it.
	Rack *string `json:"rack" db:"rack"`
	// The time at which revalidations for this server were last updated by t3c.
	RevalApplyTime *time.Time `json:"revalApplyTime,omitempty" db:"revalidate_apply_time"`
	// The time at which revalidations were last queued for this server.
	RevalUpdateTime *time.Time `json:"revalUpdateTime,omitempty" db:"revalidate_update_time"`
	// If the last reval apply failed for this server
	RevalUpdateFailed bool       `json:"revalUpdateFailed" db:"revalidate_update_failed"`
	Status            string     `json:"status" db:"status"`
	StatusID          int        `json:"statusID" db:"status_id"`
	StatusLastUpdated *time.Time `json:"statusLastUpdated,omitempty" db:"status_last_updated"`
	TCPPort           *int       `json:"tcpPort" db:"tcp_port"`
	Type              string     `json:"type" db:"server_type"`
	TypeID            int        `json:"typeID" db:"server_type_id"`
	XMPPID            *string    `json:"xmppId" db:"xmpp_id"`
	// Deprecated: This property has unknown purpose and should not be used so
	// that we can get rid of it.
	XMPPPasswd *string `json:"xmppPasswd" db:"xmpp_passwd"`
}

// Downgrade downgrades to a V4 representation of a Server.
func (s ServerV50) Downgrade() ServerV4 {
	downgraded := ServerV40{
		Cachegroup:        util.Ptr(s.CacheGroup),
		CachegroupID:      util.Ptr(s.CacheGroupID),
		CDNID:             util.Ptr(s.CDNID),
		CDNName:           util.Ptr(s.CDN),
		DeliveryServices:  nil,
		DomainName:        util.Ptr(s.DomainName),
		FQDN:              util.Ptr(s.HostName + "." + s.DomainName),
		FqdnTime:          time.Time{},
		GUID:              util.CopyIfNotNil(s.GUID),
		HostName:          util.Ptr(s.HostName),
		HTTPSPort:         util.CopyIfNotNil(s.HTTPSPort),
		ID:                util.Ptr(s.ID),
		ILOIPAddress:      util.CopyIfNotNil(s.ILOIPAddress),
		ILOIPGateway:      util.CopyIfNotNil(s.ILOIPGateway),
		ILOIPNetmask:      util.CopyIfNotNil(s.ILOIPNetmask),
		ILOPassword:       util.CopyIfNotNil(s.ILOPassword),
		ILOUsername:       util.CopyIfNotNil(s.ILOUsername),
		LastUpdated:       &TimeNoMod{Time: s.LastUpdated},
		MgmtIPAddress:     util.CopyIfNotNil(s.MgmtIPAddress),
		MgmtIPGateway:     util.CopyIfNotNil(s.MgmtIPGateway),
		MgmtIPNetmask:     util.CopyIfNotNil(s.MgmtIPNetmask),
		OfflineReason:     util.CopyIfNotNil(s.OfflineReason),
		PhysLocation:      util.Ptr(s.PhysicalLocation),
		PhysLocationID:    util.Ptr(s.PhysicalLocationID),
		ProfileNames:      make([]string, len(s.Profiles)),
		Rack:              util.CopyIfNotNil(s.Rack),
		RevalPending:      util.Ptr(s.RevalidationPending()),
		Status:            util.Ptr(s.Status),
		StatusID:          util.Ptr(s.StatusID),
		TCPPort:           util.CopyIfNotNil(s.TCPPort),
		Type:              s.Type,
		TypeID:            util.Ptr(s.TypeID),
		UpdPending:        util.Ptr(s.UpdatePending()),
		XMPPID:            util.CopyIfNotNil(s.XMPPID),
		XMPPPasswd:        util.CopyIfNotNil(s.XMPPPasswd),
		Interfaces:        make([]ServerInterfaceInfoV40, len(s.Interfaces)),
		StatusLastUpdated: util.CopyIfNotNil(s.StatusLastUpdated),
		ConfigUpdateTime:  util.CopyIfNotNil(s.ConfigUpdateTime),
		ConfigApplyTime:   util.CopyIfNotNil(s.ConfigApplyTime),
		RevalUpdateTime:   util.CopyIfNotNil(s.RevalUpdateTime),
		RevalApplyTime:    util.CopyIfNotNil(s.RevalApplyTime),
	}

	copy(downgraded.ProfileNames, s.Profiles)
	for i, inf := range s.Interfaces {
		downgraded.Interfaces[i] = inf.Copy()
	}

	return downgraded
}

// UpdatePending tells whether the Server has pending updates.
func (s ServerV50) UpdatePending() bool {
	return s.ConfigApplyTime != nil && s.ConfigUpdateTime != nil && s.ConfigApplyTime.Before(*s.ConfigUpdateTime)
}

// RevalidationPending tells whether the Server has pending revalidations.
func (s ServerV50) RevalidationPending() bool {
	return s.RevalApplyTime != nil && s.RevalUpdateTime != nil && s.RevalApplyTime.Before(*s.RevalUpdateTime)
}

// ServerV5 is the representation of a Server in the latest minor version of
// version 5 of the Traffic Ops API.
type ServerV5 = ServerV50

// ServerUpdateStatusV5 is the type of each entry in the `response` property of
// the response from Traffic Ops to GET requests made to its
// /servers/{{host name}}/update_status in the latest minor API
// v5.0 endpoint.
type ServerUpdateStatusV5 ServerUpdateStatusV50

// ServerUpdateStatusV50 is the type of each entry in the `response` property of
// the response from Traffic Ops to GET requests made to its
// /servers/{{host name}}/update_status in API v5.0 endpoint.
type ServerUpdateStatusV50 struct {
	HostName string `json:"host_name"`
	// Deprecated: In APIv5 and later, this extraneous field is not calculated
	// by Traffic Ops; the information is available by comparing ConfigUpdateTime
	// to ConfigApplyTime.
	UpdatePending bool `json:"upd_pending"`
	// Deprecated: In APIv5 and later, this extraneous field is not calculated
	// by Traffic Ops; the information is available by comparing RevalUpdateTime
	// to RevalApplyTime.
	RevalPending           bool       `json:"reval_pending"`
	UseRevalPending        bool       `json:"use_reval_pending"`
	HostId                 int        `json:"host_id"`
	Status                 string     `json:"status"`
	ParentPending          bool       `json:"parent_pending"`
	ParentRevalPending     bool       `json:"parent_reval_pending"`
	ConfigUpdateTime       *time.Time `json:"config_update_time"`
	ConfigApplyTime        *time.Time `json:"config_apply_time"`
	ConfigUpdateFailed     *bool      `json:"config_update_failed"`
	RevalidateUpdateTime   *time.Time `json:"revalidate_update_time"`
	RevalidateApplyTime    *time.Time `json:"revalidate_apply_time"`
	RevalidateUpdateFailed *bool      `json:"revalidate_update_failed"`
}

func (sus ServerUpdateStatusV5) Downgrade() ServerUpdateStatusV40 {
	return ServerUpdateStatusV40{
		sus.HostName,
		sus.UpdatePending,
		sus.RevalPending,
		sus.UseRevalPending,
		sus.HostId,
		sus.Status,
		sus.ParentPending,
		sus.ParentRevalPending,
		sus.ConfigUpdateTime,
		sus.ConfigApplyTime,
		sus.RevalidateUpdateTime,
		sus.RevalidateApplyTime,
	}
}

// ServerUpdateStatusV4 is the type of each entry in the `response` property of
// the response from Traffic Ops to GET requests made to its
// /servers/{{host name}}/update_status in the latest minor API
// v4.0 endpoint.
type ServerUpdateStatusV4 ServerUpdateStatusV40

// ServerUpdateStatusV40 is the type of each entry in the `response` property of
// the response from Traffic Ops to GET requests made to its
// /servers/{{host name}}/update_status in API v4.0 endpoint.
type ServerUpdateStatusV40 struct {
	HostName      string `json:"host_name"`
	UpdatePending bool   `json:"upd_pending"`
	// Deprecated: In APIv5 and later, this extraneous field is not calculated
	// by Traffic Ops; the information is available by comparing RevalUpdateTime
	// to RevalApplyTime.
	RevalPending         bool       `json:"reval_pending"`
	UseRevalPending      bool       `json:"use_reval_pending"`
	HostId               int        `json:"host_id"`
	Status               string     `json:"status"`
	ParentPending        bool       `json:"parent_pending"`
	ParentRevalPending   bool       `json:"parent_reval_pending"`
	ConfigUpdateTime     *time.Time `json:"config_update_time"`
	ConfigApplyTime      *time.Time `json:"config_apply_time"`
	RevalidateUpdateTime *time.Time `json:"revalidate_update_time"`
	RevalidateApplyTime  *time.Time `json:"revalidate_apply_time"`
}

// Downgrade strips the Config and Revalidate timestamps from
// ServerUpdateStatusV40 to return previous versions of the struct to ensure
// previous compatibility.
func (sus ServerUpdateStatusV40) Downgrade() ServerUpdateStatus {
	return ServerUpdateStatus{
		HostName:           sus.HostName,
		UpdatePending:      sus.UpdatePending,
		RevalPending:       sus.RevalPending,
		UseRevalPending:    sus.UseRevalPending,
		HostId:             sus.HostId,
		Status:             sus.Status,
		ParentPending:      sus.ParentPending,
		ParentRevalPending: sus.ParentRevalPending,
	}
}

// ServerUpdateStatus is the type of each entry in the `response` property of
// the response from Traffic Ops to GET requests made to its
// /servers/{{host name}}/update_status API endpoint.
//
// This is a subset of Server structure information mainly relating to what
// operations t3c has done/needs to do. For most purposes, using Server
// structures will be better - especially since the basic principle of this
// type is predicated on a lie: that server host names are unique.
//
// Deprecated: ServerUpdateStatus is for use only in APIs below V4. New code
// should use ServerUpdateStatusV40 or newer.
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

// Upgrade converts the deprecated ServerUpdateStatus to a
// ServerUpdateStatusV4 struct.
func (sus ServerUpdateStatus) Upgrade() ServerUpdateStatusV4 {
	return ServerUpdateStatusV4{
		HostName:           sus.HostName,
		UpdatePending:      sus.UpdatePending,
		RevalPending:       sus.RevalPending,
		UseRevalPending:    sus.UseRevalPending,
		HostId:             sus.HostId,
		Status:             sus.Status,
		ParentPending:      sus.ParentPending,
		ParentRevalPending: sus.ParentRevalPending,
	}
}

// ServerUpdateStatusResponseV50 is the type of a response from the Traffic
// Ops API to a request to its /servers/{{host name}}/update_status endpoint
// in API version 5.0.
type ServerUpdateStatusResponseV50 struct {
	Response []ServerUpdateStatusV50 `json:"response"`
	Alerts
}

// ServerUpdateStatusResponseV5 is the type of a response from the Traffic
// Ops API to a request to its /servers/{{host name}}/update_status endpoint
// in the latest minor version of API version 5.
type ServerUpdateStatusResponseV5 = ServerUpdateStatusResponseV50

// ServerUpdateStatusResponseV40 is the type of a response from the Traffic
// Ops API to a request to its /servers/{{host name}}/update_status endpoint
// in API version 4.0.
type ServerUpdateStatusResponseV40 struct {
	Response []ServerUpdateStatusV40 `json:"response"`
	Alerts
}

// ServerUpdateStatusResponseV4 is the type of a response from the Traffic
// Ops API to a request to its /servers/{{host name}}/update_status endpoint
// in the latest minor version of API version 4.
type ServerUpdateStatusResponseV4 = ServerUpdateStatusResponseV40

// ServerPutStatus is a request to change the Status of a server, optionally
// with an explanation.
type ServerPutStatus struct {
	Status        util.JSONNameOrIDStr `json:"status"`
	OfflineReason *string              `json:"offlineReason"`
}

// ServerInfo is a stripped-down type containing a subset of information for a
// server.
//
// This is primarily only useful internally in Traffic Ops for
// constructing/examining the relationships between servers and other ATC
// objects. That is to say, for most other purposes a ServerV4 would be better
// suited.
type ServerInfo struct {
	Cachegroup   string
	CachegroupID int
	CDNID        int
	DomainName   string
	HostName     string
	ID           int
	Status       string
	Type         string
}

// ServerDetail is a type that contains a superset of the information available
// in a ServerNullable.
//
// This should NOT be used in general, it's almost always better to use a
// proper server representation. However, this structure is not deprecated
// because it is embedded in structures still in use, and its future is unclear
// and undecided.
//
// Deprecated: The API versions that use this representation have been
// deprecated, newer structures like ServerV4 should be used instead.
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

// ServersV5Response is the format of a response to a GET request to /servers in
// APIv5.
type ServersV5Response struct {
	Response []ServerV5 `json:"response"`
	Summary  struct {
		Count uint64 `json:"count"`
	} `json:"summary"`
	Alerts
}

// ServerV5Response is the format of a response to single-server operations
// using /servers in APIv5.
type ServerV5Response struct {
	Response ServerV5 `json:"response"`
	Summary  struct {
		Count uint64 `json:"count"`
	} `json:"summary"`
	Alerts
}
