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

package fixtures

import tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"

// Tenants returns a default DeliveryServiceResponse to be used for testing.
func Tenants() *tc.GetTenantsResponse {
	return &tc.GetTenantsResponse{
		Response: []tc.Tenant{
			tc.Tenant{
				ID:         1,
				Active:     true,
				Name:       "test-tenant",
				ParentName: "test-tenant-parent",
				ParentID:   2,
			},
		},
	}
}
