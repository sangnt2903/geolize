package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "geolize",
	Short: "Geolize is a service for IP geolocation",
	Long: `Geolize is a service that provides IP geolocation information.
It allows you to lookup geographical information based on IP addresses.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommands are provided, print the help message
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
