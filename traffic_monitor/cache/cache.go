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
	"fmt"
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
	DeliveryServiceStats map[tc.DeliveryServiceName]*AStat
	OutBytes             int64
	MaxKbps              int64
	Errors               []error
	Reporting            bool
	Time                 time.Time
}

// Result is the data result returned by a cache.
type Result struct {
	ID              tc.CacheName
	Error           error
	Astats          Astats
	Time            time.Time
	RequestTime     time.Duration
	Vitals          Vitals
	PollID          uint64
	PollFinished    chan<- uint64
	PrecomputedData PrecomputedData
	Available       bool
}

// HasStat returns whether the given stat is in the Result.
func (result *Result) HasStat(stat string) bool {
	computedStats := ComputedStats()
	if _, ok := computedStats[stat]; ok {
		return true // health poll has all computed stats
	}
	if _, ok := result.Astats.Ats[stat]; ok {
		return true
	}
	return false
}

// Vitals is the vitals data returned from a cache.
type Vitals struct {
	LoadAvg    float64
	BytesOut   int64
	BytesIn    int64
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
	Caches map[tc.CacheName]map[string][]ResultStatVal `json:"caches"`
}

// Filter filters whether stats and caches should be returned from a data set.
type Filter interface {
	UseStat(name string) bool
	UseCache(name tc.CacheName) bool
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
			return combinedState.IsAvailable // if the cache is missing, default to false
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
			return info.System.AstatsLoad
		},
		"system.configReloadRequests": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return info.System.ConfigLoadRequest
		},
		"system.configReloads": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return info.System.ConfigReloads
		},
		"system.inf.name": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return info.System.InfName
		},
		"system.inf.speed": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return info.System.InfSpeed
		},
		"system.lastReload": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return info.System.LastReload
		},
		"system.lastReloadRequest": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return info.System.LastReloadRequest
		},
		"system.notAvailable": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return info.System.NotAvailable
		},
		"system.proc.loadavg": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return info.System.ProcLoadavg
		},
		"system.proc.net.dev": func(info ResultInfo, serverInfo tc.TrafficServer, serverProfile tc.TMProfile, combinedState tc.IsAvailable) interface{} {
			return info.System.ProcNetDev
		},
	}
}

// Handle handles results fetched from a cache, parsing the raw Reader data and passing it along to a chan for further processing.
func (handler Handler) Handle(id string, rdr io.Reader, format string, reqTime time.Duration, reqEnd time.Time, reqErr error, pollID uint64, pollFinished chan<- uint64) {
	log.Debugf("poll %v %v (format '%v') handle start\n", pollID, time.Now(), format)
	result := Result{
		ID:           tc.CacheName(id),
		Time:         reqEnd,
		RequestTime:  reqTime,
		PollID:       pollID,
		PollFinished: pollFinished,
	}

	if reqErr != nil {
		log.Warnf("%v handler given error '%v'\n", id, reqErr) // error here, in case the thing that called Handle didn't error
		result.Error = reqErr
		handler.resultChan <- result
		return
	}

	statDecoder, ok := StatsTypeDecoders[format]
	if !ok {
		log.Errorf("Handler cache '%s' stat type '%s' not found! Returning handle error for this cache poll.\n", id, format)
		result.Error = fmt.Errorf("handler stat type %s missing", format)
		handler.resultChan <- result
		return
	}

	decodeErr := error(nil)
	if decodeErr, result.Astats.Ats, result.Astats.System = statDecoder.Parse(result.ID, rdr); decodeErr != nil {
		log.Warnf("%s decode error '%v'\n", id, decodeErr)
		result.Error = decodeErr
		handler.resultChan <- result
		return
	}

	if result.Astats.System.ProcNetDev == "" {
		log.Warnf("Handler cache %s procnetdev empty\n", id)
	}
	if result.Astats.System.InfSpeed == 0 {
		log.Warnf("Handler cache %s inf.speed empty\n", id)
	}

	result.Available = true

	if handler.Precompute() {
		result.PrecomputedData = statDecoder.Precompute(result.ID, handler.ToData.Get(), result.Astats.Ats, result.Astats.System)
	}
	result.PrecomputedData.Reporting = true
	result.PrecomputedData.Time = result.Time

	handler.resultChan <- result
}
