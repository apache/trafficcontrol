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

	"github.com/apache/trafficcontrol/lib/go-tc"
)

// GetFederationsFederationResolversByID retrieves all Federation Resolvers belonging to Federation of ID.
func (to *Session) GetFederationsFederationResolversByID(id *int) (tc.FederationFederationResolversResponse, ReqInf, error) {
	ffrs, inf, err := to.getFederationsFederationResolvers(id)
	return ffrs, inf, err
}

// AssignFederationsFederationResolver creates the Federation Resolver 'fr'.
func (to *Session) AssignFederationsFederationResolver(fedID int, ids []int) (tc.Alerts, ReqInf, error) {
	/*var reqInf = ReqInf{CacheHitStatus: CacheHitStatusMiss}
	var alerts tc.Alerts

	req, err := json.Marshal(ffr)
	if err != nil {
		return alerts, reqInf, err
	}

	var resp *http.Response
	resp, reqInf.RemoteAddr, err = to.request(http.MethodPost, fmt.Sprintf("%s/federations/%d/federation_resolvers", apiBase, ffr.FederationID), req)
	if err != nil {
		return alerts, reqInf, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, err
	*/
	return tc.Alerts{}, ReqInf{}, nil
}

func (to *Session) getFederationsFederationResolvers(id *int) (tc.FederationFederationResolversResponse, ReqInf, error) {
	var path = fmt.Sprintf("%s/federations/%d/federation_resolvers", apiBase, id)

	data := tc.FederationFederationResolversResponse{}
	inf, err := get(to, path, &data)
	return data, inf, err
}
