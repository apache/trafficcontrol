package fetcher

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
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/handler"
)

type Fetcher interface {
	Fetch(id string, url string, host string, pollId uint64, pollFinishedChan chan<- uint64)
}

type HttpFetcher struct {
	Client    *http.Client
	UserAgent string
	Headers   map[string]string
	Handler   handler.Handler
}

type Result struct {
	Source string
	Data   []byte
	Error  error
}

func (f HttpFetcher) Fetch(id string, url string, host string, pollId uint64, pollFinishedChan chan<- uint64) {
	log.Debugf("poll %v %v fetch start\n", pollId, time.Now())
	req, err := http.NewRequest("GET", url, nil)
	// TODO: change this to use f.Headers. -jse
	req.Header.Set("User-Agent", f.UserAgent)
	req.Header.Set("Connection", "keep-alive")
	req.Host = host
	startReq := time.Now()
	response, err := f.Client.Do(req)
	reqEnd := time.Now()
	reqTime := reqEnd.Sub(startReq)
	defer func() {
		if response != nil && response.Body != nil {
			ioutil.ReadAll(response.Body) // TODO determine if necessary
			response.Body.Close()
		}
	}()

	if err == nil && response == nil {
		err = fmt.Errorf("err nil and response nil")
	}
	if err == nil && response != nil && (response.StatusCode < 200 || response.StatusCode > 299) {
		err = fmt.Errorf("bad status: %v", response.StatusCode)
	}
	if err != nil {
		err = fmt.Errorf("id %v url %v fetch error: %v", id, url, err)
	}

	if err == nil && response != nil {
		log.Debugf("poll %v %v fetch end\n", pollId, time.Now())
		f.Handler.Handle(id, response.Body, reqTime, reqEnd, err, pollId, pollFinishedChan)
	} else {
		f.Handler.Handle(id, nil, reqTime, reqEnd, err, pollId, pollFinishedChan)
	}
}
