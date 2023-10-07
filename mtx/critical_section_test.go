package mtx

import (
	"testing"
	"time"
)

func TestNewCriticalSection(t *testing.T) {
	if NewCriticalSection() == nil {
		t.Errorf("init failed")
	}
}

func TestCriticalSection_Check(t *testing.T) {
	cs := NewCriticalSection()

	if !cs.Check("aa") {
		t.Errorf("check failed #1")
	}

	cs.Join("aa", "bb")

	if cs.Check("aa") {
		t.Errorf("check failed #1")
	}
}

func TestCriticalSection_Join(t *testing.T) {
	cs := NewCriticalSection()
	cs.Join("aa")
	isUnlocked := false
	go func() {
		time.Sleep(100 * time.Millisecond)
		isUnlocked = true
		cs.ReleaseOne("aa")
	}()

	if cs.Check("aa") {
		t.Errorf("lock failed")
	}

	cs.Join("aa")
	if isUnlocked == false {
		t.Errorf("join failed")
	}

}

func TestCriticalSection_Join2(t *testing.T) {
	cs := NewCriticalSection()
	cs.Join("aa", "bb", "aa")
	isUnlocked := false
	go func() {
		time.Sleep(20 * time.Millisecond)
		isUnlocked = true
		cs.ReleaseOne("aa")
		cs.ReleaseOne("bb")
	}()

	if cs.Check("aa") {
		t.Errorf("lock failed")
	}

	if cs.Check("bb") {
		t.Errorf("lock failed")
	}

	cs.Join("aa", "bb")
	if isUnlocked == false {
		t.Errorf("join failed")
	}
}

func TestCriticalSection_JoinWarning(t *testing.T) {
	cs := NewCriticalSection()
	cs.Join("aa")
	isUnlocked := false
	go func() {
		time.Sleep(7 * time.Second)
		isUnlocked = true
		cs.ReleaseOne("aa")
	}()

	if cs.Check("aa") {
		t.Errorf("lock failed")
	}

	cs.Join("aa")
	if isUnlocked == false {
		t.Errorf("join failed")
	}

}

func TestCriticalSection_TryLockAll(t *testing.T) {
	cs := NewCriticalSection()
	cs.Join("aa")

	if cs.TryLockAll("aa", "bb") {
		t.Errorf("check failed")
	}

	if cs.Check("aa") {
		t.Errorf("lock failed")
	}

	if cs.Locked("bb") {
		t.Errorf("incorrectly locked")
	}

	cs.Release("aa", "bb")

	if cs.Locked("aa") {
		t.Errorf("incorrectly failed")
	}

	if cs.Locked("bb") {
		t.Errorf("incorrectly locked")
	}

	if !cs.TryLockAll("aa", "bb") {
		t.Errorf("lock failed")
	}

	if cs.Check("aa") {
		t.Errorf("lock failed")
	}

	if cs.Check("bb") {
		t.Errorf("lock locked")
	}

}

func TestCriticalSection_TryLockPartial(t *testing.T) {
	cs := NewCriticalSection()
	cs.Join("aa")

	if cs.Check("aa") {
		t.Errorf("lock failed")
	}

	locked := cs.TryLockPartial("aa", "bb")
	if len(locked) != 1 {
		t.Fatalf("incorrect partial locked size")
	}

	if locked[0] != "bb" {
		t.Errorf("incorrect partial locked string")
	}

	if cs.Check("aa") {
		t.Errorf("lock failed")
	}

	if cs.Check("bb") {
		t.Errorf("lock failed")
	}

}
