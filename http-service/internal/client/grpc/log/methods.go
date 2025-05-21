package log

import (
	"context"
	"fmt"
	gen "http-service/gen"
	"time"
)

func (c *LogClient) LogDataGRPC(ctx context.Context, entry *gen.LogEntry) (*gen.LogID, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := c.LoggerClient.HandleIncomingLog(ctx, entry)
	if err != nil {
		return nil, fmt.Errorf("failed to send log to gRPC server: %w", err)
	}

	return resp.Id, nil
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
