package main

import (
	"context"
	"fmt"
	"log-service/internal/config"
	"log-service/internal/logger/server"
	"log-service/internal/signals"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()
	fmt.Println("Route to logs: ", cfg.LogsDir)

	server.RunLogServer(cfg)

	signals.WaitForShutdown(ctx, cancel)
}
