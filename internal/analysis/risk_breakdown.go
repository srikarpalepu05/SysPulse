package analysis

import (
	"fmt"
	"strings"

	"syspulse/internal/collector"
	"syspulse/internal/storage"
)

type RiskBreakdown struct {
	StartupRisk int
	ProcessRisk int
	SuspiciousStartup int
	Total int
}

func ComputeRiskBreakdown(baseline storage.BaselineSnapshot, startupEntries []collector.StartupEntry, snapshots []collector.ProcessSnapshot) RiskBreakdown {
	breakdown := RiskBreakdown{}
	startupDelta := len(startupEntries) - baseline.StartupCount
	if startupDelta > 0 {
		breakdown.StartupRisk = startupDelta * 12
	}

	processDelta := len(snapshots) - baseline.ProcessCount
	if processDelta > 0 {
		breakdown.ProcessRisk = processDelta * 6
	}

	for _, entry := range startupEntries {
		lower := strings.ToLower(entry.Command)
		if strings.Contains(lower, "powershell") || strings.Contains(lower, "cmd.exe") || strings.Contains(lower, "rundll32") || strings.Contains(lower, "regsvr32") {
			breakdown.SuspiciousStartup++
		}
	}
	breakdown.SuspiciousStartup *= 18
	breakdown.Total = breakdown.StartupRisk + breakdown.ProcessRisk + breakdown.SuspiciousStartup
	return breakdown
}

func RiskBreakdownSummary(breakdown RiskBreakdown) string {
	parts := make([]string, 0)
	if breakdown.StartupRisk > 0 {
		parts = append(parts, fmt.Sprintf("startup drift: %d", breakdown.StartupRisk))
	}
	if breakdown.ProcessRisk > 0 {
		parts = append(parts, fmt.Sprintf("process drift: %d", breakdown.ProcessRisk))
	}
	if breakdown.SuspiciousStartup > 0 {
		parts = append(parts, fmt.Sprintf("suspicious startup: %d", breakdown.SuspiciousStartup))
	}
	if len(parts) == 0 {
		return "No significant risk factors detected"
	}
	return strings.Join(parts, " | ")
}
