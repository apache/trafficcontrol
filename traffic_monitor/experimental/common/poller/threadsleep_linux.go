// +build linux

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
