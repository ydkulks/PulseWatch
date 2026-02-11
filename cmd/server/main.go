package main

import (
	"fmt"
	"log"

	"net"

	"github.com/ydkulks/PulseWatch/internal/config"
	"github.com/ydkulks/PulseWatch/internal/server"
	pb "github.com/ydkulks/PulseWatch/pkg/v1/pulsewatch"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := config.GetServerPort()

	// Create gRPC server with injected dependencies
	grpcServer := server.NewGRPCServer(cfg)

	// Start listening
	listen, err := net.Listen("tcp", cfg.ServerPort)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", cfg.ServerPort, err)
	}
	fmt.Printf("PulseWatch server listening on %s\n", cfg.ServerPort)

	// Create gRPC server instance and register service
	pulseWatchService := grpc.NewServer()
	pb.RegisterPulseWatchServer(pulseWatchService, grpcServer)

	// Start the gRPC server
	err = pulseWatchService.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
