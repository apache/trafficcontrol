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
	"strconv"
	"sync"
)

func init() {
	GlobalGeneratorFuncs["gen3s"] = NewGen3sGen
}

var Buf3Cache []byte
var Buf3Mutex sync.Mutex

type Gen3sGen struct {
	Size int64
	Pos  int64
}

func (s *Gen3sGen) ContentType() string {
	return "text/plain"
}

func (s *Gen3sGen) Read(p []byte) (n int, err error) {
	plen := len(p)

	Buf3Mutex.Lock()

	// we could be more clever about this lock but we
	// expect it to only be hit a couple of times
	// as 32k seems to be the observed max
	if len(Buf3Cache) < plen {
		Buf3Cache = make([]byte, plen)
		for index := range Buf3Cache {
			Buf3Cache[index] = byte('3')
		}
	}

	cbuf := Buf3Cache

	Buf3Mutex.Unlock()

	copy(p, cbuf[:plen])
	s.Pos += int64(plen)

	return plen, nil
}

func (s *Gen3sGen) Seek(off int64, whence int) (int64, error) {
	posnew, err := NewSeekPosFor(off, whence, s.Pos, s.Size)
	if nil == err {
		s.Pos = posnew
	}
	return s.Pos, nil
}

func NewGen3sGen(reqdat map[string]string, lastmod int64) Generator {
	sz, _ := strconv.ParseInt(reqdat[`sz`], 10, 64)
	return &Gen3sGen{Size: sz}
}
