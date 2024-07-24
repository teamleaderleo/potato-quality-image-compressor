package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/chai2010/webp"
	"golang.org/x/image/draw"
)

func main() {
    http.HandleFunc("/", handleRoot)
    http.HandleFunc("/compress", handleCompress)
    http.HandleFunc("/batch-compress", handleBatchCompress)
    http.HandleFunc("/test", handleTest)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // default port
    }

    log.Printf("Server starting on http://localhost:%s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleTest(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "test.html")
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Image Compression API")
}

func handleCompress(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    file, header, err := r.FormFile("image")
    if err != nil {
        http.Error(w, "Error retrieving the file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    quality, _ := strconv.Atoi(r.FormValue("quality"))
    if quality == 0 {
        quality = 75 // default quality
    }

    outputFormat := r.FormValue("format")
    if outputFormat == "" {
        outputFormat = "jpeg" // default format
    }

    compressedImageBytes, err := compressAndEncodeImage(file, quality, outputFormat)
    if err != nil {
        http.Error(w, "Error compressing image", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Image compressed successfully",
        "filename": header.Filename,
        "size": len(compressedImageBytes),
    })
}

func handleBatchCompress(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    err := r.ParseMultipartForm(32 << 20) // 32 MB max
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    files := r.MultipartForm.File["images"]
    quality, _ := strconv.Atoi(r.FormValue("quality"))
    if quality == 0 {
        quality = 75 // default quality
    }
    outputFormat := r.FormValue("format")
    if outputFormat == "" {
        outputFormat = "jpeg" // default format
    }

    results := make([]map[string]interface{}, 0, len(files))

    for _, fileHeader := range files {
        file, err := fileHeader.Open()
        if err != nil {
            results = append(results, map[string]interface{}{
                "filename": fileHeader.Filename,
                "error": "Error opening file",
            })
            continue
        }
        defer file.Close()

        compressedImageBytes, err := compressAndEncodeImage(file, quality, outputFormat)
        if err != nil {
            results = append(results, map[string]interface{}{
                "filename": fileHeader.Filename,
                "error": "Error compressing image",
            })
            continue
        }

        results = append(results, map[string]interface{}{
            "filename": fileHeader.Filename,
            "size": len(compressedImageBytes),
        })
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Batch compression complete",
        "results": results,
    })
}

func compressAndEncodeImage(file io.Reader, quality int, outputFormat string) ([]byte, error) {
    img, _, err := image.Decode(file)
    if err != nil {
        return nil, err
    }

    compressedImg := compressImage(img, quality)

    var buf []byte
    outputBuf := bytes.NewBuffer(buf)

    switch outputFormat {
    case "jpeg", "jpg":
        err = jpeg.Encode(outputBuf, compressedImg, &jpeg.Options{Quality: quality})
    case "png":
        err = png.Encode(outputBuf, compressedImg)
    case "webp":
        err = webp.Encode(outputBuf, compressedImg, &webp.Options{Lossless: false, Quality: float32(quality)})
    default:
        return nil, fmt.Errorf("unsupported format: %s", outputFormat)
    }

    if err != nil {
        return nil, err
    }

    return outputBuf.Bytes(), nil
}

func compressImage(img image.Image, quality int) image.Image {
    bounds := img.Bounds()
    width, height := bounds.Dx(), bounds.Dy()
    newWidth, newHeight := width, height

    // Calculate new dimensions based on quality
    scaleFactor := float64(quality) / 100.0
    if width > 1000 || height > 1000 {
        newWidth = int(float64(width) * scaleFactor)
        newHeight = int(float64(height) * scaleFactor)
    }

    newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
    draw.ApproxBiLinear.Scale(newImg, newImg.Bounds(), img, img.Bounds(), draw.Over, nil)

    return newImg
}