package gracefulshutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type ShutdownHook func()

func GracefulShutdown(
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
	<-signals
	cancel()
	for _, hook := range shutdownHooks {
		hook()
	}
}
