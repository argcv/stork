package mtx

import (
	"github.com/argcv/stork/assert"
	"sync"
	"testing"
	"time"
)

func TestNewSingleton(t *testing.T) {
	s := NewSingleton()
	i := 0
	mx := sync.Mutex{}
	go s.Acquire(func() {
		time.Sleep(10 * time.Millisecond)

		mx.Lock()
		i += 1 //
		mx.Unlock()
	})

	s.Acquire(func() {
		mx.Lock()
		i += 1
		mx.Unlock()

		time.Sleep(10 * time.Millisecond)
	})
	time.Sleep(15 * time.Millisecond)

	mx.Lock()
	assert.ExpectEQ(t, 1, i)
	mx.Unlock()
}
