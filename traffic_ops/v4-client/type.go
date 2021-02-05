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

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_TYPES = apiBase + "/types"
)

// CreateType creates a Type. There should be a very good reason for doing this.
func (to *Session) CreateType(typ tc.Type) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(API_TYPES, typ, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) UpdateTypeByIDWithHdr(id int, typ tc.Type, header http.Header) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_TYPES, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, typ, header, &alerts)
	return alerts, reqInf, err
}

// UpdateTypeByID updates a Type by ID.
// Deprecated: UpdateTypeByID will be removed in 6.0. Use UpdateTypeByIDWithHdr.
func (to *Session) UpdateTypeByID(id int, typ tc.Type) (tc.Alerts, ReqInf, error) {
	return to.UpdateTypeByIDWithHdr(id, typ, nil)
}

// GetTypesWithHdr returns a list of Types, with an http header and 'useInTable' parameters.
// If a 'useInTable' parameter is passed, the returned Types are restricted to those with
// that exact 'useInTable' property. Only exactly 1 or exactly 0 'useInTable' parameters may
// be passed; passing more will result in an error being returned.
func (to *Session) GetTypesWithHdr(header http.Header, useInTable ...string) ([]tc.Type, ReqInf, error) {
	if len(useInTable) > 1 {
		return nil, ReqInf{}, errors.New("please pass in a single value for the 'useInTable' parameter")
	}
	var data tc.TypesResponse
	reqInf, err := to.get(API_TYPES, header, &data)
	if err != nil {
		return nil, reqInf, err
	}

	var types []tc.Type
	for _, d := range data.Response {
		if useInTable != nil {
			if d.UseInTable == useInTable[0] {
				types = append(types, d)
			}
		} else {
			types = append(types, d)
		}
	}

	return types, reqInf, nil
}

// GetTypes returns a list of Types. If a 'useInTable' parameter is passed, the returned Types
// are restricted to those with that exact 'useInTable' property. Only exactly 1 or exactly 0
// 'useInTable' parameters may be passed; passing more will result in an error being returned.
// Deprecated: GetTypes will be removed in 6.0. Use GetTypesWithHdr.
func (to *Session) GetTypes(useInTable ...string) ([]tc.Type, ReqInf, error) {
	return to.GetTypesWithHdr(nil, useInTable...)
}

// GetTypeByID GETs a Type by the Type ID, and filters by http header params in the request.
func (to *Session) GetTypeByIDWithHdr(id int, header http.Header) ([]tc.Type, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_TYPES, id)
	var data tc.TypesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetTypeByID GETs a Type by the Type ID.
// Deprecated: GetTypeByID will be removed in 6.0. Use GetTypeByIDWithHdr.
func (to *Session) GetTypeByID(id int) ([]tc.Type, ReqInf, error) {
	return to.GetTypeByIDWithHdr(id, nil)
}

func (to *Session) GetTypeByNameWithHdr(name string, header http.Header) ([]tc.Type, ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", API_TYPES, name)
	var data tc.TypesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetTypeByName GETs a Type by the Type name.
// Deprecated: GetTypeByName will be removed in 6.0. Use GetTypeByNameWithHdr.
func (to *Session) GetTypeByName(name string) ([]tc.Type, ReqInf, error) {
	return to.GetTypeByNameWithHdr(name, nil)
}

// DeleteTypeByID DELETEs a Type by ID.
func (to *Session) DeleteTypeByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_TYPES, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
