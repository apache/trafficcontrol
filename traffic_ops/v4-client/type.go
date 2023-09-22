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

// apiTypes is the API version-relative path to the /types API endpoint.
const apiTypes = "/types"

// CreateType creates the given Type. There should be a very good reason for doing this.
func (to *Session) CreateType(typ tc.Type, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiTypes, opts, typ, &alerts)
	return alerts, reqInf, err
}

// UpdateType replaces the Type identified by 'id' with the one provided.
func (to *Session) UpdateType(id int, typ tc.Type, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", apiTypes, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, opts, typ, &alerts)
	return alerts, reqInf, err
}

// GetTypes returns a list of Types, with an http header and 'useInTable' parameters.
// If a 'useInTable' parameter is passed, the returned Types are restricted to those with
// that exact 'useInTable' property. Only exactly 1 or exactly 0 'useInTable' parameters may
// be passed; passing more will result in an error being returned.
func (to *Session) GetTypes(opts RequestOptions) (tc.TypesResponse, toclientlib.ReqInf, error) {
	var data tc.TypesResponse
	reqInf, err := to.get(apiTypes, opts, &data)
	return data, reqInf, err
}

// DeleteType deletes the Type with the given ID.
func (to *Session) DeleteType(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", apiTypes, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, opts, &alerts)
	return alerts, reqInf, err
}
