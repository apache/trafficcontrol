package tcdata

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
	"reflect"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/v7/lib/go-tc"
)

type topologyTestCase struct {
	testCaseDescription string
	tc.Topology
}

func (r *TCData) CreateTestTopologies(t *testing.T) {
	var (
		postResponse *tc.TopologyResponse
		err          error
	)
	for _, topology := range r.TestData.Topologies {
		if postResponse, _, err = TOSession.CreateTopology(topology); err != nil {
			t.Fatalf("could not CREATE topology: %v", err)
		}
		postResponse.Response.LastUpdated = nil
		if !reflect.DeepEqual(topology, postResponse.Response) {
			t.Fatalf("Topology in response should be the same as the one POSTed. expected: %v, actual: %v", topology, postResponse.Response)
		}
	}
}

func (r *TCData) DeleteTestTopologies(t *testing.T) {
	for _, top := range r.TestData.Topologies {
		delResp, _, err := TOSession.DeleteTopology(top.Name)
		if err != nil {
			t.Fatalf("cannot DELETE topology: %v - %v", err, delResp)
		}
		deleteLog, _, err := TOSession.GetLogsByLimit(1)
		if err != nil {
			t.Fatalf("unable to get latest audit log entry")
		}
		if len(deleteLog) != 1 {
			t.Fatalf("log entry length - expected: 1, actual: %d", len(deleteLog))
		}
		if !strings.Contains(*deleteLog[0].Message, top.Name) {
			t.Errorf("topology deletion audit log entry - expected: message containing topology name '%s', actual: %s", top.Name, *deleteLog[0].Message)
		}

		topology, _, err := TOSession.GetTopology(top.Name)
		if err == nil {
			t.Fatalf("expected error trying to GET deleted topology: %s, actual: nil", top.Name)
		}
		if topology != nil {
			t.Fatalf("expected nil trying to GET deleted topology: %s, actual: non-nil", top.Name)
		}
	}
}
