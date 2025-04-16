package grpc

import (
	"bytes"
	"context"
	"io"
	"runtime"
	"time"

	"github.com/teamleaderleo/potato-quality-image-compressor/internal/api"
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/metrics"
	pb "github.com/teamleaderleo/potato-quality-image-compressor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Adapter implements the gRPC server for image compression
// by reusing the existing service implementation
type Adapter struct {
	pb.UnimplementedImageCompressionServiceServer
	service *api.Service
}

// NewAdapter creates a new gRPC adapter with the given service
func NewAdapter(service *api.Service) *Adapter {
	return &Adapter{
		service: service,
	}
}

// RegisterServer registers the adapter with a gRPC server
func RegisterServer(grpcServer *grpc.Server, service *api.Service) {
	// Create adapter
	adapter := NewAdapter(service)
	
	// Register services
	pb.RegisterImageCompressionServiceServer(grpcServer, adapter)
	
	// Register health service
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, healthServer)
}

// CompressImage handles gRPC compression requests by delegating to the service
func (a *Adapter) CompressImage(ctx context.Context, req *pb.CompressImageRequest) (*pb.CompressImageResponse, error) {
	timer := metrics.NewTimer("grpc-compress")
	defer timer.ObserveDuration()

	status := "success"
	defer func() {
		metrics.GetRequestCounter().WithLabelValues("grpc-compress", status).Inc()
	}()
	
	// Create an HTTP-like structure to reuse the service implementation
	imgData := bytes.NewReader(req.ImageData)
	
	// Process the image using our existing service logic
	result, err := a.service.CompressImage(
		ctx,
		req.Filename,
		imgData,
		string(req.Format),
		int(req.Quality),
		string(req.Strategy),
	)
	
	if err != nil {
		return &pb.CompressImageResponse{
			Error: err.Error(),
		}, nil
	}
	
	// Convert result to protobuf response
	return &pb.CompressImageResponse{
		ImageData:        result.Data,
		Format:           req.Format,
		OriginalSize:     int64(result.OriginalSize),
		CompressedSize:   int64(result.CompressedSize),
		CompressionRatio: result.CompressionRatio,
		ProcessingTimeMs: result.ProcessingTime.Milliseconds(),
		Filename:         req.Filename,
	}, nil
}

// BatchCompressImages handles multiple image compression requests
func (a *Adapter) BatchCompressImages(ctx context.Context, req *pb.BatchCompressRequest) (*pb.BatchCompressResponse, error) {
	timer := metrics.NewTimer("grpc-batch-compress")
	defer timer.ObserveDuration()

	status := "success"
	defer func() {
		metrics.GetRequestCounter().WithLabelValues("grpc-batch-compress", status).Inc()
	}()
	
	startTime := time.Now()
	
	if len(req.Requests) == 0 {
		return &pb.BatchCompressResponse{
			Responses: []*pb.CompressImageResponse{},
		}, nil
	}
	
	// Convert protobuf requests to BatchRequest objects
	batchRequests := make([]api.BatchRequest, len(req.Requests))
	for i, protoReq := range req.Requests {
		batchRequests[i] = api.BatchRequest{
			Filename:  protoReq.Filename,
			Data:      bytes.NewReader(protoReq.ImageData),
			Format:    string(protoReq.Format),
			Quality:   int(protoReq.Quality),
			Algorithm: string(protoReq.Strategy),
		}
	}
	
	// Process using the unified batch processor
	batchResponse := a.service.ProcessBatchRequests(ctx, batchRequests)
	
	// Convert results to protobuf responses
	responses := make([]*pb.CompressImageResponse, 0, len(batchResponse.Results))
	for _, result := range batchResponse.Results {
		resp := &pb.CompressImageResponse{
			ImageData:        result.Data,
			Format:           result.Format,
			OriginalSize:     int64(result.OriginalSize),
			CompressedSize:   int64(result.CompressedSize),
			CompressionRatio: result.CompressionRatio,
			ProcessingTimeMs: result.ProcessingTime.Milliseconds(),
			Filename:         result.Filename,
		}
		responses = append(responses, resp)
	}
	
	// Handle errors
	if len(batchResponse.ProcessingErrors) > 0 {
		status = "partial_success"
		
		// Add error responses
		for _, err := range batchResponse.ProcessingErrors {
			errorResp := &pb.CompressImageResponse{
				Error:    err.Error.Error(),
				Filename: err.Filename,
			}
			responses = append(responses, errorResp)
		}
	}
	
	return &pb.BatchCompressResponse{
		Responses:             responses,
		TotalProcessingTimeMs: time.Since(startTime).Milliseconds(),
	}, nil
}

// StreamCompressImages handles streaming compression requests
func (a *Adapter) StreamCompressImages(stream pb.ImageCompressionService_StreamCompressImagesServer) error {
	for {
		// Receive the next request
		req, err := stream.Recv()
		if err == io.EOF {
			// End of stream
			return nil
		}
		if err != nil {
			return err
		}
		
		// Process the image
		resp, _ := a.CompressImage(stream.Context(), req)
		
		// Send the response
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}

// GetServiceStats returns statistics about the service
func (a *Adapter) GetServiceStats(ctx context.Context, req *pb.ServiceStatsRequest) (*pb.ServiceStatsResponse, error) {
	// Get basic stats about worker pool
	workerCount := a.service.GetWorkerCount()
	busyWorkers := a.service.GetBusyWorkerCount()
	
	// Get memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return &pb.ServiceStatsResponse{
		WorkerCount:      int32(workerCount),
		BusyWorkers:      int32(busyWorkers),
		MemoryUsageBytes: int64(m.Alloc),
	}, nil
}