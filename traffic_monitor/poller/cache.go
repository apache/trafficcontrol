package poller

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
	"bytes"
	"io"
	"math/rand"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/config"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/handler"
)

type CachePoller struct {
	Config         CachePollerConfig
	ConfigChannel  chan CachePollerConfig
	TickChan       chan uint64
	GlobalContexts map[string]interface{}
	Handler        handler.Handler
}

type PollConfig struct {
	URL      string
	URLv6    string
	Host     string
	Timeout  time.Duration
	Format   string
	PollType string
}

type CachePollerConfig struct {
	Urls            map[string]PollConfig
	Interval        time.Duration
	NoKeepAlive     bool
	PollingProtocol config.PollingProtocol
}

// NewCache creates and returns a new CachePoller.
// If tick is false, CachePoller.TickChan() will return nil.
func NewCache(
	tick bool,
	handler handler.Handler,
	cfg config.Config,
	appData config.StaticAppData,
) CachePoller {
	var tickChan chan uint64
	if tick {
		tickChan = make(chan uint64)
	}
	return CachePoller{
		TickChan:      tickChan,
		ConfigChannel: make(chan CachePollerConfig),
		Config: CachePollerConfig{
			PollingProtocol: cfg.CachePollingProtocol,
		},
		GlobalContexts: GetGlobalContexts(cfg, appData),
		Handler:        handler,
	}
}

var pollNum uint64

type CachePollInfo struct {
	NoKeepAlive     bool
	Interval        time.Duration
	ID              string
	PollingProtocol config.PollingProtocol
	PollConfig
}

func (p CachePoller) Poll() {
	killChans := map[string]chan<- struct{}{}
	for newConfig := range p.ConfigChannel {
		deletions, additions := diffConfigs(p.Config, newConfig)
		for _, id := range deletions {
			killChan := killChans[id]
			go func() { killChan <- struct{}{} }() // go - we don't want to wait for old polls to die.
			delete(killChans, id)
		}
		for _, info := range additions {
			kill := make(chan struct{})
			killChans[info.ID] = kill

			if _, ok := pollers[info.PollType]; !ok {
				if info.PollType != "" { // don't warn for missing parameters
					log.Warnln("CachePoller.Poll: poll type '" + info.PollType + "' not found, using default poll type '" + DefaultPollerType + "'")
				}
				info.PollType = DefaultPollerType
			}
			pollerObj := pollers[info.PollType]

			pollerCfg := PollerConfig{
				Timeout:     info.Timeout,
				NoKeepAlive: info.NoKeepAlive,
				PollerID:    info.ID,
			}
			pollerCtx := interface{}(nil)
			if pollerObj.Init != nil {
				pollerCtx = pollerObj.Init(pollerCfg, p.GlobalContexts[info.PollType])
			}
			go poller(info.Interval, info.ID, info.PollingProtocol, info.URL, info.URLv6, info.Host, info.Format, p.Handler, pollerObj.Poll, pollerCtx, kill)
		}
		p.Config = newConfig
	}
}

// TODO iterationCount and/or p.TickChan?
func poller(
	interval time.Duration,
	id string,
	pollingProtocol config.PollingProtocol,
	url string,
	url6 string,
	host string,
	format string,
	handler handler.Handler,
	pollFunc PollerFunc,
	pollCtx interface{},
	die <-chan struct{},
) {
	pollSpread := time.Duration(rand.Float64()*float64(interval/time.Nanosecond)) * time.Nanosecond
	time.Sleep(pollSpread)
	tick := time.NewTicker(interval)
	lastTime := time.Now()
	oscillateProtocols := false
	if pollingProtocol == config.Both {
		oscillateProtocols = true
	}
	usingIPv4 := pollingProtocol != config.IPv6Only
	for {
		select {
		case <-tick.C:
			if (usingIPv4 && url == "") || (!usingIPv4 && url6 == "") {
				usingIPv4 = !usingIPv4
				continue
			}

			realInterval := time.Now().Sub(lastTime)
			if realInterval > interval+(time.Millisecond*100) {
				log.Debugf("Intended Duration: %v Actual Duration: %v\n", interval, realInterval)
			}
			lastTime = time.Now()

			pollID := atomic.AddUint64(&pollNum, 1)
			pollFinishedChan := make(chan uint64)
			log.Debugf("poll %v %v start\n", pollID, time.Now())
			pollUrl := url
			if !usingIPv4 {
				pollUrl = url6
			}

			bts, reqEnd, reqTime, err := pollFunc(pollCtx, pollUrl, host, pollID)
			rdr := io.Reader(nil)
			if bts != nil {
				rdr = bytes.NewReader(bts) // TODO change handler to take bytes? Benchmark?
			}

			log.Debugf("poll %v %v poller end\n", pollID, time.Now())
			go handler.Handle(id, rdr, format, reqTime, reqEnd, err, pollID, usingIPv4, pollCtx, pollFinishedChan)

			if oscillateProtocols {
				usingIPv4 = !usingIPv4
			}

			<-pollFinishedChan
		case <-die:
			tick.Stop()
			return
		}
	}
}

// diffConfigs takes the old and new configs, and returns a list of deleted IDs, and a list of new polls to do
func diffConfigs(old CachePollerConfig, new CachePollerConfig) ([]string, []CachePollInfo) {
	deletions := []string{}
	additions := []CachePollInfo{}

	if old.Interval != new.Interval || old.NoKeepAlive != new.NoKeepAlive {
		for id, _ := range old.Urls {
			deletions = append(deletions, id)
		}
		for id, pollCfg := range new.Urls {
			additions = append(additions, CachePollInfo{
				Interval:        new.Interval,
				NoKeepAlive:     new.NoKeepAlive,
				ID:              id,
				PollingProtocol: new.PollingProtocol,
				PollConfig:      pollCfg,
			})
		}
		return deletions, additions
	}

	for id, oldPollCfg := range old.Urls {
		newPollCfg, newIdExists := new.Urls[id]
		if !newIdExists {
			deletions = append(deletions, id)
		} else if newPollCfg != oldPollCfg {
			deletions = append(deletions, id)
			additions = append(additions, CachePollInfo{
				Interval:        new.Interval,
				NoKeepAlive:     new.NoKeepAlive,
				ID:              id,
				PollingProtocol: new.PollingProtocol,
				PollConfig:      newPollCfg,
			})
		}
	}

	for id, newPollCfg := range new.Urls {
		_, oldIdExists := old.Urls[id]
		if !oldIdExists {
			additions = append(additions, CachePollInfo{
				Interval:        new.Interval,
				NoKeepAlive:     new.NoKeepAlive,
				ID:              id,
				PollingProtocol: new.PollingProtocol,
				PollConfig:      newPollCfg,
			})
		}
	}

	return deletions, additions
}

func stacktrace() []byte {
	initialBufSize := 1024
	buf := make([]byte, initialBufSize)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, len(buf)*2)
	}
}
