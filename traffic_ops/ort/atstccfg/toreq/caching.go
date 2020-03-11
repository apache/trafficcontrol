package toreq

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
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/config"
)

// GetRetry attempts to get the given object from tempDir/cacheFileName, retrying with exponential backoff up to cfg.NumRetries.
func GetRetry(cfg config.TCCfg, cacheFileName string, obj interface{}, getter func(obj interface{}) error) error {
	start := time.Now()
	currentRetry := 0
	for {
		err := getter(obj)
		if err == nil {
			break
		}
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			// if the server returned a 404, retrying won't help
			return errors.New("getting uncached: " + err.Error())
		}
		if currentRetry == cfg.NumRetries {
			return errors.New("getting uncached: " + err.Error())
		}

		sleepSeconds := config.RetryBackoffSeconds(currentRetry)
		log.Warnf("getting '%v', sleeping for %v seconds: %v\n", cacheFileName, sleepSeconds, err)
		currentRetry++
		time.Sleep(time.Second * time.Duration(sleepSeconds)) // TODO make backoff configurable?
	}

	log.Infof("GetCachedJSON %v (miss) took %v\n", cacheFileName, time.Since(start).Round(time.Millisecond))
	return nil
}
