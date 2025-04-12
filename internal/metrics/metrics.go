package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
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
)

func Init() error {
    if err := prometheus.Register(requestCounter); err != nil {
        return fmt.Errorf("failed to register request counter: %w", err)
    }
    if err := prometheus.Register(requestDuration); err != nil {
        return fmt.Errorf("failed to register request duration: %w", err)
    }
    
    http.Handle("/metrics", promhttp.Handler())
    return nil
}

func GetRequestCounter() *prometheus.CounterVec {
	return requestCounter
}

func GetRequestDuration() *prometheus.HistogramVec {
	return requestDuration
}

func NewTimer(endpoint string) *prometheus.Timer {
    return prometheus.NewTimer(requestDuration.WithLabelValues(endpoint))
}