package v13

import (
	"github.com/apache/trafficcontrol/lib/go-tc/v13"
)

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

func (to *Session) GetDomains() ([]v13.Domain, ReqInf, error) {
	var data v13.DomainsResponse
	inf, err := get(to, "/api/1.3/cdns/domains", &data)
	if err != nil {
		return nil, inf, err
	}
	return data.Response, inf, nil
}
