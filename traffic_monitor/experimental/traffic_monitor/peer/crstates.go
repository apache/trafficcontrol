package peer

import (
	"encoding/json"
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
