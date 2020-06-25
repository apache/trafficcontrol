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
	"net"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
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

// MissLocation ...
type MissLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// MatchSet ...
type MatchSet struct {
	Protocol  string      `json:"protocol"`
	MatchList []MatchList `json:"matchlist"`
}

// MatchList ...
type MatchList struct {
	Regex     string `json:"regex"`
	MatchType string `json:"match-type"`
}

// BypassDestination ...
type BypassDestination struct {
	FQDN string `json:"fqdn"`
	Type string `json:"type"`
	Port int    `json:"Port"`
}

// TTLS ...
type TTLS struct {
	Arecord    int `json:"A"`
	SoaRecord  int `json:"SOA"`
	NsRecord   int `json:"NS"`
	AaaaRecord int `json:"AAAA"`
}

// TrafficRouter ...
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

// TrafficRouterConfigMap ...
type TrafficRouterConfigMap struct {
	TrafficServer   map[string]TrafficServer
	TrafficMonitor  map[string]TrafficMonitor
	TrafficRouter   map[string]TrafficRouter
	CacheGroup      map[string]TMCacheGroup
	DeliveryService map[string]TRDeliveryService
	Config          map[string]interface{}
	Stat            map[string]interface{}
}

// TRDeliveryService ...
// TODO JvD: move to deliveryservice.go ??
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

// StaticDNS ...
type StaticDNS struct {
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

// LegacyTrafficServer ...
type LegacyTrafficServer struct {
	Profile          string              `json:"profile"`
	IP               string              `json:"ip"`
	ServerStatus     string              `json:"status"`
	CacheGroup       string              `json:"cacheGroup"`
	IP6              string              `json:"ip6"`
	Port             int                 `json:"port"`
	HTTPSPort        int                 `json:"httpsPort,omitempty"`
	HostName         string              `json:"hostName"`
	FQDN             string              `json:"fqdn"`
	InterfaceName    string              `json:"interfaceName"`
	Type             string              `json:"type"`
	HashID           string              `json:"hashId"`
	DeliveryServices []tsdeliveryService `json:"deliveryServices,omitempty"` // the deliveryServices key does not exist on mids
}

// GetDefaultAddress returns the ipv4 and ipv6 service addresses of the interface.
func (i *InterfaceInfo) GetDefaultAddress() (string, string) {
	var ipv4 string
	var ipv6 string
	for _, ip := range i.IPAddresses {
		if ip.ServiceAddress {
			address, _, err := net.ParseCIDR(ip.Address)
			if err != nil {
				log.Warnf("Unable to parse ipaddress %v on interface %v: %v", ip.Address, i.Name, err)
			} else if address == nil {
				log.Warnf("Unable to parse ipaddress %v on interface %v", ip.Address, i.Name)
				continue
			}
			if address.To4() != nil {
				ipv4 = ip.Address
			} else if address.To16() != nil {
				ipv6 = ip.Address
			} else {
				log.Warnf("Invalid address %v on interface %v", address, i.Name)
			}

			if ipv4 != "" && ipv6 != "" {
				break
			}
		}
	}
	return ipv4, ipv6
}

// GetVIPInterface returns the primary interface specified by the `Monitor` property of an Interface. First interface marked as `Monitor` is returned.
func GetVIPInterface(ts TrafficServer) InterfaceInfo {
	for _, interf := range ts.Interfaces {
		if interf.Monitor {
			return interf
		}
	}
	return InterfaceInfo{}
}

// ToLegacyServer converts a TrafficServer to LegacyTrafficServer.
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

// TrafficServer ...
type TrafficServer struct {
	Profile          string              `json:"profile"`
	ServerStatus     string              `json:"status"`
	CacheGroup       string              `json:"cacheGroup"`
	Port             int                 `json:"port"`
	HostName         string              `json:"hostName"`
	FQDN             string              `json:"fqdn"`
	Interfaces       []InterfaceInfo     `json:"interfaces"`
	HTTPSPort        int                 `json:"httpsPort,omitempty"`
	Type             string              `json:"type"`
	HashID           string              `json:"hashId"`
	DeliveryServices []tsdeliveryService `json:"deliveryServices,omitempty"` // the deliveryServices key does not exist on mids
}

// ServerIPAddress is the data associated with a server's interface's IP address.
type IPAddress struct {
	Address        string  `json:"address" db:"address"`
	Gateway        *string `json:"gateway" db:"gateway"`
	ServiceAddress bool    `json:"serviceAddress" db:"service_address"`
}

// ServerInterfaceInfo is the data associated with a server's interface.
type InterfaceInfo struct {
	IPAddresses  []IPAddress `json:"ipAddresses" db:"ip_addresses"`
	MaxBandwidth *uint64     `json:"maxBandwidth" db:"max_bandwidth"`
	Monitor      bool        `json:"monitor" db:"monitor"`
	MTU          *uint64     `json:"mtu" db:"mtu"`
	Name         string      `json:"name" db:"name"`
}

type tsdeliveryService struct {
	Xmlid  string   `json:"xmlId"`
	Remaps []string `json:"remaps"`
}
