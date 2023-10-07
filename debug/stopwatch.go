package debug

import (
	"fmt"
	"github.com/argcv/stork/log"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var (
	defaultStopWatch = NewStopwatch()
)

type Stopwatch struct {
	TimePoints []time.Time
	Labels     map[string]int
	LabelsRev  map[int]string
	verbose    bool
	mx         sync.Mutex
}

func NewStopwatch() *Stopwatch {
	return &Stopwatch{
		TimePoints: []time.Time{time.Now()},
		Labels:     map[string]int{},
		LabelsRev:  map[int]string{},
		verbose:    false,
		mx:         sync.Mutex{},
	}
}

func (sw *Stopwatch) addTimePoint(label string) time.Time {
	sw.mx.Lock()
	defer sw.mx.Unlock()
	now := time.Now()
	idx := len(sw.TimePoints)
	sw.TimePoints = append(sw.TimePoints, now)

	labelSuffix := ""
	if len(label) > 0 {
		sw.Labels[label] = idx
		sw.LabelsRev[idx] = label
		labelSuffix = fmt.Sprintf(" label:[%v]", label)
	}

	if sw.verbose {
		log.Infof("New time point: %v%v", now, labelSuffix)
	}
	return now
}

func (sw *Stopwatch) AddTimePoint() time.Time {
	return sw.addTimePoint("")
}

func (sw *Stopwatch) AddTimePointWithLabel(label string) time.Time {
	return sw.addTimePoint(label)
}

func (sw *Stopwatch) unsafeGetTimeByLabel(label string) (time.Time, bool) {
	if val, ok := sw.Labels[label]; ok {
		return sw.TimePoints[val], ok
	} else {
		return time.Unix(0, 0), false
	}
}

func (sw *Stopwatch) GetTimeByLabel(label string) (time.Time, bool) {
	sw.mx.Lock()
	defer sw.mx.Unlock()
	return sw.unsafeGetTimeByLabel(label)
}

func (sw *Stopwatch) unsafeGetLabelByIndex(idx int) string {
	if val, ok := sw.LabelsRev[idx]; ok {
		return val
	} else {
		return ""
	}
}

// return elapsed time
// it will return 2 results, the first is the duration before the previous one
// the second is the duration of all
func (sw *Stopwatch) TimeElapsed() (d1 time.Duration, d2 time.Duration) {
	sw.mx.Lock()
	defer sw.mx.Unlock()
	d1 = 0
	d2 = 0
	ltp := len(sw.TimePoints)
	if ltp == 0 {
		if sw.verbose {
			log.Errorf("Invalid time points!!, the counter is zero")
		}
		return
	}

	if ltp > 1 {
		d1 = sw.TimePoints[ltp-1].Sub(sw.TimePoints[ltp-2])
	}
	d2 = sw.TimePoints[ltp-1].Sub(sw.TimePoints[0])
	if sw.verbose {
		log.Infof("Duration [%v] before prev[%v]: %v, duration of all: %v",
			sw.unsafeGetLabelByIndex(ltp-1),
			sw.unsafeGetLabelByIndex(ltp-2),
			d1, d2)
	}
	return
}

func (sw *Stopwatch) TimeElapsedAfterLabel(label string) (dur time.Duration, err error) {
	sw.mx.Lock()
	defer sw.mx.Unlock()
	dur = 0
	ltp := len(sw.TimePoints)
	if ltp == 0 {
		err = errors.New("Invalid time points!!, the counter is zero")
		if sw.verbose {
			log.Errorf("%v", err)
		}
		return
	}

	ctime := sw.TimePoints[ltp-1]
	target, ok := sw.unsafeGetTimeByLabel(label)

	if !ok {
		err = errors.New("Invalid label!!, target label not found")
		if sw.verbose {
			log.Errorf("%v", err)
		}
		return
	}

	dur = ctime.Sub(target)
	if sw.verbose {
		log.Infof("Duration [%v] before label[%v]: %v",
			sw.unsafeGetLabelByIndex(ltp-1),
			label,
			dur)
	}
	return
}

func (sw *Stopwatch) PrintAll() error {
	sw.mx.Lock()
	defer sw.mx.Unlock()
	if len(sw.TimePoints) == 0 {
		err := errors.New("Empty time points!!??")
		if sw.verbose {
			log.Errorf("%v", err)
		}
		return err
	}

	prev := sw.TimePoints[0]
	first := sw.TimePoints[0]
	iprev := 0
	lprev := sw.unsafeGetLabelByIndex(0)
	log.Infof("--------- Print All start --------")
	for i, tp := range sw.TimePoints {
		if i == 0 {
			continue
		}
		label := sw.unsafeGetLabelByIndex(i)
		log.Infof("Dur Gap: %v:%v => %v:%v -- %v;%v",
			iprev, lprev,
			i, label,
			tp.Sub(prev), tp.Sub(first))
		prev = tp
		lprev = label
		iprev = i
	}
	log.Infof("--------- Print All end   --------")
	return nil
}

func (sw *Stopwatch) Reset() {
	sw.mx.Lock()
	defer sw.mx.Unlock()
	sw.TimePoints = []time.Time{time.Now()}
	sw.Labels = map[string]int{}
	sw.LabelsRev = map[int]string{}
}

func AddTimePoint() time.Time {
	return defaultStopWatch.AddTimePoint()
}

func AddTimePointWithLabel(label string) time.Time {
	return defaultStopWatch.AddTimePointWithLabel(label)
}

func Verbose() {
	defaultStopWatch.verbose = true
}

func Reset() {
	defaultStopWatch.Reset()
}

func TimeElapsed() (d1 time.Duration, d2 time.Duration) {
	return defaultStopWatch.TimeElapsed()
}

func TimeElapsedAfterLabel(label string) (dur time.Duration, err error) {
	return defaultStopWatch.TimeElapsedAfterLabel(label)
}

func PrintAll() error {
	return defaultStopWatch.PrintAll()
}
