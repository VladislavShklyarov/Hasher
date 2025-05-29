package main

import (
	"business-service/internal/config"
	"business-service/internal/server"
	"business-service/internal/signals"
	"context"
	"log"
)

func main() {
	cfg := config.Load()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server.RunBusinessServer(cfg)

	signals.WaitForShutdown(ctx, cancel)

	log.Println("Shutting down server...")

}
