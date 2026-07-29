package analysis

import (
	"fmt"

	"syspulse/internal/collector"
	"syspulse/internal/storage"
)

type DriftFinding struct {
	Type     string
	Message  string
	Severity string
	Count    int
	Baseline int
	Current  int
}

func DetectDrift(baseline storage.BaselineSnapshot, startupEntries []collector.StartupEntry, snapshots []collector.ProcessSnapshot) []DriftFinding {
	findings := make([]DriftFinding, 0)

	startupDelta := len(startupEntries) - baseline.StartupCount
	if startupDelta != 0 {
		severity := "low"
		switch {
		case startupDelta >= 4 || startupDelta <= -4:
			severity = "high"
		case startupDelta >= 2 || startupDelta <= -2:
			severity = "medium"
		}

		findings = append(findings, DriftFinding{
			Type:     "startup",
			Message:  fmt.Sprintf("startup count changed from %d to %d", baseline.StartupCount, len(startupEntries)),
			Severity: severity,
			Baseline: baseline.StartupCount,
			Current:  len(startupEntries),
			Count:    startupDelta,
		})
	}

	processDelta := len(snapshots) - baseline.ProcessCount
	if processDelta != 0 {
		severity := "low"
		switch {
		case processDelta >= 10 || processDelta <= -10:
			severity = "high"
		case processDelta >= 5 || processDelta <= -5:
			severity = "medium"
		}

		findings = append(findings, DriftFinding{
			Type:     "process",
			Message:  fmt.Sprintf("active process count changed from %d to %d", baseline.ProcessCount, len(snapshots)),
			Severity: severity,
			Baseline: baseline.ProcessCount,
			Current:  len(snapshots),
			Count:    processDelta,
		})
	}

	return findings
}
