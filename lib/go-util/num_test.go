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
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestJSONNameOrIDStr_UnmarshalJSON(t *testing.T) {
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

func TestJSONNameOrIDStr_MarshalJSON(t *testing.T) {
	type marshalTestStruct = struct {
		Val JSONNameOrIDStr `json:"val"`
	}
	type testCase = struct {
		input          marshalTestStruct
		expectedOutput string
		expectErr      bool
	}

	testName := "test"
	testID := 7

	testCases := []testCase{
		{
			input: marshalTestStruct{
				Val: JSONNameOrIDStr{
					Name: &testName,
				},
			},
			expectedOutput: `{"val":"` + testName + `"}`,
		},
		{
			input: marshalTestStruct{
				Val: JSONNameOrIDStr{
					ID: &testID,
				},
			},
			expectedOutput: fmt.Sprintf(`{"val":%d}`, testID),
		},
		{
			input: marshalTestStruct{
				Val: JSONNameOrIDStr{
					Name: &testName,
					ID:   &testID,
				},
			},
			expectedOutput: fmt.Sprintf(`{"val":%d}`, testID),
		},
		{
			input: marshalTestStruct{
				Val: JSONNameOrIDStr{},
			},
			expectedOutput: ``,
			expectErr:      true,
		},
	}

	for _, testCase := range testCases {
		bts, err := json.Marshal(testCase.input)
		if testCase.expectErr {
			if err == nil {
				t.Errorf("Expected marshaling %+v to produce an error", testCase.input)
			}
		} else if err != nil {
			t.Errorf("Unexpected error marshalling %+v: %v", testCase.input, err)
		} else if string(bts) != testCase.expectedOutput {
			t.Errorf("Incorrect marshal output, expected: '%s', got: '%s'", testCase.expectedOutput, string(bts))
		}
	}
}

func TestToNumeric(t *testing.T) {
	number := "34.59354233"
	val, success := ToNumeric("34.59354233")
	if !success {
		t.Errorf("expected ToNumeric to succeed for string %s", number)
	}
	if val != 34.59354233 {
		t.Errorf("expected ToNumeric to return %s, got %f", number, val)
	}
	_, success = ToNumeric("Not a number")
	if success {
		t.Error("expected ToNumeric to fail to convert a non-numeric string")
	}

	validInputs := []interface{}{
		uint8(12),
		uint16(12),
		uint32(12),
		uint64(12),
		int8(12),
		int16(12),
		int32(12),
		int64(12),
		float32(12),
		float64(12),
		int(12),
		uint(12),
	}
	for _, input := range validInputs {
		val, success = ToNumeric(input)
		if !success {
			t.Errorf("Expected to be able to convert %T to a float64", input)
		} else if val != float64(val) {
			t.Errorf("Incorrect conversion - went from %v (%T) to %v (float64)", input, input, val)
		}
	}

	_, success = ToNumeric(1 + 2i)
	if success {
		t.Error("Expected to be unable to convert complex numbers to a numeric")
	}
}

func TestJSONIntStr_UnmarshalJSON(t *testing.T) {
	type unmarshalTestStruct = struct {
		Test JSONIntStr `json:"test"`
	}

	data := []byte(`{"test": null}`)
	var uts unmarshalTestStruct
	if err := json.Unmarshal(data, &uts); err == nil {
		t.Error("Expected an error unmarshalling 'null' as an int or string int")
	}

	data = []byte(`{"test": []}`)
	if err := json.Unmarshal(data, &uts); err == nil {
		t.Error("Expected an error unmarshalling an array as an int or string int")
	}

	data = []byte(`{"test": ""}`)
	if err := json.Unmarshal(data, &uts); err == nil {
		t.Error("Expected an error unmarshalling an empty string as an int or string int")
	}

	data = []byte(`{"test": true}`)
	if err := json.Unmarshal(data, &uts); err == nil {
		t.Error("Expected an error unmarshalling a boolean as an int or string int")
	}

	data = []byte(`{"test": 12.1}`)
	if err := json.Unmarshal(data, &uts); err == nil {
		t.Error("Expected an error unmarshalling a floating-point number as an int or string int")
	}

	data = []byte(`{"test": "12.1"}`)
	if err := json.Unmarshal(data, &uts); err == nil {
		t.Error("Expected an error unmarshalling a string containing a floating-point number as an int or string int")
	}

	data = []byte(`{"test": 12}`)
	if err := json.Unmarshal(data, &uts); err != nil {
		t.Errorf("Unexpected error unmarshalling an integer as an int or string int: %v", err)
	}

	data = []byte(`{"test": "12"}`)
	if err := json.Unmarshal(data, &uts); err != nil {
		t.Errorf("Unexpected error unmarshalling a string containing an integer as an int or string int: %v", err)
	}
}

func ExampleJSONIntStr_ToInt64() {
	var a JSONIntStr = 5
	fmt.Printf("%d (%T)\n", a, a)
	fmt.Printf("%d (%T)\n", a.ToInt64(), a.ToInt64())
	// Output: 5 (util.JSONIntStr)
	// 5 (int64)
}

func ExampleJSONIntStr_String() {
	var a JSONIntStr = 5
	fmt.Println(a)
	// Output: 5
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

func TestHashInts(t *testing.T) {
	ints := []int{1, 3, 2}
	hash := HashInts(ints, false)
	hashAgain := HashInts(ints, false)
	if !bytes.Equal(hash, hashAgain) {
		t.Errorf("Expected hashing the same things to yield the same hash, got '%v' first and '%v' the second time", hash, hashAgain)
	}
	sortedHash := HashInts(ints, true)
	if bytes.Equal(hash, sortedHash) {
		t.Error("Expected hashing with sort first to yield a different hash than without the sort")
	}
}

func ExampleIntSliceToMap() {
	ints := []int{1, 2, 3}
	fmt.Printf("%+v", IntSliceToMap(ints))
	// Output: map[1:{} 2:{} 3:{}]
}
