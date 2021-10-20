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

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

// GetFederationFederationResolversByID retrieves all Federation Resolvers belonging to Federation of ID.
func (to *Session) GetFederationFederationResolversByID(id int) (tc.FederationFederationResolversResponse, ReqInf, error) {
	var (
		path   = fmt.Sprintf("%s/federations/%d/federation_resolvers", apiBase, id)
		reqInf = ReqInf{CacheHitStatus: CacheHitStatusMiss}
		resp   tc.FederationFederationResolversResponse
	)
	httpResp, remoteAddr, err := to.request(http.MethodGet, path, nil)
	reqInf.RemoteAddr = remoteAddr
	if err != nil {
		return resp, reqInf, err
	}
	defer httpResp.Body.Close()

	err = json.NewDecoder(httpResp.Body).Decode(&resp)

	return resp, reqInf, err
}

// AssignFederationFederationResolver creates the Federation Resolver 'fr'.
func (to *Session) AssignFederationFederationResolver(fedID int, resolverIDs []int, replace bool) (tc.AssignFederationFederationResolversResponse, ReqInf, error) {
	var (
		path = fmt.Sprintf("%s/federations/%d/federation_resolvers", apiBase, fedID)
		req  = tc.AssignFederationResolversRequest{
			Replace:        replace,
			FedResolverIDs: resolverIDs,
		}
		reqInf = ReqInf{CacheHitStatus: CacheHitStatusMiss}
		resp   tc.AssignFederationFederationResolversResponse
	)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return resp, reqInf, err
	}

	httpResp, remoteAddr, err := to.request(http.MethodPost, path, reqBody)
	reqInf.RemoteAddr = remoteAddr
	if err != nil {
		return resp, reqInf, err
	}
	defer httpResp.Body.Close()

	err = json.NewDecoder(httpResp.Body).Decode(&resp)

	return resp, reqInf, err
}
