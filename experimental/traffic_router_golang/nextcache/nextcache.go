package nextcache

import (
	"sync/atomic"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

// New creates and returns a new NextCacher. The returned NextCacher is safe for use by multiple goroutines.
func New(dses []tc.DeliveryServiceName) NextCacher {
	m := make(map[tc.DeliveryServiceName]*uint64, len(dses))
	for _, ds := range dses {
		i := uint64(0)
		m[ds] = &i
	}
	return nextCacher(m)
}

// NextCacher is the interface that wraps the NextCache method.
//
// NextCache returns the next cache to use for the given delivery service. This is neither pure nor idempotent, and successive calls will return different numbers. The underlying mechanism may not be aware of the number of caches, and the returned number MAY exceed the number of caches. Typically, callers should mod the returned number by the size of their cache list, to determine the cache to use. Returns false if the given delivery service is not found.
type NextCacher interface {
	NextCache(tc.DeliveryServiceName) (uint64, bool)
}

type nextCacher map[tc.DeliveryServiceName]*uint64

func (c nextCacher) NextCache(ds tc.DeliveryServiceName) (uint64, bool) {
	m := (map[tc.DeliveryServiceName]*uint64)(c)
	i, ok := m[ds]
	if !ok {
		return 0, false
	}
	return atomic.AddUint64(i, 1), true
}
