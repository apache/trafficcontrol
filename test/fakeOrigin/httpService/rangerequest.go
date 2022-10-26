package httpService

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
	"fmt"
	"strconv"
	"strings"
)

type httpRange struct {
	startOffset uint64
	length      uint64
}

const multipartSeparator = "3d6b6a416f9b5fakeOrigin"

func parseRange(rawRange string, size uint64) ([]httpRange, error) {
	var out []httpRange
	if rawRange == "" {
		return nil, nil
	}
	if !strings.HasPrefix(rawRange, "bytes=") {
		return nil, errors.New("invalid range")
	}
	for _, r := range strings.Split(strings.TrimPrefix(rawRange, "bytes="), ",") {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}
		if !strings.Contains(r, "-") {
			return nil, errors.New("invalid range")
		}
		offsets := strings.Split(r, "-")
		if len(offsets) > 2 || len(offsets) <= 0 {
			return nil, errors.New("invalid range")
		}
		start := strings.TrimSpace(offsets[0])
		startI, err := strconv.ParseUint(start, 10, 64)
		if err != nil || startI >= size {
			return nil, errors.New("invalid range")
		}
		end := strings.TrimSpace(offsets[1])
		var endI uint64
		if end != "" {
			endI, err = strconv.ParseUint(end, 10, 64)
			if err != nil || endI < startI {
				return nil, errors.New("invalid range")
			}
			// clip to size of body
			if endI >= size {
				endI = size - 1
			}
		}
		var iout httpRange
		if start == "" {
			// relative to end of file
			iout.startOffset = size - endI
			iout.length = size - iout.startOffset
		} else {
			iout.startOffset = startI
			if end == "" {
				// all to end of file
				iout.length = size - startI
			} else {
				// mid-range
				iout.length = endI - startI + 1
			}
		}
		out = append(out, iout)
	}

	return out, nil
}

func getContentRangeHeader(start, length, totalSize uint64) string {
	var end uint64
	if start+length > totalSize-1 {
		end = totalSize - 1
	} else {
		end = start + length
	}
	return fmt.Sprintf("bytes=%d-%d/%d", start, end, totalSize)
}

func clipToRange(ranges []httpRange, obody []byte, contentHeader string) ([]byte, map[string]string, error) {
	totalSize := uint64(len(obody))
	// need to use existing instead of appending since we have to lookup 206 time mismatches
	// this also means we have to start dealing with slice string values

	if len(ranges) == 0 {
		return nil, nil, errors.New("no ranges supplied")
	} else if len(ranges) == 1 {
		// single part ranges
		r := ranges[0]
		b := obody[r.startOffset : r.startOffset+r.length]
		// Update response code and other headers based on isTimeMatch
		return b, map[string]string{"Content-Range": getContentRangeHeader(r.startOffset, r.startOffset+r.length-1, totalSize)}, nil
	}
	// multipart
	var b []byte
	for i := range ranges {
		b = append(b, []byte("--"+multipartSeparator+"\n")...)
		b = append(b, []byte("Content-Type: "+contentHeader+"\n")...)
		b = append(b, []byte("Content-Range: "+getContentRangeHeader(ranges[i].startOffset, ranges[i].startOffset+ranges[i].length-1, totalSize)+"\n\n")...)
		b = append(b, obody[ranges[i].startOffset:ranges[i].startOffset+ranges[i].length]...)
		b = append(b, []byte("\n")...)
	}
	b = append(b, []byte("--"+multipartSeparator+"--\n")...)
	//w.Header().Set("Content-Type", "multipart/byteranges; boundary="+multipartSeparator)
	return b, map[string]string{"Content-Type": "multipart/byteranges; boundary=" + multipartSeparator}, nil
}
