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

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/handler"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/log"
	"github.com/davecheney/gmx"
)

type Fetcher interface {
	Fetch(string, string, uint64, chan<- uint64)
}

type HttpFetcher struct {
	Client  *http.Client
	Headers map[string]string
	Handler handler.Handler
	Counters
}

type Result struct {
	Source string
	Data   []byte
	Error  error
}

type Counters struct {
	Success *gmx.Counter
	Fail    *gmx.Counter
	Pending *gmx.Gauge
}

func (f HttpFetcher) Fetch(id string, url string, pollId uint64, pollFinishedChan chan<- uint64) {
	log.Debugf("poll %v %v fetch start\n", pollId, time.Now())
	req, err := http.NewRequest("GET", url, nil)
	// TODO: change this to use f.Headers. -jse
	req.Header.Set("User-Agent", "traffic_monitor/1.0") // TODO change to 2.0?
	req.Header.Set("Connection", "keep-alive")
	if f.Pending != nil {
		f.Pending.Inc()
	}
	response, err := f.Client.Do(req)
	if f.Pending != nil {
		f.Pending.Dec()
	}
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
		err = fmt.Errorf("fetch error: %v", err)
	}

	if err == nil && response != nil {
		if f.Success != nil {
			f.Success.Inc()
		}
		log.Debugf("poll %v %v fetch end\n", pollId, time.Now())
		f.Handler.Handle(id, response.Body, err, pollId, pollFinishedChan)
	} else {
		if f.Fail != nil {
			f.Fail.Inc()
		}
		f.Handler.Handle(id, nil, err, pollId, pollFinishedChan)
	}
}
