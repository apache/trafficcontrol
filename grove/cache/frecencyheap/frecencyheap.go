package frecencyheap

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
	"container/heap"
	"sync"
	"time"
)

func New(hitPerc float64) *Heap {
	return &Heap{
		l:       []*heapObj{},
		lElems:  map[string]*heapObj{},
		hitPerc: hitPerc,
	}
}

func frecency(o *heapObj, hitPerc float64) uint64 {
	return uint64(float64(o.lastHit.UnixNano()) * (1.0 + (float64(o.hits) * hitPerc)))
}

// Heap implements a "frecency" cache algorithm, considering both recency (LRU) and frequency, when removing an element from the cache.
type Heap struct {
	l       []*heapObj
	lElems  map[string]*heapObj
	hitPerc float64
	m       sync.RWMutex
}

type heapObj struct {
	key     string
	size    uint64
	hits    uint64
	lastHit time.Time
}

func (h Heap) Len() int      { return len(h.l) }
func (h Heap) Swap(i, j int) { h.l[i], h.l[j] = h.l[j], h.l[i] }

func (h Heap) Less(i, j int) bool {
	// TODO verify this isn't backwards, and removing the most-recently-used by accident.
	return frecency(h.l[i], h.hitPerc) < frecency(h.l[j], h.hitPerc)
}

func (h *Heap) Push(x interface{}) {
	h.l = append(h.l, x.(*heapObj))
}

func (h *Heap) Pop() interface{} {
	if len(h.l) == 0 {
		return (*heapObj)(nil)
	}
	old := h.l
	n := len(old)
	x := old[n-1]
	h.l = old[0 : n-1]
	return x
}

// Add adds the key to the Heap, with the given size. Returns the size of the the old size, or 0 if no key existed.
func (c *Heap) Add(key string, size uint64) uint64 {
	c.m.Lock()
	defer c.m.Unlock()
	if elem, ok := c.lElems[key]; ok {
		elem.hits += 1
		elem.lastHit = time.Now()
		heap.Init(c) // TODO determine if iterating to find the element, and calling heap.Fix() is faster.
		oldSize := elem.size
		elem.size = size
		return oldSize
	}
	c.Push(&heapObj{key: key, size: size, hits: 1, lastHit: time.Now()})
	return 0
}

// RemoveOldest returns the key, size, and true if the Heap is nonempty; else false.
func (c *Heap) RemoveOldest() (string, uint64, bool) {
	c.m.Lock()
	defer c.m.Unlock()
	obj := c.Pop().(*heapObj)
	if obj == nil {
		return "", 0, false
	}

	return obj.key, obj.size, true
}

// Keys returns a string array of the keys
func (c *Heap) Keys() []string {
	c.m.RLock()
	defer c.m.RUnlock()
	keys := make([]string, len(c.l))
	for i := 0; i < len(c.l); i++ {
		keys[i] = c.l[i].key
	}
	return keys
}
