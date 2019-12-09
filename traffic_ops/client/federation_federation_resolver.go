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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

// GetFederationsFederationResolversByID retrieves all Federation Resolvers belonging to Federation of ID.
func (to *Session) GetFederationsFederationResolversByID(id int) (tc.FederationFederationResolversResponse, ReqInf, error) {
	path := fmt.Sprintf("%s/federations/%d/federation_resolvers", apiBase, id)

	data := tc.FederationFederationResolversResponse{}
	inf, err := get(to, path, &data)
	return data, inf, err
}

// AssignFederationsFederationResolver creates the Federation Resolver 'fr'.
func (to *Session) AssignFederationsFederationResolver(fedID int, ids []int, replace bool) (tc.Alerts, ReqInf, error) {
	var (
		req = tc.AssignFederationResolversRequest{
			Replace:        replace,
			FedResolverIDs: ids,
		}
		reqInf = ReqInf{CacheHitStatus: CacheHitStatusMiss}
		resp   = tc.FederationFederationResolversResponse
	)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return resp, reqInf, err
	}

	path := fmt.Sprintf("%s/federations/%d/federation_resolvers", apiBase, fedID) //, req)
	httpResp, remoteAddr, err := to.request(http.MethodPost, path, reqBody)
	if err != nil {
		return resp, reqInf, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(httpResp.Body).Decode(&resp)

	return resp, reqInf, err
}
