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
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"strconv"
	"time"
)

type TextStampSeeker struct {
	Size int64
	Pos  int64
	Seed int64
	Rnd  int64
}

func (s *TextStampSeeker) Read(p []byte) (n int, err error) {
	return ReadBlock(p, s.Pos, s.Size, ReadBlockSize, func(ct int64) []byte {
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
}

func (s *TextStampSeeker) Seek(off int64, whence int) (int64, error) {
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
	DebugLogf("Seeking to %d via %d: %d\n", off, whence, newPos)
	if newPos < 0 {
		return s.Pos, errors.New(`unable to seek before file`)
	}
	s.Pos = newPos
	return s.Pos, nil
}

func TextStamp(w http.ResponseWriter, r *http.Request, reqdat map[string]string) {
	seed, _ := strconv.ParseInt(reqdat[`rnd`], 10, 64)
	lastmod, _ := strconv.ParseInt(reqdat[`lm`], 10, 64)
	sz, _ := strconv.ParseInt(reqdat[`sz`], 10, 64)
	DebugLogf("Serving req of size %d.\n", sz)
	w.Header()[`Content-Type`] = []string{`text/plain`}
	http.ServeContent(w, r, ``, time.Unix(lastmod, 0), &TextStampSeeker{Size: sz, Seed: lastmod, Rnd: seed})
}
