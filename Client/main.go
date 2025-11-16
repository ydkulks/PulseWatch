package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	pb "github.com/ydkulks/PulseWatch/pulsewatch"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var name = flag.String("name", "ydkulks", "Name to greet")

func main() {
	flag.Parse()
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %V", err)
	}
	defer conn.Close()

	client := pb.NewPulseWatchClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := client.GetPulse(ctx, &pb.PulseRequest{Name: *name})
	if err != nil {
		log.Fatalf("Failed to get pulse : %V", err)
	}
	log.Printf("Unary response : %s", response.GetMessage())

	// Server streaming
	log.Println("Starting server streaming...")
	stream, err := client.ServerStreamPulse(ctx, &pb.PulseRequest{Name: *name})
	if err != nil {
		log.Fatalf("Failed to call ServerStreamPulse: %v", err)
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive: %v", err)
		}
		log.Printf("Server stream response: %s", resp.GetMessage())
	}
	log.Println("Server streaming done.")

	// Client streaming
	log.Println("Starting client streaming...")
	stream2, err := client.ClientStreamPulse(ctx)
	if err != nil {
		log.Fatalf("Failed to call ClientStreamPulse: %v", err)
	}
	names := []string{"Alice", "Bob", "Charlie"}
	for _, n := range names {
		if err := stream2.Send(&pb.PulseRequest{Name: n}); err != nil {
			log.Fatalf("Failed to send: %v", err)
		}
	}
	resp2, err := stream2.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to receive: %v", err)
	}
	log.Printf("Client stream response: %s", resp2.GetMessage())
	log.Println("Client streaming done.")

	// Bidirectional streaming
	log.Println("Starting bidirectional streaming...")
	stream3, err := client.BidiStreamPulse(ctx)
	if err != nil {
		log.Fatalf("Failed to call BidiStreamPulse: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			resp, err := stream3.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive: %v", err)
			}
			log.Printf("Bidi response: %s", resp.GetMessage())
		}
	}()
	names2 := []string{"Dave", "Eve", "Frank"}
	for _, n := range names2 {
		if err := stream3.Send(&pb.PulseRequest{Name: n}); err != nil {
			log.Fatalf("Failed to send: %v", err)
		}
	}
	stream3.CloseSend()
	<-waitc
	log.Println("Bidirectional streaming done.")
}
