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
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

// GetFederationResolvers retrieves Federation Resolvers from Traffic Ops.
func (to *Session) GetFederationResolvers(params url.Values, header http.Header) ([]tc.FederationResolver, toclientlib.ReqInf, error) {
	var path = "federation_resolvers/"
	if len(params) > 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}

	var data struct {
		Response []tc.FederationResolver `json:"response"`
	}
	inf, err := to.get(path, header, &data)
	return data.Response, inf, err
}

// GetFederationResolverByID retrieves a single Federation Resolver identified by ID.
func (to *Session) GetFederationResolverByID(ID uint, header http.Header) (tc.FederationResolver, toclientlib.ReqInf, error) {
	vals := url.Values{}
	vals.Set("id", strconv.FormatUint(uint64(ID), 10))
	feds, reqInf, err := to.GetFederationResolvers(vals, nil)
	if err != nil {
		return tc.FederationResolver{}, reqInf, err
	}
	if len(feds) != 1 {
		return tc.FederationResolver{}, reqInf, fmt.Errorf("Traffic Ops returned %d Federation Resolvers with ID '%d'", len(feds), ID)
	}
	return feds[0], reqInf, nil
}

// CreateFederationResolver creates the Federation Resolver 'fr'.
func (to *Session) CreateFederationResolver(fr tc.FederationResolver, header http.Header) (tc.FederationResolverResponseV4, toclientlib.ReqInf, error) {
	var response tc.FederationResolverResponseV4
	reqInf, err := to.post("federation_resolvers/", fr, header, &response)
	return response, reqInf, err
}

// DeleteFederationResolver deletes the Federation Resolver identified by 'id'.
func (to *Session) DeleteFederationResolver(id uint) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	path := fmt.Sprintf("/federation_resolvers?id=%d", id)
	reqInf, err := to.del(path, nil, &alerts)
	return alerts, reqInf, err
}
