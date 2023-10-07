package mtx

import (
	"github.com/argcv/stork/assert"
	"testing"
	"time"
)

func TestNewWaitGroupWithState(t *testing.T) {
	wg := NewWaitGroupWithState()
	wg.State()
	assert.ExpectEQ(t, int64(0), wg.State())
	for i := 0; i < 10; i ++ {
		wg.Add(1)
		go func() {
			time.Sleep(1 * time.Millisecond)
			assert.ExpectLT(t, int64(-1), wg.Done())
		}()
	}
	assert.ExpectLT(t, int64(-1), wg.State())
	wg.Wait()
	assert.ExpectEQ(t, int64(0), wg.State())
}
