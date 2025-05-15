package server

import (
	"business-service/internal/client/grpc"
)

func newBusinessLogicManager(logClient *grpc.LogClient) *BusinessLogicManager {
	return &BusinessLogicManager{
		logClient: logClient,
	}
}
