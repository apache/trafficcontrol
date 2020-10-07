package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-log"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

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

package main

import (
"bytes"
"encoding/json"
"fmt"
"github.com/apache/trafficcontrol/lib/go-log"
"github.com/apache/trafficcontrol/traffic_ops/app/bin/checks/DnssecRefresh/config"
"io/ioutil"
"net/http"
"net/http/cookiejar"
)

func main() {
	cfg, err := config.GetCfg()
	config.ErrCheck(err)
	log.Debugln("Including DEBUG messages in output. Config is:")
	config.PrintConfig(cfg) // only if DEBUG logging is set.
	body := &config.Creds{
		User: cfg.TOUser,
		Password: cfg.TOPass,
	}
	loginUrl := fmt.Sprintf("%s%s", cfg.TOUrl, "/api/2.0/user/login")
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(body)
	config.ErrCheck(err)
	req, _ := http.NewRequest("POST", loginUrl, buf)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar, Transport: cfg.Tr}

	log.Debugf("Posting to: %s", loginUrl)

	res, err := client.Do(req)
	config.ErrCheck(err)
	defer config.Dclose(res.Body)
	refreshUrl := fmt.Sprintf("%s%s", cfg.TOUrl, "/api/2.0/cdns/dnsseckeys/refresh")
	resp, _ := http.NewRequest("GET", refreshUrl, buf)

	log.Debugf("Get req to: %s", refreshUrl)

	refresh, _ := client.Do(resp)
	respData, _ := ioutil.ReadAll(refresh.Body)
	var response config.ToResponse
	config.ErrCheck(json.Unmarshal(respData, &response))
	log.Debugln(response.Response)
}
