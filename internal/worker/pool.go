// Package worker provides goroutines pool management for concurrent job processing
package worker

import (
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/teamleaderleo/potato-quality-image-compressor/internal/metrics"
)

// Job represents a generic worker job
type Job interface {
	ID() string
	
	Process() (JobResult, error)
}

// JobResult represents the generic result of a job
type JobResult interface {
	ID() string
}

// Pool manages a pool of worker goroutines for parallel processing
type Pool struct {
	jobs           chan jobWrapper
	workerCount    int
	busyWorkers    int32 // atomic counter for metrics
	wg             sync.WaitGroup
	shuttingDown   bool
	shutdownMutex  sync.Mutex
	metricsEnabled bool
}

// jobWrapper wraps a job with its result channel
type jobWrapper struct {
	job    Job
	result chan<- JobResult
	err    chan<- error
}

// NewPool creates a new worker pool with the specified number of workers
func NewPool(workerCount int, jobQueueSize int, enableMetrics bool) *Pool {
	if workerCount <= 0 {
		workerCount = 1
	}
	if jobQueueSize <= 0 {
		jobQueueSize = workerCount * 2
	}

	pool := &Pool{
		jobs:           make(chan jobWrapper, jobQueueSize),
		workerCount:    workerCount,
		metricsEnabled: enableMetrics,
	}

	// Start the workers
	pool.wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go pool.worker(i)
	}

	return pool
}

// worker is the goroutine function that processes jobs
func (p *Pool) worker(id int) {
	defer p.wg.Done()

	for wrapper := range p.jobs {
		if p.metricsEnabled {
			atomic.AddInt32(&p.busyWorkers, 1)
			(*metrics.GetWorkerGauge()).Inc()
		}

		startTime := time.Now()

		// Process the job
		result, err := wrapper.job.Process()
		
		// Send the result
		if err != nil {
			wrapper.err <- err
		} else {
			wrapper.result <- result
		}

		jobTime := time.Since(startTime)

		if p.metricsEnabled {
			atomic.AddInt32(&p.busyWorkers, -1)
			(*metrics.GetWorkerGauge()).Dec()
			(*metrics.GetJobDuration()).Observe(jobTime.Seconds())
		}
	}
}

// Submit adds a job to the worker pool
func (p *Pool) Submit(job Job, resultChan chan<- JobResult, errChan chan<- error) error {
	p.shutdownMutex.Lock()
	if p.shuttingDown {
		p.shutdownMutex.Unlock()
		return io.ErrClosedPipe
	}
	p.shutdownMutex.Unlock()

	p.jobs <- jobWrapper{
		job:    job,
		result: resultChan,
		err:    errChan,
	}
	
	return nil
}

// Shutdown gracefully shuts down the worker pool
func (p *Pool) Shutdown() {
	p.shutdownMutex.Lock()
	if p.shuttingDown {
		p.shutdownMutex.Unlock()
		return
	}
	p.shuttingDown = true
	p.shutdownMutex.Unlock()

	close(p.jobs)
	p.wg.Wait()
}

// BusyWorkerCount returns the number of busy workers
func (p *Pool) BusyWorkerCount() int {
	return int(atomic.LoadInt32(&p.busyWorkers))
}

// TotalWorkerCount returns the total number of workers
func (p *Pool) TotalWorkerCount() int {
	return p.workerCount
}