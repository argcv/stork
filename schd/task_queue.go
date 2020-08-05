package schd

import (
	"runtime"
	"sync"

	"github.com/argcv/stork/mtx"
)

type taskQueueWorker struct {
	//id string

	q *TaskQueue // global environment
	i int        // current index
	m sync.Mutex // mutex for current worker

	stop    chan struct{}
	stopAck chan struct{}
}

func (w *taskQueueWorker) work() {
	defer close(w.stopAck)

	running := true
	for running {
		select {
		case f, ok := <-w.q.c:
			if ok {
				//log.Infof("exec: #1 %v", w.id)
				f() // perform current job
				runtime.Gosched()
				w.q.wg.Done() // done
			} else {
				running = false
				// this channel is closed
				w.stopAck <- struct{}{}
				//log.Infof("quit #1 %v", w.id)
			}
		case <-w.stop:
			running = false
			w.stopAck <- struct{}{}
			//log.Infof("quit #2 %v", w.id)
		}
	}
}

func newTaskQueueWorker(q *TaskQueue, i int) *taskQueueWorker {
	return &taskQueueWorker{
		q: q,
		i: i,

		stop:    make(chan struct{}),
		stopAck: make(chan struct{}),
	}
}

/**
 * Task Queue: it is used to help us schedule a
 * sequence of executions
 *
 * The default num workers is 1
 * If this number is 1, it could also treated as a
 * strict FIFO Queue
 *
 */
type TaskQueue struct {
	numWorkers int
	workers    []*taskQueueWorker

	c  chan func()
	wg mtx.WaitGroupWithState
	sd *mtx.SingletonDesc
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		numWorkers: 1,
		c:          make(chan func()),
		wg:         mtx.NewWaitGroupWithState(),
		sd:         mtx.NewSingleton(),
	}
}

func (q *TaskQueue) SetNumWorkers(numWorkers int) int {
	if numWorkers > 0 {
		q.numWorkers = numWorkers
	}
	return q.GetNumWorkers()
}

func (q *TaskQueue) GetNumWorkers() int {
	return q.numWorkers
}

func (q *TaskQueue) Perform() {
	go q.sd.Acquire(func() {
		for q.wg.State() > 0 {
			// launch workers
			q.workers = make([]*taskQueueWorker, q.numWorkers)
			for i := 0; i < q.numWorkers; i++ {
				//log.Infof("start worker... %v", i)
				cw := newTaskQueueWorker(q, i)
				//cw.id = xid.New().String()
				q.workers[i] = cw
				go cw.work()
			}

			// waiting for job's finished
			q.wg.Wait()

			for _, w := range q.workers {
				w.stop <- struct{}{}
				close(w.stop)
			}

			for _, w := range q.workers {
				<-w.stopAck // wait for completion
			}
			//log.Infof("stop workers")
		}
	})
}

// add a new task
func (q *TaskQueue) Enqueue(f func()) {
	// to announce a new job is comming
	q.wg.Add(1)
	// try launch the worker
	q.Perform()
	// enqueue the task
	q.c <- f
}

// return current work loader
func (q *TaskQueue) State() int64 {
	return q.wg.State()
}

func (q *TaskQueue) Flush() {
	q.Perform()
	q.wg.Wait()
}

// It's OK to leave a Go channel open forever and never close it.
// When the channel is no longer used, it will be garbage collected.
// -- <https://stackoverflow.com/questions/8593645>
// However we could provide a close + wait interface, which is used
// to indicate its finishing
func (q *TaskQueue) Close() {
	q.Flush()
	close(q.c)
}
