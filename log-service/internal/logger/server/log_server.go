package server

import (
	"google.golang.org/grpc"
	"log"
	"log-service/gen"
	"log-service/internal/clients/kafka"
	"log-service/internal/config"
	"net"
)

func RunLogServer(cfg *config.Config) {
	lis, err := net.Listen("tcp", cfg.LoggerAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Starting log server on :%s", cfg.LoggerAddr)

	logManager := NewLogManager(cfg)

	go kafka.StartKafka(logManager.LogChanel)

	if err := StartGRPCServer(lis, logManager); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func StartGRPCServer(listener net.Listener, loggerServer gen.LoggerServer) error {
	s := grpc.NewServer()
	gen.RegisterLoggerServer(s, loggerServer)
	return s.Serve(listener)
}
