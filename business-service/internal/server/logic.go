package server

import (
	gen "business-service/gen/business"
	genLog "business-service/gen/logger"
	"business-service/internal/client/grpc"
	"context"
	"fmt"
	"log"
	"time"
)

type BusinessLogicManager struct {
	gen.UnimplementedBusinessLogicServer
	logClient *grpc.LogClient
}

func (blm *BusinessLogicManager) Process(ctx context.Context, req *gen.Request) (*gen.Response, error) {
	log.Printf("Revieved request: %+v", req)

	vv := &gen.VariableValue{
		Var:   "x",
		Value: 5,
	}
	resp := &gen.Response{
		Items: []*gen.VariableValue{vv},
	}

	entry := &genLog.LogEntry{
		ServiceName:   "business-service",
		Level:         "INFO",
		Message:       fmt.Sprintf("Processed request with result: %+v", resp),
		Metadata:      nil,
		TimestampSend: time.Now().UnixMilli(),
	}

	responseLog, err := blm.logClient.LoggerClient.HandleIncomingLog(ctx, entry)
	if err != nil {
		fmt.Println("Something went wrong during logging: ", err.Error())
	}
	if responseLog != nil {
		fmt.Println("Got response from log-service: " + responseLog.String())
	}
	return resp, nil

}
