package app

import (
	"http-service/internal/client/grpc/business"
	logClient "http-service/internal/client/grpc/log"
)

// Dependency Inversion Principle

type Clients struct {
	LogClient      logClient.LogClientInterface
	BusinessClient *business.BusinessClient
}
