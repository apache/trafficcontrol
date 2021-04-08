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
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

// APITenants is the API version-relative path to the /tenants API endpoint.
const APITenants = "/tenants"

// APITenantID is the API version-relative path to the /tenants/{{ID}} API endpoint.
const APITenantID = APITenants + "/%d"

// GetTenants retrieves all Tenants stored in Traffic Ops.
func (to *Session) GetTenants(header http.Header) ([]tc.Tenant, toclientlib.ReqInf, error) {
	var data tc.GetTenantsResponse
	reqInf, err := to.get(APITenants, header, &data)
	return data.Response, reqInf, err
}

// GetTenantByID retrieves the Tenant with the given ID.
func (to *Session) GetTenantByID(id int, header http.Header) (tc.Tenant, toclientlib.ReqInf, error) {
	var data tc.GetTenantsResponse
	reqInf, err := to.get(fmt.Sprintf("%s?id=%v", APITenants, id), header, &data)
	if err != nil {
		return tc.Tenant{}, reqInf, err
	}
	if reqInf.StatusCode == http.StatusNotModified {
		return tc.Tenant{}, reqInf, nil
	}
	return data.Response[0], reqInf, nil
}

// GetTenantByName retrieves the Tenant with the given Name
func (to *Session) GetTenantByName(name string, header http.Header) (tc.Tenant, toclientlib.ReqInf, error) {
	var data tc.GetTenantsResponse
	query := APITenants + "?name=" + url.QueryEscape(name)
	reqInf, err := to.get(query, header, &data)
	if err != nil {
		return tc.Tenant{}, reqInf, err
	}
	if reqInf.StatusCode == http.StatusNotModified {
		return tc.Tenant{}, reqInf, nil
	}
	if len(data.Response) > 0 {
		return data.Response[0], reqInf, nil
	}
	return tc.Tenant{}, reqInf, errors.New("no tenant found with name " + name)
}

// CreateTenant creates the Tenant it's passed.
func (to *Session) CreateTenant(t tc.Tenant) (tc.TenantResponse, error) {
	if t.ParentID == 0 && t.ParentName != "" {
		tenant, _, err := to.GetTenantByName(t.ParentName, nil)
		if err != nil {
			return tc.TenantResponse{}, err
		}
		t.ParentID = tenant.ID
	}

	var data tc.TenantResponse
	_, err := to.post(APITenants, t, nil, &data)
	return data, err
}

// UpdateTenant replaces the Tenant identified by 'id' with the one provided.
func (to *Session) UpdateTenant(id int, t tc.Tenant, header http.Header) (tc.TenantResponse, toclientlib.ReqInf, error) {
	var data tc.TenantResponse
	reqInf, err := to.put(fmt.Sprintf(APITenantID, id), t, header, &data)
	return data, reqInf, err
}

// DeleteTenant deletes the Tenant matching the ID it's passed.
func (to *Session) DeleteTenant(id int) (tc.DeleteTenantResponse, error) {
	var data tc.DeleteTenantResponse
	_, err := to.del(fmt.Sprintf(APITenantID, id), nil, &data)
	return data, err
}
