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
	"time"

	"github.com/apache/trafficcontrol/v8/traffic_monitor/config"
)

const DefaultPollerType = PollerTypeHTTP

// Poller is a particular type of cache stat poller. Examples are HTTP, TCP, or NFS files. It only polls and returns bytes received, it does not do any parsing. For stat parsing, see traffic_monitor/cache/stats_types.go.
type PollerType struct {
	GlobalInit PollerGlobalInitFunc
	Init       PollerInitFunc
	Poll       PollerFunc
}

// PollerConfig is the data given to cache pollers when they're initialized.
type PollerConfig struct {
	Timeout     time.Duration
	NoKeepAlive bool
	PollerID    string
}

// PollerGlobalInit performs global initialization, and returns a global context object.
// For example, a poller might create a threadsafe client object, which can be used concurrently by all pollers.
type PollerGlobalInitFunc func(cfg config.Config, appData config.StaticAppData) interface{}

// PollerInit performs initialization for a specific poller. It takes the global context created by the poller's GlobalInit.
// For example, a poller might create a template client in its GlobalInit, and then use that to create a non-threadsafe client for each concurrent poller.
type PollerInitFunc func(cfg PollerConfig, globalCtx interface{}) interface{}

// pollers holds the functions for polling caches. This is not const, because Go doesn't allow constant maps. This is populated on startup, and MUST NOT be modified after startup.
var pollers = map[string]PollerType{}

// PollerFunc polls a cache. It takes the global context created by this Poller's GlobalInit, and the poller-specific context created by this poller's Init. It returns the response bytes, the time the request finished, the length of time the request took, and any error.
// If the PollerFunc needs the global context object, the Init func should embed it in the context object it returns. If Init is nil, the global context will be given to the poller.
type PollerFunc func(ctx interface{}, url string, host string, pollID uint64) ([]byte, time.Time, time.Duration, error)

// AddPollerType adds a poller with the given name, and the given init and poll funcs. The globalInit and init funcs may be nil; poller MUST NOT be nil.
func AddPollerType(name string, globalInit PollerGlobalInitFunc, init PollerInitFunc, poller PollerFunc) {
	pollers[name] = PollerType{GlobalInit: globalInit, Init: init, Poll: poller}
}

// GetGlobalContexts returns the global contexts corresponding to the registered pollers
func GetGlobalContexts(cfg config.Config, appData config.StaticAppData) map[string]interface{} {
	ctxs := map[string]interface{}{}
	for pollerName, poller := range pollers {
		if poller.GlobalInit != nil {
			ctxs[pollerName] = poller.GlobalInit(cfg, appData)
		}
	}
	return ctxs
}
