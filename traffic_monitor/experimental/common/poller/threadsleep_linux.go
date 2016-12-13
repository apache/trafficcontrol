// +build linux

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
	"errors"
	"golang.org/x/sys/unix"
	"time"
)

// ThreadSleep sleeps using the POSIX syscall `nanosleep`. Note this does not sleep the goroutine, but the operating system thread itself. This should only be called by a goroutine which has previously called `LockOSThread`. This exists due to a bug with `time.Sleep` getting progressively slower as the app runs, and should be removed if the bug in Go is fixed.
func ThreadSleep(d time.Duration) {
	if d < 0 {
		d = 0
	}
	t := unix.Timespec{}
	leftover := unix.NsecToTimespec(d.Nanoseconds())
	err := errors.New("")
	for err != nil && (leftover.Sec != 0 || leftover.Nsec != 0) {
		t = leftover
		err = unix.Nanosleep(&t, &leftover)
	}
}
