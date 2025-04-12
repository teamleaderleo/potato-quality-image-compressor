package main

import (
	"fmt"

	"log"
	"net/http"
	"os"


    "github.com/teamleaderleo/potato-quality-image-compressor/internal/api"
    "github.com/teamleaderleo/potato-quality-image-compressor/internal/metrics"
)

func main() {
    metrics.Init()

    service := api.NewService()

    http.HandleFunc("/", handleRoot)
    http.HandleFunc("/compress", service.HandleCompress)
    http.HandleFunc("/batch-compress", service.HandleBatchCompress)

    if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
        api.StartLambda()
    } else {
        port := os.Getenv("PORT")
        if port == "" {
            port = "8080" // default port
        }

        log.Printf("Server starting on http://localhost:%s", port)
        log.Fatal(http.ListenAndServe(":"+port, nil))
    }
}

// func handleTest(w http.ResponseWriter, r *http.Request) {
//     http.ServeFile(w, r, "test.html")
// }

func handleRoot(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Image Compression API")
}