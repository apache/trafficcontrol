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
	"testing"
)

func TestDeliveryTenantsEp(t *testing.T) {
	test_helper.Context(t, "Given the need to test that DeliveryServices uses the correct URL")

	ep := tenantsEp()
	expected := "/api/1.2/tenants"
	if ep != expected {
		test_helper.Error(t, "Should get back %s for \"tenantsEp\", got: %s", expected, ep)
	} else {
		test_helper.Success(t, "Should be able to get the correct tenants endpoint")
	}
}
