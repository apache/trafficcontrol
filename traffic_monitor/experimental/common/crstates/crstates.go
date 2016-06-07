package crstates

type Cache struct {
	Name      string `json:"name,omitempty"`
	Available bool   `json:"isAvailable"`
}

type DeliveryService struct {
	Name              string   `json:"name,omitempty"`
	DisabledLocations []string `json:"disabledLocations"`
	Available         bool     `json:"isAvailable"`
}

type CRStates struct {
	Caches           map[string]Cache           `json:"caches"`
	DeliveryServices map[string]DeliveryService `json:"deliveryServices"`
}
