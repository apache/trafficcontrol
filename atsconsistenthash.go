package grove

import (
	"fmt"
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
	return int(f + 0.5)
}

func (h *SimpleATSConsistentHash) Insert(node *ATSConsistentHashNode, weight float64) error {
	for i := 0; i < round(float64(h.Replicas)*weight); i++ {
		hashStr := fmt.Sprintf("%d-%s", i, node)
		hashKey := ConsistentHash(hashStr)
		h.NodeMap.Insert(hashKey, node)
	}
	return nil
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

	hashVal := ConsistentHash(name)
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
	// return ATSConsistentHashNode{}, nil, false // TODO implement
	wrapped := false
	iter := h.NodeMap.LowerBound(hashVal)
	if iter == nil {
		wrapped = true
		iter = h.NodeMap.First()
	}
	return iter, wrapped
}
