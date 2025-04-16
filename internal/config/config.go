package config

import (
	"os"
	"runtime"
	"strconv"
	"time"
)

// AppConfig represents the application configuration
type AppConfig struct {
	Server        ServerConfig
	Compression   CompressionConfig
	Worker        WorkerConfig
	Metrics       MetricsConfig
	HttpEnabled   bool
	GrpcEnabled   bool
	GrpcPort      string
	ShutdownDelay time.Duration
}

// ServerConfig represents HTTP server configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// CompressionConfig represents image compression configuration
type CompressionConfig struct {
	DefaultQuality         int
	DefaultFormat          string
	DefaultAlgorithm       string
	MaxUploadSize          int64
	MaxBatchSize           int
	ImageProcessingTimeout time.Duration
	BatchProcessingTimeout time.Duration
}

// WorkerConfig represents worker pool configuration
type WorkerConfig struct {
	WorkerCount  int
	JobQueueSize int
}

// MetricsConfig represents metrics configuration
type MetricsConfig struct {
	Enabled            bool
	UpdateInterval     time.Duration
	MetricsEndpoint    string
	PrometheusEnabled  bool
}

// ServiceConfig contains configuration for the compression service
type ServiceConfig struct {
	WorkerCount            int
	JobQueueSize           int
	DefaultQuality         int
	DefaultFormat          string
	DefaultAlgorithm       string
	EnableMetrics          bool
	ImageProcessingTimeout time.Duration
	BatchProcessingTimeout time.Duration
	MaxUploadSize          int64
	MaxBatchSize           int
}

// CreateServiceConfig creates a ServiceConfig from AppConfig
func (c AppConfig) CreateServiceConfig() ServiceConfig {
	return ServiceConfig{
		WorkerCount:            c.Worker.WorkerCount,
		JobQueueSize:           c.Worker.JobQueueSize,
		DefaultQuality:         c.Compression.DefaultQuality,
		DefaultFormat:          c.Compression.DefaultFormat,
		DefaultAlgorithm:       c.Compression.DefaultAlgorithm,
		EnableMetrics:          c.Metrics.Enabled,
		ImageProcessingTimeout: c.Compression.ImageProcessingTimeout,
		BatchProcessingTimeout: c.Compression.BatchProcessingTimeout,
		MaxUploadSize:          c.Compression.MaxUploadSize,
		MaxBatchSize:           c.Compression.MaxBatchSize,
	}
}

// LoadConfig loads the application configuration from environment variables
func LoadConfig() AppConfig {
	return AppConfig{
		Server: ServerConfig{
			Port:         getEnvWithDefault("PORT", "8080"),
			ReadTimeout:  getDurationWithDefault("READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationWithDefault("WRITE_TIMEOUT", 60*time.Second),
			IdleTimeout:  getDurationWithDefault("IDLE_TIMEOUT", 120*time.Second),
		},
		Compression: CompressionConfig{
			DefaultQuality:         getIntWithDefault("DEFAULT_QUALITY", 75),
			DefaultFormat:          getEnvWithDefault("DEFAULT_FORMAT", "webp"),
			DefaultAlgorithm:       getEnvWithDefault("DEFAULT_ALGORITHM", "scale"),
			MaxUploadSize:          getInt64WithDefault("MAX_UPLOAD_SIZE", 32<<20), // 32 MB
			MaxBatchSize:           getIntWithDefault("MAX_BATCH_SIZE", 50),
			ImageProcessingTimeout: getDurationWithDefault("IMAGE_PROCESSING_TIMEOUT", 30*time.Second),
			BatchProcessingTimeout: getDurationWithDefault("BATCH_PROCESSING_TIMEOUT", 5*time.Minute),
		},
		Worker: WorkerConfig{
			WorkerCount:  getIntWithDefault("WORKER_COUNT", runtime.NumCPU()),
			JobQueueSize: getIntWithDefault("JOB_QUEUE_SIZE", runtime.NumCPU()*4),
		},
		Metrics: MetricsConfig{
			Enabled:           getBoolWithDefault("METRICS_ENABLED", true),
			UpdateInterval:    getDurationWithDefault("METRICS_UPDATE_INTERVAL", 10*time.Second),
			MetricsEndpoint:   getEnvWithDefault("METRICS_ENDPOINT", "/metrics"),
			PrometheusEnabled: getBoolWithDefault("PROMETHEUS_ENABLED", true),
		},
		HttpEnabled:   getBoolWithDefault("HTTP_ENABLED", true),
		GrpcEnabled:   getBoolWithDefault("GRPC_ENABLED", false),
		GrpcPort:      getEnvWithDefault("GRPC_PORT", "9090"),
		ShutdownDelay: getDurationWithDefault("SHUTDOWN_DELAY", 30*time.Second),
	}
}

// Helper functions to get environment variables with defaults

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntWithDefault(key string, defaultValue int) int {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}
	
	value, err := strconv.Atoi(strValue)
	if err != nil {
		return defaultValue
	}
	
	return value
}

func getInt64WithDefault(key string, defaultValue int64) int64 {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}
	
	value, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		return defaultValue
	}
	
	return value
}

func getBoolWithDefault(key string, defaultValue bool) bool {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}
	
	value, err := strconv.ParseBool(strValue)
	if err != nil {
		return defaultValue
	}
	
	return value
}

func getDurationWithDefault(key string, defaultValue time.Duration) time.Duration {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}
	
	value, err := time.ParseDuration(strValue)
	if err != nil {
		return defaultValue
	}
	
	return value
}