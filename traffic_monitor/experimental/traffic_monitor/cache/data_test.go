package cache

import (
	"errors"
	"fmt"
	dsdata "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/deliveryservicedata"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
	"math/rand"
	"reflect"
	"testing"
	"time"
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
		a[enum.CacheName(randStr())] = AvailableStatus{Available: randBool(), Status: randStr()}
	}
	return a
}

func TestAvailableStatusesCopy(t *testing.T) {
	num := 100
	for i := 0; i < num; i++ {
		a := randAvailableStatuses()
		b := a.Copy()

		if !reflect.DeepEqual(a, b) {
			t.Errorf("expected a and b DeepEqual, actual copied map not equal", a, b)
		}

		// verify a and b don't point to the same map
		a[enum.CacheName(randStr())] = AvailableStatus{Available: randBool(), Status: randStr()}
		if reflect.DeepEqual(a, b) {
			t.Errorf("expected a != b, actual a and b point to the same map", a)
		}
	}
}

func randStrIfaceMap() map[string]interface{} {
	m := map[string]interface{}{}
	num := 5
	for i := 0; i < num; i++ {
		m[randStr()] = randStr()
	}
	return m
}

func randAstats() Astats {
	return Astats{
		Ats:    randStrIfaceMap(),
		System: randAstatsSystem(),
	}
}

func randAstatsSystem() AstatsSystem {
	return AstatsSystem{
		InfName:           randStr(),
		InfSpeed:          rand.Int(),
		ProcNetDev:        randStr(),
		ProcLoadavg:       randStr(),
		ConfigLoadRequest: rand.Int(),
		LastReloadRequest: rand.Int(),
		ConfigReloads:     rand.Int(),
		LastReload:        rand.Int(),
		AstatsLoad:        rand.Int(),
	}
}

func randVitals() Vitals {
	return Vitals{
		LoadAvg:    rand.Float64(),
		BytesOut:   rand.Int63(),
		BytesIn:    rand.Int63(),
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
		TpsTotal:    dsdata.StatInt{Value: rand.Int63(), StatMeta: randStatMeta()},
	}
}

func randStatCommon() dsdata.StatCommon {
	cachesReporting := map[enum.CacheName]bool{}
	num := 5
	for i := 0; i < num; i++ {
		cachesReporting[enum.CacheName(randStr())] = randBool()
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

func randDsStat() dsdata.Stat {
	num := 5
	cacheGroups := map[enum.CacheGroupName]dsdata.StatCacheStats{}
	types := map[enum.CacheType]dsdata.StatCacheStats{}
	caches := map[enum.CacheName]dsdata.StatCacheStats{}
	cachesTime := map[enum.CacheName]time.Time{}
	for i := 0; i < num; i++ {
		cacheGroups[enum.CacheGroupName(randStr())] = randStatCacheStats()
		types[enum.CacheType(randStr())] = randStatCacheStats()
		cachesTime[enum.CacheName(randStr())] = time.Now()
	}

	return dsdata.Stat{
		CommonStats:        randStatCommon(),
		CacheGroups:        cacheGroups,
		Types:              types,
		Caches:             caches,
		CachesTimeReceived: cachesTime,
		TotalStats:         randStatCacheStats(),
	}
}

func randDsStats() map[enum.DeliveryServiceName]dsdata.Stat {
	num := 5
	a := map[enum.DeliveryServiceName]dsdata.Stat{}
	for i := 0; i < num; i++ {
		a[enum.DeliveryServiceName(randStr())] = randDsStat()
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
		OutBytes:             rand.Int63(),
		MaxKbps:              rand.Int63(),
		Errors:               randErrs(),
		Reporting:            randBool(),
	}
}

func randResult() Result {
	return Result{
		ID:              enum.CacheName(randStr()),
		Error:           fmt.Errorf(randStr()),
		Astats:          randAstats(),
		Time:            time.Now(),
		RequestTime:     time.Millisecond * time.Duration(rand.Int()),
		Vitals:          randVitals(),
		PollID:          uint64(rand.Int63()),
		PollFinished:    make(chan uint64),
		PrecomputedData: randPrecomputedData(),
		Available:       randBool(),
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
		a[enum.CacheName(randStr())] = randResultSlice()
	}
	return a
}

func TestResultHistoryCopy(t *testing.T) {
	num := 5
	for i := 0; i < num; i++ {
		a := randResultHistory()
		b := a.Copy()

		if !reflect.DeepEqual(a, b) {
			t.Errorf("expected a and b DeepEqual, actual copied map not equal", a, b)
		}

		// verify a and b don't point to the same map
		a[enum.CacheName(randStr())] = randResultSlice()
		if reflect.DeepEqual(a, b) {
			t.Errorf("expected a != b, actual a and b point to the same map", a)
		}
	}
}
