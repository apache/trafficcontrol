package endpoint

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

//go:generate jsonenums -type=Type

// Type models the supported types of endpoints
type Type int

// Type models the supported types of endpoints
const (
	InvalidType Type = iota + 1
	Vod
	Live
	Event
	Static
	Dir
	Testing
)

func (e Type) String() string {
	switch e {
	case InvalidType:
		return "invalid type"
	case Vod:
		return "vod"
	case Live:
		return "live"
	case Event:
		return "event"
	case Static:
		return "static"
	case Dir:
		return "dir"
	case Testing:
		return "testing"
	default:
		return "invalid type"
	}
}
