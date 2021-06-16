package t3cutil

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

// Utility functions.

import (
	"errors"
	"net/url"
	"strings"

	"github.com/apache/trafficcontrol/cache-config/t3c-generate/toreq"
	"github.com/apache/trafficcontrol/lib/go-log"
)

func TOConnect(tccfg *TCCfg) (*TCCfg, error) {
	toClient, err := toreq.New(
		tccfg.TOURL,
		tccfg.TOUser,
		tccfg.TOPass,
		tccfg.TOInsecure,
		tccfg.TOTimeoutMS,
		tccfg.UserAgent)

	if err != nil {
		return nil, errors.New("failed to connect to traffic ops: " + err.Error())
	}

	if toClient.FellBack() {
		log.Warnln("Traffic Ops does not support the latest version supported by this app! Falling back to previous major Traffic Ops API version!")
	}

	tccfg.TOClient = toClient

	return tccfg, nil
}

func ValidateURL(u *url.URL) error {
	if u == nil {
		return errors.New("nil url")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("scheme expected 'http' or 'https', actual '" + u.Scheme + "'")
	}
	if strings.TrimSpace(u.Host) == "" {
		return errors.New("no host")
	}
	return nil
}
