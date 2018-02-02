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

func TestStatusTransition(t *testing.T) {
	bad := errors.New("")
	var validTests = [][]error{
		//      Dra  Sub  Rej  Pen  Com
		[]error{nil, nil, bad, bad, bad}, // Draft
		[]error{nil, nil, nil, nil, bad}, // Submitted
		[]error{nil, nil, nil, bad, bad}, // Rejected
		[]error{bad, bad, bad, nil, nil}, // Pending
		[]error{bad, bad, bad, bad, nil}, // Complete
	}

	// test all transitions
	for i := range validTests {
		from := RequestStatus(i)
		for j, exp := range validTests[i] {
			to := RequestStatus(j)
			if exp == nil {
				continue
			}

			exp = errors.New("invalid transition from " + from.Name() + " to " + to.Name())
			if from.ValidTransition(to).Error() != exp.Error() {
				t.Errorf("%s -> %s : expected %v, got %v", from.Name(), to.Name(), exp, from.ValidTransition(to))
			}
		}
	}

	type tester struct {
		st   RequestStatus
		name string
	}
	tests := []tester{
		tester{RequestStatus(-1), "INVALID"},
		tester{RequestStatus(9999), "INVALID"},
		tester{RequestStatusDraft, "draft"},
		tester{RequestStatusSubmitted, "submitted"},
		tester{RequestStatusRejected, "rejected"},
		tester{RequestStatusPending, "pending"},
		tester{RequestStatusComplete, "complete"},
	}

	for _, tst := range tests {
		got := tst.st.Name()
		exp := tst.name
		if got != exp {
			t.Errorf("%v: expected %s, got %s", tst.st, exp, got)
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
	err = json.Unmarshal(b, &r)
	if err != nil {
		t.Errorf("Error unmarshalling %s: %s", b, err.Error())
	}
	if r != RequestStatusDraft {
		t.Errorf("expected %v, got %v", RequestStatusDraft, r)
	}

}
