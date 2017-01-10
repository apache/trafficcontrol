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

package integration

import (
	"encoding/json"
	"testing"

	traffic_ops "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

//TestCachegroupResults compares the results of the Cachegroup api and Cachegroup client
func TestCachegroups(t *testing.T) {
	//Get Cachegroups Data from API
	resp, err := Request(*to, "GET", "/api/1.2/cachegroups.json", nil)
	if err != nil {
		t.Errorf("Could not get cachegroups.json reponse was: %v\n", err)
	}

	defer resp.Body.Close()
	var apiCgRes traffic_ops.CacheGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiCgRes); err != nil {
		t.Errorf("Could not decode Cachegroup json.  Error is: %v\n", err)
	}
	apiCgs := apiCgRes.Response

	clientCgs, err := to.CacheGroups()
	if err != nil {
		t.Errorf("Could not get Cachegroups from client.  Error is: %v\n", err)
	}

	if len(apiCgs) != len(clientCgs) {
		t.Errorf("Array lengths from client and API are different...API = %d, Client = %d\n", len(apiCgs), len(clientCgs))
	}

	for _, apiCg := range apiCgs {
		matchFound := false
		for _, clientCg := range clientCgs {
			if clientCg.Name != apiCg.Name {
				continue
			}
			matchFound = true
			//compare API results and client results
			if apiCg.Name != clientCg.Name {
				t.Errorf("Cachegroup Name from client and API are different...API = %s, Client = %s\n", apiCg.Name, clientCg.Name)
			}
			if apiCg.ShortName != clientCg.ShortName {
				t.Errorf("Cachegroup ShortName from client and API are different...API = %s, Client = %s\n", apiCg.ShortName, clientCg.ShortName)
			}
			if apiCg.LastUpdated != clientCg.LastUpdated {
				t.Errorf("Cachegroup Last Updated from client and API are different...API = %s, Client = %s\n", apiCg.Name, clientCg.Name)
			}
			if apiCg.Latitude != clientCg.Latitude {
				t.Errorf("Cachegroup Latitude from client and API are different...API = %f, Client = %f\n", apiCg.Latitude, clientCg.Latitude)
			}
			if apiCg.Longitude != clientCg.Longitude {
				t.Errorf("Cachegroup Longitude from client and API are different...API = %f, Client = %f\n", apiCg.Longitude, clientCg.Longitude)
			}
			if apiCg.ParentName != clientCg.ParentName {
				t.Errorf("Cachegroup ParentName from client and API are different...API = %s, Client = %s\n", apiCg.ParentName, clientCg.ParentName)
			}
			if apiCg.Type != clientCg.Type {
				t.Errorf("Cachegroup Type from client and API are different...API = %s, Client = %s\n", apiCg.Type, clientCg.Type)
			}
		}
		if !matchFound {
			t.Errorf("A match for %s from the API was not found in the client results\n", apiCg.Name)
		}
	}
}
