package client

import (
	"encoding/json"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_v11_OSVERSIONS = "/api/1.1/osversions"
)

// GetOSVersions GET all available Operating System (OS) versions for ISO generation,
// as well as the name of the directory where the "kickstarter" files are found.
// Structure of returned map:
//  key:   Name of OS
//  value: Directory where the ISO source can be found
func (to *Session) GetOSVersions() (map[string]string, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_v11_OSVERSIONS, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.OSVersionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Versions, reqInf, nil
}
