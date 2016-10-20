package peer

import (
	"encoding/json"
	"sync"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
)

type Crstates struct {
	Caches          map[enum.CacheName]IsAvailable               `json:"caches"`
	Deliveryservice map[enum.DeliveryServiceName]Deliveryservice `json:"deliveryServices"`
}

func NewCrstates() Crstates {
	return Crstates{
		Caches:          map[enum.CacheName]IsAvailable{},
		Deliveryservice: map[enum.DeliveryServiceName]Deliveryservice{},
	}
}

func (a Crstates) Copy() Crstates {
	b := NewCrstates()
	for k, v := range a.Caches {
		b.Caches[k] = v
	}
	for k, v := range a.Deliveryservice {
		b.Deliveryservice[k] = v
	}
	return b
}

func (a Crstates) CopyDeliveryservices() map[enum.DeliveryServiceName]Deliveryservice {
	b := map[enum.DeliveryServiceName]Deliveryservice{}
	for k, v := range a.Deliveryservice {
		b[k] = v
	}
	return b
}

func (a Crstates) CopyCaches() map[enum.CacheName]IsAvailable {
	b := map[enum.CacheName]IsAvailable{}
	for k, v := range a.Caches {
		b[k] = v
	}
	return b
}

type IsAvailable struct {
	IsAvailable bool `json:"isAvailable"`
}

type Deliveryservice struct {
	DisabledLocations []enum.CacheName `json:"disabledLocations"`
	IsAvailable       bool             `json:"isAvailable"`
}

func CrstatesUnMarshall(body []byte) (Crstates, error) {
	var crStates Crstates

	err := json.Unmarshal(body, &crStates)
	return crStates, err
}

func CrstatesMarshall(states Crstates) ([]byte, error) {
	return json.Marshal(states)
}

// This could be made lock-free, if the performance was necessary
// TODO add separate locks for Caches and Deliveryservice maps?
type CRStatesThreadsafe struct {
	crStates *Crstates
	m        *sync.RWMutex
}

func NewCRStatesThreadsafe() CRStatesThreadsafe {
	crs := NewCrstates()
	return CRStatesThreadsafe{m: &sync.RWMutex{}, crStates: &crs}
}

// TODO add GetCaches, GetDeliveryservices?
func (t *CRStatesThreadsafe) Get() Crstates {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.crStates.Copy()
}

// TODO add GetCaches, GetDeliveryservices?
func (t *CRStatesThreadsafe) GetDeliveryServices() map[enum.DeliveryServiceName]Deliveryservice {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.crStates.CopyDeliveryservices()
}

func (t *CRStatesThreadsafe) GetCache(name enum.CacheName) IsAvailable {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.crStates.Caches[name]
}

func (t *CRStatesThreadsafe) GetCaches() map[enum.CacheName]IsAvailable {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.crStates.CopyCaches()
}

func (t *CRStatesThreadsafe) GetDeliveryService(name enum.DeliveryServiceName) Deliveryservice {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.crStates.Deliveryservice[name]
}

func (o *CRStatesThreadsafe) Set(newCRStates Crstates) {
	o.m.Lock()
	*o.crStates = newCRStates
	o.m.Unlock()
}

func (o *CRStatesThreadsafe) SetCache(cacheName enum.CacheName, available IsAvailable) {
	o.m.Lock()
	o.crStates.Caches[cacheName] = available
	o.m.Unlock()
}

func (o *CRStatesThreadsafe) DeleteCache(name enum.CacheName) {
	o.m.Lock()
	delete(o.crStates.Caches, name)
	o.m.Unlock()
}

func (o *CRStatesThreadsafe) SetDeliveryService(name enum.DeliveryServiceName, ds Deliveryservice) {
	o.m.Lock()
	o.crStates.Deliveryservice[name] = ds
	o.m.Unlock()
}

func (o *CRStatesThreadsafe) SetDeliveryServices(deliveryServices map[enum.DeliveryServiceName]Deliveryservice) {
	o.m.Lock()
	o.crStates.Deliveryservice = deliveryServices
	o.m.Unlock()
}

func (o *CRStatesThreadsafe) DeleteDeliveryService(name enum.DeliveryServiceName) {
	o.m.Lock()
	delete(o.crStates.Deliveryservice, name)
	o.m.Unlock()
}

// This could be made lock-free, if the performance was necessary
type CRStatesPeersThreadsafe struct {
	crStates map[enum.TrafficMonitorName]Crstates
	m        *sync.RWMutex
}

func NewCRStatesPeersThreadsafe() CRStatesPeersThreadsafe {
	return CRStatesPeersThreadsafe{m: &sync.RWMutex{}, crStates: map[enum.TrafficMonitorName]Crstates{}}
}

func (t *CRStatesPeersThreadsafe) Get() map[enum.TrafficMonitorName]Crstates {
	t.m.RLock()
	m := map[enum.TrafficMonitorName]Crstates{}
	for k, v := range t.crStates {
		m[k] = v.Copy()
	}
	t.m.RUnlock()
	return m
}

func (o *CRStatesPeersThreadsafe) Set(peerName enum.TrafficMonitorName, peerState Crstates) {
	o.m.Lock()
	o.crStates[peerName] = peerState
	o.m.Unlock()
}
