package peer

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
	"io"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"

	jsoniter "github.com/json-iterator/go"
)

// Handler handles peer Traffic Monitor data, taking a raw reader, parsing the data, and passing a result object to the ResultChannel. This fulfills the common `Handler` interface.
type Handler struct {
	ResultChannel chan Result
}

// NewHandler returns a new peer Handler.
func NewHandler() Handler {
	return Handler{ResultChannel: make(chan Result)}
}

// Result contains the data parsed from polling a peer Traffic Monitor.
type Result struct {
	ID           tc.TrafficMonitorName
	Available    bool
	Errors       []error
	PeerStates   tc.CRStates
	PollID       uint64
	PollFinished chan<- uint64
	Time         time.Time
}

// Handle handles a response from a polled Traffic Monitor peer, parsing the data and forwarding it to the ResultChannel.
func (handler Handler) Handle(id string, r io.Reader, format string, reqTime time.Duration, reqEnd time.Time, err error, pollID uint64, usingIPv4 bool, pollCtx interface{}, pollFinished chan<- uint64) {
	result := Result{
		ID:           tc.TrafficMonitorName(id),
		Available:    false,
		Errors:       []error{},
		PollID:       pollID,
		PollFinished: pollFinished,
		Time:         reqEnd,
	}

	if err != nil {
		log.Warnf("%s handler given error '%s'\n", id, err.Error())
		result.Errors = append(result.Errors, err)
	}

	if r != nil {
		json := jsoniter.ConfigFastest // TODO make configurable?
		err = json.NewDecoder(r).Decode(&result.PeerStates)
		if err == nil {
			result.Available = true
		} else {
			result.Errors = append(result.Errors, err)
		}
	}

	handler.ResultChannel <- result
}
