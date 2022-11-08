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
	"errors"
	"hash/fnv"
	"io"
	"net/http"
	"strconv"
	"time"
)

type BinStampSeeker struct {
	Size int64
	Pos  int64
	Seed int64
	Rnd  int64
}

func (s *BinStampSeeker) Read(p []byte) (n int, err error) {
	return ReadBlock(p, s.Pos, s.Size, 8, func(ct int64) []byte {
		ct += s.Seed
		b := make([]byte, 8)
		if s.Rnd == 0 {
			for i := 0; i < len(b); i++ {
				b[i] = byte(ct & 0xf)
				ct >>= 8
			}
		} else {
			h := fnv.New64()
			binary.LittleEndian.PutUint64(b, uint64(ct))
			h.Write(b)
			binary.LittleEndian.PutUint64(b, uint64(s.Rnd))
			h.Write(b)
			return h.Sum(nil)
		}
		return b
	})
}

func (s *BinStampSeeker) Seek(off int64, whence int) (int64, error) {
	var newPos int64
	switch whence {
	case io.SeekStart:
		newPos = off
	case io.SeekEnd:
		newPos = s.Size - off
	case io.SeekCurrent:
		newPos += off
		if newPos > s.Size {
			newPos = s.Size
		}
	}
	if newPos < 0 {
		return s.Pos, errors.New(`unable to seek before file`)
	}
	s.Pos = newPos
	return s.Pos, nil
}

func BinStamp(w http.ResponseWriter, r *http.Request, reqdat map[string]string) {
	seed, _ := strconv.ParseInt(reqdat[`rnd`], 10, 64)
	lastmod, _ := strconv.ParseInt(reqdat[`lm`], 10, 64)
	sz, _ := strconv.ParseInt(reqdat[`sz`], 10, 64)
	w.Header()[`Content-Type`] = []string{`application/octet-stream`}
	http.ServeContent(w, r, ``, time.Unix(lastmod, 0), &BinStampSeeker{Size: sz, Seed: lastmod, Rnd: seed})
}
