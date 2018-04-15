package icache

import (
	"github.com/apache/incubator-trafficcontrol/grove/cacheobj"
)

// TODO change to return errors

type Cache interface {
	Add(key string, val *cacheobj.CacheObj) bool
	Capacity() uint64
	Get(key string) (*cacheobj.CacheObj, bool)
	Peek(key string) (*cacheobj.CacheObj, bool)
	Keys() []string
	Size() uint64
	Close()
}
