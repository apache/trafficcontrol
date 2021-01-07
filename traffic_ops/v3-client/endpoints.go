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

const apiBaseStr = "/api/"

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

// APIBase returns the base API string for HTTP requests, such as /api/3.1.
// If UseLatestSupportedAPI
func (sn *Session) APIBase() string {
	if sn.latestSupportedAPI == "" {
		// this will be the case for a Session initalized with any creation func other than New,
		// or if New was called with ClientOps.ForceLatestAPI.
		return apiBaseStr + apiVersions()[0]
	}
	return apiBaseStr + sn.latestSupportedAPI
}

// APIVersion is the version of the Traffic Ops API this client will use for requests.
// If the client was created with any function except New, or with UseLatestSupportedAPI false,
// this will be LatestAPIVersion().
// Otherwise, it will be the version dynamically determined to be the latest the Traffic Ops Server supports.
func (sn *Session) APIVersion() string {
	if sn.latestSupportedAPI != "" {
		return sn.latestSupportedAPI
	}
	return sn.LatestAPIVersion()
}

// LatestAPIVersion() returns the latest Traffic Ops API version this client supports.
func (sn *Session) LatestAPIVersion() string {
	return apiVersions()[0]
}
