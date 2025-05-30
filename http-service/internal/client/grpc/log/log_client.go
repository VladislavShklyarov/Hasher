package log

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gen "http-service/gen"
	"http-service/internal/config"
	"log"
	"time"
)

const connectionError = "rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp 127.0.0.1:9090: connect: connection refused\""

type LogClient struct {
	LoggerClient gen.LoggerClient
}

func CreateLogClient(cfg *config.Config) *LogClient {
	conn, err := grpc.NewClient(cfg.LoggerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("failed to connect to log server: %v", err)
		return nil
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
		fmt.Println(err.Error())
		fmt.Println()

		if err.Error() == connectionError {
			log.Println("Log server connection error")
			return nil
		} else {
			log.Printf("log server connected, but test message failed: %v", err)
			return nil
		}
	}
	log.Printf("Log server handshake successful, log ID: %v", resp.GetId())

	return client

}
