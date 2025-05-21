package app

import (
	"context"
	"http-service/gen"
)

type BusinessClientInterface interface {
	Process(ctx context.Context, req *gen.OperationRequest) (*gen.OperationResponse, error)
}

type LogClientInterface interface {
	ReadLogGRPC(id, filename string) (*gen.LogReadingResponse, error)
	DeleteLogGRPC(id, filename string) (*gen.LogDeletionResponse, error)
	LogDataGRPC(ctx context.Context, entry *gen.LogEntry) (id *gen.LogID, err error)
}

// Dependency Inversion Principle

type Clients struct {
	LogClient      LogClientInterface
	BusinessClient BusinessClientInterface
}
