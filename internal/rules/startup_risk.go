package rules

import (
	"strings"

	"syspulse/internal/collector"
)

type StartupRisk struct {
	Name     string
	Source   string
	Command  string
	Score    int
	Severity string
	Reason   string
}

func ScoreStartupEntries(entries []collector.StartupEntry) []StartupRisk {
	risks := make([]StartupRisk, 0)

	for _, entry := range entries {
		score := 0
		reasons := make([]string, 0)
		lower := strings.ToLower(entry.Command)

		if entry.Source == "Startup Folder" {
			score += 10
			reasons = append(reasons, "startup folder entry")
		}

		if strings.Contains(lower, "powershell") || strings.Contains(lower, "cmd.exe") || strings.Contains(lower, "rundll32") || strings.Contains(lower, "regsvr32") {
			score += 30
			reasons = append(reasons, "launches script or system helper")
		}

		if strings.Contains(lower, "telemetry") || strings.Contains(lower, "update") || strings.Contains(lower, "installer") {
			score += 15
			reasons = append(reasons, "common telemetry or updater pattern")
		}

		if strings.Contains(lower, "msedge") || strings.Contains(lower, "chrome") || strings.Contains(lower, "firefox") {
			score += 10
			reasons = append(reasons, "browser-related startup activity")
		}

		if entry.Command == "" {
			score += 20
			reasons = append(reasons, "missing or incomplete command target")
		}

		if score == 0 {
			continue
		}

		severity := "low"
		switch {
		case score >= 60:
			severity = "high"
		case score >= 35:
			severity = "medium"
		default:
			severity = "low"
		}

		risks = append(risks, StartupRisk{
			Name:     entry.Name,
			Source:   entry.Source,
			Command:  entry.Command,
			Score:    score,
			Severity: severity,
			Reason:   strings.Join(reasons, "; "),
		})
	}

	return risks
}
