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
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

// Creates a new Content Invalidation Job
func (to *Session) CreateInvalidationJob(job tc.InvalidationJobInput) (tc.Alerts, ReqInf, error) {
	remoteAddr := (net.Addr)(nil)
	reqBody, err := json.Marshal(job)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, apiBase+`/jobs`, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, err
}

// GetJobs returns a list of Jobs.
// If deliveryServiceID or userID are not nil, only jobs for that delivery service or belonging to
// that user are returned. Both deliveryServiceID and userID may be nil.
//
// Deprecated, use GetInvalidationJobs instead
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

// Returns a list of Content Invalidation Jobs visible to your Tenant, filtered according to
// ds and user
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
func (to *Session) GetInvalidationJobs(ds *interface{}, user *interface{}) ([]tc.InvalidationJob, ReqInf, error) {
	path := apiBase + "/jobs"
	if ds != nil || user != nil {
		path += "?"

		if ds != nil {
			d := *ds
			switch t := d.(type) {
			case uint:
				path += "dsId=" + strconv.FormatUint(uint64(d.(uint)), 10)
			case float64:
				path += "dsId=" + strconv.FormatInt(int64(d.(float64)), 10)
			case int:
				path += "dsId=" + strconv.FormatInt(int64(d.(int)), 10)
			case string:
				path += "deliveryService=" + d.(string)
			case tc.DeliveryServiceNullable:
				if d.(tc.DeliveryServiceNullable).XMLID != nil {
					path += "deliveryService=" + *d.(tc.DeliveryServiceNullable).XMLID
				} else if d.(tc.DeliveryServiceNullable).ID != nil {
					path += "dsId=" + strconv.FormatInt(int64(*d.(tc.DeliveryServiceNullable).ID), 10)
				} else {
					return nil, ReqInf{}, errors.New("No non-nil identifier on passed Delivery Service!")
				}
			default:
				return nil, ReqInf{}, fmt.Errorf("Invalid type for argument 'ds': %T*", t)
			}

			if user != nil {
				path += "&"
			}
		}

		if user != nil {
			u := *user
			switch t := u.(type) {
			case uint:
				path += "userId=" + strconv.FormatUint(uint64(u.(uint)), 10)
			case float64:
				path += "userId=" + strconv.FormatInt(int64(u.(float64)), 10)
			case int:
				path += "userId=" + strconv.FormatInt(int64(u.(int64)), 10)
			case string:
				path += "createdBy=" + u.(string)
			case tc.User:
				if u.(tc.User).Username != nil {
					path += "createdBy=" + *u.(tc.User).Username
				} else if u.(tc.User).ID != nil {
					path += "userId=" + strconv.FormatInt(int64(*u.(tc.User).ID), 10)
				} else {
					return nil, ReqInf{}, errors.New("No non-nil identifier on passed User!")
				}
			case tc.UserCurrent:
				if u.(tc.UserCurrent).UserName != nil {
					path += "createdBy=" + *u.(tc.UserCurrent).UserName
				} else if u.(tc.UserCurrent).ID != nil {
					path += "userId=" + strconv.FormatInt(int64(*u.(tc.UserCurrent).ID), 10)
				} else {
					return nil, ReqInf{}, errors.New("No non-nil identifier on passed UserCurrent!")
				}
			default:
				return nil, ReqInf{}, fmt.Errorf("Invalid type for argument 'user': %T*", t)
			}
		}
	}

	resp, remoteAddr, err := to.request(http.MethodGet, path, nil)
	reqInf := ReqInf{CacheHitStatusMiss, remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	data := struct {
		Response []tc.InvalidationJob `json:"response"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, err
}
