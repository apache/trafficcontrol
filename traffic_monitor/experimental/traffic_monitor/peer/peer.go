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
	"encoding/json"
	"io"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
)

// Handler handles peer Traffic Monitor data, taking a raw reader, parsing the data, and passing a result object to the ResultChannel. This fulfills the common `Handler` interface.
type Handler struct {
	ResultChannel chan Result
	Notify        int
}

// NewHandler returns a new peer Handler.
func NewHandler() Handler {
	return Handler{ResultChannel: make(chan Result)}
}

// Result contains the data parsed from polling a peer Traffic Monitor.
type Result struct {
	ID           enum.TrafficMonitorName
	Available    bool
	Errors       []error
	PeerStates   Crstates
	PollID       uint64
	PollFinished chan<- uint64
	Time         time.Time
}

// Handle handles a response from a polled Traffic Monitor peer, parsing the data and forwarding it to the ResultChannel.
func (handler Handler) Handle(id string, r io.Reader, reqTime time.Duration, err error, pollID uint64, pollFinished chan<- uint64) {
	result := Result{
		ID:           enum.TrafficMonitorName(id),
		Available:    false,
		Errors:       []error{},
		PollID:       pollID,
		PollFinished: pollFinished,
		Time:         time.Now(),
	}

	if err != nil {
		result.Errors = append(result.Errors, err)
	}

	if r != nil {
		dec := json.NewDecoder(r)
		err = dec.Decode(&result.PeerStates)

		if err == nil {
			result.Available = true
		} else {
			result.Errors = append(result.Errors, err)
		}
	}

	handler.ResultChannel <- result
}
