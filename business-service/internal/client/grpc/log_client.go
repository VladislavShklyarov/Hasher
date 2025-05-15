package grpc

import (
	gen "business-service/gen/logger"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type LogClient struct {
	LoggerClient gen.LoggerClient
}

func CreateLogClient() *LogClient {
	conn, err := grpc.NewClient("Localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("failed to connect to log service: %v", err)
	}

	client := &LogClient{
		LoggerClient: gen.NewLoggerClient(conn),
	}

	resp, err := LogHandshake(client)

	if err != nil {
		log.Fatalf("log service connected, but test message failed: %v", err)
	}
	log.Printf("Handshake successful, log ID: %v", resp.GetId())

	return client

}

func LogHandshake(client *LogClient) (*gen.LogCreationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	testEntry := &gen.LogEntry{
		TimestampSend: time.Now().UnixMilli(),
		Message:       "Test handshake log",
		ServiceName:   "business-service",
		Level:         "DEBUG",
	}

	return client.LoggerClient.HandleIncomingLog(ctx, testEntry)
}
