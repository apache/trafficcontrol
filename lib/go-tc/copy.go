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

func coalesce[T any](p *T, def T) T {
	if p == nil {
		return def
	}
	return *p
}

// coalesceToDefault coalesces a pointer to the type to which it points. It
// returns the "zero value" of its input's pointed-to type when the input is
// nil. This is equivalent to:
//
//	var x T
//	result := coalesceToDefault(p, x)
//
// ... but can be done on one line without knowing the type of `p`.
func coalesceToDefault[T any](p *T) T {
	var ret T
	if p != nil {
		ret = *p
	}
	return ret
}

// copyIfNotNil makes a deep copy of p - unless it's nil, in which case it just
// returns nil.
func copyIfNotNil[T any](p *T) *T {
	if p == nil {
		return nil
	}
	q := *p
	return &q
}
