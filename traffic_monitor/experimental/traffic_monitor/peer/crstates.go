package peer

import (
	"encoding/json"
)

type Crstates struct {
	Caches          map[string]IsAvailable     `json:"caches"`
	Deliveryservice map[string]Deliveryservice `json:"deliveryServices"`
}

type IsAvailable struct {
	IsAvailable bool `json:"isAvailable"`
}

type Deliveryservice struct {
	DisabledLocations []string `json:"disabledLocations"`
	IsAvailable       bool     `json:"isAvailable"`
}

func CrStatesUnMarshall(body []byte) (Crstates, error) {
	var crStates Crstates

	err := json.Unmarshal(body, &crStates)
	return crStates, err
}

func CrStatesMarshall(states Crstates) ([]byte, error) {
	return json.Marshal(states)
}
