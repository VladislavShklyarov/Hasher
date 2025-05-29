package server

import (
	"business-service/gen"
	grpcClient "business-service/internal/clients/grpc/log"
	"business-service/internal/config"
	"google.golang.org/grpc"
	"log"
	"net"
)

func RunBusinessServer(cfg *config.Config) {
	lis, err := net.Listen("tcp", cfg.BusinessAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Starting business logic server on :%s", cfg.BusinessAddr)

	logClient := grpcClient.CreateLogClient(cfg)

	if err := StartGRPCServer(lis, newBusinessLogicManager(logClient)); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func StartGRPCServer(listener net.Listener, businessServer gen.BusinessLogicServer) error {
	s := grpc.NewServer()
	gen.RegisterBusinessLogicServer(s, businessServer)
	return s.Serve(listener)
}
