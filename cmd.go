package main

import "github.com/spf13/cobra"

var cmd = &cobra.Command{
	Use:   "poddle",
	Short: "podcast manager/server/app",
	Long:  `A tool for managing a library of podcasts.`,
	Run: func(c *cobra.Command, args []string) {
		start()
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start a server",
	Long:  `Start a web server to serve the Poddle API and app.`,
	Run: func(c *cobra.Command, args []string) {
	},
}

func init() {
	cmd.AddCommand(startCmd)
}
