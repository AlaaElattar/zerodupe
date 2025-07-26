package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"zerodupe/pkg/client"

	"github.com/spf13/cobra"
)

var (
	downloadServer   string
	downloadToken    string
	downloadOutput   string
	downloadFileName string
)

var downloadCmd = &cobra.Command{
	Use:   "download <filehash>",
	Short: "Download a file from the server by its hash",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fileHash := args[0]

		c := client.NewClient(downloadServer)
		c.SetToken(downloadToken)

		fmt.Printf("Downloading file with hash %s from %s\n", fileHash, downloadServer)

		outputPath := filepath.Join(downloadOutput, downloadFileName)
		if err := c.DownloadFile(fileHash, downloadOutput, downloadFileName); err != nil {
			log.Fatalf("Failed to download file: %v", err)
		}

		fmt.Printf("File downloaded successfully to: %s\n", outputPath)
	},
}

func init() {
	downloadCmd.Flags().StringVar(&downloadServer, "server", "http://localhost:8080", "Server URL")
	downloadCmd.Flags().StringVar(&downloadToken, "token", "", "JWT authentication token")
	downloadCmd.Flags().StringVarP(&downloadOutput, "output", "o", ".", "Output directory")
	downloadCmd.Flags().StringVarP(&downloadFileName, "name", "n", "", "Output file name (default: file hash)")
	downloadCmd.MarkFlagRequired("token")
}
