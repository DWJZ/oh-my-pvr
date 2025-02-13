package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	// File upload endpoint with CORS enabled.
	http.HandleFunc("/upload", enableCORS(uploadHandler))

	// Start the server.
	fmt.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("Failed to start server:", err)
	}
}

// enableCORS is a middleware function that sets CORS headers.
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight (OPTIONS) requests.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler.
		next(w, r)
	}
}

// uploadHandler processes file uploads.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r == nil {
		http.Error(w, "Request is nil", http.StatusBadRequest)
		return
	}

	// Parse the multipart form with a generous max memory allocation.
	if err := r.ParseMultipartForm(10 << 28); err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusInternalServerError)
		return
	}

	// Read the upload directory from the request parameters.
	// The client can send this as a form field or query parameter.
	uploadDir := r.FormValue("uploadDir")
	if uploadDir == "" {
		// Fallback to default if not provided.
		uploadDir = "./uploads"
	}

	// Clean the path to ensure no directory traversal occurs.
	cleanUploadDir := filepath.Clean(uploadDir)

	// Convert to an absolute path for clarity.
	absUploadDir, err := filepath.Abs(cleanUploadDir)
	if err != nil {
		http.Error(w, "Failed to get absolute path", http.StatusInternalServerError)
		return
	}
	log.Printf("Saving files to: %s", absUploadDir)

	// Create the upload directory if it doesn't exist.
	if err := os.MkdirAll(absUploadDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// Retrieve all uploaded files from the "file" field.
	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	// Use a WaitGroup to handle concurrent file uploads.
	var wg sync.WaitGroup
	wg.Add(len(files))

	// Slices to store upload results and errors.
	results := make([]string, len(files))
	errors := make([]error, len(files))

	// Process each file concurrently.
	for i, fileHeader := range files {
		go func(index int, header *multipart.FileHeader) {
			defer wg.Done()

			// Open the uploaded file.
			file, err := header.Open()
			if err != nil {
				errors[index] = err
				return
			}
			defer file.Close()

			// Create the destination file path.
			destPath := filepath.Join(absUploadDir, filepath.Base(header.Filename))
			log.Printf("Saving file: %s", destPath)

			// Create the destination file.
			dst, err := os.Create(destPath)
			if err != nil {
				errors[index] = err
				return
			}
			defer dst.Close()

			// Copy the uploaded file contents to the destination file.
			if _, err := io.Copy(dst, file); err != nil {
				errors[index] = err
				return
			}

			// Record the success message.
			results[index] = fmt.Sprintf("File uploaded successfully: %s", header.Filename)
		}(i, fileHeader)
	}

	// Wait for all file uploads to finish.
	wg.Wait()

	// Check for any errors during the upload process.
	for _, err := range errors {
		if err != nil {
			http.Error(w, "Failed to upload some files", http.StatusInternalServerError)
			return
		}
	}

	// Return a successful response with the results.
	w.WriteHeader(http.StatusOK)
	for _, result := range results {
		fmt.Fprintln(w, result)
	}
}
