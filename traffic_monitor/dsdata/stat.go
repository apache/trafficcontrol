package dsdata

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
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/srvhttp"
)

// Filter encapsulates functions to filter a given set of Stats, e.g. from HTTP query parameters.
// TODO combine with cache.Filter?
type Filter interface {
	UseStat(name string) bool
	UseDeliveryService(name tc.DeliveryServiceName) bool
	WithinStatHistoryMax(int) bool
}

// StatName is the name of a stat.
type StatName string

// StatOld is the old JSON representation of a stat, from Traffic Monitor 1.0.
type StatOld struct {
	Time  int64       `json:"time"`
	Value interface{} `json:"value"`
	Span  int         `json:"span,omitempty"`  // TODO set? remove?
	Index int         `json:"index,omitempty"` // TODO set? remove?
}

// StatsOld is the old JSON representation of stats, from Traffic Monitor 1.0. It is designed to be serialized and returns from an API, and includes stat history for each delivery service, as well as data common to most endpoints.
type StatsOld struct {
	DeliveryService map[tc.DeliveryServiceName]map[StatName][]StatOld `json:"deliveryService"`
	tc.CommonAPIData
}

// StatsReadonly is a read-only interface for delivery service Stats, designed to be passed to multiple goroutine readers.
type StatsReadonly interface {
	Get(tc.DeliveryServiceName) (StatReadonly, bool)
	JSON(Filter, url.Values) StatsOld
}

// StatReadonly is a read-only interface for a delivery service Stat, designed to be passed to multiple goroutine readers.
type StatReadonly interface {
	Copy() *Stat
	Common() StatCommonReadonly
	CacheGroup(name tc.CacheGroupName) (*StatCacheStats, bool)
	Type(name tc.CacheType) (*StatCacheStats, bool)
	Total() *StatCacheStats
}

// StatCommonReadonly is a read-only interface for a delivery service's common Stat data, designed to be passed to multiple goroutine readers.
type StatCommonReadonly interface {
	Copy() StatCommon
	CachesConfigured() StatInt
	CachesReportingNames() []tc.CacheName
	Error() StatString
	Status() StatString
	Healthy() StatBool
	Available() StatBool
	CachesAvailable() StatInt
}

// StatMeta includes metadata about a particular stat.
type StatMeta struct {
	Time int64 `json:"time"`
}

// StatFloat is a float stat, combined with its metadata
type StatFloat struct {
	StatMeta
	Value float64 `json:"value"`
}

// StatBool is a boolean stat, combined with its metadata
type StatBool struct {
	StatMeta
	Value bool `json:"value"`
}

// StatInt is an integer stat, combined with its metadata
type StatInt struct {
	StatMeta
	Value int64 `json:"value"`
}

// StatString is a string stat, combined with its metadata
type StatString struct {
	StatMeta
	Value string `json:"value"`
}

// StatCommon contains stat data common to most delivery service stats.
type StatCommon struct {
	CachesConfiguredNum StatInt               `json:"caches_configured"`
	CachesReporting     map[tc.CacheName]bool `json:"caches_reporting"`
	ErrorStr            StatString            `json:"error_string"`
	StatusStr           StatString            `json:"status"`
	IsHealthy           StatBool              `json:"is_healthy"`
	IsAvailable         StatBool              `json:"is_available"`
	CachesAvailableNum  StatInt               `json:"caches_available"`
	CachesDisabled      []string              `json:"disabled_locations"`
}

// Copy returns a deep copy of this StatCommon object.
func (a StatCommon) Copy() StatCommon {
	b := a
	for k, v := range a.CachesReporting {
		b.CachesReporting[k] = v
	}
	b.CachesDisabled = make([]string, len(a.CachesDisabled), len(a.CachesDisabled))
	for i, v := range a.CachesDisabled {
		b.CachesDisabled[i] = v
	}
	return b
}

// CachesConfigured returns the number of caches configured for this delivery service stat. It is part of the StatCommonReadonly interface.
func (a StatCommon) CachesConfigured() StatInt {
	return a.CachesConfiguredNum
}

// CacheReporting returns the number of caches reporting for this delivery service stat. It is part of the StatCommonReadonly interface.
func (a StatCommon) CacheReporting(name tc.CacheName) (bool, bool) {
	c, ok := a.CachesReporting[name]
	return c, ok
}

// CachesReportingNames returns the list of caches reporting for this delivery service stat. It is part of the StatCommonReadonly interface.
func (a StatCommon) CachesReportingNames() []tc.CacheName {
	names := make([]tc.CacheName, 0, len(a.CachesReporting))
	for name := range a.CachesReporting {
		names = append(names, name)
	}
	return names
}

// Error returns the error string of this delivery service stat. It is part of the StatCommonReadonly interface.
func (a StatCommon) Error() StatString {
	return a.ErrorStr
}

// Status returns the status string of this delivery service stat. It is part of the StatCommonReadonly interface.
func (a StatCommon) Status() StatString {
	return a.StatusStr
}

// Healthy returns whether this delivery service is considered healthy by this stat. It is part of the StatCommonReadonly interface.
func (a StatCommon) Healthy() StatBool {
	return a.IsHealthy
}

// Available returns whether this delivery service is considered available by this stat. It is part of the StatCommonReadonly interface.
func (a StatCommon) Available() StatBool {
	return a.IsAvailable
}

// CachesAvailable returns the number of caches available to the delivery service in this stat. It is part of the StatCommonReadonly interface.
func (a StatCommon) CachesAvailable() StatInt {
	return a.CachesAvailableNum
}

// StatCacheStats is all the stats generated by a cache.
// This may also be used for aggregate stats, for example, the summary of all cache stats for a cache group, or delivery service.
// Each stat is an array, in case there are multiple data points at different times. However, a single data point i.e. a single array member is common.
type StatCacheStats struct {
	OutBytes    StatInt    `json:"out_bytes"`
	IsAvailable StatBool   `json:"is_available"`
	Status5xx   StatInt    `json:"status_5xx"`
	Status4xx   StatInt    `json:"status_4xx"`
	Status3xx   StatInt    `json:"status_3xx"`
	Status2xx   StatInt    `json:"status_2xx"`
	InBytes     StatFloat  `json:"in_bytes"`
	Kbps        StatFloat  `json:"kbps"`
	Tps5xx      StatFloat  `json:"tps_5xx"`
	Tps4xx      StatFloat  `json:"tps_4xx"`
	Tps3xx      StatFloat  `json:"tps_3xx"`
	Tps2xx      StatFloat  `json:"tps_2xx"`
	ErrorString StatString `json:"error_string"`
	TpsTotal    StatFloat  `json:"tps_total"`
}

// Sum adds the given cache stats to this cache stats. Numeric values are summed; strings are appended.
func (a StatCacheStats) Sum(b StatCacheStats) StatCacheStats {
	return StatCacheStats{
		OutBytes:    StatInt{Value: a.OutBytes.Value + b.OutBytes.Value},
		IsAvailable: StatBool{Value: a.IsAvailable.Value || b.IsAvailable.Value},
		Status5xx:   StatInt{Value: a.Status5xx.Value + b.Status5xx.Value},
		Status4xx:   StatInt{Value: a.Status4xx.Value + b.Status4xx.Value},
		Status3xx:   StatInt{Value: a.Status3xx.Value + b.Status3xx.Value},
		Status2xx:   StatInt{Value: a.Status2xx.Value + b.Status2xx.Value},
		InBytes:     StatFloat{Value: a.InBytes.Value + b.InBytes.Value},
		Kbps:        StatFloat{Value: a.Kbps.Value + b.Kbps.Value},
		Tps5xx:      StatFloat{Value: a.Tps5xx.Value + b.Tps5xx.Value},
		Tps4xx:      StatFloat{Value: a.Tps4xx.Value + b.Tps4xx.Value},
		Tps3xx:      StatFloat{Value: a.Tps3xx.Value + b.Tps3xx.Value},
		Tps2xx:      StatFloat{Value: a.Tps2xx.Value + b.Tps2xx.Value},
		ErrorString: StatString{Value: a.ErrorString.Value + b.ErrorString.Value},
		TpsTotal:    StatFloat{Value: a.TpsTotal.Value + b.TpsTotal.Value},
	}
}

// Stat represents a complete delivery service stat, for a given poll, or at the time requested.
type Stat struct {
	CommonStats StatCommon
	CacheGroups map[tc.CacheGroupName]*StatCacheStats
	Types       map[tc.CacheType]*StatCacheStats
	Caches      map[tc.CacheName]*StatCacheStats
	TotalStats  StatCacheStats
}

// ErrNotProcessedStat indicates a stat received is not used by Traffic Monitor, nor returned by any API endpoint. Receiving this error indicates the stat has been discarded.
var ErrNotProcessedStat = errors.New("This stat is not used.")

// NewStat returns a new delivery service Stat, initializing pointer members.
func NewStat() *Stat {
	return &Stat{
		CacheGroups: map[tc.CacheGroupName]*StatCacheStats{},
		Types:       map[tc.CacheType]*StatCacheStats{},
		CommonStats: StatCommon{CachesReporting: map[tc.CacheName]bool{}},
		Caches:      map[tc.CacheName]*StatCacheStats{},
	}
}

// Copy performs a deep copy of this Stat. It does not modify, and is thus safe for multiple goroutines.
func (a Stat) Copy() *Stat {
	// TODO sync.Pool. Better yet, remove copy usage
	b := &Stat{
		CommonStats: a.CommonStats.Copy(),
		TotalStats:  a.TotalStats,
		CacheGroups: make(map[tc.CacheGroupName]*StatCacheStats, len(a.CacheGroups)),
		Types:       make(map[tc.CacheType]*StatCacheStats, len(a.Types)),
		Caches:      make(map[tc.CacheName]*StatCacheStats, len(a.Caches)),
	}
	for k, v := range a.CacheGroups {
		b.CacheGroups[k] = v
	}
	for k, v := range a.Types {
		b.Types[k] = v
	}
	for k, v := range a.Caches {
		b.Caches[k] = v
	}
	return b
}

// Common returns the common stat data for this stat. It is part of the StatCommonReadonly interface.
func (a *Stat) Common() StatCommonReadonly {
	return a.CommonStats
}

// CacheGroup returns the data for the given cachegroup in this stat. It is part of the StatCommonReadonly interface.
func (a *Stat) CacheGroup(name tc.CacheGroupName) (*StatCacheStats, bool) {
	c, ok := a.CacheGroups[name]
	return c, ok
}

// Type returns the aggregated data for the given cache type in this stat. It is part of the StatCommonReadonly interface.
func (a *Stat) Type(name tc.CacheType) (*StatCacheStats, bool) {
	t, ok := a.Types[name]
	return t, ok
}

// Total returns the aggregated total data in this stat. It is part of the StatCommonReadonly interface.
func (a *Stat) Total() *StatCacheStats {
	return &a.TotalStats
}

// Stats is the JSON-serialisable representation of delivery service Stats. It maps delivery service names to individual stat objects.
type Stats struct {
	DeliveryService map[tc.DeliveryServiceName]*Stat `json:"deliveryService"`
	Time            time.Time                        `json:"-"`
}

// Copy performs a deep copy of this Stats object.
func (s *Stats) Copy() *Stats {
	b := NewStats(len(s.DeliveryService))
	for k, v := range s.DeliveryService {
		b.DeliveryService[k] = v.Copy()
	}
	b.Time = s.Time
	return b
}

// Get returns the stats for the given delivery service, and whether it exists.
func (s Stats) Get(name tc.DeliveryServiceName) (StatReadonly, bool) {
	ds, ok := s.DeliveryService[name]
	return ds, ok
}

// JSON returns an object formatted as expected to be serialized to JSON and served.
func (s Stats) JSON(filter Filter, params url.Values) StatsOld {
	// TODO fix to be the time calculated, not the time requested
	now := s.Time.UnixNano() / int64(time.Millisecond) // Traffic Monitor 1.0 API is 'ms since the epoch'
	jsonObj := &StatsOld{
		CommonAPIData:   srvhttp.GetCommonAPIData(params, time.Now()),
		DeliveryService: map[tc.DeliveryServiceName]map[StatName][]StatOld{},
	}

	for deliveryService, stat := range s.DeliveryService {
		if !filter.UseDeliveryService(deliveryService) {
			continue
		}
		jsonObj.DeliveryService[deliveryService] = map[StatName][]StatOld{}
		jsonObj = addCommonData(jsonObj, &stat.CommonStats, deliveryService, now, filter)
		for cacheGroup, cacheGroupStats := range stat.CacheGroups {
			jsonObj = addStatCacheStats(jsonObj, cacheGroupStats, deliveryService, "location."+string(cacheGroup)+".", now, filter)
		}
		for cacheType, typeStats := range stat.Types {
			jsonObj = addStatCacheStats(jsonObj, typeStats, deliveryService, "type."+cacheType.String()+".", now, filter)
		}
		jsonObj = addStatCacheStats(jsonObj, &stat.TotalStats, deliveryService, "total.", now, filter)
	}
	return *jsonObj
}

// NewStats creates a new Stats object, initializing any pointer members.
// TODO rename to just 'New'?
func NewStats(size int) *Stats {
	return &Stats{DeliveryService: make(map[tc.DeliveryServiceName]*Stat, size)}
}

// LastStats includes the previously recieved stats for DeliveryServices and Caches, the stat itself, when it was received, and the stat value per second.
type LastStats struct {
	DeliveryServices map[tc.DeliveryServiceName]*LastDSStat
	Caches           map[tc.CacheName]*LastStatsData
}

// NewLastStats returns a new LastStats object, initializing internal pointer values.
func NewLastStats(dsLen, cacheLen int) *LastStats {
	// TODO add map size params?
	return &LastStats{DeliveryServices: map[tc.DeliveryServiceName]*LastDSStat{}, Caches: map[tc.CacheName]*LastStatsData{}}
}

// Copy performs a deep copy of this LastStats object.
func (a *LastStats) Copy() *LastStats {
	b := NewLastStats(len(a.DeliveryServices), len(a.Caches))
	for k, v := range a.DeliveryServices {
		b.DeliveryServices[k] = v.Copy()
	}
	for k, v := range a.Caches {
		b.Caches[k] = v
	}
	return b
}

// LastDSStat maps and aggregates the last stats received for the given delivery service to caches, cache groups, types, and total.
// TODO figure a way to associate this type with StatHTTP, with which its members correspond.
type LastDSStat struct {
	Caches      map[tc.CacheName]*LastStatsData
	CacheGroups map[tc.CacheGroupName]*LastStatsData
	Type        map[tc.CacheType]*LastStatsData
	Total       LastStatsData
	Available   bool
}

// Copy performs a deep copy of this LastDSStat object.
func (a LastDSStat) Copy() *LastDSStat {
	b := &LastDSStat{
		CacheGroups: make(map[tc.CacheGroupName]*LastStatsData, len(a.CacheGroups)),
		Type:        make(map[tc.CacheType]*LastStatsData, len(a.Type)),
		Caches:      make(map[tc.CacheName]*LastStatsData, len(a.Caches)),
		Total:       a.Total,
		Available:   a.Available,
	}
	for k, v := range a.CacheGroups {
		b.CacheGroups[k] = v
	}
	for k, v := range a.Type {
		b.Type[k] = v
	}
	for k, v := range a.Caches {
		b.Caches[k] = v
	}
	return b
}

// LastStatsData contains the last stats and per-second calculations for bytes and status codes received from a cache.
// TODO sync.Pool?
type LastStatsData struct {
	Bytes     LastStatData
	Status2xx LastStatData
	Status3xx LastStatData
	Status4xx LastStatData
	Status5xx LastStatData
}

// Sum returns the Sum() of each member data with the given LastStatsData corresponding members
func (a *LastStatsData) Sum(b *LastStatsData) {
	a.Bytes.PerSec += b.Bytes.PerSec
	a.Bytes.Stat += b.Bytes.Stat
	a.Status2xx.PerSec += b.Status2xx.PerSec
	a.Status2xx.Stat += b.Status2xx.Stat
	a.Status3xx.PerSec += b.Status3xx.PerSec
	a.Status3xx.Stat += b.Status3xx.Stat
	a.Status4xx.PerSec += b.Status4xx.PerSec
	a.Status4xx.Stat += b.Status4xx.Stat
	a.Status5xx.PerSec += b.Status5xx.PerSec
	a.Status5xx.Stat += b.Status5xx.Stat
}

// LastStatData contains the value, time it was received, and per-second calculation since the previous stat, for a stat from a cache.
type LastStatData struct {
	PerSec float64
	Stat   int64
	Time   time.Time
}

func addCommonData(s *StatsOld, c *StatCommon, deliveryService tc.DeliveryServiceName, t int64, filter Filter) *StatsOld {
	add := func(name string, val interface{}) {
		if filter.UseStat(name) {
			s.DeliveryService[deliveryService][StatName(name)] = []StatOld{StatOld{Time: t, Value: val}}
		}
	}
	add("caches-configured", fmt.Sprintf("%d", c.CachesConfiguredNum.Value))
	add("caches-reporting", fmt.Sprintf("%d", len(c.CachesReporting)))
	add("error-string", c.ErrorStr.Value)
	add("status", c.StatusStr.Value)
	add("isHealthy", fmt.Sprintf("%t", c.IsHealthy.Value))
	add("isAvailable", fmt.Sprintf("%t", c.IsAvailable.Value))
	add("caches-available", fmt.Sprintf("%d", c.CachesAvailableNum.Value))
	add("disabledLocations", c.CachesDisabled)
	return s
}

func addStatCacheStats(s *StatsOld, c *StatCacheStats, deliveryService tc.DeliveryServiceName, prefix string, t int64, filter Filter) *StatsOld {
	add := func(name, val string) {
		if filter.UseStat(name) {
			// This is for compatibility with the Traffic Monitor 1.0 API.
			// TODO abstract this? Or deprecate and remove it?
			if name == "isAvailable" || name == "error-string" {
				s.DeliveryService[deliveryService][StatName("location."+prefix+name)] = []StatOld{StatOld{Time: t, Value: val}}
			} else {
				s.DeliveryService[deliveryService][StatName(prefix+name)] = []StatOld{StatOld{Time: t, Value: val}}
			}
		}
	}
	add("out_bytes", strconv.Itoa(int(c.OutBytes.Value)))
	add("isAvailable", fmt.Sprintf("%t", c.IsAvailable.Value))
	add("status_5xx", strconv.Itoa(int(c.Status5xx.Value)))
	add("status_4xx", strconv.Itoa(int(c.Status4xx.Value)))
	add("status_3xx", strconv.Itoa(int(c.Status3xx.Value)))
	add("status_2xx", strconv.Itoa(int(c.Status2xx.Value)))
	add("in_bytes", strconv.Itoa(int(c.InBytes.Value)))
	add("kbps", strconv.Itoa(int(c.Kbps.Value)))
	add("tps_5xx", fmt.Sprintf("%f", c.Tps5xx.Value))
	add("tps_4xx", fmt.Sprintf("%f", c.Tps4xx.Value))
	add("tps_3xx", fmt.Sprintf("%f", c.Tps3xx.Value))
	add("tps_2xx", fmt.Sprintf("%f", c.Tps2xx.Value))
	add("error-string", c.ErrorString.Value)
	add("tps_total", fmt.Sprintf("%f", c.TpsTotal.Value))
	return s
}
