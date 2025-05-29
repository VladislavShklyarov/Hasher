package kafka

import (
	"context"
	"encoding/base64"
	kafka "github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

func PublishAlgoGraph(broker, topic, path string) {
	writer := NewKafkaWriter(broker, topic)
	encoded, err := encodeFile(path)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("image"),
		Value: []byte(encoded),
	})

	if err != nil {
		log.Printf("Ошибка при отправке картинки: %v\n", err)
	} else {
		log.Println("Картинка успешно отправлено")
	}

	writer.Close()

}

func encodeFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Ошибка чтения файла: %v\n", err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), err
}
