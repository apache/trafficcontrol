// (C) Copyright 2012, Jeramey Crawford <jeramey@antihe.ro>. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sha512_crypt

import "testing"

var sha512Crypt = New()

func TestGenerate(t *testing.T) {
	data := []struct {
		salt []byte
		key  []byte
		out  string
		cost int
	}{
		{
			[]byte("$6$saltstring"),
			[]byte("Hello world!"),
			"$6$saltstring$svn8UoSVapNtMuq1ukKS4tPQd8iKwSMHWjl/O817G3uBnIFNjn" +
				"QJuesI68u4OTLiBFdcbYEdFCoEOfaS35inz1",
			RoundsDefault,
		},
		{
			[]byte("$6$rounds=10000$saltstringsaltstring"),
			[]byte("Hello world!"),
			"$6$rounds=10000$saltstringsaltst$OW1/O6BYHV6BcXZu8QVeXbDWra3Oeqh" +
				"0sbHbbMCVNSnCM/UrjmM0Dp8vOuZeHBy/YTBmSK6H9qs/y3RnOaw5v.",
			10000,
		},
		{
			[]byte("$6$rounds=5000$toolongsaltstring"),
			[]byte("This is just a test"),
			"$6$rounds=5000$toolongsaltstrin$lQ8jolhgVRVhY4b5pZKaysCLi0QBxGoN" +
				"eKQzQ3glMhwllF7oGDZxUhx1yxdYcz/e1JSbq3y6JMxxl8audkUEm0",
			5000,
		},
		{
			[]byte("$6$rounds=1400$anotherlongsaltstring"),
			[]byte("a very much longer text to encrypt.  " +
				"This one even stretches over more" +
				"than one line."),
			"$6$rounds=1400$anotherlongsalts$POfYwTEok97VWcjxIiSOjiykti.o/pQs" +
				".wPvMxQ6Fm7I6IoYN3CmLs66x9t0oSwbtEW7o7UmJEiDwGqd8p4ur1",
			1400,
		},
		{
			[]byte("$6$rounds=77777$short"),
			[]byte("we have a short salt string but not a short password"),
			"$6$rounds=77777$short$WuQyW2YR.hBNpjjRhpYD/ifIw05xdfeEyQoMxIXbkv" +
				"r0gge1a1x3yRULJ5CCaUeOxFmtlcGZelFl5CxtgfiAc0",
			77777,
		},
		{
			[]byte("$6$rounds=123456$asaltof16chars.."),
			[]byte("a short string"),
			"$6$rounds=123456$asaltof16chars..$BtCwjqMJGx5hrJhZywWvt0RLE8uZ4o" +
				"PwcelCjmw2kSYu.Ec6ycULevoBK25fs2xXgMNrCzIMVcgEJAstJeonj1",
			123456,
		},
		{
			[]byte("$6$rounds=10$roundstoolow"),
			[]byte("the minimum number is still observed"),
			"$6$rounds=1000$roundstoolow$kUMsbe306n21p9R.FRkW3IGn.S9NPN0x50Yh" +
				"H1xhLsPuWGsUSklZt58jaTfF4ZEQpyUNGc0dqbpBYYBaHHrsX.",
			1000,
		},
	}

	for i, d := range data {
		hash, err := sha512Crypt.Generate(d.key, d.salt)
		if err != nil {
			t.Fatal(err)
		}
		if hash != d.out {
			t.Errorf("Test %d failed\nExpected: %s, got: %s", i, d.out, hash)
		}

		cost, err := sha512Crypt.Cost(hash)
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
		hash, err := sha512Crypt.Generate(d, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err = sha512Crypt.Verify(hash, d); err != nil {
			t.Errorf("Test %d failed: %s", i, d)
		}
	}
}
