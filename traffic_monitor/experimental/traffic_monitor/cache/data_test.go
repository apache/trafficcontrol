package cache

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
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
