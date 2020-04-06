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

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_monitor/srvhttp"
	"github.com/apache/trafficcontrol/traffic_monitor/todata"
)

// Handler is a cache handler, which fulfills the common/handler `Handler` interface.
type Handler struct {
	resultChan chan Result
	ToData     *todata.TODataThreadsafe
}

func (h Handler) ResultChan() <-chan Result {
	return h.resultChan
}

// NewHandler returns a new cache handler. Note this handler does NOT precomputes stat data before calling ResultChan, and Result.Precomputed will be nil
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

// Result is the data result returned by a cache.
// type Result struct {
// 	ID              string// 	Error           error
// 	Astats          Astats
// 	Time            time.Time
// 	RequestTime     time.Duration
// 	Vitals          Vitals
// 	PollID          uint64
// 	UsingIPv4       bool
// 	PollFinished    chan<- uint64
// 	PrecomputedData PrecomputedData
// 	Available       bool
// }

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
	ID            string
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

// Vitals is the vitals data returned from a cache.
type Vitals struct {
	// LoadAvg is the one-minute "loadavg" of the cache server.
	LoadAvg    float64
	BytesOut   uint64
	BytesIn    uint64
	KbpsOut    int64
	MaxKbpsOut int64
}

// Stat is a generic stat, including the untyped value and the time the stat was taken.
type Stat struct {
	Time  int64       `json:"time"`
	Value interface{} `json:"value"`
}

// Stats is designed for returning via the API. It contains result history for each cache, as well as common API data.
type Stats struct {
	srvhttp.CommonAPIData
	Caches map[string]map[string][]ResultStatVal `json:"caches"`
}

// Filter filters whether stats and caches should be returned from a data set.
type Filter interface {
	UseStat(name string) bool
	UseCache(name string) bool
	WithinStatHistoryMax(int) bool
}

const nsPerMs = 1000000

type StatComputeFunc func(resultInfo ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{}

// ComputedStats returns a map of cache stats which are computed by Traffic Monitor (rather than returned literally from ATS), mapped to the func to compute them.
func ComputedStats() map[string]StatComputeFunc {
	return map[string]StatComputeFunc{
		"availableBandwidthInKbps": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return info.Vitals.MaxKbpsOut - info.Vitals.KbpsOut
		},

		"availableBandwidthInMbps": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return (info.Vitals.MaxKbpsOut - info.Vitals.KbpsOut) / 1000
		},
		"bandwidth": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return info.Vitals.KbpsOut
		},
		"error-string": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			if info.Error != nil {
				return info.Error.Error()
			}
			return "false"
		},
		"isAvailable": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return combinedState // if the cache is missing, default to false
		},
		"isHealthy": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			if tc.CacheStatusFromString(serverInfo.ServerStatus) == tc.CacheStatusAdminDown {
				return true
			}
			return combinedState.IsAvailable
		},
		"kbps": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return info.Vitals.KbpsOut
		},
		"gbps": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return float64(info.Vitals.KbpsOut) / 1000000.0
		},
		"loadavg": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return info.Vitals.LoadAvg
		},
		"maxKbps": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return info.Vitals.MaxKbpsOut
		},
		"queryTime": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return info.RequestTime.Nanoseconds() / nsPerMs
		},
		"stateUrl": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return serverProfile.Parameters.HealthPollingURL
		},
		"status": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return serverInfo.ServerStatus
		},
		"system.astatsLoad": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			// return info.System.AstatsLoad
			return 0
		},
		"system.configReloadRequests": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			// return info.System.ConfigLoadRequest
			return 0
		},
		"system.configReloads": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			// return info.System.ConfigReloads
			return 0
		},
		"system.inf.name": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			// return info.System.InfName
			return 0
		},
		"system.inf.speed": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			// return info.System.InfSpeed
			return 0
		},
		"system.lastReload": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			// return info.System.LastReload
			return 0
		},
		"system.lastReloadRequest": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			// return info.System.LastReloadRequest
			return 0
		},
		"system.notAvailable": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			// return info.System.NotAvailable
			return 0
		},
		"system.proc.loadavg": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			// return info.System.ProcLoadavg
			return 0
		},
		"system.proc.net.dev": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			// return info.System.ProcNetDev
			return 0
		},
	}
}

// Handle handles results fetched from a cache, parsing the raw Reader data and passing it along to a chan for further processing.
func (handler Handler) Handle(id string, rdr io.Reader, format string, reqTime time.Duration, reqEnd time.Time, reqErr error, pollID uint64, usingIPv4 bool, pollFinished chan<- uint64) {
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
		log.Warnf("%v handler given error '%v'\n", id, reqErr) // error here, in case the thing that called Handle didn't error
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

	stats, miscStats, err := decoder.Parse(string(result.ID), rdr)
	if err != nil {
		log.Warnf("%s decode error '%v'", id, err)
		result.Error = err
		handler.resultChan <- result
		return
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
