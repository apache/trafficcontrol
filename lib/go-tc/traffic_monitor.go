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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ThresholdPrefix is the prefix of all Names of Parameters used to define
// monitoring thresholds.
const ThresholdPrefix = "health.threshold."

// These are the Names of Parameters with special meaning to Traffic Monitor.
// They are used as thresholds in aggregate across all of a monitored cache
// server's network interfaces. Documentation on specific Parameter meanings
// can be found in the ATC official documentation.
const (
	AvailableBandwidthInKbpsThresholdName = ThresholdPrefix + "availableBandwidthInKbps"
	AvailableBandwidthInMbpsThresholdName = ThresholdPrefix + "availableBandwidthInMbps"
	BandwidthThresholdName                = ThresholdPrefix + "bandwidth"
	GbpsThresholdName                     = ThresholdPrefix + "gbps"
	KbpsThresholdName                     = ThresholdPrefix + "kbps"
	LoadavgThresholdName                  = ThresholdPrefix + "loadavg"
	MaxKbpsThresholdName                  = ThresholdPrefix + "maxKbps"
)

// TMConfigResponse is the response to requests made to the
// cdns/{{Name}}/configs/monitoring endpoint of the Traffic Ops API.
type TMConfigResponse struct {
	Response TrafficMonitorConfig `json:"response"`
}

// LegacyTMConfigResponse was the response to requests made to the
// cdns/{{Name}}/configs/montoring endpoint of the Traffic Ops API in older API
// versions.
type LegacyTMConfigResponse struct {
	Response LegacyTrafficMonitorConfig `json:"response"`
}

// TrafficMonitorConfig is the full set of information needed by Traffic
// Monitor to do its job of monitoring health and statistics of cache servers.
type TrafficMonitorConfig struct {
	// TrafficServers is the set of all Cache Servers which should be monitored
	// by the Traffic Monitor.
	TrafficServers []TrafficServer `json:"trafficServers,omitempty"`
	// CacheGroups is a collection of Cache Group Names associated with their
	// geographic coordinates, for use in determining Cache Group availability.
	CacheGroups []TMCacheGroup `json:"cacheGroups,omitempty"`
	// Config is a mapping of arbitrary configuration parameters to their
	// respective values. While there is no defined structure to this map,
	// certain configuration parameters have specifically defined semantics
	// which may be found in the Traffic Monitor documentation.
	Config map[string]interface{} `json:"config,omitempty"`
	// TrafficMonitors is the set of ALL Traffic Monitors (including whatever
	// Traffic Monitor requested the endpoint (if indeed one did)) which is
	// used to determine quorum and peer polling behavior.
	TrafficMonitors []TrafficMonitor `json:"trafficMonitors,omitempty"`
	// DeliveryServices is the set of all Delivery Services within the
	// monitored CDN, which are used to determine Delivery Service-level
	// statistics.
	DeliveryServices []TMDeliveryService `json:"deliveryServices,omitempty"`
	// Profiles is the set of Profiles in use by any and all monitored cache
	// servers (those given in TrafficServers), which are stored here to
	// avoid potentially lengthy reiteration.
	Profiles []TMProfile `json:"profiles,omitempty"`
}

// ToLegacyConfig converts TrafficMonitorConfig to LegacyTrafficMonitorConfig.
func (tmc *TrafficMonitorConfig) ToLegacyConfig() LegacyTrafficMonitorConfig {
	var servers []LegacyTrafficServer
	for _, s := range tmc.TrafficServers {
		servers = append(servers, s.ToLegacyServer())
	}

	legacy := LegacyTrafficMonitorConfig{
		CacheGroups:      tmc.CacheGroups,
		Config:           tmc.Config,
		TrafficMonitors:  tmc.TrafficMonitors,
		Profiles:         tmc.Profiles,
		DeliveryServices: tmc.DeliveryServices,
		TrafficServers:   servers,
	}
	return legacy
}

// LegacyTrafficMonitorConfig represents TrafficMonitorConfig for ATC versions before 5.0.
type LegacyTrafficMonitorConfig struct {
	TrafficServers   []LegacyTrafficServer  `json:"trafficServers,omitempty"`
	CacheGroups      []TMCacheGroup         `json:"cacheGroups,omitempty"`
	Config           map[string]interface{} `json:"config,omitempty"`
	TrafficMonitors  []TrafficMonitor       `json:"trafficMonitors,omitempty"`
	DeliveryServices []TMDeliveryService    `json:"deliveryServices,omitempty"`
	Profiles         []TMProfile            `json:"profiles,omitempty"`
}

// TrafficMonitorConfigMap is a representation of a TrafficMonitorConfig using
// unique values as map keys.
type TrafficMonitorConfigMap struct {
	// TrafficServer is a map of cache server hostnames to TrafficServer
	// objects.
	//
	// WARNING: Cache server hostnames are NOT guaranteed to be unique, so when
	// more than one cache server with the same hostname exists, the two CANNOT
	// coexist within this structure; only one will appear and, when
	// constructed using TrafficMonitorTransformToMap, which cache server does
	// exist in the generated map is undefined.
	TrafficServer map[string]TrafficServer
	// CacheGroup is a map of Cache Group Names to TMCacheGroup objects.
	CacheGroup map[string]TMCacheGroup
	// Config is a mapping of arbitrary configuration parameters to their
	// respective values. While there is no defined structure to this map,
	// certain configuration parameters have specifically defined semantics
	// which may be found in the Traffic Monitor documentation.
	Config map[string]interface{}
	// TrafficMonitor is a map of Traffic Monitor hostnames to TrafficMonitor
	// objects.
	//
	// WARNING: Traffic Monitor hostnames are NOT guaranteed to be unique, so
	// when more than one cache server with the same hostname exists, the two
	// CANNOT coexist within this structure; only one will appear and, when
	// constructed using TrafficMonitorTransformToMap, which Traffic Monitor
	// does exist in the generated map is undefined.
	TrafficMonitor map[string]TrafficMonitor
	// DeliveryService is a map of Delivery Service XMLIDs to TMDeliveryService
	// objects.
	DeliveryService map[string]TMDeliveryService
	// Profile is a map of Profile Names to TMProfile objects.
	Profile map[string]TMProfile
}

// LegacyTrafficMonitorConfigMap ...
type LegacyTrafficMonitorConfigMap struct {
	TrafficServer   map[string]LegacyTrafficServer
	CacheGroup      map[string]TMCacheGroup
	Config          map[string]interface{}
	TrafficMonitor  map[string]TrafficMonitor
	DeliveryService map[string]TMDeliveryService
	Profile         map[string]TMProfile
}

// TrafficMonitor is a structure containing enough information about a Traffic
// Monitor instance for another Traffic Monitor to use it for quorums and peer
// polling.
type TrafficMonitor struct {
	// Port is the port on which the Traffic Monitor listens for HTTP requests.
	Port int `json:"port"`
	// IP6 is the Traffic Monitor's IPv6 address.
	IP6 string `json:"ip6"`
	// IP is the Traffic Monitor's IPv4 address.
	IP string `json:"ip"`
	// HostName is the Traffic Monitor's hostname.
	HostName string `json:"hostName"`
	// FQDN is the Fully Qualified Domain Name of the Traffic Monitor server.
	FQDN string `json:"fqdn"`
	// Profile is the Name of the Profile used by the Traffic Monitor.
	Profile string `json:"profile"`
	// Location is the Name of the Cache Group to which the Traffic Monitor
	// belongs - called "Location" for legacy reasons.
	Location string `json:"location"`
	// ServerStatus is the Name of the Status of the Traffic Monitor.
	ServerStatus string `json:"status"`
}

// TMCacheGroup contains all of the information about a Cache Group necessary
// for Traffic Monitor to do its job of monitoring health and statistics.
type TMCacheGroup struct {
	Name        string                `json:"name"`
	Coordinates MonitoringCoordinates `json:"coordinates"`
}

// MonitoringCoordinates holds a coordinate pair for inclusion as a field in
// TMCacheGroup.
type MonitoringCoordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// TMDeliveryService is all of the information about a Delivery Service
// necessary for Traffic Monitor to do its job of monitoring health and
// statistics.
type TMDeliveryService struct {
	XMLID              string `json:"xmlId"`
	TotalTPSThreshold  int64  `json:"TotalTpsThreshold"`
	ServerStatus       string `json:"status"`
	TotalKbpsThreshold int64  `json:"TotalKbpsThreshold"`
}

// TMProfile is primarily a collection of the Parameters with special meaning
// to Traffic Monitor for a Profile of one of the monitored cache servers
// and/or other Traffic Monitors, along with some identifying information for
// the Profile.
type TMProfile struct {
	Parameters TMParameters `json:"parameters"`
	Name       string       `json:"name"`
	Type       string       `json:"type"`
}

// TMParameters is a structure containing all of the Parameters with special
// meaning to Traffic Monitor. For specifics regarding each Parameter, refer to
// the official documentation.
// TODO change TO to return this struct, so a custom UnmarshalJSON isn't necessary.
type TMParameters struct {
	HealthConnectionTimeout int    `json:"health.connection.timeout"`
	HealthPollingURL        string `json:"health.polling.url"`
	HealthPollingFormat     string `json:"health.polling.format"`
	HealthPollingType       string `json:"health.polling.type"`
	HistoryCount            int    `json:"history.count"`
	MinFreeKbps             int64
	Thresholds              map[string]HealthThreshold `json:"health_threshold"`
	AggregateThresholds     map[string]HealthThreshold `json:"health_threshold_aggregate"`
}

const DefaultHealthThresholdComparator = "<"

// HealthThreshold describes some value against which to compare health
// measurements to determine if a cache server is healthy.
type HealthThreshold struct {
	// Val is the actual, numeric, threshold value.
	Val float64
	// Comparator is the comparator used to compare the Val to the monitored
	// value. One of '=', '>', '<', '>=', or '<=' - other values are invalid.
	Comparator string // TODO change to enum?
}

// String implements the fmt.Stringer interface.
func (t HealthThreshold) String() string {
	return fmt.Sprintf("%s%f", t.Comparator, t.Val)
}

// strToThreshold takes a string like ">=42" and returns a HealthThreshold with
// a Val of `42` and a Comparator of `">="`. If no comparator exists,
// `DefaultHealthThresholdComparator` is used. If the string does not match
// "(>|<|)(=|)\d+" an error is returned.
func strToThreshold(s string) (HealthThreshold, error) {
	comparators := []string{"=", ">", "<", ">=", "<="}
	for _, comparator := range comparators {
		if strings.HasPrefix(s, comparator) {
			valStr := s[len(comparator):]
			val, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				return HealthThreshold{}, fmt.Errorf("invalid threshold: NaN")
			}
			return HealthThreshold{Val: val, Comparator: comparator}, nil
		}
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return HealthThreshold{}, fmt.Errorf("invalid threshold: NaN")
	}
	return HealthThreshold{Val: val, Comparator: DefaultHealthThresholdComparator}, nil
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (params *TMParameters) UnmarshalJSON(bytes []byte) (err error) {
	raw := map[string]interface{}{}
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	if vi, ok := raw["health.connection.timeout"]; ok {
		if v, ok := vi.(float64); !ok {
			return fmt.Errorf("Unmarshalling TMParameters health.connection.timeout expected integer, got %v", vi)
		} else {
			params.HealthConnectionTimeout = int(v)
		}
	}

	if vi, ok := raw["health.polling.url"]; ok {
		if v, ok := vi.(string); !ok {
			return fmt.Errorf("Unmarshalling TMParameters health.polling.url expected string, got %v", vi)
		} else {
			params.HealthPollingURL = v
		}
	}

	if vi, ok := raw["health.polling.format"]; ok {
		if v, ok := vi.(string); !ok {
			return fmt.Errorf("Unmarshalling TMParameters health.polling.format expected string, got %v", vi)
		} else {
			params.HealthPollingFormat = v
		}
	}

	if vi, ok := raw["health.polling.type"]; ok {
		if v, ok := vi.(string); !ok {
			return fmt.Errorf("Unmarshalling TMParameters health.polling.type expected string, got %v", vi)
		} else {
			params.HealthPollingType = v
		}
	}

	if vi, ok := raw["history.count"]; ok {
		if v, ok := vi.(float64); !ok {
			return fmt.Errorf("Unmarshalling TMParameters history.count expected integer, got %v", vi)
		} else {
			params.HistoryCount = int(v)
		}
	}

	params.Thresholds = map[string]HealthThreshold{}
	params.AggregateThresholds = map[string]HealthThreshold{}
	for k, v := range raw {
		switch k {
		case AvailableBandwidthInKbpsThresholdName:
			fallthrough
		case AvailableBandwidthInMbpsThresholdName:
			fallthrough
		case BandwidthThresholdName:
			fallthrough
		case KbpsThresholdName:
			fallthrough
		case GbpsThresholdName:
			fallthrough
		case LoadavgThresholdName:
			fallthrough
		case MaxKbpsThresholdName:
			stat := k[len(ThresholdPrefix):]
			vStr := fmt.Sprintf("%v", v)
			if t, err := strToThreshold(vStr); err != nil {
				return fmt.Errorf("Unmarshalling TMParameters `%s` parameter value not of the form `(>|)(=|)\\d+`: stat '%s' value '%v'", ThresholdPrefix, k, v)
			} else {
				params.AggregateThresholds[stat] = t
			}
		default:
			if strings.HasPrefix(k, ThresholdPrefix) {
				stat := k[len(ThresholdPrefix):]
				vStr := fmt.Sprintf("%v", v) // allows string or numeric JSON types. TODO check if a type switch is faster.
				if t, err := strToThreshold(vStr); err != nil {
					return fmt.Errorf("Unmarshalling TMParameters `%s` parameter value not of the form `(>|)(=|)\\d+`: stat '%s' value '%v'", ThresholdPrefix, k, v)
				} else {
					params.Thresholds[stat] = t
				}
			}
		}
	}
	return nil
}

func TrafficMonitorTransformToMap(tmConfig *TrafficMonitorConfig) (*TrafficMonitorConfigMap, error) {
	var tm TrafficMonitorConfigMap

	tm.TrafficServer = make(map[string]TrafficServer, len(tmConfig.TrafficServers))
	tm.CacheGroup = make(map[string]TMCacheGroup, len(tmConfig.CacheGroups))
	tm.Config = make(map[string]interface{}, len(tmConfig.Config))
	tm.TrafficMonitor = make(map[string]TrafficMonitor, len(tmConfig.TrafficMonitors))
	tm.DeliveryService = make(map[string]TMDeliveryService, len(tmConfig.DeliveryServices))
	tm.Profile = make(map[string]TMProfile, len(tmConfig.Profiles))

	for _, trafficServer := range tmConfig.TrafficServers {
		tm.TrafficServer[trafficServer.HostName] = trafficServer
	}

	for _, cacheGroup := range tmConfig.CacheGroups {
		tm.CacheGroup[cacheGroup.Name] = cacheGroup
	}

	for parameterKey, parameterVal := range tmConfig.Config {
		tm.Config[parameterKey] = parameterVal
	}

	for _, trafficMonitor := range tmConfig.TrafficMonitors {
		tm.TrafficMonitor[trafficMonitor.HostName] = trafficMonitor
	}

	for _, deliveryService := range tmConfig.DeliveryServices {
		tm.DeliveryService[deliveryService.XMLID] = deliveryService
	}

	for _, profile := range tmConfig.Profiles {
		bwThreshold := profile.Parameters.Thresholds["availableBandwidthInKbps"]
		profile.Parameters.MinFreeKbps = int64(bwThreshold.Val)
		tm.Profile[profile.Name] = profile
	}

	return &tm, MonitorConfigValid(&tm)
}

func MonitorConfigValid(cfg *TrafficMonitorConfigMap) error {
	if cfg == nil {
		return errors.New("MonitorConfig is nil")
	}
	if len(cfg.TrafficServer) == 0 {
		return errors.New("MonitorConfig.TrafficServer empty (is the monitoring.json an empty object?)")
	}
	if len(cfg.CacheGroup) == 0 {
		return errors.New("MonitorConfig.CacheGroup empty")
	}
	if len(cfg.TrafficMonitor) == 0 {
		return errors.New("MonitorConfig.TrafficMonitor empty")
	}
	// TODO uncomment this, when TO is fixed to include DeliveryServices.
	// See https://github.com/apache/trafficcontrol/issues/3528
	// if len(cfg.DeliveryService) == 0 {
	// 	return errors.New("MonitorConfig.DeliveryService empty")
	// }
	if len(cfg.Profile) == 0 {
		return errors.New("MonitorConfig.Profile empty")
	}

	if intervalI, ok := cfg.Config["peers.polling.interval"]; !ok {
		return errors.New(`MonitorConfig.Config["peers.polling.interval"] missing, peers.polling.interval parameter required`)
	} else if _, ok := intervalI.(float64); !ok {
		return fmt.Errorf(`MonitorConfig.Config["peers.polling.interval"] '%v' not a number, parameter peers.polling.interval must be a number`, intervalI)
	}

	if intervalI, ok := cfg.Config["health.polling.interval"]; !ok {
		return errors.New(`MonitorConfig.Config["health.polling.interval"] missing, health.polling.interval parameter required`)
	} else if _, ok := intervalI.(float64); !ok {
		return fmt.Errorf(`MonitorConfig.Config["health.polling.interval"] '%v' not a number, parameter health.polling.interval must be a number`, intervalI)
	}

	return nil
}

func LegacyTrafficMonitorTransformToMap(tmConfig *LegacyTrafficMonitorConfig) (*LegacyTrafficMonitorConfigMap, error) {
	var tm LegacyTrafficMonitorConfigMap

	tm.TrafficServer = make(map[string]LegacyTrafficServer)
	tm.CacheGroup = make(map[string]TMCacheGroup)
	tm.Config = make(map[string]interface{})
	tm.TrafficMonitor = make(map[string]TrafficMonitor)
	tm.DeliveryService = make(map[string]TMDeliveryService)
	tm.Profile = make(map[string]TMProfile)

	for _, trafficServer := range tmConfig.TrafficServers {
		tm.TrafficServer[trafficServer.HostName] = trafficServer
	}

	for _, cacheGroup := range tmConfig.CacheGroups {
		tm.CacheGroup[cacheGroup.Name] = cacheGroup
	}

	for parameterKey, parameterVal := range tmConfig.Config {
		tm.Config[parameterKey] = parameterVal
	}

	for _, trafficMonitor := range tmConfig.TrafficMonitors {
		tm.TrafficMonitor[trafficMonitor.HostName] = trafficMonitor
	}

	for _, deliveryService := range tmConfig.DeliveryServices {
		tm.DeliveryService[deliveryService.XMLID] = deliveryService
	}

	for _, profile := range tmConfig.Profiles {
		bwThreshold := profile.Parameters.Thresholds["availableBandwidthInKbps"]
		profile.Parameters.MinFreeKbps = int64(bwThreshold.Val)
		tm.Profile[profile.Name] = profile
	}

	return &tm, LegacyMonitorConfigValid(&tm)
}

func LegacyMonitorConfigValid(cfg *LegacyTrafficMonitorConfigMap) error {
	if cfg == nil {
		return errors.New("MonitorConfig is nil")
	}
	if len(cfg.TrafficServer) == 0 {
		return errors.New("MonitorConfig.TrafficServer empty (is the monitoring.json an empty object?)")
	}
	if len(cfg.CacheGroup) == 0 {
		return errors.New("MonitorConfig.CacheGroup empty")
	}
	if len(cfg.TrafficMonitor) == 0 {
		return errors.New("MonitorConfig.TrafficMonitor empty")
	}
	// TODO uncomment this, when TO is fixed to include DeliveryServices.
	// See https://github.com/apache/trafficcontrol/issues/3528
	// if len(cfg.DeliveryService) == 0 {
	// 	return errors.New("MonitorConfig.DeliveryService empty")
	// }
	if len(cfg.Profile) == 0 {
		return errors.New("MonitorConfig.Profile empty")
	}

	if intervalI, ok := cfg.Config["peers.polling.interval"]; !ok {
		return errors.New(`MonitorConfig.Config["peers.polling.interval"] missing, peers.polling.interval parameter required`)
	} else if _, ok := intervalI.(float64); !ok {
		return fmt.Errorf(`MonitorConfig.Config["peers.polling.interval"] '%v' not a number, parameter peers.polling.interval must be a number`, intervalI)
	}

	if intervalI, ok := cfg.Config["health.polling.interval"]; !ok {
		return errors.New(`MonitorConfig.Config["health.polling.interval"] missing, health.polling.interval parameter required`)
	} else if _, ok := intervalI.(float64); !ok {
		return fmt.Errorf(`MonitorConfig.Config["health.polling.interval"] '%v' not a number, parameter health.polling.interval must be a number`, intervalI)
	}

	return nil
}

type HealthData struct {
	TotalOffline uint64                 `json:"totalOffline"`
	TotalOnline  uint64                 `json:"totalOnline"`
	CacheGroups  []HealthDataCacheGroup `json:"cachegroups"`
}

type HealthDataCacheGroup struct {
	Offline int64          `json:"offline"`
	Online  int64          `json:"online"`
	Name    CacheGroupName `json:"name"`
}
