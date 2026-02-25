package main

import (
	"fmt"
	"log"

	"net"

	"github.com/ydkulks/PulseWatch/internal/config"
	"github.com/ydkulks/PulseWatch/internal/logger"
	"github.com/ydkulks/PulseWatch/internal/transport/grpc"
	proto "github.com/ydkulks/PulseWatch/proto/v1/pulsewatch"
	"google.golang.org/grpc"
)

func main() {
	logFile := logger.Init()
	defer logger.Close(logFile)

	cfg := config.GetServerPort()

	// Start listening with http server
	listen, err := net.Listen("tcp", cfg.ServerPort)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", cfg.ServerPort, err)
	}
	fmt.Printf("PulseWatch server listening on %s\n", cfg.ServerPort)

	// Create PulseWatch gRPC server with injected dependencies
	pulseWatchServer := server.NewGRPCServer(cfg)

	// Create gRPC server instance and register service
	grpcServer := grpc.NewServer()
	proto.RegisterPulseWatchServer(grpcServer, pulseWatchServer)

	// Start the gRPC server
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
