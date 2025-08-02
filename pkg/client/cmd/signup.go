package cmd

import (
	"fmt"
	"log"
	"zerodupe/pkg/client"

	"github.com/spf13/cobra"
)

var (
	server          string
	username        string
	password        string
	confirmPassword string
)

var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Create a new user account",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("username: %q, password: %q, confirm_password: %q\n", username, password, confirmPassword)
		if username == "" || password == "" || confirmPassword == "" {
			log.Fatal("username, password, and confirm password are required")
		}
		c := client.NewClient(server)
		err := c.Signup(username, password, confirmPassword)
		if err != nil {
			log.Fatalf("Failed to signup: %v", err)
		}
		fmt.Println("Signup successful.")
	},
}

func init() {
	signupCmd.Flags().StringVar(&server, "server", "http://localhost:8080", "Server URL")
	signupCmd.Flags().StringVar(&username, "username", "", "Username")
	signupCmd.Flags().StringVar(&password, "password", "", "Password")
	signupCmd.Flags().StringVar(&confirmPassword, "confirm-password", "", "Confirm Password")

	signupCmd.MarkFlagRequired("username")
	signupCmd.MarkFlagRequired("password")
	signupCmd.MarkFlagRequired("confirm-password")
}
