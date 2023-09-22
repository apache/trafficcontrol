package ds

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
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
	"regexp"
	"testing"

	tc_log "github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/dsdata"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/health"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/threadsafe"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/test"
)

func checkLogOutput(t *testing.T, buffer *bytes.Buffer, toData todata.TOData, caches []tc.CacheName) {
	output := buffer.String()
	for _, cacheName := range caches {
		match, err := regexp.MatchString(string(cacheName)+".*not found in DeliveryServices", output)
		if err != nil {
			t.Fatalf("cannot match cache name %s against output", cacheName)
		}
		if toData.ServerTypes[cacheName] != tc.CacheTypeMid && !match {
			t.Fatalf("expected to find info-level message about cache %s with cache type %s not being found in DeliveryServices, no message found", cacheName, toData.ServerTypes[cacheName])
		}
		if toData.ServerTypes[cacheName] == tc.CacheTypeMid && match {
			t.Fatalf("expected not to find info-level message about cache %s with cache type %s not being found in DeliveryServices, but message was found", cacheName, toData.ServerTypes[cacheName])
		}
	}
}

func TestCreateStats(t *testing.T) {
	toData := getMockTOData()
	combinedCRStates := peer.NewCRStatesThreadsafe()
	lastStatsThs := threadsafe.NewLastStats()
	maxEvents := uint64(4)
	events := health.NewThreadsafeEvents(maxEvents)
	localCRStates := peer.NewCRStatesThreadsafe()

	dses := []tc.DeliveryServiceName{}
	for ds, _ := range toData.DeliveryServiceServers {
		dses = append(dses, ds)
	}

	caches := []tc.CacheName{}
	for cache, _ := range toData.ServerDeliveryServices {
		caches = append(caches, cache)
	}

	for _, cache := range caches {
		combinedCRStates.AddCache(cache, tc.IsAvailable{IsAvailable: true})
		localCRStates.AddCache(cache, tc.IsAvailable{IsAvailable: true})
	}

	precomputeds := randCachesPrecomputedData(caches, toData)

	monitorConfig := getMockMonitorConfig(dses)

	lastStatsVal := lastStatsThs.Get()
	lastStatsCopy := lastStatsVal.Copy()

	dsStats := CreateStats(precomputeds, toData, combinedCRStates.Get(), lastStatsCopy, monitorConfig, events, localCRStates)

	serverDeliveryServices := toData.ServerDeliveryServices
	toData.ServerDeliveryServices = map[tc.CacheName][]tc.DeliveryServiceName{} // temporarily unassign servers to generate warnings about caches not assigned to delivery services
	buffer := bytes.NewBuffer(make([]byte, 0, 10000))
	tc_log.Info = log.New(buffer, "TestAddAvailabilityDataNotFoundInDeliveryService", log.Lshortfile)
	_ = CreateStats(precomputeds, toData, combinedCRStates.Get(), lastStatsCopy, monitorConfig, events, localCRStates)
	checkLogOutput(t, buffer, toData, caches)
	toData.ServerDeliveryServices = serverDeliveryServices

	lastStatsThs.Set(*lastStatsCopy)

	cgMap := map[tc.CacheGroupName]struct{}{}
	for _, cg := range toData.ServerCachegroups {
		cgMap[cg] = struct{}{}
	}

	tpMap := map[tc.CacheType]struct{}{}
	for _, tp := range toData.ServerTypes {
		tpMap[tp] = struct{}{}
	}

	caMap := map[tc.CacheName]struct{}{}
	for ca, _ := range toData.ServerDeliveryServices {
		caMap[ca] = struct{}{}
	}

	for dsName, dsStat := range dsStats.DeliveryService {
		for cgName, cgStat := range dsStat.CacheGroups {
			if _, ok := cgMap[cgName]; !ok {
				t.Fatalf("CreateStats cachegroup expected: %+v, actual: %+v", cgMap, cgName)
			}

			var cgExpected cache.DSStat
			for pCache, pData := range precomputeds {
				if toData.ServerCachegroups[pCache] != cgName {
					continue
				}

				if pDataDS, ok := pData.DeliveryServiceStats[string(dsName)]; ok {
					cgExpected.InBytes += pDataDS.InBytes
					cgExpected.OutBytes += pDataDS.OutBytes
					cgExpected.Status2xx += pDataDS.Status2xx
					cgExpected.Status3xx += pDataDS.Status3xx
					cgExpected.Status4xx += pDataDS.Status4xx
					cgExpected.Status5xx += pDataDS.Status5xx
				}
			}

			if errStr := compareAStatToStatCacheStats(&cgExpected, cgStat); errStr != "" {
				t.Fatalf("CreateStats cachegroup " + string(cgName) + ": " + errStr)
			}

		}

		for tpName, tpStat := range dsStat.Types {
			if _, ok := tpMap[tpName]; !ok {
				t.Fatalf("CreateStats type expected: %+v, actual: %+v", tpMap, tpName)
			}

			var tpExpected cache.DSStat
			for pCache, pData := range precomputeds {
				if toData.ServerTypes[pCache] != tpName {
					continue
				}

				if pDataDS, ok := pData.DeliveryServiceStats[string(dsName)]; ok {
					tpExpected.InBytes += pDataDS.InBytes
					tpExpected.OutBytes += pDataDS.OutBytes
					tpExpected.Status2xx += pDataDS.Status2xx
					tpExpected.Status3xx += pDataDS.Status3xx
					tpExpected.Status4xx += pDataDS.Status4xx
					tpExpected.Status5xx += pDataDS.Status5xx
				}
			}

			if errStr := compareAStatToStatCacheStats(&tpExpected, tpStat); errStr != "" {
				t.Fatalf("CreateStats type " + string(tpName) + ": " + errStr)
			}
		}

		for caName, caStat := range dsStat.Caches {
			if _, ok := caMap[caName]; !ok {
				t.Fatalf("CreateStats cache expected: %+v, actual: %+v", caMap, caName)
			}

			var caExpected cache.DSStat
			for pCache, pData := range precomputeds {
				if pCache != caName {
					continue
				}

				if pDataDS, ok := pData.DeliveryServiceStats[string(dsName)]; ok {
					caExpected.InBytes += pDataDS.InBytes
					caExpected.OutBytes += pDataDS.OutBytes
					caExpected.Status2xx += pDataDS.Status2xx
					caExpected.Status3xx += pDataDS.Status3xx
					caExpected.Status4xx += pDataDS.Status4xx
					caExpected.Status5xx += pDataDS.Status5xx
				}
			}

			if errStr := compareAStatToStatCacheStats(&caExpected, caStat); errStr != "" {
				t.Fatalf("CreateStats cache " + string(caName) + ": " + errStr)
			}
		}

		{
			cmStat := dsStat.CommonStats

			if int(cmStat.CachesConfiguredNum.Value) != len(toData.DeliveryServiceServers[dsName]) {
				t.Fatalf("CreateStats CommonStats.CachesConfiguredNum expected: %+v actual: %+v", len(toData.DeliveryServiceServers[dsName]), dsStat.CommonStats.CachesConfiguredNum.Value)
			}

			for caName, reporting := range cmStat.CachesReporting {
				if _, ok := caMap[caName]; !ok {
					t.Fatalf("CreateStats CommonStats.CachesReporting '%+v' not in test caches", caName)
				}
				if !reporting {
					t.Fatalf("CreateStats len(CommonStats.CachesReporting[%+v] expected: true actual: false", caName)
				}
			}

			if cmStat.ErrorStr.Value != "" {
				t.Fatalf("CreateStats CommonStats.ErrorStr expected: '' actual: %+v", cmStat.ErrorStr.Value)
			}

			if cmStat.StatusStr.Value != "" {
				t.Fatalf("CreateStats CommonStats.StatusStr expected: '' actual: '%+v'", cmStat.StatusStr.Value)
			}
		}
	}

	if len(lastStatsCopy.DeliveryServices) != len(toData.DeliveryServiceServers) {
		t.Fatalf("CreateStats len(LastStats.DeliveryServices) expected: %+v actual: %+v", len(toData.DeliveryServiceServers), len(lastStatsCopy.DeliveryServices))
	}

	if len(lastStatsCopy.Caches) != len(toData.ServerDeliveryServices) {
		t.Fatalf("CreateStats len(LastStats.Caches) expected: %+v actual: %+v", len(toData.ServerDeliveryServices), len(lastStatsCopy.Caches))
	}

}

// compareAStatToStatCacheStats compares the two stats, and returns an error string, which is empty of both are equal.
// The fields in StatCacheStats but not AStat are ignored.
func compareAStatToStatCacheStats(expected *cache.DSStat, actual *dsdata.StatCacheStats) string {
	if actual.InBytes.Value != float64(expected.InBytes) {
		return fmt.Sprintf("InBytes expected: \n%+v, actual: \n%+v", expected.InBytes, actual.InBytes.Value)
	}

	if actual.OutBytes.Value != int64(expected.OutBytes) {
		return fmt.Sprintf("OutBytes expected: \n%+v, actual: \n%+v", expected.OutBytes, actual.OutBytes.Value)
	}

	if actual.Status2xx.Value != int64(expected.Status2xx) {
		return fmt.Sprintf("Status2xx expected: \n%+v, actual: \n%+v", expected.Status2xx, actual.Status2xx.Value)
	}

	if actual.Status3xx.Value != int64(expected.Status3xx) {
		return fmt.Sprintf("Status3xx expected: \n%+v, actual: \n%+v", expected.Status3xx, actual.Status3xx.Value)
	}

	if actual.Status4xx.Value != int64(expected.Status4xx) {
		return fmt.Sprintf("Status4xx expected: \n%+v, actual: \n%+v", expected.Status4xx, actual.Status4xx.Value)
	}

	if actual.Status5xx.Value != int64(expected.Status5xx) {
		return fmt.Sprintf("Status5xx expected: \n%+v, actual: \n%+v", expected.Status5xx, actual.Status5xx.Value)
	}

	if actual.ErrorString.Value != "" {
		return fmt.Sprintf("ErrorString expected: empty, actual: %+v", actual.ErrorString.Value)
	}

	return ""
}

func getMockMonitorDSNoThresholds(name tc.DeliveryServiceName) tc.TMDeliveryService {
	return tc.TMDeliveryService{
		XMLID:              string(name),
		TotalTPSThreshold:  math.MaxInt64,
		ServerStatus:       string(tc.CacheStatusReported),
		TotalKbpsThreshold: math.MaxInt64,
	}
}

func getMockMonitorDSLowThresholds(name tc.DeliveryServiceName) tc.TMDeliveryService {
	return tc.TMDeliveryService{
		XMLID:              string(name),
		TotalTPSThreshold:  1,
		ServerStatus:       string(tc.CacheStatusReported),
		TotalKbpsThreshold: 1,
	}
}

func getMockMonitorConfig(dses []tc.DeliveryServiceName) tc.TrafficMonitorConfigMap {
	mc := tc.TrafficMonitorConfigMap{
		TrafficServer:   map[string]tc.TrafficServer{},
		CacheGroup:      map[string]tc.TMCacheGroup{},
		Config:          map[string]interface{}{},
		TrafficMonitor:  map[string]tc.TrafficMonitor{},
		DeliveryService: map[string]tc.TMDeliveryService{},
		Profile:         map[string]tc.TMProfile{},
	}

	tmDSes := map[string]tc.TMDeliveryService{}
	for _, ds := range dses {
		tmDSes[string(ds)] = getMockMonitorDSNoThresholds(ds)
	}
	mc.DeliveryService = tmDSes

	return mc
}

func getMockTOData() todata.TOData {
	numCaches := 100
	numDSes := 100
	numCacheDSes := numDSes / 3
	numCGs := 20

	types := []tc.CacheType{tc.CacheTypeEdge, tc.CacheTypeEdge, tc.CacheTypeEdge, tc.CacheTypeEdge, tc.CacheTypeEdge, tc.CacheTypeMid}

	caches := []tc.CacheName{}
	for i := 0; i < numCaches; i++ {
		caches = append(caches, tc.CacheName(test.RandStr()))
	}

	dses := []tc.DeliveryServiceName{}
	for i := 0; i < numDSes; i++ {
		dses = append(dses, tc.DeliveryServiceName(test.RandStr()))
	}

	cgs := []tc.CacheGroupName{}
	for i := 0; i < numCGs; i++ {
		cgs = append(cgs, tc.CacheGroupName(test.RandStr()))
	}

	serverDSes := map[tc.CacheName][]tc.DeliveryServiceName{}
	for _, ca := range caches {
		for i := 0; i < numCacheDSes; i++ {
			serverDSes[ca] = append(serverDSes[ca], dses[test.RandIntn(len(dses))])
		}
	}

	dsServers := map[tc.DeliveryServiceName][]tc.CacheName{}
	for server, dses := range serverDSes {
		for _, ds := range dses {
			dsServers[ds] = append(dsServers[ds], server)
		}
	}

	serverCGs := map[tc.CacheName]tc.CacheGroupName{}
	for _, cache := range caches {
		serverCGs[cache] = cgs[test.RandIntn(len(cgs))]
	}

	serverTypes := map[tc.CacheName]tc.CacheType{}
	for _, cache := range caches {
		serverTypes[cache] = types[test.RandIntn(len(types))]
	}

	tod := todata.New()
	tod.DeliveryServiceServers = dsServers
	tod.ServerDeliveryServices = serverDSes
	tod.ServerTypes = serverTypes
	tod.ServerCachegroups = serverCGs
	return *tod
}

func randCachesPrecomputedData(caches []tc.CacheName, toData todata.TOData) map[tc.CacheName]cache.PrecomputedData {
	prc := map[tc.CacheName]cache.PrecomputedData{}
	for _, ca := range caches {
		prc[ca] = randPrecomputedData(toData)
	}
	return prc
}

func randPrecomputedData(toData todata.TOData) cache.PrecomputedData {
	dsStats := randDsStats(toData)
	dsTotal := uint64(0)
	for _, stat := range dsStats {
		dsTotal += stat.OutBytes
	}
	return cache.PrecomputedData{
		DeliveryServiceStats: dsStats,
		OutBytes:             dsTotal,
		MaxKbps:              test.RandInt64(),
		Errors:               randErrs(),
		Reporting:            true,
	}
}

func randDsStats(toData todata.TOData) map[string]*cache.DSStat {
	a := map[string]*cache.DSStat{}
	for ds, _ := range toData.DeliveryServiceServers {
		a[string(ds)] = randAStat()
	}
	return a
}

func randAStat() *cache.DSStat {
	return &cache.DSStat{
		InBytes:   uint64(test.RandIntn(1000)),
		OutBytes:  uint64(test.RandIntn(1000)),
		Status2xx: uint64(test.RandIntn(1000)),
		Status3xx: uint64(test.RandIntn(1000)),
		Status4xx: uint64(test.RandIntn(1000)),
		Status5xx: uint64(test.RandIntn(1000)),
	}
}

func randErrs() []error {
	if test.RandBool() {
		return []error{}
	}
	num := 5
	errs := []error{}
	for i := 0; i < num; i++ {
		errs = append(errs, errors.New(test.RandStr()))
	}
	return errs
}

func TestAddLastStatsToStatCacheStatsNilVals(t *testing.T) {
	// test that addLastStatsToStatCacheStats doesn't panic with nil values
	addLastStatsToStatCacheStats(nil, nil)
	addLastStatsToStatCacheStats(&dsdata.StatCacheStats{}, nil)
	addLastStatsToStatCacheStats(nil, &dsdata.LastStatsData{})
}
