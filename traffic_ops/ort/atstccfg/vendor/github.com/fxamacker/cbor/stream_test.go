// Copyright (c) Faye Amacker. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root for license information.

package cbor

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"time"
)

func TestDecoder(t *testing.T) {
	var buf bytes.Buffer
	for i := 0; i < 5; i++ {
		for _, tc := range unmarshalTests {
			buf.Write(tc.cborData)
		}
	}
	decoder := NewDecoder(&buf)
	bytesRead := 0
	for i := 0; i < 5; i++ {
		for _, tc := range unmarshalTests {
			var v interface{}
			if err := decoder.Decode(&v); err != nil {
				t.Fatalf("Decode() returned error %v", err)
			}
			if tm, ok := tc.emptyInterfaceValue.(time.Time); ok {
				if vt, ok := v.(time.Time); !ok || !tm.Equal(vt) {
					t.Errorf("Decode() = %v (%T), want %v (%T)", v, v, tc.emptyInterfaceValue, tc.emptyInterfaceValue)
				}
			} else if !reflect.DeepEqual(v, tc.emptyInterfaceValue) {
				t.Errorf("Decode() = %v (%T), want %v (%T)", v, v, tc.emptyInterfaceValue, tc.emptyInterfaceValue)
			}
			bytesRead += len(tc.cborData)
			if decoder.NumBytesRead() != bytesRead {
				t.Errorf("NumBytesRead() = %v, want %v", decoder.NumBytesRead(), bytesRead)
			}
		}
	}
	// no more data
	var v interface{}
	err := decoder.Decode(&v)
	if v != nil {
		t.Errorf("Decode() = %v (%T), want nil (no more data)", v, v)
	}
	if err != io.EOF {
		t.Errorf("Decode() returned error %v, want io.EOF (no more data)", err)
	}
}

func TestDecoderUnmarshalTypeError(t *testing.T) {
	var buf bytes.Buffer
	for i := 0; i < 5; i++ {
		for _, tc := range unmarshalTests {
			for j := 0; j < len(tc.wrongTypes)*2; j++ {
				buf.Write(tc.cborData)
			}
		}
	}
	decoder := NewDecoder(&buf)
	bytesRead := 0
	for i := 0; i < 5; i++ {
		for _, tc := range unmarshalTests {
			for _, typ := range tc.wrongTypes {
				v := reflect.New(typ)
				if err := decoder.Decode(v.Interface()); err == nil {
					t.Errorf("Decode(0x%x) didn't return an error, want UnmarshalTypeError", tc.cborData)
				} else if _, ok := err.(*UnmarshalTypeError); !ok {
					t.Errorf("Decode(0x%x) returned wrong error type %T, want UnmarshalTypeError", tc.cborData, err)
				}
				bytesRead += len(tc.cborData)
				if decoder.NumBytesRead() != bytesRead {
					t.Errorf("NumBytesRead() = %v, want %v", decoder.NumBytesRead(), bytesRead)
				}

				var vi interface{}
				if err := decoder.Decode(&vi); err != nil {
					t.Errorf("Decode() returned error %v", err)
				}
				if tm, ok := tc.emptyInterfaceValue.(time.Time); ok {
					if vt, ok := vi.(time.Time); !ok || !tm.Equal(vt) {
						t.Errorf("Decode() = %v (%T), want %v (%T)", vi, vi, tc.emptyInterfaceValue, tc.emptyInterfaceValue)
					}
				} else if !reflect.DeepEqual(vi, tc.emptyInterfaceValue) {
					t.Errorf("Decode() = %v (%T), want %v (%T)", vi, vi, tc.emptyInterfaceValue, tc.emptyInterfaceValue)
				}
				bytesRead += len(tc.cborData)
				if decoder.NumBytesRead() != bytesRead {
					t.Errorf("NumBytesRead() = %v, want %v", decoder.NumBytesRead(), bytesRead)
				}
			}
		}
	}
	// no more data
	var v interface{}
	err := decoder.Decode(&v)
	if v != nil {
		t.Errorf("Decode() = %v (%T), want nil (no more data)", v, v)
	}
	if err != io.EOF {
		t.Errorf("Decode() returned error %v, want io.EOF (no more data)", err)
	}
}

func TestDecoderStructTag(t *testing.T) {
	type strc struct {
		A string `json:"x" cbor:"a"`
		B string `json:"y" cbor:"b"`
		C string `json:"z"`
	}
	want := strc{
		A: "A",
		B: "B",
		C: "C",
	}
	cborData := hexDecode("a36161614161626142617a6143") // {"a":"A", "b":"B", "z":"C"}

	var v strc
	dec := NewDecoder(bytes.NewReader(cborData))
	if err := dec.Decode(&v); err != nil {
		t.Errorf("Decode() returned error %v", err)
	}
	if !reflect.DeepEqual(v, want) {
		t.Errorf("Decode() = %+v (%T), want %+v (%T)", v, v, want, want)
	}
}

func TestEncoder(t *testing.T) {
	var want bytes.Buffer
	var w bytes.Buffer
	em, err := CanonicalEncOptions().EncMode()
	if err != nil {
		t.Errorf("EncMode() returned an error %v", err)
	}
	encoder := em.NewEncoder(&w)
	for _, tc := range marshalTests {
		for _, value := range tc.values {
			want.Write(tc.cborData)

			if err := encoder.Encode(value); err != nil {
				t.Fatalf("Encode() returned error %v", err)
			}
		}
	}
	if !bytes.Equal(w.Bytes(), want.Bytes()) {
		t.Errorf("Encoding mismatch: got %v, want %v", w.Bytes(), want.Bytes())
	}
}

func TestEncoderError(t *testing.T) {
	testcases := []struct {
		name         string
		value        interface{}
		wantErrorMsg string
	}{
		{"channel cannot be marshaled", make(chan bool), "cbor: unsupported type: chan bool"},
		{"function cannot be marshaled", func(i int) int { return i * i }, "cbor: unsupported type: func(int) int"},
		{"complex cannot be marshaled", complex(100, 8), "cbor: unsupported type: complex128"},
	}
	var w bytes.Buffer
	encoder := NewEncoder(&w)
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := encoder.Encode(&tc.value)
			if err == nil {
				t.Errorf("Encode(%v) didn't return an error, want error %q", tc.value, tc.wantErrorMsg)
			} else if _, ok := err.(*UnsupportedTypeError); !ok {
				t.Errorf("Encode(%v) error type %T, want *UnsupportedTypeError", tc.value, err)
			} else if err.Error() != tc.wantErrorMsg {
				t.Errorf("Encode(%v) error %q, want %q", tc.value, err.Error(), tc.wantErrorMsg)
			}
		})
	}
	if w.Len() != 0 {
		t.Errorf("Encoder's writer has %d bytes of data, want empty data", w.Len())
	}
}

func TestIndefiniteByteString(t *testing.T) {
	want := hexDecode("5f42010243030405ff")
	var w bytes.Buffer
	encoder := NewEncoder(&w)
	if err := encoder.StartIndefiniteByteString(); err != nil {
		t.Fatalf("StartIndefiniteByteString() returned error %v", err)
	}
	if err := encoder.Encode([]byte{1, 2}); err != nil {
		t.Fatalf("Encode() returned error %v", err)
	}
	if err := encoder.Encode([3]byte{3, 4, 5}); err != nil {
		t.Fatalf("Encode() returned error %v", err)
	}
	if err := encoder.EndIndefinite(); err != nil {
		t.Fatalf("EndIndefinite() returned error %v", err)
	}
	if !bytes.Equal(w.Bytes(), want) {
		t.Errorf("Encoding mismatch: got %v, want %v", w.Bytes(), want)
	}
}

func TestIndefiniteByteStringError(t *testing.T) {
	var w bytes.Buffer
	encoder := NewEncoder(&w)
	if err := encoder.StartIndefiniteByteString(); err != nil {
		t.Fatalf("StartIndefiniteByteString() returned error %v", err)
	}
	if err := encoder.Encode([]int{1, 2}); err == nil {
		t.Errorf("Encode() didn't return an error")
	} else if err.Error() != "cbor: cannot encode item type slice for indefinite-length byte string" {
		t.Errorf("Encode() returned error %q, want %q", err.Error(), "cbor: cannot encode item type slice for indefinite-length byte string")
	}
	if err := encoder.Encode("hello"); err == nil {
		t.Errorf("Encode() didn't return an error")
	} else if err.Error() != "cbor: cannot encode item type string for indefinite-length byte string" {
		t.Errorf("Encode() returned error %q, want %q", err.Error(), "cbor: cannot encode item type string for indefinite-length byte string")
	}
}

func TestIndefiniteTextString(t *testing.T) {
	want := hexDecode("7f657374726561646d696e67ff")
	var w bytes.Buffer
	encoder := NewEncoder(&w)
	if err := encoder.StartIndefiniteTextString(); err != nil {
		t.Fatalf("StartIndefiniteTextString() returned error %v", err)
	}
	if err := encoder.Encode("strea"); err != nil {
		t.Fatalf("Encode() returned error %v", err)
	}
	if err := encoder.Encode("ming"); err != nil {
		t.Fatalf("Encode() returned error %v", err)
	}
	if err := encoder.EndIndefinite(); err != nil {
		t.Fatalf("EndIndefinite() returned error %v", err)
	}
	if !bytes.Equal(w.Bytes(), want) {
		t.Errorf("Encoding mismatch: got %v, want %v", w.Bytes(), want)
	}
}

func TestIndefiniteTextStringError(t *testing.T) {
	var w bytes.Buffer
	encoder := NewEncoder(&w)
	if err := encoder.StartIndefiniteTextString(); err != nil {
		t.Fatalf("StartIndefiniteTextString() returned error %v", err)
	}
	if err := encoder.Encode([]byte{1, 2}); err == nil {
		t.Errorf("Encode() didn't return an error")
	} else if err.Error() != "cbor: cannot encode item type slice for indefinite-length text string" {
		t.Errorf("Encode() returned error %q, want %q", err.Error(), "cbor: cannot encode item type slice for indefinite-length text string")
	}
}

func TestIndefiniteArray(t *testing.T) {
	want := hexDecode("9f018202039f0405ffff")
	var w bytes.Buffer
	encoder := NewEncoder(&w)
	if err := encoder.StartIndefiniteArray(); err != nil {
		t.Fatalf("StartIndefiniteArray() returned error %v", err)
	}
	if err := encoder.Encode(1); err != nil {
		t.Fatalf("Encode() returned error %v", err)
	}
	if err := encoder.Encode([]int{2, 3}); err != nil {
		t.Fatalf("Encode() returned error %v", err)
	}
	if err := encoder.StartIndefiniteArray(); err != nil {
		t.Fatalf("StartIndefiniteArray() returned error %v", err)
	}
	if err := encoder.Encode(4); err != nil {
		t.Fatalf("Encode() returned error %v", err)
	}
	if err := encoder.Encode(5); err != nil {
		t.Fatalf("Encode() returned error %v", err)
	}
	if err := encoder.EndIndefinite(); err != nil {
		t.Fatalf("EndIndefinite() returned error %v", err)
	}
	if err := encoder.EndIndefinite(); err != nil {
		t.Fatalf("EndIndefinite() returned error %v", err)
	}
	if !bytes.Equal(w.Bytes(), want) {
		t.Errorf("Encoding mismatch: got %v, want %v", w.Bytes(), want)
	}
}

func TestIndefiniteMap(t *testing.T) {
	want := hexDecode("bf61610161629f0203ffff")
	var w bytes.Buffer
	em, err := EncOptions{Sort: SortCanonical}.EncMode()
	if err != nil {
		t.Errorf("EncMode() returned an error %v", err)
	}
	encoder := em.NewEncoder(&w)
	if err := encoder.StartIndefiniteMap(); err != nil {
		t.Fatalf("StartIndefiniteMap() returned error %v", err)
	}
	if err := encoder.Encode("a"); err != nil {
		t.Fatalf("Encode() returned error %v", err)
	}
	if err := encoder.Encode(1); err != nil {
		t.Fatalf("Encode() returned error %v", err)
	}
	if err := encoder.Encode("b"); err != nil {
		t.Fatalf("Encode() returned error %v", err)
	}
	if err := encoder.StartIndefiniteArray(); err != nil {
		t.Fatalf("StartIndefiniteArray() returned error %v", err)
	}
	if err := encoder.Encode(2); err != nil {
		t.Fatalf("Encode() returned error %v", err)
	}
	if err := encoder.Encode(3); err != nil {
		t.Fatalf("Encode() returned error %v", err)
	}
	if err := encoder.EndIndefinite(); err != nil {
		t.Fatalf("EndIndefinite() returned error %v", err)
	}
	if err := encoder.EndIndefinite(); err != nil {
		t.Fatalf("EndIndefinite() returned error %v", err)
	}
	if !bytes.Equal(w.Bytes(), want) {
		t.Errorf("Encoding mismatch: got %v, want %v", w.Bytes(), want)
	}
}

func TestIndefiniteLengthError(t *testing.T) {
	var w bytes.Buffer
	encoder := NewEncoder(&w)
	if err := encoder.StartIndefiniteByteString(); err != nil {
		t.Fatalf("StartIndefiniteByteString() returned error %v", err)
	}
	if err := encoder.EndIndefinite(); err != nil {
		t.Fatalf("EndIndefinite() returned error %v", err)
	}
	if err := encoder.EndIndefinite(); err == nil {
		t.Fatalf("EndIndefinite() didn't return an error")
	}
}

func TestEncoderStructTag(t *testing.T) {
	type strc struct {
		A string `json:"x" cbor:"a"`
		B string `json:"y" cbor:"b"`
		C string `json:"z"`
	}
	v := strc{
		A: "A",
		B: "B",
		C: "C",
	}
	want := hexDecode("a36161614161626142617a6143") // {"a":"A", "b":"B", "z":"C"}

	var w bytes.Buffer
	encoder := NewEncoder(&w)
	if err := encoder.Encode(v); err != nil {
		t.Errorf("Encode(%+v) returned error %v", v, err)
	}
	if !bytes.Equal(w.Bytes(), want) {
		t.Errorf("Encoding mismatch: got %v, want %v", w.Bytes(), want)
	}
}

func TestRawMessage(t *testing.T) {
	type strc struct {
		A RawMessage  `cbor:"a"`
		B *RawMessage `cbor:"b"`
		C *RawMessage `cbor:"c"`
	}
	cborData := hexDecode("a361610161628202036163f6") // {"a": 1, "b": [2, 3], "c": nil},
	r := RawMessage(hexDecode("820203"))
	want := strc{
		A: RawMessage([]byte{0x01}),
		B: &r,
	}
	var v strc
	if err := Unmarshal(cborData, &v); err != nil {
		t.Fatalf("Unmarshal(0x%x) returned error %v", cborData, err)
	}
	if !reflect.DeepEqual(v, want) {
		t.Errorf("Unmarshal(0x%x) returned v %v, want %v", cborData, v, want)
	}
	b, err := Marshal(v)
	if err != nil {
		t.Fatalf("Marshal(%+v) returned error %v", v, err)
	}
	if !bytes.Equal(b, cborData) {
		t.Errorf("Marshal(%+v) = 0x%x, want 0x%x", v, b, cborData)
	}
}

func TestNullRawMessage(t *testing.T) {
	r := RawMessage(nil)
	wantCborData := []byte{0xf6}
	b, err := Marshal(r)
	if err != nil {
		t.Errorf("Marshal(%+v) returned error %v", r, err)
	}
	if !bytes.Equal(b, wantCborData) {
		t.Errorf("Marshal(%+v) = 0x%x, want 0x%x", r, b, wantCborData)
	}
}

func TestEmptyRawMessage(t *testing.T) {
	var r RawMessage
	wantCborData := []byte{0xf6}
	b, err := Marshal(r)
	if err != nil {
		t.Errorf("Marshal(%+v) returned error %v", r, err)
	}
	if !bytes.Equal(b, wantCborData) {
		t.Errorf("Marshal(%+v) = 0x%x, want 0x%x", r, b, wantCborData)
	}
}

func TestNilRawMessageUnmarshalCBORError(t *testing.T) {
	wantErrorMsg := "cbor.RawMessage: UnmarshalCBOR on nil pointer"
	var r *RawMessage
	cborData := hexDecode("01")
	if err := r.UnmarshalCBOR(cborData); err == nil {
		t.Errorf("UnmarshalCBOR() didn't return error")
	} else if err.Error() != wantErrorMsg {
		t.Errorf("UnmarshalCBOR() returned error %q, want %q", err.Error(), wantErrorMsg)
	}
}
