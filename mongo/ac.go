// Mongo Auto Committer is a simple tool to submit
// a batch of mongo request automatically
// It was based on Oliver's bulk_processor.go for Elastic
// see: https://git.io/fjK56
package mongo

import (
	"errors"
	"fmt"
	"github.com/argcv/stork/log"
	"sync"
	"time"
)

/**
 *
 */
type AutoCommitBuilder struct {
	client        *Client // mongo client
	name          string  // name of processor
	coll          string
	numWorkers    int           // # of workers (>= 1)
	bulkActions   int           // # of requests after which to commit
	flushInterval time.Duration // periodic flush interval
	verbose       bool          // verbose
}

// Name is an optional name to identify this bulk processor.
func (s *AutoCommitBuilder) Name(name string) *AutoCommitBuilder {
	s.name = name
	return s
}

// Collection defined the target coll
func (s *AutoCommitBuilder) Collection(coll string) *AutoCommitBuilder {
	s.coll = coll
	return s
}

// Workers is the number of concurrent workers allowed to be
// executed. Defaults to 1 and must be greater or equal to 1.
func (s *AutoCommitBuilder) Workers(num int) *AutoCommitBuilder {
	s.numWorkers = num
	return s
}

// BulkActions specifies when to flush based on the number of actions
// currently added. Defaults to 1000 and can be set to -1 to be disabled.
func (s *AutoCommitBuilder) BulkActions(bulkActions int) *AutoCommitBuilder {
	s.bulkActions = bulkActions
	return s
}

func (s *AutoCommitBuilder) FlushInterval(interval time.Duration) *AutoCommitBuilder {
	s.flushInterval = interval
	return s
}

func (s *AutoCommitBuilder) Verbose(verbose bool) *AutoCommitBuilder {
	s.verbose = verbose
	return s
}

func (s *AutoCommitBuilder) Start() (mac *AutoCommitter, err error) {
	mac = newAutoCommitter(s)
	err = mac.Start()
	if err != nil {
		return nil, err
	}
	return
}

type AutoCommitter struct {
	client        *Client       // mongo client
	name          string        // name of committer # note: this is just a label for verbose saver
	coll          string        // name of coll
	numWorkers    int           // # of workers (>= 1)
	bulkActions   int           // # of requests after which to commit
	flushInterval time.Duration // periodic flush interval
	verbose       bool          // verbose
	workerWg      sync.WaitGroup
	workers       []*mongoAutoCommitWorker
	docsInsert    chan interface{}
	docsUpdate    chan []interface{}
	docsUpsert    chan []interface{}

	// todo: Remove, RemoveAll, UpdateAll

	flusherStopC chan struct{}

	startedMu sync.Mutex // guards the following block
	started   bool
}

func newAutoCommitter(builder *AutoCommitBuilder) *AutoCommitter {
	return &AutoCommitter{
		client:        builder.client,
		name:          builder.name,
		coll:          builder.coll,
		numWorkers:    builder.numWorkers,
		bulkActions:   builder.bulkActions,
		flushInterval: builder.flushInterval,
		verbose:       builder.verbose,
	}
}

func (p *AutoCommitter) Start() error {
	p.startedMu.Lock()
	defer p.startedMu.Unlock()

	if p.started {
		return nil
	}

	// We must have at least one worker.
	if p.numWorkers < 1 {
		p.numWorkers = 1
	}

	p.docsInsert = make(chan interface{})
	p.docsUpdate = make(chan []interface{})
	p.docsUpsert = make(chan []interface{})

	// Create and start up workers.
	p.workers = make([]*mongoAutoCommitWorker, p.numWorkers)
	for i := 0; i < p.numWorkers; i++ {
		p.workerWg.Add(1)
		p.workers[i] = newMongoAutoCommitWorker(p, i)
		go p.workers[i].work()
	}

	// Start the ticker for flush (if enabled)
	if int64(p.flushInterval) > 0 {
		p.flusherStopC = make(chan struct{})
		go p.flusher(p.flushInterval)
	}

	p.started = true

	return nil
}

// Stop is an alias for Close.
func (p *AutoCommitter) Stop() error {
	return p.Close()
}

// Close stops the bulk processor previously started with Do.
// If it is already stopped, this is a no-op and nil is returned.
//
// By implementing Close, BulkProcessor implements the io.Closer interface.
func (p *AutoCommitter) Close() error {
	p.startedMu.Lock()
	defer p.startedMu.Unlock()

	// Already stopped? Do nothing.
	if !p.started {
		return nil
	}

	// Stop flusher (if enabled)
	if p.flusherStopC != nil {
		p.flusherStopC <- struct{}{}
		<-p.flusherStopC
		close(p.flusherStopC)
		p.flusherStopC = nil
	}

	// Stop all workers.
	close(p.docsInsert)
	close(p.docsUpdate)
	close(p.docsUpsert)
	p.workerWg.Wait()

	p.started = false

	return nil
}

// Add adds a single request to commit by the BulkProcessorService.
//
// The caller is responsible for setting the index and type on the request.
func (p *AutoCommitter) Insert(doc interface{}) (e error) {
	if p.started {
		p.docsInsert <- doc
	} else {
		e = errors.New(fmt.Sprintf("AutoCommitter-%s(%s)_is_closed", p.name, p.coll))
	}
	return
}

// Add adds a single request to commit by the BulkProcessorService.
//
// The caller is responsible for setting the index and type on the request.
func (p *AutoCommitter) Update(pair []interface{}) (e error) {
	if p.started {
		p.docsUpdate <- pair
	} else {
		e = errors.New(fmt.Sprintf("AutoCommitter-%s(%s)_is_closed", p.name, p.coll))
	}
	return
}

// Add adds a single request to commit by the BulkProcessorService.
//
// The caller is responsible for setting the index and type on the request.
func (p *AutoCommitter) Upsert(pair []interface{}) (e error) {
	if p.started {
		p.docsUpsert <- pair
	} else {
		e = errors.New(fmt.Sprintf("AutoCommitter-%s(%s)_is_closed", p.name, p.coll))
	}
	return
}

// Flush manually asks all workers to commit their outstanding requests.
// It returns only when all workers acknowledge completion.
func (p *AutoCommitter) Flush() error {
	if p.verbose {
		log.Info(fmt.Sprintf("AutoCommitter-%s(%s) a new flush is comming", p.name, p.coll))
	}
	for _, w := range p.workers {
		w.flushC <- struct{}{}
		<-w.flushAckC // wait for completion
	}
	if p.verbose {
		log.Info(fmt.Sprintf("AutoCommitter-%s(%s) a new flush is finished", p.name, p.coll))
	}
	return nil
}

// flusher is a single goroutine that periodically asks all workers to
// commit their outstanding bulk requests. It is only started if
// FlushInterval is greater than 0.
func (p *AutoCommitter) flusher(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C: // Periodic flush
			p.Flush()
		case <-p.flusherStopC:
			p.flusherStopC <- struct{}{}
			return
		}
	}
}

type mongoAutoCommitWorker struct {
	p         *AutoCommitter
	i         int
	flushC    chan struct{}
	flushAckC chan struct{}

	docMu     sync.Mutex // guards the following block
	docInsert []interface{}
	docUpdate []interface{}
	docUpsert []interface{}
}

func newMongoAutoCommitWorker(p *AutoCommitter, i int) *mongoAutoCommitWorker {
	return &mongoAutoCommitWorker{
		p:         p,
		i:         i,
		flushC:    make(chan struct{}),
		flushAckC: make(chan struct{}),
	}
}

func (w *mongoAutoCommitWorker) bulkActions() int {
	return w.p.bulkActions
}

func (w *mongoAutoCommitWorker) insert(doc interface{}) {
	w.docMu.Lock()
	defer w.docMu.Unlock()
	w.docInsert = append(w.docInsert, doc)
}

func (w *mongoAutoCommitWorker) update(pair []interface{}) {
	w.docMu.Lock()
	defer w.docMu.Unlock()
	w.docUpdate = append(w.docUpdate, pair...)
}

func (w *mongoAutoCommitWorker) upsert(pair []interface{}) {
	w.docMu.Lock()
	defer w.docMu.Unlock()
	w.docUpsert = append(w.docUpsert, pair...)
}

// work waits for bulk requests and manual flush calls on the respective
// channels and is invoked as a goroutine when the bulk processor is started.
func (w *mongoAutoCommitWorker) work() {
	defer func() {
		w.p.workerWg.Done()
		close(w.flushAckC)
		close(w.flushC)
	}()

	var stop bool
	for !stop {
		select {
		case req, open := <-w.p.docsInsert:
			if open {
				// Received a new request
				w.insert(req)
				if w.commitRequired() {
					w.commit()
				}
			} else {
				// Channel closed: Stop.
				stop = true
				w.commit()
			}
		case req, open := <-w.p.docsUpdate:
			if open {
				// Received a new request
				w.update(req)
				if w.commitRequired() {
					w.commit()
				}
			} else {
				// Channel closed: Stop.
				stop = true
				w.commit()
			}
		case req, open := <-w.p.docsUpsert:
			if open {
				// Received a new request
				w.upsert(req)
				if w.commitRequired() {
					w.commit()
				}
			} else {
				// Channel closed: Stop.
				stop = true
				w.commit()
			}

		case <-w.flushC:
			// Commit outstanding requests
			if w.capacity() > 0 {
				w.commit()
			}
			w.flushAckC <- struct{}{}
		}
	}
}

// commit commits the bulk requests in the given service,
// invoking callbacks as specified.
func (w *mongoAutoCommitWorker) commit() (err error) {
	w.docMu.Lock()
	defer w.docMu.Unlock()
	if w.capacity() > 0 {
		if w.p.verbose {
			log.Info(fmt.Sprintf("AutoCommitter-%s(%s)-%d commiting size: %d", w.p.name, w.p.coll, w.i, len(w.docInsert)))
		}
		session := w.p.client.Session.Copy()
		defer session.Close()
		db := w.p.client.Db
		collName := w.p.coll
		coll := session.DB(db).C(collName)
		bulk := coll.Bulk()
		bulk.Unordered()

		if len(w.docInsert) > 0 {
			bulk.Insert(w.docInsert...)
		}

		if len(w.docUpdate) > 0 {
			bulk.Update(w.docUpdate...)
		}

		if len(w.docUpsert) > 0 {
			bulk.Upsert(w.docUpsert...)
		}

		w.docInsert = []interface{}{}
		w.docUpdate = []interface{}{}
		w.docUpsert = []interface{}{}

		if _, err = bulk.Run(); err != nil {
			log.Info(fmt.Sprintf("AutoCommitter-%s(%s)-%d >>ERROR<<: %v", w.p.name, w.p.coll, w.i, err.Error()))
		}
		if w.p.verbose {
			log.Info(fmt.Sprintf("AutoCommitter-%s(%s)-%d committed size: %d", w.p.name, w.p.coll, w.i, len(w.docInsert)))
		}
	} else {
		if w.p.verbose {
			log.Info(fmt.Sprintf("AutoCommitter-%s(%s)-%d committed nothing", w.p.name, w.p.coll, w.i))
		}
	}
	return
}

func (w *mongoAutoCommitWorker) commitRequired() bool {
	if w.bulkActions() >= 0 && w.capacity() >= w.bulkActions() {
		return true
	}
	return false
}

func (w *mongoAutoCommitWorker) capacity() int {
	return len(w.docInsert) + len(w.docUpdate) + len(w.docUpsert)
}
