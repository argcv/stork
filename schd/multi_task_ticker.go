package schd

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/argcv/stork/log"
)

type runningState struct {
	bucket map[int]bool
	mx     *sync.RWMutex
}

func newRunningState() *runningState {
	return &runningState{
		bucket: map[int]bool{},
		mx:     &sync.RWMutex{},
	}
}

func (rs *runningState) add(id int) bool {
	rs.mx.Lock()
	defer rs.mx.Unlock()
	if _, ok := rs.bucket[id]; ok {
		// add failed
		return false
	} else {
		rs.bucket[id] = true
		return true
	}
}

func (rs *runningState) remove(id int) {
	rs.mx.Lock()
	defer rs.mx.Unlock()
	rs.bucket[id] = true
	delete(rs.bucket, id)
}

func (rs *runningState) check(id int) bool {
	rs.mx.RLock()
	defer rs.mx.RUnlock()
	_, ok := rs.bucket[id]
	return ok
}

func (rs *runningState) stat() int {
	rs.mx.RLock()
	defer rs.mx.RUnlock()
	return len(rs.bucket)
}

type MultiTaskTickerFunc func(ctx context.Context, param interface{})

type MultiTaskTicker struct {
	period    time.Duration
	nWorkers  int
	params    []interface{}
	rs        *runningState
	cancel    context.CancelFunc
	wg        *sync.WaitGroup
	m         *sync.Mutex
	isStarted bool
}

func NewMultiTaskTicker() *MultiTaskTicker {
	return &MultiTaskTicker{
		period:    100 * time.Millisecond,
		nWorkers:  10,
		params:    nil,
		rs:        newRunningState(),
		wg:        &sync.WaitGroup{},
		m:         &sync.Mutex{},
		isStarted: false,
	}
}

func (mtt *MultiTaskTicker) SetPeriod(period time.Duration) *MultiTaskTicker {
	mtt.period = period
	return mtt
}

func (mtt *MultiTaskTicker) SetNumWorkers(nWorkers int) *MultiTaskTicker {
	mtt.nWorkers = nWorkers
	return mtt
}

func (mtt *MultiTaskTicker) GetNumWorkers() int {
	return mtt.nWorkers
}

func (mtt *MultiTaskTicker) AddTask(params ...interface{}) *MultiTaskTicker {
	mtt.params = append(mtt.params, params...)
	return mtt
}

func (mtt *MultiTaskTicker) SetTasks(params ...interface{}) *MultiTaskTicker {
	mtt.params = params
	return mtt
}

func (mtt *MultiTaskTicker) Start(ctx context.Context, f MultiTaskTickerFunc) (err error) {
	mtt.m.Lock()
	defer mtt.m.Unlock()
	if mtt.isStarted {
		return errors.New("already_started")
	}
	mtt.isStarted = true
	var cctx context.Context
	cctx, mtt.cancel = context.WithCancel(ctx)
	wkr := NewTaskQueue()
	wkr.SetNumWorkers(mtt.nWorkers + 1)
	wkr.Enqueue(func() {
		mtt.wg.Add(1)
		defer mtt.wg.Done()
		ticker := time.NewTicker(mtt.period)
		for {
			select {
			case <-cctx.Done():
				log.Infof("canceled...")
				return
			case <-ticker.C:
				params := mtt.params
				for id := range params {
					cid := id
					param := params[id]
					if mtt.rs.add(cid) {
						// start it
						mtt.wg.Add(1)
						wkr.Enqueue(func() {
							defer mtt.wg.Done()
							f(cctx, param)
							mtt.rs.remove(cid)
						})
					} else {
						// task adding failed
					}
				}
			}
		}
	})
	return
}

func (mtt *MultiTaskTicker) Stop(ctx context.Context) (err error) {
	mtt.m.Lock()
	defer mtt.m.Unlock()
	if !mtt.isStarted {
		return errors.New("not_started")
	}
	stop := make(chan bool)
	go func() {
		mtt.cancel()
		mtt.wg.Wait()
		mtt.isStarted = false
		stop <- true
	}()
	select {
	case <-ctx.Done():
		return errors.New("timeout")
	case <-stop:
	}
	return
}

func (mtt *MultiTaskTicker) StopWait() {
	mtt.wg.Wait()
	return
}

func (mtt *MultiTaskTicker) IsStarted() bool {
	mtt.m.Lock()
	defer mtt.m.Unlock()
	return mtt.isStarted
}
