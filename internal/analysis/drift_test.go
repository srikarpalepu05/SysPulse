package analysis
package analysis

import (
	"testing"
	"time"

	"syspulse/internal/collector"
	"syspulse/internal/storage"
)

func TestDetectDriftFlagsChangesFromBaseline(t *testing.T) {
	baseline := storage.BaselineSnapshot{
		GeneratedAt:  time.Now().Add(-1 * time.Hour),
		StartupCount: 3,
		ProcessCount: 25,
		Summary:      "baseline",
	}

	startup := []collector.StartupEntry{
		{Name: "App1", Command: "C:/Temp/app1.exe", Source: "Startup Folder"},
		{Name: "App2", Command: "C:/Temp/app2.exe", Source: "HKCU Run"},
		{Name: "App3", Command: "C:/Temp/app3.exe", Source: "HKLM Run"},
		{Name: "App4", Command: "C:/Temp/app4.exe", Source: "Startup Folder"},
	}

	snapshots := []collector.ProcessSnapshot{
		{Name: "chrome", PID: 101, CPU: 12, MemoryMB: 200},
		{Name: "code", PID: 102, CPU: 20, MemoryMB: 320},
		{Name: "discord", PID: 103, CPU: 30, MemoryMB: 400},
		{Name: "helper", PID: 104, CPU: 8, MemoryMB: 210},
		{Name: "helper2", PID: 105, CPU: 9, MemoryMB: 240},
		{Name: "helper3", PID: 106, CPU: 10, MemoryMB: 260},
	}

	findings := DetectDrift(baseline, startup, snapshots)
	if len(findings) == 0 {
		t.Fatal("expected drift findings for changed startup and process counts")
	}

	seen := false
	for _, finding := range findings {
		if finding.Type == "startup" || finding.Type == "process" {
			seen = true
		}
	}
	if !seen {
		t.Fatal("expected startup or process drift detection among findings")
	}
}

func TestDetectDriftReturnsEmptyWhenStateMatchesBaseline(t *testing.T) {
	baseline := storage.BaselineSnapshot{
		GeneratedAt:  time.Now().Add(-30 * time.Minute),
		StartupCount: 2,
		ProcessCount: 3,
		Summary:      "baseline",
	}

	startup := []collector.StartupEntry{
		{Name: "App1", Command: "C:/Temp/app1.exe", Source: "HKCU Run"},
		{Name: "App2", Command: "C:/Temp/app2.exe", Source: "Startup Folder"},
	}

	snapshots := []collector.ProcessSnapshot{
		{Name: "chrome", PID: 101, CPU: 12, MemoryMB: 200},
		{Name: "code", PID: 102, CPU: 18, MemoryMB: 260},
		{Name: "explorer", PID: 103, CPU: 10, MemoryMB: 150},
	}

	findings := DetectDrift(baseline, startup, snapshots)
	if len(findings) != 0 {
		t.Fatalf("expected no drift when state matches baseline, got %d findings", len(findings))
	}
}
