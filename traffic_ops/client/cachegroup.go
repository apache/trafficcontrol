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
)

// CacheGroupResponse ...
type CacheGroupResponse struct {
	Version  string       `json:"version"`
	Response []CacheGroup `json:"response"`
}

// CacheGroup ...
type CacheGroup struct {
	Name        string  `json:"name"`
	ShortName   string  `json:"shortName"`
	Latitude    float64 `json:"latitude,string"`
	Longitude   float64 `json:"longitude,string"`
	ParentName  string  `json:"parentCachegroupName,omitempty"`
	Type        string  `json:"typeName,omitempty"`
	LastUpdated string  `json:"lastUpdated,omitempty"`
}

// CacheGroups gets the CacheGroups in an array of CacheGroup structs
// (note CacheGroup used to be called location)
func (to *Session) CacheGroups() ([]CacheGroup, error) {
	body, err := to.getBytes("/api/1.1/cachegroups.json")
	if err != nil {
		return nil, err
	}
	cgList, err := cgUnmarshall(body)
	if err != nil {
		return nil, err
	}
	return cgList.Response, nil
}

func cgUnmarshall(body []byte) (*CacheGroupResponse, error) {
	var data CacheGroupResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
