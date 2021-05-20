package main

import (
	"os"

	"github.com/JieTrancender/nsq_to_consumer/cmd"
)

func main() {
	// cmd.Execute()
	if err := cmd.NsqConsumer(cmd.NsqConsumerSettings()).Execute(); err != nil {
		os.Exit(1)
	}
}
