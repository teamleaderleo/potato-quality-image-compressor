package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

const (
	url            = "http://localhost:8080/compress"
	imagePath      = "test/assets/test.png"
	totalRequests  = 100
	concurrency    = 16
	requestTimeout = 15 * time.Second
)

type httpResult struct {
	duration time.Duration
	err      error
}

func createRequest(image []byte) (*http.Request, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("image", path.Base(imagePath))
	if err != nil {
		return nil, err
	}
	_, err = part.Write(image)
	if err != nil {
		return nil, err
	}

	writer.WriteField("quality", "80")
	writer.WriteField("format", "webp")
	writer.WriteField("algorithm", "scale")

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func main() {
	image, err := os.ReadFile(imagePath)
	if err != nil {
		panic(fmt.Errorf("failed to read image: %v", err))
	}

	fmt.Printf("Starting %d HTTP requests...\n", totalRequests)

	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)
	results := make([]httpResult, totalRequests)

	client := &http.Client{
		Timeout: requestTimeout,
	}

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()

			req, err := createRequest(image)
			if err != nil {
				results[i] = httpResult{err: err}
				return
			}

			start := time.Now()
			resp, err := client.Do(req)
			duration := time.Since(start)

			if err != nil {
				results[i] = httpResult{duration: duration, err: err}
				return
			}
			defer resp.Body.Close()
			io.Copy(io.Discard, resp.Body)

			if resp.StatusCode != http.StatusOK {
				results[i] = httpResult{duration: duration, err: fmt.Errorf("HTTP %d", resp.StatusCode)}
				return
			}

			results[i] = httpResult{duration: duration}
		}(i)
	}

	wg.Wait()

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
