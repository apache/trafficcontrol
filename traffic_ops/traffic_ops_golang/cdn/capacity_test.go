package cdn

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
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

func TestGetStatsFromServiceInterface(t *testing.T) {
	data1 := tc.ServerStats{
		Interfaces: nil,
		Stats: map[string][]tc.ResultStatVal{
			"kbps": {
				{Val: 24.5},
			},
			"maxKbps": {
				{Val: 66.8},
			},
		},
	}

	kbps, maxKbps, err := getStats(data1)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err.Error())
	}
	if kbps != 24.5 || maxKbps != 66.8 {
		t.Errorf("Expected kbps to be 24.5, got %v; Expected maxKbps to be 66.8, got %v", kbps, maxKbps)
	}
}
