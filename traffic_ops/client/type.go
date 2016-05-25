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

import "encoding/json"

// TypeResponse ...
type TypeResponse struct {
	Version  string `json:"version"`
	Response []Type `json:"response"`
}

// Type contains information about a given Type in Traffic Ops.
type Type struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Types gets an array of Types.
func (to *Session) Types() ([]Type, error) {
	url := "/api/1.2/types.json"
	resp, err := to.request(url, nil)
	if err != nil {
		return nil, err
	}

	var data TypeResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Response, nil
}
