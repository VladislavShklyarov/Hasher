package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"log"
	"log-service/gen"
	"log-service/internal/config"
	"log-service/internal/logger/CRUD"
	"strings"
	"time"
)

func publishOperationResult(broker, topic string, payload []byte) {
	writer := NewKafkaWriter(broker, topic)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("operation"),
		Value: payload,
	})

	if err != nil {
		log.Printf("Failed to publish operation: %v\n", err)
	} else {
		log.Println("Operation send successfully.")
	}

	writer.Close()

}

func StartKafka(ch chan *gen.LogEntry) {

	cfg := config.Load()

	file, err := CRUD.OpenFile("http_logs.json")
	if err != nil {
		fmt.Println("Failed to open file: ", err)
	}
	for msg := range ch {
		log, err := CRUD.FindLog(file, msg.Message.Path)
		if err != nil {
			fmt.Println("Failed to fing log: ", err)
		}

		operations, err := extractOperations(log)
		result := msg.Message.Result

		payload, err := proto.Marshal(&gen.StructuredMessage{
			Method: "POST",
			Path:   "from log-service",
			Body:   operations,
			Result: result,
		})
		publishOperationResult(cfg.KafkaBroker, cfg.KafkaTopic, payload)
	}
}

func extractBody(log string) string {
	start := strings.Index(log, `"body":[`)
	if start == -1 {
		fmt.Println("body not found")
		return ""
	}
	start += len(`"body":`)

	// Ищем закрывающую скобку после начала
	end := start
	brackets := 1 // уже встретили первую [
	for i := start + 1; i < len(log); i++ {
		switch log[i] {
		case '[':
			brackets++
		case ']':
			brackets--
			if brackets == 0 {
				end = i + 1
				break
			}
		}
	}

	bodyRaw := log[start:end]
	return bodyRaw
}

func extractOperations(log string) ([]*gen.Operation, error) {
	bodyJson := extractBody(log)
	if bodyJson == "" {
		return nil, fmt.Errorf("body not found")
	}

	var ops []*gen.Operation
	if err := json.Unmarshal([]byte(bodyJson), &ops); err != nil {
		return nil, err
	}
	return ops, nil
}
