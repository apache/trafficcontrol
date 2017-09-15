package thread

import (
	"sync"

	cacheobj "github.com/apache/incubator-trafficcontrol/grove/cacheobj"
)

type Getter interface {
	Get(key string, actualGet func() *cacheobj.CacheObj, canUse func(*cacheobj.CacheObj) bool) *cacheobj.CacheObj
}

type CanUseCachedFunc func(obj *cacheobj.CacheObj) bool

type getter struct {
	objs map[string][]GetterObj
	m    sync.Mutex
}

type GetterObj struct {
	c      chan *cacheobj.CacheObj
	canUse CanUseCachedFunc
}

func NewGetter() Getter {
	return &getter{objs: map[string][]GetterObj{}}
}

func (g *getter) Get(key string, actualGet func() *cacheobj.CacheObj, canUse func(*cacheobj.CacheObj) bool) *cacheobj.CacheObj {
	getChan := g.atomicGetOrCreateGetter(key, actualGet, canUse)
	return <-getChan
}

// GetterBuffer is the size of the getter chan. This is buffered for the getter, not for clients. We don't want the getter to block, wait for a client to read, write another, and block again. We want to write a batch of objects for clients to process, to reduce the back-and-forth stutter between goroutines.
// TODO test increase. Make configurable?
const GetterBuffer = 10

func (g *getter) atomicGetOrCreateGetter(key string, actualGet func() *cacheobj.CacheObj, canUse func(*cacheobj.CacheObj) bool) <-chan *cacheobj.CacheObj {
	g.m.Lock()
	defer g.m.Unlock()
	getterObjs, ok := g.objs[key]
	if !ok {
		getterObjs = make([]GetterObj, 0, 1)
		go gogetter(g, key, actualGet, canUse)
	}
	c := make(chan *cacheobj.CacheObj, GetterBuffer)
	getterObj := GetterObj{c: c, canUse: canUse}
	getterObjs = append(getterObjs, getterObj)
	g.objs[key] = getterObjs
	return c
}

// atomicPopGetter atomically gets a GetterObj and returns it, and whether it's the last one, i.e. whether the gogetter should exit. This MUST NOT be called by the same gogetter goroutine after returning true.
func (g *getter) atomicPopGetter(key string) (GetterObj, bool) {
	g.m.Lock()
	defer g.m.Unlock()
	objs := g.objs[key]
	if len(objs) == 1 {
		delete(g.objs, key)
		return objs[0], true
	} else {
		g.objs[key] = objs[1:] // if this panics, you called this function after it returned true and there were no more getters
		return objs[0], false
	}
}

func gogetter(g *getter, key string, actualGet func() *cacheobj.CacheObj, canReuse func(*cacheobj.CacheObj) bool) {
	getter, noRemainingGetters := g.atomicPopGetter(key)
	for {
		obj := actualGet()
		// don't need to check if the first can use, since the object isn't 'cached' i.e. given to anyone else yet.
		getter.c <- obj

		// TODO process all getters, and create a list of "getters who couldn't use", and re-queue them? That would optimize for some requests being cachable and some not, for the same key. The current routine is optimized for an uncacheable key being uncacheable for all requestors.
		for {
			if noRemainingGetters {
				return
			}
			getter, noRemainingGetters = g.atomicPopGetter(key)
			if getter.canUse(obj) {
				getter.c <- obj
			} else {
				break
			}
		}
	}
}
