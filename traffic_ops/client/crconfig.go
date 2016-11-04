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

import "fmt"

// CRConfigRaw Deprecated: use GetCRConfig instead
func (to *Session) CRConfigRaw(cdn string) ([]byte, error) {
	bytes, _, err := to.GetCRConfig(cdn)
	return bytes, err
}

// GetCRConfig returns the raw JSON bytes of the CRConfig from Traffic Ops, and whether the bytes were from the client's internal cache.
func (to *Session) GetCRConfig(cdn string) ([]byte, CacheHitStatus, error) {
	url := fmt.Sprintf("/CRConfig-Snapshots/%s/CRConfig.json", cdn)
	return to.getBytesWithTTL(url, tmPollingInterval)
}
