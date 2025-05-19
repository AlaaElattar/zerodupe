package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"zerodupe/packages/client"
)

func main() {

	// zerodupe upload -server http://myhost:8080 file.txt
	// zerodupe download -server http://myhost:8080 -o ./downloads -n custom_name.txt HASH

	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	uploadServerURL := uploadCmd.String("server", "http://localhost:8080", "Server URL")

	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	downloadServerURL := downloadCmd.String("server", "http://localhost:8080", "Server URL")

	downloadOutput := downloadCmd.String("o", ".", "Output directory")
	downloadFileName := downloadCmd.String("n", "", "Output file name (default: file hash)")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "upload":
		uploadCmd.Parse(os.Args[2:])
		if uploadCmd.NArg() < 1 {
			fmt.Println("Error: Missing file path")
			fmt.Println("Usage: zerodupe upload [flags] <filepath>")
			uploadCmd.PrintDefaults()
			os.Exit(1)
		}
		filePath := uploadCmd.Arg(0)
		uploadFile(*uploadServerURL, filePath)

	case "download":
		downloadCmd.Parse(os.Args[2:])
		if downloadCmd.NArg() < 1 {
			fmt.Println("Error: Missing file hash")
			fmt.Println("Usage: zerodupe download [flags] <filehash>")
			downloadCmd.PrintDefaults()
			os.Exit(1)
		}
		fileHash := downloadCmd.Arg(0)
		downloadFile(*downloadServerURL, fileHash, *downloadOutput, *downloadFileName)

	default:
		printUsage()
		os.Exit(1)

	}

}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  zerodupe upload [flags] <filepath>")
	fmt.Println("  zerodupe download [flags] <filehash>")
	fmt.Println("\nFlags for upload:")
	fmt.Println("  -server string    Server URL (default \"http://localhost:8080\")")
	fmt.Println("\nFlags for download:")
	fmt.Println("  -server string    Server URL (default \"http://localhost:8080\")")
	fmt.Println("  -o string         Output directory (default \".\")")
	fmt.Println("  -n string         Output file name (default: file hash)")
}
func uploadFile(serverURL, filePath string) {
	fmt.Printf("Uploading file %s to %s\n", filePath, serverURL)

	client := client.NewClient(serverURL)

	err := client.UploadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}
}

func downloadFile(serverURL, fileHash, outputDir, fileName string) {
	fmt.Printf("Downloading file with hash %s from %s\n", fileHash, serverURL)

	client := client.NewClient(serverURL)

	err := client.DownloadFile(fileHash, outputDir, fileName)
	if err != nil {
		log.Fatalf("Failed to download file: %v", err)
	}
}
