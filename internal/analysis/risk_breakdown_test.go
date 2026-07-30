package analysis

import (
	"testing"

	"syspulse/internal/collector"
	"syspulse/internal/storage"
)

func TestRiskBreakdownTracksMainRiskSources(t *testing.T) {
	baseline := storage.BaselineSnapshot{StartupCount: 1, ProcessCount: 2}
	startup := []collector.StartupEntry{
		{Name: "bad", Command: "powershell -enc ...", Source: "Startup Folder"},
		{Name: "safe", Command: "C:/Program Files/Good/app.exe", Source: "HKCU Run"},
	}
	snapshots := []collector.ProcessSnapshot{
		{Name: "chrome", PID: 1, CPU: 40, MemoryMB: 500},
		{Name: "worker", PID: 2, CPU: 30, MemoryMB: 300},
		{Name: "worker2", PID: 3, CPU: 20, MemoryMB: 250},
	}

	breakdown := ComputeRiskBreakdown(baseline, startup, snapshots)
	if breakdown.Total <= 0 {
		t.Fatal("expected some risk contribution in breakdown")
	}
	if breakdown.SuspiciousStartup <= 0 {
		t.Fatal("expected suspicious startup penalty")
	}
}

func TestRiskBreakdownSummaryReturnsNoRiskWhenClean(t *testing.T) {
	summary := RiskBreakdownSummary(RiskBreakdown{})
	if summary != "No significant risk factors detected" {
		t.Fatalf("unexpected summary: %q", summary)
	}
}
