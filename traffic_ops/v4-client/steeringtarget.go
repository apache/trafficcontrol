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
	"errors"
	"fmt"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func (to *Session) CreateSteeringTarget(st tc.SteeringTargetNullable) (tc.Alerts, ReqInf, error) {
	if st.DeliveryServiceID == nil {
		return tc.Alerts{}, ReqInf{CacheHitStatus: CacheHitStatusMiss}, errors.New("missing delivery service id")
	}
	alerts := tc.Alerts{}
	route := fmt.Sprintf("%s/steering/%d/targets", apiBase, *st.DeliveryServiceID)
	reqInf, err := to.post(route, st, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) UpdateSteeringTargetWithHdr(st tc.SteeringTargetNullable, header http.Header) (tc.Alerts, ReqInf, error) {
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss}
	if st.DeliveryServiceID == nil {
		return tc.Alerts{}, reqInf, errors.New("missing delivery service id")
	}
	if st.TargetID == nil {
		return tc.Alerts{}, reqInf, errors.New("missing target id")
	}
	route := fmt.Sprintf("%s/steering/%d/targets/%d", apiBase, *st.DeliveryServiceID, *st.TargetID)
	alerts := tc.Alerts{}
	reqInf, err := to.put(route, st, header, &alerts)
	return alerts, reqInf, err
}

// Deprecated: UpdateSteeringTarget will be removed in 6.0. Use UpdateSteeringTargetWithHdr.
func (to *Session) UpdateSteeringTarget(st tc.SteeringTargetNullable) (tc.Alerts, ReqInf, error) {
	return to.UpdateSteeringTargetWithHdr(st, nil)
}

func (to *Session) GetSteeringTargets(dsID int) ([]tc.SteeringTargetNullable, ReqInf, error) {
	route := fmt.Sprintf("%s/steering/%d/targets", apiBase, dsID)
	data := struct {
		Response []tc.SteeringTargetNullable `json:"response"`
	}{}
	reqInf, err := to.get(route, nil, &data)
	return data.Response, reqInf, err
}

func (to *Session) DeleteSteeringTarget(dsID int, targetID int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/steering/%d/targets/%d", apiBase, dsID, targetID)
	alerts := tc.Alerts{}
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
