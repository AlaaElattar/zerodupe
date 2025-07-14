package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"zerodupe/pkg/client"
)

var (
	refreshServer string
	refreshToken  string
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh the access token using a refresh token",
	Run: func(cmd *cobra.Command, args []string) {
		if refreshToken == "" {
			log.Fatal("refresh token is required")
		}
		c := client.NewClient(refreshServer)
		resp, err := c.RefreshToken(refreshToken)
		if err != nil {
			log.Fatalf("Failed to refresh token: %v", err)
		}
		fmt.Println("Token refreshed successfully")
		fmt.Printf("New access token: %s\n", resp.AccessToken)
		fmt.Printf("New refresh token: %s\n", resp.RefreshToken)
	},
}

func init() {
	refreshCmd.Flags().StringVar(&refreshServer, "server", "http://localhost:8080", "Server URL")
	refreshCmd.Flags().StringVar(&refreshToken, "token", "", "Refresh token")
	refreshCmd.MarkFlagRequired("token")
}
