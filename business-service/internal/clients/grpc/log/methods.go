package log

import (
	"business-service/gen"
	"context"
	"fmt"
	"time"
)

func (c *LogClient) LogDataGRPC(ctx context.Context, entry *gen.LogEntry) (*gen.LogID, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := c.LoggerClient.HandleIncomingLog(ctx, entry)
	if err != nil {
		return nil, fmt.Errorf("failed to send log to gRPC server: %w", err)
	}

	return resp.Id, nil
}
