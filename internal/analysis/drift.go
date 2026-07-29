package analysis
package analysis

import (
	"fmt"
	"strings"

	"syspulse/internal/collector"
	"syspulse/internal/storage"
)

type DriftFinding struct {
	Type      string
	Message   string
	Severity  string
	Count     int
	Baseline  int
	Current   int
}

func DetectDrift(baseline storage.BaselineSnapshot, startupEntries []collector.StartupEntry, snapshots []collector.ProcessSnapshot) []DriftFinding {
	findings := make([]DriftFinding, 0)

	startupDelta := len(startupEntries) - baseline.StartupCount
	if startupDelta != 0 {
		severity := "low"
		if startupDelta >= 2 || startupDelta <= -2 {
			severity = "medium"
		}
		if startupDelta >= 4 || startupDelta <= -4 {
			severity = "high"
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
		if processDelta >= 5 || processDelta <= -5 {
			severity = "medium"
		}
		if processDelta >= 10 || processDelta <= -10 {
			severity = "high"
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

	startupNames := make(map[string]bool)
	for _, entry := range startupEntries {
		startupNames[strings.ToLower(strings.TrimSpace(entry.Name))] = true
	}

	if len(startupEntries) > baseline.StartupCount {
		newEntries := 0
		for _, entry := range startupEntries {
			name := strings.ToLower(strings.TrimSpace(entry.Name))
			if name == "" {
				continue
			}
			if !strings.Contains(strings.ToLower(baseline.Summary), name) {
				newEntries++
			}
		}
		if newEntries > 0 {
			findings = append(findings, DriftFinding{
				Type:     "startup",
				Message:  fmt.Sprintf("%d new startup entries detected compared to the baseline", newEntries),
				Severity: "medium",
				Count:    newEntries,
			})
		}
	}

	return findings
}
