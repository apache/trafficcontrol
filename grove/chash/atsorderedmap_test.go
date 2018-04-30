package chash

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

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
