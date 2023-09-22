package client

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

import (
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiCDNLocks is the API version-relative path for the /cdn_locks API endpoint.
const apiCDNLocks = "/cdn_locks"

// CreateCDNLock creates a CDN Lock.
func (to *Session) CreateCDNLock(cdnLock tc.CDNLock, opts RequestOptions) (tc.CDNLockCreateResponse, toclientlib.ReqInf, error) {
	var response tc.CDNLockCreateResponse
	reqInf, err := to.post(apiCDNLocks, opts, cdnLock, &response)
	return response, reqInf, err
}

// GetCDNLocks retrieves the CDN locks based on the passed in parameters.
func (to *Session) GetCDNLocks(opts RequestOptions) (tc.CDNLocksGetResponse, toclientlib.ReqInf, error) {
	var data tc.CDNLocksGetResponse
	reqInf, err := to.get(apiCDNLocks, opts, &data)
	return data, reqInf, err
}

// DeleteCDNLocks deletes the CDN lock of a particular(requesting) user.
func (to *Session) DeleteCDNLocks(opts RequestOptions) (tc.CDNLockDeleteResponse, toclientlib.ReqInf, error) {
	var data tc.CDNLockDeleteResponse
	reqInf, err := to.del(apiCDNLocks, opts, &data)
	return data, reqInf, err
}
