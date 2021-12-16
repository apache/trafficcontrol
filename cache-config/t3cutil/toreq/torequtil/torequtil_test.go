package torequtil

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
	"errors"
	"net"
	"net/http"
	"strings"
	"testing"
)

func TestGetRetry(t *testing.T) {
	getterCalled := 0
	getter := func(obj interface{}) error {
		getterCalled++
		return errors.New("something")
	}

	numRetries := 2
	err := GetRetry(numRetries, "foo", nil, getter)
	if err == nil {
		t.Fatal("GetRetry expected error from f, actual nil")
	}
	if errStr := err.Error(); !strings.Contains(errStr, "something") {
		t.Errorf("GetRetry expected error from getter 'something', actual '" + errStr + "'")
	}
	if getterCalled != numRetries+1 { // +1 because the first call isn't a retry.
		t.Errorf("GetRetry expected to call getter numRetries %v +1 times, actual %v\n", numRetries, getterCalled)
	}
}

func TestRetryBackoffSeconds(t *testing.T) {
	// Just test that it's greater. We don't want a brittle test that tests exactly how we know the function increases.
	currentRetry := 0
	newRetry := RetryBackoffSeconds(currentRetry)
	if newRetry <= currentRetry {
		t.Errorf("RetryBackoffSeconds expected greater than current retry %v actual %v", currentRetry, newRetry)
	}
	newerRetry := RetryBackoffSeconds(newRetry)
	if newerRetry <= newRetry {
		t.Errorf("RetryBackoffSeconds expected greater than new retry %v actual %v", currentRetry, newRetry)
	}
}

func TestMaybeIPStr(t *testing.T) {
	if is := MaybeIPStr(nil); is != "" {
		t.Errorf("MaybeIPStr(nil) expected '', actual '%v'", is)
	}
	addr := &net.IPAddr{IP: net.ParseIP("192.0.2.1")}
	if is := MaybeIPStr(addr); is != "192.0.2.1" {
		t.Errorf("MaybeIPStr(nil) expected '192.0.2.1', actual '%v'", is)
	}
}

func TestMaybeHdrStr(t *testing.T) {
	if is := MaybeHdrStr(nil, "Age"); is != "" {
		t.Errorf("MaybeHdrStr(nil) expected '', actual '%v'", is)
	}
	hdr := http.Header{"Age": {"1001"}}
	if is := MaybeHdrStr(hdr, "Age"); is != "1001" {
		t.Errorf("MaybeIPStr(val) expected '1001', actual '%v'", is)
	}
}
