package main

import (
	"context"
	"dashboard-service/internal/config"
	"dashboard-service/internal/kafka"
	"dashboard-service/internal/kafka/consumer"
	"dashboard-service/internal/signals"
	wscd "dashboard-service/internal/ws"
	"github.com/gorilla/websocket"
	"log"
)

func main() {

	cfg := config.Load()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kafkaConn, err := kafka.Connect(ctx, cfg)
	if err != nil {
		log.Println(err)
	}
	defer kafkaConn.Close()
	clients := &wscd.Clients{Clients: make(map[*websocket.Conn]bool)}

	go wscd.StartWebSocket(clients, cfg)
	go consumer.StartAll(clients, ctx, cfg.KafkaAddr)

	log.Println("dashboard-service is running...")

	signals.WaitForShutdown(ctx, cancel)

}
