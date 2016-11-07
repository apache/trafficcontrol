/*

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

// ParamResponse ...
type ParamResponse struct {
	Version  string      `json:"version"`
	Response []Parameter `json:"response"`
}

// Parameter ...
type Parameter struct {
	Name        string `json:"name"`
	ConfigFile  string `json:"configFile"`
	Value       string `json:"Value"`
	LastUpdated string `json:"lastUpdated"`
}

// Parameters gets an array of parameter structs for the profile given
func (to *Session) Parameters(profileName string) ([]Parameter, error) {
	url := fmt.Sprintf("/api/1.2/parameters/profile/%s.json", profileName)
	resp, err := to.request(url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data ParamResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Response, nil
}
