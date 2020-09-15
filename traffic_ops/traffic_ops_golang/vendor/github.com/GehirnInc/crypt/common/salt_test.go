// (C) Copyright 2012, Jeramey Crawford <jeramey@antihe.ro>. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package common

import (
	"testing"
	"fmt"
	"strings"
)

var _Salt = &Salt{
	MagicPrefix: []byte("$foo$"),
	SaltLenMin:  1,
	SaltLenMax:  8,
	RoundsDefault: 5,
	RoundsMin: 1,
	RoundsMax: 10,
}

func TestGenerateSalt(t *testing.T) {
	salt := _Salt.Generate(0)
	if len(salt) != len(_Salt.MagicPrefix)+1 {
		t.Errorf("Expected len 1, got len %d", len(salt))
	}

	for i := 1; i <= 8; i++ {
		salt = _Salt.Generate(i)
		if len(salt) != len(_Salt.MagicPrefix)+i {
			t.Errorf("Expected len %d, got len %d", i, len(salt))
		}
	}

	salt = _Salt.Generate(9)
	if len(salt) != len(_Salt.MagicPrefix)+8 {
		t.Errorf("Expected len 8, got len %d", len(salt))
	}
}

func TestGenerateSaltWRounds(t *testing.T) {
	rounds := 7
	expectPrefix := fmt.Sprintf("%srounds=%d$", _Salt.MagicPrefix, 7)

	salt := _Salt.GenerateWRounds(10, rounds)
	if !strings.HasPrefix(string(salt), expectPrefix) {
		t.Errorf("Expected it has prefix \"%s\", but missing it", expectPrefix)
	}
}
