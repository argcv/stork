package sigr

import (
	"fmt"
	"github.com/argcv/stork/log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

var onStopService = struct {
	m       sync.Mutex
	handers map[string]func()
	state   int32
}{
	m:       sync.Mutex{},
	handers: map[string]func(){},
	state:   0,
}

var autoIncId uint64 = 0
var quitDirectly = true

func SetQuitDirectly(setting bool) {
	quitDirectly = setting
}

func handlerNameExists(name string) bool {
	onStopService.m.Lock()
	defer onStopService.m.Unlock()
	_, ok := onStopService.handers[name]
	return ok
}

func RegisterOnStopFuncAutoName(f func()) (name string) {
	atomic.AddUint64(&autoIncId, 1)
	name = fmt.Sprintf("$%d", autoIncId)
	for handlerNameExists(name) {
		atomic.AddUint64(&autoIncId, 1)
		name = fmt.Sprintf("$%d", autoIncId)
	}
	RegisterOnStopFunc(name, f)
	return
}

func RegisterOnStopFunc(name string, f func()) {
	// register a new function on signal int(interrupt) and term(terminate)
	onStopService.m.Lock()
	defer onStopService.m.Unlock()
	if atomic.CompareAndSwapInt32(&onStopService.state, 0, 1) {
		log.Debug("OnStopService Initialized..")
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			wg.Done()
			log.Debugf("sig: waiting...")
			var sig os.Signal = <-sigs

			log.Debugf("sig: [%v] Processing...", sig)
			for k, v := range onStopService.handers {
				log.Debugf("Processing task [%v]", k)
				v()
			}
			signal.Stop(sigs)

			if quitDirectly {
				log.Debugf("sig: [%v] Processed, quitting directly", sig)
			} else {
				log.Debugf("sig: [%v] Processed, please try again to terminate the process", sig)
			}

			atomic.CompareAndSwapInt32(&onStopService.state, 1, 0)
			if quitDirectly {
				if p, e := os.FindProcess(syscall.Getpid()); e == nil {
					e = p.Signal(sig)
					if e != nil {
						log.Errorf("sig: pid[%v] send sig %v failed: %v", p.Pid, sig, e)
					}
				}
			}
		}()
		wg.Wait()
	}
	onStopService.handers[name] = f
}

func UnregisterOnStopFunc(name string) {
	onStopService.m.Lock()
	defer onStopService.m.Unlock()
	delete(onStopService.handers, name)
}
