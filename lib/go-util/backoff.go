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

// Used to return an exponentially increasing time.Duration value that may be used
// to sleep between retries.

package util

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// DefaultFactor may be used by applications for the factor argument.
const DefaultFactor = 2.0

// ConstantBackoffDuration is a fallback duration that may be used by
// an application along with NewConstantBackoff()
const ConstantBackoffDuration = 30 * time.Second

// Implementation of a Backoff that returns a constant time duration
// used as a fallback.
type constantBackoff struct{ d time.Duration }

func NewConstantBackoff(d time.Duration) Backoff {
	return &constantBackoff{d}
}
func (b *constantBackoff) BackoffDuration() time.Duration { return b.d }
func (b *constantBackoff) Reset()                         {}

type backoff struct {
	attempt float64
	Factor  float64
	Min     time.Duration
	Max     time.Duration
	rgen    *rand.Rand
}

type Backoff interface {
	BackoffDuration() time.Duration
	Reset()
}

func NewBackoff(min time.Duration, max time.Duration, factor float64) (Backoff, error) {

	// verify arguments and set defaults if necessary.
	if min < 1 {
		return nil, fmt.Errorf("'min: %v, is invalid.  min must be greater than '1'", min)
	}
	if max <= min {
		return nil, fmt.Errorf("'max: %v, is invalid.  max must be greater than 'min'", max)
	}
	if factor <= 1.0 {
		return nil, fmt.Errorf("'factor: %v, is invalid.  factor must be greater than '1'", factor)
	}

	src := rand.NewSource(time.Now().UTC().UnixNano())

	return &backoff{
		attempt: 0,
		Factor:  factor,
		Min:     min,
		Max:     max,
		rgen:    rand.New(src),
	}, nil
}

// generate random jitter
func (b *backoff) jitter(durFloat float64, minFloat float64) float64 {
	return b.rgen.Float64() * (durFloat - minFloat)
}

func (b *backoff) Reset() {
	b.attempt = 0
}

// Calculate and return  backoff time duration
func (b *backoff) BackoffDuration() time.Duration {

	minFloat := float64(b.Min)
	durFloat := minFloat * math.Pow(b.Factor, b.attempt)
	b.attempt++

	// add jitter
	durFloat += b.jitter(durFloat, minFloat)

	// reached the max duration return max
	if durFloat >= float64(b.Max) {
		return b.Max
	}

	dur := time.Duration(durFloat)

	return dur
}
