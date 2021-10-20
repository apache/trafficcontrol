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
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

func (to *Session) CreateSteeringTarget(st tc.SteeringTargetNullable) (tc.Alerts, ReqInf, error) {
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss}
	if st.DeliveryServiceID == nil {
		return tc.Alerts{}, reqInf, errors.New("missing delivery service id")
	}
	reqBody, err := json.Marshal(st)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}

	resp := (*http.Response)(nil)
	if resp, reqInf.RemoteAddr, err = to.request(http.MethodPost, apiBase+`/steering/`+strconv.Itoa(int(*st.DeliveryServiceID))+`/targets`, reqBody); err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	alerts := tc.Alerts{}
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, err
}

func (to *Session) UpdateSteeringTarget(st tc.SteeringTargetNullable) (tc.Alerts, ReqInf, error) {
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss}
	if st.DeliveryServiceID == nil {
		return tc.Alerts{}, reqInf, errors.New("missing delivery service id")
	}
	if st.TargetID == nil {
		return tc.Alerts{}, reqInf, errors.New("missing target id")
	}

	reqBody, err := json.Marshal(st)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp := (*http.Response)(nil)
	resp, reqInf.RemoteAddr, err = to.request(http.MethodPut, apiBase+`/steering/`+strconv.Itoa(int(*st.DeliveryServiceID))+`/targets/`+strconv.Itoa(int(*st.TargetID)), reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	alerts := tc.Alerts{}
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, err
}

func (to *Session) GetSteeringTargets(dsID int) ([]tc.SteeringTargetNullable, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, apiBase+`/steering/`+strconv.Itoa(dsID)+`/targets`, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	data := struct {
		Response []tc.SteeringTargetNullable `json:"response"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, err
}

func (to *Session) DeleteSteeringTarget(dsID int, targetID int) (tc.Alerts, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodDelete, apiBase+`/steering/`+strconv.Itoa(dsID)+`/targets/`+strconv.Itoa(targetID), nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	alerts := tc.Alerts{}
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, err
}
