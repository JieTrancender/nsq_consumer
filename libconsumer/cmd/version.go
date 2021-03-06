package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/JieTrancender/nsq_consumer/libconsumer/cmd/instance"
	"github.com/JieTrancender/nsq_consumer/libconsumer/common/cli"
)

func genVersionCmd(settings instance.Settings) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show current version info",
		Run: cli.RunWith(
			func(_ *cobra.Command, args []string) error {
				consumer, err := instance.NewConsumer(settings.Name, settings.IndexPrefix, settings.Version)
				if err != nil {
					return fmt.Errorf("error initializing consumer: %s", err)
				}

				fmt.Printf("%s version %s Arch(%s) runtime(%s)\n", consumer.Info.Consumer, consumer.Info.Version, runtime.GOARCH, runtime.Version())
				return nil
			}),
	}
}
