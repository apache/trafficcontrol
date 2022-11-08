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
	"errors"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func init() {
	GlobalForwarderFuncs["posevt"] = NewEvtForwardGen
}

type RangeBE struct {
	Begin int64
	End   int64
}

func parseRange(rangestr string, bsize int64) (RangeBE, error) {
	var beg int64 = -1
	var end int64 = -1

	const bstr = `bytes=`
	bpos := strings.Index(rangestr, bstr)
	if 0 <= bpos {
		rstr := rangestr[bpos+len(bstr):]
		fpos := strings.Index(rstr, `,`)
		if 0 < fpos {
			rstr = rstr[0:fpos]
		}
		fields := strings.Split(rstr, `-`)
		if len(fields) == 2 {
			v0, ok0 := strconv.ParseInt(fields[0], 10, 64)
			v1, ok1 := strconv.ParseInt(fields[1], 10, 64)

			if nil == ok0 && nil == ok1 {
				beg = v0
				end = v1 + 1 // convert from min/max to begin/end
			} else if nil == ok0 {
				beg = v0
				end = bsize
			} else if nil == ok1 {
				beg = bsize - v1
				end = bsize
			}
		}
	}

	if beg < 0 || end < beg {
		DebugLogf("bad range %d %d\n", beg, end)
		return RangeBE{}, errors.New("bad range passed in")
	} else {
		DebugLogf("parsed range %d %d\n", beg, end)
		return RangeBE{beg, end}, nil
	}
}

func (rng *RangeBE) Contains(val int64) bool {
	return rng.Begin <= val && val < rng.End
}

type EvtForward struct {
	generator Generator
	posevt    int64
}

func (ss *EvtForward) ContentType() string {
	return ss.generator.ContentType()
}

func (ss *EvtForward) Read(bufout []byte) (n int, err error) {

	posstart, _ := ss.generator.Seek(0, io.SeekCurrent)
	bytes, err := ss.generator.Read(bufout)

	errout := err

	if nil == err && math.MaxInt64 != ss.posevt {
		// determine if we need to cut the buffer short
		failbytes := ss.posevt - posstart
		if 0 <= failbytes && failbytes < int64(bytes) {

			// fix up the underlying seeker< dump an EOF
			ss.generator.Seek(posstart+failbytes, io.SeekStart)
			bytes = int(failbytes)
			errout = io.EOF
		}
	}

	return bytes, errout
}

func (ss *EvtForward) Seek(off int64, whence int) (int64, error) {
	return ss.generator.Seek(off, whence)
}

func NewEvtForwardGen(w http.ResponseWriter, r *http.Request, gen Generator, reqdat map[string]string, updated int64) Generator {

	var posevt int64 = math.MaxInt64

	posstr := reqdat[`posevt`]
	posfields := strings.Split(posstr, `.`)

	if 2 <= len(posfields) {
		if posevt = EvalNumber(posfields[0], 0); 0 < posevt {
			key := posfields[1]

			// This one doesn't depend on range requests
			if key == `close` {
				return &EvtForward{generator: gen, posevt: posevt}
			}

			// shift off the byte pos field and key
			posfields = posfields[2:]

			if rreq, ok := r.Header[`Range`]; ok {
				if sz, err := strconv.ParseInt(reqdat[`sz`], 10, 64); nil == err {
					if rangebe, err := parseRange(rreq[0], sz); nil == err {
						switch key {
						case `sc`:
							if rangebe.Contains(posevt) {
								code, _ := strconv.Atoi(posfields[0])
								DebugLogf("Sending posevt code %d\n", code)
								w.WriteHeader(code)
								return nil
							}
						case `etags`:
							if 2 <= len(posfields) {
								var etag string
								if rangebe.End < posevt {
									etag = posfields[0]
								} else {
									etag = posfields[1]
								}

								DebugLogf("Setting posevt etag %s\n", etag)
								w.Header().Set(`ETag`, etag)
								posevt = math.MaxInt64
							}
						}
					}
				}
			}
		}
	}

	return gen
}
