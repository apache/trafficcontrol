package middleware

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
)

type ReqObj struct {
	Body    []byte
	Code    int
	Headers http.Header
}

type CacheKey string

// RWR contains the in-flight read-while-writer state.
//
// It has a map of cache keys, to queues.
// Each unique cache key has potentially multiple queues of readers, each for a different start time.
//
// Ordinarily, if a client doesn't send any Cache-Control, it will use the oldest (and thus soonest-to-return) queue.
//
// But, if a client sends a 'Cache-Control: max-age=', for long requests, that max could be longer than the request has been going on.
//
// In practical terms, ClientA could do a GET, ClientB could do a POST, and then while ClientA's request and Read-While-Writer queue is still ongoing, ClientB then does a GET max-age=x.
//
// In that scenario, we need to start a new request queue, to avoid giving ClientB data from before its POST.
//
// Hence, each cache key can actually have multiple queues, each with their own start time.
// A new requestor with no max-age will simply get the first.
// But a requestor with a max-age will get the earliest older than its requested max age.
//
type RWR struct {
	reqs        map[CacheKey][]*RWRReqQueue
	m           *sync.Mutex
	nextQueueID uint64
}

func NewRWR() *RWR {
	// now := time.Now()
	// inGC := int32(0)
	// nowP := (unsafe.Pointer)(&now)
	return &RWR{
		reqs: map[CacheKey][]*RWRReqQueue{},
		m:    &sync.Mutex{},
		// CacheTime: cacheTime,
		// InGC:      &inGC,
		// LastGC:    &nowP,
	}
}

type RWRReqQueue struct {
	Start   time.Time
	Waiters []chan<- ReqObj
	ID      uint64 // queues need an ID, so we know which queue to delete
}

// GetOrMakeQueue returns a reader from the queue if it exists, or creates it if it doesn't and returns a QueueWriter.
// Either the returned reader chan or queueWriter will be nil, but never both.
// If a reader chan is returned, the caller must read from it, and write its object to its client.
// If a QueueWriter is returned, the caller must execute its request as normal, write its data to its client, and then write that same data to the QueueWriter.
//
// If another request is happening, chan ReqObj will not be nil. In which case, the caller must read from it to get the object when it's ready.
// If the chan is nil, then no queue existed, and the caller must do the request itself, and then write what it gets via rwr.WriteQueue.
//
func (rwr *RWR) GetOrMakeQueue(cacheKey CacheKey, thisReqStart time.Time, maxAge time.Duration) (RWRQueueReader, RWRQueueWriter) {
	newQueue := &RWRReqQueue{
		ID:      atomic.AddUint64(&rwr.nextQueueID, 1),
		Start:   thisReqStart,
		Waiters: []chan<- ReqObj{},
	}
	reqChan := make(chan ReqObj, 1) // buffer 1, so the WriteQueue writes don't block

	log.Infoln("RWR GetOrMakeQueue locking")
	rwr.m.Lock()
	log.Infoln("RWR GetOrMakeQueue locked")

	log.Infoln("RWR GetOrMakeQueue starting for loop")
	for _, reqQueue := range rwr.reqs[cacheKey] {
		if time.Since(reqQueue.Start) > maxAge {
			continue
		}
		reqQueue.Waiters = append(reqQueue.Waiters, reqChan)
		log.Infoln("RWR GetOrMakeQueue unlocking and returning reqChan")
		rwr.m.Unlock()
		log.Infoln("RWR GetOrMakeQueue in-for unlocked")

		return &rwrQueueReader{
			ch:    reqChan,
			start: reqQueue.Start,
		}, nil
	}
	log.Infoln("RWR GetOrMakeQueue done for loop, unlocking and returning queueWriter")

	rwr.reqs[cacheKey] = append(rwr.reqs[cacheKey], newQueue)
	rwr.m.Unlock()
	log.Infoln("RWR GetOrMakeQueue done-for unlocked")

	return nil, &rwrQueueWriter{
		rwr:      rwr,
		cacheKey: cacheKey,
		queue:    newQueue,
	}
}

type RWRQueueReader interface {
	// Read blocks until the original request finishes and the original request handler writes to its RWRQueueWriter.
	// Read may only be called once. Subsequent calls will block forever.
	Read() ReqObj
	// Start is the time that the original requestor which created the queue started.
	Start() time.Time
}

type rwrQueueReader struct {
	ch    <-chan ReqObj
	start time.Time
}

func (qr *rwrQueueReader) Start() time.Time { return qr.start }
func (qr *rwrQueueReader) Read() ReqObj     { return <-qr.ch }

type RWRQueueWriter interface {
	Write(obj ReqObj)
}

type rwrQueueWriter struct {
	rwr      *RWR
	cacheKey CacheKey
	queue    *RWRReqQueue
}

func (qw *rwrQueueWriter) Write(obj ReqObj) {
	log.Infoln("RWR rwrQueueWriter.Write starting, locking")
	qw.rwr.m.Lock()
	log.Infoln("RWR rwrQueueWriter.Write starting, locked")

	queueIndex := -1
	oldQueues := qw.rwr.reqs[qw.cacheKey]
	for i := 0; i < len(oldQueues); i++ {
		if oldQueues[i].ID == qw.queue.ID {
			queueIndex = i
			break
		}
	}
	if queueIndex == -1 {
		log.Errorln("RWR.Write was called, but the queue didn't exist in RWR! This is a critical code error, and should never happen!")
		log.Errorf("RWR rwrQueueWriter.Write crit err our queue id: %v\n", qw.queue.ID)
		for i := 0; i < len(oldQueues); i++ {
			log.Errorf("RWR rwrQueueWriter.Write crit err queue id: %v\n", oldQueues[i].ID)
		}
		log.Errorln("RWR rwrQueueWriter.Write crit err unlocking")
		qw.rwr.m.Unlock()
		log.Errorln("RWR rwrQueueWriter.Write crit err returning")
		return
	}

	for i := queueIndex; i < len(oldQueues)-1; i++ {
		oldQueues[i] = oldQueues[i+1]
	}
	oldQueues = oldQueues[:len(oldQueues)-1]

	// TODO add debug prints of IDs to verify

	if len(oldQueues) == 0 {
		delete(qw.rwr.reqs, qw.cacheKey)
	} else {
		qw.rwr.reqs[qw.cacheKey] = oldQueues
	}

	log.Infoln("RWR rwrQueueWriter.Write removed queue, unlocking")
	qw.rwr.m.Unlock()
	log.Infoln("RWR rwrQueueWriter.Write removed queue, unlocked")

	log.Infof("RWR rwr.WriteQueue writing to %v readers\n", len(qw.queue.Waiters))
	for _, ch := range qw.queue.Waiters {
		ch <- obj
	}
}
