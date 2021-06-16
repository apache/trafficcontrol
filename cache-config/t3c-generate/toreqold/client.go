// toreqold implements a Traffic Ops client vendored one version back.
//
// This should be used for all requests, unless they require an endpoint or field added in the latest version.
//
// If a feature in the latest Traffic Ops is required, toreqnew should be used with a fallback to this client if the Traffic Ops is not the latest (which can be determined by the bool returned by all toreqnew.TOClient funcs).
//
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
	"net/url"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/cache-config/t3c-generate/torequtil"
	"github.com/apache/trafficcontrol/lib/go-log"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v3-client"
)

type TOClient struct {
	C          *toclient.Session
	NumRetries int
}

const isFallback = true

// New logs into Traffic Ops, returning the TOClient which contains the logged-in client.
func New(url *url.URL, user string, pass string, insecure bool, timeout time.Duration, userAgent string) (*TOClient, error) {
	log.Infoln("URL: '" + url.String() + "' User: '" + user + "' Pass len: '" + strconv.Itoa(len(pass)) + "'")

	toURLStr := url.Scheme + "://" + url.Host
	log.Infoln("TO URL string: '" + toURLStr + "'")
	log.Infoln("TO URL: '" + url.String() + "'")

	opts := toclient.ClientOpts{}
	opts.Insecure = insecure
	opts.UserAgent = userAgent
	opts.RequestTimeout = timeout
	toClient, inf, err := toclient.Login(toURLStr, user, pass, opts)
	if err != nil {
		return nil, errors.New("Logging in to Traffic Ops '" + torequtil.MaybeIPStr(inf.RemoteAddr) + "': " + err.Error())
	}

	return &TOClient{C: toClient}, nil
}
