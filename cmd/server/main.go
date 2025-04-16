package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/api"
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/config"
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/grpc"
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/metrics"
	"google.golang.org/grpc/reflection"
	grpcServer "google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

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

	// Initialize metrics if enabled
	if cfg.Metrics.Enabled {
		if err := metrics.Init(); err != nil {
			log.Fatalf("Failed to initialize metrics: %v", err)
		}
	}

	// Create service with configuration
	serviceConfig := cfg.CreateServiceConfig()
	service := api.NewServiceWithConfig(serviceConfig)
	defer service.Shutdown()

	// Determine running mode
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Lambda mode
		api.StartLambda()
		return
	}

	// Track if any server is started
	serversStarted := false
	
	// Setup and start HTTP server if enabled
	var httpServer *http.Server
	if cfg.HttpEnabled {
		httpServer = setupHTTPServer(service, cfg)
		go startHTTPServer(httpServer)
		serversStarted = true
	}

	// Setup and start gRPC server if enabled
	var grpcSrv *grpcServer.Server
	if cfg.GrpcEnabled {
		grpcSrv = setupGRPCServer(service, cfg)
		go startGRPCServer(grpcSrv, cfg.GrpcPort)
		serversStarted = true
	}
	
	// If no servers were started, log a warning
	if !serversStarted {
		log.Println("Warning: No servers were started. Application will only process background tasks.")
	}

	// Wait for context cancellation (shutdown signal)
	<-ctx.Done()
	
	// Graceful shutdown
	log.Println("Shutting down...")
	
	// Shutdown HTTP server if it was started
	if cfg.HttpEnabled && httpServer != nil {
		shutdownHTTPServer(httpServer, cfg.ShutdownDelay)
	}
	
	// Shutdown gRPC server if it was started
	if cfg.GrpcEnabled && grpcSrv != nil {
		shutdownGRPCServer(grpcSrv)
	}
	
	log.Println("Shutdown completed")
}

// setupHTTPServer creates and configures the HTTP server
func setupHTTPServer(service *api.Service, cfg config.AppConfig) *http.Server {
	// Set up HTTP server and routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/compress", service.HandleCompress)
	mux.HandleFunc("/batch-compress", service.HandleBatchCompress)
	
	// Add Prometheus metrics endpoint if enabled
	if cfg.Metrics.PrometheusEnabled {
		mux.Handle(cfg.Metrics.MetricsEndpoint, promhttp.Handler())
	}

	// Create server with configured timeouts
	return &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
}

// startHTTPServer starts the HTTP server and logs any error
func startHTTPServer(server *http.Server) {
	log.Printf("HTTP server starting on http://localhost%s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("HTTP server failed to start: %v", err)
	}
}

// setupGRPCServer creates and configures the gRPC server
func setupGRPCServer(service *api.Service, cfg config.AppConfig) *grpcServer.Server {
	// Create new gRPC server
	grpcSrv := grpcServer.NewServer()
	
	// Register services
	grpc.RegisterServer(grpcSrv, service)
	
	// Enable reflection for debugging
	reflection.Register(grpcSrv)
	
	return grpcSrv
}

// startGRPCServer starts the gRPC server
func startGRPCServer(server *grpcServer.Server, port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}
	
	log.Printf("gRPC server starting on port %s", port)
	if err := server.Serve(listener); err != nil {
		log.Printf("gRPC server error: %v", err)
	}
}

// shutdownHTTPServer gracefully shuts down the HTTP server
func shutdownHTTPServer(server *http.Server, timeout time.Duration) {
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), timeout)
	defer shutdownCancel()
	
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	} else {
		log.Println("HTTP server gracefully stopped")
	}
}

// shutdownGRPCServer gracefully shuts down the gRPC server
func shutdownGRPCServer(server *grpcServer.Server) {
	log.Println("Stopping gRPC server...")
	server.GracefulStop()
	log.Println("gRPC server gracefully stopped")
}

// handleRoot handles the root endpoint
func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Image Compression API")
}