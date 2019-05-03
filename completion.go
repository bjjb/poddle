package main

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate shell completion",
	Long: `Generate a shell completion script for the given shell. See '-h' for
available shells.`,
}

var bashCompletionCommand = &cobra.Command{
	Use:   "bash",
	Short: "Generate bash shell completion",
	Long: `Generate a completion script for the bash shell. For example, to
enable bash completion on a Debian based Linux system, run

  $ sudo poddle completion bash > /etc/bash_completion.d/poddle.bash

Adjust as needed for your OS.
`,
	Run: func(c *cobra.Command, args []string) {
		c.Parent().Parent().GenBashCompletion(os.Stdout)
	},
}

func init() {
	completionCmd.AddCommand(bashCompletionCommand)
	cmd.AddCommand(completionCmd)
}
