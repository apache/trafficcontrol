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
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/tc-health-client/config"
	"github.com/apache/trafficcontrol/v8/tc-health-client/util"
)

// ParentServiceHealth is the recursive parent health polled from other
// Parent Health Services.
//
// Parents will be recursively polled.
// If a parent is a cache, its Parent Health Service will be polled and the result
// placed in its map key.
// If a parent is not a cache, its ordinary Parent Health will be placed in its map key.
// See StartParentHealthPoller.
type ParentServiceHealth struct {
	// Version is the direct version of this parent health, which applies to
	// direct RecursiveParentHealth which are ParentHealth.
	//
	// This does not apply to RecursiveParentHealth which are ParentServiceHealth;
	// rather, their own Version applies.
	Version string `json:"version"`
	// ParentServiceHealthPollResults map[HostName]RecursiveParentHealth `json:"parent_health_poll_results"`

	ParentServiceHealthPollResults map[string]RecursiveParentHealth `json:"parent_service_health_poll_results"`

	// Since is the RFC3339Nano-formatted time that this health result was polled.
	// Note this is the time of the service doing the polling, not the service sending the result.
	Since time.Time `json:"since"`
}

// RecursiveParentHealth is either a poll from a direct parent poll of a non-cache
// or a ParentServiceHealth from a parent service poll of a cache.
//
// Either ParentHealth or ParentServiceHealth will not be nil, but never both.
//
// This would be more elegant as an interface, but a struct with two pointers is easier to deserialize.
type RecursiveParentHealth struct {
	ParentHealthL4      *ParentHealthPollResult `json:"parent_health_l4"`
	ParentHealthL7      *ParentHealthPollResult `json:"parent_health_l7"`
	ParentServiceHealth *ParentServiceHealth    `json:"parent_service_health"` // debug
}

func NewParentServiceHealth() *ParentServiceHealth {
	return &ParentServiceHealth{
		Version:                        ParentHealthVersion,
		ParentServiceHealthPollResults: map[string]RecursiveParentHealth{},
	}
}

// Serialize serializes the ParentServiceHealth for network or disk output.
func (ph *ParentServiceHealth) Serialize(prettyPrint bool) ([]byte, error) {
	if prettyPrint {
		return json.MarshalIndent(ph, "", "  ")
	}
	return json.Marshal(ph)
}

// DeserializeParentServiceHealth deserializes the ParentHealth from network or disk output,
// as previously serialized with Serialize.
func DeserializeParentServiceHealth(bts []byte) (*ParentServiceHealth, error) {
	ph := &ParentServiceHealth{}
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

// NewParentServiceHealthPtr is a convenience func for NewAtomicPtr(NewParentServiceHealth()).
func NewParentServiceHealthPtr() *util.AtomicPtr[ParentServiceHealth] {
	return util.NewAtomicPtr(NewParentServiceHealth())
}

// StartParentServiceHealthPoller polls parents for health every pollInterval.
// Note only parents which are caches are polled.
// Returns a done channel, which will terminate the polling loop when written to.
func StartParentServiceHealthPoller(pi *ParentInfo, numWorkers int, pollInterval time.Duration, updateHealthSignal func()) chan<- struct{} {
	// TODO share HTTP client, which should reuse connections

	// TODO dynamically get pollInterval every loop
	//      Which will require atomic/safe config loading
	doneChan := make(chan struct{})
	go parentServiceHealthPoll(pi, numWorkers, pollInterval, doneChan, updateHealthSignal)
	return doneChan
}

func parentServiceHealthPoll(pi *ParentInfo, numWorkers int, pollInterval time.Duration, doneChan <-chan struct{}, updateHealthSignal func()) {
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
		doPollParentServiceHealth(workCh, pi, pollInterval)
		close(workCh) // closing the channel will cause the worker goroutines to return and stop
		updateHealthSignal()
		log.Infof("poll-status poll=parent-health-service ms=%v\n", int(time.Since(start)/time.Millisecond))
		time.Sleep(pollInterval)
	}
}

func doPollParentServiceHealth(workCh chan func(), pi *ParentInfo, pollInterval time.Duration) {
	log.Debugf("doPollParentHealth starting\n")
	// TODO add config to not poll at all

	cfg := pi.Cfg.Get()
	parentFQDNs := pi.GetCacheParents()

	log.Debugf("about to poll service cache parents: %+v\n", parentFQDNs)

	newParentServiceHealth := pollParentServices(workCh, cfg, parentFQDNs)

	// logParentServiceHealthChanges(pi.ParentServiceHealthLog, pi.ParentserviceHealthPtr.Get(), newParentServiceHealth)

	pi.ParentServiceHealth.Set(newParentServiceHealth)
}

// logParentServiceHealthChanges writes parent health changes to the Event log.
// func logParentServiceHealthChanges(plog io.Writer, oldPH *ParentServiceHealth, newPH *ParentServiceHealth) {
// 	for parentFQDN, newHealth := range newPH.ParentServiceHealthPollResults {
// 		oldHealth, ok := oldPH.ParentServiceHealthPollResults[parentFQDN]
// 		// We could remove this if check, if the user needed to log every poll,
// 		// even if it was the same (for example, if their logging system needed that).
// 		// TODO add config option?
// 		// TODO add event log to config
// 		if !ok || oldHealth.Healthy != newHealth.Healthy {
// 			logStr := time.Now().UTC().Format(time.RFC3339) + ` ParentHealth '` + parentFQDN + `' healthy=` + strconv.FormatBool(newHealth.Healthy) + ` reason='` + newHealth.UnhealthyReason + `'` + "\n"
// 			bytesWritten, err := io.WriteString(plog, logStr)
// 			if err != nil {
// 				log.Errorln("writing to parent health log: " + err.Error())
// 			} else if bytesWritten != len(logStr) {
// 				log.Errorf("writing to parent health log: no error but wrote %v/%v bytes\n", bytesWritten, len(logStr))
// 			}
// 		}
// 	}

// 	// Healthy bool
// 	// UnhealthyReason string
// }

// ParentServicePollTimeoutService is the timeout for PollParent.
// TODO make configurable.
var ParentServicePollTimeoutService = 2 * time.Second

// ParentServiceHealthPollResult is the health data from polling a parent.
// This is currently just a boolean, but may include other data in the future.
// TODO delete
type ParentServiceHealthPollResult struct {
	Healthy bool
	// UnhealthyReason is a human-readable string as to why the parent was unhealthy.
	// This will be empty if Healthy is true.
	UnhealthyReason string // TODO rename to HealthChangeReason?
	// TODO add poll type (L4 vs L7)?
}

// ParentServiceHealthErrFQDN contains the health, any error, and the parent fqdn.
// This is used to return on a result chan from a poll worker.
type ParentServiceHealthErrFQDN struct {
	Health RecursiveParentHealth
	Err    error
	FQDN   string
}

func pollParentServices(workCh chan func(), cfg *config.Cfg, parentFQDNs []string) *ParentServiceHealth {
	// Note this could be made parallel with a pool of goroutine workers
	// if it mattered for performance
	results := map[string]RecursiveParentHealth{}

	resultCh := make(chan ParentServiceHealthErrFQDN) // TODO buffer?

	// write to work chan in a goroutine, because the work chan buffer won't be large enough that it'll block before we can start reading from it
	go func() {
		for _, parentFQDN := range parentFQDNs {
			if parentShortHostName := util.HostNameToShort(parentFQDN); parentShortHostName == cfg.HostName {
				// log.Debugf("pollParentServices got self %+v (%v), not polling self\n", parentFQDN, parentShortHostName)
				continue
			}

			workCh <- func() {
				const scheme = "http" // TODO always https and support client certs? Or make configurable?
				parentURI := scheme + `://` + parentFQDN + `:` + strconv.Itoa(cfg.ParentHealthServicePort)
				parentServiceHealth, err := pollParentService(parentURI, ParentServicePollTimeoutService)
				resultCh <- ParentServiceHealthErrFQDN{
					Health: RecursiveParentHealth{ParentServiceHealth: parentServiceHealth},
					Err:    err,
					FQDN:   parentFQDN,
				}
			}
		}
	}()

	for i := 0; i < len(parentFQDNs); i++ {
		result := <-resultCh
		if result.Err != nil {
			log.Errorln("parent service health poll for '" + result.FQDN + "' failed, parent health will be determined purely from native parent health; error: " + result.Err.Error())
		} else {
			results[result.FQDN] = result.Health
		}
	}

	return &ParentServiceHealth{
		ParentServiceHealthPollResults: results,
	}
}

func pollParentService(parentURI string, timeout time.Duration) (*ParentServiceHealth, error) {
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

	req, err := http.NewRequest(http.MethodGet, parentURI, nil)
	if err != nil {
		// TODO log?
		return nil, errors.New("making request: " + err.Error())
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.New("requesting: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.New("code " + strconv.Itoa(resp.StatusCode) + " reading: " + err.Error())
		}
		bts = bytes.ReplaceAll(bts, []byte("\n"), []byte(`\n`)) // we don't want newlines in the error log
		return nil, errors.New("code " + strconv.Itoa(resp.StatusCode) + " body: " + string(bts))
	}

	parentServiceHealth := &ParentServiceHealth{}
	if err := json.NewDecoder(resp.Body).Decode(parentServiceHealth); err != nil {
		return nil, errors.New("code " + strconv.Itoa(resp.StatusCode) + " decoding: " + err.Error())
	}

	if err := resp.Body.Close(); err != nil {
		return nil, errors.New("closing body: " + err.Error())
	}

	return parentServiceHealth, nil
}
