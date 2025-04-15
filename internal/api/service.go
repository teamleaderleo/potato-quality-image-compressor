package api

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
    "io"

	"github.com/teamleaderleo/potato-quality-image-compressor/internal/compression"
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/metrics"
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/worker"
)

// Service handles the API endpoints for image compression
type Service struct {
	workerPool     *worker.Pool
	processor      *compression.ImageProcessor
	defaultQuality int
	defaultFormat  string
	defaultAlgorithm string
}

// ServiceConfig contains configuration for the Service
type ServiceConfig struct {
	WorkerCount     int
	JobQueueSize    int
	DefaultQuality  int
	DefaultFormat   string
	DefaultAlgorithm string
	EnableMetrics   bool
}

// DefaultServiceConfig returns the default service configuration
func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{
		WorkerCount:     runtime.NumCPU(),
		JobQueueSize:    runtime.NumCPU() * 4,
		DefaultQuality:  75,
		DefaultFormat:   "webp",
		DefaultAlgorithm: "scale",
		EnableMetrics:   true,
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
		workerPool:      workerPool,
		processor:       processor,
		defaultQuality:  config.DefaultQuality,
		defaultFormat:   config.DefaultFormat,
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

// HandleBatchCompress handles batch image compression requests
func (s *Service) HandleBatchCompress(w http.ResponseWriter, r *http.Request) {
	timer := metrics.NewTimer("batch-compress")
	defer timer.ObserveDuration()

	status := "success"
	defer func() {
		metrics.GetRequestCounter().WithLabelValues("batch-compress", status).Inc()
	}()

	if r.Method != http.MethodPost {
		status = "method_not_allowed"
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // 32 MB max
	if err != nil {
		status = "bad_request"
		http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		status = "bad_request"
		http.Error(w, "No images provided", http.StatusBadRequest)
		return
	}

	// Parse parameters
	quality, format, algorithm := s.parseParameters(r)

	// Create a buffer for the zip file
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Process each file concurrently
	resultChan := make(chan struct {
		filename string
		result   CompressionResult
		err      error
	}, len(files))
	
	var wg sync.WaitGroup
	for _, fileHeader := range files {
		wg.Add(1)
		
		go func(fh *multipart.FileHeader) {
			defer wg.Done()
			
			file, err := fh.Open()
			if err != nil {
				resultChan <- struct {
					filename string
					result   CompressionResult
					err      error
				}{
					filename: fh.Filename,
					err:      err,
				}
				return
			}
			defer file.Close()

			fileBytes, err := io.ReadAll(file)
			if err != nil {
				resultChan <- struct {
					filename string
					result   CompressionResult
					err      error
				}{
					filename: fh.Filename,
					err:      err,
				}
				return
			}

			// Process the image
			result, err := s.CompressImageDirect(
				fh.Filename,
				bytes.NewReader(fileBytes),
				format,
				quality,
				algorithm,
			)
			
			resultChan <- struct {
				filename string
				result   CompressionResult
				err      error
			}{
				filename: fh.Filename,
				result:   result,
				err:      err,
			}
		}(fileHeader)
	}
	
	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results and add to zip
	for r := range resultChan {
		if r.err != nil {
			log.Printf("Error processing %s: %v", r.filename, r.err)
			continue
		}
		
		result := r.result
		if result.Error != nil {
			log.Printf("Error compressing %s: %v", r.filename, result.Error)
			continue
		}

		zipFile, err := zipWriter.Create(fmt.Sprintf("%s.%s", 
			filepath.Base(r.filename), 
			format,
		))
		if err != nil {
			log.Printf("Error creating zip entry: %v", err)
			continue
		}
		
		_, err = zipFile.Write(result.Data)
		if err != nil {
			log.Printf("Error writing to zip: %v", err)
			continue
		}
		
		// Record metrics
		metrics.RecordCompressionRatio(
			format,
			result.AlgorithmUsed,
			result.OriginalSize,
			result.CompressedSize,
		)
	}

	// Close the zip writer
	err = zipWriter.Close()
	if err != nil {
		status = "zip_error"
		http.Error(w, "Error creating zip file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=compressed_images.zip")
	_, err = w.Write(buf.Bytes())
	if err != nil {
		log.Printf("Error writing zip response: %v", err)
	}
}

// parseParameters parses and validates request parameters
func (s *Service) parseParameters(r *http.Request) (int, string, string) {
	// Parse quality
	quality, err := strconv.Atoi(r.FormValue("quality"))
	if err != nil || quality < 1 || quality > 100 {
		quality = s.defaultQuality
	}

	// Parse format
	format := r.FormValue("format")
	if format == "" {
		format = s.defaultFormat
	}

	// Parse algorithm
	algorithm := r.FormValue("algorithm")
	if algorithm == "" {
		algorithm = s.defaultAlgorithm
	}

	return quality, format, algorithm
}

// Shutdown gracefully shuts down the service
func (s *Service) Shutdown() {
	if s.workerPool != nil {
		s.workerPool.Shutdown()
	}
}