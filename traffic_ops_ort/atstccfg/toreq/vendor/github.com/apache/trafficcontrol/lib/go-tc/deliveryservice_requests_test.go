package tc

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
	"bytes"
	"encoding/json"
	"errors"
	"testing"
)

func TestStatus(t *testing.T) {
	type tester struct {
		st   RequestStatus
		name string
	}
	tests := []tester{
		tester{RequestStatus("foo"), "invalid"},
		tester{RequestStatusDraft, "draft"},
		tester{RequestStatusSubmitted, "submitted"},
		tester{RequestStatusRejected, "rejected"},
		tester{RequestStatusPending, "pending"},
		tester{RequestStatusComplete, "complete"},
	}

	for _, tst := range tests {
		v, _ := RequestStatusFromString(tst.name)
		if tst.name != string(v) {
			t.Errorf("%v: expected %s, got %s", tst, tst.name, string(v))
		}
	}
}

func TestStatusTransition(t *testing.T) {
	bad := errors.New("bad error")
	var validTests = [][]error{
		// To:  Dra  Sub  Rej  Pen  Com   // From:
		[]error{nil, nil, bad, bad, bad}, // Draft
		[]error{nil, nil, nil, nil, nil}, // Submitted
		[]error{bad, bad, bad, bad, bad}, // Rejected
		[]error{bad, bad, bad, nil, nil}, // Pending
		[]error{bad, bad, bad, bad, bad}, // Complete
	}

	// test all transitions
	for i := range validTests {
		from := RequestStatuses[i]
		for j, exp := range validTests[i] {
			to := RequestStatuses[j]
			if exp != nil {
				if from == RequestStatusRejected || from == RequestStatusComplete {
					exp = errors.New(string(from) + " request cannot be changed")
				} else {
					exp = errors.New("invalid transition from " + string(from) + " to " + string(to))
				}
			}
			got := from.ValidTransition(to)
			if got == exp {
				continue
			}

			if got != nil && exp != nil && got.Error() == exp.Error() {
				continue
			}

			t.Errorf("%s -> %s : expected %++v, got %++v", string(from), string(to), exp, got)
		}
	}
}

func TestRequestStatusJSON(t *testing.T) {
	b, err := json.Marshal(RequestStatusDraft)
	if err != nil {
		t.Errorf("Error marshalling %v: %s", RequestStatusDraft, err.Error())
	}

	exp := []byte(`"draft"`)
	if !bytes.Equal(exp, b) {
		t.Errorf("expected %s, got %s", exp, string(b))
	}

	var r RequestStatus
	err = json.Unmarshal([]byte(b), &r)
	if err != nil {
		t.Errorf("Error unmarshalling %s: %v", b, err)
	}
	if r != RequestStatusDraft {
		t.Errorf("expected %v, got %v", RequestStatusDraft, r)
	}
}
