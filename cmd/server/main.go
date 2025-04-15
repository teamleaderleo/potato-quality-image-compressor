package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/teamleaderleo/potato-quality-image-compressor/internal/api"
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/metrics"
)

func main() {
	// Initialize context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		sig := <-signalChan
		log.Printf("Received signal: %v, initiating shutdown", sig)
		cancel()
	}()

	// Initialize metrics
	if err := metrics.Init(); err != nil {
		log.Fatalf("Failed to initialize metrics: %v", err)
	}

	// Create service with default configuration
	service := api.NewService()
	defer service.Shutdown()

	// Set up HTTP server and routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/compress", service.HandleCompress)
	mux.HandleFunc("/batch-compress", service.HandleBatchCompress)

	// Determine running mode
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Lambda mode
		api.StartLambda()
	} else {
		// Server mode
		port := getEnvWithDefault("PORT", "8080")
		server := createServer(port, mux)
		
		// Start server in a goroutine
		go startServer(server)

		// Wait for context cancellation (shutdown signal)
		<-ctx.Done()
		
		// Graceful shutdown with 30s timeout
		log.Println("Shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()
		
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
		
		log.Println("Server gracefully stopped")
	}
}

// createServer creates an HTTP server with configured timeouts
func createServer(port string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

// startServer starts the HTTP server and logs any error
func startServer(server *http.Server) {
	log.Printf("Server starting on http://localhost%s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// getEnvWithDefault returns an environment variable or the default value if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// handleRoot handles the root endpoint
func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Image Compression API")
}