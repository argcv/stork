package schd

import (
	"fmt"
	"github.com/argcv/stork/assert"
	"sync"
	"testing"
	"time"
)

func TestNewTaskQueue(t *testing.T) {
	fq := NewTaskQueue()

	prev := -1
	for i := 0; i < 100; i++ {
		ci := i
		fq.Enqueue(func() {
			assert.ExpectEQ(t, 1, ci-prev, fmt.Sprintf("%v vs. %v", ci, prev))
			//time.Sleep(1 * time.Millisecond)
			assert.ExpectEQ(t, 1, ci-prev, fmt.Sprintf("%v vs. %v", ci, prev))
			t.Logf("From %v to %v", prev, ci)
			prev = ci
		})
	}

	t.Logf("Close..")
	fq.Close()
	assert.ExpectEQ(t, int64(0), fq.State())
	assert.ExpectEQ(t, 99, prev)
}

func TestTaskQueue_SetNumWorkers(t *testing.T) {
	fq := NewTaskQueue()

	assert.ExpectEQ(t, 10, fq.SetNumWorkers(10))

	prev := 0
	mx := sync.Mutex{}

	// batch 1
	for i := 0; i < 100; i++ {
		ci := i
		fq.Enqueue(func() {
			mx.Lock()
			prev += 1
			t.Logf("a' ci: %v prev: %v", ci, prev)
			mx.Unlock()
			//time.Sleep(20 * time.Millisecond)
			mx.Lock()
			prev += 1
			t.Logf("b' ci: %v prev: %v", ci, prev)
			mx.Unlock()
		})
	}

	time.Sleep(100 * time.Millisecond)

	// sleep, and another batch
	for i := 0; i < 100; i++ {
		ci := i
		fq.Enqueue(func() {
			mx.Lock()
			prev += 1
			t.Logf("a' ci: %v prev: %v", ci, prev)
			mx.Unlock()
			time.Sleep(20 * time.Millisecond)
			mx.Lock()
			prev += 1
			t.Logf("b' ci: %v prev: %v", ci, prev)
			mx.Unlock()
		})
	}

	fq.Flush()

	// flush and another batch
	for i := 0; i < 100; i++ {
		ci := i
		fq.Enqueue(func() {
			mx.Lock()
			prev += 1
			t.Logf("a' ci: %v prev: %v", ci, prev)
			mx.Unlock()
			time.Sleep(20 * time.Millisecond)
			mx.Lock()
			prev += 1
			t.Logf("b' ci: %v prev: %v", ci, prev)
			mx.Unlock()
		})
	}

	t.Logf("Staring to Close...")
	fq.Close()
	t.Logf("Closed")
	assert.ExpectEQ(t, 600, prev)
}

func TestTaskQueue_Close(t *testing.T) {
	fq := NewTaskQueue()
	fq.Close()
}
