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
)

func TestStatusTransition(t *testing.T) {
	bad := errors.New("")
	var validTests = [][]error{
		[]error{nil, nil, bad, bad, bad}, //statusDraft
		[]error{nil, nil, nil, nil, bad}, //statusSubmitted
		[]error{nil, nil, nil, bad, bad}, //statusRejected
		[]error{bad, bad, bad, nil, nil}, //statusPending
		[]error{bad, bad, bad, bad, nil}, //statusComplete
	}

	// test all proper transitions
	for i := range validTests {
		from := requestStatus(i)
		for j, exp := range validTests[i] {
			to := requestStatus(j)
			if exp == nil {
				continue
			}

			exp = errors.New("invalid transition from " + from.name() + " to " + to.name())
			if from.validTransition(to).Error() != exp.Error() {
				t.Errorf("%s -> %s : expected %v, got %v", from.name(), to.name(), exp, from.validTransition(to))
			}
		}
	}

	// test names including out of range
	i := -1
	if requestStatus(i).name() != "INVALID" {
		t.Errorf("%d should be INVALID", i)
	}

	i = len(statusNames) + 1
	if requestStatus(i).name() != "INVALID" {
		t.Errorf("%d should be INVALID", i)
	}

	if statusDraft.name() != "draft" {
		t.Errorf("%d should be draft", int(statusDraft))
	}

	if statusSubmitted.name() != "submitted" {
		t.Errorf("%d should be submitted", int(statusSubmitted))
	}

	if statusRejected.name() != "rejected" {
		t.Errorf("%d should be rejected", int(statusRejected))
	}

	if statusPending.name() != "pending" {
		t.Errorf("%d should be pending", int(statusPending))
	}

	if statusComplete.name() != "complete" {
		t.Errorf("%d should be complete", int(statusComplete))
	}
}
