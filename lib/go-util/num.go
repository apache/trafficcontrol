package util

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
	"crypto/sha512"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
)

// MSPerNS is the number of *nanoseconds* in a *millisecond* - not the other
// way around, as the name might imply.
const MSPerNS = int64(1000000)

// ToNumeric returns a float for any numeric type, and false if the interface
// does not hold a numeric type. This allows converting unknown numeric types
// (for example, from JSON) in a single line.
//
// TODO try to parse string stats as numbers? Also, JSON numbers are defined by
// the JSON spec to be IEEE double-precision floating point numbers and as such
// the encoding/json package always decodes them as float64s before coercing to
// requested types, so this may not be needed.
func ToNumeric(v interface{}) (float64, bool) {
	switch i := v.(type) {
	case uint8:
		return float64(i), true
	case uint16:
		return float64(i), true
	case uint32:
		return float64(i), true
	case uint64:
		return float64(i), true
	case int8:
		return float64(i), true
	case int16:
		return float64(i), true
	case int32:
		return float64(i), true
	case int64:
		return float64(i), true
	case float32:
		return float64(i), true
	case float64:
		return i, true
	case int:
		return float64(i), true
	case uint:
		return float64(i), true
	case string:
		n, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
		if err != nil {
			return 0.0, false
		}
		return n, true
	default:
		return 0.0, false
	}
}

// JSONIntStr unmarshals JSON strings or numbers into an int.
// This is designed to handle backwards-compatibility for old Perl endpoints which accept both. Please do not use this for new endpoints or new APIs, APIs should be well-typed.
type JSONIntStr int64

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (i *JSONIntStr) UnmarshalJSON(d []byte) error {
	if len(d) == 0 {
		return errors.New("empty object")
	}
	if d[0] == '"' {
		d = d[1 : len(d)-1] // strip JSON quotes
	}
	err := error(nil)
	di, err := strconv.ParseInt(string(d), 10, 64)
	if err != nil {
		return errors.New("not an integer")
	}
	*i = JSONIntStr(di)
	return nil
}

// ToInt64 returns the int64 value of the JSONIntStr.
func (i JSONIntStr) ToInt64() int64 {
	return int64(i)
}

// String implements the fmt.Stringer interface by returning the JSONIntStr
// encoded in base 10 into a string.
func (i JSONIntStr) String() string {
	return strconv.FormatInt(int64(i), 10)
}

// JSONNameOrIDStr is designed to handle backwards-compatibility for old Perl endpoints which accept both. Please do not use this for new endpoints or new APIs, APIs should be well-typed.
// NOTE: this differs from JSONIntStr in that this field could be 1 of 3 options:
//  1. string representing an integer
//  2. string representing a unique name
//  3. integer
type JSONNameOrIDStr struct {
	Name *string
	ID   *int
}

// MarshalJSON implements the encoding/json.Marshaler interface.
func (i JSONNameOrIDStr) MarshalJSON() ([]byte, error) {
	if i.ID != nil {
		return json.Marshal(*i.ID)
	}
	if i.Name != nil {
		return json.Marshal(*i.Name)
	}
	return nil, errors.New("either Name or ID must be non-nil")
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (i *JSONNameOrIDStr) UnmarshalJSON(d []byte) error {
	if len(d) == 0 {
		return errors.New("empty object")
	}
	quoted := false
	if d[0] == '"' {
		quoted = true
		d = d[1 : len(d)-1] // strip JSON quotes
	}
	di, err := strconv.ParseInt(string(d), 10, strconv.IntSize)
	if err != nil {
		if quoted {
			// if quoted, assume it is a name
			name := string(d)
			i.Name = &name
			return nil
		}
		return errors.New("expected an integer value: " + err.Error())
	}
	conv := int(di)
	i.ID = &conv
	return nil
}

// BytesLenSplit splits the given byte array into an n-length arrays. If n > len(s), returns a slice with a single []byte containing all of s. If n <= 0, returns an empty slice.
func BytesLenSplit(s []byte, n int) [][]byte {
	ss := [][]byte{}
	if n <= 0 {
		return ss
	}
	if n > len(s) {
		n = len(s)
	}
	for i := 0; i+n <= len(s); i += n {
		ss = append(ss, s[i:i+n])
	}
	rem := len(s) % n
	if rem != 0 {
		ss = append(ss, s[n*(len(s)/n):])
	}
	return ss
}

// HashInts returns a SHA512 hash of ints.
// If sortIntsBeforeHashing, the ints are sorted before before hashing. Sorting is done in a copy, the input ints slice is not modified.
func HashInts(ints []int, sortIntsBeforeHashing bool) []byte {
	sortedInts := ints
	if sortIntsBeforeHashing {
		sortedInts = make([]int, 0, len(ints))
		for _, in := range ints {
			sortedInts = append(sortedInts, in)
		}
		sort.Ints(sortedInts)
	}

	buf := make([]byte, binary.MaxVarintLen64*len(sortedInts))
	currBuf := buf
	for _, i := range sortedInts {
		n := binary.PutVarint(currBuf, int64(i))
		currBuf = currBuf[n:]
	}
	bts := sha512.Sum512(buf)
	return bts[:]
}

// IntSliceToMap creates an int set from an array.
func IntSliceToMap(s []int) map[int]struct{} {
	m := map[int]struct{}{}
	for _, v := range s {
		m[v] = struct{}{}
	}
	return m
}
