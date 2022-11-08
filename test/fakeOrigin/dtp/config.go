package dtp

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
	"fmt"
	"net/http"
	"time"
)

type Log struct {
	Access          bool   `json:"access"`
	Path            string `json:"path"`
	RequestHeaders  bool   `json:"request_headers"`
	ResponseHeaders bool   `json:"response_headers"`
}

type Timeout struct {
	Read  time.Duration `json:"read"`
	Write time.Duration `json:"write"`
	Idle  time.Duration `json:"idle"`
}

type Config struct {
	Debug         bool          `json:"debug"`
	EnablePprof   bool          `json:"enable_pprof"`
	Log           Log           `json:"log"`
	Timeout       Timeout       `json:"timeout"`
	StallDuration time.Duration `json:"stall_duration"`
}

func NewConfig() Config {
	var cfg Config
	cfg.Log.Access = true
	cfg.Log.Path = "dtp.log"
	cfg.Log.RequestHeaders = false
	cfg.Log.ResponseHeaders = false
	cfg.EnablePprof = false
	cfg.Timeout.Read = time.Duration(10) * time.Second
	cfg.Timeout.Write = time.Duration(10) * time.Second
	cfg.Timeout.Idle = time.Duration(10) * time.Second
	cfg.StallDuration = time.Duration(0)
	return cfg
}

var GlobalConfig = NewConfig()

// handle api configuration endpoint
func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	{
		dbghdrs := r.URL.Query().Get("debug")
		if dbghdrs != "" {
			fmt.Println("processing debug", dbghdrs)
			GlobalConfig.Debug = (dbghdrs == "true")
			fmt.Println("debugging:", GlobalConfig.Debug)
		}
	}
	{
		reqhdrs := r.URL.Query().Get("request_headers")
		if reqhdrs != "" {
			GlobalConfig.Log.RequestHeaders = (reqhdrs == "true")
			fmt.Println("req header logging:", GlobalConfig.Log.RequestHeaders)
		}
	}
	{
		reshdrs := r.URL.Query().Get("response_headers")
		if reshdrs != "" {
			GlobalConfig.Log.ResponseHeaders = (reshdrs == "true")
			fmt.Println("resp header logging:", GlobalConfig.Log.ResponseHeaders)
		}
	}
	{
		hdrs := r.URL.Query().Get("all_headers")
		if hdrs != "" {
			if hdrs == "true" {
				GlobalConfig.Log.RequestHeaders = true
				GlobalConfig.Log.ResponseHeaders = true
				fmt.Println("req/resp header logging:", true)
			} else {
				GlobalConfig.Log.RequestHeaders = false
				GlobalConfig.Log.ResponseHeaders = false
				fmt.Println("req/resp header logging:", false)
			}
		}
	}
	{
		stallhdr := r.URL.Query().Get("stall_duration")
		if stallhdr != "" {
			dur, err := time.ParseDuration(stallhdr)
			if nil == err {
				GlobalConfig.StallDuration = dur
				fmt.Println("stall_duration:", dur)
			} else {
				fmt.Println("error setting stall_duration", err)
			}
		}
	}

	// default is to dump current config
	bytes, _ := json.Marshal(&GlobalConfig)
	w.Write(bytes)
}
