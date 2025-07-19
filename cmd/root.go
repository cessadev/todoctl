package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "todoctl",
	Short: "Task manager from the terminal",
	Long:  `A CLI tool to manage your daily tasks`,
	Run: func(cmd *cobra.Command, args []string) {
		/** If executed without subcommands, display help */
		cmd.Help()
	},
}

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
