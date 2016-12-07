// +build !linux

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
