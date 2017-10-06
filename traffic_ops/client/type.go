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
	"errors"

	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

// Types gets an array of Types.
// optional parameter: userInTable
func (to *Session) Types(useInTable ...string) ([]tc.Type, error) {

	if len(useInTable) > 1 {
		return nil, errors.New("Please pass in a single value for the 'useInTable' parameter")
	}

	url := "/api/1.2/types.json"
	resp, err := to.request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data tc.TypeResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var types []tc.Type
	for _, d := range data.Response {
		if useInTable != nil {
			if d.UseInTable == useInTable[0] {
				types = append(types, d)
			}
		} else {
			types = append(types, d)
		}
	}

	return types, nil
}
