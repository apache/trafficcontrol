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

package datareq

import (
	"net/http"
	"net/url"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/srvhttp"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/threadsafe"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"
	jsoniter "github.com/json-iterator/go"
)

// msPerNs is the number of milliseconds in a nanosecond.
const msPerNs = 1000000

// for whatever reason, all stats are prefixed with this
const statPrefix = "ats."

// StatSummaryStat represents a summary of a stat's values and changes over a
// period of time.
type StatSummaryStat struct {
	// Average is the stat's value's arithmetic mean within the time period.
	Average float64 `json:"average"`
	// DataPointCount is the number of measurements of the stat that have been
	// cataloged in the stat's history.
	DataPointCount int64 `json:"dpCount"`
	// End is the final value of the stat at the time of the most recent
	// measurement.
	End float64 `json:"end"`
	// EndTime is the time of the most recent measurement, and defines the end
	// point of the time period.
	EndTime int64 `json:"endTime"`
	// High is the stat's maximum value within the time period.
	High float64 `json:"high"`
	// Low is the stat's minimum value within the time period.
	Low float64 `json:"low"`
	// Start is the initial value of the stat at the time of the first
	// measurement.
	Start float64 `json:"start"`
	// StartTime is the time of the first measurement, and defines the beginning
	// of the time period.
	StartTime int64 `json:"startTime"`
}

// CacheStatSummary is a summary of a cache server's measured statistics and the
// measured statistics of its network interfaces.
type CacheStatSummary struct {
	// InterfaceStats is a map of network interface names to a map of statistic
	// names to their summaries.
	InterfaceStats map[string]map[string]StatSummaryStat `json:"interfaceStats"`
	// Stats is a map of statistic names to summaries of those statistics.
	Stats map[string]StatSummaryStat `json:"stats"`
}

type StatSummary struct {
	Caches map[string]CacheStatSummary `json:"caches"`
	tc.CommonAPIData
}

func srvStatSummary(params url.Values, errorCount threadsafe.Uint, path string, toData todata.TODataThreadsafe, statResultHistory threadsafe.ResultStatHistory) ([]byte, int) {
	filter, err := NewCacheStatFilter(path, params, toData.Get().ServerTypes)
	if err != nil {
		HandleErr(errorCount, path, err)
		return []byte(err.Error()), http.StatusBadRequest
	}

	json := jsoniter.ConfigFastest
	bytes, err := json.Marshal(createStatSummary(statResultHistory, filter, params))
	return WrapErrCode(errorCount, path, bytes, err)
}

func createStatSummary(statResultHistory threadsafe.ResultStatHistory, filter cache.Filter, params url.Values) StatSummary {
	ss := StatSummary{
		Caches:        map[string]CacheStatSummary{},
		CommonAPIData: srvhttp.GetCommonAPIData(params, time.Now()),
	}

	statResultHistory.Range(func(cacheName string, stats threadsafe.CacheStatHistory) bool {
		if !filter.UseCache(tc.CacheName(cacheName)) {
			return true
		}

		var cacheStats CacheStatSummary

		ssStats := map[string]StatSummaryStat{}
		stats.Stats.Range(func(statName string, statHistory []tc.ResultStatVal) bool {
			if !filter.UseStat(statName) || len(statHistory) == 0 {
				return true
			}

			ssStat := StatSummaryStat{
				EndTime:        statHistory[0].Time.UnixNano() / msPerNs,
				DataPointCount: 0,
				StartTime:      statHistory[len(statHistory)-1].Time.UnixNano() / msPerNs,
			}

			oldestVal, isOldestValNumeric := util.ToNumeric(statHistory[len(statHistory)-1].Val)
			newestVal, isNewestValNumeric := util.ToNumeric(statHistory[0].Val)
			if !isOldestValNumeric || !isNewestValNumeric {
				return true // skip non-numeric stats
			}

			ssStat.End = newestVal
			ssStat.High = newestVal
			ssStat.Low = newestVal
			ssStat.Start = oldestVal

			for _, val := range statHistory {
				fVal, ok := util.ToNumeric(val.Val)
				if !ok {
					log.Warnf("threshold stat %v value %v is not a number, cannot use.", statName, val.Val)
					return true
				}

				for i := uint64(0); i < val.Span; i++ {
					ssStat.DataPointCount++
					ssStat.Average -= ssStat.Average / float64(ssStat.DataPointCount)
					ssStat.Average += fVal / float64(ssStat.DataPointCount)
				}
				if fVal < ssStat.Low {
					ssStat.Low = fVal
				}
				if fVal > ssStat.High {
					ssStat.High = fVal
				}
			}
			ssStats[statPrefix+statName] = ssStat
			return true
		})

		cacheStats.Stats = ssStats

		infStats := make(map[string]map[string]StatSummaryStat, len(stats.Interfaces))
		for infName, infStatHistory := range stats.Interfaces {
			if _, ok := infStats[infName]; ok {
				log.Warnf("Somehow found duplicate interface '%s' in stat history for cache server '%s'", infName, cacheName)
				continue
			}
			infStatMap := map[string]StatSummaryStat{}

			infStatHistory.Range(func(statName string, statHistory []tc.ResultStatVal) bool {
				if !filter.UseInterfaceStat(statName) || len(statHistory) == 0 {
					return true
				}

				ssStat := StatSummaryStat{
					EndTime:        statHistory[0].Time.UnixNano() / msPerNs,
					DataPointCount: 0,
					StartTime:      statHistory[len(statHistory)-1].Time.UnixNano() / msPerNs,
				}

				oldestVal, isOldestValNumeric := util.ToNumeric(statHistory[len(statHistory)-1].Val)
				newestVal, isNewestValNumeric := util.ToNumeric(statHistory[0].Val)
				if !isOldestValNumeric || !isNewestValNumeric {
					return true // skip non-numeric stats
				}

				ssStat.End = newestVal
				ssStat.High = newestVal
				ssStat.Low = newestVal
				ssStat.Start = oldestVal

				for _, val := range statHistory {
					fVal, ok := util.ToNumeric(val.Val)
					if !ok {
						log.Warnf("threshold stat %v value %v is not a number, cannot use.", statName, val.Val)
						return true
					}

					for i := uint64(0); i < val.Span; i++ {
						ssStat.DataPointCount++
						ssStat.Average -= ssStat.Average / float64(ssStat.DataPointCount)
						ssStat.Average += fVal / float64(ssStat.DataPointCount)
					}
					if fVal < ssStat.Low {
						ssStat.Low = fVal
					}
					if fVal > ssStat.High {
						ssStat.High = fVal
					}
				}
				infStatMap[statPrefix+statName] = ssStat
				return true
			})

			infStats[infName] = infStatMap
		}

		cacheStats.InterfaceStats = infStats
		ss.Caches[cacheName] = cacheStats
		return true
	})
	return ss
}
