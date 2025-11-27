package server

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
)

type Process struct {
	Pid int32
}

func (p Process) ProcessExists() (bool, error) {
	exists, err := process.PidExists(p.Pid)
	if err != nil {
		return false, errors.New("Error checking process existence")
	}
	return exists, nil
}

func GetMemInfo() error {
	v, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("failed to get virtual memory: %v", err)
	}
	log.Println("Memory Usage")
	log.Println("Mem - Total: ", v.Total)
	log.Println("Mem - Used: ", v.Used)
	log.Println("Mem - UsedPercent: ", v.UsedPercent)
	log.Println("Mem - Available: ", v.Available)
	log.Println("Mem - Cached: ", v.Cached)
	log.Println("Mem - Free: ", v.Free)
	return nil
}

type OSMetrics struct {
	Host          string
	CpuPercentage []float32
}

func (p Process) GetOSMetrics() (OSMetrics, error) {
	hostUsers, err := host.Users()
	if err != nil {
		return OSMetrics{}, fmt.Errorf("failed to get host users: %v", err)
	}
	if len(hostUsers) == 0 {
		return OSMetrics{}, fmt.Errorf("no host users found")
	}
	var interval time.Duration = time.Second * 1
	cpus, err := cpu.Percent(interval, false)
	if err != nil {
		return OSMetrics{}, fmt.Errorf("failed to get CPU percent: %v", err)
	}

	// Typecasting float64 to float32
	cpuPercents := make([]float32, len(cpus))
	for i, v := range cpus {
		cpuPercents[i] = float32(v)
	}
	return OSMetrics{
		Host:          hostUsers[0].User,
		CpuPercentage: cpuPercents,
	}, nil
}

type ProcessMetrics struct {
	Name          string
	Status        []string
	MemoryInfo    uint64
	CpuPercentage float32
}

func (p Process) GetProcessMetrics() (ProcessMetrics, error) {
	// Get process info
	currentProcess, err := process.NewProcess(p.Pid)
	if err != nil {
		return ProcessMetrics{}, fmt.Errorf("process not found: %v", err)
	}
	// Get process status
	status, err := currentProcess.Status()
	if err != nil {
		return ProcessMetrics{}, fmt.Errorf("failed to get process status: %v", err)
	}
	memoryInfo, err := currentProcess.MemoryInfo()
	if err != nil {
		return ProcessMetrics{}, fmt.Errorf("failed to get memory info: %v", err)
	}
	name, err := currentProcess.Name()
	if err != nil {
		return ProcessMetrics{}, fmt.Errorf("failed to get process name: %v", err)
	}
	cpuPercentage, err := currentProcess.CPUPercent()
	if err != nil {
		return ProcessMetrics{}, fmt.Errorf("failed to get CPU percentage: %v", err)
	}

	// Typecasting float64 to float32
	cpuPercent := float32(cpuPercentage)
	return ProcessMetrics{
		Name:          name,
		Status:        status,
		MemoryInfo:    memoryInfo.RSS,
		CpuPercentage: cpuPercent,
	}, nil
}
