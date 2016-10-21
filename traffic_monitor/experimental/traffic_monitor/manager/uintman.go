package manager

import (
	"sync/atomic"
)

type UintThreadsafe struct {
	val *uint64
}

func NewUintThreadsafe() UintThreadsafe {
	v := uint64(0)
	return UintThreadsafe{val: &v}
}

func (u *UintThreadsafe) Get() uint64 {
	return atomic.LoadUint64(u.val)
}

func (u *UintThreadsafe) Set(v uint64) {
	atomic.StoreUint64(u.val, v)
}

// Inc increments the internal uint64.
// TODO make sure everything using this uses the value it returns, not a separate Get
func (u *UintThreadsafe) Inc() uint64 {
	return atomic.AddUint64(u.val, 1)
}
