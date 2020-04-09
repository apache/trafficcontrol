// (C) Copyright 2012, Jeramey Crawford <jeramey@antihe.ro>. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sha256_crypt

import "testing"

var sha256Crypt = New()

func TestGenerate(t *testing.T) {
	data := []struct {
		salt []byte
		key  []byte
		out  string
		cost int
	}{
		{
			[]byte("$5$saltstring"),
			[]byte("Hello world!"),
			"$5$saltstring$5B8vYYiY.CVt1RlTTf8KbXBH3hsxY/GNooZaBBGWEc5",
			RoundsDefault,
		},
		{
			[]byte("$5$rounds=10000$saltstringsaltstring"),
			[]byte("Hello world!"),
			"$5$rounds=10000$saltstringsaltst$3xv.VbSHBb41AL9AvLeujZkZRBAwqFM" +
				"z2.opqey6IcA",
			10000,
		},
		{
			[]byte("$5$rounds=5000$toolongsaltstring"),
			[]byte("This is just a test"),
			"$5$rounds=5000$toolongsaltstrin$Un/5jzAHMgOGZ5.mWJpuVolil07guHPv" +
				"OW8mGRcvxa5",
			5000,
		},
		{
			[]byte("$5$rounds=1400$anotherlongsaltstring"),
			[]byte("a very much longer text to encrypt.  " +
				"This one even stretches over more" +
				"than one line."),
			"$5$rounds=1400$anotherlongsalts$Rx.j8H.h8HjEDGomFU8bDkXm3XIUnzyx" +
				"f12oP84Bnq1",
			1400,
		},
		{
			[]byte("$5$rounds=77777$short"),
			[]byte("we have a short salt string but not a short password"),
			"$5$rounds=77777$short$JiO1O3ZpDAxGJeaDIuqCoEFysAe1mZNJRs3pw0KQRd/",
			77777,
		},
		{
			[]byte("$5$rounds=123456$asaltof16chars.."),
			[]byte("a short string"),
			"$5$rounds=123456$asaltof16chars..$gP3VQ/6X7UUEW3HkBn2w1/Ptq2jxPy" +
				"zV/cZKmF/wJvD",
			123456,
		},
		{
			[]byte("$5$rounds=10$roundstoolow"),
			[]byte("the minimum number is still observed"),
			"$5$rounds=1000$roundstoolow$yfvwcWrQ8l/K0DAWyuPMDNHpIVlTQebY9l/g" +
				"L972bIC",
			1000,
		},
	}

	for i, d := range data {
		hash, err := sha256Crypt.Generate(d.key, d.salt)
		if err != nil {
			t.Fatal(err)
		}
		if hash != d.out {
			t.Errorf("Test %d failed\nExpected: %s, got: %s", i, d.out, hash)
		}

		cost, err := sha256Crypt.Cost(hash)
		if err != nil {
			t.Fatal(err)
		}
		if cost != d.cost {
			t.Errorf("Test %d failed\nExpected: %d, got: %d", i, d.cost, cost)
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
		hash, err := sha256Crypt.Generate(d, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err = sha256Crypt.Verify(hash, d); err != nil {
			t.Errorf("Test %d failed: %s", i, d)
		}
	}
}
