package tmagent

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
	"errors"
	"net/http"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/tc-health-client/config"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/tmclient"
)

type TrafficMonitorHealth struct {
	// PollTime is the local time this TM health was obtained
	PollTime      time.Time
	CacheStatuses map[tc.CacheName]tc.IsAvailable
}

func NewTrafficMonitorHealth() *TrafficMonitorHealth {
	return &TrafficMonitorHealth{
		CacheStatuses: map[tc.CacheName]tc.IsAvailable{},
	}
}

// StartTrafficMonitorHealthPoller
// The main polling function that keeps the parents list current if
// with any changes to the trafficserver 'parent.config' or 'strategies.yaml'.
// Also, it keeps parent status current with the the trafficserver HostStatus
// subsystem.  Finally, on each poll cycle a trafficmonitor is queried to check
// that all parents used by this trafficserver are available for use based upon
// the trafficmonitors idea from it's health protocol.  Parents are marked up or
// down in the trafficserver subsystem based upon that hosts current status and
// the status that trafficmonitor health protocol has determined for a parent.
func StartTrafficMonitorHealthPoller(pi *ParentInfo, updateHealthSignal func()) chan<- struct{} {
	log.Infoln("Traffic Monitor Health Poller started")
	doneChan := make(chan struct{})
	go loopPollAndUpdateCacheStatus(pi, doneChan, updateHealthSignal)
	return doneChan
}

func loopPollAndUpdateCacheStatus(pi *ParentInfo, doneChan <-chan struct{}, updateHealthSignal func()) {
	cfg := pi.Cfg.Get()
	toLoginDispersion := config.GetTOLoginDispersion(cfg.TMPollingInterval, cfg.TOLoginDispersionFactor)
	for {
		select {
		case <-doneChan:
			return
		default:
			break
		}

		pollingInterval := pi.Cfg.Get().TMPollingInterval
		if toLoginDispersion <= 0 {
			cfg := pi.Cfg.Get()
			toLoginDispersion = config.GetTOLoginDispersion(cfg.TMPollingInterval, cfg.TOLoginDispersionFactor)
		} else {
			toLoginDispersion -= pollingInterval
		}
		doTrafficOpsReq := toLoginDispersion <= 0

		log.Infoln("service-status service=tm-health event=\"starting\"")
		start := time.Now()
		doPollAndUpdateCacheStatus(pi, doTrafficOpsReq)
		updateHealthSignal()
		log.Infof("poll-status poll=tm-health ms=%v\n", int(time.Since(start)/time.Millisecond))
		time.Sleep(pollingInterval)
	}
}

func doPollAndUpdateCacheStatus(pi *ParentInfo, doTrafficOpsReq bool) {
	cfg := pi.Cfg.Get()

	// check for parent and strategies config file updates, and trafficserver
	// host status changes.  If an error is encountered reading data the current
	// parents lists and hoststatus remains unchanged.
	// TODO move reading ATS config files to its own poller service
	if err := pi.UpdateParentInfo(cfg); err != nil {
		log.Errorf("could not load new ATS parent info: %s\n", err.Error())
	} else {
		// log.Debugf("updated parent info, total number of parents: %d\n", len(pi.Parents))
		// TODO track map len
		log.Debugf("tm-agent total_parents=%v\n", len(pi.GetParents()))
	}

	// read traffic manager cache statuses.
	crStates, err := pi.GetCacheStatuses()
	now := time.Now() // get the current poll time
	if err != nil {
		log.Errorf("poll-status %v\n", err.Error())
		if err := pi.GetTOData(cfg); err != nil {
			log.Errorf("update event=\"could not update the list of trafficmonitors, keeping the old config\": %v", err.Error())
		} else {
			log.Infoln("service-status service=tm-health event=\"updated TrafficMonitor statuses from TrafficOps\"")
		}

		// log the poll state data if enabled
		if cfg.EnablePollStateLog {
			err = pi.WritePollState()
			if err != nil {
				log.Errorf("could not write the poll state log: %s\n", err.Error())
			}
		}
		return
	}

	pi.TrafficMonitorHealth.Set(&TrafficMonitorHealth{
		PollTime:      now,
		CacheStatuses: crStates.Caches,
	})

	// periodically update the TrafficMonitor list and statuses
	if doTrafficOpsReq {
		// TODO move to its own TO poller
		if err = pi.GetTOData(cfg); err != nil {
			log.Errorf("update event=\"could not update the list of trafficmonitors, keeping the old config\": %v", err.Error())
		} else {
			log.Infoln("service-status service=tm-health event=\"updated TrafficMonitor statuses from TrafficOps\"")
		}
	}

	// log the poll state data if enabled
	if cfg.EnablePollStateLog {
		if err = pi.WritePollState(); err != nil {
			log.Errorf("could not write the poll state log: %s\n", err.Error())
		}
	}
}

// Queries a traffic monitor that is monitoring the trafficserver instance running on a host to
// obtain the availability, health, of a parent used by trafficserver.
func (pi *ParentInfo) GetCacheStatuses() (tc.CRStates, error) {
	cfg := pi.Cfg.Get()

	tmHostName, err := pi.findATrafficMonitor()
	if err != nil {
		return tc.CRStates{}, errors.New("monitor=finding a trafficmonitor: " + err.Error())
	}
	tmc := tmclient.New("http://"+tmHostName, cfg.TORequestTimeout)

	// Use a proxy to query TM if the ProxyURL is set
	if cfg.ParsedProxyURL != nil {
		tmc.Transport = &http.Transport{Proxy: http.ProxyURL(cfg.ParsedProxyURL)}
	}

	return tmc.CRStates(false)
}
