package grove

import (
	"testing"
)

func TestSimpleOrderedMapUInt64Node(t *testing.T) {
	vals := map[uint64]*ATSConsistentHashNode{
		1: &ATSConsistentHashNode{Name: "foo"},
		2: &ATSConsistentHashNode{Name: "bar"},
		3: &ATSConsistentHashNode{Name: "baz"},
	}

	m := NewSimpleOrderedMapUint64Node()

	for k, v := range vals {
		m.Insert(k, v)
	}

	count := 0
	for i := m.First(); i != nil; i = i.Next() {
		count++
		k := i.Key()
		val := i.Val()

		// fmt.Printf("OrderedMapUint64Node %v : %v\n", k, val)

		valsVal, exists := vals[k]
		if !exists {
			t.Errorf("hash key %v was not inserted!", k)
		}
		if valsVal != val {
			t.Errorf("hash key %v value expected %v actual %v", k, valsVal, val)
		}
		delete(vals, k)
	}
	if len(vals) != 0 {
		t.Errorf("hash entries expected %+v actual nil", vals)
	}
}
