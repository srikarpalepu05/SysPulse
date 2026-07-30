package analysis

import (
	"fmt"
	"strings"

	"syspulse/internal/collector"
	"syspulse/internal/storage"
)

type AnomalyFinding struct {
	ProcessName string
	PID         int32
	Severity    string
	Reason      string
}

func DetectAnomalies(baseline storage.BaselineSnapshot, snapshots []collector.ProcessSnapshot) []AnomalyFinding {
	findings := make([]AnomalyFinding, 0)

	for _, snapshot := range snapshots {
		lower := strings.ToLower(snapshot.Name)
		if strings.Contains(lower, "powershell") || strings.Contains(lower, "cmd") || strings.Contains(lower, "mshta") || strings.Contains(lower, "rundll32") {
			if snapshot.CPU > 10 || snapshot.MemoryMB > 200 {
				findings = append(findings, AnomalyFinding{
					ProcessName: snapshot.Name,
					PID:         snapshot.PID,
					Severity:    "high",
					Reason:      "process matches a suspicious execution pattern and is consuming resources",
				})
				continue
			}
		}

		if snapshot.CPU >= 35 || snapshot.MemoryMB >= 900 {
			findings = append(findings, AnomalyFinding{
				ProcessName: snapshot.Name,
				PID:         snapshot.PID,
				Severity:    "medium",
				Reason:      fmt.Sprintf("process resource usage is unusually high for a background process (%0.1f%% CPU, %0.1f MB)", snapshot.CPU, snapshot.MemoryMB),
			})
		}
	}

	if baseline.ProcessCount > 0 && len(snapshots) > baseline.ProcessCount+10 {
		findings = append(findings, AnomalyFinding{
			ProcessName: "system",
			Severity:    "high",
			Reason:      fmt.Sprintf("process count is %d, above the baseline of %d", len(snapshots), baseline.ProcessCount),
		})
	}

	return findings
}
