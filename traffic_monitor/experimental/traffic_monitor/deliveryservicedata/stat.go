package deliveryservicedata // TODO rename?

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
	"net/url"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/srvhttp"
)

// Filter encapsulates functions to filter a given set of Stats, e.g. from HTTP query parameters.
// TODO combine with cache.Filter?
type Filter interface {
	UseStat(name string) bool
	UseDeliveryService(name enum.DeliveryServiceName) bool
	WithinStatHistoryMax(int) bool
}

// StatName is the name of a stat.
type StatName string

// StatOld is the old JSON representation of a stat, from Traffic Monitor 1.0.
type StatOld struct {
	Time  int64  `json:"time"`
	Value string `json:"value"`
	Span  int    `json:"span,omitempty"`  // TODO set? remove?
	Index int    `json:"index,omitempty"` // TODO set? remove?
}

// StatsOld is the old JSON representation of stats, from Traffic Monitor 1.0. It is designed to be serialized and returns from an API, and includes stat history for each delivery service, as well as data common to most endpoints.
type StatsOld struct {
	DeliveryService map[enum.DeliveryServiceName]map[StatName][]StatOld `json:"deliveryService"`
	srvhttp.CommonAPIData
}

// StatsReadonly is a read-only interface for delivery service Stats, designed to be passed to multiple goroutine readers.
type StatsReadonly interface {
	Get(enum.DeliveryServiceName) (StatReadonly, bool)
	JSON(Filter, url.Values) StatsOld
}

// StatReadonly is a read-only interface for a delivery service Stat, designed to be passed to multiple goroutine readers.
type StatReadonly interface {
	Copy() Stat
	Common() StatCommonReadonly
	CacheGroup(name enum.CacheGroupName) (StatCacheStats, bool)
	Type(name enum.CacheType) (StatCacheStats, bool)
	Total() StatCacheStats
}

// StatCommonReadonly is a read-only interface for a delivery service's common Stat data, designed to be passed to multiple goroutine readers.
type StatCommonReadonly interface {
	Copy() StatCommon
	CachesConfigured() StatInt
	CachesReportingNames() []enum.CacheName
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
	CachesConfiguredNum StatInt                 `json:"caches_configured"`
	CachesReporting     map[enum.CacheName]bool `json:"caches_reporting"`
	ErrorStr            StatString              `json:"error_string"`
	StatusStr           StatString              `json:"status"`
	IsHealthy           StatBool                `json:"is_healthy"`
	IsAvailable         StatBool                `json:"is_available"`
	CachesAvailableNum  StatInt                 `json:"caches_available"`
}

// Copy returns a deep copy of this StatCommon object.
func (a StatCommon) Copy() StatCommon {
	b := a
	for k, v := range a.CachesReporting {
		b.CachesReporting[k] = v
	}
	return b
}

// CachesConfigured returns the number of caches configured for this delivery service stat. It is part of the StatCommonReadonly interface.
func (a StatCommon) CachesConfigured() StatInt {
	return a.CachesConfiguredNum
}

// CacheReporting returns the number of caches reporting for this delivery service stat. It is part of the StatCommonReadonly interface.
func (a StatCommon) CacheReporting(name enum.CacheName) (bool, bool) {
	c, ok := a.CachesReporting[name]
	return c, ok
}

// CachesReportingNames returns the list of caches reporting for this delivery service stat. It is part of the StatCommonReadonly interface.
func (a StatCommon) CachesReportingNames() []enum.CacheName {
	names := make([]enum.CacheName, 0, len(a.CachesReporting))
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
	CommonStats        StatCommon
	CacheGroups        map[enum.CacheGroupName]StatCacheStats
	Types              map[enum.CacheType]StatCacheStats
	Caches             map[enum.CacheName]StatCacheStats
	CachesTimeReceived map[enum.CacheName]time.Time
	TotalStats         StatCacheStats
}

// ErrNotProcessedStat indicates a stat received is not used by Traffic Monitor, nor returned by any API endpoint. Receiving this error indicates the stat has been discarded.
var ErrNotProcessedStat = errors.New("This stat is not used.")

// NewStat returns a new delivery service Stat, initializing pointer members.
func NewStat() *Stat {
	return &Stat{
		CacheGroups:        map[enum.CacheGroupName]StatCacheStats{},
		Types:              map[enum.CacheType]StatCacheStats{},
		CommonStats:        StatCommon{CachesReporting: map[enum.CacheName]bool{}},
		Caches:             map[enum.CacheName]StatCacheStats{},
		CachesTimeReceived: map[enum.CacheName]time.Time{},
	}
}

// Copy performs a deep copy of this Stat. It does not modify, and is thus safe for multiple goroutines.
func (a Stat) Copy() Stat {
	b := Stat{
		CommonStats:        a.CommonStats.Copy(),
		TotalStats:         a.TotalStats,
		CacheGroups:        map[enum.CacheGroupName]StatCacheStats{},
		Types:              map[enum.CacheType]StatCacheStats{},
		Caches:             map[enum.CacheName]StatCacheStats{},
		CachesTimeReceived: map[enum.CacheName]time.Time{},
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
	for k, v := range a.CachesTimeReceived {
		b.CachesTimeReceived[k] = v
	}
	return b
}

// Common returns the common stat data for this stat. It is part of the StatCommonReadonly interface.
func (a Stat) Common() StatCommonReadonly {
	return a.CommonStats
}

// CacheGroup returns the data for the given cachegroup in this stat. It is part of the StatCommonReadonly interface.
func (a Stat) CacheGroup(name enum.CacheGroupName) (StatCacheStats, bool) {
	c, ok := a.CacheGroups[name]
	return c, ok
}

// Type returns the aggregated data for the given cache type in this stat. It is part of the StatCommonReadonly interface.
func (a Stat) Type(name enum.CacheType) (StatCacheStats, bool) {
	t, ok := a.Types[name]
	return t, ok
}

// Total returns the aggregated total data in this stat. It is part of the StatCommonReadonly interface.
func (a Stat) Total() StatCacheStats {
	return a.TotalStats
}
