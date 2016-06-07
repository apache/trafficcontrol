package cache

import (
	"encoding/json"
	"io"
)

type Astats struct {
	Ats    map[string]interface{} `json:"ats"`
	System AstatsSystem           `json:"system"`
}

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

type AstatsAdapter struct{}

func Unmarshal(body []byte) (Astats, error) {
	var aStats Astats
	err := json.Unmarshal(body, &aStats)
	return aStats, err
}

func (AstatsAdapter) Transform(r io.Reader) ([]Astats, error) {
	dec := json.NewDecoder(r)
	var as []Astats

	for {
		var a Astats
		if err := dec.Decode(&a); err == io.EOF {
			return as, nil
		} else if err != nil {
			return as, err
		}
		as = append(as, a)
	}
}
