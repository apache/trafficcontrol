package peer

import (
	"encoding/json"
	"sync"
)

type Crstates struct {
	Caches          map[string]IsAvailable     `json:"caches"`
	Deliveryservice map[string]Deliveryservice `json:"deliveryServices"`
}

func NewCrstates() Crstates {
	return Crstates{
		Caches:          map[string]IsAvailable{},
		Deliveryservice: map[string]Deliveryservice{},
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

func (a Crstates) CopyDeliveryservices() map[string]Deliveryservice {
	b := map[string]Deliveryservice{}
	for k, v := range a.Deliveryservice {
		b[k] = v
	}
	return b
}

func (a Crstates) CopyCaches() map[string]IsAvailable {
	b := map[string]IsAvailable{}
	for k, v := range a.Caches {
		b[k] = v
	}
	return b
}

type IsAvailable struct {
	IsAvailable bool `json:"isAvailable"`
}

type Deliveryservice struct {
	DisabledLocations []string `json:"disabledLocations"`
	IsAvailable       bool     `json:"isAvailable"`
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
	m        *sync.Mutex
}

func NewCRStatesThreadsafe() CRStatesThreadsafe {
	crs := NewCrstates()
	return CRStatesThreadsafe{m: &sync.Mutex{}, crStates: &crs}
}

// TODO add GetCaches, GetDeliveryservices?
func (t CRStatesThreadsafe) Get() Crstates {
	t.m.Lock()
	defer func() {
		t.m.Unlock()
	}()
	return t.crStates.Copy()
}

// TODO add GetCaches, GetDeliveryservices?
func (t CRStatesThreadsafe) GetDeliveryServices() map[string]Deliveryservice {
	t.m.Lock()
	defer func() {
		t.m.Unlock()
	}()
	return t.crStates.CopyDeliveryservices()
}

func (t CRStatesThreadsafe) GetCache(name string) IsAvailable {
	t.m.Lock()
	defer func() {
		t.m.Unlock()
	}()
	return t.crStates.Caches[name]
}

func (t CRStatesThreadsafe) GetCaches() map[string]IsAvailable {
	t.m.Lock()
	defer func() {
		t.m.Unlock()
	}()
	return t.crStates.CopyCaches()
}

func (t CRStatesThreadsafe) GetDeliveryService(name string) Deliveryservice {
	t.m.Lock()
	defer func() {
		t.m.Unlock()
	}()
	return t.crStates.Deliveryservice[name]
}

func (o CRStatesThreadsafe) Set(newCRStates Crstates) {
	o.m.Lock()
	*o.crStates = newCRStates
	o.m.Unlock()
}

func (o CRStatesThreadsafe) SetCache(cacheName string, available IsAvailable) {
	o.m.Lock()
	o.crStates.Caches[cacheName] = available
	o.m.Unlock()
}

func (o CRStatesThreadsafe) DeleteCache(name string) {
	o.m.Lock()
	delete(o.crStates.Caches, name)
	o.m.Unlock()
}

func (o CRStatesThreadsafe) SetDeliveryService(name string, ds Deliveryservice) {
	o.m.Lock()
	o.crStates.Deliveryservice[name] = ds
	o.m.Unlock()
}

func (o CRStatesThreadsafe) SetDeliveryServices(deliveryServices map[string]Deliveryservice) {
	o.m.Lock()
	o.crStates.Deliveryservice = deliveryServices
	o.m.Unlock()
}

// This could be made lock-free, if the performance was necessary
type CRStatesPeersThreadsafe struct {
	crStates map[string]Crstates // TODO change string to type?
	m        *sync.Mutex
}

func NewCRStatesPeersThreadsafe() CRStatesPeersThreadsafe {
	return CRStatesPeersThreadsafe{m: &sync.Mutex{}, crStates: map[string]Crstates{}}
}

func (t CRStatesPeersThreadsafe) Get() map[string]Crstates {
	t.m.Lock()
	m := map[string]Crstates{}
	for k, v := range t.crStates {
		m[k] = v.Copy()
	}
	t.m.Unlock()
	return m
}

// TODO use reader-writer mutex
func (o CRStatesPeersThreadsafe) Set(peerName string, peerState Crstates) {
	o.m.Lock()
	o.crStates[peerName] = peerState
	o.m.Unlock()
}
