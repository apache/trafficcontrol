package diskcache

import (
	"errors"

	"github.com/apache/incubator-trafficcontrol/grove/cacheobj"
	"github.com/apache/incubator-trafficcontrol/grove/config"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"

	"github.com/dchest/siphash"
)

// MultiDiskCache is a disk cache using multiple files. It exists primarily to allow caching across multiple physical disks, but may be used for other purposes. For example, it may be more performant to use multiple files, or it may be advantageous to keep each remap rule in its own file. Keys are evenly distributed across the given files via consistent hashing.
type MultiDiskCache []*DiskCache

func NewMulti(files []config.CacheFile) (*MultiDiskCache, error) {
	caches := make([]*DiskCache, len(files), len(files))
	for i, file := range files {
		cache, err := New(file.Path, file.Bytes)
		if err != nil {
			return nil, errors.New("creating disk cache '" + file.Path + "': " + err.Error())
		}
		caches[i] = cache
	}

	mdc := MultiDiskCache(caches)
	return &mdc, nil
}

// KeyIdx gets the consistent-hashed index of which DiskCache the key is mapped to.
func (c *MultiDiskCache) keyIdx(key string) int {
	return int(siphash.Hash(0, 0, []byte(key)) % uint64(len(*c)))
}

func (c *MultiDiskCache) Add(key string, val *cacheobj.CacheObj) bool {
	i := c.keyIdx(key)
	log.Debugf("MultiDiskCache.Add key '%+v' size '%+v' mapped to %+v\n", key, val.Size, i)
	return (*c)[i].Add(key, val)
}

func (c *MultiDiskCache) Get(key string) (*cacheobj.CacheObj, bool) {
	i := c.keyIdx(key)
	log.Debugf("MultiDiskCache.Get key '%+v' mapped to %+v\n", key, i)
	return (*c)[i].Get(key)
}

func (c *MultiDiskCache) Size() uint64 {
	sum := uint64(0)
	for _, cache := range *c {
		sum += cache.Size()
	}
	return sum
}

func (c *MultiDiskCache) Close() {
	for _, cache := range *c {
		cache.Close()
	}
}
