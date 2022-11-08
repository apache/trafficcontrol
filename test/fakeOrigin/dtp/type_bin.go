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
	"hash/fnv"
	"strconv"
)

func init() {
	GlobalGeneratorFuncs["bin"] = NewBinGen
}

type BinGen struct {
	Size     int64
	Pos      int64
	Seed     int64
	Rnd      int64
	BufCache []byte
}

func (bs *BinGen) ContentType() string {
	return "application/octet-stream"
}

func (bs *BinGen) Read(p []byte) (n int, err error) {
	sz, err := ReadBlock(p, bs.Pos, bs.Size, 8, func(ct int64) []byte {
		val := ct + bs.Seed
		if bs.Rnd == 0 {
			for index := 0; index < 8; index++ {
				bs.BufCache[index] = byte(val & 0xff)
				val >>= 8
			}
			return bs.BufCache
		} else {
			hash := fnv.New64()
			binary.LittleEndian.PutUint64(bs.BufCache, uint64(val))
			hash.Write(bs.BufCache)
			binary.LittleEndian.PutUint64(bs.BufCache, uint64(bs.Rnd))
			hash.Write(bs.BufCache)
			return hash.Sum(nil)
		}
	})

	bs.Pos += int64(sz)
	return sz, err
}

func (s *BinGen) Seek(off int64, whence int) (int64, error) {
	posnew, err := NewSeekPosFor(off, whence, s.Pos, s.Size)
	if nil == err {
		s.Pos = posnew
	}
	return s.Pos, nil
}

func NewBinGen(reqdat map[string]string, lastmod int64) Generator {
	rnd, _ := strconv.ParseInt(reqdat[`rnd`], 10, 64)
	sz, _ := strconv.ParseInt(reqdat[`sz`], 10, 64)
	return &BinGen{Size: sz, Seed: lastmod, Rnd: rnd, BufCache: make([]byte, 8)}
}
