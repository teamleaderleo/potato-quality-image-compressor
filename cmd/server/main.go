package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/chai2010/webp"
    "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/image/draw"

    "github.com/teamleaderleo/potato-quality-image-compressor/internal/api"
    "github.com/teamleaderleo/potato-quality-image-compressor/internal/metrics"
)

func main() {
    metrics.Init()

    _ = api.NewService() 

    http.HandleFunc("/", handleRoot)
    http.HandleFunc("/compress", handleCompress)
    http.HandleFunc("/batch-compress", handleBatchCompress)

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

func handleCompress(w http.ResponseWriter, r *http.Request) {
    timer := prometheus.NewTimer(requestDuration.WithLabelValues("compress"))
    defer timer.ObserveDuration()

    status := "success"
    defer func() {
        requestCounter.WithLabelValues("compress", status).Inc()
    }()

    if r.Method != http.MethodPost {
        status = "method_not_allowed"
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    file, header, err := r.FormFile("image")
    if err != nil {
        status = "bad_request"
        http.Error(w, "Error retrieving the file: "+err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    quality, err := strconv.Atoi(r.FormValue("quality"))
    if err != nil || quality < 1 || quality > 100 {
        quality = 75 // default quality
    }

    outputFormat := r.FormValue("format")
    if outputFormat == "" {
        outputFormat = "webp" // default format
    }

    compressedImageBytes, err := compressAndEncodeImage(file, quality, outputFormat)
    if err != nil {
        status = "compression_error"
        http.Error(w, "Error compressing image: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", fmt.Sprintf("image/%s", outputFormat))
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.%s", filepath.Base(header.Filename), outputFormat))
    _, err = w.Write(compressedImageBytes)
    if err != nil {
        log.Printf("Error writing response: %v", err)
    }
}

func handleBatchCompress(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    err := r.ParseMultipartForm(32 << 20) // 32 MB max
    if err != nil {
        http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
        return
    }

    files := r.MultipartForm.File["images"]
    quality, err := strconv.Atoi(r.FormValue("quality"))
    if err != nil || quality < 1 || quality > 100 {
        quality = 75 // default quality
    }
    outputFormat := r.FormValue("format")
    if outputFormat == "" {
        outputFormat = "webp" // default format
    }

    buf := new(bytes.Buffer)
    zipWriter := zip.NewWriter(buf)

    for _, fileHeader := range files {
        file, err := fileHeader.Open()
        if err != nil {
            log.Printf("Error opening file %s: %v", fileHeader.Filename, err)
            continue
        }

        compressedImageBytes, err := compressAndEncodeImage(file, quality, outputFormat)
        file.Close()
        if err != nil {
            log.Printf("Error compressing file %s: %v", fileHeader.Filename, err)
            continue
        }

        zipFile, err := zipWriter.Create(fmt.Sprintf("%s.%s", filepath.Base(fileHeader.Filename), outputFormat))
        if err != nil {
            log.Printf("Error creating zip entry for %s: %v", fileHeader.Filename, err)
            continue
        }

        _, err = zipFile.Write(compressedImageBytes)
        if err != nil {
            log.Printf("Error writing to zip for %s: %v", fileHeader.Filename, err)
        }
    }

    zipWriter.Close()

    w.Header().Set("Content-Type", "application/zip")
    w.Header().Set("Content-Disposition", "attachment; filename=compressed_images.zip")
    _, err = w.Write(buf.Bytes())
    if err != nil {
        log.Printf("Error writing zip response: %v", err)
    }
}

func compressAndEncodeImage(file io.Reader, quality int, outputFormat string) ([]byte, error) {
    img, _, err := image.Decode(file)
    if err != nil {
        return nil, fmt.Errorf("error decoding image: %v", err)
    }

    compressedImg := compressImage(img, quality)

    var buf bytes.Buffer
    switch outputFormat {
    case "webp":
        err = webp.Encode(&buf, compressedImg, &webp.Options{Quality: float32(quality)})
    case "jpeg", "jpg":
        err = jpeg.Encode(&buf, compressedImg, &jpeg.Options{Quality: quality})
    case "png":
        err = png.Encode(&buf, compressedImg)
    default:
        return nil, fmt.Errorf("unsupported format: %s", outputFormat)
    }

    if err != nil {
        return nil, fmt.Errorf("error encoding image to %s: %v", outputFormat, err)
    }

    return buf.Bytes(), nil
}

func compressImage(img image.Image, quality int) image.Image {
    // If quality is 100, return the original image
    if quality == 100 {
        return img
    }

    bounds := img.Bounds()
    width, height := bounds.Dx(), bounds.Dy()
    
    // Calculate new dimensions based on quality
    scaleFactor := float64(quality) / 100.0
    newWidth := int(float64(width) * scaleFactor)
    newHeight := int(float64(height) * scaleFactor)

    // Ensure minimum dimensions
    if newWidth < 10 {
        newWidth = 10
    }
    if newHeight < 10 {
        newHeight = 10
    }

    newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
    draw.ApproxBiLinear.Scale(newImg, newImg.Bounds(), img, img.Bounds(), draw.Over, nil)

    return newImg
}