package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/ydkulks/PulseWatch/pkg/v1/pulsewatch"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedPulseWatchServer
}

func (s *server) GetPulse(ctx context.Context, req *pb.GetPulseRequest) (*pb.GetPulseResponse, error) {
	return &pb.GetPulseResponse{Message: "Hello " + req.Name}, nil
}

func main() {
	// Create a HTTP listen
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	fmt.Println("Server listening on port ", listen.Addr().String())

	// Create a gRPC server
	pulseWatchService := grpc.NewServer()
	// Register the service
	pb.RegisterPulseWatchServer(pulseWatchService, &server{})

	// Start the gRPC server
	err = pulseWatchService.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
