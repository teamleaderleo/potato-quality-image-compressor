package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chai2010/webp"
	"golang.org/x/image/draw"
)

func main() {
    http.HandleFunc("/", handleRoot)
    http.HandleFunc("/upload", handleUpload)
    http.HandleFunc("/upload.html", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "upload.html")
    })

    log.Println("Server starting on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/upload.html", http.StatusSeeOther)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    err := r.ParseMultipartForm(10 << 20) // 10 MB max
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    file, handler, err := r.FormFile("image")
    if err != nil {
        http.Error(w, "Error retrieving the file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    log.Printf("Received file: %s\n", handler.Filename)

    img, format, err := image.Decode(file)
    if err != nil {
        http.Error(w, "Error decoding image", http.StatusBadRequest)
        return
    }

    quality, _ := strconv.Atoi(r.FormValue("quality"))
    if quality == 0 {
        quality = 75 // default quality
    }

    outputFormat := r.FormValue("format")
    if outputFormat == "" {
        outputFormat = format // default to input format
    }

    compressedImg := compressImage(img, quality)

    if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
        http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
        return
    }

    // Get the file extension from the original filename
    origExt := filepath.Ext(handler.Filename)
    baseFileName := strings.TrimSuffix(handler.Filename, origExt)

    // Determine the new file extension based on the output format
    var newExt string
    switch outputFormat {
    case "jpeg", "jpg":
        newExt = ".jpg"
    case "png":
        newExt = ".png"
    case "webp":
        newExt = ".webp"
    default:
        newExt = origExt // Use original extension if format is not recognized
    }

    outputPath := filepath.Join("./uploads", "compressed_"+baseFileName+newExt)
    out, err := os.Create(outputPath)
    if err != nil {
        http.Error(w, "Error creating the file", http.StatusInternalServerError)
        return
    }
    defer out.Close()

    switch outputFormat {
    case "jpeg", "jpg":
        jpeg.Encode(out, compressedImg, &jpeg.Options{Quality: quality})
    case "png":
        png.Encode(out, compressedImg)
    case "webp":
        err = webp.Encode(out, compressedImg, &webp.Options{Lossless: false, Quality: float32(quality)})
        if err != nil {
            http.Error(w, "Error encoding WebP image", http.StatusInternalServerError)
            return
        }
    default:
        http.Error(w, "Unsupported output format", http.StatusBadRequest)
        return
    }

    log.Printf("File compressed and saved as: %s\n", outputPath)
    fmt.Fprintf(w, "File compressed and saved as: %s", outputPath)
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