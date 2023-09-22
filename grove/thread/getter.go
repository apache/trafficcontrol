package thread

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
	"sync"

	cacheobj "github.com/apache/trafficcontrol/v8/grove/cacheobj"
)

type Getter interface {
	Get(key string, actualGet func() *cacheobj.CacheObj, canUse func(*cacheobj.CacheObj) bool, reqID uint64) (*cacheobj.CacheObj, uint64)
}

type GetterResp struct {
	CacheObj *cacheobj.CacheObj
	GetReqID uint64
}

func NewGetter() Getter {
	return &getter{waiters: map[string][]chan GetterResp{}}
}

// getter implements Getter, and does a fan-in so only one real request is made to the parent at any given time, and then that object is given to all concurrent requesters.
//
// When a request for a key with no-one currently processing it comes in, that requestor becomes the Author. Subsequent requests become Waiters.
// The initial Author inserts a new constructed (but empty) slice into the waiters map.
// Then, when other requests come in, they see that waiters[key] exists, and add themselves to it, and block reading from their chan.
// Then, when the Author gets its response, it iterates over the Waiters and sends the response to all of them, at the same time (with the same lock, atomically) clearing the waiters for the next request that comes in.
//
// If the Author response can't be used, all Waiters make their own requests.
// Note this assumes an uncacheable response for one request is likely uncacheable for all, and it's faster and less load on the origin if so.
// If it's likely the author request is uncacheable, but a different waiter is cacheable for all other waiters, this will be more network, more origin load, and more work. If that's the case for you, consider creating another type that fulfills the Getter interface, and making the Getter configurable.
type getter struct {
	// waiters is a map of cache keys to chans for getters.
	waiters  map[string][]chan GetterResp
	waitersM sync.Mutex
}

func (g *getter) Get(key string, actualGet func() *cacheobj.CacheObj, canUse func(*cacheobj.CacheObj) bool, reqID uint64) (*cacheobj.CacheObj, uint64) {
	isAuthor := false
	// Buffered for performance, so the author can iterate over all wait chans without blocking.
	// Note this is unused if isAuthor becomes true.
	getChan := make(chan GetterResp, 1)

	g.waitersM.Lock()
	if _, ok := g.waiters[key]; !ok {
		isAuthor = true
		g.waiters[key] = []chan GetterResp{}
	} else {
		g.waiters[key] = append(g.waiters[key], getChan)
	}
	g.waitersM.Unlock()

	if isAuthor {
		obj := actualGet()
		waitResp := GetterResp{CacheObj: obj, GetReqID: reqID}

		g.waitersM.Lock()
		for _, waitChan := range g.waiters[key] {
			waitChan <- waitResp
		}
		delete(g.waiters, key)
		g.waitersM.Unlock()

		return obj, reqID
	}

	if waitResp := <-getChan; canUse(waitResp.CacheObj) {
		return waitResp.CacheObj, waitResp.GetReqID
	}

	// if the Author response can't be used, all Waiters make their own requests
	return actualGet(), reqID
}
