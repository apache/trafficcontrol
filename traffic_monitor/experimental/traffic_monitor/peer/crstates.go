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

// CRStatesThreadsafe provides safe access for multiple goroutines to read a single Crstates object, with a single goroutine writer.
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

// Get returns the internal Crstates object for reading.
// TODO add GetCaches, GetDeliveryservices?
func (t *CRStatesThreadsafe) Get() Crstates {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.crStates.Copy()
}

// GetDeliveryServices returns the internal Crstates delivery services map for reading.
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

func (t *CRStatesThreadsafe) Set(newCRStates Crstates) {
	t.m.Lock()
	*t.crStates = newCRStates
	t.m.Unlock()
}

func (t *CRStatesThreadsafe) SetCache(cacheName enum.CacheName, available IsAvailable) {
	t.m.Lock()
	t.crStates.Caches[cacheName] = available
	t.m.Unlock()
}

func (t *CRStatesThreadsafe) DeleteCache(name enum.CacheName) {
	t.m.Lock()
	delete(t.crStates.Caches, name)
	t.m.Unlock()
}

func (t *CRStatesThreadsafe) SetDeliveryService(name enum.DeliveryServiceName, ds Deliveryservice) {
	t.m.Lock()
	t.crStates.Deliveryservice[name] = ds
	t.m.Unlock()
}

func (t *CRStatesThreadsafe) SetDeliveryServices(deliveryServices map[enum.DeliveryServiceName]Deliveryservice) {
	t.m.Lock()
	t.crStates.Deliveryservice = deliveryServices
	t.m.Unlock()
}

func (t *CRStatesThreadsafe) DeleteDeliveryService(name enum.DeliveryServiceName) {
	t.m.Lock()
	delete(t.crStates.Deliveryservice, name)
	t.m.Unlock()
}

// CRStatesPeersThreadsafe provides safe access for multiple goroutines to read a map of Traffic Monitor peers to their returned Crstates, with a single goroutine writer.
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

func (t *CRStatesPeersThreadsafe) Set(peerName enum.TrafficMonitorName, peerState Crstates) {
	t.m.Lock()
	t.crStates[peerName] = peerState
	t.m.Unlock()
}
