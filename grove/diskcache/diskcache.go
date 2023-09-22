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
	"bytes"
	"encoding/gob"
	"errors"
	"sync/atomic"
	"time"

	"github.com/apache/trafficcontrol/v8/grove/cacheobj"
	"github.com/apache/trafficcontrol/v8/grove/lru"

	"github.com/apache/trafficcontrol/v8/lib/go-log"

	bolt "go.etcd.io/bbolt"
)

type DiskCache struct {
	db           *bolt.DB
	sizeBytes    uint64
	maxSizeBytes uint64
	lru          *lru.LRU
}

const BucketName = "b"

func New(path string, cacheSizeBytes uint64) (*DiskCache, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return nil, errors.New("opening database '" + path + "': " + err.Error())
	}

	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(BucketName)); err != nil {
			return errors.New("creating bucket: " + err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("creating bucket for database '" + path + "': " + err.Error())
	}

	return &DiskCache{db: db, maxSizeBytes: cacheSizeBytes, lru: lru.NewLRU(), sizeBytes: 0}, nil
}

// ResetAfterRestart rebuilds the LRU with an arbirtrary order and sets sizeBytes. This seems crazy, but it is better than doing nothing, sice gc is based on the LRU and sizeBytes. In the future, we may want to periodically sync the LRU to disk, but we'll still need to iterate over all keys in the disk DB to avoid orphaning objects.
// Note: this assumes the LRU is empty. Don't run twice
func (c *DiskCache) ResetAfterRestart() {
	go c.db.View(func(tx *bolt.Tx) error {
		log.Infof("Starting cache recovery from disk for: %s... ", c.db.Path())
		size := 0
		b := tx.Bucket([]byte(BucketName))

		cursor := b.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			c.lru.Add(string(k), uint64(len(v)))
			size += len(v)
		}

		atomic.AddUint64(&c.sizeBytes, uint64(size))
		log.Infof("Cache recovery from disk for %s done (%d bytes). ", c.db.Path(), c.sizeBytes)
		return nil
	})
}

// Add takes a key and value to add. Returns whether an eviction occurred
// The size is taken to fulfill the Cache interface, but the DiskCache doesn't use it.
// Instead, we compute size from the serialized bytes stored to disk.
//
// Note DiskCache.Add does garbage collection in a goroutine, and thus it is not possible to determine eviction without impacting performance. This always returns false.
func (c *DiskCache) Add(key string, val *cacheobj.CacheObj) bool {
	log.Debugf("DiskCache Add CALLED key '%+v' size '%+v'\n", key, val.Size)
	eviction := false

	buf := bytes.Buffer{}
	if err := gob.NewEncoder(&buf).Encode(val); err != nil {
		log.Errorln("DiskCache.Add encoding cache object: " + err.Error())
		return eviction
	}
	valBytes := buf.Bytes()

	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		if b == nil {
			return errors.New("bucket does not exist")
		}
		return b.Put([]byte(key), valBytes)
	})
	if err != nil {
		log.Errorln("DiskCache.Add inserting '" + key + "' in database: " + err.Error())
		return eviction
	}

	c.lru.Add(key, uint64(len(valBytes)))

	newSizeBytes := atomic.AddUint64(&c.sizeBytes, uint64(len(valBytes)))
	if newSizeBytes > c.maxSizeBytes {
		go c.gc(newSizeBytes)
	}

	log.Debugf("DiskCache Add SUCCESS key '%+v' size '%+v' valBytes '%+v' c.sizeBytes '%+v'\n", key, val.Size, len(valBytes), c.sizeBytes)
	return eviction
}

// gc does garbage collection, deleting stored entries until the DiskCache's size is less than maxSizeBytes. This is threadsafe, and should be called in a goroutine to avoid blocking the caller.
// The given cacheSizeBytes must be `c.Size()`; it's passed here, because gc should be called immediately after an insert updates the size, so it saves an atomic instruction to pass rather than calling Size() again.
func (c *DiskCache) gc(cacheSizeBytes uint64) {
	for cacheSizeBytes > c.maxSizeBytes {
		log.Debugf("DiskCache.gc cacheSizeBytes %+v > c.maxSizeBytes %+v\n", cacheSizeBytes, c.maxSizeBytes)
		key, sizeBytes, exists := c.lru.RemoveOldest() // TODO change lru to use strings
		if !exists {
			// should never happen
			log.Errorf("sizeBytes %v > %v maxSizeBytes, but LRU is empty!? Setting cache size to 0!\n", cacheSizeBytes, c.maxSizeBytes)
			atomic.StoreUint64(&c.sizeBytes, 0)
			return
		}

		log.Debugf("DiskCache.gc deleting key '" + key + "'")
		err := c.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(BucketName))
			if b == nil {
				return errors.New("bucket does not exist")
			}
			b.Delete([]byte(key))

			return b.Delete([]byte(key))
		})
		if err != nil {
			log.Errorln("removing '" + key + "' from cache: " + err.Error())
		}

		cacheSizeBytes = atomic.AddUint64(&c.sizeBytes, ^uint64(sizeBytes-1)) // subtract sizeBytes
	}
}

// Get takes a key, and returns its value, and whether it was found, and updates the lru-ness and hitcount
func (c *DiskCache) Get(key string) (*cacheobj.CacheObj, bool) {
	val, found := c.Peek(key)
	if found {
		c.lru.Add(key, val.Size) // TODO directly call c.ll.MoveToFront
		log.Debugln("DiskCache.Get getting '" + key + "' from cache and updating LRU")
		atomic.AddUint64(&val.HitCount, 1)
		return val, true
	}
	return nil, false

}

// Peek takes a key, and returns its value, and whether it was found, without changing the lru-ness or hitcount
func (c *DiskCache) Peek(key string) (*cacheobj.CacheObj, bool) {
	log.Debugln("DiskCache.Get key '" + key + "'")
	valBytes := []byte(nil)

	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		if b == nil {
			return errors.New("bucket does not exist")
		}
		valBytes = b.Get([]byte(key))
		return nil
	})
	if err != nil {
		log.Errorln("DiskCache.Peek getting '" + key + "' from cache: " + err.Error())
		return nil, false
	}

	if valBytes == nil {
		log.Debugln("DiskCache.Peek key '" + key + "' CACHE MISS")
		return nil, false
	}

	buf := bytes.NewBuffer(valBytes)
	val := cacheobj.CacheObj{}
	if err := gob.NewDecoder(buf).Decode(&val); err != nil {
		log.Errorln("DiskCache.Peek decoding '" + key + "' from cache: " + err.Error())
		return nil, false
	}

	log.Debugln("DiskCache.Peek key '" + key + "' CACHE HIT")
	return &val, true
}

func (c *DiskCache) Size() uint64 {
	return atomic.LoadUint64(&c.sizeBytes)
}

func (c *DiskCache) Close() {
	c.db.Close()
}

func (c *DiskCache) Keys() []string {
	return c.lru.Keys()

}

func (c *DiskCache) Capacity() uint64 {
	return c.maxSizeBytes
}
