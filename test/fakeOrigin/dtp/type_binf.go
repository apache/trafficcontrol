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
)

func init() {
	GlobalGeneratorFuncs["binf"] = NewBinFastGen
}

type xorshift128plusstate struct {
	s0 uint64
	s1 uint64
}

// https://en.wikipedia.org/wiki/Xorshift
func (state *xorshift128plusstate) xorshift128plus() uint64 {
	xx := state.s0
	yy := state.s1
	state.s0 = yy
	xx ^= xx << 23                               // a
	state.s1 = xx ^ yy ^ (xx >> 17) ^ (yy >> 26) // b, c
	val := state.s1 + yy
	return val
}

func (state *xorshift128plusstate) fixup() {
	if state.s0 == 0 && state.s1 == 0 {
		state.s0 = 1
	}
}

type BinFastGen struct {
	Size       int64
	Pos        int64
	Rnd        int64
	Blockbytes int64
}

func (bs *BinFastGen) ContentType() string {
	return "application/octet-stream"
}

func (bs *BinFastGen) Read(bufout []byte) (n int, err error) {

	maxbytes := len(bufout)
	posbeg := bs.Pos
	posend := posbeg + int64(maxbytes)
	if bs.Size < posend {
		posend = bs.Size
	}

	blocknum := posbeg / bs.Blockbytes
	posfile := blocknum * bs.Blockbytes

	lenout := posend - posbeg
	posout := posfile - posbeg

	var state xorshift128plusstate

	for posout < lenout {
		if (posfile % bs.Blockbytes) == 0 {
			blocknum = posfile / bs.Blockbytes
			state.s0 = uint64(blocknum + bs.Rnd)
			state.s1 = uint64(bs.Blockbytes)
			state.fixup()
		}

		nextnum := state.xorshift128plus()

		if posout <= -8 {
			posout += 8
		} else {
			for posn := 0; posn < 8 && posout < lenout; posn++ {
				if 0 <= posout {
					bufout[posout] = byte(nextnum & 0xff)
					bs.Pos++
				}
				nextnum >>= 8
				posout++
			}
		}
		posfile += 8
	}

	return int(lenout), nil
}

func (bs *BinFastGen) Seek(off int64, whence int) (int64, error) {
	posnew, err := NewSeekPosFor(off, whence, bs.Pos, bs.Size)
	if nil == err {
		bs.Pos = posnew
	}
	return bs.Pos, nil
}

func NewBinFastGen(reqdat map[string]string, lastmod int64) Generator {
	rnd, _ := strconv.ParseInt(reqdat[`rnd`], 10, 64)
	sz, _ := strconv.ParseInt(reqdat[`sz`], 10, 64)
	bs, _ := strconv.ParseInt(reqdat[`bs`], 10, 64)
	if bs == 0 || (bs%8) != 0 {
		bs = 1024
	}
	return &BinFastGen{Size: sz, Rnd: rnd, Blockbytes: bs}
}
