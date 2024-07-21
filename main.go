package main

import (
	"fmt" // formatted I/O
	"io" // basic I/O operations
	"log" // logging
	"net/http" // HTTP client and server implementations
	"os" // OS functionality
)

func main() {
	http.HandleFunc("/", handleRoot) // sets up the root route
	http.HandleFunc("/upload", handleUpload) // sets up the upload route

	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil)) // starts the server on port 8080
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Potato Quality Image Compressor!")
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from the request
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the uploads directory if it doesn't exist
	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
		http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	dst, err := os.Create(fmt.Sprintf("./uploads/%s", handler.Filename))
	if err != nil {
		http.Error(w, "Error creating the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error copying the file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", handler.Filename)
}