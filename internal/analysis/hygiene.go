package analysis

import (
	"strings"

	"syspulse/internal/collector"
	"syspulse/internal/storage"
)

type HygieneScore struct {
	Score  int
	Status string
	Reason string
}

func ComputeHygieneScore(baseline storage.BaselineSnapshot, startupEntries []collector.StartupEntry, snapshots []collector.ProcessSnapshot) HygieneScore {
	score := 100

	startupDelta := len(startupEntries) - baseline.StartupCount
	if startupDelta > 0 {
		score -= startupDelta * 12
	}
	if startupDelta < 0 {
		score += -startupDelta * 3
	}

	processDelta := len(snapshots) - baseline.ProcessCount
	if processDelta > 0 {
		score -= processDelta * 6
	}
	if processDelta < 0 {
		score += -processDelta * 2
	}

	for _, entry := range startupEntries {
		lower := strings.ToLower(entry.Command)
		if strings.Contains(lower, "powershell") || strings.Contains(lower, "cmd.exe") || strings.Contains(lower, "rundll32") || strings.Contains(lower, "regsvr32") {
			score -= 18
		}
	}

	for _, snapshot := range snapshots {
		if snapshot.CPU >= 35 || snapshot.MemoryMB >= 900 {
			score -= 10
		}
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	status := "healthy"
	reason := "No major system hygiene issues detected."
	if score < 70 {
		status = "warning"
		reason = "System drift or risky startup behavior detected."
	}
	if score < 40 {
		status = "critical"
		reason = "System is significantly outside the known-good baseline."
	}

	return HygieneScore{Score: score, Status: status, Reason: reason}
}
