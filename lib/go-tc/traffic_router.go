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
	"time"
)

// CoverageZonePollingPrefix is the prefix of all Names of Parameters used to define
// coverage zone polling parameters.
const CoverageZonePollingPrefix = "coveragezone.polling."

// CoverageZonePollingURL is the Name of the Parameter that defines the URL that
// Traffic Router will repeatedly poll for a coverage zone file.
const CoverageZonePollingURL = CoverageZonePollingPrefix + "url"

// CoverageZoneLocation represents a single "location" in a coverage zone file.
type CoverageZoneLocation struct {
	Network  []string `json:"network,omitempty"`
	Network6 []string `json:"network6,omitempty"`
}

// GetFirstIPAddressOfType returns the first IP address (or CIDR-notation
// subnet) that can be found within the location that meets the criteria
// specified by its argument - "true" means an IPv4 address should be returned,
// while "false" means it should be IPv6. If no addresses (or subnets) are found
// that meet that criteria, an empty string is returned.
func (c *CoverageZoneLocation) GetFirstIPAddressOfType(isIPv4 bool) string {
	var network []string
	if isIPv4 {
		network = c.Network
	} else {
		network = c.Network6
	}
	if len(network) < 1 {
		return ""
	}
	return network[0]
}

// CoverageZoneFile is used for unmarshalling a Coverage Zone File.
type CoverageZoneFile struct {
	CoverageZones map[string]CoverageZoneLocation `json:"coverageZones,omitempty"`
}

// X_MM_CLIENT_IP is an optional HTTP header that causes Traffic Router to use its value
// as the client IP address.
const X_MM_CLIENT_IP = "X-MM-Client-IP"

// SOA (Start of Authority record) defines the SOA record for the CDN's
// top-level domain.
type SOA struct {
	Admin              *string   `json:"admin,omitempty"`
	AdminTime          time.Time `json:"-"`
	ExpireSeconds      *string   `json:"expire,omitempty"`
	ExpireSecondsTime  time.Time `json:"-"`
	MinimumSeconds     *string   `json:"minimum,omitempty"`
	MinimumSecondsTime time.Time `json:"-"`
	RefreshSeconds     *string   `json:"refresh,omitempty"`
	RefreshSecondsTime time.Time `json:"-"`
	RetrySeconds       *string   `json:"retry,omitempty"`
	RetrySecondsTime   time.Time `json:"-"`
}

// MatchSet structures are a list of MatchList structures with an associated
// Protocol.
type MatchSet struct {
	Protocol  string      `json:"protocol"`
	MatchList []MatchList `json:"matchlist"`
}

// A MatchList is actually just a single match item in a match list.
type MatchList struct {
	Regex     string `json:"regex"`
	MatchType string `json:"match-type"`
}

// A LegacyTrafficServer is a representation of a cache server containing a
// subset of the information available in a server structure that conveys all
// the information important for Traffic Router and Traffic Monitor to handle
// it.
//
// Deprecated: The configuration versions that use this structure to represent
// a cache server are deprecated, new code should use TrafficServer instead.
type LegacyTrafficServer struct {
	CacheGroup       string              `json:"cacheGroup"`
	DeliveryServices []TSDeliveryService `json:"deliveryServices,omitempty"` // the deliveryServices key does not exist on mids
	FQDN             string              `json:"fqdn"`
	HashID           string              `json:"hashId"`
	HostName         string              `json:"hostName"`
	HTTPSPort        int                 `json:"httpsPort,omitempty"`
	InterfaceName    string              `json:"interfaceName"`
	IP               string              `json:"ip"`
	IP6              string              `json:"ip6"`
	Port             int                 `json:"port"`
	Profile          string              `json:"profile"`
	ServerStatus     string              `json:"status"`
	Type             string              `json:"type"`
}

// Upgrade upgrades the LegacyTrafficServer into its modern-day equivalent.
//
// Note that the DeliveryServices slice is a "shallow" copy of the original, so
// making changes to the original slice will affect the upgraded copy.
//
// Deprecated: LegacyTrafficServer is deprecated.
func (s LegacyTrafficServer) Upgrade() TrafficServer {
	upgraded := TrafficServer{
		CacheGroup:       s.CacheGroup,
		DeliveryServices: s.DeliveryServices,
		FQDN:             s.FQDN,
		HashID:           s.HashID,
		HostName:         s.HostName,
		HTTPSPort:        s.HTTPSPort,
		Interfaces: []ServerInterfaceInfo{
			{
				MaxBandwidth: nil,
				MTU:          nil,
				Monitor:      false,
				Name:         s.InterfaceName,
				IPAddresses:  []ServerIPAddress{},
			},
		},
		Port:         s.Port,
		Profile:      s.Profile,
		ServerStatus: s.ServerStatus,
		Type:         s.Type,
	}

	if s.IP != "" {
		upgraded.Interfaces[0] = ServerInterfaceInfo{
			MaxBandwidth: nil,
			MTU:          nil,
			Monitor:      false,
			Name:         s.InterfaceName,
			IPAddresses: []ServerIPAddress{
				{
					Address:        s.IP,
					Gateway:        nil,
					ServiceAddress: true,
				},
			},
		}
	}

	if s.IP6 != "" {
		upgraded.Interfaces[0] = ServerInterfaceInfo{
			MaxBandwidth: nil,
			MTU:          nil,
			Monitor:      false,
			Name:         s.InterfaceName,
			IPAddresses: append(upgraded.Interfaces[0].IPAddresses, ServerIPAddress{
				Address:        s.IP6,
				Gateway:        nil,
				ServiceAddress: true,
			}),
		}
	}

	return upgraded
}

// GetVIPInterface returns the primary interface specified by the `Monitor` property of an Interface. First interface marked as `Monitor` is returned.
func GetVIPInterface(ts TrafficServer) ServerInterfaceInfo {
	for _, interf := range ts.Interfaces {
		if interf.Monitor {
			return interf
		}
	}
	return ServerInterfaceInfo{}
}

// ToLegacyServer converts a TrafficServer to LegacyTrafficServer.
//
// Deprecated: LegacyTrafficServer is deprecated.
func (ts *TrafficServer) ToLegacyServer() LegacyTrafficServer {
	vipInterface := GetVIPInterface(*ts)
	ipv4, ipv6 := vipInterface.GetDefaultAddressOrCIDR()

	return LegacyTrafficServer{
		Profile:          ts.Profile,
		IP:               ipv4,
		ServerStatus:     ts.ServerStatus,
		CacheGroup:       ts.CacheGroup,
		IP6:              ipv6,
		Port:             ts.Port,
		HTTPSPort:        ts.HTTPSPort,
		HostName:         ts.HostName,
		FQDN:             ts.FQDN,
		InterfaceName:    vipInterface.Name,
		Type:             ts.Type,
		HashID:           ts.HashID,
		DeliveryServices: ts.DeliveryServices,
	}
}

// TrafficServer represents a cache server for use by Traffic Monitor and
// Traffic Router instances.
type TrafficServer struct {
	CacheGroup       string                `json:"cachegroup"`
	DeliveryServices []TSDeliveryService   `json:"deliveryServices,omitempty"` // the deliveryServices key does not exist on mids
	FQDN             string                `json:"fqdn"`
	HashID           string                `json:"hashId"`
	HostName         string                `json:"hostName"`
	HTTPSPort        int                   `json:"httpsPort,omitempty"`
	Interfaces       []ServerInterfaceInfo `json:"interfaces"`
	Port             int                   `json:"port"`
	Profile          string                `json:"profile"`
	ServerStatus     string                `json:"status"`
	Type             string                `json:"type"`
}

// IPv4 gets the server's IPv4 address if one exists, otherwise an empty
// string.
//
// Note: This swallows errors from the legacy data conversion process.
func (ts *TrafficServer) IPv4() string {
	if ts == nil {
		return ""
	}
	lid, err := InterfaceInfoToLegacyInterfaces(ts.Interfaces)
	if err != nil || lid.IPAddress == nil {
		return ""
	}
	return *lid.IPAddress
}

// IPv6 gets the server's IPv6 address if one exists, otherwise an empty
// string.
//
// Note: This swallows errors from the legacy data conversion process.
func (ts *TrafficServer) IPv6() string {
	if ts == nil {
		return ""
	}
	lid, err := InterfaceInfoToLegacyInterfaces(ts.Interfaces)
	if err != nil || lid.IP6Address == nil {
		return ""
	}
	return *lid.IP6Address
}

// TSDeliveryService represents a delivery service as assigned to a TrafficServer in the TMConfig.
type TSDeliveryService struct {
	XmlId string `json:"xmlId"`
}
