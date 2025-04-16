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

	// Process metrics
	memoryUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "image_compression_memory_bytes",
			Help: "Current memory usage of the image compression service process",
		},
	)

	cpuUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "image_compression_cpu_percent",
			Help: "Current CPU usage percentage of the image compression service process",
		},
	)

	// System metrics
	systemMemoryUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "system_memory_bytes_used",
			Help: "Current system memory usage in bytes",
		},
	)

	systemMemoryPercent = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "system_memory_percent_used",
			Help: "Current system memory usage percentage",
		},
	)

	systemCPUUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "system_cpu_percent",
			Help: "Current system CPU usage percentage",
		},
	)

	// Throughput metrics
	throughputImages = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "image_compression_images_processed_total",
			Help: "Total number of images processed",
		},
	)

	throughputBytes = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "image_compression_bytes_processed_total",
			Help: "Total number of bytes processed",
		},
	)
)

// Init registers all metrics and sets up the HTTP handler
func Init() error {
	// Request metrics
	if err := prometheus.Register(requestCounter); err != nil {
		return fmt.Errorf("failed to register request counter: %w", err)
	}
	if err := prometheus.Register(requestDuration); err != nil {
		return fmt.Errorf("failed to register request duration: %w", err)
	}
	if err := prometheus.Register(jobDuration); err != nil {
		return fmt.Errorf("failed to register job duration: %w", err)
	}
	
	// Performance metrics
	if err := prometheus.Register(compressionRatio); err != nil {
		return fmt.Errorf("failed to register compression ratio: %w", err)
	}
	if err := prometheus.Register(workerGauge); err != nil {
		return fmt.Errorf("failed to register worker gauge: %w", err)
	}
	
	// Process resource metrics
	if err := prometheus.Register(memoryUsage); err != nil {
		return fmt.Errorf("failed to register memory usage: %w", err)
	}
	if err := prometheus.Register(cpuUsage); err != nil {
		return fmt.Errorf("failed to register CPU usage: %w", err)
	}
	
	// System resource metrics
	if err := prometheus.Register(systemMemoryUsage); err != nil {
		return fmt.Errorf("failed to register system memory usage: %w", err)
	}
	if err := prometheus.Register(systemMemoryPercent); err != nil {
		return fmt.Errorf("failed to register system memory percentage: %w", err)
	}
	if err := prometheus.Register(systemCPUUsage); err != nil {
		return fmt.Errorf("failed to register system CPU usage: %w", err)
	}
	
	// Throughput metrics
	if err := prometheus.Register(throughputImages); err != nil {
		return fmt.Errorf("failed to register image throughput: %w", err)
	}
	if err := prometheus.Register(throughputBytes); err != nil {
		return fmt.Errorf("failed to register byte throughput: %w", err)
	}

	http.Handle("/metrics", promhttp.Handler())
	return nil
}

// Timer functions

// NewTimer creates a new timer for measuring request duration
func NewTimer(endpoint string) *prometheus.Timer {
	return prometheus.NewTimer(requestDuration.WithLabelValues(endpoint))
}

// Process metrics update functions

// UpdateMemoryUsage updates the process memory usage metric
func UpdateMemoryUsage(bytesUsed uint64) {
	memoryUsage.Set(float64(bytesUsed))
}

// UpdateCPUUsage updates the process CPU usage metric
func UpdateCPUUsage(percentUsed float64) {
	cpuUsage.Set(percentUsed)
}

// System metrics update functions

// UpdateSystemMemoryUsage updates the system memory usage metric
func UpdateSystemMemoryUsage(bytesUsed uint64) {
	systemMemoryUsage.Set(float64(bytesUsed))
}

// UpdateSystemMemoryPercent updates the system memory percentage metric
func UpdateSystemMemoryPercent(percentUsed float64) {
	systemMemoryPercent.Set(percentUsed)
}

// UpdateSystemCPUUsage updates the system CPU usage metric
func UpdateSystemCPUUsage(percentUsed float64) {
	systemCPUUsage.Set(percentUsed)
}

// Performance metrics

// RecordCompressionRatio records the compression ratio metric
func RecordCompressionRatio(format, algorithm string, originalSize, compressedSize int) {
	if originalSize > 0 {
		ratio := float64(compressedSize) / float64(originalSize)
		compressionRatio.WithLabelValues(format, algorithm).Observe(ratio)
		
		// Update throughput metrics
		throughputImages.Inc()
		throughputBytes.Add(float64(originalSize))
	}
}

// Getter functions

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

// GetCPUUsage returns the CPU usage metric
func GetCPUUsage() *prometheus.Gauge {
	return &cpuUsage
}