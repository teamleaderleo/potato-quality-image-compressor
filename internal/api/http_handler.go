package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/teamleaderleo/potato-quality-image-compressor/internal/metrics"
)

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