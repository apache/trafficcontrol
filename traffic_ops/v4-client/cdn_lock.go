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
	"fmt"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

// apiCDNLocks is the API version-relative path for the /cdn_locks API endpoint.
const apiCDNLocks = "/cdn_locks"
// apiAdminCDNLocks is the API version-relative path for the /cdn_locks/admin API endpoint.
const apiAdminCDNLocks = "/cdn_locks/admin"

// CreateCdnLock creates a CDN Lock.
func (to *Session) CreateCdnLock(cdnLock tc.CdnLock, opts RequestOptions) (tc.CdnLockCreateResponse, toclientlib.ReqInf, error) {
	var response tc.CdnLockCreateResponse
	var alerts tc.Alerts
	reqInf, err := to.post(apiCDNLocks, opts, cdnLock, &alerts)
	response.Response = cdnLock
	response.Alerts = alerts
	return response, reqInf, err
}

// GetCdnLocks retrieves the CDN locks based on the passed in parameters.
func (to *Session) GetCdnLocks(opts RequestOptions) (tc.CdnLocksGetResponse, toclientlib.ReqInf, error) {
	var data tc.CdnLocksGetResponse
	reqInf, err := to.get(fmt.Sprintf(apiCDNLocks), opts, &data)
	return data, reqInf, err
}

// DeleteCdnLocks deletes the CDN lock of a particular(requesting) user.
func (to *Session) DeleteCdnLocks(opts RequestOptions) (tc.CdnLockDeleteResponse, toclientlib.ReqInf, error) {
	var data tc.CdnLockDeleteResponse
	reqInf, err := to.del(fmt.Sprintf(apiCDNLocks), opts, &data)
	return data, reqInf, err
}

// AdminDeleteCdnLocks hits the endpoint an admin user would use to delete somebody else's lock.
func (to *Session) AdminDeleteCdnLocks(opts RequestOptions) (tc.CdnLockDeleteResponse, toclientlib.ReqInf, error) {
	var data tc.CdnLockDeleteResponse
	reqInf, err := to.del(fmt.Sprintf(apiAdminCDNLocks), opts, &data)
	return data, reqInf, err
}