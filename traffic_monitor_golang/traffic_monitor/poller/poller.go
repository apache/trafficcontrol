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
	"math/rand"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/fetcher"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/handler"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/towrap" // TODO move to common
)

type Poller interface {
	Poll()
}

type HttpPoller struct {
	Config          HttpPollerConfig
	ConfigChannel   chan HttpPollerConfig
	FetcherTemplate fetcher.HttpFetcher // FetcherTemplate has all the constant settings, and is copied to create fetchers with custom HTTP client timeouts.
	TickChan        chan uint64
}

type PollConfig struct {
	URL     string
	Host    string
	Timeout time.Duration
	Handler handler.Handler
}

type HttpPollerConfig struct {
	Urls        map[string]PollConfig
	Interval    time.Duration
	NoKeepAlive bool
}

// NewHTTP creates and returns a new HttpPoller.
// If tick is false, HttpPoller.TickChan() will return nil.
func NewHTTP(
	interval time.Duration,
	tick bool,
	httpClient *http.Client,
	fetchHandler handler.Handler,
	userAgent string,
) HttpPoller {
	var tickChan chan uint64
	if tick {
		tickChan = make(chan uint64)
	}
	return HttpPoller{
		TickChan:      tickChan,
		ConfigChannel: make(chan HttpPollerConfig),
		Config: HttpPollerConfig{
			Interval: interval,
		},
		FetcherTemplate: fetcher.HttpFetcher{
			Handler:   fetchHandler,
			Client:    httpClient,
			UserAgent: userAgent,
		},
	}
}

type MonitorCfg struct {
	CDN string
	Cfg tc.TrafficMonitorConfigMap
}

type MonitorConfigPoller struct {
	Session          towrap.ITrafficOpsSession
	SessionChannel   chan towrap.ITrafficOpsSession
	ConfigChannel    chan MonitorCfg
	OpsConfigChannel chan handler.OpsConfig
	Interval         time.Duration
	IntervalChan     chan time.Duration
	OpsConfig        handler.OpsConfig
}

// Creates and returns a new HttpPoller.
// If tick is false, HttpPoller.TickChan() will return nil
func NewMonitorConfig(interval time.Duration) MonitorConfigPoller {
	return MonitorConfigPoller{
		Interval:       interval,
		SessionChannel: make(chan towrap.ITrafficOpsSession),
		// ConfigChannel MUST have a buffer size 1, to make the nonblocking writeConfig work
		ConfigChannel:    make(chan MonitorCfg, 1),
		OpsConfigChannel: make(chan handler.OpsConfig),
		IntervalChan:     make(chan time.Duration),
	}
}

// writeConfig writes the given config to the Config chan. This is nonblocking, and immediately returns.
// Because readers only ever want the latest config, if nobody has read the previous write, we remove it. Since the config chan is buffered size 1, this function is therefore asynchronous.
func (p MonitorConfigPoller) writeConfig(cfg MonitorCfg) {
	for {
		select {
		case p.ConfigChannel <- cfg:
			return // return after successfully writing.
		case <-p.ConfigChannel:
			// if the channel buffer was full, read, then loop and try to write again
		}
	}
}

func (p MonitorConfigPoller) Poll() {
	tick := time.NewTicker(p.Interval)
	defer tick.Stop()
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("MonitorConfigPoller panic: %v\n", err)
		} else {
			log.Errorf("MonitorConfigPoller failed without panic\n")
		}
		os.Exit(1) // The Monitor can't run without a MonitorConfigPoller
	}()
	for {
		// Every case MUST be asynchronous and non-blocking, to prevent livelocks. If a chan must be written to, it must either be buffered AND remove existing values, or be written to in a goroutine.
		select {
		case opsConfig := <-p.OpsConfigChannel:
			log.Infof("MonitorConfigPoller: received new opsConfig: %v\n", opsConfig)
			p.OpsConfig = opsConfig
		case session := <-p.SessionChannel:
			log.Infof("MonitorConfigPoller: received new session: %v\n", session)
			p.Session = session
		case i := <-p.IntervalChan:
			if i == p.Interval {
				continue
			}
			log.Infof("MonitorConfigPoller: received new interval: %v\n", i)
			if i < 0 {
				log.Errorf("MonitorConfigPoller: received negative interval: %v; ignoring\n", i)
				continue
			}
			p.Interval = i
			tick.Stop()
			tick = time.NewTicker(p.Interval)
		case <-tick.C:
			if p.Session == nil || p.OpsConfig.CdnName == "" {
				log.Warnln("MonitorConfigPoller: skipping this iteration, Session is nil")
				continue
			}
			monitorConfig, err := p.Session.TrafficMonitorConfigMap(p.OpsConfig.CdnName)
			if err != nil {
				log.Errorf("MonitorConfigPoller: %s\n %v\n", err, monitorConfig)
				continue
			}
			p.writeConfig(MonitorCfg{CDN: p.OpsConfig.CdnName, Cfg: *monitorConfig})
		}
	}
}

var debugPollNum uint64

type HTTPPollInfo struct {
	NoKeepAlive bool
	Interval    time.Duration
	Timeout     time.Duration
	ID          string
	URL         string
	Host        string
	Handler     handler.Handler
}

func (p HttpPoller) Poll() {
	// iterationCount := uint64(0)
	// iterationCount++ // on tick<:
	// case p.TickChan <- iterationCount:
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

			fetcher := p.FetcherTemplate
			if info.Timeout != 0 || info.NoKeepAlive { // if the timeout isn't explicitly set, use the template value.
				c := *fetcher.Client
				fetcher.Client = &c // copy the client, so we don't change other fetchers.
				if info.Timeout != 0 {
					fetcher.Client.Timeout = info.Timeout
				}
				if info.NoKeepAlive {
					transportI := http.DefaultTransport
					transport, ok := transportI.(*http.Transport)
					if !ok {
						log.Errorf("failed to set NoKeepAlive for '%v': http.DefaultTransport expected type *http.Transport actual %T\n", info.URL, transportI)
					} else {
						transport.DisableKeepAlives = info.NoKeepAlive
						fetcher.Client.Transport = transport
						log.Infof("Setting transport.DisableKeepAlives %v for %v\n", transport.DisableKeepAlives, info.URL)
					}
				}
			}
			go poller(info.Interval, info.ID, info.URL, info.Host, fetcher, kill)
		}
		p.Config = newConfig
	}
}

func mustDie(die <-chan struct{}) bool {
	select {
	case <-die:
		return true
	default:
	}
	return false
}

// TODO iterationCount and/or p.TickChan?
func poller(interval time.Duration, id string, url string, host string, fetcher fetcher.Fetcher, die <-chan struct{}) {
	pollSpread := time.Duration(rand.Float64()*float64(interval/time.Nanosecond)) * time.Nanosecond
	time.Sleep(pollSpread)
	tick := time.NewTicker(interval)
	lastTime := time.Now()
	for {
		select {
		case <-tick.C:
			realInterval := time.Now().Sub(lastTime)
			if realInterval > interval+(time.Millisecond*100) {
				log.Debugf("Intended Duration: %v Actual Duration: %v\n", interval, realInterval)
			}
			lastTime = time.Now()

			pollId := atomic.AddUint64(&debugPollNum, 1)
			pollFinishedChan := make(chan uint64)
			log.Debugf("poll %v %v start\n", pollId, time.Now())
			go fetcher.Fetch(id, url, host, pollId, pollFinishedChan) // TODO persist fetcher, with its own die chan?
			<-pollFinishedChan
		case <-die:
			tick.Stop()
			return
		}
	}
}

// diffConfigs takes the old and new configs, and returns a list of deleted IDs, and a list of new polls to do
func diffConfigs(old HttpPollerConfig, new HttpPollerConfig) ([]string, []HTTPPollInfo) {
	deletions := []string{}
	additions := []HTTPPollInfo{}

	if old.Interval != new.Interval || old.NoKeepAlive != new.NoKeepAlive {
		for id, _ := range old.Urls {
			deletions = append(deletions, id)
		}
		for id, pollCfg := range new.Urls {
			additions = append(additions, HTTPPollInfo{
				Interval:    new.Interval,
				NoKeepAlive: new.NoKeepAlive,
				ID:          id,
				URL:         pollCfg.URL,
				Host:        pollCfg.Host,
				Timeout:     pollCfg.Timeout,
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
			additions = append(additions, HTTPPollInfo{
				Interval:    new.Interval,
				NoKeepAlive: new.NoKeepAlive,
				ID:          id,
				URL:         newPollCfg.URL,
				Host:        newPollCfg.Host,
				Timeout:     newPollCfg.Timeout,
			})
		}
	}

	for id, newPollCfg := range new.Urls {
		_, oldIdExists := old.Urls[id]
		if !oldIdExists {
			additions = append(additions, HTTPPollInfo{
				Interval:    new.Interval,
				NoKeepAlive: new.NoKeepAlive,
				ID:          id,
				URL:         newPollCfg.URL,
				Host:        newPollCfg.Host,
				Timeout:     newPollCfg.Timeout,
			})
		}
	}

	return deletions, additions
}
