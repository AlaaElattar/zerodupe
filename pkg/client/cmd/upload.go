package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"zerodupe/pkg/client"
)

var (
	uploadServer       string
	uploadToken        string
	uploadRefreshToken string
)

var uploadCmd = &cobra.Command{
	Use:   "upload <filepath>",
	Short: "Upload a file to the server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		fmt.Printf("Uploading file %s to %s\n", filePath, uploadServer)

		c := client.NewClient(uploadServer)
		c.SetToken(uploadToken)

		err := c.ExecuteWithAuth(func() error {
			return c.UploadFile(filePath)
		})
		if err != nil {
			log.Fatalf("Failed to upload file: %v", err)
		}
		fmt.Println("File uploaded successfully.")
	},
}

func init() {
	uploadCmd.Flags().StringVar(&uploadServer, "server", "http://localhost:8080", "Server URL")
	uploadCmd.Flags().StringVar(&uploadToken, "token", "", "JWT authentication token")
	uploadCmd.Flags().StringVar(&uploadRefreshToken, "refresh-token", "", "Refresh token")
	uploadCmd.MarkFlagRequired("token")
}
