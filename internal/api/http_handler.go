package api

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	
	"strconv"
	"sync"
	
    "io"

	
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/metrics"
	
)

// BatchResult represents a result from batch image processing
type BatchResult struct {
	filename string
	result   CompressionResult
	err      error
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

	// Create channel for results
	resultChan := make(chan BatchResult, len(files))
	
	var wg sync.WaitGroup
	for _, fileHeader := range files {
		wg.Add(1)
		
		go func(fh *multipart.FileHeader) {
			defer wg.Done()
			
			file, err := fh.Open()
			if err != nil {
				resultChan <- BatchResult{
					filename: fh.Filename,
					err:      err,
				}
				return
			}
			defer file.Close()

			fileBytes, err := io.ReadAll(file)
			if err != nil {
				resultChan <- BatchResult{
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
			
			resultChan <- BatchResult{
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