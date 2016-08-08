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

func (u UintThreadsafe) Get() uint64 {
	return atomic.LoadUint64(u.val)
}

func (u UintThreadsafe) Set(v uint64) {
	atomic.StoreUint64(u.val, v)
}

func (u UintThreadsafe) Inc() {
	atomic.AddUint64(u.val, 1)
}
