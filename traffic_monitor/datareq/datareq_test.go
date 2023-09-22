package datareq

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
	"math"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/config"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/test"

	jsoniter "github.com/json-iterator/go"
)

type AvailabilityType string

const (
	Random      AvailabilityType = "Random"
	Available                    = "Available"
	Unavailable                  = "Unavailable"
)

func getMockStaticAppData() config.StaticAppData {
	return config.StaticAppData{
		StartTime:      time.Now(),
		GitRevision:    "1234abc",
		FreeMemoryMB:   99999999,
		Version:        "0.1",
		WorkingDir:     "/usr/sbin/",
		Name:           "traffic_monitor",
		BuildTimestamp: time.Now().Format(time.RFC3339),
		Hostname:       "monitor01",
		UserAgent:      "traffic_monitor/0.1",
	}
}

func getMockLastHealthTimes() map[tc.CacheName]time.Duration {
	mockTimes := map[tc.CacheName]time.Duration{}
	numCaches := 10
	for i := 0; i < numCaches; i++ {
		mockTimes[tc.CacheName(test.RandStr())] = time.Duration(test.RandInt())
	}
	return mockTimes
}

func getMockCRStatesDeliveryService() tc.CRStatesDeliveryService {
	return tc.CRStatesDeliveryService{
		DisabledLocations: []tc.CacheGroupName{},
		IsAvailable:       test.RandBool(),
	}
}

func getMockPeerStates() peer.CRStatesThreadsafe {
	ps := peer.NewCRStatesThreadsafe()

	numCaches := 10
	for i := 0; i < numCaches; i++ {
		ps.SetCache(tc.CacheName(test.RandStr()), tc.IsAvailable{IsAvailable: test.RandBool()})
	}

	numDSes := 10
	for i := 0; i < numDSes; i++ {
		ps.SetDeliveryService(tc.DeliveryServiceName(test.RandStr()), getMockCRStatesDeliveryService())
	}
	return ps
}

func getRandDuration() time.Duration {
	return time.Duration(test.RandInt64())
}

func getResult(name tc.TrafficMonitorName, availabilityType AvailabilityType) peer.Result {
	peerStates := getMockPeerStates()

	availability := true

	if availabilityType == Random {
		availability = test.RandBool()
	} else if availabilityType == Unavailable {
		availability = false
	}

	return peer.Result{
		ID:         name,
		Available:  availability,
		Errors:     []error{errors.New(test.RandStr())},
		PeerStates: peerStates.Get(),
		PollID:     test.RandUint64(),
		// PollFinished chan<- uint64,
		Time: time.Now(),
	}
}

func getMockCRStatesPeers(quorumMin int, numPeers int, availabilityType AvailabilityType) peer.CRStatesPeersThreadsafe {
	ps := peer.NewCRStatesPeersThreadsafe(quorumMin)

	ps.SetTimeout(getRandDuration())

	randPeers := map[tc.TrafficMonitorName]struct{}{}
	for i := 0; i < numPeers; i++ {
		randPeers[tc.TrafficMonitorName(test.RandStr())] = struct{}{}
	}

	for peer, _ := range randPeers {
		ps.Set(getResult(peer, availabilityType))
	}

	ps.SetPeers(randPeers)

	return ps
}

func TestOptimisticQuorum(t *testing.T) {
	quorumMin := 1 // start with quorum enabled
	numPeers := 10

	// happy path; all peers available, quorum enabled
	peerStates := getMockCRStatesPeers(quorumMin, numPeers, Available)

	if !peerStates.OptimisticQuorumEnabled() {
		t.Fatalf("Optimistic quorum not enabled but should be; peers=%d, quorumMin=%d", quorumMin, numPeers)
	}

	optimisticQuorum, peersAvailable, peerCount, minimum := peerStates.HasOptimisticQuorum()

	if !optimisticQuorum {
		t.Fatalf("number of peers available (%d/%d) is less than the minimum number of %d required for optimistic peer quorum", peersAvailable, peerCount, minimum)
	}

	// no peers available
	peerStates = getMockCRStatesPeers(quorumMin, numPeers, Unavailable)

	optimisticQuorum, peersAvailable, peerCount, minimum = peerStates.HasOptimisticQuorum()

	if optimisticQuorum {
		t.Fatalf("optimistic quorum should be false; number of peers available (%d/%d) is less than the minimum number of %d required for optimistic peer quorum", peersAvailable, peerCount, minimum)
	}

	// optimistic quorum disabled
	quorumMin = 0
	peerStates = getMockCRStatesPeers(quorumMin, numPeers, Available)

	if peerStates.OptimisticQuorumEnabled() {
		t.Fatalf("Optimistic quorum enabled and should not be; peers=%d, quorumMin=%d", quorumMin, numPeers)
	}

	optimisticQuorum, peersAvailable, peerCount, minimum = peerStates.HasOptimisticQuorum()

	if !optimisticQuorum {
		t.Fatalf("optimistic quorum should be false; number of peers available (%d/%d) is less than the minimum number of %d required for optimistic peer quorum", peersAvailable, peerCount, minimum)
	}

	// optimistic quorum enabled but with a minimum greater than the number of peers; this config leads to 503s for any request
	quorumMin = 10
	numPeers = 4
	peerStates = getMockCRStatesPeers(quorumMin, numPeers, Available)

	if !peerStates.OptimisticQuorumEnabled() {
		t.Fatalf("Optimistic quorum not enabled but should be; peers=%d, quorumMin=%d", quorumMin, numPeers)
	}

	optimisticQuorum, peersAvailable, peerCount, minimum = peerStates.HasOptimisticQuorum()

	if optimisticQuorum {
		t.Fatalf("optimistic quorum should be false; number of peers available (%d/%d) is less than the minimum number of %d required for optimistic peer quorum", peersAvailable, peerCount, minimum)
	}
}

func TestGetStats(t *testing.T) {
	appData := getMockStaticAppData()
	pollingInterval := 5 * time.Second
	lastHealthTimes := getMockLastHealthTimes()
	fetchCount := uint64(test.RandInt())
	healthIteration := uint64(test.RandInt())
	errCount := uint64(test.RandInt())
	crStatesPeers := getMockCRStatesPeers(1, 10, Random)

	statsBts, err := getStats(appData, pollingInterval, lastHealthTimes, fetchCount, healthIteration, errCount, crStatesPeers)
	if err != nil {
		t.Fatalf("expected getStats error: nil, actual: %+v\n", err)
	}

	jsonStats := JSONStats{}
	json := jsoniter.ConfigFastest // TODO make configurable
	if err := json.Unmarshal(statsBts, &jsonStats); err != nil {
		t.Fatalf("expected getStats bytes: Stats JSON, actual: error decoding: %+v\n", err)
	}
	st := jsonStats.Stats

	if st.GitRevision != appData.GitRevision {
		t.Fatalf("expected getStats GitRevision '%+v', actual: '%+v'\n", appData.GitRevision, st.GitRevision)
	}
	if st.ErrorCount != errCount {
		t.Fatalf("expected getStats ErrorCount '%+v', actual: '%+v'\n", errCount, st.ErrorCount)
	}
	if st.Uptime < uint64(time.Since(appData.StartTime)/time.Second) {
		t.Fatalf("expected getStats Uptime > '%+v', actual: '%+v'\n", appData.StartTime, st.Uptime)
	}
	if tooFar := time.Since(appData.StartTime); st.Uptime > uint64(tooFar/time.Second-10) {
		t.Fatalf("expected getStats Uptime < '%+v', actual: '%+v'\n", tooFar, st.Uptime)
	}
	if st.FreeMemoryMB != appData.FreeMemoryMB {
		t.Fatalf("expected getStats FreeMemoryMB '%+v', actual: '%+v'\n", appData.FreeMemoryMB, st.FreeMemoryMB)
	}
	if st.Version != appData.Version {
		t.Fatalf("expected getStats Version '%+v', actual: '%+v'\n", appData.Version, st.Version)
	}
	if st.DeployDir != appData.WorkingDir {
		t.Fatalf("expected getStats DeployDir '%+v', actual: '%+v'\n", appData.WorkingDir, st.DeployDir)
	}
	if st.FetchCount != fetchCount {
		t.Fatalf("expected getStats FetchCount '%+v', actual: '%+v'\n", fetchCount, st.FetchCount)
	}

	slowestCache, slowestCacheTime := getLongestPoll(lastHealthTimes)

	if st.SlowestCache != string(slowestCache) {
		t.Fatalf("expected getStats SlowestCache '%+v', actual: '%+v'\n", slowestCache, st.SlowestCache)
	}
	if st.Name != appData.Name {
		t.Fatalf("expected getStats Name '%+v', actual: '%+v'\n", appData.Name, st.Name)
	}
	if st.BuildTimestamp != appData.BuildTimestamp {
		t.Fatalf("expected getStats BuildTimestamp '%+v', actual: '%+v'\n", appData.BuildTimestamp, st.BuildTimestamp)
	}
	if st.QueryIntervalTarget != int(pollingInterval/time.Millisecond) {
		t.Fatalf("expected getStats QueryIntervalTarget '%+v', actual: '%+v'\n", pollingInterval/time.Millisecond, st.QueryIntervalTarget)
	}
	if st.QueryIntervalActual != int(slowestCacheTime/time.Millisecond) {
		t.Fatalf("expected getStats QueryIntervalActual '%+v', actual: '%+v'\n", slowestCacheTime/time.Millisecond, st.QueryIntervalActual)
	}
	if st.QueryIntervalDelta != int((slowestCacheTime-pollingInterval)/time.Millisecond) {
		t.Fatalf("expected getStats QueryIntervalActual '%+v', actual: '%+v'\n", (slowestCacheTime - pollingInterval), st.QueryIntervalDelta)
	}
	if st.LastQueryInterval != int(math.Max(float64(slowestCacheTime), float64(pollingInterval))/float64(time.Millisecond)) {
		t.Fatalf("expected getStats LastQueryInterval expected '%+v', actual: '%+v'\n", int(math.Max(float64(slowestCacheTime), float64(pollingInterval))), st.LastQueryInterval)
	}
	if st.Microthreads <= 0 {
		t.Fatalf("expected getStats Microthreads >0, actual: '%+v'\n", st.Microthreads)
	}
	if st.LastGC == "" {
		t.Fatalf("expected getStats LastGC nonempty, actual: '%+v'\n", st.LastGC)
	}
	if st.MemAllocBytes <= 0 {
		t.Fatalf("expected getStats MemAllocBytes >0, actual: '%+v'\n", st.MemAllocBytes)
	}
	if st.MemTotalBytes <= 0 {
		t.Fatalf("expected getStats MemTotalBytes >0, actual: '%+v'\n", st.MemTotalBytes)
	}
	if st.MemSysBytes <= 0 {
		t.Fatalf("expected getStats MemSysBytes >0, actual: '%+v'\n", st.MemSysBytes)
	}
	if st.GCCPUFraction == 0.0 {
		t.Fatalf("expected getStats GCCPUFraction != 0, actual: '%+v'\n", st.GCCPUFraction)
	}

	oldestPolledPeer, oldestPolledPeerTime := oldestPeerPollTime(crStatesPeers.GetQueryTimes(), crStatesPeers.GetPeersOnline())
	if st.OldestPolledPeer != string(oldestPolledPeer) {
		t.Fatalf("expected getStats OldestPolledPeer '%+v', actual: '%+v'\n", oldestPolledPeer, st.OldestPolledPeer)
	}

	oldestPolledPeerTimeMS := time.Now().Sub((oldestPolledPeerTime)).Nanoseconds() / util.MSPerNS
	if st.OldestPolledPeerMs > oldestPolledPeerTimeMS+10 || st.OldestPolledPeerMs < oldestPolledPeerTimeMS-10 {
		t.Fatalf("expected getStats OldestPolledPeerMs '%+v', actual: '%+v'\n", oldestPolledPeerTimeMS, st.OldestPolledPeerMs)
	}

	queryInterval95thPercentile := getCacheTimePercentile(lastHealthTimes, 0.95).Nanoseconds() / util.MSPerNS

	if st.QueryInterval95thPercentile != queryInterval95thPercentile {
		t.Fatalf("expected getStats QueryInterval95thPercentile '%+v', actual: '%+v'\n", queryInterval95thPercentile, st.QueryInterval95thPercentile)
	}
}
