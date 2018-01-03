package availableservers

import (
	"errors"
	"fmt"
	"sync"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

type AvailableServersMap map[tc.DeliveryServiceName]map[tc.CacheGroupName][]tc.CacheName

// AvailableServers provides access to the currently available servers, by Delivery Service and Cache Group. This is safe for access by multiple goroutines.
type AvailableServers struct {
	p **AvailableServersMap
	m *sync.RWMutex
}

func New() AvailableServers {
	mp := &AvailableServersMap{}
	return AvailableServers{p: &mp, m: &sync.RWMutex{}}
}

func (a *AvailableServers) Get(ds tc.DeliveryServiceName, cg tc.CacheGroupName) ([]tc.CacheName, error) {
	a.m.RLock()
	s := *a.p
	a.m.RUnlock()

	cgs, ok := (*s)[ds]
	if !ok {
		return nil, errors.New("deliveryservice not found")
	}
	cs, ok := cgs[cg]
	if !ok {
		return nil, errors.New("cachegroup not found")
	}
	return cs, nil
}

func (a *AvailableServers) Set(m AvailableServersMap) {
	a.m.Lock()
	defer a.m.Unlock()
	*a.p = &m
}

// TODO put in _test.go file
func Test() {
	a := New()

	as := map[tc.DeliveryServiceName]map[tc.CacheGroupName][]tc.CacheName{}
	as[tc.DeliveryServiceName("dsOne")] = map[tc.CacheGroupName][]tc.CacheName{}

	cs := as[tc.DeliveryServiceName("dsOne")]
	cs[tc.CacheGroupName("cgOne")] = []tc.CacheName{"cacheOne", "cacheTwo"}

	fmt.Printf("testAvailableServers as %+v\n", as)

	a.Set(as)

	newCs, err := a.Get(tc.DeliveryServiceName("dsOne"), tc.CacheGroupName("cgOne"))

	if err != nil {
		fmt.Println("testAvailableServers err ", err.Error())
	}
	fmt.Println("testAvailableServers caches ", newCs)
}
