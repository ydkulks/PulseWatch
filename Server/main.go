package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	pb "github.com/ydkulks/PulseWatch/pulsewatch"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedPulseWatchServer
}

func (s *server) GetPulse(pulseContext context.Context, pulseReq *pb.PulseRequest) (*pb.PulseResponse, error) {
	log.Printf("Received: %v", pulseReq.Name)
	return &pb.PulseResponse{Message: "Hello " + pulseReq.Name}, nil
}

func (s *server) ServerStreamPulse(pulseReq *pb.PulseRequest, stream grpc.ServerStreamingServer[pb.PulseResponse]) error {
	for i := 0; i < 10; i++ {
		log.Printf("Received: %v", pulseReq.Name)
		message := fmt.Sprintf("Hello %s, count: %d", pulseReq.Name, i)
		if err := stream.Send(&pb.PulseResponse{Message: message}); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) ClientStreamPulse(stream grpc.ClientStreamingServer[pb.PulseRequest, pb.PulseResponse]) error {
	var names []string
	for {
		pulseReq, err := stream.Recv()
		if err == io.EOF {
			message := fmt.Sprintf("Goodbye %s", strings.Join(names, ", "))
			return stream.SendAndClose(&pb.PulseResponse{Message: message})
		}
		if err != nil {
			return err
		}
		names = append(names, pulseReq.Name)
		log.Printf("Received: %v", pulseReq.Name)
	}
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Server listening on port ", listen.Addr().String())

	pulseServer := grpc.NewServer()
	pb.RegisterPulseWatchServer(pulseServer, &server{})

	if err := pulseServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve : %s", err)
	}
}
