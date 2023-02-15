package main

import (
    "bytes"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "os"
)

func main() {
    file, err := os.Open("test.txt")
    if err != nil {
        fmt.Println("Failed to open file:", err)
        return
    }
    defer file.Close()

    // Create a buffer to store the HTTP request body.
    requestBody := &bytes.Buffer{}

    // Create a new multipart writer to write the HTTP request body.
    writer := multipart.NewWriter(requestBody)

    // Create a new part for the file data.
    part, err := writer.CreateFormFile("file", "test.txt")
    if err != nil {
        fmt.Println("Failed to create form file:", err)
        return
    }

    // Use a buffer to read the file in small chunks.
    buffer := make([]byte, 1024)
    for {
        // Read a chunk of data from the file.
        bytesRead, err := file.Read(buffer)
        if err != nil && err != io.EOF {
            fmt.Println("Failed to read from file:", err)
            return
        }

        // If we've reached the end of the file, break out of the loop.
        if bytesRead == 0 {
            break
        }

        // Write the chunk of data to the HTTP request body.
        part.Write(buffer[:bytesRead])
    }

    // Close the multipart writer to finalize the request body.
    writer.Close()

    // Create a new HTTP request to upload the file.
    req, err := http.NewRequest("POST", "http://localhost:5555/upload", requestBody)
    if err != nil {
        fmt.Println("Failed to create HTTP request:", err)
        return
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())

    // Send the HTTP request and print the response.
    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        fmt.Println("Failed to upload file:", err)
        return
    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        fmt.Println("Failed to upload file: ", res.Status)
        return
    }

    fmt.Println("File uploaded successfully.")
}
