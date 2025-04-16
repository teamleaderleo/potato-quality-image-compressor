package api

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/teamleaderleo/potato-quality-image-compressor/internal/metrics"
)

const (
	maxUploadSize = 32 << 20 // 32 MB
	maxBatchSize  = 50       // Maximum number of files in batch
)

// BatchProcessError represents an error during batch processing
type BatchProcessError struct {
	Filename string
	Error    error
}

// BatchResult represents a result from batch image processing
type BatchResult struct {
	filename string
	result   CompressionResult
}

// HandleBatchCompress handles batch image compression requests
func (s *Service) HandleBatchCompress(w http.ResponseWriter, r *http.Request) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

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

	// Limit the max upload size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	err := r.ParseMultipartForm(maxUploadSize)
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

	if len(files) > maxBatchSize {
		status = "bad_request"
		http.Error(w, fmt.Sprintf("Too many files. Maximum is %d", maxBatchSize), http.StatusBadRequest)
		return
	}

	// Parse parameters
	quality, format, algorithm := s.parseParameters(r)

	// Process the batch of images
	results, processingErrors := s.processBatchImages(ctx, files, format, quality, algorithm)

	// If all files failed, return an error
	if len(results) == 0 && len(processingErrors) > 0 {
		status = "batch_processing_failed"
		errMsg := fmt.Sprintf("Failed to process all %d images", len(processingErrors))
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// Create zip file from results
	zipData, err := createZipFromResults(results, format)
	if err != nil {
		status = "zip_error"
		http.Error(w, "Error creating zip file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Record metrics for all successful compressions
	for _, result := range results {
		metrics.RecordCompressionRatio(
			format,
			result.result.AlgorithmUsed,
			result.result.OriginalSize,
			result.result.CompressedSize,
		)
	}

	// Send the response
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=compressed_images.zip")
	_, err = w.Write(zipData)
	if err != nil {
		log.Printf("Error writing zip response: %v", err)
	}

	// Log processing errors for debugging
	if len(processingErrors) > 0 {
		for _, procErr := range processingErrors {
			log.Printf("Error processing %s: %v", procErr.Filename, procErr.Error)
		}
	}
}

// processBatchImages processes multiple images concurrently
func (s *Service) processBatchImages(
	ctx context.Context,
	files []*multipart.FileHeader,
	format string,
	quality int,
	algorithm string,
) ([]BatchResult, []BatchProcessError) {
	var (
		results          []BatchResult
		processingErrors []BatchProcessError
		resultsMutex     sync.Mutex
		errorsMutex      sync.Mutex
		wg               sync.WaitGroup
		// Limit concurrency to number of workers
		sem = make(chan struct{}, s.workerPool.TotalWorkerCount())
	)

	for _, fileHeader := range files {
		wg.Add(1)

		go func(fh *multipart.FileHeader) {
			defer wg.Done()

			// Acquire semaphore slot
			sem <- struct{}{}
			defer func() { <-sem }()

			// Check if context is cancelled
			if ctx.Err() != nil {
				errorsMutex.Lock()
				processingErrors = append(processingErrors, BatchProcessError{
					Filename: fh.Filename,
					Error:    ctx.Err(),
				})
				errorsMutex.Unlock()
				return
			}

			file, err := fh.Open()
			if err != nil {
				errorsMutex.Lock()
				processingErrors = append(processingErrors, BatchProcessError{
					Filename: fh.Filename,
					Error:    err,
				})
				errorsMutex.Unlock()
				return
			}
			defer file.Close()

			fileBytes, err := io.ReadAll(file)
			if err != nil {
				errorsMutex.Lock()
				processingErrors = append(processingErrors, BatchProcessError{
					Filename: fh.Filename,
					Error:    err,
				})
				errorsMutex.Unlock()
				return
			}

			// Create a new context with a timeout for each image
			fileCtx, fileCancel := context.WithTimeout(ctx, 30*time.Second)
			defer fileCancel()

			// Process the image
			result, err := s.CompressImage(
				fileCtx,
				fh.Filename,
				bytes.NewReader(fileBytes),
				format,
				quality,
				algorithm,
			)

			if err != nil {
				errorsMutex.Lock()
				processingErrors = append(processingErrors, BatchProcessError{
					Filename: fh.Filename,
					Error:    err,
				})
				errorsMutex.Unlock()
				return
			}

			resultsMutex.Lock()
			results = append(results, BatchResult{
				filename: fh.Filename,
				result:   result,
			})
			resultsMutex.Unlock()
		}(fileHeader)
	}

	// Wait for all processing to complete
	wg.Wait()

	return results, processingErrors
}

// createZipFromResults creates a zip file from batch processing results
func createZipFromResults(results []BatchResult, format string) ([]byte, error) {
	if len(results) == 0 {
		return nil, errors.New("no valid results to create zip file")
	}

	// Create a buffer for the zip file
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Add each result to the zip
	for _, r := range results {
		if r.result.Error != nil {
			continue // Skip failed results
		}

		zipFile, err := zipWriter.Create(fmt.Sprintf("%s.%s",
			filepath.Base(r.filename),
			format,
		))
		if err != nil {
			return nil, fmt.Errorf("creating zip entry: %w", err)
		}

		_, err = zipFile.Write(r.result.Data)
		if err != nil {
			return nil, fmt.Errorf("writing to zip: %w", err)
		}
	}

	// Close the zip writer
	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("closing zip writer: %w", err)
	}

	return buf.Bytes(), nil
}
