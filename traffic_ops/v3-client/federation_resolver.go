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

import "fmt"
import "net/http"
import "net/url"
import "strconv"

import "github.com/apache/trafficcontrol/v8/lib/go-tc"
import "github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"

func (to *Session) getFederationResolvers(id *uint, ip *string, t *string, header http.Header) ([]tc.FederationResolver, toclientlib.ReqInf, error) {
	var vals = url.Values{}
	if id != nil {
		vals.Set("id", strconv.FormatUint(uint64(*id), 10))
	}
	if ip != nil {
		vals.Set("ipAddress", *ip)
	}
	if t != nil {
		vals.Set("type", *t)
	}

	path := "/federation_resolvers"
	if len(vals) > 0 {
		path = fmt.Sprintf("%s?%s", path, vals.Encode())
	}

	var data struct {
		Response []tc.FederationResolver `json:"response"`
	}
	inf, err := to.get(path, header, &data)
	return data.Response, inf, err
}

func (to *Session) GetFederationResolversWithHdr(header http.Header) ([]tc.FederationResolver, toclientlib.ReqInf, error) {
	return to.getFederationResolvers(nil, nil, nil, header)
}

// GetFederationResolvers retrieves all Federation Resolvers from Traffic Ops
// Deprecated: GetFederationResolvers will be removed in 6.0. Use GetFederationResolversWithHdr.
func (to *Session) GetFederationResolvers() ([]tc.FederationResolver, toclientlib.ReqInf, error) {
	return to.GetFederationResolversWithHdr(nil)
}

func (to *Session) GetFederationResolverByIDWithHdr(ID uint, header http.Header) (tc.FederationResolver, toclientlib.ReqInf, error) {
	var fr tc.FederationResolver
	frs, inf, err := to.getFederationResolvers(&ID, nil, nil, header)
	if len(frs) > 0 {
		fr = frs[0]
	}
	return fr, inf, err
}

// GetFederationResolverByID retrieves a single Federation Resolver identified by ID.
// Deprecated: GetFederationResolverByID will be removed in 6.0. Use GetFederationResolverByIDWithHdr.
func (to *Session) GetFederationResolverByID(ID uint) (tc.FederationResolver, toclientlib.ReqInf, error) {
	return to.GetFederationResolverByIDWithHdr(ID, nil)
}

func (to *Session) GetFederationResolverByIPAddressWithHdr(ip string, header http.Header) (tc.FederationResolver, toclientlib.ReqInf, error) {
	var fr tc.FederationResolver
	frs, inf, err := to.getFederationResolvers(nil, &ip, nil, header)
	if len(frs) > 0 {
		fr = frs[0]
	}
	return fr, inf, err
}

// GetFederationResolverByIPAddress retrieves the Federation Resolver that uses the IP address or
// CIDR-notation subnet 'ip'.
// Deprecated: GetFederationResolverByIPAddress will be removed in 6.0. Use GetFederationResolverByIPAddressWithHdr.
func (to *Session) GetFederationResolverByIPAddress(ip string) (tc.FederationResolver, toclientlib.ReqInf, error) {
	return to.GetFederationResolverByIPAddressWithHdr(ip, nil)
}

func (to *Session) GetFederationResolversByTypeWithHdr(t string, header http.Header) ([]tc.FederationResolver, toclientlib.ReqInf, error) {
	return to.getFederationResolvers(nil, nil, &t, header)
}

// GetFederationResolversByType gets all Federation Resolvers that are of the Type named 't'.
// Deprecated: GetFederationResolversByType will be removed in 6.0. Use GetFederationResolversByTypeWithHdr.
func (to *Session) GetFederationResolversByType(t string) ([]tc.FederationResolver, toclientlib.ReqInf, error) {
	return to.GetFederationResolversByTypeWithHdr(t, nil)
}

// CreateFederationResolver creates the Federation Resolver 'fr'.
func (to *Session) CreateFederationResolver(fr tc.FederationResolver) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post("/federation_resolvers", fr, nil, &alerts)
	return alerts, reqInf, err
}

// DeleteFederationResolver deletes the Federation Resolver identified by 'id'.
func (to *Session) DeleteFederationResolver(id uint) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	path := fmt.Sprintf("/federation_resolvers?id=%d", id)
	reqInf, err := to.del(path, nil, &alerts)
	return alerts, reqInf, err
}
