package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func uploadFile(filename string, targetUrl string) error {
	// Open the file to be uploaded
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new HTTP POST request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	// Send the HTTP POST request to the server
	req, err := http.NewRequest("POST", targetUrl, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload file: %v", resp.Status)
	}

	return nil
}

func main() {
	err := uploadFile("test.txt", "http://localhost:8080/upload")
	if err != nil {
		fmt.Println(err)
	}
}
