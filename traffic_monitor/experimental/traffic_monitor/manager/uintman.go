package manager

import (
	"sync/atomic"
)

// UintThreadsafe provides safe access for multiple goroutines readers and a single writer to a stored uint.
type UintThreadsafe struct {
	val *uint64
}

// NewUintThreadsafe returns a new single-writer-multiple-reader threadsafe uint
func NewUintThreadsafe() UintThreadsafe {
	v := uint64(0)
	return UintThreadsafe{val: &v}
}

// Get gets the internal uint. This is safe for multiple readers
func (u *UintThreadsafe) Get() uint64 {
	return atomic.LoadUint64(u.val)
}

// Set sets the internal uint. This MUST NOT be called by multiple goroutines.
func (u *UintThreadsafe) Set(v uint64) {
	atomic.StoreUint64(u.val, v)
}

// Inc increments the internal uint64.
// TODO make sure everything using this uses the value it returns, not a separate Get
func (u *UintThreadsafe) Inc() uint64 {
	return atomic.AddUint64(u.val, 1)
}
