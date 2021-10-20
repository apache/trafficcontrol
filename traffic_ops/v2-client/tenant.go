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
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	tc "github.com/apache/trafficcontrol/v6/lib/go-tc"
)

const API_TENANTS = apiBase + "/tenants"
const API_TENANT_ID = API_TENANTS + "/%v"

// Tenants gets an array of Tenants.
func (to *Session) Tenants() ([]tc.Tenant, ReqInf, error) {
	var data tc.GetTenantsResponse
	reqInf, err := get(to, API_TENANTS, &data)
	if err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// Tenant gets the Tenant identified by the passed integral, unique identifer - which
// must be passed as a string.
func (to *Session) Tenant(id string) (*tc.Tenant, ReqInf, error) {
	var data tc.GetTenantsResponse
	reqInf, err := get(to, fmt.Sprintf("%s?id=%v", API_TENANTS, id), &data)
	if err != nil {
		return nil, reqInf, err
	}

	return &data.Response[0], reqInf, nil
}

// TenantByName gets the Tenant with the name it's passed.
func (to *Session) TenantByName(name string) (*tc.Tenant, ReqInf, error) {
	var data tc.GetTenantsResponse
	query := API_TENANTS + "?name=" + url.QueryEscape(name)
	reqInf, err := get(to, query, &data)
	if err != nil {
		return nil, reqInf, err
	}

	var ten *tc.Tenant
	if len(data.Response) > 0 {
		ten = &data.Response[0]
	} else {
		err = errors.New("no tenant found with name " + name)
	}
	return ten, reqInf, err
}

// CreateTenant creates the Tenant it's passed.
func (to *Session) CreateTenant(t *tc.Tenant) (*tc.TenantResponse, error) {
	if t.ParentID == 0 && t.ParentName != "" {
		tenant, _, err := to.TenantByName(t.ParentName)
		if err != nil {
			return nil, err
		}
		if tenant == nil {
			return nil, errors.New("no tenant with name " + t.ParentName)
		}
		t.ParentID = tenant.ID
	}

	var data tc.TenantResponse
	jsonReq, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	_, err = post(to, API_TENANTS, jsonReq, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// UpdateTenant updates the Tenant matching the ID it's passed with
// the Tenant it is passed.
func (to *Session) UpdateTenant(id string, t *tc.Tenant) (*tc.TenantResponse, error) {
	var data tc.TenantResponse
	jsonReq, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	_, err = put(to, fmt.Sprintf(API_TENANT_ID, id), jsonReq, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// DeleteTenant deletes the Tenant matching the ID it's passed.
func (to *Session) DeleteTenant(id string) (*tc.DeleteTenantResponse, error) {
	var data tc.DeleteTenantResponse
	_, err := del(to, fmt.Sprintf(API_TENANT_ID, id), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
