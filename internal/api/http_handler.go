package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"net/http"
	"time"

	"github.com/teamleaderleo/potato-quality-image-compressor/internal/metrics"
)

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

// HandleBatchCompress handles batch image compression requests via HTTP
func (s *Service) HandleBatchCompress(w http.ResponseWriter, r *http.Request) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), s.batchProcessingTimeout)
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
	r.Body = http.MaxBytesReader(w, r.Body, s.maxUploadSize)

	err := r.ParseMultipartForm(s.maxUploadSize)
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

	if len(files) > s.maxBatchSize {
		status = "bad_request"
		http.Error(w, fmt.Sprintf("Too many files. Maximum is %d", s.maxBatchSize), http.StatusBadRequest)
		return
	}

	// Parse parameters
	quality, format, algorithm := s.parseParameters(r)

	// Convert multipart files to batch requests
	requests, err := ConvertFilesToBatchRequests(files, format, quality, algorithm)
	if err != nil {
		status = "batch_preparation_failed"
		http.Error(w, "Error preparing batch: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Process the batch of images using the unified processor
	batchResponse := s.ProcessBatchRequests(ctx, requests)
	
	// If all files failed, return an error
	if len(batchResponse.Results) == 0 && len(batchResponse.ProcessingErrors) > 0 {
		status = "batch_processing_failed"
		errMsg := fmt.Sprintf("Failed to process all %d images", len(batchResponse.ProcessingErrors))
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// Create zip file from results
	zipData, err := CreateZipFromResults(batchResponse.Results)
	if err != nil {
		status = "zip_error"
		http.Error(w, "Error creating zip file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Record metrics for all successful compressions
	for _, result := range batchResponse.Results {
		metrics.RecordCompressionRatio(
			result.Format,
			result.AlgorithmUsed,
			result.OriginalSize,
			result.CompressedSize,
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
	if len(batchResponse.ProcessingErrors) > 0 {
		for _, procErr := range batchResponse.ProcessingErrors {
			log.Printf("Error processing %s: %v", procErr.Filename, procErr.Error)
		}
	}
}