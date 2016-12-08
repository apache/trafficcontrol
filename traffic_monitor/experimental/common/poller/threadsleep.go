// +build !linux

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

package poller

import (
	"runtime"
	"time"
)

// ThreadSleep actually busywaits for the given duration. This is becuase Go doesn't have Mac and Windows nanosleep syscalls, and `Sleep` sleeps for progressively longer than requested.
func ThreadSleep(d time.Duration) {
	// TODO fix to not busywait on Mac, Windows. We can't simply Sleep, because Sleep gets progressively slower as the app runs, due to a Go runtime issue. If this is changed, you MUST verify the poll doesn't get slower after the app runs for several days.
	end := time.Now().Add(d)
	for end.After(time.Now()) {
		runtime.Gosched()
	}
}
