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
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/dsdata"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/test"
)

func randAvailableStatuses() AvailableStatuses {
	a := AvailableStatuses{}
	num := 100
	for i := 0; i < num; i++ {
		cacheName := test.RandStr()
		a[cacheName] = AvailableStatus{
			Available: AvailableTuple{
				IPv4: test.RandBool(),
				IPv6: test.RandBool(),
			},
			Status: test.RandStr(),
		}
	}
	return a
}

func randStrIfaceMap() map[string]interface{} {
	m := map[string]interface{}{}
	num := 5
	for i := 0; i < num; i++ {
		m[test.RandStr()] = test.RandStr()
	}
	return m
}

func randStats() (Statistics, map[string]interface{}) {
	return randStatistics(), randStrIfaceMap()
}

func randStatistics() Statistics {
	return Statistics{
		Loadavg: Loadavg{
			One:              test.RandFloat64(),
			Five:             test.RandFloat64(),
			Fifteen:          test.RandFloat64(),
			CurrentProcesses: test.RandUint64(),
			TotalProcesses:   test.RandUint64(),
			LatestPID:        test.RandInt64(),
		},
		Interfaces: map[string]Interface{
			test.RandStr(): Interface{
				Speed:    test.RandInt64(),
				BytesIn:  test.RandUint64(),
				BytesOut: test.RandUint64(),
			},
		},
	}
}

func randVitals() Vitals {
	return Vitals{
		LoadAvg:    test.RandFloat64(),
		BytesOut:   test.RandUint64(),
		BytesIn:    test.RandUint64(),
		KbpsOut:    test.RandInt64(),
		MaxKbpsOut: test.RandInt64(),
	}
}

func randStatMeta() dsdata.StatMeta {
	return dsdata.StatMeta{Time: test.RandInt64()}
}

func randStatCacheStats() dsdata.StatCacheStats {
	return dsdata.StatCacheStats{
		OutBytes:    dsdata.StatInt{Value: test.RandInt64(), StatMeta: randStatMeta()},
		IsAvailable: dsdata.StatBool{Value: test.RandBool(), StatMeta: randStatMeta()},
		Status5xx:   dsdata.StatInt{Value: test.RandInt64(), StatMeta: randStatMeta()},
		Status4xx:   dsdata.StatInt{Value: test.RandInt64(), StatMeta: randStatMeta()},
		Status3xx:   dsdata.StatInt{Value: test.RandInt64(), StatMeta: randStatMeta()},
		Status2xx:   dsdata.StatInt{Value: test.RandInt64(), StatMeta: randStatMeta()},
		InBytes:     dsdata.StatFloat{Value: test.RandFloat64(), StatMeta: randStatMeta()},
		Kbps:        dsdata.StatFloat{Value: test.RandFloat64(), StatMeta: randStatMeta()},
		Tps5xx:      dsdata.StatFloat{Value: test.RandFloat64(), StatMeta: randStatMeta()},
		Tps4xx:      dsdata.StatFloat{Value: test.RandFloat64(), StatMeta: randStatMeta()},
		Tps3xx:      dsdata.StatFloat{Value: test.RandFloat64(), StatMeta: randStatMeta()},
		Tps2xx:      dsdata.StatFloat{Value: test.RandFloat64(), StatMeta: randStatMeta()},
		ErrorString: dsdata.StatString{Value: test.RandStr(), StatMeta: randStatMeta()},
		TpsTotal:    dsdata.StatFloat{Value: test.RandFloat64(), StatMeta: randStatMeta()},
	}
}

func randStatCommon() dsdata.StatCommon {
	cachesReporting := map[tc.CacheName]bool{}
	num := 5
	for i := 0; i < num; i++ {
		cachesReporting[tc.CacheName(test.RandStr())] = test.RandBool()
	}
	return dsdata.StatCommon{
		CachesConfiguredNum: dsdata.StatInt{Value: test.RandInt64(), StatMeta: randStatMeta()},
		CachesReporting:     cachesReporting,
		ErrorStr:            dsdata.StatString{Value: test.RandStr(), StatMeta: randStatMeta()},
		StatusStr:           dsdata.StatString{Value: test.RandStr(), StatMeta: randStatMeta()},
		IsHealthy:           dsdata.StatBool{Value: test.RandBool(), StatMeta: randStatMeta()},
		IsAvailable:         dsdata.StatBool{Value: test.RandBool(), StatMeta: randStatMeta()},
		CachesAvailableNum:  dsdata.StatInt{Value: test.RandInt64(), StatMeta: randStatMeta()},
	}
}

func randAStat() *DSStat {
	return &DSStat{
		InBytes:   test.RandUint64(),
		OutBytes:  test.RandUint64(),
		Status2xx: test.RandUint64(),
		Status3xx: test.RandUint64(),
		Status4xx: test.RandUint64(),
		Status5xx: test.RandUint64(),
	}
}

func randDsStats() map[string]*DSStat {
	num := 5
	a := map[string]*DSStat{}
	for i := 0; i < num; i++ {
		a[test.RandStr()] = randAStat()
	}
	return a
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

func randPrecomputedData() PrecomputedData {
	return PrecomputedData{
		DeliveryServiceStats: randDsStats(),
		OutBytes:             test.RandUint64(),
		MaxKbps:              test.RandInt64(),
		Errors:               randErrs(),
		Reporting:            test.RandBool(),
	}
}

func randResult() Result {
	stats, misc := randStats()
	return Result{
		ID:              test.RandStr(),
		Error:           fmt.Errorf(test.RandStr()),
		Statistics:      stats,
		Time:            time.Now(),
		RequestTime:     time.Millisecond * time.Duration(test.RandInt()),
		Vitals:          randVitals(),
		PollID:          uint64(test.RandInt64()),
		PollFinished:    make(chan uint64),
		PrecomputedData: randPrecomputedData(),
		Available:       test.RandBool(),
		Miscellaneous:   misc,
	}
}

func randResultSlice() []Result {
	a := []Result{}
	num := 5
	for i := 0; i < num; i++ {
		a = append(a, randResult())
	}
	return a
}

func randResultHistory() ResultHistory {
	a := ResultHistory{}
	num := 5
	for i := 0; i < num; i++ {
		a[tc.CacheName(test.RandStr())] = randResultSlice()
	}
	return a
}

func TestResultHistoryCopy(t *testing.T) {
	num := 5
	for i := 0; i < num; i++ {
		a := randResultHistory()
		b := a.Copy()

		if !reflect.DeepEqual(a, b) {
			t.Errorf("expected a and b DeepEqual, actual copied map not equal: a: %+v b: %+v", a, b)
		}

		// verify a and b don't point to the same map
		a[tc.CacheName(test.RandStr())] = randResultSlice()
		if reflect.DeepEqual(a, b) {
			t.Errorf("expected a != b, actual a and b point to the same map: %+v", a)
		}
	}
}

func TestAvailableStatusesCopy(t *testing.T) {
	num := 100
	for i := 0; i < num; i++ {
		a := randAvailableStatuses()
		b := a.Copy()

		if !reflect.DeepEqual(a, b) {
			t.Errorf("expected a and b DeepEqual, actual copied map not equal: a: %v b: %v", a, b)
		}

		cacheName := test.RandStr()
		a[cacheName] = AvailableStatus{
			Available: AvailableTuple{
				test.RandBool(),
				test.RandBool(),
			},
			Status: test.RandStr(),
		}

		if reflect.DeepEqual(a, b) {
			t.Errorf("expected a != b, actual a and b point to the same map: a: %+v", a)
		}
	}
}
