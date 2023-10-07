// mutex
package mtx

import (
	"github.com/argcv/stork/cntr"
	"github.com/argcv/stork/log"
	"runtime"
	"sync"
	"time"
)

/**
 * CriticalSection Provides a thread safe set, so that
 * we can define a set of lockers in one object
 *
 *
 */
type CriticalSection struct {
	m  sync.Mutex // global mutex
	mj sync.Mutex // mutex for joining
	d  map[string]bool
}

func NewCriticalSection() *CriticalSection {
	return &CriticalSection{
		m:  sync.Mutex{},
		mj: sync.Mutex{},
		d:  map[string]bool{},
	}
}

/**
 * check and release one id
 */
func (c *CriticalSection) ReleaseOne(entry string) {
	c.m.Lock()
	defer c.m.Unlock()
	delete(c.d, entry)
}

/**
 * This is a non-blocking method.
 * Given a set of entries, if all of them are released
 * lock all and return true
 * otherwise do nothing and return false
 */
func (c *CriticalSection) TryLockAll(entries ...string) bool {
	c.m.Lock()
	defer c.m.Unlock()

	// check entries, make sure all of them are unlocked
	for _, entry := range entries {
		if c.unsafeIsLocked(entry) {
			return false
		}
	}

	// lock all
	for _, entry := range entries {
		c.d[entry] = true
	}
	return true
}

/**
 * This is a non-blocking method.
 * Given a set of entries
 * try to lock some of the entries, and return the locked items
 */
func (c *CriticalSection) TryLockPartial(entries ...string) (locked []string) {
	c.m.Lock()
	defer c.m.Unlock()

	// check entries, make sure all of them are unlocked
	for _, entry := range entries {
		if !c.unsafeIsLocked(entry) {
			// is NOT locked

			// lock it
			c.d[entry] = true

			// add to locked list
			locked = append(locked, entry)
		}
	}

	return
}

/**
 * Caution: This is a blocking method
 * Wait until all the functions are locked
 * Maybe it is not a good idea..?
 */
func (c *CriticalSection) Join(entries ...string) {
	c.mj.Lock()
	defer c.mj.Unlock()

	var locked []string
	timeIn := time.Now()
	timeLast := timeIn

	entries = cntr.DistinctStrings(entries...)

	for len(locked) < len(entries) {
		//log.Infof("from: %v", len(locked))
		clocked := c.TryLockPartial(entries...)
		locked = append(locked, clocked...)
		//log.Infof("to: %v", len(locked))
		runtime.Gosched()

		timeCurr := time.Now()
		if timeCurr.Sub(timeLast) > 3*time.Second {
			mCaptured := map[string]bool{}
			var waiting []string
			for _, entry := range locked {
				mCaptured[entry] = true
			}
			for _, entry := range entries {
				if _, ok := mCaptured[entry]; !ok {
					waiting = append(waiting, entry)
				}
			}

			n := 3
			if len(waiting) < 3 {
				n = len(waiting)
			}
			log.Warnf("Possible Deadlock!!! Start Time: %v, obtaining size: %v, missing: %v... in total %v entries",
				timeIn,
				len(entries),
				waiting[:n],
				len(waiting),
			)

			timeLast = timeCurr
		}
	}
}

/**
 * release all the entries
 */
func (c *CriticalSection) Release(entries ...string) {
	for _, entry := range entries {
		c.ReleaseOne(entry)
	}
}

/**
 * return true if all the entries are unlocked
 */
func (c *CriticalSection) Check(entries ...string) bool {
	c.m.Lock()
	defer c.m.Unlock()
	for _, entry := range entries {
		if c.unsafeIsLocked(entry) {
			return false
		}
	}
	return true
}

/**
 * return true if at least 1 of the entries are unlocked
 * NOTE: **NOT** ALL of them are locked
 */
func (c *CriticalSection) Locked(entries ...string) bool {
	return !c.Check(entries...)
}

/**
 * an internal function, a simple helper to check the status of
 * data set. return true if this entry already exists right now
 */
func (c *CriticalSection) unsafeIsLocked(entry string) bool {
	_, ok := c.d[entry]
	return ok
}
