package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "Zerodupe Client",
	Short: "Client for Zerodupe deduplication file storage system",
}

func init() {
	rootCmd.AddCommand(signupCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(refreshCmd)
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(downloadCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
