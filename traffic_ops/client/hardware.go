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

type hardwareResponse struct {
	Version  string     `json:"version"`
	Response []Hardware `json:"response"`
}

// Hardware ...
type Hardware struct {
	ID          string `json:"serverId"`
	HostName    string `json:"serverHostName"`
	LastUpdated string `json:"lastUpdated"`
	Value       string `json:"val"`
	Description string `json:"description"`
}

// Hardware gets an array of Hardware
func (to *Session) Hardware() ([]Hardware, error) {
	body, err := to.getBytes("/api/1.1/hwinfo.json")
	if err != nil {
		return nil, err
	}
	hardwareList, err := hardwareUnmarshall(body)
	return hardwareList.Response, err
}

func hardwareUnmarshall(body []byte) (hardwareResponse, error) {

	var data hardwareResponse
	err := json.Unmarshal(body, &data)
	return data, err
}
