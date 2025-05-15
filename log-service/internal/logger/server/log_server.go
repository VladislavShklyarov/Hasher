package server

import (
	"google.golang.org/grpc"
	"log"
	gen "log-service/gen/logger"
	"net"
)

func RunLogServer() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Starting log service on :9090")
	if err := StartGRPCServer(lis, NewLogManager()); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func StartGRPCServer(listener net.Listener, loggerServer gen.LoggerServer) error {
	s := grpc.NewServer()
	gen.RegisterLoggerServer(s, loggerServer)
	return s.Serve(listener)
}
