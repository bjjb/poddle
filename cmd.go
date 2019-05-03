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

func init() {
	cmd.PersistentFlags().StringP("database", "D", "sqlite3::memory:", "database DSN")
	cmd.PersistentFlags().String("ffmpeg-path", "ffmpeg", "path to an ffmpeg executable")
	cmd.PersistentFlags().String("search-backend", "itunes", "preferred search backend")
}
