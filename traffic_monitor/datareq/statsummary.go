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

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/traffic_monitor/srvhttp"
	"github.com/apache/trafficcontrol/traffic_monitor/threadsafe"
	"github.com/apache/trafficcontrol/traffic_monitor/todata"

	"github.com/json-iterator/go"
)

type StatSummary struct {
	Caches map[string]map[string]StatSummaryStat `json:"caches"`
	srvhttp.CommonAPIData
}

type StatSummaryStat struct {
	DataPointCount int64   `json:"dpCount"`
	Start          float64 `json:"start"`
	End            float64 `json:"end"`
	High           float64 `json:"high"`
	Low            float64 `json:"low"`
	Average        float64 `json:"average"`
	StartTime      int64   `json:"startTime"`
	EndTime        int64   `json:"endTime"`
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
	statPrefix := "ats."
	ss := StatSummary{
		Caches:        map[string]map[string]StatSummaryStat{},
		CommonAPIData: srvhttp.GetCommonAPIData(params, time.Now()),
	}

	statResultHistory.Range(func(cacheName string, stats threadsafe.ResultStatValHistory) bool {
		if !filter.UseCache(cacheName) {
			return true
		}
		ssStats := map[string]StatSummaryStat{}
		stats.Range(func(statName string, statHistory []cache.ResultStatVal) bool {
			if !filter.UseStat(statName) {
				return true
			}
			if len(statHistory) == 0 {
				return true
			}
			ssStat := StatSummaryStat{}
			msPerNs := int64(1000000)
			ssStat.StartTime = statHistory[len(statHistory)-1].Time.UnixNano() / msPerNs
			ssStat.EndTime = statHistory[0].Time.UnixNano() / msPerNs
			oldestVal, isOldestValNumeric := util.ToNumeric(statHistory[len(statHistory)-1].Val)
			newestVal, isNewestValNumeric := util.ToNumeric(statHistory[0].Val)
			if !isOldestValNumeric || !isNewestValNumeric {
				return true // skip non-numeric stats
			}
			ssStat.Start = oldestVal
			ssStat.End = newestVal
			ssStat.High = newestVal
			ssStat.Low = newestVal
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
		ss.Caches[cacheName] = ssStats
		return true
	})
	return ss
}
