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

import "github.com/apache/incubator-trafficcontrol/traffic_ops/client"

// Users returns a default UserResponse to be used for testing.
func Users() *client.UserResponse {
	return &client.UserResponse{
		Response: []client.User{
			client.User{
				Username:     "bsmith",
				PublicSSHKey: "some-ssh-key",
				Role:         "3",
				RoleName:     "operations",
				Email:        "bobsmith@email.com",
				FullName:     "Bob Smith",
			},
		},
	}
}
