package sigr

import (
	"fmt"
	"github.com/argcv/stork/assert"
	"github.com/argcv/stork/log"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"testing"
	"time"
)

// hmm...
// how to let them appeared in codecov.io?
func TestRegisterOnStopFunc(t *testing.T) {
	log.Verbose()

	if os.Getenv("BE_CRASHER") == "1" {
		log.Verbose()
		isCalled1 := false
		isCalled2 := false

		wg := &sync.WaitGroup{}

		wg.Add(2)

		SetQuitDirectly(false)
		RegisterOnStopFunc("f1", func() {
			log.Infof("called: f1")
			isCalled1 = true
			wg.Done()
		})

		//isCalled1 = true
		//isCalled2 = true

		RegisterOnStopFuncAutoName(func() {
			log.Infof("called: f2")
			isCalled2 = true
			wg.Done()
		})

		log.Infof("Killing")
		if p, e := os.FindProcess(syscall.Getpid()); e == nil {
			log.Infof("pid: %v", p.Pid)
			_ = p.Signal(syscall.SIGINT)
		}
		log.Infof("Killed")
		wg.Wait()
		assert.ExpectTrue(t, isCalled1, fmt.Sprintf("Not called!!! %v", isCalled1))
		assert.ExpectTrue(t, isCalled2, fmt.Sprintf("Not called!!! %v", isCalled2))

		if isCalled1 && isCalled2 {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestRegisterOnStopFunc")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")

	wg := &sync.WaitGroup{}
	wg.Add(1)
	cstop := make(chan bool, 1)
	go func() {
		wg.Done()
		output, err := cmd.CombinedOutput()
		t.Logf("\n----------- OUTPUT ---------------\n%v\n----------- OUTPUT ---------------\n", string(output))
		ee, ok := err.(*exec.ExitError)
		t.Logf("err: %v | %v", err, ee)
		assert.ExpectFalse(t, ok, )
		assert.ExpectEQ(t, nil, err)

		cstop <- true
	}()
	wg.Wait()
	select {
	case stop := <-cstop:
		t.Logf("finished: %v", stop)
	case <-time.After(3 * time.Second):
		t.Errorf("timeout in 3 seconds")
	}
}
