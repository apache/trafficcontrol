package util

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at
   http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

// Package rfc contains functions implementing RFC 7234, 2616, and other RFCs.
// When changing functions, be sure they still conform to the corresponding RFC.
// When adding symbols, document the RFC and section they correspond to.

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestJSONNameOrIDStr(t *testing.T) {
	type testCase struct {
		input     string
		expected  JSONNameOrIDStr
		expectErr bool
	}

	testCases := []testCase{
		{
			`"foo"`,
			JSONNameOrIDStr{Name: StrPtr("foo")},
			false,
		},
		{
			`"1"`,
			JSONNameOrIDStr{ID: IntPtr(1)},
			false,
		},
		{
			`1`,
			JSONNameOrIDStr{ID: IntPtr(1)},
			false,
		},
		{
			`"-1"`,
			JSONNameOrIDStr{ID: IntPtr(-1)},
			false,
		},
		{
			`-1`,
			JSONNameOrIDStr{ID: IntPtr(-1)},
			false,
		},
		{
			`1.234`,
			JSONNameOrIDStr{},
			true,
		},
		{
			`false`,
			JSONNameOrIDStr{},
			true,
		},
	}

	for _, testCase := range testCases {
		actual := JSONNameOrIDStr{}
		err := json.Unmarshal([]byte(testCase.input), &actual)
		if testCase.expectErr && err == nil {
			t.Errorf("expected: err, actual: %+v", actual)
			continue
		} else if !testCase.expectErr && err != nil {
			t.Errorf("expected: nil error, actual: %v", err)
			continue
		}
		if !reflect.DeepEqual(testCase.expected, actual) {
			t.Errorf("expected: %+v, actual: %+v", testCase.expected, actual)
		}
	}
}

func TestToNumeric(t *testing.T) {
	var number interface{} = "34.59354233"
	val, success := ToNumeric(number)
	if !success {
		t.Errorf("expected ToNumeric to succeed for string %v", number)
	}
	if val != 34.59354233 {
		t.Errorf("expected ToNumeric to return %v, got %v", number, val)
	}
}

func TestBytesLenSplit(t *testing.T) {
	{
		b := []byte("abcdefghijklmnopqrstuvwxyz1234567890_-+=")
		n := 3
		expected := [][]byte{
			[]byte("abc"),
			[]byte("def"),
			[]byte("ghi"),
			[]byte("jkl"),
			[]byte("mno"),
			[]byte("pqr"),
			[]byte("stu"),
			[]byte("vwx"),
			[]byte("yz1"),
			[]byte("234"),
			[]byte("567"),
			[]byte("890"),
			[]byte("_-+"),
			[]byte("="),
		}
		actual := BytesLenSplit(b, n)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("BytesLenSplit expected: %+v actual: %+v\n", expected, actual)
		}
	}
	{
		b := []byte("abcdefghijklmnopqrstuvwxyz1234567890_-+=")
		n := 500
		expected := [][]byte{[]byte("abcdefghijklmnopqrstuvwxyz1234567890_-+=")}
		actual := BytesLenSplit(b, n)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("BytesLenSplit expected: %+v actual: %+v\n", expected, actual)
		}
	}
	{
		b := []byte("abcdefghijklmnopqrstuvwxyz1234567890_-+=")
		n := len(b) - 1
		expected := [][]byte{
			[]byte("abcdefghijklmnopqrstuvwxyz1234567890_-+"),
			[]byte("="),
		}
		actual := BytesLenSplit(b, n)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("BytesLenSplit expected: %+v actual: %+v\n", expected, actual)
		}
	}
	{
		b := []byte("abcdefghijklmnopqrstuvwxyz1234567890_-+=")
		n := len(b) - 2
		expected := [][]byte{
			[]byte("abcdefghijklmnopqrstuvwxyz1234567890_-"),
			[]byte("+="),
		}
		actual := BytesLenSplit(b, n)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("BytesLenSplit expected: %+v actual: %+v\n", expected, actual)
		}
	}
	{
		b := []byte("abcdefghijklmnopqrstuvwxyz1234567890_-+=")
		n := 20
		expected := [][]byte{
			[]byte("abcdefghijklmnopqrst"),
			[]byte("uvwxyz1234567890_-+="),
		}
		actual := BytesLenSplit(b, n)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("BytesLenSplit expected: %+v actual: %+v\n", expected, actual)
		}
	}
	{
		b := []byte("abcdefghijklmnopqrstuvwxyz1234567890_-+=")
		n := 0
		expected := [][]byte{}
		actual := BytesLenSplit(b, n)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("BytesLenSplit expected: %+v actual: %+v\n", expected, actual)
		}
	}
	{
		b := []byte("abcdefghijklmnopqrstuvwxyz1234567890_-+=")
		n := -1
		expected := [][]byte{}
		actual := BytesLenSplit(b, n)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("BytesLenSplit expected: %+v actual: %+v\n", expected, actual)
		}
	}
	{
		b := []byte("abcdefghijklmnopqrstuvwxyz1234567890_-+=")
		n := -30
		expected := [][]byte{}
		actual := BytesLenSplit(b, n)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("BytesLenSplit expected: %+v actual: %+v\n", expected, actual)
		}
	}
	{
		b := []byte("abcdefghijklmnopqrstuvwxyz1234567890_-+=")
		n := 1000000000
		expected := [][]byte{[]byte("abcdefghijklmnopqrstuvwxyz1234567890_-+=")}
		actual := BytesLenSplit(b, n)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("BytesLenSplit expected: %+v actual: %+v\n", expected, actual)
		}
	}
	{
		b := []byte("abcdefghijklmnopqrstuvwxyz1234567890_-+=")
		n := 1
		expected := [][]byte{
			[]byte("a"),
			[]byte("b"),
			[]byte("c"),
			[]byte("d"),
			[]byte("e"),
			[]byte("f"),
			[]byte("g"),
			[]byte("h"),
			[]byte("i"),
			[]byte("j"),
			[]byte("k"),
			[]byte("l"),
			[]byte("m"),
			[]byte("n"),
			[]byte("o"),
			[]byte("p"),
			[]byte("q"),
			[]byte("r"),
			[]byte("s"),
			[]byte("t"),
			[]byte("u"),
			[]byte("v"),
			[]byte("w"),
			[]byte("x"),
			[]byte("y"),
			[]byte("z"),
			[]byte("1"),
			[]byte("2"),
			[]byte("3"),
			[]byte("4"),
			[]byte("5"),
			[]byte("6"),
			[]byte("7"),
			[]byte("8"),
			[]byte("9"),
			[]byte("0"),
			[]byte("_"),
			[]byte("-"),
			[]byte("+"),
			[]byte("="),
		}
		actual := BytesLenSplit(b, n)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("BytesLenSplit expected: %+v actual: %+v\n", expected, actual)
		}
	}
}
