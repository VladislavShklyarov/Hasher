package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func WaitForShutdown(ctx context.Context, cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-sig:
		println("Получен сигнал:", s.String())
		cancel()
	case <-ctx.Done():
		return
	}
}
