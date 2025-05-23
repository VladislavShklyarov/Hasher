package server

import (
	gen "business-service/gen/business"
	grpcClient "business-service/internal/client/grpc"
	"google.golang.org/grpc"
	"log"
	"net"
)

func RunBusinessServer() {
	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Starting business logic service on :9091")

	logClient := grpcClient.CreateLogClient()

	if err := StartGRPCServer(lis, newBusinessLogicManager(logClient)); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func StartGRPCServer(listener net.Listener, businessServer gen.BusinessLogicServer) error {
	s := grpc.NewServer()
	gen.RegisterBusinessLogicServer(s, businessServer)
	return s.Serve(listener)
}
