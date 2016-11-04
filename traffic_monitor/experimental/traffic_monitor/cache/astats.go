package cache

import (
	"encoding/json"
)

// Astats contains ATS data returned from the Astats ATS plugin. This includes generic stats, as well as fixed system stats.
type Astats struct {
	Ats    map[string]interface{} `json:"ats"`
	System AstatsSystem           `json:"system"`
}

// AstatsSystem represents fixed system stats returne from ATS by the Astats plugin.
type AstatsSystem struct {
	InfName           string `json:"inf.name"`
	InfSpeed          int    `json:"inf.speed"`
	ProcNetDev        string `json:"proc.net.dev"`
	ProcLoadavg       string `json:"proc.loadavg"`
	ConfigLoadRequest int    `json:"configReloadRequests"`
	LastReloadRequest int    `json:"lastReloadRequest"`
	ConfigReloads     int    `json:"configReloads"`
	LastReload        int    `json:"lastReload"`
	AstatsLoad        int    `json:"astatsLoad"`
}

// Unmarshal unmarshalls the given bytes, which must be JSON Astats data, into an Astats object.
func Unmarshal(body []byte) (Astats, error) {
	var aStats Astats
	err := json.Unmarshal(body, &aStats)
	return aStats, err
}
