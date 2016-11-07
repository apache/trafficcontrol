package handler

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
	"encoding/json"
	"io"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/log"
)

const (
	NOTIFY_NEVER = iota
	NOTIFY_CHANGE
	NOTIFY_ALWAYS
)

type Handler interface {
	Handle(string, io.Reader, time.Duration, error, uint64, chan<- uint64)
}

type OpsConfigFileHandler struct {
	Content          interface{}
	ResultChannel    chan interface{}
	OpsConfigChannel chan OpsConfig
}

type OpsConfig struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Url          string `json:"url"`
	Insecure     bool   `json:"insecure"`
	CdnName      string `json:"cdnName"`
	HttpListener string `json:"httpListener"`
}

func (handler OpsConfigFileHandler) Listen() {
	for {
		result := <-handler.ResultChannel
		var toc OpsConfig

		err := json.Unmarshal(result.([]byte), &toc)

		if err != nil {
			log.Errorf("unmarshalling JSON: %s\n", err)
		} else {
			handler.OpsConfigChannel <- toc
		}
	}
}
