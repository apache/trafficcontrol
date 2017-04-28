package grove

// import (
// 	"sync"
// 	"time"
// )

// // KeyLimiter is a simultaneous request limiter. It only allows a given number of simultaneous function calls per key. The function given should not return until its task is completed, that is, it should not call a goroutine and return before completing. That would defeat the purpose of a limiter.
// type KeyLimiter interface {
// 	// Limit calls the given function, allowing only the Limiter's maximum number of simultaneous calls for the given key.
// 	Limit(key string, f func())
// }

// // SoftLimiter is a fast, soft limiter. It may exceed the given max, but should generally fall within it. In exchange, it's faster than possible for a hard limiter.
// type SoftLimiter struct {
// 	reqMax          uint64
// 	requestors      map[string]*uint64
// 	requestorsM     sync.RWMutex
// 	lastCleanupM    sync.RWMutex
// 	lastCleanup     time.Time
// 	cleanupInterval time.Duration
// }

// // NewSoftLimiter creates a new SoftLimiter. The max is the soft limit of simultaneous functions to allow. The expectedExecutionMax is the maximum time any function is expected to take to finish. If the expectedExecutionMax is too low, the soft limit will be exceeded by more; if it's too high, more memory than necessary will be used.
// func NewSoftLimiter(max uint64, expectedExecutionMax time.Duration) KeyLimiter {
// 	return &SoftLimiter{
// 		reqMax:          max,
// 		requestors:      map[string]*int{},
// 		lastCleanup:     time.Now(),
// 		cleanupInterval: time.Second, // TODO make configurable? Make dynamically average function length?
// 	}
// }

// func (l *SoftLimiter) Limit(key string, f func()) {

// 	// !! if you don't understand this, don't change it!! Lock-free programming is perilous.

// 	l.requestorsM.RLock()
// 	requestorCount, ok := l.requestors[key]
// 	l.requestorsM.RUnlock()
// 	if !ok {
// 		newCount := 0
// 		l.requestorsM.Lock()
// 		if existingCount, ok := l.requestors[key]; !ok { // check if someone inserted after we RUnlocked.
// 			l.requestors[key] = &newCount
// 			requestorCount = &newCount
// 		} else {
// 			requestorCount = existingCount
// 		}
// 		l.requestorsM.Unlock()
// 	}

// 	for {
// 		currCount := atomic.LoadUint64(requestorCount)
// 		if currCount > l.reqMax {
// 			runtime.Gosched()
// 			continue
// 		}
// 		ok := atomic.CompareAndSwapUint64(requestorCount, currCount, currCount+1)
// 		if !ok {
// 			continue
// 		}
// 		break
// 	}

// 	f()

// 	atomic.AddUint64(requestorCount, ^uint64(0)) // decrement
// 	cleanup()
// }

// // cleanup swaps the requestors map with a new one, effectively deleting old URLs. This allows us to delete, without guaranteeing there are no other concurrent requests adding to the counters in the map, which in turn allows us to use fast atomics instead of slow mutexes.
// // This has the side-effect of making reqMax a soft limit, not a hard limit. For example if the max is 10, and there are 5 simultaneous requests when cleanup() is called and the map is swapped, and 10 more immediately come in, there are now 15 simultaneous requests, despite the "max" of 10. If the map is again swapped and 10 more come in, there are now 25. But in practice, as long as cleanup is called on a time-basis rather than a request-basis, it shouldn't greatly exceed the max in practice.
// func (l *SoftLimiter) cleanup() {

// 	// !! if you don't understand this, don't change it!! Lock-free programming is perilous.

// 	l.lastCleanupM.Rlock()
// 	lastCleanup = l.lastCleanup
// 	l.lastCleanupM.RUnlock()
// 	if time.Since(lastCleanup) < l.cleanupInterval {
// 		return
// 	}
// 	l.requestorsM.Lock()
// 	l.lastCleanupM.Lock()
// 	requestors = map[string]*int{}
// 	lastCleanup := time.Now()
// 	l.lastCleanupM.Unlock()
// 	l.requestorsM.Unlock()
// }

// // // RateLimitedServe will only make one request to the origin for all simultaneous connections. The key should be a unique request identifier, such as the method, URL, and query parameters.
// // // For example, if a million simultaneous requests occur for the same key which is not in the cache, only one request will be made to the origin, and the other requests will block until the original succeeds.
// // // TODO add a configurable number of simultaneous requets.
// // // TODO put the above in the documentation
// // func RateLimitedServe(
// // 	parent http.HandlerFunc,
// // 	cache Cache,
// // 	w http.ResponseWriter,
// // 	r *http.Request,
// // 	reqTime time.Time,
// // 	key string,
// // ) {

// // 	// TODO: queue readers

// // 	teeWriter := NewHTTPResponseWriterTee(w)
// // 	reqHeader := http.Header{}
// // 	copyHeader(r.Header, &reqHeader) // copy before ServeHTTP invalidates the request

// // 	h.parent.ServeHTTP(teeWriter, r)
// // 	respTime := time.Now() // TODO get response time as soon as it's returned. This is used to estimate latency per RFC 7234, it needs to be _immediately_ after the origin responds.

// // 	// signal readers

// // 	h.TryCache(key, reqHeader, teeWriter.Bytes, teeWriter.Code, teeWriter.WrittenHeader, reqTime, respTime)
// // }

// // HardLimiter allows simultaneous calls to Limit, and will only execute the given maximum simultaneously, causing others to wait until fewer than the maximum are executing before being called.
// type HardLimiter struct {
// 	reqMax      uint64
// 	requestors  map[string]*int
// 	requestorsM sync.RWMutex
// }

// // NewHardLimiter creates a new HardLimiter. The max is the limit of simultaneous functions to allow to execute.
// func NewHardLimiter(max uint64) KeyLimiter {
// 	return &SoftLimiter{
// 		reqMax:      max,
// 		requestoros: map[string]*int{},
// 	}
// }

// // Limit calls the given function, allowing only the RateLimiter's maximum number of simultaneous calls for the given key.
// func (l *HardLimiter) Limit(key string, f func()) {
// 	// Do a readlock first, which is faster than a writelock. This optimizes for the case of exceeding the rate limit.
// 	// TODO add config option to skip this, i.e. optimizing for few simultaneous connections
// 	// TODO add config option to sleep(0) and/or backoff (if exponential there should be a low cap, e.g. 1s)
// 	for {
// 		l.requestorsM.RLock()
// 		requestorCount := l.requestors[key]
// 		l.requestorsM.RUnlock()
// 		if requestorCount > l.reqMax {
// 			runtime.Gosched()
// 			continue
// 		}
// 		break
// 	}

// 	requestorCount = l.reqMax + 1
// 	for {
// 		l.requestorsM.Lock()
// 		requestorCount = l.requestors[key]
// 		if requestorCount > l.reqMax {
// 			l.requestorsM.Unlock()
// 			runtime.Gosched()
// 			continue
// 		}
// 		l.requestors[key] = requestorCount + 1
// 		l.requestorsM.Unlock()
// 		break
// 	}

// 	f()

// 	requestorsM.Lock()
// 	requestorCount := requestors[key]
// 	requestors[key] = requestorCount - 1
// 	if requestors[key] == 0 {
// 		delete(l.requestors, key)
// 	}
// 	requestorsM.Unlock()
// }
