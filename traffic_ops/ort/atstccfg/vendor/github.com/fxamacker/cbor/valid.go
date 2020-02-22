// Copyright (c) Faye Amacker. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root for license information.

package cbor

import (
	"encoding/binary"
	"errors"
	"io"
	"strconv"
)

// SyntaxError is a description of a CBOR syntax error.
type SyntaxError struct {
	msg string
}

func (e *SyntaxError) Error() string { return e.msg }

// SemanticError is a description of a CBOR semantic error.
type SemanticError struct {
	msg string
}

func (e *SemanticError) Error() string { return e.msg }

// valid checks whether CBOR data is complete and well-formed.
func valid(data []byte) (rest []byte, err error) {
	if len(data) == 0 {
		return nil, io.EOF
	}
	offset, _, err := validInternal(data, 0, 1)
	if err != nil {
		return nil, err
	}
	return data[offset:], nil
}

const (
	maxNestingLevel = 32
)

// validInternal checks data's well-formedness and returns data's next offset, max depth, and error.
func validInternal(data []byte, off int, depth int) (int, int, error) {
	if depth > maxNestingLevel {
		return 0, 0, errors.New("cbor: reached max depth " + strconv.Itoa(maxNestingLevel))
	}

	off, t, ai, val, err := validHead(data, off)
	if err != nil {
		return 0, 0, err
	}

	if ai == 31 {
		if t == cborTypeByteString || t == cborTypeTextString {
			return validIndefiniteString(data, off, t, depth)
		}
		return validIndefiniteArrOrMap(data, off, t, depth)
	}

	dataLen := len(data)

	switch t {
	case cborTypeByteString, cborTypeTextString:
		valInt := int(val)
		if valInt < 0 {
			// Detect integer overflow
			return 0, 0, errors.New("cbor: " + t.String() + " length " + strconv.FormatUint(val, 10) + " is too large, causing integer overflow")
		}
		if dataLen-off < valInt { // valInt+off may overflow integer
			return 0, 0, io.ErrUnexpectedEOF
		}
		off += valInt
	case cborTypeArray, cborTypeMap:
		valInt := int(val)
		if valInt < 0 {
			// Detect integer overflow
			return 0, 0, errors.New("cbor: " + t.String() + " length " + strconv.FormatUint(val, 10) + " is too large, causing integer overflow")
		}
		count := 1
		if t == cborTypeMap {
			count = 2
		}
		maxDepth := depth
		for j := 0; j < count; j++ {
			for i := 0; i < valInt; i++ {
				var d int
				if off, d, err = validInternal(data, off, depth+1); err != nil {
					return 0, 0, err
				}
				if d > maxDepth {
					maxDepth = d // Save max depth
				}
			}
		}
		depth = maxDepth
	case cborTypeTag:
		// Scan nested tag numbers to avoid recursion.
		for {
			if dataLen == off { // Tag number must be followed by tag content.
				return 0, 0, io.ErrUnexpectedEOF
			}
			if cborType(data[off]&0xe0) != cborTypeTag {
				break
			}
			if off, _, _, _, err = validHead(data, off); err != nil {
				return 0, 0, err
			}
			depth++
		}
		// Check tag content.
		return validInternal(data, off, depth)
	}
	return off, depth, nil
}

// validIndefiniteString checks indefinite length byte/text string's well-formedness and returns data's next offset, max depth, and error.
func validIndefiniteString(data []byte, off int, t cborType, depth int) (int, int, error) {
	var err error
	dataLen := len(data)
	for {
		if dataLen == off {
			return 0, 0, io.ErrUnexpectedEOF
		}
		if data[off] == 0xff {
			off++
			break
		}
		// Peek ahead to get next type and indefinite length status.
		nt := cborType(data[off] & 0xe0)
		if t != nt {
			return 0, 0, &SyntaxError{"cbor: wrong element type " + nt.String() + " for indefinite-length " + t.String()}
		}
		if (data[off] & 0x1f) == 31 {
			return 0, 0, &SyntaxError{"cbor: indefinite-length " + t.String() + " chunk is not definite-length"}
		}
		if off, depth, err = validInternal(data, off, depth); err != nil {
			return 0, 0, err
		}
	}
	return off, depth, nil
}

// validIndefiniteArrOrMap checks indefinite length array/map's well-formedness and returns data's next offset, max depth, and error.
func validIndefiniteArrOrMap(data []byte, off int, t cborType, depth int) (int, int, error) {
	var err error
	maxDepth := depth
	dataLen := len(data)
	i := 0
	for {
		if dataLen == off {
			return 0, 0, io.ErrUnexpectedEOF
		}
		if data[off] == 0xff {
			off++
			break
		}
		var d int
		if off, d, err = validInternal(data, off, depth+1); err != nil {
			return 0, 0, err
		}
		if d > maxDepth {
			maxDepth = d
		}
		i++
	}
	if t == cborTypeMap && i%2 == 1 {
		return 0, 0, &SyntaxError{"cbor: unexpected \"break\" code"}
	}
	return off, maxDepth, nil
}

func validHead(data []byte, off int) (_ int, t cborType, ai byte, val uint64, err error) {
	dataLen := len(data) - off
	if dataLen == 0 {
		return 0, 0, 0, 0, io.ErrUnexpectedEOF
	}

	t = cborType(data[off] & 0xe0)
	ai = data[off] & 0x1f
	val = uint64(ai)
	off++

	if ai < 24 {
		return off, t, ai, val, nil
	}
	if ai == 24 {
		if dataLen < 2 {
			return 0, 0, 0, 0, io.ErrUnexpectedEOF
		}
		val = uint64(data[off])
		off++
		if t == cborTypePrimitives && val < 32 {
			return 0, 0, 0, 0, &SyntaxError{"cbor: invalid simple value " + strconv.Itoa(int(val)) + " for type " + t.String()}
		}
		return off, t, ai, val, nil
	}
	if ai == 25 {
		if dataLen < 3 {
			return 0, 0, 0, 0, io.ErrUnexpectedEOF
		}
		val = uint64(binary.BigEndian.Uint16(data[off : off+2]))
		off += 2
		return off, t, ai, val, nil
	}
	if ai == 26 {
		if dataLen < 5 {
			return 0, 0, 0, 0, io.ErrUnexpectedEOF
		}
		val = uint64(binary.BigEndian.Uint32(data[off : off+4]))
		off += 4
		return off, t, ai, val, nil
	}
	if ai == 27 {
		if dataLen < 9 {
			return 0, 0, 0, 0, io.ErrUnexpectedEOF
		}
		val = binary.BigEndian.Uint64(data[off : off+8])
		off += 8
		return off, t, ai, val, nil
	}
	if ai == 31 {
		switch t {
		case cborTypePositiveInt, cborTypeNegativeInt, cborTypeTag:
			return 0, 0, 0, 0, &SyntaxError{"cbor: invalid additional information " + strconv.Itoa(int(ai)) + " for type " + t.String()}
		case cborTypePrimitives: // 0xff (break code) should not be outside validIndefinite().
			return 0, 0, 0, 0, &SyntaxError{"cbor: unexpected \"break\" code"}
		}
		return off, t, ai, val, nil
	}
	// ai == 28, 29, 30
	return 0, 0, 0, 0, &SyntaxError{"cbor: invalid additional information " + strconv.Itoa(int(ai)) + " for type " + t.String()}
}
