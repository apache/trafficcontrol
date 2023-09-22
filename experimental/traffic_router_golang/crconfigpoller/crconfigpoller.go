package crconfigpoller

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
 *
 */

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/cgsrch"
	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/crconfig"
	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/crconfigregex"
	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/fetch"
	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/nextcache"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

func valid(new *tc.CRConfig, old *tc.CRConfig) error {
	if old == nil {
		return nil
	}
	if new == nil {
		return errors.New("CRConfig is nil")
	}

	if old.Stats.CDNName != nil && new.Stats.CDNName == nil {
		return errors.New("CDN name of new CRConfig is null")
	}
	if old.Stats.CDNName != nil && *new.Stats.CDNName != *old.Stats.CDNName {
		return errors.New("CDN name of new CRConfig '" + *new.Stats.CDNName + "' doesn't match old '" + *old.Stats.CDNName + "'")
	}

	if old.Stats.DateUnixSeconds != nil && new.Stats.DateUnixSeconds == nil {
		return errors.New("Date of new CRConfig is null")
	}
	if old.Stats.DateUnixSeconds != nil && *new.Stats.DateUnixSeconds < *old.Stats.DateUnixSeconds {
		return errors.New("Date of new CRConfig " + strconv.FormatInt(*new.Stats.DateUnixSeconds, 10) + " less than old " + strconv.FormatInt(*old.Stats.DateUnixSeconds, 10))
	}

	return nil
}

// createCGSearcher creates and returns the Cache Group Searcher, which can be used to efficiently find the nearest cachegroup to a given point.
func createCGSearcher(crc *tc.CRConfig) (cgsrch.Ths, error) {
	ths := cgsrch.NewThs()
	cgSearcher, err := cgsrch.Create(crc)
	if err != nil {
		return ths, errors.New("creating searcher: " + err.Error())
	}
	ths.Set(cgSearcher)
	return ths, nil
}

// createNextCacher creates and returns a NextCacher, which can be used to get the next cache to use for each Delivery Service (e.g. via round-robin, consistent hash, etc).
func createNextCacher(crc *tc.CRConfig) nextcache.NextCacher {
	dses := make([]tc.DeliveryServiceName, 0, len(crc.DeliveryServices))
	for ds, _ := range crc.DeliveryServices {
		dses = append(dses, tc.DeliveryServiceName(ds))
	}
	return nextcache.New(dses)
}

// TODO implement HTTP poller
func Start(fetcher fetch.Fetcher, interval time.Duration) (crconfig.Ths, crconfigregex.Ths, cgsrch.Ths, nextcache.Ths, error) {
	thsCrcRgx := crconfigregex.NewThs()
	thsCrc := crconfig.NewThs()
	thsCGSearcher := cgsrch.NewThs()
	thsNextCacher := nextcache.NewThs()
	prevBts := []byte{}
	prevCrc := (*tc.CRConfig)(nil)

	get := func() {
		newBts, err := fetcher.Fetch()
		if err != nil {
			fmt.Println("ERROR CRConfig read error: " + err.Error())
			return
		}

		if bytes.Equal(newBts, prevBts) {
			fmt.Println("INFO CRConfig unchanged.")
			return
		}

		fmt.Println("INFO CRConfig changed.")
		crc := &tc.CRConfig{}
		if err := json.Unmarshal(newBts, crc); err != nil {
			fmt.Println("ERROR CRConfig unmarshalling: " + err.Error())
			return
		}

		if err := valid(crc, prevCrc); err != nil {
			fmt.Println("ERROR not using invalid new CRConfig: " + err.Error())
			return
		}

		crcRgx, err := crconfigregex.Get(crc)
		if err != nil {
			fmt.Println("ERROR not using invalid new CRConfig: failed to get Regexes " + err.Error())
			return
		}
		cgSearcher, err := cgsrch.Create(crc)
		if err != nil {
			fmt.Println("ERROR not using invalid new CRConfig: failed to create Cachegroup searcher: " + err.Error())
		}
		nextCacher := createNextCacher(crc)

		thsNextCacher.Set(nextCacher)
		thsCGSearcher.Set(cgSearcher)
		thsCrc.Set(crc)
		thsCrcRgx.Set(&crcRgx)
		prevBts = newBts
		prevCrc = crc
		fmt.Println("INFO CRConfig set new")
	}

	get()

	go func() {
		for {
			time.Sleep(interval)
			get()
		}
	}()
	return thsCrc, thsCrcRgx, thsCGSearcher, thsNextCacher, nil
}
