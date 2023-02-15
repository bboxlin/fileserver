package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// http.HandleFunc("/upload", handleUpload) 
	http.ListenAndServe(":5555", nil)
}


func handleUpload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Determine the target filename based on the original filename.
	targetFilename := filepath.Join("uploads", header.Filename)

	// Create a new file on disk to save the uploaded file.
	outFile, err := os.Create(targetFilename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Copy the contents of the uploaded file to the new file on disk.
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success response.
	fmt.Fprintln(w, "File uploaded successfully.")
}

// func handleUpload(w http.ResponseWriter, r *http.Request) {
// 	file, header, err := r.FormFile("file")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	// Create a temporary file to save the uploaded file.
// 	tempFile, err := ioutil.TempFile("", "")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer tempFile.Close()

// 	// Write the contents of the uploaded file to the temporary file.
// 	_, err = io.Copy(tempFile, file)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Determine the target filename based on the original filename.
// 	targetFilename := filepath.Join("uploads", header.Filename)

// 	// Move the temporary file to the target filename.
// 	err = os.Rename(tempFile.Name(), targetFilename)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Return a success response.
// 	fmt.Fprintln(w, "File uploaded successfully.")
// }
