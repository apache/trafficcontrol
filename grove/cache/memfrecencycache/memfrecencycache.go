package memfrecencycache

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
	"sync"
	"sync/atomic"

	"github.com/apache/trafficcontrol/grove/cache/frecencyheap"
	"github.com/apache/trafficcontrol/grove/cacheobj"
	"github.com/apache/trafficcontrol/lib/go-log"
)

// DefaultHitWeight is the recommended default hit percentage. This is not used within this package, but may be passed to New if a good default is desired.
const DefaultHitWeight = 0.01

// MemCache is a threadsafe memory cache with a soft byte limit, enforced via a Frecency Heap.
type MemCacheFrecency struct {
	heap         *frecencyheap.Heap            // threadsafe.
	cache        map[string]*cacheobj.CacheObj // mutexed: MUST NOT access without locking cacheM. TODO test performance of sync.Map
	cacheM       sync.RWMutex                  // TODO test performance of one mutex for frecency-heap+cache
	sizeBytes    uint64                        // atomic: MUST NOT access without sync.atomic
	maxSizeBytes uint64                        // constant: MUST NOT be modified after creation
	gcChan       chan<- uint64
}

// New creates a new frecency memory cache.
// The bytes is the softe byte cap for the cache.
// The hitWeight is the percentage that the hit count will weigh against recently-used, when considering for cache eviction. For a reasonable default, use DefaultHitWeight.
func New(bytes uint64, hitWeight float64) *MemCacheFrecency {
	log.Debugf("MemCacheFrecency.New: creating cache with %d capacity.", bytes)
	gcChan := make(chan uint64, 1)
	c := &MemCacheFrecency{
		heap:         frecencyheap.New(hitWeight),
		cache:        map[string]*cacheobj.CacheObj{},
		maxSizeBytes: bytes,
		gcChan:       gcChan,
	}
	go c.gcManager(gcChan)
	return c
}

func (c *MemCacheFrecency) Get(key string) (*cacheobj.CacheObj, bool) {
	c.cacheM.RLock()
	obj, ok := c.cache[key]
	if ok {
		c.heap.Add(key, obj.Size) // TODO directly call c.ll.MoveToFront
		atomic.AddUint64(&obj.HitCount, 1)
	}
	c.cacheM.RUnlock()
	return obj, ok
}

func (c *MemCacheFrecency) Peek(key string) (*cacheobj.CacheObj, bool) {
	c.cacheM.RLock()
	obj, ok := c.cache[key]
	c.cacheM.RUnlock()
	return obj, ok
}

func (c *MemCacheFrecency) Add(key string, val *cacheobj.CacheObj) bool {
	c.cacheM.Lock()
	c.cache[key] = val
	c.cacheM.Unlock()
	oldSize := c.heap.Add(key, val.Size)
	log.Errorln("DEBUG MemCacheFrecency.gc added key '" + key + "'")
	sizeChange := val.Size - oldSize
	if sizeChange == 0 {
		return false
	}
	newSizeBytes := atomic.AddUint64(&c.sizeBytes, sizeChange)
	if newSizeBytes <= c.maxSizeBytes {
		return false
	}
	c.doGC(newSizeBytes)
	return false // TODO remove eviction from interface; it's unnecessary and expensive
}

func (c *MemCacheFrecency) Size() uint64 { return atomic.LoadUint64(&c.sizeBytes) }
func (c *MemCacheFrecency) Close()       {}

// doGC kicks off garbage collection if it isn't already. Does not block.
func (c *MemCacheFrecency) doGC(cacheSizeBytes uint64) {
	select {
	case c.gcChan <- cacheSizeBytes:
	default: // don't block if GC is already running
	}
}

// gcManager is the garbage collection manager function, designed to be run in a goroutine. Never returns.
func (c *MemCacheFrecency) gcManager(gcChan <-chan uint64) {
	for cacheSizeBytes := range gcChan {
		c.gc(cacheSizeBytes)
	}
}

// gc executes garbage collection, until the cache size is under the max. This should be called in a singleton manager goroutine, so only one goroutine is ever doing garbage collection at any time.
func (c *MemCacheFrecency) gc(cacheSizeBytes uint64) {
	for cacheSizeBytes > c.maxSizeBytes {
		log.Debugf("MemCacheFrecency.gc cacheSizeBytes %+v > c.maxSizeBytes %+v\n", cacheSizeBytes, c.maxSizeBytes)
		key, sizeBytes, exists := c.heap.RemoveOldest() // TODO change heap to use strings
		if !exists {
			// should never happen
			log.Errorf("MemCacheFrecency.gc sizeBytes %v > %v maxSizeBytes, but Frecency-Heap is empty!? Setting cache size to 0!\n", cacheSizeBytes, c.maxSizeBytes)
			atomic.StoreUint64(&c.sizeBytes, 0)
			return
		}

		log.Errorln("DEBUG MemCacheFrecency.gc deleting key '" + key + "'")
		log.Debugln("MemCacheFrecency.gc deleting key '" + key + "'")
		c.cacheM.Lock()
		delete(c.cache, key)
		c.cacheM.Unlock()

		cacheSizeBytes = atomic.AddUint64(&c.sizeBytes, ^uint64(sizeBytes-1)) // subtract sizeBytes
	}
}

func (c *MemCacheFrecency) Keys() []string {
	return c.heap.Keys()
}

func (c *MemCacheFrecency) Capacity() uint64 {
	return c.maxSizeBytes
}
