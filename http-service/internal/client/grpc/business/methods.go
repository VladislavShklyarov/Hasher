package business

import (
	"context"
	"fmt"
	gen "http-service/gen/business"
	"time"
)

func (c *BusinessClient) Process(ctx context.Context, req *gen.Request) (*gen.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.BusinessLogicClient.Process(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Process: %w", err)
	}

	return resp, nil
}
