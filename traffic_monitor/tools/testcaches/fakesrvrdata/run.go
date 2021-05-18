package fakesrvrdata

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
	"math/rand"
	"time"
)

type BytesPerSec struct { // TODO change to PerMin? PerHour? (to allow, e.g. one 5xx per hour)
	Min FakeRemap
	Max FakeRemap
}

// runValidate verifies the FakeServerData Remaps match the RemapIncrements
func runValidate(s *FakeServerData, remapIncrements map[string]BytesPerSec) error {
	for r := range s.ATS.Remaps {
		if _, ok := remapIncrements[r]; !ok {
			return errors.New("remap increments missing server remap '" + r + "'")
		}
	}
	for r, bps := range remapIncrements {
		if _, ok := s.ATS.Remaps[r]; !ok {
			return errors.New("remap increments has remap not in server '" + r + "'")
		}
		if bps.Min.InBytes > bps.Max.InBytes || bps.Min.InBytes < 0 {
			return errors.New("invalid remap increments InBytes: must be Max >= Min >= 0)")
		}
		if bps.Min.OutBytes > bps.Max.OutBytes || bps.Min.OutBytes < 0 {
			return errors.New("invalid remap increments OutBytes: must be Max >= Min >= 0)")
		}
		if bps.Min.Status2xx > bps.Max.Status2xx || bps.Min.Status2xx < 0 {
			return errors.New("invalid remap increments Status2xx: must be Max >= Min >= 0)")
		}
		if bps.Min.Status3xx > bps.Max.Status3xx || bps.Min.Status3xx < 0 {
			return errors.New("invalid remap increments Status3xx: must be Max >= Min >= 0)")
		}
		if bps.Min.Status4xx > bps.Max.Status4xx || bps.Min.Status4xx < 0 {
			return errors.New("invalid remap increments Status4xx: must be Max >= Min >= 0)")
		}
		if bps.Min.Status5xx > bps.Max.Status5xx || bps.Min.Status5xx < 0 {
			return errors.New("invalid remap increments Status5xx: must be Max >= Min >= 0)")
		}
	}
	return nil
}

// Run takes a FakeServerData and a config, and starts running it, incrementing stats per the config. Returns a Threadsafe accessor to the running FakeServerData pointer, whose value may be accessed, but MUST NOT be modified.
// TODO add increments for Rcv,SndPackets, ProcLoadAvg variance, ConfigReloads
func Run(s FakeServerData, remapIncrements map[string]BytesPerSec) (Ths, error) {
	// TODO seed rand? Param?
	if err := runValidate(&s, remapIncrements); err != nil {
		return Ths{}, errors.New("invalid configuration: " + err.Error())
	}
	ths := NewThs()
	ths.Set(&s)

	go run(ths, remapIncrements)
	return ths, nil
}

type IncrementChanT struct {
	RemapName   string
	BytesPerSec BytesPerSec
}

// run starts a goroutine incrementing the FakeServerData's values according to the remapIncrements. Never returns.
func run(srvrThs Ths, remapIncrements map[string]BytesPerSec) {
	tickSecs := uint64(1) // adjustable for performance (i.e. a higher number is less CPU work)

	ticker := time.NewTicker(time.Second * time.Duration(tickSecs))

	for {
		select {
		case srvrThs.GetIncrementsChan <- remapIncrements:
		case newIncrement := <-srvrThs.IncrementChan:
			remapIncrements[newIncrement.RemapName] = newIncrement.BytesPerSec
		case <-ticker.C:
			srvr := srvrThs.Get()
			newRemaps := copyRemaps(srvr.ATS.Remaps)
			for remap, increments := range remapIncrements {
				srvrRemap := newRemaps[remap]

				addInBytes := increments.Min.InBytes * tickSecs
				if increments.Min.InBytes != increments.Max.InBytes {
					addInBytes += uint64(rand.Int63n(int64((increments.Max.InBytes - increments.Min.InBytes) * tickSecs)))
				}
				srvrRemap.InBytes += addInBytes
				srvr.System.ProcNetDev.RcvBytes += addInBytes

				addOutBytes := increments.Min.OutBytes * tickSecs
				if increments.Min.OutBytes != increments.Max.OutBytes {
					addOutBytes += uint64(rand.Int63n(int64((increments.Max.OutBytes - increments.Min.OutBytes) * tickSecs)))
				}
				srvrRemap.OutBytes += addOutBytes
				srvr.System.ProcNetDev.SndBytes += addOutBytes

				srvrRemap.Status2xx += increments.Min.Status2xx * tickSecs
				if increments.Min.Status2xx != increments.Max.Status2xx {
					srvrRemap.Status2xx += uint64(rand.Int63n(int64((increments.Max.Status2xx - increments.Min.Status2xx) * tickSecs)))
				}

				srvrRemap.Status3xx += increments.Min.Status3xx * tickSecs
				if increments.Min.Status3xx != increments.Max.Status3xx {
					srvrRemap.Status3xx += uint64(rand.Int63n(int64((increments.Max.Status3xx - increments.Min.Status3xx) * tickSecs)))
				}

				srvrRemap.Status4xx += increments.Min.Status4xx * tickSecs
				if increments.Min.Status4xx != increments.Max.Status4xx {
					srvrRemap.Status4xx += uint64(rand.Int63n(int64((increments.Max.Status4xx - increments.Min.Status4xx) * tickSecs)))
				}

				srvrRemap.Status5xx += increments.Min.Status5xx * tickSecs
				if increments.Min.Status5xx != increments.Max.Status5xx {
					srvrRemap.Status5xx += uint64(rand.Int63n(int64((increments.Max.Status5xx - increments.Min.Status5xx) * tickSecs)))
				}

				newRemaps[remap] = srvrRemap
			}
			srvr.ATS.Remaps = newRemaps
			srvrThs.Set(srvr)
		}
	}
}
