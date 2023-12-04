package cache

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
	"io"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"
)

// Handler is a cache handler, which fulfills the common/handler `Handler` interface.
type Handler struct {
	resultChan chan Result
	ToData     *todata.TODataThreadsafe
}

func (h Handler) ResultChan() <-chan Result {
	return h.resultChan
}

// NewHandler returns a new cache handler. Note this handler does NOT precompute stat data before calling ResultChan, and Result.Precomputed will be nil.
func NewHandler() Handler {
	return Handler{resultChan: make(chan Result)}
}

// NewPrecomputeHandler constructs a new cache Handler, which precomputes stat data and populates result.Precomputed before passing to ResultChan.
func NewPrecomputeHandler(toData todata.TODataThreadsafe) Handler {
	return Handler{resultChan: make(chan Result), ToData: &toData}
}

// Precompute returns whether this handler precomputes data before passing the result to the ResultChan
func (handler Handler) Precompute() bool {
	return handler.ToData != nil
}

// PrecomputedData represents data parsed and pre-computed from the Result.
type PrecomputedData struct {
	DeliveryServiceStats map[string]*DSStat
	// This is the total bytes transmitted by all interfaces on the Cache
	// Server.
	OutBytes uint64
	// MaxKbps is the maximum bandwidth of all interfaces on the Cache Server,
	// each one calculated as the speed of the interface in Kbps.
	MaxKbps   int64
	Errors    []error
	Reporting bool
	Time      time.Time
}

// Result is a result of polling a cache server for statistics.
type Result struct {
	// Available indicates whether or not the cache server should be considered
	// "available" based on its status as configured in Traffic Ops, the cache
	// server's own reported availability (if applicable), and the polled
	// vitals and statistics as compared to threshold values.
	Available bool
	// Error holds what error - if any - caused the statistic polling to fail.
	Error error
	// ID is the fully qualified domain name of the cache server being polled.
	// (This is assumed to be unique even though that isn't necessarily true)
	ID string
	// Miscellaneous contains the stats that were not directly gathered into
	// Statistics, but were still found in the stats polling payload. Their
	// contents are NOT guaranteed in ANY way.
	Miscellaneous map[string]interface{}
	// PollFinished is a channel to which data should be sent to indicate that
	// polling has been completed and a Result has been produced.
	PollFinished chan<- uint64
	// PollID is a unique identifier for the specific polling instance that
	// produced this Result.
	PollID          uint64
	PrecomputedData PrecomputedData
	// RequestTime holds the elapsed duration between making a statistics
	// polling request and either receiving a result or giving up.
	RequestTime time.Duration
	// Statistics holds the parsed statistic data returned by the cache server.
	Statistics Statistics
	// Time is the time at which the result has been obtained.
	Time time.Time
	// UsingIPv4 indicates whether IPv4 can/should be/was used by the polling
	// instance that produced this Result. If ``false'', it may be assumed that
	// IPv6 was used instead.
	UsingIPv4 bool
	// Vitals holds the parsed health information returned by the cache server.
	Vitals Vitals
	// InterfaceVitals holds the parsed health information returned by the cache server per interface.
	InterfaceVitals map[string]Vitals
}

// HasStat returns whether the given stat is in the Result.
func (result *Result) HasStat(stat string) bool {
	computedStats := ComputedStats()
	if _, ok := computedStats[stat]; ok {
		return true // health poll has all computed stats
	}
	if _, ok := result.Miscellaneous[stat]; ok {
		return true
	}
	return false
}

// InterfacesNames returns the names of all network interfaces used by the cache
// server that was monitored to obtain the result.
func (result *Result) InterfacesNames() []string {
	interfaceNames := make([]string, 0, len(result.Statistics.Interfaces))
	for name, _ := range result.Statistics.Interfaces {
		interfaceNames = append(interfaceNames, name)
	}
	return interfaceNames
}

// Interfaces returns the interfaces assigned to this result.
func (result *Result) Interfaces() map[string]Interface {
	return result.Statistics.Interfaces
}

// Vitals is the vitals data returned from a cache.
type Vitals struct {
	// LoadAvg is the one-minute "loadavg" of the cache server.
	LoadAvg    float64
	BytesOut   uint64
	BytesIn    uint64
	KbpsOut    int64
	MaxKbpsOut int64
}

// Stat is a generic stat, including the untyped value and the time the stat was
// taken.
type Stat struct {
	Time  int64       `json:"time"`
	Value interface{} `json:"value"`
}

// Filter filters whether stats and caches should be returned from a data set.
type Filter interface {
	UseCache(tc.CacheName) bool
	UseInterfaceStat(string) bool
	UseStat(string) bool
	WithinStatHistoryMax(uint64) bool
}

const nsPerMs = 1000000

// StatComputeFunc functions calculate a specific statistic given a set of
// polling results, server and profile information, whether or not the server
// is available, and the name of the specific network interface for which stats
// will be computed.
type StatComputeFunc func(ResultInfo, tc.TrafficServer, tc.TMProfile, tc.IsAvailable) interface{}

// ComputedStats returns a map of cache stats which are computed by Traffic
// Monitor (rather than returned literally from ATS), mapped to the function to
// compute them.
func ComputedStats() map[string]StatComputeFunc {
	return map[string]StatComputeFunc{
		"availableBandwidthInKbps": func(info ResultInfo, _ tc.TrafficServer, _ tc.TMProfile, _ tc.IsAvailable) interface{} {
			return info.Vitals.MaxKbpsOut - info.Vitals.KbpsOut
		},
		"availableBandwidthInMbps": func(info ResultInfo, _ tc.TrafficServer, _ tc.TMProfile, _ tc.IsAvailable) interface{} {
			return (info.Vitals.MaxKbpsOut - info.Vitals.KbpsOut) / 1000.0
		},
		tc.StatNameBandwidth: func(info ResultInfo, _ tc.TrafficServer, _ tc.TMProfile, _ tc.IsAvailable) interface{} {
			return info.Vitals.KbpsOut
		},
		tc.StatNameKBPS: func(info ResultInfo, _ tc.TrafficServer, _ tc.TMProfile, _ tc.IsAvailable) interface{} {
			return info.Vitals.KbpsOut
		},
		"gbps": func(info ResultInfo, _ tc.TrafficServer, _ tc.TMProfile, _ tc.IsAvailable) interface{} {
			return float64(info.Vitals.KbpsOut) / 1000000.0
		},
		tc.StatNameMaxKBPS: func(info ResultInfo, _ tc.TrafficServer, _ tc.TMProfile, _ tc.IsAvailable) interface{} {
			return info.Vitals.MaxKbpsOut
		},
		"loadavg": func(info ResultInfo, _ tc.TrafficServer, _ tc.TMProfile, _ tc.IsAvailable) interface{} {
			return info.Vitals.LoadAvg
		},
		"queryTime": func(info ResultInfo, _ tc.TrafficServer, _ tc.TMProfile, _ tc.IsAvailable) interface{} {
			return info.RequestTime.Nanoseconds() / nsPerMs
		},
		"stateUrl": func(_ ResultInfo, _ tc.TrafficServer, serverProfile tc.TMProfile, _ tc.IsAvailable) interface{} {
			return serverProfile.Parameters.HealthPollingURL
		},
		"status": func(_ ResultInfo, serverInfo tc.TrafficServer, _ tc.TMProfile, _ tc.IsAvailable) interface{} {
			return serverInfo.ServerStatus
		},
		"error-string": func(info ResultInfo, _ tc.TrafficServer, _ tc.TMProfile, _ tc.IsAvailable) interface{} {
			if info.Error != nil {
				return info.Error.Error()
			}
			return "false"
		},
		"isAvailable": func(_ ResultInfo, _ tc.TrafficServer, _ tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return combinedState // if the cache is missing, default to false
		},
		"isHealthy": func(_ ResultInfo, serverInfo tc.TrafficServer, _ tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			if tc.CacheStatusFromString(serverInfo.ServerStatus) == tc.CacheStatusAdminDown {
				return true
			}
			return combinedState.IsAvailable
		},

		// These are back-up values for when a statistics format doesn't
		// support reporting these stats - which would make sense because five
		// of them are pre-parsed in Statistics structures already, and I'm not
		// sure what the rest of them are even for. None of these are
		// documented anywhere. The values in comments are the ones that astats
		// parsers will give back (because it won't get this far).
		"system.astatsLoad": func(ResultInfo, tc.TrafficServer, tc.TMProfile, tc.IsAvailable) interface{} {
			// return info.System.AstatsLoad
			return float64(0)
		},
		"system.configReloadRequests": func(ResultInfo, tc.TrafficServer, tc.TMProfile, tc.IsAvailable) interface{} {
			// return info.System.ConfigLoadRequest
			return float64(0)
		},
		"system.configReloads": func(ResultInfo, tc.TrafficServer, tc.TMProfile, tc.IsAvailable) interface{} {
			// return info.System.ConfigReloads
			return float64(0)
		},
		"system.inf.name": func(ResultInfo, tc.TrafficServer, tc.TMProfile, tc.IsAvailable) interface{} {
			// return info.System.InfName
			return ""
		},
		"system.inf.speed": func(ResultInfo, tc.TrafficServer, tc.TMProfile, tc.IsAvailable) interface{} {
			// return info.System.InfSpeed
			return float64(0)
		},
		"system.lastReload": func(ResultInfo, tc.TrafficServer, tc.TMProfile, tc.IsAvailable) interface{} {
			// return info.System.LastReload
			return float64(0)
		},
		"system.lastReloadRequest": func(ResultInfo, tc.TrafficServer, tc.TMProfile, tc.IsAvailable) interface{} {
			// return info.System.LastReloadRequest
			return ""
		},
		"system.notAvailable": func(ResultInfo, tc.TrafficServer, tc.TMProfile, tc.IsAvailable) interface{} {
			// return info.System.NotAvailable
			return ""
		},
		"system.proc.loadavg": func(ResultInfo, tc.TrafficServer, tc.TMProfile, tc.IsAvailable) interface{} {
			// return info.System.ProcLoadavg
			return float64(0)
		},
		"system.proc.net.dev": func(ResultInfo, tc.TrafficServer, tc.TMProfile, tc.IsAvailable) interface{} {
			// return info.System.ProcNetDev
			return float64(0)
		},
	}
}

// Handle handles results fetched from a cache, parsing the raw Reader data and passing it along to a chan for further processing.
func (handler Handler) Handle(id string, rdr io.Reader, format string, reqTime time.Duration, reqEnd time.Time, reqErr error, pollID uint64, usingIPv4 bool, pollCtx interface{}, pollFinished chan<- uint64) {
	log.Debugf("poll %v %v (format '%v') handle start\n", pollID, time.Now(), format)
	result := Result{
		ID:           id,
		Time:         reqEnd,
		RequestTime:  reqTime,
		PollID:       pollID,
		UsingIPv4:    usingIPv4,
		PollFinished: pollFinished,
	}

	if reqErr != nil {
		log.Warnf("%s handler given error: %s", id, reqErr.Error()) // error here, in case the thing that called Handle didn't error
		result.Error = reqErr
		handler.resultChan <- result
		return
	}

	decoder, err := GetDecoder(format)
	if err != nil {
		log.Errorln(err.Error())
		result.Error = err
		handler.resultChan <- result
		return
	}

	stats, miscStats, err := decoder.Parse(result.ID, rdr, pollCtx)
	if err != nil {
		log.Warnf("%s decode error '%v'", id, err)
		result.Error = err
		handler.resultChan <- result
		return
	}
	if val, ok := miscStats["plugin.system_stats.timestamp_ms"]; ok {
		valInt, valErr := parseNumericStat(val)
		if valErr != nil {
			log.Errorln("parse error: ", valErr)
			result.Error = valErr
			handler.resultChan <- result
			return
		}
		result.Time = time.UnixMilli(int64(valInt))
	}

	result.Statistics = stats
	result.Miscellaneous = miscStats

	result.Available = true

	if handler.Precompute() {
		result.PrecomputedData = decoder.Precompute(result.ID, handler.ToData.Get(), result.Statistics, result.Miscellaneous)
	}
	result.PrecomputedData.Reporting = true
	result.PrecomputedData.Time = result.Time

	handler.resultChan <- result
}
