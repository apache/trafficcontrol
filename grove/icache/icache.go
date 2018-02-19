package icache

import (
	"github.com/apache/incubator-trafficcontrol/grove/cacheobj"
)

// TODO change to return errors

type Cache interface {
	Add(key string, val *cacheobj.CacheObj) bool
	Get(key string) (*cacheobj.CacheObj, bool)
	Size() uint64
	Close()
}
