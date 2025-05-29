package main

import (
	"context"
	"log-service/internal/config"
	"log-service/internal/logger/server"
	"log-service/internal/signals"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()

	server.RunLogServer(cfg)

	signals.WaitForShutdown(ctx, cancel)
}
