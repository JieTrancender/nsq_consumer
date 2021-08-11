package main

import (
	"os"

	"github.com/JieTrancender/nsq_consumer/cmd"
)

func main() {
	if err := cmd.NsqConsumer(cmd.NsqConsumerSettings()).Execute(); err != nil {
		os.Exit(1)
	}
}
