package server

import (
	"business-service/internal/client/grpc/log"
	blm "business-service/internal/server/handlers"
)

func newBusinessLogicManager(logClient *log.LogClient) *blm.BusinessLogicManager {
	return &blm.BusinessLogicManager{
		GRPCClient: logClient,
	}
}
