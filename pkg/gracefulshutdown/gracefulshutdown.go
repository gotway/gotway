package gracefulshutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gotway/gotway/pkg/log"
)

type ShutdownHook func()

func GracefulShutdown(
	logger log.Logger,
	cancel context.CancelFunc,
	shutdownHooks ...ShutdownHook,
) {
	signals := make(chan os.Signal)
	signal.Notify(
		signals,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)
	s := <-signals
	logger.Info("received signal ", s.String())
	cancel()
	for _, hook := range shutdownHooks {
		hook()
	}
}
