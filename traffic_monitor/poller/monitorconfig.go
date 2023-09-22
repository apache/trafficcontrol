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
	"os"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/handler"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/towrap" // TODO move to common
)

type MonitorCfg struct {
	CDN string
	Cfg tc.TrafficMonitorConfigMap
}

type MonitorConfigPoller struct {
	Session          towrap.TrafficOpsSessionThreadsafe
	SessionChannel   chan towrap.TrafficOpsSessionThreadsafe
	ConfigChannel    chan MonitorCfg
	OpsConfigChannel chan handler.OpsConfig
	Interval         time.Duration
	IntervalChan     chan time.Duration
	OpsConfig        handler.OpsConfig
}

// NewMonitorConfig Creates and returns a new MonitorConfigPoller.
func NewMonitorConfig(interval time.Duration) MonitorConfigPoller {
	return MonitorConfigPoller{
		Interval:       interval,
		SessionChannel: make(chan towrap.TrafficOpsSessionThreadsafe),
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
		log.Errorf("%s\n", stacktrace())
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
			if !p.Session.Initialized() || p.OpsConfig.CdnName == "" {
				log.Warnln("MonitorConfigPoller: skipping this iteration, Session is nil")
				continue
			}
			monitorConfig, err := p.Session.TrafficMonitorConfigMap(p.OpsConfig.CdnName)
			if err != nil {
				log.Errorf("MonitorConfigPoller: %s\n %v\n", err, monitorConfig)
				continue
			}
			// poll the CRConfig so that it is synchronized with the TMConfig
			if _, err := p.Session.CRConfigRaw(p.OpsConfig.CdnName); err != nil {
				log.Errorf("MonitorConfigPoller: error getting CRConfig: %v", err)
				continue
			}
			p.writeConfig(MonitorCfg{CDN: p.OpsConfig.CdnName, Cfg: *monitorConfig})
		}
	}
}
