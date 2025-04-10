package main

import (
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

func InitMetrics() {
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
	
	http.Handle("/metrics", promhttp.Handler())
}