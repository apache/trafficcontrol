package lru

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
	"container/list"
	"sync"
)

type LRU struct {
	l      *list.List
	lElems map[string]*list.Element
	m      sync.RWMutex
}

type listObj struct {
	key  string
	size uint64
}

func NewLRU() *LRU {
	return &LRU{l: list.New(), lElems: map[string]*list.Element{}}
}

// Add adds the key to the LRU, with the given size. Returns the size of the the old size, or 0 if no key existed.
func (c *LRU) Add(key string, size uint64) uint64 {
	c.m.Lock()
	defer c.m.Unlock()
	if elem, ok := c.lElems[key]; ok {
		c.l.MoveToFront(elem)
		oldSize := elem.Value.(*listObj).size
		elem.Value.(*listObj).size = size
		return oldSize
	}
	c.lElems[key] = c.l.PushFront(&listObj{key, size})
	return 0
}

// RemoveOldest returns the key, size, and true if the LRU is nonempty; else false.
func (c *LRU) RemoveOldest() (string, uint64, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	elem := c.l.Back()
	if elem == nil {
		return "", 0, false
	}
	c.l.Remove(elem)
	obj := elem.Value.(*listObj)
	delete(c.lElems, obj.key)
	return obj.key, obj.size, true
}

// Keys returns a string array of the keys
func (c *LRU) Keys() []string {
	c.m.RLock()
	defer c.m.RUnlock()
	arr := make([]string, 0)
	for e := c.l.Back(); e != nil; e = e.Prev() {
		object := e.Value.(*listObj)
		arr = append(arr, object.key)
	}
	return arr
}
