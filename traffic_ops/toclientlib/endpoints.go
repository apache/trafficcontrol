package toclientlib

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

const apiBaseStr = "/api/"

// APIBase returns the base API string for HTTP requests, such as /api/3.1.
func (to *TOClient) APIBase() string {
	return apiBaseStr + to.APIVersion()
}

// APIVersion is the version of the Traffic Ops API this client will use for requests.
// If the client was created with any function except Login, or with UseLatestSupportedAPI false,
// this will be LatestAPIVersion().
// Otherwise, it will be the version dynamically determined to be the latest the Traffic Ops Server supports.
func (to *TOClient) APIVersion() string {
	if to.latestSupportedAPI != "" {
		return to.latestSupportedAPI
	}
	return to.LatestAPIVersion()
}

// LatestAPIVersion returns the latest Traffic Ops API version this client supports.
func (to *TOClient) LatestAPIVersion() string {
	return to.apiVersions[0]
}
