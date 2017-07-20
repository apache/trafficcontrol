package grove

import (
	"testing"
)

func TestConsistentHash(t *testing.T) {
	vals := []string{
		"basdfasdf",
		"b6klmnfji",
		"b6klmnfjio0i89jzxcvl;kj3t0o9vdf",
		"b62mkl;xfioje09y7=gsdfl;ikjertj",
		"https://foo.example.net/bar/baz/asdf/",
		"https://bar.example.com/a/b/c/d/e/f/g",
	}

	hashes := map[uint64]string{}

	for _, val := range vals {
		h := ConsistentHash(val)
		if hval, ok := hashes[h]; ok {
			t.Errorf("TestConsistentHash expected no collison, actual %v and %v both hash to %v", val, hval, h)
		}
		hashes[h] = val
	}

}
