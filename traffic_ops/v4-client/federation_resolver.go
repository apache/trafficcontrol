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

// apiFederationResolvers is the API-version-relative path to the
// /federation_resolvers endpoint.
const apiFederationResolvers = "/federation_resolvers"

// GetFederationResolvers retrieves Federation Resolvers from Traffic Ops.
func (to *Session) GetFederationResolvers(opts RequestOptions) (tc.FederationResolversResponse, toclientlib.ReqInf, error) {
	var data tc.FederationResolversResponse
	inf, err := to.get(apiFederationResolvers, opts, &data)
	return data, inf, err
}

// CreateFederationResolver creates the Federation Resolver 'fr'.
func (to *Session) CreateFederationResolver(fr tc.FederationResolver, opts RequestOptions) (tc.FederationResolverResponse, toclientlib.ReqInf, error) {
	var response tc.FederationResolverResponse
	reqInf, err := to.post(apiFederationResolvers, opts, fr, &response)
	return response, reqInf, err
}

// DeleteFederationResolver deletes the Federation Resolver identified by 'id'.
func (to *Session) DeleteFederationResolver(id uint, opts RequestOptions) (tc.FederationResolverResponse, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("id", strconv.FormatUint(uint64(id), 10))
	var alerts tc.FederationResolverResponse
	reqInf, err := to.del(apiFederationResolvers, opts, &alerts)
	return alerts, reqInf, err
}
