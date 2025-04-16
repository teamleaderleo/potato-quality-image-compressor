package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"
	
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/compression"
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/worker"
)

// Common errors
var (
	ErrInvalidResultType = errors.New("invalid result type")
	ErrProcessingTimeout = errors.New("processing timeout")
)

// CompressionResult represents the result of a compression operation
type CompressionResult struct {
	Data             []byte
	Error            error
	ProcessingTime   time.Duration
	OriginalSize     int
	CompressedSize   int
	CompressionRatio float64
	AlgorithmUsed    string
}

// CompressImage processes an image directly and returns the result
// Used by the gRPC adapter for synchronous processing
func (s *Service) CompressImage(
	ctx context.Context,
	filename string,
	input io.Reader,
	format string,
	quality int,
	algorithm string,
) (CompressionResult, error) {
	// Create a new context with a timeout if not already set
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}

	// Read the input data
	inputData, err := io.ReadAll(input)
	if err != nil {
		return CompressionResult{Error: fmt.Errorf("reading input: %w", err)}, err
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
	if err := s.workerPool.Submit(job, resultChan, errChan); err != nil {
		return CompressionResult{Error: fmt.Errorf("submitting job: %w", err)}, err
	}

	// Wait for result, error, or context cancellation
	select {
	case result := <-resultChan:
		// Type assertion
		compressionResult, ok := result.(*compression.CompressionResult)
		if !ok {
			return CompressionResult{Error: ErrInvalidResultType}, ErrInvalidResultType
		}
		
		// Return the result
		return CompressionResult{
			Data:             compressionResult.Data(),
			Error:            nil,
			ProcessingTime:   time.Since(startTime),
			OriginalSize:     compressionResult.OriginalSize(),
			CompressedSize:   compressionResult.CompressedSize(),
			CompressionRatio: compressionResult.CompressionRatio(),
			AlgorithmUsed:    compressionResult.AlgorithmUsed(),
		}, nil
		
	case err := <-errChan:
		return CompressionResult{Error: fmt.Errorf("processing job: %w", err)}, err

	case <-ctx.Done():
		return CompressionResult{Error: ctx.Err()}, ctx.Err()
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

// GetServiceHealth returns the health status of the service
func (s *Service) GetServiceHealth() bool {
	// Check if worker pool is operational
	return s.workerPool != nil
}
