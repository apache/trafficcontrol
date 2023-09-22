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

package client

import (
	"net/http"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func (to *Session) SteeringWithHdr(header http.Header) ([]tc.Steering, toclientlib.ReqInf, error) {
	data := struct {
		Response []tc.Steering `json:"response"`
	}{}
	reqInf, err := to.get(`/steering`, header, &data)
	return data.Response, reqInf, err
}

// Deprecated: Steering will be removed in 6.0. Use SteeringWithHdr.
func (to *Session) Steering() ([]tc.Steering, toclientlib.ReqInf, error) {
	return to.SteeringWithHdr(nil)
}
