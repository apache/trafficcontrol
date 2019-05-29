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

// TrafficServer ...
type TrafficServer struct {
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

type tsdeliveryService struct {
	Xmlid  string   `json:"xmlId"`
	Remaps []string `json:"remaps"`
}
