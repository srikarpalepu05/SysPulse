package analysis

import (
	"testing"

	"syspulse/internal/collector"
	"syspulse/internal/storage"
)

func TestHygieneScoreReturnsLowScoreOnRiskyState(t *testing.T) {
	baseline := storage.BaselineSnapshot{StartupCount: 2, ProcessCount: 5}
	startup := []collector.StartupEntry{
		{Name: "bad1", Command: "powershell -enc ...", Source: "Startup Folder"},
		{Name: "bad2", Command: "cmd.exe /c install", Source: "HKCU Run"},
	}
	snapshots := []collector.ProcessSnapshot{
		{Name: "chrome", PID: 1, CPU: 25, MemoryMB: 500},
		{Name: "svchost", PID: 2, CPU: 50, MemoryMB: 900},
		{Name: "worker", PID: 3, CPU: 30, MemoryMB: 400},
		{Name: "worker2", PID: 4, CPU: 30, MemoryMB: 420},
		{Name: "worker3", PID: 5, CPU: 30, MemoryMB: 430},
		{Name: "worker4", PID: 6, CPU: 30, MemoryMB: 470},
	}

	score := ComputeHygieneScore(baseline, startup, snapshots)
	if score.Score >= 70 {
		t.Fatalf("expected low hygiene score for risky state, got %d", score.Score)
	}
	if score.Status == "healthy" {
		t.Fatal("expected unhealthy status for risky state")
	}
}

func TestHygieneScoreReturnsHealthyStateForCleanBaseline(t *testing.T) {
	baseline := storage.BaselineSnapshot{StartupCount: 2, ProcessCount: 3}
	startup := []collector.StartupEntry{
		{Name: "App1", Command: "C:/Program Files/App1/app.exe", Source: "HKCU Run"},
		{Name: "App2", Command: "C:/Program Files/App2/app.exe", Source: "Startup Folder"},
	}
	snapshots := []collector.ProcessSnapshot{
		{Name: "chrome", PID: 1, CPU: 12, MemoryMB: 200},
		{Name: "code", PID: 2, CPU: 18, MemoryMB: 260},
		{Name: "explorer", PID: 3, CPU: 10, MemoryMB: 150},
	}

	score := ComputeHygieneScore(baseline, startup, snapshots)
	if score.Score < 75 {
		t.Fatalf("expected healthy hygiene score for clean state, got %d", score.Score)
	}
	if score.Status != "healthy" {
		t.Fatalf("expected healthy status, got %s", score.Status)
	}
}

func TestSecurityHygieneSummaryBuildsFriendlyStatus(t *testing.T) {
	summary := SecurityHygieneSummary(82, "healthy")
	if summary == "" {
		t.Fatal("expected non-empty summary")
	}
	if summary != "System hygiene is healthy (score: 82/100)" {
		t.Fatalf("unexpected summary value: %q", summary)
	}
}
