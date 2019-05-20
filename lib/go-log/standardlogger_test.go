package log

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
	"log"
	"strings"
	"testing"
	"time"
)

func TestStandardLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	realLogPrefix := "a-real-log-prefix: "
	logger := log.New(buf, realLogPrefix, log.Lshortfile)

	sPrefix := "a-standard-logger-prefix: "
	sLogger := StandardLogger(logger, sPrefix)

	msg := "a-message"
	sLogger.Printf(msg + "\n") // newline, to verify a double-newline isn't printed.

	actual := buf.String()

	if !strings.HasPrefix(actual, realLogPrefix) {
		t.Errorf("expected prefix '%v' actual '%v'\n", realLogPrefix, actual)
	}

	actualFields := strings.Fields(actual)
	if len(actualFields) < 4 {
		t.Fatalf("expected fields >4 (prefix, line, time, msg} actual %v '''%v'''\n", len(actualFields), actual)
	}

	timeField := actualFields[2]
	if len(timeField) > 0 {
		timeField = timeField[:len(timeField)-1] // timestamp ends with :, strip it for parsing
	}

	actualTime, err := time.Parse(time.RFC3339Nano, timeField)
	if err != nil {
		t.Fatalf("expected 3rd field is RFC3339 nano format, actual '%v'\n", timeField)
	}
	if actualTime.After(time.Now().Add(time.Minute)) || actualTime.Before(time.Now().Add(time.Minute*-1)) {
		t.Errorf("expected 3rd field is RFC3339 nano format around now '%v', actual '%v'\n", time.Now(), actualTime)
	}
}

func TestStandardLoggerNil(t *testing.T) {
	// test that a nil logger doesn't panic
	StandardLogger(nil, "").Println("foo")
}
