package main

import (
	"log"
	"zerodupe/client"
)

func main() {
	client := client.NewClient("http://localhost:8080")

	client.SetUploadDir("client/upload-dir")
	client.SetDownloadDir("client/download-dir")
	err := client.StartServer()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	err = client.DownloadFile("test_5chunks.dat")
	if err != nil {
		log.Fatalf("Failed to download file: %v", err)
		return
	}

}