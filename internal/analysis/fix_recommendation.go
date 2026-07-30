package analysis

import "fmt"

func FixRecommendation(breakdown RiskBreakdown, score int) string {
	if score >= 80 {
		return "No urgent action needed; continue monitoring and keep the baseline stable."
	}
	if breakdown.SuspiciousStartup > 0 {
		return "Review and remove unexpected startup entries, especially scripts or command shells."
	}
	if breakdown.ProcessRisk > 0 {
		return "Inspect newly running background processes and check for unknown or resource-heavy services."
	}
	if breakdown.StartupRisk > 0 {
		return "Check recent startup additions and remove unnecessary apps from startup permissions."
	}
	return "Review the current system drift and verify whether the changes are legitimate."
}

func FixRecommendationSummary(breakdown RiskBreakdown, score int) string {
	return fmt.Sprintf("Recommended fix: %s", FixRecommendation(breakdown, score))
}
