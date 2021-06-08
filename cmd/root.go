package cmd

import (
	"fmt"
	"os"

	"github.com/JieTrancender/nsq_to_consumer/cmd/instance"
	"github.com/JieTrancender/nsq_to_consumer/internal/version"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const Name = "NsqConsumer"

type NsqConsumerRootCmd struct {
	cobra.Command
	RunCmd     *cobra.Command
	VersionCmd *cobra.Command
}

var (
	RootCmd *NsqConsumerRootCmd
)

func NsqConsumerSettings() instance.Settings {
	var runFlags = pflag.NewFlagSet(Name, pflag.ExitOnError)
	runFlags.BoolP("version", "v", false, "show version")
	runFlags.StringArray("etcd-endpoints", []string{"127.0.0.1:2379"}, "etcd endpoints(may be given multi time)")
	runFlags.String("etcd-path", "/config/nsq_consumer/default", "etcd path")
	runFlags.String("etcd-username", "root", "etcd username")
	runFlags.String("etcd-password", "root", "etcd password")
	return instance.Settings{
		RunFlags: runFlags,
		Name:     Name,
	}
}

func NsqConsumer(settings instance.Settings) *NsqConsumerRootCmd {
	command := genRootCmdWithSettings(settings)
	return command
}

func genRootCmdWithSettings(settings instance.Settings) *NsqConsumerRootCmd {
	if settings.IndexPrefix == "" {
		settings.IndexPrefix = settings.Name
	}

	rootCmd := &NsqConsumerRootCmd{}
	rootCmd.Use = settings.Name

	rootCmd.RunCmd = genRunCmd(settings)
	rootCmd.VersionCmd = GenVersionCmd(settings)

	// Root command is an alias for run
	rootCmd.Run = rootCmd.RunCmd.Run

	rootCmd.Flags().AddFlagSet(rootCmd.RunCmd.Flags())

	// Register subcommands common to all consumers
	rootCmd.AddCommand(rootCmd.RunCmd)
	rootCmd.AddCommand(rootCmd.VersionCmd)

	return rootCmd
}

func genRunCmd(settings instance.Settings) *cobra.Command {
	name := settings.Name
	runCmd := cobra.Command{
		Use:   "run",
		Short: "Run " + name,
		Run: func(cmd *cobra.Command, args []string) {
			isVersion, _ := cmd.Flags().GetBool("version")
			if isVersion {
				fmt.Println(version.String())
				os.Exit(0)
			}

			err := instance.Run(settings)
			if err != nil {
				os.Exit(1)
			}
		},
	}

	if settings.RunFlags != nil {
		runCmd.Flags().AddFlagSet(settings.RunFlags)
	}

	return &runCmd
}
