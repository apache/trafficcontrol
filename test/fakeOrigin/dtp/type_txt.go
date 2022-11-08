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
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"strconv"
)

func init() {
	GlobalGeneratorFuncs[`txt`] = NewTxtGen
}

type TxtGen struct {
	Size int64
	Pos  int64
	Seed int64
	Rnd  int64
}

func (s *TxtGen) ContentType() string {
	return "text/plain"
}

func (s *TxtGen) Read(p []byte) (n int, err error) {
	sz, err := ReadBlock(p, s.Pos, s.Size, ReadBlockSize, func(ct int64) []byte {
		if s.Rnd == 0 {
			return []byte(fmt.Sprintf("%31d\n", ct+s.Seed))
		} else {
			h := fnv.New64()
			b := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, uint64(ct))
			h.Write(b)
			binary.LittleEndian.PutUint64(b, uint64(s.Seed))
			h.Write(b)
			binary.LittleEndian.PutUint64(b, uint64(s.Rnd))
			h.Write(b)
			return []byte(fmt.Sprintf("%31d\n", h.Sum64()&^ByteMask))
		}
	})

	s.Pos += int64(sz)
	return sz, err
}

func (s *TxtGen) Seek(off int64, whence int) (int64, error) {
	posnew, err := NewSeekPosFor(off, whence, s.Pos, s.Size)
	if nil == err {
		s.Pos = posnew
	}
	return s.Pos, nil
}

func NewTxtGen(reqdat map[string]string, lastmod int64) Generator {
	seed, _ := strconv.ParseInt(reqdat[`rnd`], 10, 64)
	sz, _ := strconv.ParseInt(reqdat[`sz`], 10, 64)
	return &TxtGen{Size: sz, Seed: lastmod, Rnd: seed}
}
