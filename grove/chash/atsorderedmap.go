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
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

// func NewATSHashRing(vals []string) HashRing {

// }

// TODO move Weighted to a separate struct, `ATSWeightedHashRing`
// func NewATSHashRingWeighted(vals map[string]int) HashRing {

// }

// func (h *ATSHashRing) Add(val string) {

// }

// func (h *ATSHashRing) AddWeighted(val string, weight int) {

// }

// ATSConsistentHashNode is an ATS ParentRecord
type ATSConsistentHashNode struct {
	Available bool
	Name      string
	ProxyURL  *url.URL
	Transport *http.Transport
	// pRecord fields (ParentSelection.h)
	Hostname  string
	Port      int
	FailedAt  time.Time
	FailCount int
	UpAt      int
	Scheme    string
	Index     int
	Weight    float64
}

func (n ATSConsistentHashNode) String() string {
	return n.Name
}

type OrderedMapUint64NodeIterator interface {
	Val() *ATSConsistentHashNode
	Key() uint64
	Next() OrderedMapUint64NodeIterator
	NextWrap() OrderedMapUint64NodeIterator
	Index() int
}

type OrderedMapUint64Node interface {
	Insert(key uint64, val *ATSConsistentHashNode)
	String() string
	InsertBulk(keys []uint64, vals []*ATSConsistentHashNode) error
	First() OrderedMapUint64NodeIterator
	Last() OrderedMapUint64NodeIterator
	At(index int) (uint64, *ATSConsistentHashNode)
	LowerBound(val uint64) OrderedMapUint64NodeIterator
}

type SimpleOrderedMapUInt64Node struct {
	M map[uint64]*ATSConsistentHashNode
	O []uint64
}

func NewSimpleOrderedMapUint64Node() OrderedMapUint64Node {
	return &SimpleOrderedMapUInt64Node{M: map[uint64]*ATSConsistentHashNode{}, O: []uint64{}}
}

type SortableUint64 []uint64

func (a SortableUint64) Len() int           { return len(a) }
func (a SortableUint64) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortableUint64) Less(i, j int) bool { return a[i] < a[j] }

func (m *SimpleOrderedMapUInt64Node) Insert(key uint64, val *ATSConsistentHashNode) {
	m.M[key] = val
	m.O = append(m.O, key)
	sort.Sort(SortableUint64(m.O))
}

func (m *SimpleOrderedMapUInt64Node) InsertBulk(keys []uint64, vals []*ATSConsistentHashNode) error {
	if len(keys) != len(vals) {
		return fmt.Errorf("SimpleOrderedMapUInt64Node InsertBulk failed - len(keys) != len(vals)")
	}

	for i := 0; i < len(keys); i++ {
		// fmt.Println("InsertBulk " + strconv.FormatUint(keys[i], 10) + ": " + vals[i].String() + " " + vals[i].ProxyURL.String())
		m.M[keys[i]] = vals[i]
	}

	m.O = nil // clear, in case there were previous inserts
	for k := range m.M {
		m.O = append(m.O, k)
	}
	sort.Sort(SortableUint64(m.O))
	return nil
}

func (m *SimpleOrderedMapUInt64Node) LowerBound(key uint64) OrderedMapUint64NodeIterator {
	// TODO change to binary search
	for i := 0; i < len(m.O); i++ {
		if m.O[i] >= key {
			return NewSimpleOrderedMapUint64NodeIterator(key, m.M[m.O[i]], i, m)
		}
	}
	return nil
}

func (m *SimpleOrderedMapUInt64Node) String() string {
	s := ""
	i := m.First()
	for {
		if i == nil {
			return s
		}
		s += strconv.FormatUint(i.Key(), 10) + ": " + i.Val().String() + " " + i.Val().ProxyURL.String() + "\n"
		i = i.Next()
	}
}

// First returns the iterator to the first element in the map. Returns nil if the map is empty
func (m *SimpleOrderedMapUInt64Node) First() OrderedMapUint64NodeIterator {
	if len(m.O) == 0 {
		return nil
	}
	i := 0
	key := m.O[0]
	val := m.M[key]
	return NewSimpleOrderedMapUint64NodeIterator(key, val, i, m)
}

// Last returns the iterator to the last element in the map. Returns nil if the map is empty
func (m *SimpleOrderedMapUInt64Node) Last() OrderedMapUint64NodeIterator {
	if len(m.O) == 0 {
		return nil
	}
	i := len(m.O) - 1
	key := m.O[i]
	val := m.M[key]
	return NewSimpleOrderedMapUint64NodeIterator(key, val, i, m)
}

func (m *SimpleOrderedMapUInt64Node) At(i int) (uint64, *ATSConsistentHashNode) {
	key := m.O[i]
	val := m.M[key]
	return key, val
}

func NewSimpleOrderedMapUint64NodeIterator(key uint64, val *ATSConsistentHashNode, index int, m *SimpleOrderedMapUInt64Node) OrderedMapUint64NodeIterator {
	return &SimpleOrderedMapUint64NodeIterator{key: key, val: val, index: index, m: m}
}

type SimpleOrderedMapUint64NodeIterator struct {
	key   uint64
	val   *ATSConsistentHashNode
	index int
	m     *SimpleOrderedMapUInt64Node
}

func (i *SimpleOrderedMapUint64NodeIterator) Val() *ATSConsistentHashNode { return i.val }
func (i *SimpleOrderedMapUint64NodeIterator) Key() uint64                 { return i.key }

func (i *SimpleOrderedMapUint64NodeIterator) Next() OrderedMapUint64NodeIterator {
	next := i.index + 1
	if next >= len(i.m.M) {
		return nil
	}
	key, val := i.m.At(next)
	return NewSimpleOrderedMapUint64NodeIterator(key, val, next, i.m)
}

func (i *SimpleOrderedMapUint64NodeIterator) NextWrap() OrderedMapUint64NodeIterator {
	next := i.Next()
	if next == nil {
		return i.m.First()
	}
	return next
}

func (i *SimpleOrderedMapUint64NodeIterator) Prev() OrderedMapUint64NodeIterator {
	prevI := i.index - 1
	if prevI <= len(i.m.M) {
		return nil
	}
	key, val := i.m.At(prevI)
	return NewSimpleOrderedMapUint64NodeIterator(key, val, prevI, i.m)
}

func (i *SimpleOrderedMapUint64NodeIterator) Index() int {
	return i.index
}
