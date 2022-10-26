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
	"net/http"
	"strconv"
	"time"
)

func init() {
	GlobalForwarderFuncs["delay"] = NewDelayForwardGen
}

type DelayForward struct {
	Generator Generator
	Latency   time.Duration // per request latency
	FirstTime bool
}

func (ss *DelayForward) ContentType() string {
	return ss.Generator.ContentType()
}

func (ss *DelayForward) Read(bufout []byte) (bytes int, err error) {
	if ss.FirstTime {
		ss.FirstTime = false
		bytes, err = ss.Generator.Read(bufout)
	} else {
		timer := time.NewTimer(ss.Latency)
		bytes, err = ss.Generator.Read(bufout)
		<-timer.C
	}

	return bytes, err
}

func (ss *DelayForward) Seek(off int64, whence int) (int64, error) {
	return ss.Generator.Seek(off, whence)
}

func NewDelayForwardGen(
	w http.ResponseWriter,
	r *http.Request,
	gen Generator,
	reqdat map[string]string, updated int64,
) Generator {
	delaystr := reqdat[`delay`]
	latency, err := time.ParseDuration(delaystr)
	if err != nil {
		seed, _ := strconv.ParseInt(reqdat[`rnd`], 10, 64)
		latencyns := EvalNumber(delaystr, seed)
		latency = time.Duration(latencyns)
		reqdat[`delay`] = latency.String() // save this for next time
	}
	return &DelayForward{
		Generator: gen,
		Latency:   latency,
		FirstTime: true,
	}
}
