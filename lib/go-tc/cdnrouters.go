package tc

type CDNRouting struct {
	StaticRoute       float64 `json:"staticRoute"`
	Geo               float64 `json:"geo"`
	Err               float64 `json:"err"`
	Fed               float64 `json:"fed"`
	CZ                float64 `json:"cz"`
	DeepCZ            float64 `json:"deepCz"`
	RegionalAlternate float64 `json:"regionalAlternate"`
	DSR               float64 `json:"dsr"`
	Miss              float64 `json:"miss"`
	RegionalDenied    float64 `json:"regionalDenied"`
}
