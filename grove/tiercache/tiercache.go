package tiercache

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
	"github.com/apache/trafficcontrol/v8/grove/cacheobj"
	"github.com/apache/trafficcontrol/v8/grove/icache"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
)

// TierCache wraps two icache.Caches and implements icache.Cache. Adding adds to both caches. Getting tries the first cache first, and if it isn't found, it then tries the second. Thus, the first should be smaller and faster (memory) and the second should be larger and slower (disk).
//
// This is more suitable for more less-frequently-requested items, with a separate cache for fewer frequently-requested items.
//
// An alternative implementation would be to Add to the first only, and when an object is evicted from the first, then store in the second. That would be more efficient for any given frequency, but less efficient for known infrequent objects.
type TierCache struct {
	first  icache.Cache
	second icache.Cache
}

// New creates a new TierCache with the given first and second caches to use.
func New(first, second icache.Cache) *TierCache {
	return &TierCache{first: first, second: second}
}

// Get returns the object if it's in the first cache. Else, it returns the object from the second cache. Else, false.
func (c *TierCache) Get(key string) (*cacheobj.CacheObj, bool) {
	log.Debugln("TierCache.Get '" + key + "' calling")
	v, ok := c.first.Get(key)
	log.Debugf("TierCache.Get '"+key+"' FOUND FIRST: %+v\n", ok)
	if !ok {
		v, ok = c.second.Get(key)
		if ok {
			// if it was in second but not first, add back to first (LRU behavior)
			c.first.Add(key, v)
		}
		log.Debugf("TierCache.Get '"+key+"' FOUND SECOND: %+v\n", ok)
	}
	return v, ok
}

// Peek returns the object if it's in the first cache. Else, it returns the object from the second cache. Else, false.
// Peek does not change the lru-ness, or the first cache
func (c *TierCache) Peek(key string) (*cacheobj.CacheObj, bool) {
	v, ok := c.first.Peek(key)
	if !ok {
		v, ok = c.second.Peek(key)
	}
	return v, ok
}

// Add adds to both internal caches. Returns whether either reported an eviction.
func (c *TierCache) Add(key string, val *cacheobj.CacheObj) bool {
	aevict := c.first.Add(key, val)
	bevict := c.second.Add(key, val)
	return aevict || bevict
}

// Size returns the size of the second cache. This is because, since all objects are added to both, they are presumed to have the same content, and the second is presumed to be larger.
//
// For example, if the first is a memory cache and the second is a disk cache, it's most useful to report the size used on disk.
func (c *TierCache) Size() uint64 {
	return c.second.Size()
}

func (c *TierCache) Close() {
	c.first.Close()
	c.second.Close()
}

// Keys returns the keys of the second tier only. The first is just an accelerator.
func (c *TierCache) Keys() []string {
	return c.second.Keys()
}

// Capacity returns the maximum size in bytes of the cache
func (c *TierCache) Capacity() uint64 { return c.second.Capacity() }
