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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func (to *Session) CreateSteeringTarget(st tc.SteeringTargetNullable) (tc.Alerts, toclientlib.ReqInf, error) {
	if st.DeliveryServiceID == nil {
		return tc.Alerts{}, toclientlib.ReqInf{CacheHitStatus: toclientlib.CacheHitStatusMiss}, errors.New("missing delivery service id")
	}
	alerts := tc.Alerts{}
	route := fmt.Sprintf("/steering/%d/targets", *st.DeliveryServiceID)
	reqInf, err := to.post(route, st, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) UpdateSteeringTargetWithHdr(st tc.SteeringTargetNullable, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	reqInf := toclientlib.ReqInf{CacheHitStatus: toclientlib.CacheHitStatusMiss}
	if st.DeliveryServiceID == nil {
		return tc.Alerts{}, reqInf, errors.New("missing delivery service id")
	}
	if st.TargetID == nil {
		return tc.Alerts{}, reqInf, errors.New("missing target id")
	}
	route := fmt.Sprintf("/steering/%d/targets/%d", *st.DeliveryServiceID, *st.TargetID)
	alerts := tc.Alerts{}
	reqInf, err := to.put(route, st, header, &alerts)
	return alerts, reqInf, err
}

// Deprecated: UpdateSteeringTarget will be removed in 6.0. Use UpdateSteeringTargetWithHdr.
func (to *Session) UpdateSteeringTarget(st tc.SteeringTargetNullable) (tc.Alerts, toclientlib.ReqInf, error) {
	return to.UpdateSteeringTargetWithHdr(st, nil)
}

func (to *Session) GetSteeringTargetsWithHdr(dsID int, header http.Header) ([]tc.SteeringTargetNullable, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("/steering/%d/targets", dsID)
	data := struct {
		Response []tc.SteeringTargetNullable `json:"response"`
	}{}
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

func (to *Session) GetSteeringTargets(dsID int) ([]tc.SteeringTargetNullable, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("/steering/%d/targets", dsID)
	data := struct {
		Response []tc.SteeringTargetNullable `json:"response"`
	}{}
	reqInf, err := to.get(route, nil, &data)
	return data.Response, reqInf, err
}

func (to *Session) DeleteSteeringTarget(dsID int, targetID int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("/steering/%d/targets/%d", dsID, targetID)
	alerts := tc.Alerts{}
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
