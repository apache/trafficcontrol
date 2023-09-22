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
	"bytes"
	"errors"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/tc-health-client/config"
)

// ParentStatus contains the trafficserver 'HostStatus' fields that
// are necessary to interface with the trafficserver 'traffic_ctl' command.
type ParentStatus struct {
	Fqdn                 string
	ActiveReason         bool
	LocalReason          bool
	ManualReason         bool
	LastTmPoll           int64
	UnavailablePollCount int
	MarkUpPollCount      int
}

// used to get the overall parent availablity from the
// HostStatus markdown reasons.  all markdown reasons
// must be true for a parent to be considered available.
func (p ParentStatus) available(reasonCode string) bool {
	rc := false

	switch reasonCode {
	case "active":
		rc = p.ActiveReason
	case "local":
		rc = p.LocalReason
	case "manual":
		rc = p.ManualReason
	}
	return rc
}

// used to log that a parent's status is either UP or
// DOWN based upon the HostStatus reason codes.  to
// be considered UP, all reason codes must be 'true'.
func (p ParentStatus) Status() string {
	if !p.ActiveReason {
		return "DOWN"
	} else if !p.LocalReason {
		return "DOWN"
	} else if !p.ManualReason {
		return "DOWN"
	}
	return "UP"
}

type StatusReason int

// these are the HostStatus reason codes used within
// trafficserver.
const (
	ACTIVE StatusReason = iota
	LOCAL
	MANUAL
)

// used for logging a parent's HostStatus reason code
// setting.
func (s StatusReason) String() string {
	switch s {
	case ACTIVE:
		return "ACTIVE"
	case LOCAL:
		return "LOCAL"
	case MANUAL:
		return "MANUAL"
	}
	return "UNDEFINED"
}

// MarkdownServer contains symbols for manipulating a running MarkdownService
// started by StartMarkdownService.
type MarkdownServer struct {
	// Shutdown signals the MarkdownService to stop its goroutines and release all resources.
	//
	// It is not buffered. Writes will block until the markdown service is done with any
	// ongoing markdown operation. Once a write returns, the MarkdownService will begin
	// shutdown.
	// TODO: add signal to callers when shutdown has finished.
	Shutdown chan<- struct{}

	// UpdateHealth signals to the markdown service that health status has changed, and it should process markdowns again.
	UpdateHealth func()
}

func StartMarkdownService(pi *ParentInfo, pollInterval time.Duration) *MarkdownServer {
	doneCh := make(chan struct{})
	updateCh := make(chan struct{})
	updateHealthF := func() {
		select {
		case updateCh <- struct{}{}:
		default:
		}
	}

	go markdownServicePoll(pi, pollInterval, doneCh, updateCh)
	return &MarkdownServer{Shutdown: doneCh, UpdateHealth: updateHealthF}
}

func markdownServicePoll(pi *ParentInfo, pollInterval time.Duration, doneChan <-chan struct{}, updateChan <-chan struct{}) {
	for {
		select {
		case <-doneChan:
			return
		default:
			<-updateChan
			start := time.Now()
			doMarkdown(pi)
			log.Infof("poll-status poll=markdown ms=%v\n", int(time.Since(start)/time.Millisecond))
			time.Sleep(pollInterval)
		}
	}
}

// decideHealthy is the algorithm for deciding whether a cache is healthy.
// It takes the cache hostname, and all health results, and returns a single boolean decision.
func decideHealthy(
	markdownMethods map[config.HealthMethod]struct{},
	parentFQDN string,
	tmHealth *TrafficMonitorHealth,
	parentHealthL4 *ParentHealth,
	parentHealthL7 *ParentHealth,
	parentServiceHealth *ParentServiceHealth,
) bool {
	// TODO use hostname:port? Parents can be healthy on one port/service but not another.
	// The ATS markdown command can't mark down per-port as of this writing, but
	// we can at least calculate it, to log that data, and be able to easily
	// mark down host:port when ATS adds that.

	// note the CacheStatuses has v4 and v6 data,
	// but ATS also can't mark down hosts per IP version yet,
	// so we just use the generic/old non-IP-versioned IsAvailable field.

	// TODO make decision algorithm/pessimism/ratio/etc configurable

	tmHealthy := (*bool)(nil)
	l4Healthy := (*bool)(nil)
	l7Healthy := (*bool)(nil)
	recursiveHealthy := (*bool)(nil)

	l4Health, hasL4Health := parentHealthL4.ParentHealthPollResults[parentFQDN]
	l7Health, hasL7Health := parentHealthL7.ParentHealthPollResults[parentFQDN]

	if _, use := markdownMethods[config.HealthMethodTrafficMonitor]; use && tmHealth.CacheStatuses != nil {
		parentHostName := parseFqdn(parentFQDN)
		tmHealthy = util.BoolPtr(tmHealth.CacheStatuses[tc.CacheName(parentHostName)].IsAvailable)
	}
	if _, use := markdownMethods[config.HealthMethodParentL4]; use && parentHealthL4 != nil && hasL4Health && len(parentHealthL4.ParentHealthPollResults) > 0 {
		l4Healthy = util.BoolPtr(l4Health.Healthy)
	}
	if _, use := markdownMethods[config.HealthMethodParentL7]; use && parentHealthL7 != nil && hasL7Health && len(parentHealthL7.ParentHealthPollResults) > 0 {
		l7Healthy = util.BoolPtr(l7Health.Healthy)
	}
	if _, use := markdownMethods[config.HealthMethodParentService]; use && parentServiceHealth != nil && parentServiceHealth.ParentServiceHealthPollResults != nil {
		recursiveHealthy = util.BoolPtr(decideRecursiveHealthy(parentServiceHealth.ParentServiceHealthPollResults[parentFQDN].ParentServiceHealth))
	}

	printbp := func(bp *bool) string {
		if bp == nil {
			return "nil"
		}
		return strconv.FormatBool(*bp)
	}

	// TODO wait for all polls to have results, before marking down?
	//      Maybe it doesn't matter for pessimistic health?

	log.Infof("decide-healthy monitored_host=%v tm=%v l4=%v l7=%v svc=%v\n", parentFQDN, printbp(tmHealthy), printbp(l4Healthy), printbp(l7Healthy), printbp(recursiveHealthy))

	// nilOrTrue is a helper func that returns true if the given *bool is nil or true.
	// This is used for pessimistic health, but if the health doesn't exist for the parent
	// for a particular health type poller, we don't want to mark it unhealthy.
	nilOrTrue := func(bp *bool) bool {
		return bp == nil || *bp
	}

	// pessimistic: if any health mechanism is unhealthy, consider the host unhealthy
	return nilOrTrue(tmHealthy) && nilOrTrue(l4Healthy) && nilOrTrue(l7Healthy) && nilOrTrue(recursiveHealthy)
}

// decideRecursiveHealthy is the algorithm for deciding whether a parent is healthy,
// based on a heuristic of whether it can get to most parents
// It takes the cache hostname, and all health results, and returns a single boolean decision.
//
// Note this is a heuristic because ATS doesn't currently have the ability to mark down
// a parent hostname for a single remap.
// If and when ATS has the ability to mark down parents per-remap, we can mark down parents
// just for remaps whose origins aren't available on that parent.
func decideRecursiveHealthy(parentServiceHealth *ParentServiceHealth) bool {
	// TODO get direct vs final parent data from strategies.yaml and parent.config
	//      Because we really want to heuristically check only the final parents (i.e. origins)
	//      are available on our direct parent, not any other origins it has that aren't assigned
	//      to this cache.

	// TODO make configurable
	const ratioToConsiderHealthy = 0.5

	parentHealthy := map[string]bool{}

	// TODO this is currently pessimistic: if any health mechanism is unhealthy, consider unhealthy.
	//      make health decision configurable.

	if parentServiceHealth == nil || parentServiceHealth.ParentServiceHealthPollResults == nil {
		log.Warnln("decideRecursiveHealthy got nil parents, returning healthy")
		return true
	}

	for parentFQDN, recursiveHealth := range parentServiceHealth.ParentServiceHealthPollResults {
		if recursiveHealth.ParentHealthL4 != nil && !recursiveHealth.ParentHealthL4.Healthy {
			parentHealthy[parentFQDN] = false
			continue
		}
		if recursiveHealth.ParentHealthL7 != nil && !recursiveHealth.ParentHealthL7.Healthy {
			parentHealthy[parentFQDN] = false
			continue
		}
		if !decideRecursiveHealthy(recursiveHealth.ParentServiceHealth.ParentServiceHealthPollResults[parentFQDN].ParentServiceHealth) {
			parentHealthy[parentFQDN] = false
			continue
		}
		parentHealthy[parentFQDN] = true
	}

	numParents := 0
	numParentsHealthy := 0
	for _, healthy := range parentHealthy {
		numParents++
		if healthy {
			numParentsHealthy++
		}
	}
	healthyRatio := float64(numParentsHealthy) / float64(numParents)
	return healthyRatio >= ratioToConsiderHealthy
}

type parent struct {
	host string
	fqdn string
}

// getParentFQDNs returns the FQDN of all parents in all health types.
// This is necessary, because some health polls might not have parents, so we need to combine them all
func getParentFQDNs(pi *ParentInfo, tmh *TrafficMonitorHealth, l4h *ParentHealth, l7h *ParentHealth, sh *ParentServiceHealth) []string {

	// put FQDns in a set first, to remove duplicates
	fqdnSet := map[string]struct{}{}

	for cacheName, _ := range tmh.CacheStatuses {
		hostName := string(cacheName)
		fqdn, ok := pi.ParentHostFQDNs.Load(hostName)
		if !ok {
			// this should be normal. Many parents in TM won't be parents of this cache
			continue
		}
		fqdnSet[fqdn] = struct{}{}
	}

	for fqdn, _ := range l4h.ParentHealthPollResults {
		fqdnSet[fqdn] = struct{}{}
	}
	for fqdn, _ := range l7h.ParentHealthPollResults {
		fqdnSet[fqdn] = struct{}{}
	}
	for fqdn, _ := range sh.ParentServiceHealthPollResults {
		// we don't need recursive health, just the direct parent. We only need direct parents to mark down.
		fqdnSet[fqdn] = struct{}{}
	}

	fqdns := make([]string, 0, len(fqdnSet))
	for fqdn, _ := range fqdnSet {
		fqdns = append(fqdns, fqdn)
	}
	return fqdns
}

// HealthSafetyRatio is the ratio of unhealthy parents after which no parents will be marked down.
//
// This is a safety mechanism: if for any reason most or all parents are marked down, something
// is seriously wrong, possibly with the health code itself, and therefore don't mark any parents down,
const HealthSafetyRatio = 0.3 // TODO make configurable?

func doMarkdown(pi *ParentInfo) {
	cfg := pi.Cfg.Get()
	var pv ParentStatus

	tmHealth := pi.TrafficMonitorHealth.Get()
	parentHealthL4 := pi.ParentHealthL4.Get()
	parentHealthL7 := pi.ParentHealthL7.Get()
	parentServiceHealth := pi.ParentServiceHealth.Get()

	parentFQDNs := getParentFQDNs(pi, tmHealth, parentHealthL4, parentHealthL7, parentServiceHealth)
	unhealthyNum := 0
	newCacheHealth := map[string]bool{} // map[fqdn]healthy
	for _, fqdn := range parentFQDNs {
		newAvailable := decideHealthy(pi.MarkdownMethods, fqdn, tmHealth, parentHealthL4, parentHealthL7, parentServiceHealth)
		if !newAvailable {
			unhealthyNum++
		}
		newCacheHealth[fqdn] = newAvailable
	}

	unhealthyParentsExceedSafetyRatio := (float64(unhealthyNum) / float64(len(newCacheHealth))) > HealthSafetyRatio
	log.Infof("markdown unhealthy_parents=%v total_parents=%v safety_ratio=%v\n", unhealthyNum, len(newCacheHealth), HealthSafetyRatio)
	if unhealthyParentsExceedSafetyRatio {
		log.Errorf("unhealthy_parents=%v total_parents=%v exceed safety_ratio=%v!! Marking all parents up!!\n", unhealthyNum, len(newCacheHealth), HealthSafetyRatio)
	}

	for _, fqdn := range parentFQDNs {
		isAvailable := tc.IsAvailable{}
		{
			hostName := parseFqdn(fqdn)
			isAvailable = tmHealth.CacheStatuses[tc.CacheName(hostName)]
		}
		parentStatus, ok := pi.LoadParentStatus(fqdn)
		if !ok {
			continue // TODO warn? error? is this normal?
		}

		// update the polling time
		parentStatus.LastTmPoll = tmHealth.PollTime.Unix()

		newAvailable := newCacheHealth[fqdn]
		if unhealthyParentsExceedSafetyRatio {
			newAvailable = true
		}
		// newAvailable := isAvailable.IsAvailable
		oldAvailable := parentStatus.available(cfg.ReasonCode)
		if oldAvailable != newAvailable {
			// do not mark down if the configuration disables mark downs.
			if !cfg.EnableActiveMarkdowns && !newAvailable {
				log.Infof("markdown monitored_host=%v host_status=%v event=\"TM reports host is not available\"", fqdn, pv.Status())
			} else {
				if newParentStatus, err := markParent(cfg, parentStatus, isAvailable.Status, newAvailable); err != nil {
					log.Errorln(err.Error())
				} else {
					log.Infoln("TM reports that '" + parentStatus.Fqdn + "' is not available so marked DOWN")
					parentStatus = newParentStatus
				}
			}
		}

		// if the host is available clear the unavailable poll count if not 0.
		if oldAvailable && newAvailable {
			if parentStatus.UnavailablePollCount > 0 {
				log.Debugf("resetting the UnavailablePollCount for '%s' from %d to 0",
					fqdn, parentStatus.UnavailablePollCount)
				parentStatus.UnavailablePollCount = 0
			}
		}

		{
			parentHostName := parseFqdn(parentStatus.Fqdn)
			pi.ParentHostFQDNs.Store(parentHostName, parentStatus.Fqdn)
		}

		// even if the status wasn't updated, we always need to update the poll time on the ParentStatus
		pi.StoreParentStatus(parentStatus.Fqdn, parentStatus)
	}

}

// markParent is used to mark a parent as up or down in the trafficserver HostStatus subsystem.
// It takes a parent, modifies it by marking available and setting reasons and other data,
// and returns the modified/marked parent.
func markParent(cfg *config.Cfg, pv ParentStatus, cacheStatus string, available bool) (ParentStatus, error) {
	var hostAvailable bool

	hostStatus := pv.Status()
	hostName := pv.Fqdn
	activeReason := pv.ActiveReason
	localReason := pv.LocalReason
	unavailablePollCount := pv.UnavailablePollCount
	markUpPollCount := pv.MarkUpPollCount

	log.Debugf("hostName: %s, UnavailablePollCount: %d, available: %v", hostName, unavailablePollCount, available)

	if !available { // unavailable
		unavailablePollCount += 1
		if unavailablePollCount < cfg.UnavailablePollThreshold {
			log.Infof("markdown monitored_host=%v host_status=UNAVAILABLE event=\"TM indicates host is unavailable but the UnavailablePollThreshold has not been reached\"", hostName)
			hostAvailable = true
		} else {
			// marking the host down
			if err := execTrafficCtl(pv.Fqdn, false, cfg.ReasonCode, cfg.TrafficServerBinDir); err != nil {
				log.Errorln(err)
				return ParentStatus{}, err
			}

			hostAvailable = false
			// reset the poll counts
			markUpPollCount = 0
			unavailablePollCount = 0
			log.Infof("marked monitored_host=%v host_status=%v event=\"%v\"\n", hostName, hostStatus, cacheStatus)
		}
	} else { // available
		// marking the host up
		markUpPollCount += 1
		if markUpPollCount < cfg.MarkUpPollThreshold {
			log.Infof("monitored_host=%v event=\"TM indicates host is available but the MarkUpPollThreshold has not been reached\"", hostName)
			hostAvailable = false
		} else {
			if err := execTrafficCtl(pv.Fqdn, true, cfg.ReasonCode, cfg.TrafficServerBinDir); err != nil {
				log.Errorln(err)
				return ParentStatus{}, err
			}
			hostAvailable = true
			// reset the poll counts
			unavailablePollCount = 0
			markUpPollCount = 0
			log.Infof("markdown monitored_host=%v host_status=%v event=\"%v\"\n", hostName, hostStatus, cacheStatus)
		}
	}

	// update parent info
	reason := cfg.ReasonCode
	switch reason {
	case "active":
		activeReason = hostAvailable
	case "local":
		localReason = hostAvailable
	}
	// save updates
	pv.ActiveReason = activeReason
	pv.LocalReason = localReason
	pv.UnavailablePollCount = unavailablePollCount
	pv.MarkUpPollCount = markUpPollCount
	log.Debugf("Updated parent status: %v", pv)
	return pv, nil
}

func execTrafficCtl(fqdn string, available bool, reason string, atsBinDir string) error {
	tc := filepath.Join(atsBinDir, TrafficCtl)

	var status string
	if available {
		status = "up"
	} else {
		status = "down"
	}

	cmd := exec.Command(tc, "host", status, "--reason", reason, fqdn)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return errors.New("marking " + fqdn + " " + status + ": " + TrafficCtl + " error: " + err.Error())
	}

	return nil
}
