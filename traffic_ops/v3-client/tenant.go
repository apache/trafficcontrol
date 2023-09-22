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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// API_TENANTS is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
const API_TENANTS = apiBase + "/tenants"

// API_TENANT_ID is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
const API_TENANT_ID = API_TENANTS + "/%s"

const APITenants = "/tenants"

const APITenantID = APITenants + "/%v"

func (to *Session) TenantsWithHdr(header http.Header) ([]tc.Tenant, toclientlib.ReqInf, error) {
	var data tc.GetTenantsResponse
	reqInf, err := to.get(APITenants, header, &data)
	return data.Response, reqInf, err
}

// Tenants gets an array of Tenants.
// Deprecated: Tenants will be removed in 6.0. Use TenantsWithHdr.
func (to *Session) Tenants() ([]tc.Tenant, toclientlib.ReqInf, error) {
	return to.TenantsWithHdr(nil)
}

func (to *Session) TenantWithHdr(id string, header http.Header) (*tc.Tenant, toclientlib.ReqInf, error) {
	var data tc.GetTenantsResponse
	reqInf, err := to.get(fmt.Sprintf("%s?id=%v", APITenants, id), header, &data)
	if err != nil {
		return nil, reqInf, err
	}
	if reqInf.StatusCode == http.StatusNotModified {
		return nil, reqInf, nil
	}
	return &data.Response[0], reqInf, nil
}

// Tenant gets the Tenant identified by the passed integral, unique identifer - which
// must be passed as a string.
// Deprecated: Tenant will be removed in 6.0. Use TenantWithHdr.
func (to *Session) Tenant(id string) (*tc.Tenant, toclientlib.ReqInf, error) {
	return to.TenantWithHdr(id, nil)
}

func (to *Session) TenantByNameWithHdr(name string, header http.Header) (*tc.Tenant, toclientlib.ReqInf, error) {
	var data tc.GetTenantsResponse
	query := APITenants + "?name=" + url.QueryEscape(name)
	reqInf, err := to.get(query, header, &data)
	if err != nil {
		return nil, reqInf, err
	}
	if reqInf.StatusCode == http.StatusNotModified {
		return nil, reqInf, nil
	}
	if len(data.Response) > 0 {
		return &data.Response[0], reqInf, nil
	} else {
		return nil, reqInf, errors.New("no tenant found with name " + name)
	}
}

// TenantByName gets the Tenant with the name it's passed.
// Deprecated: TenantByName will be removed in 6.0. Use TenantByNameWithHdr.
func (to *Session) TenantByName(name string) (*tc.Tenant, toclientlib.ReqInf, error) {
	return to.TenantByNameWithHdr(name, nil)
}

// CreateTenant creates the Tenant it's passed.
func (to *Session) CreateTenant(t *tc.Tenant) (*tc.TenantResponse, error) {
	if t.ParentID == 0 && t.ParentName != "" {
		tenant, _, err := to.TenantByNameWithHdr(t.ParentName, nil)
		if err != nil {
			return nil, err
		}
		if tenant == nil {
			return nil, errors.New("no tenant with name " + t.ParentName)
		}
		t.ParentID = tenant.ID
	}

	var data tc.TenantResponse
	_, err := to.post(APITenants, t, nil, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (to *Session) UpdateTenantWithHdr(id string, t *tc.Tenant, header http.Header) (*tc.TenantResponse, toclientlib.ReqInf, error) {
	var data tc.TenantResponse
	reqInf, err := to.put(fmt.Sprintf(APITenantID, id), t, header, &data)
	if err != nil {
		return nil, reqInf, err
	}
	return &data, reqInf, nil
}

// UpdateTenant updates the Tenant matching the ID it's passed with
// the Tenant it is passed.
// Deprecated: UpdateTenant will be removed in 6.0. Use UpdateTenantWithHdr.
func (to *Session) UpdateTenant(id string, t *tc.Tenant) (*tc.TenantResponse, error) {
	data, _, err := to.UpdateTenantWithHdr(id, t, nil)
	return data, err
}

// DeleteTenant deletes the Tenant matching the ID it's passed.
func (to *Session) DeleteTenant(id string) (*tc.DeleteTenantResponse, error) {
	var data tc.DeleteTenantResponse
	_, err := to.del(fmt.Sprintf(APITenantID, id), nil, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
