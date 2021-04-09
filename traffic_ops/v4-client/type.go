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
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	// APITypes is the API version-relative path to the /types API endpoint.
	APITypes = "/types"
)

// CreateType creates the given Type. There should be a very good reason for doing this.
func (to *Session) CreateType(typ tc.Type) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APITypes, typ, nil, &alerts)
	return alerts, reqInf, err
}

// UpdateType replaces the Type identified by 'id' with the one provided.
func (to *Session) UpdateType(id int, typ tc.Type, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APITypes, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, typ, header, &alerts)
	return alerts, reqInf, err
}

// GetTypes returns a list of Types, with an http header and 'useInTable' parameters.
// If a 'useInTable' parameter is passed, the returned Types are restricted to those with
// that exact 'useInTable' property. Only exactly 1 or exactly 0 'useInTable' parameters may
// be passed; passing more will result in an error being returned.
func (to *Session) GetTypes(header http.Header, useInTable ...string) ([]tc.Type, toclientlib.ReqInf, error) {
	if len(useInTable) > 1 {
		return nil, toclientlib.ReqInf{}, errors.New("please pass in a single value for the 'useInTable' parameter")
	}
	var data tc.TypesResponse
	reqInf, err := to.get(APITypes, header, &data)
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

// GetTypeByID retrieves the Type with the given ID.
func (to *Session) GetTypeByID(id int, header http.Header) ([]tc.Type, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APITypes, id)
	var data tc.TypesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetTypeByName retrieves the Type with the given Name.
func (to *Session) GetTypeByName(name string, header http.Header) ([]tc.Type, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", APITypes, name)
	var data tc.TypesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// DeleteType deletes the Type with the given ID.
func (to *Session) DeleteType(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APITypes, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
