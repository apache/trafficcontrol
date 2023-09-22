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

// GetFederationFederationResolversByID retrieves all Federation Resolvers belonging to Federation of ID.
func (to *Session) GetFederationFederationResolversByID(id int) (tc.FederationFederationResolversResponse, toclientlib.ReqInf, error) {
	path := fmt.Sprintf("/federations/%d/federation_resolvers", id)
	resp := tc.FederationFederationResolversResponse{}
	reqInf, err := to.get(path, nil, &resp)
	return resp, reqInf, err
}

// AssignFederationFederationResolver creates the Federation Resolver 'fr'.
func (to *Session) AssignFederationFederationResolver(fedID int, resolverIDs []int, replace bool) (tc.AssignFederationFederationResolversResponse, toclientlib.ReqInf, error) {
	path := fmt.Sprintf("/federations/%d/federation_resolvers", fedID)
	req := tc.AssignFederationResolversRequest{
		Replace:        replace,
		FedResolverIDs: resolverIDs,
	}
	resp := tc.AssignFederationFederationResolversResponse{}
	reqInf, err := to.post(path, req, nil, &resp)
	return resp, reqInf, err
}
