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

	conn, err := dialer.DialContext(ctx, "tcp", cfg.KafkaAddr)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к Kafka: %w", err)
	}

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
		log.Fatalf("failed to create topics: %v", err)
	}

	fmt.Println("Topics created successfully")

	return conn, nil
}
