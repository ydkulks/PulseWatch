package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"net"

	pulse "github.com/ydkulks/PulseWatch/internal/server"
	pb "github.com/ydkulks/PulseWatch/pkg/v1/pulsewatch"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedPulseWatchServer
}

func (s *server) GetPulse(ctx context.Context, req *pb.GetPulseRequest) (*pb.GetPulseResponse, error) {
	return &pb.GetPulseResponse{Message: "Hello " + req.Name}, nil
}

func (s *server) WatchProcess(req *pb.WatchProcessRequest, res grpc.ServerStreamingServer[pb.WatchProcessResponse]) error {
	process := pulse.Process{Pid: req.Pid}
	log.Println("Process ID: ", process)
	for {
		exists, err := process.ProcessExists()
		log.Println("Process Status: ", exists)
		if err != nil {
			log.Printf("Error checking process existence: %v", err)
			return err
		}

		if !exists {
			return nil
		}

		osMetrics, err := process.GetOSMetrics()
		if err != nil {
			log.Printf("Error getting OS metrics: %v", err)
			return err
		}
		processMetrics, err := process.GetProcessMetrics()
		if err != nil {
			log.Printf("Error getting process metrics: %v", err)
			return err
		}
		if err := res.Send(&pb.WatchProcessResponse{
			Status: exists,
			Pid:    req.Pid,
			OsMetrics: &pb.OsMetrics{
				Host:          osMetrics.Host,
				CpuPercentage: osMetrics.CpuPercentage,
			},
			ProcessMetrics: &pb.ProcessMetrics{
				Name:          processMetrics.Name,
				Status:        processMetrics.Status,
				MemoryInfo:    processMetrics.MemoryInfo,
				CpuPercentage: processMetrics.CpuPercentage,
			},
		}); err != nil {
			return err
		}
		time.Sleep(time.Second * 5)
	}
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
