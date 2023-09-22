// Package toreqold calls the previous Traffic Ops API major version.
//
// This should never be imported by anything except toreq.
//
// The toreq.Client will automatically fall back to the older client if necessary.
package toreqold

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
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil/toreq/torequtil"
	"github.com/apache/trafficcontrol/v8/lib/go-log"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

type TOClient struct {
	c          *toclient.Session
	NumRetries int
}

func (cl *TOClient) URL() string {
	return cl.c.URL
}

func (cl *TOClient) SetURL(newURL string) {
	cl.c.URL = newURL
}

func (cl *TOClient) HTTPClient() *http.Client {
	return cl.c.Client
}

func (cl *TOClient) APIVersion() string {
	return cl.c.APIVersion()
}

// New logs into Traffic Ops, returning the TOClient which contains the logged-in client.
func New(url *url.URL, user string, pass string, insecure bool, timeout time.Duration, userAgent string) (*TOClient, error) {
	log.Infoln("URL: '" + url.String() + "' User: '" + user + "' Pass len: '" + strconv.Itoa(len(pass)) + "'")

	toURLStr := url.Scheme + "://" + url.Host
	log.Infoln("TO URL string: '" + toURLStr + "'")
	log.Infoln("TO URL: '" + url.String() + "'")

	opts := toclient.Options{}
	opts.Insecure = insecure
	opts.UserAgent = userAgent
	opts.RequestTimeout = timeout
	toClient, inf, err := toclient.Login(toURLStr, user, pass, opts)
	if err != nil {
		return nil, errors.New("Logging in to Traffic Ops '" + torequtil.MaybeIPStr(inf.RemoteAddr) + "': " + err.Error())
	}

	return &TOClient{c: toClient}, nil
}
