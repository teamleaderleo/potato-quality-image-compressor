package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/teamleaderleo/potato-quality-image-compressor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

const (
	address        = "localhost:9000"
	imagePath      = "assets/test.png"
	totalRequests  = 100
	concurrency    = 16
	requestTimeout = 15 * time.Second
)

type grpcResult struct {
	duration         time.Duration
	originalSize     int64
	compressedSize   int64
	compressionRatio float64
	err              error
}

func sendCompressRequest(client pb.ImageCompressionServiceClient, imageData []byte) grpcResult {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	start := time.Now()
	resp, err := client.CompressImage(ctx, &pb.CompressImageRequest{
		ImageData: imageData,
		Quality:   80,
		Format:    "webp",
		Strategy:  "scale",
		Filename:  "test.png",
	})
	duration := time.Since(start)

	if err != nil {
		return grpcResult{duration: duration, err: err}
	}

	return grpcResult{
		duration:         duration,
		originalSize:     resp.OriginalSize,
		compressedSize:   resp.CompressedSize,
		compressionRatio: resp.CompressionRatio,
	}
}

func main() {
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		log.Fatalf("Failed to read image: %v", err)
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewImageCompressionServiceClient(conn)

	fmt.Printf("Starting %d gRPC requests...\n", totalRequests)

	results := make([]grpcResult, totalRequests)
	sem := make(chan struct{}, concurrency)

	for i := 0; i < totalRequests; i++ {
		sem <- struct{}{}
		go func(index int) {
			defer func() { <-sem }()
			results[index] = sendCompressRequest(client, imageData)
		}(i)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}

	var successCount, errorCount int
	var totalTime time.Duration
	var minTime, maxTime time.Duration

	for i, res := range results {
		if res.err != nil {
			errorCount++
			continue
		}
		successCount++
		totalTime += res.duration
		if i == 0 || res.duration < minTime {
			minTime = res.duration
		}
		if res.duration > maxTime {
			maxTime = res.duration
		}
	}

	fmt.Printf("Completed %d/%d requests\n", successCount, totalRequests)
	if successCount > 0 {
		avg := totalTime / time.Duration(successCount)
		fmt.Printf("Average: %v | Min: %v | Max: %v\n", avg, minTime, maxTime)
	}
	fmt.Printf("Errors: %d\n", errorCount)
}
