package kafka

import (
	"context"
	"dashboard-service/internal/config"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func Connect(ctx context.Context, cfg *config.Config) (*kafka.Conn, error) {
	dialer := kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	var conn *kafka.Conn
	var err error

	// Попытки подключения с retry
	for i := 0; i < 10; i++ {
		conn, err = dialer.DialContext(ctx, "tcp", cfg.KafkaBroker)
		if err == nil {
			break
		}
		log.Printf("Retry %d: could not connect to Kafka at %s (%s), retrying...\n", i+1, cfg.KafkaBroker, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("Kafka is not available after several retries: %w", err)
	}

	// Попытка создать топики
	topics := []kafka.TopicConfig{
		{
			Topic:             cfg.LogTopic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
		{
			Topic:             cfg.BizTopic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = conn.CreateTopics(topics...)
	if err != nil {
		// Не аварийно завершаем, Kafka может сама создать топики при публикации
		log.Printf("warning: failed to create topics (they may already exist): %v", err)
	} else {
		log.Println("Kafka topics created successfully")
	}

	return conn, nil
}
