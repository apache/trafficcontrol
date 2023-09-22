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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiTenants is the API version-relative path to the /tenants API endpoint.
const apiTenants = "/tenants"

// apiTenantID is the API version-relative path to the /tenants/{{ID}} API endpoint.
const apiTenantID = apiTenants + "/%d"

// GetTenants retrieves all Tenants stored in Traffic Ops.
func (to *Session) GetTenants(opts RequestOptions) (tc.GetTenantsResponseV5, toclientlib.ReqInf, error) {
	var data tc.GetTenantsResponseV5
	reqInf, err := to.get(apiTenants, opts, &data)
	return data, reqInf, err
}

// CreateTenant creates the Tenant it's passed.
func (to *Session) CreateTenant(t tc.TenantV5, opts RequestOptions) (tc.TenantResponseV5, toclientlib.ReqInf, error) {
	if (t.ParentID == nil || *t.ParentID == 0) && (t.ParentName != nil || *t.ParentName != "") {
		parentOpts := NewRequestOptions()
		parentOpts.QueryParameters.Set("name", *t.ParentName)
		tenant, reqInf, err := to.GetTenants(parentOpts)
		if err != nil {
			return tc.TenantResponseV5{Alerts: tenant.Alerts}, reqInf, err
		}
		if len(tenant.Response) < 1 {
			return tc.TenantResponseV5{Alerts: tenant.Alerts}, reqInf, fmt.Errorf("no Tenant could be found for Parent Tenant '%s'", *t.ParentName)
		}
		t.ParentID = tenant.Response[0].ID
	}

	var data tc.TenantResponseV5
	reqInf, err := to.post(apiTenants, opts, t, &data)
	return data, reqInf, err
}

// UpdateTenant replaces the Tenant identified by 'id' with the one provided.
func (to *Session) UpdateTenant(id int, t tc.TenantV5, opts RequestOptions) (tc.TenantResponseV5, toclientlib.ReqInf, error) {
	var data tc.TenantResponseV5
	reqInf, err := to.put(fmt.Sprintf(apiTenantID, id), opts, t, &data)
	return data, reqInf, err
}

// DeleteTenant deletes the Tenant matching the ID it's passed.
func (to *Session) DeleteTenant(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var data tc.Alerts
	reqInf, err := to.del(fmt.Sprintf(apiTenantID, id), opts, &data)
	return data, reqInf, err
}
