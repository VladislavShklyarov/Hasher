package log

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gen "http-service/gen"
	"log"
	"time"
)

type LogClient struct {
	LoggerClient gen.LoggerClient
}

//type LogicClientInterface interface {
//	Process(ctx context.Context, req *gen.R) (*gen.Response, error)
//}
//
//var _ LogClientInterface = (*LogClient)(nil) // compile-time check

func CreateLogClient() *LogClient {
	conn, err := grpc.NewClient("Localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("failed to connect to log server: %v", err)
	}

	client := &LogClient{
		LoggerClient: gen.NewLoggerClient(conn),
	}

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
		ServiceName: "HTTP-server",
		Level:       "DEBUG",
	}

	resp, err := client.LoggerClient.HandleIncomingLog(ctx, testEntry)

	if err != nil {
		log.Fatalf("log server connected, but test message failed: %v", err)
	}
	log.Printf("Log server handshake successful, log ID: %v", resp.GetId())

	return client

}
