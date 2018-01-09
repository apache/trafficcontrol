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
	"sync/atomic"
	"time"
)

type fetcher struct {
	fetchers []Fetcher
	pos      *uint64
}

func NewRoundRobin(fetchers []Fetcher) Fetcher {
	p := uint64(0)
	return fetcher{fetchers: fetchers, pos: &p}
}

func (f fetcher) Fetch() ([]byte, error) {
	nextI := atomic.AddUint64(f.pos, 1)
	fetcher := f.fetchers[nextI%uint64(len(f.fetchers))]

	return fetcher.Fetch() // TODO round-robin retry on error?
}

func NewHTTPRoundRobin(hosts []string, path string, timeout time.Duration, userAgent string) Fetcher {
	fetchers := []Fetcher{}
	for _, host := range hosts {
		fetchers = append(fetchers, NewHTTP(host+path, timeout, userAgent))
	}
	return NewRoundRobin(fetchers)
}

func NewFileRoundRobin(paths []string, timeout time.Duration, userAgent string) Fetcher {
	fetchers := []Fetcher{}
	for _, path := range paths {
		fetchers = append(fetchers, NewFile(path))
	}
	return NewRoundRobin(fetchers)
}
