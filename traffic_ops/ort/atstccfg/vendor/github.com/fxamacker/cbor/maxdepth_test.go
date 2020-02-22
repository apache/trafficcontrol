// Copyright (c) Faye Amacker. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root for license information.

package cbor

import (
	"testing"
)

func TestDepth(t *testing.T) {
	testCases := []struct {
		name      string
		cborData  []byte
		wantDepth int
	}{
		{"uint", hexDecode("00"), 1},                                                          // 0
		{"int", hexDecode("20"), 1},                                                           // -1
		{"bool", hexDecode("f4"), 1},                                                          // false
		{"nil", hexDecode("f6"), 1},                                                           // nil
		{"float", hexDecode("fa47c35000"), 1},                                                 // 100000.0
		{"byte string", hexDecode("40"), 1},                                                   // []byte{}
		{"indefinite length byte string", hexDecode("5f42010243030405ff"), 1},                 // []byte{1, 2, 3, 4, 5}
		{"text string", hexDecode("60"), 1},                                                   // ""
		{"indefinite length text string", hexDecode("7f657374726561646d696e67ff"), 1},         // "streaming"
		{"empty array", hexDecode("80"), 1},                                                   // []
		{"indefinite length empty array", hexDecode("9fff"), 1},                               // []
		{"array", hexDecode("98190102030405060708090a0b0c0d0e0f101112131415161718181819"), 2}, // [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25]
		{"indefinite length array", hexDecode("9f0102030405060708090a0b0c0d0e0f101112131415161718181819ff"), 2}, // [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25]
		{"nested array", hexDecode("8301820203820405"), 3},                                                      // [1,[2,3],[4,5]]
		{"indefinite length nested array", hexDecode("83018202039f0405ff"), 3},                                  // [1,[2,3],[4,5]]
		{"array and map", hexDecode("826161a161626163"), 3},                                                     // [a", {"b": "c"}]
		{"indefinite length array and map", hexDecode("826161bf61626163ff"), 3},                                 // [a", {"b": "c"}]
		{"empty map", hexDecode("a0"), 1},                                                                       // {}
		{"indefinite length empty map", hexDecode("bfff"), 1},                                                   // {}
		{"map", hexDecode("a201020304"), 2},                                                                     // {1:2, 3:4}
		{"nested map", hexDecode("a26161016162820203"), 3},                                                      // {"a": 1, "b": [2, 3]}
		{"indefinite length nested map", hexDecode("bf61610161629f0203ffff"), 3},                                // {"a": 1, "b": [2, 3]}
		{"tag", hexDecode("c074323031332d30332d32315432303a30343a30305a"), 1},                                   // 0("2013-03-21T20:04:00Z")
		{"tagged map", hexDecode("d864a26161016162820203"), 3},                                                  // 100({"a": 1, "b": [2, 3]})
		{"tagged map and array", hexDecode("d864a26161016162d865d866820203"), 4},                                // 100({"a": 1, "b": 101(102([2, 3]))})
		{"nested tag", hexDecode("d864d865d86674323031332d30332d32315432303a30343a30305a"), 3},                  // 100(101(102("2013-03-21T20:04:00Z")))
		{"32-level array", hexDecode("820181818181818181818181818181818181818181818181818181818181818101"), 32},
		{"32-level indefinite length array", hexDecode("9f0181818181818181818181818181818181818181818181818181818181818101ff"), 32},
		{"32-level map", hexDecode("a10181818181818181818181818181818181818181818181818181818181818101"), 32},
		{"32-level indefinite length map", hexDecode("bf0181818181818181818181818181818181818181818181818181818181818101ff"), 32},
		{"32-level tag", hexDecode("d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d86474323031332d30332d32315432303a30343a30305a"), 32}, // 100(100(...("2013-03-21T20:04:00Z")))
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, depth, err := validInternal(tc.cborData, 0, 1)
			if err != nil {
				t.Errorf("valid(0x%x) returned error %v", tc.cborData, err)
			}
			if depth != tc.wantDepth {
				t.Errorf("valid(0x%x) returned depth %d, want %d", tc.cborData, depth, tc.wantDepth)
			}
		})
	}
}

func TestDepthError(t *testing.T) {
	testCases := []struct {
		name         string
		cborData     []byte
		wantErrorMsg string
	}{
		{"33-level array", hexDecode("82018181818181818181818181818181818181818181818181818181818181818101"), "cbor: reached max depth 32"},
		{"33-level indefinite length array", hexDecode("9f018181818181818181818181818181818181818181818181818181818181818101ff"), "cbor: reached max depth 32"},
		{"33-level map", hexDecode("a1018181818181818181818181818181818181818181818181818181818181818101"), "cbor: reached max depth 32"},
		{"33-level indefinite length map", hexDecode("bf018181818181818181818181818181818181818181818181818181818181818101ff"), "cbor: reached max depth 32"},
		{"33-level tag", hexDecode("d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d864d86474323031332d30332d32315432303a30343a30305a"), "cbor: reached max depth 32"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := validInternal(tc.cborData, 0, 1)
			if err == nil {
				t.Errorf("valid(0x%x) didn't return an error, want %q", tc.cborData, tc.wantErrorMsg)
			} else if err.Error() != tc.wantErrorMsg {
				t.Errorf("valid(0x%x) returned error %q, want %q", tc.cborData, err.Error(), tc.wantErrorMsg)
			}
		})
	}
}
