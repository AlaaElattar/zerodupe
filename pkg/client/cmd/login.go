package cmd

import (
	"fmt"
	"log"
	"zerodupe/pkg/client"

	"github.com/spf13/cobra"
)

var (
	loginServer   string
	loginUsername string
	loginPassword string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate a user",
	Run: func(cmd *cobra.Command, args []string) {
		if loginUsername == "" || loginPassword == "" {
			log.Fatal("username and password are required")
		}
		c := client.NewClient(loginServer)
		resp, err := c.Login(loginUsername, loginPassword)
		if err != nil {
			log.Fatalf("Failed to login: %v", err)
		}
		c.SetToken(resp.AccessToken)
		fmt.Println("Login successful.")
		fmt.Printf("Access token: %s\n", resp.AccessToken)
		fmt.Printf("Refresh token: %s\n", resp.RefreshToken)
	},
}

func init() {
	loginCmd.Flags().StringVar(&loginServer, "server", "http://localhost:8080", "Server URL")
	loginCmd.Flags().StringVar(&loginUsername, "username", "", "Username")
	loginCmd.Flags().StringVar(&loginPassword, "password", "", "Password")
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("password")
}
