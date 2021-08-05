package cmd

import (
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/cmd/instance"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/consumer"
	"github.com/spf13/cobra"
)

type NsqConsumerRootCmd struct {
	cobra.Command
	RunCmd     *cobra.Command
	VersionCmd *cobra.Command
}

func GenRootCmdWithSettings(ct consumer.Creator, settings instance.Settings) *NsqConsumerRootCmd {
	if settings.IndexPrefix == "" {
		settings.IndexPrefix = settings.Name
	}

	rootCmd := &NsqConsumerRootCmd{}
	rootCmd.Use = settings.Name

	rootCmd.RunCmd = genRunCmd(settings, ct)
	rootCmd.VersionCmd = genVersionCmd(settings)

	// Root command is an alias for run
	rootCmd.Run = rootCmd.RunCmd.Run

	rootCmd.Flags().AddFlagSet(rootCmd.RunCmd.Flags())

	// Register subcommands common to all consumers
	rootCmd.AddCommand(rootCmd.RunCmd)
	rootCmd.AddCommand(rootCmd.VersionCmd)

	return rootCmd
}
