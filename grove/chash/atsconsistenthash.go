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
	"math"
	"strconv"

	"github.com/dchest/siphash"
)

// This is specifically designed to match the Apache Traffic Server Parent Selection Consistent Hash, so that Grove deployed alongside ATS will hash to the same parent (mid-tier) caches, and thus result in the same mids caching the same content

// Transliterated from https://raw.githubusercontent.com/apache/trafficserver/master/lib/ts/HashSip.cc
// via https://github.com/floodyberry/siphash
// via https://131002.net/siphash/

// type HashRing interface {
// 	Get(key string) string
// 	// Add(val string)
// 	// AddWeighted(val string, weight int)
// }

// type ATSHashRing struct {
// }

type ATSConsistentHash interface {
	Insert(node *ATSConsistentHashNode, weight float64) error
	String() string
	// Lookup returns the found node, its map iterator, and whether the lookup wrapped
	Lookup(name string) (OrderedMapUint64NodeIterator, bool, error)
	LookupHash(hashVal uint64) (OrderedMapUint64NodeIterator, bool)
	LookupIter(OrderedMapUint64NodeIterator) (OrderedMapUint64NodeIterator, bool)
	First() OrderedMapUint64NodeIterator // debug
}

const DefaultSimpleATSConsistentHashReplicas = 1024

type SimpleATSConsistentHash struct {
	Replicas int
	NodeMap  OrderedMapUint64Node
}

func NewSimpleATSConsistentHash(replicas int) ATSConsistentHash {
	return &SimpleATSConsistentHash{Replicas: replicas, NodeMap: NewSimpleOrderedMapUint64Node()}
}

func round(f float64) int {
	if math.Abs(f) < 0.5 {
		return 0
	}
	return int(f + math.Copysign(0.5, f))
}

func (h *SimpleATSConsistentHash) String() string {
	return h.NodeMap.String()
}

func (h *SimpleATSConsistentHash) Insert(node *ATSConsistentHashNode, weight float64) error {
	numInserts := round(float64(h.Replicas) * weight)
	keys := make([]uint64, numInserts)
	vals := make([]*ATSConsistentHashNode, numInserts)
	for i := 0; i < numInserts; i++ {
		hashStr := ""
		if node.ProxyURL != nil {
			hashStr = strconv.Itoa(i) + "-" + node.ProxyURL.Hostname()
		} else {
			hashStr = strconv.Itoa(i) + "-" + node.Name
		}
		hashKey := siphash.Hash(0, 0, []byte(hashStr))
		keys[i] = hashKey
		vals[i] = node
	}
	err := h.NodeMap.InsertBulk(keys, vals)
	return err
}

func (h *SimpleATSConsistentHash) First() OrderedMapUint64NodeIterator {
	return h.NodeMap.First()
}

// Lookup returns the found node, its map iterator, and whether the lookup wrapped
func (h *SimpleATSConsistentHash) Lookup(name string) (OrderedMapUint64NodeIterator, bool, error) {
	iter := OrderedMapUint64NodeIterator(nil)

	if name == "" {
		// (*iter)++;
		return nil, false, fmt.Errorf("lookup name is empty")
	}

	hashVal := siphash.Hash(0, 0, []byte(name))
	iter = h.NodeMap.LowerBound(hashVal)

	wrapped := false
	if iter == nil {
		wrapped = true
		iter = h.NodeMap.First()
	}

	if wrapped && iter == nil {
		return nil, false, fmt.Errorf("not found")
	}

	return iter, wrapped, nil

}

func (h *SimpleATSConsistentHash) LookupIter(i OrderedMapUint64NodeIterator) (OrderedMapUint64NodeIterator, bool) {
	wrapped := false
	if i == nil {
		i = h.NodeMap.First()
		wrapped = true
	} else {
		i = i.Next()
	}
	if i == nil {
		i = h.NodeMap.First()
		wrapped = true
	}
	return i, wrapped
}

func (h *SimpleATSConsistentHash) LookupHash(hashVal uint64) (OrderedMapUint64NodeIterator, bool) {
	wrapped := false
	iter := h.NodeMap.LowerBound(hashVal)
	if iter == nil {
		wrapped = true
		iter = h.NodeMap.First()
	}
	return iter, wrapped
}
