package util

import (
	"sync/atomic"
)

type AtomicBool int32

func NewAtomicBool(b bool) *AtomicBool {
	a := new(AtomicBool)
	a.Set(b)
	return a
}

func (a *AtomicBool) Set(b bool) {
	if b {
		atomic.StoreInt32((*int32)(a), 1)
	} else {
		atomic.StoreInt32((*int32)(a), 0)
	}
}

func (a *AtomicBool) Enabled() bool {
	return atomic.LoadInt32((*int32)(a)) == 1
}

func (a *AtomicBool) Disabled() bool {
	return atomic.LoadInt32((*int32)(a)) == 0
}
