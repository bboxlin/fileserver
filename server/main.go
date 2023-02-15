package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Set a maximum upload size of 5MB
	const maxUploadSize = 5 * 1024 * 1024
	if r.ContentLength > maxUploadSize {
		http.Error(w, "file too large", http.StatusBadRequest)
		return
	}

	// Parse the multipart form data
	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the file from the multipart form
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file on disk to save the uploaded file
	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Use a 4KB buffer to read the request body in small chunks
	const bufferSize = 4 * 1024
	buffer := make([]byte, bufferSize)
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if n == 0 {
			break
		}

		// Write the chunk to the new file on disk
		_, err = f.Write(buffer[:n])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintln(w, "File uploaded successfully")
}

func main() {
	http.HandleFunc("/upload", uploadFileHandler)
	http.ListenAndServe(":8080", nil)
}
