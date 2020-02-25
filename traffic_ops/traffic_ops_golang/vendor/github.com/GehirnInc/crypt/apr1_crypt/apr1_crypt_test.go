// (C) Copyright 2012, Jeramey Crawford <jeramey@antihe.ro>. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apr1_crypt

import "testing"

var apr1Crypt = New()

func TestGenerate(t *testing.T) {
	data := []struct {
		salt []byte
		key  []byte
		out  string
	}{
		{
			[]byte("$apr1$$"),
			[]byte("abcdefghijk"),
			"$apr1$$NTjzQjNZnhYRPxN6ryN191",
		},
		{
			[]byte("$apr1$an overlong salt$"),
			[]byte("abcdefgh"),
			"$apr1$an overl$iroRZrWCEoQojCkf6p8LC0",
		},
		{
			[]byte("$apr1$12345678$"),
			[]byte("Lorem ipsum dolor sit amet"),
			"$apr1$12345678$/DpfgRGBHG8N0cbkmw0Fk/",
		},
		{
			[]byte("$apr1$deadbeef$"),
			[]byte("password"),
			"$apr1$deadbeef$NWLhx1Ai4ScyoaAboTFco.",
		},
		{
			[]byte("$apr1$$"),
			[]byte("missing salt"),
			"$apr1$$EcorjwkoQz4mYcksVEk6j0",
		},
		{
			[]byte("$apr1$holy-moly-batman$"),
			[]byte("1234567"),
			"$apr1$holy-mol$/WX0350ZUEkvQkrrVJsrU.",
		},
		{
			[]byte("$apr1$asdfjkl;$"),
			[]byte("A really long password. " +
				"Longer than a password has any righ" +
				"t to be. Hey bub, don't mess with t" +
				"his password."),
			"$apr1$asdfjkl;$2MbDUb/Bj6qcIIf38PXzp0",
		},
	}
	for i, d := range data {
		hash, err := apr1Crypt.Generate(d.key, d.salt)
		if err != nil {
			t.Fatal(err)
		}
		if hash != d.out {
			t.Errorf("Test %d failed\nExpected: %s, got: %s", i, d.out, hash)
		}
	}
}

func TestVerify(t *testing.T) {
	data := [][]byte{
		[]byte("password"),
		[]byte("12345"),
		[]byte("That's amazing! I've got the same combination on my luggage!"),
		[]byte("And change the combination on my luggage!"),
		[]byte("         random  spa  c    ing."),
		[]byte("94ajflkvjzpe8u3&*j1k513KLJ&*()"),
	}
	for i, d := range data {
		hash, err := apr1Crypt.Generate(d, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err = apr1Crypt.Verify(hash, d); err != nil {
			t.Errorf("Test %d failed: %s", i, d)
		}
	}
}
