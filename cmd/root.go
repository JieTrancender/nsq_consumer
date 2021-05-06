package cmd

import (
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/JieTrancender/nsq_to_consumer/internal/lg"
	"github.com/JieTrancender/nsq_to_consumer/nsq_consumer"

	"github.com/judwhite/go-svc"

	"github.com/spf13/cobra"
)

type program struct {
	once     sync.Once
	consumer *nsq_consumer.NsqConsumer
}

var rootCmd = &cobra.Command{
	Use:   "kbm",
	Short: "kbm means keyboard man service.",
	Run: func(cmd *cobra.Command, args []string) {
		prg := &program{}
		if err := svc.Run(prg, syscall.SIGINT); err != nil {
			logFatal("%s", err)
		}
	},
}

func (p *program) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil
}

func (p *program) Start() error {
	opts := nsq_consumer.NewOptions()
	consumer, err := nsq_consumer.NewNsqConsumer(opts)
	if err != nil {
		logFatal("new nsq consumer fail, err: %s", err)
	}
	p.consumer = consumer

	signalChan := make(chan os.Signal, 1)
	go func() {
		for range signalChan {
			p.once.Do(func() {
				p.consumer.Exit()
			})
		}
	}()
	signal.Notify(signalChan, syscall.SIGTERM)

	go func() {
		err := p.consumer.Main()
		if err != nil {
			_ = p.Stop()
			os.Exit(1)
		}
	}()

	return nil
}

func (p *program) Stop() error {
	p.once.Do(func() {
		p.consumer.Exit()
	})
	return nil
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
