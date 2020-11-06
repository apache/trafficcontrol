// toreqnew implements a Traffic Ops client for features in the latest version.
//
// This should only be used if an endpoint or field needed for config gen is in the latest.
//
// Users should always check the returned bool, and if it's false, call the vendored toreq client and set the proper defaults for the new feature(s).
//
// All TOClient functions should check for 404 or 503 and return a bool false if so.
//
package toreqnew

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"errors"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"

	toclient "github.com/apache/trafficcontrol/traffic_ops/v3-client"
	"github.com/apache/trafficcontrol/traffic_ops_ort/atstccfg/torequtil"
)

type TOClient struct {
	C          *toclient.Session
	NumRetries int
}

// New returns a TOClient with the given credentials.
// Note it does not actually log in or try to make a request. Rather, it assumes the cookies are valid for a session. No external communication is made.
func New(cookies string, url *url.URL, user string, pass string, insecure bool, timeout time.Duration, userAgent string) (*TOClient, error) {
	log.Infoln("URL: '" + url.String() + "' User: '" + user + "' Pass len: '" + strconv.Itoa(len(pass)) + "'")

	useCache := false
	toClient := toclient.NewNoAuthSession(url.String(), insecure, userAgent, useCache, timeout)
	toClient.UserName = user
	toClient.Password = pass

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, errors.New("making cookie jar: " + err.Error())
	}
	toClient.Client.Jar = jar
	toClient.Client.Jar.SetCookies(url, torequtil.StringToCookies(cookies))

	return &TOClient{C: toClient}, nil
}

// GetCDNDeliveryServices returns the deliveryservices, whether this client's version is unsupported by the server, and any error.
// Note if the server returns a 404 or 503, this returns false and a nil error.
// Users should check the "not supported" bool, and use the vendored TOClient if it's set, and set proper defaults for the new feature(s).
func (cl *TOClient) GetCDNDeliveryServices(cdnID int) ([]tc.DeliveryServiceNullableV30, bool, error) {
	deliveryServices := []tc.DeliveryServiceNullableV30{}
	unsupported := false
	err := torequtil.GetRetry(cl.NumRetries, "cdn_"+strconv.Itoa(cdnID)+"_deliveryservices", &deliveryServices, func(obj interface{}) error {
		params := url.Values{}
		params.Set("cdn", strconv.Itoa(cdnID))
		toDSes, reqInf, err := cl.C.GetDeliveryServicesV30WithHdr(nil, params)
		if err != nil {
			if IsUnsupportedErr(err) {
				unsupported = true
				return nil
			}
			return errors.New("getting delivery services from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		dses := obj.(*[]tc.DeliveryServiceNullableV30)
		*dses = toDSes
		return nil
	})
	if unsupported {
		return nil, true, nil
	}
	if err != nil {
		return nil, false, errors.New("getting delivery services: " + err.Error())
	}
	return deliveryServices, false, nil
}

func (cl *TOClient) GetTopologies() ([]tc.Topology, bool, error) {
	topologies := []tc.Topology{}
	unsupported := false
	err := torequtil.GetRetry(cl.NumRetries, "topologies", &topologies, func(obj interface{}) error {
		toTopologies, reqInf, err := cl.C.GetTopologies()
		if err != nil {
			if IsUnsupportedErr(err) {
				unsupported = true
				return nil
			}
			return errors.New("getting topologies from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		topologies := obj.(*[]tc.Topology)
		*topologies = toTopologies
		return nil
	})
	if unsupported {
		return nil, true, nil
	}
	if err != nil {
		return nil, false, errors.New("getting topologies: " + err.Error())
	}
	return topologies, false, nil
}

func (cl *TOClient) GetServerUpdateStatus(cacheHostName tc.CacheName) (tc.ServerUpdateStatus, bool, error) {
	status := tc.ServerUpdateStatus{}
	unsupported := false
	err := torequtil.GetRetry(cl.NumRetries, "server_update_status_"+string(cacheHostName), &status, func(obj interface{}) error {
		toStatus, reqInf, err := cl.C.GetServerUpdateStatus(string(cacheHostName))
		if err != nil {
			if IsUnsupportedErr(err) {
				unsupported = true
				return nil
			}
			return errors.New("getting server update status from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		status := obj.(*tc.ServerUpdateStatus)
		*status = toStatus
		return nil
	})
	if unsupported {
		return tc.ServerUpdateStatus{}, true, nil
	}
	if err != nil {
		return tc.ServerUpdateStatus{}, false, errors.New("getting server update status: " + err.Error())
	}
	return status, false, nil
}

func (cl *TOClient) GetServers() ([]tc.ServerV30, bool, error) {
	servers := []tc.ServerV30{}
	unsupported := false
	err := torequtil.GetRetry(cl.NumRetries, "servers", &servers, func(obj interface{}) error {
		toServers, reqInf, err := cl.C.GetServersWithHdr(nil, nil)
		if err != nil {
			if IsUnsupportedErr(err) {
				unsupported = true
				return nil
			}
			return errors.New("getting servers from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		servers := obj.(*[]tc.ServerV30)
		*servers = toServers.Response
		return nil
	})
	if unsupported {
		return nil, true, nil
	}
	if err != nil {
		return nil, false, errors.New("getting servers: " + err.Error())
	}
	return servers, false, nil
}

func IsUnsupportedErr(err error) bool {
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "not found") || strings.Contains(errStr, "not impl")
}
