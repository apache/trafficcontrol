package fetch

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
 *
 */

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type httpFetcher struct {
	url       string
	timeout   time.Duration
	userAgent string
}

func NewHTTP(url string, timeout time.Duration, userAgent string) Fetcher {
	return httpFetcher{url: url, timeout: timeout}
}

func (f httpFetcher) Fetch() ([]byte, error) {
	client := http.Client{
		Timeout: f.timeout,
	}
	req, err := http.NewRequest("GET", f.url, nil)
	if err != nil {
		// TODO round-robin retry on error?
		return nil, errors.New("HTTP creating request '" + f.url + "': " + err.Error())
	}

	req.Header.Set("User-Agent", f.userAgent)

	resp, err := client.Do(req)
	if err != nil {
		// TODO round-robin retry on error?
		return nil, errors.New("HTTP request '" + f.url + "': " + err.Error())
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("HTTP reading response body '" + f.url + "': " + err.Error())
	}
	return b, nil
}
