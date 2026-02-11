package repository

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
)

type ProcessRepository interface {
	ProcessExists(pid int32) (bool, error)
	GetProcessMetrics(pid int32) (*ProcessMetrics, error)
	GetOSMetrics() (*OSMetrics, error)
	GetSystemMemory() (*SystemMemory, error)
}

type processRepository struct{}

func NewProcessRepository() ProcessRepository {
	return &processRepository{}
}

type OSMetrics struct {
	Host          string
	CpuPercentage []float32
	Timestamp     time.Time
}

type ProcessMetrics struct {
	Name          string
	Status        []string
	MemoryInfo    uint64
	CpuPercentage float32
	Timestamp     time.Time
}

type SystemMemory struct {
	Total       uint64
	Used        uint64
	UsedPercent float32
	Available   uint64
	Cached      uint64
	Free        uint64
	Timestamp   time.Time
}

func (r *processRepository) ProcessExists(pid int32) (bool, error) {
	exists, err := process.PidExists(pid)
	if err != nil {
		return false, fmt.Errorf("failed to check process existence for PID %d: %w", pid, err)
	}
	return exists, nil
}

func (r *processRepository) GetProcessMetrics(pid int32) (*ProcessMetrics, error) {
	p, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("failed to create process object for PID %d: %w", pid, err)
	}

	status, err := p.Status()
	if err != nil {
		return nil, fmt.Errorf("failed to get process status for PID %d: %w", pid, err)
	}

	memoryInfo, err := p.MemoryInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info for PID %d: %w", pid, err)
	}

	name, err := p.Name()
	if err != nil {
		return nil, fmt.Errorf("failed to get process name for PID %d: %w", pid, err)
	}

	cpuPercentage, err := p.CPUPercent()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU percentage for PID %d: %w", pid, err)
	}

	return &ProcessMetrics{
		Name:          name,
		Status:        status,
		MemoryInfo:    memoryInfo.RSS,
		CpuPercentage: float32(cpuPercentage),
		Timestamp:     time.Now(),
	}, nil
}

func (r *processRepository) GetOSMetrics() (*OSMetrics, error) {
	hostUsers, err := host.Users()
	if err != nil {
		return nil, fmt.Errorf("failed to get host users: %w", err)
	}

	if len(hostUsers) == 0 {
		return nil, fmt.Errorf("no host users found")
	}

	interval := time.Second * 1
	cpus, err := cpu.Percent(interval, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU percent: %w", err)
	}

	cpuPercents := make([]float32, len(cpus))
	for i, v := range cpus {
		cpuPercents[i] = float32(v)
	}

	return &OSMetrics{
		Host:          hostUsers[0].User,
		CpuPercentage: cpuPercents,
		Timestamp:     time.Now(),
	}, nil
}

func (r *processRepository) GetSystemMemory() (*SystemMemory, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get virtual memory: %w", err)
	}

	return &SystemMemory{
		Total:       v.Total,
		Used:        v.Used,
		UsedPercent: float32(v.UsedPercent),
		Available:   v.Available,
		Cached:      v.Cached,
		Free:        v.Free,
		Timestamp:   time.Now(),
	}, nil
}
