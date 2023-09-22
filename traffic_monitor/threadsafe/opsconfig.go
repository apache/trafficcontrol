package threadsafe

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
	"sync"

	"github.com/apache/trafficcontrol/v8/traffic_monitor/handler"
)

// OpsConfig provides safe access for multiple reader goroutines and a single writer to a stored OpsConfig object.
// This could be made lock-free, if the performance was necessary
type OpsConfig struct {
	opsConfig *handler.OpsConfig
	m         *sync.RWMutex
}

// NewOpsConfig returns a new single-writer-multiple-reader OpsConfig
func NewOpsConfig() OpsConfig {
	return OpsConfig{m: &sync.RWMutex{}, opsConfig: &handler.OpsConfig{}}
}

// Get gets the internal OpsConfig object. This MUST NOT be modified. If modification is necessary, copy the object.
func (o *OpsConfig) Get() handler.OpsConfig {
	o.m.RLock()
	defer o.m.RUnlock()
	return *o.opsConfig
}

// Set sets the internal OpsConfig object. This MUST NOT be called from multiple goroutines.
func (o *OpsConfig) Set(newOpsConfig handler.OpsConfig) {
	o.m.Lock()
	*o.opsConfig = newOpsConfig
	o.m.Unlock()
}
