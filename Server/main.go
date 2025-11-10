package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/ydkulks/PulseWatch/pulsewatch"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedPulseWatchServer
}

func (s *server) GetPulse(pulseContext context.Context, pulseReq *pb.PulseRequest) (*pb.PulseResponse, error) {
	log.Printf("Received: %v", pulseReq.Name)
	return &pb.PulseResponse{Message: "Hello " + pulseReq.Name},nil
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Server listening on port ", 8080)

	pulseServer := grpc.NewServer()
	pb.RegisterPulseWatchServer(pulseServer, &server{})

	if err := pulseServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve : %s", err)
	}
}

