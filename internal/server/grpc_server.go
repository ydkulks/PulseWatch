package server

import (
	"context"
	"fmt"
	"log"

	"github.com/ydkulks/PulseWatch/internal/config"
	"github.com/ydkulks/PulseWatch/internal/repository"
	"github.com/ydkulks/PulseWatch/internal/service"
	pb "github.com/ydkulks/PulseWatch/pkg/v1/pulsewatch"
)

type GRPCServer struct {
	pb.UnimplementedPulseWatchServer
	config  *config.Config
	service service.PulseWatchService
}

func NewGRPCServer(cfg *config.Config) *GRPCServer {
	repo := repository.NewProcessRepository()
	svc := service.NewPulseWatchService(cfg, repo)

	return &GRPCServer{
		config:  cfg,
		service: svc,
	}
}

func (s *GRPCServer) GetPulse(ctx context.Context, req *pb.GetPulseRequest) (*pb.GetPulseResponse, error) {
	log.Printf("GetPulse request received: name=%s", req.Name)

	response, err := s.service.GetPulse(ctx, req)
	if err != nil {
		log.Printf("Error in GetPulse: %v", err)
		return nil, fmt.Errorf("internal server error: %w", err)
	}

	log.Printf("GetPulse response sent: %s", response.Message)
	return response, nil
}

func (s *GRPCServer) WatchProcess(req *pb.WatchProcessRequest, stream pb.PulseWatch_WatchProcessServer) error {
	log.Printf("WatchProcess request received: pid=%d", req.Pid)

	watcher, err := s.service.WatchProcess(stream.Context(), req)
	if err != nil {
		log.Printf("Error creating process watcher: %v", err)
		return fmt.Errorf("failed to create process watcher: %w", err)
	}

	responseChan := watcher.Start(stream.Context(), req.Pid)

	for response := range responseChan {
		if response == nil {
			log.Printf("Process %d stopped", req.Pid)
			return nil
		}

		if err := stream.Send(response); err != nil {
			log.Printf("Error sending stream response: %v", err)
			return fmt.Errorf("failed to send stream response: %w", err)
		}
	}

	log.Printf("WatchProcess stream ended for pid=%d", req.Pid)
	return nil
}
