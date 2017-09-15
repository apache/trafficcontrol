package simplelru

import (
	"container/list"
)

// EvictCallback is used to get a callback when a cache entry is evicted
type EvictCallbackStr func(key string, value interface{})

// LRUStr implements a non-thread safe fixed size LRU cache with string keys
type LRUStr struct {
	size, total uint64
	evictList   *list.List
	items       map[string]*list.Element
	onEvict     EvictCallbackStr
}

// entry is used to hold a value in the evictList
type entryStr struct {
	key   string
	value interface{}
	size  uint64
}

// NewLRU constructs a count constrained LRU of the given size
func NewLRUStr(size int, onEvict EvictCallbackStr) (*LRUStr, error) {
	return NewLRUStrLarge(uint64(size), onEvict)
}

// NewLRUofType constructs an LRU of the given type and size
func NewLRUStrLarge(size uint64, onEvict EvictCallbackStr) (*LRUStr, error) {
	c := &LRUStr{
		size:      size,
		total:     0,
		evictList: list.New(),
		items:     make(map[string]*list.Element),
		onEvict:   onEvict,
	}
	return c, nil
}

// Purge is used to completely clear the cache
func (c *LRUStr) Purge() {
	for k, v := range c.items {
		if c.onEvict != nil {
			c.onEvict(k, v.Value.(*entryStr).value)
		}
		delete(c.items, k)
	}
	c.total = 0
	c.evictList.Init()
}

// Add adds a value to the cache.  Returns true if an eviction occurred.
func (c *LRUStr) Add(key string, value interface{}) bool {
	return c.AddSize(key, value, 1)
}

// AddSize adds a value to the cache with an associated size.  Returns true if an eviction occurred.
func (c *LRUStr) AddSize(key string, value interface{}, size uint64) bool {
	// Check for existing item
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		ent.Value.(*entryStr).value = value
		return false
	}

	// Add new item
	ent := &entryStr{key, value, size}
	entry := c.evictList.PushFront(ent)
	c.items[key] = entry
	c.total += size

	evict := c.total > c.size
	// Verify size not exceeded
	for c.total > c.size {
		c.removeOldest()
	}
	return evict
}

// Get looks up a key's value from the cache.
func (c *LRUStr) Get(key string) (value interface{}, ok bool) {
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		return ent.Value.(*entryStr).value, true
	}
	return
}

// Check if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (c *LRUStr) Contains(key string) (ok bool) {
	_, ok = c.items[key]
	return ok
}

// Returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *LRUStr) Peek(key string) (value interface{}, ok bool) {
	if ent, ok := c.items[key]; ok {
		return ent.Value.(*entryStr).value, true
	}
	return nil, ok
}

// Remove removes the provided key from the cache, returning if the
// key was contained.
func (c *LRUStr) Remove(key string) bool {
	if ent, ok := c.items[key]; ok {
		c.removeElement(ent)
		c.total -= ent.Value.(*entryStr).size
		return true
	}
	return false
}

// RemoveOldest removes the oldest item from the cache.
func (c *LRUStr) RemoveOldest() (string, interface{}, bool) {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
		kv := ent.Value.(*entryStr)
		return kv.key, kv.value, true
	}
	return "", nil, false
}

// GetOldest returns the oldest entry
func (c *LRUStr) GetOldest() (string, interface{}, bool) {
	ent := c.evictList.Back()
	if ent != nil {
		kv := ent.Value.(*entryStr)
		return kv.key, kv.value, true
	}
	return "", nil, false
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *LRUStr) Keys() []string {
	keys := make([]string, len(c.items))
	i := 0
	for ent := c.evictList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*entryStr).key
		i++
	}
	return keys
}

// Len returns the number of items in the cache.
func (c *LRUStr) Len() int {
	return c.evictList.Len()
}

// Size returns the total bytes of all items in the cache.
func (c *LRUStr) Size() uint64 {
	return c.total
}

// removeOldest removes the oldest item from the cache.
func (c *LRUStr) removeOldest() {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
	}
}

// removeElement is used to remove a given list element from the cache
func (c *LRUStr) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	kv := e.Value.(*entryStr)
	c.total -= kv.size
	delete(c.items, kv.key)
	if c.onEvict != nil {
		c.onEvict(kv.key, kv.value)
	}
}
