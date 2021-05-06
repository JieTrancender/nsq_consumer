package cmd

import (
	"github.com/JieTrancender/nsq_to_consumer/internal/lg"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kbm",
	Short: "kbm means keyboard man service.",
	Run: func(cmd *cobra.Command, args []string) {
		logInfo("%s", "Hello World!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logFatal("execute fail, err: %v", err)
	}
}

func logFatal(f string, args ...interface{}) {
	lg.LogFatal("[nsq_consumer] ", f, args...)
}

func logInfo(f string, args ...interface{}) {
	lg.LogInfo("[nsq_consumer] ", f, args...)
}
