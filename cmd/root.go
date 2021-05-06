package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kbm",
	Short: "kbm means keyboard man service.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello World!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
