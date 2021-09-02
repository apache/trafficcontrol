package util

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
	"net/http/httptest"
	"testing"
)

func TestInterceptor_WriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	interceptor := Interceptor{
		W: w,
	}

	interceptor.WriteHeader(http.StatusAccepted)
	result := w.Result()
	if result.StatusCode != http.StatusAccepted {
		t.Errorf("Incorrect status code written, expected: %d, got: %d", http.StatusAccepted, result.StatusCode)
	}
}

func TestInterceptor_Write(t *testing.T) {
	w := httptest.NewRecorder()
	interceptor := Interceptor{
		W: w,
	}

	const data = "some data"
	n, err := interceptor.Write([]byte(data))
	if err != nil {
		t.Errorf("Unexpected error writing a response: %v", err)
	}
	if n != len([]byte(data)) {
		t.Errorf("Incorrect number of bytes written, expected: %d, got: %d", len([]byte(data)), n)
	}
	if n != interceptor.ByteCount {
		t.Errorf("Incorrect reported total number of bytes written, expected: %d, got: %d", n, interceptor.ByteCount)
	}

	result := w.Result()
	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected default status code to be %d, got: %d", http.StatusOK, result.StatusCode)
	}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	} else if string(body) != data {
		t.Errorf("Incorrect response body, expected: '%s', got: '%s'", data, string(body))
	}

	w = httptest.NewRecorder()
	interceptor.W = w
	interceptor.ByteCount = 0
	interceptor.Code = http.StatusAccepted

	n, err = interceptor.Write([]byte(data))
	if err != nil {
		t.Errorf("Unexpected error writing a response: %v", err)
	}
	if n != len([]byte(data)) {
		t.Errorf("Incorrect number of bytes written, expected: %d, got: %d", len([]byte(data)), n)
	}
	if n != interceptor.ByteCount {
		t.Errorf("Incorrect reported total number of bytes written, expected: %d, got: %d", n, interceptor.ByteCount)
	}

	result = w.Result()
	// TODO: The interceptor currently only overrwrites the Code if it's zero -
	// but in the event that a Write is called without first calling WriteHeader
	// like this, the Code is ignored and will not match the actual response code.
	// Should it always override the Code? Or should it first write out a header
	// if that hasn't been done yet? In any case, there's nothing to prevent
	// calls from manipulating the Code after a response has been written.
	// if result.StatusCode != http.StatusAccepted {
	// 	t.Errorf("Incorrect status code, expected: %d, got: %d", http.StatusAccepted, result.StatusCode)
	// }
	body, err = ioutil.ReadAll(result.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	} else if string(body) != data {
		t.Errorf("Incorrect response body, expected: '%s', got: '%s'", data, string(body))
	}
}

func TestBodyInterceptor_WriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	interceptor := BodyInterceptor{
		W: w,
	}

	interceptor.WriteHeader(http.StatusAccepted)
	result := w.Result()
	if result.StatusCode != http.StatusAccepted {
		t.Errorf("Incorrect status code written, expected: %d, got: %d", http.StatusAccepted, result.StatusCode)
	}
}

func ExampleInterceptor_Header() {
	i := Interceptor{W: httptest.NewRecorder()}
	i.W.Header().Add("test", "quest")
	fmt.Println(i.Header().Get("test"))
	// Output: quest
}

func TestBodyInterceptor_Write(t *testing.T) {
	w := httptest.NewRecorder()
	interceptor := BodyInterceptor{
		W: w,
	}

	const data = "some data"
	dataLen := len([]byte(data))
	const moreData = " some more data"
	moreDataLen := len([]byte(moreData))

	n, err := interceptor.Write([]byte(data))
	if err != nil {
		t.Errorf("Unexpected error writing a response: %v", err)
	}
	if n != dataLen {
		t.Errorf("Incorrect number of bytes written, expected: %d, got: %d", dataLen, n)
	}
	body := interceptor.Body()
	if string(body) != data {
		t.Errorf("Incorrect cached body, expected: '%s', got: '%s'", data, string(body))
	}
	m, err := interceptor.Write([]byte(moreData))
	if err != nil {
		t.Errorf("Unexpected error writing a response: %v", err)
	}
	if m != moreDataLen {
		t.Errorf("Incorrect number of bytes written, expected: %d, got: %d", moreDataLen, m)
	}
	body = interceptor.Body()
	if string(body) != data+moreData {
		t.Errorf("Incorrect cached body, expected: '%s', got: '%s'", data+moreData, string(body))
	}
	total, err := interceptor.RealWrite(nil)
	if err != nil {
		t.Errorf("Unexpected error writing a *real* response: %v", err)
	}
	if total != n+m {
		t.Errorf("Incorrect total number of bytes written, expected: %d, got: %d", n+m, total)
	}

	result := w.Result()
	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected default status code to be %d, got: %d", http.StatusOK, result.StatusCode)
	}
	body, err = ioutil.ReadAll(result.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	} else if string(body) != data+moreData {
		t.Errorf("Incorrect response body, expected: '%s', got: '%s'", data+moreData, string(body))
	}
}

func ExampleBodyInterceptor_Header() {
	i := BodyInterceptor{W: httptest.NewRecorder()}
	i.W.Header().Add("test", "quest")
	fmt.Println(i.Header().Get("test"))
	// Output: quest
}
