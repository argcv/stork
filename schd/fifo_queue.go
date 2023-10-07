package schd

type FifoQueue struct {
	q *TaskQueue
}

func NewFifoQueue() *FifoQueue {
	return &FifoQueue{
		q: NewTaskQueue(),
	}
}

func (q *FifoQueue) Enqueue(f func()) {
	q.q.Enqueue(f)
}

// flush will wait
func (q *FifoQueue) Flush() {
	q.q.Flush()
}

// closing is NOT required in Golang
// It just used to tell the receiver
// everything are sent
func (q *FifoQueue) Close() {
	q.q.Close()
}
