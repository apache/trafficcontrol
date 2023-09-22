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
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiASNs is the API version-relative path for the /asns API endpoint.
const apiASNs = "/asns"

// CreateASN creates the passed ASN.
func (to *Session) CreateASN(asn tc.ASN, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiASNs, opts, asn, &alerts)
	return alerts, reqInf, err
}

// UpdateASN updates the ASN identified by id by replacing it with the passed
// ASN.
func (to *Session) UpdateASN(id int, entity tc.ASN, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("id", strconv.Itoa(id))
	var alerts tc.Alerts
	reqInf, err := to.put(apiASNs, opts, entity, &alerts)
	return alerts, reqInf, err
}

// GetASNs retrieves ASNs from Traffic Ops.
func (to *Session) GetASNs(opts RequestOptions) (tc.ASNsResponse, toclientlib.ReqInf, error) {
	var data tc.ASNsResponse
	reqInf, err := to.get(apiASNs, opts, &data)
	return data, reqInf, err
}

// DeleteASN deletes the ASN with the given ID.
func (to *Session) DeleteASN(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("id", strconv.Itoa(id))
	var alerts tc.Alerts
	reqInf, err := to.del(apiASNs, opts, &alerts)
	return alerts, reqInf, err
}
