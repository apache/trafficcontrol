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
import "encoding/json"
import "fmt"
import "net/http"
import "net/url"
import "strconv"

import "github.com/apache/trafficcontrol/v6/lib/go-tc"

func (to *Session) getFederationResolvers(id *uint, ip *string, t *string) ([]tc.FederationResolver, ReqInf, error) {
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

	var path = apiBase + "/federation_resolvers"
	if len(vals) > 0 {
		path = fmt.Sprintf("%s?%s", path, vals.Encode())
	}

	var data struct {
		Response []tc.FederationResolver `json:"response"`
	}
	inf, err := get(to, path, &data)
	return data.Response, inf, err
}

// GetFederationResolvers retrieves all Federation Resolvers from Traffic Ops
func (to *Session) GetFederationResolvers() ([]tc.FederationResolver, ReqInf, error) {
	return to.getFederationResolvers(nil, nil, nil)
}

// GetFederationResolverByID retrieves a single Federation Resolver identified by ID.
func (to *Session) GetFederationResolverByID(ID uint) (tc.FederationResolver, ReqInf, error) {
	var fr tc.FederationResolver
	frs, inf, err := to.getFederationResolvers(&ID, nil, nil)
	if len(frs) > 0 {
		fr = frs[0]
	}
	return fr, inf, err
}

// GetFederationResolverByIPAddress retrieves the Federation Resolver that uses the IP address or
// CIDR-notation subnet 'ip'.
func (to *Session) GetFederationResolverByIPAddress(ip string) (tc.FederationResolver, ReqInf, error) {
	var fr tc.FederationResolver
	frs, inf, err := to.getFederationResolvers(nil, &ip, nil)
	if len(frs) > 0 {
		fr = frs[0]
	}
	return fr, inf, err
}

// GetFederationResolversByType gets all Federation Resolvers that are of the Type named 't'.
func (to *Session) GetFederationResolversByType(t string) ([]tc.FederationResolver, ReqInf, error) {
	return to.getFederationResolvers(nil, nil, &t)
}

// CreateFederationResolver creates the Federation Resolver 'fr'.
func (to *Session) CreateFederationResolver(fr tc.FederationResolver) (tc.Alerts, ReqInf, error) {
	var reqInf = ReqInf{CacheHitStatus: CacheHitStatusMiss}
	var alerts tc.Alerts

	req, err := json.Marshal(fr)
	if err != nil {
		return alerts, reqInf, err
	}

	var resp *http.Response
	resp, reqInf.RemoteAddr, err = to.request(http.MethodPost, apiBase+"/federation_resolvers", req)
	if err != nil {
		return alerts, reqInf, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, err
}

// DeleteFederationResolver deletes the Federation Resolver identified by 'id'.
func (to *Session) DeleteFederationResolver(id uint) (tc.Alerts, ReqInf, error) {
	var reqInf = ReqInf{CacheHitStatus: CacheHitStatusMiss}
	var alerts tc.Alerts

	var path = fmt.Sprintf("%s/federation_resolvers?id=%d", apiBase, id)
	var resp *http.Response
	var err error
	resp, reqInf.RemoteAddr, err = to.request(http.MethodDelete, path, nil)
	if err != nil {
		return alerts, reqInf, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, err
}
