package log

import "testing"

func TestInfo(t *testing.T) {
	SetLevel(ERROR)
	Info("AA", "BB", "CC")
	SetLevel(INFO)
	Info("DD", "EE", "FF")
	Verbose()
	IfDebug(func() {
		Debugf("in debug...")
	})

	IfDebug(func() {
		Debugf("in debug...")
	})

	Debug("a DEBUG message")
	Info("a INFO message")
	Warn("a WARN message")
	Error("a ERROR message")
	Fatal("a FATAL message")
	Debugf("a DEBUG message")
	Infof("a INFO message")
	Warnf("a WARN message")
	Errorf("a ERROR message")
	Fatalf("a FATAL message")

	Debugd(0, "a DEBUG message")
	Infod(0, "a INFO message")
	Warnd(0, "a WARN message")
	Errord(0, "a ERROR message")
	Fatald(0, "a FATAL message")
	Quiet()
	IfDebug(func() {
		t.Fatalf("what happened!!")
	})
	Fatal("Should Be Disabled")

}
