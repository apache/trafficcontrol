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
	"encoding/base64"
	"errors"
	"net"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	toclient "github.com/apache/trafficcontrol/traffic_ops/client"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/torequtil"
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
func (cl *TOClient) GetCDNDeliveryServices(cdnID int) ([]tc.DeliveryServiceNullable, bool, error) {
	deliveryServices := []tc.DeliveryServiceNullable{}
	unsupported := false
	err := torequtil.GetRetry(cl.NumRetries, "cdn_"+strconv.Itoa(cdnID)+"_deliveryservices", &deliveryServices, func(obj interface{}) error {
		toDSes, reqInf, err := cl.C.GetDeliveryServicesByCDNID(cdnID)
		if err != nil {
			if errIsUnsupported(err) {
				unsupported = true
				return nil
			}
			return errors.New("getting delivery services from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		dses := obj.(*[]tc.DeliveryServiceNullable)
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

func (cl *TOClient) GetServers() ([]tc.Server, bool, error) {
	servers := []tc.Server{}
	unsupported := false
	err := torequtil.GetRetry(cl.NumRetries, "servers", &servers, func(obj interface{}) error {
		toServers, reqInf, err := cl.C.GetServers()
		if err != nil {
			if errIsUnsupported(err) {
				unsupported = true
				return nil
			}
			return errors.New("getting servers from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		servers := obj.(*[]tc.Server)
		*servers = toServers
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

func (cl *TOClient) GetServerByHostName(serverHostName string) (tc.Server, bool, error) {
	server := tc.Server{}
	unsupported := false
	err := torequtil.GetRetry(cl.NumRetries, "server-name-"+serverHostName, &server, func(obj interface{}) error {
		toServers, reqInf, err := cl.C.GetServerByHostName(serverHostName)
		if err != nil {
			if errIsUnsupported(err) {
				unsupported = true
				return nil
			}
			return errors.New("getting server name '" + serverHostName + "' from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		} else if len(toServers) < 1 {
			return errors.New("getting server name '" + serverHostName + "' from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': no servers returned")
		}
		server := obj.(*tc.Server)
		*server = toServers[0]
		return nil
	})
	if unsupported {
		return tc.Server{}, true, nil
	}
	if err != nil {
		return tc.Server{}, false, errors.New("getting server name '" + serverHostName + "': " + err.Error())
	}
	return server, false, nil
}

func (cl *TOClient) GetCacheGroups() ([]tc.CacheGroupNullable, bool, error) {
	cacheGroups := []tc.CacheGroupNullable{}
	unsupported := false
	err := torequtil.GetRetry(cl.NumRetries, "cachegroups", &cacheGroups, func(obj interface{}) error {
		toCacheGroups, reqInf, err := cl.C.GetCacheGroupsNullable()
		if err != nil {
			if errIsUnsupported(err) {
				unsupported = true
				return nil
			}
			return errors.New("getting cachegroups from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		cacheGroups := obj.(*[]tc.CacheGroupNullable)
		*cacheGroups = toCacheGroups
		return nil
	})
	if unsupported {
		return nil, true, nil
	}
	if err != nil {
		return nil, false, errors.New("getting cachegroups: " + err.Error())
	}
	return cacheGroups, false, nil
}

func (cl *TOClient) GetServerCapabilitiesByID(serverIDs []int) (map[int]map[atscfg.ServerCapability]struct{}, error) {
	serverIDsStr := ""
	if len(serverIDs) > 0 {
		sortIDsInHash := true
		serverIDsStr = base64.RawURLEncoding.EncodeToString((util.HashInts(serverIDs, sortIDsInHash)))
	}

	serverCaps := map[int]map[atscfg.ServerCapability]struct{}{}
	err := torequtil.GetRetry(cl.NumRetries, "server_capabilities_s_"+serverIDsStr, &serverCaps, func(obj interface{}) error {
		// TODO add list of IDs to API+Client
		toServerCaps, reqInf, err := cl.C.GetServerServerCapabilities(nil, nil, nil)
		if err != nil {
			return errors.New("getting server caps from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		serverCaps := obj.(*map[int]map[atscfg.ServerCapability]struct{})

		for _, sc := range toServerCaps {
			if sc.ServerID == nil {
				log.Errorln("Traffic Ops returned Server Capability with nil server id! Skipping!")
			}
			if sc.ServerCapability == nil {
				log.Errorln("Traffic Ops returned Server Capability with nil capability! Skipping!")
			}
			if _, ok := (*serverCaps)[*sc.ServerID]; !ok {
				(*serverCaps)[*sc.ServerID] = map[atscfg.ServerCapability]struct{}{}
			}
			(*serverCaps)[*sc.ServerID][atscfg.ServerCapability(*sc.ServerCapability)] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("getting server server capabilities: " + err.Error())
	}
	return serverCaps, nil
}

func (cl *TOClient) GetDeliveryServiceRequiredCapabilitiesByID(dsIDs []int) (map[int]map[atscfg.ServerCapability]struct{}, error) {
	dsIDsStr := ""
	if len(dsIDs) > 0 {
		sortIDsInHash := true
		dsIDsStr = base64.RawURLEncoding.EncodeToString((util.HashInts(dsIDs, sortIDsInHash)))
	}

	dsCaps := map[int]map[atscfg.ServerCapability]struct{}{}
	err := torequtil.GetRetry(cl.NumRetries, "ds_capabilities_d_"+dsIDsStr, &dsCaps, func(obj interface{}) error {
		// TODO add list of IDs to API+Client
		toDSCaps, reqInf, err := cl.C.GetDeliveryServicesRequiredCapabilities(nil, nil, nil)
		if err != nil {
			return errors.New("getting ds caps from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		dsCaps := obj.(*map[int]map[atscfg.ServerCapability]struct{})

		for _, sc := range toDSCaps {
			if sc.DeliveryServiceID == nil {
				log.Errorln("Traffic Ops returned Delivery Service Capability with nil ds id! Skipping!")
			}
			if sc.RequiredCapability == nil {
				log.Errorln("Traffic Ops returned Delivery Service Capability with nil capability! Skipping!")
			}
			if (*dsCaps)[*sc.DeliveryServiceID] == nil {
				(*dsCaps)[*sc.DeliveryServiceID] = map[atscfg.ServerCapability]struct{}{}
			}
			(*dsCaps)[*sc.DeliveryServiceID][atscfg.ServerCapability(*sc.RequiredCapability)] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("getting ds server capabilities: " + err.Error())
	}
	return dsCaps, nil
}

func errIsUnsupported(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "not found") || strings.Contains(errStr, "not impl")
}

// MaybeIPStr returns the addr string if it isn't nil, or the empty string if it is.
// This is intended for logging, to allow logging with one line, whether addr is nil or not.
func MaybeIPStr(addr net.Addr) string {
	if addr != nil {
		return addr.String()
	}
	return ""
}
