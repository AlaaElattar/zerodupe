package main

import (
	"log"
	"zerodupe/client"
)

func main() {
	client := client.NewClient("http://localhost:8080")
	err := client.StartServer()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	err = client.DownloadFile("test2.txt")
	if err != nil {
		log.Fatalf("Failed to download file: %v", err)
		return
	}

}
