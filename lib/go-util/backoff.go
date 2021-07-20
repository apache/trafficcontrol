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

// Used to return an exponentially increasing time.Duration value that may be used
// to sleep between retries.

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// DefaultFactor may be used by applications for the factor argument to
// NewBackoff.
const DefaultFactor = 2.0

// ConstantBackoffDuration is a fallback duration that may be used by
// an application along with NewConstantBackoff().
const ConstantBackoffDuration = 30 * time.Second

// Implementation of a Backoff that returns a constant time duration
// used as a fallback.
type constantBackoff struct{ d time.Duration }

// NewConstantBackoff returns a Backoff that does not change its duration.
//
// This is roughly equivalent to calling NewBackoff(d, d+1*time.Second, 1.0)
// (the max duration doesn't matter when the factor is 1), but is more
// efficient in terms of both CPU load/time and memory used, and so should be
// done instead.
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

// A Backoff is a definition of how long to wait between attempting something,
// which is normally a network action, to avoid overloading requested
// resources.
type Backoff interface {
	// BackoffDuration returns the time that should be waited before attempting
	// the action again. Normally, this duration will grow exponentially, but
	// in the case of a Backoff constructed using NewConstantBackoff, it will
	// be the same every time.
	BackoffDuration() time.Duration
	// Reset clears any incrementing of the duration returned by
	// BackoffDuration, so that the next call will yield the same duration as
	// the first call.
	Reset()
}

// NewBackoff constructs and returns a Backoff that starts with a duration at
// min and increments it exponentially according to the "factor" up to a
// maximum defined by the passed max.
//
// The rate of increase is defined to be nmfⁿ⁻¹ where n is the number of
// attempts that have already been made, m is the minimum duration, and f is
// the factor. The duration for any given attempt n (starting at zero) is
// defined to be mfⁿ+j where j is a randomly generated "jitter" that is added
// to each duration. The "jitter" will be a number of nanoseconds between the
// zero and the magnitude of the difference between the min and the factor. If
// the factor, treated as a number of nanoseconds, is greater than the minimum
// duration, this jitter will *subtract* from the resulting duration, and it
// will be between this difference (exclusive) and zero (inclusive). If the
// factor is less than the min it will *add* to the resulting duration, and it
// will be between 0 (inclusive) and the difference (exclusive).
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

// generate random jitter.
func (b *backoff) jitter(durFloat float64, minFloat float64) float64 {
	return b.rgen.Float64() * (durFloat - minFloat)
}

func (b *backoff) Reset() {
	b.attempt = 0
}

// Calculate and return  backoff time duration.
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
