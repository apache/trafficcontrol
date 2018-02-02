package request

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
	"testing"

	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

func TestStatusTransition(t *testing.T) {
	bad := errors.New("")
	var validTests = [][]error{
		[]error{nil, nil, bad, bad, bad}, // Draft
		[]error{nil, nil, nil, nil, bad}, // Submitted
		[]error{nil, nil, nil, bad, bad}, // Rejected
		[]error{bad, bad, bad, nil, nil}, // Pending
		[]error{bad, bad, bad, bad, nil}, // Complete
	}

	// test all proper transitions
	for i := range validTests {
		from := tc.RequestStatus(i)
		for j, exp := range validTests[i] {
			to := tc.RequestStatus(j)
			if exp == nil {
				continue
			}

			exp = errors.New("invalid transition from " + from.Name() + " to " + to.Name())
			if from.ValidTransition(to).Error() != exp.Error() {
				t.Errorf("%s -> %s : expected %v, got %v", from.Name(), to.Name(), exp, from.ValidTransition(to))
			}
		}
	}

	// test names including out of range
	i := -1
	if tc.RequestStatus(i).Name() != "INVALID" {
		t.Errorf("%d should be INVALID", i)
	}

	i = len(tc.RequestStatusNames) + 1
	if tc.RequestStatus(i).Name() != "INVALID" {
		t.Errorf("%d should be INVALID", i)
	}

	if tc.RequestStatusDraft.Name() != "draft" {
		t.Errorf("%d should be draft", int(tc.RequestStatusDraft))
	}

	if tc.RequestStatusSubmitted.Name() != "submitted" {
		t.Errorf("%d should be submitted", int(tc.RequestStatusSubmitted))
	}

	if tc.RequestStatusRejected.Name() != "rejected" {
		t.Errorf("%d should be rejected", int(tc.RequestStatusRejected))
	}

	if tc.RequestStatusPending.Name() != "pending" {
		t.Errorf("%d should be pending", int(tc.RequestStatusPending))
	}

	if tc.RequestStatusComplete.Name() != "complete" {
		t.Errorf("%d should be complete", int(tc.RequestStatusComplete))
	}
}
