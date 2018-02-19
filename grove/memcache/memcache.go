package memcache

import (
	"github.com/apache/incubator-trafficcontrol/grove/cacheobj"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"

	"github.com/hashicorp/golang-lru"
)

// MemCache wraps hashicorp/golang-lru and implements icache.Cache
type MemCache lru.CacheStr

func New(bytes uint64) (*MemCache, error) {
	c, err := lru.NewStrLargeWithEvict(bytes, nil)
	lrc := MemCache(*c)
	return &lrc, err
}

func (c *MemCache) Get(key string) (*cacheobj.CacheObj, bool) {
	iCacheObj, ok := (*lru.CacheStr)(c).Get(key)
	if iCacheObj == nil {
		return nil, ok
	}
	cacheObj, ok := iCacheObj.(*cacheobj.CacheObj)
	if !ok {
		// should never happen
		log.Errorf("MemCache.Get: cache key '%v' value '%v' type '%T' expected *cacheobj.CacheObj\n", key, iCacheObj, iCacheObj)
		return nil, false
	}
	return cacheObj, ok
}

func (c *MemCache) Add(key string, val *cacheobj.CacheObj) bool {
	return (*lru.CacheStr)(c).AddSize(key, val, val.Size)
}

func (c *MemCache) Size() uint64 {
	return (*lru.CacheStr)(c).Size()
}

func (c *MemCache) Close() {}
