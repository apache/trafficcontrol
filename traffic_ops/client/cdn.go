/*
   Copyright 2015 Comcast Cable Communications Management, LLC

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package client

import (
	"encoding/json"
	"fmt"
)

// CDNResponse ...
type CDNResponse struct {
	Version  string `json:"version"`
	Response []CDN  `json:"response"`
}

// CDN ...
type CDN struct {
	Name        string `json:"name"`
	LastUpdated string `json:"lastUpdated"`
}

// Cdns gets an array of CDNs
func (to *Session) Cdns() ([]CDN, error) {
	body, err := to.getBytes("/api/1.2/cdns.json")
	if err != nil {
		return nil, err
	}

	var cdn CDNResponse
	if err := json.Unmarshal(body, &cdn); err != nil {
		return nil, err
	}
	return cdn.Response, err
}

// CdnName gets an array of CDNs
func (to *Session) CdnName(name string) ([]CDN, error) {
	url := fmt.Sprintf("/api/1.2/cdns/name/%s.json", name)
	body, err := to.getBytes(url)
	if err != nil {
		return nil, err
	}

	var cdn CDNResponse
	if err := json.Unmarshal(body, &cdn); err != nil {
		return nil, err
	}
	return cdn.Response, err
}
