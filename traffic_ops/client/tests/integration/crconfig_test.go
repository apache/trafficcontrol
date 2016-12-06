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

import "testing"

//TestCachegroupResults compares the results of the Cachegroup api and Cachegroup client
func TestGetCrConfig(t *testing.T) {
	//Get a CDN from the to client
	cdn, err := GetCdn()
	if err != nil {
		t.Errorf("TestGetCrConfig -- Could not get CDNs from TO...%v\n", err)
	}

	crConfig, cacheHitStatus, err := to.GetCRConfig(cdn.Name)
	if err != nil {
		t.Errorf("Could not get CrConfig for %s.  Error is...%v\n", cdn.Name, err)
	}

	if cacheHitStatus == "" {
		t.Error("cacheHitStatus is empty...")
	}

	if len(crConfig) == 0 {
		t.Error("Raw CrConfig reponse was 0...")
	}
}
