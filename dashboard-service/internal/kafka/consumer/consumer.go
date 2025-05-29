package consumer

import (
	"context"
	wscd "dashboard-service/internal/ws"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaMessageHandler interface {
	Handle(msg []byte) []byte
	Broadcast(data []byte, clients *wscd.Clients)
}

func StartAll(clients *wscd.Clients, ctx context.Context, broker string) {
	StartConsumer(ctx, broker, "alg_graph_pic", "dashbord-service", &BizHandler{}, clients)
	StartConsumer(ctx, broker, "operation_log", "dashbord-service", &LogHandler{}, clients)
}

func StartConsumer(
	ctx context.Context,
	broker,
	topic string,
	groupID string,
	handler KafkaMessageHandler,
	clients *wscd.Clients) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	go func() {
		defer r.Close()
		for {
			m, err := r.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					log.Printf("%s: context canceled, shutting down", topic)
				} else {
					log.Printf("%s: error: %v", topic, err)
				}
				break
			}
			data := handler.Handle(m.Value)
			if data != nil {
				handler.Broadcast(data, clients)
			}
		}
	}()
}
