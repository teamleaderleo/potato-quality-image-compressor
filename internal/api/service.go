package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/teamleaderleo/potato-quality-image-compressor/internal/compression"
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/metrics"
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

// Service handles the API endpoints for image compression
type Service struct {
	workerPool       *worker.Pool
	processor        *compression.ImageProcessor
	defaultQuality   int
	defaultFormat    string
	defaultAlgorithm string
}

// ServiceConfig contains configuration for the Service
type ServiceConfig struct {
	WorkerCount      int
	JobQueueSize     int
	DefaultQuality   int
	DefaultFormat    string
	DefaultAlgorithm string
	EnableMetrics    bool
}

// DefaultServiceConfig returns the default service configuration
func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{
		WorkerCount:      runtime.NumCPU(),
		JobQueueSize:     runtime.NumCPU() * 4,
		DefaultQuality:   75,
		DefaultFormat:    "webp",
		DefaultAlgorithm: "scale",
		EnableMetrics:    true,
	}
}

// NewService creates a new service with default configuration
func NewService() *Service {
	return NewServiceWithConfig(DefaultServiceConfig())
}

// NewServiceWithConfig creates a new service with the given configuration
func NewServiceWithConfig(config ServiceConfig) *Service {
	// Create worker pool
	workerPool := worker.NewPool(config.WorkerCount, config.JobQueueSize, config.EnableMetrics)

	// Create image processor
	processor := compression.NewImageProcessor()

	// Set default algorithm if specified
	if config.DefaultAlgorithm != "scale" {
		processor.SetDefaultAlgorithm(config.DefaultAlgorithm)
	}

	// Start a goroutine to periodically update memory usage metrics
	if config.EnableMetrics {
		go monitorMemoryUsage()
	}

	return &Service{
		workerPool:       workerPool,
		processor:        processor,
		defaultQuality:   config.DefaultQuality,
		defaultFormat:    config.DefaultFormat,
		defaultAlgorithm: config.DefaultAlgorithm,
	}
}

// monitorMemoryUsage periodically updates memory usage metrics
func monitorMemoryUsage() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	var m runtime.MemStats
	for range ticker.C {
		runtime.ReadMemStats(&m)
		metrics.UpdateMemoryUsage(m.Alloc)
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

// HandleCompress handles single image compression requests via HTTP
func (s *Service) HandleCompress(w http.ResponseWriter, r *http.Request) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()
	
	timer := metrics.NewTimer("compress")
	defer timer.ObserveDuration()

	status := "success"
	defer func() {
		metrics.GetRequestCounter().WithLabelValues("compress", status).Inc()
	}()

	if r.Method != http.MethodPost {
		status = "method_not_allowed"
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		status = "bad_request"
		http.Error(w, "Error retrieving the file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file into memory to get original size for metrics
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		status = "read_error"
		http.Error(w, "Error reading file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse parameters
	quality, format, algorithm := s.parseParameters(r)

	// Process the image using the core CompressImage method
	result, err := s.CompressImage(
		ctx,
		header.Filename,
		bytes.NewReader(fileBytes),
		format,
		quality,
		algorithm,
	)

	if err != nil {
		status = "compression_error"
		http.Error(w, "Error compressing image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Record metrics
	metrics.RecordCompressionRatio(
		format,
		result.AlgorithmUsed,
		result.OriginalSize,
		result.CompressedSize,
	)

	// Set response headers
	w.Header().Set("Content-Type", fmt.Sprintf("image/%s", format))
	w.Header().Set("Content-Disposition", fmt.Sprintf(
		"attachment; filename=%s.%s",
		filepath.Base(header.Filename),
		format,
	))

	// Send the response
	_, err = w.Write(result.Data)
	if err != nil {
		log.Printf("Error writing response: %v", err)
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