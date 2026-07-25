package collector

import (
	"sort"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

type ProcessSnapshot struct {
	Name      string
	PID       int32
	CPU       float64
	MemoryMB  float64
	Timestamp time.Time
}

func Collect() ([]ProcessSnapshot, error) {
	allProcesses, err := process.Processes()
	if err != nil {
		return nil, err
	}

	snapshots := make([]ProcessSnapshot, 0, len(allProcesses))

	for _, proc := range allProcesses {
		name, err := proc.Name()
		if err != nil || name == "" {
			continue
		}

		cpu, err := proc.CPUPercent()
		if err != nil {
			cpu = 0
		}

		pid := proc.Pid

		memInfo, err := proc.MemoryInfo()
		memMB := 0.0
		if err == nil && memInfo != nil {
			memMB = float64(memInfo.RSS) / 1024 / 1024
		}

		snapshots = append(snapshots, ProcessSnapshot{
			Name:      name,
			PID:       pid,
			CPU:       cpu,
			MemoryMB:  memMB,
			Timestamp: time.Now(),
		})
	}

	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].CPU > snapshots[j].CPU
	})

	return snapshots, nil
}
