package log

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gen "http-service/gen/logger"
	"io"
	"log"
	"net/http"
	"time"
)

func CreateLogClient() *LogClient {
	conn, err := grpc.NewClient("Localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("failed to connect to log service: %v", err)
	}

	client := &LogClient{
		LoggerClient: gen.NewLoggerClient(conn),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	testEntry := &gen.LogEntry{
		TimestampSend: time.Now().UnixMilli(),
		Message:       "Test handshake log",
		ServiceName:   "HTTP-service",
		Level:         "DEBUG",
	}

	resp, err := client.LoggerClient.HandleIncomingLog(ctx, testEntry)

	if err != nil {
		log.Fatalf("log service connected, but test message failed: %v", err)
	}
	log.Printf("Handshake successful, log ID: %v", resp.GetId())

	return client

}

func (c *LogClient) LogDataGRPC(r *http.Request) (id string, err error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read request body: %w", err)
	}

	entry := &gen.LogEntry{
		TimestampSend: time.Now().UnixMilli(),
		Message:       fmt.Sprintf("Incoming request: %s %s\nBody: %s", r.Method, r.URL.Path, string(bodyBytes)),
		ServiceName:   "business-service",
		Level:         "INFO",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := c.LoggerClient.HandleIncomingLog(ctx, entry)

	if err != nil {
		return "", fmt.Errorf("failed to send log to gRPC service: %w", err)
	}

	return resp.Id.GetId(), nil
}

func (c *LogClient) ReadLogGRPC(id string, filename string) (*gen.LogReadingResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	readLogResponse, err := c.LoggerClient.ReadLog(ctx, &gen.LogInfo{
		Id:       id,
		Filename: filename,
	})

	fmt.Println(readLogResponse)

	if err != nil {
		return nil, fmt.Errorf("failed to call ReadLog: %w", err)
	}

	return readLogResponse, err

}

func (c *LogClient) DeleteLogGRPC(id string, filename string) (*gen.LogDeletionResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	deleteResponse, err := c.LoggerClient.DeleteLog(ctx, &gen.LogInfo{
		Id:       id,
		Filename: filename,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to call ReadLog: %w", err)

	}

	return deleteResponse, err

}
