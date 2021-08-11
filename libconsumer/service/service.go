package service

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/JieTrancender/nsq_consumer/libconsumer/logp"
)

// HandleSignals manages OS signals that ask the service/daemon to stop.
// The stopFunction should break the loop in the Consumer so that the service shut downs gracefully
func HandleSignals(stopFunction func(), cancel context.CancelFunc) {
	var callback sync.Once
	logger := logp.NewLogger("service")

	// On termination signals, gracefully stop the Consumer
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		sig := <-sigc

		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			logger.Debug("Received sigterm/sigint, stopping")
		case syscall.SIGHUP:
			logger.Debug("Received sighup, stopping")
		}

		cancel()
		callback.Do(stopFunction)
	}()

	// Handle the Windows service events
	go ProcessWindowsControlEvents(func() {
		logger.Debug("Received svc stop/shutdown request")
		callback.Do(stopFunction)
	})
}
