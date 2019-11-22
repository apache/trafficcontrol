package v14

import (
	"fmt"
	"testing"
)

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

func TestFederationFederationResolvers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants, DeliveryServices, CDNFederations, FederationResolvers}, func() {
		GetTestFederationFederationResolvers(t)
		AssignTestFederationFederationResolvers(t)
	})
}

func GetTestFederationFederationResolvers(t *testing.T) {
	for _, f := range testData.Federations {
		data, _, err := TOSession.GetFederationsFederationResolversByID(f.ID)
		if err != nil {
			t.Fatalf("Error getting federation_federation_resolvers by fed ID %d", f.ID)
		}
		fmt.Println("data", data)
	}
}

func AssignTestFederationFederationResolvers(t *testing.T) {
}
