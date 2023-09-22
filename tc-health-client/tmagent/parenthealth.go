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
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/tc-health-client/sar"
	"github.com/apache/trafficcontrol/v8/tc-health-client/util"
)

const ParentHealthVersion = "1.0"

type ParentHealth struct {
	// Version is the direct version of this parent health, which applies to
	// direct RecursiveParentHealth which are ParentHealth.
	//
	// This does not apply to RecursiveParentHealth which are ParentServiceHealth;
	// rather, their own Version applies.
	//
	// Note this is the version of the service sending the result, not the service doing the polling.
	Version string `json:"version"`

	// Since is the RFC3339Nano-formatted time that this health result was polled.
	// Note this is the time of the service doing the polling, not the service sending the result.
	Since string `json:"since"`

	ParentHealthPollResults map[string]ParentHealthPollResult `json:"parent_health_poll_results"`
}

// ParentHealthPollResult is the health data from polling a parent.
// This is currently just a boolean, but may include other data in the future.
type ParentHealthPollResult struct {
	Healthy bool `json:"healthy"`
	// UnhealthyReason is a human-readable string as to why the parent was unhealthy.
	// This will be empty if Healthy is true.
	UnhealthyReason string `json:"unhealthy_reason,omitempty"` // TODO rename to HealthChangeReason?
	// TODO add poll type (L4 vs L7)?
}

func NewParentHealth() *ParentHealth {
	return &ParentHealth{
		Version:                 ParentHealthVersion,
		ParentHealthPollResults: map[string]ParentHealthPollResult{},
	}
}

// Serialize serializes the ParentHealth for network or disk output.
func (ph *ParentHealth) Serialize(prettyPrint bool) ([]byte, error) {
	if prettyPrint {
		return json.MarshalIndent(ph, "", "  ")
	}
	return json.Marshal(ph)
}

// DeserializeParentHealth deserializes the ParentHealth from network or disk output,
// as previously serialized with Serialize.
func DeserializeParentHealth(bts []byte) (*ParentHealth, error) {
	ph := &ParentHealth{}
	err := json.Unmarshal(bts, ph)
	if err != nil {
		return nil, errors.New("unmarshalling json: " + err.Error())
	}

	// this can be smarter if and when we have a new version that can be converted from an old.
	if ph.Version != ParentHealthVersion {
		return nil, errors.New("incompatible version '" + ph.Version + "'")
	}

	return ph, nil
}

// NewParentHealthPtr is a convenience func for NewAtomicPtr(NewParentHealth()).
func NewParentHealthPtr() *util.AtomicPtr[ParentHealth] {
	return util.NewAtomicPtr(NewParentHealth())
}

type ParentHealthPollType string

const ParentHealthPollTypeL4 = ParentHealthPollType("l4")
const ParentHealthPollTypeL7 = ParentHealthPollType("l7")
const ParentHealthPollTypeInvalid = ParentHealthPollType("")

func (pt ParentHealthPollType) String() string { return string(pt) }

func ParentHealthPollTypeFromStr(st string) ParentHealthPollType {
	switch st {
	case string(ParentHealthPollTypeL4):
		return ParentHealthPollTypeL4
	case string(ParentHealthPollTypeL7):
		return ParentHealthPollTypeL7
	default:
		return ParentHealthPollTypeInvalid
	}
}

// StartParentHealthPoller polls parents for health every pollInterval.
// Returns a done channel, which will terminate the polling loop when written to.
func StartParentHealthPoller(pi *ParentInfo, numWorkers int, pollInterval time.Duration, pollType ParentHealthPollType, updateHealthSignal func()) (chan<- struct{}, error) {
	// validate pollType immediately so we can assume it's valid in the rest of the poller
	if ParentHealthPollTypeFromStr(pollType.String()) == ParentHealthPollTypeInvalid {
		return nil, errors.New("invalid poll type")
	}

	// TODO dynamically get pollInterval every loop
	//      Which will require atomic/safe config loading
	doneChan := make(chan struct{})
	go parentHealthPoll(pi, pollInterval, pollType, numWorkers, doneChan, updateHealthSignal)
	return doneChan, nil
}

// startWorker starts a worker goroutine, which reads from workCh and calls the function it passes.
// The goroutine stops when workCh is closed.
func startWorker(workCh <-chan func()) {
	go func() {
		for f := range workCh {
			f()
		}
	}()
}

// const numParentHealthPollWorkers = 10 // TODO make configurable?

func parentHealthPoll(pi *ParentInfo, pollInterval time.Duration, pollType ParentHealthPollType, numWorkers int, doneChan <-chan struct{}, updateHealthSignal func()) {
	for {
		select {
		case <-doneChan:
			return
		default:
			break
		}

		start := time.Now()
		workCh := make(chan func(), 10) // TODO make buffer configurable?
		for i := 0; i < numWorkers; i++ {
			startWorker(workCh)
		}
		doPollParentHealth(workCh, pi, pollInterval, pollType)
		close(workCh) // closing the channel will cause the worker goroutines to return and stop
		updateHealthSignal()
		log.Infof("poll-status poll=parent-health-"+pollType.String()+" ms=%v\n", int(time.Since(start)/time.Millisecond))
		time.Sleep(pollInterval)
	}
}

func doPollParentHealth(workCh chan func(), pi *ParentInfo, pollInterval time.Duration, pollType ParentHealthPollType) {
	parentFQDNs := pi.GetParents()
	newParentHealth := pollParents(workCh, parentFQDNs, pollType)
	parentHealthPtr := pi.ParentHealthL4
	if pollType == ParentHealthPollTypeL7 {
		parentHealthPtr = pi.ParentHealthL7
	}

	logParentHealthChanges(pi.ParentHealthLog, pollType, parentHealthPtr.Get(), newParentHealth)

	if pollType == ParentHealthPollTypeL4 {
		pi.ParentHealthL4.Set(newParentHealth)
	} else if pollType == ParentHealthPollTypeL7 {
		pi.ParentHealthL7.Set(newParentHealth)
	}
}

// logParentHealthChanges writes parent health changes to plog.
func logParentHealthChanges(plog io.Writer, pollType ParentHealthPollType, oldPH *ParentHealth, newPH *ParentHealth) {
	for parentFQDN, newHealth := range newPH.ParentHealthPollResults {
		oldHealth, ok := oldPH.ParentHealthPollResults[parentFQDN]
		// We could remove this if check, if the user needed to log every poll,
		// even if it was the same (for example, if their logging system needed that).
		// TODO add config option?
		// TODO add event log to config
		if !ok || oldHealth.Healthy != newHealth.Healthy {
			logStr := time.Now().UTC().Format(time.RFC3339) + ` ParentHealth type=` + pollType.String() + ` parent=` + parentFQDN + ` healthy=` + strconv.FormatBool(newHealth.Healthy) + ` reason='` + newHealth.UnhealthyReason + `'` + "\n"
			bytesWritten, err := io.WriteString(plog, logStr)
			if err != nil {
				log.Errorln("writing to parent health log: " + err.Error())
			} else if bytesWritten != len(logStr) {
				log.Errorf("writing to parent health log: no error but wrote %v/%v bytes\n", bytesWritten, len(logStr))
			}
		}
	}
}

type ParentHealthPollResultAndFQDN struct {
	ParentHealthPollResult
	ParentFQDN string
}

func pollParents(workCh chan func(), parentFQDNs []string, pollType ParentHealthPollType) *ParentHealth {
	if pollType == ParentHealthPollTypeL4 {
		return pollParentsL4(workCh, parentFQDNs, pollType)
	}
	if pollType == ParentHealthPollTypeL7 {
		return pollParentsL7(workCh, parentFQDNs, pollType)
	}
	// should never happen, the poll start function should validate the poll type
	log.Errorf("pollParent got unknown poll type '%v', defaulting to L7!\n", pollType)
	return pollParentsL7(workCh, parentFQDNs, pollType)
}

func pollParentsL7(workCh chan func(), parentFQDNs []string, pollType ParentHealthPollType) *ParentHealth {
	var timeout = 2 * time.Second // TODO make configurable.
	const parentPort = 80         // TODO make configurable

	// Note this could be made parallel with a pool of goroutine workers
	// if it mattered for performance
	results := map[string]ParentHealthPollResult{}

	resultCh := make(chan ParentHealthPollResultAndFQDN) // TODO buffer?

	// writing the work to workers needs to be in its own goroutine,
	// because the channel buffer isn't large enough to hold all parents,
	// and we don't know how many parents there are when workers are created and need passed the work chan.
	// So, we need this writing goroutine to happen while we concurrently read work as it finishes in this thread.
	go func() {
		for _, parentFQDN := range parentFQDNs {
			workCh <- func() {
				resultCh <- pollParentL7(parentFQDN, parentPort, timeout)
			}
		}
	}()

	// read work as it finishes
	for i := 0; i < len(parentFQDNs); i++ {
		resultFQDN := <-resultCh
		results[resultFQDN.ParentFQDN] = resultFQDN.ParentHealthPollResult
	}
	return &ParentHealth{ParentHealthPollResults: results}
}

func pollParentL7(parentFQDN string, parentPort int, timeout time.Duration) ParentHealthPollResultAndFQDN {
	// TODO add service port, instead of assuming 80/443
	// TODO reuse http.Client

	httpClient := http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			// TODO more granular timeout options?
			// ConnectTimeout:        ParentPollTimeout,
			// RequestTimeout:        ParentPollTimeout,
			ResponseHeaderTimeout: timeout,
		},
	}

	req, err := http.NewRequest(http.MethodGet, "http://"+parentFQDN+":"+strconv.Itoa(parentPort), nil)
	if err != nil {
		// TODO log?
		return ParentHealthPollResultAndFQDN{
			ParentHealthPollResult: ParentHealthPollResult{
				Healthy:         false,
				UnhealthyReason: "error making request: " + err.Error(),
			},
			ParentFQDN: parentFQDN,
		}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return ParentHealthPollResultAndFQDN{
			ParentHealthPollResult: ParentHealthPollResult{
				Healthy:         false,
				UnhealthyReason: "error requesting: " + err.Error(),
			},
			ParentFQDN: parentFQDN,
		}
	}
	// Immediately close the body, we don't care what was in it.
	// But we do care that we're able to successfully stream the entire HTTP response.
	if err := resp.Body.Close(); err != nil {
		return ParentHealthPollResultAndFQDN{
			ParentHealthPollResult: ParentHealthPollResult{
				Healthy:         false,
				UnhealthyReason: "error reading body: " + err.Error(),
			},
			ParentFQDN: parentFQDN,
		}
	}

	return ParentHealthPollResultAndFQDN{
		ParentHealthPollResult: ParentHealthPollResult{
			Healthy: true,
		},
		ParentFQDN: parentFQDN,
	}
}

func pollParentsL4(workCh chan func(), parentFQDNs []string, pollType ParentHealthPollType) *ParentHealth {
	timeout := 2 * time.Second // TODO make configurable.
	parentPort := 80           // TODO make configurable

	// Note this could be made parallel with a pool of goroutine workers
	// if it mattered for performance
	results := map[string]ParentHealthPollResult{}

	hosts := []sar.HostPort{}
	for _, parentFQDN := range parentFQDNs {
		hosts = append(hosts, sar.HostPort{Host: parentFQDN, Port: parentPort})
	}

	sarResults, err := sar.MultiSAR(log.LLog(), hosts, timeout)
	if err != nil {
		for _, parentFQDN := range parentFQDNs {
			results[parentFQDN] = ParentHealthPollResult{
				Healthy:         false,
				UnhealthyReason: "l4 SAR error: " + err.Error(),
			}
		}
		return &ParentHealth{
			ParentHealthPollResults: results,
		}
	}

	for _, rs := range sarResults {
		if rs.Err != nil {
			results[rs.Host] = ParentHealthPollResult{
				Healthy:         false,
				UnhealthyReason: rs.Err.Error(),
			}
			continue
		}
		results[rs.Host] = ParentHealthPollResult{
			Healthy: true,
		}
	}
	return &ParentHealth{
		ParentHealthPollResults: results,
	}
}
