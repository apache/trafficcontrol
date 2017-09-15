// This package provides a simple LRU cache. It is based on the
// LRU implementation in groupcache:
// https://github.com/golang/groupcache/tree/master/lru
package lru

import (
	"sync"

	"github.com/hashicorp/golang-lru/simplelru"
)

// CacheStr is a thread-safe fixed size LRU cache with string keys
type CacheStr struct {
	lru  *simplelru.LRUStr
	lock sync.RWMutex
}

// New creates an LRU of the given size
func NewStr(size int) (*CacheStr, error) {
	return NewStrWithEvict(size, nil)
}

// NewWithEvict constructs a fixed size cache with the given eviction
// callback.
func NewStrWithEvict(size int, onEvicted func(key string, value interface{})) (*CacheStr, error) {
	return NewStrLargeWithEvict(uint64(size), onEvicted)
}

// NewLargeWithEvict constructs a fixed size cache with the given eviction
// callback.
func NewStrLargeWithEvict(size uint64, onEvicted func(key string, value interface{})) (*CacheStr, error) {
	lru, err := simplelru.NewLRUStrLarge(size, simplelru.EvictCallbackStr(onEvicted))
	if err != nil {
		return nil, err
	}
	c := &CacheStr{
		lru: lru,
	}
	return c, nil
}

// Purge is used to completely clear the cache
func (c *CacheStr) Purge() {
	c.lock.Lock()
	c.lru.Purge()
	c.lock.Unlock()
}

// Add adds a value to the cache.  Returns true if an eviction occurred.
func (c *CacheStr) Add(key string, value interface{}) bool {
	return c.AddSize(key, value, 1)
}

// AddSize adds a value to the cache.  Returns true if an eviction occurred.
func (c *CacheStr) AddSize(key string, value interface{}, size uint64) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.lru.AddSize(key, value, size)
}

// Get looks up a key's value from the cache.
func (c *CacheStr) Get(key string) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.lru.Get(key)
}

// Check if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (c *CacheStr) Contains(key string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Contains(key)
}

// Returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *CacheStr) Peek(key string) (interface{}, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Peek(key)
}

// ContainsOrAdd checks if a key is in the cache  without updating the
// recent-ness or deleting it for being stale,  and if not, adds the value.
// Returns whether found and whether an eviction occurred.
func (c *CacheStr) ContainsOrAdd(key string, value interface{}) (ok, evict bool) {
	return c.ContainsOrAddSize(key, value, 1)
}

// ContainsOrAddSize checks if a key is in the cache  without updating the
// recent-ness or deleting it for being stale,  and if not, adds the value.
// Returns whether found and whether an eviction occurred.
func (c *CacheStr) ContainsOrAddSize(key string, value interface{}, size uint64) (ok, evict bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.lru.Contains(key) {
		return true, false
	} else {
		evict := c.lru.AddSize(key, value, size)
		return false, evict
	}
}

// Remove removes the provided key from the cache.
func (c *CacheStr) Remove(key string) {
	c.lock.Lock()
	c.lru.Remove(key)
	c.lock.Unlock()
}

// RemoveOldest removes the oldest item from the cache.
func (c *CacheStr) RemoveOldest() {
	c.lock.Lock()
	c.lru.RemoveOldest()
	c.lock.Unlock()
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *CacheStr) Keys() []string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Keys()
}

// Len returns the number of items in the cache.
func (c *CacheStr) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Len()
}

// Size return the size of the cache.
func (c *CacheStr) Size() uint64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Size()
}
