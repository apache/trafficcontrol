package peer

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
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

func TestCrStates(t *testing.T) {
	t.Log("Running Peer Tests")

	text, err := ioutil.ReadFile("crstates.json")
	if err != nil {
		t.Log(err)
	}
	crStates, err := tc.CRStatesUnMarshall(text)
	if err != nil {
		t.Log(err)
	}
	fmt.Println(len(crStates.Caches), "caches found")
	for cacheName, crState := range crStates.Caches {
		t.Logf("%v -> %v", cacheName, crState.IsAvailable)
	}

	fmt.Println(len(crStates.DeliveryService), "deliveryservices found")
	for dsName, deliveryService := range crStates.DeliveryService {
		t.Logf("%v -> %v (len:%v)", dsName, deliveryService.IsAvailable, len(deliveryService.DisabledLocations))
	}

}
