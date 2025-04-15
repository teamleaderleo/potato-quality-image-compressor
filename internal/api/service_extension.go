package api

import (
	"fmt"
	"io"
	"bytes"
	"time"
	
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/compression"
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/worker"
)

// CompressionResult represents the result of a compression operation
type CompressionResult struct {
	Data            []byte
	Error           error
	ProcessingTime  time.Duration
	OriginalSize    int
	CompressedSize  int
	CompressionRatio float64
	AlgorithmUsed   string
}

// CompressImageDirect processes an image directly and returns the result
// Used by the gRPC adapter for synchronous processing
func (s *Service) CompressImageDirect(
	filename string,
	input io.Reader,
	format string,
	quality int,
	algorithm string,
) (CompressionResult, error) {
	// Read the input data
	inputData, err := io.ReadAll(input)
	if err != nil {
		return CompressionResult{Error: err}, err
	}
	
	// Create a job
	job := compression.NewCompressionJob(
		filename,
		bytes.NewReader(inputData),
		format,
		quality,
		algorithm,
		s.processor,
	)
	
	// Create channels for result and error
	resultChan := make(chan worker.JobResult, 1)
	errChan := make(chan error, 1)
	
	// Start time measurement
	startTime := time.Now()
	
	// Submit job to worker pool
	err = s.workerPool.Submit(job, resultChan, errChan)
	if err != nil {
		return CompressionResult{Error: err}, err
	}
	
	// Wait for result or error
	select {
	case result := <-resultChan:
		// Type assertion
		compressionResult, ok := result.(*compression.CompressionResult)
		if !ok {
			err := fmt.Errorf("invalid result type")
			return CompressionResult{Error: err}, err
		}
		
		// Return the result
		return CompressionResult{
			Data:            compressionResult.Data(),
			Error:           nil,
			ProcessingTime:  time.Since(startTime),
			OriginalSize:    compressionResult.OriginalSize(),
			CompressedSize:  compressionResult.CompressedSize(),
			CompressionRatio: compressionResult.CompressionRatio(),
			AlgorithmUsed:   compressionResult.AlgorithmUsed(),
		}, nil
		
	case err := <-errChan:
		return CompressionResult{Error: err}, err
		
	case <-time.After(30 * time.Second):
		err := fmt.Errorf("processing timeout")
		return CompressionResult{Error: err}, err
	}
}

// GetWorkerCount returns the total number of workers
func (s *Service) GetWorkerCount() int {
	return s.workerPool.TotalWorkerCount()
}

// GetBusyWorkerCount returns the number of busy workers
func (s *Service) GetBusyWorkerCount() int {
	return s.workerPool.BusyWorkerCount()
}