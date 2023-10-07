package schd

import (
	"context"
	"testing"
	"time"
)

func TestMultiTaskTicker_Start(t *testing.T) {
	c1 := 0
	c2 := 0
	mtt := NewMultiTaskTicker()
	mtt.SetNumWorkers(3)
	mtt.SetPeriod(10 * time.Millisecond)
	mtt.AddTask("1", "2")
	mtt.SetTasks("1", "2")
	st := mtt.Start(context.Background(), func(ctx context.Context, label interface{}) {
		switch label {
		case "1":
			c1 += 1
			t.Logf("c1 => %v", c1)
			time.Sleep(1 * time.Second)
		case "2":
			c2 += 2
			t.Logf("c2 => %v", c2)
		default:
			t.Errorf("Unknown label: %v", label)
		}
	})
	if st != nil {
		t.Errorf("failed: %v", st)
	}
	time.Sleep(100 * time.Millisecond)
	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	st = mtt.Stop(ctx)
	if st != nil {
		t.Errorf("failed: %v", st)
	}
	t.Logf("c1: %v c2: %v", c1, c2)
}
