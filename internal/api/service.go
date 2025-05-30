package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/teamleaderleo/potato-quality-image-compressor/internal/compression"
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/config"
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
	Filename         string  // Added to store the original filename
	Format           string  // Added to store the output format
}

// Service handles the API endpoints for image compression
type Service struct {
	workerPool             *worker.Pool
	processor              *compression.ImageProcessor
	defaultQuality         int
	defaultFormat          string
	defaultAlgorithm       string
	imageProcessingTimeout time.Duration
	batchProcessingTimeout time.Duration
	maxUploadSize          int64
	maxBatchSize           int
}

// NewServiceWithConfig creates a new service with the given configuration
func NewServiceWithConfig(config config.ServiceConfig) *Service {
	// Create worker pool
	workerPool := worker.NewPool(config.WorkerCount, config.JobQueueSize, config.EnableMetrics)

	// Create image processor
	processor := compression.NewImageProcessor()

	// Set default algorithm if specified
	if config.DefaultAlgorithm != "scale" {
		processor.SetDefaultAlgorithm(config.DefaultAlgorithm)
	}

	// Start resource monitoring if metrics are enabled
	if config.EnableMetrics {
		StartResourceMonitor(10 * time.Second)
	}

	return &Service{
		workerPool:             workerPool,
		processor:              processor,
		defaultQuality:         config.DefaultQuality,
		defaultFormat:          config.DefaultFormat,
		defaultAlgorithm:       config.DefaultAlgorithm,
		imageProcessingTimeout: config.ImageProcessingTimeout,
		batchProcessingTimeout: config.BatchProcessingTimeout,
		maxUploadSize:          config.MaxUploadSize,
		maxBatchSize:           config.MaxBatchSize,
	}
}



// CompressImage processes an image directly and returns the result
// Core implementation used by both HTTP and gRPC handlers
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
		ctx, cancel = context.WithTimeout(ctx, s.imageProcessingTimeout)
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
		
		// Return the result with the new fields
		return CompressionResult{
			Data:             compressionResult.Data(),
			Error:            nil,
			ProcessingTime:   time.Since(startTime),
			OriginalSize:     compressionResult.OriginalSize(),
			CompressedSize:   compressionResult.CompressedSize(),
			CompressionRatio: compressionResult.CompressionRatio(),
			AlgorithmUsed:    compressionResult.AlgorithmUsed(),
			Filename:         filename,
			Format:           format,
		}, nil
		
	case err := <-errChan:
		return CompressionResult{Error: fmt.Errorf("processing job: %w", err)}, err

	case <-ctx.Done():
		return CompressionResult{Error: ctx.Err()}, ctx.Err()
	}
}

// parseParameters parses and validates request parameters
func (s *Service) parseParameters(r *http.Request) (int, string, string) {
	// Parse quality
	quality, err := validateQuality(r.FormValue("quality"), s.defaultQuality)
	if err != nil {
		quality = s.defaultQuality
	}

	// Parse format
	format := validateFormat(r.FormValue("format"), s.defaultFormat)

	// Parse algorithm
	algorithm := validateAlgorithm(r.FormValue("algorithm"), s.defaultAlgorithm)

	return quality, format, algorithm
}

// validateQuality validates the quality parameter
func validateQuality(qualityStr string, defaultQuality int) (int, error) {
	if qualityStr == "" {
		return defaultQuality, nil
	}

	quality, err := strconv.Atoi(qualityStr)
	if err != nil {
		return defaultQuality, err
	}

	if quality < 1 || quality > 100 {
		return defaultQuality, fmt.Errorf("quality must be between 1 and 100")
	}

	return quality, nil
}

// validateFormat validates the format parameter
func validateFormat(format, defaultFormat string) string {
	if format == "" {
		return defaultFormat
	}

	// Add validation logic for supported formats
	validFormats := map[string]bool{
		"webp": true,
		"jpeg": true,
		"jpg":  true,
		"png":  true,
	}
	
	if !validFormats[format] {
		return defaultFormat
	}
	
	return format
}

// validateAlgorithm validates the algorithm parameter
func validateAlgorithm(algorithm, defaultAlgorithm string) string {
	if algorithm == "" {
		return defaultAlgorithm
	}
	
	// Add validation logic for supported algorithms
	validAlgorithms := map[string]bool{
		"scale":      true,
		"qualitymod": true,
	}
	
	if !validAlgorithms[algorithm] {
		return defaultAlgorithm
	}
	
	return algorithm
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

// Shutdown gracefully shuts down the service
func (s *Service) Shutdown() {
	if s.workerPool != nil {
		s.workerPool.Shutdown()
	}
}