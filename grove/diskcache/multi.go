package diskcache

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
	"errors"

	"github.com/apache/trafficcontrol/v8/grove/cacheobj"
	"github.com/apache/trafficcontrol/v8/grove/config"

	"github.com/apache/trafficcontrol/v8/lib/go-log"

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
		cache.ResetAfterRestart() // should this be optional?
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

func (c *MultiDiskCache) Peek(key string) (*cacheobj.CacheObj, bool) {
	i := c.keyIdx(key)
	log.Debugf("MultiDiskCache.Get key '%+v' mapped to %+v\n", key, i)
	return (*c)[i].Peek(key)
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

func (c *MultiDiskCache) Keys() []string {
	// TODO Fix this - each cache is an independent LRU, and the below doesn't make sense.
	arr := make([]string, 0)
	for _, cache := range *c {
		arr = append(arr, cache.Keys()...)
	}
	return arr
}

func (c *MultiDiskCache) Capacity() uint64 {
	sum := uint64(0)
	for _, cache := range *c {
		sum += cache.Capacity()
	}
	return sum
}
