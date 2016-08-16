package manager

import (
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/peer"
	"sync"
)

// This could be made lock-free, if the performance was necessary
// TODO add separate locks for Caches and Deliveryservice maps?
type CRStatesThreadsafe struct {
	crStates *peer.Crstates
	m        *sync.Mutex
}

func NewCRStatesThreadsafe() CRStatesThreadsafe {
	crs := peer.NewCrstates()
	return CRStatesThreadsafe{m: &sync.Mutex{}, crStates: &crs}
}

// TODO add GetCaches, GetDeliveryservices?
func (t CRStatesThreadsafe) Get() peer.Crstates {
	t.m.Lock()
	defer func() {
		t.m.Unlock()
	}()
	return t.crStates.Copy()
}

// TODO add GetCaches, GetDeliveryservices?
func (t CRStatesThreadsafe) GetDeliveryServices() map[string]peer.Deliveryservice {
	t.m.Lock()
	defer func() {
		t.m.Unlock()
	}()
	return t.crStates.CopyDeliveryservices()
}

func (t CRStatesThreadsafe) GetCache(name string) peer.IsAvailable {
	t.m.Lock()
	defer func() {
		t.m.Unlock()
	}()
	return t.crStates.Caches[name]
}

func (t CRStatesThreadsafe) GetCaches() map[string]peer.IsAvailable {
	t.m.Lock()
	defer func() {
		t.m.Unlock()
	}()
	return t.crStates.CopyCaches()
}

func (t CRStatesThreadsafe) GetDeliveryService(name string) peer.Deliveryservice {
	t.m.Lock()
	defer func() {
		t.m.Unlock()
	}()
	return t.crStates.Deliveryservice[name]
}

func (o CRStatesThreadsafe) Set(newCRStates peer.Crstates) {
	o.m.Lock()
	*o.crStates = newCRStates
	o.m.Unlock()
}

func (o CRStatesThreadsafe) SetCache(cacheName string, available peer.IsAvailable) {
	o.m.Lock()
	o.crStates.Caches[cacheName] = available
	o.m.Unlock()
}

func (o CRStatesThreadsafe) DeleteCache(name string) {
	o.m.Lock()
	delete(o.crStates.Caches, name)
	o.m.Unlock()
}

func (o CRStatesThreadsafe) SetDeliveryService(name string, ds peer.Deliveryservice) {
	o.m.Lock()
	o.crStates.Deliveryservice[name] = ds
	o.m.Unlock()
}

func (o CRStatesThreadsafe) SetDeliveryServices(deliveryServices map[string]peer.Deliveryservice) {
	o.m.Lock()
	o.crStates.Deliveryservice = deliveryServices
	o.m.Unlock()
}

// func (s CRStates) AddDeliveryServiceDisableLocation(cacheName string) {
// 	o.m.Lock()
// 	o.crStates.Deliveryservice[name].DisabledLocations[cacheName] = struct{}{}
// 	o.m.Unlock()
// }

// func (s CRStates) RemoveDeliveryServiceDisableLocation(cacheName string) {
// 	o.m.Lock()
// 	delete(o.crStates.Deliveryservice[name].DisabledLocations, cacheName)
// 	o.m.Unlock()
// }

// This could be made lock-free, if the performance was necessary
type CRStatesPeersThreadsafe struct {
	crStates map[string]peer.Crstates // TODO change string to type?
	m        *sync.Mutex
}

func NewCRStatesPeersThreadsafe() CRStatesPeersThreadsafe {
	return CRStatesPeersThreadsafe{m: &sync.Mutex{}, crStates: map[string]peer.Crstates{}}
}

func (t CRStatesPeersThreadsafe) Get() map[string]peer.Crstates {
	t.m.Lock()
	m := map[string]peer.Crstates{}
	for k, v := range t.crStates {
		m[k] = v.Copy()
	}
	t.m.Unlock()
	return m
}

// TODO use reader-writer mutex
func (o CRStatesPeersThreadsafe) Set(peerName string, peerState peer.Crstates) {
	o.m.Lock()
	o.crStates[peerName] = peerState
	o.m.Unlock()
}
