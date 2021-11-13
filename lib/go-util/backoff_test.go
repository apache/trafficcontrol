package util

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
	"time"
)

func TestNewBackoff(t *testing.T) {
	min := 500 * time.Millisecond
	max := 600 * time.Millisecond
	fac := DefaultFactor

	bo, err := NewBackoff(min, max, fac)

	if err != nil {
		t.Errorf("un-expected error. min: %v, max: %v, fac: %v - %v", min, max, fac, err)
	}

	if bo == nil {
		t.Errorf("un-expected error bo is nil, min: %v, max: %v, fac: %v - %v", min, max, fac, err)
	}

	bo, err = NewBackoff(0, 1, 1.1)

	if err == nil {
		t.Errorf("Expected errors passing invalid minimum to NewBackOff()")
	}

	if bo != nil {
		t.Errorf("Unexpected error passing invalid minimum to NewBackOff(), 'bo' should be nil")
	}

	_, err = NewBackoff(1, 1, 1.1)
	if err == nil {
		t.Error("Expected passing a maximum that's not greater than a minimum to return an error")
	}

	_, err = NewBackoff(1, 2, 1.0)
	if err == nil {
		t.Error("Expected passing an invalid factor to return an error")
	}
}

func TestReset(t *testing.T) {
	bo, err := NewBackoff(1, 2000000000, DefaultFactor)

	if err != nil {
		t.Errorf("un-expected error: %v", err)
	}

	val1 := bo.BackoffDuration()
	val2 := bo.BackoffDuration()

	if val2 < val1 {
		t.Errorf("expected val1: %v is less than val2: %v", val1, val2)
	}

	bo.Reset()

	val3 := bo.BackoffDuration()

	if val3 > val2 {
		t.Errorf("expected val3: %v is less than val2: %v", val3, val2)
	}
}

func TestBackoffDuration(t *testing.T) {
	min := 100 * time.Nanosecond
	max := 20000 * time.Nanosecond
	fac := DefaultFactor
	bo, err := NewBackoff(min, max, fac)

	if err != nil {
		t.Errorf("un-expected error: %v", err)
	}
	dur := min

	val := bo.BackoffDuration()
	if val < min {
		t.Errorf("unexpected duration calculation, val: %v is less than min: %v", val, min)
	}

	// 8 iterations with the default settings.
	for i := 0; i < 8; i++ {
		val = bo.BackoffDuration()
		if val < dur {
			t.Errorf("unexpected duration calculation, iteration: %v, val: %v  <= dur: %v", i, val, dur)
		}
		dur = val
	}

	// after 8 calls with the default settings, val should be '1m0s', DefaultMaxMS
	if val != max {
		t.Errorf("unexpected duration calculation, val: %v  != : max: %v", val, max)
	}

	fac = 1.000001
	bo, err = NewBackoff(min, max, fac)
	for i := 0; i < 8; i++ {
		val = bo.BackoffDuration()
		if val < min {
			t.Errorf("backoff duration generated as less than the minimum")
		}
	}

	fac = 100000
	bo, err = NewBackoff(min, max, fac)
	for i := 0; i < 100; i++ {
		val = bo.BackoffDuration()
		if val > max {
			t.Errorf("backoff duration generated as greater than the maximum")
		}
	}
}

func TestNewConstantBackoff(t *testing.T) {
	bo := NewConstantBackoff(ConstantBackoffDuration)

	if bo.BackoffDuration() != ConstantBackoffDuration {
		t.Errorf("unexepected duration, return value != %v", ConstantBackoffDuration)
	}
}
