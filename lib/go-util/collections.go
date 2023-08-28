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

// CopyMap makes a "deep-ish" copy of the map passed to it. This will only
// deeply copy the map itself; this means that if the map values are references
// or structures containing references, those are only being shallowly copied!
func CopyMap[T comparable, U any](original map[T]U) map[T]U {
	newMap := make(map[T]U, len(original))
	for k, v := range original {
		newMap[k] = v
	}

	return newMap
}
