/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2019 Yu Jing
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */
package mail

import (
	"github.com/argcv/stork/log"
	"github.com/pkg/errors"
	"sync"
	"testing"
	"time"
)

func TestNewSMTPCloser(t *testing.T) {
	log.Verbose()
	c := NewSMTPCloser(func() {
		t.Errorf("really closed")
	})
	c.Delay = 1 * time.Second
	c.Activate()
	time.Sleep(500 * time.Millisecond)
	c.Activate()
	time.Sleep(500 * time.Millisecond)
	c.Activate()
	c.Cancel()
	time.Sleep(1500 * time.Millisecond)

	c.WithLock(&sync.Mutex{})
	c.LockGlobal()
	isLocked := true
	wg := &sync.WaitGroup{}
	wg2 := &sync.WaitGroup{}
	wg.Add(1)
	wg2.Add(1)
	go func() {
		wg.Done()
		c.LockGlobal()
		if isLocked {
			t.Errorf("is locked!!!!")
		}
		wg2.Done()
	}()
	wg.Wait()
	isLocked = false
	c.UnlockGlobal()
	wg2.Wait()

	ed := errors.New("expected error")
	er := c.SafeExec(func() error {
		return ed
	})
	if ed != er {
		t.Errorf("Unexpected error: %v vs. %v", ed, er)
	}

}
