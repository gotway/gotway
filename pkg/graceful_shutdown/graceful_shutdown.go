package gracefulshutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type Stoppable interface{ Stop() }

func GracefulShutdown(
	cancel context.CancelFunc,
	stoppables ...Stoppable,
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
	for _, s := range stoppables {
		s.Stop()
	}
}
