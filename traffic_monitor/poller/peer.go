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
	"net/url"
	"sync/atomic"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/config"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/handler"
)

type PeerPoller struct {
	Config         PeerPollerConfig
	ConfigChannel  chan PeerPollerConfig
	GlobalContexts map[string]interface{}
	Handler        handler.Handler
}

type PeerPollConfig struct {
	URLs     []string
	Timeout  time.Duration
	Format   string
	PollType string
}

func (c PeerPollConfig) Equals(other PeerPollConfig) bool {
	if len(c.URLs) != len(other.URLs) {
		return false
	}
	for i, v := range c.URLs {
		if v != other.URLs[i] {
			return false
		}
	}
	return c.Timeout == other.Timeout && c.Format == other.Format && c.PollType == other.PollType
}

type PeerPollerConfig struct {
	Urls        map[string]PeerPollConfig
	Interval    time.Duration
	NoKeepAlive bool
}

// NewPeer creates and returns a new PeerPoller.
func NewPeer(
	handler handler.Handler,
	cfg config.Config,
	appData config.StaticAppData,
) PeerPoller {
	return PeerPoller{
		ConfigChannel:  make(chan PeerPollerConfig),
		GlobalContexts: GetGlobalContexts(cfg, appData),
		Handler:        handler,
	}
}

type PeerPollInfo struct {
	NoKeepAlive bool
	Interval    time.Duration
	ID          string
	PeerPollConfig
}

func (p PeerPoller) Poll() {
	killChans := map[string]chan<- struct{}{}
	for newConfig := range p.ConfigChannel {
		deletions, additions := diffPeerConfigs(p.Config, newConfig)
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
			go peerPoller(info.Interval, info.ID, info.URLs, info.Format, p.Handler, pollerObj.Poll, pollerCtx, kill)
		}
		p.Config = newConfig
	}
}

func peerPoller(
	interval time.Duration,
	id string,
	urls []string,
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
	urlI := rand.Intn(len(urls)) // start at a random URL index in order to help spread load
	for {
		select {
		case <-tick.C:

			realInterval := time.Now().Sub(lastTime)
			if realInterval > interval+(time.Millisecond*100) {
				log.Debugf("Intended Duration: %v Actual Duration: %v\n", interval, realInterval)
			}
			lastTime = time.Now()

			pollID := atomic.AddUint64(&pollNum, 1)
			pollFinishedChan := make(chan uint64)
			log.Debugf("peer poll %v %v start\n", pollID, time.Now())

			urlString := urls[urlI]
			urlI = (urlI + 1) % len(urls)
			urlParsed, err := url.Parse(urlString)
			if err != nil {
				// this should never happen because TM creates the URL
				log.Errorf("parsing peer poller URL %s: %s", urlString, err.Error())
			}
			host := urlParsed.Host
			bts, reqEnd, reqTime, err := pollFunc(pollCtx, urlString, host, pollID)
			rdr := io.Reader(nil)
			if bts != nil {
				rdr = bytes.NewReader(bts) // TODO change handler to take bytes? Benchmark?
			}

			log.Debugf("peer poll %v %v poller end\n", pollID, time.Now())
			go handler.Handle(id, rdr, format, reqTime, reqEnd, err, pollID, false, pollCtx, pollFinishedChan)

			<-pollFinishedChan
		case <-die:
			tick.Stop()
			return
		}
	}
}

// diffPeerConfigs takes the old and new configs, and returns a list of deleted IDs, and a list of new polls to do
func diffPeerConfigs(old PeerPollerConfig, new PeerPollerConfig) ([]string, []PeerPollInfo) {
	deletions := []string{}
	additions := []PeerPollInfo{}

	if old.Interval != new.Interval || old.NoKeepAlive != new.NoKeepAlive {
		for id, _ := range old.Urls {
			deletions = append(deletions, id)
		}
		for id, pollCfg := range new.Urls {
			additions = append(additions, PeerPollInfo{
				Interval:       new.Interval,
				NoKeepAlive:    new.NoKeepAlive,
				ID:             id,
				PeerPollConfig: pollCfg,
			})
		}
		return deletions, additions
	}

	for id, oldPollCfg := range old.Urls {
		newPollCfg, newIdExists := new.Urls[id]
		if !newIdExists {
			deletions = append(deletions, id)
		} else if !newPollCfg.Equals(oldPollCfg) {
			deletions = append(deletions, id)
			additions = append(additions, PeerPollInfo{
				Interval:       new.Interval,
				NoKeepAlive:    new.NoKeepAlive,
				ID:             id,
				PeerPollConfig: newPollCfg,
			})
		}
	}

	for id, newPollCfg := range new.Urls {
		_, oldIdExists := old.Urls[id]
		if !oldIdExists {
			additions = append(additions, PeerPollInfo{
				Interval:       new.Interval,
				NoKeepAlive:    new.NoKeepAlive,
				ID:             id,
				PeerPollConfig: newPollCfg,
			})
		}
	}

	return deletions, additions
}
