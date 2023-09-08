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
	"time"

	jsoniter "github.com/json-iterator/go"
)

// ThresholdPrefix is the prefix of all Names of Parameters used to define
// monitoring thresholds.
const ThresholdPrefix = "health.threshold."

// These are the names of statistics that can be used in thresholds for server
// health.
const (
	StatNameKBPS      = "kbps"
	StatNameMaxKBPS   = "maxKbps"
	StatNameBandwidth = "bandwidth"
)

// TMConfigResponse is the response to requests made to the
// cdns/{{Name}}/configs/monitoring endpoint of the Traffic Ops API.
type TMConfigResponse struct {
	Response TrafficMonitorConfig `json:"response"`
	Alerts
}

// LegacyTMConfigResponse was the response to requests made to the
// cdns/{{Name}}/configs/monitoring endpoint of the Traffic Ops API in older
// API versions.
//
// Deprecated: New code should use TMConfigResponse instead.
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
	// Topologies is the set of topologies defined in Traffic Ops, consisting
	// of just the EDGE_LOC-type cachegroup nodes.
	Topologies map[string]CRConfigTopology `json:"topologies,omitempty"`
}

const healthThresholdAvailableBandwidthInKbps = "availableBandwidthInKbps"
const healthThresholdLoadAverage = "loadavg"
const healthThresholdQueryTime = "queryTime"

// ToLegacyConfig converts TrafficMonitorConfig to LegacyTrafficMonitorConfig.
//
// Deprecated: LegacyTrafficMonitoryConfig is deprecated. New code should just
// use TrafficMonitorConfig instead.
func (tmc *TrafficMonitorConfig) ToLegacyConfig() LegacyTrafficMonitorConfig {
	var servers []LegacyTrafficServer
	for _, s := range tmc.TrafficServers {
		servers = append(servers, s.ToLegacyServer())
	}

	for profileIndex, profile := range tmc.Profiles {
		thresholds := profile.Parameters.Thresholds
		if _, exists := thresholds[healthThresholdAvailableBandwidthInKbps]; exists {
			tmc.Profiles[profileIndex].Parameters.AvailableBandwidthInKbps = thresholds[healthThresholdAvailableBandwidthInKbps].String()
			delete(tmc.Profiles[profileIndex].Parameters.Thresholds, healthThresholdAvailableBandwidthInKbps)
		}
		if _, exists := thresholds[healthThresholdLoadAverage]; exists {
			tmc.Profiles[profileIndex].Parameters.LoadAverage = thresholds[healthThresholdLoadAverage].String()
			delete(tmc.Profiles[profileIndex].Parameters.Thresholds, healthThresholdLoadAverage)
		}
		if _, exists := thresholds[healthThresholdQueryTime]; exists {
			tmc.Profiles[profileIndex].Parameters.QueryTime = thresholds[healthThresholdQueryTime].String()
			delete(tmc.Profiles[profileIndex].Parameters.Thresholds, healthThresholdQueryTime)
		}
	}

	legacy := LegacyTrafficMonitorConfig{
		CacheGroups:      tmc.CacheGroups,
		Config:           tmc.Config,
		TrafficMonitors:  tmc.TrafficMonitors,
		Profiles:         tmc.Profiles,
		DeliveryServices: tmc.DeliveryServices,
		TrafficServers:   servers,
		Topologies:       tmc.Topologies,
	}
	return legacy
}

// LegacyTrafficMonitorConfig represents TrafficMonitorConfig for ATC versions
// before 5.0.
type LegacyTrafficMonitorConfig struct {
	TrafficServers   []LegacyTrafficServer       `json:"trafficServers,omitempty"`
	CacheGroups      []TMCacheGroup              `json:"cacheGroups,omitempty"`
	Config           map[string]interface{}      `json:"config,omitempty"`
	TrafficMonitors  []TrafficMonitor            `json:"trafficMonitors,omitempty"`
	DeliveryServices []TMDeliveryService         `json:"deliveryServices,omitempty"`
	Profiles         []TMProfile                 `json:"profiles,omitempty"`
	Topologies       map[string]CRConfigTopology `json:"topologies,omitempty"`
}

// Upgrade converts a legacy TM Config to the newer structure.
//
// Deprecated: LegacyTrafficMonitoryConfig is deprecated. New code should just
// use TrafficMonitorConfig instead.
func (s *LegacyTrafficMonitorConfig) Upgrade() *TrafficMonitorConfig {
	upgraded := TrafficMonitorConfig{
		CacheGroups:      s.CacheGroups,
		Config:           s.Config,
		DeliveryServices: s.DeliveryServices,
		Profiles:         s.Profiles,
		TrafficMonitors:  s.TrafficMonitors,
		TrafficServers:   make([]TrafficServer, 0, len(s.TrafficServers)),
		Topologies:       s.Topologies,
	}
	for _, ts := range s.TrafficServers {
		upgraded.TrafficServers = append(upgraded.TrafficServers, ts.Upgrade())
	}
	return &upgraded
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
	// Topology is a map of Topology names to CRConfigTopology structs.
	Topology map[string]CRConfigTopology
}

// ToLegacy converts a Stats to a LegacyStats.
//
// This returns a list of descriptions of which - if any - cache servers were
// skipped in the conversion and why, as well as the converted LegacyStats.
//
// This creates a "shallow" copy of most properties of the Stats.
//
// Deprecated: LegacyStats is deprecated. New code should just use Stats
// instead.
func (s *Stats) ToLegacy(monitorConfig TrafficMonitorConfigMap) ([]string, LegacyStats) {
	legacyStats := LegacyStats{
		CommonAPIData: s.CommonAPIData,
		Caches:        make(map[CacheName]map[string][]ResultStatVal, len(s.Caches)),
	}
	skippedCaches := []string{}

	for cacheName, cache := range s.Caches {
		ts, ok := monitorConfig.TrafficServer[cacheName]
		if !ok {
			skippedCaches = append(skippedCaches, "Cache "+cacheName+" does not exist in the "+
				"TrafficMonitorConfigMap")
			continue
		}
		legacyInterface, err := InterfaceInfoToLegacyInterfaces(ts.Interfaces)
		if err != nil {
			skippedCaches = append(skippedCaches, "Cache "+cacheName+": unable to convert to legacy "+
				"interfaces: "+err.Error())
			continue
		}
		if legacyInterface.InterfaceName == nil {
			skippedCaches = append(skippedCaches, "Cache "+cacheName+": computed legacy interface "+
				"does not have a name")
			continue
		}
		monitorInterfaceStats, ok := cache.Interfaces[*legacyInterface.InterfaceName]
		if !ok && len(cache.Interfaces) > 0 {
			skippedCaches = append(skippedCaches, "Cache "+cacheName+" does not contain interface "+
				*legacyInterface.InterfaceName)
			continue
		}
		length := len(monitorInterfaceStats) + len(cache.Stats)
		legacyStats.Caches[CacheName(cacheName)] = make(map[string][]ResultStatVal, length)
		for statName, stat := range cache.Stats {
			legacyStats.Caches[CacheName(cacheName)][statName] = stat
		}
		for statName, stat := range monitorInterfaceStats {
			legacyStats.Caches[CacheName(cacheName)][statName] = stat
		}
	}

	return skippedCaches, legacyStats
}

// ServerStats is a representation of cache server statistics as present in the
// TM API.
type ServerStats struct {
	// Interfaces contains statistics specific to each monitored interface
	// of the cache server.
	Interfaces map[string]map[string][]ResultStatVal `json:"interfaces"`
	// Stats contains statistics regarding the cache server in general.
	Stats map[string][]ResultStatVal `json:"stats"`
}

// Stats is designed for returning via the API. It contains result history
// for each cache, as well as common API data.
type Stats struct {
	CommonAPIData
	// Caches is a map of cache server hostnames to groupings of statistics
	// regarding each cache server and all of its separate network interfaces.
	Caches map[string]ServerStats `json:"caches"`
}

// LegacyStats is designed for returning via the API. It contains result
// history for each cache server, as well as common API data.
//
// Deprecated: This structure is incapable of representing interface-level
// stats, so new code should use Stats instead.
type LegacyStats struct {
	CommonAPIData
	Caches map[CacheName]map[string][]ResultStatVal `json:"caches"`
}

// CommonAPIData contains generic data common to most Traffic Monitor API
// endpoints.
type CommonAPIData struct {
	QueryParams string `json:"pp"`
	DateStr     string `json:"date"`
}

// ResultStatVal is the value of an individual stat returned from a poll.
// JSON values are all strings, for the TM1.0 /publish/CacheStats API.
type ResultStatVal struct {
	// Span is the number of polls this stat has been the same. For example,
	// if History is set to 100, and the last 50 polls had the same value for
	// this stat (but none of the previous 50 were the same), this stat's map
	// value slice will actually contain 51 entries, and the first entry will
	// have the value, the time of the last poll, and a Span of 50.
	// Assuming the poll time is every 8 seconds, users will then know, looking
	// at the Span, that the value was unchanged for the last 50*8=400 seconds.
	Span uint64 `json:"span"`
	// Time is the time this stat was returned.
	Time time.Time   `json:"time"`
	Val  interface{} `json:"value"`
}

// MarshalJSON implements the encoding/json.Marshaler interface.
func (t *ResultStatVal) MarshalJSON() ([]byte, error) {
	v := struct {
		Val  string `json:"value"`
		Time int64  `json:"time"`
		Span uint64 `json:"span"`
	}{
		Val:  fmt.Sprintf("%v", t.Val),
		Time: t.Time.UnixNano() / 1000000, // ms since the epoch
		Span: t.Span,
	}
	json := jsoniter.ConfigFastest // TODO make configurable
	return json.Marshal(&v)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (t *ResultStatVal) UnmarshalJSON(data []byte) error {
	v := struct {
		Val  string `json:"value"`
		Time int64  `json:"time"`
		Span uint64 `json:"span"`
	}{}
	json := jsoniter.ConfigFastest // TODO make configurable
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	t.Time = time.Unix(0, v.Time*1000000)
	t.Val = v.Val
	t.Span = v.Span
	return nil
}

// Valid returns a non-nil error if the configuration map is invalid.
//
// A configuration map is considered invalid if:
//
// - it is nil,
//
// - it has no CacheGroups,
//
// - it has no Profiles,
//
// - it has no Traffic Monitors,
//
// - it has no Traffic Servers,
//
// - the Config mapping has no 'peers.polling.interval' key,
//
// - or the Config mapping has no 'health.polling.interval' key.
func (cfg *TrafficMonitorConfigMap) Valid() error {
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
	if len(cfg.DeliveryService) == 0 {
		return errors.New("MonitorConfig.DeliveryService empty")
	}
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

// LegacyTrafficMonitorConfigMap is a representation of a
// LegacyTrafficMonitorConfig using unique values as map keys.
//
// Deprecated: This structure is incapable of representing per-interface
// configuration information for servers, so new code should use
// TrafficMonitorConfigMap instead.
type LegacyTrafficMonitorConfigMap struct {
	TrafficServer   map[string]LegacyTrafficServer
	CacheGroup      map[string]TMCacheGroup
	Config          map[string]interface{}
	TrafficMonitor  map[string]TrafficMonitor
	DeliveryService map[string]TMDeliveryService
	Profile         map[string]TMProfile
}

// Upgrade returns a TrafficMonitorConfigMap that is equivalent to this legacy
// configuration map.
//
// Note that all fields except TrafficServer are "shallow" copies, so modifying
// the original will impact the upgraded copy.
//
// Deprecated: LegacyTrafficMonitorConfigMap is deprecated.
func (c *LegacyTrafficMonitorConfigMap) Upgrade() *TrafficMonitorConfigMap {
	upgraded := TrafficMonitorConfigMap{
		CacheGroup:      c.CacheGroup,
		Config:          c.Config,
		DeliveryService: c.DeliveryService,
		Profile:         c.Profile,
		TrafficMonitor:  c.TrafficMonitor,
		TrafficServer:   make(map[string]TrafficServer, len(c.TrafficServer)),
	}

	for k, ts := range c.TrafficServer {
		upgraded.TrafficServer[k] = ts.Upgrade()
	}
	return &upgraded
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
	Location string `json:"cachegroup"`
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
	XMLID              string   `json:"xmlId"`
	TotalTPSThreshold  int64    `json:"TotalTpsThreshold"`
	ServerStatus       string   `json:"status"`
	TotalKbpsThreshold int64    `json:"TotalKbpsThreshold"`
	Topology           string   `json:"topology"`
	Type               string   `json:"type"`
	HostRegexes        []string `json:"hostRegexes,omitempty"`
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
// meaning to Traffic Monitor.
//
// For specifics regarding each Parameter, refer to the official documentation.
//
// TODO change TO to return this struct, so a custom UnmarshalJSON isn't
// necessary.
type TMParameters struct {
	HealthConnectionTimeout int    `json:"health.connection.timeout"`
	HealthPollingURL        string `json:"health.polling.url"`
	HealthPollingFormat     string `json:"health.polling.format"`
	HealthPollingType       string `json:"health.polling.type"`
	HistoryCount            int    `json:"history.count"`
	MinFreeKbps             int64
	// HealthThresholdJSONParameters contains the Parameters contained in the
	// Thresholds field, formatted as individual string Parameters, rather than as
	// a JSON object.
	Thresholds map[string]HealthThreshold `json:"health_threshold,omitempty"`
	HealthThresholdJSONParameters
}

// HealthThresholdJSONParameters contains Parameters whose Thresholds must be met in order for
// Caches using the Profile containing these Parameters to be marked as Healthy.
type HealthThresholdJSONParameters struct {
	// AvailableBandwidthInKbps is The total amount of bandwidth that servers using this profile are
	// allowed, in Kilobits per second. This is a string and using comparison operators to specify
	// ranges, e.g. ">10" means "more than 10 kbps".
	AvailableBandwidthInKbps string `json:"health.threshold.availableBandwidthInKbps,omitempty"`
	// LoadAverage is the UNIX loadavg at which the server should be marked "unhealthy".
	LoadAverage string `json:"health.threshold.loadavg,omitempty"`
	// QueryTime is the highest allowed length of time for completing health queries (after
	// connection has been established) in milliseconds.
	QueryTime string `json:"health.threshold.queryTime,omitempty"`
}

// DefaultHealthThresholdComparator is the comparator used for health
// thresholds when one is not explicitly provided in the Value of a Parameter
// used to define a Threshold.
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

// StrToThreshold takes a string like ">=42" and returns a HealthThreshold with
// a Val of `42` and a Comparator of `">="`. If no comparator exists,
// `DefaultHealthThresholdComparator` is used. If the string does not match
// "(>|<|)(=|)\d+" an error is returned.
func StrToThreshold(s string) (HealthThreshold, error) {
	// The order of these is important - don't re-order without considering the
	// consequences.
	comparators := []string{">=", "<=", ">", "<", "="}
	for _, comparator := range comparators {
		if strings.HasPrefix(s, comparator) {
			valStr := s[len(comparator):]
			val, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				return HealthThreshold{}, fmt.Errorf("invalid threshold: NaN (%v)", err)
			}
			return HealthThreshold{Val: val, Comparator: comparator}, nil
		}
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return HealthThreshold{}, fmt.Errorf("invalid threshold: NaN (%v)", err)
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

	params.Thresholds = make(map[string]HealthThreshold, len(raw))
	for k, v := range raw {
		if strings.HasPrefix(k, ThresholdPrefix) {
			stat := k[len(ThresholdPrefix):]
			vStr := fmt.Sprintf("%v", v) // allows string or numeric JSON types. TODO check if a type switch is faster.
			if t, err := StrToThreshold(vStr); err != nil {
				return fmt.Errorf("Unmarshalling TMParameters `%s` parameter value not of the form `(>|)(=|)\\d+`: stat '%s' value '%v': %v", ThresholdPrefix, k, v, err)
			} else {
				params.Thresholds[stat] = t
			}
		}
	}
	return nil
}

// TrafficMonitorTransformToMap converts the given TrafficMonitorConfig to a
// TrafficMonitorConfigMap.
//
// This also implicitly calls Valid on the TrafficMonitorConfigMap before
// returning it, and gives back whatever value that returns as the error return
// value.
func TrafficMonitorTransformToMap(tmConfig *TrafficMonitorConfig) (*TrafficMonitorConfigMap, error) {
	var tm TrafficMonitorConfigMap

	tm.TrafficServer = make(map[string]TrafficServer, len(tmConfig.TrafficServers))
	tm.CacheGroup = make(map[string]TMCacheGroup, len(tmConfig.CacheGroups))
	tm.Config = make(map[string]interface{}, len(tmConfig.Config))
	tm.TrafficMonitor = make(map[string]TrafficMonitor, len(tmConfig.TrafficMonitors))
	tm.DeliveryService = make(map[string]TMDeliveryService, len(tmConfig.DeliveryServices))
	tm.Profile = make(map[string]TMProfile, len(tmConfig.Profiles))
	tm.Topology = tmConfig.Topologies

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

	return &tm, tm.Valid()
}

// LegacyTrafficMonitorTransformToMap converts the given
// LegacyTrafficMonitorConfig to a LegacyTrafficMonitorConfigMap.
//
// This also implicitly calls LegacyMonitorConfigValid on the
// LegacyTrafficMonitorConfigMap before returning it, and gives back whatever
// value that returns as the error return value.
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

// LegacyMonitorConfigValid checks the validity of the passed
// LegacyTrafficMonitorConfigMap, returning an error if it is invalid.
//
// Deprecated: LegacyTrafficMonitorConfigMap is deprecated, new code should use
// TrafficMonitorConfigMap instead.
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

// HealthData is a representation of all of the health information for a CDN
// or Delivery Service.
//
// This is the type of the `response` property of responses from Traffic Ops to
// GET requests made to its /deliveryservices/{{ID}}/health and /cdns/health
// API endpoints.
type HealthData struct {
	TotalOffline uint64                 `json:"totalOffline"`
	TotalOnline  uint64                 `json:"totalOnline"`
	CacheGroups  []HealthDataCacheGroup `json:"cachegroups"`
}

// HealthDataCacheGroup holds health information specific to a particular Cache
// Group.
type HealthDataCacheGroup struct {
	Offline int64          `json:"offline"`
	Online  int64          `json:"online"`
	Name    CacheGroupName `json:"name"`
}
