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

func TestSimpleATSConsistentHashLookup(t *testing.T) {
	replicas := 10

	names := []string{"foo", "bar", "baz"}

	nodes := []*ATSConsistentHashNode{}
	for _, name := range names {
		nodes = append(nodes, &ATSConsistentHashNode{Name: name})
	}

	inNodes := func(s string) bool {
		for _, node := range nodes {
			if node.Name == s {
				return true
			}
		}
		return false
	}

	h := NewSimpleATSConsistentHash(replicas)
	for _, node := range nodes {
		h.Insert(node, 1.0) // TODO test weights
	}

	lookup0 := "lookupasdf"
	lookup1 := "lookupjkl;aasdfqeroipuzxcn;"
	i, _, err := h.Lookup(lookup0)
	if err != nil {
		t.Errorf("ATSConsistentHash.Lookup expected nil error, actual %v", err)
	}
	lookup0Val := i.Val().Name
	// fmt.Printf("ATSConsistentHash.Lookup0 got %v\n", i.Val().Name)
	if !inNodes(i.Val().Name) {
		t.Errorf("ATSConsistentHash.Lookup expected in %+v actual %v", names, i.Val().Name)
	}

	i, _, err = h.Lookup(lookup1)
	if err != nil {
		t.Errorf("ATSConsistentHash.Lookup expected nil error, actual %v", err)
	}
	lookup1Val := i.Val().Name
	// fmt.Printf("ATSConsistentHash.Lookup1 got %v\n", i.Val().Name)
	if !inNodes(i.Val().Name) {
		t.Errorf("ATSConsistentHash.Lookup expected in %+v actual %v", names, i.Val().Name)
	}

	i, _, err = h.Lookup(lookup0)
	if err != nil {
		t.Errorf("ATSConsistentHash.Lookup expected nil error, actual %v", err)
	}
	if i.Val().Name != lookup0Val {
		t.Errorf("ATSConsistentHash.Lookup expected consistent %v actual %v", lookup0Val, i.Val().Name)
	}
	// fmt.Printf("ATSConsistentHash.Lookup0 got %v\n", i.Val().Name)
	if !inNodes(i.Val().Name) {
		t.Errorf("ATSConsistentHash.Lookup expected in %+v actual %v", names, i.Val().Name)
	}

	i, _, err = h.Lookup(lookup1)
	if err != nil {
		t.Errorf("ATSConsistentHash.Lookup expected nil error, actual %v", err)
	}
	if i.Val().Name != lookup1Val {
		t.Errorf("ATSConsistentHash.Lookup expected consistent %v actual %v", lookup1Val, i.Val().Name)
	}
	// fmt.Printf("ATSConsistentHash.Lookup1 got %v\n", i.Val().Name)
	if !inNodes(i.Val().Name) {
		t.Errorf("ATSConsistentHash.Lookup expected in %+v actual %v", names, i.Val().Name)
	}

}
