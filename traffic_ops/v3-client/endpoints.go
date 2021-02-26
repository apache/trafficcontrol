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

// DEPRECATED: All new code should us Session.APIBase().
// This isn't public, but only exists for deprecated public constants. It should be removed when they are.
const apiBase = "/api/3.1"

// apiVersions is the list of minor API versions in this client's major version.
// This should be all minor versions from 0 up to the latest minor in Traffic Control
// as of this client code.
//
// Versions are ordered latest-first.
func apiVersions() []string {
	return []string{
		"3.1",
		"3.0",
	}
}
