package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "image_compression_requests_total",
			Help: "Total number of image compression requests",
		},
		[]string{"endpoint", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "image_compression_duration_seconds",
			Help:    "Time taken to process image compression request",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)

	jobDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "image_compression_job_duration_seconds",
			Help:    "Time taken to process a single image compression job",
			Buckets: prometheus.DefBuckets,
		},
	)

	compressionRatio = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "image_compression_ratio",
			Help:    "Ratio of compressed image size to original image size",
			Buckets: []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0},
		},
		[]string{"format", "algorithm"},
	)

	workerGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "image_compression_busy_workers",
			Help: "Number of busy workers in the worker pool",
		},
	)

	memoryUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "image_compression_memory_bytes",
			Help: "Current memory usage of the image compression service",
		},
	)
)

// Init registers all metrics and sets up the HTTP handler
func Init() error {
	if err := prometheus.Register(requestCounter); err != nil {
		return fmt.Errorf("failed to register request counter: %w", err)
	}
	if err := prometheus.Register(requestDuration); err != nil {
		return fmt.Errorf("failed to register request duration: %w", err)
	}
	if err := prometheus.Register(jobDuration); err != nil {
		return fmt.Errorf("failed to register job duration: %w", err)
	}
	if err := prometheus.Register(compressionRatio); err != nil {
		return fmt.Errorf("failed to register compression ratio: %w", err)
	}
	if err := prometheus.Register(workerGauge); err != nil {
		return fmt.Errorf("failed to register worker gauge: %w", err)
	}
	if err := prometheus.Register(memoryUsage); err != nil {
		return fmt.Errorf("failed to register memory usage: %w", err)
	}

	http.Handle("/metrics", promhttp.Handler())
	return nil
}

// GetRequestCounter returns the request counter metric
func GetRequestCounter() *prometheus.CounterVec {
	return requestCounter
}

// GetRequestDuration returns the request duration metric
func GetRequestDuration() *prometheus.HistogramVec {
	return requestDuration
}

// GetJobDuration returns the job duration metric
func GetJobDuration() *prometheus.Histogram {
	return &jobDuration
}

// GetCompressionRatio returns the compression ratio metric
func GetCompressionRatio() *prometheus.HistogramVec {
	return compressionRatio
}

// GetWorkerGauge returns the worker gauge metric
func GetWorkerGauge() *prometheus.Gauge {
	return &workerGauge
}

// GetMemoryUsage returns the memory usage metric
func GetMemoryUsage() *prometheus.Gauge {
	return &memoryUsage
}

// NewTimer creates a new timer for measuring request duration
func NewTimer(endpoint string) *prometheus.Timer {
	return prometheus.NewTimer(requestDuration.WithLabelValues(endpoint))
}

// RecordCompressionRatio records the compression ratio metric
func RecordCompressionRatio(format, algorithm string, originalSize, compressedSize int) {
	if originalSize > 0 {
		ratio := float64(compressedSize) / float64(originalSize)
		compressionRatio.WithLabelValues(format, algorithm).Observe(ratio)
	}
}

// UpdateMemoryUsage updates the memory usage metric (called periodically)
func UpdateMemoryUsage(bytesUsed uint64) {
	memoryUsage.Set(float64(bytesUsed))
}