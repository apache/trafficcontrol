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
	"net"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

// CreateJob creates a Job.
func (to *Session) CreateJob(job tc.JobRequest) (tc.Alerts, ReqInf, error) {
	remoteAddr := (net.Addr)(nil)
	reqBody, err := json.Marshal(job)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, apiBase+`/user/current/jobs`, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	alerts := tc.Alerts{}
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, err
}

// GetJobs returns a list of Jobs.
// If deliveryServiceID or userID are not nil, only jobs for that delivery service or belonging to that user are returned. Both deliveryServiceID and userID may be nil.
func (to *Session) GetJobs(deliveryServiceID *int, userID *int) ([]tc.Job, ReqInf, error) {
	path := apiBase + "/jobs"
	if deliveryServiceID != nil || userID != nil {
		path += "?"
		if deliveryServiceID != nil {
			path += "dsId=" + strconv.Itoa(*deliveryServiceID)
			if userID != nil {
				path += "&"
			}
		}
		if userID != nil {
			path += "userId=" + strconv.Itoa(*userID)
		}
	}

	resp, remoteAddr, err := to.request(http.MethodGet, path, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	data := struct {
		Response []tc.Job `json:"response"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, err
}
