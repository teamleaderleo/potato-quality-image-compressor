package api

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"
	"strings"
	"sync"
)

// BatchRequest represents a single image compression request in a batch
type BatchRequest struct {
	Filename  string
	Data      io.Reader
	Format    string
	Quality   int
	Algorithm string
}

// BatchResponse represents a batch processing response
type BatchResponse struct {
	Results          []CompressionResult
	ProcessingErrors []BatchProcessError
}

// BatchProcessError represents an error during batch processing
type BatchProcessError struct {
	Filename string
	Error    error
}

// ProcessBatchRequests processes multiple image compression requests concurrently
// This is the exported method that both HTTP and gRPC handlers will use
func (s *Service) ProcessBatchRequests(
	ctx context.Context,
	requests []BatchRequest,
) BatchResponse {
	var (
		results          []CompressionResult
		processingErrors []BatchProcessError
		resultsMutex     sync.Mutex
		errorsMutex      sync.Mutex
		wg               sync.WaitGroup
		// Limit concurrency to number of workers
		sem = make(chan struct{}, s.workerPool.TotalWorkerCount())
	)

	for _, req := range requests {
		wg.Add(1)

		go func(request BatchRequest) {
			defer wg.Done()

			// Acquire semaphore slot
			sem <- struct{}{}
			defer func() { <-sem }()

			// Check if context is cancelled
			if ctx.Err() != nil {
				errorsMutex.Lock()
				processingErrors = append(processingErrors, BatchProcessError{
					Filename: request.Filename,
					Error:    ctx.Err(),
				})
				errorsMutex.Unlock()
				return
			}

			// Create a new context with a timeout for each image
			fileCtx, fileCancel := context.WithTimeout(ctx, s.imageProcessingTimeout)
			defer fileCancel()

			// Process the image
			result, err := s.CompressImage(
				fileCtx,
				request.Filename,
				request.Data,
				request.Format,
				request.Quality,
				request.Algorithm,
			)

			if err != nil {
				errorsMutex.Lock()
				processingErrors = append(processingErrors, BatchProcessError{
					Filename: request.Filename,
					Error:    err,
				})
				errorsMutex.Unlock()
				return
			}

			resultsMutex.Lock()
			results = append(results, result)
			resultsMutex.Unlock()
		}(req)
	}

	// Wait for all processing to complete
	wg.Wait()

	return BatchResponse{
		Results:          results,
		ProcessingErrors: processingErrors,
	}
}

// ConvertFilesToBatchRequests converts multipart file headers to batch requests
// Helper function for HTTP handler
func ConvertFilesToBatchRequests(
	files []*multipart.FileHeader,
	format string,
	quality int,
	algorithm string,
) ([]BatchRequest, error) {
	requests := make([]BatchRequest, 0, len(files))

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("opening file %s: %w", fileHeader.Filename, err)
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("reading file %s: %w", fileHeader.Filename, err)
		}

		requests = append(requests, BatchRequest{
			Filename:  fileHeader.Filename,
			Data:      bytes.NewReader(fileBytes),
			Format:    format,
			Quality:   quality,
			Algorithm: algorithm,
		})
	}

	return requests, nil
}

// CreateZipFromResults creates a zip file from compression results
func CreateZipFromResults(results []CompressionResult) ([]byte, error) {
	if len(results) == 0 {
		return nil, errors.New("no valid results to create zip file")
	}

	// Create a buffer for the zip file
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Create a map to track name counts in case of duplicates
	nameCounts := make(map[string]int)

	// Add each result to the zip
	for _, result := range results {
		if result.Error != nil {
			continue // Skip failed results
		}

		// Generate a unique filename for the zip entry
		baseName := strings.TrimSuffix(filepath.Base(result.Filename), filepath.Ext(result.Filename))
		ext := result.Format

		count := nameCounts[baseName]
		nameCounts[baseName]++

		var zipFilename string
		if count == 0 {
			zipFilename = fmt.Sprintf("%s.%s", baseName, ext)
		} else {
			zipFilename = fmt.Sprintf("%s_%d.%s", baseName, count, ext)
		}

		// Create a zip file header
		zipHeader := &zip.FileHeader{
			Name:     zipFilename,
			Method:   zip.Deflate,
			Modified: time.Now(),
		}

		zipFile, err := zipWriter.CreateHeader(zipHeader)

		if err != nil {
			return nil, fmt.Errorf("creating zip entry: %w", err)
		}

		_, err = zipFile.Write(result.Data)
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