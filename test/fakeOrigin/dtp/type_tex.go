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
	"encoding/base64"
	"fmt"
	"image/color"
	"image/png"
	"io"
	"strconv"
	"strings"
	"sync"
)

// repeated buffer from a tilable 64x64 plasma png texture
const texturebytes = `
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
`

func texturePng() io.Reader {
	return base64.NewDecoder(base64.StdEncoding, strings.NewReader(texturebytes))
}

var Texture []byte

func init() {

	img, err := png.Decode(texturePng())
	if nil != err {
		fmt.Println("tex: error decoding texture", err)
		return
	}

	bounds := img.Bounds()
	size := bounds.Size()
	area := size.X * size.Y
	if 0 < area {
		Texture = make([]byte, area)

		tindex := 0
		for row := bounds.Min.Y; row < bounds.Max.Y; row++ {
			for col := bounds.Min.X; col < bounds.Max.X; col++ {
				pgray := img.At(row, col).(color.Gray)
				Texture[tindex] = byte(pgray.Y)
				tindex++
			}
		}

		GlobalGeneratorFuncs["tex"] = NewTexGen
	}
}

var TextureCache []byte
var TextureCap int
var TexMutex sync.Mutex

type TexGen struct {
	Size int64
	Pos  int64
	Rnd  int64
}

func (s *TexGen) ContentType() string {
	return "application/octet-stream"
}

func (s *TexGen) Read(pbuf []byte) (n int, err error) {
	plen := len(pbuf)
	var leftover int = int(s.Size - s.Pos)

	if leftover < plen {
		plen = leftover
	}

	// fix up the texture cache
	TexMutex.Lock()

	texlen := len(Texture)

	// we could be more clever about this lock but we
	// expect it to only be hit a couple of times
	// as 32k seems to be the observed max
	if TextureCap < plen {
		numtextures := plen / texlen
		if (plen % texlen) != 0 {
			numtextures++
		}

		// reset the texture capacity
		TextureCap = numtextures * texlen

		// pad the number required to handle any offset
		numtextures++

		// create new texture cache buffer
		tcsize := numtextures * texlen
		TextureCache = make([]byte, tcsize)

		for indexb := 0; indexb < tcsize; indexb += texlen {
			copy(TextureCache[indexb:indexb+texlen], Texture)
		}
	}

	cbuf := TextureCache

	TexMutex.Unlock()

	offset := ((s.Pos + s.Rnd) % int64(texlen))

	copy(pbuf, cbuf[offset:offset+int64(plen)])

	s.Pos += int64(plen)

	return plen, nil
}

func (s *TexGen) Seek(off int64, whence int) (int64, error) {
	posnew, err := NewSeekPosFor(off, whence, s.Pos, s.Size)
	if nil == err {
		s.Pos = posnew
	}
	return s.Pos, nil
}

func NewTexGen(reqdat map[string]string, latmod int64) Generator {
	seed, _ := strconv.ParseInt(reqdat[`rnd`], 10, 64)
	sz, _ := strconv.ParseInt(reqdat[`sz`], 10, 64)
	return &TexGen{Size: sz, Rnd: seed}
}
