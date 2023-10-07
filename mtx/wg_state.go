package mtx

import (
	"sync"
	"sync/atomic"

	"github.com/argcv/stork/log"
)

// Compare with waiting group
// it will return current state
// aka **How many workers are still working**
type WaitGroupWithState interface {
	Add(delta int64) int64
	Done() int64
	State() int64
	Wait()
}

type waitGroupWithStateImpl struct {
	st int64
	cv *sync.Cond
}

func NewWaitGroupWithState() WaitGroupWithState {
	return &waitGroupWithStateImpl{
		st: 0,
		cv: sync.NewCond(&sync.Mutex{}),
	}
}

func (wg *waitGroupWithStateImpl) Add(delta int64) int64 {
	newSt := atomic.AddInt64(&(wg.st), delta)
	if newSt < 0 {
		log.Fatalf("ERROR: status is lower than 0!!! (%v)", newSt)
	}
	return newSt
}

// minus one, return current value
func (wg *waitGroupWithStateImpl) Done() int64 {
	newSt := wg.Add(-1)
	wg.cv.Broadcast()
	return newSt
}

func (wg *waitGroupWithStateImpl) State() int64 {
	return atomic.LoadInt64(&(wg.st))
}

func (wg *waitGroupWithStateImpl) Wait() {
	for wg.State() > 0 {
		wg.cv.L.Lock()
		wg.cv.Wait()
		wg.cv.L.Unlock()
	}
}
