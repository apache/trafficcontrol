package main

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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/traffic_ops/app/bin/checks/DnssecRefresh/config"
)

func main() {
	cfg, err := config.GetCfg()
	config.ErrCheck(err)
	//for the -h --help option
	if cfg == (config.Cfg{}) {
		os.Exit(0)
	}
	log.Debugln("Including DEBUG messages in output. Config is:")
	config.PrintConfig(cfg) // only if DEBUG logging is set.
	body := &config.Creds{
		User:     cfg.TOUser,
		Password: cfg.TOPass,
	}
	loginUrl := cfg.TOUrl + "/api/4.0/user/login"
	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(body)
	config.ErrCheck(err)
	req, err := http.NewRequest(http.MethodPost, loginUrl, buf)
	config.ErrCheck(err)
	jar, err := cookiejar.New(nil)
	config.ErrCheck(err)
	client := &http.Client{Jar: jar, Transport: cfg.Transport, Timeout: 5 * time.Second}

	log.Debugf("Posting to: %s", loginUrl)

	res, err := client.Do(req)
	config.ErrCheck(err)
	defer config.Dclose(res.Body)
	refreshUrl := cfg.TOUrl + "/api/4.0/cdns/dnsseckeys/refresh"
	resp, err := http.NewRequest(http.MethodPut, refreshUrl, nil)
	config.ErrCheck(err)
	log.Debugf("Get req to: %s", refreshUrl)

	refresh, err := client.Do(resp)
	config.ErrCheck(err)
	respData, err := ioutil.ReadAll(refresh.Body)
	config.ErrCheck(err)
	defer config.Dclose(refresh.Body)

	if refresh.StatusCode < 200 || 299 < refresh.StatusCode {
		log.Errorln(string(respData))
		os.Exit(1)
	}
	response := config.ToResponse{}
	config.ErrCheck(json.Unmarshal(respData, &response))
	log.Debugln(response.Response)
}
