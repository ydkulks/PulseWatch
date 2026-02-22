package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ydkulks/PulseWatch/internal/config"
	"github.com/ydkulks/PulseWatch/internal/repository"
	proto "github.com/ydkulks/PulseWatch/proto/v1/pulsewatch"
)

type PulseWatchService interface {
	GetPulse(ctx context.Context, req *proto.GetPulseRequest) (*proto.GetPulseResponse, error)
	WatchProcess(ctx context.Context, req *proto.WatchProcessRequest) (*ProcessWatcher, error)
}

type pulseWatchService struct {
	config *config.Config
	repo   repository.ProcessRepository
}

func NewPulseWatchService(cfg *config.Config, repo repository.ProcessRepository) PulseWatchService {
	return &pulseWatchService{
		config: cfg,
		repo:   repo,
	}
}

type ProcessWatcher struct {
	config *config.Config
	repo   repository.ProcessRepository
	ticker *time.Ticker
}

func (s *pulseWatchService) GetPulse(ctx context.Context, req *proto.GetPulseRequest) (*proto.GetPulseResponse, error) {
	return &proto.GetPulseResponse{
		Message: fmt.Sprintf("Hello %s - PulseWatch service is running", req.Name),
	}, nil
}

func (s *pulseWatchService) WatchProcess(ctx context.Context, req *proto.WatchProcessRequest) (*ProcessWatcher, error) {
	if req.Pid <= 0 {
		return nil, fmt.Errorf("invalid PID: %d", req.Pid)
	}

	exists, err := s.repo.ProcessExists(req.Pid)
	if err != nil {
		return nil, fmt.Errorf("failed to check process existence: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("process with PID %d does not exist", req.Pid)
	}

	return &ProcessWatcher{
		config: s.config,
		repo:   s.repo,
		ticker: time.NewTicker(5 * time.Second),
	}, nil
}

func (pw *ProcessWatcher) Start(ctx context.Context, pid int32) <-chan *proto.WatchProcessResponse {
	responseChan := make(chan *proto.WatchProcessResponse)

	go func() {
		defer close(responseChan)
		defer pw.ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-pw.ticker.C:
				response, err := pw.getProcessResponse(pid)
				if err != nil {
					return
				}

				select {
				case responseChan <- response:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return responseChan
}

func (pw *ProcessWatcher) getProcessResponse(pid int32) (*proto.WatchProcessResponse, error) {
	exists, err := pw.repo.ProcessExists(pid)
	if err != nil {
		return nil, fmt.Errorf("failed to check process existence: %w", err)
	}

	if !exists {
		return nil, nil
	}

	osMetrics, err := pw.repo.GetOSMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to get OS metrics: %w", err)
	}

	processMetrics, err := pw.repo.GetProcessMetrics(pid)
	if err != nil {
		return nil, fmt.Errorf("failed to get process metrics: %w", err)
	}

	return &proto.WatchProcessResponse{
		Status: exists,
		Pid:    pid,
		OsMetrics: &proto.OsMetrics{
			Host:          osMetrics.Host,
			CpuPercentage: osMetrics.CpuPercentage,
		},
		ProcessMetrics: &proto.ProcessMetrics{
			Name:          processMetrics.Name,
			Status:        processMetrics.Status,
			MemoryInfo:    processMetrics.MemoryInfo,
			CpuPercentage: processMetrics.CpuPercentage,
		},
	}, nil
}
