package cmd

import (
	"fmt"
	"os"

	"github.com/JieTrancender/nsq_consumer/internal/version"
	"github.com/JieTrancender/nsq_consumer/libconsumer/cmd/instance"
	"github.com/JieTrancender/nsq_consumer/libconsumer/consumer"
	"github.com/spf13/cobra"
)

func genRunCmd(settings instance.Settings, ct consumer.Creator) *cobra.Command {
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

			err := instance.Run(settings, ct)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	if settings.RunFlags != nil {
		runCmd.Flags().AddFlagSet(settings.RunFlags)
	}

	return &runCmd
}
