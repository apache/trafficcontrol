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

import "encoding/json"

// ProfileResponse ...
type ProfileResponse struct {
	Response []Profile `json:"response"`
}

// Profile ...
type Profile struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	LastUpdated string `json:"lastUpdated"`
}

// Profiles gets an array of Profiles
func (to *Session) Profiles() ([]Profile, error) {
	url := "/api/1.2/profiles.json"
	resp, err := to.request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data ProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Response, nil
}
