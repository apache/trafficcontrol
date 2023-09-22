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
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// Creates a new Content Invalidation Job
func (to *Session) CreateInvalidationJob(job tc.InvalidationJobInput) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(`/jobs`, job, nil, &alerts)
	return alerts, reqInf, err
}

// Deletes a Content Invalidation Job
func (to *Session) DeleteInvalidationJob(jobID uint64) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.del(fmt.Sprintf("/jobs?id=%d", jobID), nil, &alerts)
	return alerts, reqInf, err

}

// Updates a Content Invalidation Job
func (to *Session) UpdateInvalidationJob(job tc.InvalidationJob) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.put(fmt.Sprintf(`/jobs?id=%d`, *job.ID), job, nil, &alerts)
	return alerts, reqInf, err
}

// GetJobs returns a list of Jobs.
// If deliveryServiceID or userID are not nil, only jobs for that delivery service or belonging to
// that user are returned. Both deliveryServiceID and userID may be nil.
//
// Deprecated, use GetInvalidationJobs instead
func (to *Session) GetJobs(deliveryServiceID *int, userID *int) ([]tc.Job, toclientlib.ReqInf, error) {
	params := url.Values{}
	if deliveryServiceID != nil {
		params.Add("dsId", strconv.Itoa(*deliveryServiceID))
	}
	if userID != nil {
		params.Add("userId", strconv.Itoa(*userID))
	}
	path := "/jobs?" + params.Encode()
	data := struct {
		Response []tc.Job `json:"response"`
	}{}
	reqInf, err := to.get(path, nil, &data)
	return data.Response, reqInf, err
}

// GetInvalidationJobs is deprecated, use GetInvalidationJobsWithHdr instead.
func (to *Session) GetInvalidationJobs(ds *interface{}, user *interface{}) ([]tc.InvalidationJob, toclientlib.ReqInf, error) {
	return to.GetInvalidationJobsWithHdr(ds, user, nil)
}

// GetInvalidationJobsWithHdr returns a list of Content Invalidation Jobs visible to your Tenant,
// filtered according to ds and user.
//
// Either or both of ds and user may be nil, but when they are not they cause filtering of the
// returned jobs by Delivery Service and Traffic Ops user, respectively.
//
// ds may be a uint, int, or float64 indicating the integral, unique identifier of the desired
// Delivery Service (in the case of a float64 the fractional part is dropped, e.g. 3.45 -> 3), or
// it may be a string, in which case it should be the xml_id of the desired Delivery Service, or it
// may be an actual tc.DeliveryService or tc.DeliveryServiceNullable structure.
//
// Likewise, user may be a uint, int or float64 indicating the integral, unique identifier of the
// desired user (in the case of a float64 the fractional part is dropped, e.g. 3.45 -> 3), or it may
// be a string, in which case it should be the username of the desired user, or it may be an actual
// tc.User or tc.UserCurrent structure.
func (to *Session) GetInvalidationJobsWithHdr(ds *interface{}, user *interface{}, hdr http.Header) ([]tc.InvalidationJob, toclientlib.ReqInf, error) {
	const DSIDKey = "dsId"
	const DSKey = "deliveryService"
	const UserKey = "userId"
	const CreatedKey = "createdBy"

	params := url.Values{}
	if ds != nil {
		d := *ds
		switch t := d.(type) {
		case uint:
			params.Add(DSIDKey, strconv.FormatUint(uint64(d.(uint)), 10))
		case float64:
			params.Add(DSIDKey, strconv.FormatInt(int64(d.(float64)), 10))
		case int:
			params.Add(DSIDKey, strconv.FormatInt(int64(d.(int)), 10))
		case string:
			params.Add(DSKey, d.(string))
		case tc.DeliveryServiceNullable:
			if d.(tc.DeliveryServiceNullable).XMLID != nil {
				params.Add(DSKey, *d.(tc.DeliveryServiceNullable).XMLID)
			} else if d.(tc.DeliveryServiceNullable).ID != nil {
				params.Add(DSIDKey, strconv.FormatInt(int64(*d.(tc.DeliveryServiceNullable).ID), 10))
			} else {
				return nil, toclientlib.ReqInf{}, errors.New("no non-nil identifier on passed Delivery Service")
			}
		default:
			return nil, toclientlib.ReqInf{}, fmt.Errorf("invalid type for argument 'ds': %T*", t)
		}
	}
	if user != nil {
		u := *user
		switch t := u.(type) {
		case uint:
			params.Add(UserKey, strconv.FormatUint(uint64(u.(uint)), 10))
		case float64:
			params.Add(UserKey, strconv.FormatInt(int64(u.(float64)), 10))
		case int:
			params.Add(UserKey, strconv.FormatInt(u.(int64), 10))
		case string:
			params.Add(CreatedKey, u.(string))
		case tc.User:
			if u.(tc.User).Username != nil {
				params.Add(CreatedKey, *u.(tc.User).Username)
			} else if u.(tc.User).ID != nil {
				params.Add(UserKey, strconv.FormatInt(int64(*u.(tc.User).ID), 10))
			} else {
				return nil, toclientlib.ReqInf{}, errors.New("no non-nil identifier on passed User")
			}
		case tc.UserCurrent:
			if u.(tc.UserCurrent).UserName != nil {
				params.Add(CreatedKey, *u.(tc.UserCurrent).UserName)
			} else if u.(tc.UserCurrent).ID != nil {
				params.Add(UserKey, strconv.FormatInt(int64(*u.(tc.UserCurrent).ID), 10))
			} else {
				return nil, toclientlib.ReqInf{}, errors.New("no non-nil identifier on passed UserCurrent")
			}
		default:
			return nil, toclientlib.ReqInf{}, fmt.Errorf("invalid type for argument 'user': %T*", t)
		}
	}
	path := "/jobs"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	data := struct {
		Response []tc.InvalidationJob `json:"response"`
	}{}
	reqInf, err := to.get(path, hdr, &data)
	return data.Response, reqInf, err
}
