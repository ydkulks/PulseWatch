package main

import (
	"context"
	"flag"
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
	log.Printf("Response : %s", response.GetMessage())
}
