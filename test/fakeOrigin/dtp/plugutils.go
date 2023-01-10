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
)

func ReadBlock(
	pbuf []byte,
	at int64,
	max int64,
	blk_size int64,
	blk_gen func(ct int64) []byte,
) (n int, err error) {

	sz := len(pbuf)
	if (max - at) < int64(sz) {
		sz = int(max - at)
	}
	DebugLogf("Reading at %d/%d%%%d for %d.\n", at, max, blk_size, sz)
	if sz == 0 {
		return 0, io.EOF
	}

	pidx := 0

	blk_ct := at / blk_size
	blk_off := at % blk_size

	for {
		blk_buf := blk_gen(blk_ct)
		blk_len := int64(len(blk_buf))

		DebugLogf("Got block %d: %v\n", blk_ct, blk_buf)

		DebugLogf("Using block %d: %v\n", blk_ct, blk_buf)
		DebugLogf("Copying into buffer at %d:%d from %d:%d.\n", pidx, sz, blk_off, blk_len)
		ncopied := copy(pbuf[pidx:sz], blk_buf[blk_off:])
		DebugLogf("Copied %d bytes into buffer at %d from %d.\n",
			ncopied, pidx, blk_off)
		if ncopied == 0 {
			break
		}
		pidx += ncopied
		DebugLogf("Now at %d/%d\n", pidx, sz)
		if sz <= pidx {
			break
		}
		blk_off = 0 /* Always zero after the first partial block */
		blk_ct++
	}
	DebugLogf("%d/%d written.\n", at+int64(sz), max)
	return sz, nil
}

func NewSeekPosFor(
	off int64,
	whence int,
	posold int64,
	size int64,
) (int64, error) {

	var posnew int64
	switch whence {
	case io.SeekStart:
		posnew = off
	case io.SeekEnd:
		posnew = size - off
	case io.SeekCurrent:
		posnew = posold + off
		if size < posnew {
			posnew = size
		}
	}

	DebugLogf("Seeking to %d via %d: %d\n", off, whence, posnew)
	if posnew < 0 {
		return posold, errors.New(`unable to seek before file`)
	}
	return posnew, nil
}
