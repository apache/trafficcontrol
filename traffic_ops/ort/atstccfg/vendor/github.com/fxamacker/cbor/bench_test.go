// Copyright (c) Faye Amacker. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root for license information.

package cbor

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
)

const rounds = 100

type claims struct {
	Iss string `cbor:"1,keyasint"`
	Sub string `cbor:"2,keyasint"`
	Aud string `cbor:"3,keyasint"`
	Exp int    `cbor:"4,keyasint"`
	Nbf int    `cbor:"5,keyasint"`
	Iat int    `cbor:"6,keyasint"`
	Cti []byte `cbor:"7,keyasint"`
}

type coseHeader struct {
	Alg int    `cbor:"1,keyasint,omitempty"`
	Kid []byte `cbor:"4,keyasint,omitempty"`
	IV  []byte `cbor:"5,keyasint,omitempty"`
}

type macedCOSE struct {
	_           struct{} `cbor:",toarray"`
	Protected   []byte
	Unprotected coseHeader
	Payload     []byte
	Tag         []byte
}

type coseKey struct {
	Kty       int        `cbor:"1,keyasint,omitempty"`
	Kid       []byte     `cbor:"2,keyasint,omitempty"`
	Alg       int        `cbor:"3,keyasint,omitempty"`
	KeyOpts   int        `cbor:"4,keyasint,omitempty"`
	IV        []byte     `cbor:"5,keyasint,omitempty"`
	CrvOrNOrK RawMessage `cbor:"-1,keyasint,omitempty"` // K for symmetric keys, Crv for elliptic curve keys, N for RSA modulus
	XOrE      RawMessage `cbor:"-2,keyasint,omitempty"` // X for curve x-coordinate, E for RSA public exponent
	Y         RawMessage `cbor:"-3,keyasint,omitempty"` // Y for curve y-cooridate
	D         []byte     `cbor:"-4,keyasint,omitempty"`
}

type attestationObject struct {
	AuthnData []byte     `cbor:"authData"`
	Fmt       string     `cbor:"fmt"`
	AttStmt   RawMessage `cbor:"attStmt"`
}

type SenMLRecord struct {
	BaseName    string  `cbor:"-2,keyasint,omitempty"`
	BaseTime    float64 `cbor:"-3,keyasint,omitempty"`
	BaseUnit    string  `cbor:"-4,keyasint,omitempty"`
	BaseValue   float64 `cbor:"-5,keyasint,omitempty"`
	BaseSum     float64 `cbor:"-6,keyasint,omitempty"`
	BaseVersion int     `cbor:"-1,keyasint,omitempty"`
	Name        string  `cbor:"0,keyasint,omitempty"`
	Unit        string  `cbor:"1,keyasint,omitempty"`
	Time        float64 `cbor:"6,keyasint,omitempty"`
	UpdateTime  float64 `cbor:"7,keyasint,omitempty"`
	Value       float64 `cbor:"2,keyasint,omitempty"`
	ValueS      string  `cbor:"3,keyasint,omitempty"`
	ValueB      bool    `cbor:"4,keyasint,omitempty"`
	ValueD      string  `cbor:"8,keyasint,omitempty"`
	Sum         float64 `cbor:"5,keyasint,omitempty"`
}

type T1 struct {
	T    bool
	UI   uint
	I    int
	F    float64
	B    []byte
	S    string
	Slci []int
	Mss  map[string]string
}

type T2 struct {
	T    bool              `cbor:"1,keyasint"`
	UI   uint              `cbor:"2,keyasint"`
	I    int               `cbor:"3,keyasint"`
	F    float64           `cbor:"4,keyasint"`
	B    []byte            `cbor:"5,keyasint"`
	S    string            `cbor:"6,keyasint"`
	Slci []int             `cbor:"7,keyasint"`
	Mss  map[string]string `cbor:"8,keyasint"`
}

type T3 struct {
	_    struct{} `cbor:",toarray"`
	T    bool
	UI   uint
	I    int
	F    float64
	B    []byte
	S    string
	Slci []int
	Mss  map[string]string
}

var decodeBenchmarks = []struct {
	name          string
	cborData      []byte
	decodeToTypes []reflect.Type
}{
	{"bool", hexDecode("f5"), []reflect.Type{typeIntf, typeBool}},                                                                                                                                                                                           // true
	{"positive int", hexDecode("1bffffffffffffffff"), []reflect.Type{typeIntf, typeUint64}},                                                                                                                                                                 // uint64(18446744073709551615)
	{"negative int", hexDecode("3903e7"), []reflect.Type{typeIntf, typeInt64}},                                                                                                                                                                              // int64(-1000)
	{"float", hexDecode("fbc010666666666666"), []reflect.Type{typeIntf, typeFloat64}},                                                                                                                                                                       // float64(-4.1)
	{"bytes", hexDecode("581a0102030405060708090a0b0c0d0e0f101112131415161718191a"), []reflect.Type{typeIntf, typeByteSlice}},                                                                                                                               // []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26}
	{"bytes indef len", hexDecode("5f410141024103410441054106410741084109410a410b410c410d410e410f4110411141124113411441154116411741184119411aff"), []reflect.Type{typeIntf, typeByteSlice}},                                                                 // []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26}
	{"text", hexDecode("782b54686520717569636b2062726f776e20666f78206a756d7073206f76657220746865206c617a7920646f67"), []reflect.Type{typeIntf, typeString}},                                                                                                 // "The quick brown fox jumps over the lazy dog"
	{"text indef len", hexDecode("7f61546168616561206171617561696163616b612061626172616f6177616e61206166616f61786120616a6175616d617061736120616f61766165617261206174616861656120616c6161617a617961206164616f6167ff"), []reflect.Type{typeIntf, typeString}}, // "The quick brown fox jumps over the lazy dog"
	{"array", hexDecode("981a0102030405060708090a0b0c0d0e0f101112131415161718181819181a"), []reflect.Type{typeIntf, typeIntSlice}},                                                                                                                          // []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26}
	{"array indef len", hexDecode("9f0102030405060708090a0b0c0d0e0f101112131415161718181819181aff"), []reflect.Type{typeIntf, typeIntSlice}},                                                                                                                // []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26}
	{"map", hexDecode("ad616161416162614261636143616461446165614561666146616761476168614861696149616a614a616c614c616d614d616e614e"), []reflect.Type{typeIntf, typeMapStringIntf, typeMapStringString}},                                                      // map[string]string{"a": "A", "b": "B", "c": "C", "d": "D", "e": "E", "f": "F", "g": "G", "h": "H", "i": "I", "j": "J", "l": "L", "m": "M", "n": "N"}}
	{"map indef len", hexDecode("bf616161416162614261636143616461446165614561666146616761476168614861696149616a614a616b614b616c614c616d614d616e614eff"), []reflect.Type{typeIntf, typeMapStringIntf, typeMapStringString}},                                  // map[string]string{"a": "A", "b": "B", "c": "C", "d": "D", "e": "E", "f": "F", "g": "G", "h": "H", "i": "I", "j": "J", "l": "L", "m": "M", "n": "N"}}
}

var encodeBenchmarks = []struct {
	name     string
	cborData []byte
	values   []interface{}
}{
	{"bool", hexDecode("f5"), []interface{}{true}},
	{"positive int", hexDecode("1bffffffffffffffff"), []interface{}{uint64(18446744073709551615)}},
	{"negative int", hexDecode("3903e7"), []interface{}{int64(-1000)}},
	{"float", hexDecode("fbc010666666666666"), []interface{}{float64(-4.1)}},
	{"bytes", hexDecode("581a0102030405060708090a0b0c0d0e0f101112131415161718191a"), []interface{}{[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26}}},
	{"text", hexDecode("782b54686520717569636b2062726f776e20666f78206a756d7073206f76657220746865206c617a7920646f67"), []interface{}{"The quick brown fox jumps over the lazy dog"}},
	{"array", hexDecode("981a0102030405060708090a0b0c0d0e0f101112131415161718181819181a"), []interface{}{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26}}},
	{"map", hexDecode("ad616161416162614261636143616461446165614561666146616761476168614861696149616a614a616c614c616d614d616e614e"), []interface{}{map[string]string{"a": "A", "b": "B", "c": "C", "d": "D", "e": "E", "f": "F", "g": "G", "h": "H", "i": "I", "j": "J", "l": "L", "m": "M", "n": "N"}}},
}

func BenchmarkUnmarshal(b *testing.B) {
	for _, bm := range decodeBenchmarks {
		for _, t := range bm.decodeToTypes {
			name := "CBOR " + bm.name + " to Go " + t.String()
			if t.Kind() == reflect.Struct {
				name = "CBOR " + bm.name + " to Go " + t.Kind().String()
			}
			b.Run(name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					vPtr := reflect.New(t).Interface()
					if err := Unmarshal(bm.cborData, vPtr); err != nil {
						b.Fatal("Unmarshal:", err)
					}
				}
			})
		}
	}
	var moreBenchmarks = []struct {
		name         string
		cborData     []byte
		decodeToType reflect.Type
	}{
		// Unmarshal CBOR map with string key to map[string]interface{}.
		{
			"CBOR map to Go map[string]interface{}",
			hexDecode("a86154f56255691bffffffffffffffff61493903e76146fbc0106666666666666142581a0102030405060708090a0b0c0d0e0f101112131415161718191a6153782b54686520717569636b2062726f776e20666f78206a756d7073206f76657220746865206c617a7920646f6764536c6369981a0102030405060708090a0b0c0d0e0f101112131415161718181819181a634d7373ad6163614361656145616661466167614761686148616e614e616d614d61616141616261426164614461696149616a614a616c614c"),
			reflect.TypeOf(map[string]interface{}{}),
		},
		// Unmarshal CBOR map with string key to struct.
		{
			"CBOR map to Go struct",
			hexDecode("a86154f56255491bffffffffffffffff61493903e76146fbc0106666666666666142581a0102030405060708090a0b0c0d0e0f101112131415161718191a6153782b54686520717569636b2062726f776e20666f78206a756d7073206f76657220746865206c617a7920646f6764536c6369981a0102030405060708090a0b0c0d0e0f101112131415161718181819181a634d7373ad6163614361656145616661466167614761686148616e614e616d614d61616141616261426164614461696149616a614a616c614c"),
			reflect.TypeOf(T1{}),
		},
		// Unmarshal CBOR map with integer key, such as COSE Key and SenML, to map[int]interface{}.
		{
			"CBOR map to Go map[int]interface{}",
			hexDecode("a801f5021bffffffffffffffff033903e704fbc01066666666666605581a0102030405060708090a0b0c0d0e0f101112131415161718191a06782b54686520717569636b2062726f776e20666f78206a756d7073206f76657220746865206c617a7920646f6707981a0102030405060708090a0b0c0d0e0f101112131415161718181819181a08ad61646144616661466167614761686148616d614d616e614e6161614161626142616361436165614561696149616a614a616c614c"),
			reflect.TypeOf(map[int]interface{}{}),
		},
		// Unmarshal CBOR map with integer key, such as COSE Key and SenML, to struct.
		{
			"CBOR map to Go struct keyasint",
			hexDecode("a801f5021bffffffffffffffff033903e704fbc01066666666666605581a0102030405060708090a0b0c0d0e0f101112131415161718191a06782b54686520717569636b2062726f776e20666f78206a756d7073206f76657220746865206c617a7920646f6707981a0102030405060708090a0b0c0d0e0f101112131415161718181819181a08ad61646144616661466167614761686148616d614d616e614e6161614161626142616361436165614561696149616a614a616c614c"),
			reflect.TypeOf(T2{}),
		},
		// Unmarshal CBOR array of known sequence of data types, such as signed/maced/encrypted CWT, to []interface{}.
		{
			"CBOR array to Go []interface{}",
			hexDecode("88f51bffffffffffffffff3903e7fbc010666666666666581a0102030405060708090a0b0c0d0e0f101112131415161718191a782b54686520717569636b2062726f776e20666f78206a756d7073206f76657220746865206c617a7920646f67981a0102030405060708090a0b0c0d0e0f101112131415161718181819181aad616261426163614361646144616561456166614661696149616e614e616161416167614761686148616a614a616c614c616d614d"),
			reflect.TypeOf([]interface{}{}),
		},
		// Unmarshal CBOR array of known sequence of data types, such as signed/maced/encrypted CWT, to struct.
		{
			"CBOR array to Go struct toarray",
			hexDecode("88f51bffffffffffffffff3903e7fbc010666666666666581a0102030405060708090a0b0c0d0e0f101112131415161718191a782b54686520717569636b2062726f776e20666f78206a756d7073206f76657220746865206c617a7920646f67981a0102030405060708090a0b0c0d0e0f101112131415161718181819181aad616261426163614361646144616561456166614661696149616e614e616161416167614761686148616a614a616c614c616d614d"),
			reflect.TypeOf(T3{}),
		},
	}
	for _, bm := range moreBenchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				vPtr := reflect.New(bm.decodeToType).Interface()
				if err := Unmarshal(bm.cborData, vPtr); err != nil {
					b.Fatal("Unmarshal:", err)
				}
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	for _, bm := range decodeBenchmarks {
		for _, t := range bm.decodeToTypes {
			name := "CBOR " + bm.name + " to Go " + t.String()
			if t.Kind() == reflect.Struct {
				name = "CBOR " + bm.name + " to Go " + t.Kind().String()
			}
			buf := bytes.NewReader(bm.cborData)
			decoder := NewDecoder(buf)
			b.Run(name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					vPtr := reflect.New(t).Interface()
					if err := decoder.Decode(vPtr); err != nil {
						b.Fatal("Decode:", err)
					}
					buf.Seek(0, 0) //nolint:errcheck
				}
			})
		}
	}
}

func BenchmarkDecodeStream(b *testing.B) {
	var cborData []byte
	for _, bm := range decodeBenchmarks {
		for i := 0; i < len(bm.decodeToTypes); i++ {
			cborData = append(cborData, bm.cborData...)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bytes.NewReader(cborData)
		decoder := NewDecoder(buf)
		for j := 0; j < rounds; j++ {
			for _, bm := range decodeBenchmarks {
				for _, t := range bm.decodeToTypes {
					vPtr := reflect.New(t).Interface()
					if err := decoder.Decode(vPtr); err != nil {
						b.Fatal("Decode:", err)
					}
				}
			}
			buf.Seek(0, 0) //nolint:errcheck
		}
	}
}

func BenchmarkMarshal(b *testing.B) {
	for _, bm := range encodeBenchmarks {
		for _, v := range bm.values {
			name := "Go " + reflect.TypeOf(v).String() + " to CBOR " + bm.name
			if reflect.TypeOf(v).Kind() == reflect.Struct {
				name = "Go " + reflect.TypeOf(v).Kind().String() + " to CBOR " + bm.name
			}
			b.Run(name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					if _, err := Marshal(v); err != nil {
						b.Fatal("Marshal:", err)
					}
				}
			})
		}
	}
	// Marshal map[string]interface{} to CBOR map
	m1 := map[string]interface{}{ //nolint:dupl
		"T":    true,
		"UI":   uint(18446744073709551615),
		"I":    -1000,
		"F":    -4.1,
		"B":    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
		"S":    "The quick brown fox jumps over the lazy dog",
		"Slci": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
		"Mss":  map[string]string{"a": "A", "b": "B", "c": "C", "d": "D", "e": "E", "f": "F", "g": "G", "h": "H", "i": "I", "j": "J", "l": "L", "m": "M", "n": "N"},
	}
	// Marshal struct to CBOR map
	v1 := T1{ //nolint:dupl
		T:    true,
		UI:   18446744073709551615,
		I:    -1000,
		F:    -4.1,
		B:    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
		S:    "The quick brown fox jumps over the lazy dog",
		Slci: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
		Mss:  map[string]string{"a": "A", "b": "B", "c": "C", "d": "D", "e": "E", "f": "F", "g": "G", "h": "H", "i": "I", "j": "J", "l": "L", "m": "M", "n": "N"},
	}
	// Marshal map[int]interface{} to CBOR map
	m2 := map[int]interface{}{ //nolint:dupl
		1: true,
		2: uint(18446744073709551615),
		3: -1000,
		4: -4.1,
		5: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
		6: "The quick brown fox jumps over the lazy dog",
		7: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
		8: map[string]string{"a": "A", "b": "B", "c": "C", "d": "D", "e": "E", "f": "F", "g": "G", "h": "H", "i": "I", "j": "J", "l": "L", "m": "M", "n": "N"},
	}
	// Marshal struct keyasint, such as COSE Key and SenML
	v2 := T2{ //nolint:dupl
		T:    true,
		UI:   18446744073709551615,
		I:    -1000,
		F:    -4.1,
		B:    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
		S:    "The quick brown fox jumps over the lazy dog",
		Slci: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
		Mss:  map[string]string{"a": "A", "b": "B", "c": "C", "d": "D", "e": "E", "f": "F", "g": "G", "h": "H", "i": "I", "j": "J", "l": "L", "m": "M", "n": "N"},
	}
	// Marshal []interface to CBOR array.
	slc := []interface{}{
		true,
		uint(18446744073709551615),
		-1000,
		-4.1,
		[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
		"The quick brown fox jumps over the lazy dog",
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
		map[string]string{"a": "A", "b": "B", "c": "C", "d": "D", "e": "E", "f": "F", "g": "G", "h": "H", "i": "I", "j": "J", "l": "L", "m": "M", "n": "N"},
	}
	// Marshal struct toarray to CBOR array, such as signed/maced/encrypted CWT.
	v3 := T3{ //nolint:dupl
		T:    true,
		UI:   18446744073709551615,
		I:    -1000,
		F:    -4.1,
		B:    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
		S:    "The quick brown fox jumps over the lazy dog",
		Slci: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
		Mss:  map[string]string{"a": "A", "b": "B", "c": "C", "d": "D", "e": "E", "f": "F", "g": "G", "h": "H", "i": "I", "j": "J", "l": "L", "m": "M", "n": "N"},
	}
	var moreBenchmarks = []struct {
		name  string
		value interface{}
	}{
		{"Go map[string]interface{} to CBOR map", m1},
		{"Go struct to CBOR map", v1},
		{"Go map[int]interface{} to CBOR map", m2},
		{"Go struct keyasint to CBOR map", v2},
		{"Go []interface{} to CBOR map", slc},
		{"Go struct toarray to CBOR array", v3},
	}
	for _, bm := range moreBenchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if _, err := Marshal(bm.value); err != nil {
					b.Fatal("Marshal:", err)
				}
			}
		})
	}
}

func BenchmarkMarshalCanonical(b *testing.B) {
	type strc struct {
		A string `cbor:"a"`
		B string `cbor:"b"`
		C string `cbor:"c"`
		D string `cbor:"d"`
		E string `cbor:"e"`
		F string `cbor:"f"`
		G string `cbor:"g"`
		H string `cbor:"h"`
		I string `cbor:"i"`
		J string `cbor:"j"`
		L string `cbor:"l"`
		M string `cbor:"m"`
		N string `cbor:"n"`
	}
	for _, bm := range []struct {
		name     string
		cborData []byte
		values   []interface{}
	}{
		{"map", hexDecode("ad616161416162614261636143616461446165614561666146616761476168614861696149616a614a616c614c616d614d616e614e"), []interface{}{map[string]string{"a": "A", "b": "B", "c": "C", "d": "D", "e": "E", "f": "F", "g": "G", "h": "H", "i": "I", "j": "J", "l": "L", "m": "M", "n": "N"}, strc{A: "A", B: "B", C: "C", D: "D", E: "E", F: "F", G: "G", H: "H", I: "I", J: "J", L: "L", M: "M", N: "N"}}},
	} {
		for _, v := range bm.values {
			name := "Go " + reflect.TypeOf(v).String() + " to CBOR " + bm.name
			if reflect.TypeOf(v).Kind() == reflect.Struct {
				name = "Go " + reflect.TypeOf(v).Kind().String() + " to CBOR " + bm.name
			}
			b.Run(name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					if _, err := Marshal(v); err != nil {
						b.Fatal("Marshal:", err)
					}
				}
			})
			// Canonical encoding
			name = "Go " + reflect.TypeOf(v).String() + " to CBOR " + bm.name + " canonical"
			if reflect.TypeOf(v).Kind() == reflect.Struct {
				name = "Go " + reflect.TypeOf(v).Kind().String() + " to CBOR " + bm.name + " canonical"
			}
			em, _ := EncOptions{Sort: SortCanonical}.EncMode()
			b.Run(name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					if _, err := em.Marshal(v); err != nil {
						b.Fatal("Marshal:", err)
					}
				}
			})
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	for _, bm := range encodeBenchmarks {
		for _, v := range bm.values {
			name := "Go " + reflect.TypeOf(v).String() + " to CBOR " + bm.name
			if reflect.TypeOf(v).Kind() == reflect.Struct {
				name = "Go " + reflect.TypeOf(v).Kind().String() + " to CBOR " + bm.name
			}
			b.Run(name, func(b *testing.B) {
				encoder := NewEncoder(ioutil.Discard)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					if err := encoder.Encode(v); err != nil {
						b.Fatal("Encode:", err)
					}
				}
			})
		}
	}
}

func BenchmarkEncodeStream(b *testing.B) {
	for i := 0; i < b.N; i++ {
		encoder := NewEncoder(ioutil.Discard)
		for i := 0; i < rounds; i++ {
			for _, bm := range encodeBenchmarks {
				for _, v := range bm.values {
					if err := encoder.Encode(v); err != nil {
						b.Fatal("Encode:", err)
					}
				}
			}
		}
	}
}

func BenchmarkUnmarshalCOSE(b *testing.B) {
	// Data from https://tools.ietf.org/html/rfc8392#appendix-A section A.2
	testCases := []struct {
		name     string
		cborData []byte
	}{
		{"128-Bit Symmetric Key", hexDecode("a42050231f4c4d4d3051fdc2ec0a3851d5b3830104024c53796d6d6574726963313238030a")},
		{"256-Bit Symmetric Key", hexDecode("a4205820403697de87af64611c1d32a05dab0fe1fcb715a86ab435f1ec99192d795693880104024c53796d6d6574726963323536030a")},
		{"ECDSA P256 256-Bit Key", hexDecode("a72358206c1382765aec5358f117733d281c1c7bdc39884d04a45a1e6c67c858bc206c1922582060f7f1a780d8a783bfb7a2dd6b2796e8128dbbcef9d3d168db9529971a36e7b9215820143329cce7868e416927599cf65a34f3ce2ffda55a7eca69ed8919a394d42f0f2001010202524173796d6d657472696345434453413235360326")},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var v coseKey
				if err := Unmarshal(tc.cborData, &v); err != nil {
					b.Fatal("Unmarshal:", err)
				}
			}
		})
	}
}

func BenchmarkMarshalCOSE(b *testing.B) {
	// Data from https://tools.ietf.org/html/rfc8392#appendix-A section A.2
	testCases := []struct {
		name     string
		cborData []byte
	}{
		{"128-Bit Symmetric Key", hexDecode("a42050231f4c4d4d3051fdc2ec0a3851d5b3830104024c53796d6d6574726963313238030a")},
		{"256-Bit Symmetric Key", hexDecode("a4205820403697de87af64611c1d32a05dab0fe1fcb715a86ab435f1ec99192d795693880104024c53796d6d6574726963323536030a")},
		{"ECDSA P256 256-Bit Key", hexDecode("a72358206c1382765aec5358f117733d281c1c7bdc39884d04a45a1e6c67c858bc206c1922582060f7f1a780d8a783bfb7a2dd6b2796e8128dbbcef9d3d168db9529971a36e7b9215820143329cce7868e416927599cf65a34f3ce2ffda55a7eca69ed8919a394d42f0f2001010202524173796d6d657472696345434453413235360326")},
	}
	for _, tc := range testCases {
		var v coseKey
		if err := Unmarshal(tc.cborData, &v); err != nil {
			b.Fatal("Unmarshal:", err)
		}
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if _, err := Marshal(v); err != nil {
					b.Fatal("Marshal:", err)
				}
			}
		})
	}
}

func BenchmarkUnmarshalCWTClaims(b *testing.B) {
	// Data from https://tools.ietf.org/html/rfc8392#appendix-A section A.1
	cborData := hexDecode("a70175636f61703a2f2f61732e6578616d706c652e636f6d02656572696b77037818636f61703a2f2f6c696768742e6578616d706c652e636f6d041a5612aeb0051a5610d9f0061a5610d9f007420b71")
	for i := 0; i < b.N; i++ {
		var v claims
		if err := Unmarshal(cborData, &v); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}

func BenchmarkUnmarshalCWTClaimsWithDupMapKeyOpt(b *testing.B) {
	// Data from https://tools.ietf.org/html/rfc8392#appendix-A section A.1
	cborData := hexDecode("a70175636f61703a2f2f61732e6578616d706c652e636f6d02656572696b77037818636f61703a2f2f6c696768742e6578616d706c652e636f6d041a5612aeb0051a5610d9f0061a5610d9f007420b71")
	dm, _ := DecOptions{DupMapKey: DupMapKeyEnforcedAPF}.DecMode()
	for i := 0; i < b.N; i++ {
		var v claims
		if err := dm.Unmarshal(cborData, &v); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}

func BenchmarkMarshalCWTClaims(b *testing.B) {
	// Data from https://tools.ietf.org/html/rfc8392#appendix-A section A.1
	cborData := hexDecode("a70175636f61703a2f2f61732e6578616d706c652e636f6d02656572696b77037818636f61703a2f2f6c696768742e6578616d706c652e636f6d041a5612aeb0051a5610d9f0061a5610d9f007420b71")
	var v claims
	if err := Unmarshal(cborData, &v); err != nil {
		b.Fatal("Unmarshal:", err)
	}
	for i := 0; i < b.N; i++ {
		if _, err := Marshal(v); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}

func BenchmarkUnmarshalSenML(b *testing.B) {
	// Data from https://tools.ietf.org/html/rfc8428#section-6
	cborData := hexDecode("87a721781b75726e3a6465763a6f773a3130653230373361303130383030363a22fb41d303a15b00106223614120050067766f6c7461676501615602fb405e066666666666a3006763757272656e74062402fb3ff3333333333333a3006763757272656e74062302fb3ff4cccccccccccda3006763757272656e74062202fb3ff6666666666666a3006763757272656e74062102f93e00a3006763757272656e74062002fb3ff999999999999aa3006763757272656e74060002fb3ffb333333333333")
	for i := 0; i < b.N; i++ {
		var v []SenMLRecord
		if err := Unmarshal(cborData, &v); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}

func BenchmarkMarshalSenML(b *testing.B) {
	// Data from https://tools.ietf.org/html/rfc8428#section-6
	cborData := hexDecode("87a721781b75726e3a6465763a6f773a3130653230373361303130383030363a22fb41d303a15b00106223614120050067766f6c7461676501615602fb405e066666666666a3006763757272656e74062402fb3ff3333333333333a3006763757272656e74062302fb3ff4cccccccccccda3006763757272656e74062202fb3ff6666666666666a3006763757272656e74062102f93e00a3006763757272656e74062002fb3ff999999999999aa3006763757272656e74060002fb3ffb333333333333")
	var v []SenMLRecord
	if err := Unmarshal(cborData, &v); err != nil {
		b.Fatal("Unmarshal:", err)
	}
	for i := 0; i < b.N; i++ {
		if _, err := Marshal(v); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}

func BenchmarkMarshalSenMLShortestFloat16(b *testing.B) {
	// Data from https://tools.ietf.org/html/rfc8428#section-6
	cborData := hexDecode("87a721781b75726e3a6465763a6f773a3130653230373361303130383030363a22fb41d303a15b00106223614120050067766f6c7461676501615602fb405e066666666666a3006763757272656e74062402fb3ff3333333333333a3006763757272656e74062302fb3ff4cccccccccccda3006763757272656e74062202fb3ff6666666666666a3006763757272656e74062102f93e00a3006763757272656e74062002fb3ff999999999999aa3006763757272656e74060002fb3ffb333333333333")
	var v []SenMLRecord
	if err := Unmarshal(cborData, &v); err != nil {
		b.Fatal("Unmarshal:", err)
	}
	em, _ := EncOptions{ShortestFloat: ShortestFloat16}.EncMode()
	for i := 0; i < b.N; i++ {
		if _, err := em.Marshal(v); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}

func BenchmarkUnmarshalWebAuthn(b *testing.B) {
	// Data generated from Yubico security key
	cborData := hexDecode("a363666d74686669646f2d7532666761747453746d74a26373696758483046022100e7ab373cfbd99fcd55fd59b0f6f17fef5b77a20ddec3db7f7e4d55174e366236022100828336b4822125fb56541fb14a8a273876acd339395ec2dad95cf41c1dd2a9ae637835638159024e3082024a30820132a0030201020204124a72fe300d06092a864886f70d01010b0500302e312c302a0603550403132359756269636f2055324620526f6f742043412053657269616c203435373230303633313020170d3134303830313030303030305a180f32303530303930343030303030305a302c312a302806035504030c2159756269636f205532462045452053657269616c203234393431343937323135383059301306072a8648ce3d020106082a8648ce3d030107034200043d8b1bbd2fcbf6086e107471601468484153c1c6d3b4b68a5e855e6e40757ee22bcd8988bf3befd7cdf21cb0bf5d7a150d844afe98103c6c6607d9faae287c02a33b3039302206092b0601040182c40a020415312e332e362e312e342e312e34313438322e312e313013060b2b0601040182e51c020101040403020520300d06092a864886f70d01010b05000382010100a14f1eea0076f6b8476a10a2be72e60d0271bb465b2dfbfc7c1bd12d351989917032631d795d097fa30a26a325634e85721bc2d01a86303f6bc075e5997319e122148b0496eec8d1f4f94cf4110de626c289443d1f0f5bbb239ca13e81d1d5aa9df5af8e36126475bfc23af06283157252762ff68879bcf0ef578d55d67f951b4f32b63c8aea5b0f99c67d7d814a7ff5a6f52df83e894a3a5d9c8b82e7f8bc8daf4c80175ff8972fda79333ec465d806eacc948f1bab22045a95558a48c20226dac003d41fbc9e05ea28a6bb5e10a49de060a0a4f6a2676a34d68c4abe8c61874355b9027e828ca9e064b002d62e8d8cf0744921753d35e3c87c5d5779453e7768617574684461746158c449960de5880e8c687434170f6476605b8fe4aeb9a28632c7995cf3ba831d976341000000000000000000000000000000000000000000408903fd7dfd2c9770e98cae0123b13a2c27828a106349bc6277140e7290b7e9eb7976aa3c04ed347027caf7da3a2fa76304751c02208acfc4e7fc6c7ebbc375c8a5010203262001215820ad7f7992c335b90d882b2802061b97a4fabca7e2ee3e7a51e728b8055e4eb9c7225820e0966ba7005987fece6f0e0e13447aa98cec248e4000a594b01b74c1cb1d40b3")
	for i := 0; i < b.N; i++ {
		var v attestationObject
		if err := Unmarshal(cborData, &v); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}

func BenchmarkMarshalWebAuthn(b *testing.B) {
	// Data generated from Yubico security key
	cborData := hexDecode("a363666d74686669646f2d7532666761747453746d74a26373696758483046022100e7ab373cfbd99fcd55fd59b0f6f17fef5b77a20ddec3db7f7e4d55174e366236022100828336b4822125fb56541fb14a8a273876acd339395ec2dad95cf41c1dd2a9ae637835638159024e3082024a30820132a0030201020204124a72fe300d06092a864886f70d01010b0500302e312c302a0603550403132359756269636f2055324620526f6f742043412053657269616c203435373230303633313020170d3134303830313030303030305a180f32303530303930343030303030305a302c312a302806035504030c2159756269636f205532462045452053657269616c203234393431343937323135383059301306072a8648ce3d020106082a8648ce3d030107034200043d8b1bbd2fcbf6086e107471601468484153c1c6d3b4b68a5e855e6e40757ee22bcd8988bf3befd7cdf21cb0bf5d7a150d844afe98103c6c6607d9faae287c02a33b3039302206092b0601040182c40a020415312e332e362e312e342e312e34313438322e312e313013060b2b0601040182e51c020101040403020520300d06092a864886f70d01010b05000382010100a14f1eea0076f6b8476a10a2be72e60d0271bb465b2dfbfc7c1bd12d351989917032631d795d097fa30a26a325634e85721bc2d01a86303f6bc075e5997319e122148b0496eec8d1f4f94cf4110de626c289443d1f0f5bbb239ca13e81d1d5aa9df5af8e36126475bfc23af06283157252762ff68879bcf0ef578d55d67f951b4f32b63c8aea5b0f99c67d7d814a7ff5a6f52df83e894a3a5d9c8b82e7f8bc8daf4c80175ff8972fda79333ec465d806eacc948f1bab22045a95558a48c20226dac003d41fbc9e05ea28a6bb5e10a49de060a0a4f6a2676a34d68c4abe8c61874355b9027e828ca9e064b002d62e8d8cf0744921753d35e3c87c5d5779453e7768617574684461746158c449960de5880e8c687434170f6476605b8fe4aeb9a28632c7995cf3ba831d976341000000000000000000000000000000000000000000408903fd7dfd2c9770e98cae0123b13a2c27828a106349bc6277140e7290b7e9eb7976aa3c04ed347027caf7da3a2fa76304751c02208acfc4e7fc6c7ebbc375c8a5010203262001215820ad7f7992c335b90d882b2802061b97a4fabca7e2ee3e7a51e728b8055e4eb9c7225820e0966ba7005987fece6f0e0e13447aa98cec248e4000a594b01b74c1cb1d40b3")
	var v attestationObject
	if err := Unmarshal(cborData, &v); err != nil {
		b.Fatal("Unmarshal:", err)
	}
	for i := 0; i < b.N; i++ {
		if _, err := Marshal(v); err != nil {
			b.Fatal("Marshal:", err)
		}
	}
}

func BenchmarkUnmarshalCOSEMAC(b *testing.B) {
	// Data from https://tools.ietf.org/html/rfc8392#appendix-A section A.4
	cborData := hexDecode("d83dd18443a10104a1044c53796d6d65747269633235365850a70175636f61703a2f2f61732e6578616d706c652e636f6d02656572696b77037818636f61703a2f2f6c696768742e6578616d706c652e636f6d041a5612aeb0051a5610d9f0061a5610d9f007420b7148093101ef6d789200")

	for i := 0; i < b.N; i++ {
		var v macedCOSE
		if err := Unmarshal(cborData, &v); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}

func BenchmarkUnmarshalCOSEMACWithTag(b *testing.B) {
	// Data from https://tools.ietf.org/html/rfc8392#appendix-A section A.4
	cborData := hexDecode("d83dd18443a10104a1044c53796d6d65747269633235365850a70175636f61703a2f2f61732e6578616d706c652e636f6d02656572696b77037818636f61703a2f2f6c696768742e6578616d706c652e636f6d041a5612aeb0051a5610d9f0061a5610d9f007420b7148093101ef6d789200")

	// Register tag CBOR Web Token (CWT) 61 and COSE_Mac0 17 with macedCOSE type
	tags := NewTagSet()
	if err := tags.Add(TagOptions{EncTag: EncTagRequired, DecTag: DecTagRequired}, reflect.TypeOf(macedCOSE{}), 61, 17); err != nil {
		b.Fatal("TagSet.Add:", err)
	}

	// Create DecMode with tags
	dm, _ := DecOptions{}.DecModeWithTags(tags)

	for i := 0; i < b.N; i++ {
		var v macedCOSE
		if err := dm.Unmarshal(cborData, &v); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}
func BenchmarkMarshalCOSEMAC(b *testing.B) {
	// Data from https://tools.ietf.org/html/rfc8392#appendix-A section A.4
	cborData := hexDecode("d83dd18443a10104a1044c53796d6d65747269633235365850a70175636f61703a2f2f61732e6578616d706c652e636f6d02656572696b77037818636f61703a2f2f6c696768742e6578616d706c652e636f6d041a5612aeb0051a5610d9f0061a5610d9f007420b7148093101ef6d789200")

	var v macedCOSE
	if err := Unmarshal(cborData, &v); err != nil {
		b.Fatal("Unmarshal():", err)
	}

	for i := 0; i < b.N; i++ {
		if _, err := Marshal(v); err != nil {
			b.Fatal("Marshal():", v, err)
		}
	}
}

func BenchmarkMarshalCOSEMACWithTag(b *testing.B) {
	// Data from https://tools.ietf.org/html/rfc8392#appendix-A section A.4
	cborData := hexDecode("d83dd18443a10104a1044c53796d6d65747269633235365850a70175636f61703a2f2f61732e6578616d706c652e636f6d02656572696b77037818636f61703a2f2f6c696768742e6578616d706c652e636f6d041a5612aeb0051a5610d9f0061a5610d9f007420b7148093101ef6d789200")

	// Register tag CBOR Web Token (CWT) 61 and COSE_Mac0 17 with macedCOSE type
	tags := NewTagSet()
	if err := tags.Add(TagOptions{EncTag: EncTagRequired, DecTag: DecTagRequired}, reflect.TypeOf(macedCOSE{}), 61, 17); err != nil {
		b.Fatal("TagSet.Add:", err)
	}

	// Create EncMode with tags.
	dm, _ := DecOptions{}.DecModeWithTags(tags)
	em, _ := EncOptions{}.EncModeWithTags(tags)

	var v macedCOSE
	if err := dm.Unmarshal(cborData, &v); err != nil {
		b.Fatal("Unmarshal():", err)
	}

	for i := 0; i < b.N; i++ {
		if _, err := em.Marshal(v); err != nil {
			b.Fatal("Marshal():", v, err)
		}
	}
}
