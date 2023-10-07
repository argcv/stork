package schd

import (
	"github.com/argcv/stork/assert"
	"testing"
	"time"
)

func TestNewFifoQueue(t *testing.T) {
	fq := NewFifoQueue()

	prev := -1
	for i := 0; i < 100; i++ {
		ci := i
		fq.Enqueue(func() {
			assert.ExpectEQ(t, 1, ci-prev)
			time.Sleep(1 * time.Millisecond)
			assert.ExpectEQ(t, 1, ci-prev)
			t.Logf("From %v to %v", prev, ci)
			prev = ci
		})
	}

	t.Logf("Close..")
	fq.Close()
	assert.ExpectEQ(t, 99, prev)
}
