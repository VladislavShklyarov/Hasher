package log

import (
	"business-service/gen"
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
		log.Fatalf("failed to connect to log server: %v", err)
	}

	client := &LogClient{
		LoggerClient: gen.NewLoggerClient(conn),
	}

	resp, err := LogHandshake(client)

	if err != nil {
		log.Fatalf("log server connected, but test message failed: %v", err)
	}
	log.Printf("Handshake successful, log ID: %v", resp.GetId())

	return client

}

func LogHandshake(client *LogClient) (*gen.LogCreationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	testEntry := &gen.LogEntry{
		TimestampSend: time.Now().UnixMilli(),
		Message: &gen.StructuredMessage{
			Method: "POST",
			Path:   "/process",
			Body: []*gen.Operation{
				{
					Type:  "Test type",
					Op:    "Test operation",
					Var:   "Test variable",
					Left:  "Test Left",
					Right: "Test Right",
				},
			},
			Result: &gen.OperationResponse{
				Items: []*gen.VariableValue{
					{
						Var:   "Test var",
						Value: 999,
					},
				},
			},
		},
		Metadata: map[string]string{
			"test": "true",
		},
		ServiceName: "business-server",
		Level:       "DEBUG",
	}

	return client.LoggerClient.HandleIncomingLog(ctx, testEntry)
}
