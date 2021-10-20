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
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
	"github.com/apache/trafficcontrol/v6/traffic_monitor/dsdata"
)

func randBool() bool {
	return rand.Int()%2 == 0
}

func randStr() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_"
	num := 100
	s := ""
	for i := 0; i < num; i++ {
		s += string(chars[rand.Intn(len(chars))])
	}
	return s
}

func randAvailableStatuses() AvailableStatuses {
	a := AvailableStatuses{}
	num := 100
	for i := 0; i < num; i++ {
		cacheName := randStr()
		a[cacheName] = AvailableStatus{
			Available: AvailableTuple{
				IPv4: randBool(),
				IPv6: randBool(),
			},
			Status: randStr(),
		}
	}
	return a
}

func randStrIfaceMap() map[string]interface{} {
	m := map[string]interface{}{}
	num := 5
	for i := 0; i < num; i++ {
		m[randStr()] = randStr()
	}
	return m
}

func randStats() (Statistics, map[string]interface{}) {
	return randStatistics(), randStrIfaceMap()
}

func randStatistics() Statistics {
	return Statistics{
		Loadavg: Loadavg{
			One:              rand.Float64(),
			Five:             rand.Float64(),
			Fifteen:          rand.Float64(),
			CurrentProcesses: rand.Uint64(),
			TotalProcesses:   rand.Uint64(),
			LatestPID:        rand.Int63(),
		},
		Interfaces: map[string]Interface{
			randStr(): Interface{
				Speed:    rand.Int63(),
				BytesIn:  rand.Uint64(),
				BytesOut: rand.Uint64(),
			},
		},
	}
}

func randVitals() Vitals {
	return Vitals{
		LoadAvg:    rand.Float64(),
		BytesOut:   rand.Uint64(),
		BytesIn:    rand.Uint64(),
		KbpsOut:    rand.Int63(),
		MaxKbpsOut: rand.Int63(),
	}
}

func randStatMeta() dsdata.StatMeta {
	return dsdata.StatMeta{Time: rand.Int63()}
}

func randStatCacheStats() dsdata.StatCacheStats {
	return dsdata.StatCacheStats{
		OutBytes:    dsdata.StatInt{Value: rand.Int63(), StatMeta: randStatMeta()},
		IsAvailable: dsdata.StatBool{Value: randBool(), StatMeta: randStatMeta()},
		Status5xx:   dsdata.StatInt{Value: rand.Int63(), StatMeta: randStatMeta()},
		Status4xx:   dsdata.StatInt{Value: rand.Int63(), StatMeta: randStatMeta()},
		Status3xx:   dsdata.StatInt{Value: rand.Int63(), StatMeta: randStatMeta()},
		Status2xx:   dsdata.StatInt{Value: rand.Int63(), StatMeta: randStatMeta()},
		InBytes:     dsdata.StatFloat{Value: rand.Float64(), StatMeta: randStatMeta()},
		Kbps:        dsdata.StatFloat{Value: rand.Float64(), StatMeta: randStatMeta()},
		Tps5xx:      dsdata.StatFloat{Value: rand.Float64(), StatMeta: randStatMeta()},
		Tps4xx:      dsdata.StatFloat{Value: rand.Float64(), StatMeta: randStatMeta()},
		Tps3xx:      dsdata.StatFloat{Value: rand.Float64(), StatMeta: randStatMeta()},
		Tps2xx:      dsdata.StatFloat{Value: rand.Float64(), StatMeta: randStatMeta()},
		ErrorString: dsdata.StatString{Value: randStr(), StatMeta: randStatMeta()},
		TpsTotal:    dsdata.StatFloat{Value: rand.Float64(), StatMeta: randStatMeta()},
	}
}

func randStatCommon() dsdata.StatCommon {
	cachesReporting := map[tc.CacheName]bool{}
	num := 5
	for i := 0; i < num; i++ {
		cachesReporting[tc.CacheName(randStr())] = randBool()
	}
	return dsdata.StatCommon{
		CachesConfiguredNum: dsdata.StatInt{Value: rand.Int63(), StatMeta: randStatMeta()},
		CachesReporting:     cachesReporting,
		ErrorStr:            dsdata.StatString{Value: randStr(), StatMeta: randStatMeta()},
		StatusStr:           dsdata.StatString{Value: randStr(), StatMeta: randStatMeta()},
		IsHealthy:           dsdata.StatBool{Value: randBool(), StatMeta: randStatMeta()},
		IsAvailable:         dsdata.StatBool{Value: randBool(), StatMeta: randStatMeta()},
		CachesAvailableNum:  dsdata.StatInt{Value: rand.Int63(), StatMeta: randStatMeta()},
	}
}

func randAStat() *DSStat {
	return &DSStat{
		InBytes:   rand.Uint64(),
		OutBytes:  rand.Uint64(),
		Status2xx: rand.Uint64(),
		Status3xx: rand.Uint64(),
		Status4xx: rand.Uint64(),
		Status5xx: rand.Uint64(),
	}
}

func randDsStats() map[string]*DSStat {
	num := 5
	a := map[string]*DSStat{}
	for i := 0; i < num; i++ {
		a[randStr()] = randAStat()
	}
	return a
}
func randErrs() []error {
	if randBool() {
		return []error{}
	}
	num := 5
	errs := []error{}
	for i := 0; i < num; i++ {
		errs = append(errs, errors.New(randStr()))
	}
	return errs
}

func randPrecomputedData() PrecomputedData {
	return PrecomputedData{
		DeliveryServiceStats: randDsStats(),
		OutBytes:             rand.Uint64(),
		MaxKbps:              rand.Int63(),
		Errors:               randErrs(),
		Reporting:            randBool(),
	}
}

func randResult() Result {
	stats, misc := randStats()
	return Result{
		ID:              randStr(),
		Error:           fmt.Errorf(randStr()),
		Statistics:      stats,
		Time:            time.Now(),
		RequestTime:     time.Millisecond * time.Duration(rand.Int()),
		Vitals:          randVitals(),
		PollID:          uint64(rand.Int63()),
		PollFinished:    make(chan uint64),
		PrecomputedData: randPrecomputedData(),
		Available:       randBool(),
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
		a[tc.CacheName(randStr())] = randResultSlice()
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
		a[tc.CacheName(randStr())] = randResultSlice()
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

		cacheName := randStr()
		a[cacheName] = AvailableStatus{
			Available: AvailableTuple{
				randBool(),
				randBool(),
			},
			Status: randStr(),
		}

		if reflect.DeepEqual(a, b) {
			t.Errorf("expected a != b, actual a and b point to the same map: a: %+v", a)
		}
	}
}
