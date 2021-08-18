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

// MissLocation is a geographic location that will be used in the event that a
// client request cannot be localized (a "Geo Miss").
type MissLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
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

// A BypassDestination is an FQDN, IP address, or CNAME to use for redirection
// in the event that bypass becomes necessary.
type BypassDestination struct {
	FQDN string `json:"fqdn"`
	Type string `json:"type"`
	Port int    `json:"Port"`
}

// TTLS is a structure that contains Time-To-Live values for different types of
// DNS records (in seconds).
type TTLS struct {
	Arecord    int `json:"A"`
	SoaRecord  int `json:"SOA"`
	NsRecord   int `json:"NS"`
	AaaaRecord int `json:"AAAA"`
}

// TrafficRouter is a subset of a server that contains all the information a
// Traffic Router instance may need to know about one of its peers.
type TrafficRouter struct {
	Port         int    `json:"port"`
	IP6          string `json:"ip6"`
	IP           string `json:"ip"`
	FQDN         string `json:"fqdn"`
	Profile      string `json:"profile"`
	Location     string `json:"location"`
	ServerStatus string `json:"status"`
	APIPort      int    `json:"apiPort"`
}

// TrafficRouterConfig is the json unmarshalled without any changes
// note all structs are local to this file _except_ the TrafficRouterConfig struct.
type TrafficRouterConfig struct {
	TrafficServers   []TrafficServer        `json:"trafficServers,omitempty"`
	TrafficMonitors  []TrafficMonitor       `json:"trafficMonitors,omitempty"`
	TrafficRouters   []TrafficRouter        `json:"trafficRouters,omitempty"`
	CacheGroups      []TMCacheGroup         `json:"cacheGroups,omitempty"`
	DeliveryServices []TRDeliveryService    `json:"deliveryServices,omitempty"`
	Stats            map[string]interface{} `json:"stats,omitempty"`
	Config           map[string]interface{} `json:"config,omitempty"`
}

// TrafficRouterConfigMap is a set of mappings of various objects important to
// Traffic Router to the structural data a Traffic Router needs to know about
// them.
//
// In general, the objects herein contained do not strictly correspond to other
// structures that represent the same concept defined in this package - Traffic
// Router only gets exactly what information it needs, as keeping the size of
// this structure low is important for performance.
//
// TODO: Why isn't this used instead of legacy-named things like CRConfig?
type TrafficRouterConfigMap struct {
	TrafficServer   map[string]TrafficServer
	TrafficMonitor  map[string]TrafficMonitor
	TrafficRouter   map[string]TrafficRouter
	CacheGroup      map[string]TMCacheGroup
	DeliveryService map[string]TRDeliveryService
	Config          map[string]interface{}
	Stat            map[string]interface{}
}

// TRDeliveryService structures contain all the information about a Delivery
// Service that are important to Traffic Routers.
type TRDeliveryService struct {
	XMLID             string            `json:"xmlId"`
	Domains           []string          `json:"domains"`
	RoutingName       string            `json:"routingName"`
	MissLocation      MissLocation      `json:"missCoordinates"`
	CoverageZoneOnly  bool              `json:"coverageZoneOnly"`
	MatchSets         []MatchSet        `json:"matchSets"`
	TTL               int               `json:"ttl"`
	TTLs              TTLS              `json:"ttls"`
	BypassDestination BypassDestination `json:"bypassDestination"`
	StatcDNSEntries   []StaticDNS       `json:"statitDnsEntries"`
	Soa               SOA               `json:"soa"`
}

// A StaticDNS is a somewhat arbitrary static DNS entry of some kind.
type StaticDNS struct {
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
	Name  string `json:"name"`
	Type  string `json:"type"`
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
	DeliveryServices []tsdeliveryService `json:"deliveryServices,omitempty"` // the deliveryServices key does not exist on mids
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
	ipv4, ipv6 := vipInterface.GetDefaultAddress()

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
	CacheGroup       string                `json:"cacheGroup"`
	DeliveryServices []tsdeliveryService   `json:"deliveryServices,omitempty"` // the deliveryServices key does not exist on mids
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

type tsdeliveryService struct {
	Xmlid  string   `json:"xmlId"`
	Remaps []string `json:"remaps"`
}
