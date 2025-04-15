package compression

import (
	"bytes"
	"io"
	"time"
	
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/worker"
)

// CompressionJob represents a job to compress an image
type CompressionJob struct {
	id        string
	input     io.Reader
	format    string
	quality   int
	algorithm string
	processor *ImageProcessor
}

// NewCompressionJob creates a new compression job
func NewCompressionJob(id string, input io.Reader, format string, quality int, algorithm string, processor *ImageProcessor) *CompressionJob {
	return &CompressionJob{
		id:        id,
		input:     input,
		format:    format,
		quality:   quality,
		algorithm: algorithm,
		processor: processor,
	}
}

// ID returns the job identifier
func (j *CompressionJob) ID() string {
	return j.id
}

// Process executes the compression job
func (j *CompressionJob) Process() (worker.JobResult, error) {
	startTime := time.Now()
	
	// Read the entire input
	inputData, err := io.ReadAll(j.input)
	if err != nil {
		return nil, err
	}
	
	// Get the algorithm to use
	var algorithm CompressionAlgorithm
	if a, ok := j.processor.GetAlgorithm(j.algorithm); ok {
		algorithm = a
	} else {
		// Use default algorithm if requested one isn't available
		algorithm = j.processor.GetDefaultAlgorithm()
	}
	
	// Process the image
	data, err := j.processor.ProcessImage(bytes.NewReader(inputData), j.format, j.quality, algorithm)
	if err != nil {
		return nil, err
	}
	
	// Create the result
	result := &CompressionResult{
		id:           j.id,
		data:         data,
		jobTime:      time.Since(startTime),
		algorithmUsed: algorithm.Name(),
		originalSize:  len(inputData),
		compressedSize: len(data),
	}
	
	return result, nil
}

// CompressionResult represents the result of a compression job
type CompressionResult struct {
	id            string
	data          []byte
	jobTime       time.Duration
	algorithmUsed string
	originalSize  int
	compressedSize int
}

// ID returns the job identifier
func (r *CompressionResult) ID() string {
	return r.id
}

// Data returns the compressed image data
func (r *CompressionResult) Data() []byte {
	return r.data
}

// JobTime returns the time taken to process the job
func (r *CompressionResult) JobTime() time.Duration {
	return r.jobTime
}

// AlgorithmUsed returns the name of the algorithm used
func (r *CompressionResult) AlgorithmUsed() string {
	return r.algorithmUsed
}

// OriginalSize returns the original image size in bytes
func (r *CompressionResult) OriginalSize() int {
	return r.originalSize
}

// CompressedSize returns the compressed image size in bytes
func (r *CompressionResult) CompressedSize() int {
	return r.compressedSize
}

// CompressionRatio returns the compression ratio (compressed/original)
func (r *CompressionResult) CompressionRatio() float64 {
	if r.originalSize == 0 {
		return 0
	}
	return float64(r.compressedSize) / float64(r.originalSize)
}