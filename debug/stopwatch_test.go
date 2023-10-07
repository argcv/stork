package debug

import (
	"testing"
	"time"
)

func TestNewStopwatch(t *testing.T) {
	Verbose()
	Reset()
	time.Sleep(10 * time.Millisecond)
	AddTimePointWithLabel("aa")
	time.Sleep(20 * time.Millisecond)
	AddTimePoint()
	time.Sleep(20 * time.Millisecond)
	AddTimePoint()
	TimeElapsed()
	_, _ = TimeElapsedAfterLabel("aa")

	_ = PrintAll()
}
