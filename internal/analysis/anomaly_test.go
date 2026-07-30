package analysis

import (
	"testing"

	"syspulse/internal/collector"
	"syspulse/internal/storage"
)

func TestDetectAnomaliesFlagsSuspiciousProcessUsage(t *testing.T) {
	baseline := storage.BaselineSnapshot{ProcessCount: 4}
	snapshots := []collector.ProcessSnapshot{
		{Name: "powershell", PID: 301, CPU: 40, MemoryMB: 600},
		{Name: "chrome", PID: 302, CPU: 12, MemoryMB: 180},
	}

	findings := DetectAnomalies(baseline, snapshots)
	if len(findings) == 0 {
		t.Fatal("expected anomaly findings")
	}
	if findings[0].Severity != "high" {
		t.Fatalf("expected high severity anomaly, got %s", findings[0].Severity)
	}
}

func TestDetectAnomaliesDoesNotFlagNormalProcessState(t *testing.T) {
	baseline := storage.BaselineSnapshot{ProcessCount: 4}
	snapshots := []collector.ProcessSnapshot{
		{Name: "chrome", PID: 101, CPU: 10, MemoryMB: 160},
		{Name: "code", PID: 102, CPU: 12, MemoryMB: 180},
		{Name: "explorer", PID: 103, CPU: 8, MemoryMB: 120},
	}

	findings := DetectAnomalies(baseline, snapshots)
	if len(findings) != 0 {
		t.Fatalf("expected no anomaly findings, got %d", len(findings))
	}
}
