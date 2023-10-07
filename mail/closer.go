package mail

import (
	"github.com/argcv/stork/log"
	"sync"
	"time"
)

type SMTPCloseSch struct {
	Delay   time.Duration
	OnClose func()

	lastUpdate   time.Time
	stopped      bool
	locker       *sync.Mutex
	lockerGlobal *sync.Mutex
	canceled     bool
}

func (c *SMTPCloseSch) SafeExec(f func() error) error {
	c.locker.Lock()
	defer c.locker.Unlock()
	return f()
}

func (c *SMTPCloseSch) WithLock(gm *sync.Mutex) *SMTPCloseSch {
	c.locker.Lock()
	defer c.locker.Unlock()
	prevGm := c.lockerGlobal
	if prevGm != nil {
		prevGm.Lock()
	}
	c.lockerGlobal = gm
	if prevGm != nil {
		prevGm.Unlock()
	}
	return c
}

func (c *SMTPCloseSch) SetDelaySeconds(sec int) *SMTPCloseSch {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.Delay = time.Duration(sec) * time.Second
	return c
}

func (c *SMTPCloseSch) Cancel() {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.canceled = true
}

func (c *SMTPCloseSch) Activate() {
	c.locker.Lock()
	defer c.locker.Unlock()
	if ! c.stopped {
		c.lastUpdate = time.Now()
	} else {
		c.stopped = false
		c.canceled = false
		go func() {
			log.Debugf("New closer session..")
			for {
				log.Debugf("Closer: wait and check")
				time.Sleep(c.Delay + 1*time.Millisecond)
				log.Debugf("Closer: lock and check")
				c.locker.Lock()
				log.Debugf("Closer: locked, check")
				now := time.Now()
				if now.After(c.lastUpdate.Add(c.Delay)) {
					log.Debugf("Closer: close and jump out")
					// stop & close
					c.LockGlobal()

					if c.canceled {
						log.Debugf("canceled")
					} else if c.OnClose != nil {
						c.OnClose()
					}
					c.stopped = true
					c.UnlockGlobal()
					log.Debugf("Closer: closed")

					c.locker.Unlock()
					return
				}
				c.locker.Unlock()
			}
		}()
	}

}

func (c *SMTPCloseSch) LockGlobal() {
	locker := c.lockerGlobal
	if locker != nil {
		locker.Lock()
	}
}

func (c *SMTPCloseSch) UnlockGlobal() {
	locker := c.lockerGlobal
	if locker != nil {
		locker.Unlock()
	}
}

func NewSMTPCloser(onClose func()) *SMTPCloseSch {
	return &SMTPCloseSch{
		Delay:        5 * time.Second, // 30s in default
		OnClose:      onClose,
		lastUpdate:   time.Now(),
		stopped:      true,
		canceled:     false,
		locker:       &sync.Mutex{},
		lockerGlobal: nil,
	}
}
