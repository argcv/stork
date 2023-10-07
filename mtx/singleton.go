package mtx

import (
	"sync/atomic"
)

/**
 * SingletonDesc provides a simple way to **SKIP** one function
 * If another function is in executing
 *
 */
type SingletonDesc struct {
	cnt int32
}

func NewSingleton() *SingletonDesc {
	return &SingletonDesc{
		cnt: 0,
	}
}

// We can call Acquire multiple times, however, there is only
// 1 running callback
func (s *SingletonDesc) Acquire(f func()) {
	if atomic.CompareAndSwapInt32(&s.cnt, 0, 1) {
		defer atomic.StoreInt32(&s.cnt, 0)
		f()
	}
}
