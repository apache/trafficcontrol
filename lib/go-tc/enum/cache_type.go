package enum

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

import "strings"

// CacheType is the type (or tier) of a CDN cache.
type CacheType string

const OriginLocationType = "ORG_LOC"

const (
	// CacheTypeEdge represents an edge cache.
	CacheTypeEdge = CacheType("EDGE")
	// CacheTypeMid represents a mid cache.
	CacheTypeMid = CacheType("MID")
	// CacheTypeInvalid represents an cache type enumeration. Note this is the default construction for a CacheType.
	CacheTypeInvalid = CacheType("")
)

// String returns a string representation of this cache type.
func (t CacheType) String() string {
	switch t {
	case CacheTypeEdge:
		return "EDGE"
	case CacheTypeMid:
		return "MID"
	default:
		return "INVALIDCACHETYPE"
	}
}

// CacheTypeFromString returns a cache type object from its string representation, or CacheTypeInvalid if the string is not a valid type.
func CacheTypeFromString(s string) CacheType {
	s = strings.ToLower(s)
	if strings.HasPrefix(s, "edge") {
		return CacheTypeEdge
	}
	if strings.HasPrefix(s, "mid") {
		return CacheTypeMid
	}
	return CacheTypeInvalid
}
