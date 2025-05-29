package business

import (
	"context"
	"fmt"
	gen "http-service/gen"
	"time"
)

func (c *BusinessClient) Process(ctx context.Context, req *gen.OperationRequest) (*gen.OperationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()

	resp, err := c.GRPCClient.Process(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Process: %w", err)
	}

	return resp, nil
}
