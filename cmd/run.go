package cmd

import (
	"fmt"
	"os"

	"github.com/JieTrancender/nsq_to_consumer/cmd/instance"
	"github.com/JieTrancender/nsq_to_consumer/internal/version"
	"github.com/spf13/cobra"
)

func GenRunCmd(settings instance.Settings) *cobra.Command {
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
