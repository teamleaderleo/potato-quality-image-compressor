package api

import (
	"bytes"
	"fmt"
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/compression"
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/metrics"
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/worker"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"time"
)

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
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()

			var m runtime.MemStats
			for range ticker.C {
				runtime.ReadMemStats(&m)
				metrics.UpdateMemoryUsage(m.Alloc)
			}
		}()
	}

	return &Service{
		workerPool:       workerPool,
		processor:        processor,
		defaultQuality:   config.DefaultQuality,
		defaultFormat:    config.DefaultFormat,
		defaultAlgorithm: config.DefaultAlgorithm,
	}
}

// HandleCompress handles single image compression requests
func (s *Service) HandleCompress(w http.ResponseWriter, r *http.Request) {
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

	// Process the image
	result, err := s.CompressImageDirect(
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

// Shutdown gracefully shuts down the service
func (s *Service) Shutdown() {
	if s.workerPool != nil {
		s.workerPool.Shutdown()
	}
}